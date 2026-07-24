// Package toggle holds the site's example gsx components for ui/toggle.
package toggle

import (
	"github.com/gsxhq/gsxui/ui"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Basic rounds up shadcn's own scattered toggle-*.tsx examples (toggle-demo,
// -outline, -lg, -disabled) into one row: default variant with an icon +
// label (mirrors toggle-demo's Bookmark+"Bookmark" shape, swapped for a
// Bold text-formatting glyph — ui/icon has no Bookmark), an outline-variant
// icon-only toggle, sm/lg sizes, one server-rendered already-pressed
// (pressed={true}), and one disabled.
component Basic() {
	<div class="flex flex-wrap items-center gap-4">
		<ui.Toggle aria-label="Toggle bold">
			<icon.Bold/>
			Bold
		</ui.Toggle>
		<ui.Toggle variant="outline" aria-label="Toggle italic">
			<icon.Italic/>
		</ui.Toggle>
		<ui.Toggle size="sm" aria-label="Toggle underline (small)">
			<icon.Underline/>
		</ui.Toggle>
		<ui.Toggle size="lg" aria-label="Toggle underline (large)">
			<icon.Underline/>
		</ui.Toggle>
		<ui.Toggle pressed={true} aria-label="Toggle bold (pressed)">
			<icon.Bold/>
		</ui.Toggle>
		<ui.Toggle disabled aria-label="Toggle bold (disabled)">
			<icon.Bold/>
		</ui.Toggle>
	</div>
}
