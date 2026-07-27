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
	<span data-gsxui-avatar { withSlot("avatar", attrs)... }>
		{ children }
	</span>
}

component AvatarImage(src string, alt string, attrs gsx.Attrs) {
	<img
		data-gsxui-avatar-image
		src={src}
		alt={alt}
		{ withSlot("avatar-image", attrs)... }
	/>
}

component AvatarFallback(children gsx.Node, attrs gsx.Attrs) {
	<span
		data-gsxui-avatar-fallback
		{ withSlot("avatar-fallback", attrs)... }
	>
		{ children }
	</span>
}
