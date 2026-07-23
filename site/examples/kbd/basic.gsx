// Package kbd holds the site's example gsx components for ui/kbd.
package kbd

import "github.com/gsxhq/gsxui/ui"

// Basic renders a KbdGroup showing a compound shortcut alongside a single
// standalone Kbd.
component Basic() {
	<div class="flex items-center gap-4">
		<ui.KbdGroup>
			<ui.Kbd>Ctrl</ui.Kbd>
			<ui.Kbd>B</ui.Kbd>
		</ui.KbdGroup>
		<ui.Kbd>Esc</ui.Kbd>
	</div>
}
