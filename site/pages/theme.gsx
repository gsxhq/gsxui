package pages

import (
	"encoding/json"
	"strings"

	"github.com/gsxhq/gsxui/internal/preset"
	"github.com/gsxhq/gsxui/ui"
	"github.com/gsxhq/gsxui/ui/icon"
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
	Styles      []string `json:"styles"`
	BaseColors  []string `json:"baseColors"`
	Themes      []string `json:"themes"`
	Radii       []string `json:"radii"`
	MenuAccents []string `json:"menuAccents"`
	Fonts       []string `json:"fonts"`
}

type themePaletteSchema struct {
	BaseColors  []themePaletteChoiceSchema            `json:"baseColors"`
	Themes      map[string][]themePaletteChoiceSchema `json:"themes"`
	Radii       []themeRadiusChoiceSchema             `json:"radii"`
	MenuAccents []themeMenuAccentChoiceSchema         `json:"menuAccents"`
	// ChartColors is flat by hue name (not nested by base color like
	// Resolved): the chart color axis overlays independently of base
	// color, and its picker reuses Themes[baseColor]'s own choices for
	// display, so this table only needs to answer "what are chart-1..5 for
	// this name" regardless of which base color's picker view showed it.
	ChartColors map[string]themePaletteResolvedSchema            `json:"chartColors"`
	Resolved    map[string]map[string]themePaletteResolvedSchema `json:"resolved"`
	// Fonts backs BOTH the body (FontSans) and heading (FontHeading)
	// pickers — same list, mirrors catalog.go's FontChoices doc comment.
	Fonts            []themeFontChoiceSchema     `json:"fonts"`
	DefaultSelection themePaletteSelectionSchema `json:"defaultSelection"`
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

type themeMenuAccentChoiceSchema struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}

// themeFontChoiceSchema is the font picker's choice shape: Stack is the
// literal CSS font-family value theme-state.js's selectFont/selectFontHeading
// commit into resolved.fontSans/fontHeading — see catalog.go's FontChoice.
type themeFontChoiceSchema struct {
	Name  string `json:"name"`
	Title string `json:"title"`
	Stack string `json:"stack"`
}

type themePaletteResolvedSchema struct {
	Light preset.ThemeValues `json:"light"`
	Dark  preset.ThemeValues `json:"dark"`
}

type themePaletteSelectionSchema struct {
	BaseColor   string `json:"baseColor"`
	Theme       string `json:"theme"`
	Radius      string `json:"radius"`
	ChartColor  string `json:"chartColor"`
	MenuAccent  string `json:"menuAccent"`
	FontSans    string `json:"fontSans"`
	FontHeading string `json:"fontHeading"`
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
			Styles:      styles,
			BaseColors:  transport.BaseColors,
			Themes:      transport.Themes,
			Radii:       transport.Radii,
			MenuAccents: transport.MenuAccents,
			Fonts:       transport.Fonts,
		},
	}
}

func themePaletteSchemaValue() themePaletteSchema {
	selection := preset.DefaultPaletteSelection()
	schema := themePaletteSchema{
		BaseColors:       themePaletteChoiceSchemaValues(preset.BaseColorChoices()),
		Themes:           make(map[string][]themePaletteChoiceSchema),
		Radii:            themeRadiusChoiceSchemaValues(preset.RadiusChoices()),
		MenuAccents:      themeMenuAccentChoiceSchemaValues(preset.MenuAccentChoices()),
		ChartColors:      make(map[string]themePaletteResolvedSchema),
		Resolved:         make(map[string]map[string]themePaletteResolvedSchema),
		Fonts:            themeFontChoiceSchemaValues(preset.FontChoices()),
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
			// Resolved swatches always use the default chart color and
			// menu accent ("neutral"/"subtle"): this matrix backs the base
			// color/theme pickers' hue previews, not the independent chart
			// color or menu-accent controls.
			resolved, err := preset.ResolvePalette(preset.StyleNova, preset.PaletteSelection{
				BaseColor:   baseColor.Name,
				Theme:       theme.Name,
				Radius:      selection.Radius,
				ChartColor:  selection.ChartColor,
				MenuAccent:  selection.MenuAccent,
				FontSans:    selection.FontSans,
				FontHeading: selection.FontHeading,
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
	for _, chartColor := range preset.ChartColorChoices() {
		light, dark, err := preset.ResolveChartColor(chartColor.Name)
		if err != nil {
			panic(err)
		}
		schema.ChartColors[chartColor.Name] = themePaletteResolvedSchema{Light: light, Dark: dark}
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

func themeFontChoiceSchemaValues(choices []preset.FontChoice) []themeFontChoiceSchema {
	values := make([]themeFontChoiceSchema, len(choices))
	for i, choice := range choices {
		values[i] = themeFontChoiceSchema{Name: choice.Name, Title: choice.Title, Stack: choice.Stack}
	}
	return values
}

func themeMenuAccentChoiceSchemaValues(choices []preset.MenuAccentChoice) []themeMenuAccentChoiceSchema {
	values := make([]themeMenuAccentChoiceSchema, len(choices))
	for i, choice := range choices {
		values[i] = themeMenuAccentChoiceSchema{Name: choice.Name, Title: choice.Title}
	}
	return values
}

func themePaletteSelectionSchemaFromPreset(selection preset.PaletteSelection) themePaletteSelectionSchema {
	return themePaletteSelectionSchema{
		BaseColor:   selection.BaseColor,
		Theme:       selection.Theme,
		Radius:      selection.Radius,
		ChartColor:  selection.ChartColor,
		MenuAccent:  selection.MenuAccent,
		FontSans:    selection.FontSans,
		FontHeading: selection.FontHeading,
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

// themeFontPickerChoices adapts FontChoices() to ThemePicker's generic
// choice shape. No SwatchStyle: a font isn't a hue or a radius, so the
// picker's swatch circle simply stays hidden for every option (renderSwatch
// in theme.js already hides it whenever a choice carries neither Swatch nor
// SwatchStyle).
func themeFontPickerChoices(choices []preset.FontChoice) []themePickerChoice {
	values := make([]themePickerChoice, len(choices))
	for i, choice := range choices {
		values[i] = themePickerChoice{Title: choice.Title, Value: choice.Name}
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

// themeStylePressedAttr and themeStyleButtonClass render the style panel's
// JS-less default: web/theme.js's own render() re-syncs both aria-pressed
// and these two classes against state.resolved.style on load and on every
// click (data-theme-style is generic there — no style name is hardcoded),
// so this only has to get the page's FIRST paint right. preset.StyleNova
// is the real system-wide default (preset.Default's own fallback, and
// theme-state.js's schema.defaults.nova), not merely StyleChoices()'s
// first entry (that's vega, the plan's canonical listing order).
func themeStylePressedAttr(name string) string {
	if name == string(preset.StyleNova) {
		return "true"
	}
	return "false"
}

func themeStyleButtonClass(name string) string {
	if name == string(preset.StyleNova) {
		return "flex flex-col items-start rounded-xl border border-primary bg-accent/50 p-4 text-left transition-colors hover:bg-accent"
	}
	return "flex flex-col items-start rounded-xl border border-border p-4 text-left transition-colors hover:bg-accent"
}

const tabBtnBase = "rounded-md border border-border px-3 py-1.5 text-sm font-medium transition-colors"

// themePageChoices backs the preview panel's page tab pair — the two
// site/stylepreview/gallery.gsx.src page groups (theme-creator-parity
// Task 4c), a content-navigation control over the preview document, not a
// theme axis, so unlike Style/Palette it needs no preset.go catalog: the
// two names and titles are hardcoded here the same way Mode's light/dark
// tab pair already is below.
type themePageChoice struct {
	Name  string
	Title string
}

func themePageChoices() []themePageChoice {
	return []themePageChoice{
		{Name: "components", Title: "Components"},
		{Name: "product", Title: "Product"},
	}
}

// themePagePressedAttr and themePageTabClass are the page tab pair's
// JS-less first-paint default, the same pattern as
// themeMenuAccentPressedAttr/themeMenuAccentTabClass above: web/theme.js's
// render() re-syncs both against state.page afterward. "components" is
// createThemeState's own default (web/theme-state.js).
func themePagePressedAttr(name string) string {
	if name == "components" {
		return "true"
	}
	return "false"
}

func themePageTabClass(name string) string {
	if name == "components" {
		return tabBtnBase + " bg-accent text-accent-foreground"
	}
	return tabBtnBase + " text-muted-foreground hover:bg-accent hover:text-accent-foreground"
}

// themeMenuAccentPressedAttr and themeMenuAccentTabClass are the menu
// accent tab pair's JS-less first-paint default, the same pattern as
// themeStylePressedAttr/themeStyleButtonClass above: web/theme.js's
// render() re-syncs both against state.selection.menuAccent afterward.
// "subtle" is DefaultPaletteSelection()'s own default.
func themeMenuAccentPressedAttr(name string) string {
	if name == "subtle" {
		return "true"
	}
	return "false"
}

func themeMenuAccentTabClass(name string) string {
	if name == "subtle" {
		return tabBtnBase + " bg-accent text-accent-foreground"
	}
	return tabBtnBase + " text-muted-foreground hover:bg-accent hover:text-accent-foreground"
}

const themeImportPlaceholder = `:root {
  --primary: oklch(0.6 0.2 280);
}
.dark {
  --primary: oklch(0.7 0.2 280);
}`

component (t Theme) Page() {
	<siteLayout title="Theme" active="" mode={layoutWorkspace}>
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
	<div class={ editorClass }>
		<script type="application/json" data-theme-schema>@{ themeEditorSchemaValue() }</script>
		<div>
			<h1 class="text-3xl font-semibold tracking-tight">Theme editor</h1>
			<p class="mt-2 max-w-2xl text-sm text-muted-foreground">
				Pick one of 8 component styles, then shape the semantic theme it renders with: base color, accent, chart color,
				menu accent, radius, and a body/heading font pair, all previewed live in light and dark. Undo/redo, the address
				bar, and share codes track every change, and a foreign <code>theme.css</code> imports straight in. The gallery
				renders the exact component source a project receives from <code>gsxui add</code> under the selected style.
			</p>
		</div>
		<div class={ gridClass }>
			<section data-theme-style-panel class="flex min-w-0 flex-col gap-3">
				<div class="flex items-center justify-between gap-3">
					<h2 class="text-sm font-medium uppercase tracking-wide text-muted-foreground">Style</h2>
					<div class="flex items-center gap-1">
						<ui.Button data-theme-undo variant="outline" size="icon-sm" aria-label="Undo" title="Undo (⌘Z)" disabled>
							<icon.Undo/>
						</ui.Button>
						<ui.Button data-theme-redo variant="outline" size="icon-sm" aria-label="Redo" title="Redo (⌘⇧Z)" disabled>
							<icon.Redo/>
						</ui.Button>
						<ui.Button data-theme-reset variant="outline" size="sm">Reset</ui.Button>
					</div>
				</div>
				<div class="grid grid-cols-2 gap-3 sm:grid-cols-4">
					{ for _, choice := range preset.StyleChoices() {
						<button
							type="button"
							data-theme-style={choice.Name}
							aria-pressed={themeStylePressedAttr(choice.Name)}
							class={ themeStyleButtonClass(choice.Name) }
						>
							<span class="block font-medium">{ choice.Title }</span>
							<span class="mt-1 block text-xs text-muted-foreground">{ choice.Description }</span>
						</button>
					} }
				</div>
				<p class="text-xs text-muted-foreground">
					All 8 styles render the full component catalogue. <code>gsxui add</code> and <code>gsxui apply</code> vendor
					each component's exact source for the selected style.
				</p>
			</section>
			<div
				data-theme-preview-panel
				class={ previewPanelClass }
			>
				<div class="flex flex-wrap items-center justify-between gap-3">
					<div>
						<h2 class="font-medium">Component preview</h2>
						<p data-theme-preview-status class="text-xs text-muted-foreground">Connecting to preview…</p>
					</div>
					<div class="flex items-center gap-2">
						<div class="flex items-center gap-2">
							{ for _, choice := range themePageChoices() {
								<button
									type="button"
									data-theme-page-tab={choice.Name}
									aria-pressed={themePagePressedAttr(choice.Name)}
									class={ themePageTabClass(choice.Name) }
								>
									{ choice.Title }
								</button>
							} }
						</div>
						<ui.Button data-theme-preview-retry variant="outline" size="sm" class="hidden">Retry</ui.Button>
					</div>
				</div>
				<iframe
					data-theme-preview-frame
					title="Theme preview"
					src={previewURL}
					class={ iframeClass }
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
				class={ controlsPanelClass }
			>
				<section class="flex flex-col gap-3 border-t border-border pt-6">
					<h2 class="text-sm font-medium uppercase tracking-wide text-muted-foreground">Mode and palette</h2>
					<div class="flex flex-wrap items-end justify-between gap-4">
						<div class="flex flex-col gap-1.5">
							<span class="text-xs text-muted-foreground">Mode</span>
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
						<div class="flex flex-col gap-1.5">
							<span class="text-xs text-muted-foreground">Menu accent</span>
							<div class="flex items-center gap-2">
								{ for _, choice := range preset.MenuAccentChoices() {
									<button
										type="button"
										data-theme-menu-accent-tab={choice.Name}
										aria-pressed={themeMenuAccentPressedAttr(choice.Name)}
										class={ themeMenuAccentTabClass(choice.Name) }
									>
										{ choice.Title }
									</button>
								} }
							</div>
						</div>
					</div>
					<div class="grid grid-cols-2 gap-3">
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
						<ThemePicker
							name="chartColor"
							label="Chart color"
							selected="neutral"
							choices={themePickerChoices(mustThemePickerChoices("neutral"))}
						/>
						<ThemePicker
							name="fontSans"
							label="Body font"
							selected="geist"
							choices={themeFontPickerChoices(preset.FontChoices())}
						/>
						<ThemePicker
							name="fontHeading"
							label="Heading font"
							selected="geist"
							choices={themeFontPickerChoices(preset.FontChoices())}
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
