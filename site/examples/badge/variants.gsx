package badge

import uibadge "github.com/gsxhq/gsxui/ui/badge"

// Variants renders every Badge variant side by side.
component Variants() {
	<div class="flex flex-wrap items-center gap-2">
		<uibadge.Badge>Default</uibadge.Badge>
		<uibadge.Badge variant="secondary">Secondary</uibadge.Badge>
		<uibadge.Badge variant="destructive">Destructive</uibadge.Badge>
		<uibadge.Badge variant="outline">Outline</uibadge.Badge>
		<uibadge.Badge variant="ghost">Ghost</uibadge.Badge>
		<uibadge.Badge variant="link">Link</uibadge.Badge>
	</div>
}
