// switch.gsx backs the Switch component. "switch" is a Go keyword, so this
// could never be its own package — one of the reasons ui/ is a single flat
// package (see docs/jsx-parity.md packaging entry); as a file basename and
// CLI name (`gsxui add switch`) it is fine. The component is `ui.Switch`.
package sera

import "github.com/gsxhq/gsx"

// Switch is the shadcn/ui Switch, ported as a real native
// <input type="checkbox" role="switch">: form-native, zero JS, browser
// :checked/:disabled truth replaces Radix's button-role Root + separate
// Thumb span (ledger ADAPT). role="switch" preserves the switch semantics
// a plain checkbox input doesn't carry on its own (ARIA maps
// input[type=checkbox][role=switch] correctly; the checked state itself is
// still native).
//
// The Thumb is no longer a sibling element targeted through Radix's
// ancestor size state — it is this element's own ::before
// pseudo-element (MECHANISM: thumb-span→before:). A generated pseudo-
// element renders nothing without an explicit content utility, unlike a
// real child element, so before:content-[''] is required (not present on
// the Radix Thumb span, which needs no content at all) — see
// docs/jsx-parity.md.
component Switch(attrs gsx.Attrs) {
	<input
		type="checkbox"
		role="switch"
		class={
			"data-checked:bg-primary data-unchecked:bg-input data-checked:border-primary border data-unchecked:border-input/50 focus-visible:border-ring focus-visible:ring-ring/30 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive dark:aria-invalid:border-destructive/50 shrink-0 rounded-none focus-visible:ring-2 aria-invalid:ring-2 data-[size=default]:h-4.5 data-[size=default]:w-8.25 data-[size=sm]:h-3.5 data-[size=sm]:w-6.25 before:bg-background dark:data-unchecked:before:bg-foreground dark:data-checked:before:bg-primary-foreground group-data-[size=default]/switch:before:size-3.5 group-data-[size=sm]/switch:before:size-2.5 data-checked:before:translate-x-[calc(100%+2px)] data-unchecked:before:translate-x-0.25 inline-flex"
		}
		{ attrs... }
		data-gsxui-slot-switch
	/>
}
