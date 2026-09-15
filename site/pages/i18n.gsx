package pages

import (
	"github.com/gsxhq/gsxui/site/examples"
	"github.com/gsxhq/gsxui/site/hl"
)

// I18n is the /docs/i18n page: translating the strings components write
// themselves, and the Calendar locale.
type I18n struct{}

var i18nTOCItems = []docTOCItem{
	{ID: "get-started", Title: "Get started", Depth: 2},
	{ID: "how-it-works", Title: "How it works", Depth: 2},
	{ID: "calendar", Title: "Calendar", Depth: 2},
	{ID: "leftover-english", Title: "Finding leftover English", Depth: 2},
	{ID: "not-translated", Title: "Not translated", Depth: 2},
}

func localizedCalendarExample() examples.Example {
	for _, ex := range examples.For("calendar") {
		if ex.Name == "localized" {
			return ex
		}
	}
	panic("calendar example \"localized\" is not registered")
}

component (p I18n) Page() {
	{{ localized := localizedCalendarExample() }}
	<siteLayout title="Internationalization" active="i18n" mode={layoutDocs} toc={i18nTOCItems}>
		<div data-doc="i18n" class="flex max-w-3xl flex-col gap-10 py-10">
			<div class="flex flex-col gap-4">
				<h1 class="text-3xl font-semibold tracking-tight">Internationalization</h1>
				<p class="text-muted-foreground">
					Text you pass to a component is yours. The few strings a component writes itself are
					<code>i18n.T</code> messages, and one line in <code>gsx.toml</code> routes all of them through your
					translator.
				</p>
			</div>
			<section class="flex flex-col gap-3">
				<docHeading item={i18nTOCItems[0]}/>
				<p>
					Add <code>ui/i18n/translate.go</code> beside the vendored <code>ui/i18n/i18n.go</code> and register it, using
					your module path:
				</p>
				<pre><code>{ hl.Node("snippets/i18n-renderer.toml") }</code></pre>
				<p>Write the translator. It receives the render context, so the request's locale is in reach:</p>
				<pre><code>{ hl.Node("snippets/i18n-translate.go") }</code></pre>
			</section>
			<section class="flex flex-col gap-3">
				<docHeading item={i18nTOCItems[1]}/>
				<ul class="list-disc space-y-2 pl-6">
					<li>
						<code>i18n.T</code> is a string type in <code>ui/i18n/i18n.go</code>, its own package, vendored with any
						component that uses it.
					</li>
					<li>
						The English text is the message id: <code>i18n.T("Close")</code>, <code>i18n.T("Next slide")</code>
						, <code>i18n.T("Toggle Sidebar")</code>.
					</li>
					<li>Without a renderer, a message renders as its English text.</li>
					<li>With one, gsx calls your function everywhere a message renders, in text and in attributes.</li>
					<li>Children, props and <code>attrs</code> never pass through <code>T</code>.</li>
					<li>
						<code>gsxui add --overwrite</code> rewrites only the files it vendored,
						so <code>translate.go</code> survives.
					</li>
				</ul>
			</section>
			<section class="flex flex-col gap-3">
				<docHeading item={i18nTOCItems[2]}/>
				<p>
					Calendar composes month and weekday names, the caption and each day's label. Pass a
					<code>locale</code>; the zero value is English.
				</p>
				<ul class="list-disc space-y-2 pl-6">
					<li>
						<code>Caption</code> and <code>DayLabel</code> substitute <code>{ "{month}" }</code>
						, <code>{ "{year}" }</code>, <code>{ "{weekday}" }</code> and <code>{ "{day}" }</code>. Other text is
						literal.
					</li>
					<li>
						<code>Digits</code> is ten runes replacing 0-9 in visible text. Form values and data attributes stay ASCII.
					</li>
					<li>Fill every field. An empty field renders empty; only the all-zero value means English.</li>
					<li>The client reads the same values from the root, so navigation writes the same text the server did.</li>
				</ul>
				<div class="border rounded-lg p-8 bg-background">
					{ localized.Node }
				</div>
				<div class="relative" data-site-example>
					<pre
						class="overflow-x-auto rounded-2xl bg-muted/50 px-4 py-3.5 font-mono text-sm"
					><code>{ hl.Node(localized.SourcePath) }</code></pre>
					<button
						type="button"
						data-site-copy
						class="absolute right-2 top-2 rounded-md border border-border bg-background px-2 py-1 text-xs text-muted-foreground transition-colors hover:bg-accent hover:text-accent-foreground"
					>
						Copy
					</button>
				</div>
			</section>
			<section class="flex flex-col gap-3">
				<docHeading item={i18nTOCItems[3]}/>
				<p>
					Register a translator that wraps every message in markers, render each page, and search the output for English
					outside the markers. Anything found is text a caller passed in, one of the overridable defaults below, or a
					gsxui bug worth an issue.
				</p>
			</section>
			<section class="flex flex-col gap-3">
				<docHeading item={i18nTOCItems[4]}/>
				<ul class="list-disc space-y-2 pl-6">
					<li>
						Defaults a caller can already override by passing the same attribute stay
						English: <code>aria-label</code>{ " " }on Spinner, Breadcrumb, Pagination and its Previous/Next links,
						SidebarRail's <code>aria-label</code> and{ " " }
						<code>title</code>, and Carousel's <code>aria-roledescription</code>. Pass a translated value in your
						markup.
					</li>
					<li>
						When no <code>ui.Toaster</code> is mounted, <code>toaster.js</code> creates its own region with an
						English <code>aria-label</code>. Mount <code>ui.Toaster</code> and the landmark is a message like the rest.
					</li>
				</ul>
			</section>
		</div>
	</siteLayout>
}
