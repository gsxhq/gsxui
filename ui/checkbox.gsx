package ui

import "github.com/gsxhq/gsx"

// Checkbox is the shadcn/ui Checkbox, ported as a real native <input
// type="checkbox">: form-native, zero JS, :checked/:disabled browser truth
// replaces Radix's button-role + hidden-input + Indicator part (ledger
// ADAPT). The Indicator's CheckIcon child becomes a checked:bg-[url(...)]
// data-URI background — a void, childless element opts into fallthrough via
// the explicit { attrs... } spread, same as every other component here.
//
// NOTE: the data-URI payloads are base64 — deliberately, after every
// richer encoding lost to some layer of the toolchain in turn: literal
// spaces are class-token boundaries (torn by tailwind-merge); Tailwind's
// `_` whitespace escape is NOT converted inside url() values, so it
// reached the browser as invalid XML (<svg_xmlns=...>) and the checkmark
// silently never painted; and percent-encoding with parens (the dark
// stroke's oklch(...)) broke the postcss parse of Tailwind's emitted CSS
// under vite. Base64 is [A-Za-z0-9+/=] only — nothing for any layer to
// split, convert, or mis-parse. Pinned by
// TestCheckboxDataURIDecodesToValidSVG, which base64-decodes the rendered
// URIs and XML-parses them. See docs/jsx-parity.md.
//
// The dark:checked:* trio mirrors shadcn's explicit
// dark:data-[state=checked]:bg-primary — NOT redundant: the dark custom
// variant (:is(.dark *)) carries class specificity that beats a bare
// :checked, so without it dark:bg-input/30 wins over checked:bg-primary in
// dark mode. The dark check URI strokes the dark --primary-foreground
// value (oklch(0.205 0 0)) because primary flips near-white there and the
// light URI's white stroke would vanish; both strokes are static text — a
// data-URI can't read CSS variables — so custom themes that move
// primary-foreground still need the ledgered currentColor-mask follow-up.
component Checkbox(attrs gsx.Attrs) {
	<input
		type="checkbox"
		data-slot="checkbox"
		class="peer size-4 shrink-0 appearance-none rounded-[4px] border border-input transition-shadow outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 dark:bg-input/30 checked:bg-primary checked:border-primary checked:bg-[url('data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAyNCAyNCIgZmlsbD0ibm9uZSIgc3Ryb2tlPSJ3aGl0ZSIgc3Ryb2tlLXdpZHRoPSIzIiBzdHJva2UtbGluZWNhcD0icm91bmQiIHN0cm9rZS1saW5lam9pbj0icm91bmQiPjxwYXRoIGQ9Ik0yMCA2IDkgMTdsLTUtNSIvPjwvc3ZnPg==')] checked:bg-center checked:bg-no-repeat checked:bg-[length:14px_14px] dark:checked:bg-primary dark:checked:border-primary dark:checked:bg-[url('data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAyNCAyNCIgZmlsbD0ibm9uZSIgc3Ryb2tlPSJva2xjaCgwLjIwNSAwIDApIiBzdHJva2Utd2lkdGg9IjMiIHN0cm9rZS1saW5lY2FwPSJyb3VuZCIgc3Ryb2tlLWxpbmVqb2luPSJyb3VuZCI+PHBhdGggZD0iTTIwIDYgOSAxN2wtNS01Ii8+PC9zdmc+')]"
		{ attrs... }
	/>
}
