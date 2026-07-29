# Component Migration Waves Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Migrate the remaining 51 components onto the slot axis, so the whole catalogue ships as generated `(canonical × recipe)` source and `assets/css/styles/default.css` is retired.

**Architecture:** Each component is migrated independently: its rules are lifted out of `default.css` into per-style recipe stylesheets, its shape is declared, and its `.gsx` gains accessor calls. Migrations run as parallel agents in isolated worktrees, because the only shared artifacts are `default.css` (disjoint deletions) and `recipes.json` (regenerated), and the controller merges them serially. Acceptance for every component is a computed-style sweep proving rendering did not change.

**Tech Stack:** Go 1.24+, `github.com/gsxhq/gsx`, Tailwind CLI, Playwright + Chromium.

**Spec:** `docs/superpowers/specs/2026-07-29-slot-axis-and-component-migration-design.md`

## Global Constraints

- Recipe class prefix is exactly `gsxui-recipe-`; it must never appear in a generated consumer artifact.
- Registry/CSS identity stays kebab-case. Go type and receiver are **derived** (`input-group` → `inputGroupRecipe`, receiver `inputGroup`).
- Styles are strictly conformant: every style implements every `(slot, dimension, value)`.
- `make audit` must pass — it runs the layer gate over six stylesheets, per style, using the implication model.
- **A component is not migrated until its computed-style sweep is clean.** "No test failed" is not evidence; two visible regressions already passed a 313-test suite.
- Intermediate commits need not be green; the final state of each wave must be.
- `gofmt -l .` prints nothing. In gsx markup, comments are `{/* */}`.
- Do **not** bump `github.com/gsxhq/gsx`. The available version bundles a URL-sanitization redesign and is deferred until after this migration.

## What migration actually is (read before Task 1)

This is not "move CSS rules". Three facts shape every task:

1. **47 of 53 components have no `class=` attribute at all.** They are styled entirely through `data-gsxui-slot-*` markers. Migration *adds* class attributes carrying accessor calls to elements that have none. Button, having had classes already, was the least representative pilot possible.

2. **103 rules are relational** — `:has()` and combinators. A recipe slot emits a flat class list on one element and cannot express "when this card contains a footer". These must be re-expressed as Tailwind variants on the owning element. The dominant form is `:has(> svg)`, which becomes `has-[>svg]:*` — Button's own recipe already does this (`has-[>svg]:px-2`). Zero rules need `group/*` markers.

3. **Some rules will not translate.** The escape hatch is to leave the rule in CSS, in `@layer utilities`, with specificity ≥ (0,1,0) — the layer gate enforces exactly that. Using the hatch is legitimate; using it silently is not. Every retained rule needs a comment saying why.

## File Structure

Per migrated component `<c>`:

- Create: `registry/canonical/shapes/<c>.go` — the shape (hand-written, pure data)
- Create: `registry/styles/nova/<c>.css`, `registry/styles/maia/<c>.css` — recipe rules
- Create: `registry/canonical/<c>.gsx` — structural source with accessor calls
- Generated: `registry/canonical/<c>_recipe.gen.go`, `registry/generated/<style>/<c>.gsx`, `ui/<c>.gsx`, `site/stylepreview/<style>/<c>.gsx`
- Modify: `assets/css/styles/default.css` — delete the component's migrated rules
- Modify: `registry/canonical/shapes/shapes.go` — register the shape

Shared, regenerated: `registry/generated/recipes.json`.

---

### Task 1: The migration playbook, proven on Card

**Files:**
- Create: `docs/superpowers/plans/migration-playbook.md`
- Create: `registry/canonical/shapes/card.go`, `registry/styles/{nova,maia}/card.css`, `registry/canonical/card.gsx`
- Modify: `assets/css/styles/default.css`, `registry/canonical/shapes/shapes.go`
- Test: `jstest/specs/layer-precedence.spec.ts` (add Card pins)

**Interfaces:**
- Consumes: `recipe.Shape`/`Slot`/`Dimension`, `shapes.All()`, the layer gate, the sweep harness
- Produces: `docs/superpowers/plans/migration-playbook.md` — the procedure every later task follows verbatim

Card is the right pilot: 7 slots, **zero dimensions**, 13 rules, and it exercises the relational case (`:has([data-gsxui-slot-card-footer])`) without also exercising switch generation.

- [ ] **Step 1: Capture the sweep baseline**

```bash
make sweep-baseline
cp -r jstest/.tmp/sweep-baseline /tmp/card-baseline
```

This is the acceptance evidence. Without it the migration cannot be verified.

- [ ] **Step 2: Inventory Card's rules**

```bash
grep -n "data-gsxui-slot-card" assets/css/styles/default.css
```

Record every rule with its line number and classify each as: **slot base** (plain marker rule), **relational** (`:has()`, combinator), or **untranslatable** (goes to the escape hatch). Card's rules run 200–236 plus one stray at 2442 — check for strays on every component, they will not all be contiguous.

- [ ] **Step 3: Declare the shape**

`registry/canonical/shapes/card.go`:

```go
package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Card is a pure container: seven slots, no dimensions. Every slot carries only
// a base rule, which is the common shape across the catalogue — most components
// vary structurally rather than by variant.
var Card = recipe.Shape{
	Component: "card",
	Slots: []recipe.Slot{
		{Name: "", Base: true},
		{Name: "header", Base: true},
		{Name: "title", Base: true},
		{Name: "description", Base: true},
		{Name: "action", Base: true},
		{Name: "content", Base: true},
		{Name: "footer", Base: true},
	},
}
```

Register it in `shapes.go`'s `all` map: `Card.Component: Card,`.

- [ ] **Step 4: Write the recipe stylesheets**

`registry/styles/nova/card.css` carries one rule per slot, in `@layer components`, using only `@apply`. Copy the utilities verbatim from the rules inventoried in Step 2 — do not "improve" them; a visual change here is indistinguishable from a migration bug in the sweep.

Relational rules fold into the owning slot's base as variants:

```css
/* was:  :where([data-gsxui-slot-card]):has([data-gsxui-slot-card-footer]) */
.gsxui-recipe-card {
  @apply … has-[[data-gsxui-slot-card-footer]]:pb-0;
}
```

`registry/styles/maia/card.css` must implement the identical shape. If Maia's `default.css` rules are the same as Nova's, the file is the same — that is expected and not a smell. Conformance is checked both ways, so a missing slot fails the build naming it.

- [ ] **Step 5: Add accessor calls to the canonical**

```bash
cp ui/card.gsx registry/canonical/card.gsx
```

Change `package ui` to `package canonical`. Add the binding to `shapes.go`'s sibling file or alongside the shape:

```go
var card = cardRecipe{c: recipe.Component{Shape: shapes.Card}}
```

Then add a `class={ … }` attribute to each marked element. `ui/card.gsx` has **none today**, so these are new:

```gsx
<div class={ card.Root() } data-gsxui-slot-card>
	<div class={ card.Header() } data-gsxui-slot-card-header>
```

Leave the `data-gsxui-slot-*` markers exactly as they are — they are the style contract and the layer gate's join key.

- [ ] **Step 6: Delete the migrated rules from default.css**

Remove exactly the rules inventoried in Step 2 that you re-expressed. Leave anything you deliberately did not migrate, move it into the `@layer utilities` block, and comment why.

- [ ] **Step 7: Generate and check**

```bash
go run ./cmd/stylegen
go tool gsx generate
go build ./... && go test ./... -count=1
go run ./cmd/stylegen --check   # exit 0
make audit                       # the layer gate, per style
```

`make audit` failing here is the gate working: it means a rule left in `@layer components` is now dead. Fix the rule, do not weaken the gate.

- [ ] **Step 8: THE ACCEPTANCE TEST — sweep**

```bash
make sweep-compare
```

Every difference is a finding. Card's computed styles must be identical before and after: the whole point is that presentation moved layers without changing. Investigate each diff; do not accept any without a written reason.

- [ ] **Step 9: Add regression pins**

Add two pins to `jstest/specs/layer-precedence.spec.ts` for properties Card would visibly lose — its border radius and its header border. Follow the existing carousel/input-group pins. The sweep catches drift between runs; pins catch it in CI forever.

- [ ] **Step 10: Write the playbook**

Write `docs/superpowers/plans/migration-playbook.md` recording Steps 1–9 as a reusable procedure, plus everything Card taught you: which relational forms translated and how, whether rules were contiguous, whether the escape hatch was needed, and how long it took. Every later task follows this document rather than re-deriving the procedure.

- [ ] **Step 11: Commit**

```bash
git add -A
git commit -m "feat: migrate Card to the slot axis"
```

---

### Task 2: Badge and Alert — the dimension case

**Files:** per-component set as in Task 1, for `badge` and `alert`
**Interfaces:** Consumes the playbook from Task 1

Badge (20 rules) and Alert (18) are the smallest components with **per-slot dimensions** (`[data-variant]`). They prove switch generation on a slot that is not the root — the case Button could not exercise.

- [ ] **Step 1: Follow the playbook for Badge**

Its shape has a root slot with a `variant` dimension:

```go
var Badge = recipe.Shape{
	Component: "badge",
	Slots: []recipe.Slot{{
		Name: "", Base: true,
		Dimensions: []recipe.Dimension{
			{Name: "variant", Default: "default", Values: []string{
				"default", "secondary", "destructive", "outline"}},
		},
	}},
}
```

Read the real values out of `default.css`'s `[data-gsxui-slot-badge][data-variant="…"]` rules rather than trusting this list.

- [ ] **Step 2: Verify the generated switch**

Read `registry/generated/nova/badge.gsx` and confirm the `default:` arm carries the **declared default's** utilities, not `""`. That property is what makes a misspelled variant render the default instead of an unstyled element, and it is pinned for Button by `TestResolveDefaultArmCarriesDeclaredDefault`.

- [ ] **Step 3: Follow the playbook for Alert**

- [ ] **Step 4: Run the full gate set and sweep for both**

```bash
go run ./cmd/stylegen --check && make audit && make sweep-compare
```

- [ ] **Step 5: Update the playbook**

Record what dimensions changed about the procedure — in particular how you derived the value list and whether the style contract's axis values agreed with the shape's (the agreement test will tell you).

- [ ] **Step 6: Commit**

```bash
git add -A
git commit -m "feat: migrate Badge and Alert to the slot axis"
```

---

### Task 3: The parallel migration protocol

**Files:**
- Create: `docs/superpowers/plans/migration-parallel-protocol.md`
- Modify: `Makefile` (a `sweep-baseline-ref` target)

**Interfaces:**
- Produces: the procedure the controller and every wave agent follow

Only two artifacts collide when migrations run concurrently:

| Artifact | Per-component? | Collides? |
|---|---|---|
| `shapes/<c>.go`, `canonical/<c>.gsx`, `styles/*/<c>.css`, `ui/<c>.gsx` | yes | no |
| `assets/css/styles/default.css` | no — each migration **deletes** its own rules | yes |
| `registry/generated/recipes.json` | no — regenerated whole | yes |
| `registry/canonical/shapes/shapes.go` | no — one line added per component | yes (trivially) |

- [ ] **Step 1: Write the protocol**

Record this as the rule, with the reasoning:

**Each wave agent works in its own git worktree** (`isolation: "worktree"`), performing the complete playbook including generation, gate and sweep. This is deliberate: an agent that cannot run `stylegen`, `make audit` and `make sweep-compare` cannot verify its own work, and unverified migrations are precisely what this project has repeatedly caught late.

**The controller merges serially.** `default.css` deletions are disjoint line ranges and merge cleanly in the common case; `shapes.go` gains one line per component; `recipes.json` is regenerated after each merge rather than merged. On any conflict the controller re-runs `go run ./cmd/stylegen` and `make audit` before proceeding to the next agent.

**The sweep baseline is taken once per wave, before any agent starts**, from the wave's base commit — not per agent, or each agent measures against a different tree.

- [ ] **Step 2: Add the baseline-from-ref target**

```make
sweep-baseline-ref:
	@test -n "$(REF)" || (echo "usage: make sweep-baseline-ref REF=<commit>" && exit 1)
	git worktree add /tmp/gsxui-sweep-ref $(REF)
	ln -sf $(CURDIR)/node_modules /tmp/gsxui-sweep-ref/node_modules
	cd /tmp/gsxui-sweep-ref && SWEEP_OUT=$(CURDIR)/jstest/.tmp/sweep-baseline \
		npx playwright test --config jstest/playwright.config.ts jstest/specs/layer-precedence.spec.ts
	git worktree remove --force /tmp/gsxui-sweep-ref
```

Add both target names to the existing `.PHONY` line.

- [ ] **Step 3: Commit**

```bash
git add docs/superpowers/plans/migration-parallel-protocol.md Makefile
git commit -m "docs: define the parallel migration protocol"
```

---

### Task 4: Wave 1 — the trivial components (parallel)

**Components:** `skeleton` (1 rule), `spinner` (1), `progress` (2), `avatar` (3), `label` (3), `aspect-ratio`, `separator` (7), `kbd` (7)

These are single- or few-rule components. Wave 1 exists to exercise the **parallel protocol itself** on work where a mistake is cheap, before the protocol carries anything large.

- [ ] **Step 1: Take the wave baseline**

```bash
make sweep-baseline-ref REF=$(git rev-parse HEAD)
```

- [ ] **Step 2: Dispatch one agent per component, in parallel, each in its own worktree**

Each agent receives: the playbook, its single component name, and the wave base commit. No agent is given more than one component — the point is that a failure is attributable.

- [ ] **Step 3: Merge serially, regenerating between each**

```bash
go run ./cmd/stylegen && go build ./... && go test ./... -count=1 && make audit
```

after every merge, not once at the end. A failure then names the component that caused it.

- [ ] **Step 4: Sweep the whole wave**

```bash
make sweep-compare
```

- [ ] **Step 5: Record what the protocol got wrong**

Update `migration-parallel-protocol.md` with every conflict, mis-merge or surprise. Wave 1's real deliverable is a protocol that survived contact.

- [ ] **Step 6: Commit**

---

### Task 5: Wave 2 — mid-size components (parallel)

**Components:** `accordion`, `breadcrumb`, `collapsible`, `empty`, `field`, `item`, `input`, `textarea`, `checkbox`, `radio`, `switch`, `slider`, `toggle`, `tabs`, `table`, `tooltip`, `hover-card`, `popover`, `scroll-area`

Includes the first **hyphenated** components (`hover-card`, `scroll-area`), which exercise the derived Go identifier (`hoverCard`, `scrollArea`) end to end for the first time — the path Button could not test.

- [ ] **Step 1: Baseline, dispatch, merge, sweep**

Follow `docs/superpowers/plans/migration-parallel-protocol.md` exactly, as revised by Wave 1: take one baseline for the wave from its base commit, dispatch one agent per component in its own worktree, merge serially running `go run ./cmd/stylegen && go build ./... && go test ./... -count=1 && make audit` after **each** merge, then `make sweep-compare` across the whole wave.

- [ ] **Step 2: Confirm the derived identifiers appear correctly**

```bash
grep -n "hoverCardRecipe\|scrollAreaRecipe" registry/canonical/*_recipe.gen.go
grep -c "gsxui-recipe-hover-card" registry/styles/nova/hover-card.css
```

The Go identifier is camel; the CSS identity stays kebab. Both must be true simultaneously.

- [ ] **Step 3: Commit**

---

### Task 6: Wave 3 — composite and menu components (parallel)

**Components:** `alert-dialog`, `dialog`, `drawer`, `sheet`, `dropdown`, `context-menu`, `menubar`, `navigation-menu`, `select`, `native-select`, `combobox`, `command`, `input-group`, `input-otp`, `button-group`, `toggle-group`, `pagination`, `carousel`, `resizable`, `sonner`, `calendar`

These compose other components, so their `@layer components` rules are the ones most likely to be demoted by an already-migrated dependency. Expect the layer gate to fire here; that is why it exists.

- [ ] **Step 1: Resolve Sonner's naming exception first**

Sonner's registry name is `sonner` but its style-contract slots are `toaster`, `toast`, `toast-*`, which breaks the `<RegistryName>-<relative>` mapping the agreement test joins on. Decide: rename the contract slots, or record a documented exception in `agreement_test.go`. Do this **before** dispatching the wave, because every Sonner agent would otherwise hit it.

- [ ] **Step 2: Baseline, dispatch, merge, sweep**

Follow `docs/superpowers/plans/migration-parallel-protocol.md`: one wave baseline from the base commit, one agent per component in its own worktree, serial merge with `go run ./cmd/stylegen && go build ./... && go test ./... -count=1 && make audit` after each, then `make sweep-compare` for the wave.

- [ ] **Step 3: Commit**

---

### Task 7: Wave 4 — Sidebar

**Component:** `sidebar` alone — 38 slots, composes Button, and its rules live in both `default.css` and `assets/css/foundation.css`.

Sidebar gets its own wave because it is larger than most waves and because it is the component the layer gate was widened to cover.

- [ ] **Step 1: Inventory rules across BOTH stylesheets**

```bash
grep -n "data-gsxui-slot-sidebar" assets/css/styles/default.css assets/css/foundation.css
```

`foundation.css` rules are foundation *mechanics* — they may legitimately stay. Decide per rule, and record the reasoning; this is the judgement the gate cannot make.

- [ ] **Step 2: Follow the playbook, then gate and sweep**

- [ ] **Step 3: Commit**

---

### Task 8: Retire default.css and close out

**Files:** `assets/css/styles/default.css`, `Makefile`, `docs/`

- [ ] **Step 1: Confirm what remains**

```bash
grep -c "data-gsxui-slot" assets/css/styles/default.css
```

Anything left is either a deliberate escape-hatch rule (in `@layer utilities`, commented) or an unmigrated component. Enumerate and justify each.

- [ ] **Step 2: Full verification**

```bash
go build ./... && make audit && go run ./cmd/stylegen --check && go test ./... -count=1
gofmt -l . && git status --short
make sweep-compare        # against the pre-migration baseline
npx playwright test --config jstest/playwright.config.ts
```

- [ ] **Step 3: Update the specs**

Mark the migration complete in `2026-07-29-slot-axis-and-component-migration-design.md`, and revise `docs/theme-system-roadmap.md`, which still describes the CSS-only architecture as the outcome for everything except Button.

- [ ] **Step 4: Commit**

---

## Coverage

53 components total: 4 already migrated or piloted here (`button`, `card`,
`badge`, `alert`), Wave 1 has 8, Wave 2 has 19, Wave 3 has 21, Wave 4 has 1
(`sidebar`). Every component appears in exactly one wave; verified by
enumeration against `ui/*.gsx`.

## Notes for the implementer

**The sweep is the acceptance test, not the suite.** Every component's migration is accepted on computed styles being unchanged. Two visible regressions — non-round carousel arrows, wrong input-group type scale — already passed build, tests, drift check, `make audit` and a 313-test browser suite. Do not substitute a green suite for a clean sweep.

**Copy utilities verbatim.** Improving a utility while migrating it makes a visual change indistinguishable from a migration bug in the sweep diff. Migrate first; improve in a separate commit with its own sweep.

**The gate firing is success.** When `make audit` fails on a composite component in Wave 3, it has found a rule that a migrated dependency demoted. Fix the rule — move it to `@layer utilities` with specificity ≥ (0,1,0), or fold it into the owning slot. Never weaken the gate to pass.

**One component per agent.** A failed wave must be attributable to a component. Agents are cheap; a wave you cannot bisect is not.

**Rules are not always contiguous.** Card's live at 200–236 *and* 2442. Grep for the marker across the whole file, not just around the first hit.
