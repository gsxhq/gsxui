# Init Conventions Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Uniform init conventions plus four shared helpers across ui/*.js — pure refactor, zero behavior change.

**Architecture:** `ui/gsxui.js` gains `uid`, `wireGroupLabels`, `once`, `hasLiveTabStop`. All 12 modules registering with core `init()` converge on a named `initRoot` passed directly, living in a uniform `// --- init ---` section at the end of each file. toaster.js's unrelated private `init()` is renamed `bootstrap()`.

**Tech Stack:** ui/*.js (vendored library, dependency-free ESM), Playwright jstest.

## Global Constraints

- ZERO behavior change. The full Playwright suite is the regression net; any spec that changes (other than the one new uid assertion) is a red flag, not a fixture update.
- No htmx references in ui/.
- Registered entry point is always `function initRoot(el)`; internal helpers keep their domain names.
- Every registering file ends with a `// --- init ---...` section (match each file's existing horizontal-rule comment style) containing `initRoot` + the single `init(selector, initRoot)` call. avatar's window-load fallback also lives there.
- menubar/toggle-group keep SEMANTIC bail via the shared `hasLiveTabStop` (NOT `once()` — morph-stripped tabindex must re-normalize). carousel's autoplay uses `once()`.
- Playwright: `npx playwright test --config jstest/playwright.config.ts <files>`; known flakes dialog.spec.ts:453/503, rtl.spec.ts:15, home-showcase.spec.ts:16 — retry in isolation.
- Commit trailer:
  `Claude-Session: https://claude.ai/code/session_01GJjsg9jx6Cxp5B87A1RnH3`

---

### Task 1: Core helpers + contract comment + uid spec

**Files:**
- Modify: `ui/gsxui.js`
- Test: `jstest/specs/dynamic-init.spec.ts` (one new test)

**Interfaces:**
- Produces (all exported): `uid(prefix) -> string`; `wireGroupLabels(root, groupSelector, labelSelector, idPrefix) -> void`; `once(fn) -> (el) => void`; `hasLiveTabStop(els) -> boolean`.

- [ ] **Step 1: Write the failing uid test**

Append to jstest/specs/dynamic-init.spec.ts's describe block:

```ts
  test("uid never repeats across prefixes", async ({ page }) => {
    await page.goto("/x/select/basic");
    const ids = await page.evaluate(async () => {
      const { uid } = await import("/ui/gsxui.js");
      return [uid("gsxui-a"), uid("gsxui-b"), uid("gsxui-a"), uid("gsxui-b")];
    });
    expect(new Set(ids).size).toBe(4);
    expect(ids[0]).not.toBe(ids[2]);
  });
```

- [ ] **Step 2: Run to verify it fails** (`uid` not exported).

Run: `npx playwright test --config jstest/playwright.config.ts jstest/specs/dynamic-init.spec.ts`

- [ ] **Step 3: Implement the helpers in ui/gsxui.js**

After the init section, matching the file's comment voice:

```js
// --- shared init helpers --------------------------------------------------

// One counter for every generated id in the library. Uniqueness within the
// document is the entire contract — which module the number came from never
// matters, so modules share it instead of each keeping a private counter.
let uidCounter = 0;
export function uid(prefix) {
  return `${prefix}-${++uidCounter}`;
}

// Group→label aria wiring shared by listbox-shaped components (select,
// combobox): every group under root that lacks aria-labelledby gets pointed
// at its own label, generating the label's id on demand. Skip-if-present
// per group keeps it idempotent under re-init.
export function wireGroupLabels(root, groupSelector, labelSelector, idPrefix) {
  for (const group of root.querySelectorAll(groupSelector)) {
    if (group.getAttribute("aria-labelledby")) continue;
    const label = group.querySelector(labelSelector);
    if (!label) continue;
    if (!label.id) label.id = uid(idPrefix);
    group.setAttribute("aria-labelledby", label.id);
  }
}

// Initializers re-run on ANY mutation under a match (see init() above), so
// they split into two kinds of work: REFLECTION (recompute state from the
// DOM — must re-run unguarded, that is the whole point of re-init) and
// RESOURCE BINDS (timers, listeners, per-element observers — must happen
// once per element or they stack). once(fn) is the guard for the second
// kind: the returned function runs fn at most once per element, forever.
// It is NOT for semantic "a live value already exists" bails — those must
// re-arm after a morph strips the value (see hasLiveTabStop below).
export function once(fn) {
  const ran = new WeakSet();
  return (el) => {
    if (ran.has(el)) return;
    ran.add(el);
    fn(el);
  };
}

// Roving-tabindex bail shared by menubar and toggle-group: "has normalize
// already picked an entry tab stop, or has the user moved it?" Checks the
// tabindex ATTRIBUTE, never the IDL property — a plain <button> reports
// .tabIndex === 0 even with no attribute at all, which would make this
// always-true on server-fresh DOM and prevent the initial roving collapse.
// Server markup renders no tabindex attributes; interaction and normalize
// write them; a morph back to server state strips them — exactly the reset
// that must re-trigger normalize, which is why this is a semantic check
// and not a once() guard.
export function hasLiveTabStop(els) {
  return [...els].some((el) => el.getAttribute("tabindex") === "0");
}
```

Also extend the init() contract comment with two sentences: the registered
function is by convention named `initRoot` and lives in each module's
trailing `--- init ---` section; the shared observer is instantiated by
whichever barrel import calls init() first (avatar today).

- [ ] **Step 4: Run the spec to verify it passes**, plus `jstest/specs/smoke.spec.ts`.

- [ ] **Step 5: Commit**

```bash
git add ui/gsxui.js jstest/specs/dynamic-init.spec.ts
git commit -m "refactor(ui): shared init helpers — uid, wireGroupLabels, once, hasLiveTabStop"
```

---

### Task 2: Converge the listbox pair + command (uid/wireGroupLabels consumers)

**Files:**
- Modify: `ui/select.js`, `ui/combobox.js`, `ui/command.js`

**Interfaces:**
- Consumes: `uid`, `wireGroupLabels` from Task 1.

Per module (read each file first; keep all behavior identical):

- **select.js**: delete the private `let uid = 0`; replace the group/label wiring loop inside `initRoot` with `wireGroupLabels(root, "[data-gsxui-slot-select-group]", "[data-gsxui-slot-select-label]", "gsxui-select-label")`; replace remaining `` `gsxui-select-content-${++uid}` `` uses with `uid("gsxui-select-content")`; move `initRoot` + `init(...)` into a trailing `// --- init ---` section (the rest of the file's sections keep their order). Careful: the toggle handler's `content.id` assignment uses the counter — convert it too.
- **combobox.js**: same treatment — private counter deleted, wiring loop → `wireGroupLabels(root, "[data-gsxui-slot-combobox-group]", "[data-gsxui-slot-combobox-label]", "gsxui-combobox-label")` (verify the exact slot selectors in the file; use what the current loop uses), other `${++uid}` sites → `uid("gsxui-combobox-list")` / `uid("gsxui-combobox-item")` as they are today; `initRoot` + registration to the trailing section.
- **command.js**: private counter → `uid("gsxui-command-item")`; replace `init("[data-gsxui-slot-command]", (root) => filter(root))` with a named `function initRoot(root) { filter(root); }` + `init("[data-gsxui-slot-command]", initRoot)` in a trailing section, with a one-line comment that initRoot exists to carry the convention while `filter` stays the domain name.

- [ ] **Step 1: Convert select.js; run `npx playwright test --config jstest/playwright.config.ts jstest/specs/select.spec.ts jstest/specs/site-boost.spec.ts jstest/specs/dynamic-init.spec.ts`** → PASS
- [ ] **Step 2: Convert combobox.js; run the combobox-covering specs (grep jstest/specs for combobox; composites-style-contract + smoke at minimum)** → PASS
- [ ] **Step 3: Convert command.js; run `jstest/specs/command.spec.ts`** → PASS
- [ ] **Step 4: Commit**

```bash
git add ui/select.js ui/combobox.js ui/command.js
git commit -m "refactor(ui): select/combobox/command adopt shared uid + wireGroupLabels"
```

---

### Task 3: Converge the remaining nine + toaster rename

**Files:**
- Modify: `ui/avatar.js`, `ui/calendar.js`, `ui/carousel.js`, `ui/input-otp.js`, `ui/menubar.js`, `ui/toggle-group.js`, `ui/resizable.js`, `ui/toaster.js`

Per module:

- **avatar.js**: `sweep` stays as the per-image helper; add `function initRoot(img) { sweep(img); }`? NO — avatar registers per-image already with `sweep` doing exactly initRoot's job: RENAME `sweep` to `initRoot` (update the window-load fallback loop to call `initRoot`), trailing init section holds initRoot + `init(...)` + the window-load fallback with its existing comment.
- **calendar.js**: name the inline arrow: `function initRoot(root) { captureDefaults(root); const { year, month } = currentMonth(root); reconcileToday(root, year, month); }`, registered in the trailing section. Keep captureDefaults's WeakMap first-call-wins guard untouched.
- **carousel.js**: rename `initCarousel` → `initRoot`; replace the `autoplayBound` WeakSet inside `initAutoplay` with module-level `const bindAutoplay = once(initAutoplayInner)` — concretely: rename current `initAutoplay` internals so the once() wrapper owns the guard (delete the WeakSet + its guard lines; wrap: `const initAutoplay = once(function bindAutoplay(root) { ...existing body minus guard... });`). Keep the per-step idempotency audit comment, updated to point at once().
- **input-otp.js**: `function initRoot(root) { stamp(root); recompute(root); }` in trailing section.
- **menubar.js**: rename `normalize` → keep `normalize` as domain helper BUT registered entry becomes `function initRoot(bar) { normalize(bar); }`? NO — simpler: rename `normalize` to `initRoot` directly (it has no other callers — verify with grep; if interaction handlers also call `normalize`, KEEP `normalize` as the domain name and add `initRoot` calling it — read the file and report which case applied). Replace the guard with `if (hasLiveTabStop(enabledTriggersOf(bar))) return;` and delete the duplicated attribute-vs-IDL comment (point at gsxui.js's helper comment instead).
- **toggle-group.js**: same treatment as menubar (same verify-callers caveat).
- **resizable.js**: `function initRoot(handle) { handle.style.touchAction = "none"; syncHandleAria(handle); }` in trailing section, keeping the touch-action comment.
- **toaster.js**: rename private `init` → `bootstrap` (definition + both readyState call sites), adding: `// bootstrap, not the core init(): toast cards are cloned from <template>s, not server DOM matching a stable selector, so the selector-based self-healing model does not fit.`

- [ ] **Step 1: Convert module-by-module, running each one's covering spec after** (calendar → calendar.spec.ts; carousel → dynamic-init.spec.ts carousel test + any carousel spec; menubar/toggle-group → roving-tabindex.spec.ts; toaster → toaster.spec.ts; avatar/input-otp/resizable → smoke.spec.ts + the two style-contract sweeps)
- [ ] **Step 2: Full suite** `npx playwright test --config jstest/playwright.config.ts` → green (flakes in isolation)
- [ ] **Step 3: Grep gates**: `grep -rn "\\+\\+uid" ui/` → nothing; `grep -rn "function sweep\|function initCarousel\|autoplayBound" ui/` → nothing; `grep -c "function initRoot" ui/*.js` → one per registering module; `grep -rn -i htmx ui/` → only the pre-existing toaster/calendar mentions.
- [ ] **Step 4: Commit**

```bash
git add ui/avatar.js ui/calendar.js ui/carousel.js ui/input-otp.js ui/menubar.js ui/toggle-group.js ui/resizable.js ui/toaster.js
git commit -m "refactor(ui): uniform initRoot convention across behavior modules"
```

---

### Task 4: Gates + finish

- [ ] **Step 1:** `go build ./... && go test ./... -count=1` green; `make audit`; gofmt clean; full Playwright green.
- [ ] **Step 2:** Controller: final review, then superpowers:finishing-a-development-branch (branch stacks on htmx4-boost / PR #20).
