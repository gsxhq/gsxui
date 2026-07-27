# JS test layer — real-browser behavior tests in CI

gsxui ships 21 behavior modules totalling ~4,500 lines of JavaScript
(`ui/*.js`). None of it is tested. Every guarantee those modules make is
currently held up by a manual browser pass at the end of a batch.

That has already cost us. Tier 4 Batch B shipped two defects that reached a
live browser before anyone saw them:

- **The hook-prefix collision.** `dropdown`, `context-menu` and `menubar`
  shared a `data-gsxui-menu-*` attribute prefix. `ui/gsxui.js` keys its
  delegation registry by `${type}:${capture}` only and dispatches to *every*
  handler whose selector matches, so one click on a dropdown checkbox item
  ran two modules' handlers, fired two `gsxui:change` events, and left the
  state unchanged. This was live on every page shipping a dropdown.
- **The one-way click toggle.** `navigation-menu` opened on `focusin`, which
  fires mid-gesture, so a mouse click found the panel already open and the
  two handlers fought. The panel opened and would not close.

Both are cheap to catch automatically and expensive to catch by hand.

This document specifies the test layer. Scope: all 21 behavior modules,
delivered in three phases, with CI gating from the end of phase 1.

## 1. Why a real browser, not jsdom

Measured against jsdom 29.1.1, not assumed:

```
popover property:   undefined
showPopover:        undefined
dialog showModal:   undefined
getAnimations:      undefined
el.matches(":open"):           false   (no throw)
el.matches(":focus-visible"):  false   (no throw)
```

14 of the 21 modules use the popover or dialog APIs, and the same 14 read
geometry. jsdom cannot run them.

The silent `false` is the disqualifying part. Under jsdom a test asserting
"the closed menu is hidden" passes while asserting nothing, and
`:focus-visible` — the exact predicate the navigation-menu fix turns on —
always reports false. A test layer that cannot fail is worse than no test
layer, because it reads as coverage.

**Decision: Playwright Test against real Chromium, one runner for
everything.** No jsdom tier. A split runner would put the vacuous-test risk
back in the repo for the sake of a faster inner loop on the seven simplest
modules.

Firefox and WebKit are deferred, not rejected. Chromium-only is what CI runs
in phase 1; adding WebKit is a follow-up once the suite is stable, and is
worth doing because Safari's popover implementation is the one most likely
to diverge from ours.

## 2. Harness

`jstest/harness` is a Go `package main` that imports `site/examples` and
serves:

| route | serves |
|---|---|
| `/x/<component>/<example>` | one example, rendered through `ex.Node.Render`, in a minimal HTML shell |
| `/f/<name>` | a test-only fixture (§4), same shell |
| `/ui/*.js` | the real `ui/` modules, byte-for-byte, as native ES modules |
| `/shim/*.js` | the same modules, except `gsxui.js` is replaced by a recording `on()` (§5) |
| `/static/site.css` | the compiled stylesheet (§3) |
| `-manifest <path>` | flag, not a route: writes the component × example index and exits |

Three properties earn this over the alternatives:

**Real markup.** Pages render through the same `gsx.Node` values the site
serves, so a class-string or attribute change in a `.gsx` file reaches the
tests immediately. Hand-written HTML fixtures would drift silently — a
component could stop emitting the attribute its own tests bind to and stay
green.

**Real JS, no bundler.** Every import in `ui/` is relative, so
`<script type="module" src="/ui/index.js">` loads the shipped source
directly. Nothing between the test and the code under test — no Vite
transform, no bundle step to keep in sync, no build cache to invalidate.

**One example per page.** The site stacks several examples per component
page, which makes selectors ambiguous and couples tests to site layout. One
example per URL removes both problems.

Routes are keyed by the **registered** component name, not the directory
name: `/x/navigation-menu/basic`, not `/x/navigationmenu/basic`. The
directory names drop hyphens because Go package names must
(`navigationmenu`, `contextmenu`, `nativeselect`, `switchctl`), but
`examples.For` is keyed by the registered name and that is what the tests
should read.

The harness binds `127.0.0.1:7799` — deliberately clear of the dev loop's
7777 and Vite's 5173, so a running `make site-dev` does not collide with a
test run.

## 3. CSS

`@tailwindcss/cli` (new devDependency) compiles `web/site.css` to a temp
file that the harness serves at `/static/site.css`.

`web/site.css` already declares `@source "../ui/**/*.gsx"`, so every class
any component can emit is in the output. Fixture files add one more
`@source` line.

Real CSS is not optional. Two of the defect classes we have actually shipped
are invisible without it:

- **Ghost boxes.** An author `display` utility on a popover beats the UA
  stylesheet's `display: none` for the closed state, leaving an invisible
  but hit-testable box. This shipped in both `dialog` and `sidebar`. The
  assertion is a computed style, so it needs the compiled stylesheet.
- **Flex geometry.** `resizable`'s panel sizing depends on `flex-basis`
  resolution, which is why percentage bases failed and `0px` bases worked.

No Vite in the test path, and nothing written into the working tree — the
stylesheet lands in a gitignored temp dir.

## 4. Fixtures for scenarios examples do not cover

Most behavior has an example already: `dropdown` ships `basic`,
`checkboxes`, `radios`, `submenu` and `destructive`; `sidebar` ships
`basic`, `persisted` and `variants`. Where a test needs a state the showcase
does not demonstrate — a sidebar rendered collapsed, a menu whose first item
is disabled — the scenario gets a `.gsx` file under `jstest/fixtures/`,
registered into the harness's own registry.

Fixtures are `.gsx`, never hand-written HTML. They compose the real
components and go through `gsx generate` like everything else, so they
inherit the same drift protection the examples have.

They live in a separate registry from `site/examples` so test-only
scenarios never appear on the showcase site.

## 5. Test taxonomy

### Invariants — parameterized over every example

Roughly four tests covering all 103 examples across 53 components, because
each one is generated per example from the manifest:

1. **Clean load.** Zero console errors, zero unhandled rejections.
2. **No ghost boxes.** Every `[popover]` not in the `:open` state computes
   `display: none`.
3. **Selector disjointness.** Detailed below.
4. **No duplicate ids** within a rendered example.

**Selector disjointness** is the mechanical form of the hook-prefix
collision. It runs in two parts:

*Recording.* The page loads `/shim/index.js`. The shim's `on()` registers
nothing; it records `(type, capture, selector, module)`, taking `module`
from `new Error().stack` — the calling module's URL is on the second frame
in Chromium. `emit()` passes through unchanged.

*Checking.* For every example page, every recorded selector is run against
every element in the document. An element matched by selectors belonging to
two different modules, for the same `(type, capture)` pair, is a failure.

This is the real check, not a proxy for it: it tests the actual condition
that broke — two modules' handlers both firing for one element and one event
— against real rendered markup, rather than grepping source for attribute
prefixes. Grepping cannot see selectors held in variables, which is how
`CONTENT_SELECTOR` and friends are written throughout `ui/`.

Genuine overlaps, if any turn up, go in an explicit allowlist keyed by
module pair and selector. An allowlist entry is a decision someone made and
can be reviewed; a tolerance threshold is not.

### Per-component specs

One `.spec.ts` per behavior module, 21 in all. Each covers that module's
public contract: the state it reflects into the DOM, the events it emits,
its keyboard model, and its open/close lifecycle where it has one.

Two specs open with the regressions already paid for, written red-first
against the reverted fix:

- `dropdown.spec.ts` — on `/x/dropdown/checkboxes`, opening the menu and
  clicking one checkbox item fires **exactly one** `gsxui:change` and flips
  `aria-checked`. Against the shared-prefix version this sees two events and
  an unchanged attribute.
- `navigation-menu.spec.ts` — on `/x/navigation-menu/basic`, a real
  `page.mouse.click()` on a trigger opens it and a second click closes it;
  a `Tab` to the same trigger opens it, since keyboard focus is always
  focus-visible. Against the `focusin` version the second click does
  nothing.

## 6. Timing

**No global animation disabling.** The obvious move — inject
`* { transition: none !important }` for determinism — is wrong here.
`dialog.js` awaits finite own Web Animations, while `sonner.js` waits on
`transitionend` with a 600ms fallback cap. With animations off, dialog loses
the exit boundary under test and Sonner falls through to its cap. That hides
real timing bugs and makes the suite slower.

Playwright's retrying assertions handle the waiting instead. Tests assert on
settled state and let the runner poll.

One environment note carried over from the manual passes: occluded Chrome
tabs freeze animation clocks, leaving transitions permanently pending. This
does not apply to headless Chromium under Playwright, and frozen-transition
readings must not be encoded as expectations.

## 7. Wiring

`playwright.config.ts` lives at `jstest/`:

- `globalSetup` runs `go run ./jstest/harness -manifest <tmp>/examples.json`
  and the Tailwind CLI build. Playwright runs `globalSetup` before workers
  import spec files, so specs read the manifest synchronously with
  `readFileSync`. This is what lets `npx playwright test` work on its own,
  rather than only through `make`.
- `webServer` starts the harness on 7799 with
  `reuseExistingServer: !process.env.CI`.
- Chromium only. `retries: 1` and `forbidOnly` on CI, `retries: 0` locally.
  Trace on first retry.

Make targets:

- `make test` stays Go-only and fast. A Go-only edit should not pay for a
  browser boot.
- `make test-js` runs the browser suite.
- `make check` becomes `test` + `test-js` + the existing generated-file and
  gofmt checks.

CI: the existing `test` job in `.github/workflows/fly-deploy.yml` gains Node
22 (matching the Dockerfile), a cached Playwright Chromium, and a
`make test-js` step. `deploy` already depends on `test`, so it gates on both
without a new job or a new dependency edge.

`npm ci` in `site/Dockerfile` will install the new devDependencies.
`@playwright/test` does not download browsers on install — only
`npx playwright install` does — so the image build is unaffected beyond a
marginally slower `npm ci`.

Specs are TypeScript. Playwright transpiles `.ts` with no tsconfig required,
and the locator API is close to unusable without types. This matches
`vite.config.ts`, already the repo's one TypeScript file.

## 8. Phases

**Phase 1 — infrastructure and invariants.** Harness, CSS build, Playwright
config, Make targets, CI wiring, and the four invariant sweeps. CI is green
and genuinely useful at the end of this phase: the ghost-box and
selector-collision classes are both covered across all 103 examples before a
single per-component spec exists.

**Phase 2 — the high-risk eight.** `dropdown`, `context-menu`, `menubar`,
`navigation-menu`, `dialog`, `select`, `combobox`, `sidebar`. Opens with the
two regression tests from §5.

**Phase 3 — the remaining thirteen.** `resizable` (including flex geometry),
`command`, `carousel`, `sonner`, `popover`, `hover-card`, `tooltip`,
`tabs`, `toggle`, `toggle-group`, `input-otp`, `slider`, `avatar`.

## 9. Out of scope

- **Visual regression / screenshot diffing.** A separate discipline with its
  own flake profile and its own storage story. The layer specified here
  asserts behavior and computed style, not pixels.
- **Firefox and WebKit.** Deferred as §1 records.
- **Testing `web/site.js` and `web/theme.js`.** Site code, not library code.
  The harness deliberately does not load them.
- **Go-side rendering tests.** Already covered by
  `site/examples/examples_test.go` and the registry tests.
