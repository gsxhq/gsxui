package dropdown_test

import (
	"context"
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui/dropdown"
)

func render(t *testing.T, n gsx.Node) string {
	t.Helper()
	var sb strings.Builder
	if err := n.Render(context.Background(), &sb); err != nil {
		t.Fatal(err)
	}
	return sb.String()
}

func TestDropdownStructure(t *testing.T) {
	got := render(t, dropdown.DropdownMenu(
		gsx.Fragment(
			dropdown.DropdownMenuTrigger(gsx.Raw("Open"), nil),
			dropdown.DropdownMenuContent(gsx.Fragment(
				dropdown.DropdownMenuLabel(gsx.Raw("Actions"), nil),
				dropdown.DropdownMenuSeparator(nil),
				dropdown.DropdownMenuItem("", gsx.Fragment(
					gsx.Raw("Edit"),
					dropdown.DropdownMenuShortcut(gsx.Raw("⌘E"), nil),
				), nil),
				dropdown.DropdownMenuItem("destructive", gsx.Raw("Delete"), nil),
			), nil),
		),
		nil,
	))
	for _, want := range []string{
		`data-gsxui-dropdown`,         // root hook
		`class="contents"`,            // root is layout-neutral
		`data-gsxui-dropdown-trigger`, // trigger hook
		`aria-haspopup="menu"`,        // trigger a11y: server-rendered
		`aria-expanded="false"`,       // trigger a11y: initial state
		`data-slot="dropdown-menu-content"`,
		`data-gsxui-dropdown-content`,
		`popover="auto"`,      // top-layer, light-dismiss, free Esc
		`role="menu"`,         // content a11y
		`data-state="closed"`, // server-rendered initial state
		`data-slot="dropdown-menu-label"`, ">Actions<",
		`data-slot="dropdown-menu-separator"`, `role="separator"`,
		`data-slot="dropdown-menu-item"`, `data-gsxui-dropdown-item`,
		`role="menuitem"`, `tabindex="-1"`,
		`data-variant="default"`, ">Edit<",
		`data-slot="dropdown-menu-shortcut"`, ">⌘E<",
		`data-variant="destructive"`, ">Delete<",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestDropdownItemVariants(t *testing.T) {
	cases := map[string]string{
		"":            "focus:bg-accent",
		"destructive": "data-[variant=destructive]:text-destructive",
	}
	for variant, wantClass := range cases {
		got := render(t, dropdown.DropdownMenuItem(variant, gsx.Raw("x"), nil))
		if !strings.Contains(got, wantClass) {
			t.Errorf("variant %q: missing %q\nin: %s", variant, wantClass, got)
		}
	}
	// Zero value renders the shadcn default stamp.
	got := render(t, dropdown.DropdownMenuItem("", gsx.Raw("x"), nil))
	if !strings.Contains(got, `data-variant="default"`) {
		t.Errorf("zero-value variant should stamp data-variant=\"default\"\nin: %s", got)
	}
}

func TestDropdownContentCallerClassMerges(t *testing.T) {
	// Caller z-10 must WIN over base z-50 via tailwind-merge — and base
	// structural classes must survive.
	got := render(t, dropdown.DropdownMenuContent(gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "z-10"}}))
	if strings.Contains(got, "z-50") {
		t.Errorf("base z-50 should be dropped by caller z-10\nin: %s", got)
	}
	for _, want := range []string{"z-10", "rounded-md", "bg-popover"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestDropdownAttrsFallThrough(t *testing.T) {
	got := render(t, dropdown.DropdownMenuContent(gsx.Raw("x"), gsx.Attrs{
		{Key: "id", Value: "my-menu"},
		{Key: "aria-label", Value: "Actions"},
	}))
	for _, want := range []string{`id="my-menu"`, `aria-label="Actions"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestDropdownPopoverAttrOverridable(t *testing.T) {
	// popover is a regular attribute with a value — attrs fallthrough can
	// override it like any other attribute (e.g. a caller opting a specific
	// menu out of the default "auto" light-dismiss behavior).
	got := render(t, dropdown.DropdownMenuContent(gsx.Raw("x"), gsx.Attrs{{Key: "popover", Value: "manual"}}))
	if strings.Contains(got, `popover="auto"`) {
		t.Errorf("caller popover=manual should replace the default auto\nin: %s", got)
	}
	if !strings.Contains(got, `popover="manual"`) {
		t.Errorf("missing overridden popover=\"manual\"\nin: %s", got)
	}
}

func TestDropdownPinned(t *testing.T) {
	// Exact full-render pin for DropdownMenuItem's default variant, verified
	// token-by-token against shadcn's DropdownMenuItem
	// (registry/new-york-v4/ui/dropdown-menu.tsx) and docs/jsx-parity.md's
	// ADAPT: the inset prop is dropped, so data-[inset]:pl-8 is dropped
	// along with it.
	got := render(t, dropdown.DropdownMenuItem("", gsx.Raw("Edit"), nil))
	want := `<div data-slot="dropdown-menu-item" data-gsxui-dropdown-item data-variant="default" role="menuitem" tabindex="-1" class="relative flex cursor-default items-center gap-2 rounded-sm px-2 py-1.5 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50 data-[variant=destructive]:text-destructive data-[variant=destructive]:focus:bg-destructive/10 data-[variant=destructive]:focus:text-destructive dark:data-[variant=destructive]:focus:bg-destructive/20 [&amp;_svg]:pointer-events-none [&amp;_svg]:shrink-0 [&amp;_svg:not([class*=&#39;size-&#39;])]:size-4 [&amp;_svg:not([class*=&#39;text-&#39;])]:text-muted-foreground data-[variant=destructive]:*:[svg]:text-destructive!">Edit</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestDropdownContentPinned(t *testing.T) {
	// Exact full-render pin for DropdownMenuContent, verified token-by-token
	// against shadcn's DropdownMenuContent classes plus the popover/role/
	// data-state hooks that replace Radix's Portal+Content wiring.
	got := render(t, dropdown.DropdownMenuContent(gsx.Raw("x"), nil))
	want := `<div data-slot="dropdown-menu-content" data-gsxui-dropdown-content popover="auto" role="menu" data-state="closed" class="z-50 max-h-(--radix-dropdown-menu-content-available-height) min-w-[8rem] origin-(--radix-dropdown-menu-content-transform-origin) overflow-x-hidden overflow-y-auto rounded-md border bg-popover p-1 text-popover-foreground shadow-md data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95 data-[state=open]:animate-in data-[state=open]:fade-in-0 data-[state=open]:zoom-in-95">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}
