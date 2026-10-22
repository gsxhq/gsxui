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
	{ID: "calendar-from-cldr", Title: "Filling the locale from CLDR", Depth: 3},
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
					<li>
						Renderers bind when <code>gsx generate</code> runs in your module, so only vendored components translate.
						Importing <code>github.com/gsxhq/gsxui/ui</code> directly renders English.
					</li>
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
					<li>
						<code>MonthsFormat</code> is the month inside a date, for day labels: Russian <code>29 июня</code>,
						caption{ " " }<code>июнь</code>. Leave it unset where both forms match.
					</li>
					<li>Fill every other field. An empty field renders empty; only the all-zero value means English.</li>
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
					Every field is CLDR data, so the server can fill it from a CLDR source instead of a typed table. With{ " " }
					<code>Intl</code> on the server (a Cloudflare Worker has V8's), each field is one call on a date in the
					Gregorian calendar and UTC:
				</p>
				<ul class="list-disc space-y-2 pl-6">
					<li>
						<code>Months</code>: the <code>month</code> part of <code>formatToParts</code> for{ " " }
						<code>{ `{month: "long", year: "numeric"}` }</code>, the caption's own form: Russian{ " " }<code>июнь</code>
						, Japanese <code>6</code> (its <code>月</code> is in <code>Caption</code>).
					</li>
					<li>
						<code>MonthsFormat</code>: the <code>month</code> part for the day label's{ " " }
						<code>{ `{weekday: "long", month: "long", day: "numeric", year: "numeric"}` }</code>: Russian{ " " }
						<code>июня</code>.
					</li>
					<li><code>MonthsShort</code>: <code>{ `{month: "short"}` }</code>.</li>
					<li>
						<code>Weekdays</code> and <code>WeekdaysShort</code>: <code>{ `{weekday: "long"}` }</code> and{ " " }
						<code>{ `{weekday: "narrow"}` }</code>.
					</li>
					<li>
						<code>Caption</code> and <code>DayLabel</code>: the same two formats' parts, with each <code>month</code>
						,{ " " }
						<code>year</code>, <code>weekday</code> and <code>day</code> part replaced by its placeholder. Russian's
						caption is <code>{ "{month} {year} г." }</code>
					</li>
					<li>
						<code>Digits</code>
						: <code>{ `new Intl.NumberFormat(lang, {useGrouping: false}).format(1234567890)` }</code>{ " " }with its
						last digit moved to the front; empty when that is ASCII.
					</li>
					<li>
						<code>weekStartsOn</code>: <code>new Intl.Locale(lang).maximize().getWeekInfo().firstDay</code>, where 7 is
						Sunday.
					</li>
				</ul>
				<p>
					Taking each name from the pattern it goes into keeps name and pattern consistent in every language, and a
					locale filled from the server's own data follows CLDR as it changes.
				</p>
			</section>
			<section class="flex flex-col gap-3">
				<docHeading item={i18nTOCItems[4]}/>
				<p>
					Register a translator that wraps every message in markers, render each page, and search the output for English
					outside the markers. Anything found is text a caller passed in, or a gsxui bug worth an issue.
				</p>
			</section>
			<section class="flex flex-col gap-3">
				<docHeading item={i18nTOCItems[5]}/>
				<p>
					When no <code>ui.Toaster</code> is mounted, <code>toaster.js</code> creates its own region with an
					English <code>aria-label</code>. Mount <code>ui.Toaster</code> and the landmark is a message like the rest.
				</p>
			</section>
		</div>
	</siteLayout>
}
