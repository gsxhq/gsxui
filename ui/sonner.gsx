package ui

import "github.com/gsxhq/gsx"

// Toaster is the sole server-rendered surface of the sonner port. shadcn's
// own sonner.tsx renders nothing but a re-themed <Sonner> passthrough — the
// toast library owns 100% of the toast DOM and ships it from a non-Tailwind
// stylesheet — so there is no toast-card markup to port. gsxui reconstructs
// that look as plain Tailwind classes built entirely by ui/sonner.js: this
// component ships only the always-present, positioned region the JS appends
// each toast <li> into.
//
// Mount ONCE per page (typically the root layout, same convention as
// shadcn's <Toaster/> in app/layout.tsx). v1 ships only the default
// bottom-right position — the other five sonner positions are a ledgered
// gap (docs/jsx-parity.md ## sonner).
//
// The <section> is the aria landmark ("Notifications"); the <ol> is
// sonner.js's mount point (data-gsxui-toaster) and carries the fixed
// bottom-right stacking region. pointer-events-none lets clicks fall
// through the empty gutter; each constructed toast re-enables
// pointer-events on itself.
component Toaster(attrs gsx.Attrs) {
	<section aria-label="Notifications" tabindex="-1">
		<ol
			data-slot="toaster"
			data-gsxui-toaster
			class="pointer-events-none fixed z-100 flex flex-col gap-2 p-6 bottom-0 right-0"
			{ attrs... }
		></ol>
	</section>
}
