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
					Every gsxui component renders correctly in a right-to-left document without a separate RTL variant to install
					or opt into. Set <code>dir="rtl"</code> and the components already vendored into your project adapt on their
					own.
				</p>
			</div>
			<section class="flex flex-col gap-3">
				<docHeading item={rtlTOCItems[0]}/>
				<p>
					Set <code>dir="rtl"</code> on <code>&lt;html&gt;</code> for a fully right-to-left page, or on any subtree that
					needs it — a single Arabic or Hebrew panel inside an otherwise left-to-right app, for example. No component
					prop, build flag, or extra CSS import is required either way.
				</p>
				<pre><code>{ `<html lang="ar" dir="rtl">
  ...
</html>` }</code></pre>
			</section>
			<section class="flex flex-col gap-3">
				<docHeading item={rtlTOCItems[1]}/>
				<p>
					Four mechanisms carry direction through the component set, all driven by the ancestor's
					resolved <code>dir</code> — none of it is per-component configuration:
				</p>
				<ul class="list-disc space-y-2 pl-6">
					<li>
						<strong>Logical Tailwind classes.</strong> Components are written with <code>ms-*</code>/<code>me-*</code>,
						<code>ps-*</code>/<code>pe-*</code>, <code>start-*</code>/<code>end-*</code>, and
						<code>text-start</code>/<code>text-end</code> instead of physical <code>ml-*</code>/<code>mr-*</code>,
						<code>left-*</code>/<code>right-*</code>. The browser resolves inline-start/inline-end against the element's
						own direction, so spacing and alignment mirror for free.
					</li>
					<li>
						<strong>Directional icons.</strong> Chevrons and arrows that encode a left/right meaning (accordion carets,
						breadcrumb separators, pagination's prev/next, carousel arrows) carry
						<code>rtl:rotate-180</code>, flipping only under an RTL ancestor.
					</li>
					<li>
						<strong>Direction-aware floating positioning.</strong> Popover, dropdown-menu, select, tooltip, and the rest
						of the floating-UI family resolve their placement in JS (<code>ui/gsxui.js</code>'s <code>isRTL(el)</code>,
						read from computed style at position time), so "start"-aligned content opens on the correct physical side
						under either direction.
					</li>
					<li>
						<strong>Mirrored keyboard semantics.</strong> Roving-focus and open/close arrow keys in dropdown-menu,
						context-menu, menubar, tabs, toggle-group, carousel, resizable, and calendar follow the WAI-ARIA convention
						of mirroring by meaning, not by physical key — ArrowLeft still means "toward the next item" under RTL the
						way ArrowRight does under LTR, and ArrowLeft still opens a submenu the way ArrowRight does under LTR.
					</li>
				</ul>
				<p>
					Slider's fill gradient is direction-aware too, filling from the correct inline-start edge.
					<code>input-otp</code> is the one deliberate exception: its digit group stays pinned
					<code>dir="ltr"</code>, because a code like <code>482915</code> should read left-to-right even inside an RTL
					form — this matches how phone numbers and codes read in real Arabic and Hebrew UIs.
				</p>
				<p>
					Sheet, Drawer, and Sidebar are the other deliberate exception: their <code>side="left"</code> /
					<code>side="right"</code> prop stays <strong>physical</strong> under RTL, matching shadcn's own
					<code>data-side</code> contract — a sidebar you explicitly place with <code>side="right"</code> stays on the
					visual right whether the document is LTR or RTL. Everything inside those components — header, content, footer
					spacing, the toggle icon — is still logical and mirrors normally; only the outer <code>side</code> placement
					is fixed.
				</p>
			</section>
			<section class="flex flex-col gap-3">
				<docHeading item={rtlTOCItems[2]}/>
				<p>
					An Arabic sign-in card, composed from unmodified <code>Card</code>, <code>Label</code>,
					<code>Input</code>, and <code>Button</code> — the same components the rest of the docs use, wrapped in
					<code>dir="rtl"</code>.
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
					The site's own Latin type is Geist. For Arabic content, pair it with
					<a href="https://fonts.google.com/noto/specimen/Noto+Sans+Arabic" target="_blank" rel="noreferrer">Noto Sans Arabic</a>
					(a UI sans matching Geist's weight range) or
					<a href="https://fonts.google.com/noto/specimen/Noto+Naskh+Arabic" target="_blank" rel="noreferrer">Noto Naskh Arabic</a>
					(a book-style naskh better suited to longer prose). Load the Arabic font as a <code>lang="ar"</code> fallback
					rather than replacing Geist outright, so Latin text — brand names, code, page furniture — keeps rendering in
					Geist inside an RTL document.
				</p>
			</section>
		</div>
	</siteLayout>
}
