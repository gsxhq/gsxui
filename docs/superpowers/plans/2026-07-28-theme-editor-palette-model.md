# Theme Editor Palette Model Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace raw semantic-color inputs with shadcn-style Base Color, Theme, and named Radius controls while preserving exact preset v1 transport.

**Architecture:** Go owns an immutable palette catalog and resolves each Base Color + Theme + Radius selection into the existing complete preset. The server serializes that catalog to the browser; browser state keeps committed and transient preview values separate, while unmatched imported presets appear as Custom.

**Tech Stack:** Go 1.26.1, GSX, browser-native JavaScript, Popover and native Radio components, Playwright, Vite.

## Global Constraints

- Work only in `/Users/jackieli/personal/gsxhq/gsxui/.worktrees/css-only-theme-architecture`.
- Do not stop or replace the user's process on port 5173.
- Do not hand-edit generated `.x.go`; run `go tool gsx generate`.
- Preset schema v1, canonical field ordering, and JSON/CSS/share round-trip
  rules remain unchanged.
- Audited shadcn Neutral values intentionally replace differing overlapping
  built-in defaults; update the pinned default JSON/CSS/share goldens.
- The checked-in catalog must not read `/Users/jackieli/personal/shadcn-ui` at build or test time.
- Base Color supplies the neutral/surface layer; Theme may override only the
  approved brand set, including the deliberately layered secondary pair.
- gsxui-only status and overlay tokens remain present.
- Do not add a color wheel, palette generator, chart tokens, or advanced per-token inputs.
- Every task uses red-green TDD and ends in its own commit.

---

## File map

- `internal/preset/catalog.go`: exported copy-returning catalog API, composition, and exact reverse matching.
- `internal/preset/catalog_data.go`: package-private copied definitions with source provenance.
- `internal/preset/catalog_test.go`: catalog completeness, ownership, immutability, uniqueness, and matching.
- `site/pages/theme.gsx`: browser schema and editor composition.
- `site/pages/theme_picker.gsx`: accessible Popover + Radio picker markup.
- `web/theme-state.js`: committed selection, transient preview, custom matching, and unchanged transport.
- `web/theme.js`: picker interaction, rendering, imports, iframe preview, and artifact output.
- `web/theme-state.test.js`: deterministic state transitions and transport invariants.
- `jstest/specs/theme-editor.spec.ts`: picker UX, hover preview, custom imports, and existing iframe behavior.

### Task 1: Immutable Go palette catalog

**Files:**
- Create: `internal/preset/catalog.go`
- Create: `internal/preset/catalog_data.go`
- Create: `internal/preset/catalog_test.go`
- Modify: `internal/preset/preset.go`
- Modify: `internal/preset/{preset_test,json_test,css_test}.go`
- Modify: `internal/preset/testdata/default-nova.json`
- Modify: `internal/preset/testdata/default-theme.css`
- Modify: `assets/css/themes/default.css`
- Modify: `site/pages/theme_defaults_test.go`

**Interfaces:**
- Produces: `type PaletteChoice struct { Name, Title, Swatch string }`
- Produces: `type RadiusChoice struct { Name, Title, Value string }`
- Produces: `type PaletteSelection struct { BaseColor, Theme, Radius string }`
- Produces: `const CustomChoice = "custom"`
- Produces: `func BaseColorChoices() []PaletteChoice`
- Produces: `func ThemeChoices(baseColor string) ([]PaletteChoice, error)`
- Produces: `func RadiusChoices() []RadiusChoice`
- Produces: `func DefaultPaletteSelection() PaletteSelection`
- Produces: `func ResolvePalette(style Style, selection PaletteSelection) (Preset, error)`
- Produces: `func MatchPalette(preset Preset) PaletteSelection`

- [ ] **Step 1: Write failing catalog contract tests**

Create `internal/preset/catalog_test.go` with:

```go
func TestPaletteCatalogCardinality(t *testing.T) {
	if got := len(BaseColorChoices()); got != 7 {
		t.Fatalf("base colors = %d, want 7", got)
	}
	themes, err := ThemeChoices("neutral")
	if err != nil {
		t.Fatal(err)
	}
	if got := len(themes); got != 18 {
		t.Fatalf("neutral themes = %d, want selected base plus 17 accents", got)
	}
}
```

Add table tests for the exact base names `neutral`, `stone`, `zinc`, `mauve`,
`olive`, `mist`, `taupe`; the 17 accent names from the design spec; and radius
name/value pairs. Mutate returned slices and confirm a second call is unchanged.

For every style, base, available theme, and radius, call `ResolvePalette`, then
require `Validate` succeeds and `MatchPalette` returns the same names. Hash the
complete light/dark maps once per Base Color + Theme pair and reject duplicate
palette resolutions independently from radius.

- [ ] **Step 2: Run tests and verify the catalog API is missing**

```bash
go test ./internal/preset -run 'TestPalette' -count=1
```

Expected: compile failure because the catalog API is undefined.

- [ ] **Step 3: Define the catalog API and ownership sets**

In `catalog.go`, define the interfaces above and these package-private sets:

```go
var baseLayerTokens = map[string]bool{
	"background": true, "foreground": true,
	"card": true, "card-foreground": true,
	"popover": true, "popover-foreground": true,
	"secondary": true, "secondary-foreground": true,
	"muted": true, "muted-foreground": true,
	"accent": true, "accent-foreground": true,
	"destructive": true, "border": true, "input": true, "ring": true,
	"sidebar": true, "sidebar-foreground": true,
	"sidebar-accent": true, "sidebar-accent-foreground": true,
	"sidebar-border": true, "sidebar-ring": true,
}

var themeLayerTokens = map[string]bool{
	"primary": true, "primary-foreground": true,
	"secondary": true, "secondary-foreground": true,
	"sidebar-primary": true, "sidebar-primary-foreground": true,
}
```

The overlap at `secondary` and `secondary-foreground` is deliberate: the Theme
layer wins, exactly like shadcn's ordered merge. Catalog validation must run
deterministically from tests and reject unknown, missing, duplicate, or
out-of-layer keys before resolution.

- [ ] **Step 4: Port the audited shadcn values into package-private data**

In `catalog_data.go`, add a provenance comment naming:

```text
/Users/jackieli/personal/shadcn-ui/apps/v4/registry/themes.ts
```

For each of the seven base entries, copy its overlapping light/dark surface
values into the Base Color layer and its brand values into a same-named Theme
definition. Copy the same brand subset for all 17 accent definitions. This
makes `stone + stone` reproduce shadcn Stone and `stone + blue` reproduce
shadcn's ordered Stone/Blue merge. Do not copy chart keys or radius from
shadcn. Use the exact strings from the audited file; do not convert or
normalize OKLCH.

Define named radii exactly:

```go
var radiusCatalog = []RadiusChoice{
	{Name: "none", Title: "None", Value: "0"},
	{Name: "small", Title: "Small", Value: "0.45rem"},
	{Name: "medium", Title: "Medium", Value: "0.625rem"},
	{Name: "large", Title: "Large", Value: "0.875rem"},
}
```

- [ ] **Step 5: Implement resolution and reverse matching**

Start from a deep clone of the canonical token defaults, apply the selected
base mode maps, then the selected theme maps, and set the selected radius.
The selected base name is itself the neutral Theme choice. Accent names are
valid with every base.

`MatchPalette` validates the preset, compares every token string in both modes
without normalization, and matches radius independently. Return
`CustomChoice` for palette fields or radius when no exact match exists.
Preserve `preset.Style` outside the matching decision.

Make `Default(style)` resolve:

```go
PaletteSelection{BaseColor: "neutral", Theme: "neutral", Radius: "medium"}
```

without changing its public output shape.

- [ ] **Step 6: Update the intentional default fixtures and assertions**

Regenerate `default-nova.json` with `CanonicalJSON(Default(StyleNova))` and
regenerate both CSS copies with `ThemeCSS(Default(StyleNova))`. Update exact
default assertions in preset and site tests to the audited Neutral strings.
Do not loosen byte-for-byte golden checks.

- [ ] **Step 7: Run catalog and existing transport tests**

```bash
go test ./internal/preset -run 'Test(Palette|Default|NovaAndMaia|Canonical|ThemeCSS|Share)' -count=1
go test ./site/pages -run 'TestTheme(EditorRendersPresetGroupNamesAndDefaults|DefaultsExposeEverySidebarTokenInTheEditor)' -count=1
node --test web/theme-state.test.js
```

Expected: PASS after intentionally updating built-in-default transport
goldens; canonical ordering and round-trip tests remain unchanged.

- [ ] **Step 8: Commit the catalog**

```bash
git add internal/preset assets/css/themes/default.css site/pages/theme_defaults_test.go
git commit -m "feat: add curated theme palette catalog"
```

### Task 2: Palette-aware browser state without transport changes

**Files:**
- Modify: `site/pages/theme.gsx`
- Modify generated: `site/pages/theme.x.go`
- Modify: `web/theme-state.js`
- Modify: `web/theme-state.test.js`
- Modify: `site/pages/theme_schema_test.go`

**Interfaces:**
- Consumes: catalog API from Task 1
- Produces browser schema field: `palette`
- Produces state functions: `selectBaseColor`, `selectTheme`, `selectRadius`,
  `previewBaseColor`, `previewTheme`, `clearPalettePreview`,
  `previewPreset`, and `replacePreset`
- Removes: `applyField` and invalid per-token `drafts`

- [ ] **Step 1: Write failing Go tests for the serialized palette schema**

Extend `TestThemeEditorSchemaMatchesPresetAuthority` to parse the `palette`
field from `[data-theme-schema]` JSON and require:

```go
type paletteSchemaProbe struct {
	BaseColors []struct{ Name string `json:"name"` } `json:"baseColors"`
	Themes     map[string][]struct{ Name string `json:"name"` } `json:"themes"`
	Radii      []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"radii"`
}
```

Assert 7 base colors, 18 Theme choices for `neutral`, 4 radii, and a complete
resolved light/dark map for `neutral + blue`.

- [ ] **Step 2: Write failing state tests**

In `web/theme-state.test.js`, add tests that:

- create state matched to `neutral + neutral + medium`;
- select `stone + blue` and verify committed values change;
- preview `rose`, verify `previewPreset(state)` changes while
  `canonicalJSON(state.resolved)` does not;
- clear preview and recover the committed preset;
- replace with one changed token and get `custom/custom`;
- replace with radius `1rem` and get radius `custom`; and
- selecting a built-in base from Custom commits that base plus its same-named
  Theme while preserving a Custom radius;
- selecting Blue from Custom commits `neutral + blue`; and
- selecting a named radius from Custom preserves the exact imported colors.

- [ ] **Step 3: Run focused tests and verify failures**

```bash
go test ./site/pages -run TestThemeEditorSchemaMatchesPresetAuthority -count=1
node --test web/theme-state.test.js
```

Expected: FAIL because the schema and state transitions do not exist.

- [ ] **Step 4: Serialize the Go catalog into the existing editor schema**

Add JSON-tagged palette schema structs in `theme.gsx`. Serialize:

- copy-returning Base Color choices;
- per-base available Theme choices;
- named radii;
- complete resolved light/dark themes for each valid Base Color + Theme pair;
  and
- the default selection.

Do not add palette names to preset JSON. Keep `defaults`, `canonicalDefaults`,
token names, and radius units for import/transport validation.

- [ ] **Step 5: Refactor browser state around selections**

State becomes:

```js
{
  resolved,
  selection: { baseColor, theme, radius },
  previewResolved: null,
  mode: "light",
}
```

Selection functions clone state, resolve palette choices from
`schema.palette.resolved`, retain the current style, and clear transient
preview. Palette functions replace light/dark values without replacing the
current radius; radius functions replace the radius without replacing the
current light/dark values. From Custom, Base Color selects its same-named Theme
and Theme selects against Neutral. `previewPreset(state)` returns
`previewResolved ?? resolved`. `replacePreset` deep-clones an imported preset
and exact-matches it against schema resolutions. `resetThemeState` restores the
style's canonical default and recomputes its selection.

Keep `canonicalJSON`, `themeCSS`, share encoding, command strings, strict JSON
import, and CSS import byte-compatible. `importThemeCSS` must finish through
`replacePreset`.

- [ ] **Step 6: Run state, Go, and generated checks**

```bash
go tool gsx generate
go test ./site/pages -run TestThemeEditorSchemaMatchesPresetAuthority -count=1
node --test web/theme-state.test.js
go run ./internal/generatedcheck/cmd
```

Expected: PASS.

- [ ] **Step 7: Commit palette-aware state**

```bash
git add site/pages/theme.gsx site/pages/theme.x.go site/pages/theme_schema_test.go web/theme-state.js web/theme-state.test.js
git commit -m "refactor: model theme editor palette selections"
```

### Task 3: Accessible picker UI and transient preview

**Files:**
- Create: `site/pages/theme_picker.gsx`
- Create generated: `site/pages/theme_picker.x.go`
- Modify: `site/pages/theme.gsx`
- Modify generated: `site/pages/theme.x.go`
- Modify: `web/theme.js`
- Modify: `jstest/specs/theme-editor.spec.ts`
- Modify: `site/pages/pages_test.go`
- Modify: `site/pages/theme_defaults_test.go`

**Interfaces:**
- Consumes: Task 2 state functions and `schema.palette`
- Produces markup hooks: `data-theme-picker`, `data-theme-picker-trigger`,
  `data-theme-choice`, `data-theme-choice-swatch`, and
  `data-theme-selection-value`
- Preserves every existing transport, preview, style, mode, and Retry hook

- [ ] **Step 1: Write failing server tests for picker markup**

Extend `TestThemePageRoute` to require one picker for each of
`baseColor`, `theme`, and `radius`, and reject:

```go
for _, forbidden := range []string{
	`data-theme-var=`,
	`data-theme-field="light.`,
	`data-theme-field="dark.`,
} {
	if strings.Contains(body, forbidden) {
		t.Errorf("theme page still exposes raw token input %q", forbidden)
	}
}
```

Require radio semantics and named option text for Neutral, Blue, and Medium.

- [ ] **Step 2: Write failing Playwright tests**

Add tests to `jstest/specs/theme-editor.spec.ts` for:

- no per-token inputs;
- all seven Base Color choices;
- the selected base plus 17 Theme choices;
- changing Base Color and Theme updates iframe CSS variables;
- desktop pointer hover changes iframe only, then restores on pointer leave;
- hover leaves downloaded/canonical JSON unchanged;
- click commits and changes JSON/share/commands;
- custom JSON and CSS imports show Custom with exact values retained; and
- a built-in selection atomically replaces Custom.

- [ ] **Step 3: Run the tests and verify old UI failures**

```bash
go test ./site/pages -run TestThemePageRoute -count=1
npx playwright test --config jstest/playwright.config.ts jstest/specs/theme-editor.spec.ts
```

Expected: FAIL because raw token fields still render and no picker hooks exist.

- [ ] **Step 4: Build a reusable server-rendered picker**

Create `theme_picker.gsx` using `ui.Popover`, `ui.PopoverTrigger`,
`ui.PopoverContent`, and native `ui.Radio`. Each picker trigger shows its label,
selected title, and supplemental swatch. Each option has a stable value, radio
state, text label, and swatch.

Render only the selected base plus accent choices in the Theme picker. When the
palette is Custom, render Neutral plus the accent choices because Neutral is
the documented fallback for a Theme selection. Render Custom as the trigger
value when matched state is custom; do not render Custom as a selectable radio
option.

- [ ] **Step 5: Replace token groups with Base Color, Theme, and Radius**

Remove `ThemeGroups`, token rows, light/dark field panels, radius text input,
and field-error spans from `theme.gsx`. Replace the old raw-field rendering
tests in `theme_defaults_test.go` with catalog/picker default assertions. Keep:

- Nova/Maia style selection;
- light/dark preview mode;
- Reset;
- picker controls;
- JSON/CSS transport sections;
- share/install output;
- status, manual-copy fallback, and the one iframe.

- [ ] **Step 6: Wire committed and transient interactions**

In `theme.js`:

- render selected trigger labels, radio checks, and swatches;
- replace the same-named Theme option's value, label, and swatch from the
  server schema whenever Base Color changes;
- commit native radio `change` events through Task 2 state functions;
- on fine-pointer `pointerenter`, apply a transient preview;
- on picker `pointerleave`, `toggle`, or dismissal, clear preview;
- send `previewPreset(state)` to the iframe;
- generate every artifact only from `state.resolved`; and
- render Custom when reverse matching fails.

Do not infer color names from CSS strings in the browser.

- [ ] **Step 7: Generate and run focused tests**

```bash
go tool gsx generate
go test ./site/pages -run TestThemePageRoute -count=1
node --test web/theme-state.test.js
npx playwright test --config jstest/playwright.config.ts jstest/specs/theme-editor.spec.ts
go run ./internal/generatedcheck/cmd
```

Expected: PASS.

- [ ] **Step 8: Commit the picker UI**

```bash
git add site/pages web/theme.js jstest/specs/theme-editor.spec.ts
git commit -m "feat: replace raw theme fields with palette pickers"
```

### Task 4: Palette integration gate and adversarial review

**Files:**
- Review all Task 1–3 files
- Modify only when verification exposes a defect

**Interfaces:**
- Produces a clean, independently reviewed palette subsystem

- [ ] **Step 1: Run static and authoritative checks**

```bash
gopls check -severity=hint internal/preset/catalog.go internal/preset/catalog_data.go internal/preset/catalog_test.go site/pages/theme.x.go site/pages/theme_picker.x.go site/pages/pages_test.go
node --check web/theme-state.js
node --check web/theme.js
git diff --check
make check
```

Expected: no new diagnostics and the full repository gate passes.

- [ ] **Step 2: Verify the production bundle**

```bash
npm run build
git restore -- site/dist/.gitkeep
git status --short
```

Expected: the dedicated preview entry remains lightweight and no build output
is left in the worktree.

- [ ] **Step 3: Request independent adversarial review**

The reviewer must probe catalog ownership and mutation, ambiguous reverse
matches, Custom import round trips, hover/commit separation, keyboard radio
behavior, light/dark iframe values, Reset, malformed imports, handshake Retry,
and exact preset transport. The review passes only with no Critical or
Important findings.

- [ ] **Step 4: Confirm process and worktree cleanup**

```bash
lsof -nP -iTCP:7799 -sTCP:LISTEN || true
git status --short --branch
```

Expected: no harness remains on 7799, port 5173 is untouched, and the worktree
is clean.
