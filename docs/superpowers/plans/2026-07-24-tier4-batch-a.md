# Tier 4 Batch A Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Port shadcn's `resizable`, `combobox` and `sidebar` to gsxui, taking
coverage from 47 to 50 of shadcn's 61.

**Architecture:** Server-rendered gsx components carrying shadcn's class
strings token-for-token, with behavior in a companion `ui/<name>.js` that
progressively enhances a correct no-JS first paint. All three reuse shipped
machinery — resizable is self-contained pointer/keyboard work, combobox
composes the popover anchoring and select value model, sidebar composes
sheet/tooltip/collapsible. None of the three owns persistence: each takes
state as a parameter and emits `gsxui:change` when it changes.

**Tech Stack:** Go + gsx templating, Tailwind v4 utility classes, vanilla ES
modules, `go test` render pins.

**Design spec:** `docs/superpowers/specs/2026-07-24-tier4-batch-a-design.md`
— read it before Task 1. Its §0 (reference sources), §3 (sidebar mobile),
and §4 (components render state, they do not own it) are binding.

## Global Constraints

Every task's requirements implicitly include this section.

- **Markup reference is `~/personal/shadcn-ui/apps/v4/registry/new-york-v4/ui/<name>.tsx`.**
  It is the only remaining inline-utility form of the markup. Do NOT port
  from `registry/bases/{radix,base,aria}/` — those emit `.cn-*` named classes
  with the utilities extracted into `registry/styles/style-*.css`.
- **Visual reference is `~/personal/shadcn-ui/apps/v4/registry/styles/style-nova.css`.**
  Where it and `new-york-v4` disagree, nova wins; `bases/base/ui/<name>.tsx`
  is the tiebreak for markup-structure questions nova's CSS cannot answer.
  New ports adopt nova metric and interaction tokens from the start.
- **Standing house exceptions** (do not "fix" these to match nova):
  border→`ring-1` swaps are NOT adopted; keep `data-[orientation=…]:` not
  `data-horizontal:`; keep `focus-visible:ring-[3px]` not `ring-3`; keep
  `data-[state=…]` not `data-open:`.
- **Class carryover is token-for-token.** Every dropped, added or altered
  token is ledgered in `docs/jsx-parity.md` under the component's heading as
  GAP / ADAPT / FINDING, matching the existing entries' style.
- **No component chooses a persistence mechanism.** Take state as a
  parameter, reflect it in the DOM, emit an event. Event naming follows
  `ui/gsxui.js`'s `emit()`: `gsxui:change` with a detail payload for
  state-carrying components, `gsxui:open`/`gsxui:close` for overlays,
  `gsxui:select` for selection.
- **`registry.Deps` is derived by go/parser over the generated `.x.go`.** A
  dependency that exists only in a JS `import` is invisible to vendoring and
  is therefore forbidden. If two components need the same JS, duplicate it.
- **`registry.HasJS(name)` requires `ui/<name>.js` named exactly.**
- **Any change under `site/examples/**/*.gsx` requires `make highlight` in
  the same commit.** `site/hl` pin tests fail otherwise and the deploy
  silently does not happen.
- **gsx authoring gotchas:** `//` inside markup renders as literal page text
  (Go doc comments only); markup comments containing `<tag>` break the
  parser; `{{ x := }}` scopes function-wide; `<script>`/`<style>` bodies are
  raw text; inline `if/else` is not valid as a whole attribute value (only
  inside a `class={a,b,c}` comma-list); example directory names cannot
  contain hyphens.
- **Component parameter order** is params in declaration order, then
  `children gsx.Node`, then `attrs gsx.Attrs`. Parts with no children take
  only `attrs`.
- **Never start or kill processes.** No `gsx dev`, no killing anything on a
  port. If you need a build, use `go tool gsx generate` and `go test ./...`.
- **Verification commands:** `go tool gsx generate` after any `.gsx` edit,
  then `go test ./...`. `make check` adds `node --check` on every `ui/*.js`
  plus gofmt and generated-file drift checks.

---

## File Structure

| file | responsibility |
|---|---|
| `docs/superpowers/plans/2026-07-24-tier4-source-map.md` | Traced source of truth for all three components: every class string, ARIA attribute, and keyboard behavior, each marked read-from-source or `derived-not-read`. |
| `ui/resizable.gsx` | `ResizablePanelGroup` / `ResizablePanel` / `ResizableHandle` markup. |
| `ui/resizable.js` | Pointer-drag and keyboard resize; emits `gsxui:change`. |
| `ui/resizable_test.go` | Render pins. |
| `ui/combobox.gsx` | The 12 combobox parts. |
| `ui/combobox.js` | Filter-on-type, highlight, value model, form bridge; emits `gsxui:select`. |
| `ui/combobox_test.go` | Render pins. |
| `ui/sidebar.gsx` | The ~25 sidebar parts, including the dual desktop/mobile tree. |
| `ui/sidebar.js` | Toggle, `Cmd/Ctrl+B`; emits `gsxui:change`. |
| `ui/sidebar_test.go` | Render pins. |
| `site/examples/resizable/*.gsx` | Example pages. |
| `site/examples/combobox/*.gsx` | Example pages. |
| `site/examples/sidebar/*.gsx` | Example pages, including the cookie round-trip. |
| `docs/jsx-parity.md` | Per-component divergence ledger (append three sections). |
| `docs/component-roadmap.md` | Tier 4 status update (Task 4). |

---

## Task 0: Source map

**Files:**
- Create: `docs/superpowers/plans/2026-07-24-tier4-source-map.md`

**Interfaces:**
- Consumes: nothing.
- Produces: the document every later task cites for class strings and
  behavioral claims. Later tasks reference it by section heading:
  `## resizable`, `## combobox`, `## sidebar`.

This task writes no code. Its deliverable is a document that makes Tasks 1–3
mechanical.

**Method — this is the part that matters.** Every behavioral claim must be
traced to real source, not inferred from the `.tsx` wrapper. The `.tsx` files
only show classes; the ARIA attributes and keyboard handling come from the
underlying library's built output. Where a library is not present in the
local checkout, mark the claim `derived-not-read` explicitly rather than
asserting it — this is the convention
`docs/superpowers/plans/2026-07-24-tier3-source-map-wrapped.md` established.

- [ ] **Step 1: Locate the libraries**

```bash
ls ~/personal/shadcn-ui/node_modules/react-resizable-panels/dist/ 2>/dev/null
ls ~/personal/shadcn-ui/node_modules/@base-ui/react/ 2>/dev/null
find ~/personal/shadcn-ui -maxdepth 6 -type d \( -name 'react-resizable-panels' -o -name '*base-ui*' \) -not -path '*/.git/*' 2>/dev/null | head
```

Record for each: found (with path) or absent. Absent means every claim about
that library is `derived-not-read`.

- [ ] **Step 2: Write `## resizable`**

Cover, from `new-york-v4/ui/resizable.tsx` and the four
`registry/new-york-v4/examples/resizable-*.tsx` demos:

- The three parts' full class strings, verbatim.
- The group's orientation mechanism: `aria-[orientation=vertical]:flex-col`.
- **The handle's inverted orientation.** The handle's base is `w-px` (a
  vertical rule) and its `aria-[orientation=horizontal]:` variants switch it
  to `h-px w-full` (a horizontal rule). A horizontal rule belongs in a
  *vertical* (flex-col) group. State the mapping explicitly as a table:
  group orientation → handle `aria-orientation` value. Confirm against
  `react-resizable-panels` if present; mark `derived-not-read` if not.
- `defaultSize` format in the current demos (`"50%"`, a percentage string).
- Keyboard behavior of the library's `Separator`: which keys, what step size,
  Home/End semantics. Confirm or mark `derived-not-read`.
- Nova's `.cn-resizable-handle-icon` rule, quoted.

- [ ] **Step 3: Write `## combobox`**

Cover, from `new-york-v4/ui/combobox.tsx`:

- Every part's full class string, verbatim, including the long
  `ComboboxContent` string.
- The `ComboboxInput` composition: `InputGroup` > `InputGroupInput` +
  `InputGroupAddon align="inline-end"` > trigger/clear, and the
  `showTrigger`/`showClear` gating.
- **Resolve the filter question.** Read `@base-ui/react`'s Combobox filter
  default if the package is present. Report: does it rank (returning a score)
  or filter (returning a boolean)? Is it collator-based? Case/accent
  handling? If the package is absent, say so and mark `derived-not-read` —
  do NOT guess.
- The ARIA anatomy: what roles/attributes Base UI stamps on input, list,
  items, and the active-descendant relationship. Trace or mark
  `derived-not-read`.
- Nova's `.cn-combobox-*` rules, quoted.

- [ ] **Step 4: Write `## sidebar`**

Cover, from `new-york-v4/ui/sidebar.tsx`:

- The six constants (`SIDEBAR_COOKIE_NAME`, `SIDEBAR_COOKIE_MAX_AGE`,
  `SIDEBAR_WIDTH`, `SIDEBAR_WIDTH_MOBILE`, `SIDEBAR_WIDTH_ICON`,
  `SIDEBAR_KEYBOARD_SHORTCUT`) with their values.
- Every one of the ~25 parts: its element, `data-slot`, and full class
  string, verbatim. This is the bulk of the document.
- The root `Sidebar`'s three render branches (`collapsible="none"`, mobile
  Sheet, desktop) with each branch's markup.
- Every `data-*` the desktop root stamps and which descendant selectors
  consume it (`group-data-[collapsible=…]`, `peer-data-[variant=…]`,
  `[[data-side=left][data-state=collapsed]_&]`, etc.). A table.
- `SidebarMenuButton`'s tooltip gating and its variant/size axes.
- Nova's `.cn-sidebar-*` rules, quoted.

- [ ] **Step 5: Commit**

```bash
git add docs/superpowers/plans/2026-07-24-tier4-source-map.md
git commit -m "docs(plan): Tier 4 Batch A source map — resizable, combobox, sidebar"
```

---

## Task 1: resizable

**Files:**
- Create: `ui/resizable.gsx`, `ui/resizable.js`, `ui/resizable_test.go`
- Create: `site/examples/resizable/basic.gsx`, `site/examples/resizable/vertical.gsx`, `site/examples/resizable/handle.gsx`
- Modify: `ui/index.js` (add `import "./resizable.js";` in alphabetical position)
- Modify: `docs/jsx-parity.md` (append `## resizable`)
- Modify: `site/main.go` or the site's component registry — follow how `carousel` is wired; grep for `"carousel"` under `site/` and mirror every hit.

**Interfaces:**
- Consumes: nothing from other tasks, and — per decision 6 below — **no
  `ui/icon` dependency**. `resizable` has no `Deps` at all.
- Produces:

```go
func ResizablePanelGroup(orientation string, children gsx.Node, attrs gsx.Attrs) gsx.Node
func ResizablePanel(defaultSize, minSize, maxSize string, children gsx.Node, attrs gsx.Attrs) gsx.Node
func ResizableHandle(orientation string, withHandle bool, attrs gsx.Attrs) gsx.Node
```

`orientation` is `"horizontal"` (or `""`, the Go zero value) or `"vertical"`,
and on **both** the group and the handle it names the **group's** orientation
— the author never reasons about the handle's inverted `aria-orientation`,
the component derives it. `defaultSize`/`minSize`/`maxSize` are percentage
strings like `"50%"`; `""` means unset.

### Binding decisions

1. **Sizes are server-rendered inline styles.** `defaultSize` becomes
   `style="flex-basis: 50%"` on the panel. A panel with `defaultSize=""`
   gets `flex-basis` omitted and relies on `flex-1`. The split is therefore
   correct on first paint with JS disabled. `minSize`/`maxSize` are stamped
   as `data-min-size` / `data-max-size` and read by JS during drag only.
2. **Handle `aria-orientation` is inverted from the group.** Group
   `horizontal` → handle `aria-orientation="vertical"`; group `vertical` →
   handle `aria-orientation="horizontal"`. This is forced by shadcn's own
   class string (`aria-[orientation=horizontal]:h-px aria-[orientation=horizontal]:w-full`
   is a horizontal rule, which only belongs in a flex-col group). Confirm
   against the source map's table; if the map found the library disagrees,
   STOP and report rather than silently choosing.
3. **`gsxui:change`, not `autoSaveId`.** On drag end and on keyboard commit,
   emit `gsxui:change` on the group with
   `{ sizes: [<percent number>, …] }` — one entry per panel, in DOM order.
   Do NOT write `localStorage`. See spec §4.
4. **Keyboard step is 10 percentage points**, Home/End drive the panel
   before the handle to its min/max (or 0%/100% when unset). If the source
   map found the library's real step, use that value instead and note it.
5. **`gsx.Attrs` merge**: `class` passed by the caller merges through the
   house class-merge, same as every other component. Do not hand-concatenate.
6. **The grip is a nova pill, not new-york-v4's icon-in-a-box.** The source
   map (`## resizable`, "Handle icon") established via the `bases/base`
   structure tiebreak that nova's `.cn-resizable-handle-icon` applies to an
   **empty** `div` — `{withHandle && <div className="cn-resizable-handle-icon z-10 flex shrink-0" />}`
   — and the rule is `@apply bg-border h-6 w-1 rounded-lg`. So `withHandle`
   renders exactly:

```html
<div class="z-10 flex shrink-0 h-6 w-1 rounded-lg bg-border"></div>
```

   NOT new-york-v4's bordered `h-4 w-3` box containing a `size-2.5`
   `GripVerticalIcon`. Nova wins on visual disagreement, and this is a
   structural delta rather than a token swap, so it is called out here
   explicitly. **Consequence: do not import `ui/icon`** — `resizable` has no
   dependencies. Ledger it as an ADAPT citing the source map.

- [ ] **Step 1: Write the failing test**

Create `ui/resizable_test.go`. Pin the group and panel exactly; assert the
handle's inverted orientation both ways.

```go
package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestResizablePanelGroupPinned(t *testing.T) {
	got := render(t, ui.ResizablePanelGroup("", gsx.Raw("x"), nil))
	want := `<div data-slot="resizable-panel-group" data-gsxui-resizable aria-orientation="horizontal" class="flex h-full w-full aria-[orientation=vertical]:flex-col">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestResizablePanelGroupVertical(t *testing.T) {
	got := render(t, ui.ResizablePanelGroup("vertical", gsx.Raw("x"), nil))
	if !strings.Contains(got, `aria-orientation="vertical"`) {
		t.Errorf("want vertical group orientation\nin: %s", got)
	}
}

func TestResizablePanelDefaultSizeIsInlineFlexBasis(t *testing.T) {
	// Server-rendered geometry: the split is correct on first paint with
	// JS disabled, so defaultSize must land as a real inline style.
	got := render(t, ui.ResizablePanel("50%", "", "", gsx.Raw("x"), nil))
	if !strings.Contains(got, `style="flex-basis: 50%"`) {
		t.Errorf("want inline flex-basis\nin: %s", got)
	}
}

func TestResizablePanelUnsizedHasNoFlexBasis(t *testing.T) {
	got := render(t, ui.ResizablePanel("", "", "", gsx.Raw("x"), nil))
	if strings.Contains(got, "flex-basis") {
		t.Errorf("unsized panel must not stamp flex-basis\nin: %s", got)
	}
}

func TestResizableHandleOrientationIsInverted(t *testing.T) {
	// The handle's aria-orientation names the RULE, not the group: a
	// horizontal rule (h-px w-full) divides a vertical (flex-col) group.
	// Callers pass the group's orientation; the component inverts.
	h := render(t, ui.ResizableHandle("horizontal", false, nil))
	if !strings.Contains(h, `aria-orientation="vertical"`) {
		t.Errorf("horizontal group wants a vertical rule\nin: %s", h)
	}
	v := render(t, ui.ResizableHandle("vertical", false, nil))
	if !strings.Contains(v, `aria-orientation="horizontal"`) {
		t.Errorf("vertical group wants a horizontal rule\nin: %s", v)
	}
}

func TestResizableHandleWithHandleRendersNovaPill(t *testing.T) {
	// nova's grip is an empty pill, not new-york-v4's icon-in-a-box.
	got := render(t, ui.ResizableHandle("horizontal", true, nil))
	for _, want := range []string{
		`role="separator"`,
		`tabindex="0"`,
		"h-6 w-1 rounded-lg bg-border",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
	if strings.Contains(got, "svg") {
		t.Errorf("nova's grip carries no icon glyph\nin: %s", got)
	}
}

func TestResizableHandleWithoutHandleHasNoGrip(t *testing.T) {
	got := render(t, ui.ResizableHandle("horizontal", false, nil))
	if strings.Contains(got, "h-6 w-1 rounded-lg bg-border") {
		t.Errorf("withHandle=false must not render the grip\nin: %s", got)
	}
}
```

- [ ] **Step 2: Run the test to verify it fails**

```bash
go test ./ui/ -run TestResizable -v
```

Expected: FAIL — `undefined: ui.ResizablePanelGroup`.

- [ ] **Step 3: Write `ui/resizable.gsx`**

Port the three parts from the source map's `## resizable`, class strings
verbatim. Open with a Go doc comment in the house style (see
`ui/toggle-group.gsx`) recording: the inverted-`aria-orientation` ADAPT with
its reasoning, the server-rendered `flex-basis` MECHANISM, and the
`autoSaveId`-dropped-for-`gsxui:change` GAP.

Then regenerate:

```bash
go tool gsx generate
```

- [ ] **Step 4: Run the test to verify it passes**

```bash
go test ./ui/ -run TestResizable -v
```

Expected: PASS. If a pinned string mismatches, fix the **component** to match
the source map, not the test to match the component — unless the source map
itself is wrong, in which case fix the map and say so in the report.

- [ ] **Step 5: Write `ui/resizable.js`**

Follow `ui/carousel.js`'s shape: a module with a `data-gsxui-resizable` root
selector, delegated listeners registered through `ui/gsxui.js`'s helpers, and
`emit()` for events.

Required behavior:

- `pointerdown` on `[data-slot="resizable-handle"]` → `setPointerCapture` on
  the handle, record the group's content-box size and the two adjacent
  panels' current pixel sizes.
- `pointermove` → convert the delta to a percentage of the group's content
  box, clamp against both neighbours' `data-min-size`/`data-max-size`, write
  both `flex-basis` values.
- `pointerup` / `lostpointercapture` → release, then `emit(group,
  "gsxui:change", { sizes })`.
- `keydown` on a focused handle: ArrowLeft/ArrowRight when the group is
  horizontal, ArrowUp/ArrowDown when vertical, ±10 points; Home/End to
  min/max. Each commit emits `gsxui:change`.
- Keep `aria-valuenow` / `aria-valuemin` / `aria-valuemax` on the handle in
  sync with the panel before it.

Add the barrel import in alphabetical position:

```js
import "./resizable.js";
```

- [ ] **Step 6: Verify JS syntax and the full suite**

```bash
node --check ui/resizable.js
go test ./...
```

Expected: both clean. `TestRegistry*` should pick up `resizable` with
`HasJS` true and **empty `Deps`** (no `icon` — see decision 6).

- [ ] **Step 7: Write the three site examples**

`site/examples/resizable/basic.gsx` (horizontal two-panel, nested vertical
group inside the right panel — port `resizable-demo-with-handle.tsx`),
`vertical.gsx` (port `resizable-vertical.tsx`), `handle.gsx` (the
`withHandle` grip variant). Wire them the way `carousel`'s examples are wired
— grep `site/` for `"carousel"` and mirror every hit.

- [ ] **Step 8: Regenerate highlighting and test**

```bash
go tool gsx generate
make highlight
go test ./...
```

`make highlight` is mandatory — `site/hl` pin tests fail without it and the
deploy silently does not happen.

- [ ] **Step 9: Ledger the divergences**

Append a `## resizable` section to `docs/jsx-parity.md` in the existing
style: the inverted-orientation ADAPT, the nova-pill grip ADAPT (decision 6,
citing the source map's `bases/base` tiebreak), the server-rendered
`flex-basis` MECHANISM, and GAPs for `autoSaveId` (replaced by
`gsxui:change`), collapsible panels, and the imperative panel API.

- [ ] **Step 10: Commit**

```bash
git add ui/resizable.gsx ui/resizable.x.go ui/resizable.js ui/resizable_test.go ui/index.js site/ docs/jsx-parity.md
git commit -m "feat(ui): add resizable — flex-basis panels with pointer and keyboard drag"
```

---

## Task 2: combobox

**Files:**
- Create: `ui/combobox.gsx`, `ui/combobox.js`, `ui/combobox_test.go`
- Create: `site/examples/combobox/basic.gsx`, `site/examples/combobox/groups.gsx`, `site/examples/combobox/clear.gsx`, `site/examples/combobox/form.gsx`
- Modify: `ui/index.js`, `docs/jsx-parity.md`, the site wiring (mirror `carousel`).

**Interfaces:**
- Consumes: shipped `ui.InputGroup` / `ui.InputGroupInput` /
  `ui.InputGroupAddon` / `ui.InputGroupButton`; `ui/icon`'s `Check`,
  `ChevronDown`, `X`; the popover discrete-transition class block used
  byte-identically by `ui/dropdown.gsx` and `ui/select.gsx`.
- Produces:

```go
func Combobox(name, value string, children gsx.Node, attrs gsx.Attrs) gsx.Node
func ComboboxInput(placeholder string, showTrigger, showClear, disabled bool, children gsx.Node, attrs gsx.Attrs) gsx.Node
func ComboboxTrigger(attrs gsx.Attrs) gsx.Node
func ComboboxClear(attrs gsx.Attrs) gsx.Node
func ComboboxContent(children gsx.Node, attrs gsx.Attrs) gsx.Node
func ComboboxList(children gsx.Node, attrs gsx.Attrs) gsx.Node
func ComboboxItem(value string, selected bool, children gsx.Node, attrs gsx.Attrs) gsx.Node
func ComboboxGroup(children gsx.Node, attrs gsx.Attrs) gsx.Node
func ComboboxLabel(children gsx.Node, attrs gsx.Attrs) gsx.Node
func ComboboxEmpty(children gsx.Node, attrs gsx.Attrs) gsx.Node
func ComboboxSeparator(attrs gsx.Attrs) gsx.Node
func ComboboxValue(children gsx.Node, attrs gsx.Attrs) gsx.Node
```

### Binding decisions

1. **Port the component, not the legacy recipe.**
   `registry/new-york-v4/examples/combobox-demo.tsx` is the old
   Popover+Command composition and is NOT the reference. The live docs page
   renders `combobox-basic`, `-clear`, `-groups`, `-invalid`, `-disabled`,
   `-input-group`, all of which are the component. Nova shipping
   `.cn-combobox-*` confirms it.
2. **No shared JS with `command`.** `combobox.js` gets its own matcher.
   Importing `command.js`'s scorer would create a dependency invisible to
   `registry.Deps` and silently break vendoring. Duplicate if needed.
3. **Filter semantics: an `Intl.Collator`-backed boolean `contains`.**
   RESOLVED — the source map correctly marked this `derived-not-read`
   (`@base-ui/react` is absent from the local checkout), and the controller
   then resolved it from Base UI's published docs and issue tracker: the
   default filter is `contains` from Base UI's `useFilter()` hook, which
   does case- and accent-insensitive matching via `Intl.Collator` with
   `sensitivity: "base"`. It returns a **boolean**, not a score — there is
   no ranking and no result reordering.

   Implement it as the collator scan that Base UI's `useFilter` mirrors
   from React Aria:

```js
const collator = new Intl.Collator(undefined, { usage: "search", sensitivity: "base" });

function contains(string, substring) {
  if (substring.length === 0) return true;
  const s = string.normalize("NFC");
  const q = substring.normalize("NFC");
  for (let i = 0; i <= s.length - q.length; i++) {
    if (collator.compare(s.slice(i, i + q.length), q) === 0) return true;
  }
  return false;
}
```

   Because items are never reordered, `combobox.js` only hides and shows
   them — do NOT port `command.js`'s DOM-reordering pass. Ledger this as a
   `web-verified` ADAPT (docs + issue tracker, not the package source) and
   record the divergence from cmdk's `command-score` ranking.
4. **Value model and form bridge come from `ui/select.js`'s pattern**: a
   hidden `sr-only` real control carries the value so non-JS form posts and
   `FormData` both work, populated by JS at init. `Combobox`'s `name`
   parameter names it. Emit `gsxui:select` with `{ value }` on the root and
   a native `change` on the bridge — identical to `select.js:64-66`, so the
   two components are interchangeable to a consumer.
5. **Open/close uses the shipped discrete-transition block**, copied
   byte-identically from `ui/select.gsx`'s content, and stamps
   `data-state="open"` synchronously BEFORE `showPopover()` — the
   flash-avoidance rule (`docs/jsx-parity.md ## animations`).
6. **Deferred parts, ledgered as GAP, not built:** `ComboboxChips`,
   `ComboboxChip`, `ComboboxChipsInput`, `ChipRemove` (multi-select);
   `ComboboxCollection` (a React data-mapping helper — gsx callers write a
   `for` loop); `useComboboxAnchor` (a React ref hook — gsx callers pass an
   anchor selector).

- [ ] **Step 1: Write the failing test**

Create `ui/combobox_test.go`.

```go
package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestComboboxItemSelectedStampsIndicator(t *testing.T) {
	got := render(t, ui.ComboboxItem("next.js", true, gsx.Raw("Next.js"), nil))
	for _, want := range []string{
		`data-slot="combobox-item"`,
		`data-value="next.js"`,
		`aria-selected="true"`,
		`data-slot="combobox-item-indicator"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
}

func TestComboboxItemUnselected(t *testing.T) {
	got := render(t, ui.ComboboxItem("next.js", false, gsx.Raw("Next.js"), nil))
	if !strings.Contains(got, `aria-selected="false"`) {
		t.Errorf("want aria-selected=false\nin: %s", got)
	}
}

func TestComboboxRootRendersFormBridge(t *testing.T) {
	// Non-JS form posts must carry the value: a real named control is
	// server-rendered, mirroring ui/select.gsx's bridge.
	got := render(t, ui.Combobox("framework", "next.js", gsx.Raw("x"), nil))
	for _, want := range []string{
		`data-slot="combobox"`,
		`name="framework"`,
		`value="next.js"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
}

func TestComboboxInputComposesInputGroup(t *testing.T) {
	got := render(t, ui.ComboboxInput("Search framework...", true, false, false, nil, nil))
	for _, want := range []string{
		`data-slot="input-group"`,
		`data-slot="input-group-input"`,
		`data-slot="combobox-trigger"`,
		`role="combobox"`,
		`aria-expanded="false"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
	if strings.Contains(got, `data-slot="combobox-clear"`) {
		t.Errorf("showClear=false must not render the clear button\nin: %s", got)
	}
}

func TestComboboxInputShowClear(t *testing.T) {
	got := render(t, ui.ComboboxInput("", false, true, false, nil, nil))
	if !strings.Contains(got, `data-slot="combobox-clear"`) {
		t.Errorf("want the clear button\nin: %s", got)
	}
	if strings.Contains(got, `data-slot="combobox-trigger"`) {
		t.Errorf("showTrigger=false must not render the trigger\nin: %s", got)
	}
}

func TestComboboxContentCarriesDiscreteTransitionBlock(t *testing.T) {
	// The popover family's shared exit-animation mechanism — must be
	// byte-identical to ui/select.gsx's content, not re-derived.
	got := render(t, ui.ComboboxContent(gsx.Raw("x"), nil))
	for _, want := range []string{
		`popover="auto"`,
		`data-slot="combobox-content"`,
		"transition-discrete",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
}

func TestComboboxEmptyIsHiddenUntilListIsEmpty(t *testing.T) {
	got := render(t, ui.ComboboxEmpty(gsx.Raw("No framework found."), nil))
	for _, want := range []string{
		`data-slot="combobox-empty"`,
		"group-data-empty/combobox-content:flex",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
}
```

- [ ] **Step 2: Run the test to verify it fails**

```bash
go test ./ui/ -run TestCombobox -v
```

Expected: FAIL — `undefined: ui.ComboboxItem`.

- [ ] **Step 3: Write `ui/combobox.gsx`**

Port all twelve parts from the source map's `## combobox`, class strings
verbatim. Doc comment records: decision 2 (no shared JS with `command`, with
the `registry.Deps` reasoning), decision 3 (the filter ADAPT and whether it
was read or derived), decision 4 (the form bridge), and the deferred parts as
GAPs.

```bash
go tool gsx generate
```

- [ ] **Step 4: Run the test to verify it passes**

```bash
go test ./ui/ -run TestCombobox -v
```

Expected: PASS.

- [ ] **Step 5: Write `ui/combobox.js`**

Model it on `ui/select.js` (347 lines) — same value model, same bridge, same
open/close. The combobox-specific additions:

- `input` event on the input filters `[data-slot="combobox-item"]`, hiding
  non-matches and toggling `data-empty` on the content so
  `ComboboxEmpty` appears.
- Highlight tracking with `data-highlighted` and
  `aria-activedescendant` on the input — focus stays in the input the whole
  time, the `command.js` model.
- ArrowUp/ArrowDown move the highlight, Enter commits, Escape closes.
- Trigger toggles open; clear resets value and input text, then emits.
- On commit: set the bridge value, emit `gsxui:select` with `{ value }` plus
  a native `change` on the bridge, close.

Add to `ui/index.js` in alphabetical position:

```js
import "./combobox.js";
```

- [ ] **Step 6: Verify**

```bash
node --check ui/combobox.js
go test ./...
```

- [ ] **Step 7: Write the four site examples**

`basic.gsx`, `groups.gsx` (grouped items with labels and a separator),
`clear.gsx` (`showClear`), `form.gsx` (inside a `<form>` demonstrating that
the bridge carries the value on submit). Wire as `carousel`'s examples are.

- [ ] **Step 8: Regenerate and test**

```bash
go tool gsx generate
make highlight
go test ./...
```

- [ ] **Step 9: Ledger**

Append `## combobox` to `docs/jsx-parity.md`: the filter ADAPT (naming
whether it was read or `derived-not-read`), the no-shared-JS-with-`command`
decision and its reasoning, the form-bridge MECHANISM, and GAPs for chips /
multi-select, `Collection`, and `useComboboxAnchor`.

- [ ] **Step 10: Commit**

```bash
git add ui/combobox.gsx ui/combobox.x.go ui/combobox.js ui/combobox_test.go ui/index.js site/ docs/jsx-parity.md
git commit -m "feat(ui): add combobox — filtered listbox on the shipped popover machinery"
```

---

## Task 3: sidebar

**Files:**
- Create: `ui/sidebar.gsx`, `ui/sidebar.js`, `ui/sidebar_test.go`
- Create: `site/examples/sidebar/basic.gsx`, `site/examples/sidebar/variants.gsx`, `site/examples/sidebar/persisted.gsx`
- Modify: `ui/index.js`, `docs/jsx-parity.md`, the site wiring (mirror `carousel`).

**Interfaces:**
- Consumes: shipped `ui.Sheet*`, `ui.Tooltip*`, `ui.Button`, `ui.Input`,
  `ui.Separator`, `ui.Skeleton`; `ui/icon`'s `PanelLeft`.
- Produces (the full part set — every one of these must exist):

```go
func SidebarProvider(open bool, children gsx.Node, attrs gsx.Attrs) gsx.Node
func Sidebar(side, variant, collapsible string, children gsx.Node, attrs gsx.Attrs) gsx.Node
func SidebarTrigger(attrs gsx.Attrs) gsx.Node
func SidebarRail(attrs gsx.Attrs) gsx.Node
func SidebarInset(children gsx.Node, attrs gsx.Attrs) gsx.Node
func SidebarInput(attrs gsx.Attrs) gsx.Node
func SidebarHeader(children gsx.Node, attrs gsx.Attrs) gsx.Node
func SidebarFooter(children gsx.Node, attrs gsx.Attrs) gsx.Node
func SidebarSeparator(attrs gsx.Attrs) gsx.Node
func SidebarContent(children gsx.Node, attrs gsx.Attrs) gsx.Node
func SidebarGroup(children gsx.Node, attrs gsx.Attrs) gsx.Node
func SidebarGroupLabel(children gsx.Node, attrs gsx.Attrs) gsx.Node
func SidebarGroupAction(children gsx.Node, attrs gsx.Attrs) gsx.Node
func SidebarGroupContent(children gsx.Node, attrs gsx.Attrs) gsx.Node
func SidebarMenu(children gsx.Node, attrs gsx.Attrs) gsx.Node
func SidebarMenuItem(children gsx.Node, attrs gsx.Attrs) gsx.Node
func SidebarMenuButton(isActive bool, variant, size, tooltip string, children gsx.Node, attrs gsx.Attrs) gsx.Node
func SidebarMenuAction(showOnHover bool, children gsx.Node, attrs gsx.Attrs) gsx.Node
func SidebarMenuBadge(children gsx.Node, attrs gsx.Attrs) gsx.Node
func SidebarMenuSkeleton(showIcon bool, attrs gsx.Attrs) gsx.Node
func SidebarMenuSub(children gsx.Node, attrs gsx.Attrs) gsx.Node
func SidebarMenuSubItem(children gsx.Node, attrs gsx.Attrs) gsx.Node
func SidebarMenuSubButton(size string, isActive bool, children gsx.Node, attrs gsx.Attrs) gsx.Node
```

`side` is `"left"` (or `""`) / `"right"`; `variant` is `"sidebar"` (or `""`)
/ `"floating"` / `"inset"`; `collapsible` is `"offcanvas"` (or `""`) /
`"icon"` / `"none"`.

### Binding decisions

1. **`SidebarProvider` takes `open bool` and stamps
   `data-state="expanded"|"collapsed"` server-side.** No cookie, no
   `localStorage`, no media query. See spec §3 and §4.
2. **`sidebar.js` emits `gsxui:change` with `{ open }` on the provider** and
   flips `data-state` so the component works standalone. It must NOT write
   `document.cookie`. The cookie round-trip is Task 3's `persisted.gsx`
   example, not component code.
3. **Mobile is both trees, CSS-gated.** `Sidebar` renders the desktop tree
   (which already carries shadcn's own `hidden … md:block` root and
   `hidden … md:flex` container) AND a `Sheet`-based mobile tree marked
   `md:hidden` with the reference's `data-mobile="true"`,
   `--sidebar-width: 18rem` override, and `sr-only` header/title/description.
   No media-query JS. `collapsible="none"` still short-circuits to the
   single flat `div` branch, as the reference does.
4. **`children` therefore renders twice.** This is a known, accepted cost.
   Document it prominently in the component doc comment and ledger it as a
   GAP: sidebar content must not contain `id` attributes, because they would
   be duplicated and the document invalid. Use classes and `data-*`.
5. **No context substitute.** `data-state` / `data-collapsible` /
   `data-variant` / `data-side` are stamped on the desktop root and consumed
   by descendants through `group-data-*` / `peer-data-*` selectors — already
   pure CSS in the reference, so nothing needs threading. `SidebarTrigger`
   and `SidebarRail` resolve the provider from the DOM (`closest`), not a
   handle.
6. **`Cmd/Ctrl+B` toggles**, matching `SIDEBAR_KEYBOARD_SHORTCUT`. Register
   once at module scope, guard against repeat-fire, and `preventDefault`.

- [ ] **Step 1: Write the failing test**

Create `ui/sidebar_test.go`.

```go
package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestSidebarProviderStampsServerState(t *testing.T) {
	// The whole point of the design: state arrives as a parameter and is
	// rendered, so there is no flash and no hydration step.
	open := render(t, ui.SidebarProvider(true, gsx.Raw("x"), nil))
	if !strings.Contains(open, `data-state="expanded"`) {
		t.Errorf("want expanded\nin: %s", open)
	}
	closed := render(t, ui.SidebarProvider(false, gsx.Raw("x"), nil))
	if !strings.Contains(closed, `data-state="collapsed"`) {
		t.Errorf("want collapsed\nin: %s", closed)
	}
}

func TestSidebarProviderCarriesWidthVars(t *testing.T) {
	got := render(t, ui.SidebarProvider(true, gsx.Raw("x"), nil))
	for _, want := range []string{"--sidebar-width", "--sidebar-width-icon", `data-slot="sidebar-wrapper"`} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
}

func TestSidebarRendersBothTrees(t *testing.T) {
	// Mobile is CSS-gated, not JS-swapped: the Sheet tree and the desktop
	// tree both exist in the DOM, gated md:hidden / hidden md:block.
	got := render(t, ui.Sidebar("", "", "", gsx.Raw("CONTENT"), nil))
	if strings.Count(got, "CONTENT") != 2 {
		t.Errorf("want children rendered in both trees, got %d\nin: %s", strings.Count(got, "CONTENT"), got)
	}
	for _, want := range []string{`data-mobile="true"`, "md:hidden", "md:block"} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
}

func TestSidebarCollapsibleNoneIsFlat(t *testing.T) {
	// The reference short-circuits to one plain div; no gap, no container,
	// no Sheet, children rendered exactly once.
	got := render(t, ui.Sidebar("", "", "none", gsx.Raw("CONTENT"), nil))
	if strings.Count(got, "CONTENT") != 1 {
		t.Errorf("collapsible=none renders children once, got %d\nin: %s", strings.Count(got, "CONTENT"), got)
	}
	if strings.Contains(got, `data-mobile="true"`) {
		t.Errorf("collapsible=none must not render the Sheet tree\nin: %s", got)
	}
}

func TestSidebarStampsVariantSideCollapsible(t *testing.T) {
	got := render(t, ui.Sidebar("right", "floating", "icon", gsx.Raw("x"), nil))
	for _, want := range []string{`data-side="right"`, `data-variant="floating"`} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
}

func TestSidebarMenuButtonActiveAndTooltip(t *testing.T) {
	got := render(t, ui.SidebarMenuButton(true, "", "", "Inbox", gsx.Raw("Inbox"), nil))
	for _, want := range []string{
		`data-slot="sidebar-menu-button"`,
		`data-active="true"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
}

func TestSidebarMenuSkeletonShowIcon(t *testing.T) {
	with := render(t, ui.SidebarMenuSkeleton(true, nil))
	without := render(t, ui.SidebarMenuSkeleton(false, nil))
	if strings.Count(with, `data-slot="skeleton"`) <= strings.Count(without, `data-slot="skeleton"`) {
		t.Errorf("showIcon must add a skeleton\n with: %s\nwithout: %s", with, without)
	}
}
```

- [ ] **Step 2: Run the test to verify it fails**

```bash
go test ./ui/ -run TestSidebar -v
```

Expected: FAIL — `undefined: ui.SidebarProvider`.

- [ ] **Step 3: Write `ui/sidebar.gsx`**

Port all parts from the source map's `## sidebar`, class strings verbatim.
This is the largest file in the batch; work through the map's part table in
order. The doc comment must record decisions 1–6, with the duplicate-`id`
warning from decision 4 stated prominently.

```bash
go tool gsx generate
```

- [ ] **Step 4: Run the test to verify it passes**

```bash
go test ./ui/ -run TestSidebar -v
```

Expected: PASS.

- [ ] **Step 5: Write `ui/sidebar.js`**

- Click on `[data-slot="sidebar-trigger"]` or `[data-slot="sidebar-rail"]`
  → resolve the provider with `closest('[data-slot="sidebar-wrapper"]')`,
  flip `data-state`, `emit(provider, "gsxui:change", { open })`.
- On mobile the trigger must instead open the Sheet tree. Resolve which
  tree is visible via `getComputedStyle` / `offsetParent` rather than a
  hard-coded breakpoint, so a consumer's custom `md` still works.
- `Cmd/Ctrl+B` at module scope toggles the first provider on the page,
  `preventDefault`, no repeat-fire.
- No `document.cookie`. No `localStorage`.

Add to `ui/index.js` in alphabetical position:

```js
import "./sidebar.js";
```

- [ ] **Step 6: Verify**

```bash
node --check ui/sidebar.js
go test ./...
```

- [ ] **Step 7: Write the three site examples**

- `basic.gsx` — provider + sidebar with header, grouped menu, footer, inset
  content, trigger.
- `variants.gsx` — the three `variant`s and the `collapsible` modes.
- `persisted.gsx` — **the cookie round-trip as a documented pattern**: a Go
  handler snippet reading `sidebar_state` into the `open` parameter, plus the
  three-line `gsxui:change` listener that writes it back. This is what
  replaces shipping the cookie in the component; it must read as a complete,
  copyable recipe.

- [ ] **Step 8: Regenerate and test**

```bash
go tool gsx generate
make highlight
go test ./...
```

- [ ] **Step 9: Ledger**

Append `## sidebar` to `docs/jsx-parity.md`: the both-trees ADAPT with the
duplicate-`id` GAP, the no-context MECHANISM, and the
persistence-is-the-consumer's ADAPT covering the dropped cookie.

- [ ] **Step 10: Commit**

```bash
git add ui/sidebar.gsx ui/sidebar.x.go ui/sidebar.js ui/sidebar_test.go ui/index.js site/ docs/jsx-parity.md
git commit -m "feat(ui): add sidebar — CSS-gated dual tree, server-rendered state"
```

---

## Task 4: Roll-up

**Files:**
- Modify: `docs/component-roadmap.md`, `README.md`,
  `.superpowers/sdd/progress.md`

**Interfaces:**
- Consumes: Tasks 1–3 shipped.
- Produces: nothing code-level.

- [ ] **Step 1: Update the roadmap**

In `docs/component-roadmap.md` § Tier 4: mark resizable, combobox and sidebar
SHIPPED with a deferred-sub-features table in the shape Tier 3 uses. **Also
delete the duplicated rows** — the current Tier 4 table lists `menubar`,
`calendar`, `resizable`, `sidebar` and `chart` twice (lines 86–90 repeat
82–85). Add a short note recording the §0 registry-fork finding and that no
re-sync of shipped components is implied.

- [ ] **Step 2: Prune the README backlog**

`README.md` § Post-v1 backlog still lists three items that have shipped:
custom listbox select, the popover exit-animation strategy, and gsx syntax
highlighting. Remove those three. Leave the genuinely open ones: dropdown
checkbox/radio items + submenus, tooltip delay-groups, CSS anchor-positioning
migration, the checkbox `currentColor` mask, local `gsxui theme`, icon
search. Add the three new components to the § Components list.

- [ ] **Step 3: Full-suite verification**

```bash
make check
```

Expected: clean — tests pass, no `.x.go` drift, no untracked generated files,
every `ui/*.js` passes `node --check`, gofmt clean.

- [ ] **Step 4: Commit**

```bash
git add docs/component-roadmap.md README.md .superpowers/sdd/progress.md
git commit -m "docs: Tier 4 Batch A roll-up — roadmap status, README backlog prune"
```

- [ ] **Step 5: Whole-batch review**

Dispatch a review of the full batch diff on the strongest model, as Tiers 2
and 3 did. Its job is the between-task orphans no single task's reviewer can
see: class tokens that drifted between components, ledger entries that
contradict shipped code, barrel ordering, `Deps` correctness.

- [ ] **Step 6: Live browser verification**

Required before calling the batch shipped (spec §6). Side-by-side against the
matching `ui.shadcn.com` pages:

- **resizable** — drag both orientations, nested groups, keyboard resize with
  the handle focused, min/max clamping, `aria-valuenow` tracking.
- **combobox** — type to filter, empty state, keyboard highlight and commit,
  clear button, form submission carrying the value.
- **sidebar** — collapse and expand in all three `variant`s and both
  `collapsible` modes, `Cmd/Ctrl+B`, the rail, and the mobile tree at a
  narrow viewport.

Known environment limit: hidden or occluded Chrome tabs freeze animation
clocks, leave CSS transitions permanently `pending`, and never run
smooth-scroll or rAF callbacks. DOM-state assertions still hold. Do not
report those artifacts as bugs.

---

## Self-review

**Spec coverage.** §0 reference sources → Global Constraints + Task 4 Step 1.
§1 resizable → Task 1 (layout, drag, keyboard, ARIA, deferred list all
covered). §2 combobox → Task 2 (parts, reuse, filter question, deferred list).
§3 sidebar → Task 3 (state, mobile, parts, deferred). §4 cross-cutting → the
"No component chooses a persistence mechanism" constraint plus binding
decisions in Tasks 1 and 3. §5 execution → task order matches. §6
verification → Task 4 Steps 5–6. §7 out-of-scope → nothing in the plan
touches the 47 shipped components, the menu family, calendar, or chart.

**Placeholder scan.** No TBD/TODO. Every code step carries real code. The one
conditional is Task 2's filter semantics, which is resolved by Task 0's
finding with a concrete fallback specified — not left open.

**Type consistency.** `ResizableHandle(orientation, withHandle, attrs)` is
called with the group's orientation in every reference. `ComboboxItem(value,
selected, children, attrs)` matches its test. `SidebarProvider(open bool, …)`
matches its test and both `sidebar.js` and `persisted.gsx`. Event names are
`gsxui:change` (resizable, sidebar) and `gsxui:select` (combobox) throughout,
matching `ui/gsxui.js`'s `emit()`.
