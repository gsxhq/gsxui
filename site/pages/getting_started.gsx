package pages

import "github.com/gsxhq/gsxui/site/hl"

// GettingStarted is the /docs/getting-started page: install → init → add →
// a minimal first page, expanded from README.md with real CLI output
// (copied from internal/cli/init.go / add.go's actual printed strings, not
// invented).
type GettingStarted struct{}

component (g GettingStarted) Page() {
	<Layout title="Getting Started" active="getting-started">
		<div data-doc="getting-started" class="flex max-w-3xl flex-col gap-10 py-10">
			<div class="flex flex-col gap-4">
				<h1 class="text-3xl font-semibold tracking-tight">Getting Started</h1>
				<p class="text-muted-foreground">
					gsxui components are copy-in: the CLI vendors real <code>.gsx</code> source into your own
					module, so what you build against is code you own and can edit — not a package you import
					and can't touch.
				</p>
			</div>
			<section class="flex flex-col gap-3">
				<h2>1. Install the CLI</h2>
				<pre><code>{ hl.Node("snippets/install.sh") }</code></pre>
			</section>
			<section class="flex flex-col gap-3">
				<h2>2. Initialize your project</h2>
				<p>In your project (a Go module):</p>
				<pre><code>{ hl.Node("snippets/init.sh") }</code></pre>
				<pre><code>{ hl.Node("snippets/init.output") }</code></pre>
				<p>
					This vendors the theme tokens (<code>web/gsxui.css</code>), the JS runtime and behavior
					barrel (<code>web/gsxui/</code>), and the class merger (<code>ui/merge/merge.go</code>), then
					points <code>gsx.toml</code>'s <code>class_merger</code> at it — the seam that makes
					caller-class-merge work (see <a href={Theming{} |> url}>Theming</a>).
					It also <code>go get</code> <code>gsx</code> and <code>tailwind-merge-go</code>, and installs
					the <code>gsx</code> tool via <code>go get -tool</code>.
				</p>
				<p>
					<code>web/gsxui.css</code> begins with <code>@import "tailwindcss"</code> and <code>@import
					"tw-animate-css"</code> — your Tailwind build resolves both from npm, so make sure they're
					installed: <code>npm install tailwindcss @tailwindcss/vite tw-animate-css</code>.
					Without <code>tw-animate-css</code> every <code>animate-in</code>/<code>animate-out</code> class
					the components carry (dialog, dropdown, tooltip) is silently inert.
				</p>
			</section>
			<section class="flex flex-col gap-3">
				<h2>3. Add components</h2>
				<pre><code>{ hl.Node("snippets/add.sh") }</code></pre>
				<pre><code>{ hl.Node("snippets/add.output") }</code></pre>
				<p>
					<code>card</code> has no dependencies of its own, but a component that does
					(e.g. <code>native-select</code>, which needs <code>icon</code>) pulls its dependency in
					automatically — <code>gsxui add native-select</code> vendors <code>icon</code> too. You own every
					file this writes: <code>gsxui add</code> never touches one you've already modified unless
					you pass <code>--overwrite</code>. After upgrading the <code>gsxui</code> binary,
					re-run <code>gsxui add &lt;name&gt; --overwrite</code> to refresh vendored components —
					that discards local edits to those files.
				</p>
			</section>
			<section class="flex flex-col gap-3">
				<h2>4. Your first page</h2>
				<p>
					A tiny two-file app: <code>home.gsx</code> renders a <code>Card</code> around
					a <code>Button</code>, and <code>main.go</code> serves it.
				</p>
				<pre><code>{ hl.Node("snippets/first-page.gsx") }</code></pre>
				<pre><code>{ hl.Node("snippets/first-main.go") }</code></pre>
				<p>
					Compile the <code>.gsx</code> file to plain Go, then run it:
				</p>
				<pre><code>{ hl.Node("snippets/generate.sh") }</code></pre>
				<p>
					(silent on success — it writes <code>home.x.go</code> next
					to <code>home.gsx</code> and exits 0)
				</p>
				<pre><code>{ hl.Node("snippets/run.sh") }</code></pre>
				<pre><code>{ hl.Node("snippets/run.output") }</code></pre>
				<p>
					Open <code>http://localhost:8080</code> — a styled Card with a Button inside, rendered with
					gsxui's default light theme.
					Next: <a href={Theming{} |> url}>restyle it</a>.
				</p>
			</section>
		</div>
	</Layout>
}
