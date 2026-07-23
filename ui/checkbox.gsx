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
// percent-encoded (%20) — never literal spaces, and never Tailwind's
// underscore escape. A literal space would split the token (HTML's class
// attribute is whitespace-delimited, so any tool walking the class list —
// the configured tailwind-merge included — treats a bare space as a token
// boundary, tearing the SVG apart). And `_` does NOT become a space here:
// Tailwind v4 deliberately preserves underscores inside url() values, so an
// underscored URI reaches the browser verbatim as invalid XML
// (<svg_xmlns=...>) and the checkmark silently never paints. %20 satisfies
// both layers: one unbroken class token, standard percent-decoding in the
// browser. Pinned by TestCheckboxDataURIDecodesToValidSVG, which decodes
// the rendered URI and XML-parses it. See docs/jsx-parity.md.
component Checkbox(attrs gsx.Attrs) {
	<input type="checkbox" data-slot="checkbox" class="peer size-4 shrink-0 appearance-none rounded-[4px] border border-input shadow-xs transition-shadow outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 dark:bg-input/30 checked:bg-primary checked:border-primary checked:bg-[url('data:image/svg+xml;charset=utf-8,%3Csvg%20xmlns=%22http://www.w3.org/2000/svg%22%20viewBox=%220%200%2024%2024%22%20fill=%22none%22%20stroke=%22white%22%20stroke-width=%223%22%20stroke-linecap=%22round%22%20stroke-linejoin=%22round%22%3E%3Cpath%20d=%22M20%206%209%2017l-5-5%22/%3E%3C/svg%3E')] checked:bg-center checked:bg-no-repeat checked:bg-[length:12px_12px]" { attrs... }/>
}
