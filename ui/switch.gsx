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
		{ attrs... }
		data-gsxui-slot-switch
	/>
}
