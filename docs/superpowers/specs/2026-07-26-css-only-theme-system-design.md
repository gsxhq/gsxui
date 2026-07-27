# CSS-only theme system and configurator architecture

**Status:** architecture approved. The structural Phase 1 boundary is
implemented with the original packed slot encoding; the presence-marker
revision below is approved and in progress. The configurator roadmap remains
unscheduled.

This document defines the boundary we want before replacing the current
theme editor. It is deliberately an architecture and roadmap, not a promise
to ship a catalogue of named component styles.

The central decision is:

> One canonical component template per semantic component. Themes provide
> values. Styles provide presentation in CSS. A style never forks the
> component's markup or behavior.

That makes a future second style a CSS asset, not another copy of all 53
component templates. It also makes today's single style worth refactoring:
the editor, CLI, site, tests, and consumer project all operate on the same
explicit theme/style contracts instead of reconstructing them independently.

## 1. Problem

The current `/theme` page is a token form:

- 20 light/dark custom-property text inputs;
- import, copy, and download of one generated `gsxui.css`;
- a small inline preview;
- no preset model shared with the CLI;
- no visual controls, undo, responsive canvas, or full component coverage.

The component sources also combine semantic rendering and the current Nova
presentation in Tailwind class strings. `assets/gsxui.css` combines Tailwind
setup, theme values, global base rules, and interaction-critical component
CSS. `web/site.css` then copies some of that content for the site. These
boundaries make a second presentation style a repository-wide component
fork or another large class-string retarget.

The editor gap is therefore not primarily a page-design problem. It is an
ownership problem:

```
component template ─┬─ semantics, accessibility, behavior
                    └─ current visual treatment

single CSS file ────┬─ Tailwind integration
                    ├─ theme values
                    ├─ component mechanics
                    └─ current visual treatment
```

A richer editor built on that model would merely put better controls in
front of the same coupling.

## 2. Goals

1. Preserve one canonical `.gsx` implementation of each semantic component.
2. Make component presentation replaceable with CSS alone.
3. Make the theme, style, editor, CLI, and site consume one versioned preset
   model.
4. Keep caller-supplied Tailwind utilities as the final per-instance
   override.
5. Give CSS hooks and style packs a validated public contract.
6. Make CLI application deterministic, ownership-aware, and recoverable.
7. Give the editor the same direct-manipulation quality as shadcn Create:
   visual controls, a real component canvas, import/export, responsive and
   light/dark inspection, and clear new/existing-project handoff.
8. Preserve the current Nova rendering and behavior through the
   architectural migration.

## 3. Non-goals

- We do not commit to shipping a second named style.
- We do not commit to tracking shadcn's style names or release cadence.
- We do not add a style picker while only one validated style exists.
- We do not put markup, accessibility, JavaScript behavior, dependencies,
  or component APIs in a style pack.
- We do not support arbitrary CSS pasted into a preset.
- We do not preserve the current combined `web/gsxui.css` layout through a
  compatibility shim. gsxui is early enough for one clean migration.
- We do not build a local `gsxui theme` web server. The remote editor and
  the CLI share a preset format; the CLI does not need to host the editor.
- We do not promise every axis visible in shadcn Create. gsxui exposes axes
  it can own faithfully. For example, a chart axis requires a chart
  component and an icon-library axis is not a CSS style.

## 4. Reference: what shadcn does

The reference was inspected at shadcn-ui commit
[`7774cd7dcee1e98d0815aa6e829f33a7fc952fdf`](https://github.com/shadcn-ui/ui/commit/7774cd7dcee1e98d0815aa6e829f33a7fc952fdf).
Its current system is a useful hybrid:

- `apps/v4/registry/bases/{aria,base,radix}` contains separate structural
  component families.
- `apps/v4/registry/styles/style-*.css` contains scoped style rules over
  `.cn-*` hooks using `@apply`.
- The Create preview loads the styles together, selects one with a scoped
  body class, and injects theme variables separately.
- `packages/shadcn/src/styles/create-style-map.ts` parses the CSS with
  PostCSS and a selector parser.
- `transform-style-map.ts` uses the resulting map, `ts-morph`, and
  `tailwind-merge` to flatten the selected style's utilities into the TSX
  source delivered to consumers.
- Preset codes efficiently encode a known set of configurator choices.

The valuable ideas are the separation of structural bases, style hooks,
theme values, preview scoping, parser-backed validation, and one preset
driving preview and install.

gsxui deliberately diverges at the delivery boundary. We keep the style as
CSS in the consumer project instead of flattening it back into every
component template. Flattening would recreate the coupling this refactor is
meant to remove and would require a Go/GSX source transformer for every
future style. The trade-off is that the style selectors become a public
contract and the consumer must import one additional CSS asset. We accept
that trade because it gives gsxui one component source and independently
replaceable presentation.

## 5. Ownership model

Every piece of UI belongs to exactly one layer.

| layer | owns | must not own |
|---|---|---|
| component `.gsx` | semantic HTML, public API, accessibility, behavior/state attributes, stable styling hooks, caller class | product styling, palette, density, decorative shape |
| component `.js` | interaction and reflected state | visual values or class-string styling |
| foundation CSS | Tailwind integration and mechanics required for the component to function | a particular visual style |
| theme CSS | semantic token values for light/dark and typed global design values | component selectors |
| style CSS | component presentation over the public hooks | behavior, markup assumptions outside the contract, global page styling |
| caller utilities | intentional per-instance changes | library defaults |

The cascade is:

```
Tailwind theme/base
    ↓
gsxui foundation
    ↓
gsxui style in @layer components
    ↓
caller utilities
```

Style rules live in `@layer components`. Tailwind's utilities layer then
lets a caller's `class="rounded-none"` override the style pack's default
rounding without `!important` or a merge step. Style rules use low-specificity
selectors and never use `!important`.

### Foundation versus presentation

Foundation is intentionally narrow. It includes things such as:

- hiding a closed popover despite author CSS;
- native dialog/top-layer mechanics;
- the grid/overflow mechanics required for an accordion close animation;
- browser-specific control normalization needed for the control to work;
- Tailwind token mapping and the dark custom variant.

It does not include a card's padding, a button's height, menu corner radius,
surface shadows, palette choices, or typography scale. Those are style.
Accessible focus indication belongs to style and is a conformance
requirement for every style pack; the foundation is not a visual fallback
style.

When a rule is disputed, the test is counterfactual: if changing it can
produce a coherent visual style without breaking semantics, interaction, or
accessibility, it is style. If removing it makes the primitive malfunction
regardless of visual treatment, it is foundation.

## 6. The component styling contract

### Namespaced slot markers

Components expose one stable, namespaced presence attribute per styling
role:

```gsx
<button
	data-variant={variant |> default("default")}
	data-size={size |> default("default")}
	{ attrs... }
	data-gsxui-slot-button
>
	{ children }
</button>
```

The style pack targets those attributes:

```css
@layer components {
  [data-gsxui-slot-button] {
    @apply inline-flex h-8 items-center justify-center rounded-lg px-2.5;
  }

  [data-gsxui-slot-button][data-variant="outline"] {
    @apply border border-border bg-background;
  }
}
```

The `data-gsxui-slot-<name>` prefix replaces the current generic `data-slot`
as the public styling namespace. Existing component-specific behavior hooks
such as `data-gsxui-calendar-*` remain behavior hooks; styles may read
reflected state, but JavaScript must never toggle presentation classes.

Each role has its own valueless presence attribute instead of sharing a
space-separated scalar value. CSS uses only exact presence selectors such as
`[data-gsxui-slot-button]`; value and token operators are not part of the
contract. Composed components therefore rely on GSX's normal fallthrough
attribute merge. For example, `AlertDialogAction` adds its role to `Button`:

```gsx
<Button { attrs... } data-gsxui-slot-alert-dialog-action>
	{ children }
</Button>
```

`Button` forwards that distinct key and forces its own invariant marker
after the spread, so the result is:

```html
<button
  data-gsxui-slot-alert-dialog-action
  data-gsxui-slot-button
>...</button>
```

This representation deliberately matches GSX's attribute semantics.
`class` and `style` aggregate, while other duplicate scalar keys are
last-wins. Packing several roles into one scalar would require a private
token-merging helper in every vendored project. Giving every role a distinct
key lets ordinary forwarding accumulate roles without parsing, allocation,
deduplication, an internal support package, or a gsxui-specific merge API.
It also prevents every style pack from duplicating Button, Input, Label,
Separator, and Dialog rules into their composed parts.

The contract has these rules:

1. A slot marker is the valueless presence attribute
   `data-gsxui-slot-<name>`; `<name>` is globally unique within gsxui and
   uses kebab-case.
2. A slot identifies a semantic component part, not one CSS declaration.
3. A component places its own marker after its fallthrough spread so callers
   cannot remove the component's styling identity.
4. A composed component adds another distinct marker; GSX forwarding keeps
   every styling role on the final element without a helper.
5. Public presentation axes are reflected as attributes with explicit
   values, such as `data-variant`, `data-size`, `aria-invalid`, or
   `data-state`.
6. Component templates contain no library-owned presentation utilities
   after migration. Their `class` output is the caller-supplied class.
7. Inline styles are permitted only for dynamic values computed by behavior,
   such as a slider's live fill percentage or a resizable panel's flex
   value. Style CSS consumes those custom properties.
8. A style may select a slot, its declared states, its pseudo-elements, and
   documented descendants. It may not depend on undocumented DOM depth or
   site/page selectors.

### First-class contract manifest

The hooks are an API and must not be inferred by a grep. A typed style
contract is the source of truth for:

- component name;
- slots it emits;
- public variant, size, state, and orientation values;
- which slot/state combinations require style coverage;
- documented parent/child slot relationships;
- mechanism rules owned by foundation.

The manifest is serialized for the editor and validators, but authored as
typed Go data so the registry and Go tests consume it directly. Contract
tests render the registered component examples with an HTML parser and prove
both directions:

- every emitted `data-gsxui-slot-<name>` marker and declared state exists in
  the contract;
- every contract slot and declared state value is emitted by at least one
  registered rendering fixture;
- every rendered marker is valueless and every style selector uses exact
  presence matching rather than a value operator;
- declared composition relationships are co-located on one rendered element,
  such as `data-gsxui-slot-alert-dialog-action` and
  `data-gsxui-slot-button`;
- a caller-supplied value for a component's own marker cannot suppress or
  replace the valueless invariant marker, while a distinct caller role
  survives forwarding.

The source tree, embedded manifest, generators, and fresh `gsxui add` output
contain no `ui/slots.go`, `ui/internal/slotattr`, or dependency edge to either.
Automatically deleting those paths from an already-vendored external project
is a separate destructive migration policy and is not implied by this source
contract.

The CSS validator parses selectors into an AST. It rejects unknown slots or
state values, selectors outside the gsxui namespace, forbidden global
selectors, malformed `@apply`, and rules outside the expected cascade
layer. It verifies coverage of every contract entry marked required,
including selector lists. There is no regex selector parser and no
tolerance threshold.

## 7. CSS assets

The repository has one authored source for each asset:

```text
assets/css/
  index.css
  foundation.css
  themes/
    default.css
  styles/
    default.css
```

`index.css` imports Tailwind, `tw-animate-css`, foundation, the selected
theme destination, and the selected style destination in the defined
order. The site builds from these same sources; it does not copy token or
mechanism blocks into `web/site.css`.

`gsxui init` vendors the selected result as:

```text
web/gsxui/
  index.css
  foundation.css
  theme.css
  style.css
  index.js
  ...component behavior modules
```

The configured CSS path becomes the entry file
`web/gsxui/index.css`. Its sibling filenames are fixed by this contract so
imports remain deterministic. The project's Tailwind entry imports
`web/gsxui/index.css` and keeps its project-specific `@source` declarations
outside the vendored files.

- `foundation.css` changes only with a gsxui upgrade.
- `theme.css` is generated from the project's preset and is project-owned.
- `style.css` is a vendored, unscoped production style pack and is
  project-owned once manually edited.
- `index.css` is stable wiring.

Only one unscoped style pack ships into a production project by default.
There is no runtime style identifier in component markup and no unused
style catalogue in the application bundle.

### Theme versus style

A theme supplies semantic values: background, foreground, surfaces,
actions, borders, rings, radius, and typography. A style decides how those
values are used by components: height, spacing, shape, border/ring
treatment, shadows, and type scale.

A style may only use semantic color tokens, `currentColor`, or transparent
color. It must not bake a palette into component selectors. This preserves
the full theme × style matrix rather than creating style-specific themes.

### Structural variation

If a proposed visual treatment requires different roles, elements, behavior,
dependencies, or public props, it is not a style.

- A useful optional semantic part can be added once to the canonical
  component contract.
- A fundamentally different primitive is a different component.
- A systematic alternative primitive family could become a future base,
  like shadcn's Radix/Base UI/ARIA split, but no base abstraction is added
  until gsxui has a real second structural family.

There are no style-specific `.gsx` templates.

## 8. Preset model

The preset is the shared source of truth for editor state and CLI output.
It is durable, human-readable, and versioned JSON:

```json
{
  "$schema": "https://ui.gsxhq.dev/schemas/preset-v1.json",
  "schemaVersion": 1,
  "style": "default",
  "radius": "0.625rem",
  "typography": {
    "sans": "Geist, ui-sans-serif, system-ui, sans-serif",
    "mono": "Geist Mono, ui-monospace, monospace"
  },
  "theme": {
    "light": {
      "background": "oklch(1 0 0)",
      "foreground": "oklch(0 0 0)",
      "primary": "oklch(0 0 0)"
    },
    "dark": {
      "background": "oklch(0.145 0 0)",
      "foreground": "oklch(0.985 0 0)",
      "primary": "oklch(0.922 0 0)"
    }
  }
}
```

The real schema enumerates every required semantic token in canonical
order; the shortened example above is illustrative. Within schema version
1, unknown fields, unknown styles, missing tokens, and duplicate JSON keys
are errors. A client rejects a future schema version it does not understand
instead of partially applying it.

The preset contains structured values, not raw CSS:

- color tokens must parse as one CSS color value;
- radius must parse as one non-negative CSS length;
- typography is a structured font-family list and is escaped when emitted;
- token names come from the schema and cannot be supplied by the preset.

The Go CLI performs authoritative parser-backed validation. The browser uses
native CSS parsing for immediate field feedback and the same generated
schema for names, types, required values, and canonical ordering. Imported
CSS is parsed into declarations and converted to a preset; regex block
extraction is retired.

Canonical JSON uses the schema's field and token order and a trailing
newline. Import → normalize → export is deterministic and idempotent.

### Preset references and transport

`gsxui.preset.json` at the project root is the persisted project preset.
The CLI also accepts a file, stdin, URL, or share code and resolves all of
them to the same validated model before touching the filesystem.

The initial share code is `gsxui:v1:` followed by base64url-encoded
canonical preset JSON. A share URL carries that code as one query
parameter. The code is a transport envelope, not another preset format. A
future compressed transport can use a new transport version without
changing preset schema version 1.

## 9. Resolver and data flow

One pure resolver turns a validated preset into all derived artifacts:

```text
preset JSON
    ↓ parse + validate + normalize
resolved preset
    ├─ editor control state
    ├─ preview theme variables
    ├─ selected preview style
    ├─ generated theme.css
    ├─ selected production style.css
    └─ canonical preset JSON / share code
```

The renderer and CLI do not independently assemble CSS strings. Golden
tests pin the resolver's normalized preset, `theme.css`, style identity, and
share-code round trip.

The `style` field exists from schema version 1 so adding a real second pack
does not require a new preset shape. While only `default` exists, the editor
does not show a decorative one-option style control.

## 10. CLI application model

The intended commands are:

```text
gsxui init --preset <file-or-code>
gsxui apply --preset <file-or-code>
gsxui apply --preset <file-or-code> --only theme
gsxui apply --preset <file-or-code> --only style
```

Full apply writes `gsxui.preset.json`, `theme.css`, and `style.css`.
`--only theme` changes the preset's theme/radius/typography fields and
`theme.css` while retaining the project's resolved style. `--only style`
changes the style field and `style.css` while retaining the theme.

`gsxui.json` records a content hash for each CLI-managed CSS artifact.
Application follows these ownership rules:

- absent file: create it;
- file equals its recorded managed hash: replace it;
- file differs from the recorded hash: refuse and print a focused diff;
- `--overwrite`: explicitly replace a modified managed file.

All inputs and every target conflict are validated before any target is
changed. Writes use a transaction directory, same-filesystem renames, a
rollback snapshot, and a small recovery journal. A validation, conflict, or
injected write failure must leave the project byte-identical. If the
process is interrupted between renames, the next CLI invocation completes
rollback before doing new work.

The command reports exactly which preset axes and files will change.
Applying an already-applied preset is a no-op.

## 11. Editor architecture

The new editor is a configurator over the preset model, not a CSS text
generator with a preview attached.

### State

One immutable resolved preset is the application state. Every control
dispatches a typed change. Undo and redo store preset revisions, not DOM or
CSS snapshots. The last valid preset and editor layout persist locally.
Reset restores the selected built-in starting preset.

Invalid field text may remain in the focused control for correction, but it
does not enter resolved state, preview state, history, export, or share
links.

### Controls

The first complete editor exposes the axes gsxui owns:

- visual color controls with an exact CSS-value escape hatch;
- radius;
- sans and mono typography stacks;
- light and dark values with linked/unlinked editing;
- style only after more than one validated style exists;
- reset, undo, redo, deterministic shuffle, import, export, and share.

Shuffle accepts a visible seed and produces a valid preset from defined
palettes/constraints. The same seed and starting preset always produce the
same result; it is not uncontrolled randomness.

### Preview

The component canvas runs in an iframe so editor chrome cannot leak into
component CSS and component top-layer behavior remains realistic. It uses
the real example registry rather than hand-written preview replicas and
supports:

- representative overview compositions plus every component example;
- light and dark;
- desktop, tablet, and mobile viewport sizes;
- interactive states through the real behavior modules;
- keyboard operation;
- a reload/error surface when the preview cannot render.

Theme values are injected into the iframe from resolved preset state.
For editor use only, the build derives a scoped bundle for every validated
style:

```css
@layer components {
  body[data-gsxui-preview-style="default"]
    [data-gsxui-slot-button] {
    /* compiled default style declarations */
  }
}
```

The scoping transform is parser-backed and generated from the same unscoped
production pack. Scoped preview CSS is never hand-maintained and never
vendored into consumer projects.

### Handoff

The primary action opens a code panel with:

- **New project:** the complete init command carrying the preset code.
- **Existing project:** apply commands for full, theme-only, and style-only
  changes, with the files each command owns.
- **Files:** canonical `gsxui.preset.json`, `theme.css`, and `style.css`
  downloads for manual workflows.

Copy actions have visible success and failure feedback. Imported or shared
presets show validation errors at the exact field; the editor never silently
drops an unknown value.

## 12. Verification strategy

The refactor is complete only when these layers agree.

### Component contract

- Go render tests cover every declared slot and state value.
- Contract tests reject undeclared emitted hooks and declared hooks without
  a fixture.
- Existing accessibility and behavior tests continue to run against the
  canonical components.

### CSS contract

- A real CSS and selector parser validates every style and foundation pack.
- Unknown/misspelled slots and state values fail with the selector and
  source location.
- Every required contract entry has selector coverage.
- Tailwind compiles every `@apply`; no candidate is silently missing.
- Style selectors cannot escape the gsxui namespace.
- Semantic-token dependency checks reject literal component palettes.

### Rendering

- Current browser behavior tests run with foundation and style loaded as
  separate assets.
- Computed-style tests pin interaction-critical geometry, hidden/open
  visibility, focus indication, disabled/invalid state, and caller utility
  precedence.
- Fixed visual baselines cover representative components in light/dark and
  desktop/mobile. They test the current Nova appearance during migration
  and every theme × style combination if a second style is ever added.
- The site and consumer harness compile from the same authored CSS assets.

### Preset and CLI

- Schema, malformed input, unknown version, invalid CSS values, and hostile
  import cases.
- Canonical JSON and share-code round trips.
- Golden `theme.css` output.
- Filesystem tests for absent, identical, modified, overwrite, partial-axis,
  no-op, injected failure, recovery-journal, and byte-identical rollback.

### Generated artifacts

- `.x.go` files are regenerated from `.gsx`, never edited.
- Generated scoped editor styles, serialized contract/schema, and CSS
  artifacts have drift checks in `make check`.

Before a subsystem is merged, an adversarial reviewer builds throwaway
consumer projects that exercise modified managed files, hostile imported
CSS, missing style rules, broken style packs, caller overrides, top-layer
components, and interrupted CLI application.

## 13. Migration

The migration is a clean architectural cut:

1. Define the complete contract and asset split.
2. Move the current Nova presentation from all components into the default
   style pack; move only invariant mechanics into foundation.
3. Change components to emit namespaced slots and state axes.
4. Make the site and browser harness consume the shared assets.
5. Change `gsxui init` and configuration to vendor the split layout.
6. Remove the combined `assets/gsxui.css`, copied site blocks, visual
   library classes in `.gsx`, and generic `data-slot` hooks in the same
   change.

The work can be developed in component slices, but it does not merge with
two styling architectures active. The merge gate is full current visual
and behavioral equivalence across the component corpus.

Existing early-stage consumers receive a documented one-time migration:
replace the old CSS entry path, re-run init with explicit overwrite after
reviewing the diff, and regenerate components. There is no runtime adapter
or selector alias layer.

## 14. Decision gates

These are future triggers, not open architectural questions:

- **Second style:** begin only when a concrete second treatment is worth
  maintaining. It must use unchanged component markup, behavior,
  foundation, and theme; pass the full contract and visual matrix; and
  replace only `style.css` in a consumer project. Only then expose the style
  picker.
- **Structural base:** begin only when a real alternative primitive family
  cannot share the semantic/behavior contract. It is designed as its own
  project; style packs do not absorb the difference.
- **Additional preset axis:** add only when it has one owner and deterministic
  CLI output. A source/dependency choice is not mislabeled as CSS style.
- **Compressed share transport:** introduce a new transport version only if
  measured preset URLs exceed practical browser/share limits. The JSON
  schema remains independent.

## 15. Consequences

This refactor is large because every component's visual classes move and
the selectors become an API. That is the cost of correcting the boundary
once while the project is early.

In return:

- a theme changes values without changing component CSS;
- a style changes component CSS without changing templates or behavior;
- a caller changes one instance with ordinary utilities;
- the editor previews the exact artifacts the CLI installs;
- a second style, if it ever earns its way in, is an additive CSS pack and
  a test-matrix expansion rather than a component-library fork.
