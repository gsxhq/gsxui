# Inline style generation and theme editor Button pilot

**Date:** 2026-07-28

**Status:** Approved pilot design

**Review gate:** No component after Button adopts this architecture until the
pilot's authored source, generated source, editor UX, CLI ownership behavior,
and measured render cost have been reviewed.

## 1. Purpose

gsxui currently has two proven pieces that do not yet form one product:

- GSX can merge Tailwind classes with the project-configured
  `class_merger`, and gsxui supplies a genuine Tailwind-aware implementation
  through `github.com/gsxhq/gsxui/merge.Merge`.
- The CSS-only theme branch separates semantic theme values from component
  presentation and exposes a basic browser theme editor.

The remaining developer-experience gap is source ownership. A component copied
into a user's project is easier to understand and customize when its concrete
Tailwind presentation is visible beside its markup and behavior. The current
separate component stylesheet makes a copied `.gsx` file incomplete on its own.

This pilot tests a shadcn-style delivery model adapted to GSX:

1. maintain one style-neutral Button structure;
2. maintain Nova and Maia Button recipes separately;
3. compile one selected recipe into concrete Tailwind classes in copied GSX;
4. preview and export that selection through the public web theme editor; and
5. stop after Button for an explicit architecture review.

This is a vertical product slice, not a compiler-only experiment. The result
must be visible and usable from `/theme` and through a throwaway consumer
project.

## 2. Relationship to the CSS-only design

`2026-07-26-css-only-theme-system-design.md` remains the baseline for:

- semantic light and dark theme tokens;
- the distinction between behavior, foundation mechanics, theme values, and
  component presentation;
- parser-backed preset validation;
- iframe preview isolation;
- deterministic import, export, and share transport; and
- transactional CLI mutation.

This pilot overrides that design for Button presentation and delivery only:

- Button presentation is baked into copied `.gsx` source instead of remaining
  in a vendored production `style.css`;
- style selection therefore regenerates Button source;
- the browser editor still derives its preview from the same authored recipe,
  but a consumer does not need the recipe stylesheet; and
- a full style migration is an explicit source overwrite, not a CSS-file swap.

No repository-wide reversal is approved by this document. Other components
remain on their current styling path until the Button review concludes.

## 3. Goals

The pilot is successful when:

1. `/theme` visibly offers Nova and Maia.
2. Changing style updates an isolated, interactive Button preview.
3. Light/dark theme edits and radius edits affect both styles.
4. One canonical Button structure produces concrete Nova and Maia `.gsx`
   source without duplicating markup or behavior.
5. The copied Button contains readable Tailwind utilities and no internal
   recipe tokens or styling support helpers.
6. Caller classes override generated defaults through `merge.Merge`.
7. Preset JSON, share codes, CLI handoff commands, and copied Button source all
   agree on the selected style.
8. Existing local Button source is never silently overwritten.
9. Applying Maia is rejected when installed components outside the Button
   pilot would create a falsely mixed-style project.
10. The whole slice is verified in real browser and throwaway consumer builds.

## 4. Non-goals

The pilot does not:

- migrate a second component;
- promise Nova/Maia support for the complete component catalogue;
- add a third style or alternative primitive base;
- intelligently merge upstream recipe changes with user-edited component
  source;
- retain a runtime style switch in consumer applications;
- implement undo/redo, deterministic shuffle, responsive viewport controls, or
  the all-component editor browser;
- edit or apply typography/font choices;
- infer styling intent from arbitrary concrete Tailwind utilities; or
- introduce regex-based CSS or GSX rewriting.

## 5. Architectural overview

```text
canonical Button.gsx
  semantic markup, behavior, state, recipe tokens
                     │
                     ├──────────────┐
                     │              │
              Nova recipe      Maia recipe
                     │              │
                     └──────┬───────┘
                            │
                 parser-backed resolver
                            │
             selected concrete Button.gsx
                            │
             ┌──────────────┴──────────────┐
             │                             │
      gsxui add/apply                 generated fixtures
             │                             │
      user's ui/button.gsx        editor + consumer tests
```

The recipe resolver runs before consumer GSX generation. GSX's normal compiler
then wires `merge.Merge` into the concrete class expressions. Style selection
is build/install-time; caller conflict resolution remains render-time.

The pilot's authored and generated sources have one unambiguous location:

```text
ui/button.gsx                         canonical tokenized Button
registry/styles/nova/button.css       authored Nova recipe
registry/styles/maia/button.css       authored Maia recipe
registry/generated/nova/button.gsx    committed generated consumer source
registry/generated/maia/button.gsx    committed generated consumer source
```

The embedded registry includes `registry/generated`; `gsxui add` does not
transform the canonical file opportunistically in a consumer project.

## 6. Canonical Button source

### 6.1 Responsibilities

The canonical Button owns:

- `<button>` versus `<a>` selection;
- public parameters and defaults;
- disabled semantics;
- accessible and behavioral attributes;
- stable slot, variant, and size markers;
- caller attribute fallthrough; and
- the mapping from public variant/size values to semantic recipe roles.

It must not contain Nova- or Maia-specific concrete utilities.

### 6.2 Recipe tokens

Recipe tokens are internal, namespaced class tokens. The pilot uses one token
per complete semantic role, for example:

```text
gsxui-recipe-button
gsxui-recipe-button-variant-default
gsxui-recipe-button-variant-outline
gsxui-recipe-button-size-default
gsxui-recipe-button-size-icon
```

Rules:

1. Tokens are whole class tokens, never substrings.
2. Tokens are string literals inside GSX class expressions.
3. Tokens cannot be assembled dynamically.
4. Every token used by canonical Button must exist exactly once in every
   validated style.
5. Every Button token declared by a style must be used.
6. Tokens never appear in consumer output.
7. Unrelated caller classes and invariant non-recipe classes are preserved.

The canonical component uses an inline `switch` for variant and size mapping so
the generated consumer file does not need `variantClass`, `sizeClass`, or
another private styling helper:

```gsx
class={
  "gsxui-recipe-button",
  switch variant {
  case "destructive": "gsxui-recipe-button-variant-destructive"
  case "outline": "gsxui-recipe-button-variant-outline"
  case "secondary": "gsxui-recipe-button-variant-secondary"
  case "ghost": "gsxui-recipe-button-variant-ghost"
  case "link": "gsxui-recipe-button-variant-link"
  default: "gsxui-recipe-button-variant-default"
  },
  switch size {
  case "xs": "gsxui-recipe-button-size-xs"
  case "sm": "gsxui-recipe-button-size-sm"
  case "lg": "gsxui-recipe-button-size-lg"
  case "icon": "gsxui-recipe-button-size-icon"
  case "icon-xs": "gsxui-recipe-button-size-icon-xs"
  case "icon-sm": "gsxui-recipe-button-size-icon-sm"
  case "icon-lg": "gsxui-recipe-button-size-icon-lg"
  default: "gsxui-recipe-button-size-default"
  }
}
```

Button's `<a>` and `<button>` branches receive identical generated class
expressions. A structural test pins that equality. The pilot deliberately
accepts visible repetition over hiding styling in private vendored helpers;
the review will decide whether that remains readable in the concrete output.

### 6.3 Persistent markers

`data-gsxui-slot-button`, `data-variant`, and `data-size` remain in rendered
HTML. They are useful semantic and composition hooks even though the selected
Button presentation is inline. Bare slot markers rely on GSX's presence
semantics and pass through attribute bags without `"true"`/`"false"` values.

JavaScript still reflects behavior state through attributes, never by toggling
presentation classes.

## 7. Style recipe sources

Nova and Maia each author a Button recipe map as CSS using exact, simple class
selectors and `@apply`:

```css
@layer components {
  .gsxui-recipe-button {
    @apply ...;
  }

  .gsxui-recipe-button-variant-outline {
    @apply ...;
  }
}
```

This format serves two consumers:

- the resolver extracts the ordered Tailwind utility string for generated GSX;
- the editor build scopes the same recipe rules for internal preview use.

The supported pilot grammar is intentionally strict:

- one exact recipe class is the subject of each rule;
- recipe selectors cannot be combined with element, ID, attribute, descendant,
  pseudo-class, or pseudo-element selectors;
- declarations consist of one or more `@apply` directives;
- ordinary CSS declarations are rejected for recipe rules;
- duplicate recipe keys are errors;
- non-recipe rules are not accepted in the Button recipe file; and
- parsing uses a real CSS parser with source positions.

State such as hover, focus-visible, dark mode, ARIA state, child SVGs, and
ancestor composition is expressed as Tailwind variants inside `@apply`, not as
extra CSS selectors. This ensures the exact utility can be emitted into
consumer GSX.

Nova preserves the current compact Button treatment. Maia is deliberately
different: more generous heights and padding, softer rounding, and its own
border/shadow treatment while using the same semantic color tokens and public
variant/size axes.

## 8. Parser-backed GSX transformation

The resolver parses canonical Button through the GSX parser and walks class
expression string literals. It does not rewrite source text with regular
expressions.

For every literal:

1. tokenize it with HTML class whitespace semantics;
2. replace each whole recipe token with the selected ordered utility list;
3. retain all non-recipe tokens in their authored position;
4. reject a recipe-looking token missing from the selected style;
5. preserve source structure, comments, imports, public signature, and
   formatting; and
6. format the resulting GSX with the canonical GSX formatter.

The output is deterministic. The same canonical source, style recipe, schema,
and formatter version produce byte-identical `.gsx`.

The transform is narrow by design. It does not interpret or normalize arbitrary
Tailwind utilities, alter caller expressions, or attempt to reverse a concrete
component back into recipe tokens.

## 9. Generated artifact and ownership model

### 9.1 Registry artifacts

Nova and Maia Button output are committed generated registry artifacts. CI
regenerates them and rejects drift. They are the exact bytes
`gsxui add button` offers after package/import rewriting.

Generated `.x.go` continues to come only from the selected `.gsx`; it is never
hand-edited.

### 9.2 Project state

`gsxui.preset.json` is the durable selected preset. Its pilot schema contains:

- `schemaVersion`;
- `style`, restricted to `nova` or `maia`;
- radius;
- the complete ordered light/dark semantic theme maps.

`gsxui.json` continues to own project paths and records managed-artifact hashes.
It does not duplicate preset values.

### 9.3 Add

`gsxui add button`:

1. loads and validates the project preset;
2. resolves canonical Button with the selected style;
3. rewrites only package/import paths;
4. shows a focused conflict if the target exists and differs;
5. writes only after all validation succeeds; and
6. runs GSX generation.

`--diff` is read-only and shows generated selected source against the local
file. `--overwrite` is required to replace a modified file.

### 9.4 Apply

A full `gsxui apply --preset ...` is explicit replacement:

1. resolve and validate the entire incoming preset;
2. discover installed managed components;
3. reject Maia if any discovered styled component other than Button would
   remain unsupported;
4. list files and axes that will change;
5. warn that managed component source will be overwritten;
6. require confirmation unless an explicit non-interactive confirmation flag
   is supplied;
7. reject locally modified managed files unless `--overwrite` is supplied;
8. update Button, preset, and theme artifact transactionally; and
9. leave the project byte-identical on failure.

`--only theme` retains the selected style and does not reinstall Button.
`--only style` retains semantic theme values but still follows the full
component overwrite checks.

There is no automatic three-way merge. The CLI tells users to commit or stash
before replacement.

## 10. Web editor vertical slice

### 10.1 Browser/project boundary

The public `/theme` route is a preset builder. It never writes directly to a
user's filesystem. It produces:

- canonical preset JSON;
- a share code and share URL;
- theme CSS for manual use;
- downloads; and
- exact `gsxui init`/`gsxui apply` commands.

Only the CLI mutates a project.

### 10.2 State

One validated preset is the editor's resolved state. Focused invalid field text
may exist as draft UI state, but it does not update:

- the iframe;
- exported JSON or CSS;
- share codes; or
- handoff commands.

The browser uses native CSS value parsing for immediate feedback and a
generated schema for allowed fields, ordering, and styles. The Go resolver is
authoritative when the CLI consumes the preset.

### 10.3 Controls

The pilot editor exposes:

- a visible Nova/Maia style picker with descriptions;
- light/dark mode;
- every required semantic color value;
- radius;
- exact CSS-value text entry;
- reset to the selected built-in preset;
- preset JSON import/export;
- compatible theme CSS import/export;
- share-code copy;
- generated init/apply command copy; and
- visible success, validation, clipboard-failure, and preview-error feedback.

Undo/redo, shuffle, responsive viewport controls, and catalogue browsing are
post-review features.

### 10.4 Preview

The preview is a same-origin iframe using the site's real document head,
Tailwind build, theme bootstrap, and behavior modules.

Its Button matrix includes:

- every variant;
- every size;
- text and icon content;
- disabled Button;
- enabled link Button;
- focus-visible and invalid examples;
- Button Group adjacency; and
- representative composed uses where Button receives additional classes.

Changing style changes the iframe's selected recipe. Changing theme values
injects semantic variables into the iframe without rebuilding Button. Light and
dark use the same mechanism as production.

The editor may use scoped internal recipe CSS to switch instantly, but that CSS
is generated from the exact recipe map used for concrete GSX. It is not a
second hand-authored implementation. End-to-end browser fixtures also build
and render concrete generated Nova and Maia Button source so equivalence is
measured rather than assumed.

The editor labels Maia as a Button pilot and does not imply catalogue-wide Maia
support.

## 11. Preset transport

Canonical JSON:

- uses schema field order and theme-token order;
- ends with one newline;
- rejects duplicate keys, missing keys, unknown keys, unknown styles, and
  unsupported future schema versions; and
- round-trips byte-identically after normalization.

The initial share transport remains:

```text
gsxui:v1:<base64url canonical preset JSON>
```

The editor can load a share code from its URL and reports field-specific errors
without partially applying invalid state.

## 12. Error handling and transactions

Errors carry the source layer and location where possible:

- recipe CSS path, selector, and line/column;
- canonical GSX path and token location;
- preset field path;
- target project file; or
- preview document/route.

The resolver produces all intended bytes before filesystem mutation. CLI writes
use a transaction directory, same-filesystem replacement, rollback snapshot,
and recovery journal. Conflict detection happens before the first rename.

An interrupted or failed invocation must either complete or restore the
previous project before another mutation begins.

## 13. Verification

### 13.1 Recipe and transform

- CSS-parser tests for valid recipes and every rejected grammar form.
- Exact set equality between canonical Button recipe tokens and each style.
- Unknown, missing, duplicate, dynamically assembled, and unused token tests.
- Golden Nova and Maia Button `.gsx`.
- Generated output contains no `gsxui-recipe-` token.
- Output parses, formats idempotently, generates, and passes `gopls check`.
- `<a>` and `<button>` class expressions resolve identically.

### 13.2 Runtime

- Existing Button semantics and accessibility tests pass for both generated
  styles.
- Caller `rounded-none`, height, padding, and color utilities replace generated
  defaults through `merge.Merge`.
- Composition classes survive and conflicts resolve last-wins.
- Render benchmarks compare the current Button path with generated inline
  classes. The review records absolute and relative cost; it does not hide a
  material regression.

### 13.3 Editor

- Style picker changes only style and retains theme state.
- Theme edits change both Nova and Maia previews.
- Invalid drafts do not change resolved output.
- Light/dark, reset, import, export, share, and handoff commands round-trip.
- Iframe theme/style synchronization works before first paint and after lazy
  load.
- Clipboard failure exposes a manual-copy fallback.
- Preview route errors produce a visible recoverable state.

### 13.4 Browser

- Computed-style assertions pin representative geometry, radius, borders,
  focus, disabled state, and caller overrides.
- Visual snapshots cover Nova/Maia × light/dark.
- Keyboard and link behavior use the real Button implementation.
- Concrete generated-consumer fixtures match the editor recipe preview for
  pinned properties.

### 13.5 CLI and consumer

- Init/add use the selected style.
- `--diff` is non-mutating.
- Existing unmodified, modified, overwrite, no-op, unsupported-component, and
  injected-failure cases are pinned.
- Theme-only apply preserves Button bytes.
- Style-only apply preserves theme values.
- A throwaway Nova project and Maia Button-only project both generate, build,
  render, and pass focused browser checks.

The authoritative repository gate remains `make ci` against the required GSX
core version.

## 14. Review gate

When the vertical slice passes, review pauses before any generalization. The
review answers:

1. Is canonical tokenized GSX understandable to maintainers?
2. Are Nova and Maia recipe files easier to maintain than duplicated GSX?
3. Is generated Button source the code users would want to own and edit?
4. Is inline switch repetition preferable to private styling helpers?
5. Does the editor accurately communicate style versus theme?
6. Are add/diff/apply ownership rules safe and unsurprising?
7. Is render-time Tailwind merging cost acceptable with concrete defaults?
8. Did any CSS-only capability prove more valuable than source locality?

Only an explicit approval after that review authorizes:

- migrating the remaining components;
- removing the production component style pack;
- exposing Maia as catalogue-wide;
- or updating the roadmap from pilot to committed architecture.
