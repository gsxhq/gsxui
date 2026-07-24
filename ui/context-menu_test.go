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
	for _, want := range []string{"z-10", "rounded-md", "bg-popover"} {
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
	want := `<div data-slot="context-menu-item" data-gsxui-contextmenu-item data-variant="default" role="menuitem" tabindex="-1" class="relative flex cursor-default items-center gap-2 rounded-sm px-2 py-1.5 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50 data-[variant=destructive]:text-destructive data-[variant=destructive]:focus:bg-destructive/10 data-[variant=destructive]:focus:text-destructive dark:data-[variant=destructive]:focus:bg-destructive/20 [&amp;_svg]:pointer-events-none [&amp;_svg]:shrink-0 [&amp;_svg:not([class*=&#39;size-&#39;])]:size-4 [&amp;_svg:not([class*=&#39;text-&#39;])]:text-muted-foreground data-[variant=destructive]:*:[svg]:text-destructive!">Back</div>`
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
	want := `<div data-slot="context-menu-content" data-gsxui-contextmenu-content popover="auto" role="menu" tabindex="-1" data-state="closed" class="z-50 max-h-96 min-w-[8rem] origin-top-left overflow-x-hidden overflow-y-auto rounded-md border bg-popover p-1 text-popover-foreground shadow-md data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95 data-[state=open]:animate-in data-[state=open]:fade-in-0 data-[state=open]:zoom-in-95">x</div>`
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
