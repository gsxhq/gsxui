# Theme editor curated palette model

**Date:** 2026-07-28

**Status:** Approved design

## Purpose

The current theme editor renders every semantic CSS token as a text input.
That exposes OKLCH serialization rather than a design decision. Users should
choose a coherent palette and see its effect; they should not be expected to
type values such as `oklch(0.488 0.243 264.376)`.

The primary editor adopts the current shadcn/create model:

- **Base Color** selects the neutral surface family;
- **Theme** selects the brand/accent family; and
- **Radius** selects a named geometry choice.

Exact light/dark CSS values remain part of the preset, preview, CLI handoff,
and import/export formats. They are no longer the primary editing UI.

This is a separate subsystem from the site-layout design. The two changes may
be executed in sequence, but the palette catalog and preset semantics remain
testable without the workspace layout.

## Reference model

The design follows the source under
`/Users/jackieli/personal/shadcn-ui/apps/v4`:

- `registry/themes.ts` owns named base and accent definitions;
- `registry/base-colors.ts` identifies the neutral families;
- `base-color-picker.tsx` and `theme-picker.tsx` expose those choices;
- `registry/config.ts` resolves a valid Base Color + Theme combination; and
- `design-system-provider.tsx` composes the selected values for the preview.

gsxui ports the model into Go-owned preset data. The browser does not import
shadcn code, and builds do not depend on the sibling checkout.

## User-facing controls

### Base Color

The initial catalog contains the same seven neutral families as the current
shadcn/create source:

- Neutral
- Stone
- Zinc
- Mauve
- Olive
- Mist
- Taupe

Selecting a base color replaces surface, foreground, card, popover, muted,
accent, border, input, ring, and sidebar-neutral values for both light and dark
modes.

### Theme

The Theme picker offers the selected base family as the neutral/default choice
plus the current shadcn accent families:

- Amber
- Blue
- Cyan
- Emerald
- Fuchsia
- Green
- Indigo
- Lime
- Orange
- Pink
- Purple
- Red
- Rose
- Sky
- Teal
- Violet
- Yellow

An accent theme overrides brand-related values on top of the chosen base:
primary, primary foreground, secondary, secondary foreground, sidebar primary,
and sidebar primary foreground. gsxui has no chart-token contract, so shadcn
chart values are not added in this slice.

Each base family therefore contributes two layers from the audited shadcn
theme entry: its surface values form the Base Color layer, while its
brand-related values form the same-named Theme choice. For example, Stone +
Stone reproduces the overlapping values of shadcn Stone, while Stone + Blue
uses Stone surfaces with Blue brand values.

The existing gsxui-only status and overlay tokens remain stable unless a
future catalog explicitly owns them:

- success
- info
- warning
- overlay
- contrast
- destructive foreground

### Radius

Radius uses named choices:

- None (`0`)
- Small (`0.45rem`)
- Medium (`0.625rem`, the built-in default)
- Large (`0.875rem`)

The selected preset still stores the resolved CSS length. An imported valid
length outside the named set appears as Custom.

shadcn can expose both Default and Medium because its configuration transports
the radius choice name: Default is an empty override, while Medium explicitly
selects `0.625rem`. Preset schema v1 transports only the resolved CSS length,
so those two labels would be indistinguishable after export and import. gsxui
therefore exposes the four unique resolved choices and calls the built-in
`0.625rem` value Medium instead of adding hidden editor metadata.

## Catalog ownership

`internal/preset` owns the authoritative catalog and composition rules.
Package-private storage is exposed through copy-returning functions so callers
cannot mutate shared definitions.

The catalog distinguishes:

```text
baseColorDefinition
  name
  title
  light overrides
  dark overrides

themeDefinition
  name
  title
  light overrides
  dark overrides

radiusChoice
  name
  title
  value
```

Resolution starts from the complete canonical gsxui token set, applies the
selected base definition, applies the selected theme definition, and finally
sets radius. The result must pass the existing `preset.Validate`.

Definitions contain only known tokens. Base definitions must resolve every
surface token they claim to supply in both modes. The same-named base Theme
definitions and the 17 accent Theme definitions may override only the
documented brand set. The secondary pair deliberately exists in both layers:
the Theme layer wins, matching shadcn's ordered merge. Startup/unit validation
rejects incomplete, duplicate, unknown, or out-of-scope definitions. It also
rejects two valid catalog combinations that resolve to the same complete
palette, keeping reverse matching unambiguous. The browser never repairs
catalog data.

The values are checked into gsxui with a provenance comment pointing to the
audited shadcn registry source. Generation and tests never read
`/Users/jackieli/personal/shadcn-ui`.

## Preset and transport semantics

Preset schema v1 remains self-contained and unchanged:

- `style` stores Nova or Maia;
- `radius` stores the resolved CSS length; and
- `theme.light` / `theme.dark` store every exact semantic value.

`baseColor` and `theme` are editor selections, not new required transport
fields. This preserves CLI compatibility and ensures shared/imported presets
remain reproducible even if the built-in catalog changes later.

The audited shadcn Neutral values intentionally become the canonical built-in
gsxui defaults for the overlapping token set. Existing gsxui-only status and
overlay values remain unchanged. This changes the pinned default JSON/CSS/share
golden values where the old defaults differed, but it does not change preset
schema v1, canonical field ordering, validation rules, or round-trip behavior.

When loading preset state, the editor compares its exact light/dark values
against catalog resolutions:

- one exact match selects its Base Color and Theme names;
- no exact match selects Custom;
- matching ignores no fields and performs no color normalization; and
- radius matching is independent from palette matching.

Selecting either color picker while the palette is Custom replaces the
arbitrary imported palette with a complete catalog palette. Selecting a Base
Color from Custom pairs it with that base's same-named Theme; selecting a Theme
from Custom pairs it with Neutral. These rules are explicit and deterministic
rather than relying on hidden pre-import state. Reset restores the canonical
built-in selection for the active style.

Palette and radius replacement remain independent. Choosing a built-in palette
preserves an exact Custom radius, and choosing a named radius preserves an
exact Custom palette. Only Reset replaces both.

JSON, share code, share URL, CSS, download, and CLI command outputs continue to
serialize the resolved preset exactly as they do today.

## Picker behavior

The customizer uses accessible picker triggers showing:

- the control label;
- selected name;
- a representative swatch; and
- selected radio state inside the opened choice list.

Desktop pointer hover previews a choice without committing it. Leaving the
picker restores the committed preset; clicking commits. Keyboard focus never
commits a transient value. On narrow/touch layouts, selection commits on
activation without hover behavior.

The preview override is transient UI state only. It must not alter JSON/CSS
export, URLs, commands, local committed state, or Reset behavior.

The light/dark mode switch remains because each named combination resolves both
modes. Style selection remains independent: Nova and Maia consume the same
resolved semantic palette.

## Custom and advanced workflows

Per-token text inputs are removed from the main editor.

CSS and JSON import remain the precision path for custom themes. A successful
custom import:

- updates the exact preset and iframe;
- marks the palette and/or radius as Custom when unmatched;
- remains exportable and shareable without loss; and
- can be replaced deliberately by selecting a built-in choice.

Import textareas remain in the existing transport section for this slice.
Moving them into a dialog or action menu is a separate UX decision.

## Accessibility

- Base Color, Theme, and Radius have explicit accessible names.
- Each option exposes radio semantics and the selected state.
- Swatches are supplemental; names never depend on color perception.
- Hover preview has an equivalent commit action and never traps keyboard focus.
- Custom is a state label, not a selectable empty palette.
- Existing validation and status messages remain live-region announcements.

## Verification

Go tests establish:

1. every catalog combination resolves to a valid complete preset;
2. returned definitions cannot mutate catalog storage;
3. base themes affect only their owned surface tokens;
4. accent themes affect only the documented brand tokens;
5. gsxui-only tokens remain present and valid;
6. the four unique named radii resolve exactly;
7. exact preset matching returns the correct names;
8. one changed token produces Custom; and
9. preset v1 canonical JSON and CSS round-trip rules remain unchanged, with
   intentionally updated built-in-default goldens.

Browser and state tests establish:

1. no per-token OKLCH inputs appear in the primary editor;
2. the seven Base Color choices and all Theme choices are available;
3. changing either choice updates both the committed preset and iframe;
4. light/dark preview modes use their corresponding resolved values;
5. desktop hover previews and then restores without changing exports;
6. clicking commits and updates exports/share/commands;
7. custom CSS/JSON import shows Custom without losing values;
8. selecting a built-in choice replaces Custom atomically;
9. named and custom radius states behave equivalently; and
10. existing malformed-import, handshake, Retry, and transactional state tests
    continue to pass.

## Non-goals

This slice does not add a free-form color wheel, generate a palette from one
brand color, normalize equivalent CSS color syntaxes, add chart tokens, expose
individual advanced token fields, edit typography/icons/primitives, or change
the preset schema version.
