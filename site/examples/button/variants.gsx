package button

import uibutton "github.com/gsxhq/gsxui/ui/button"

// Variants renders every Button variant side by side.
component Variants() {
	<div class="flex flex-wrap items-center gap-3">
		<uibutton.Button>Default</uibutton.Button>
		<uibutton.Button variant="secondary">Secondary</uibutton.Button>
		<uibutton.Button variant="destructive">Destructive</uibutton.Button>
		<uibutton.Button variant="outline">Outline</uibutton.Button>
		<uibutton.Button variant="ghost">Ghost</uibutton.Button>
		<uibutton.Button variant="link">Link</uibutton.Button>
	</div>
}
