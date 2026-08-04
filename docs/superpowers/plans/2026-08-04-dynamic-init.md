# Self-Healing Component Init Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Library-level `init(selector, fn)` in the delegation core so component init work re-runs for dynamically added or morphed DOM — fixing state desync (e.g. select's checked-item reflection) under htmx swaps, boosted morph navigation, and raw innerHTML injection.

**Architecture:** `ui/gsxui.js` gains an init registry plus ONE shared MutationObserver (childList + subtree + attributes + characterData) that schedules microtask-debounced, idempotent re-inits of registered roots; mutations produced during an init pass are discarded (loop prevention). Ten modules convert their one-time load scans to `init(...)`. No htmx-specific code; the site's boost is just the first consumer.

**Tech Stack:** ui/gsxui.js (vendored library core), ui/*.js behavior modules, Playwright jstest harness.

## Global Constraints

- Library code (`ui/`) must stay dependency-free native ESM; no htmx references anywhere in `ui/`.
- The core deliberately allows RE-init of the same element (morph-reset roots need it). True once-per-element semantics are each stateful component's own job (WeakSet/data-flag inside its initializer).
- The observer must not start unless `init()` is registered at least once.
- Loop prevention is mandatory: an init pass's own DOM writes must not re-trigger the observer (re-entrancy flag + `takeRecords()` discard after the pass). This tradeoff is documented in the spec.
- Comments in several modules call late-added elements an "accepted limitation" (carousel.js:236-239, resizable.js:253-255, input-otp.js:152-155, command.js:224-225) — those comments must be updated by the conversion, not left stale.
- Playwright: `npx playwright test --config jstest/playwright.config.ts <files>`. Known parallel-load flakes: dialog.spec.ts:453/503 and rtl.spec.ts:15 — retry in isolation before blaming a change.
- After editing any `.gsx` run `go tool gsx generate`; gsx markup comments are `{/* */}`.
- Commit trailer:
  `Claude-Session: https://claude.ai/code/session_01GJjsg9jx6Cxp5B87A1RnH3`

---

### Task 1: Core `init()` + shared MutationObserver + harness proof

**Files:**
- Modify: `ui/gsxui.js` (new section after the `on`/`dispatch` block)
- Test: `jstest/specs/dynamic-init.spec.ts` (new)

**Interfaces:**
- Produces: `export function init(selector, fn)` in ui/gsxui.js. Semantics: runs `fn(el)` now for every `document.querySelectorAll(selector)` match; thereafter (a) added element nodes matching selector (self or descendant) get `fn`, (b) attribute/characterData/childList mutations at-or-under an existing match schedule that root's re-init; scheduling is microtask-debounced per (element, fn); mutations during the flush are discarded.
- Consumes: nothing.

- [ ] **Step 1: Write the failing Playwright spec**

`jstest/specs/dynamic-init.spec.ts` — drive the core directly through the browser on any harness page (the harness serves real modules; `/x/select/basic` renders a select example):

```ts
import { expect, test } from "@playwright/test";

// The library promise under test: init(selector, fn) runs for current
// matches, for injected matches, and re-runs when a match's subtree is
// mutated back to server state (the morph case) — with init-pass writes
// not re-triggering the observer.
test.describe("gsxui dynamic init", () => {
  test("init runs now, on injection, and on mutation — without looping", async ({ page }) => {
    await page.goto("/x/select/basic");
    const counts = await page.evaluate(async () => {
      const { init } = await import("/ui/gsxui.js");
      const calls: string[] = [];
      // Idempotent marker initializer: writes DOM (data-inited) — the
      // observer must NOT re-trigger on that own write.
      init("[data-di-probe]", (el) => {
        calls.push(el.id);
        el.dataset.inited = "true";
      });
      // (a) current match
      const a = document.createElement("div");
      a.id = "pre";
      a.setAttribute("data-di-probe", "");
      document.body.append(a);
      await new Promise((r) => setTimeout(r, 20)); // observer flush
      // (b) injected subtree (descendant match)
      const wrap = document.createElement("div");
      wrap.innerHTML = '<section><div id="inj" data-di-probe></div></section>';
      document.body.append(wrap);
      await new Promise((r) => setTimeout(r, 20));
      // (c) mutation under an existing match re-inits it
      document.getElementById("inj")!.textContent = "server reset";
      await new Promise((r) => setTimeout(r, 20));
      return calls;
    });
    // "pre" was appended before registration? No — registration ran first,
    // so: no initial matches, then pre (injection), inj (injection),
    // inj again (mutation). The own-write (data-inited) must not add more.
    expect(counts).toEqual(["pre", "inj", "inj"]);
  });

  test("initial matches run at registration", async ({ page }) => {
    await page.goto("/x/select/basic");
    const ran = await page.evaluate(async () => {
      const { init } = await import("/ui/gsxui.js");
      let n = 0;
      init("[data-gsxui-slot-select]", () => n++);
      return n;
    });
    expect(ran).toBeGreaterThan(0); // the page's select root
  });
});
```

Adjust the import URL to how the harness serves the module (`/ui/gsxui.js` — same URL the shell's `<script type="module" src="/ui/index.js">` graph uses). Record any locator/URL adjustment in your report.

- [ ] **Step 2: Run to verify it fails**

Run: `npx playwright test --config jstest/playwright.config.ts jstest/specs/dynamic-init.spec.ts`
Expected: FAIL — `init` is not exported.

- [ ] **Step 3: Implement in ui/gsxui.js**

Add after the `dispatch` function, matching the file's comment style:

```js
// --- dynamic init ---------------------------------------------------------
//
// on() needs no re-scan for late DOM — document-level delegation already
// covers it. Init WORK (reflecting server-rendered state, per-instance
// wiring) does not get that for free: a module's load-time scan runs once.
// init(selector, fn) closes the gap: fn runs for current matches at
// registration, for matches added later (htmx swaps, morphs, innerHTML),
// and again when a match's subtree is morphed back to server state.
//
// Contract: fn is IDEMPOTENT. The core deliberately re-runs it on mutated
// roots — a component needing true once-per-element wiring guards inside
// its own fn (WeakSet / data flag). Mutations caused by an init pass
// itself are discarded (flushing + takeRecords below); a concurrent
// external mutation landing in that exact window is missed until the next
// mutation re-heals it.

const inits = []; // [{ selector, fn, scheduled: WeakSet-per-flush }]
let observer = null;
let pending = null; // Map<fn, Set<Element>> while a flush is queued
let flushing = false;

export function init(selector, fn) {
  inits.push({ selector, fn });
  for (const el of document.querySelectorAll(selector)) fn(el);
  if (!observer) {
    observer = new MutationObserver(collect);
    observer.observe(document.documentElement, {
      childList: true,
      subtree: true,
      attributes: true,
      characterData: true,
    });
  }
}

function schedule(fn, el) {
  if (!pending) {
    pending = new Map();
    queueMicrotask(flush);
  }
  let set = pending.get(fn);
  if (!set) pending.set(fn, (set = new Set()));
  set.add(el);
}

function collect(records) {
  if (flushing) return;
  for (const record of records) {
    if (record.type === "childList") {
      for (const node of record.addedNodes) {
        if (!(node instanceof Element)) continue;
        for (const { selector, fn } of inits) {
          if (node.matches(selector)) schedule(fn, node);
          for (const el of node.querySelectorAll(selector)) schedule(fn, el);
        }
      }
    }
    const target =
      record.target instanceof Element
        ? record.target
        : record.target.parentElement;
    if (!target) continue;
    for (const { selector, fn } of inits) {
      const root = target.closest(selector);
      if (root) schedule(fn, root);
    }
  }
}

function flush() {
  const batch = pending;
  pending = null;
  if (!batch) return;
  flushing = true;
  try {
    for (const [fn, els] of batch) {
      for (const el of els) if (el.isConnected) fn(el);
    }
  } finally {
    observer.takeRecords(); // discard mutations our own inits produced
    flushing = false;
  }
}
```

Note the subtlety the implementation must preserve: `flushing` alone does not stop record delivery (observer callbacks fire on a microtask AFTER ours), so the `takeRecords()` drain inside the flush is the actual discard; the `flushing` guard covers synchronous delivery edge cases. If during implementation you find a case where the callback still fires with stale records after the drain, say so in your report rather than papering over it.

- [ ] **Step 4: Run the spec to verify it passes**

Run: `npx playwright test --config jstest/playwright.config.ts jstest/specs/dynamic-init.spec.ts`
Expected: PASS (both tests).

- [ ] **Step 5: Sanity sweep of untouched behavior**

Run: `npx playwright test --config jstest/playwright.config.ts jstest/specs/smoke.spec.ts jstest/specs/select.spec.ts`
Expected: PASS — registering nothing changes nothing.

- [ ] **Step 6: Commit**

```bash
git add ui/gsxui.js jstest/specs/dynamic-init.spec.ts
git commit -m "feat(ui): init(selector, fn) — self-healing init for dynamic DOM"
```

---

### Task 2: Convert select.js + boosted-navigation regression specs

**Files:**
- Modify: `ui/select.js:182` (the `for (const root of document.querySelectorAll("[data-gsxui-slot-select]")) init(root)` scan — note the local function is already named `init`; rename the local to `initRoot` to import the core's `init` cleanly)
- Test: `jstest/specs/site-boost.spec.ts` (extend)

**Interfaces:**
- Consumes: `init(selector, fn)` from Task 1 (`import { on, init, ... } from "./gsxui.js"` — extend select.js's existing import list).
- Produces: select reflection self-heals; spec names `"select reflects server-checked item after boosted navigation"` and `"…after same-page boosted navigation"`.

- [ ] **Step 1: Extend site-boost.spec.ts (failing first)**

```ts
  test("select reflects server-checked item after boosted navigation", async ({ page }) => {
    await page.goto("/site/components/button");
    await page.getByRole("link", { name: "select", exact: true }).last().click();
    await expect(page).toHaveURL(/\/components\/select$/);
    await expect(
      page.locator("[data-gsxui-slot-select-trigger]").first(),
    ).toContainText("Apple");
  });

  test("select reflects server-checked item after same-page boosted navigation", async ({ page }) => {
    await page.goto("/site/components/select");
    await expect(page.locator("[data-gsxui-slot-select-trigger]").first()).toContainText("Apple");
    await page.getByRole("link", { name: "select", exact: true }).last().click();
    await expect(page.locator("[data-gsxui-slot-select-trigger]").first()).toContainText("Apple");
  });
```

Adjust URL prefixes/locators to how site-boost.spec.ts already addresses site pages (it uses `page.goto("/components/button")` — mirror that; the sidebar link's accessible name is lowercase). `.last()` picks the visible aside link over the hidden mobile-popover duplicate — verify against the DOM and note what you used.

- [ ] **Step 2: Run to verify the new specs fail**

Run: `npx playwright test --config jstest/playwright.config.ts jstest/specs/site-boost.spec.ts`
Expected: the two new tests FAIL (trigger shows "Select a fruit"); existing four still pass.

- [ ] **Step 3: Convert select.js**

- Rename select.js's local `function init(root)` → `function initRoot(root)` (one call site in the scan).
- Add `init` to the `./gsxui.js` import.
- Replace the scan (`for (const root of document.querySelectorAll("[data-gsxui-slot-select]")) initRoot(root);`) with `init("[data-gsxui-slot-select]", initRoot);`
- Audit initRoot for idempotency (it reflects checked state via applyValue with `{silent: true}` — re-running is safe; confirm and say so in the report).

- [ ] **Step 4: Run to verify green**

Run: `npx playwright test --config jstest/playwright.config.ts jstest/specs/site-boost.spec.ts jstest/specs/select.spec.ts jstest/specs/dynamic-init.spec.ts`
Expected: all PASS.

- [ ] **Step 5: Commit**

```bash
git add ui/select.js jstest/specs/site-boost.spec.ts
git commit -m "fix(ui): select reflection self-heals under dynamic DOM"
```

---

### Task 3: Convert the remaining scan modules

**Files:**
- Modify: `ui/avatar.js:15-25`, `ui/toggle-group.js:107`, `ui/command.js:222-226`, `ui/menubar.js:116-117`, `ui/combobox.js:296-297`, `ui/input-otp.js:152-159`, `ui/calendar.js:982,1064-1067`, `ui/resizable.js:252-262`, `ui/carousel.js:235-250`
- Test: existing jstest specs (calendar, command, carousel via style/contract specs, toaster etc.) + `jstest/specs/dynamic-init.spec.ts` idempotency addition for carousel

Conversion pattern for each (import `init` from `./gsxui.js`, replace the loop):

| module | selector | initializer | idempotency note |
|---|---|---|---|
| avatar.js | `[data-gsxui-slot-avatar-image]` | wrap existing `sweep` body per-element: `(img) => { if (img.complete) sync(img, img.naturalWidth > 0); }` — keep the `window load` sweep | pure reflection, safe |
| toggle-group.js | `[data-gsxui-slot-toggle-group]` | `normalize` | reflection, safe |
| command.js | `[data-gsxui-slot-command]` | `filter` | reflection, safe |
| menubar.js | `[data-gsxui-slot-menubar]` | `normalize` | sets first active trigger; re-run resets roving focus to first trigger on morph — correct (server reset) |
| combobox.js | `[data-gsxui-slot-combobox]` | local `init` → rename `initRoot` | reflection, safe |
| input-otp.js | `[data-gsxui-slot-input-otp]` | `(root) => { stamp(root); recompute(root); }` | reflection, safe |
| calendar.js | `ROOT` | `(root) => { captureDefaults(root); const { year, month } = currentMonth(root); reconcileToday(root, year, month); }` (merge both loops into one initializer) | captureDefaults re-run on morph re-captures server state — correct; confirm captureDefaults has no once-only assumption by reading it |
| resizable.js | `HANDLE_SELECTOR` | `(handle) => { handle.style.touchAction = "none"; }` (keep the comment about touch-action) | style write, idempotent |
| carousel.js | `[data-gsxui-slot-carousel]` | existing block wrapped as `initCarousel(root)` | STATEFUL: `resizeObserver.observe(viewport)` double-observe is a no-op (same target ⇒ one entry — verify against MDN and say so), `root.gsxuiCarousel` reassign safe, `recompute` reflection; `initAutoplay` MUST be audited — if it binds timers/listeners per call, add an internal `if (root.gsxuiAutoplayBound) return;`-style guard (read the function; report what you found) |

Also update the now-false "accepted limitation" comments in carousel.js, resizable.js, input-otp.js, command.js (and toggle-group.js if it carries one) to describe the new self-healing behavior.

- [ ] **Step 1: Add the carousel idempotency spec (failing or green — it pins behavior)**

Append to `jstest/specs/dynamic-init.spec.ts`:

```ts
  test("carousel init is idempotent under re-init", async ({ page }) => {
    await page.goto("/x/carousel/basic");
    await page.evaluate(async () => {
      const { init } = await import("/ui/gsxui.js");
      // Force a re-init pass over existing carousels by touching a root.
      document.querySelector("[data-gsxui-slot-carousel]")!.setAttribute("data-poke", "1");
      await new Promise((r) => setTimeout(r, 20));
    });
    const before = await page
      .locator("[data-gsxui-slot-carousel-item]")
      .first()
      .boundingBox();
    await page.locator("[data-gsxui-slot-carousel-next]").first().click();
    await page.waitForTimeout(400); // one smooth-scroll settle
    const after = await page
      .locator("[data-gsxui-slot-carousel-item]")
      .first()
      .boundingBox();
    expect(after!.x).toBeLessThan(before!.x); // advanced exactly one step is
    // asserted by carousel's own existing specs; here we assert it MOVED
    // (a double-bound autoplay/observer would break those existing specs).
  });
```

Verify the example URL and slot names against the carousel example fixtures (`site/examples/carousel/`), adjust, and note adjustments. If carousel has no `/x/carousel/basic` route, use the manifest (`jstest/.tmp/examples.json` at runtime) pattern other specs use — copy their approach.

- [ ] **Step 2: Convert all nine modules per the table**

One module at a time; after each, run its closest covering spec (calendar → calendar.spec.ts, command → command.spec.ts, carousel → the carousel-covering spec, others → smoke.spec.ts + the two style-contract specs that sweep selectors).

- [ ] **Step 3: Full jstest suite**

Run: `npx playwright test --config jstest/playwright.config.ts`
Expected: all pass (known flakes in isolation).

- [ ] **Step 4: Commit**

```bash
git add ui/avatar.js ui/toggle-group.js ui/command.js ui/menubar.js ui/combobox.js ui/input-otp.js ui/calendar.js ui/resizable.js ui/carousel.js jstest/specs/dynamic-init.spec.ts
git commit -m "feat(ui): convert load-time init scans to self-healing init()"
```

---

### Task 4: Gates + real-browser verification

- [ ] **Step 1: Gates**

```bash
go build ./... && gofmt -l . | grep -v node_modules
go test ./... -count=1 > tmp/gates.out 2>&1; echo "exit:$?"; grep -E "^--- FAIL|^FAIL" tmp/gates.out
make audit && make verify-generated && make verify-generated-styles
npx playwright test --config jstest/playwright.config.ts
```

All green (flakes verified in isolation).

- [ ] **Step 2: Real-browser check (dev server)**

`npm run dev`, then in a browser: /components/button → sidebar-click select → trigger reads "Apple" immediately; open the popup — no flash, Apple ✓ checked; same-page sidebar re-click → still "Apple". Also spot-check a carousel and calendar page after boosted navigation. Report observations.

- [ ] **Step 3: Finish**

Controller: final whole-branch review (this branch also carries the htmx4-boost feature), then superpowers:finishing-a-development-branch.
