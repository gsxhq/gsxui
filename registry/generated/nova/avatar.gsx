package ui

import "github.com/gsxhq/gsx"

// Avatar and its parts are the shadcn/ui Avatar. Radix's client-side
// load-state machinery (image loading/loaded/error, driving which of
// image/fallback is mounted) is replaced by delegation (ADAPT, see
// docs/jsx-parity.md): AvatarFallback always renders — no hidden attribute,
// since load state isn't known at render time — and AvatarImage carries the
// data-gsxui-avatar-image hook; ui/avatar/avatar.js toggles display on the
// image's native load/error events. Requires the avatar behavior module
// (ui/avatar/avatar.js).

component Avatar(children gsx.Node, attrs gsx.Attrs) {
	<span
		data-gsxui-avatar
		class={ "relative flex size-8 shrink-0 overflow-hidden rounded-full select-none" }
		{ attrs... }
		data-gsxui-slot-avatar
	>
		{ children }
	</span>
}

component AvatarImage(src string, alt string, attrs gsx.Attrs) {
	<img
		data-gsxui-avatar-image
		src={src}
		alt={alt}
		class={ "absolute inset-0 aspect-square size-full" }
		{ attrs... }
		data-gsxui-slot-avatar-image
	/>
}

component AvatarFallback(children gsx.Node, attrs gsx.Attrs) {
	<span
		data-gsxui-avatar-fallback
		class={ "flex size-full items-center justify-center rounded-full bg-muted text-sm text-muted-foreground" }
		{ attrs... }
		data-gsxui-slot-avatar-fallback
	>
		{ children }
	</span>
}
