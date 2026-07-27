package pages

import "github.com/gsxhq/gsxui/site/hl"

// Theming documents the four-file CSS boundary installed by gsxui init.
type Theming struct{}

component (t Theming) Page() {
	<Layout title="Theming" active="theming">
		<div data-doc="theming" class="flex max-w-3xl flex-col gap-10 py-10">
			<div class="flex flex-col gap-4">
				<h1 class="text-3xl font-semibold tracking-tight">Theming</h1>
				<p class="text-muted-foreground">
					gsxui keeps behavior, semantic values, and component presentation in separate CSS files. A theme changes variables only; every component keeps the same canonical GSX markup.
				</p>
			</div>

			<section class="flex flex-col gap-3">
				<h2>The four CSS files</h2>
				<pre><code>{ hl.Node("snippets/theme-entry.css") }</code></pre>
				<ul class="list-disc space-y-2 pl-6">
					<li><code>index.css</code> is the one entry your app imports.</li>
					<li><code>foundation.css</code> owns accessibility and behavior-critical mechanics such as hidden states, positioning, and interaction geometry.</li>
					<li><code>theme.css</code> owns semantic light/dark variables, including sidebar, status, overlay, contrast, and radius tokens.</li>
					<li><code>style.css</code> owns the replaceable component presentation: density, borders, shadows, typography, and visual variants.</li>
				</ul>
				<p>
					To recolor the default style, replace only <code>theme.css</code>. Keep the other three files unchanged.
				</p>
			</section>

			<section class="flex flex-col gap-3">
				<h2>Edit semantic variables</h2>
				<p>
					The variables are shadcn-compatible <code>:root</code> and <code>.dark</code> blocks. The theme editor imports and exports exactly this variables-only file.
				</p>
				<pre><code>{ hl.Node("snippets/theme-restyle.css") }</code></pre>
			</section>

			<section class="flex flex-col gap-3">
				<h2>Stable component tokens</h2>
				<p>
					Components expose space-separated part tokens through <code>data-gsxui-slot</code>. Because an element can compose several parts, selectors use token membership:
				</p>
				<pre><code>{ hl.Node("snippets/theme-slot.css") }</code></pre>
				<p>
					Use the same membership form in project CSS. Equality selectors are incorrect for composed elements.
				</p>
			</section>

			<section class="flex flex-col gap-3">
				<h2>Caller utilities win</h2>
				<p>
					The default style is in <code>@layer components</code>. Tailwind utilities are emitted later, so an ordinary caller class such as <code>h-12</code> or <code>rounded-full</code> overrides the style without <code>!important</code> or a class merger fighting component-owned utility strings.
				</p>
				<pre><code>{ hl.Node("snippets/theme-merge.gsx") }</code></pre>
				<p>
					Fallthrough attributes still carry ids, ARIA, data, and HTMX attributes to the rendered element:
				</p>
				<pre><code>{ hl.Node("snippets/theme-attrs.gsx") }</code></pre>
			</section>

			<section class="flex flex-col gap-3">
				<h2>Breaking migration</h2>
				<ol class="list-decimal space-y-2 pl-6">
					<li>Change the CSS entry from <code>web/gsxui.css</code> to <code>web/gsxui/index.css</code>.</li>
					<li>Review the four-file diff, then run <code>gsxui init --overwrite</code>.</li>
					<li>Run <code>gsxui add &lt;component&gt; --overwrite</code> for each vendored component you want to refresh.</li>
					<li>Replace intentional project <code>data-slot</code> selectors with <code>[data-gsxui-slot~="&lt;token&gt;"]</code>.</li>
				</ol>
				<p>
					This is a one-time breaking migration. There is no legacy selector or combined-file compatibility layer.
				</p>
			</section>
		</div>
	</Layout>
}
