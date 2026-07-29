# Theme system and configurator roadmap

**Status:** Phase 1 shipped; later phases are uncommitted possibilities.

This is dependency ordering, not a release promise. Work starts only when a
phase is explicitly prioritized. A later phase does not imply that gsxui
will ship multiple component styles.

The architecture and acceptance criteria live in
[`docs/superpowers/specs/2026-07-26-css-only-theme-system-design.md`](superpowers/specs/2026-07-26-css-only-theme-system-design.md).

## Architectural outcome

```text
canonical component markup + behavior
                  │ stable data-gsxui-slot-<name> contract
                  ▼
          foundation mechanics
                  +
             theme values
                  +
       one CSS presentation style
                  ▼
          caller utility overrides
```

Production projects install one unscoped style. The editor may build scoped
copies for preview. Structural differences are separate components or a
future structural base, never style-specific `.gsx` templates.

**Exception: Button.** The diagram above still holds for every component
except Button. Button's canonical source is compiled against a typed style
recipe at generation time, so `ui/button.gsx` ships concrete Tailwind
utilities baked into its generated `.gsx` source rather than a
`.gsxui-recipe-button*` class consumed from a components-layer stylesheet.
The rest of the component set still follows the CSS-layer path described
here. See
[`docs/superpowers/specs/2026-07-29-typed-recipe-model-design.md`](superpowers/specs/2026-07-29-typed-recipe-model-design.md)
for the compiled-presentation architecture and its invariants.

## Phase 1 — correct the styling boundary (shipped)

Shipped as one CSS-only boundary:

- Introduce the typed slot/state contract.
- Replace generic `data-slot` with one `data-gsxui-slot-<name>` presence
  attribute per semantic role.
- Split the CSS source of truth into index, foundation, theme, and default
  style assets.
- Move current Nova presentation out of every `.gsx` file and into
  `@layer components`. (Superseded for Button only: Button's presentation
  is now compiled back into concrete utilities in generated `.gsx` source —
  see the exception noted above and the 2026-07-29 typed recipe model spec.)
- Leave only semantics, accessibility, behavior state, stable hooks,
  dynamic mechanism values, and caller classes in components.
- Make the site, browser harness, and `gsxui init` consume the same assets.
- Remove the old combined/copy-based styling path in the same change.

The shipped gate covers the exact component/style contract, real finite
examples, runtime-created states, foundation-only behavior, light/dark and
responsive visual baselines, and caller utility precedence.

## Phase 2 — make CSS a checked API

**Relative size:** medium to large.

- Add parser-backed CSS and selector validation.
- Check emitted component hooks against the typed contract in both
  directions.
- Require coverage for every declared style slot/state.
- Reject unknown hooks, global selector leakage, literal component palettes,
  invalid `@apply`, and undeclared state values with source locations.
- Generate and drift-check the serialized editor contract and scoped
  preview CSS.

**Exit gate:** deleting or misspelling a required rule produces a precise CI
failure, and a style cannot compile while escaping the component contract.

## Phase 3 — one preset and safe CLI application

**Relative size:** medium.

- Add the versioned `gsxui.preset.json` schema and pure resolver.
- Generate deterministic `theme.css` and select `style.css` from the same
  resolved preset the editor uses.
- Add init/apply flows for full, theme-only, and style-only application.
- Record managed hashes and refuse to overwrite user-modified files without
  explicit `--overwrite`.
- Add transactional writes, rollback, and interrupted-transaction recovery.
- Add canonical preset/share-code import and export.

**Exit gate:** all inputs and conflicts are checked before mutation; no-op
apply is byte-identical; every tested failure rolls back byte-for-byte.

## Phase 4 — replace the current editor

**Relative size:** large.

- Build visual controls for semantic colors, radius, and typography while
  retaining exact values.
- Use one immutable preset state with undo, redo, deterministic shuffle,
  reset, local persistence, import/export, and sharing.
- Render the real component example registry in an isolated iframe.
- Cover full component browsing, light/dark, responsive sizes, keyboard
  interaction, and preview failure states.
- Add new-project, existing-project, theme-only, style-only, and manual-file
  handoff.
- Make controls, dialogs, validation, clipboard feedback, and the preview
  keyboard-accessible.

**Exit gate:** edit → share/export → import → CLI apply round-trips to the
same canonical preset and CSS, with browser tests covering state history and
the full handoff.

## Phase 5 — prove the extension point if a real style arises

**Relative size:** large per style; dormant without a concrete candidate.

- Author a genuinely distinct second CSS pack over the unchanged contract.
- Run the complete component, accessibility, behavior, computed-style, and
  light/dark/responsive visual matrix.
- Prove applying it changes only `style.css`.
- Add it to the editor's generated scoped bundle.
- Expose the style picker only after both packs pass every gate.

**Exit gate:** the second style requires no `.gsx`, JavaScript, foundation,
theme, component API, or dependency fork.

If that proof fails, the candidate is structural work and is designed as a
component or base instead. The architecture remains valid with one style.
