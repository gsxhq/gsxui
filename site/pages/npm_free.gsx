package pages

import "github.com/gsxhq/gsxui/site/hl"

// NpmFree is the /docs/npm-free page: running gsxui without Vite — or npm —
// with real CLI output (copied from internal/cli/init.go's actual printed
// strings, not invented; a drift test pins the init summary snippet).
type NpmFree struct{}

var npmFreeTOCItems = []docTOCItem{
	{ID: "initialize", Title: "Initialize without a scaffold", Depth: 2},
	{ID: "serve-scripts", Title: "Serve the scripts", Depth: 2},
	{ID: "build-css", Title: "Build the CSS", Depth: 2},
}

component (n NpmFree) Page() {
	<siteLayout title="npm-free" active="npm-free" mode={layoutDocs} toc={npmFreeTOCItems}>
		<div data-doc="npm-free" class="flex max-w-3xl flex-col gap-10 py-10">
			<div class="flex flex-col gap-4">
				<h1 class="text-3xl font-semibold tracking-tight">npm-free</h1>
				<p class="text-muted-foreground">
					The recommended setup is the gsx scaffold with gsxui on top — <code>gsx init</code> then
					<code>gsxui init</code>, as <a href={GettingStarted{} |> url}>Getting Started</a> shows. But nothing in
					gsxui requires Vite — or npm: the component behaviors are dependency-free native ES modules, and the one
					CSS dependency (<code>tw-animate-css</code>) is vendored as <code>web/gsxui/animate.css</code> so the
					stylesheet builds without <code>node_modules</code>.
				</p>
			</div>
			<section class="flex flex-col gap-3">
				<docHeading item={npmFreeTOCItems[0]}/>
				<p>
					If <code>gsxui init</code> finds neither <code>vite.config.ts</code> nor <code>web/main.js</code>, it
					initializes in npm-free mode: no npm commands run and no <code>package.json</code> is written. Everything
					else matches the <a href={GettingStarted{} |> url}>standard setup</a> — the same CSS and JS entries are
					vendored, and the Go tooling is installed the same way.
				</p>
				<pre><code>{ hl.Node("snippets/nonvite-init.output") }</code></pre>
			</section>
			<section class="flex flex-col gap-3">
				<docHeading item={npmFreeTOCItems[1]}/>
				<p>
					You own serving and CSS building. Serve the vendored JS directory statically and load the barrel
					with one module script tag — no bundler required, any bundler welcome:
				</p>
				<pre><code>{ hl.Node("snippets/nonvite-serve.go") }</code></pre>
			</section>
			<section class="flex flex-col gap-3">
				<docHeading item={npmFreeTOCItems[2]}/>
				<p>
					Build the CSS entry with any Tailwind v4 tool — gsxui never installs or manages your build tooling.
					True to this page's title, the <a href="https://tailwindcss.com/docs/installation/tailwind-cli">standalone
					Tailwind CLI</a> is the natural fit: a single binary
					from <a href="https://github.com/tailwindlabs/tailwindcss/releases/latest">GitHub releases</a>, no npm
					involved. If you already have npm, its <code>@tailwindcss/cli</code> package takes the same flags. Then
					link the output from your pages:
				</p>
				<pre><code>{ hl.Node("snippets/nonvite-css.sh") }</code></pre>
			</section>
		</div>
	</siteLayout>
}
