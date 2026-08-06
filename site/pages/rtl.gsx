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
					Every gsxui component renders correctly in a right-to-left document. Set <code>dir="rtl"</code> and the
					components you've already vendored adapt on their own — no RTL variant to install or opt into.
				</p>
			</div>
			<section class="flex flex-col gap-3">
				<docHeading item={rtlTOCItems[0]}/>
				<p>
					Set <code>dir="rtl"</code> on <code>&lt;html&gt;</code> for a fully right-to-left page, or on any subtree — a
					single Arabic or Hebrew panel inside a left-to-right app, for example:
				</p>
				<pre><code>{ `<html lang="ar" dir="rtl">
  ...
</html>` }</code></pre>
			</section>
			<section class="flex flex-col gap-3">
				<docHeading item={rtlTOCItems[1]}/>
				<p>
					Four mechanisms carry direction through the component set, all driven by the ancestor's
					resolved <code>dir</code>:
				</p>
				<ul class="list-disc space-y-2 pl-6">
					<li>
						<strong>Logical Tailwind classes.</strong> Components use <code>ms-*</code>/<code>me-*</code>,
						<code>start-*</code>/<code>end-*</code>, and <code>text-start</code>/<code>text-end</code> instead of
						physical <code>ml-*</code>/<code>left-*</code>, so spacing and alignment mirror for free.
					</li>
					<li>
						<strong>Directional icons.</strong> Chevrons and arrows that encode a left/right meaning carry
						<code>rtl:rotate-180</code>, flipping only under an RTL ancestor.
					</li>
					<li>
						<strong>Direction-aware floating positioning.</strong> Popover, dropdown-menu, select, tooltip, and the rest
						of the floating-UI family resolve their placement in JS at position time, so "start"-aligned content opens
						on the correct physical side.
					</li>
					<li>
						<strong>Mirrored keyboard semantics.</strong> Arrow keys in menus, tabs, carousel, calendar, and the other
						roving-focus components mirror by meaning per WAI-ARIA — under RTL, ArrowLeft means "toward the next item"
						and opens submenus, the way ArrowRight does under LTR.
					</li>
				</ul>
				<p>
					Two deliberate exceptions. <code>input-otp</code>'s digit group stays pinned <code>dir="ltr"</code>: a code
					like <code>482915</code> reads left-to-right even inside an RTL form, matching real Arabic and Hebrew UIs. And
					Sheet, Drawer, and Sidebar keep their <code>side="left"</code>/<code>side="right"</code> prop
					<strong>physical</strong>, matching shadcn's <code>data-side</code> contract — a <code>side="right"</code>
					sidebar stays on the visual right under either direction, while everything inside it still mirrors normally.
				</p>
			</section>
			<section class="flex flex-col gap-3">
				<docHeading item={rtlTOCItems[2]}/>
				<p>
					An Arabic sign-in card: unmodified <code>Card</code>, <code>Label</code>, <code>Input</code>,
					and <code>Button</code>, wrapped in <code>dir="rtl"</code>.
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
				<p>
					More RTL variants live on their own component pages: <a href="/components/calendar">Calendar</a>,
					<a href="/components/pagination">Pagination</a>, and <a href="/components/sidebar">Sidebar</a> each register
					an "RTL" example alongside their other demos.
				</p>
			</section>
			<section class="flex flex-col gap-3">
				<docHeading item={rtlTOCItems[3]}/>
				<p>
					The site's own Latin type is Geist. For Arabic content, pair it
					with <a href="https://fonts.google.com/noto/specimen/Noto+Sans+Arabic" target="_blank" rel="noreferrer">Noto Sans Arabic</a>
					(a UI sans)
					or <a href="https://fonts.google.com/noto/specimen/Noto+Naskh+Arabic" target="_blank" rel="noreferrer">Noto Naskh Arabic</a>
					(better for longer prose). Load the Arabic font as a fallback rather than replacing Geist outright, so Latin
					text — brand names, code — keeps rendering in Geist inside an RTL document.
				</p>
			</section>
		</div>
	</siteLayout>
}
