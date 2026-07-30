package stylegen

import (
	"bytes"
	"errors"
	"fmt"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	gsxast "github.com/gsxhq/gsx/ast"
	gsxparser "github.com/gsxhq/gsx/parser"
	parse "github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/css"

	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/merge"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// slotMarkerPrefix is the attribute-name prefix every slot marker carries.
const slotMarkerPrefix = "data-gsxui-slot-"

// designSpecRef points at the invariant this file enforces.
const designSpecRef = "docs/superpowers/specs/2026-07-29-typed-recipe-model-design.md §9"

// ComponentComposedMarkers returns every data-gsxui-slot-* marker on an element
// that renders through a migrated component. Those markers are the ones whose
// presentation now comes from the utilities layer, so a components-layer rule
// against them is dead.
func ComponentComposedMarkers(root string) ([]string, error) {
	byMarker, err := composedMarkers(root)
	if err != nil {
		return nil, err
	}
	markers := make([]string, 0, len(byMarker))
	for marker := range byMarker {
		markers = append(markers, marker)
	}
	sort.Strings(markers)
	return markers, nil
}

// composedMarkers maps each composed marker to every migrated component whose
// compiled utilities land on the same element. A marker CAN have more than one
// owner: input-group.gsx stamps data-gsxui-slot-input-group-control onto both
// an <Input> and a <Textarea>, and once both are migrated neither is "the"
// owner — both bring compiled utilities to elements carrying that marker, and
// the gate's question ("does a components-layer rule lose to a migrated
// component's utilities") has to be answered against each of them
// independently, not by picking one or merging their utility sets (merging
// would answer "is this dead under the union", which is not what a reader can
// act on — they need to know WHICH component's utilities are winning). Owners
// are kept sorted so the diagnostic is deterministic regardless of gsx.Glob or
// map iteration order.
func composedMarkers(root string) (map[string][]string, error) {
	byMarker := map[string][]string{}
	addOwner := func(marker, component string) {
		for _, owner := range byMarker[marker] {
			if owner == component {
				return
			}
		}
		byMarker[marker] = append(byMarker[marker], component)
		sort.Strings(byMarker[marker])
	}
	// migrated maps a GSX component tag to the canonical component it renders.
	// The tag is the exported name — matching case-insensitively would make the
	// plain HTML <button> element look like <Button>, which it is not.
	migrated := map[string]string{}
	for name, shape := range shapes.All() {
		migrated[accessorName(name)] = name // "button" matches <ui.Button> and <Button>
		// A migrated component's own slot markers sit on the very elements that
		// carry its compiled utilities, so they are composed by definition.
		for _, slot := range shape.Slots {
			addOwner(slotMarker(name, slot.Name), name)
		}
	}

	paths, err := filepath.Glob(filepath.Join(root, "ui", "*.gsx"))
	if err != nil {
		return nil, err
	}
	sort.Strings(paths)

	for _, path := range paths {
		src, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		file, err := gsxparser.ParseFile(token.NewFileSet(), path, src, 0)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", path, err)
		}
		gsxast.Inspect(file, func(node gsxast.Node) bool {
			element, ok := node.(*gsxast.Element)
			if !ok {
				return true
			}
			component, ok := migrated[strings.TrimPrefix(element.Tag, "ui.")]
			if !ok {
				return true
			}
			for _, attr := range element.Attrs {
				marker := attrName(attr)
				if !strings.HasPrefix(marker, slotMarkerPrefix) {
					continue
				}
				addOwner(marker, component)
			}
			return true
		})
	}
	return byMarker, nil
}

// slotMarker names the marker attribute a component renders for one of its
// declared slots. The root slot carries the bare component marker.
func slotMarker(component, slot string) string {
	if slot == "" {
		return slotMarkerPrefix + component
	}
	return slotMarkerPrefix + component + "-" + slot
}

// attrName reports the attribute name of any named GSX attribute. Spread and
// markup attributes carry no marker name, so they report "".
func attrName(attr gsxast.Attr) string {
	switch a := attr.(type) {
	case *gsxast.StaticAttr:
		return a.Name
	case *gsxast.BoolAttr:
		return a.Name
	case *gsxast.ExprAttr:
		return a.Name
	case *gsxast.ClassAttr:
		return a.Name
	case *gsxast.EmbeddedAttr:
		return a.Name
	case *gsxast.MarkupAttr:
		return a.Name
	default:
		return ""
	}
}

// layerCheckedStylesheets is every individually named authored stylesheet the
// layer gate reads, spelled out one by one on purpose. A blanket glob would
// silently pull in whatever a future commit drops under assets/css or web/ — a
// file nobody decided was in scope.
// TestLayerCheckedStylesheetsCoverEveryAuthoredStylesheet fails when a new
// stylesheet appears, so adding one is a deliberate line in this list.
//
// Paths are slash-separated and relative to the repo root.
var layerCheckedStylesheets = []string{
	"assets/css/foundation.css",
	"assets/css/index.css",
	"assets/css/styles/default.css",
	"assets/css/themes/default.css",
	"web/site.css",
	"web/site-button.css",
}

// layerCheckedStylesheetDirs is every directory whose ENTIRE .css contents are
// in scope by construction, so the gate globs it instead of naming its files.
//
// There is exactly one, and the reason it is safe to glob is the reason the list
// above is not: assets/css/styles/default/ holds the default style's components
// layer split one file per component, and a file appearing there is a component's
// rules and nothing else. 41 components have yet to migrate; each migration
// deletes one of these files, and before the split each was a block in the
// middle of a single 3193-line sheet that every parallel migration had to edit.
// Naming them individually would mean a second edit, in this file, for every
// migration — churn that buys no decision, because "is this in scope" has one
// answer for the whole directory.
//
// The audited-by-default property the explicit list exists for is intact:
// TestLayerCheckedStylesheetsCoverEveryAuthoredStylesheet still compares the
// full walk of assets/css and web/ against the explicit list PLUS these
// directories, so a new stylesheet anywhere else — a new top-level sheet, a new
// style directory, a new web/ sheet — still fails until someone decides about it.
var layerCheckedStylesheetDirs = []string{
	"assets/css/styles/default",
}

// layerCheckedSheets resolves the explicit list and the globbed directories into
// the concrete set of stylesheets to read, slash-separated and relative to root.
//
// Order is the explicit list first, then each directory's files sorted by name.
// The gate judges every rule against the compiled utilities independently, so
// this order affects only the order diagnostics come out in — the CASCADE order
// of the split files is the aggregator's @import order, which
// internal/cli's TestComposedStyleCSSPreservesImportOrder is what pins.
func layerCheckedSheets(root string) ([]string, error) {
	sheets := slices.Clone(layerCheckedStylesheets)
	for _, dir := range layerCheckedStylesheetDirs {
		entries, err := os.ReadDir(filepath.Join(root, filepath.FromSlash(dir)))
		if err != nil {
			return nil, fmt.Errorf("layer gate: reading %s: %w", dir, err)
		}
		var found []string
		for _, entry := range entries {
			if entry.IsDir() || filepath.Ext(entry.Name()) != ".css" {
				return nil, fmt.Errorf("layer gate: %s/%s is not a stylesheet; this directory holds only the default style's per-component sheets", dir, entry.Name())
			}
			found = append(found, dir+"/"+entry.Name())
		}
		if len(found) == 0 {
			return nil, fmt.Errorf("layer gate: %s holds no stylesheets", dir)
		}
		sort.Strings(found)
		sheets = append(sheets, found...)
	}
	return sheets, nil
}

// exemptionKey identifies ONE deliberate violation: a file, the selector that
// loses the cascade, and the single utility or property it loses on. Keying by
// file alone — which is what this used to do — exempts every future rule anyone
// drops into that file too, so the first genuinely wrong rule added beside a
// deliberate one is invisible. 52 components are about to migrate through these
// files; a whole-file waiver is a hole that widens by itself.
type exemptionKey struct {
	file string
	// style is the installable style whose compiled utilities win. It is part of
	// the key because contest is a per-style fact: Maia's Button is `rounded-4xl`
	// and `h-9` where Nova's is `rounded-lg` and `h-8`, so the same authored rule
	// is dead under one style and live under the other. A key without this
	// dimension would let a decision taken about Nova silently forgive a Maia
	// violation nobody looked at — which is the whole failure P1b describes, just
	// moved into the exemption list.
	style     string
	selector  string
	contested string // utility name (`bg-primary`) or CSS property (`background-color`)
}

// layerCheckExemption records why one violation is deliberate. The reason is
// mandatory and is printed in every diagnostic that mentions the exemption, so
// a reader who trips over a suppressed rule learns why it is suppressed without
// going looking.
type layerCheckExemption struct {
	key    exemptionKey
	reason string
}

// layerCheckExemptions names the individual components-layer violations that are
// deliberate. Nothing else is exempt: a rule that loses the cascade by accident
// is exactly what this gate exists to catch, and an exemption is the one way to
// say "this one is on purpose".
//
// Every entry is SELF-INVALIDATING: CheckLayerPrecedence fails when an entry no
// longer names a violation that really occurs, so an exemption whose rule was
// deleted, reflowed, renamed, or rewritten to win the cascade fails the build
// instead of quietly widening its own scope.
//
// The Dialog-content fallback group is gone entirely: AlertDialog, Drawer and
// Sheet all migrated to the slot axis and carry the utilities they used to
// borrow through a marker on their own recipe classes, so the fallbacks the
// exemptions covered no longer exist to be exempted.
// The chain is slices.Concat rather than nested appends: every migration wave
// adds a group here, and a nested append chain put each addition on the same
// two lines, so every wave collided with every other one.
var layerCheckExemptions = slices.Concat(
	siteButtonFallbackExemptions(),
	toggleMarkerFallbackExemptions(),
	menuInsetDisabledFallbackExemptions(),
	calendarCaptionVisuallyHiddenExemptions(),
	comboboxItemDisabledFallbackExemptions(),
	sidebarFoundationMechanicsExemptions(),
	sidebarStyleOverridableExemptions(),
)

// comboboxItemDisabledFallbackReason is the decision behind every entry
// comboboxItemDisabledFallbackExemptions lists.
//
// assets/css/styles/default/combobox.css keeps ONE rule in @layer components
// on purpose — spec section 10b — `[data-gsxui-slot-combobox-item][data-
// disabled] { opacity-50 }`. data-disabled is stamped by the CALLER through
// attrs (ComboboxItem declares no `disabled` param at all), so the rule must
// stay caller-overridable, and @layer utilities — the layer Tailwind's own
// scanned classes live in — would flip that the wrong way. This is the same
// call assets/css/styles/default/menu.css's data-disabled rule makes; see
// menuInsetDisabledFallbackReason.
//
// Unlike menu's, though, the contest the gate reports here is a FALSE
// POSITIVE of its granularity, the calendar-caption shape:
// componentUtilities builds one utility set per COMPONENT, not per slot, and
// Combobox's union contains opacity-* only because combobox-content's
// popover transition emits opacity-0/opacity-100. Nothing on combobox-item
// sets opacity. The identical rule on Command's own item
// (assets/css/styles/default/command.css) is not reported at all, for exactly
// that reason — Command has no opacity utility anywhere.
const comboboxItemDisabledFallbackReason = "assets/css/styles/default/combobox.css's data-disabled/opacity-50 " +
	"rule is caller-attrs-driven per spec section 10b (ComboboxItem declares no disabled param; the caller " +
	"stamps data-disabled through attrs) and must stay in @layer components to remain caller-overridable. " +
	"The reported contest is an artifact of the gate comparing against Combobox's whole-component utility " +
	"union: the only opacity-* utilities Combobox emits are combobox-content's popover-transition " +
	"opacity-0/opacity-100, on a different element entirely. The same rule on command-item is not reported " +
	"only because Command emits no opacity utility anywhere."

func comboboxItemDisabledFallbackExemptions() []layerCheckExemption {
	const selector = ":where([data-gsxui-slot-combobox-item])[data-disabled]"
	out := make([]layerCheckExemption, 0, 2)
	for _, style := range []string{"maia", "nova"} {
		out = append(out, layerCheckExemption{
			key: exemptionKey{
				file:      "assets/css/styles/default/combobox.css",
				style:     style,
				selector:  selector,
				contested: "opacity-50",
			},
			reason: comboboxItemDisabledFallbackReason,
		})
	}
	return out
}

// calendarCaptionVisuallyHiddenReason is the decision behind every entry
// calendarCaptionVisuallyHiddenExemptions lists.
//
// assets/css/foundation.css's @layer base carries ONE visually-hidden rule
// shared by eight screen-reader labels across the catalogue
// (breadcrumb-ellipsis-label, pagination-ellipsis-label, dialog-close-label,
// sheet-close-label, command-dialog-header, combobox-bridge,
// carousel-control-label, and calendar-caption under
// data-caption-layout="dropdown"). Calendar's migration made
// calendar-caption the first marker in that list to sit on a component with
// compiled utilities, and the gate reports position/width/height/padding as
// contested.
//
// It is a false positive of the gate's granularity, not a real cascade loss.
// componentUtilities builds ONE utility set per COMPONENT, not per slot, so
// the union it compares against is every utility Calendar renders anywhere:
// the padding comes from the root's p-2, the width from w-fit/w-full, the
// height from month-caption's h-(--cell-size), the position from months'
// `relative` and nav's `absolute`. None of them is on the caption element.
// Calendar's caption slot renders exactly `text-sm font-medium select-none`
// and takes no attrs spread of its own, so nothing a caller or a style can
// put on that element contests the visually-hidden rule.
//
// The gate's usual fix does not apply either: moving this rule to @layer
// utilities would mean either splitting one a11y rule eight ways in
// foundation.css, or duplicating it into every style — and it must keep
// working with NO style loaded at all (jstest's foundation mode), which a
// style-owned rule cannot do.
const calendarCaptionVisuallyHiddenReason = "assets/css/foundation.css's shared visually-hidden rule for " +
	"screen-reader labels covers calendar-caption[data-caption-layout=\"dropdown\"], which since Calendar's " +
	"migration sits on a component with compiled utilities. The contest is an artifact of the gate comparing " +
	"against Calendar's whole-component utility union rather than the caption slot's own three utilities " +
	"(text-sm font-medium select-none) — no Calendar rule sets position, width, height or padding on the " +
	"caption element, and the caption slot accepts no caller attrs. The rule is foundation and must apply " +
	"with no style loaded, so it cannot move into a style's @layer utilities."

// calendarCaptionVisuallyHiddenExemptions spells out the four properties the
// shared rule is reported as losing, one (property, style) pair at a time,
// following siteButtonFallbackExemptions' enumeration discipline.
func calendarCaptionVisuallyHiddenExemptions() []layerCheckExemption {
	const selector = ":where( [data-gsxui-slot-calendar-caption][data-caption-layout=\"dropdown\"] )"
	out := make([]layerCheckExemption, 0, 8)
	for _, property := range []string{"height", "padding", "position", "width"} {
		for _, style := range []string{"maia", "nova"} {
			out = append(out, layerCheckExemption{
				key: exemptionKey{
					file:      "assets/css/foundation.css",
					style:     style,
					selector:  selector,
					contested: property,
				},
				reason: calendarCaptionVisuallyHiddenReason,
			})
		}
	}
	return out
}

// menuInsetDisabledFallbackExemptions spells out the two violations the
// retained menu.css rules commit, one (selector, utility, style) triple at a
// time, following siteButtonFallbackExemptions' own enumeration discipline.
// menuInsetDisabledFallbackReason is the decision behind both entries
// menuInsetDisabledFallbackExemptions lists. assets/css/styles/default/menu.css
// keeps a data-inset/pl-8 rule and a data-disabled/opacity-50 rule in @layer
// components on purpose — spec section 10b — because both data-inset and
// data-disabled are stamped by the CALLER through attrs on DropdownMenuItem/
// ContextMenuItem/MenubarItem and their Label/SubTrigger/CheckboxItem/RadioItem
// siblings, not by the component itself. Moving either rule to @layer
// utilities (this gate's usual fix) would make it win over a caller's plain
// utility unconditionally instead of losing to it: @layer utilities is the
// SAME layer Tailwind's own scanned classes live in, and this file is
// @imported before "tailwindcss" is scanned in every consumer, so the ordering
// would flip the wrong way. Keeping it in @layer components is what lets a
// caller's own plain `pl-*`/`opacity-*` utility (which lives in @layer
// utilities) win regardless of specificity, exactly as it did before this
// migration. The one real cost is that these two rules now also sit in the
// same layer as DropdownMenu/ContextMenu/Menubar's own compiled recipe
// utilities, which legitimately beat them there (a real, accepted loss, not a
// caller-facing regression) — jstest/specs/layer-precedence.spec.ts pins the
// caller-override case these two rules exist to protect.
const menuInsetDisabledFallbackReason = "assets/css/styles/default/menu.css's data-inset/data-disabled rules are " +
	"caller-attrs-driven per spec section 10b and must stay in @layer components to remain caller-" +
	"overridable; losing the cascade against DropdownMenu/ContextMenu/Menubar's own compiled recipe " +
	"utilities in that same layer is the design, not a defect. jstest/specs/layer-precedence.spec.ts's " +
	"menu caller-override pins cover the real rendering paths these rules must not break."

func menuInsetDisabledFallbackExemptions() []layerCheckExemption {
	both := []string{"maia", "nova"}
	pairs := []struct {
		selector, contested string
		styles              []string
	}{
		{
			selector: ":where( [data-gsxui-slot-dropdown-menu-item],[data-gsxui-slot-context-menu-item]," +
				"[data-gsxui-slot-menubar-item],[data-gsxui-slot-dropdown-menu-label]," +
				"[data-gsxui-slot-context-menu-label],[data-gsxui-slot-menubar-label]," +
				"[data-gsxui-slot-dropdown-menu-sub-trigger],[data-gsxui-slot-context-menu-sub-trigger]," +
				"[data-gsxui-slot-menubar-sub-trigger] )[data-inset]",
			contested: "pl-8",
			styles:    both,
		},
		{
			selector: ":where( [data-gsxui-slot-dropdown-menu-item],[data-gsxui-slot-context-menu-item]," +
				"[data-gsxui-slot-menubar-item],[data-gsxui-slot-dropdown-menu-checkbox-item]," +
				"[data-gsxui-slot-dropdown-menu-radio-item],[data-gsxui-slot-context-menu-checkbox-item]," +
				"[data-gsxui-slot-context-menu-radio-item],[data-gsxui-slot-menubar-checkbox-item]," +
				"[data-gsxui-slot-menubar-radio-item] )[data-disabled]",
			contested: "opacity-50",
			styles:    both,
		},
	}
	out := make([]layerCheckExemption, 0, len(pairs)*2)
	for _, pair := range pairs {
		if len(pair.styles) == 0 {
			panic("exemption for " + pair.selector + " / " + pair.contested + " names no style")
		}
		for _, style := range pair.styles {
			out = append(out, layerCheckExemption{
				key: exemptionKey{
					file:      "assets/css/styles/default/menu.css",
					style:     style,
					selector:  pair.selector,
					contested: pair.contested,
				},
				reason: menuInsetDisabledFallbackReason,
			})
		}
	}
	return out
}

// toggleMarkerFallbackReason is the decision behind every entry
// toggleMarkerFallbackExemptions lists. ui/toggle-group.gsx's
// ToggleGroupItem stamps data-gsxui-slot-toggle onto its own <button>
// alongside data-gsxui-slot-toggle-group-item, and has no class attribute
// of its own — every bit of its base appearance came entirely from
// Toggle's marker-keyed rules, not from a class Toggle's own migration
// adds. This :where()-wrapped fallback in assets/css/styles/default.css is
// authored at zero specificity in @layer components precisely so Toggle's
// own compiled recipe class beats it wherever a real <ui.Toggle> renders
// (both apply the identical values, so there is no visible difference
// either way) — the same "authored to lose on purpose" design as
// web/site-button.css's fallback above, just for a sibling component
// (ToggleGroup) composing an unmigrated marker instead of docs markup.
// Moving it to @layer utilities instead (this gate's usual fix) was tried
// first and rejected on the earlier Toggle migration attempt: it broke a
// real caller override, because @layer utilities is the SAME named layer
// Tailwind's own scanned utility classes live in, and this file is
// @imported after "tailwindcss" in every consumer — so these rules would
// always win over a caller's plain utility class instead of losing to it,
// the opposite of every other escape hatch in this file.
const toggleMarkerFallbackReason = "assets/css/styles/default.css's Toggle marker fallback is docs-and-" +
	"unmigrated-consumer-only presentation for ToggleGroupItem, which composes " +
	"data-gsxui-slot-toggle without going through <ui.Toggle>'s own accessor calls. It is " +
	"authored at zero specificity in @layer components precisely so a real Toggle's own " +
	"compiled utilities win wherever <ui.Toggle> renders; losing the cascade there is the " +
	"design, not a defect. jstest/specs/layer-precedence.spec.ts's Toggle caller-override pin " +
	"covers the real rendering path this fallback must not break."

// toggleMarkerFallbackExemptions spells out every violation the Toggle
// marker fallback commits, one (selector, utility, style) triple at a
// time, following siteButtonFallbackExemptions' own enumeration
// discipline exactly (see its own doc comment for why enumerating beats a
// whole-file waiver).
func toggleMarkerFallbackExemptions() []layerCheckExemption {
	both := []string{"maia", "nova"}
	pairs := []struct {
		selector, contested string
		styles              []string
	}{
		{selector: ":where([data-gsxui-slot-toggle])", contested: "rounded-lg", styles: both},
		{selector: ":where([data-gsxui-slot-toggle])", contested: "text-sm", styles: both},
		{selector: ":where([data-gsxui-slot-toggle]) svg:not([class*=\"size-\"])", contested: "size-4", styles: both},
		{selector: ":where([data-gsxui-slot-toggle]):focus-visible", contested: "border-ring", styles: both},
		{selector: ":where([data-gsxui-slot-toggle]):hover", contested: "bg-muted", styles: both},
		{selector: ":where([data-gsxui-slot-toggle])[aria-invalid=\"true\"]", contested: "border-destructive", styles: both},
		{selector: ":where([data-gsxui-slot-toggle])[data-size=\"default\"]", contested: "h-8", styles: both},
		{selector: ":where([data-gsxui-slot-toggle])[data-size=\"default\"]", contested: "min-w-8", styles: both},
		{selector: ":where([data-gsxui-slot-toggle])[data-size=\"default\"]:has(>svg)", contested: "px-2", styles: both},
		{selector: ":where([data-gsxui-slot-toggle])[data-size=\"lg\"]", contested: "h-9", styles: both},
		{selector: ":where([data-gsxui-slot-toggle])[data-size=\"lg\"]", contested: "min-w-9", styles: both},
		{selector: ":where([data-gsxui-slot-toggle])[data-size=\"lg\"]:has(>svg)", contested: "px-2", styles: both},
		{selector: ":where([data-gsxui-slot-toggle])[data-size=\"sm\"]", contested: "h-7", styles: both},
		{selector: ":where([data-gsxui-slot-toggle])[data-size=\"sm\"]", contested: "min-w-7", styles: both},
		{selector: ":where([data-gsxui-slot-toggle])[data-size=\"sm\"]", contested: "rounded-[min(var(--radius-md),12px)]", styles: both},
		{selector: ":where([data-gsxui-slot-toggle])[data-size=\"sm\"]", contested: "text-[0.8rem]", styles: both},
		{selector: ":where([data-gsxui-slot-toggle])[data-size=\"sm\"] svg:not([class*=\"size-\"])", contested: "size-3.5", styles: both},
		{selector: ":where([data-gsxui-slot-toggle])[data-size=\"sm\"]:has(>svg)", contested: "px-1.5", styles: both},
		{selector: ":where([data-gsxui-slot-toggle])[data-state=\"on\"]", contested: "bg-accent", styles: both},
		{selector: ":where([data-gsxui-slot-toggle])[data-variant=\"outline\"]:hover", contested: "bg-accent", styles: both},
	}
	out := make([]layerCheckExemption, 0, len(pairs)*2)
	for _, pair := range pairs {
		if len(pair.styles) == 0 {
			panic("exemption for " + pair.selector + " / " + pair.contested + " names no style")
		}
		for _, style := range pair.styles {
			out = append(out, layerCheckExemption{
				key: exemptionKey{
					file:      "assets/css/styles/default.css",
					style:     style,
					selector:  pair.selector,
					contested: pair.contested,
				},
				reason: toggleMarkerFallbackReason,
			})
		}
	}
	return out
}

// siteButtonFallbackReason is the single decision behind every entry below. It
// is attached to each one because a reader who trips over a suppressed rule
// should learn why it is suppressed from the diagnostic, not from a comment
// somewhere else in the tree.
const siteButtonFallbackReason = "web/site-button.css is docs-and-demos-only presentation for raw markup that carries " +
	"[data-gsxui-slot-button] WITHOUT rendering through <ui.Button> — hand-written examples and the " +
	"theme preview. It is authored at zero specificity in @layer components precisely so the compiled " +
	"Button's own utilities beat it wherever a real Button is rendered; losing the cascade there is the " +
	"design, not a defect. jstest/specs/basic-demo-presentation.spec.ts:182 pins that both halves still " +
	"hold in the browser."

// siteButtonFallbackExemptions spells out every violation the fallback really
// commits, one (selector, utility, style) triple at a time.
//
// This list is long on purpose. The old exemption was the single string
// "web/site-button.css", which forgave these violations AND every violation
// anyone adds to that file afterwards — including, during the migration wave, a
// genuinely wrong one. Enumerating them is the cost of the file no longer being
// a blind spot: a new dead rule in site-button.css is not on this list, so the
// gate reports it, and a rule that leaves the list is caught as stale.
//
// THE STYLE DIMENSION IS DECLARED, NOT DERIVED. Each pair names the styles under
// which it is a deliberate violation, and the expansion below turns one line
// into one key per style. Two things are deliberately NOT done: the styles are
// not enumerated from registry/styles (a new style would then inherit a waiver
// nobody granted, silently, which is the P1b failure moved into this list), and
// they are not read back from what the check currently finds (that turns a
// precise exemption into the blanket waiver this list replaced). Naming them
// makes each one a decision: `bg-primary` is contested under both styles because
// both set it, `rounded-lg` only under Nova because Maia is `rounded-4xl`, and a
// pair whose style list stops matching reality goes stale and fails the build.
func siteButtonFallbackExemptions() []layerCheckExemption {
	// Every pair below happens to be contested under both styles: the fallback
	// contests a Tailwind GROUP (border-radius, height, background-color), and
	// both Button recipes set every one of those groups even where they disagree
	// on the value — `rounded-lg` against `rounded-4xl`, `h-8` against `h-9`. A
	// pair contested under only one style is written with that style alone, and
	// the staleness check is what forces the question to be answered honestly.
	both := []string{"maia", "nova"}
	pairs := []struct {
		selector, contested string
		styles              []string
	}{
		{selector: ".dark :where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button][data-variant=\"destructive\"])", contested: "bg-destructive/60", styles: both},
		{selector: ".dark :where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button][data-variant=\"ghost\"]):hover", contested: "bg-accent/50", styles: both},
		{selector: ".dark :where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button][data-variant=\"outline\"])", contested: "bg-input/30", styles: both},
		{selector: ".dark :where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button][data-variant=\"outline\"])", contested: "border-input", styles: both},
		{selector: ".dark :where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button][data-variant=\"outline\"]):hover", contested: "bg-input/50", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button]) svg:not([class*=\"size-\"])", contested: "size-4", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])", contested: "border-transparent", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])", contested: "rounded-lg", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])", contested: "text-sm", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button]):focus-visible", contested: "border-ring", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[aria-invalid=\"true\"]", contested: "border-destructive", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-size=\"default\"]", contested: "gap-1.5", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-size=\"default\"]", contested: "h-8", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-size=\"default\"]", contested: "px-2.5", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-size=\"default\"]:has(>svg)", contested: "px-2", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-size=\"icon-lg\"]", contested: "size-9", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-size=\"icon-sm\"]", contested: "rounded-[min(var(--radius-md),12px)]", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-size=\"icon-sm\"]", contested: "size-7", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-size=\"icon-xs\"] svg:not([class*=\"size-\"])", contested: "size-3", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-size=\"icon-xs\"]", contested: "rounded-[min(var(--radius-md),10px)]", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-size=\"icon-xs\"]", contested: "size-6", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-size=\"icon\"]", contested: "size-8", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-size=\"lg\"]", contested: "gap-1.5", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-size=\"lg\"]", contested: "h-9", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-size=\"lg\"]", contested: "px-2.5", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-size=\"lg\"]:has(>svg)", contested: "px-2", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-size=\"sm\"] svg:not([class*=\"size-\"])", contested: "size-3.5", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-size=\"sm\"]", contested: "gap-1", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-size=\"sm\"]", contested: "h-7", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-size=\"sm\"]", contested: "px-2.5", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-size=\"sm\"]", contested: "rounded-[min(var(--radius-md),12px)]", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-size=\"sm\"]", contested: "text-[0.8rem]", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-size=\"sm\"]:has(>svg)", contested: "px-1.5", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-size=\"xs\"] svg:not([class*=\"size-\"])", contested: "size-3", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-size=\"xs\"]", contested: "gap-1", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-size=\"xs\"]", contested: "h-6", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-size=\"xs\"]", contested: "px-2", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-size=\"xs\"]", contested: "rounded-[min(var(--radius-md),10px)]", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-size=\"xs\"]", contested: "text-xs", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-size=\"xs\"]:has(>svg)", contested: "px-1.5", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-variant=\"default\"]", contested: "bg-primary", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-variant=\"default\"]", contested: "text-primary-foreground", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-variant=\"default\"]:hover", contested: "bg-primary/90", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-variant=\"destructive\"]", contested: "bg-destructive", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-variant=\"destructive\"]", contested: "text-contrast", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-variant=\"destructive\"]:hover", contested: "bg-destructive/90", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-variant=\"ghost\"]:hover", contested: "bg-accent", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-variant=\"ghost\"]:hover", contested: "text-accent-foreground", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-variant=\"link\"]", contested: "text-primary", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-variant=\"outline\"]", contested: "bg-background", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-variant=\"outline\"]", contested: "border-border", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-variant=\"outline\"]:hover", contested: "bg-accent", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-variant=\"outline\"]:hover", contested: "text-accent-foreground", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-variant=\"secondary\"]", contested: "bg-secondary", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-variant=\"secondary\"]", contested: "text-secondary-foreground", styles: both},
		{selector: ":where(body:not([data-theme-button-preview])) :where([data-gsxui-slot-button])[data-variant=\"secondary\"]:hover", contested: "bg-secondary/80", styles: both},
	}
	out := make([]layerCheckExemption, 0, len(pairs))
	for _, pair := range pairs {
		if len(pair.styles) == 0 {
			panic("exemption for " + pair.selector + " / " + pair.contested + " names no style")
		}
		for _, style := range pair.styles {
			out = append(out, layerCheckExemption{
				key: exemptionKey{
					file:      "web/site-button.css",
					style:     style,
					selector:  pair.selector,
					contested: pair.contested,
				},
				reason: siteButtonFallbackReason,
			})
		}
	}
	return out
}

// exemptionIndex maps each exempted violation to its reason. Building it here
// rather than declaring a map literal keeps the reason mandatory at the type
// level: an entry without one cannot be written down.
func exemptionIndex(exemptions []layerCheckExemption) map[exemptionKey]string {
	index := make(map[exemptionKey]string, len(exemptions))
	for _, exemption := range exemptions {
		index[exemption.key] = exemption.reason
	}
	return index
}

// staleExemptions reports every exemption that no longer forgives anything. An
// exemption that stops matching is not harmless: it is a decision whose subject
// has changed underneath it, and the next rule written in its place inherits a
// waiver nobody granted. Failing on it is what makes the exemption list
// self-invalidating rather than accumulating.
func staleExemptions(found []violation, exemptions []layerCheckExemption) []string {
	occurring := make(map[exemptionKey]struct{}, len(found))
	for _, v := range found {
		occurring[v.key] = struct{}{}
	}
	var stale []string
	for _, exemption := range exemptions {
		if _, ok := occurring[exemption.key]; ok {
			continue
		}
		stale = append(stale, fmt.Sprintf(
			"%s: under style %s, the exemption for %s\n  on %q no longer names a violation that occurs — the rule was\n"+
				"  deleted, reflowed, or fixed to win the cascade, or this style never contested it.\n"+
				"  Re-point it, drop %s from its style list, or delete it.\n"+
				"  Its reason was: %s",
			exemption.key.file, exemption.key.style, exemption.key.selector, exemption.key.contested,
			exemption.key.style, exemption.reason))
	}
	return stale
}

// normalizeSelectorKey makes an exemption's recorded selector comparable to the
// one the parser hands back, so that whitespace in the stylesheet does not
// decide whether an exemption applies.
func normalizeSelectorKey(selector string) string {
	return normalizeSpace(strings.ReplaceAll(selector, "\n", " "))
}

// CheckLayerPrecedence enforces the layer invariant of the design spec §9:
//
//	A rule overriding compiled component presentation must live in
//	@layer utilities AND carry specificity >= (0,1,0).
//
// Both halves fail silently on their own — the rule parses, generation and
// make audit stay green, and the override simply never applies in the browser.
// Every violation is reported, not just the first.
//
// It also fails on a STALE exemption. An exemption is a decision about one
// specific rule; when that rule stops existing the decision has no subject, and
// leaving the entry behind hands a waiver to whatever is written next in its
// place. Reporting it here is what keeps the exemption list from being a
// one-way ratchet.
func CheckLayerPrecedence(root string) error {
	found, err := allLayerViolations(root)
	if err != nil {
		return err
	}
	exempt := exemptionIndex(layerCheckExemptions)
	var reported []string
	for _, v := range found {
		if _, ok := exempt[v.key]; ok {
			continue
		}
		reported = append(reported, v.message)
	}
	reported = append(reported, staleExemptions(found, layerCheckExemptions)...)
	if len(reported) == 0 {
		return nil
	}
	return fmt.Errorf(
		"%d layer-precedence problem(s) — a rule that cannot win the cascade, or an\n"+
			"exemption that no longer forgives one:\n\n%s",
		len(reported), strings.Join(reported, "\n\n"))
}

// allLayerViolations runs the check ONCE PER INSTALLABLE STYLE and returns every
// violation any of them produces.
//
// Reading only registry/styles/<DefaultStyle> was a hole with a shipping user:
// Maia is installable, its Button utilities differ substantially from Nova's,
// and a shared structural rule can therefore be dead under Maia while `make
// audit` is green. The styles' utility sets are deliberately NOT unioned — a
// union answers "is this rule dead under some style" and throws away WHICH,
// which is exactly the fact whoever has to fix the rule needs.
func allLayerViolations(root string) ([]violation, error) {
	markers, err := composedMarkers(root)
	if err != nil {
		return nil, err
	}
	styles, err := styleNames(root)
	if err != nil {
		return nil, err
	}
	var out []violation
	for _, style := range styles {
		sets, err := componentUtilities(root, style, markers)
		if err != nil {
			return nil, err
		}
		// One resolver per style: the same utility name compiles to the same CSS
		// in every style, but the SETS differ, and a shared resolver would compile
		// utilities no style under test actually renders.
		properties := &utilityPropertyResolver{root: root, sets: sets}
		found, err := layerViolations(root, style, markers, sets, properties)
		if err != nil {
			return nil, err
		}
		out = append(out, found...)
	}
	return out, nil
}

// layerViolations collects every violation in every checked stylesheet under ONE
// style, exempt or not. Exemption is applied by the caller so that the same
// collection pass can answer both "what is broken" and "is every exemption still
// forgiving something real".
func layerViolations(
	root, style string,
	markers map[string][]string,
	sets map[string]utilitySets,
	properties *utilityPropertyResolver,
) ([]violation, error) {
	var out []violation
	sheets, err := layerCheckedSheets(root)
	if err != nil {
		return nil, err
	}
	for _, relative := range sheets {
		path := filepath.Join(root, filepath.FromSlash(relative))
		src, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		rules, err := layeredRules(path, src)
		if err != nil {
			return nil, err
		}
		for _, rule := range rules {
			found, err := rule.violations(relative, style, markers, sets, properties)
			if err != nil {
				return nil, err
			}
			out = append(out, found...)
		}
	}
	return out, nil
}

// componentUtilities collects, per migrated component that owns at least one
// composed marker, every utility that component can render. Two sets are kept
// because a component's `[&_svg…]:` utilities land on its descendants, not on
// the marker element itself.
type utilitySets struct {
	own        []string
	descendant []string
}

func componentUtilities(root, style string, markers map[string][]string) (map[string]utilitySets, error) {
	wanted := map[string]struct{}{}
	for _, components := range markers {
		for _, component := range components {
			wanted[component] = struct{}{}
		}
	}
	sets := make(map[string]utilitySets, len(wanted))
	for component := range wanted {
		path := filepath.Join(root, "registry", "styles", style, component+".css")
		src, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read %s %s recipe: %w", style, component, err)
		}
		parsed, err := recipe.ParseStyle(path, src)
		if err != nil {
			return nil, fmt.Errorf("parse %s %s recipe: %w", style, component, err)
		}
		var set utilitySets
		seen := map[string]struct{}{}
		for _, rule := range parsed.Rules() {
			for _, utility := range rule.Utilities {
				if _, dup := seen[utility]; dup {
					continue
				}
				seen[utility] = struct{}{}
				if inner, ok := descendantUtility(utility); ok {
					set.descendant = append(set.descendant, inner)
					continue
				}
				set.own = append(set.own, utility)
			}
		}
		sort.Strings(set.own)
		sort.Strings(set.descendant)
		sets[component] = set
	}
	return sets, nil
}

// descendantUtility unwraps Tailwind's arbitrary-descendant variant, so
// `[&_svg:not([class*='size-'])]:size-4` is compared as `size-4` against a rule
// that selects a descendant of the marker.
func descendantUtility(utility string) (string, bool) {
	if !strings.HasPrefix(utility, "[&") {
		return "", false
	}
	depth := 0
	for i, r := range utility {
		switch r {
		case '[':
			depth++
		case ']':
			depth--
			if depth == 0 {
				rest := utility[i+1:]
				if strings.HasPrefix(rest, ":") && len(rest) > 1 {
					return rest[1:], true
				}
				return "", false
			}
		}
	}
	return "", false
}

// layerRule is one ruleset from an authored stylesheet together with the
// cascade layer it resolved in, the utilities it applies, and the properties it
// sets directly. Both halves matter: foundation.css is written almost entirely
// in plain declarations, so a gate that only saw @apply would read those files
// and find nothing in them.
type layerRule struct {
	layer      string
	selector   string
	line       int
	utilities  []string
	properties []layerDeclaration
}

// layerDeclaration is one plain declaration together with the at-rules enclosing
// it. The at-rules are part of WHEN the declaration applies, and a declaration
// inside a nested @media inside a ruleset applies in a different state from its
// siblings outside that @media — so they are recorded per declaration, not per
// rule.
type layerDeclaration struct {
	property string
	atRules  []string
}

// violation is one rule losing the cascade on one contested thing. Contests are
// reported one at a time rather than grouped per rule so that an exemption can
// name the exact violation it forgives instead of the whole rule.
type violation struct {
	key     exemptionKey
	message string
}

// layeredRules parses one authored stylesheet and returns every ruleset that
// carries at least one @apply or one declaration, tagged with its enclosing
// @layer.
func layeredRules(filename string, src []byte) ([]layerRule, error) {
	parser := css.NewParser(parse.NewInputBytes(src), false)

	type frame struct {
		layer      string
		atRule     string
		selector   string
		line       int
		rule       bool
		utilities  []string
		properties []layerDeclaration
	}
	var stack []frame
	var rules []layerRule
	// innermost returns the innermost enclosing ruleset frame. Declarations and
	// @apply inside a nested @media/@supports belong to the ruleset around them,
	// and Tailwind's own output nests @supports inside rulesets too.
	innermost := func() int {
		for i := len(stack) - 1; i >= 0; i-- {
			if stack[i].rule {
				return i
			}
		}
		return -1
	}
	for {
		grammar, _, data := parser.Next()
		switch grammar {
		case css.ErrorGrammar:
			if err := parser.Err(); errors.Is(err, io.EOF) {
				return rules, nil
			} else if err != nil {
				return nil, fmt.Errorf("%s: %w", filename, err)
			}
			return nil, fmt.Errorf("%s: malformed CSS at offset %d", filename, parser.Offset())

		case css.BeginAtRuleGrammar:
			f := frame{}
			prelude := strings.TrimSpace(recipe.SelectorText(parser.Values()))
			if bytes.Equal(data, []byte("@layer")) {
				f.layer = prelude
			} else {
				f.atRule = normalizeSpace(string(data) + " " + prelude)
			}
			stack = append(stack, f)

		case css.BeginRulesetGrammar:
			stack = append(stack, frame{
				selector: recipe.SelectorText(parser.Values()),
				line:     lineAt(src, parser.Offset()),
				rule:     true,
			})

		case css.EndAtRuleGrammar, css.EndRulesetGrammar:
			if len(stack) == 0 {
				return nil, fmt.Errorf("%s: unbalanced block at offset %d", filename, parser.Offset())
			}
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if !top.rule || (len(top.utilities) == 0 && len(top.properties) == 0) {
				continue
			}
			layer := ""
			for i := len(stack) - 1; i >= 0; i-- {
				if stack[i].layer != "" {
					layer = stack[i].layer
					break
				}
			}
			rules = append(rules, layerRule{
				layer:      layer,
				selector:   top.selector,
				line:       top.line,
				utilities:  top.utilities,
				properties: top.properties,
			})

		case css.DeclarationGrammar:
			index := innermost()
			if index < 0 {
				continue
			}
			property := strings.TrimSpace(string(data))
			if property == "" {
				continue
			}
			var atRules []string
			for _, f := range stack {
				if f.atRule != "" {
					atRules = append(atRules, f.atRule)
				}
			}
			stack[index].properties = append(stack[index].properties, layerDeclaration{
				property: strings.ToLower(property),
				atRules:  atRules,
			})

		case css.AtRuleGrammar:
			if !bytes.Equal(data, []byte("@apply")) {
				continue
			}
			index := innermost()
			if index < 0 {
				continue
			}
			applied := strings.Fields(recipe.SelectorText(parser.Values()))
			stack[index].utilities = append(stack[index].utilities, applied...)
		}
	}
}

func lineAt(src []byte, offset int) int {
	if offset > len(src) {
		offset = len(src)
	}
	return bytes.Count(src[:offset], []byte("\n")) + 1
}

// violations reports every way this rule loses to compiled component
// presentation. A rule can name several markers via a selector list, so each
// complex selector is judged on its own.
func (r layerRule) violations(
	filename, style string,
	markers map[string][]string,
	sets map[string]utilitySets,
	properties *utilityPropertyResolver,
) ([]violation, error) {
	var out []violation
	for _, complex := range splitSelectorList(r.selector) {
		marker, components, descendant, ok := composedTarget(complex, markers)
		if !ok {
			continue
		}
		// A marker can have several owners (input-group-control renders
		// through both Input and Textarea once both are migrated). The rule
		// is contested if it is contested under ANY owner, so each owner is
		// checked independently and never merged into one combined utility
		// set — merging would answer "dead under the union", which hides
		// which component is actually doing the winning. mergeContests then
		// groups identical contests back together so a property/utility
		// contested by two owners is reported once, naming both.
		perOwner := make(map[string][]contest, len(components))
		for _, component := range components {
			set := sets[component]
			against := set.own
			if descendant {
				against = set.descendant
			}
			// Utilities are compared by name through the repo's Tailwind merger,
			// which already knows `hover:bg-x` and `bg-x` are different utilities —
			// so the variant blindness this gate had was only ever in the property
			// path below.
			contests := contestedUtilities(r.utilities, against)
			if len(r.properties) > 0 {
				owned, err := properties.statePropertiesOf(against)
				if err != nil {
					return nil, err
				}
				found, err := contestedProperties(r.properties, complex, marker, descendant, owned)
				if err != nil {
					return nil, fmt.Errorf(
						"%s:%d: style %s: selector %s: %w\n"+
							"  The layer gate converts this hazard into a build failure, so it must not\n"+
							"  pass a rule whose state it cannot model — teach the model this form in\n"+
							"  internal/stylegen/state.go, or rewrite the selector — see %s.",
						filepath.ToSlash(filename), r.line, style, normalizeSelectorKey(complex), err, designSpecRef)
				}
				contests = append(contests, found...)
			}
			if len(contests) > 0 {
				perOwner[component] = contests
			}
		}
		spec := selectorSpecificity(complex)
		for _, contested := range mergeContests(perOwner, components) {
			key := exemptionKey{
				file:      filepath.ToSlash(filename),
				style:     style,
				selector:  normalizeSelectorKey(complex),
				contested: contested.name,
			}
			owners := strings.Join(contested.owners, ", ")
			var message string
			switch {
			case r.layer != "utilities":
				layer := r.layer
				if layer == "" {
					layer = "no layer"
				}
				message = fmt.Sprintf(
					"%s:%d: under style %s, %s %s in @layer %s, but that\n"+
						"  element renders through %s, whose own utilities win the layer ordering.\n"+
						"  Move this rule to @layer utilities and give it specificity >= (0,1,0) —\n"+
						"  see %s.",
					filepath.ToSlash(filename), r.line, style, normalizeSelectorKey(complex),
					contested.describe(), layer, owners, designSpecRef)
			case !spec.beatsPlainUtility():
				message = fmt.Sprintf(
					"%s:%d: under style %s, %s %s in @layer utilities, but its\n"+
						"  specificity is %s — inside one layer the cascade falls back to specificity,\n"+
						"  and a plain utility class from %s scores (0,1,0). Drop the :where() wrapper\n"+
						"  so the selector carries specificity >= (0,1,0) — see %s.",
					filepath.ToSlash(filename), r.line, style, normalizeSelectorKey(complex),
					contested.describe(), spec, owners, designSpecRef)
			default:
				continue
			}
			out = append(out, violation{key: key, message: message})
		}
	}
	return out, nil
}

// contest is one thing a rule loses on: either a utility it applies or a
// property it declares. The two are kept distinguishable because
// `@apply rounded-full` is a class and `display: flex` is a property, and
// telling a reader "applies display" would send them looking for a utility that
// does not exist. owners names every composing component whose utilities
// contest this one thing — usually one, but a shared marker (see
// composedMarkers) can make it more than one, and the diagnostic must name all
// of them, not just the first found.
type contest struct {
	name     string
	property bool
	owners   []string
}

func (c contest) describe() string {
	if c.property {
		return "sets " + c.name
	}
	return "applies " + c.name
}

// mergeContests folds each owner's independently-computed contests back
// together, keyed by (name, property) so a rule contested by two owners
// reports once naming both rather than twice naming one apiece. order fixes
// the iteration order (composedTarget already returns owners sorted), so the
// owners list attached to a merged contest is deterministic run to run.
func mergeContests(perOwner map[string][]contest, order []string) []contest {
	index := map[string]int{}
	var merged []contest
	for _, owner := range order {
		for _, c := range perOwner[owner] {
			key := c.name
			if c.property {
				key = "prop:" + key
			} else {
				key = "util:" + key
			}
			if i, ok := index[key]; ok {
				merged[i].owners = append(merged[i].owners, owner)
				continue
			}
			c.owners = []string{owner}
			index[key] = len(merged)
			merged = append(merged, c)
		}
	}
	return merged
}

// contestedProperties reports which of a rule's own declarations name a CSS
// property the component's compiled utilities take over IN A STATE THIS RULE
// GUARANTEES. Unlike utilities, which the Tailwind merger can compare by name, a
// raw `display: flex` has no class name to compare — so the utilities are
// resolved to the (state, property) pairs they actually declare, and the
// comparison happens there.
//
// The relation is IMPLICATION, in one direction only. Matching on the property
// alone made the oracle variant-blind: a component utility `hover:bg-accent`
// made every unconditional `background-color` rule look contested, though the
// two never compete. Matching on property AND EQUAL STATE — the correction that
// followed — swung past the answer the other way and passed a `:hover` rule
// sitting under an unconditional utility, which is dead every time it applies.
// Only "the authored state implies the generated state" gets both.
func contestedProperties(
	declared []layerDeclaration,
	complex, marker string,
	descendant bool,
	owned ownedProperties,
) ([]contest, error) {
	var contests []contest
	seen := map[string]struct{}{}
	for _, declaration := range declared {
		state, ok, err := ruleState(declaration.atRules, complex, marker, descendant)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}
		if !owned.shadows(declaration.property, state) {
			continue
		}
		if _, dup := seen[declaration.property]; dup {
			continue
		}
		seen[declaration.property] = struct{}{}
		contests = append(contests, contest{name: declaration.property, property: true})
	}
	sort.Slice(contests, func(i, j int) bool { return contests[i].name < contests[j].name })
	return contests, nil
}

// ruleState derives an authored rule's state the same way a compiled utility's
// is derived: everything the selector requires of the subject beyond the marker
// itself, plus the at-rules the declaration sits in.
//
// The descendant case is the one place the two sides cannot be lined up. A
// component's descendant presentation is written as `[&_svg…]:size-4`, and the
// gate deliberately unwraps that variant to compare `size-4` against a rule
// selecting a descendant of the marker — which means the descendant compound is
// the TARGETING, already accounted for by that unwrapping, and Tailwind never
// compiled it. So only the enclosing context is compared and the descendant
// compound's own conditions are dropped. Dropping them WEAKENS the authored
// predicate, and a weaker predicate implies strictly less, so this can only ever
// miss a contest — never invent one.
func ruleState(atRules []string, complex, marker string, descendant bool) (statePredicate, bool, error) {
	if descendant {
		var ancestors []string
		for _, compound := range splitCompounds(complex) {
			if containsAttributeName(compound, marker) {
				break
			}
			ancestors = append(ancestors, compound)
		}
		state, err := buildPredicate(atRules, scopeAncestors(ancestors), "")
		return state, err == nil, err
	}
	ancestors, suffix, ok := subjectConditions(complex, func(compound string) (string, bool) {
		return removeMarkerAttribute(compound, marker)
	})
	if !ok {
		return statePredicate{}, false, nil
	}
	state, err := buildPredicate(atRules, ancestors, suffix)
	return state, err == nil, err
}

// contestedUtilities reports which of a rule's utilities collide with one the
// component itself renders. Collision is decided by the repo's own Tailwind
// merger — the same authority that decides which class wins everywhere else —
// so no property table has to be maintained here. Identical utilities are not
// contested: restating a value the component already sets changes nothing.
func contestedUtilities(applied, componentUtilities []string) []contest {
	var contested []contest
	for _, utility := range applied {
		for _, own := range componentUtilities {
			if own == utility {
				continue
			}
			if merge.Merge([]string{own, utility}) == utility {
				contested = append(contested, contest{name: utility})
				break
			}
		}
	}
	return contested
}

// composedTarget reports the composed marker a complex selector targets, and
// whether the selected element is the marker itself or a descendant of it.
func composedTarget(complex string, markers map[string][]string) (marker string, components []string, descendant, ok bool) {
	compounds := splitCompounds(complex)
	names := make([]string, 0, len(markers))
	for name := range markers {
		names = append(names, name)
	}
	// Map iteration order is random, and one compound can legitimately name two
	// markers owned by different components — `[data-gsxui-slot-card][data-gsxui-slot-sidebar-inset]`.
	// Picking whichever came out of the map first would make the diagnostic, and
	// therefore the gate's verdict, vary run to run. Longest marker first is the
	// most specific reading of the compound; the name breaks remaining ties.
	sort.Slice(names, func(i, j int) bool {
		if len(names[i]) != len(names[j]) {
			return len(names[i]) > len(names[j])
		}
		return names[i] < names[j]
	})
	for i := len(compounds) - 1; i >= 0; i-- {
		for _, name := range names {
			if !containsAttributeName(compounds[i], name) {
				continue
			}
			// markers[name] is already sorted by composedMarkers' addOwner.
			return name, markers[name], i != len(compounds)-1, true
		}
	}
	return "", nil, false, false
}

// containsAttributeName reports whether selector names the given attribute,
// requiring a full-token match so data-gsxui-slot-card never matches
// data-gsxui-slot-card-header.
func containsAttributeName(selector, name string) bool {
	for offset := 0; ; {
		index := strings.Index(selector[offset:], name)
		if index < 0 {
			return false
		}
		index += offset
		after := index + len(name)
		if after >= len(selector) || !isAttributeNameByte(selector[after]) {
			if index == 0 || !isAttributeNameByte(selector[index-1]) {
				return true
			}
		}
		offset = index + 1
	}
}

func isAttributeNameByte(b byte) bool {
	return b == '-' || b == '_' ||
		(b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9')
}

// splitSelectorList splits a selector list on its top-level commas, ignoring
// commas nested inside :where(...)/:is(...) argument lists and attribute values.
func splitSelectorList(selector string) []string {
	parts := splitSelectorListKeepingEmpty(selector)
	out := parts[:0]
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

// splitSelectorListKeepingEmpty is the same split with the empty parts left in.
// An empty branch is not noise: inside `:is()`/`:where()` it is a branch that
// matches everything, which makes the whole group vacuous, and dropping it would
// turn "no condition at all" into "this list of conditions".
func splitSelectorListKeepingEmpty(selector string) []string {
	var parts []string
	depth := 0
	start := 0
	var quote byte
	for i := 0; i < len(selector); i++ {
		c := selector[i]
		switch {
		case quote != 0:
			if c == quote {
				quote = 0
			}
		case c == '"' || c == '\'':
			quote = c
		case c == '(' || c == '[':
			depth++
		case c == ')' || c == ']':
			depth--
		case c == ',' && depth == 0:
			parts = append(parts, selector[start:i])
			start = i + 1
		}
	}
	return append(parts, selector[start:])
}

// splitCompounds splits one complex selector into its compound selectors,
// dropping the combinators between them. The last compound is the subject.
func splitCompounds(complex string) []string {
	var compounds []string
	depth := 0
	start := 0
	var quote byte
	flush := func(end int) {
		if trimmed := strings.TrimSpace(complex[start:end]); trimmed != "" {
			compounds = append(compounds, trimmed)
		}
	}
	for i := 0; i < len(complex); i++ {
		c := complex[i]
		switch {
		case quote != 0:
			if c == quote {
				quote = 0
			}
		case c == '"' || c == '\'':
			quote = c
		case c == '(' || c == '[':
			depth++
		case c == ')' || c == ']':
			depth--
		case depth == 0 && (c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '>' || c == '+' || c == '~'):
			flush(i)
			start = i + 1
		}
	}
	flush(len(complex))
	if len(compounds) == 0 {
		return []string{strings.TrimSpace(complex)}
	}
	return compounds
}

// specificity is a CSS specificity triple (ids, classes, types).
type specificity struct {
	a, b, c int
}

func (s specificity) String() string { return fmt.Sprintf("(%d,%d,%d)", s.a, s.b, s.c) }

func (s specificity) less(other specificity) bool {
	if s.a != other.a {
		return s.a < other.a
	}
	if s.b != other.b {
		return s.b < other.b
	}
	return s.c < other.c
}

func (s specificity) add(other specificity) specificity {
	return specificity{s.a + other.a, s.b + other.b, s.c + other.c}
}

// beatsPlainUtility reports whether the selector can out-rank a plain Tailwind
// utility class, which scores exactly (0,1,0).
func (s specificity) beatsPlainUtility() bool {
	return !s.less(specificity{0, 1, 0})
}

// selectorSpecificity computes a complex selector's specificity per CSS
// Selectors Level 4: :where() contributes nothing, and :is()/:not()/:has()
// contribute the specificity of their most specific argument.
func selectorSpecificity(selector string) specificity {
	var total specificity
	for i := 0; i < len(selector); {
		switch c := selector[i]; {
		case c == '#':
			i = skipIdent(selector, i+1)
			total.a++
		case c == '.':
			i = skipIdent(selector, i+1)
			total.b++
		case c == '[':
			i = skipBalanced(selector, i, '[', ']')
			total.b++
		case c == ':':
			i, total = pseudoSpecificity(selector, i, total)
		case c == '*':
			i++
		case isIdentStartByte(c):
			i = skipIdent(selector, i)
			total.c++
		default:
			i++
		}
	}
	return total
}

// pseudoSpecificity consumes one pseudo-class or pseudo-element starting at the
// ':' in index and folds its contribution into total.
func pseudoSpecificity(selector string, index int, total specificity) (int, specificity) {
	i := index + 1
	element := false
	if i < len(selector) && selector[i] == ':' {
		element = true
		i++
	}
	nameStart := i
	i = skipIdent(selector, i)
	name := strings.ToLower(selector[nameStart:i])
	if i < len(selector) && selector[i] == '(' {
		end := skipBalanced(selector, i, '(', ')')
		args := selector[i+1 : max(end-1, i+1)]
		i = end
		switch name {
		case "where":
			return i, total
		case "is", "not", "has", "matches", "any":
			var best specificity
			for _, arg := range splitSelectorList(args) {
				if candidate := selectorSpecificity(arg); best.less(candidate) {
					best = candidate
				}
			}
			return i, total.add(best)
		}
	}
	if element {
		total.c++
	} else {
		total.b++
	}
	return i, total
}

func isIdentStartByte(b byte) bool {
	return b == '-' || b == '_' || b == '\\' || b >= 0x80 ||
		(b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

func skipIdent(s string, i int) int {
	for i < len(s) && isAttributeNameByte(s[i]) {
		i++
	}
	return i
}

// skipBalanced returns the index just past the bracket that closes the one at
// open, ignoring brackets inside strings.
func skipBalanced(s string, open int, left, right byte) int {
	depth := 0
	var quote byte
	for i := open; i < len(s); i++ {
		c := s[i]
		switch {
		case quote != 0:
			if c == quote {
				quote = 0
			}
		case c == '"' || c == '\'':
			quote = c
		case c == left:
			depth++
		case c == right:
			depth--
			if depth == 0 {
				return i + 1
			}
		}
	}
	return len(s)
}

// sidebarFoundationMechanicsReason is the decision behind every entry
// sidebarFoundationMechanicsExemptions lists.
//
// assets/css/foundation.css carries Sidebar's state-machine mechanics — the
// display/position/width/height/flex arithmetic that makes one semantic tree
// render as either a mobile Sheet or a fixed desktop column. Sidebar's
// migration to the slot axis gave those same elements compiled recipe
// utilities, and the gate now reports 44 (selector, property) pairs as
// contested under each style.
//
// Every one of them is the same false positive Calendar's caption already
// documents, at a much larger scale: componentUtilities builds ONE utility set
// per COMPONENT, not per slot, so the union it compares each foundation rule
// against is every utility Sidebar renders ANYWHERE across its 32 slots. The
// `display` it reports against the wrapper comes from menu-action's
// collapsed-icon `hidden`; the `width` from menu-sub-button's `min-w-0`; the
// `content` against group-action::after from the rail's own
// `hover:after:bg-sidebar-border`. Checked slot by slot against
// registry/styles/{nova,maia}/sidebar.css, NO recipe rule sets any of these
// properties on the element the foundation rule targets — which is the whole
// reason the mechanics stayed in foundation and the presentation moved.
//
// The gate's usual fix does not apply. These rules must keep working with no
// style loaded at all (jstest's foundation mode), which a style-owned
// @layer utilities rule cannot do; and promoting them into foundation's own
// @layer utilities block would make them beat a caller's plain utility instead
// of losing to it — site/examples/sidebar/basic.gsx's
// `class="min-h-[32rem]"` on the wrapper is exactly that caller, and it must
// keep winning over foundation's `min-height: 100svh`.
const sidebarFoundationMechanicsReason = "assets/css/foundation.css holds Sidebar's state-machine mechanics " +
	"(display/position/size/flex), which since Sidebar's migration sit on a component with compiled " +
	"utilities. The contest is an artifact of the gate comparing against Sidebar's whole-component utility " +
	"union rather than the individual slot's own utilities — no rule in registry/styles/{nova,maia}/" +
	"sidebar.css sets the reported property on the element the foundation rule targets. These rules must " +
	"apply with no style loaded, so they cannot move into a style's @layer utilities, and promoting them " +
	"inside foundation would make them beat caller utilities they are required to lose to."

// sidebarRailConditionalOverrideReason covers the two entries the gate is
// literally right about, and which are still deliberate.
//
// foundation.css's rail base sets `transform: translateX(-50%)` and its
// ::after sets `left: 50%`. Sidebar's recipe overrides both — but only under
// an offcanvas desktop ancestor, as `[…[data-collapsible=offcanvas] &]:
// [transform:translateX(0)]` and `…:after:left-full`. That is exactly what the
// pre-migration rules in assets/css/styles/default/sidebar.css did from the
// same layer by source order; the migration moved the winner into
// @layer utilities and left the loser where it was, so the unconditional
// foundation value still applies in every non-offcanvas state. The gate models
// neither the variant condition nor the fact that the base value is still
// reachable, so it reports a total loss where there is a conditional one.
// jstest/specs/layer-precedence.spec.ts pins both halves.
const sidebarRailConditionalOverrideReason = "assets/css/foundation.css's rail transform and ::after offset are " +
	"overridden by registry/styles/{nova,maia}/sidebar.css ONLY under an offcanvas desktop ancestor, which " +
	"is what the pre-migration rules did by source order in the same layer. The unconditional foundation " +
	"value still applies in every other state; the gate models neither the variant condition nor the " +
	"surviving base value. jstest/specs/layer-precedence.spec.ts pins both."

// sidebarFoundationMechanicsExemptions spells out every violation Sidebar's
// migration exposes in foundation.css, one (selector, property, style) triple
// at a time, following siteButtonFallbackExemptions' enumeration discipline.
func sidebarFoundationMechanicsExemptions() []layerCheckExemption {
	both := []string{"maia", "nova"}
	pairs := []struct {
		selector, contested, reason string
		styles                      []string
	}{
		{selector: ":is( [data-gsxui-slot-sidebar-desktop][data-collapsible=\"icon\"] ) :is( [data-gsxui-slot-sidebar-menu-button-tooltip-content]:popover-open )", contested: "display", styles: both},
		{selector: ":where( [data-gsxui-slot-sidebar-desktop][data-collapsible=\"icon\"] ) :where([data-gsxui-slot-sidebar-content])", contested: "overflow", styles: both},
		{selector: ":where( [data-gsxui-slot-sidebar-desktop][data-collapsible=\"icon\"] )>:where([data-gsxui-slot-sidebar-container])", contested: "width", styles: both},
		{selector: ":where( [data-gsxui-slot-sidebar-desktop][data-collapsible=\"icon\"] )>:where([data-gsxui-slot-sidebar-gap])", contested: "width", styles: both},
		{selector: ":where( [data-gsxui-slot-sidebar-desktop][data-collapsible=\"offcanvas\"] )>:where([data-gsxui-slot-sidebar-gap])", contested: "width", styles: both},
		{selector: ":where( [data-gsxui-slot-sidebar-desktop][data-side=\"left\"][data-collapsible=\"offcanvas\"] )>:where([data-gsxui-slot-sidebar-container])", contested: "left", styles: both},
		{selector: ":where( [data-gsxui-slot-sidebar-desktop][data-side=\"right\"][data-collapsible=\"offcanvas\"] )>:where([data-gsxui-slot-sidebar-container])", contested: "right", styles: both},
		{selector: ":where([data-gsxui-slot-sidebar-container])", contested: "display", styles: both},
		{selector: ":where([data-gsxui-slot-sidebar-container])", contested: "height", styles: both},
		{selector: ":where([data-gsxui-slot-sidebar-container])", contested: "position", styles: both},
		{selector: ":where([data-gsxui-slot-sidebar-container])", contested: "width", styles: both},
		{selector: ":where([data-gsxui-slot-sidebar-content])", contested: "display", styles: both},
		{selector: ":where([data-gsxui-slot-sidebar-content])", contested: "flex", styles: both},
		{selector: ":where([data-gsxui-slot-sidebar-content])", contested: "flex-direction", styles: both},
		{selector: ":where([data-gsxui-slot-sidebar-content])", contested: "overflow", styles: both},
		{selector: ":where([data-gsxui-slot-sidebar-desktop])", contested: "display", styles: both},
		{selector: ":where([data-gsxui-slot-sidebar-desktop]) :where([data-gsxui-slot-sidebar-rail])", contested: "display", styles: both},
		{selector: ":where([data-gsxui-slot-sidebar-desktop][data-side=\"left\"])>:where([data-gsxui-slot-sidebar-container])", contested: "left", styles: both},
		{selector: ":where([data-gsxui-slot-sidebar-desktop][data-side=\"right\"])>:where([data-gsxui-slot-sidebar-container])", contested: "right", styles: both},
		{selector: ":where([data-gsxui-slot-sidebar-gap])", contested: "flex", styles: both},
		{selector: ":where([data-gsxui-slot-sidebar-gap])", contested: "position", styles: both},
		{selector: ":where([data-gsxui-slot-sidebar-gap])", contested: "width", styles: both},
		{selector: ":where([data-gsxui-slot-sidebar-group-action])::after", contested: "content", styles: both},
		{selector: ":where([data-gsxui-slot-sidebar-inner])", contested: "display", styles: both},
		{selector: ":where([data-gsxui-slot-sidebar-inner])", contested: "flex-direction", styles: both},
		{selector: ":where([data-gsxui-slot-sidebar-inner])", contested: "height", styles: both},
		{selector: ":where([data-gsxui-slot-sidebar-inner])", contested: "width", styles: both},
		{selector: ":where([data-gsxui-slot-sidebar-inset])", contested: "display", styles: both},
		{selector: ":where([data-gsxui-slot-sidebar-inset])", contested: "flex", styles: both},
		{selector: ":where([data-gsxui-slot-sidebar-inset])", contested: "flex-direction", styles: both},
		{selector: ":where([data-gsxui-slot-sidebar-inset])", contested: "width", styles: both},
		{selector: ":where([data-gsxui-slot-sidebar-menu-action])::after", contested: "content", styles: both},
		{selector: ":where([data-gsxui-slot-sidebar-mobile-content]) :where([data-gsxui-slot-sidebar-rail])", contested: "display", styles: both},
		{selector: ":where([data-gsxui-slot-sidebar-rail])", contested: "display", styles: both},
		{selector: ":where([data-gsxui-slot-sidebar-rail])", contested: "position", styles: both},
		{selector: ":where([data-gsxui-slot-sidebar-rail])", contested: "width", styles: both},
		{selector: ":where([data-gsxui-slot-sidebar-rail])::after", contested: "content", styles: both},
		{selector: ":where([data-gsxui-slot-sidebar-wrapper])", contested: "display", styles: both},
		{selector: ":where([data-gsxui-slot-sidebar-wrapper])", contested: "width", styles: both},
		{selector: ":where([data-gsxui-slot-sidebar][data-collapsible=\"none\"])", contested: "display", styles: both},
		{selector: ":where([data-gsxui-slot-sidebar][data-collapsible=\"none\"])", contested: "flex-direction", styles: both},
		{selector: ":where([data-gsxui-slot-sidebar][data-collapsible=\"none\"])", contested: "height", styles: both},
		{selector: ":where([data-gsxui-slot-sidebar][data-collapsible=\"none\"])", contested: "width", styles: both},

		// The two the gate is literally right about; see the reason above.
		{selector: ":where([data-gsxui-slot-sidebar-rail])", contested: "transform", reason: sidebarRailConditionalOverrideReason, styles: both},
		{selector: ":where([data-gsxui-slot-sidebar-rail])::after", contested: "left", reason: sidebarRailConditionalOverrideReason, styles: both},
	}
	out := make([]layerCheckExemption, 0, len(pairs)*2)
	for _, pair := range pairs {
		if len(pair.styles) == 0 {
			panic("exemption for " + pair.selector + " / " + pair.contested + " names no style")
		}
		reason := pair.reason
		if reason == "" {
			reason = sidebarFoundationMechanicsReason
		}
		for _, style := range pair.styles {
			out = append(out, layerCheckExemption{
				key: exemptionKey{
					file:      "assets/css/foundation.css",
					style:     style,
					selector:  pair.selector,
					contested: pair.contested,
				},
				reason: reason,
			})
		}
	}
	return out
}

// sidebarStyleOverridableReason is the decision behind the two entries
// sidebarStyleOverridableExemptions lists.
//
// assets/css/styles/default/sidebar.css keeps three rule groups in
// @layer components on purpose — spec section 10b — because a STYLE is
// expected to redefine them, and a style can only do that from this layer.
// jstest/specs/sidebar-style-contract.spec.ts's "custom style owns floating
// and inset collapsed density without foundation shrink" injects exactly such
// a stylesheet and requires all three to take effect; migrating them onto the
// recipe made that test fail, which is how the retention was decided rather
// than guessed.
//
// Two of the three set raw custom properties the gate does not model. The
// other two — the floating/inset container padding and the floating inner
// border — it reports, correctly, as sitting in the same layer as Sidebar's
// own compiled utilities. That is an accepted cost, not a caller-facing
// regression: registry/styles/{nova,maia}/sidebar.css deliberately sets no
// padding on sidebar-container and no border width or colour on
// sidebar-inner, so nothing of Sidebar's own actually shadows either rule and
// both are live. This is the same shape as menuInsetDisabledFallbackReason.
const sidebarStyleOverridableReason = "assets/css/styles/default/sidebar.css's floating/inset container padding " +
	"and floating inner border are style-overridable per spec section 10b and must stay in @layer components " +
	"so a consumer stylesheet can beat them; Sidebar's own recipe sets neither property on those slots, so " +
	"neither rule is shadowed in practice. jstest/specs/sidebar-style-contract.spec.ts's custom-density test " +
	"is the contract these two exist to keep."

func sidebarStyleOverridableExemptions() []layerCheckExemption {
	both := []string{"maia", "nova"}
	pairs := []struct {
		selector, contested string
		styles              []string
	}{
		{
			selector:  ":where( [data-gsxui-slot-sidebar-desktop]:is( [data-variant=\"floating\"],[data-variant=\"inset\"] ) )>:where([data-gsxui-slot-sidebar-container])",
			contested: "p-2",
			styles:    both,
		},
		{
			selector:  ":where( [data-gsxui-slot-sidebar-desktop][data-variant=\"floating\"] ) :where([data-gsxui-slot-sidebar-inner])",
			contested: "border",
			styles:    both,
		},
	}
	out := make([]layerCheckExemption, 0, len(pairs)*2)
	for _, pair := range pairs {
		if len(pair.styles) == 0 {
			panic("exemption for " + pair.selector + " / " + pair.contested + " names no style")
		}
		for _, style := range pair.styles {
			out = append(out, layerCheckExemption{
				key: exemptionKey{
					file:      "assets/css/styles/default/sidebar.css",
					style:     style,
					selector:  pair.selector,
					contested: pair.contested,
				},
				reason: sidebarStyleOverridableReason,
			})
		}
	}
	return out
}
