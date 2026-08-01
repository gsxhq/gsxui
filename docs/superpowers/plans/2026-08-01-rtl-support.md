# RTL Support Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Every gsxui component renders and behaves correctly under `dir="rtl"`, following shadcn/ui's official RTL conventions, with RTL docs-site examples where shadcn has them.

**Architecture:** Three layers. (1) CSS: convert remaining physical-direction Tailwind classes in `ui/*.gsx` to logical ones (`ms-/me-/ps-/pe-/start-/end-`), `rtl:rotate-180` on directional icons. (2) JS: one `isRTL(el)` helper in `ui/gsxui.js`; `position()` mirrors alignment/sides when the anchor is RTL; direction-sensitive keyboard/pointer handlers flip. (3) Docs: an RTL page with an Arabic login-card demo plus RTL examples for Calendar, Pagination, Sidebar.

**Tech Stack:** Go + gsx components (`ui/*.gsx` compiled to `.x.go`), vanilla JS modules (`ui/*.js`), Tailwind v4, Playwright specs in `jstest/specs/`.

**Spec:** `docs/superpowers/specs/2026-08-01-rtl-support-design.md`

## Global Constraints

- Follow shadcn conventions exactly: `left-*`→`start-*`, `right-*`→`end-*`, `ml/mr/pl/pr`→`ms/me/ps/pe`, `text-left/right`→`text-start/end`, `border-l/r`→`border-s/e`, `rounded-l/r-*`→`rounded-s/e-*`, `slide-in-from-left/right`→`slide-in-from-start/end`; `rtl:` variants only for icon flips and gradients that can't be expressed logically.
- Do NOT convert genuinely physical usage: `left-1/2` centering (dialog, tooltip arrow, carousel vertical buttons), symmetric `translate-x` animation math, `rounded-lg`.
- gsx comments use `{/* */}` inside markup, never `//`.
- After editing any `.gsx`: run `make generate` (or the project's gsx build target per Makefile) so `.x.go` regenerates; run `gsx fmt` and confirm idempotency.
- Many `ui/*_test.go` files assert exact class strings — update expectations in the same commit as the class change; never weaken an assertion to a substring match to avoid updating it.
- No `git add -A`; add files explicitly. Commit after every task.
- Playwright: run via the project's configured runner (`cd jstest && npx playwright test <spec>`); check `jstest/playwright.config.ts` for the dev-server setup before assuming a URL.

---

### Task 1: CSS logical-class pass over `ui/*.gsx`

**Files:**
- Modify: `ui/combobox.gsx`, `ui/command.gsx`, `ui/context-menu.gsx`, `ui/dropdown-menu.gsx`, `ui/menubar.gsx`, `ui/select.gsx`, `ui/dialog.gsx`, `ui/sheet.gsx`, `ui/drawer.gsx`, `ui/table.gsx`, `ui/toast.gsx`, `ui/native-select.gsx`, `ui/navigation-menu.gsx`, `ui/input-group.gsx`, `ui/input-otp.gsx`, `ui/button-group.gsx`, `ui/toggle-group.gsx`, `ui/calendar.gsx`, `ui/sidebar.gsx`, `ui/carousel.gsx`, `ui/tooltip.gsx`
- Modify: corresponding `ui/*_test.go` class-string expectations
- Regenerate: corresponding `ui/*.x.go`

**Interfaces:**
- Produces: markup whose only remaining physical-direction classes are the deliberate exceptions listed below. Later tasks (7, 10) rely on sheet/drawer/sidebar using `start/end` + `data-side` selectors.

Exact conversions (grep-verified inventory as of 2026-08-01; re-grep before editing):

| File | Old → New |
|---|---|
| combobox.gsx:335,343 | `pr-8 pl-1.5`→`pe-8 ps-1.5`; `right-2`→`end-2` |
| command.gsx:164 | `ml-auto`→`ms-auto` |
| context-menu.gsx:156,173,212,229,257,291 | `pr-8 pl-1.5`→`pe-8 ps-1.5`; `right-2`→`end-2`; `ml-auto`→`ms-auto` (×2 incl. `[&>svg:last-child]:ml-auto`) |
| dropdown-menu.gsx:157,174,214,231,258,293 | same set as context-menu |
| menubar.gsx:317,334,367,384,418,469 | same set as context-menu |
| select.gsx:170,176 | `pr-8 pl-1.5`→`pe-8 ps-1.5`; `right-2`→`end-2` |
| dialog.gsx:37,67 | `right-2`→`end-2`; `sm:text-left`→`sm:text-start` (keep `top-1/2 left-1/2 -translate-x-1/2` centering) |
| sheet.gsx:35,39,53 | `left-0 right-auto`→`start-0 end-auto`; `border-r`→`border-e`; `slide-out-to-left/slide-in-from-left`→`…-start`; mirror for the right arm (`end-0 start-auto border-s …-end`); `right-3`→`end-3` |
| drawer.gsx:34,36 | same pattern as sheet + `rounded-r-xl`→`rounded-e-xl`, `rounded-l-xl`→`rounded-s-xl`, `md:…text-left`→`text-start` |
| table.gsx:47,58 | `text-left`→`text-start`; `pr-0`→`pe-0` |
| toast.gsx:122 | `-right-1.5`→`-end-1.5` |
| native-select.gsx:42 | `[&>svg]:right-2.5`→`[&>svg]:end-2.5` |
| navigation-menu.gsx:137,~162,~187 | `ml-1`→`ms-1`; `pr-2.5`→`pe-2.5` |
| input-group.gsx:50,56 | `pr-2`→`pe-2`; `pl-2`→`ps-2` |
| input-otp.gsx:134 | `border-r`→`border-e`; `first:rounded-l-lg first:border-l`→`first:rounded-s-lg first:border-s`; `last:rounded-r-lg`→`last:rounded-e-lg` |
| button-group.gsx (rule strings near :20) | `rounded-l-none`→`rounded-s-none`, `border-l-0`→`border-s-0`, `rounded-r-none`→`rounded-e-none` |
| toggle-group.gsx:140,142 | `border-l-0/first:border-l`→`border-s-0/first:border-s`; `first:rounded-l-lg last:rounded-r-lg`→`first:rounded-s-lg last:rounded-e-lg` |
| calendar.gsx:922 (and pr/pl in caption per :140,:190 comments' live class strings) | `rounded-l-md`→`rounded-s-md`, `rounded-r-md`→`rounded-e-md` (all range arms); `pr-1 pl-2`→`pe-1 ps-2`; `pr-1.5`/`pl-1.5` variants likewise |
| sidebar.gsx:119,156,273,295,309,337 | `border-r/border-l`→`border-e/border-s` is NOT enough here — these are keyed on `data-side=left/right`; convert per Task 7 instead. In THIS task convert only the side-agnostic ones: `text-left`→`text-start` (:273), `pr-8`→`pe-8` (:273), `right-1`→`end-1` (:295,:309), `border-l`→`border-s` (:337) |
| carousel.gsx:70,104,150-174 | `-ml-4/-ml-1`→`-ms-4/-ms-1`; `pl-4 -scroll-ml-4`→`ps-4 -scroll-ms-4`; horizontal prev/next `-left-12`→`-start-12`, `-right-12`→`-end-12` (keep vertical `left-1/2 -translate-x-1/2`) |
| tooltip.gsx:22 | keep — `data-[side=left/right]` translate arms are physical viewport geometry set by `position()`; keep `left-1/2` arrow centering |

Steps:

- [ ] **Step 1: Re-run the inventory grep** to catch drift since this plan was written:

```bash
grep -nE '\b(ml|mr|pl|pr)-[a-z0-9.\[\]]+|\b(left|right)-[0-9a-z.\[\]/]+|text-(left|right)\b|rounded-(l|r)(-[a-z]+)?\b|border-(l|r)\b|slide-(in-from|out-to)-(left|right)' ui/*.gsx | grep -v rounded-lg
```

- [ ] **Step 2: Apply the table above** file by file (class strings only; leave code comments describing history untouched).
- [ ] **Step 3: Regenerate + format**: run the Makefile's gsx generate target, then `gsx fmt` twice, diffing to confirm idempotency.
- [ ] **Step 4: Run `go test ./ui/...`** — update every class-string expectation that fails to the new logical string (verify each diff is exactly the intended conversion, nothing else).
- [ ] **Step 5: Re-run the Step 1 grep** — remaining hits must all be documented exceptions (dialog/tooltip/carousel-vertical centering, `data-[side=…]` translate arms, sidebar side-keyed rules deferred to Task 7). List them in the commit message.
- [ ] **Step 6: Commit** `git commit -m "feat(rtl): convert physical direction classes to logical properties"`

### Task 2: Directional icon flips (`rtl:rotate-180`)

**Files:**
- Modify: `ui/breadcrumb.gsx` (separator chevron), `ui/pagination.gsx` (prev/next chevrons), `ui/dropdown-menu.gsx` + `ui/context-menu.gsx` + `ui/menubar.gsx` (sub-trigger chevron: the `[&>svg:last-child]` targeted icon), `ui/calendar.gsx` (month prev/next chevrons), `ui/carousel.gsx` (horizontal prev/next arrow icons), `ui/sidebar.gsx` (SidebarTrigger icon), `ui/command.gsx`/`ui/combobox.gsx` if they render directional chevrons (inspect; the select caret is vertical — skip), plus example/demo `.gsx` files under `site/examples/` that hardcode `ArrowRight`-style icons only if they are part of the component recipe (leave pure demo content alone).
- Modify: matching `ui/*_test.go` expectations; regenerate `.x.go`.

**Interfaces:**
- Produces: every icon whose meaning is "forward/back along reading direction" carries `rtl:rotate-180`. Vertical chevrons (select, accordion, collapsible) unchanged.

- [ ] **Step 1: Find candidates**: `grep -n 'chevron-right\|chevron-left\|arrow-right\|arrow-left\|ChevronRight\|ChevronLeft' ui/*.gsx site/icons -r`
- [ ] **Step 2: Add `rtl:rotate-180`** to each directional icon's class (on the `<svg>`/icon element or via the parent's `[&>svg…]` arm, matching where the size classes already live).
- [ ] **Step 3: Regenerate, `gsx fmt`, `go test ./ui/...`**, update expectations.
- [ ] **Step 4: Commit** `git commit -m "feat(rtl): flip directional icons with rtl:rotate-180"`

### Task 3: `isRTL` helper + RTL-aware `position()`

**Files:**
- Modify: `ui/gsxui.js` (helper + `applyPlacement`)
- Test: `jstest/specs/positioning.spec.ts` (add RTL cases)

**Interfaces:**
- Produces: `export function isRTL(el)` in `ui/gsxui.js` — Tasks 4–6 import it. `position(content, anchor, opts)` API unchanged; under an RTL anchor, `align:"start"/"end"` and `side:"left"/"right"` resolve logically.

- [ ] **Step 1: Write failing Playwright test** in `jstest/specs/positioning.spec.ts` (follow the file's existing fixture pattern for mounting a trigger+content pair):

```ts
test("dropdown content start-aligns to the trigger's right edge under dir=rtl", async ({ page }) => {
  // reuse the spec file's existing dropdown fixture page, then:
  await page.evaluate(() => document.documentElement.setAttribute("dir", "rtl"));
  await page.getByRole("button", { name: /open/i }).click();
  const t = await page.locator("[data-gsxui-slot-dropdown-menu-trigger]").boundingBox();
  const c = await page.locator("[data-gsxui-slot-dropdown-menu-content]").boundingBox();
  expect(Math.abs((c!.x + c!.width) - (t!.x + t!.width))).toBeLessThan(1.5);
});
```

- [ ] **Step 2: Run it, confirm it fails** (content left-aligns: `c.x ≈ t.x`).
- [ ] **Step 3: Implement** in `ui/gsxui.js`:

```js
// isRTL reports the resolved direction at el — computed style, so it
// honors both dir attributes and CSS `direction`.
export function isRTL(el) {
  return getComputedStyle(el).direction === "rtl";
}
```

In `applyPlacement`, resolve logical options physically before the existing math (anchor rects have no element when a virtual rect is passed — resolve from `content`, which lives in the same direction context):

At the top of `applyPlacement` (before the existing `vertical`/flip math), resolve logical → physical:

```js
const rtl = isRTL(content);
let { side, align, alignOffset } = o;
const vertical = side === "top" || side === "bottom";
if (rtl) {
  if (vertical) {
    // horizontal cross-axis: mirror alignment and its offset
    if (align === "start") align = "end";
    else if (align === "end") align = "start";
    alignOffset = -alignOffset;
  } else {
    // horizontal main axis: mirror the preferred side
    side = side === "left" ? "right" : "left";
  }
}
```

Then use these locals in place of `o.side`/`o.align`/`o.alignOffset` throughout the rest of the function. Keep `data-side` reporting the PLACED physical side (animations key on physical geometry).

- [ ] **Step 4: Run the new test — PASS; run the whole positioning spec — no regressions.**
- [ ] **Step 5: Add one more RTL case**: context-menu submenu prefers physical LEFT under RTL (`side:"right"` authored → placed left). Assert `data-side="left"` and `sub.x + sub.width <= parent.x + 2`.
- [ ] **Step 6: Commit** `git commit -m "feat(rtl): direction-aware position() with isRTL helper"`

### Task 4: Menu keyboard flips (dropdown, context-menu, menubar)

**Files:**
- Modify: `ui/dropdown-menu.js` (~lines 153,163), `ui/context-menu.js` (~224,234), `ui/menubar.js` (~236,321,342)
- Test: `jstest/specs/dropdown.spec.ts`, `jstest/specs/menus-style-contract.spec.ts` or a new `jstest/specs/rtl.spec.ts` (create it here; Tasks 5–8 add to it)

**Interfaces:**
- Consumes: `isRTL` from `ui/gsxui.js`.
- Produces: `jstest/specs/rtl.spec.ts` exists with a `describe("rtl")` block and a shared `setRTL(page)` helper other tasks reuse:

```ts
export async function setRTL(page: Page) {
  await page.evaluate(() => document.documentElement.setAttribute("dir", "rtl"));
}
```

- [ ] **Step 1: Write failing test** in new `jstest/specs/rtl.spec.ts`: open a dropdown with a submenu, focus the sub-trigger, press `ArrowLeft` → submenu opens (RTL "into the menu" key); press `ArrowRight` → submenu closes, focus returns to sub-trigger.
- [ ] **Step 2: Run, confirm fails.**
- [ ] **Step 3: Implement**: in each menu module, at the keydown handlers, swap the meaning of ArrowLeft/ArrowRight when `isRTL(content)`:

```js
const openKey = isRTL(content) ? "ArrowLeft" : "ArrowRight";
const closeKey = isRTL(content) ? "ArrowRight" : "ArrowLeft";
if (e.key === openKey) { /* existing ArrowRight body */ }
if (e.key === closeKey /* + existing guards */) { /* existing ArrowLeft body */ }
```

Menubar additionally: top-level roving (`{ ArrowLeft: -1, ArrowRight: 1 }` at ~236) negates direction under RTL; the "switch to next/previous top-level menu" branches (~321/~342) swap likewise.

- [ ] **Step 4: Run rtl.spec.ts + existing dropdown/menubar specs — all pass.**
- [ ] **Step 5: Commit** `git commit -m "feat(rtl): mirror menu arrow-key semantics under rtl"`

### Task 5: Horizontal roving focus + carousel

**Files:**
- Modify: `ui/tabs.js` (~29), `ui/toggle-group.js` (same `{ArrowRight:1, ArrowLeft:-1}` pattern), `ui/carousel.js` (~197 keyboard; ~85-114 scroll math; prev/next buttons ~183-188)
- Test: `jstest/specs/rtl.spec.ts`

**Interfaces:**
- Consumes: `isRTL` from `ui/gsxui.js`, `setRTL` from rtl.spec.ts.

- [ ] **Step 1: Write failing tests**: (a) tabs under RTL: focus first tab, press `ArrowLeft` → focus moves to the SECOND tab (visually left = forward in RTL); (b) carousel under RTL: click "next" → the second item's bounding box moves into the viewport from the LEFT (assert `scrollLeft` decreased or item 2 is visible).
- [ ] **Step 2: Run, confirm fail.**
- [ ] **Step 3: Implement**:
  - tabs/toggle-group/menubar roving: `let dir = {ArrowRight: 1, ArrowLeft: -1}[e.key]; if (dir && isRTL(list)) dir = -dir;`
  - carousel: horizontal only — in RTL, `viewport.scrollLeft` is 0-or-negative (CSSOM spec). Normalize when computing position and delta: `const pos = vertical ? viewport.scrollTop : Math.abs(viewport.scrollLeft);` and when scrolling, negate the horizontal delta under RTL: `viewport.scrollBy({ left: isRTL(viewport) ? -delta : delta, behavior: "smooth" })`. Keyboard `{ArrowLeft:-1, ArrowRight:1}` negates under RTL. Prev/next button handlers keep logical meaning (prev = previous item) — no change beyond the math above.
- [ ] **Step 4: Run rtl.spec.ts + existing carousel/tabs coverage — pass.**
- [ ] **Step 5: Commit** `git commit -m "feat(rtl): mirror roving focus and carousel scroll math"`

### Task 6: Slider fill + resizable drag

**Files:**
- Modify: `ui/slider.gsx` (gradient), `ui/resizable.js` (~220 keyboard + pointer delta)
- Test: `jstest/specs/rtl.spec.ts`; `ui/slider_test.go` expectations

**Interfaces:**
- Consumes: `isRTL`.

- [ ] **Step 1: Slider** — native `<input type=range>` flips value/keys natively under RTL; only the painted gradient is physical. In `ui/slider.gsx:82`, add `rtl:` arms next to both existing track-gradient classes:

```
rtl:[&::-webkit-slider-runnable-track]:bg-[linear-gradient(to_left,var(--primary)_0%,var(--primary)_var(--fill,0%),var(--muted)_var(--fill,0%),var(--muted)_100%)]
rtl:[&::-moz-range-track]:bg-[linear-gradient(to_left,...same...)]
```

Regenerate, `gsx fmt`, update `ui/slider_test.go`, `go test ./ui/...`.
- [ ] **Step 2: Write failing resizable test**: under RTL, dragging the handle 50px toward the viewport-left GROWS the first (inline-start, visually right) panel; `ArrowLeft` likewise grows it (native-direction convention: keys mirror).
- [ ] **Step 3: Implement resizable**: horizontal only — negate pointer `dx` and swap `ArrowLeft/ArrowRight` steps when `isRTL(group)`:

```js
const rtlFlip = !vertical && isRTL(group) ? -1 : 1;
// pointer: delta = (event.clientX - startX) * rtlFlip;
// keys: { ArrowLeft: -STEP * rtlFlip, ArrowRight: STEP * rtlFlip }
```

- [ ] **Step 4: Run rtl.spec.ts + existing resizable coverage — pass. Slider visual check via existing style specs.**
- [ ] **Step 5: Commit** `git commit -m "feat(rtl): mirror slider fill and resizable drag math"`

### Task 7: Sidebar, sheet, drawer logical sides

**Files:**
- Modify: `ui/sidebar.gsx` (side-keyed selectors at :119,:156 and the container inset/width rules), `ui/sidebar.js` (any left/right ternaries — grep `'left'`/`'right'`), `ui/sheet.gsx`/`ui/drawer.gsx` (verify Task 1's logical conversion covers them; fix leftovers)
- Test: `jstest/specs/rtl.spec.ts`, existing `jstest/specs/sidebar-page.spec.ts`, `jstest/specs/sheet.spec.ts`

**Interfaces:**
- Produces: `Sidebar`/`Sheet`/`Drawer` keep their `side="left|right"` API with PHYSICAL meaning (matching shadcn: `data-side` is a physical placement), while interior paddings/borders/text alignment are logical and the trigger icon flips. Document this in each component's header comment.

- [ ] **Step 1: Write failing test**: sidebar with default side under `dir="rtl"`: trigger icon is rotated (has `rtl:rotate-180` from Task 2); sheet (mobile sidebar) with `side="left"` under RTL still opens from the physical left, but its close button and header text sit at logical positions (close button on visual LEFT = `end-*` under RTL). Assert close-button x < content midpoint.
- [ ] **Step 2: Implement**: follow shadcn's sidebar RTL migration: replace any JS left/right ternaries in `ui/sidebar.js` with CSS keyed on `data-side` (the `.gsx` already uses `data-side` selectors — extend rather than rewrite); inside-the-rail spacing/borders from Task 1's deferred list (`sidebar.gsx:119,:156`) stay `data-side`-keyed and PHYSICAL (they describe which physical edge the rail borders), so the actual change set here is: verify + convert only content-level classes (menu badges/actions already done in Task 1), pass `dir` through to the mobile `SheetContent` if the sheet reads it, and confirm `SidebarTrigger`'s icon flip.
- [ ] **Step 3: Run rtl.spec.ts + sidebar/sheet specs — pass. `go test ./ui/...` after any `.gsx` change.**
- [ ] **Step 4: Commit** `git commit -m "feat(rtl): sidebar/sheet/drawer direction handling"`

### Task 8: input-otp stays LTR

**Files:**
- Modify: `ui/input-otp.gsx` (slot group element)
- Test: `ui/input-otp_test.go` (if present; else add assertion to existing test file), `jstest/specs/rtl.spec.ts`

**Interfaces:** none downstream.

- [ ] **Step 1: Check shadcn's input-otp RTL behavior** (fetch their input-otp source/docs). Expected: digit groups render LTR even in RTL documents (numeric-code convention).
- [ ] **Step 2: Write failing Go test**: rendered group markup contains `dir="ltr"` on the slot-group container.
- [ ] **Step 3: Implement**: add `dir="ltr"` to the group element in `ui/input-otp.gsx`. Note: Task 1 converted its `border-r/rounded-l/r` to logical — with `dir="ltr"` forced, logical == physical, so no visual change in either context. Regenerate, `gsx fmt`, `go test ./ui/...`.
- [ ] **Step 4: Add rtl.spec.ts case**: OTP inside a `dir="rtl"` page — first slot's bounding box is the LEFTMOST. Run — pass.
- [ ] **Step 5: Commit** `git commit -m "feat(rtl): keep input-otp digit order ltr"`

### Task 9: Docs — RTL page, Arabic login demo, per-component RTL examples

**Files:**
- Create: `site/pages/rtl.gsx` (+ generated `.x.go`) — follow `site/pages/getting_started.gsx` structure for a prose page; wire into nav/router wherever getting_started is registered (grep for `getting_started` or `GettingStarted` in `site/`)
- Create: `site/examples/rtl/login.gsx` — Arabic login card (Card + Label + Input + Button, `dir="rtl"`, Arabic strings: title `تسجيل الدخول`, email label `البريد الإلكتروني`, password `كلمة المرور`, submit `دخول`, hint `ليس لديك حساب؟ إنشاء حساب`)
- Create: `site/examples/calendar/rtl.gsx`, `site/examples/pagination/rtl.gsx`, `site/examples/sidebar/rtl.gsx` (or the sidebar example dir's naming convention — check `ls site/examples/sidebar*`)
- Modify: `site/examples/calendar.go`, `site/examples/pagination.go`, `site/examples/sidebar.go` — `Register(...)` a new `Example{Name: "rtl", Title: "RTL", Node: ..., SourcePath: ".../rtl.gsx"}`; create `site/examples/rtl.go` registering the login demo under the rtl page's component key.

**Interfaces:**
- Consumes: components as fixed in Tasks 1–8.
- Produces: docs pages Task 10's specs can screenshot/assert against.

- [ ] **Step 1: Study one existing example end-to-end** (`site/examples/pagination/basic.gsx` + `pagination.go` + how `site/pages/component.gsx` renders registered examples) before writing anything.
- [ ] **Step 2: Write the RTL examples**: each wraps its demo in a `dir="rtl"` container with Arabic content — pagination: page numbers with Arabic-Indic numerals optional but keep Latin digits (shadcn does); calendar: default month grid under `dir="rtl"` (week start/nav flip); sidebar: `side="right"` + `dir="rtl"` variant matching shadcn's sidebar RTL section.
- [ ] **Step 3: Write `site/pages/rtl.gsx`**: sections — Get started (set `dir="rtl"` on `<html>` or any subtree; components adapt automatically), How it works (logical properties, `rtl:rotate-180` icons, direction-aware positioning JS), Try it out (embed the login demo), Fonts (recommend Noto Sans Arabic / Noto Naskh, pairs with the site's Geist).
- [ ] **Step 4: Build and eyeball**: run the site dev server, load `/rtl` and the three component pages, verify layout in both themes. Fix issues found.
- [ ] **Step 5: `go test ./site/...`, `gsx fmt` idempotency.**
- [ ] **Step 6: Commit** `git commit -m "docs(rtl): rtl guide page and rtl examples for calendar, pagination, sidebar"`

### Task 10: Full-suite verification gates

**Files:**
- Test: everything; no new source files (fix regressions where found)

- [ ] **Step 1: `go build ./... && go test ./...`** — green.
- [ ] **Step 2: `gsx fmt` over the repo, confirm zero diff (idempotency gate).**
- [ ] **Step 3: Full Playwright suite** in `jstest/` including `rtl.spec.ts` — green; update style-visual snapshots ONLY where a diff is an intended RTL-neutral change (logical classes must not shift LTR rendering — a snapshot diff in LTR means a Task 1 conversion error; fix the class, not the snapshot).
- [ ] **Step 4: Manual smoke** via dev server: dropdown/select/tooltip alignment, sheet sides, slider fill, carousel, sidebar — under both `dir` values.
- [ ] **Step 5: Commit any fixes; final commit** `git commit -m "test(rtl): full verification pass"`
