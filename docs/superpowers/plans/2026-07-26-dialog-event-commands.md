# Dialog Event Commands Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make gsxui dialog behavior request-event-first and DOM-state-only while preserving native `<dialog>` methods and exact-ID Invoker Commands as first-class escape hatches.

**Architecture:** Built-in gsxui controls dispatch bubbling, cancelable `gsxui:request-open` and `gsxui:request-close` events on the exact dialog they own. Document-delegated default handlers defer one microtask so any ancestor listener may cancel, then perform the native transition. `data-state` plus the native `open` attribute are the entire state model; `beforetoggle` stamps impending state and `toggle` reconciles confirmed state, ARIA, and existing `gsxui:open` / `gsxui:close` notifications. Native `commandfor`, `showModal()`, `show()`, `close()`, and `requestClose()` remain unwrapped platform APIs.

**Tech Stack:** Vanilla ES modules and CustomEvent, native HTMLDialogElement/Invoker Commands/ToggleEvent/Web Animations APIs, Playwright Chromium, gsx-generated Go examples.

## Global Constraints

- Write and run each browser regression before changing production behavior.
- Retain no dialog registry, counter, timer, `Map`, `WeakMap`, `Set`, `WeakSet`, or per-instance JavaScript object.
- Persist all identity and state in authored/generated DOM attributes.
- Proximity resolution must never cross a nested `[data-gsxui-dialog]` boundary.
- External declarative controls target an authored dialog `id` with native `commandfor`; generated ARIA IDs are not a public targeting API.
- There is no public gsxui imperative helper. Programmatic callers use the native dialog methods.
- Built-in close requests animate; native `close` and `command="close"` remain immediate.
- Only finite animations attached to the dialog element itself may delay close.
- Never hand-edit generated `.x.go` or highlighted example output.
- Browser tests assert behavior and DOM contracts, never utility classes,
  class strings, or computed visual styles. Use controlled Web Animations for
  transition behavior without pinning CSS implementation.
- `make check` must pass end to end before completion.

---

### Task 1: Request-event default actions and DOM-only dialog state

**Files:**

- Create: `jstest/specs/dialog.spec.ts`
- Modify: `ui/dialog.js`

**Interfaces:**

- Consumes: `on()` and `emit()` from `ui/gsxui.js`, existing `data-gsxui-dialog-*` hooks.
- Produces: `gsxui:request-open`, `gsxui:request-close`, reason detail, cancelable microtask-deferred defaults, finite own-animation close.

- [ ] **Step 1: Write the failing request-event regressions**

Use `/x/dialog/basic` for the real rendered component and dynamically append
extra dialog roots where a multi-instance or nested-root probe is required.
Add tests that:

1. attach a listener before clicking the proximity trigger and assert exactly
   one `gsxui:request-open` targeted at that dialog with
   `{ reason: "trigger" }`;
2. dispatch `gsxui:request-open` from a descendant, then
   `gsxui:request-close` from a descendant, and assert native `open` follows;
3. register a `document` listener after `ui/index.js` has loaded, call
   `preventDefault()`, and prove it cancels each default action;
4. click a `data-gsxui-dialog-close` control and assert
   `{ reason: "close-button" }`;
5. press Escape and click outside the content box, asserting reasons
   `cancel` and `backdrop`;
6. during close, assert `data-state="closed"` while native `open` is still
   present, then one `gsxui:close` notification after `open` disappears;
7. dispatch `gsxui:request-open` during that exit window and prove the
   pending close aborts.

Use a helper inside the spec, not production code:

```ts
async function dispatch(page, selector, type, detail = {}) {
  return page.locator(selector).evaluate(
    (element, { type, detail }) =>
      element.dispatchEvent(
        new CustomEvent(type, { bubbles: true, cancelable: true, detail }),
      ),
    { type, detail },
  );
}
```

- [ ] **Step 2: Verify RED**

Run:

```bash
npx playwright test --config jstest/playwright.config.ts jstest/specs/dialog.spec.ts
```

Expected: new tests fail because built-in controls do not emit request events,
request events have no default action, and the existing close uses a timer cap.

- [ ] **Step 3: Replace the counter with collision-checked DOM identity**

Implement random IDs without retained JavaScript state:

```js
function generatedId(prefix) {
  for (;;) {
    const words = crypto.getRandomValues(new Uint32Array(4));
    const suffix = [...words]
      .map((word) => word.toString(16).padStart(8, "0"))
      .join("");
    const id = `gsxui-${prefix}-${suffix}`;
    if (!document.getElementById(id)) return id;
  }
}

function ensureId(element, prefix) {
  if (!element.id) element.id = generatedId(prefix);
  return element.id;
}
```

Keep authored `id`, `aria-labelledby`, `aria-describedby`, and
`aria-controls` untouched. Restrict title, description, trigger, and content
queries to elements owned by the same nearest `[data-gsxui-dialog]` root.

- [ ] **Step 4: Add private request dispatch and default actions**

Built-in modules may use a private local dispatcher:

```js
function request(dialog, type, detail) {
  return emit(dialog, type, detail);
}
```

Register delegated handlers for both request events on
`dialog[data-gsxui-dialog-content]`. The handlers queue their default action:

```js
on("gsxui:request-open", DIALOG, (event, dialog) => {
  queueMicrotask(() => {
    if (!event.defaultPrevented) performOpen(dialog);
  });
});

on("gsxui:request-close", DIALOG, (event, dialog) => {
  queueMicrotask(() => {
    if (!event.defaultPrevented) void performClose(dialog);
  });
});
```

`performOpen` wires ARIA, stamps `data-state="open"`, and calls
`showModal()` only when native `open` is absent. That same state stamp cancels
an in-flight close.

`performClose` is idempotent when native `open` is absent or
`data-state="closed"` already marks an exit in progress. Stamp closed, then
wait only for finite own animations:

```js
const animations = dialog.getAnimations().filter((animation) => {
  const endTime = animation.effect?.getComputedTiming().endTime;
  return typeof endTime === "number" && Number.isFinite(endTime);
});
await Promise.allSettled(animations.map((animation) => animation.finished));
if (dialog.open && dialog.dataset.state === "closed") dialog.close();
```

Do not estimate duration and do not use a timeout.

- [ ] **Step 5: Route built-in controls through requests**

- proximity trigger → `gsxui:request-open`, reason `trigger`;
- close control → `gsxui:request-close`, reason `close-button`;
- valid light-dismiss click → `gsxui:request-close`, reason `backdrop`;
- native `cancel` → first `preventDefault()`, then
  `gsxui:request-close`, reason `cancel`.

Keep the existing static-dialog backdrop guard and coordinate check.

- [ ] **Step 6: Add identity and ownership regressions**

Extend the spec to render two copied roots plus a nested root. Assert:

- generated dialog/title/description IDs are unique and stable after repeated
  opens;
- every generated ARIA reference resolves to the intended element;
- authored IDs and authored ARIA relationships are unchanged;
- a proximity trigger opens only the dialog owned directly by its nearest
  root, never the nested root;
- `aria-expanded` changes only on triggers that own that dialog.

Add an infinite animation to a child and then to the dialog itself; neither
may prevent a close request from completing. Keep the finite exit-animation
assertion above so the test also proves gsxui does not skip real exits.

- [ ] **Step 7: Verify GREEN and commit**

Run:

```bash
npx playwright test --config jstest/playwright.config.ts jstest/specs/dialog.spec.ts
node --check ui/dialog.js
git diff --check
```

Expected: all focused dialog tests pass and JavaScript parses.

Commit:

```bash
git add ui/dialog.js jstest/specs/dialog.spec.ts
git commit -m "feat: drive dialogs through request events"
```

---

### Task 2: Native Invoker Commands and native dialog lifecycle

**Files:**

- Modify: `jstest/specs/dialog.spec.ts`
- Modify: `site/examples/dialog/events.gsx`
- Generated: `site/examples/dialog/events.x.go`
- Generated: `site/hl/blocks.gen.go`
- Modify: `ui/dialog.js`

**Interfaces:**

- Consumes: native `commandfor`, `command`, `beforetoggle`, `toggle`,
  `showModal()`, `show()`, `close()`, and `requestClose()`.
- Produces: exact-ID native control over many dialogs, before-paint state
  stamps, external invoker ARIA synchronization, unchanged post-transition
  notifications.

- [ ] **Step 1: Write failing native-path regressions**

Change the Events example to give its content the authored ID
`events-dialog` and demonstrate native invokers:

```gsx
<ui.Button commandfor="events-dialog" command="show-modal">Open</ui.Button>
...
<ui.Button
  variant="outline"
  commandfor="events-dialog"
  command="request-close"
>
  Close
</ui.Button>
```

In Playwright, append a second authored-ID dialog and assert:

1. each `commandfor="…" command="show-modal"` button opens only its exact
   target;
2. `command="request-close"` produces `gsxui:request-close` with reason
   `cancel` through the native cancel path, stamps `data-state="closed"`,
   then animates closed;
3. `command="close"` closes immediately without being converted into a
   gsxui request;
4. direct `showModal()` and `show()` expose `data-state="open"` no later than
   `beforetoggle`;
5. direct `close()` exposes `data-state="closed"` and emits exactly one
   `gsxui:close`;
6. direct `requestClose()` enters the same cancelable animated path as
   Escape;
7. request events remain cancelable without preventing these explicitly
   native immediate methods.

- [ ] **Step 2: Verify RED**

Run:

```bash
go tool gsx generate
npx playwright test --config jstest/playwright.config.ts jstest/specs/dialog.spec.ts
```

Expected: exact native targeting opens correctly in Chromium, while
before-paint state and external invoker ARIA assertions fail until lifecycle
handling is added. The generated example diff is expected and must remain
unstaged until reviewed.

- [ ] **Step 3: Add `beforetoggle` and reconcile `toggle`**

Capture-delegate `beforetoggle` on the dialog selector. For an impending open:

- resolve the owning root;
- wire accessibility relationships;
- stamp `data-state="open"` before paint.

For an impending close, stamp `data-state="closed"`. Keep `toggle` as the
confirmed-state source for:

- final `data-state`;
- proximity trigger `aria-expanded`;
- exact authored-ID `[commandfor][command="show-modal"]` invoker
  `aria-expanded`;
- `gsxui:open` / `gsxui:close`.

Use `CSS.escape(dialog.id)` when querying authored command targets. Do not
assign a generated ID merely to discover external controls; native external
control requires an authored ID.

- [ ] **Step 4: Generate, highlight, verify, and commit**

Run:

```bash
go tool gsx generate
make highlight
npx playwright test --config jstest/playwright.config.ts jstest/specs/dialog.spec.ts
go test ./site/examples/dialog ./site/hl -count=1
gopls check -severity=hint site/examples/dialog/events.x.go
git diff --check
```

Expected: all native lifecycle tests pass; generated Go and highlighted
source exactly match the authored example.

Commit:

```bash
git add ui/dialog.js jstest/specs/dialog.spec.ts \
  site/examples/dialog/events.gsx site/examples/dialog/events.x.go \
  site/hl/blocks.gen.go
git commit -m "feat: support native dialog invoker commands"
```

---

### Task 3: Sibling behavior migration, contract documentation, and full gate

**Files:**

- Modify: `ui/command.js`
- Modify: `jstest/specs/dialog.spec.ts`
- Modify: `docs/jsx-parity.md`

**Interfaces:**

- Consumes: `gsxui:request-open`, `gsxui:request-close`, existing command
  palette hotkey and item activation.
- Produces: no sibling direct state manipulation, accurate Phase 1 ledger,
  authoritative repository verification.

- [ ] **Step 1: Write failing command-palette request regressions**

On the command-dialog example, listen for request events and assert:

- Cmd/Ctrl-K while closed emits one `gsxui:request-open` on the command
  dialog and opens it;
- Cmd/Ctrl-K while open emits one `gsxui:request-close` and follows the
  animated close path;
- activating an item with `data-href` requests close before navigation.

Use reason strings `shortcut` and `select` for these sibling-originated
requests. The default dialog behavior must not branch on either value.

- [ ] **Step 2: Verify RED**

Run:

```bash
npx playwright test --config jstest/playwright.config.ts jstest/specs/dialog.spec.ts
```

Expected: command behavior still imports `requestClose` and manipulates
`data-state` / `showModal()` directly, so request-event assertions fail.

- [ ] **Step 3: Migrate `ui/command.js`**

Remove the `requestClose` import from `dialog.js`. Import `emit` only from
`gsxui.js` as today and dispatch requests on the exact command dialog:

```js
emit(dialog, dialog.open ? "gsxui:request-close" : "gsxui:request-open", {
  reason: "shortcut",
});
```

For a navigable selected item, emit `gsxui:request-close` with reason
`select` before navigation. Do not stamp dialog state or call native methods
from command.js.

- [ ] **Step 4: Reconcile the parity ledger**

Update the dialog and JavaScript sections of `docs/jsx-parity.md`:

- replace `requestClose` and 600 ms cap claims with finite own-animation
  completion;
- record the two request events, stable built-in reasons, and cancelation;
- record DOM-only state/identity and collision-checked random IDs;
- document proximity for local controls and authored-ID native `commandfor`
  for external/multi-dialog controls;
- document native methods as the programmatic API and explicitly state that
  gsxui adds no imperative helper.

Do not alter the deferred `DialogFooter` decision.

- [ ] **Step 5: Focused verification and commit**

Run:

```bash
npx playwright test --config jstest/playwright.config.ts jstest/specs/dialog.spec.ts
node --check ui/dialog.js
node --check ui/command.js
git diff --check
```

Expected: focused behavior passes and no stale `requestClose` import or timer
claim remains.

Commit:

```bash
git add ui/command.js jstest/specs/dialog.spec.ts docs/jsx-parity.md
git commit -m "docs: finalize dialog request contract"
```

- [ ] **Step 6: Run the authoritative repository gate**

Run after all generated changes are committed, because `make check` verifies
generated output against `HEAD`:

```bash
make check
```

Expected: uncached Go vet/tests, the full Playwright suite, generated `.x.go`
drift, JavaScript syntax, and formatting all pass.

- [ ] **Step 7: Independent adversarial review**

Have a fresh reviewer inspect the spec, plan, diff, and tests, and build
throwaway browser probes for:

- a late `window` listener canceling a request;
- two native commandfor triggers targeting two dialogs;
- an infinite child animation and an infinite dialog animation not blocking
  close;
- rapid reopen during a finite exit animation;
- HTMX-style dynamically inserted markup working without initialization.

Address every confirmed finding, rerun the focused test and `make check`, and
stop with a clean worktree.
