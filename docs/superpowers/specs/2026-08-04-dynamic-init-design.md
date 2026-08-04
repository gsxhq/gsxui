# Self-healing component init for dynamic DOM

Date: 2026-08-04
Status: approved (option A from the htmx4-boost debugging session)

## Problem

Each behavior module runs a load-time init scan
(`for (const root of document.querySelectorAll(...)) init(root)`) exactly
once, at module evaluation. The delegation core's guarantee — "elements
added later just work" — holds for event handling but not for init work.
Any DOM swapped in later (htmx morphs/swaps, innerHTML, our boosted site
navigation) gets event handling but no init.

Reproduced concretely: boosted-navigate to /components/select — the Basic
example's trigger shows the placeholder instead of the server-checked
"Apple"; the popup shows Apple ✓. ~10 components carry load-time scans:
avatar, calendar, carousel, combobox, command, input-otp, menubar,
resizable, select (see `grep "document.querySelectorAll" ui/*.js`).

This is a library defect, not a site one: htmx users are gsxui's core
audience and they swap DOM constantly. The fix lives in `ui/gsxui.js` and
the component modules — shipped to users via `gsxui add`.

## Design

### Core API (ui/gsxui.js)

```js
export function init(selector, fn)
```

- Registers `fn` as the initializer for elements matching `selector`, and
  immediately runs it for every current match (preserving today's
  load-time behavior).
- `fn` must be idempotent: running it twice on the same element must be
  harmless. Existing init functions already are (they reflect state), and
  the conversion audits each one.

### One shared MutationObserver

A single observer in gsxui.js watches `document.documentElement` with
`{ childList: true, subtree: true, attributes: true, characterData: true }`:

- **Added element nodes**: match the node itself and its descendants
  against every registered selector; schedule init for each match.
- **Attribute / characterData / childList mutations on existing nodes**:
  find the nearest ancestor (or self) matching a registered selector;
  schedule that root's re-init. This covers morphs that keep DOM nodes but
  reset their attributes or text to server state (e.g. boosted navigation
  to the same page, idiomorph node reuse).
- Scheduling is **microtask-debounced per (element, initializer)**: many
  mutation records for one morph collapse into one re-init pass.
- **Loop prevention**: mutations produced while the scheduled inits are
  running are discarded (`takeRecords()` after the pass, plus a
  re-entrancy flag). Init functions mutate the DOM they own; without this
  the observer would re-trigger on its own writes. Tradeoff: a genuine
  external mutation landing in exactly that window is missed — accepted
  and documented; the next mutation re-heals.
- The observer starts on first `init()` registration, so pages using only
  event behaviors pay nothing.

### Component conversion

Every load-time scan converts mechanically:

```js
// before
for (const root of document.querySelectorAll("[data-gsxui-slot-select]"))
  initRoot(root);
// after
init("[data-gsxui-slot-select]", initRoot);
```

One module at a time, auditing each initializer for idempotency (e.g.
carousel binds per-instance state — must guard against double-binding via
a WeakSet or data flag *inside its own initializer* where true
once-per-element semantics are needed; the core deliberately does not
impose once-per-element, because morph-reset roots need re-init).

### What does not change

- The event delegation core (`on`/`emit`) is untouched.
- No htmx-specific code anywhere in `ui/` — this is generic dynamic-DOM
  support; the site's htmx boost is just its first consumer.
- Vendored-file flow: `gsxui add` already ships ui/*.js; no CLI changes.

## Testing

- Playwright (jstest): extend site-boost.spec.ts — boosted-navigate from
  /components/button to /components/select, assert the Basic trigger reads
  "Apple" without any click; boosted-navigate select→select (same page,
  morph-reuse path) and assert the trigger still reads "Apple".
- Harness spec for raw dynamic DOM (no htmx): a page that `innerHTML`-injects
  a select with a checked item after load, asserting reflection runs — this
  pins the library guarantee independently of htmx.
- Idempotency regression: carousel (stateful init) injected twice /
  re-inited must not double-bind (assert one set of controls behavior,
  e.g. next-button advances exactly one slide).
- All existing jstest specs stay green — they are the regression net for
  the conversion of ~10 modules.

## Out of scope

- Init for components that never had load-time scans.
- Teardown/destroy lifecycle (nothing needs it today; morphs keep nodes).
- htmx-specific integration events (htmx:afterSwap etc.).
