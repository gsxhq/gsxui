// switch.gsx backs the Switch component. "switch" is a Go keyword, so this
// could never be its own package — one of the reasons ui/ is a single flat
// package (see docs/jsx-parity.md packaging entry); as a file basename and
// CLI name (`gsxui add switch`) it is fine. The component is `ui.Switch`.
package ui

import "github.com/gsxhq/gsx"

// Switch is the shadcn/ui Switch, ported as a real native
// <input type="checkbox" role="switch">: form-native, zero JS, browser
// :checked/:disabled truth replaces Radix's button-role Root + separate
// Thumb span (ledger ADAPT). role="switch" preserves the switch semantics
// a plain checkbox input doesn't carry on its own (ARIA maps
// input[type=checkbox][role=switch] correctly; the checked state itself is
// still native).
//
// The Thumb is no longer a sibling element to target via Radix's
// group-data-[size]/switch selector — it is this element's own ::before
// pseudo-element (MECHANISM: thumb-span→before:). A generated pseudo-
// element renders nothing without an explicit content utility, unlike a
// real child element, so before:content-[''] is required (not present on
// the Radix Thumb span, which needs no content at all) — see
// docs/jsx-parity.md.
component Switch(attrs gsx.Attrs) {
	<input
		type="checkbox"
		role="switch"
		data-slot="switch"
		class="peer inline-flex shrink-0 items-center appearance-none rounded-full border border-transparent transition-all outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 h-[1.15rem] w-8 bg-input checked:bg-primary dark:bg-input/80 dark:checked:bg-primary before:pointer-events-none before:block before:size-4 before:rounded-full before:bg-background before:transition-transform before:content-[''] checked:before:translate-x-[calc(100%-2px)] dark:before:bg-foreground dark:checked:before:bg-primary-foreground"
		{ attrs... }
	/>
}
