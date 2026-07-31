package pages

import (
	"encoding/json"
	"strings"

	"github.com/gsxhq/gsxui/internal/preset"
	"github.com/gsxhq/gsxui/ui"
)

// Theme is the /theme route: a live editor over the shadcn-compatible
// CSS custom properties every gsxui component reads. Entirely client-side
// once loaded (web/theme.js) — the Go side only server-renders the default
// light/dark values so the page works before any JS runs.
type Theme struct{}

type themeEditorSchema struct {
	Schema            string                     `json:"schema"`
	SchemaVersion     int                        `json:"schemaVersion"`
	Transport         themeTransportSchema       `json:"transport"`
	TokenNames        []string                   `json:"tokenNames"`
	RadiusUnits       []string                   `json:"radiusUnits"`
	Styles            []string                   `json:"styles"`
	Defaults          map[string]json.RawMessage `json:"defaults"`
	CanonicalDefaults map[string]string          `json:"canonicalDefaults"`
	Palette           themePaletteSchema         `json:"palette"`
}

type themeTransportSchema struct {
	FullPrefix    string                      `json:"fullPrefix"`
	CompactPrefix string                      `json:"compactPrefix"`
	Compact       themeCompactTransportSchema `json:"compact"`
}

type themeCompactTransportSchema struct {
	Styles     []string `json:"styles"`
	BaseColors []string `json:"baseColors"`
	Themes     []string `json:"themes"`
	Radii      []string `json:"radii"`
}

type themePaletteSchema struct {
	BaseColors       []themePaletteChoiceSchema                       `json:"baseColors"`
	Themes           map[string][]themePaletteChoiceSchema            `json:"themes"`
	Radii            []themeRadiusChoiceSchema                        `json:"radii"`
	Resolved         map[string]map[string]themePaletteResolvedSchema `json:"resolved"`
	DefaultSelection themePaletteSelectionSchema                      `json:"defaultSelection"`
}

type themePaletteChoiceSchema struct {
	Name   string `json:"name"`
	Title  string `json:"title"`
	Swatch string `json:"swatch"`
}

type themeRadiusChoiceSchema struct {
	Name  string `json:"name"`
	Title string `json:"title"`
	Value string `json:"value"`
}

type themePaletteResolvedSchema struct {
	Light preset.ThemeValues `json:"light"`
	Dark  preset.ThemeValues `json:"dark"`
}

type themePaletteSelectionSchema struct {
	BaseColor string `json:"baseColor"`
	Theme     string `json:"theme"`
	Radius    string `json:"radius"`
}

func themeEditorSchemaValue() themeEditorSchema {
	styles := preset.Styles()
	schema := themeEditorSchema{
		Schema:            preset.SchemaURL,
		SchemaVersion:     preset.SchemaVersion,
		Transport:         themeTransportSchemaValue(),
		TokenNames:        preset.TokenNames(),
		RadiusUnits:       preset.RadiusUnits(),
		Styles:            make([]string, len(styles)),
		Defaults:          make(map[string]json.RawMessage, len(styles)),
		CanonicalDefaults: make(map[string]string, len(styles)),
		Palette:           themePaletteSchemaValue(),
	}
	for i, style := range styles {
		schema.Styles[i] = string(style)
		canonical, err := preset.CanonicalJSON(preset.Default(style))
		if err != nil {
			panic(err)
		}
		schema.Defaults[string(style)] = json.RawMessage(strings.TrimSpace(string(canonical)))
		schema.CanonicalDefaults[string(style)] = string(canonical)
	}
	return schema
}

func themeTransportSchemaValue() themeTransportSchema {
	transport := preset.ShareTransportSchema()
	styles := make([]string, len(transport.Styles))
	for i, style := range transport.Styles {
		styles[i] = string(style)
	}
	return themeTransportSchema{
		FullPrefix:    transport.FullPrefix,
		CompactPrefix: transport.CompactPrefix,
		Compact: themeCompactTransportSchema{
			Styles:     styles,
			BaseColors: transport.BaseColors,
			Themes:     transport.Themes,
			Radii:      transport.Radii,
		},
	}
}

func themePaletteSchemaValue() themePaletteSchema {
	selection := preset.DefaultPaletteSelection()
	schema := themePaletteSchema{
		BaseColors:       themePaletteChoiceSchemaValues(preset.BaseColorChoices()),
		Themes:           make(map[string][]themePaletteChoiceSchema),
		Radii:            themeRadiusChoiceSchemaValues(preset.RadiusChoices()),
		Resolved:         make(map[string]map[string]themePaletteResolvedSchema),
		DefaultSelection: themePaletteSelectionSchemaFromPreset(selection),
	}
	for _, baseColor := range preset.BaseColorChoices() {
		themes, err := preset.ThemeChoices(baseColor.Name)
		if err != nil {
			panic(err)
		}
		schema.Themes[baseColor.Name] = themePaletteChoiceSchemaValues(themes)
		schema.Resolved[baseColor.Name] = make(map[string]themePaletteResolvedSchema, len(themes))
		for _, theme := range themes {
			resolved, err := preset.ResolvePalette(preset.StyleNova, preset.PaletteSelection{
				BaseColor: baseColor.Name,
				Theme:     theme.Name,
				Radius:    selection.Radius,
			})
			if err != nil {
				panic(err)
			}
			schema.Resolved[baseColor.Name][theme.Name] = themePaletteResolvedSchema{
				Light: resolved.Theme.Light,
				Dark:  resolved.Theme.Dark,
			}
		}
	}
	return schema
}

func themePaletteChoiceSchemaValues(choices []preset.PaletteChoice) []themePaletteChoiceSchema {
	values := make([]themePaletteChoiceSchema, len(choices))
	for i, choice := range choices {
		values[i] = themePaletteChoiceSchema{Name: choice.Name, Title: choice.Title, Swatch: choice.Swatch}
	}
	return values
}

func themeRadiusChoiceSchemaValues(choices []preset.RadiusChoice) []themeRadiusChoiceSchema {
	values := make([]themeRadiusChoiceSchema, len(choices))
	for i, choice := range choices {
		values[i] = themeRadiusChoiceSchema{Name: choice.Name, Title: choice.Title, Value: choice.Value}
	}
	return values
}

func themePaletteSelectionSchemaFromPreset(selection preset.PaletteSelection) themePaletteSelectionSchema {
	return themePaletteSelectionSchema{
		BaseColor: selection.BaseColor,
		Theme:     selection.Theme,
		Radius:    selection.Radius,
	}
}

func themePickerChoices(choices []preset.PaletteChoice) []themePickerChoice {
	values := make([]themePickerChoice, len(choices))
	for i, choice := range choices {
		values[i] = themePickerChoice{
			Title:  choice.Title,
			Swatch: choice.Swatch,
			Value:  choice.Name,
		}
	}
	return values
}

func themeRadiusPickerChoices(choices []preset.RadiusChoice) []themePickerChoice {
	values := make([]themePickerChoice, len(choices))
	for i, choice := range choices {
		values[i] = themePickerChoice{
			Title:       choice.Title,
			Value:       choice.Name,
			SwatchStyle: "border-radius: " + choice.Value,
		}
	}
	return values
}

func mustThemePickerChoices(baseColor string) []preset.PaletteChoice {
	choices, err := preset.ThemeChoices(baseColor)
	if err != nil {
		panic(err)
	}
	return choices
}

const tabBtnBase = "rounded-md border border-border px-3 py-1.5 text-sm font-medium transition-colors"

const themeImportPlaceholder = `:root {
  --primary: oklch(0.6 0.2 280);
}
.dark {
  --primary: oklch(0.7 0.2 280);
}`

component (t Theme) Page() {
	<siteLayout title="Theme" active="" mode={layoutWorkspace} toc={nil}>
		<themeEditor previewURL={ThemePreview{} |> url} workspace/>
	</siteLayout>
}

// ThemeEditor is the editor body without the site shell, so the browser
// harness can exercise the production controls and web/theme.js directly.
component ThemeEditor(previewURL string) {
	<themeEditor previewURL={previewURL}/>
}

component themeEditor(previewURL string, workspace bool) {
	{{
		editorClass := "flex flex-col gap-8 py-10"
		gridClass := "grid grid-cols-1 gap-8"
		previewPanelClass := "flex min-w-0 flex-col gap-4"
		controlsPanelClass := "flex min-w-0 flex-col gap-7"
		iframeClass := "w-full rounded-xl border border-border bg-background shadow-sm"
		if workspace {
			editorClass += " lg:h-full lg:min-h-0 lg:gap-4 lg:py-0"
			gridClass += " lg:min-h-0 lg:flex-1 lg:grid-cols-[minmax(22rem,28rem)_minmax(0,1fr)] lg:grid-rows-[auto_minmax(12rem,1fr)] lg:gap-x-8 lg:gap-y-6"
			previewPanelClass += " lg:col-start-2 lg:row-span-2 lg:row-start-1 lg:min-h-0"
			controlsPanelClass += " lg:col-start-1 lg:row-start-2 lg:min-h-0 lg:overflow-y-auto lg:pr-2"
			iframeClass += " h-[min(70svh,640px)] lg:h-auto lg:min-h-0 lg:flex-1"
		} else {
			gridClass += " xl:grid-cols-[minmax(0,5fr)_minmax(420px,7fr)]"
			previewPanelClass += " xl:col-start-2 xl:row-span-2 xl:row-start-1"
			controlsPanelClass += " xl:col-start-1 xl:row-start-2"
			iframeClass += " min-h-[640px]"
		}
	}}
	<div class={editorClass}>
		<script type="application/json" data-theme-schema>@{ themeEditorSchemaValue() }</script>
		<div>
			<h1 class="text-3xl font-semibold tracking-tight">Theme editor</h1>
			<p class="mt-2 max-w-2xl text-sm text-muted-foreground">
				Choose the copied component style, then edit the semantic theme it consumes. The gallery renders the exact
				Nova or Maia component source a project receives from <code>gsxui add</code>.
			</p>
		</div>
		<div class={gridClass}>
			<section data-theme-style-panel class="flex min-w-0 flex-col gap-3">
				<div class="flex items-center justify-between gap-3">
					<h2 class="text-sm font-medium uppercase tracking-wide text-muted-foreground">Style</h2>
					<ui.Button data-theme-reset variant="outline" size="sm">Reset</ui.Button>
				</div>
				<div class="grid grid-cols-1 gap-3 sm:grid-cols-2">
					<button
						type="button"
						data-theme-style="nova"
						aria-pressed="true"
						class="rounded-xl border border-primary bg-accent/50 p-4 text-left transition-colors hover:bg-accent"
					>
						<span class="block font-medium">Nova</span>
						<span class="mt-1 block text-xs text-muted-foreground">Compact, practical defaults.</span>
					</button>
					<button
						type="button"
						data-theme-style="maia"
						aria-pressed="false"
						class="rounded-xl border border-border p-4 text-left transition-colors hover:bg-accent"
					>
						<span class="block font-medium">Maia</span>
						<span class="mt-1 block text-xs text-muted-foreground">Softer geometry and roomier controls.</span>
					</button>
				</div>
				<p class="text-xs text-muted-foreground">
					Both styles render the full component catalogue. The CLI refuses an unsafe mixed-style migration once
					components are installed.
				</p>
			</section>
			<div
				data-theme-preview-panel
				class={previewPanelClass}
			>
				<div class="flex items-center justify-between gap-3">
					<div>
						<h2 class="font-medium">Component preview</h2>
						<p data-theme-preview-status class="text-xs text-muted-foreground">Connecting to preview…</p>
					</div>
					<ui.Button data-theme-preview-retry variant="outline" size="sm" class="hidden">Retry</ui.Button>
				</div>
				<iframe
					data-theme-preview-frame
					title="Theme preview"
					src={previewURL}
					class={iframeClass}
				></iframe>
				<p data-theme-status role="status" aria-live="polite" class="min-h-5 text-sm text-muted-foreground"></p>
				<textarea
					data-theme-manual-copy
					readonly
					rows="5"
					class="hidden w-full rounded-md border border-input bg-background p-3 font-mono text-xs"
				></textarea>
			</div>
			<div
				data-theme-controls-panel
				class={controlsPanelClass}
			>
				<section class="flex flex-col gap-3 border-t border-border pt-6">
					<h2 class="text-sm font-medium uppercase tracking-wide text-muted-foreground">Mode and palette</h2>
					<div class="flex flex-wrap items-end justify-between gap-4">
						<div class="flex items-center gap-2">
							<button
								type="button"
								data-theme-mode-tab="light"
								aria-pressed="true"
								class={ tabBtnBase, "bg-accent text-accent-foreground" }
							>
								Light
							</button>
							<button
								type="button"
								data-theme-mode-tab="dark"
								aria-pressed="false"
								class={ tabBtnBase, "text-muted-foreground hover:bg-accent hover:text-accent-foreground" }
							>
								Dark
							</button>
						</div>
					</div>
					<div class="grid grid-cols-1 gap-3 sm:grid-cols-3">
						<ThemePicker
							name="baseColor"
							label="Base color"
							selected="neutral"
							choices={themePickerChoices(preset.BaseColorChoices())}
						/>
						<ThemePicker
							name="theme"
							label="Theme"
							selected="neutral"
							choices={themePickerChoices(mustThemePickerChoices("neutral"))}
						/>
						<ThemePicker
							name="radius"
							label="Radius"
							selected="medium"
							choices={themeRadiusPickerChoices(preset.RadiusChoices())}
						/>
					</div>
				</section>
				<section class="flex flex-col gap-3 border-t border-border pt-6">
					<h2 class="text-sm font-medium uppercase tracking-wide text-muted-foreground">Preset JSON</h2>
					<div class="flex flex-wrap gap-2">
						<ui.Button data-theme-copy="json" variant="outline" size="sm">Copy JSON</ui.Button>
						<ui.Button data-theme-download="json" variant="outline" size="sm">Download preset.json</ui.Button>
					</div>
					<textarea
						data-theme-import="json"
						rows="6"
						placeholder="Paste a gsxui preset JSON document"
						class="w-full rounded-md border border-input bg-transparent p-3 font-mono text-xs shadow-xs outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50"
					></textarea>
					<div>
						<ui.Button data-theme-import-apply="json" variant="outline" size="sm">Apply JSON</ui.Button>
					</div>
				</section>
				<section class="flex flex-col gap-3 border-t border-border pt-6">
					<h2 class="text-sm font-medium uppercase tracking-wide text-muted-foreground">Theme CSS</h2>
					<div class="flex flex-wrap gap-2">
						<ui.Button data-theme-copy="css" variant="outline" size="sm">Copy CSS</ui.Button>
						<ui.Button data-theme-download="css" variant="outline" size="sm">Download theme.css</ui.Button>
					</div>
					<textarea
						data-theme-import="css"
						rows="6"
						placeholder={themeImportPlaceholder}
						class="w-full rounded-md border border-input bg-transparent p-3 font-mono text-xs shadow-xs outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50"
					></textarea>
					<div>
						<ui.Button data-theme-import-apply="css" variant="outline" size="sm">Apply CSS</ui.Button>
					</div>
				</section>
				<section class="flex flex-col gap-3 border-t border-border pt-6">
					<h2 class="text-sm font-medium uppercase tracking-wide text-muted-foreground">Share and install</h2>
					<div class="flex flex-wrap gap-2">
						<ui.Button data-theme-copy="share" variant="outline" size="sm">Copy share code</ui.Button>
						<ui.Button data-theme-copy="url" variant="outline" size="sm">Copy share URL</ui.Button>
					</div>
					<label class="flex flex-col gap-1.5 text-xs text-muted-foreground">
						New project
						<textarea
							data-theme-command="init"
							readonly
							rows="3"
							class="rounded-md border border-input bg-muted/40 p-3 font-mono text-xs text-foreground"
						></textarea>
					</label>
					<label class="flex flex-col gap-1.5 text-xs text-muted-foreground">
						Initialized project
						<textarea
							data-theme-command="apply"
							readonly
							rows="3"
							class="rounded-md border border-input bg-muted/40 p-3 font-mono text-xs text-foreground"
						></textarea>
					</label>
				</section>
			</div>
		</div>
	</div>
}
