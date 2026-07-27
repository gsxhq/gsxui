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
		<ui.Button size="lg">Large</ui.Button>
		<ui.Button size="icon-sm" aria-label="Small icon"><span>+</span></ui.Button>
		<ui.Button size="icon-lg" aria-label="Large icon"><span>+</span></ui.Button>
		<ui.Button aria-invalid="true">Invalid</ui.Button>
	</div>
}
