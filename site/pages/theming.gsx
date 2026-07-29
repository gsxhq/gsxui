package pages

import "github.com/gsxhq/gsxui/site/hl"

// Theming documents the four-file CSS boundary installed by gsxui init.
type Theming struct{}

var themingTOCItems = []docTOCItem{
	{ID: "css-files", Title: "The four CSS files", Depth: 2},
	{ID: "semantic-variables", Title: "Edit semantic variables", Depth: 2},
	{ID: "component-markers", Title: "Stable component markers", Depth: 2},
	{ID: "caller-utilities", Title: "Caller utilities win", Depth: 2},
	{ID: "breaking-migration", Title: "Breaking migration", Depth: 2},
}

component (t Theming) Page() {
	<siteLayout title="Theming" active="theming" mode={layoutDocs} toc={themingTOCItems}>
		<div data-doc="theming" class="flex max-w-3xl flex-col gap-10 py-10">
			<div class="flex flex-col gap-4">
				<h1 class="text-3xl font-semibold tracking-tight">Theming</h1>
				<p class="text-muted-foreground">
					gsxui keeps behavior, semantic values, and component presentation in separate CSS files. A theme changes
					variables only; every component keeps the same canonical GSX markup.
				</p>
			</div>
			<section class="flex flex-col gap-3">
				<docHeading item={themingTOCItems[0]}/>
				<pre><code>{ hl.Node("snippets/theme-entry.css") }</code></pre>
				<ul class="list-disc space-y-2 pl-6">
					<li><code>index.css</code> is the one entry your app imports.</li>
					<li>
						<code>foundation.css</code> owns accessibility and behavior-critical mechanics such as hidden states,
						positioning, and interaction geometry.
					</li>
					<li>
						<code>theme.css</code> owns semantic light/dark variables, including sidebar, status, overlay, contrast, and
						radius tokens.
					</li>
					<li>
						<code>style.css</code> owns the replaceable component presentation: density, borders, shadows, typography,
						and visual variants.
					</li>
				</ul>
				<p>
					To recolor the default style, replace only <code>theme.css</code>. Keep the other three files unchanged.
				</p>
			</section>
			<section class="flex flex-col gap-3">
				<docHeading item={themingTOCItems[1]}/>
				<p>
					The variables are shadcn-compatible <code>:root</code> and <code>.dark</code> blocks. The theme editor imports
					and exports exactly this variables-only file.
				</p>
				<pre><code>{ hl.Node("snippets/theme-restyle.css") }</code></pre>
			</section>
			<section class="flex flex-col gap-3">
				<docHeading item={themingTOCItems[2]}/>
				<p>
					Components expose one bare <code>data-gsxui-slot-&lt;name&gt;</code> attribute per semantic part. Composed
					parts forward each distinct attribute onto the same element, so selectors use exact presence matching:
				</p>
				<pre><code>{ hl.Node("snippets/theme-slot.css") }</code></pre>
				<p>
					Use the same exact presence form in project CSS. Value and token operators are not part of the contract.
				</p>
			</section>
			<section class="flex flex-col gap-3">
				<docHeading item={themingTOCItems[3]}/>
				<p>
					The default style is in <code>@layer components</code>. Tailwind utilities are emitted later, so an ordinary
					caller class such as <code>h-12</code> or <code>rounded-full</code> overrides the style
					without <code>!important</code> or a class merger fighting component-owned utility strings.
				</p>
				<pre><code>{ hl.Node("snippets/theme-merge.gsx") }</code></pre>
				<p>
					Fallthrough attributes still carry ids, ARIA, data, and HTMX attributes to the rendered element:
				</p>
				<pre><code>{ hl.Node("snippets/theme-attrs.gsx") }</code></pre>
			</section>
			<section class="flex flex-col gap-3">
				<docHeading item={themingTOCItems[4]}/>
				<ol class="list-decimal space-y-2 pl-6">
					<li>Change the CSS entry from <code>web/gsxui.css</code> to <code>web/gsxui/index.css</code>.</li>
					<li>Review the four-file diff, then run <code>gsxui init --overwrite</code>.</li>
					<li>
						Run <code>gsxui add &lt;component&gt; --overwrite</code> for each vendored component you want to refresh.
					</li>
					<li>
						Replace intentional project <code>data-slot</code> selectors with exact presence selectors such
						as <code>[data-gsxui-slot-button]</code>.
					</li>
				</ol>
				<p>
					This is a one-time breaking migration. There is no legacy selector or combined-file compatibility layer.
				</p>
			</section>
		</div>
	</siteLayout>
}
