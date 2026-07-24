package ui

import "github.com/gsxhq/gsx"

// ScrollArea is the shadcn/ui ScrollArea (registry/new-york-v4/ui/
// scroll-area.tsx), collapsed onto ONE native `overflow-auto`/
// `overflow-x-auto` div — Radix's four-part tree (Root/Viewport/ScrollBar/
// Thumb) collapses because a real `overflow: auto` element already keeps
// every native scroll input modality (wheel/trackpad/touch/keyboard)
// working for free, and the standardized `scrollbar-width`/`scrollbar-color`
// pair (CSS Scrollbars Module Level 1) styles the browser's own thumb/track
// color with zero JS — the roadmap's own stated CSS-first preference for
// this component ("CSS scrollbar-width/scrollbar-color styling first;
// Radix-style custom thumbs only if that falls short"). Root's own
// `relative` and Viewport's own `rounded-[inherit]`/focus-visible ring both
// port unchanged onto this one node — see docs/jsx-parity.md `## scroll-area`
// and the 2026-07-24 controls source map's own `## scroll-area` section for
// the full ledger.
//
// GAP (no ScrollBar/ScrollAreaThumb): the standard `scrollbar-color`
// property paints the browser's OWN thumb geometry (squared-off,
// platform-dependent width) — no radius/inset/border control at all. Full
// WebKit/Chromium shape fidelity (`rounded-full` thumb, the `p-px` track
// inset) is layered on separately as a hand-authored `::-webkit-scrollbar-*`
// block in assets/gsxui.css/web/site.css; Firefox has never implemented
// those pseudo-elements, so `scrollbar-color` is what it gets there.
//
// GAP (no ScrollAreaCorner): `::-webkit-scrollbar-corner` is WebKit/Blink
// only, the standard `scrollbar-*` properties have no corner concept
// whatsoever, and Firefox cannot style the intersection square under any
// mechanism — dropped outright for v1, matching this codebase's "drop the
// part, don't ship dead CSS" convention (e.g. dropdown's never-ported
// scroll buttons).
//
// GAP (no `type` visibility-timing prop): Radix's hover/scroll/auto/always
// scrollbar show/hide TIMING is OS/browser policy on a native scrollbar —
// e.g. macOS's own "Show scroll bars" System Settings preference overrides
// any page-level attempt to mimic Radix's show-on-hover/fade-after-600ms
// behavior. Not a fixable styling gap; ledgered as an accepted
// platform-policy override.
//
// orientation: "" (default, vertical — `overflow-auto`) | "horizontal"
// (`overflow-x-auto`) — the axis switch needs no second component the way
// Radix needs a literal `<ScrollBar orientation="horizontal"/>` element; it
// is just a different overflow utility on the same collapsed div.
// `[scrollbar-color:var(--border)_transparent]` needs no separate `dark:`
// variant: `--border` is itself a theme token that already changes value
// under `.dark` (assets/gsxui.css), so the property tracks the theme for
// free — the same "no hardcoded colors" reasoning every other component's
// `var(--token)` usage relies on.
component ScrollArea(orientation string, children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="scroll-area"
		class={
			"relative rounded-[inherit] outline-none transition-[color,box-shadow] focus-visible:ring-[3px] focus-visible:ring-ring/50 focus-visible:outline-1 [scrollbar-width:thin] [scrollbar-color:var(--border)_transparent]",
			switch orientation {
			case "horizontal":
				"overflow-x-auto"
			default:
				"overflow-auto"
			}
		}
		{ attrs... }
	>
		{ children }
	</div>
}
