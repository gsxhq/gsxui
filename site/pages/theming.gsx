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
					A theme changes CSS variables only: recolor every component, light and dark, by editing one file. Component
					markup never changes.
				</p>
			</div>
			<section class="flex flex-col gap-3">
				<docHeading item={themingTOCItems[0]}/>
				<p><code>gsxui init</code> vendors four CSS files into <code>web/gsxui/</code>:</p>
				<pre><code>{ hl.Node("snippets/theme-entry.css") }</code></pre>
				<ul class="list-disc space-y-2 pl-6">
					<li><code>index.css</code> is the one entry your app imports.</li>
					<li>
						<code>foundation.css</code> owns accessibility and behavior-critical mechanics such as hidden states,
						positioning, and interaction geometry.
					</li>
					<li>
						<code>theme.css</code> owns semantic light/dark variables, including sidebar, status, overlay, contrast,
						chart, and radius tokens, plus the mode-independent <code>--font-sans</code>/<code>--font-heading</code>
						typography pair.
					</li>
					<li>
						<code>style.css</code> owns the default style's component presentation rules: density, borders, shadows,
						typography.
					</li>
				</ul>
				<p>
					To recolor the default style, edit or replace <code>theme.css</code> only. The other three files stay
					untouched.
				</p>
			</section>
			<section class="flex flex-col gap-3">
				<docHeading item={themingTOCItems[1]}/>
				<p>
					The variables are shadcn-compatible <code>:root</code> and <code>.dark</code> blocks — nothing else lives in
					the file. The <a href={Theme{} |> url}>theme editor</a> imports and exports exactly this file, so you can
					build a theme visually and download the result.
				</p>
				<pre><code>{ hl.Node("snippets/theme-restyle.css") }</code></pre>
			</section>
			<section class="flex flex-col gap-3">
				<docHeading item={themingTOCItems[2]}/>
				<p>
					To restyle a part from your own CSS, target its marker: every semantic part carries one bare
					<code>data-gsxui-slot-&lt;name&gt;</code> attribute. A composed part carries every role it plays —
					AlertDialog's action renders a Button whose element has both <code>data-gsxui-slot-button</code>
					and <code>data-gsxui-slot-alert-dialog-action</code> — so each role stays targetable:
				</p>
				<pre><code>{ hl.Node("snippets/theme-slot.css") }</code></pre>
				<p>
					Match on bare presence only — value and operator forms are not part of the contract.
				</p>
			</section>
			<section class="flex flex-col gap-3">
				<docHeading item={themingTOCItems[3]}/>
				<p>
					Pass an ordinary utility class and it wins — no <code>!important</code>:
				</p>
				<pre><code>{ hl.Node("snippets/theme-merge.gsx") }</code></pre>
				<p>Two mechanisms make that hold:</p>
				<ul class="list-disc space-y-2 pl-6">
					<li>
						Utilities a component carries in its own markup are settled by the vendored class merger
						(<code>ui/merge/merge.go</code>): your <code>h-12</code> replaces the component's own height utility
						instead of conflicting with it.
					</li>
					<li>
						Rules in <code>style.css</code> live in <code>@layer components</code>, while your class lands in
						Tailwind's later <code>@layer utilities</code> — the later layer wins regardless of specificity.
					</li>
				</ul>
				<p>
					Fallthrough attributes carry ids, ARIA, data, and HTMX attributes to the rendered element the same way:
				</p>
				<pre><code>{ hl.Node("snippets/theme-attrs.gsx") }</code></pre>
				<p>
					Behaviour roles are opt-in data attributes: any element becomes a dialog trigger by
					carrying <code>data-gsxui-dialog-trigger</code>. The family's own Trigger components need no role attribute
					— their slot marker implies it:
				</p>
				<pre><code>{ hl.Node("snippets/theme-dataattr.gsx") }</code></pre>
			</section>
			<section class="flex flex-col gap-3">
				<docHeading item={themingTOCItems[4]}/>
				<p>
					Only projects from before the four-file split — a single <code>web/gsxui.css</code> entry
					and <code>data-slot</code> markers — need this. There is no compatibility layer; migrate in one pass:
				</p>
				<ol class="list-decimal space-y-2 pl-6">
					<li>Change the CSS entry from <code>web/gsxui.css</code> to <code>web/gsxui/index.css</code>.</li>
					<li>Review the four-file diff, then run <code>gsxui init --overwrite</code>.</li>
					<li>
						Run <code>gsxui add &lt;component&gt; --overwrite</code> for each vendored component you want to refresh.
					</li>
					<li>
						Replace project <code>data-slot</code> selectors with exact presence selectors such
						as <code>[data-gsxui-slot-button]</code>.
					</li>
				</ol>
			</section>
		</div>
	</siteLayout>
}
