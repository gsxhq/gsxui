package ui

import "github.com/gsxhq/gsx"

// ResizablePanelGroup, ResizablePanel, and ResizableHandle are the
// shadcn/ui Resizable (registry/new-york-v4/ui/resizable.tsx), which wraps
// react-resizable-panels' Group/Panel/Separator. That library is absent
// from the reference checkout (no node_modules installed at all — see the
// 2026-07-24 tier4 source map's own "Library presence check" for
// `## resizable`), so its runtime ARIA-stamping and keyboard code were
// never read; this port's dynamic behavior (drag, keyboard, ARIA sync) is
// gsxui-authored from ui/resizable.js, not ported 1:1 from the library.
//
// ADAPT (handle aria-orientation is inverted from the group's own
// orientation — forced by shadcn's own class string, not a gsxui choice):
// `orientation` on ResizablePanelGroup/ResizableHandle always names the
// GROUP's layout ("horizontal" = panels side-by-side, the default; the Go
// zero value "" also means horizontal). But `role="separator"`'s own
// `aria-orientation` describes the SEPARATOR itself, and a divider between
// side-by-side panels is a VERTICAL line — so a horizontal group's handle
// carries `aria-orientation="vertical"`, and vice versa. This is exactly
// what the handle's base (unmatched-variant) class `w-px` — a full-height,
// 1px-WIDE rule — versus its `aria-[orientation=horizontal]:h-px
// aria-[orientation=horizontal]:w-full` branch — a full-WIDTH, 1px-tall
// rule — encodes; see the source map's `## resizable` §3 mapping table for
// the full derivation (independently corroborated there by the
// `[&[aria-orientation=horizontal]>div]:rotate-90` grip-rotation rule).
// Callers reason only about the group's own orientation; ResizableHandle
// inverts it internally so the caller never has to.
//
// MECHANISM (sizes are server-rendered inline flex-basis, not JS-computed):
// react-resizable-panels itself sets pixel/percentage sizing imperatively
// after mount — ResizablePanel's own upstream class string is empty (no
// `cn()` call at all, confirmed by the source map), because the library
// supplies 100% of its layout at runtime. gsxui has no client-side layout
// pass before first paint, so ResizablePanel instead renders `defaultSize`
// (a percentage string like "50%", "" meaning unset) as a real inline
// `style="flex-basis: <defaultSize>"` — the split is correct on first
// paint with JS disabled. minSize/maxSize are NOT rendered as styles (they
// don't constrain anything visually at rest); they're stamped as
// `data-min-size`/`data-max-size`, read only by resizable.js's drag/
// keyboard clamping, and otherwise absent from the DOM if unset.
//
// NEW (flex-1 min-w-0 min-h-0, not in shadcn's source at all): since
// ResizablePanel carries no upstream class, an unsized panel here needs
// SOMETHING to make it share remaining space — `flex-1` (`flex: 1 1 0%`)
// sets flex-grow/flex-shrink to 1 while leaving flex-basis at the
// stylesheet layer, which a caller-supplied inline `style="flex-basis:
// ..."` cleanly overrides (inline beats an ordinary class rule for that
// one longhand, grow/shrink from the class still apply) — so `flex-1` is
// unconditional on every panel, sized or not. `min-w-0 min-h-0` counters
// the flexbox default `min-width:auto`/`min-height:auto` floor (a flex
// item won't shrink below its content's intrinsic size otherwise,
// regardless of flex-basis/flex-shrink) along BOTH axes, since
// ResizablePanel itself has no `orientation` param to know which axis
// matters — harmless on the axis that doesn't apply.
//
// GAP (`autoSaveId` dropped, `gsxui:change` instead): react-resizable-
// panels' own `autoSaveId` persists layout to localStorage internally, one
// more component silently owning storage — gsxui components take state as
// a parameter and emit events instead (house rule, see e.g. `## sheet`/
// `## select`'s own persistence GAPs). ResizablePanelGroup takes no
// autoSaveId; resizable.js emits `gsxui:change` on the group
// (`{ sizes: [<percent number>, ...] }`, one entry per panel in DOM order)
// on drag end and on every keyboard commit — the caller owns persisting
// that however it wants (cookie, server round-trip, localStorage from
// their own code), never gsxui itself.
//
// GAP (collapsible panels, the imperative Panel API): react-resizable-
// panels' `collapsible`/`collapsedSize` props and its ref-based imperative
// API (`panelRef.current.resize()`/`.collapse()`/`.expand()`/`.getSize()`)
// have no port here — no docs demo in scope
// (resizable-demo/-demo-with-handle/-handle/-vertical) exercises either,
// same accepted-gap shape as `## carousel`'s own `loop`/`align` GAPs.
component ResizablePanelGroup(orientation string, children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="resizable-panel-group"
		data-gsxui-resizable
		aria-orientation={orientation |> default("horizontal")}
		class="flex h-full w-full aria-[orientation=vertical]:flex-col"
		{ attrs... }
	>
		{ children }
	</div>
}

// ResizablePanel — see the package doc comment's MECHANISM/NEW entries
// above for why this carries a class at all (upstream's own ResizablePanel
// has none) and why defaultSize becomes a real inline style instead of a
// data attribute.
component ResizablePanel(defaultSize string, minSize string, maxSize string, children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="resizable-panel"
		{ if minSize != "" {
			data-min-size={minSize}
		} }
		{ if maxSize != "" {
			data-max-size={maxSize}
		} }
		{ if defaultSize != "" {
			style=css`flex-basis: @{defaultSize}`
		} }
		class="flex-1 min-w-0 min-h-0"
		{ attrs... }
	>
		{ children }
	</div>
}

// ResizableHandle — see the package doc comment's ADAPT entry above for
// the inverted-aria-orientation derivation.
//
// ADAPT (nova's empty pill, not new-york-v4's icon-in-a-box — 2026-07-24
// tier4 source map, `## resizable` §6, via the `registry/bases/base/ui/
// resizable.tsx` structure tiebreak): new-york-v4's `withHandle` renders a
// bordered `h-4 w-3` box containing a `size-2.5` `GripVerticalIcon`; nova's
// own `.cn-resizable-handle-icon` rule (`@apply bg-border h-6 w-1
// rounded-lg`) instead applies to a completely empty div — no icon glyph,
// no border, no box. That's a markup-STRUCTURE change nova's CSS alone
// can't reveal (the class name doesn't say what DOM it targets), settled
// via the bases/base tiebreak the map cites. Nova wins on visual
// disagreement per the house rule, so this port drops the
// `GripVerticalIcon` import/render entirely — the resulting dependency-free
// `Deps("resizable") == []` is a direct consequence, not a coincidence.
component ResizableHandle(orientation string, withHandle bool, attrs gsx.Attrs) {
	{{
		handleOrientation := "vertical"
		if orientation == "vertical" {
			handleOrientation = "horizontal"
		}
	}}
	<div
		data-slot="resizable-handle"
		role="separator"
		aria-orientation={handleOrientation}
		tabindex="0"
		class="relative flex w-px items-center justify-center bg-border after:absolute after:inset-y-0 after:left-1/2 after:w-1 after:-translate-x-1/2 focus-visible:ring-1 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:outline-hidden aria-[orientation=horizontal]:h-px aria-[orientation=horizontal]:w-full aria-[orientation=horizontal]:after:left-0 aria-[orientation=horizontal]:after:h-1 aria-[orientation=horizontal]:after:w-full aria-[orientation=horizontal]:after:translate-x-0 aria-[orientation=horizontal]:after:-translate-y-1/2 [&[aria-orientation=horizontal]>div]:rotate-90"
		{ attrs... }
	>
		{ if withHandle {
			<div class="z-10 flex shrink-0 h-6 w-1 rounded-lg bg-border"></div>
		} }
	</div>
}
