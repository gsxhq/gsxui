// Package avatar holds the site's example gsx components for ui/avatar.
package avatar

import uiavatar "github.com/gsxhq/gsxui/ui/avatar"

// data: image URLs must be ;base64, — gsx's image-sink sanitizer blocks percent-encoded forms (see docs/jsx-parity.md ## avatar).
const avatarSVG = "data:image/svg+xml;base64,PHN2ZyB4bWxucz0naHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmcnIHdpZHRoPSc2NCcgaGVpZ2h0PSc2NCc+PHJlY3Qgd2lkdGg9JzY0JyBoZWlnaHQ9JzY0JyBmaWxsPScjNmQyOGQ5Jy8+PC9zdmc+"

// Basic renders two Avatars side by side: one whose image loads (a data
// URI), one whose image 404s and falls back to initials — avatar.js
// toggles which of image/fallback is visible on the image's load/error.
component Basic() {
	<div class="flex items-center gap-4">
		<uiavatar.Avatar>
			<uiavatar.AvatarImage src={avatarSVG} alt="Ada Lovelace"/>
			<uiavatar.AvatarFallback>AL</uiavatar.AvatarFallback>
		</uiavatar.Avatar>
		<uiavatar.Avatar>
			<uiavatar.AvatarImage src="/broken-image.jpg" alt="Broken avatar"/>
			<uiavatar.AvatarFallback>JD</uiavatar.AvatarFallback>
		</uiavatar.Avatar>
	</div>
}
