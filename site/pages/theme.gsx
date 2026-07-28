package pages

import (
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
		<ThemeEditor/>
	</Layout>
}

// ThemeEditor is the editor body without the site Layout, so the browser
// harness can exercise the production controls and web/theme.js directly.
component ThemeEditor() {
	<div class="flex flex-col gap-6 py-10">
		<div>
			<h1 class="text-3xl font-semibold tracking-tight">Theme editor</h1>
			<p class="mt-2 max-w-2xl text-sm text-muted-foreground">
				Edit the semantic CSS custom properties gsxui's components read. Paste a tweakcn/shadcn theme's root and dark
				blocks into Import to try it, or export a variables-only <code>theme.css</code>. Your project's entry,
				foundation, and component style files stay unchanged.
			</p>
		</div>
		<div class="grid grid-cols-1 gap-8 lg:grid-cols-2">
			<div class="flex flex-col gap-6">
				{ for _, g := range ThemeGroups() {
					<section class="flex flex-col gap-3">
						<h2 class="text-sm font-medium uppercase tracking-wide text-muted-foreground">{ g.Title }</h2>
						<div class="flex flex-col gap-2">
							<div class="grid grid-cols-[minmax(0,120px)_1fr_1fr] gap-3 text-xs text-muted-foreground">
								<span></span>
								<span>Light</span>
								<span>Dark</span>
							</div>
							{ for _, v := range g.Vars {
								<div>
									<div class="grid grid-cols-[minmax(0,120px)_1fr_1fr] items-center gap-3">
										<label class="truncate font-mono text-xs text-muted-foreground" title={v.Name}>{ v.Name }</label>
										<input
											type="text"
											data-theme-var={v.Name}
											data-theme-mode="light"
											value={v.Light}
											class="h-8 w-full min-w-0 rounded-md border border-input bg-transparent px-2 font-mono text-xs shadow-xs outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50"
										/>
										<input
											type="text"
											data-theme-var={v.Name}
											data-theme-mode="dark"
											value={v.Dark}
											class="h-8 w-full min-w-0 rounded-md border border-input bg-transparent px-2 font-mono text-xs shadow-xs outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50"
										/>
									</div>
									{ if v.Name == "--radius" {
										<p class="col-start-2 col-span-2 mt-1 text-xs text-muted-foreground">
											preview only — radius is theme-invariant in exports
										</p>
									} }
								</div>
							} }
						</div>
					</section>
				} }
				<section class="flex flex-col gap-3 border-t border-border pt-6">
					<h2 class="text-sm font-medium uppercase tracking-wide text-muted-foreground">Export</h2>
					<div class="flex flex-wrap gap-2">
						<ui.Button data-theme-copy variant="outline" size="sm">Copy CSS</ui.Button>
						<ui.Button data-theme-download variant="outline" size="sm">Download theme.css</ui.Button>
					</div>
					<textarea
						data-theme-export-output
						readonly
						rows="6"
						class="hidden w-full rounded-md border border-input bg-transparent p-2 font-mono text-xs shadow-xs outline-none"
					></textarea>
				</section>
				<section class="flex flex-col gap-3 border-t border-border pt-6">
					<h2 class="text-sm font-medium uppercase tracking-wide text-muted-foreground">Import</h2>
					<p class="text-xs text-muted-foreground">
						Paste a tweakcn/shadcn-style root/dark block of --var: value; pairs.
					</p>
					<textarea
						data-theme-import
						rows="6"
						placeholder={themeImportPlaceholder}
						class="w-full rounded-md border border-input bg-transparent p-2 font-mono text-xs shadow-xs outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50"
					></textarea>
					<div>
						<ui.Button data-theme-import-apply variant="outline" size="sm">Apply</ui.Button>
					</div>
				</section>
			</div>
			<div class="flex flex-col gap-4">
				<div class="flex items-center gap-2">
					<button
						type="button"
						data-theme-tab="light"
						aria-pressed="true"
						class={ tabBtnBase, "bg-accent text-accent-foreground" }
					>
						Light
					</button>
					<button
						type="button"
						data-theme-tab="dark"
						aria-pressed="false"
						class={ tabBtnBase, "text-muted-foreground hover:bg-accent hover:text-accent-foreground" }
					>
						Dark
					</button>
				</div>
				<div data-theme-preview class="flex flex-col gap-6 rounded-xl border border-border bg-background p-6">
					<div class="flex flex-wrap items-center gap-2">
						<ui.Button>Default</ui.Button>
						<ui.Button variant="secondary">Secondary</ui.Button>
						<ui.Button variant="outline">Outline</ui.Button>
						<ui.Button variant="ghost">Ghost</ui.Button>
						<ui.Button variant="link">Link</ui.Button>
						<ui.Button variant="destructive">Destructive</ui.Button>
					</div>
					<div class="flex flex-wrap items-center gap-2">
						<ui.Badge>Default</ui.Badge>
						<ui.Badge variant="secondary">Secondary</ui.Badge>
						<ui.Badge variant="outline">Outline</ui.Badge>
						<ui.Badge variant="destructive">Destructive</ui.Badge>
					</div>
					<ui.Card class="max-w-sm">
						<ui.CardHeader>
							<ui.CardTitle>Profile</ui.CardTitle>
							<ui.CardDescription>Preview restyles live as you edit the tokens.</ui.CardDescription>
						</ui.CardHeader>
						<ui.CardContent>
							<div class="flex flex-col gap-3">
								<div class="flex flex-col gap-1.5">
									<ui.Label for="theme-preview-name">Name</ui.Label>
									<ui.Input id="theme-preview-name" placeholder="Ada Lovelace"/>
								</div>
								<div class="flex items-center gap-2">
									<ui.Checkbox id="theme-preview-terms" checked/>
									<ui.Label for="theme-preview-terms">Accept terms</ui.Label>
								</div>
							</div>
						</ui.CardContent>
					</ui.Card>
					<ui.Alert>
						<ui.AlertTitle>Heads up</ui.AlertTitle>
						<ui.AlertDescription>This alert restyles with the tokens above.</ui.AlertDescription>
					</ui.Alert>
					<ui.Alert variant="destructive">
						<ui.AlertTitle>Something went wrong</ui.AlertTitle>
						<ui.AlertDescription>The destructive variant uses --destructive.</ui.AlertDescription>
					</ui.Alert>
				</div>
			</div>
		</div>
	</div>
}
