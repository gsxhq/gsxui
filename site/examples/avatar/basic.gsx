// Package avatar holds the site's example gsx components for ui/avatar.
package avatar

import "github.com/gsxhq/gsxui/ui"

// A 64x64 violet tile with white "AL" initials — a stand-in portrait, so
// the loaded-image state is visibly an avatar and visibly distinct from
// the gray fallback. Authored as plain SVG bytes; the `|> dataURL(mime)`
// std filter at the src site assembles the base64 data: URL gsx's
// image-sink sanitizer accepts (see docs/jsx-parity.md ## avatar).
var avatarSVG = []byte("<svg xmlns='http://www.w3.org/2000/svg' width='64' height='64'><rect width='64' height='64' fill='#6d28d9'/><text x='32' y='34' text-anchor='middle' dominant-baseline='central' font-family='sans-serif' font-weight='600' font-size='26' fill='#fff'>AL</text></svg>")

// Basic renders two Avatars side by side: one whose image loads (a data
// URI), one whose image 404s and falls back to initials — avatar.js
// toggles which of image/fallback is visible on the image's load/error.
component Basic() {
	<div class="flex items-center gap-4">
		<ui.Avatar>
			<ui.AvatarImage src={avatarSVG |> dataURL("image/svg+xml")} alt="Ada Lovelace"/>
			<ui.AvatarFallback>AL</ui.AvatarFallback>
		</ui.Avatar>
		<ui.Avatar>
			<ui.AvatarImage src="/broken-image.jpg" alt="Broken avatar"/>
			<ui.AvatarFallback>JD</ui.AvatarFallback>
		</ui.Avatar>
	</div>
}
