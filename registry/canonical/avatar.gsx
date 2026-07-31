package canonical

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
	<span class={ avatar.Root() } { attrs... } data-gsxui-slot-avatar>
		{ children }
	</span>
}

component AvatarImage(src string, alt string, attrs gsx.Attrs) {
	<img
		src={src}
		alt={alt}
		class={ avatar.Image() }
		{ attrs... }
		data-gsxui-slot-avatar-image
	/>
}

component AvatarFallback(children gsx.Node, attrs gsx.Attrs) {
	<span
		class={ avatar.Fallback() }
		{ attrs... }
		data-gsxui-slot-avatar-fallback
	>
		{ children }
	</span>
}
