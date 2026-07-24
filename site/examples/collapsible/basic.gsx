// Package collapsible holds the site's example gsx components for
// ui/collapsible.
package collapsible

import (
	"github.com/gsxhq/gsxui/ui"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Basic mirrors shadcn's own collapsible demo (registry/new-york-v4/
// examples/collapsible-demo.tsx): a title row with a chevron ghost
// icon-button trigger, one always-visible list item, then two more
// revealed by expanding.
//
// ADAPT from the demo's literal shape: shadcn wraps ONLY the icon Button in
// CollapsibleTrigger (asChild), leaving the title <h4> outside it — Radix's
// Trigger is an inert wrapper, so the title/button can share a plain sibling
// <div> row with no HTML placement constraint. Native <details> imposes one
// shadcn's version doesn't have: only a summary that is a DIRECT CHILD of
// its details is recognized as the disclosure control (a summary nested a
// level deeper, e.g. inside a wrapping row <div>, renders as inert content
// and the browser falls back to its own default-labelled summary). So the
// title and button both move INSIDE CollapsibleTrigger here, and the flex
// row classes move from the (now-gone) wrapper div onto the trigger itself
// — the whole row is the clickable summary, not just the button, per
// ui/collapsible.gsx's CollapsibleTrigger doc comment.
//
// The chevron is a STYLED SPAN (ghost icon-button look), not a real
// ui.Button: activating a nested interactive element inside a <summary> is
// the element's own activation, not the summary's — a real <button> there
// swallows the click and the details never toggles (and Radix never had
// this trap: its Trigger IS the button). The summary is the one
// interactive/focusable control; the span is decoration (aria-hidden).
component Basic() {
	<ui.Collapsible class="flex w-[350px] flex-col gap-2">
		<ui.CollapsibleTrigger class="flex cursor-default items-center justify-between gap-4 px-4">
			<h4 class="text-sm font-semibold">@peduarte starred 3 repositories</h4>
			<span
				aria-hidden="true"
				class="inline-flex size-8 shrink-0 items-center justify-center rounded-md transition-colors hover:bg-accent hover:text-accent-foreground [&_svg]:pointer-events-none [&_svg:not([class*='size-'])]:size-4"
			>
				<icon.ChevronsUpDown/>
			</span>
		</ui.CollapsibleTrigger>
		<div class="rounded-md border px-4 py-2 font-mono text-sm">@radix-ui/primitives</div>
		<ui.CollapsibleContent class="flex flex-col gap-2">
			<div class="rounded-md border px-4 py-2 font-mono text-sm">@radix-ui/colors</div>
			<div class="rounded-md border px-4 py-2 font-mono text-sm">@stitches/react</div>
		</ui.CollapsibleContent>
	</ui.Collapsible>
}
