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
					gsxui components are copy-in: the CLI vendors real <code>.gsx</code> source into your own module. You can also
					import the gsxui package directly.
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
					For the unmodified npm/Vite scaffold produced by <code>gsx init --yes</code>, this is the complete
					setup. <code>gsxui init</code>:
				</p>
				<ul class="list-disc space-y-2 pl-6">
					<li>installs Tailwind CSS and wires it into <code>vite.config.ts</code> and <code>web/main.js</code></li>
					<li>vendors the gsxui CSS and JS entries into <code>web/gsxui/</code></li>
					<li>
						installs the class merger (<code>ui/merge/merge.go</code>) that lets caller classes override component
						styles (see <a href={Theming{} |> url}>Theming</a>)
					</li>
				</ul>
				<p>Rerunning it is safe — nothing is duplicated.</p>
				<div class="mt-4 flex flex-col gap-3">
					<docHeading item={gettingStartedTOCItems[2]}/>
					<p>
						If you customized the Vite config, entry file, package manager, or gsxui paths, <code>gsxui init</code>stops
						before writing anything and prints what to wire up yourself:
					</p>
					<pre><code>{ hl.Node("snippets/manual-integration") }</code></pre>
					<p>
						Not using Vite — or npm — at all? See <a href={NpmFree{} |> url}>npm-free</a>: <code>gsxui init</code>
						detects the missing scaffold and initializes without either.
					</p>
				</div>
			</section>
			<section class="flex flex-col gap-3">
				<docHeading item={gettingStartedTOCItems[3]}/>
				<pre><code>{ hl.Node("snippets/add.sh") }</code></pre>
				<pre><code>{ hl.Node("snippets/add.output") }</code></pre>
				<p>
					Dependencies come along automatically — <code>gsxui add native-select</code> also vendors <code>icon</code>.
					You own every file this writes: <code>gsxui add</code> never touches a file you've modified unless you
					pass <code>--overwrite</code>, which is also how you refresh components after upgrading
					the <code>gsxui</code> binary (discarding local edits to those files).
				</p>
			</section>
			<section class="flex flex-col gap-3">
				<docHeading item={gettingStartedTOCItems[4]}/>
				<p>
					<code>gsx init</code> already scaffolded a working app. Replace the <code>Index</code> component
					in <code>app.gsx</code> with a <code>Card</code> around a <code>Button</code>, adding
					the <code>ui</code> import:
				</p>
				<pre><code>{ hl.Node("snippets/first-page.gsx") }</code></pre>
				<p>Then start the development loop:</p>
				<pre><code>{ hl.Node("snippets/dev.sh") }</code></pre>
				<p>
					<code>gsx dev</code> watches your sources, rebuilds the server, and reloads the browser on save. Open the
					printed URL to see your Card and Button in gsxui's default theme.
					Next: <a href={Theming{} |> url}>restyle it</a>.
				</p>
			</section>
		</div>
	</siteLayout>
}
