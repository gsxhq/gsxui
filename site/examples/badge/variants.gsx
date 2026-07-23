package badge

import "github.com/gsxhq/gsxui/ui"

// Variants renders every Badge variant side by side.
component Variants() {
	<div class="flex flex-wrap items-center gap-2">
		<ui.Badge>Default</ui.Badge>
		<ui.Badge variant="secondary">Secondary</ui.Badge>
		<ui.Badge variant="destructive">Destructive</ui.Badge>
		<ui.Badge variant="outline">Outline</ui.Badge>
		<ui.Badge variant="ghost">Ghost</ui.Badge>
		<ui.Badge variant="link">Link</ui.Badge>
	</div>
}
