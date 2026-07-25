# Tier 4 Batch B Implementation Plan — the menu family

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Close shadcn's menu family — the missing item types and submenu
machinery on `dropdown` and `context-menu`, plus `menubar` and
`navigation-menu`. Coverage 50 → 52, with two existing components reaching
full parity and the README's submenu backlog entry retired.

**Architecture:** Server-rendered gsx components carrying shadcn's class
strings token-for-token, behavior in `ui/<name>.js`. Submenus are nested
`popover="auto"` elements inside their parent's content — measured as the
only mechanism that keeps the parent open under programmatic opening.
Checked state is a server-rendered parameter reflected in the DOM; the
component emits and never persists.

**Tech Stack:** Go + gsx templating, Tailwind v4 utility classes, vanilla ES
modules, `go test` render pins.

**Design spec:** `docs/superpowers/specs/2026-07-25-tier4-batch-b-design.md`
— read it before Task 1. Its §1 (the measured submenu mechanism) is binding.

## Global Constraints

Every task's requirements implicitly include this section.

- **Markup reference is `~/personal/shadcn-ui/apps/v4/registry/new-york-v4/ui/<name>.tsx`.**
  Never port from `registry/bases/{radix,base,aria}/` — those emit `.cn-*`
  named classes with the utilities extracted into `registry/styles/style-*.css`.
- **Visual reference is `~/personal/shadcn-ui/apps/v4/registry/styles/style-nova.css`.**
  Where it and `new-york-v4` disagree, nova wins; `bases/base/ui/<name>.tsx`
  is the tiebreak for markup-structure questions only.
- **Standing house exceptions** (do not "fix" toward nova): border→`ring-1`
  swaps NOT adopted; keep `data-[orientation=…]:` not `data-horizontal:`;
  keep `focus-visible:ring-[3px]` not `ring-3`; keep `data-[state=…]` not
  `data-open:`.
- **Submenu content is DOM-nested inside its parent's content, never
  portalled.** See spec §1. A non-nested popover opened programmatically
  light-dismisses its parent; submenus open on hover and ArrowRight, so
  programmatic opening is unavoidable. This is load-bearing.
- **`data-state="open"` is stamped synchronously BEFORE `showPopover()`** —
  the flash-avoidance rule.
- **Author `display` utilities on a popover MUST be gated on `:open`.** A
  bare `block`/`grid` beats the UA's closed-popover `display: none` and
  leaves hit-testable ghost boxes. This has now bitten `dialog` and
  `sidebar`; do not make it three.
- **Public JS hook attributes are `data-gsxui-*`.** `data-slot` is markup
  identity, not a behavior contract. JS may match both, but every part JS
  binds to must carry a `data-gsxui-*` hook.
- **Components render state, never own persistence.** Take state as a
  parameter, reflect it, emit via `emit()` in `ui/gsxui.js`: `gsxui:change`
  with a detail payload for state-carrying parts, `gsxui:open`/`gsxui:close`
  for overlays, `gsxui:select` for selection.
- **`registry.Deps` is derived by go/parser over the generated `.x.go`.** A
  dependency existing only in a JS `import` is invisible to vendoring and is
  forbidden; duplicate rather than create one.
- **Any `site/examples/**/*.gsx` change requires `make highlight` in the same
  commit** — `site/hl` pin tests block CI otherwise and the deploy silently
  does not happen.
- **Adding a `ui/*.gsx` breaks `internal/registry/registry_test.go`'s pinned
  component list.** Updating it is expected, not scope creep.
- **gsx authoring gotchas:** `//` inside markup renders as literal page text
  (Go doc comments only); markup comments containing `<tag>` break the
  parser; `{{ x := }}` scopes function-wide; `<script>`/`<style>` bodies are
  raw text; inline `if/else` is not valid as a whole attribute value (only
  inside a `class={a,b,c}` comma-list); example directory names cannot
  contain hyphens.
- **Component parameter order:** params in declaration order, then
  `children gsx.Node`, then `attrs gsx.Attrs`. Parts with no children take
  only `attrs`.
- **Never start or kill processes.** Use `go tool gsx generate`,
  `go test ./...`, `make highlight`, `make check` only.

---

## File Structure

| file | responsibility |
|---|---|
| `docs/superpowers/plans/2026-07-25-tier4-source-map-menus.md` | Traced source of truth: shared item types, menubar, navigation-menu. |
| `ui/dropdown.gsx` / `.js` | Gains Group, CheckboxItem, RadioGroup, RadioItem, Sub, SubTrigger, SubContent. |
| `ui/context-menu.gsx` / `.js` | Same seven parts. |
| `ui/menubar.gsx` / `.js` | New: bar, triggers with roving tabindex, open-follows-hover, full item set. |
| `ui/navigation-menu.gsx` / `.js` | New: hover mega-menu with a shared viewport. |
| `ui/*_test.go` | Render pins per component. |
| `site/examples/{dropdown,contextmenu,menubar,navigationmenu}/*.gsx` | Example pages. |
| `docs/jsx-parity.md` | Ledger entries, including the §1 evidence table. |

---

## Task 0: Source map

**Files:**
- Create: `docs/superpowers/plans/2026-07-25-tier4-source-map-menus.md`

**Interfaces:**
- Consumes: nothing.
- Produces: the document Tasks 1–4 cite. Sections: `## shared-items`,
  `## menubar`, `## navigation-menu`.

Writes no code. Follows the convention set by
`docs/superpowers/plans/2026-07-24-tier4-source-map.md` — read it first and
match its structure, including the `derived-not-read` marker for any claim
that could not be traced to real source.

- [ ] **Step 1: Locate Radix's built output**

```bash
find ~/personal/shadcn-ui -maxdepth 6 -type d -name '@radix-ui' -not -path '*/.git/*' 2>/dev/null | head
ls ~/personal/shadcn-ui/node_modules/@radix-ui/ 2>/dev/null | head -30
```

Record found-with-path or absent. Absent means every ARIA/keyboard claim
about that primitive is `derived-not-read`. Do NOT guess.

- [ ] **Step 2: Write `## shared-items`**

From `new-york-v4/ui/dropdown-menu.tsx` and `context-menu.tsx`, for each of
Group, CheckboxItem, RadioGroup, RadioItem, Sub, SubTrigger, SubContent:
full class string verbatim, element, `data-slot`, and the indicator markup.

Then **diff the two sources against each other** and report explicitly
whether the parts are identical modulo the component prefix, listing every
difference. Task 2 depends on this answer.

Also cover the ARIA each part carries (`role="menuitemcheckbox"` /
`menuitemradio`, `aria-checked`, `aria-haspopup="menu"`, `aria-expanded`)
and the keyboard model for submenus (ArrowRight/ArrowLeft, Escape, typeahead
scope). Trace to Radix's dist or mark `derived-not-read`.

- [ ] **Step 3: Write `## menubar`**

From `new-york-v4/ui/menubar.tsx`: every part's full class string. Then the
two behaviors that distinguish a menubar from N dropdowns — roving tabindex
across triggers, and open-follows-hover once a menu is open. Trace both to
Radix's `react-menubar` dist or mark `derived-not-read`.

- [ ] **Step 4: Write `## navigation-menu`**

From `new-york-v4/ui/navigation-menu.tsx`: every part's class string,
including `navigationMenuTriggerStyle`. Document how the viewport gets its
size (which CSS custom properties, set by what), and state plainly whether
reproducing it needs continuous measurement or only discrete state — Task 4's
scope depends on that answer. Quote nova's `.cn-navigation-menu-*` rules.

- [ ] **Step 5: Commit**

```bash
git add docs/superpowers/plans/2026-07-25-tier4-source-map-menus.md
git commit -m "docs(plan): Tier 4 Batch B source map — menu family"
```

---

## Task 1: Shared item types on dropdown

**Files:**
- Modify: `ui/dropdown.gsx`, `ui/dropdown.js`, `ui/dropdown_test.go`
- Create: `site/examples/dropdown/checkboxes.gsx`, `site/examples/dropdown/radios.gsx`, `site/examples/dropdown/submenu.gsx`
- Modify: `docs/jsx-parity.md`, site wiring (mirror an existing dropdown example)

**Interfaces:**
- Consumes: shipped `ui/dropdown.gsx` machinery, `ui/icon`'s `Check`,
  `Circle`, `ChevronRight`.
- Produces:

```go
func DropdownMenuGroup(children gsx.Node, attrs gsx.Attrs) gsx.Node
func DropdownMenuCheckboxItem(checked bool, value string, children gsx.Node, attrs gsx.Attrs) gsx.Node
func DropdownMenuRadioGroup(value string, children gsx.Node, attrs gsx.Attrs) gsx.Node
func DropdownMenuRadioItem(checked bool, value string, children gsx.Node, attrs gsx.Attrs) gsx.Node
func DropdownMenuSub(children gsx.Node, attrs gsx.Attrs) gsx.Node
func DropdownMenuSubTrigger(children gsx.Node, attrs gsx.Attrs) gsx.Node
func DropdownMenuSubContent(children gsx.Node, attrs gsx.Attrs) gsx.Node
```

### Binding decisions

1. **`DropdownMenuSubContent` renders nested inside its `DropdownMenuSub`,
   which itself sits inside the parent `DropdownMenuContent`.** Not
   portalled. Spec §1 — a non-nested popover opened programmatically
   light-dismisses its parent, and submenus open on hover and ArrowRight.
   Verify this end to end before reporting done.
2. **Checked state is server-rendered.** `checked bool` stamps
   `aria-checked` and renders the indicator on first paint, so the menu is
   correct with JS disabled. JS only handles interaction. Emit
   `gsxui:change` with `{ checked, value }` from a checkbox item, and
   `{ value }` on the radio group from a radio item.
3. **`data-gsxui-*` hooks are required** on every part the JS binds to:
   `data-gsxui-menu-checkbox-item`, `-radio-item`, `-sub-trigger`,
   `-sub-content`. JS may also match `data-slot`, but the hook must exist.
4. **The parent content must not clip the submenu.** Check the shipped
   `DropdownMenuContent`'s overflow; if it clips, the fix is scoped to
   allowing the nested popover to escape, and must be ledgered rather than
   silently changing the pinned class string.
5. **Any `display` utility on the sub-content is gated on `:open`.**

- [ ] **Step 1: Write the failing tests**

Add to `ui/dropdown_test.go`. Pin the server-rendered checked state in both
directions and the submenu's nesting.

```go
func TestDropdownMenuCheckboxItemCheckedServerRendered(t *testing.T) {
	on := render(t, ui.DropdownMenuCheckboxItem(true, "show-toolbar", gsx.Raw("Toolbar"), nil))
	for _, want := range []string{
		`role="menuitemcheckbox"`,
		`aria-checked="true"`,
		`data-gsxui-menu-checkbox-item`,
		`data-value="show-toolbar"`,
	} {
		if !strings.Contains(on, want) {
			t.Errorf("want %q\nin: %s", want, on)
		}
	}
	off := render(t, ui.DropdownMenuCheckboxItem(false, "show-toolbar", gsx.Raw("Toolbar"), nil))
	if !strings.Contains(off, `aria-checked="false"`) {
		t.Errorf("want aria-checked=false\nin: %s", off)
	}
}

func TestDropdownMenuRadioItemCheckedServerRendered(t *testing.T) {
	got := render(t, ui.DropdownMenuRadioItem(true, "top", gsx.Raw("Top"), nil))
	for _, want := range []string{`role="menuitemradio"`, `aria-checked="true"`, `data-value="top"`} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
}

func TestDropdownMenuSubTriggerAria(t *testing.T) {
	got := render(t, ui.DropdownMenuSubTrigger(gsx.Raw("More"), nil))
	for _, want := range []string{
		`aria-haspopup="menu"`,
		`aria-expanded="false"`,
		`data-gsxui-menu-sub-trigger`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
}

func TestDropdownMenuSubContentIsANestedPopover(t *testing.T) {
	// Spec §1: a non-nested popover opened programmatically light-dismisses
	// its parent, so submenu content MUST be a popover nested in the parent.
	got := render(t, ui.DropdownMenuSubContent(gsx.Raw("x"), nil))
	for _, want := range []string{`popover="auto"`, `role="menu"`, `data-gsxui-menu-sub-content`} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
}

func TestDropdownMenuSubNestsContentInsideParentContent(t *testing.T) {
	got := render(t, ui.DropdownMenuContent(
		ui.DropdownMenuSub(
			gsx.Group(
				ui.DropdownMenuSubTrigger(gsx.Raw("More"), nil),
				ui.DropdownMenuSubContent(gsx.Raw("INNER"), nil),
			), nil), nil))
	outer := strings.Index(got, `data-slot="dropdown-menu-content"`)
	inner := strings.Index(got, "INNER")
	if outer < 0 || inner < 0 || inner < outer {
		t.Errorf("sub-content must render INSIDE the parent content\nin: %s", got)
	}
}
```

If `gsx.Group` is not the correct helper for adjacent nodes in this codebase,
use whatever `ui/*_test.go` already uses for the same purpose — check first.

- [ ] **Step 2: Run to verify failure**

```bash
go test ./ui/ -run TestDropdownMenu -v
```

Expected: FAIL, undefined symbols.

- [ ] **Step 3: Implement in `ui/dropdown.gsx`**

Port the seven parts from the source map's `## shared-items`, class strings
verbatim. Extend the file's doc comment with the submenu ADAPT, including the
spec §1 evidence table, and the server-rendered-checked MECHANISM.

```bash
go tool gsx generate
```

- [ ] **Step 4: Run to verify pass**

```bash
go test ./ui/ -run TestDropdownMenu -v
```

- [ ] **Step 5: Extend `ui/dropdown.js`**

- Click a checkbox item → flip `aria-checked` and the indicator, emit
  `gsxui:change` with `{ checked, value }`. Do NOT close the menu (Radix
  keeps it open for checkbox items unless told otherwise — confirm against
  the source map and follow it).
- Click a radio item → set `aria-checked="true"` on it, `"false"` on its
  siblings within the same radio group, emit `gsxui:change` with `{ value }`
  on the group.
- Sub-trigger: open on `pointerenter` and on ArrowRight, close on
  ArrowLeft/Escape, close when the pointer leaves the whole sub after a short
  grace period (match the shipped `hover-card.js` grace-period precedent).
  Stamp `data-state` on the trigger and `aria-expanded` in step with it, and
  stamp `data-state="open"` BEFORE `showPopover()`.
- Keyboard focus must move into the submenu's first item on ArrowRight and
  back to the trigger on ArrowLeft.

- [ ] **Step 6: Verify**

```bash
node --check ui/dropdown.js
go test ./...
```

- [ ] **Step 7: Site examples + highlight**

Three examples: `checkboxes.gsx`, `radios.gsx`, `submenu.gsx`. Then:

```bash
go tool gsx generate
make highlight
go test ./...
```

- [ ] **Step 8: Ledger**

Extend `docs/jsx-parity.md` `## dropdown`: the nested-not-portalled ADAPT with
the evidence table, the server-rendered-checked MECHANISM, and any GAP for
sub-features not ported.

- [ ] **Step 9: Commit**

```bash
git add ui/dropdown.gsx ui/dropdown.x.go ui/dropdown.js ui/dropdown_test.go site/ docs/jsx-parity.md
git commit -m "feat(dropdown): checkbox/radio items and nested submenus"
```

---

## Task 2: Same item types on context-menu

**Files:**
- Modify: `ui/context-menu.gsx`, `ui/context-menu.js`, `ui/context-menu_test.go`
- Create: `site/examples/contextmenu/full.gsx`
- Modify: `docs/jsx-parity.md`, site wiring

**Interfaces:**
- Consumes: Task 1's decisions verbatim.
- Produces: the same seven parts under the `ContextMenu*` prefix, with
  identical signatures.

### Binding decisions

1. **Follow Task 1 exactly** unless the source map's `## shared-items`
   diff recorded a real difference between shadcn's two sources. Where it
   did, follow the map. Where it did not, the two components must be
   consistent — a gratuitous divergence between `dropdown` and
   `context-menu` is a defect.
2. **Do NOT extract shared JS into a new module.** `registry.Deps` is derived
   from the generated `.x.go`, so a JS-only shared module would be invisible
   to vendoring. Duplication between `dropdown.js` and `context-menu.js` is
   the correct outcome here and is already the established pattern between
   those two files — check how they currently duplicate and match it.
3. All of Task 1's binding decisions 1–5 apply unchanged.

- [ ] **Step 1: Write the failing tests**

Mirror Task 1's five tests against the `ContextMenu*` parts, with
`data-gsxui-menu-*` hooks and the same nesting assertion. Repeat the code —
do not write "same as Task 1"; the implementer may be reading this alone.

```go
func TestContextMenuCheckboxItemCheckedServerRendered(t *testing.T) {
	on := render(t, ui.ContextMenuCheckboxItem(true, "show-grid", gsx.Raw("Grid"), nil))
	for _, want := range []string{
		`role="menuitemcheckbox"`,
		`aria-checked="true"`,
		`data-gsxui-menu-checkbox-item`,
		`data-value="show-grid"`,
	} {
		if !strings.Contains(on, want) {
			t.Errorf("want %q\nin: %s", want, on)
		}
	}
	off := render(t, ui.ContextMenuCheckboxItem(false, "show-grid", gsx.Raw("Grid"), nil))
	if !strings.Contains(off, `aria-checked="false"`) {
		t.Errorf("want aria-checked=false\nin: %s", off)
	}
}

func TestContextMenuSubContentIsANestedPopover(t *testing.T) {
	got := render(t, ui.ContextMenuSubContent(gsx.Raw("x"), nil))
	for _, want := range []string{`popover="auto"`, `role="menu"`, `data-gsxui-menu-sub-content`} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
}
```

Add radio-item and sub-trigger equivalents in the same shape.

- [ ] **Step 2: Run to verify failure**

```bash
go test ./ui/ -run TestContextMenu -v
```

- [ ] **Step 3: Implement, generate, verify**

```bash
go tool gsx generate && go test ./ui/ -run TestContextMenu -v
```

- [ ] **Step 4: Extend `ui/context-menu.js`**

Same interaction model as Task 1's step 5.

- [ ] **Step 5: Example, highlight, ledger, commit**

```bash
go tool gsx generate && make highlight && go test ./... && make check
git add ui/context-menu.gsx ui/context-menu.x.go ui/context-menu.js ui/context-menu_test.go site/ docs/jsx-parity.md
git commit -m "feat(context-menu): checkbox/radio items and nested submenus"
```

---

## Task 3: menubar

**Files:**
- Create: `ui/menubar.gsx`, `ui/menubar.js`, `ui/menubar_test.go`
- Create: `site/examples/menubar/basic.gsx`, `site/examples/menubar/full.gsx`
- Modify: `ui/index.js`, `internal/registry/registry_test.go`, `docs/jsx-parity.md`, site wiring

**Interfaces:**
- Consumes: Task 1's item-type decisions and submenu mechanism.
- Produces: the `Menubar*` part set per the source map's `## menubar` —
  at minimum `Menubar`, `MenubarMenu`, `MenubarTrigger`, `MenubarContent`,
  `MenubarItem`, `MenubarCheckboxItem`, `MenubarRadioGroup`,
  `MenubarRadioItem`, `MenubarLabel`, `MenubarSeparator`, `MenubarShortcut`,
  `MenubarGroup`, `MenubarSub`, `MenubarSubTrigger`, `MenubarSubContent`.
  Signatures mirror Task 1's equivalents exactly.

### Binding decisions

1. **Roving tabindex across triggers** — the whole bar is one tab stop;
   ArrowLeft/ArrowRight move between triggers. Follow the shipped
   `ui/toggle-group.js` model, including its JS-normalized-at-init approach
   so the no-JS fallback leaves every trigger reachable.
2. **Open-follows-hover once open** — with one menu open, hovering a sibling
   trigger switches to it with no click. This is the only genuinely new
   interaction in the task. Sibling menus are NOT nested in one another, so
   this is a close-then-open; spec §1's nesting constraint applies only
   between a menu and its own submenus.
3. Item types reuse Task 1's decisions verbatim, including server-rendered
   checked state and `data-gsxui-*` hooks.

- [ ] **Step 1: Write the failing tests**

```go
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
```

Add an exact-match render pin for `Menubar`, `MenubarTrigger` and
`MenubarContent` in the house style — the class strings are the spec.

- [ ] **Step 2: Run to verify failure**

```bash
go test ./ui/ -run TestMenubar -v
```

- [ ] **Step 3: Implement `ui/menubar.gsx`, generate, verify**

```bash
go tool gsx generate && go test ./ui/ -run TestMenubar -v
```

- [ ] **Step 4: Write `ui/menubar.js`**

Roving tabindex, open-follows-hover, and the item/submenu interactions from
Task 1. Add `import "./menubar.js";` to `ui/index.js` in alphabetical
position.

- [ ] **Step 5: Verify, examples, highlight, ledger, commit**

```bash
node --check ui/menubar.js
go tool gsx generate && make highlight && go test ./... && make check
git add ui/menubar.gsx ui/menubar.x.go ui/menubar.js ui/menubar_test.go ui/index.js internal/registry/registry_test.go site/ docs/jsx-parity.md
git commit -m "feat(ui): add menubar — roving triggers, hover-to-switch menus"
```

---

## Task 4: navigation-menu

**Files:**
- Create: `ui/navigation-menu.gsx`, `ui/navigation-menu.js`, `ui/navigation-menu_test.go`
- Create: `site/examples/navigationmenu/basic.gsx`, `site/examples/navigationmenu/mega.gsx`
- Modify: `ui/index.js`, `internal/registry/registry_test.go`, `docs/jsx-parity.md`, site wiring

**Interfaces:**
- Consumes: nothing from Tasks 1–3; independent.
- Produces: `NavigationMenu`, `NavigationMenuList`, `NavigationMenuItem`,
  `NavigationMenuTrigger`, `NavigationMenuContent`, `NavigationMenuLink`,
  `NavigationMenuViewport`, `NavigationMenuIndicator`, plus the exported
  trigger-style helper if the source map shows one.

### Binding decisions

1. **Scope is set by the source map's `## navigation-menu` finding on
   viewport sizing.** If sizing the viewport to the active panel needs only
   discrete state or a resize observer, build it. If it needs per-frame
   tweening between panel sizes, ship the discrete version and ledger the
   morph as a GAP. Do not invent an animation the reference does not have,
   and do not silently skip one it does.
2. Hover and focus both open; `data-state` drives transitions; the popover
   family's discrete-transition block is reused byte-identically.
3. `display` utilities gated on `:open`.

- [ ] **Step 1: Write the failing tests**

```go
func TestNavigationMenuTriggerAria(t *testing.T) {
	got := render(t, ui.NavigationMenuTrigger(gsx.Raw("Products"), nil))
	for _, want := range []string{
		`aria-expanded="false"`,
		`data-gsxui-navigation-menu-trigger`,
		`data-state="closed"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
}

func TestNavigationMenuListIsAList(t *testing.T) {
	got := render(t, ui.NavigationMenuList(gsx.Raw("x"), nil))
	if !strings.Contains(got, `data-slot="navigation-menu-list"`) {
		t.Errorf("want the list slot\nin: %s", got)
	}
}
```

Add exact-match pins for the root, list, trigger and content.

- [ ] **Step 2: Run to verify failure**

```bash
go test ./ui/ -run TestNavigationMenu -v
```

- [ ] **Step 3: Implement, generate, verify**

```bash
go tool gsx generate && go test ./ui/ -run TestNavigationMenu -v
```

- [ ] **Step 4: Write `ui/navigation-menu.js`**

Hover/focus open with a grace period, `data-state` stamping, viewport sizing
per decision 1. Add the barrel import alphabetically.

- [ ] **Step 5: Verify, examples, highlight, ledger, commit**

```bash
node --check ui/navigation-menu.js
go tool gsx generate && make highlight && go test ./... && make check
git add ui/navigation-menu.gsx ui/navigation-menu.x.go ui/navigation-menu.js ui/navigation-menu_test.go ui/index.js internal/registry/registry_test.go site/ docs/jsx-parity.md
git commit -m "feat(ui): add navigation-menu — hover mega-menu with shared viewport"
```

---

## Task 5: Roll-up

**Files:**
- Modify: `docs/component-roadmap.md`, `README.md`, `CHANGELOG.md`

- [ ] **Step 1: Roadmap**

Mark `menubar` and `navigation-menu` SHIPPED with a deferred-sub-features
table. Note that `dropdown` and `context-menu` reached full parity.

- [ ] **Step 2: README**

Retire the "Dropdown checkbox/radio items + submenus" Post-v1 backlog entry —
this batch is exactly that item. Add `menubar` and `navigation-menu` to the
Components list (Overlay and Navigation respectively).

- [ ] **Step 3: CHANGELOG**

One dated section. Two `### Added` bullets for the new components, one
`### Changed` bullet noting dropdown and context-menu gained checkbox/radio
items and submenus. Same terse register as the existing entries — no process
narration.

- [ ] **Step 4: Verify**

```bash
make check
```

- [ ] **Step 5: Commit**

```bash
git add docs/component-roadmap.md README.md CHANGELOG.md
git commit -m "docs: Tier 4 Batch B roll-up"
```

- [ ] **Step 6: Whole-batch review**

Dispatch on the strongest model. Point it at the ledger's deferred minors.
Cross-cutting focus: consistency of the shared item types across dropdown,
context-menu and menubar — three implementations of the same seven parts is
exactly where drift hides.

- [ ] **Step 7: Live browser verification**

Required before calling the batch shipped (spec §7):

- A submenu open while its parent stays open — the spec §1 mechanism end to
  end, in all three of dropdown, context-menu and menubar.
- Hovering between top-level menubar menus with one already open.
- Keyboard: ArrowRight into a submenu, ArrowLeft and Escape back out, with
  focus landing somewhere sane at every step.
- Every closed menu and submenu computes `display: none` — no ghost boxes.
- navigation-menu panel switching and viewport sizing.

Assert on resolved geometry with transitions disabled; occluded tabs freeze
animation clocks and that artifact is not a defect.

---

## Self-review

**Spec coverage.** §1 submenu mechanism → Global Constraints + Task 1
decision 1 + tests in Tasks 1–2. §2 shared item types → Tasks 1 and 2. §3
menubar → Task 3. §4 navigation-menu → Task 4. §5 cross-cutting → Global
Constraints. §6 task sequence → matches. §7 verification → Task 5 steps 6–7.

**Placeholder scan.** No TBD/TODO. Every code step carries real code. The one
conditional — navigation-menu's viewport scope — is resolved by Task 0's
finding with both branches specified.

**Type consistency.** `DropdownMenuCheckboxItem(checked bool, value string, …)`
and its ContextMenu/Menubar twins share one signature shape, asserted by
tests in all three tasks. Event names are `gsxui:change` throughout, matching
`ui/gsxui.js`'s `emit()`.
