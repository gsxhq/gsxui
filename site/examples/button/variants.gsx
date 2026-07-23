package button

import "github.com/gsxhq/gsxui/ui"

// Variants renders every Button variant side by side.
component Variants() {
	<div class="flex flex-wrap items-center gap-3">
		<ui.Button>Default</ui.Button>
		<ui.Button variant="secondary">Secondary</ui.Button>
		<ui.Button variant="destructive">Destructive</ui.Button>
		<ui.Button variant="outline">Outline</ui.Button>
		<ui.Button variant="ghost">Ghost</ui.Button>
		<ui.Button variant="link">Link</ui.Button>
	</div>
}
