package sera

import "github.com/gsxhq/gsx"

// Checkbox is the shadcn/ui Checkbox, ported as a real native <input
// type="checkbox">: form-native, zero JS, :checked/:disabled browser truth
// replaces Radix's button-role + hidden-input + Indicator part (ledger
// ADAPT). The Indicator's CheckIcon child becomes a checked-state data-URI
// background in the default stylesheet.
//
// NOTE: the data-URI payloads are base64 — deliberately, after every
// richer encoding lost to some layer of the CSS toolchain. Base64 is
// [A-Za-z0-9+/=] only, so neither PostCSS nor the browser reparses payload
// punctuation. TestCheckboxDataURIDecodesToValidSVG reads the stylesheet,
// decodes both URIs, and XML-parses them.
//
// The explicit dark checked rule mirrors shadcn's dark checked state and
// follows the dark unchecked paint in source order. The dark check URI
// strokes the dark --primary-foreground
// value (oklch(0.205 0 0)) because primary flips near-white there and the
// light URI's white stroke would vanish; both strokes are static text — a
// data-URI can't read CSS variables — so custom themes that move
// primary-foreground still need the ledgered currentColor-mask follow-up.
component Checkbox(attrs gsx.Attrs) {
	<input
		type="checkbox"
		class={
			"border-input bg-transparent checked:bg-primary checked:text-primary-foreground dark:checked:bg-primary checked:border-primary aria-invalid:checked:border-primary aria-invalid:border-destructive dark:aria-invalid:border-destructive/50 focus-visible:border-ring focus-visible:ring-ring/30 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 flex size-4.5 items-center justify-center rounded-none border transition-shadow group-has-disabled/field:opacity-50 focus-visible:ring-2 aria-invalid:ring-2 shrink-0 outline-none disabled:cursor-not-allowed disabled:opacity-50"
		}
		{ attrs... }
		data-gsxui-slot-checkbox
	/>
}
