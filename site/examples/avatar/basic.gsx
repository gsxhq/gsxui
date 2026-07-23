// Package avatar holds the site's example gsx components for ui/avatar.
package avatar

import uiavatar "github.com/gsxhq/gsxui/ui/avatar"

// avatarSVG is a small inline data-URI image — no network fetch, so the
// "loaded" half of the pair renders identically in every environment. gsx's
// image-resource sanitizer only accepts base64-encoded data: URLs (the
// ";base64," marker constrains the payload charset so it can't smuggle a
// scheme break), so the SVG is base64, not percent-encoded.
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
