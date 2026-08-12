package ui

import "github.com/gsxhq/gsx"

// Avatar and its parts are the shadcn/ui Avatar. Radix's client-side
// load-state machinery (image loading/loaded/error, driving which of
// image/fallback is mounted) is replaced by delegation (ADAPT, see
// docs/jsx-parity.md): AvatarFallback always renders — no hidden attribute,
// since load state isn't known at render time — and AvatarImage carries the
// data-gsxui-slot-avatar-image hook; ui/avatar/avatar.js toggles display on the
// image's native load/error events. Requires the avatar behavior module
// (ui/avatar/avatar.js).

component Avatar(children gsx.Node, attrs gsx.Attrs) {
	<span
		class={
			"size-8 rounded-full after:rounded-full data-[size=lg]:size-10 data-[size=sm]:size-6 flex relative shrink-0 select-none after:absolute after:inset-0 after:border after:border-border after:mix-blend-darken dark:after:mix-blend-lighten"
		}
		{ attrs... }
		data-gsxui-slot-avatar
	>
		{ children }
	</span>
}

component AvatarImage(src string, alt string, attrs gsx.Attrs) {
	<img
		src={src}
		alt={alt}
		class={ "rounded-full size-full object-cover aspect-square" }
		{ attrs... }
		data-gsxui-slot-avatar-image
	/>
}

// FIX: AvatarFallback's own recipe accessor only ever resolves to colour
// (bg-muted/text-muted-foreground) plus rounded-full — never size or
// centering, the same structural/presentational split progress.gsx's own
// FIX comment documents. Without "size-full items-center justify-center"
// (upstream's literal, non-recipe classes for this element), the fallback
// span is a flex item with no explicit size: it shrinks to fit its text
// content instead of filling the Avatar circle, and that text is never
// centered within it — invisible at the default size-8 avatar (close
// enough to look right by coincidence) but glaring at any larger override
// (found reviewing theme-creator-parity Task 4's profile card, which sets
// size-16).
component AvatarFallback(children gsx.Node, attrs gsx.Attrs) {
	<span
		class={ "size-full items-center justify-center", "bg-muted text-muted-foreground rounded-full flex text-sm" }
		{ attrs... }
		data-gsxui-slot-avatar-fallback
	>
		{ children }
	</span>
}
