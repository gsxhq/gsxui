package pages

import "github.com/gsxhq/gsxui/site/hl"

// GettingStarted is the /docs/getting-started page: install → init → add →
// a first page served by the `gsx dev` loop, expanded from README.md with real CLI output
// (copied from internal/cli/init.go / add.go's actual printed strings, not
// invented).
type GettingStarted struct{}

var gettingStartedTOCItems = []docTOCItem{
	{ID: "install-cli", Title: "1. Install the CLIs", Depth: 2},
	{ID: "initialize-project", Title: "2. Initialize your project", Depth: 2},
	{ID: "manual-integration", Title: "Manual integration", Depth: 3},
	{ID: "add-components", Title: "3. Add components", Depth: 2},
	{ID: "first-page", Title: "4. Your first page", Depth: 2},
}

component (g GettingStarted) Page() {
	<siteLayout title="Getting Started" active="getting-started" mode={layoutDocs} toc={gettingStartedTOCItems}>
		<div data-doc="getting-started" class="flex max-w-3xl flex-col gap-10 py-10">
			<div class="flex flex-col gap-4">
				<h1 class="text-3xl font-semibold tracking-tight">Getting Started</h1>
				<p class="text-muted-foreground">
					gsxui components are copy-in: the CLI vendors real <code>.gsx</code> source into your own module, so what you
					build against is code you own and can edit — not a package you import and can't touch.
				</p>
			</div>
			<section class="flex flex-col gap-3">
				<docHeading item={gettingStartedTOCItems[0]}/>
				<pre><code>{ hl.Node("snippets/install.sh") }</code></pre>
			</section>
			<section class="flex flex-col gap-3">
				<docHeading item={gettingStartedTOCItems[1]}/>
				<p>Create a fresh GSX app, then initialize gsxui inside it:</p>
				<pre><code>{ hl.Node("snippets/init.sh") }</code></pre>
				<pre><code>{ hl.Node("snippets/init.output") }</code></pre>
				<p>
					For the unmodified npm/Vite scaffold produced by <code>gsx init --yes</code>, this is the complete setup.
					<code>gsxui init</code> installs Tailwind CSS, its Vite plugin, and <code>tw-animate-css</code>; registers
					<code>tailwindcss()</code> in <code>vite.config.ts</code>; and imports the gsxui CSS and behavior entries
					from <code>web/main.js</code>.
				</p>
				<p>
					This vendors the CSS entry (<code>web/gsxui/index.css</code>) plus its sibling
					<code>foundation.css</code>, <code>theme.css</code>, and <code>style.css</code>; the JS runtime and behavior
					barrel (
					<code>web/gsxui/</code>); and the class merger (<code>ui/merge/merge.go</code>), then
					points <code>gsx.toml</code>'s <code>class_merger</code> at it — the seam that makes caller-class-merge work
					(see <a href={Theming{} |> url}>Theming</a>). It
					also <code>go get</code> <code>gsx</code> and <code>tailwind-merge-go</code>, and installs
					the <code>gsx</code> tool via <code>go get -tool</code>.
				</p>
				<p>
					Rerunning <code>gsxui init</code> is safe: npm verifies the locked dependencies and the exact scaffold
					integration is not duplicated.
				</p>
				<div class="mt-4 flex flex-col gap-3">
					<docHeading item={gettingStartedTOCItems[2]}/>
					<p>
						If you changed the Vite config, entry file, package manager, or gsxui JS/CSS paths, automatic rewriting
						stops before running commands or writing files. Keep your custom structure and apply the printed
						responsibilities yourself:
					</p>
					<pre><code>{ hl.Node("snippets/manual-integration") }</code></pre>
				</div>
			</section>
			<section class="flex flex-col gap-3">
				<docHeading item={gettingStartedTOCItems[3]}/>
				<pre><code>{ hl.Node("snippets/add.sh") }</code></pre>
				<pre><code>{ hl.Node("snippets/add.output") }</code></pre>
				<p>
					<code>card</code> has no dependencies of its own, but a component that does (e.g. <code>native-select</code>,
					which needs <code>icon</code>) pulls its dependency in automatically
					— <code>gsxui add native-select</code> vendors <code>icon</code> too. You own every file this
					writes: <code>gsxui add</code> never touches one you've already modified unless you
					pass <code>--overwrite</code>. After upgrading the <code>gsxui</code> binary,
					re-run <code>gsxui add &lt;name&gt; --overwrite</code> to refresh vendored components — that discards local
					edits to those files.
				</p>
			</section>
			<section class="flex flex-col gap-3">
				<docHeading item={gettingStartedTOCItems[4]}/>
				<p>
					<code>gsx init</code> already scaffolded a working app: <code>app.gsx</code> holds
					a <code>Layout</code> and an <code>Index</code> component, and <code>main.go</code> serves it. Make it your
					first gsxui page by rendering a <code>Card</code> around a <code>Button</code> — replace
					the <code>Index</code> component in <code>app.gsx</code> and add the <code>ui</code> import:
				</p>
				<pre><code>{ hl.Node("snippets/first-page.gsx") }</code></pre>
				<p>Then start the development loop:</p>
				<pre><code>{ hl.Node("snippets/dev.sh") }</code></pre>
				<p>
					The scaffold's <code>npm run dev</code> is a one-line wrapper around the same command. <code>gsx dev</code>
					watches your sources, regenerates <code>.gsx</code> to Go, rebuilds and swaps the server, and reloads the
					browser — edit <code>app.gsx</code>, save, and the page updates.
				</p>
				<p>
					Open the printed URL — a styled Card with a Button inside, rendered with gsxui's default
					light theme. Next: <a href={Theming{} |> url}>restyle it</a>.
				</p>
			</section>
		</div>
	</siteLayout>
}
