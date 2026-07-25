package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestMenubarTriggerRovingAndAria(t *testing.T) {
	got := render(t, ui.MenubarTrigger(gsx.Raw("File"), nil))
	for _, want := range []string{
		`aria-haspopup="menu"`,
		`aria-expanded="false"`,
		`data-gsxui-menubar-trigger`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
	// Roving tabindex is JS-normalized-at-init (toggle-group.js model) — the
	// server renders NO tabindex attribute at all, a graceful no-JS fallback
	// where every trigger is its own tab stop until menubar.js collapses the
	// bar to one.
	if strings.Contains(got, "tabindex") {
		t.Errorf("MenubarTrigger should not server-render tabindex\nin: %s", got)
	}
}

func TestMenubarTriggerDisabledRendersNativeAttr(t *testing.T) {
	// A fix-round finding: menubar.js's roving-tabindex normalize() and
	// arrow-key walk filter to ENABLED triggers only (enabledTriggersOf) —
	// without it, a bar whose FIRST trigger is disabled would hand it
	// tabIndex=0 in normalize(), an unfocusable disabled button, leaving the
	// whole bar with no reachable tab stop at all. That JS-side `.disabled`
	// read depends on MenubarTrigger being a real <button> that renders the
	// native boolean attribute here — same "prove the render contract the
	// JS depends on, at the layer Go can test" shape as toggle-group_test.go's
	// own TestToggleGroupDisabledCascade.
	got := render(t, ui.MenubarTrigger(gsx.Raw("File"), gsx.Attrs{{Key: "disabled", Value: true}}))
	if !strings.Contains(got, "disabled") {
		t.Errorf("want disabled attribute\nin: %s", got)
	}
}

func TestMenubarRootIsAMenubar(t *testing.T) {
	got := render(t, ui.Menubar(gsx.Raw("x"), nil))
	for _, want := range []string{`role="menubar"`, `data-slot="menubar"`, `data-gsxui-menubar`} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
}

func TestMenubarCheckboxItemCheckedServerRendered(t *testing.T) {
	on := render(t, ui.MenubarCheckboxItem(true, "wrap", gsx.Raw("Word Wrap"), nil))
	if !strings.Contains(on, `aria-checked="true"`) || !strings.Contains(on, `role="menuitemcheckbox"`) {
		t.Errorf("want checked menuitemcheckbox\nin: %s", on)
	}
	off := render(t, ui.MenubarCheckboxItem(false, "wrap", gsx.Raw("Word Wrap"), nil))
	if !strings.Contains(off, `aria-checked="false"`) {
		t.Errorf("want aria-checked=false\nin: %s", off)
	}
}

func TestMenubarStructure(t *testing.T) {
	got := render(t, ui.Menubar(
		gsx.Fragment(
			ui.MenubarMenu(gsx.Fragment(
				ui.MenubarTrigger(gsx.Raw("File"), nil),
				ui.MenubarContent(gsx.Fragment(
					ui.MenubarItem("", gsx.Fragment(
						gsx.Raw("New Tab"),
						ui.MenubarShortcut(gsx.Raw("⌘T"), nil),
					), nil),
					ui.MenubarSeparator(nil),
					ui.MenubarItem("", gsx.Raw("Print..."), nil),
				), nil),
			), nil),
		),
		nil,
	))
	for _, want := range []string{
		`data-slot="menubar"`, `role="menubar"`, `data-gsxui-menubar`,
		`data-slot="menubar-menu"`, `data-gsxui-menubar-menu`, `class="contents"`,
		`data-slot="menubar-trigger"`, `data-gsxui-menubar-trigger`,
		`aria-haspopup="menu"`, `aria-expanded="false"`,
		`data-slot="menubar-content"`, `data-gsxui-menubar-content`,
		`popover="auto"`, `role="menu"`, `data-state="closed"`, `data-side="bottom"`,
		`data-slot="menubar-item"`, `data-gsxui-menubar-item`,
		`role="menuitem"`, `tabindex="-1"`, `data-variant="default"`,
		">New Tab<", `data-slot="menubar-shortcut"`, ">⌘T<",
		`data-slot="menubar-separator"`, `role="separator"`,
		">Print...<",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestMenubarItemVariants(t *testing.T) {
	cases := map[string]string{
		"":            "focus:bg-accent",
		"destructive": "data-[variant=destructive]:text-destructive",
	}
	for variant, wantClass := range cases {
		got := render(t, ui.MenubarItem(variant, gsx.Raw("x"), nil))
		if !strings.Contains(got, wantClass) {
			t.Errorf("variant %q: missing %q\nin: %s", variant, wantClass, got)
		}
	}
	got := render(t, ui.MenubarItem("", gsx.Raw("x"), nil))
	if !strings.Contains(got, `data-variant="default"`) {
		t.Errorf("zero-value variant should stamp data-variant=\"default\"\nin: %s", got)
	}
}

func TestMenubarContentCallerClassMerges(t *testing.T) {
	got := render(t, ui.MenubarContent(gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "z-10"}}))
	if strings.Contains(got, "z-50") {
		t.Errorf("base z-50 should be dropped by caller z-10\nin: %s", got)
	}
	for _, want := range []string{"z-10", "rounded-lg", "bg-popover"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestMenubarAttrsFallThrough(t *testing.T) {
	got := render(t, ui.MenubarContent(gsx.Raw("x"), gsx.Attrs{
		{Key: "id", Value: "my-menu"},
		{Key: "aria-label", Value: "File"},
	}))
	for _, want := range []string{`id="my-menu"`, `aria-label="File"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestMenubarPopoverAttrOverridable(t *testing.T) {
	got := render(t, ui.MenubarContent(gsx.Raw("x"), gsx.Attrs{{Key: "popover", Value: "manual"}}))
	if strings.Contains(got, `popover="auto"`) {
		t.Errorf("caller popover=manual should replace the default auto\nin: %s", got)
	}
	if !strings.Contains(got, `popover="manual"`) {
		t.Errorf("missing overridden popover=\"manual\"\nin: %s", got)
	}
}

func TestMenubarRadioItemCheckedServerRendered(t *testing.T) {
	got := render(t, ui.MenubarRadioItem(true, "top", gsx.Raw("Top"), nil))
	for _, want := range []string{`role="menuitemradio"`, `aria-checked="true"`, `data-value="top"`} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
}

func TestMenubarSubTriggerAria(t *testing.T) {
	got := render(t, ui.MenubarSubTrigger(gsx.Raw("More Tools"), nil))
	for _, want := range []string{
		`aria-haspopup="menu"`,
		`aria-expanded="false"`,
		`data-gsxui-menubar-sub-trigger`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
}

func TestMenubarSubContentIsANestedPopover(t *testing.T) {
	// Spec §1 (dropdown.gsx's own file header, reused verbatim): a
	// non-nested popover opened programmatically light-dismisses its parent,
	// so submenu content MUST be a popover nested in the parent.
	got := render(t, ui.MenubarSubContent(gsx.Raw("x"), nil))
	for _, want := range []string{`popover="auto"`, `role="menu"`, `data-gsxui-menubar-sub-content`} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
}

func TestMenubarSubNestsContentInsideParentContentPinned(t *testing.T) {
	// Exact full-render pin, house style (dropdown_test.go's own
	// TestDropdownMenuSubNestsContentInsideParentContentPinned): a
	// string-index-order check alone is a tautology here — the caller builds
	// this tree and gsx has no portal, so no realistic regression could ever
	// put "INNER" before the parent's own data-slot. Pinning the full byte
	// sequence proves the nesting AND the class strings of every part in one
	// assertion.
	got := render(t, ui.MenubarContent(
		ui.MenubarSub(
			gsx.Fragment(
				ui.MenubarSubTrigger(gsx.Raw("More Tools"), nil),
				ui.MenubarSubContent(gsx.Raw("INNER"), nil),
			), nil), nil))
	want := `<div data-slot="menubar-content" data-gsxui-menubar-content popover="auto" role="menu" tabindex="-1" data-state="closed" data-side="bottom" class="z-50 min-w-36 origin-top-left overflow-hidden rounded-lg border bg-popover p-1 text-popover-foreground shadow-md opacity-0 scale-95 transition-[opacity,scale,translate,display,overlay] transition-discrete duration-150 open:opacity-100 open:scale-100 starting:open:opacity-0 starting:open:scale-95 data-[side=bottom]:starting:open:-translate-y-2 data-[side=left]:starting:open:translate-x-2 data-[side=right]:starting:open:-translate-x-2 data-[side=top]:starting:open:translate-y-2"><div data-slot="menubar-sub" data-gsxui-menubar-sub class="contents"><div data-slot="menubar-sub-trigger" data-gsxui-menubar-sub-trigger role="menuitem" aria-haspopup="menu" aria-expanded="false" data-state="closed" tabindex="-1" class="flex cursor-default items-center gap-1.5 rounded-md px-1.5 py-1 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground data-[state=open]:bg-accent data-[state=open]:text-accent-foreground [&amp;_svg:not([class*=&#39;size-&#39;])]:size-4">More Tools<svg data-slot="icon" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="ml-auto size-4"><path d="m9 18 6-6-6-6"/></svg></div><div data-slot="menubar-sub-content" data-gsxui-menubar-sub-content popover="auto" role="menu" tabindex="-1" data-state="closed" data-side="right" class="z-50 min-w-[8rem] origin-top-left overflow-hidden rounded-lg border bg-popover p-1 text-popover-foreground shadow-lg opacity-0 scale-95 transition-[opacity,scale,translate,display,overlay] transition-discrete duration-150 open:opacity-100 open:scale-100 starting:open:opacity-0 starting:open:scale-95 data-[side=bottom]:starting:open:-translate-y-2 data-[side=left]:starting:open:translate-x-2 data-[side=right]:starting:open:-translate-x-2 data-[side=top]:starting:open:translate-y-2">INNER</div></div></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestMenubarPinned(t *testing.T) {
	// Exact full-render pin: nova metrics — byte-identical to
	// DropdownMenuItem's own pinned class (dropdown_test.go's
	// TestDropdownPinned), per the source map's own finding that Item is
	// identical across all three menu families modulo prefix, and this
	// port's own decision to apply the shared nova retarget uniformly.
	got := render(t, ui.MenubarItem("", gsx.Raw("Print..."), nil))
	want := `<div data-slot="menubar-item" data-gsxui-menubar-item data-variant="default" role="menuitem" tabindex="-1" class="relative flex cursor-default items-center gap-1.5 rounded-md px-1.5 py-1 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50 data-[variant=destructive]:text-destructive data-[variant=destructive]:focus:bg-destructive/10 data-[variant=destructive]:focus:text-destructive dark:data-[variant=destructive]:focus:bg-destructive/20 [&amp;_svg]:pointer-events-none [&amp;_svg]:shrink-0 [&amp;_svg:not([class*=&#39;size-&#39;])]:size-4 [&amp;_svg:not([class*=&#39;text-&#39;])]:text-muted-foreground data-[variant=destructive]:*:[svg]:text-destructive!">Print...</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestMenubarTriggerPinned(t *testing.T) {
	// Exact full-render pin: nova metrics (rounded-sm px-1.5 py-[2px], nova's
	// own .cn-menubar-trigger value — NOT bumped to rounded-md, unlike the
	// item-shaped parts) — data-[state=open]: kept per the standing house
	// exception (not nova's hover:/aria-expanded: rewrite, a structural
	// interaction-model swap out of a metric-only retarget's scope).
	got := render(t, ui.MenubarTrigger(gsx.Raw("File"), nil))
	want := `<button data-slot="menubar-trigger" data-gsxui-menubar-trigger type="button" aria-haspopup="menu" aria-expanded="false" data-state="closed" class="flex items-center rounded-sm px-1.5 py-[2px] text-sm font-medium outline-hidden select-none focus:bg-accent focus:text-accent-foreground data-[state=open]:bg-accent data-[state=open]:text-accent-foreground">File</button>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestMenubarContentPinned(t *testing.T) {
	// Exact full-render pin: overflow-hidden only (not dropdown's own
	// overflow-x-hidden/overflow-y-auto — MenubarContent isn't a scrollable
	// container, see the source map's own §5 note), min-w-36 (nova's genuine
	// delta from new-york-v4's min-w-[12rem]), rounded-lg (nova), border kept
	// (house exception, not nova's ring-1), data-side="bottom" server-
	// rendered statically (menubar.js always anchors below the trigger, same
	// stopgap as dropdown's own).
	got := render(t, ui.MenubarContent(gsx.Raw("x"), nil))
	want := `<div data-slot="menubar-content" data-gsxui-menubar-content popover="auto" role="menu" tabindex="-1" data-state="closed" data-side="bottom" class="z-50 min-w-36 origin-top-left overflow-hidden rounded-lg border bg-popover p-1 text-popover-foreground shadow-md opacity-0 scale-95 transition-[opacity,scale,translate,display,overlay] transition-discrete duration-150 open:opacity-100 open:scale-100 starting:open:opacity-0 starting:open:scale-95 data-[side=bottom]:starting:open:-translate-y-2 data-[side=left]:starting:open:translate-x-2 data-[side=right]:starting:open:-translate-x-2 data-[side=top]:starting:open:translate-y-2">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestMenubarCheckboxItemPinned(t *testing.T) {
	// Exact full-render pin, checked=true: byte-identical to
	// DropdownMenuCheckboxItem's own pinned class — new-york-v4's one real
	// delta (rounded-xs vs dropdown/context's rounded-sm) is erased by nova,
	// which unifies all four onto rounded-md (source map's own finding).
	// This is a DELIBERATE OVERRIDE of nova's own menubar-specific rule, not
	// agreement with it — nova's own .cn-menubar-checkbox-item actually puts
	// the indicator on the LEFT (see the component's own doc comment).
	got := render(t, ui.MenubarCheckboxItem(true, "show-toolbar", gsx.Raw("Toolbar"), nil))
	want := `<div data-slot="menubar-checkbox-item" data-gsxui-menubar-checkbox-item role="menuitemcheckbox" data-value="show-toolbar" aria-checked="true" data-state="checked" tabindex="-1" class="relative flex cursor-default items-center gap-1.5 rounded-md py-1 pr-8 pl-1.5 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50 [&amp;_svg]:pointer-events-none [&amp;_svg]:shrink-0 [&amp;_svg:not([class*=&#39;size-&#39;])]:size-4"><span class="pointer-events-none absolute right-2 hidden size-4 items-center justify-center [[data-state=checked]_&amp;]:flex"><svg data-slot="icon" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="size-4"><path d="M20 6 9 17l-5-5"/></svg></span>Toolbar</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestMenubarRadioItemPinned(t *testing.T) {
	// Exact full-render pin, checked=true — same byte-identical-to-dropdown
	// class as CheckboxItem's own pinned test above.
	got := render(t, ui.MenubarRadioItem(true, "top", gsx.Raw("Top"), nil))
	want := `<div data-slot="menubar-radio-item" data-gsxui-menubar-radio-item role="menuitemradio" data-value="top" aria-checked="true" data-state="checked" tabindex="-1" class="relative flex cursor-default items-center gap-1.5 rounded-md py-1 pr-8 pl-1.5 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50 [&amp;_svg]:pointer-events-none [&amp;_svg]:shrink-0 [&amp;_svg:not([class*=&#39;size-&#39;])]:size-4"><span class="pointer-events-none absolute right-2 hidden size-4 items-center justify-center [[data-state=checked]_&amp;]:flex"><svg data-slot="icon" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="size-2 fill-current"><circle cx="12" cy="12" r="10"/></svg></span>Top</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestMenubarSubTriggerPinned(t *testing.T) {
	// Exact full-render pin: nova adds gap-1.5 (missing from new-york-v4's
	// own MenubarSubTrigger entirely, same "nova fixes an asymmetry"
	// harmonization already ledgered for context-menu's own SubTrigger),
	// rounded-sm->rounded-md/px-2->px-1.5/py-1.5->py-1 (nova), outline-hidden
	// (NOT new-york-v4's own outline-none — a deliberate ADAPT, see the
	// component's own doc comment), data-[state=open]: kept (house
	// exception), trailing icon ported as size-4 (not new-york-v4's literal
	// h-4 w-4 — a no-op spelling divergence, same call as context-menu's own
	// SubTrigger icon). NOT byte-identical to DropdownMenuSubTrigger's/
	// ContextMenuSubTrigger's own pinned class, unlike the CheckboxItem/
	// RadioItem tests above: this class also lacks
	// [&_svg]:pointer-events-none/[&_svg]:shrink-0/
	// [&_svg:not([class*='text-'])]:text-muted-foreground, so its chevron
	// renders in the default foreground color, not muted — see the
	// component's own doc comment for why (menubar.tsx's own source never
	// had those three tokens, and nova is silent on them for every
	// SubTrigger in this menu family).
	got := render(t, ui.MenubarSubTrigger(gsx.Raw("More Tools"), nil))
	want := `<div data-slot="menubar-sub-trigger" data-gsxui-menubar-sub-trigger role="menuitem" aria-haspopup="menu" aria-expanded="false" data-state="closed" tabindex="-1" class="flex cursor-default items-center gap-1.5 rounded-md px-1.5 py-1 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground data-[state=open]:bg-accent data-[state=open]:text-accent-foreground [&amp;_svg:not([class*=&#39;size-&#39;])]:size-4">More Tools<svg data-slot="icon" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="ml-auto size-4"><path d="m9 18 6-6-6-6"/></svg></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestMenubarSubContentPinned(t *testing.T) {
	// Exact full-render pin: min-w-[8rem] (nova's own min-w-32 == 8rem, no
	// genuine delta from new-york-v4 here — unlike dropdown's own SubContent,
	// which nova DOES narrow), rounded-lg (nova), shadow-lg (matches both),
	// border kept (house exception), origin-top-left/discrete-transition
	// ADAPT (same mechanism as every other popover-based Content in this
	// codebase — see the component's own doc comment for the animate-out
	// ruling this makes moot).
	got := render(t, ui.MenubarSubContent(gsx.Raw("x"), nil))
	want := `<div data-slot="menubar-sub-content" data-gsxui-menubar-sub-content popover="auto" role="menu" tabindex="-1" data-state="closed" data-side="right" class="z-50 min-w-[8rem] origin-top-left overflow-hidden rounded-lg border bg-popover p-1 text-popover-foreground shadow-lg opacity-0 scale-95 transition-[opacity,scale,translate,display,overlay] transition-discrete duration-150 open:opacity-100 open:scale-100 starting:open:opacity-0 starting:open:scale-95 data-[side=bottom]:starting:open:-translate-y-2 data-[side=left]:starting:open:translate-x-2 data-[side=right]:starting:open:-translate-x-2 data-[side=top]:starting:open:translate-y-2">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestMenubarLabelPinned(t *testing.T) {
	// Exact full-render pin: nova's own .cn-menubar-label genuinely differs
	// from DropdownMenuLabel's already-shipped text-xs — MenubarLabel is
	// text-sm, a real per-component nova value, not a copy error.
	got := render(t, ui.MenubarLabel(gsx.Raw("People"), nil))
	want := `<div data-slot="menubar-label" class="px-1.5 py-1 text-sm font-medium">People</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestMenubarSeparatorPinned(t *testing.T) {
	got := render(t, ui.MenubarSeparator(nil))
	want := `<div data-slot="menubar-separator" role="separator" class="-mx-1 my-1 h-px bg-border"></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestMenubarShortcutPinned(t *testing.T) {
	got := render(t, ui.MenubarShortcut(gsx.Raw("⌘T"), nil))
	want := `<span data-slot="menubar-shortcut" class="ml-auto text-xs tracking-widest text-muted-foreground">⌘T</span>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestMenubarRootPinned(t *testing.T) {
	// Exact full-render pin for the visible bar itself: nova's h-8/gap-0.5/
	// rounded-lg/p-[3px] metrics over new-york-v4's h-9/gap-1/rounded-md/p-1;
	// bg-background is KEPT (a color, out of scope for a metric-only
	// retarget — see the component's own doc comment for why this is a
	// narrower ruling than navigation-menu-trigger's own bg-background drop);
	// shadow-xs is DROPPED (nova's own .cn-menubar rule omits it entirely,
	// same SHADOW-PRESENCE category as the other 12 already-shipped shadow
	// removals in the nova density retarget, see docs/jsx-parity.md ## nova
	// density).
	got := render(t, ui.Menubar(gsx.Raw("x"), nil))
	want := `<div data-slot="menubar" data-gsxui-menubar role="menubar" class="flex h-8 items-center gap-0.5 rounded-lg border bg-background p-[3px]">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestMenubarMenuPinned(t *testing.T) {
	got := render(t, ui.MenubarMenu(gsx.Raw("x"), nil))
	want := `<div data-slot="menubar-menu" data-gsxui-menubar-menu class="contents">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestMenubarGroupPinned(t *testing.T) {
	// No prior test rendered MenubarGroup at all — a11y-only wrapper,
	// role="group" added here (derived-not-read, same call as
	// DropdownMenuGroup/ContextMenuGroup), no data-gsxui-* hook since
	// nothing in menubar.js binds to it.
	got := render(t, ui.MenubarGroup(gsx.Raw("x"), nil))
	want := `<div data-slot="menubar-group" role="group">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestMenubarRadioGroupPinned(t *testing.T) {
	// No prior test rendered MenubarRadioGroup or asserted
	// data-gsxui-menubar-radio-group at all, despite it being load-bearing:
	// menubar.js's radio-item click handler scopes its "clear every OTHER
	// item" sibling walk to `closest("[data-gsxui-menubar-radio-group]")`
	// (menubar.js's own radio-item click handler) — without this hook
	// rendering, two independent MenubarRadioGroups on the same page would
	// cross-clear each other's selections.
	got := render(t, ui.MenubarRadioGroup("benoit", gsx.Raw("x"), nil))
	want := `<div data-slot="menubar-radio-group" data-gsxui-menubar-radio-group role="group" data-value="benoit">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}
