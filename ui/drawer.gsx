package ui

import "github.com/gsxhq/gsx"

// Drawer and its parts are the shadcn/ui Drawer (registry/new-york-v4/ui/
// drawer.tsx): upstream is a bare `vaul` passthrough, structurally
// identical in shape to dialog.tsx/sheet.tsx. This port follows
// ui/sheet.gsx's own composition pattern exactly rather than deriving from
// vaul: Drawer composes ui.Dialog directly (data-slot override "dialog" ->
// "drawer", the same mechanism as Sheet/AlertDialog), so drawer -> dialog
// derives and ui/dialog.js is pulled in transitively — trigger/content
// wiring by proximity, Esc-to-close, data-state stamping, and the
// getAnimations-gated exit-animation wait are all reused unmodified, zero
// new dialog.js code (HasJS("drawer") is false, same conclusion as
// ## sheet's own Registry note). DrawerContent, like SheetContent, does
// NOT compose ui.DialogContent or ui.SheetContent — the centered-card
// recipe and each side-anchored recipe target the same CSS properties with
// materially different values, so there is no single class string that
// merges them (see ui/sheet.gsx's own header comment). DrawerContent
// renders its own <dialog data-slot="drawer-content" data-gsxui-dialog-
// content data-state="closed" data-side="...">, still using the exact same
// data-gsxui-dialog-content contract dialog.js selects on. Unlike
// Dialog/Sheet, upstream drawer.tsx's own DrawerContent never injects a
// close X button (no showCloseButton-equivalent prop) — dismissal is
// backdrop-click/Esc, or an explicit DrawerClose the caller places
// wherever the design wants it (see the footer Cancel button in
// site/examples/drawer/basic.gsx), so DrawerContent has no
// hideCloseButton param.
//
// direction (vaul's own vocabulary, distinct from Sheet's "side") is the Go
// param: default "bottom" (vaul's conventional anchor, the mobile
// bottom-sheet pattern — a real behavioral difference from Sheet's own
// "right" default), also stamped as data-side (reusing Sheet's own internal
// attribute name for any future shared tooling/CSS that keys off it
// generically, rather than inventing a second attribute). Direction
// selection happens server-side via a Go switch in class={}, the same
// idiom as Sheet's side switch — dialog.js never reads data-side, the
// slide direction is fully determined by which static class block was
// selected.
//
// Nova adopted (docs/superpowers/plans/2026-07-24-tier3-source-map-wrapped.md
// `## drawer`): rounded-*-xl on all four directions, including left/right's
// free (non-anchored) edge — new-york-v4 only rounds top/bottom, left/right
// get a plain border with no rounding at all; nova rounds every direction.
// bg-popover/text-popover-foreground is the one deliberate NOT-adopted-
// elsewhere exception: unlike every retargeted component, drawer has no
// prior new-york-v4-based gsxui version to stay consistent with, so there
// is no existing bg-background baseline to preserve. Backdrop is identical
// to Sheet's own (bg-black/10 + supports-backdrop-filter:backdrop:backdrop-
// blur-xs, duration-200 both directions).
//
// Per-direction class strings apply Sheet's own already-solved <dialog>-
// vs-UA-defaults fixes (m-0, the opposite-edge -auto utility, open:flex)
// to drawer.tsx's values. TOP/BOTTOM carry their own explicit
// max-h-[80vh] (author-origin, already beats Chrome's UA max-height safety
// net) so they do NOT need Sheet's max-h-none escape hatch; LEFT/RIGHT use
// h-full like Sheet's own left/right and DO need max-h-none, copied
// verbatim from Sheet's fix. TOP/BOTTOM also carry Sheet's own w-full
// max-w-none pair (the plain-<div> auto-width edge-to-edge stretch does not
// reproduce on a native <dialog>, see ui/sheet.gsx's own three-part ADAPT).
// These strings are transcribed-not-verified per the map's own caveat —
// Sheet's six ADAPTs were only found by rendering in a real browser tab;
// the controller runs that same verification pass for all four drawer
// directions before this task is marked complete.
//
// Handle bar: upstream renders an unconditional inline <div> (not a named
// sub-component), visible only for the bottom direction via a
// group-data-[vaul-drawer-direction=bottom]/drawer-content:block override
// on a base `hidden` class. gsxui has no vaul underneath and direction is
// server-known, so the visibility rule is ported as a Go `if` on the
// already-resolved direction instead of a client CSS group-data selector —
// the div (data-slot="drawer-handle") is rendered only when direction is
// "bottom", and simply absent (not merely hidden) otherwise. h-1 (nova; the
// h-2 new-york-v4 handle is thinner). Kept as pure decoration: v1 ships no
// drag gesture, so a "drag me" affordance that doesn't actually drag is a
// real, if minor, UX mismatch — accepted GAP (visual parity over silently
// dropping it), ledgered in docs/jsx-parity.md `## drawer`.
//
// MECHANISM: DrawerHeader's upstream class also carries a direction-
// conditional text alignment (centered for bottom/top at every breakpoint,
// left-aligned at md+ for left/right) via
// group-data-[vaul-drawer-direction=...]/drawer-content. gsxui has no vaul
// underneath, but DrawerContent already stamps data-side (always non-empty
// — direction |> default("bottom")) and carries the named group class
// group/drawer-content, so the same selector shape ports directly with
// only the attribute/group name swapped: data-side replaces
// data-vaul-drawer-direction, drawer-content (already DrawerContent's own
// data-slot) is the group name. DrawerHeader needs no direction param of
// its own — the selector reads the ancestor <dialog>'s stamped attribute
// at the CSS layer, the same group-data-[...]/name idiom `## item`/
// `## field`/`## input-group`/`## tabs`/`## toggle-group` already use
// elsewhere in this codebase.
//
// GAP: drag-to-dismiss, snap points, and background scaling (all vaul
// gesture/physics features) are not ported — v1 replaces vaul's live-
// transform drag entirely with the same <dialog> + Tailwind-keyframe
// slide-in/out architecture Sheet already uses, per the roadmap's own
// Tier 3 listing ("sheet variant; v1 without vaul's drag-to-dismiss
// gesture, ledger the gap").
component Drawer(children gsx.Node, attrs gsx.Attrs) {
	<Dialog data-slot="drawer" { attrs... }>{ children }</Dialog>
}

// DrawerTrigger renders its own <button>, same reasoning and the same
// button-in-button HTML trap as SheetTrigger/DialogTrigger (see
// ui/dialog.gsx's DialogTrigger doc comment and docs/jsx-parity.md
// `## dialog` FINDING): its children must be phrasing content, never a
// Button or other interactive element. For a styled trigger, skip the
// wrapper and put the data attribute on the Button itself:
// <ui.Button data-gsxui-dialog-trigger>Open</ui.Button>.
component DrawerTrigger(children gsx.Node, attrs gsx.Attrs) {
	<button
		data-slot="drawer-trigger"
		data-gsxui-dialog-trigger
		type="button"
		aria-haspopup="dialog"
		aria-expanded="false"
		{ attrs... }
	>
		{ children }
	</button>
}

// DrawerContent renders the native <dialog> directly (see this file's own
// header comment for why it does not compose ui.DialogContent or
// ui.SheetContent). direction |> default("bottom") stamps data-side on the
// element; the raw direction switch below selects one of the four static
// class blocks, falling through to the bottom case (its default:) for both
// "" and "bottom". The handle-bar if-condition spells out
// direction == "" || direction == "bottom" explicitly rather than reusing
// the |> default(...) pipeline: an if-condition is a plain Go boolean
// expression, not the pipe-aware attribute-expression form the |> operator
// is scoped to.
component DrawerContent(direction string, children gsx.Node, attrs gsx.Attrs) {
	<dialog
		data-slot="drawer-content"
		data-gsxui-dialog-content
		data-state="closed"
		data-side={direction |> default("bottom")}
		class={
			"group/drawer-content fixed z-50 m-0 open:flex flex-col gap-4 bg-popover text-popover-foreground text-sm shadow-lg transition ease-in-out duration-200 data-[state=closed]:animate-out data-[state=open]:animate-in backdrop:bg-black/10 backdrop:duration-200 supports-backdrop-filter:backdrop:backdrop-blur-xs data-[state=open]:backdrop:animate-in data-[state=open]:backdrop:fade-in-0 data-[state=closed]:backdrop:animate-out data-[state=closed]:backdrop:fade-out-0",
			switch direction {
			case "top":
				"inset-x-0 top-0 bottom-auto w-full max-w-none h-auto mb-24 max-h-[80vh] rounded-b-xl border-b data-[state=closed]:slide-out-to-top data-[state=open]:slide-in-from-top"
			case "left":
				"inset-y-0 left-0 right-auto h-full max-h-none w-3/4 rounded-r-xl border-r sm:max-w-sm data-[state=closed]:slide-out-to-left data-[state=open]:slide-in-from-left"
			case "right":
				"inset-y-0 right-0 left-auto h-full max-h-none w-3/4 rounded-l-xl border-l sm:max-w-sm data-[state=closed]:slide-out-to-right data-[state=open]:slide-in-from-right"
			default:
				"inset-x-0 bottom-0 top-auto w-full max-w-none h-auto mt-24 max-h-[80vh] rounded-t-xl border-t data-[state=closed]:slide-out-to-bottom data-[state=open]:slide-in-from-bottom"
			}
		}
		{ attrs... }
	>
		{ if direction == "" || direction == "bottom" {
			<div data-slot="drawer-handle" class="mx-auto mt-4 h-1 w-[100px] shrink-0 rounded-full bg-muted"></div>
		} }
		{ children }
	</dialog>
}

// DrawerHeader's text alignment is direction-conditional, ported faithfully
// via the data-side/group-drawer-content selector — see this file's own
// header comment MECHANISM note. gap-0.5 is nova at every breakpoint (no
// md: bump; new-york-v4's own DrawerHeader bumps to md:gap-1.5).
component DrawerHeader(children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="drawer-header"
		class="flex flex-col gap-0.5 p-4 group-data-[side=bottom]/drawer-content:text-center group-data-[side=top]/drawer-content:text-center md:text-left"
		{ attrs... }
	>
		{ children }
	</div>
}

// DrawerFooter is byte-identical to SheetFooter's own (upstream drawer.tsx
// and sheet.tsx happen to share this class string exactly).
component DrawerFooter(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="drawer-footer" class="mt-auto flex flex-col gap-2 p-4" { attrs... }>{ children }</div>
}

component DrawerTitle(children gsx.Node, attrs gsx.Attrs) {
	<h2 data-slot="drawer-title" class="font-medium text-foreground" { attrs... }>{ children }</h2>
}

component DrawerDescription(children gsx.Node, attrs gsx.Attrs) {
	<p data-slot="drawer-description" class="text-sm text-muted-foreground" { attrs... }>{ children }</p>
}

component DrawerClose(children gsx.Node, attrs gsx.Attrs) {
	<button data-slot="drawer-close" data-gsxui-dialog-close type="button" { attrs... }>{ children }</button>
}
