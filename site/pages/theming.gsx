package pages

// Theming is the /docs/theming page: the token model (20 shadcn-compatible
// CSS custom properties, light + .dark), how to restyle a project's vendored
// gsxui.css, and the three-part customization story from
// docs/jsx-parity.md's dialog MECHANISM entry and the button/badge WIN
// entries — caller class merge, attrs fallthrough, and the data-attribute
// idiom for behavior attachment.
type Theming struct{}

const themeTokensCSS = `@theme inline {
  --radius-sm: calc(var(--radius) - 4px);
  --radius-md: calc(var(--radius) - 2px);
  --radius-lg: var(--radius);
  --radius-xl: calc(var(--radius) + 4px);
  --color-background: var(--background);
  --color-foreground: var(--foreground);
  --color-card: var(--card);
  --color-card-foreground: var(--card-foreground);
  --color-popover: var(--popover);
  --color-popover-foreground: var(--popover-foreground);
  --color-primary: var(--primary);
  --color-primary-foreground: var(--primary-foreground);
  --color-secondary: var(--secondary);
  --color-secondary-foreground: var(--secondary-foreground);
  --color-muted: var(--muted);
  --color-muted-foreground: var(--muted-foreground);
  --color-accent: var(--accent);
  --color-accent-foreground: var(--accent-foreground);
  --color-destructive: var(--destructive);
  --color-destructive-foreground: var(--destructive-foreground);
  --color-border: var(--border);
  --color-input: var(--input);
  --color-ring: var(--ring);
}

:root {
  --radius: 0.625rem;
  --background: oklch(1 0 0);
  --foreground: oklch(0% 0 0);
  --card: oklch(1 0 0);
  --card-foreground: oklch(0% 0 0);
  --popover: oklch(1 0 0);
  --popover-foreground: oklch(0% 0 0);
  --primary: oklch(0% 0 0);
  --primary-foreground: oklch(0.985 0 0);
  --secondary: oklch(0.97 0 0);
  --secondary-foreground: oklch(0.205 0 0);
  --muted: oklch(0.97 0 0);
  --muted-foreground: oklch(0.556 0 0);
  --accent: oklch(0.97 0 0);
  --accent-foreground: oklch(0.205 0 0);
  --destructive: oklch(0.577 0.245 27.325);
  --destructive-foreground: oklch(0.97 0.01 17);
  --border: oklch(0.922 0 0);
  --input: oklch(0.922 0 0);
  --ring: oklch(0.708 0 0);
}

.dark {
  --background: oklch(0.145 0 0);
  --foreground: oklch(0.985 0 0);
  --primary: oklch(0.922 0 0);
  --primary-foreground: oklch(0.205 0 0);
  --border: oklch(1 0 0 / 10%);
  --ring: oklch(0.556 0 0);
  /* …and the rest of the pairs above, each with a dark-mode value */
}`

const themeRestyleSnippet = `/* web/gsxui.css — yours after 'gsxui init'; edit the values in place */
:root {
  --radius: 0.5rem;                  /* squarer corners, every component */
  --primary: oklch(0.55 0.2 265);    /* brand blue instead of near-black */
  --primary-foreground: oklch(0.98 0 0);
}

.dark {
  --primary: oklch(0.72 0.15 265);
}`

const themeMergeSnippet = `<ui.Button class="h-12">Tall</ui.Button>`

const themeAttrsSnippet = `<ui.Button
	id="submit"
	aria-label="Submit the form"
	data-testid="submit-btn"
	hx-post="/submit"
	hx-target="#result"
>
	Submit
</ui.Button>`

const themeDataAttrSnippet = `<ui.Dialog>
	<ui.Button variant="outline" data-gsxui-dialog-trigger>
		Open
	</ui.Button>
	<ui.DialogContent>
		...
	</ui.DialogContent>
</ui.Dialog>`

component (t Theming) Page() {
	<Layout title="Theming" active="theming">
		<div data-doc="theming" class="flex max-w-3xl flex-col gap-10 py-10">
			<div class="flex flex-col gap-4">
				<h1 class="text-3xl font-semibold tracking-tight">Theming</h1>
				<p class="text-muted-foreground">
					Every gsxui component is styled entirely through Tailwind utilities that resolve to a fixed set of CSS custom properties — never a literal color. Change the properties and every component restyles at once; no per-component theme prop, no rebuild of gsxui itself.
				</p>
			</div>
			<section class="flex flex-col gap-3">
				<h2>The token model</h2>
				<p>
					<code>gsxui init</code> vendors <code>web/gsxui.css</code> with 19 color tokens plus <code>
						--radius
					</code> (20 total) defined twice — once in <code>:root</code> for light mode, once in <code>.dark</code>
					for dark mode (toggled by a <code>.dark</code> class anywhere up the tree, via Tailwind's <code>
						@custom-variant dark
					</code>
					) — and an
					<code>@theme inline</code> block that maps each one onto Tailwind's own color scale, so <code>
						bg-primary
					</code>
					, <code>text-muted-foreground</code>
					,
					<code>border-input</code>
					, and friends all resolve to a token, not a hard-coded value:
				</p>
				<pre><code>{ themeTokensCSS }</code></pre>
				<p>
					The eight paired tokens (
					<code>background</code>
					/
					<code>foreground</code>
					,
					<code>card</code>
					, <code>popover</code>
					, <code>primary</code>
					,
					<code>secondary</code>
					, <code>muted</code>
					, <code>accent</code>
					, each with a matching <code>-foreground</code>
					, plus <code>destructive</code>
					/
					<code>destructive-foreground</code>
					) cover every surface + text combination a component draws; <code>border</code>
					, <code>input</code>
					, and <code>ring</code> cover outlines and focus rings; <code>radius</code>
					drives every rounded corner via the derived <code>--radius-sm</code>
					…
					<code>--radius-xl</code> scale.
				</p>
			</section>
			<section class="flex flex-col gap-3">
				<h2>How to restyle</h2>
				<p>
					<code>web/gsxui.css</code> is vendored, not imported — it's yours the moment <code>
						gsxui init
					</code> writes it. Restyling is editing the values inside <code>:root</code> and <code>.dark</code> directly:
				</p>
				<pre><code>{ themeRestyleSnippet }</code></pre>
				<p>
					Because the variable names (
					<code>--primary</code>
					,
					<code>--primary-foreground</code>
					, …) match shadcn/ui's own convention exactly, the file is <strong>tweakcn-compatible</strong>
					: generate a theme at
					<a href="https://tweakcn.com" target="_blank" rel="noreferrer">tweakcn.com</a>
					(or any other shadcn theme tool) and paste its <code>:root</code>
					/
					<code>.dark</code> blocks over gsxui's own — no renaming, no translation layer.
				</p>
			</section>
			<section class="flex flex-col gap-6">
				<h2>Customizing components</h2>
				<div class="flex flex-col gap-3">
					<h3>Caller class merge: a conflicting utility wins</h3>
					<p>
						Every component's fallthrough <code>attrs</code> can carry a
						<code>class</code>
						, and it doesn't just get appended — <code>gsx.toml</code>
						's
						<code>class_merger</code> (vendored to <code>ui/merge/merge.go</code> by
						<code>gsxui init</code>
						, backed by <code>tailwind-merge-go</code>
						) resolves conflicts the way Tailwind itself would: whichever utility comes last in the same category wins, structural classes that aren't in that category are untouched.
					</p>
					<pre><code>{ themeMergeSnippet }</code></pre>
					<p>
						<code>Button</code>
						's default size class is
						<code>h-9 px-4 py-2 has-[&gt;svg]:px-3</code>
						. The caller's
						<code>h-12</code> is in the same height category as <code>h-9</code>
						, so it drops <code>h-9</code> and wins; <code>px-4 py-2</code> and the structural base classes (
						<code>inline-flex</code>
						,
						<code>items-center</code>
						, <code>rounded-md</code>
						, …) survive because nothing the caller passed conflicts with them.
					</p>
				</div>
				<div class="flex flex-col gap-3">
					<h3>Attrs fallthrough: id, aria-*, data-*, hx-*</h3>
					<p>
						Beyond <code>class</code>
						, every attribute a caller passes that isn't one of the component's own named parameters lands on the rendered element untouched — ids, ARIA attributes, arbitrary
						<code>data-*</code>
						, and HTMX's <code>hx-*</code> attributes all pass straight through:
					</p>
					<pre><code>{ themeAttrsSnippet }</code></pre>
				</div>
				<div class="flex flex-col gap-3">
					<h3>Data-attribute idiom: attaching behavior to your own markup</h3>
					<p>
						Interactive components (dialog, dropdown, tabs, tooltip, …) don't use React's <code>asChild</code>
						/Slot pattern — gsx has no dynamic tag-swapping. Instead, each interactive component's
						<code>data-gsxui-*</code> attribute is its public contract, and fallthrough <code>
							attrs
						</code> deliver it to <em>any</em> element or component, no cloning and no wrapper required. A plain styled
						<code>Button</code> becomes a dialog trigger just by carrying the attribute:
					</p>
					<pre><code>{ themeDataAttrSnippet }</code></pre>
					<p>
						The same idiom covers every interactive component's public hooks —
						<code>data-gsxui-dialog-close</code>
						,
						<code>data-gsxui-dropdown-trigger</code>
						, and so on — see each component's page for its specific attribute names.
					</p>
				</div>
			</section>
		</div>
	</Layout>
}
