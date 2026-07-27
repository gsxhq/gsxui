package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestContextMenuStructure(t *testing.T) {
	got := render(t, ui.ContextMenu(
		gsx.Fragment(
			ui.ContextMenuTrigger(gsx.Raw("Right click here"), nil),
			ui.ContextMenuContent(gsx.Fragment(
				ui.ContextMenuLabel(gsx.Raw("Actions"), nil),
				ui.ContextMenuSeparator(nil),
				ui.ContextMenuItem("", gsx.Fragment(
					gsx.Raw("Back"),
					ui.ContextMenuShortcut(gsx.Raw("⌘["), nil),
				), nil),
				ui.ContextMenuItem("destructive", gsx.Raw("Delete"), nil),
			), nil),
		),
		nil,
	))
	for _, want := range []string{
		`data-gsxui-contextmenu`, // root hook
		`class="contents"`,       // root is layout-neutral
		`data-slot="context-menu-trigger"`,
		`data-gsxui-contextmenu-trigger`, // trigger hook
		`data-slot="context-menu-content"`,
		`data-gsxui-contextmenu-content`,
		`popover="auto"`,      // top-layer, light-dismiss, free Esc
		`role="menu"`,         // content a11y
		`data-state="closed"`, // server-rendered initial state
		`data-slot="context-menu-label"`, ">Actions<",
		`data-slot="context-menu-separator"`, `role="separator"`,
		`data-slot="context-menu-item"`, `data-gsxui-contextmenu-item`,
		`role="menuitem"`, `tabindex="-1"`,
		`data-variant="default"`, ">Back<",
		`data-slot="context-menu-shortcut"`, ">⌘[<",
		`data-variant="destructive"`, ">Delete<",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
	// No static data-side stamp — cursor-positioned, no fixed anchor side
	// (unlike dropdown/popover/hover-card's own static stamps).
	if strings.Contains(got, `data-side`) {
		t.Errorf("context-menu content should not stamp data-side\nin: %s", got)
	}
}

func TestContextMenuItemVariants(t *testing.T) {
	cases := map[string]string{
		"":            "focus:bg-accent",
		"destructive": "data-[variant=destructive]:text-destructive",
	}
	for variant, wantClass := range cases {
		got := render(t, ui.ContextMenuItem(variant, gsx.Raw("x"), nil))
		if !strings.Contains(got, wantClass) {
			t.Errorf("variant %q: missing %q\nin: %s", variant, wantClass, got)
		}
	}
	// Zero value renders the shadcn default stamp.
	got := render(t, ui.ContextMenuItem("", gsx.Raw("x"), nil))
	if !strings.Contains(got, `data-variant="default"`) {
		t.Errorf("zero-value variant should stamp data-variant=\"default\"\nin: %s", got)
	}
}

func TestContextMenuContentCallerClassMerges(t *testing.T) {
	// Caller z-10 must WIN over base z-50 via tailwind-merge — and base
	// structural classes must survive.
	got := render(t, ui.ContextMenuContent(gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "z-10"}}))
	if strings.Contains(got, "z-50") {
		t.Errorf("base z-50 should be dropped by caller z-10\nin: %s", got)
	}
	for _, want := range []string{"z-10", "rounded-lg", "bg-popover"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestContextMenuAttrsFallThrough(t *testing.T) {
	got := render(t, ui.ContextMenuContent(gsx.Raw("x"), gsx.Attrs{
		{Key: "id", Value: "my-menu"},
		{Key: "aria-label", Value: "Actions"},
	}))
	for _, want := range []string{`id="my-menu"`, `aria-label="Actions"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestContextMenuPopoverAttrOverridable(t *testing.T) {
	// popover is a regular attribute with a value — attrs fallthrough can
	// override it like any other attribute (e.g. a caller opting a specific
	// menu out of the default "auto" light-dismiss behavior).
	got := render(t, ui.ContextMenuContent(gsx.Raw("x"), gsx.Attrs{{Key: "popover", Value: "manual"}}))
	if strings.Contains(got, `popover="auto"`) {
		t.Errorf("caller popover=manual should replace the default auto\nin: %s", got)
	}
	if !strings.Contains(got, `popover="manual"`) {
		t.Errorf("missing overridden popover=\"manual\"\nin: %s", got)
	}
}

func TestContextMenuPinned(t *testing.T) {
	// Exact full-render pin for ContextMenuItem's default variant, verified
	// token-by-token against shadcn's ContextMenuItem
	// (registry/new-york-v4/ui/context-menu.tsx) and docs/jsx-parity.md's
	// ADAPT: the inset prop is dropped, so data-[inset]:pl-8 is dropped
	// along with it — the resulting class is byte-identical to
	// DropdownMenuItem's own pinned class (see dropdown_test.go's
	// TestDropdownPinned), a coincidence of the two shadcn sources sharing
	// every other token.
	got := render(t, ui.ContextMenuItem("", gsx.Raw("Back"), nil))
	want := `<div data-slot="context-menu-item" data-gsxui-contextmenu-item data-variant="default" role="menuitem" tabindex="-1" class="relative flex cursor-default items-center gap-1.5 rounded-md px-1.5 py-1 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50 data-[variant=destructive]:text-destructive data-[variant=destructive]:focus:bg-destructive/10 data-[variant=destructive]:focus:text-destructive dark:data-[variant=destructive]:focus:bg-destructive/20 [&amp;_svg]:pointer-events-none [&amp;_svg]:shrink-0 [&amp;_svg:not([class*=&#39;size-&#39;])]:size-4 [&amp;_svg:not([class*=&#39;text-&#39;])]:text-muted-foreground data-[variant=destructive]:*:[svg]:text-destructive!">Back</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestContextMenuContentPinned(t *testing.T) {
	// Exact full-render pin for ContextMenuContent, verified token-by-token
	// against shadcn's ContextMenuContent classes plus the popover/role/
	// data-state hooks that replace Radix's Portal+Content wiring. No
	// data-side — cursor-positioned, no fixed anchor side (see the
	// component's own doc comment).
	got := render(t, ui.ContextMenuContent(gsx.Raw("x"), nil))
	want := `<div data-slot="context-menu-content" data-gsxui-contextmenu-content popover="auto" role="menu" tabindex="-1" data-state="closed" class="z-50 max-h-96 min-w-36 origin-top-left overflow-x-hidden overflow-y-auto rounded-lg border bg-popover p-1 text-popover-foreground shadow-md opacity-0 scale-95 transition-[opacity,scale,translate,display,overlay] transition-discrete duration-150 open:opacity-100 open:scale-100 starting:open:opacity-0 starting:open:scale-95 data-[side=bottom]:starting:open:-translate-y-2 data-[side=left]:starting:open:translate-x-2 data-[side=right]:starting:open:-translate-x-2 data-[side=top]:starting:open:translate-y-2">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestContextMenuTriggerPinned(t *testing.T) {
	// Exact full-render pin for ContextMenuTrigger — a plain AREA div, not
	// a button (unlike DropdownMenuTrigger), so it carries no aria-haspopup/
	// aria-expanded of its own.
	got := render(t, ui.ContextMenuTrigger(gsx.Raw("Right click here"), nil))
	want := `<div data-slot="context-menu-trigger" data-gsxui-contextmenu-trigger>Right click here</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

// --- Task 2 (Tier 4 Batch B): shared-items parity with dropdown.gsx.
// Mirrors dropdown_test.go's own equivalent tests, ContextMenu* prefix,
// data-gsxui-contextmenu-* hooks — namespaced to this component, NOT a
// prefix shared with dropdown's own data-gsxui-dropdown-* equivalents (fix
// round 1: gsxui.js's shared delegation registry dispatches to every
// handler whose selector matches, regardless of module, so an identical
// selector registered by both dropdown.js and context-menu.js double-fires
// on one event — see ui/context-menu.gsx's own doc comments on these parts
// and docs/jsx-parity.md's ## dropdown ledger MECHANISM entry).

func TestContextMenuCheckboxItemCheckedServerRendered(t *testing.T) {
	on := render(t, ui.ContextMenuCheckboxItem(true, "show-toolbar", gsx.Raw("Toolbar"), nil))
	for _, want := range []string{
		`role="menuitemcheckbox"`,
		`aria-checked="true"`,
		`data-gsxui-contextmenu-checkbox-item`,
		`data-value="show-toolbar"`,
	} {
		if !strings.Contains(on, want) {
			t.Errorf("want %q\nin: %s", want, on)
		}
	}
	off := render(t, ui.ContextMenuCheckboxItem(false, "show-toolbar", gsx.Raw("Toolbar"), nil))
	if !strings.Contains(off, `aria-checked="false"`) {
		t.Errorf("want aria-checked=false\nin: %s", off)
	}
}

func TestContextMenuRadioItemCheckedServerRendered(t *testing.T) {
	got := render(t, ui.ContextMenuRadioItem(true, "top", gsx.Raw("Top"), nil))
	for _, want := range []string{`role="menuitemradio"`, `aria-checked="true"`, `data-value="top"`} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
}

func TestContextMenuSubTriggerAria(t *testing.T) {
	got := render(t, ui.ContextMenuSubTrigger(gsx.Raw("More"), nil))
	for _, want := range []string{
		`aria-haspopup="menu"`,
		`aria-expanded="false"`,
		`data-gsxui-contextmenu-sub-trigger`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
}

func TestContextMenuSubContentIsANestedPopover(t *testing.T) {
	// Spec §1 (dropdown.gsx's own file header, reused verbatim here): a
	// non-nested popover opened programmatically light-dismisses its
	// parent, so submenu content MUST be a popover nested in the parent.
	got := render(t, ui.ContextMenuSubContent(gsx.Raw("x"), nil))
	for _, want := range []string{`popover="auto"`, `role="menu"`, `data-gsxui-contextmenu-sub-content`} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
}

func TestContextMenuSubNestsContentInsideParentContentPinned(t *testing.T) {
	// Exact full-render pin, house style (TestContextMenuContentPinned /
	// dropdown_test.go's own TestDropdownMenuSubNestsContentInsideParentContentPinned):
	// proves the nesting AND the class strings of every part in one
	// assertion — Content wrapping Sub wrapping SubTrigger+SubContent, with
	// SubContent's own markup appearing literally inside Content's closing
	// </div>. Note ContextMenuContent itself stamps no data-side (see the
	// component's own doc comment) while the nested SubContent still stamps
	// data-side="right" — the submenu trigger IS a fixed anchor, unlike the
	// cursor-positioned top-level content.
	got := render(t, ui.ContextMenuContent(
		ui.ContextMenuSub(
			gsx.Fragment(
				ui.ContextMenuSubTrigger(gsx.Raw("More"), nil),
				ui.ContextMenuSubContent(gsx.Raw("INNER"), nil),
			), nil), nil))
	want := `<div data-slot="context-menu-content" data-gsxui-contextmenu-content popover="auto" role="menu" tabindex="-1" data-state="closed" class="z-50 max-h-96 min-w-36 origin-top-left overflow-x-hidden overflow-y-auto rounded-lg border bg-popover p-1 text-popover-foreground shadow-md opacity-0 scale-95 transition-[opacity,scale,translate,display,overlay] transition-discrete duration-150 open:opacity-100 open:scale-100 starting:open:opacity-0 starting:open:scale-95 data-[side=bottom]:starting:open:-translate-y-2 data-[side=left]:starting:open:translate-x-2 data-[side=right]:starting:open:-translate-x-2 data-[side=top]:starting:open:translate-y-2"><div data-slot="context-menu-sub" data-gsxui-contextmenu-sub class="contents"><div data-slot="context-menu-sub-trigger" data-gsxui-contextmenu-sub-trigger role="menuitem" aria-haspopup="menu" aria-expanded="false" data-state="closed" tabindex="-1" class="flex cursor-default items-center gap-1.5 rounded-md px-1.5 py-1 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground data-[state=open]:bg-accent data-[state=open]:text-accent-foreground [&amp;_svg]:pointer-events-none [&amp;_svg]:shrink-0 [&amp;_svg:not([class*=&#39;size-&#39;])]:size-4 [&amp;_svg:not([class*=&#39;text-&#39;])]:text-muted-foreground">More<svg aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="ml-auto" data-gsxui-slot="icon"><path d="m9 18 6-6-6-6"/></svg></div><div data-slot="context-menu-sub-content" data-gsxui-contextmenu-sub-content popover="auto" role="menu" tabindex="-1" data-state="closed" data-side="right" class="z-50 min-w-[96px] origin-top-left overflow-hidden rounded-lg border bg-popover p-1 text-popover-foreground shadow-lg opacity-0 scale-95 transition-[opacity,scale,translate,display,overlay] transition-discrete duration-150 open:opacity-100 open:scale-100 starting:open:opacity-0 starting:open:scale-95 data-[side=bottom]:starting:open:-translate-y-2 data-[side=left]:starting:open:translate-x-2 data-[side=right]:starting:open:-translate-x-2 data-[side=top]:starting:open:translate-y-2">INNER</div></div></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestContextMenuCheckboxItemPinned(t *testing.T) {
	// Exact full-render pin, checked=true: byte-identical class to
	// DropdownMenuCheckboxItem's own pinned class (source map ## shared-items
	// §1 — CheckboxItem is byte-identical modulo component-name prefix), only
	// data-slot differs.
	got := render(t, ui.ContextMenuCheckboxItem(true, "show-toolbar", gsx.Raw("Toolbar"), nil))
	want := `<div data-slot="context-menu-checkbox-item" data-gsxui-contextmenu-checkbox-item role="menuitemcheckbox" data-value="show-toolbar" aria-checked="true" data-state="checked" tabindex="-1" class="relative flex cursor-default items-center gap-1.5 rounded-md py-1 pr-8 pl-1.5 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50 [&amp;_svg]:pointer-events-none [&amp;_svg]:shrink-0 [&amp;_svg:not([class*=&#39;size-&#39;])]:size-4"><span class="pointer-events-none absolute right-2 hidden size-4 items-center justify-center [[data-state=checked]_&amp;]:flex"><svg aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="size-4" data-gsxui-slot="icon"><path d="M20 6 9 17l-5-5"/></svg></span>Toolbar</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestContextMenuRadioItemPinned(t *testing.T) {
	// Exact full-render pin, checked=true — same byte-identical-modulo-prefix
	// class as DropdownMenuRadioItem's own pinned class.
	got := render(t, ui.ContextMenuRadioItem(true, "top", gsx.Raw("Top"), nil))
	want := `<div data-slot="context-menu-radio-item" data-gsxui-contextmenu-radio-item role="menuitemradio" data-value="top" aria-checked="true" data-state="checked" tabindex="-1" class="relative flex cursor-default items-center gap-1.5 rounded-md py-1 pr-8 pl-1.5 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50 [&amp;_svg]:pointer-events-none [&amp;_svg]:shrink-0 [&amp;_svg:not([class*=&#39;size-&#39;])]:size-4"><span class="pointer-events-none absolute right-2 hidden size-4 items-center justify-center [[data-state=checked]_&amp;]:flex"><svg aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="size-2 fill-current" data-gsxui-slot="icon"><circle cx="12" cy="12" r="10"/></svg></span>Top</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestContextMenuSubTriggerPinned(t *testing.T) {
	// Exact full-render pin. The ONE real class-string delta this task's
	// source map found between the two shadcn sources' SubTrigger: dropdown's
	// gap-2 (context-menu's own source lacks it) is harmonized to nova's
	// gap-1.5 on BOTH (so this class is otherwise identical to
	// DropdownMenuSubTrigger's own pinned class), but the trailing
	// ChevronRightIcon here carries NO explicit size-4 (dropdown's does) —
	// preserved as a no-op divergence per the source map's own recommendation,
	// since [&_svg:not([class*='size-'])]:size-4 already stamps size-4 onto
	// it either way.
	got := render(t, ui.ContextMenuSubTrigger(gsx.Raw("More"), nil))
	want := `<div data-slot="context-menu-sub-trigger" data-gsxui-contextmenu-sub-trigger role="menuitem" aria-haspopup="menu" aria-expanded="false" data-state="closed" tabindex="-1" class="flex cursor-default items-center gap-1.5 rounded-md px-1.5 py-1 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground data-[state=open]:bg-accent data-[state=open]:text-accent-foreground [&amp;_svg]:pointer-events-none [&amp;_svg]:shrink-0 [&amp;_svg:not([class*=&#39;size-&#39;])]:size-4 [&amp;_svg:not([class*=&#39;text-&#39;])]:text-muted-foreground">More<svg aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="ml-auto" data-gsxui-slot="icon"><path d="m9 18 6-6-6-6"/></svg></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestContextMenuSubContentPinned(t *testing.T) {
	// Exact full-render pin: byte-identical (source map ## shared-items §1)
	// to DropdownMenuSubContent's own pinned class modulo data-slot — origin-
	// top-left/discrete-transition ADAPT, nova's min-w-[96px]/rounded-lg,
	// border kept (not nova's ring-1, standing house exception). Unlike the
	// top-level ContextMenuContent, this DOES stamp data-side="right" — the
	// SubTrigger is a fixed anchor, same as dropdown's own SubContent.
	got := render(t, ui.ContextMenuSubContent(gsx.Raw("x"), nil))
	want := `<div data-slot="context-menu-sub-content" data-gsxui-contextmenu-sub-content popover="auto" role="menu" tabindex="-1" data-state="closed" data-side="right" class="z-50 min-w-[96px] origin-top-left overflow-hidden rounded-lg border bg-popover p-1 text-popover-foreground shadow-lg opacity-0 scale-95 transition-[opacity,scale,translate,display,overlay] transition-discrete duration-150 open:opacity-100 open:scale-100 starting:open:opacity-0 starting:open:scale-95 data-[side=bottom]:starting:open:-translate-y-2 data-[side=left]:starting:open:translate-x-2 data-[side=right]:starting:open:-translate-x-2 data-[side=top]:starting:open:translate-y-2">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}
