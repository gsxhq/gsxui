package ui

import "github.com/gsxhq/gsx"

// Collapsible, CollapsibleTrigger, and CollapsibleContent are the shadcn/ui
// Collapsible, ported onto the same native <details>/<summary> mechanism as
// Accordion (see ui/accordion.gsx) instead of Radix's client state machine
// (WIN): a single un-grouped <details> IS a disclosure — no name attribute,
// no type="single"/collapsible/open/onOpenChange props, no JS. Unlike
// Accordion, there is no exclusive-group behavior to model (shadcn's
// collapsible.tsx is a lone Radix Root/Trigger/Content triple, not a
// Root+Item pair), so there is nothing here shaped like AccordionItem's
// `name` grouping param.
//
// shadcn's collapsible.tsx (registry/new-york-v4/ui/collapsible.tsx) is a
// bare data-slot passthrough — none of its three parts carry a single class
// string. That "no classes" fact is itself the token-for-token baseline:
// Collapsible and CollapsibleContent port with zero classes of their own.
// CollapsibleTrigger is the one ADAPT (see its own doc comment below).
//
// GAP: no data-state is stamped anywhere in this port (there is no
// collapsible.js — none exists, same as Accordion). CSS consumers who would
// reach for shadcn's `data-[state=open]`/`data-[state=closed]` selectors
// have nothing to key off; target the ancestor <details>'s native `[open]`
// attribute instead (the same substitution Accordion's trigger chevron
// uses: `[[open]_&]:...`-shaped selectors, not `data-[state=open]:...`).

// Collapsible's open bool is the explicit, server-visible initial state —
// the Go zero value (false) renders collapsed, matching shadcn's Radix
// default (defaultOpen unset). It is shadcn's `defaultOpen`, not a
// controlled `open`/`onOpenChange` pair: opening and closing thereafter are
// entirely native <details> behavior, no hydration step to reconcile.
component Collapsible(open bool, children gsx.Node, attrs gsx.Attrs) {
	<details data-slot="collapsible" open={open} { attrs... }>
		{ children }
	</details>
}

// CollapsibleTrigger drops Radix's asChild — the whole reason shadcn's own
// demo (registry/new-york-v4/examples/collapsible-demo.tsx) wraps a real
// <Button> inside <CollapsibleTrigger asChild> is to make an actual button
// element the clickable trigger while Radix's own Trigger contributes only
// behavior. Here the data-attribute idiom doesn't apply and neither does
// composing a real button: a bare <summary> already IS the clickable
// disclosure control, and activating a NESTED interactive element (a
// <button>, <a>, <input>…) inside it is that element's own activation, not
// the summary's — the click is swallowed and the details never toggles. So
// callers style NON-interactive children to look like shadcn's trigger
// button instead (the summary carries the focus/keyboard semantics; the
// visual "button" is decoration):
//
//	<ui.CollapsibleTrigger>
//		<span aria-hidden="true" class="…ghost icon-button classes…"><icon.ChevronsUpDown/></span>
//	</ui.CollapsibleTrigger>
//
// (site/examples/collapsible/basic.gsx is the full worked shape.)
//
// ADAPT: shadcn's CollapsibleTrigger carries no classes at all (nothing to
// carry token-for-token) — list-none and the webkit marker selector below
// are ADDED, not ported, for the same reason Accordion's trigger adds them:
// a real <summary> draws a native disclosure triangle marker (both the
// standard `::marker`, suppressed by `list-none`, and WebKit's separate
// `::-webkit-details-marker`) that Radix's <button>-based trigger never
// had. Without suppressing both, callers get shadcn's chevron icon PLUS a
// browser-drawn triangle.
component CollapsibleTrigger(children gsx.Node, attrs gsx.Attrs) {
	<summary data-slot="collapsible-trigger" class="list-none [&::-webkit-details-marker]:hidden" { attrs... }>
		{ children }
	</summary>
}

// CollapsibleContent is a plain div, no classes (shadcn ships none) and no
// animation block: shadcn's own collapsible.tsx has no
// data-[state=open]:animate-* pair to port in the first place (unlike
// Accordion's Radix-sourced animate-accordion-down/up, which motivated
// Accordion's CSS-only ::details-content replacement) — there is nothing to
// adapt here, so nothing is added. A caller wanting an open/close
// transition reaches for the same ::details-content technique Accordion
// documents, keyed off this component's ancestor <details>'s [open].
component CollapsibleContent(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="collapsible-content" { attrs... }>
		{ children }
	</div>
}
