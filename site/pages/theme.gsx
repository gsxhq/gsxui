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

// themeVar is one editable CSS custom property, with its default light and
// dark values.
type themeVar struct {
	Name  string
	Light string
	Dark  string
}

// themeGroup is a labeled cluster of related vars (matches the editor's
// section headings).
type themeGroup struct {
	Title string
	Vars  []themeVar
}

type themeEditorSchema struct {
	Schema          string                     `json:"schema"`
	SchemaVersion   int                        `json:"schemaVersion"`
	TransportPrefix string                     `json:"transportPrefix"`
	TokenNames      []string                   `json:"tokenNames"`
	RadiusUnits     []string                   `json:"radiusUnits"`
	Styles          []string                   `json:"styles"`
	Defaults        map[string]json.RawMessage `json:"defaults"`
	CanonicalDefaults map[string]string        `json:"canonicalDefaults"`
}

func themeEditorSchemaValue() themeEditorSchema {
	styles := preset.Styles()
	schema := themeEditorSchema{
		Schema:          preset.SchemaURL,
		SchemaVersion:   preset.SchemaVersion,
		TransportPrefix: "gsxui:v1:",
		TokenNames:      preset.TokenNames(),
		RadiusUnits:     preset.RadiusUnits(),
		Styles:          make([]string, len(styles)),
		Defaults:        make(map[string]json.RawMessage, len(styles)),
		CanonicalDefaults: make(map[string]string, len(styles)),
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

// ThemeGroups derives the editor's presentation groups and server-rendered
// defaults from the authoritative preset schema.
func ThemeGroups() []themeGroup {
	defaults := preset.Default(preset.StyleNova)
	groupIndexes := make(map[string]int)
	groups := make([]themeGroup, 0, len(preset.GroupNames()))
	for _, name := range preset.GroupNames() {
		groupIndexes[name] = len(groups)
		groups = append(groups, themeGroup{Title: name})
	}
	add := func(definition preset.TokenDefinition, light, dark string) {
		index, ok := groupIndexes[definition.Group]
		if !ok {
			index = len(groups)
			groupIndexes[definition.Group] = index
			groups = append(groups, themeGroup{Title: definition.Group})
		}
		groups[index].Vars = append(groups[index].Vars, themeVar{
			Name:  "--" + definition.Name,
			Light: light,
			Dark:  dark,
		})
	}
	for _, definition := range preset.TokenDefinitions() {
		add(
			definition,
			defaults.Theme.Light[definition.Name],
			defaults.Theme.Dark[definition.Name],
		)
	}
	radius := preset.RadiusDefinition()
	add(radius, defaults.Radius, defaults.Radius)
	return groups
}

const tabBtnBase = "rounded-md border border-border px-3 py-1.5 text-sm font-medium transition-colors"

const themeImportPlaceholder = `:root {
  --primary: oklch(0.6 0.2 280);
}
.dark {
  --primary: oklch(0.7 0.2 280);
}`

component (t Theme) Page() {
	<Layout title="Theme" active="">
		<ThemeEditor previewURL={ThemePreviewButton{} |> url}/>
	</Layout>
}

// ThemeEditor is the editor body without the site Layout, so the browser
// harness can exercise the production controls and web/theme.js directly.
component ThemeEditor(previewURL string) {
	<div class="flex flex-col gap-8 py-10">
		<script type="application/json" data-theme-schema>@{ themeEditorSchemaValue() }</script>
		<div>
			<h1 class="text-3xl font-semibold tracking-tight">Theme editor</h1>
			<p class="mt-2 max-w-2xl text-sm text-muted-foreground">
				Choose the copied component style, then edit the semantic theme it consumes. This Button pilot renders the
				exact Nova or Maia source a project receives from <code>gsxui add</code>.
			</p>
		</div>
		<div class="grid grid-cols-1 gap-8 xl:grid-cols-[minmax(0,5fr)_minmax(420px,7fr)]">
			<div class="flex min-w-0 flex-col gap-7">
				<section class="flex flex-col gap-3">
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
							<span class="flex items-center justify-between gap-2">
								<span class="font-medium">Maia</span>
								<span class="rounded-full bg-muted px-2 py-0.5 text-[11px] text-muted-foreground">Button pilot</span>
							</span>
							<span class="mt-1 block text-xs text-muted-foreground">Softer geometry and roomier controls.</span>
						</button>
					</div>
					<p class="text-xs text-muted-foreground">
						Maia currently applies only to Button. The CLI refuses an unsafe mixed-style migration once other
						components are installed.
					</p>
				</section>

				<section class="flex flex-col gap-3 border-t border-border pt-6">
					<h2 class="text-sm font-medium uppercase tracking-wide text-muted-foreground">Mode and radius</h2>
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
						<label class="flex min-w-[220px] flex-col gap-1.5 text-xs text-muted-foreground">
							Radius
							<input
								type="text"
								data-theme-field="radius"
								value={preset.Default(preset.StyleNova).Radius}
								class="h-9 rounded-md border border-input bg-transparent px-3 font-mono text-xs shadow-xs outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50"
							/>
							<span data-theme-field-error="radius" class="hidden text-destructive"></span>
						</label>
					</div>
				</section>

				{ for _, g := range ThemeGroups() {
					<section data-theme-group={g.Title} class="flex flex-col gap-3 border-t border-border pt-6">
						<h2 class="text-sm font-medium uppercase tracking-wide text-muted-foreground">{ g.Title }</h2>
						<div class="flex flex-col gap-2">
							{ for _, v := range g.Vars {
								{ if v.Name != "--radius" {
									<div class="grid grid-cols-[minmax(0,150px)_1fr] items-start gap-3">
										<label class="truncate pt-2 font-mono text-xs text-muted-foreground" title={v.Name}>{ v.Name }</label>
										<div data-theme-mode-field="light">
											<input
												type="text"
												data-theme-var={v.Name}
												data-theme-mode="light"
												data-theme-field={"light." + strings.TrimPrefix(v.Name, "--")}
												value={v.Light}
												class="h-9 w-full min-w-0 rounded-md border border-input bg-transparent px-3 font-mono text-xs shadow-xs outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50"
											/>
											<span data-theme-field-error={"light." + strings.TrimPrefix(v.Name, "--")} class="mt-1 hidden text-xs text-destructive"></span>
										</div>
										<div data-theme-mode-field="dark" hidden>
											<input
												type="text"
												data-theme-var={v.Name}
												data-theme-mode="dark"
												data-theme-field={"dark." + strings.TrimPrefix(v.Name, "--")}
												value={v.Dark}
												class="h-9 w-full min-w-0 rounded-md border border-input bg-transparent px-3 font-mono text-xs shadow-xs outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50"
											/>
											<span data-theme-field-error={"dark." + strings.TrimPrefix(v.Name, "--")} class="mt-1 hidden text-xs text-destructive"></span>
										</div>
									</div>
								} }
							} }
						</div>
					</section>
				} }

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
						<textarea data-theme-command="init" readonly rows="3" class="rounded-md border border-input bg-muted/40 p-3 font-mono text-xs text-foreground"></textarea>
					</label>
					<label class="flex flex-col gap-1.5 text-xs text-muted-foreground">
						Initialized project
						<textarea data-theme-command="apply" readonly rows="3" class="rounded-md border border-input bg-muted/40 p-3 font-mono text-xs text-foreground"></textarea>
					</label>
				</section>
			</div>

			<div class="flex min-w-0 flex-col gap-4 xl:sticky xl:top-20 xl:self-start">
				<div class="flex items-center justify-between gap-3">
					<div>
						<h2 class="font-medium">Button preview</h2>
						<p data-theme-preview-status class="text-xs text-muted-foreground">Connecting to preview…</p>
					</div>
					<ui.Button data-theme-preview-retry variant="outline" size="sm" class="hidden">Retry</ui.Button>
				</div>
				<iframe
					data-theme-preview-frame
					title="Button theme preview"
					src={previewURL}
					class="min-h-[640px] w-full rounded-xl border border-border bg-background shadow-sm"
				></iframe>
				<p data-theme-status role="status" aria-live="polite" class="min-h-5 text-sm text-muted-foreground"></p>
				<textarea
					data-theme-manual-copy
					readonly
					rows="5"
					class="hidden w-full rounded-md border border-input bg-background p-3 font-mono text-xs"
				></textarea>
			</div>
		</div>
	</div>
}
