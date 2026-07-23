// Package avatar holds the site's example gsx components for ui/avatar.
package avatar

import "github.com/gsxhq/gsxui/ui"

// data: image URLs must be ;base64, — gsx's image-sink sanitizer blocks percent-encoded forms (see docs/jsx-parity.md ## avatar).
const avatarSVG = "data:image/svg+xml;base64,PHN2ZyB4bWxucz0naHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmcnIHdpZHRoPSc2NCcgaGVpZ2h0PSc2NCc+PHJlY3Qgd2lkdGg9JzY0JyBoZWlnaHQ9JzY0JyBmaWxsPScjNmQyOGQ5Jy8+PC9zdmc+"

// Basic renders two Avatars side by side: one whose image loads (a data
// URI), one whose image 404s and falls back to initials — avatar.js
// toggles which of image/fallback is visible on the image's load/error.
component Basic() {
	<div class="flex items-center gap-4">
		<ui.Avatar>
			<ui.AvatarImage src={avatarSVG} alt="Ada Lovelace"/>
			<ui.AvatarFallback>AL</ui.AvatarFallback>
		</ui.Avatar>
		<ui.Avatar>
			<ui.AvatarImage src="/broken-image.jpg" alt="Broken avatar"/>
			<ui.AvatarFallback>JD</ui.AvatarFallback>
		</ui.Avatar>
	</div>
}
