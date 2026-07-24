package ui

import "github.com/gsxhq/gsx"

// HoverCard is the shadcn/ui HoverCard on the native popover API:
// popover="manual" puts the content in the top layer without light dismiss
// (hover/focus drive it, not outside clicks or Esc) — this is tooltip.gsx's
// mechanism, minus the arrow (hover-card has none) and anchored BELOW the
// trigger instead of above (Radix's own HoverCard default side is bottom;
// Tooltip's is top). Trigger and content are wired by proximity —
// HoverCardTrigger shows the popover inside the same HoverCard root, no
// ids. JS adds fixed-position anchoring centered below the trigger rect,
// Radix HoverCard's own open/close delays (700ms/300ms — not tooltip's flat
// 300ms-open/immediate-close), and state/event sync. Requires the
// hover-card behavior module (ui/hover-card.js).
component HoverCard(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="hover-card" data-gsxui-hovercard class="contents" { attrs... }>{ children }</div>
}

// HoverCardTrigger renders a <span> phrasing wrapper, not a <button> —
// shadcn's own @nextjs demo asChild-wraps a link-styled <Button
// variant="link">, and Radix's own HoverCardTrigger typically renders as an
// <a> (a hover card almost always previews a link's target). children
// carry the real interactive element (an <a>, or a
// <ui.Button variant="link">, styled as a link); asChild itself is not
// ported (ledgered, docs/jsx-parity.md ## hover-card) — same data-attribute-
// free composition as collapsible's trigger, since a <span> imposes no
// button-in-button trap (unlike DialogTrigger/TooltipTrigger's own
// button-shaped wrappers) and needs no data-gsxui-*-trigger attribute on
// the child at all: HoverCardTrigger's own root already carries the hook.
component HoverCardTrigger(children gsx.Node, attrs gsx.Attrs) {
	<span data-slot="hover-card-trigger" data-gsxui-hovercard-trigger { attrs... }>{ children }</span>
}

// HoverCardContent renders the popover. popover="manual" is load-bearing:
// "auto" popovers light-dismiss on outside pointerdown, which would race
// hover-card.js's own pointerout/focusout hide logic (same rationale as
// TooltipContent). data-state is server-rendered "closed" and kept in sync
// by hover-card.js. data-side="bottom" is server-rendered statically —
// hover-card.js always anchors below the trigger (Radix HoverCard's own
// default side, unlike Tooltip's top), so shadcn's
// data-[side=bottom]:slide-in-from-top-2 enter slide applies without
// Radix's runtime side tracking. origin-top replaces shadcn's Radix runtime
// transform-origin var, same substitution as PopoverContent's own (both
// are centered-below, align=center by Radix default).
component HoverCardContent(children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="hover-card-content"
		data-gsxui-hovercard-content
		popover="manual"
		data-state="closed"
		data-side="bottom"
		class="z-50 w-64 origin-top rounded-md border bg-popover p-4 text-popover-foreground shadow-md outline-hidden data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95 data-[state=open]:animate-in data-[state=open]:fade-in-0 data-[state=open]:zoom-in-95"
		{ attrs... }
	>
		{ children }
	</div>
}
