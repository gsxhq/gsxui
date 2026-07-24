package registry_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/gsxhq/gsxui/internal/registry"
)

func TestComponents(t *testing.T) {
	got, err := registry.Components()
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"accordion", "alert", "alert-dialog", "aspect-ratio", "avatar", "badge", "breadcrumb", "button", "button-group", "card", "carousel", "checkbox", "collapsible", "command", "context-menu", "dialog", "drawer", "dropdown", "empty", "field", "hover-card", "icon", "input", "input-group", "input-otp", "item", "kbd", "label", "pagination", "popover", "progress", "radio", "scroll-area", "select", "separator", "sheet", "skeleton", "slider", "spinner", "switch", "table", "tabs", "textarea", "toggle", "toggle-group", "tooltip"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
	for _, unwanted := range []string{"core", "gsxui", "index", "selectbox", "switchctl"} {
		if slicesContains(got, unwanted) {
			t.Fatalf("Components() = %v, should not contain %q", got, unwanted)
		}
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

	// select.gsx imports ui/icon (the chevron).
	deps, err = registry.Deps("select")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(deps, []string{"icon"}) {
		t.Fatalf("select deps = %v, want [icon]", deps)
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

	deps, err = registry.Deps("button")
	if err != nil {
		t.Fatal(err)
	}
	if len(deps) != 0 {
		t.Fatalf("button deps = %v, want none", deps)
	}

	// pagination.gsx imports ui/icon (ChevronLeft/ChevronRight/Ellipsis) and
	// PaginationLink calls button.gsx's package-private base/variantClass/
	// sizeClass helpers directly (flat package, no import needed for that
	// edge — resolved via declIndex, same shape as dialog's button dep).
	deps, err = registry.Deps("pagination")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(deps, []string{"button", "icon"}) {
		t.Fatalf("pagination deps = %v, want [button icon]", deps)
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

	// toggle-group.gsx has no icon import; ToggleGroupItem calls toggle.gsx's
	// package-private toggleBase/toggleVariantClass/toggleSizeClass directly
	// (flat package intra-package edge, same declIndex-resolved shape as
	// pagination's own button dep above).
	deps, err = registry.Deps("toggle-group")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(deps, []string{"toggle"}) {
		t.Fatalf("toggle-group deps = %v, want [toggle]", deps)
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

	// context-menu.gsx has no icon import and no intra-package reference to
	// another component — ContextMenu/ContextMenuTrigger/ContextMenuContent/
	// ContextMenuItem/ContextMenuLabel/ContextMenuSeparator/
	// ContextMenuShortcut are all plain elements, same shape as popover's own
	// deps entry (the site example composes nothing from another ui.*
	// component either).
	deps, err = registry.Deps("context-menu")
	if err != nil {
		t.Fatal(err)
	}
	if len(deps) != 0 {
		t.Fatalf("context-menu deps = %v, want none", deps)
	}

	// slider.gsx has no icon import and no intra-package reference to
	// another component (the site example composes nothing from another
	// ui.* component either) — same shape as toggle's/popover's/hover-
	// card's/context-menu's own deps entries above.
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
	// shape as toggle's/popover's/hover-card's/context-menu's/slider's own
	// deps entries above.
	deps, err = registry.Deps("scroll-area")
	if err != nil {
		t.Fatal(err)
	}
	if len(deps) != 0 {
		t.Fatalf("scroll-area deps = %v, want none", deps)
	}

	// input-otp.gsx imports ui/icon (InputOTPSeparator's icon.Minus) — same
	// shape as accordion's/select's/spinner's/breadcrumb's own deps entries
	// above; no intra-package reference to another component (it does NOT
	// compose ui.Input, see ui/input-otp.gsx's own doc comment).
	deps, err = registry.Deps("input-otp")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(deps, []string{"icon"}) {
		t.Fatalf("input-otp deps = %v, want [icon]", deps)
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
	// constraint, same selectbox/switchctl precedent).
	if !registry.HasJS("hover-card") {
		t.Error("hover-card should have JS")
	}
	// context-menu has its own ui/context-menu.js (cursor-positioned open on
	// contextmenu, adapted from dropdown.js's menu semantics).
	if !registry.HasJS("context-menu") {
		t.Error("context-menu should have JS")
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
	// input-otp has its own ui/input-otp.js (the entire hidden-single-input
	// mechanism: DOM-order data-index stamping, char/data-active/fake-caret
	// recompute on input/selectionchange/focus/blur, per-character pattern
	// filtering, slot-click-to-position).
	if !registry.HasJS("input-otp") {
		t.Error("input-otp should have JS")
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

	got, err = registry.Resolve([]string{"select"})
	if err != nil {
		t.Fatal(err)
	}
	want = []string{"icon", "select"}
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

	// context-menu has no deps of its own — Resolve returns just itself,
	// same shape as popover/hover-card's own entries above.
	got, err = registry.Resolve([]string{"context-menu"})
	if err != nil {
		t.Fatal(err)
	}
	want = []string{"context-menu"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}

	// slider has no deps of its own — Resolve returns just itself, same
	// shape as popover/hover-card/context-menu's own entries above.
	got, err = registry.Resolve([]string{"slider"})
	if err != nil {
		t.Fatal(err)
	}
	want = []string{"slider"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}

	// scroll-area has no deps of its own — Resolve returns just itself,
	// same shape as popover/hover-card/context-menu/slider's own entries
	// above.
	got, err = registry.Resolve([]string{"scroll-area"})
	if err != nil {
		t.Fatal(err)
	}
	want = []string{"scroll-area"}
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
