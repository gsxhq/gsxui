package pages

import (
	"github.com/gsxhq/gsxui/site/examples"
	"github.com/gsxhq/gsxui/site/hl"
)

// Rtl is the /docs/rtl page: how gsxui components adapt to dir="rtl",
// followed by a live Arabic login demo.
type Rtl struct{}

var rtlTOCItems = []docTOCItem{
	{ID: "get-started", Title: "Get started", Depth: 2},
	{ID: "how-it-works", Title: "How it works", Depth: 2},
	{ID: "try-it-out", Title: "Try it out", Depth: 2},
	{ID: "fonts", Title: "Fonts", Depth: 2},
}

component (rt Rtl) Page() {
	{{ loginExample := examples.For("rtl")[0] }}
	<siteLayout title="RTL" active="rtl" mode={layoutDocs} toc={rtlTOCItems}>
		<div data-doc="rtl" class="flex max-w-3xl flex-col gap-10 py-10">
			<div class="flex flex-col gap-4">
				<h1 class="text-3xl font-semibold tracking-tight">RTL</h1>
				<p class="text-muted-foreground">
					Opt in once and the components you vendor use logical direction classes, so one build serves{ " " }
					<code>dir="ltr"</code> and <code>dir="rtl"</code> alike.
				</p>
			</div>
			<section class="flex flex-col gap-3">
				<docHeading item={rtlTOCItems[0]}/>
				<p>Set the flag when you initialise, then add components as usual:</p>
				<pre><code>{ hl.Node("snippets/rtl-init.sh") }</code></pre>
				<p>Components vendored before the flag was set are rewritten in place, keeping your edits:</p>
				<pre><code>{ hl.Node("snippets/rtl-migrate.sh") }</code></pre>
				<p>
					Set <code>dir="rtl"</code> on <code>&lt;html&gt;</code>, or on any subtree, from the locale. An{ " " }
					<code>rtl: true</code> project renders both directions; the document's <code>dir</code> decides.
				</p>
				<p>The flag is one-way: there is no reverse migration, so opt in before you customise vendored files.</p>
				<pre><code>{ `<html lang="ar" dir="rtl">
  ...
</html>` }</code></pre>
				<p>
					Translating the strings components write themselves is covered
					on <a href={I18n{} |> url}>Internationalization</a>.
				</p>
			</section>
			<section class="flex flex-col gap-3">
				<docHeading item={rtlTOCItems[1]}/>
				<ul class="list-disc space-y-2 pl-6">
					<li>
						<strong>Logical classes, at install time.</strong> With <code>"rtl": true</code> in{ " " }
						<code>gsxui.json</code>, <code>gsxui add</code> rewrites physical classes in what it vendors:{ " " }
						<code>ml-*</code> to <code>ms-*</code>, <code>left-*</code> to <code>start-*</code>,{ " " }
						<code>text-left</code> to <code>text-start</code>, and the rest of shadcn's table.{ " " }
						<code>translate-x-*</code> and <code>space-x-*</code> gain <code>rtl:</code> companions.
					</li>
					<li>
						<strong>Directional icons.</strong> Chevrons and arrows that encode a left/right meaning carry{ " " }
						<code>rtl:rotate-180</code>.
					</li>
					<li>
						<strong>Direction-aware floating positioning.</strong> Popover, dropdown-menu, select, tooltip and the rest
						of the floating family resolve placement in JS at position time.
					</li>
					<li>
						<strong>Mirrored keyboard semantics.</strong> Arrow keys in menus, tabs, carousel, calendar and the other
						roving-focus components mirror by meaning per WAI-ARIA.
					</li>
					<li>
						<strong>The demos on this page</strong> render from the transformed components, the same output{ " " }
						<code>gsxui add</code> produces with the flag set.
					</li>
				</ul>
				<p>Not transformed:</p>
				<ul class="list-disc space-y-2 pl-6">
					<li>
						Sheet and Sidebar <code>side="left"</code>/<code>side="right"</code> and Drawer{ " " }
						<code>direction="left"</code>/<code>direction="right"</code> stay physical, matching shadcn's{ " " }
						<code>data-side</code> contract; their interiors mirror.
					</li>
					<li><code>input-otp</code>'s digit group stays pinned <code>dir="ltr"</code>.</li>
					<li>Behaviour JS carries no classes and is untouched.</li>
					<li>
						Sized logical slide utilities (<code>slide-in-from-start-2</code>) match only elements that themselves
						carry <code>dir</code> in the current tw-animate-css build; gsxui emits none.
					</li>
				</ul>
			</section>
			<section class="flex flex-col gap-3">
				<docHeading item={rtlTOCItems[2]}/>
				<p>
					An Arabic sign-in card: <code>Card</code>, <code>Label</code>, <code>Input</code>,{ " " }
					<code>NativeSelect</code> and <code>Button</code>, rendered from the transformed components, wrapped
					in <code>dir="rtl"</code>.
				</p>
				<div class="border rounded-lg p-8 bg-background">
					{ loginExample.Node }
				</div>
				<div class="relative" data-site-example>
					<pre
						class="overflow-x-auto rounded-2xl bg-muted/50 px-4 py-3.5 font-mono text-sm"
					><code>{ hl.Node(loginExample.SourcePath) }</code></pre>
					<button
						type="button"
						data-site-copy
						class="absolute right-2 top-2 rounded-md border border-border bg-background px-2 py-1 text-xs text-muted-foreground transition-colors hover:bg-accent hover:text-accent-foreground"
					>
						Copy
					</button>
				</div>
				<p class="text-sm text-muted-foreground">
					The source imports <code>uirtl</code>, this site's transformed copy of <code>ui</code>. In your project these
					are <code>ui.Card</code>, <code>ui.NativeSelect</code> and so on.
				</p>
				<p>
					More RTL variants live on their own component pages: <a href="/components/calendar">Calendar</a>
					, <a href="/components/pagination">Pagination</a>, and <a href="/components/sidebar">Sidebar</a> each register
					an "RTL" example alongside their other demos.
				</p>
			</section>
			<section class="flex flex-col gap-3">
				<docHeading item={rtlTOCItems[3]}/>
				<p>
					The site's own Latin type is Geist. For Arabic content, pair it
					with <a href="https://fonts.google.com/noto/specimen/Noto+Sans+Arabic" target="_blank" rel="noreferrer">Noto Sans Arabic</a> (a
					UI sans)
					or <a href="https://fonts.google.com/noto/specimen/Noto+Naskh+Arabic" target="_blank" rel="noreferrer">Noto Naskh Arabic</a> (better
					for longer prose). Load the Arabic font as a fallback rather than replacing Geist outright, so Latin text —
					brand names, code — keeps rendering in Geist inside an RTL document.
				</p>
			</section>
		</div>
	</siteLayout>
}
