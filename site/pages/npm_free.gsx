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
					Nothing in gsxui requires Vite — or npm. The component behaviors are dependency-free native ES modules, and
					the one CSS dependency is vendored, so the stylesheet builds without <code>node_modules</code>. (The
					recommended setup is still the scaffold from <a href={GettingStarted{} |> url}>Getting Started</a>.)
				</p>
			</div>
			<section class="flex flex-col gap-3">
				<docHeading item={npmFreeTOCItems[0]}/>
				<p>
					If <code>gsxui init</code> finds neither <code>vite.config.ts</code> nor <code>web/main.js</code>, it
					initializes in npm-free mode: no npm commands run, no <code>package.json</code> is written. The same CSS and
					JS entries are vendored and the Go tooling is installed as usual.
				</p>
				<pre><code>{ hl.Node("snippets/nonvite-init.output") }</code></pre>
			</section>
			<section class="flex flex-col gap-3">
				<docHeading item={npmFreeTOCItems[1]}/>
				<p>
					You own serving and CSS building. Serve the vendored JS directory statically and load the barrel with one
					module script tag — no bundler required, any bundler welcome:
				</p>
				<pre><code>{ hl.Node("snippets/nonvite-serve.go") }</code></pre>
			</section>
			<section class="flex flex-col gap-3">
				<docHeading item={npmFreeTOCItems[2]}/>
				<p>
					Build the CSS entry with any Tailwind v4 tool.
					The <a href="https://tailwindcss.com/docs/installation/tailwind-cli">standalone Tailwind CLI</a> — a single
					binary from <a href="https://github.com/tailwindlabs/tailwindcss/releases/latest">GitHub releases</a> — is the
					natural fit; npm's <code>@tailwindcss/cli</code> takes the same flags. Then link the output from your pages:
				</p>
				<pre><code>{ hl.Node("snippets/nonvite-css.sh") }</code></pre>
			</section>
		</div>
	</siteLayout>
}
