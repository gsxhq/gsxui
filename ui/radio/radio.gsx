package radio

import "github.com/gsxhq/gsx"

// Radio is the shadcn/ui RadioGroupItem, ported as a real native
// <input type="radio">: form-native, zero JS, browser :checked/:disabled
// truth replaces Radix's button-role + hidden-input + Indicator part
// (ledger ADAPT). shadcn's RadioGroup container (a styled grid wrapper) is
// not ported — native `name` grouping on sibling <input type="radio">
// elements already gives you the group; the layout wrapper is the caller's
// concern, same as any other flex/grid container (ledger ADAPT).
//
// The Indicator's CircleIcon child becomes a checked:bg-[url(...)] data-URI
// background (same restructuring as checkbox); the data-URI's embedded
// spaces are Tailwind's underscore escape for whitespace inside a bracketed
// arbitrary value, not literal spaces — see checkbox.gsx and
// docs/jsx-parity.md.
component Radio(attrs gsx.Attrs) {
	<input type="radio" data-slot="radio" class="peer size-4 shrink-0 appearance-none rounded-full border border-input shadow-xs transition-shadow outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 dark:bg-input/30 checked:bg-primary checked:border-primary checked:bg-[url('data:image/svg+xml;charset=utf-8,%3Csvg_xmlns=%22http://www.w3.org/2000/svg%22_viewBox=%220_0_24_24%22%3E%3Ccircle_cx=%2212%22_cy=%2212%22_r=%226%22_fill=%22white%22/%3E%3C/svg%3E')] checked:bg-center checked:bg-no-repeat checked:bg-[length:12px_12px]" { attrs... }/>
}
