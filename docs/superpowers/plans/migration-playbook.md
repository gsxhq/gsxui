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
