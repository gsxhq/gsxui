# Init conventions and shared helpers for ui/*.js

Date: 2026-08-05
Status: approved

## Problem

The per-component init code that grew during the self-healing-init migration
is inconsistent and partially duplicated. Survey findings (2026-08-05):

- Seven naming shapes for the same concept: `sweep` (avatar), `initRoot`
  (combobox, select — independently written), `normalize` (menubar,
  toggle-group), `initCarousel`, `(root) => filter(root)` (command), and
  anonymous inline arrows (calendar, input-otp, resizable). toaster.js has a
  module-private `function init()` that is NOT the core API — a grep trap.
- Duplicated logic: the group→label `aria-labelledby` wiring loop is
  near-verbatim in combobox and select; menubar and toggle-group duplicate
  the "bail on tabindex ATTRIBUTE, not IDL property" guard and its long
  rationale comment word-for-word; combobox, command, and select each
  declare a private `uid` counter for `gsxui-<module>-<kind>-${++uid}` ids.
- Idempotency machinery is ad hoc per module (nothing shared) even though
  the core's contract comment prescribes exactly two recurring shapes:
  idempotent reflection, and once-per-element resource binds.
- No file-layout convention: the `init(...)` registration call sits at a
  different position in every file.

## Design

Pure refactor — zero behavior change. Two parts.

### 1. Shared helpers in ui/gsxui.js

Exported next to `on`/`emit`/`init`/`position`:

- `uid(prefix)` → `` `${prefix}-${++counter}` `` with ONE module-level
  counter. Uniqueness within the document is the entire contract; a shared
  counter satisfies it. Replaces the three private counters.
- `wireGroupLabels(root, groupSelector, labelSelector, idPrefix)` — for
  each group under root without `aria-labelledby`: find its label, give the
  label a `uid(idPrefix)` id if missing, point the group at it. Extracted
  from the combobox/select duplicate; both call it.
- `once(fn)` → a wrapped initializer step that runs at most once per
  element (internal WeakSet). Owns the doctrine comment: initializers
  re-run on any mutation under a match, so RESOURCE BINDS (timers,
  listeners, observer handles) must be once-per-element while REFLECTION
  must re-run unguarded. Carousel's hand-rolled `autoplayBound` WeakSet
  becomes `once(bindAutoplay)`.
- `hasLiveTabStop(els)` → `els.some((el) => el.getAttribute("tabindex") === "0")`,
  carrying the attribute-vs-IDL-property rationale comment ONCE. Menubar
  and toggle-group both use it for their semantic "live roving value
  exists" bail (which is deliberately NOT `once()` — a morph that strips
  tabindex attributes must re-normalize).

### 2. Conventions across the 12 registering modules

- The function registered with core `init()` is always a named
  `initRoot(el)` — one per module, passed directly (no wrapper arrows).
  avatar/command/calendar/input-otp/resizable wrap their existing helpers
  in a named `initRoot`; menubar/toggle-group/carousel rename theirs.
  Internal helpers keep their domain names (`filter`, `recompute`,
  `captureDefaults`, …) — only the registered entry point is uniform.
- Every file ends with a uniform `// --- init ------` section containing
  `initRoot` and the single `init(selector, initRoot)` call. Load-time
  setup is always in one predictable place. avatar's window-load fallback
  sweep lives in that section too (it stays — images that settle before
  module evaluation are a real edge its comment documents).
- toaster.js's private `init()` is renamed `bootstrap()`, with a comment
  stating why it does not use the core `init()` (toast cards are cloned
  from templates, not server DOM matching a stable selector).
- The core contract comment in gsxui.js documents the `initRoot`
  convention, the reflection-vs-bind doctrine (pointing at `once`), and
  the fact that the shared observer is instantiated by whichever barrel
  import calls `init()` first (avatar today).

## Out of scope (known, deferred, documented in the core comment where apt)

- collect()'s O(registrations × records) walk — fine at 12 registrations;
  revisit only with evidence.
- The documented position()/release() listener leak for content removed
  from the DOM while open.
- toaster.js's eager module-level observers and `window.gsxui` exposure.
- Any change to the init()/observer semantics themselves.

## Testing

- Pure refactor: the full Playwright suite (dynamic-init, roving-tabindex,
  site-boost, select, calendar, command, plus style contracts) is the
  regression net and must stay green.
- One new assertion in jstest/specs/dynamic-init.spec.ts: `uid()` ids are
  unique across two different prefixes/modules on a live page (guards the
  shared-counter refactor against a future per-module-counter regression).
- `gsxui add` vendoring is unaffected (same files, no new modules).
