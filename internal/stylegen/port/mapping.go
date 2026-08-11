package port

import "strings"

// ignoredSections are upstream MARK sections with no gsxui component
// counterpart at all: seven blocks belong to an AI-chat/composer surface
// gsxui doesn't ship (Bubble, Attachment, Marker, Questionnaire, Message
// Scroller, Message), one is a chart surface with no component or shape yet
// (Chart), and one is a translucent-menu variant (Menu Translucent). "React
// Aria" is a ninth: its 30 rules are all `-aria`-suffixed duplicates of
// rules that already live in their own component's section — Ignored's
// suffix check would catch every one of them individually, but excluding
// the whole section here means a stray non-`-aria` class landing in it by
// accident can never slip through as a fabricated 55th component.
var ignoredSections = map[string]bool{
	"React Aria":       true,
	"Chart":            true,
	"Menu Translucent": true,
	"Bubble":           true,
	"Attachment":       true,
	"Marker":           true,
	"Questionnaire":    true,
	"Message Scroller": true,
	"Message":          true,
}

// sectionComponentOverrides are declared exceptions to the default
// lowercase-and-hyphenate MARK-name derivation. Verified against the full
// 54-component survey (dossier §6): "Radio Group" is the ONLY override
// across all 54 — default derivation would produce "radio-group", but
// registry/canonical/shapes/radio.go names the shape "radio" (a single
// native <input type="radio"> root slot; "group" describes the surrounding
// ui.RadioGroup composition, not this shape).
var sectionComponentOverrides = map[string]string{
	"Radio Group": "radio",
}

// SectionComponent maps an upstream MARK name to a gsxui component name.
// Returns ok=false for sections we deliberately ignore: the synthetic ""
// section (rules appearing before any MARK comment — never a component) and
// the ignoredSections table above.
func SectionComponent(mark string) (component string, ok bool) {
	if mark == "" || ignoredSections[mark] {
		return "", false
	}
	if override, has := sectionComponentOverrides[mark]; has {
		return override, true
	}
	return strings.ToLower(strings.ReplaceAll(mark, " ", "-")), true
}

// Ignored reports classes skipped by policy: React-Aria's shadow rules.
// These duplicate a component's primary rules under React Aria's own
// data-* attribute vocabulary (data-selected instead of data-checked, and
// so on) rather than Tailwind's native data-[attr=value] arbitrary variant,
// so they carry no information the primary rule doesn't already have. They
// appear both inline within a component's own MARK section (e.g. Select's
// cn-select-value-aria, cn-select-item-aria, cn-select-empty-aria) and in
// the trailing "React Aria" section (already excluded by SectionComponent
// above) — matching by suffix catches both regardless of which section a
// rule physically lives in.
func Ignored(class string) bool {
	return strings.HasSuffix(class, "-aria")
}

// styleInvariantComponents take the same recipe in every style: no upstream
// MARK section exists for them in any of the 8 style-<name>.css files, so
// they stay hand-authored residue outside the porter's reach entirely (spec
// §4.4). Confirmed against dossier §6: these are the only 4 of the 54
// canonical shapes with no upstream section at all.
var styleInvariantComponents = map[string]bool{
	"aspect-ratio": true,
	"collapsible":  true,
	"spinner":      true,
	"toaster":      true,
}

// StyleInvariant reports whether component takes the same recipe in every
// style.
func StyleInvariant(component string) bool {
	return styleInvariantComponents[component]
}

// slotOverride is a declared exception to the default class→slot
// derivation: either a fixed target slot (ok=true) or an explicit "this
// upstream class has no gsxui slot" (ok=false). Either way it is a REVIEWED
// decision, not a fallthrough: Transform's unmapped reporting only surfaces
// classes that hit no table entry at all (an unrecognized stray — real
// signal), never a class this table explicitly routed nowhere on purpose
// (silentNonMapping) — matching dossier §4.3's "explicitly ignored, listed
// [in this table, in review,] not silently dropped [from review]".
//
// variant, when set alongside ok=true, is inserted as an extra Tailwind
// variant immediately before each contributed utility's own base utility
// (see insertVariant) — the mechanism the overlay→content fusions below need:
// upstream's separate overlay rule targets a real DOM node, but gsxui's
// native <dialog> paints the same backdrop as ::backdrop, a PSEUDO-ELEMENT of
// dialog-content, so a plain slot-rename would apply the overlay's
// `bg-black/80` to the dialog BOX itself, not its backdrop — visibly wrong,
// and (for alert-dialog specifically, see below) also a real Tailwind
// conflict once merged against the box's own background utility.
type slotOverride struct {
	slot    string
	ok      bool
	variant string
	extra   []string
}

// slotOverrides holds per-component exceptions to the default class→slot
// derivation, keyed by [component][upstream class]. "At minimum" per the
// plan: select's native-popover scroll buttons have no gsxui slot at all
// (the browser handles overflow scrolling on its own — no JS needed), and
// three of the four native-<dialog>-backed compound components fuse
// upstream's separate overlay rule onto their content slot's ::backdrop.
//
// alert-dialog is the one exception: registry/canonical/shapes/alert-dialog.go
// documents that AlertDialogContent composes <ui.DialogContent> directly, so
// it inherits Dialog's own ::backdrop chrome by GSX composition and never
// declares any backdrop styling of its own (confirmed:
// registry/styles/nova/alert-dialog.css's content rule carries only the
// max-w-* narrowing, nothing backdrop-shaped). Naively fusing
// cn-alert-dialog-overlay onto alert-dialog's content slot anyway landed
// bg-black/80 next to content's own bg-popover on the SAME rule — two
// background utilities CheckConflicts correctly flagged as a real Tailwind
// conflict (merge.Merge collapses to one, "superseded"). Its overlay is
// dropped as a silent non-mapping instead: not a gap, upstream's Alert
// Dialog and Dialog sections are simply redundant here by design.
var slotOverrides = map[string]map[string]slotOverride{
	"select": {
		// Native popover: no JS scroll buttons to style.
		"cn-select-scroll-up-button":   {ok: false},
		"cn-select-scroll-down-button": {ok: false},
		// gsxui's SelectTrigger paints its icon via an arbitrary `[&_svg]:`
		// descendant selector on the trigger rule itself (see
		// registry/styles/nova/select.css), not a separate icon element/slot
		// — there is nowhere for a standalone trigger-icon rule to land.
		"cn-select-trigger-icon": {ok: false},
		// RTL-logical slide-in variants for a directional popover-open
		// animation gsxui's select doesn't implement: nova's content rule
		// uses a plain fade+scale transition with no per-side slide at all
		// (a pre-porter design choice, not a porter gap).
		"cn-select-content-logical": {ok: false},
		// Select's shape doc comment names select-group explicitly: "carries
		// no accessor — nothing to migrate."
		"cn-select-group": {ok: false},
	},
	"dialog": {
		// Native <dialog>'s ::backdrop is part of dialog-content; upstream's
		// separate overlay rule's utilities land on the content rule's own
		// ::backdrop, matching registry/styles/nova/dialog.css's existing
		// `backdrop:*` utilities exactly.
		"cn-dialog-overlay": {slot: "content", ok: true, variant: "backdrop"},
		// Upstream spells the close button's own rule "cn-dialog-close";
		// gsxui's shape calls the same slot "close-button" (paired with a
		// separate "close-icon"). Same concept (registry/styles/nova/dialog.css's
		// close-button rule is also absolute-positioned chrome), different name.
		"cn-dialog-close": {slot: "close-button", ok: true},
	},
	"alert-dialog": {
		// See the package doc comment above slotOverrides: composes
		// DialogContent, inherits its ::backdrop, declares none of its own.
		"cn-alert-dialog-overlay": {ok: false},
		// AlertDialog adds an icon/image "media" area shadcn's upstream
		// design ships that gsxui's shape (registry/canonical/shapes/alert-dialog.go:
		// content/header/footer/title/description only) doesn't declare at
		// all — a missing feature, not a naming gap.
		"cn-alert-dialog-media": {ok: false},
	},
	"drawer": {
		"cn-drawer-overlay": {slot: "content", ok: true, variant: "backdrop"},
		// Vaul's real Drawer renders a "popup" wrapper INSIDE content with
		// swipe-gesture support (group-data-[swipe-axis]/[swipe-direction]);
		// gsxui's simpler native-<dialog>-based Drawer
		// (registry/canonical/shapes/drawer.go) has no popup/swipe-handle
		// concept and no swipe gestures at all — a permanent structural
		// difference. cn-drawer-header-base/-footer-base are a "base" (no
		// bottom padding) drawer-style variant gsxui doesn't distinguish
		// from its one header/footer.
		"cn-drawer-popup":        {ok: false},
		"cn-drawer-swipe-handle": {ok: false},
		"cn-drawer-header-base":  {ok: false},
		"cn-drawer-footer-base":  {ok: false},
	},
	"sheet": {
		"cn-sheet-overlay": {slot: "content", ok: true, variant: "backdrop"},
		// Same naming mismatch as dialog's cn-dialog-close above.
		"cn-sheet-close": {slot: "close-button", ok: true},
	},
	// Components whose upstream section carries a bare root rule
	// (`cn-<component>`) but whose gsxui shape declares no root slot at all:
	// the root element is purely structural (no recipe class of its own),
	// matching each shape's own doc comment and confirmed against
	// registry/styles/nova/<component>.css, which likewise has no base rule
	// for it. Not a porting gap — the same architecture decision Accordion's
	// shape doc comment already explains.
	"accordion": {
		"cn-accordion": {ok: false},
	},
	// Button's root recipe rule needs content from TWO upstream files, not
	// one: style-<name>.css's own per-style `.cn-button` rule (theme colors,
	// radius, border — genuinely different per style, e.g. maia's
	// rounded-4xl vs nova's rounded-lg) PLUS a STYLE-INVARIANT structural
	// base string hardcoded directly in the shared Radix component
	// (apps/v4/registry/bases/radix/ui/button.tsx's cva() call: "cn-button
	// group/button inline-flex shrink-0 items-center justify-center
	// whitespace-nowrap transition-all outline-none select-none
	// disabled:pointer-events-none disabled:opacity-50
	// [&_svg]:pointer-events-none [&_svg]:shrink-0" — identical across all 8
	// styles, confirmed by reading that file directly). style-maia.css's own
	// Button MARK section genuinely never mentions these utilities at all —
	// they live outside the porter's read scope (style-<name>.css only, per
	// the plan), the same way the original hand-authored
	// registry/styles/maia/button.css needed this file as a second citation
	// ("paired with apps/v4/registry/bases/radix/ui/button.tsx"). extra
	// appends this exact, verified-invariant string (minus cn-button and
	// group/button, which are markers, not utilities) after upstream's own
	// contribution — confirmed to introduce no duplicate or Tailwind
	// conflict with maia's own rule.
	"button": {
		"cn-button": {slot: "", ok: true, extra: []string{
			"inline-flex", "shrink-0", "items-center", "justify-center",
			"whitespace-nowrap", "transition-all", "outline-none", "select-none",
			"disabled:pointer-events-none", "disabled:opacity-50",
			"[&_svg]:pointer-events-none", "[&_svg]:shrink-0",
		}},
	},
	// Radio's shape (registry/canonical/shapes/radio.go) is ONE native
	// <input type=radio> root slot — no separate item/indicator elements,
	// unlike upstream's Radix-based RadioGroupItem (a styled wrapper) plus a
	// child indicator span. Upstream's naming reflects that split
	// (radio-group is the fieldset container, radio-group-item is the
	// individual control, radio-group-indicator is its child dot), but
	// gsxui's default derivation of "cn-radio-group-item" resolves to a
	// nonexistent "item" slot (SlotFor's upstreamPrefix quirk strips only
	// "radio-group-", leaving "item"). The item's own state styling
	// (checked/focus/aria-invalid) IS what gsxui's bare native input needs —
	// confirmed identical in kind to registry/styles/nova/radio.css's own
	// root rule — so it's overridden onto root instead of left unmapped.
	// The container (cn-radio-group, RadioGroup's own layout, a different
	// composition entirely) and the indicator/indicator-icon (no child
	// element on a native input to receive them; the checked dot is a
	// radial-gradient background-image escape hatch per nova/radio.css's own
	// header) have no destination.
	"radio": {
		"cn-radio-group":                {ok: false},
		"cn-radio-group-item":           {slot: "", ok: true},
		"cn-radio-group-indicator":      {ok: false},
		"cn-radio-group-indicator-icon": {ok: false},
	},
	// table-header/table-body render as plain <thead>/<tbody> with no recipe
	// class of their own — confirmed by registry/styles/nova/table.css's own
	// header ("table-header/table-body carry no rule of their own"); other
	// slots (row) reach INTO them via arbitrary ancestor selectors instead
	// (`[[data-gsxui-slot-table-body]_&]:last:border-b-0`). Not a porting
	// gap, the pre-existing architecture.
	"table": {
		"cn-table-header": {ok: false},
		"cn-table-body":   {ok: false},
	},
	// Switch is a single native <input type=checkbox role=switch> root slot
	// (registry/canonical/shapes/switch.go); there is no separate thumb
	// element for a "thumb" slot to target. gsxui already paints the thumb
	// as root's own `::before` pseudo-element (nova/switch.css's
	// `before:content-[''] before:...`), so upstream's separate
	// cn-switch-thumb rule fuses onto root with a `before:` variant inserted
	// before each utility — the same insertVariant mechanism the dialog/
	// drawer/sheet `::backdrop` fusions use, just a different pseudo-element.
	"switch": {
		"cn-switch-thumb": {slot: "", ok: true, variant: "before"},
	},
	// Slider's upstream Radix implementation renders track/range/thumb as
	// real child <div>s; gsxui's Slider (registry/canonical/shapes/slider.go)
	// is a single native <input type=range> styled entirely through
	// vendor-prefixed pseudo-element selectors on its own root rule
	// (nova/slider.css's `[&::-webkit-slider-thumb]:`/`[&::-moz-range-thumb]:`
	// arbitrary variants) — there are no track/range/thumb DOM nodes at all
	// to receive a separate rule. A permanent structural difference, not a
	// naming gap.
	"slider": {
		"cn-slider-track": {ok: false},
		"cn-slider-range": {ok: false},
		"cn-slider-thumb": {ok: false},
	},
	// cn-progress-track duplicates cn-progress's own bare utilities verbatim
	// (both "bg-muted h-3 rounded-4xl") — upstream ships two DOM strategies
	// (native <progress> vs a custom ARIA bar) that gsxui's single native
	// element doesn't need twice; root already carries it via the ordinary
	// bare-class match. cn-progress-label/-value are a percentage-text
	// readout gsxui's minimal Progress (registry/canonical/shapes/progress.go:
	// root + indicator only) doesn't render at all — confirmed absent from
	// registry/styles/nova/progress.css too.
	"progress": {
		"cn-progress-track": {ok: false},
		"cn-progress-label": {ok: false},
		"cn-progress-value": {ok: false},
	},
	// The `-content-logical`/`-arrow-logical` sibling rule every popover-
	// family component carries (Combobox, Context Menu, Dropdown Menu, Hover
	// Card, Menubar, Popover, Select, Tooltip) is an RTL-logical
	// (inline-start/inline-end) restatement of the SAME directional
	// slide-in/position transition its own main content/arrow rule already
	// expresses physically (data-side=left/right). None of gsxui's
	// popover-family recipes implement per-side directional slide at all
	// today (confirmed: registry/styles/nova/{select,tooltip,popover,...}.css
	// use a plain fade+scale transition, no data-[side=...]:slide-in-* of
	// any kind) — a pre-porter design choice this table documents once
	// rather than once per component.
	"combobox": {
		"cn-combobox-content-logical": {ok: false},
		// Multi-select "chip" display (tag pills for selected items) and a
		// separate item-text slot: gsxui's Combobox
		// (registry/canonical/shapes/combobox.go) is single-select with no
		// chip/item-text concept at all — a missing feature, not a naming
		// gap.
		"cn-combobox-item-text":   {ok: false},
		"cn-combobox-chips":       {ok: false},
		"cn-combobox-chip":        {ok: false},
		"cn-combobox-chip-remove": {ok: false},
	},
	"context-menu": {
		"cn-context-menu-content-logical": {ok: false},
		// item-indicator is ONE shared rule upstream uses for BOTH
		// checkbox-item-indicator and radio-item-indicator; gsxui declares
		// them as two separate slots, each already carrying equivalent
		// positioning in registry/styles/nova/context-menu.css (`absolute
		// end-2 ...`), so both fall back to nova cleanly rather than an
		// arbitrary pick of one. cn-context-menu-subcontent (no hyphen,
		// distinct from -sub-content) duplicates content already in
		// cn-context-menu-sub-content's own "shadow-lg" — a redundant
		// upstream rule.
		"cn-context-menu-item-indicator": {ok: false},
		"cn-context-menu-subcontent":     {ok: false},
	},
	"dropdown-menu": {
		"cn-dropdown-menu-content-logical": {ok: false},
		// Same shared-indicator/duplicate-subcontent shape as context-menu
		// above (dropdown-menu and context-menu are near-identical Radix
		// menus).
		"cn-dropdown-menu-item-indicator": {ok: false},
		"cn-dropdown-menu-subcontent":     {ok: false},
	},
	"hover-card": {"cn-hover-card-content-logical": {ok: false}},
	"menubar":    {"cn-menubar-content-logical": {ok: false}},
	"popover": {
		"cn-popover-content-logical": {ok: false},
		// Popover grew an optional header/title/description mini-dialog
		// composition upstream; gsxui's Popover
		// (registry/canonical/shapes/popover.go: root + content only) has no
		// such sub-structure — a missing feature.
		"cn-popover-header":      {ok: false},
		"cn-popover-title":       {ok: false},
		"cn-popover-description": {ok: false},
	},
	"tooltip": {
		"cn-tooltip-content-logical": {ok: false},
		"cn-tooltip-arrow-logical":   {ok: false},
	},
	// Alert's action-button slot (an inline dismiss/act control) has no
	// counterpart in registry/canonical/shapes/alert.go (root/title/
	// description only) — a feature upstream added that gsxui's Alert
	// doesn't render.
	"alert": {"cn-alert-action": {ok: false}},
	// Avatar's status-badge overlay and AvatarGroup's "+N" count indicator
	// are both upstream features with no slot in
	// registry/canonical/shapes/avatar.go (root/image/fallback only) —
	// confirmed absent from registry/styles/nova/avatar.css too.
	"avatar": {
		"cn-avatar-badge":       {ok: false},
		"cn-avatar-group-count": {ok: false},
	},
	// Calendar's month/year <select> dropdown mode and its caption-label
	// wrapper are deliberately absent from gsxui's Calendar (see
	// registry/styles/nova/calendar.css's own header) — a pre-existing,
	// documented gap, not a porter regression.
	"calendar": {
		"cn-calendar-dropdown-root": {ok: false},
		"cn-calendar-caption-label": {ok: false},
	},
	// Checkbox's checked glyph is a background-image data URI (see
	// registry/styles/nova/checkbox.css's own header), not an <svg> child —
	// there is no child element for an "indicator" rule's `[&>svg]:size-3.5`
	// to ever reach.
	"checkbox": {"cn-checkbox-indicator": {ok: false}},
	"command": {
		// cn-command-input-group/-input-icon restate content
		// registry/styles/nova/command.css's own input-wrapper rule already
		// carries (h-9, and `[&>svg]:size-4 [&>svg]:shrink-0 [&>svg]:opacity-50`
		// respectively) — merged onto input-wrapper rather than left
		// unmapped; merge-dedup (Transform's resolveMergeOrder) drops the
		// resulting literal h-9 repeat automatically.
		"cn-command-input-group": {slot: "input-wrapper", ok: true},
		"cn-command-input-icon":  {slot: "input-wrapper", ok: true, variant: "[&>svg]"},
	},
	// gsxui's InputGroup shape has one shared "control" slot for both
	// <input> and <textarea>, and NO "button" slot at all — confirmed by
	// registry/styles/nova/input-group.css's own header, which documents
	// leaving the button size ramp, the control's rounded-none/border-0, and
	// the textarea control's py-2 "in that same escape hatch" ON PURPOSE:
	// "Those paint an element another component owns (Button, Input,
	// Textarea) and are positioned there to beat that component's own
	// migrated utilities; re-expressing them here would put a working
	// contest back in play for nothing."
	"input-group": {
		"cn-input-group-button":              {ok: false},
		"cn-input-group-button-size-xs":      {ok: false},
		"cn-input-group-button-size-icon-xs": {ok: false},
		"cn-input-group-button-size-icon-sm": {ok: false},
		"cn-input-group-input":               {ok: false},
		"cn-input-group-textarea":            {ok: false},
	},
	// A decorative animated blinking caret gsxui's simpler InputOTP doesn't
	// render (registry/canonical/shapes/input-otp.go has no caret-line slot).
	"input-otp": {"cn-input-otp-caret-line": {ok: false}},
	// item's shape (registry/canonical/shapes/item.go) declares a "variant"
	// dimension on its root but no "size" dimension at all — upstream's
	// size axis has nothing to decode onto without a shape change, which is
	// out of the porter's scope (dossier: styles never add/remove shape
	// declarations).
	"item": {
		"cn-item-size-default": {ok: false},
		"cn-item-size-sm":      {ok: false},
		"cn-item-size-xs":      {ok: false},
	},
	// The chevron icon is styled via an arbitrary `[&>svg]:` child selector
	// on the wrapper slot (registry/styles/nova/native-select.css), not a
	// separate icon element — same insertVariant fusion mechanism as
	// switch's thumb.
	"native-select": {"cn-native-select-icon": {slot: "wrapper", ok: true, variant: "[&>svg]"}},
	// viewport/positioner/popup are Radix's portal-based rendering mode for
	// NavigationMenu; gsxui's simpler NavigationMenu
	// (registry/canonical/shapes/navigation-menu.go) has no such concept —
	// confirmed absent from registry/styles/nova/navigation-menu.css too.
	"navigation-menu": {
		"cn-navigation-menu-viewport":   {ok: false},
		"cn-navigation-menu-positioner": {ok: false},
		"cn-navigation-menu-popup":      {ok: false},
	},
	// gsxui's Pagination shape only declares previous-label/next-label
	// (sr-only accessible text spans); the previous/next control itself
	// composes <ui.Button> directly and takes Button's own recipe, not a
	// Pagination-owned rule.
	"pagination": {
		"cn-pagination-previous": {ok: false},
		"cn-pagination-next":     {ok: false},
	},
	// Same naming mismatch as native-select's icon: gsxui's shape calls this
	// slot "handle-grip" (registry/canonical/shapes/resizable.go); content is
	// otherwise identical in kind to registry/styles/nova/resizable.css's own
	// handle-grip rule.
	"resizable": {"cn-resizable-handle-icon": {slot: "handle-grip", ok: true}},
	// Radix's ScrollArea renders real scrollbar/thumb child <div>s; gsxui's
	// ScrollArea (registry/canonical/shapes/scroll-area.go, one root slot)
	// uses the native scrollbar via plain CSS scrollbar-width/-color
	// properties (registry/styles/nova/scroll-area.css's own escape-hatch
	// comment) — no scrollbar/thumb DOM nodes exist to receive a rule.
	"scroll-area": {
		"cn-scroll-area-scrollbar": {ok: false},
		"cn-scroll-area-thumb":     {ok: false},
	},
	// sidebar's menu-button slot (registry/canonical/shapes/sidebar.go)
	// declares no variant/size dimensions at all — confirmed
	// registry/styles/nova/sidebar.css's own menu-button rule has no such
	// distinction either, a pre-existing shape limitation outside the
	// porter's scope to fix (would require a shape change).
	"sidebar": {
		"cn-sidebar-menu-button-variant-default": {ok: false},
		"cn-sidebar-menu-button-variant-outline": {ok: false},
		"cn-sidebar-menu-button-size-default":    {ok: false},
		"cn-sidebar-menu-button-size-sm":         {ok: false},
		"cn-sidebar-menu-button-size-lg":         {ok: false},
		// menu-action's "reveal only on hover/focus/open" interaction
		// (opacity 0 by default on desktop, opacity-100 when the menu item
		// is hovered/focused or the action's own popover is open) is
		// interaction behavior with no counterpart anywhere in upstream's
		// Sidebar MARK section at all — confirmed absent from
		// cn-sidebar-menu-action's own upstream rule. It is style-invariant
		// (an interaction pattern, not a themed value), so this appends
		// registry/styles/nova/sidebar.css's own already-verified rules
		// verbatim, required by TestSidebarRecipeUsesPresenceSelectors.
		"cn-sidebar-menu-action": {slot: "menu-action", ok: true, extra: []string{
			"md:[&[data-show-on-hover]]:opacity-0",
			"md:[[data-gsxui-slot-sidebar-menu-item]:hover>&[data-show-on-hover]]:opacity-100",
			"md:[[data-gsxui-slot-sidebar-menu-item]:focus-within>&[data-show-on-hover]]:opacity-100",
			"md:[&[data-show-on-hover][data-state=open]]:opacity-100",
			"[[data-gsxui-slot-sidebar-menu-button][data-active]~&[data-show-on-hover]]:text-sidebar-primary-foreground",
		}},
	},
}

// classRenames maps (component, upstream class) to a DIFFERENT upstream
// class name to resolve the rule AS IF it had been spelled that way —
// applied before any other resolution (SlotFor, decodeUpstreamValue,
// overrideVariant). This is for a naming mismatch a plain slotOverride can't
// express because it spans MULTIPLE upstream rules that must decode as one
// base rule PLUS its dimension values: Empty calls the icon slot "media"
// where gsxui's shape (registry/canonical/shapes/empty.go) calls it "icon",
// including the variant values ("media-default"/"media-icon" ->
// "icon-variant-default"/"icon-variant-icon", confirmed identical in kind to
// registry/styles/nova/empty.css's own icon-variant-* rules). Command's
// "cn-command-dialog" needs the same treatment against gsxui's
// "dialog-content" slot name.
var classRenames = map[string]map[string]string{
	"empty": {
		"cn-empty-media":         "cn-empty-icon",
		"cn-empty-media-default": "cn-empty-icon-variant-default",
		"cn-empty-media-icon":    "cn-empty-icon-variant-icon",
	},
	"command": {
		"cn-command-dialog": "cn-command-dialog-content",
	},
}

// renameClass substitutes a declared classRenames entry, or returns class
// unchanged when none applies.
func renameClass(component, class string) string {
	if renames, ok := classRenames[component]; ok {
		if renamed, ok := renames[class]; ok {
			return renamed
		}
	}
	return class
}

// markerRewrites declares upstream bare data-boolean variants that must be
// rewritten to gsxui's own explicit presence-selector form (`[&[<marker>]]:`)
// for a specific component. Every OTHER bare data-boolean variant elsewhere
// in the porter (data-open, data-closed, data-selected, ...) passes through
// unmodified — that is the documented, normal treatment (dossier §5.4: bare
// boolean data variants are runtime state, travel verbatim). Sidebar's
// data-active/data-show-on-hover are the one confirmed exception:
// TestSidebarRecipeUsesPresenceSelectors (internal/stylegen/button_css_test.go)
// requires exactly the bracket form, because these two markers are stamped
// with no value at all — matching them with Tailwind's bare data-active:
// shorthand does not reliably resolve to a presence check for this repo's
// Tailwind config the way the explicit `[&[data-active]]:` form does
// (already verified against the real Tailwind CLI by nova/sidebar.css,
// predating the porter). This table is deliberately narrow: it is not a
// blanket claim that bare data-*: forms are unsafe in general.
var markerRewrites = map[string]map[string]string{
	"sidebar": {
		"data-active":        "[&[data-active]]",
		"data-show-on-hover": "[&[data-show-on-hover]]",
	},
}

// rewriteMarkerVariant replaces token's leading variant segment with its
// markerRewrites substitute when one is declared for component, leaving
// every other stacked variant and the base utility untouched. Returns token
// unmodified when no rewrite applies.
func rewriteMarkerVariant(component, token string) string {
	rewrites, ok := markerRewrites[component]
	if !ok {
		return token
	}
	segments := splitVariants(token)
	rewritten, ok := rewrites[segments[0]]
	if !ok {
		return token
	}
	segments[0] = rewritten
	return strings.Join(segments, ":")
}

// upstreamPrefix returns the upstream cn-* class prefix (without "cn-") a
// gsxui component's own rules appear under. Every component shares its name
// with the upstream prefix except radio, whose section is "Radio Group" (see
// sectionComponentOverrides): its rules are cn-radio-group, cn-radio-group-
// item, and so on, never cn-radio-*.
func upstreamPrefix(component string) string {
	if component == "radio" {
		return "radio-group"
	}
	return component
}

// matchedOverride looks up a declared slotOverrides entry for (component,
// class), if any. Shared by SlotFor (the public slot/ok view) and Transform
// (which additionally needs silentNonMapping and the variant field, neither
// of which SlotFor's two-return signature — already committed and tested —
// can carry).
func matchedOverride(component, class string) (slotOverride, bool) {
	overrides, has := slotOverrides[component]
	if !has {
		return slotOverride{}, false
	}
	override, matched := overrides[class]
	return override, matched
}

// silentNonMapping reports whether (component, class) is a class the mapping
// table has EXPLICITLY reviewed and routed nowhere on purpose (a declared
// slotOverride with ok=false) — as opposed to a class that simply doesn't
// belong to this component's upstream prefix at all (an unrecognized stray).
// Transform uses this to keep the former out of its unmapped report: a
// declared "this has no gsxui slot" decision is not an open question the
// report needs to keep surfacing every run.
func silentNonMapping(component, class string) bool {
	override, matched := matchedOverride(component, class)
	return matched && !override.ok
}

// overrideVariant returns the extra Tailwind variant a declared slotOverride
// wants inserted before every utility it contributes (e.g. "backdrop" for
// the overlay→content fusions), or "" when the override carries none.
func overrideVariant(component, class string) string {
	override, matched := matchedOverride(component, class)
	if !matched || !override.ok {
		return ""
	}
	return override.variant
}

// overrideExtra returns extra utilities a declared slotOverride wants
// appended, verbatim, after the rule's own upstream contribution — for
// content that is genuinely style-invariant but lives outside
// style-<name>.css entirely (Button's cva() base string, hardcoded in the
// shared Radix component; see the "button" entry in slotOverrides).
func overrideExtra(component, class string) []string {
	override, matched := matchedOverride(component, class)
	if !matched || !override.ok {
		return nil
	}
	return override.extra
}

// SlotFor maps an upstream class within a component to a gsxui slot name.
// Returns ok=false when the class has no gsxui slot (reported, not
// dropped): either by a declared override (e.g. select's scroll buttons) or
// because the class doesn't belong to this component's upstream prefix at
// all (an unrelated stray).
func SlotFor(component, class string) (slot string, ok bool) {
	if override, matched := matchedOverride(component, class); matched {
		return override.slot, override.ok
	}

	prefix := "cn-" + upstreamPrefix(component)
	if class == prefix {
		return "", true
	}
	suffix, cut := strings.CutPrefix(class, prefix+"-")
	if !cut {
		return "", false
	}
	return suffix, true
}
