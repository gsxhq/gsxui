package pages

import (
	"github.com/gsxhq/gsxui/site/examples/showcase"
	"github.com/gsxhq/gsxui/ui"
)

// Home is the site's landing page.
type Home struct{}

const installSnippet = `gsx init app --yes
cd app
gsxui init
gsxui add button`

component (h Home) Page() {
	<siteLayout title="gsxui" active="" mode={layoutMarketing} toc={nil}>
		<section class="flex flex-col gap-4 py-10">
			<h1 class="text-4xl font-semibold tracking-tight sm:text-5xl">
				Components for modern web frontends in Go.
			</h1>
			<p class="max-w-2xl text-lg text-muted-foreground">
				gsxui is a shadcn-style component set for{ " " }
				<a
					href="https://gsxhq.github.io/"
					target="_blank"
					rel="noreferrer"
					class="underline underline-offset-4 hover:text-foreground"
				>
					gsx
				</a>
				: copy-in, type-checked, server-rendered. Own the code, style it with Tailwind, ship it.
			</p>
			<div class="flex flex-wrap items-center gap-3 pt-2">
				<ui.Button size="lg" href={Theme{} |> url}>Create</ui.Button>
				<ui.Button size="lg" variant="outline" href={GettingStarted{} |> url}>Get Started</ui.Button>
				<ui.Button
					size="lg"
					variant="outline"
					href="https://github.com/gsxhq/gsxui"
					target="_blank"
					rel="noreferrer"
				>
					View on GitHub
				</ui.Button>
				<ui.Button
					size="lg"
					variant="ghost"
					href="https://gsxhq.github.io/"
					target="_blank"
					rel="noreferrer"
				>
					gsx documentation ↗
				</ui.Button>
			</div>
		</section>
		<section class="pb-10">
			<pre
				class="overflow-x-auto rounded-lg border border-border bg-card p-4 text-sm text-card-foreground"
			><code>{ installSnippet }</code></pre>
		</section>
		<section id="components" class="flex flex-col gap-6 border-t border-border py-10">
			<div class="flex flex-wrap items-end justify-between gap-3">
				<h2 class="text-2xl font-semibold tracking-tight">Built with gsxui</h2>
				<a
					href={ComponentsIndex{} |> url}
					class="text-sm text-muted-foreground underline-offset-4 hover:text-foreground hover:underline"
				>
					Browse all components →
				</a>
			</div>
			<div class="grid gap-4 md:grid-cols-2">
				<showcase.SignInCard/>
				<showcase.SettingsCard/>
				<showcase.StatsCard/>
				<showcase.OverlaysCard/>
			</div>
		</section>
	</siteLayout>
}
