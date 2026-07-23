package pages

import (
	"github.com/gsxhq/gsxui/ui/alert"
	"github.com/gsxhq/gsxui/ui/badge"
	"github.com/gsxhq/gsxui/ui/button"
	"github.com/gsxhq/gsxui/ui/card"
	"github.com/gsxhq/gsxui/ui/checkbox"
	"github.com/gsxhq/gsxui/ui/input"
	"github.com/gsxhq/gsxui/ui/label"
)

// Theme is the /theme route: a live editor over the 20 shadcn-compatible
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

// themeGroups holds the DEFAULT light/dark values for all 20 editable
// tokens, grouped the way the editor displays them. Source of truth:
// web/site.css's :root/.dark blocks (which mirror assets/gsxui.css, the
// file `gsxui init` vendors into consumer projects) — keep all three in
// sync when shadcn/tailwind bump the base palette.
var themeGroups = []themeGroup{
	{
		Title: "Base",
		Vars: []themeVar{
			{Name: "--background", Light: "oklch(1 0 0)", Dark: "oklch(0.145 0 0)"},
			{Name: "--foreground", Light: "oklch(0% 0 0)", Dark: "oklch(0.985 0 0)"},
			{Name: "--card", Light: "oklch(1 0 0)", Dark: "oklch(0.205 0 0)"},
			{Name: "--card-foreground", Light: "oklch(0% 0 0)", Dark: "oklch(0.985 0 0)"},
			{Name: "--popover", Light: "oklch(1 0 0)", Dark: "oklch(0.205 0 0)"},
			{Name: "--popover-foreground", Light: "oklch(0% 0 0)", Dark: "oklch(0.985 0 0)"},
		},
	},
	{
		Title: "Brand",
		Vars: []themeVar{
			{Name: "--primary", Light: "oklch(0% 0 0)", Dark: "oklch(0.922 0 0)"},
			{Name: "--primary-foreground", Light: "oklch(0.985 0 0)", Dark: "oklch(0.205 0 0)"},
			{Name: "--secondary", Light: "oklch(0.97 0 0)", Dark: "oklch(0.269 0 0)"},
			{Name: "--secondary-foreground", Light: "oklch(0.205 0 0)", Dark: "oklch(0.985 0 0)"},
			{Name: "--accent", Light: "oklch(0.97 0 0)", Dark: "oklch(0.371 0 0)"},
			{Name: "--accent-foreground", Light: "oklch(0.205 0 0)", Dark: "oklch(0.985 0 0)"},
		},
	},
	{
		Title: "Feedback",
		Vars: []themeVar{
			{Name: "--muted", Light: "oklch(0.97 0 0)", Dark: "oklch(0.269 0 0)"},
			{Name: "--muted-foreground", Light: "oklch(0.556 0 0)", Dark: "oklch(0.708 0 0)"},
			{Name: "--destructive", Light: "oklch(0.577 0.245 27.325)", Dark: "oklch(0.704 0.191 22.216)"},
			{Name: "--destructive-foreground", Light: "oklch(0.97 0.01 17)", Dark: "oklch(0.58 0.22 27)"},
		},
	},
	{
		Title: "Structure",
		Vars: []themeVar{
			{Name: "--border", Light: "oklch(0.922 0 0)", Dark: "oklch(1 0 0 / 10%)"},
			{Name: "--input", Light: "oklch(0.922 0 0)", Dark: "oklch(1 0 0 / 15%)"},
			{Name: "--ring", Light: "oklch(0.708 0 0)", Dark: "oklch(0.556 0 0)"},
			{Name: "--radius", Light: "0.625rem", Dark: "0.625rem"},
		},
	},
}

// ThemeGroups returns the default theme group definitions for testing purposes.
// This allows tests to verify that Go defaults stay in sync with CSS values.
func ThemeGroups() []themeGroup {
	return themeGroups
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
		<div class="flex flex-col gap-6 py-10">
			<div>
				<h1 class="text-3xl font-semibold tracking-tight">Theme editor</h1>
				<p class="mt-2 max-w-2xl text-sm text-muted-foreground">
					Edit the 20 CSS custom properties gsxui's components read (mirrors assets/gsxui.css).
					Paste a tweakcn/shadcn theme's root and dark blocks into Import to try it, or export what you build here as a ready-to-drop-in gsxui.css.
				</p>
			</div>
			<div class="grid grid-cols-1 gap-8 lg:grid-cols-2">
				<div class="flex flex-col gap-6">
					{ for _, g := range themeGroups {
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
											<label class="truncate font-mono text-xs text-muted-foreground" title={ v.Name }>{ v.Name }</label>
											<input
												type="text"
												data-theme-var={ v.Name }
												data-theme-mode="light"
												value={ v.Light }
												class="h-8 w-full min-w-0 rounded-md border border-input bg-transparent px-2 font-mono text-xs shadow-xs outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50"
											/>
											<input
												type="text"
												data-theme-var={ v.Name }
												data-theme-mode="dark"
												value={ v.Dark }
												class="h-8 w-full min-w-0 rounded-md border border-input bg-transparent px-2 font-mono text-xs shadow-xs outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50"
											/>
										</div>
										{ if v.Name == "--radius" {
											<p class="col-start-2 col-span-2 mt-1 text-xs text-muted-foreground">preview only — radius is theme-invariant in exports</p>
										} }
									</div>
								} }
							</div>
						</section>
					} }
					<section class="flex flex-col gap-3 border-t border-border pt-6">
						<h2 class="text-sm font-medium uppercase tracking-wide text-muted-foreground">Export</h2>
						<div class="flex flex-wrap gap-2">
							<button.Button data-theme-copy variant="outline" size="sm">Copy CSS</button.Button>
							<button.Button data-theme-download variant="outline" size="sm">Download gsxui.css</button.Button>
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
							placeholder={ themeImportPlaceholder }
							class="w-full rounded-md border border-input bg-transparent p-2 font-mono text-xs shadow-xs outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50"
						></textarea>
						<div>
							<button.Button data-theme-import-apply variant="outline" size="sm">Apply</button.Button>
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
							<button.Button>Default</button.Button>
							<button.Button variant="secondary">Secondary</button.Button>
							<button.Button variant="outline">Outline</button.Button>
							<button.Button variant="ghost">Ghost</button.Button>
							<button.Button variant="link">Link</button.Button>
							<button.Button variant="destructive">Destructive</button.Button>
						</div>
						<div class="flex flex-wrap items-center gap-2">
							<badge.Badge>Default</badge.Badge>
							<badge.Badge variant="secondary">Secondary</badge.Badge>
							<badge.Badge variant="outline">Outline</badge.Badge>
							<badge.Badge variant="destructive">Destructive</badge.Badge>
						</div>
						<card.Card class="max-w-sm">
							<card.CardHeader>
								<card.CardTitle>Profile</card.CardTitle>
								<card.CardDescription>Preview restyles live as you edit the tokens.</card.CardDescription>
							</card.CardHeader>
							<card.CardContent>
								<div class="flex flex-col gap-3">
									<div class="flex flex-col gap-1.5">
										<label.Label for="theme-preview-name">Name</label.Label>
										<input.Input id="theme-preview-name" placeholder="Ada Lovelace"/>
									</div>
									<div class="flex items-center gap-2">
										<checkbox.Checkbox id="theme-preview-terms" checked/>
										<label.Label for="theme-preview-terms">Accept terms</label.Label>
									</div>
								</div>
							</card.CardContent>
						</card.Card>
						<alert.Alert>
							<alert.AlertTitle>Heads up</alert.AlertTitle>
							<alert.AlertDescription>This alert restyles with the tokens above.</alert.AlertDescription>
						</alert.Alert>
						<alert.Alert variant="destructive">
							<alert.AlertTitle>Something went wrong</alert.AlertTitle>
							<alert.AlertDescription>The destructive variant uses --destructive.</alert.AlertDescription>
						</alert.Alert>
					</div>
				</div>
			</div>
		</div>
	</Layout>
}
