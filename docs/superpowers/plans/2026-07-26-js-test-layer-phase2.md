# JS Test Layer — Phase 2 Implementation Plan

**Goal:** Add real-browser behavior specifications for the eight highest-risk
JavaScript modules and complete the fixture route promised by the JS test-layer
design.

**Architecture:** Keep Phase 1's real Chromium/native-module harness. Tests use
registered showcase examples whenever those expose the contract. Test-only
states live in generated `.gsx` fixtures served at `/f/<name>`; no hand-written
HTML enters the behavior suite. Each module gets one focused spec that asserts
DOM state, emitted events, keyboard behavior, and lifecycle.

**Tech Stack:** Go 1.26, gsx generated fixtures, Playwright Test, TypeScript,
native browser popover/dialog APIs.

**Source spec:** `docs/superpowers/specs/2026-07-25-js-test-layer-design.md`

## Global constraints

- Load the real `/ui/index.js`; no mocks, bundler, jsdom, or global animation
  disabling.
- Prefer `/x/<component>/<example>`. Add a fixture only for a state no shipped
  example demonstrates.
- Every new behavioral assertion must be red-validated against a deliberate
  local mutation, then restored before commit.
- Assert public contracts: reflected DOM/ARIA state, native form values,
  emitted events, focus, and open/close lifecycle. Do not reach into module
  locals.
- Keep each spec independently runnable with
  `npx playwright test --config jstest/playwright.config.ts <name>`.
- Never start or kill the user's dev servers. The Playwright rig owns only its
  loopback harness on port 7799.
- Any fixture `.gsx` change is followed by `go tool gsx generate`; generated
  fixture output is committed and checked for drift.

## Task 1: Generated fixture registry and `/f/<name>`

**Files**

- Create: `jstest/fixtures/registry.go`
- Create: `jstest/fixtures/select-disabled.gsx`
- Generated: `jstest/fixtures/select-disabled.x.go`
- Modify: `jstest/harness/main.go`
- Modify: `jstest/harness/harness_test.go`
- Create: `jstest/input.css`
- Modify: `jstest/global-setup.ts`

**Contract**

- `fixtures.Register(name, node)` rejects duplicate names.
- `fixtures.For(name)` returns one test-only `gsx.Node`.
- `GET /f/{name}` renders that node through the same shell and real module
  script as `/x/`.
- Unknown fixtures return 404.
- Tailwind scans fixture `.gsx` sources through a test-only CSS entry; fixture
  classes do not leak into the production stylesheet entry.

**Steps**

1. Add failing Go tests for a rendered `/f/select-disabled`, an unknown fixture
   404, and duplicate registration.
2. Run `go test ./jstest/harness ./jstest/fixtures`; confirm missing registry and
   route failures.
3. Implement the registry, route, and a select fixture whose first option is
   disabled and whose second is enabled.
4. Compile Playwright CSS from `jstest/input.css`, which imports
   `../web/site.css` and adds `@source "./fixtures/**/*.gsx"`.
5. Generate, run the Go tests, load the fixture in Chromium, and prove the route
   test red by temporarily removing its registration.
6. Commit as `test(js): add generated behavior fixtures`.

## Task 2: Dropdown and context-menu behavior

**Files**

- Create: `jstest/specs/dropdown.spec.ts`
- Create: `jstest/specs/context-menu.spec.ts`

**Dropdown contract**

- A real trigger click opens; a second closes; `aria-expanded`, `data-state`,
  and exactly one `gsxui:open`/`gsxui:close` follow.
- The paid regression: clicking one checkbox emits exactly one
  `gsxui:change`, flips `aria-checked`, and leaves the menu open.
- Radio selection unchecks its sibling, emits one value, and closes.
- Arrow/Home/End navigation skips disabled items; Enter selects; Escape closes
  and restores trigger focus.
- Submenu pointer and keyboard paths open only the owned submenu and close
  through its grace period without closing the parent prematurely.

**Context-menu contract**

- A real `contextmenu` gesture opens at the pointer, prevents the native menu,
  stamps open state, and emits once.
- Item, checkbox, radio, keyboard, Escape, and submenu behavior match the
  dropdown family while staying on context-menu-specific hooks.

**Red validation**

- Temporarily register dropdown's checkbox click handler a second time. The
  regression must report two events and unchanged checked state.
- Temporarily point one context-menu checkbox selector at the dropdown prefix.
  Its own spec must stop changing the context-menu item; selector invariants
  must also report the collision where applicable.

Commit as `test(js): cover dropdown and context-menu behavior`.

## Task 3: Menubar and navigation-menu behavior

**Files**

- Create: `jstest/specs/menubar.spec.ts`
- Create: `jstest/specs/navigation-menu.spec.ts`

**Menubar contract**

- Exactly one top-level trigger is tabbable; horizontal arrows/Home/End move
  real focus and skip disabled triggers.
- Click toggles one menu. While a menu is open, pointer movement and horizontal
  keyboard movement switch to the sibling menu without leaving two open.
- Menu-item, checkbox, radio, submenu, Escape, and open/close events follow the
  owned menubar hooks.

**Navigation-menu contract**

- The paid regression: one real mouse click opens and the second closes.
- Keyboard Tab focus opens because it is focus-visible; mouse focus alone does
  not pre-open and fight the click.
- Pointer delay/grace behavior keeps content open while moving from trigger to
  panel, switches sibling panels, and eventually closes after leaving both.
- `aria-expanded`, `data-state`, indicator position/state, exactly-once
  lifecycle events, and link `gsxui:select` are synchronized.

**Red validation**

- Restore unconditional open-on-focus in `navigation-menu.js`; the second-click
  regression must fail.
- Disable menubar's sibling-close loop; the switch test must observe two open
  popovers.

Commit as `test(js): cover menubar and navigation-menu behavior`.

## Task 4: Dialog behavior

**Files**

- Create: `jstest/specs/dialog.spec.ts`

**Contract**

- Trigger opens with `showModal`, synchronized `open`/`data-state` and
  `aria-expanded`; title/description wire `aria-labelledby`/
  `aria-describedby`.
- Close button, Escape/request-close, and backdrop coordinate click each close
  once; clicks inside the dialog do not light-dismiss.
- Closing waits for the observable transition path before native `open` leaves,
  then restores trigger focus.
- Native programmatic `showModal()`/`close()` still produce exactly one
  `gsxui:open`/`gsxui:close`.

**Red validation**

- Remove the backdrop coordinate guard; an inside-content click must close and
  fail the regression.
- Emit close directly from the close-button handler as well as toggle; the
  exactly-once assertion must report two events.

Commit as `test(js): cover dialog behavior`.

## Task 5: Select behavior

**Files**

- Create: `jstest/specs/select.spec.ts`

**Contract**

- Init populates the hidden native select bridge, reflects any server-checked
  option, and wires listbox ids/ARIA.
- Pointer and keyboard opening synchronize trigger/content state and lifecycle
  events. Escape restores trigger focus.
- Hover moves focus; mouse, touch/click, Enter, and Space commit exactly once,
  synchronize trigger label, checked item, bridge value/change event, and
  `gsxui:select`.
- Arrow/Home/End skip disabled options. Closed-trigger keys open on the correct
  item.
- Prefix typeahead and repeated-character cycling work both closed and open,
  with the one-second buffer reset.

**Fixture use**

- `/f/select-disabled` makes disabled-first entry and skip behavior explicit
  without weakening the public example.

**Red validation**

- Remove the disabled filter from the focusable item list; the fixture's first
  ArrowDown/open assertion must land on the disabled item and fail.
- Remove bridge assignment in `applyValue`; FormData/value assertions must
  fail while visible state still changes.

Commit as `test(js): cover select behavior`.

## Task 6: Combobox behavior

**Files**

- Create: `jstest/specs/combobox.spec.ts`

**Contract**

- Plain focus does not open; click, trigger, and typing do. Lifecycle state and
  events remain exactly once.
- Filtering is case/accent-insensitive contains matching, hides/shows without
  DOM reordering, and exposes empty state on content and list.
- Highlight stays in the input through `aria-activedescendant`; arrows/Home/End
  skip hidden/disabled items; Enter commits and Escape closes.
- Commit and clear synchronize visible label, item state, hidden bridge,
  native change, `gsxui:select`, and focus.
- Reopening after a commit shows all options; native form reset restores the
  server value after the browser's reset action.
- Groups receive stable `aria-labelledby` wiring.

**Red validation**

- Reopen using the committed label as the filter query; the all-options
  assertion must fail.
- Reflect reset in the reset event's microtask instead of the next task; the
  mixed/reset assertion must retain the stale checked item and fail.

Commit as `test(js): cover combobox behavior`.

## Task 7: Sidebar behavior

**Files**

- Create: `jstest/specs/sidebar.spec.ts`

**Contract**

- Desktop trigger and rail toggle the wrapper/root state, emit one
  `gsxui:change {open}`, and update the tooltip/collapsible-visible surface.
- `Ctrl/Cmd+B` toggles globally, but not from editable targets or with
  disallowed modifiers.
- Narrow viewport uses the mobile dialog path and close trigger rather than
  mutating desktop state.
- Multiple sidebar providers stay scoped to their closest wrapper.
- Persisted example's listener writes the documented cookie only for the
  sidebar wrapper event.

**Red validation**

- Remove the editable-target guard; the shortcut-in-input assertion must
  toggle and fail.
- Force desktop mode at a narrow viewport; the mobile dialog assertion must
  fail.

Commit as `test(js): cover sidebar behavior`.

## Task 8: Phase gate and documentation

**Files**

- Modify: `README.md`
- Modify: `CHANGELOG.md`
- Modify: `docs/jsx-parity.md`
- Modify: `docs/superpowers/specs/2026-07-25-js-test-layer-design.md`

**Steps**

1. Record Phase 2 completion and link the eight specs plus fixture route.
2. Keep honest scope: Phase 3's thirteen modules remain explicitly pending.
3. Run each new spec alone once, then `make check`.
4. Run `git diff --check`, confirm generated fixture output and test CSS are
   current, and inspect the full branch diff.
5. Commit as `docs(js): record high-risk behavior coverage`.

## Completion evidence

- `/f/<name>` is generated-GSX-backed and covered by Go/harness tests.
- Eight new module specs pass against real Chromium.
- Both paid regressions have committed tests and recorded red validation.
- Every new assertion family was observed failing under its named mutation.
- `make check` passes with the full Go and browser suites.
- Phase 3 remains a separate follow-up branch/plan for the remaining thirteen
  modules.
