package ui

import (
	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Accordion and its parts are the shadcn/ui Accordion, ported onto the
// native <details name="…">/<summary> pair instead of Radix's client state
// machine (ledger WIN): grouped <details> sharing a name attribute already
// give the browser exclusive-open-within-group behavior — no
// type="single"/collapsible/value/onValueChange props, no JS, no
// AccordionPrimitive.Header wrapper (a lone <summary> already establishes
// the clickable row shadcn's Header+Trigger pair built by hand).
//
// GAP (same shape as Tabs' value): grouping is a real browser mechanism
// (matching `name` attributes), not client-side context, so there is
// nothing for a root to propagate — the caller passes the same name to
// Accordion AND to every AccordionItem in the group. Accordion's own
// data-name is a readability/debugging stamp only; nothing reads it back
// (no accordion.js — there is none).
component Accordion(name string, children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="accordion" data-name={ name } { attrs... }>{ children }</div>
}

// AccordionItem's open bool is the explicit, server-visible initial state —
// the zero value (false) renders closed, matching shadcn's Radix default
// (nothing expanded until interacted with). Opening/closing thereafter is
// entirely native <details> behavior; toggling one member of a named group
// closes the previously-open sibling without any JS on our side.
component AccordionItem(name string, open bool, children gsx.Node, attrs gsx.Attrs) {
	<details data-slot="accordion-item" name={ name } open={ open } class="border-b last:border-b-0" { attrs... }>{ children }</details>
}

// AccordionTrigger drops the Radix Header wrapper (a bare block-level
// <summary> already lays out as its own row) and, with it, the Trigger's
// now-meaningless flex-1 (nothing left to flex within) and its
// [&[data-state=open]>svg]:rotate-180 selector — <summary> carries no
// data-state, native <details> truth is the [open] attribute instead. The
// rotate selector moves onto the chevron itself as an arbitrary variant
// keyed off the ancestor <details>' [open] attribute: [[open]>summary_&]
// reads "an ancestor with [open], with a summary directly beneath it, with
// this element somewhere under that summary" — validated against Tailwind
// v4's documented ancestor-selector idiom (referencing the parent when
// nesting: [.dark_&]:… ports directly to an attribute selector plus a child
// combinator) and against tailwind-merge-go (the class_merger this repo
// uses): the bracketed token round-trips merge.Merge unchanged and is
// correctly recognized as a `rotate` utility (a caller-supplied
// `rotate-90`, say, replaces it; unrelated utilities are untouched). No
// build-time Tailwind run was available to confirm the compiled CSS
// (out of scope for this sandbox) — flagged for a real-browser check.
// list-none and [&::-webkit-details-marker]:hidden strip the native
// disclosure triangle both browser families draw (the standard
// marker/::-webkit-details-marker pair), leaving only our chevron.
// pointer-events-none and translate-y-0.5 from shadcn's ChevronDownIcon are
// dropped: pointer-events-none guarded against the icon stealing the click
// from Radix's button-role Trigger, which doesn't exist here (the click
// target is the whole <summary>, native disclosure semantics, not a click
// handler an errant svg hit could steal from); translate-y-0.5 nudged the
// icon to align against the Header wrapper's row baseline, which no longer
// exists now that the icon sits directly in <summary>'s own flex row.
component AccordionTrigger(children gsx.Node, attrs gsx.Attrs) {
	<summary data-slot="accordion-trigger" class="flex items-start justify-between gap-4 rounded-md py-4 text-left text-sm font-medium transition-all outline-none hover:underline focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:pointer-events-none disabled:opacity-50 list-none [&::-webkit-details-marker]:hidden" { attrs... }>
		{ children }
		<icon.ChevronDown class="text-muted-foreground size-4 shrink-0 transition-transform duration-200 [[open]>summary_&]:rotate-180"/>
	</summary>
}

// AccordionContent drops shadcn's data-[state=open]:animate-accordion-down /
// data-[state=closed]:animate-accordion-up pair (both keyed off Radix's
// data-state, which nothing stamps here). The open/close height animation
// is instead pure CSS in gsxui.css (and web/site.css, the site's copied
// twin — TestAccordionAnimationCSSDriftPin keeps the two in sync): the
// `::details-content` pseudo-element that wraps everything but the summary
// animates its grid row 0fr -> 1fr, which is height-to-auto with no
// interpolate-size and no JS, plus a discrete content-visibility transition
// so the collapsing content stays rendered until the close finishes.
// Browser-verified both directions (rAF height sampling: 0->36px easing
// over the duration). Browsers without ::details-content ignore the rules
// and toggle instantly — a progressive enhancement over the same markup.
// That stylesheet keys on this component's data-slot attributes, and its
// min-height:0 on accordion-content is what lets the 0fr row collapse.
component AccordionContent(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="accordion-content" class="overflow-hidden text-sm" { attrs.Without("class")... }>
		<div class={ "pt-0 pb-4", attrs.Class() }>{ children }</div>
	</div>
}
