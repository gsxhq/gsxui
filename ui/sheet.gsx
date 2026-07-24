package ui

import "github.com/gsxhq/gsx"

// Sheet and its parts are the shadcn/ui Sheet
// (registry/new-york-v4/ui/sheet.tsx): a side-anchored drawer variant of
// Dialog. Sheet's root composes ui.Dialog directly (same
// data-slot-override mechanism as ui/alert-dialog.gsx's AlertDialog —
// "dialog" -> "sheet") so the sheet -> dialog dependency derives and the
// CLI vendors ui/dialog.js: trigger/content wiring by proximity, Esc-to-
// close, data-state stamping, and toggle-driven exit animations are all
// dialog.js's existing machinery, reused unmodified, and this is why sheet
// gets working exit animations at all (dialog.js stamps data-state and
// waits out getAnimations before dialog.close()).
//
// SheetContent, unlike AlertDialogContent, does NOT compose
// ui.DialogContent: DialogContent's centered-card recipe (top-[50%]/
// left-[50%]/translate-*/max-w-lg/grid) and Sheet's own side-anchored
// recipe (inset-y-0 or inset-x-0/h-full or h-auto/w-3/4/flex) target the
// same CSS properties (position offsets, sizing, display) with materially
// different values on every one of them — there is no single caller class
// string that could merge the two through tailwind-merge the way
// AlertDialogContent's content class merged cleanly against DialogContent's
// own (see docs/jsx-parity.md `## alert-dialog` ADAPT: that merge worked
// only because alert-dialog's content recipe is nearly IDENTICAL to
// dialog's). Sheet renders its own <dialog data-slot="sheet-content"
// data-gsxui-dialog-content data-state="closed" data-side="...">, still
// using the exact same data-gsxui-dialog-content contract dialog.js selects
// on. See docs/jsx-parity.md `## sheet` for the full ADAPT/GAP ledger.
component Sheet(children gsx.Node, attrs gsx.Attrs) {
	<Dialog data-slot="sheet" { attrs... }>{ children }</Dialog>
}

// SheetTrigger renders its own <button>, same reasoning and the same
// button-in-button HTML trap as DialogTrigger/AlertDialogTrigger (see
// ui/dialog.gsx's DialogTrigger doc comment and docs/jsx-parity.md
// `## dialog` FINDING): its children must be phrasing content, never a
// Button or other interactive element. For a styled trigger, skip the
// wrapper and put the data attribute on the Button itself:
// <ui.Button data-gsxui-dialog-trigger>Open</ui.Button> — the same
// documented idiom, unaffected by which dialog flavor (Dialog, AlertDialog,
// or Sheet) sits behind it, since all three key off the same
// data-gsxui-dialog-trigger contract.
component SheetTrigger(children gsx.Node, attrs gsx.Attrs) {
	<button
		data-slot="sheet-trigger"
		data-gsxui-dialog-trigger
		type="button"
		aria-haspopup="dialog"
		aria-expanded="false"
		{ attrs... }
	>
		{ children }
	</button>
}

// SheetContent renders the native <dialog> directly (see this file's own
// header comment for why it does not compose ui.DialogContent).
// `side |> default("right")` picks one of sheet.tsx's four static class
// blocks (right/left/top/bottom — no data-[side=...] selectors in the
// source to preserve, so this is a switch inside class={}, the same idiom
// as item.gsx's variant/size pair) and also stamps data-side on the
// element itself, for consistency with every other data-variant stamp in
// this codebase and for any downstream selector/JS that wants it (dialog.js
// itself does not read it — the slide direction is fully determined by
// which static class block above is selected).
//
// Six ADAPTs relative to sheet.tsx's own class string (all ledgered in
// docs/jsx-parity.md `## sheet`):
//   - the base string's unscoped `flex` becomes `open:flex`, the same
//     display-gating fix as DialogContent's own `open:grid` (`## dialog`
//     ADAPT): content stays in the DOM when closed (no Radix unmount to
//     rely on), so an ungated display utility would defeat the UA's
//     closed-dialog `display:none`.
//   - `text-foreground` is added, the same fix as DialogContent's own
//     (`## dialog` ADAPT): native <dialog> gets UA `color: CanvasText` and
//     does not inherit the themed body color.
//   - `m-0` is added: the UA style for a modal <dialog> is
//     `position:fixed; inset:0; width:fit-content; height:fit-content;
//     margin:auto`. Left unpatched, sheet.tsx's ported side classes (e.g.
//     right: `inset-y-0 right-0 w-3/4`) leave every UA value that they
//     don't explicitly override in force, including `margin:auto` — with
//     `left`(UA)/`width`(author)/`right`(author) then all non-auto, the
//     over-constrained abspos rule's "both margins auto -> center" branch
//     (CSS2.1 10.3.7) can center the sheet instead of anchoring it. `m-0`
//     forecloses that branch unconditionally; side-agnostic (every side
//     variant leaves the same kind of gap), so it sits in the base string.
//   - each side switch case below adds the ONE inset utility on the
//     opposite edge from its anchor, set to `auto` (`left-auto` for right,
//     `right-auto` for left, `bottom-auto` for top, `top-auto` for
//     bottom) — the load-bearing half of the fix, found only by rendering
//     in-browser (site/examples/sheet) per Task 3's verify requirement.
//     `m-0` alone was NOT sufficient: sheet.tsx's side classes set only
//     THREE of the four inset sides (e.g. right's `inset-y-0 right-0`
//     never touches `left`), so the UA's own base (non-`:modal`) `dialog {
//     inset: 0 }` rule's `left: 0` survives untouched — with `left`(UA-0),
//     `width`(author, e.g. `w-3/4`), and `right`(author-0) then ALL
//     concretely specified (none literally the `auto` keyword) and margins
//     pinned to 0 by the fix above, this is CSS2.1 10.3.7's genuinely
//     over-constrained case with no auto anywhere: the spec text says the
//     `left` value should be the one ignored/recomputed for `direction:
//     ltr` (i.e., `right`+`width` should win) — but empirically, in
//     Chrome, `left` won instead: the sheet rendered flush against the
//     WRONG edge (a right-side sheet opened flush-left) until `left` was
//     made genuinely `auto` via an explicit utility, not merely absent
//     from the class string. Once `left` is truly `auto` (not just
//     UA-defaulted-to-0), the box falls through to the unambiguous
//     "`left` auto, `width`+`right` given" case and anchors correctly — the
//     same reasoning applies on the block axis to the top/bottom variants,
//     where the unpatched UA `bottom`/`top` values would otherwise stretch
//     an `h-auto` box between two pinned edges instead of sizing it to
//     content.
//   - `max-h-none` is added to the base string: Chrome's UA stylesheet for
//     ANY <dialog> (`:modal` included, verified via getComputedStyle in
//     the same in-browser session) also carries `max-width`/`max-height:
//     calc(100% - 6px - 2em)` — a shrink-to-fit-viewport safety net for the
//     UA's own default centered/`fit-content`-sized dialog, never
//     overridden by sheet.tsx's ported classes because plain Radix `<div>`s
//     have no such default to fight. It legitimately clamps the LEFT/RIGHT
//     sides' `h-full` (an explicit, non-`auto` `height:100%`) short of the
//     viewport by that same ~38px unless neutralized — confirmed
//     empirically (measured rect height 38px short of `window.innerHeight`
//     before the fix, flush after). `max-width` needs no base-string
//     neutralizing: LEFT/RIGHT already carry their own explicit
//     `sm:max-w-sm`, and author-origin ALWAYS beats UA-origin regardless of
//     which is numerically smaller, so the UA `max-width` is already inert
//     there — see the next ADAPT for why TOP/BOTTOM still need their own.
//   - TOP/BOTTOM add `w-full max-w-none` (not in sheet.tsx's own source,
//     which relies on bare `inset-x-0` with `width` left `auto` to stretch
//     the box edge-to-edge — the plain-`<div>` default for an
//     absolutely-positioned box with both inline insets set). On the
//     native <dialog>, that same `inset-x-0`-with-`auto`-width did NOT
//     stretch to fill in Chrome — measured rendered width came out
//     shrink-to-fit (content width, e.g. ~447px on a 1888px-wide test
//     viewport) instead of the expected full-bleed width, a real
//     browser/engine discrepancy from the classical CSS2.1 stretch
//     algorithm this port does not attempt to fully explain, only to work
//     around: an explicit `w-full` (rather than relying on `auto`
//     resolution) sidesteps the ambiguity entirely. Because `width:100%`
//     is then itself subject to the SAME UA `max-width: calc(100% - 6px -
//     2em)` the `max-h-none` ADAPT above neutralizes on the block axis,
//     TOP/BOTTOM also need their own explicit `max-w-none` (LEFT/RIGHT
//     don't, per the previous ADAPT) to actually reach full width rather
//     than stopping ~38px short horizontally too.
//
// hideCloseButton mirrors DialogContent's own convention (zero value keeps
// shadcn's showCloseButton default of true) — the task brief's own
// signature note only calls out `side`, but sheet.tsx's showCloseButton is
// a real, pre-existing prop (not something newly grown since the brief was
// written), and DialogContent already established the bool-inversion
// pattern for exactly this prop on the sibling component the brief points
// at ("compare dialog.gsx's injected button") — see docs/jsx-parity.md
// `## sheet` NOTE for this judgment call.
component SheetContent(side string, hideCloseButton bool, children gsx.Node, attrs gsx.Attrs) {
	<dialog
		data-slot="sheet-content"
		data-gsxui-dialog-content
		data-state="closed"
		data-side={side |> default("right")}
		class={
			"fixed z-50 m-0 max-h-none open:flex flex-col gap-4 bg-background text-foreground shadow-lg transition ease-in-out data-[state=closed]:animate-out data-[state=closed]:duration-300 data-[state=open]:animate-in data-[state=open]:duration-500 backdrop:bg-black/10 supports-backdrop-filter:backdrop:backdrop-blur-xs data-[state=open]:backdrop:animate-in data-[state=open]:backdrop:fade-in-0 data-[state=open]:backdrop:duration-500 data-[state=closed]:backdrop:animate-out data-[state=closed]:backdrop:fade-out-0 data-[state=closed]:backdrop:duration-300",
			switch side {
			case "left":
				"inset-y-0 left-0 right-auto h-full w-3/4 border-r data-[state=closed]:slide-out-to-left data-[state=open]:slide-in-from-left sm:max-w-sm"
			case "top":
				"inset-x-0 top-0 bottom-auto w-full max-w-none h-auto border-b data-[state=closed]:slide-out-to-top data-[state=open]:slide-in-from-top"
			case "bottom":
				"inset-x-0 bottom-0 top-auto w-full max-w-none h-auto border-t data-[state=closed]:slide-out-to-bottom data-[state=open]:slide-in-from-bottom"
			default:
				"inset-y-0 right-0 left-auto h-full w-3/4 border-l data-[state=closed]:slide-out-to-right data-[state=open]:slide-in-from-right sm:max-w-sm"
			}
		}
		{ attrs... }
	>
		{ children }
		{ if !hideCloseButton {
			<button
				type="button"
				data-slot="sheet-close"
				data-gsxui-dialog-close
				aria-label="Close"
				class="absolute top-4 right-4 rounded-xs opacity-70 ring-offset-background transition-opacity hover:opacity-100 focus:ring-2 focus:ring-ring focus:ring-offset-2 focus:outline-hidden disabled:pointer-events-none data-[state=open]:bg-secondary"
			>
				<svg
					class="size-4"
					xmlns="http://www.w3.org/2000/svg"
					width="24"
					height="24"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
				>
					<path d="M18 6 6 18"/>
					<path d="m6 6 12 12"/>
				</svg>
			</button>
		} }
	</dialog>
}

component SheetHeader(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="sheet-header" class="flex flex-col gap-1.5 p-4" { attrs... }>{ children }</div>
}

component SheetFooter(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="sheet-footer" class="mt-auto flex flex-col gap-2 p-4" { attrs... }>{ children }</div>
}

component SheetTitle(children gsx.Node, attrs gsx.Attrs) {
	<h2 data-slot="sheet-title" class="font-semibold text-foreground" { attrs... }>{ children }</h2>
}

component SheetDescription(children gsx.Node, attrs gsx.Attrs) {
	<p data-slot="sheet-description" class="text-sm text-muted-foreground" { attrs... }>{ children }</p>
}

component SheetClose(children gsx.Node, attrs gsx.Attrs) {
	<button data-slot="sheet-close" data-gsxui-dialog-close type="button" { attrs... }>{ children }</button>
}
