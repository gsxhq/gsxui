package registry_test

import (
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/gsxhq/gsxui/internal/registry"
	"github.com/gsxhq/gsxui/internal/stylecontract"
)

func TestComponents(t *testing.T) {
	got, err := registry.Components()
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"accordion", "alert", "alert-dialog", "aspect-ratio", "avatar", "badge", "breadcrumb", "button", "button-group", "calendar", "card", "carousel", "checkbox", "collapsible", "combobox", "command", "context-menu", "dialog", "drawer", "dropdown", "empty", "field", "hover-card", "icon", "input", "input-group", "input-otp", "item", "kbd", "label", "menubar", "native-select", "navigation-menu", "pagination", "popover", "progress", "radio", "resizable", "scroll-area", "select", "separator", "sheet", "sidebar", "skeleton", "slider", "sonner", "spinner", "switch", "table", "tabs", "textarea", "toggle", "toggle-group", "tooltip"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
	for _, unwanted := range []string{"core", "gsxui", "index", "nativeselect", "switchctl"} {
		if slicesContains(got, unwanted) {
			t.Fatalf("Components() = %v, should not contain %q", got, unwanted)
		}
	}
}

func TestComponentsMatchTypedStyleContract(t *testing.T) {
	components, err := registry.Components()
	if err != nil {
		t.Fatal(err)
	}
	contractComponents := make([]string, 0, len(stylecontract.All()))
	for _, component := range stylecontract.All() {
		contractComponents = append(contractComponents, component.RegistryName)
	}
	// ui/sonner.gsx is one vendorable file backing two typed style-contract
	// components, "toaster" and "toast" (they're separate gsx components
	// rendered by gsxui itself, and the split lets each satisfy the
	// "<RegistryName>-<relative>" slot-naming rule with no exception — see
	// contracts_toast.go). registry.Components() stays file-derived, so fold
	// the pair back to the single file name "sonner" before comparing.
	foldedContractComponents := make([]string, 0, len(contractComponents))
	for _, name := range contractComponents {
		if name == "toaster" || name == "toast" {
			continue
		}
		foldedContractComponents = append(foldedContractComponents, name)
	}
	foldedContractComponents = append(foldedContractComponents, "sonner")
	slices.Sort(foldedContractComponents)
	if !slices.Equal(components, foldedContractComponents) {
		t.Fatalf("registry components = %v; typed style contract (toast/toaster folded to sonner) = %v", components, foldedContractComponents)
	}
}

func slicesContains(ss []string, s string) bool {
	for _, x := range ss {
		if x == s {
			return true
		}
	}
	return false
}

func TestDeps(t *testing.T) {
	// dialog.x.go references Button (DialogFooter's Close button) — an
	// intra-package edge with no import to scan, resolved via declIndex.
	deps, err := registry.Deps("dialog")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(deps, []string{"button"}) {
		t.Fatalf("dialog deps = %v, want [button]", deps)
	}

	// accordion.gsx imports ui/icon (AccordionTrigger's chevron).
	deps, err = registry.Deps("accordion")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(deps, []string{"icon"}) {
		t.Fatalf("accordion deps = %v, want [icon]", deps)
	}

	// native-select.gsx imports ui/icon (the chevron).
	deps, err = registry.Deps("native-select")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(deps, []string{"icon"}) {
		t.Fatalf("native-select deps = %v, want [icon]", deps)
	}

	// spinner.gsx imports ui/icon (icon.LoaderCircle).
	deps, err = registry.Deps("spinner")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(deps, []string{"icon"}) {
		t.Fatalf("spinner deps = %v, want [icon]", deps)
	}

	// breadcrumb.gsx imports ui/icon (BreadcrumbSeparator's default
	// ChevronRight, BreadcrumbEllipsis's Ellipsis/MoreHorizontal).
	deps, err = registry.Deps("breadcrumb")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(deps, []string{"icon"}) {
		t.Fatalf("breadcrumb deps = %v, want [icon]", deps)
	}

	// carousel.gsx composes Button (CarouselPrevious/CarouselNext) — an
	// intra-package edge with no import to scan, same resolution shape as
	// dialog's own Deps entry above — AND imports ui/icon
	// (CarouselPrevious/CarouselNext's ArrowLeft/ArrowRight), the ordinary
	// house default (accordion/breadcrumb/pagination/spinner all do the
	// same). Deps sorts its result, so button < icon alphabetically.
	deps, err = registry.Deps("carousel")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(deps, []string{"button", "icon"}) {
		t.Fatalf("carousel deps = %v, want [button icon]", deps)
	}

	// kbd.gsx, aspect-ratio.gsx, and progress.gsx have no icon import and no
	// intra-package reference to another component.
	deps, err = registry.Deps("kbd")
	if err != nil {
		t.Fatal(err)
	}
	if len(deps) != 0 {
		t.Fatalf("kbd deps = %v, want none", deps)
	}

	deps, err = registry.Deps("aspect-ratio")
	if err != nil {
		t.Fatal(err)
	}
	if len(deps) != 0 {
		t.Fatalf("aspect-ratio deps = %v, want none", deps)
	}

	deps, err = registry.Deps("progress")
	if err != nil {
		t.Fatal(err)
	}
	if len(deps) != 0 {
		t.Fatalf("progress deps = %v, want none", deps)
	}

	// resizable.gsx imports nothing from ui/icon — nova's own withHandle
	// grip is an empty pill div, not new-york-v4's GripVerticalIcon-in-a-box
	// (2026-07-24 tier4 source map, `## resizable` §6) — and composes no
	// other component, so it has no deps at all.
	deps, err = registry.Deps("resizable")
	if err != nil {
		t.Fatal(err)
	}
	if len(deps) != 0 {
		t.Fatalf("resizable deps = %v, want none", deps)
	}

	deps, err = registry.Deps("button")
	if err != nil {
		t.Fatal(err)
	}
	if len(deps) != 0 {
		t.Fatalf("button deps = %v, want none", deps)
	}

	// pagination.gsx imports ui/icon (ChevronLeft/ChevronRight/Ellipsis).
	// Its Button relationship is now token composition in the shared style
	// contract, not a Go code dependency for the vendored source graph.
	deps, err = registry.Deps("pagination")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(deps, []string{"icon"}) {
		t.Fatalf("pagination deps = %v, want [icon]", deps)
	}

	// button-group.gsx has no icon import; ButtonGroupSeparator calls
	// ui.Separator directly (flat package intra-package edge, same shape as
	// dialog's button dep).
	deps, err = registry.Deps("button-group")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(deps, []string{"separator"}) {
		t.Fatalf("button-group deps = %v, want [separator]", deps)
	}

	deps, err = registry.Deps("icon")
	if err != nil {
		t.Fatal(err)
	}
	if len(deps) != 0 {
		t.Fatalf("icon deps = %v, want none", deps)
	}

	// empty.gsx has no icon import and no intra-package reference to
	// another component.
	deps, err = registry.Deps("empty")
	if err != nil {
		t.Fatal(err)
	}
	if len(deps) != 0 {
		t.Fatalf("empty deps = %v, want none", deps)
	}

	// item.gsx has no icon import; ItemSeparator calls ui.Separator directly
	// (flat package intra-package edge, same shape as button-group's own
	// separator dep).
	deps, err = registry.Deps("item")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(deps, []string{"separator"}) {
		t.Fatalf("item deps = %v, want [separator]", deps)
	}

	// input-group.gsx has no icon import; InputGroupButton/InputGroupInput/
	// InputGroupTextarea call ui.Button/ui.Input/ui.Textarea directly (flat
	// package intra-package edges, same shape as dialog's button dep).
	deps, err = registry.Deps("input-group")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(deps, []string{"button", "input", "textarea"}) {
		t.Fatalf("input-group deps = %v, want [button input textarea]", deps)
	}

	// field.gsx has no icon import; FieldLabel calls ui.Label and
	// FieldSeparator calls ui.Separator directly (flat package intra-package
	// edges, same shape as item's own separator dep).
	deps, err = registry.Deps("field")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(deps, []string{"label", "separator"}) {
		t.Fatalf("field deps = %v, want [label separator]", deps)
	}

	// collapsible.gsx has no icon import and no intra-package reference to
	// another component (the site example composes ui.Button/ui/icon, but
	// internal/registry only scans ui/*.gsx, not site/examples/).
	deps, err = registry.Deps("collapsible")
	if err != nil {
		t.Fatal(err)
	}
	if len(deps) != 0 {
		t.Fatalf("collapsible deps = %v, want none", deps)
	}

	// alert-dialog.gsx has no icon import; AlertDialog composes ui.Dialog
	// and AlertDialogAction/AlertDialogCancel compose ui.Button directly
	// (flat package intra-package edges, same shape as dialog's own button
	// dep) — alert-dialog -> dialog is also what makes the CLI vendor
	// ui/dialog.js for alert-dialog (HasJS("alert-dialog") is false; it has
	// no behavior module of its own, only dialog's).
	deps, err = registry.Deps("alert-dialog")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(deps, []string{"button", "dialog"}) {
		t.Fatalf("alert-dialog deps = %v, want [button dialog]", deps)
	}

	// sheet.gsx has no icon import; Sheet composes ui.Dialog directly (flat
	// package intra-package edge, same shape as alert-dialog's own dialog
	// dep) — SheetContent renders its own <dialog> rather than composing
	// DialogContent, and SheetTrigger/SheetContent's injected close
	// button/SheetClose all render their own <button> rather than composing
	// Button, so dialog is the only edge — sheet -> dialog is also what
	// makes the CLI vendor ui/dialog.js for a sheet install (HasJS("sheet")
	// is false; it has no behavior module of its own, only dialog's).
	deps, err = registry.Deps("sheet")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(deps, []string{"dialog"}) {
		t.Fatalf("sheet deps = %v, want [dialog]", deps)
	}

	// drawer.gsx has no icon import; Drawer composes ui.Dialog directly (flat
	// package intra-package edge, same shape as sheet's own dialog dep) —
	// DrawerContent renders its own <dialog> rather than composing
	// DialogContent/SheetContent, and DrawerTrigger/DrawerClose render their
	// own <button> rather than composing Button, so dialog is the only edge
	// — drawer -> dialog is also what makes the CLI vendor ui/dialog.js for
	// a drawer install (HasJS("drawer") is false; it has no behavior module
	// of its own, only dialog's — same conclusion as sheet's own entry).
	deps, err = registry.Deps("drawer")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(deps, []string{"dialog"}) {
		t.Fatalf("drawer deps = %v, want [dialog]", deps)
	}

	// toggle.gsx has no icon import and no intra-package reference to
	// another component (the site example composes ui/icon, but
	// internal/registry only scans ui/*.gsx, not site/examples/ — same
	// shape as collapsible's own deps entry above).
	deps, err = registry.Deps("toggle")
	if err != nil {
		t.Fatal(err)
	}
	if len(deps) != 0 {
		t.Fatalf("toggle deps = %v, want none", deps)
	}

	// ToggleGroupItem composes Toggle's styling token, but the CSS pack owns
	// that token's declarations and toggle-group.js owns group interaction.
	// There is no source or behavior dependency on toggle.gsx.
	deps, err = registry.Deps("toggle-group")
	if err != nil {
		t.Fatal(err)
	}
	if len(deps) != 0 {
		t.Fatalf("toggle-group deps = %v, want none", deps)
	}

	// popover.gsx has no icon import and no intra-package reference to
	// another component — Popover/PopoverTrigger/PopoverContent are all
	// plain elements, same shape as toggle's own deps entry.
	deps, err = registry.Deps("popover")
	if err != nil {
		t.Fatal(err)
	}
	if len(deps) != 0 {
		t.Fatalf("popover deps = %v, want none", deps)
	}

	// hover-card.gsx has no icon import and no intra-package reference to
	// another component (the site example composes ui.Avatar/ui.Button, but
	// internal/registry only scans ui/*.gsx, not site/examples/ — same
	// shape as collapsible's/toggle's own deps entries).
	deps, err = registry.Deps("hover-card")
	if err != nil {
		t.Fatal(err)
	}
	if len(deps) != 0 {
		t.Fatalf("hover-card deps = %v, want none", deps)
	}

	// context-menu.gsx imports ui/icon (Tier 4 Batch B Task 2:
	// ContextMenuCheckboxItem's Check, ContextMenuRadioItem's Circle,
	// ContextMenuSubTrigger's ChevronRight) — same shape as dropdown.gsx's
	// own equivalent icon usage, no intra-package ui.* reference otherwise.
	deps, err = registry.Deps("context-menu")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(deps, []string{"icon"}) {
		t.Fatalf("context-menu deps = %v, want [icon]", deps)
	}

	// menubar.gsx imports ui/icon (Tier 4 Batch B Task 3: MenubarCheckboxItem's
	// Check, MenubarRadioItem's Circle, MenubarSubTrigger's ChevronRight) —
	// same shape as dropdown.gsx's/context-menu.gsx's own equivalent icon
	// usage, no intra-package ui.* reference otherwise.
	deps, err = registry.Deps("menubar")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(deps, []string{"icon"}) {
		t.Fatalf("menubar deps = %v, want [icon]", deps)
	}

	// navigation-menu.gsx imports ui/icon (Tier 4 Batch B Task 4:
	// NavigationMenuTrigger's ChevronDown) — same shape as dropdown.gsx's/
	// context-menu.gsx's/menubar.gsx's own equivalent icon usage, no
	// intra-package ui.* reference otherwise.
	deps, err = registry.Deps("navigation-menu")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(deps, []string{"icon"}) {
		t.Fatalf("navigation-menu deps = %v, want [icon]", deps)
	}

	// slider.gsx has no icon import and no intra-package reference to
	// another component (the site example composes nothing from another
	// ui.* component either) — same shape as toggle's/popover's/hover-
	// card's own deps entries above.
	deps, err = registry.Deps("slider")
	if err != nil {
		t.Fatal(err)
	}
	if len(deps) != 0 {
		t.Fatalf("slider deps = %v, want none", deps)
	}

	// scroll-area.gsx has no icon import and no intra-package reference to
	// another component (the site examples compose ui.Separator, but
	// internal/registry only scans ui/*.gsx, not site/examples/) — same
	// shape as toggle's/popover's/hover-card's/slider's own deps entries
	// above.
	deps, err = registry.Deps("scroll-area")
	if err != nil {
		t.Fatal(err)
	}
	if len(deps) != 0 {
		t.Fatalf("scroll-area deps = %v, want none", deps)
	}

	// input-otp.gsx imports ui/icon (InputOTPSeparator's icon.Minus) — same
	// shape as accordion's/native-select's/spinner's/breadcrumb's own deps
	// entries above; no intra-package reference to another component (it does NOT
	// compose ui.Input, see ui/input-otp.gsx's own doc comment).
	deps, err = registry.Deps("input-otp")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(deps, []string{"icon"}) {
		t.Fatalf("input-otp deps = %v, want [icon]", deps)
	}

	// select.gsx imports ui/icon (SelectTrigger's ChevronDown and SelectItem's
	// Check indicator) — same dependency-derivation shape as native-select's
	// own icon edge; the custom listbox is the second select -> icon edge the
	// controls map called out.
	deps, err = registry.Deps("select")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(deps, []string{"icon"}) {
		t.Fatalf("select deps = %v, want [icon]", deps)
	}

	// sonner.gsx imports ui/icon: the server-rendered ui.Toast card (the
	// single source of the toast <li> markup, shipped as inert per-type
	// <template>s and cloned by ui/sonner.js) renders its type glyph and the
	// close X via icon.* Go calls — so Deps is [icon], no longer empty.
	deps, err = registry.Deps("sonner")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(deps, []string{"icon"}) {
		t.Fatalf("sonner deps = %v, want [icon]", deps)
	}

	// combobox.gsx imports ui/icon (ComboboxItem's Check, ComboboxTrigger's
	// ChevronDown, ComboboxClear's X) and composes ui.InputGroup/
	// InputGroupInput/InputGroupAddon/InputGroupButton directly (flat
	// package intra-package edges, same declIndex-resolved shape as
	// input-group's own button/input/textarea deps above) — Deps sorts its
	// result, so icon < input-group alphabetically.
	deps, err = registry.Deps("combobox")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(deps, []string{"icon", "input-group"}) {
		t.Fatalf("combobox deps = %v, want [icon input-group]", deps)
	}

	// sidebar.gsx imports ui/icon (SidebarTrigger's PanelLeft) and composes
	// ui.Button (SidebarTrigger), ui.Input (SidebarInput), ui.Separator
	// (SidebarSeparator), ui.Sheet/SheetContent/SheetHeader/SheetTitle/
	// SheetDescription (Sidebar's own mobile tree), ui.Skeleton
	// (SidebarMenuSkeleton), and ui.Tooltip/TooltipContent
	// (SidebarMenuButton's tooltip branch) directly — flat package
	// intra-package edges, same declIndex-resolved shape as combobox's own
	// input-group/icon deps above. Deps sorts its result.
	deps, err = registry.Deps("sidebar")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(deps, []string{"button", "icon", "input", "separator", "sheet", "skeleton", "tooltip"}) {
		t.Fatalf("sidebar deps = %v, want [button icon input separator sheet skeleton tooltip]", deps)
	}

	// Calendar's controls compose Button's public styling token in markup
	// and CSS, not Button's Go implementation. Its only code dependencies
	// are ui/icon (nav chevrons) and NativeSelect/NativeSelectOption.
	deps, err = registry.Deps("calendar")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(deps, []string{"icon", "native-select"}) {
		t.Fatalf("calendar deps = %v, want [icon native-select]", deps)
	}

	if _, err := registry.Deps("nosuch"); err == nil || !strings.Contains(err.Error(), "gsxui list") {
		t.Fatalf("Deps(nosuch) err = %v, want error mentioning 'gsxui list'", err)
	}
	if _, err := registry.Deps("core"); err == nil || !strings.Contains(err.Error(), "gsxui list") {
		t.Fatalf("Deps(core) err = %v, want error mentioning 'gsxui list'", err)
	}
}

func TestHasJS(t *testing.T) {
	if !registry.HasJS("dropdown") {
		t.Error("dropdown should have JS")
	}
	if registry.HasJS("button") {
		t.Error("button should not have JS")
	}
	if registry.HasJS("gsxui") {
		t.Error("gsxui should not have JS")
	}
	// alert-dialog has no ui/alert-dialog.js of its own — it reuses
	// ui/dialog.js entirely (the data-gsxui-dialog-static opt-out lives in
	// dialog.js itself); vendoring dialog.js comes from the derived
	// alert-dialog -> dialog dependency (see TestDeps), not from HasJS here.
	if registry.HasJS("alert-dialog") {
		t.Error("alert-dialog should not have its own JS")
	}
	// sheet has no ui/sheet.js of its own — same reuse-dialog.js shape as
	// alert-dialog (see TestDeps' sheet entry).
	if registry.HasJS("sheet") {
		t.Error("sheet should not have its own JS")
	}
	// drawer has no ui/drawer.js of its own — same reuse-dialog.js shape as
	// sheet/alert-dialog (see TestDeps' drawer entry).
	if registry.HasJS("drawer") {
		t.Error("drawer should not have its own JS")
	}
	// toggle has its own ui/toggle.js (click flips aria-pressed/data-state).
	if !registry.HasJS("toggle") {
		t.Error("toggle should have JS")
	}
	// toggle-group has its own ui/toggle-group.js (roving tabindex, arrow-key
	// nav, click activation) — a separate behavior module from toggle.js
	// despite the toggle-group -> toggle CLASS dependency above; the two
	// components' interaction models don't overlap enough to share JS.
	if !registry.HasJS("toggle-group") {
		t.Error("toggle-group should have JS")
	}
	// popover has its own ui/popover.js (anchored positioning + state/aria
	// sync, adapted from dropdown.js).
	if !registry.HasJS("popover") {
		t.Error("popover should have JS")
	}
	// hover-card has its own ui/hover-card.js — HasJS derives from
	// <basename>.js, so the file is named ui/hover-card.js (hyphenated,
	// matching the component basename) even though the site example
	// package directory strips the hyphen to "hovercard" (Go package name
	// constraint, same nativeselect/switchctl precedent).
	if !registry.HasJS("hover-card") {
		t.Error("hover-card should have JS")
	}
	// context-menu has its own ui/context-menu.js (cursor-positioned open on
	// contextmenu, adapted from dropdown.js's menu semantics).
	if !registry.HasJS("context-menu") {
		t.Error("context-menu should have JS")
	}
	// menubar has its own ui/menubar.js — roving tabindex across triggers
	// (ui/toggle-group.js's own JS-normalized-at-init model) and
	// open-follows-hover once one menu is open, layered on dropdown.js's
	// reused item/submenu machinery (duplicated with menubar's own
	// data-gsxui-menubar-* selectors, not imported — no shared JS module,
	// same MECHANISM as context-menu.js's own).
	if !registry.HasJS("menubar") {
		t.Error("menubar should have JS")
	}
	// navigation-menu has its own ui/navigation-menu.js — hover/focus open
	// with a hover-card-shaped grace-period close, and the shared viewport's
	// own discrete-measurement + ResizeObserver sizing.
	if !registry.HasJS("navigation-menu") {
		t.Error("navigation-menu should have JS")
	}
	// slider has its own ui/slider.js (delegated `input` listener that
	// resyncs the --fill custom property while the user drags/keys the
	// thumb — the server-rendered initial --fill needs no JS at all).
	if !registry.HasJS("slider") {
		t.Error("slider should have JS")
	}
	// carousel has its own ui/carousel.js (prev/next scroll-by-one-item,
	// scroll-driven disabled-state/current-index bookkeeping, ArrowLeft/
	// ArrowRight keyboard, autoplay) — real new interactive JS, unlike
	// sheet/alert-dialog/drawer's own dialog.js reuse.
	if !registry.HasJS("carousel") {
		t.Error("carousel should have JS")
	}
	// resizable has its own ui/resizable.js (pointer drag + keyboard
	// step/Home/End resizing of the two panels adjacent to a handle,
	// aria-valuenow/-min/-max sync, gsxui:change on commit) — real new
	// interactive JS, since react-resizable-panels is absent from the
	// reference checkout entirely (nothing to reuse from a dist build that
	// was never read).
	if !registry.HasJS("resizable") {
		t.Error("resizable should have JS")
	}
	// input-otp has its own ui/input-otp.js (the entire hidden-single-input
	// mechanism: DOM-order data-index stamping, char/data-active/fake-caret
	// recompute on input/selectionchange/focus/blur, per-character pattern
	// filtering, slot-click-to-position).
	if !registry.HasJS("input-otp") {
		t.Error("input-otp should have JS")
	}
	// select has its own ui/select.js — the entire custom-listbox behavior
	// (value model, focus-aware aria-selected recompute, bespoke prefix
	// typeahead, hidden-select form bridge) layered on dropdown.js's reused
	// popover machinery; a real new interactive module of its own, unlike the
	// native-select (which has no JS at all).
	if !registry.HasJS("select") {
		t.Error("select should have JS")
	}
	if registry.HasJS("native-select") {
		t.Error("native-select should not have JS")
	}
	// sonner has its own ui/sonner.js — it clones Toaster's server-rendered
	// per-type <template>s into live toast <li>s and owns the stacking/
	// timer/pause-on-hover/promise-morph lifecycle plus adoption of
	// server-inserted rows. ui/sonner.gsx (Toaster + ui.Toast) is the
	// single source of the card markup; HasJS derives from the
	// <basename>.js match, so the file is ui/sonner.js.
	if !registry.HasJS("sonner") {
		t.Error("sonner should have JS")
	}
	// combobox has its own ui/combobox.js — filter-as-you-type (hide/show,
	// no reordering; see its own header ADAPT), a data-highlighted +
	// aria-activedescendant highlight cursor (command.js's focus-stays-in-
	// the-input model), and select.js's own value model/form bridge/
	// popover machinery restated for an input trigger. A real behavior
	// module of its own, not a reuse of select.js or command.js (## combobox
	// GAP: no shared JS with command, registry.Deps reasoning).
	if !registry.HasJS("combobox") {
		t.Error("combobox should have JS")
	}
	// sidebar has its own ui/sidebar.js — desktop data-state/data-collapsible
	// toggling, mobile-vs-desktop resolution (getComputedStyle, not a
	// hard-coded breakpoint), and the Cmd/Ctrl+B shortcut; a real behavior
	// module of its own, not a reuse of dialog.js (though the mobile tree's
	// own Sheet does ride dialog.js transitively via the sheet -> dialog
	// dependency above).
	if !registry.HasJS("sidebar") {
		t.Error("sidebar should have JS")
	}
}

func TestResolveTransitive(t *testing.T) {
	got, err := registry.Resolve([]string{"dialog"})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"button", "dialog"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}

	got, err = registry.Resolve([]string{"native-select"})
	if err != nil {
		t.Fatal(err)
	}
	want = []string{"icon", "native-select"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}

	// alert-dialog resolves transitively through dialog to button — this is
	// the chain that makes the CLI vendor ui/dialog.js (HasJS("dialog")) for
	// an alert-dialog install even though alert-dialog has no JS of its own.
	got, err = registry.Resolve([]string{"alert-dialog"})
	if err != nil {
		t.Fatal(err)
	}
	want = []string{"alert-dialog", "button", "dialog"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}

	// sheet resolves transitively through dialog to button — the same
	// vendoring chain as alert-dialog's own (TestDeps' sheet entry).
	got, err = registry.Resolve([]string{"sheet"})
	if err != nil {
		t.Fatal(err)
	}
	want = []string{"button", "dialog", "sheet"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}

	// drawer resolves transitively through dialog to button — the same
	// vendoring chain as sheet's/alert-dialog's own (TestDeps' drawer
	// entry).
	got, err = registry.Resolve([]string{"drawer"})
	if err != nil {
		t.Fatal(err)
	}
	want = []string{"button", "dialog", "drawer"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}

	// toggle has no deps of its own — Resolve returns just itself.
	got, err = registry.Resolve([]string{"toggle"})
	if err != nil {
		t.Fatal(err)
	}
	want = []string{"toggle"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}

	// popover and hover-card have no deps of their own — Resolve returns
	// just themselves, same shape as toggle's own entry above.
	got, err = registry.Resolve([]string{"popover"})
	if err != nil {
		t.Fatal(err)
	}
	want = []string{"popover"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}

	got, err = registry.Resolve([]string{"hover-card"})
	if err != nil {
		t.Fatal(err)
	}
	want = []string{"hover-card"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}

	// context-menu resolves transitively through icon (Tier 4 Batch B Task 2
	// — see this file's own TestDeps entry) — same shape as native-select's
	// own entry above; "context-menu" < "icon" lexically, so it sorts first.
	got, err = registry.Resolve([]string{"context-menu"})
	if err != nil {
		t.Fatal(err)
	}
	want = []string{"context-menu", "icon"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}

	// menubar resolves transitively through icon (Tier 4 Batch B Task 3 —
	// see this file's own TestDeps entry) — same shape as context-menu's own
	// entry above; "icon" < "menubar" lexically, so it sorts first.
	got, err = registry.Resolve([]string{"menubar"})
	if err != nil {
		t.Fatal(err)
	}
	want = []string{"icon", "menubar"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}

	// slider has no deps of its own — Resolve returns just itself, same
	// shape as popover/hover-card's own entries above.
	got, err = registry.Resolve([]string{"slider"})
	if err != nil {
		t.Fatal(err)
	}
	want = []string{"slider"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}

	// scroll-area has no deps of its own — Resolve returns just itself,
	// same shape as popover/hover-card/slider's own entries above.
	got, err = registry.Resolve([]string{"scroll-area"})
	if err != nil {
		t.Fatal(err)
	}
	want = []string{"scroll-area"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}

	// select resolves to itself + its derived icon dep, same shape as
	// native-select's own transitive entry above.
	got, err = registry.Resolve([]string{"select"})
	if err != nil {
		t.Fatal(err)
	}
	want = []string{"icon", "select"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}

	// sonner depends on icon (ui.Toast renders Lucide glyphs) — Resolve pulls
	// icon in ahead of it, same shape as select/native-select above.
	got, err = registry.Resolve([]string{"sonner"})
	if err != nil {
		t.Fatal(err)
	}
	want = []string{"icon", "sonner"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}

	// combobox resolves transitively through input-group to button/input/
	// textarea, plus its own direct icon dep — the same transitive-closure
	// shape as alert-dialog/sheet/drawer's own dialog chains above.
	got, err = registry.Resolve([]string{"combobox"})
	if err != nil {
		t.Fatal(err)
	}
	want = []string{"button", "combobox", "icon", "input", "input-group", "textarea"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}

	// sidebar resolves transitively through sheet to dialog/button, plus its
	// own direct button/icon/input/separator/skeleton/tooltip deps — the
	// same transitive-closure shape as alert-dialog/sheet/drawer/combobox's
	// own chains above.
	got, err = registry.Resolve([]string{"sidebar"})
	if err != nil {
		t.Fatal(err)
	}
	want = []string{"button", "dialog", "icon", "input", "separator", "sheet", "sidebar", "skeleton", "tooltip"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}

	// Calendar resolves to itself plus icon/native-select. NativeSelect also
	// depends on Icon, but the flattened result has no duplicate.
	got, err = registry.Resolve([]string{"calendar"})
	if err != nil {
		t.Fatal(err)
	}
	want = []string{"calendar", "icon", "native-select"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestResolveUnknown(t *testing.T) {
	if _, err := registry.Resolve([]string{"nope"}); err == nil {
		t.Fatal("want error for unknown component")
	}
}

func TestResolveRejectsNonComponentFile(t *testing.T) {
	// "core" is not a real path under the flat ui/ anymore, but must still
	// be rejected the same as any other unknown name.
	_, err := registry.Resolve([]string{"core"})
	if err == nil {
		t.Fatal("want error for core, which is not a component")
	}
	if !strings.Contains(err.Error(), `unknown component "core"`) {
		t.Fatalf("got %v, want unknown-component error", err)
	}
}

func TestDepsRejectsNonComponentFile(t *testing.T) {
	if _, err := registry.Deps("core"); err == nil {
		t.Fatal("want error for core, which is not a component")
	}
}
