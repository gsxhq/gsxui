package ui

import "github.com/gsxhq/gsx"

// Checkbox is the shadcn/ui Checkbox, ported as a real native <input
// type="checkbox">: form-native, zero JS, :checked/:disabled browser truth
// replaces Radix's button-role + hidden-input + Indicator part (ledger
// ADAPT). The Indicator's CheckIcon child becomes a checked:bg-[url(...)]
// data-URI background — a void, childless element opts into fallthrough via
// the explicit { attrs... } spread, same as every other component here.
//
// NOTE: the data-URI's embedded spaces (SVG attribute/path separators) are
// authored as Tailwind's underscore escape for whitespace inside a bracketed
// arbitrary value (`_` where the rendered CSS needs ` `) — not literal
// spaces. A literal space would split the token: HTML's class attribute is
// itself whitespace-delimited, so any tool walking the class list (the
// configured tailwind-merge included) treats a bare space as a token
// boundary, silently truncating/reordering the SVG markup. See
// docs/jsx-parity.md.
component Checkbox(attrs gsx.Attrs) {
	<input type="checkbox" data-slot="checkbox" class="peer size-4 shrink-0 appearance-none rounded-[4px] border border-input shadow-xs transition-shadow outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 dark:bg-input/30 checked:bg-primary checked:border-primary checked:bg-[url('data:image/svg+xml;charset=utf-8,%3Csvg_xmlns=%22http://www.w3.org/2000/svg%22_viewBox=%220_0_24_24%22_fill=%22none%22_stroke=%22white%22_stroke-width=%223%22_stroke-linecap=%22round%22_stroke-linejoin=%22round%22%3E%3Cpath_d=%22M20_6_9_17l-5-5%22/%3E%3C/svg%3E')] checked:bg-center checked:bg-no-repeat checked:bg-[length:12px_12px]" { attrs... }/>
}
