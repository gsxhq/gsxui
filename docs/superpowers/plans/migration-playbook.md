# Component migration playbook

**Status:** Proven once, on Card (`registry/canonical/shapes/card.go`,
`registry/canonical/card.gsx`, `registry/styles/{nova,maia}/card.css`).
This is the procedure every later component migration follows. It records
what Card actually taught, not an idealised version of the steps.

**Time:** roughly 2-3 focused hours for a 7-slot, zero-dimension component,
most of it in inventory (Step 2) and pin selection (Step 9) — the mechanical
steps (3-8) are fast once the inventory is right.

## Procedure

### Step 1 — capture the sweep baseline, before touching anything

```bash
make sweep-baseline
cp -r jstest/.tmp/sweep-baseline /tmp/<component>-baseline
```

This is the acceptance evidence, full stop. Do this first; there is no way
to reconstruct it after the fact.

### Step 2 — inventory the component's rules

```bash
grep -n "data-gsxui-slot-<component>" assets/css/styles/default.css
```

Classify every match into one of three buckets:

- **slot base** — a plain marker rule, one `@apply` block per slot.
- **relational** — a `:has()`, a combinator, or a compound selector that
  reacts to another marker or a consumer-supplied class.
- **untranslatable** — genuinely cannot be re-expressed as a Tailwind
  utility; goes to the escape hatch (§ "What actually needed the escape
  hatch" below).

**Finding: grep alone does not find every rule that matters, and does not
tell you which matches are actually yours.** For Card:

- The main block was contiguous, `default.css` lines 200-235 (10 rule
  blocks, not "13" as an earlier draft of this brief estimated — count
  actual `{ ... }` blocks, don't trust a remembered number).
- `grep -n "data-gsxui-slot-card"` also matched a line near 2404
  (`:where([data-gsxui-slot-card-content]) :where([data-gsxui-slot-calendar])`).
  **This is not a Card rule.** It styles *Calendar* (`background:
  transparent`) conditionally on Calendar being nested inside a
  `card-content` or `popover-content` ancestor. The marker it's keyed on is
  Card's, but the element it's styling and setting properties on is
  Calendar's. The test for "is this actually mine" is not "does the
  component's marker appear in the selector" — it's "does the LAST compound
  selector (the styled subject) carry your marker." Leave rules like this
  exactly where they are; migrating them would require migrating the other
  component (Calendar) too, and would attach foreign behaviour to your
  recipe's classes.
- Grep for strays on every component before starting Step 3 — they will not
  all be contiguous, and a component's rules can appear inside another
  component's block (as `card-content` did inside Calendar's) with no
  markers of their own to signal it.

### Step 3 — declare the shape

`registry/canonical/shapes/<component>.go`, one `recipe.Slot` per slot, in
document order. Register it in `shapes.go`'s `all` map.

### Step 4 — write the recipe stylesheets

`registry/styles/nova/<component>.css`, `%layer components`, one rule per
slot, `@apply` only. Copy utilities **verbatim** — a "tidied" or reordered
utility list is indistinguishable from a bug in the sweep diff.

**Finding: `@apply`-only is enforced, not a style guideline.**
`recipe.ParseStyle` rejects an "ordinary declaration in recipe rule" —
`go run ./cmd/stylegen` fails outright if a recipe rule contains anything
that isn't `@apply`. This is what forces plain-CSS properties (Card's
`container-name`/`container-type`, which have no Tailwind utility
equivalent) out to the escape hatch; it isn't a judgment call.

**Finding: both relational forms translated, and the descendant-marker form
does work.** Card exercises two distinct shapes of `:has()`:

- Same-shape as Button's proven form — `:has([data-gsxui-slot-X])` where X
  is a **descendant** marker, not `>svg`. This was the brief's flagged
  unproven case, and it does compile and apply correctly through Tailwind's
  arbitrary-value variant syntax:

  ```css
  /* was: :where([data-gsxui-slot-card]):has([data-gsxui-slot-card-footer]) { pb-0 } */
  .gsxui-recipe-card {
    @apply … has-[[data-gsxui-slot-card-footer]]:pb-0;
  }
  /* was: :where([data-gsxui-slot-card-header]):has([data-gsxui-slot-card-action]) { grid-cols-[1fr_auto] } */
  .gsxui-recipe-card-header {
    @apply … has-[[data-gsxui-slot-card-action]]:grid-cols-[1fr_auto];
  }
  ```

- A same-element compound selector on a **consumer-supplied class**, not a
  marker: `:where([data-gsxui-slot-card-header].border-b) { pb-4 }` — this
  isn't Card reacting to a descendant, it's Card's header reacting to
  itself also carrying `.border-b` (site/examples/card/compound.gsx passes
  `class="border-b"` to `CardHeader`). This is not the `:has()` case at all;
  it translates as an arbitrary same-element variant:

  ```css
  .gsxui-recipe-card-header {
    @apply … [&.border-b]:pb-4;
  }
  ```

  Watch for this distinction on later components — a relational rule keyed
  on a class rather than a marker is common (shadcn's border-b/border-t
  pattern recurs on more than just Card) and needs `[&.foo]:`, not
  `has-[...]:`.

`registry/styles/maia/<component>.css`: if Maia hasn't diverged from Nova on
this component yet (there is exactly one rule set in `default.css`, not
two), the file is byte-identical to Nova's. That's expected, not a bug —
conformance is checked in both directions and a missing slot fails loudly
naming it.

### Step 5 — add accessor calls to the canonical

```bash
cp ui/<component>.gsx registry/canonical/<component>.gsx
```

Change `package ui` to `package canonical`. Add the binding in
`registry/canonical/<component>_recipe.go` (not in `shapes/` — see the
design doc §3.1 for why shapes must stay a leaf package). Add
`class={ <component>.<Slot>() }` to every marked element, following
`button.gsx`'s attribute order: `class={...}` before `{ attrs... }`, marker
attribute last.

**Finding: this is where the real work is, not Step 4.** ~47 of 53
components have no `class` attribute at all in `ui/`, Card included — you
are inventing a place for the class to go, not editing an existing one.

### Step 6 — generate

```bash
go run ./cmd/stylegen   # regenerates ui/<component>.gsx (and per-style
                         # copies under registry/generated/ and
                         # site/stylepreview/) with the RESOLVED classes —
                         # you do not hand-write ui/<component>.gsx's
                         # Tailwind classes yourself
go tool gsx generate
```

**Finding: `ui/<component>.gsx` is generated output, not source, once a
component is on the axis.** `cmd/stylegen` resolves every accessor call in
`registry/canonical/<component>.gsx` against the active style's recipe CSS
and overwrites `ui/<component>.gsx` (plus `registry/generated/{nova,maia}/`
and `site/stylepreview/{nova,maia}/`) with literal classes. Run stylegen
*before* `gsx generate` — `gsx generate` needs the resolved `.gsx` to
compile, and needs the `<component>Recipe` type stylegen also generates
into `registry/canonical/<component>_recipe.gen.go` to exist first (if you
run `gsx generate` before `stylegen`, `registry/canonical/<component>.gsx`
fails with "undefined: <component>Recipe" and the generator emits a
`GSX_GENERATION_FAILED` stub for the whole package — harmless once you
re-run stylegen and regenerate, but confusing if you don't expect it).

### Step 7 — remove the migrated rules from `default.css`

Delete exactly what Step 2 classified as migrated. For anything genuinely
untranslatable, move it into the `@layer utilities` block (near the end of
`default.css`) with a comment explaining why, and make sure the selector is
**unwrapped** (no `:where()`) so its specificity clears the `(0,1,0)` floor
the layer gate requires — a `:where()`-wrapped rule left in `@layer
utilities` scores `(0,0,0)` and silently loses to any plain utility class in
the same layer.

### Step 8 — check the audit allowlist

`make audit` includes four rg patterns that specifically excluded
`button.gsx` (`ui -g '*.gsx' -g '!button.gsx'`) because Button was the only
component with a `class` attribute in `ui/`. **These four exclusions need
your component added too**, or `make audit` fails on the very
`class={...}` lines Step 6 just generated:

```
@! rg -n '^[[:space:]]+class=' ui -g '*.gsx' -g '!button.gsx' -g '!card.gsx'
@! rg -n '^[[:space:]]*<[^>]*class=' ui -g '*.gsx' -g '!button.gsx' -g '!card.gsx'
```

(Two more patterns at lines 76-77 in the Makefile exclude `button.gsx` for
`data-slot`/`group-`/`peer-` patterns — Card didn't trip these, but check
whether your component does.) This step is **entirely missing from the
brief's Step 7 "generate and check"** — the brief lists `stylegen --check`
and `make audit` as if `make audit` would just pass. It doesn't, the first
time, and the fix isn't "weaken the gate" (the brief's own warning about
`make audit` failing due to a dead rule doesn't apply here — this failure
is the allowlist being stale, not a real dead rule).

### Step 9 — build and check

```bash
go build ./... && go test ./... -count=1
go run ./cmd/stylegen --check   # exit 0
make audit                       # after Step 8's Makefile fix
gofmt -l .                       # empty
```

**Finding: existing hand-written Go tests pin the pre-migration render, and
migration legitimately breaks them.** `ui/<component>_test.go` may contain
an exact-string pin (`TestCardPinned` did:
`` `<div data-gsxui-slot-card>Content</div>` ``, no class attribute) and a
substring check that assumes the caller's class is the *entire* class
attribute (`` `class="py-8"` ``, not just `py-8`). Both need updating:
the pin's `want` string gets the newly-resolved class attribute added
verbatim (copy it from the actual generated output, don't hand-compute it);
the substring check switches from an anchored `class="X"` match to a bare
`strings.Contains(got, "X")`, since the recipe's classes are now merged in
before the caller's. This is a legitimate test update, not a weakening —
the computed style is what the sweep verifies, and these Go tests were
never asserting computed style, only literal markup that the migration
deliberately changes.

### Step 10 — sweep

```bash
make sweep-compare
```

Every diff is a finding, not a formality. Investigate and either fix or
write down why it's expected. Card's sweep was clean on the first attempt
(`no computed-style differences`) — the rule inventory in Step 2 was
complete enough that nothing needed a second pass.

### Step 11 — add regression pins

Add pins to `jstest/specs/layer-precedence.spec.ts` for properties the
migrated component would visibly lose if a future change got its layers
wrong — follow the Carousel/InputGroup precedent exactly (a `page.goto` to
a fixture that exercises the property, a `getComputedStyle` read, a hard
`toBe`). For Card: its border radius (needs a fixture with an actual card —
`/x/card/compound`) and its header's `.border-b` bottom padding/border
(same fixture, since `compound` is the one example that passes
`class="border-b"`).

**Finding: guess the expected pin value from the actual computed style, not
from the source utility's usual meaning.** `rounded-xl` reads as "12px" by
Tailwind's default scale, but this theme's `--radius` token makes it 14px
computed — a plausible-looking guessed value fails the pin test
immediately and for a boring reason (theme token, not a migration bug).
Run the test once, read the real number, then pin that number.

**Finding: the brief's stated Playwright baseline (329 passed / 0 failed)
was stale for this tree.** The actual pre-migration count at the starting
commit was 331 (verified by stashing all migration changes and re-running
the full suite). Re-verify the baseline count yourself rather than trusting
a number in a brief — it can drift as other work lands on the branch. Two
new pins made the final count 333; that's `<pre-migration count> + <pins
you added>`, not a fixed number.

### Step 12 — commit

```bash
git add -A
git commit -m "feat: migrate <Component> to the slot axis"
```

## What actually needed the escape hatch

One rule: the Card header's `container-name: card-header` /
`container-type: inline-size`. These are plain CSS properties with no
Tailwind utility form, so a recipe stylesheet (which only accepts `@apply`)
cannot carry them at all — this isn't a judgment call, `stylegen` refuses to
parse the file otherwise. They moved to `assets/css/styles/default.css`'s
`@layer utilities` block, unwrapped, with a comment. Nothing in the
catalogue currently reacts to this named container (no `@sm/card-header:`
or similar variant exists anywhere), so it's inert upstream setup being
preserved rather than an active concern — but it still had to move, since
leaving it under the deleted `@layer components` selector would make it
dead code the layer gate should (and does) object to.

## Dimensions on the root slot (Badge, Alert)

Badge and Alert added the first **dimension** (`variant`) to the model —
Button had dimensions but on a single-slot component where "root" and "the
only slot" are the same thing; Card had slots but no dimensions. Badge is
single-slot-with-dimension like Button; Alert is a compound component where
the dimension lives on the root slot only (`title`/`description` stay plain
`Base: true` slots, no `Dimensions`).

**Finding: derive the value list from `default.css`, but treat the style
contract as authoritative when it disagrees, and CONSULT UPSTREAM to settle
the disagreement — don't average the two or trust whichever list is easier to
read off the CSS.** For Badge specifically:

- The brief's illustrative shape (`task-2-brief.md`) quoted 4 values —
  `default, secondary, destructive, outline` — explicitly flagged as
  illustrative, and it was: incomplete.
- Reading `default.css`'s `[data-gsxui-slot-badge][data-variant="…"]` base
  rules directly gives 5 values with their own non-hover rule block —
  `default, secondary, destructive, outline, link`. `"ghost"` appears
  exactly once, sharing a hover-only compound selector with `outline`
  (`:is(a[…][data-variant="outline"]):hover, :is(a[…][data-variant="ghost"]):hover`)
  and has **no** dedicated base rule of its own.
- `internal/stylecontract/contracts_primitives.go`'s Badge axis lists 6
  values, including `ghost`.
- **This looked like a bug** (a declared variant resolving to unstyled
  except on hover) and was reported as NEEDS_CONTEXT rather than guessed.
  Checking upstream (`apps/v4/registry/new-york-v4/ui/badge.tsx` in a
  shadcn checkout) resolved it: ghost's ENTIRE upstream definition is
  `[a&]:hover:bg-accent [a&]:hover:text-accent-foreground` — hover-only and
  only when the badge renders as an anchor, by design (a ghost badge is
  meant to look like plain text until hovered as a link). gsxui's
  `default.css` is faithful to upstream, not missing a rule. **The contract's
  6 values were correct; use them.** Alert's contract and CSS agreed exactly
  (`default, destructive`) — no ambiguity, no upstream check needed.

**Rule for later components:** when the contract and `default.css` disagree
on a dimension's value list, that is not automatically a bug to fix or a
judgment call to make locally — check the corresponding upstream shadcn
component file (`/Users/jackieli/personal/shadcn-ui/apps/v4/registry/new-york-v4/ui/<component>.tsx`
in this environment) before reporting OR before trusting either source. The
contract wins once upstream confirms it; `default.css`'s own quirks
(hover-only variants, dead-looking rules) are frequently faithful ports of an
intentional upstream design, not migration bugs.

**A fourth relational form: `[a&]:hover:`.** gsxui writes
`:is(a[data-gsxui-slot-X][data-variant="Y"]):hover { @apply … }` for
several/all of Badge's variants (every one of `default, secondary,
destructive, outline, link` has one, plus the shared `ghost` case). This is
exactly upstream's `[a&]:hover:` Tailwind variant, applied only when the
component renders as an anchor. Translate it directly:

```css
/* was: :is(a[data-gsxui-slot-badge][data-variant="destructive"]):hover { bg-destructive/90 } */
.gsxui-recipe-badge-variant-destructive {
  @apply … [a&]:hover:bg-destructive/90;
}
```

It compiles cleanly under this repo's Tailwind v4 (verified directly with
the `tailwindcss` CLI the same way `jstest/support/compiled-css-audit.test.ts`
does, before committing to the syntax) — confirm this yourself on a new
component rather than trusting that it will just work, the same caution the
brief already asked for on the descendant-marker `:has()` form.

**Note: this form can be write-only.** Badge's `Badge()` component never
renders an `<a>` (no `href` prop, unlike Button) — so `[a&]:hover:` is
correctly copied but currently unreachable in this codebase, same as
upstream's own `asChild`-only Slot path. That's not a reason to drop the
utility; the recipe stylesheet has to remain faithful to `default.css`
regardless of whether every consumer happens to exercise every state.

**The `<Slot><Dimension>` accessor looked exactly like Button's, because both
declared the dimension on the root (`Name: ""`) slot:**

```go
func (r badgeRecipe) Root() string             { return r.c.SlotClass("") }
func (r badgeRecipe) Variant(value string) string { return r.c.SlotValueClass("", "variant", value) }
```

Alert's title/description slots, having no dimensions, only ever get a
`SlotClass` accessor (`Title()`, `Description()`) — same as any of Card's
plain slots. Nothing new needed in `internal/recipe` for a dimension that
happens to live on a non-root-but-single slot; the mechanism doesn't
distinguish "root" from "the component's only dimensioned slot" at all — it
was already generic over slot name.

**The generated `default:` arm.** For Badge, `registry/generated/nova/badge.gsx`'s
switch produces:

```go
switch variant {
case "secondary": "bg-secondary text-secondary-foreground [a&]:hover:bg-secondary/90"
case "destructive": "bg-destructive text-contrast focus-visible:ring-destructive/20 dark:bg-destructive/60 dark:focus-visible:ring-destructive/40 [a&]:hover:bg-destructive/90"
case "outline": "border-border text-foreground [a&]:hover:bg-accent [a&]:hover:text-accent-foreground"
case "ghost": "[a&]:hover:bg-accent [a&]:hover:text-accent-foreground"
case "link": "text-primary underline-offset-4 [a&]:hover:underline"
default: "bg-primary text-primary-foreground [a&]:hover:bg-primary/90"
}
```
confirming the `default:` arm carries `default`'s own utilities (matching
`TestResolveDefaultArmCarriesDeclaredDefault`'s pin for Button), not `""`.

## What actually needed the escape hatch (Badge, Alert)

Nothing. Every rule in both components' `default.css` blocks translated to
`@apply`-only utilities, including the `.dark` selectors (`dark:` variant)
and the parent-reaching Alert rule below. Unlike Card, no plain-CSS
properties without a Tailwind form turned up.

**A relational shape Card didn't have: a variant on one slot painting a
DIFFERENT slot.** Alert's destructive variant recolors its description child
(`:where([data-gsxui-slot-alert])[data-variant="destructive"]
:where([data-gsxui-slot-alert-description])`). This is the mirror image of
Card's descendant `:has()` form — Card reacts to a child's *presence*;
here the root reacts to its own *variant* and reaches down to paint a
specific descendant. `has-[]:` is the wrong tool (nothing is being detected,
something is being painted), so it's an ordinary arbitrary descendant
variant on the variant's own rule:

```css
.gsxui-recipe-alert-variant-destructive {
  @apply bg-card text-destructive [&_[data-gsxui-slot-alert-description]]:text-destructive/90;
}
```

## What the brief got wrong or left out

1. **Rule count.** The brief said "13 rules"; the actual count in
   `default.css` lines 200-235 is 10 distinct rule blocks. Not load-bearing,
   but don't trust a remembered count over a fresh grep.
2. **The stray at "2442" needed characterizing, not just finding.** The
   brief correctly flagged that a stray exists near that line, but didn't
   warn that it might not actually belong to your component. It doesn't for
   Card — it's a Calendar rule keyed on a Card marker as ancestor context.
   Finding a marker match is not the same as finding a rule to migrate.
3. **The audit allowlist update (this playbook's Step 8) is entirely
   missing from the brief.** The brief's Step 7 implies `make audit` just
   passes after generation; it doesn't, because four `rg` patterns in the
   Makefile hard-code `!button.gsx` as the only exclusion for a `class`
   attribute appearing in `ui/*.gsx`. Every future migration hits this.
4. **The Go test pins (Step 9 here) are undocumented fallout.** The brief
   never mentions `ui/<component>_test.go` at all, but an exact-string
   render pin from before the migration will fail the moment a `class`
   attribute is added, and a substring check anchored to `class="X"` will
   fail once classes are merged. Both are legitimate, expected updates, not
   scope creep — worth calling out explicitly so a future agent doesn't
   either skip Card's Go tests or, worse, "fix" the recipe to avoid touching
   them.
5. **The relational-form claim was directionally right but incomplete.**
   The brief flagged the descendant-marker `:has()` form as unproven (it
   proved out); it didn't anticipate the *other* relational shape Card
   actually has — a same-element compound selector on a consumer-supplied
   class (`.border-b`), which needs `[&.border-b]:`, not `has-[...]:`. Any
   later component with a similar "pass an extra class, get an extra rule"
   pattern (shadcn's border-b/border-t convention recurs) will hit this
   same distinction.
6. **The Playwright baseline count (329) doesn't match the tree.** Verified
   331 at the pre-migration commit. Re-derive it yourself; don't paste the
   brief's number into a report unchecked.
7. **(Task 2, Badge/Alert) The brief's illustrative dimension-value list was
   wrong, not just illustrative-and-approximate.** It quoted 4 Badge values;
   the real count, once the contract and upstream both agree, is 6. Don't
   spend the "derive from `default.css`" instruction as license to skip
   checking the contract AND upstream when the two disagree — `default.css`
   alone made `ghost` look like dead/broken CSS when it was a faithful,
   intentional hover-only port.
8. **Dimensions can disagree between the style contract and `default.css` for
   a legitimate design reason, not just as an error.** The brief's framing
   ("derive from `default.css`, don't trust a quoted list") is right for
   catching a stale brief, but insufficient on its own when `default.css`
   itself looks incomplete — the tiebreaker in that case is upstream shadcn,
   not a coin flip between contract and CSS.
9. **This task's shared-worktree hazard, not in the brief at all:** a
   concurrent agent splitting Sonner/Toast was working in this SAME working
   directory (not an isolated git worktree per agent), and left uncommitted
   changes to shared generated artifacts (`registry/generated/recipes.json`
   csv-adjacent files, `jstest/runtime-style-contract.json`,
   `internal/stylecontract/contract.go`, `internal/registry/registry_test.go`,
   `site/examples/style_contract_test.go`) sitting in the tree throughout this
   task. `go build`, `go test`, and the Playwright suite all ran against that
   mixed state — one Playwright test
   (`runtime-style-contract.spec.ts › real interactions cover the exact
   runtime-owned style contract`) failed for reasons entirely owned by that
   other agent's in-progress rename, unrelated to Badge/Alert. Future
   migrations sharing a worktree need to `git diff --stat` before touching
   anything and again before committing, to separate "my diff" from "their
   live WIP" — and commit only the former.

## Four small components in one pass (Label, AspectRatio, Separator, Kbd)

This wave proved the parallel protocol on cheap work (3-7 rules each) and
was the first to do more than one component per agent turn. It also
exercised the first hyphenated component end-to-end.

### The hyphenated case (AspectRatio) already had infrastructure — verify, don't build

The brief flagged AspectRatio as "the interesting one" and warned the
kebab-registry/camel-Go split had never been exercised end to end. In fact
`internal/stylegen/identifier.go` (`componentIdentifier`) and its dedicated
test (`internal/stylegen/identifier_test.go`, built around a synthetic
`input-group` shape) already existed and already covered this — a prior
wave must have built the mechanism generically before any hyphenated
component actually used it. Nothing needed writing in `internal/stylegen`;
the work was purely: declare `Component: "aspect-ratio"` in the shape, let
`GenerateAccessors` derive `aspectRatioRecipe`/`aspectRatio` itself, and
verify both halves hold simultaneously — `registry/canonical/shapes/aspect-ratio.go`
still says `Component: "aspect-ratio"`, `registry/canonical/aspect-ratio.gsx`
calls `aspectRatio.Root()`, and `registry/styles/{nova,maia}/aspect-ratio.css`
still emits `.gsxui-recipe-aspect-ratio`. **Lesson for later components:**
before assuming a brief's "this is unproven" is still true, grep for the
mechanism — a prior wave may have already generalized it past the specific
component that motivated it.

### A new failure mode: the residual-accessor-call check scans doc comments too

`internal/stylegen/resolve.go`'s final safety net (`residualAccessorCall`)
regex-matches the literal text `separator.Orientation(` against the
**entire formatted source**, comments included, not just code. Separator's
canonical doc comment described the accessor call in prose ("resolved to
concrete utilities by separator.Orientation()") and `go run ./cmd/stylegen`
failed with "resolved GSX still contains a recipe accessor call
\"separator.Orientation(\"" even though the actual `class={...}` call had
resolved correctly — the match was in the comment, not the leftover code.
**Rule for later components:** don't name an accessor call in
`receiver.Method(` form inside a canonical `.gsx` doc comment; describe the
mechanism instead ("resolved to concrete utilities at generation time").
Cheap to avoid, confusing to diagnose the first time (the error message
doesn't distinguish "found in a comment" from "found in unresolved code").

### A new relational form: reacting to context outside the component's own markers

Both Label and Kbd needed a shape `has-[...]:` and `[&.foo]:` don't cover:
a rule keyed on something that is neither a descendant/ancestor MARKER of
the same component nor a consumer-supplied class, but genuine surrounding
DOM state.

- Label: `:where([data-disabled="true"]) :where([data-gsxui-slot-label])`
  (an ancestor's own `data-disabled` attribute) and
  `:disabled ~ :where([data-gsxui-slot-label])` (a preceding sibling's
  pseudo-class) translate as `[[data-disabled=true]_&]:` and `[:disabled~&]:`
  respectively — Tailwind's arbitrary-variant `&`-placeholder syntax handles
  both an ancestor-descendant compound and a general sibling combinator.
- Kbd: `:where([data-gsxui-slot-tooltip-content]) :where([data-gsxui-slot-kbd])`
  (an ancestor from a DIFFERENT component) translates as
  `[[data-gsxui-slot-tooltip-content]_&]:`. The `.dark` variant of the same
  rule, `:where(.dark [data-gsxui-slot-tooltip-content]) :where([data-gsxui-slot-kbd])`,
  needs **no** separate `dark:` Tailwind prefix — `.dark` is already baked
  into the arbitrary selector itself (`[.dark_[data-gsxui-slot-tooltip-content]_&]:`),
  and prefixing it with `dark:` on top double-applies the condition
  (`&:is(.dark *) { .dark ... & { ... } }` in the compiled output — verified
  by compiling both forms with the Tailwind v4 CLI and diffing the result
  before choosing the syntax, same caution as the `[a&]:hover:` form in the
  Badge/Alert wave).

All of these compiled cleanly on the first attempt once checked directly
against the CLI — this repeatedly cheap verification step (never run
Playwright to check syntax; run `tailwindcss -i - --cwd . --silent` on a
throwaway file first) keeps paying for itself.

### Migrating a component X breaks an UNMIGRATED component that composes X — this recurred four times, not once

Card/Badge/Alert's playbook already covered "your shape disagrees with a
sibling's rule" (the Calendar stray). This wave hit a DIFFERENT and more
consequential shape of collateral damage: giving a component's own recipe
class real `(0,1,0)` specificity breaks an *unmigrated* consumer that
composes it directly (`<Label>`, `<Separator>` used as the underlying
element inside another component's `.gsx`), because the consumer's own
default.css rule targeting its own extra marker is still `:where()`-wrapped
at `(0,0,0)` and can no longer win the cascade regardless of source order.
This is not hypothetical or rare — it happened for every single relational
consumer both of these two components have:

- Label: `FieldLabel` (`field.gsx`) composes `<Label>` directly. Its own
  `[data-gsxui-slot-field-label] { flex w-fit gap-2 leading-snug }` rule
  contested Label's own `leading-none` on the `line-height` property.
- Separator: THREE composed consumers exist
  (`ButtonGroupSeparator`, `ItemSeparator`, `SidebarSeparator`), and TWO of
  the three had a real property conflict — `ButtonGroupSeparator`'s
  `bg-input`/`h-auto` vs Separator's own `bg-border`/`h-full`, and
  `SidebarSeparator`'s `w-auto`/`bg-sidebar-border` vs Separator's own
  `w-full`/`bg-border`. `ItemSeparator`'s own rule (`my-2`, a margin) and
  `FieldSeparator`'s (`absolute inset-0 top-1/2`, positioning) don't
  conflict with anything Separator itself sets, so those needed no change
  at all — **check each consumer's actual declared properties against the
  slot's own recipe utilities; don't assume every same-element compound
  rule needs moving just because one sibling does.**

**How this surfaces:** not from grepping default.css (a plain grep for the
migrated component's own marker doesn't show you the *other* component's
`.gsx` composing it) but from `go test ./...` itself —
`internal/stylegen`'s `layercheck_test.go` fails
(`TestCheckLayerPrecedenceAcceptsTheCurrentTree` and friends) naming the
exact selector, property, and consuming component: *"under style X, SELECTOR
applies PROPERTY in @layer components, but that element renders through
COMPONENT, whose own utilities win the layer ordering."* This is the
correct signal, not a false positive from stale test infrastructure. **Rule
for later components: after regenerating and building, run the full test
suite before moving on — this class of break only shows up there, not in
`stylegen --check` or a visual reading of default.css.**

**The fix is uniform:** move exactly the conflicting rule (not necessarily
the whole originally-contiguous block) into `@layer utilities`, unwrapped
(no `:where()`), with a comment naming which migration caused the fallout
and which properties actually conflict — same escape-hatch mechanism Card's
`container-name` rule already established, just triggered by a DIFFERENT
component's migration instead of the rule's own component's migration. Test
pins on the composing component (`ButtonGroupSeparatorDefaultPinned`,
`ItemSeparatorDefaultPinned`, `SidebarSeparatorPinned`,
`FieldLabelPinned`/`FieldSeparatorPinned`) then need the composed
component's now-real canonical class spliced into their `want` string —
exactly Card/Badge/Alert's Step 9 finding, just triggered transitively
through a compose relationship instead of directly.

### Test-pin plumbing: build one recipe-loading helper per component, reuse everywhere

Following `button_test.go`'s `novaButtonRecipe`/`canonicalButtonClass`
pattern for every migrated component (`canonicalLabelClass`,
`canonicalAspectRatioClass`, `canonicalSeparatorClass`,
`canonicalKbdClass`/`canonicalKbdGroupClass`) paid off immediately: every
transitive-fallout pin above (`FieldLabelPinned`, `ButtonGroupSeparatorDefaultPinned`,
`ItemSeparatorDefaultPinned`, `SidebarSeparatorPinned` inside
`sidebar_test.go`'s table, `FieldSeparatorNoChildrenPinned`/`WithChildrenPinned`,
and `forms_css_contract_test.go`'s `TestFormControlsExposeCSSOnlySlots`/
`TestFormControlCompositionTokenOrder`) just called the helper rather than
hand-computing or guessing the resolved utility string. `strings.Count(got,
class) != 1 || strings.Count(got, "class=") != 1` (not a bare `Contains`)
is the right assertion shape for "the caller class merged in exactly once,"
matching `assertButtonCallerAttrsOnce`'s existing pattern.

**`TestFormControlsExposeCSSOnlySlots` needed a new per-name exception, not
a blanket rule change.** That test's design ("only Button-composing parts
render a `class=` attribute; everything else in the forms family is
CSS-only") pre-dates any slot-axis migration touching a composed part other
than Button. `FieldLabel` (composes Label) and `FieldSeparator` (composes
Separator) both trip the blanket "must not render a class" check the moment
their underlying primitive migrates — the fix is one `case` per newly-migrated
composed part in the existing `switch`, asserting the exact canonical class
now legitimately renders, mirroring the pre-existing `InputGroupButton` case
for Button. Expect one new case per migrated component that some OTHER,
still-unmigrated component composes directly.

### Attribute-order surprise: gsx floats `class=` ahead of other attributes at render time

`AspectRatio`'s canonical source writes `style={...}` before `class={...}`
in the `.gsx` (following the existing style-then-class ordering already in
the file), but the actual rendered HTML puts `class=` FIRST, ahead of
`style=`. This is a gsx rendering behavior (class merging happens ahead of
other attributes), not a stylegen or migration bug — the fix is in the
test's expected string, not the source. Worth checking rendered attribute
order directly rather than assuming it matches source order when a
component's class attr isn't the first attribute in the tag.

### Doing four components in one pass: real time savings, one real risk

Investigation reuse was substantial — reading `internal/recipe/shape.go`,
`internal/stylegen/accessors.go`, and Button/Badge/Card as worked examples
once covered all four; the Tailwind-CLI-verification habit and the
recipe-test-helper pattern were written once (for Label) and then just
copied per component. The one real risk: **shared generated/registry files
(`registry/canonical/shapes/shapes.go`, `registry/generated/recipes.json`)
accumulate edits across all four components in the working tree** even
though the four are logically independent and must land as four separate,
individually-buildable commits. Committing straight from an "all four done"
working tree bakes later components' registrations into the first commit's
diff. The only reliable fix found: temporarily move the not-yet-committed
components' shape/recipe/gsx/css files out of the tree (a scratch stash
directory), truncate the shared files back to the state that existed before
this wave started plus only the components already committed, rerun
`stylegen`/`gsx generate` for real (not by hand-editing the diff) so the
shared generated files regenerate correctly for that partial set, verify
gates and tests pass at THAT checkpoint, commit, then restore the next
component's files and repeat. This is mechanical but not optional — a
diff-splitting shortcut (hand-trimming the JSON or the map literal) risks
committing a shared file in a state `stylegen` never actually produced.
