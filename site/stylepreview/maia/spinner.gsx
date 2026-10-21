package maia

import (
	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui/i18n"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Spinner is the shadcn/ui Spinner. shadcn renders lucide-react's
// Loader2Icon directly (registry/new-york-v4/ui/spinner.tsx); ui/icon's
// generated set carries the same Lucide glyph as icon.LoaderCircle (Lucide
// icon name "loader-circle" — Loader2Icon is lucide-react's alias for it),
// so this port is one icon.LoaderCircle call instead of a bespoke <svg>.
// aria-hidden="false" is required, not optional: svgIcon (ui/icon/icon.gsx)
// defaults every icon to aria-hidden="true" unless the caller supplies its
// own aria-hidden, and shadcn's Spinner is deliberately NOT hidden from
// assistive tech — role="status" and the "Loading" aria-label announce it.
// All four are authored before { attrs... }, so a caller can still override
// any of them (the same overridable-default idiom as button's
// type="button").
component Spinner(attrs gsx.Attrs) {
	<icon.LoaderCircle
		role="status"
		aria-label={i18n.T("Loading")}
		aria-hidden="false"
		class={ "size-4 animate-spin" }
		{ attrs... }
		data-gsxui-slot-spinner
	/>
}
