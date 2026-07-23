package icon

import uiicon "github.com/gsxhq/gsxui/ui/icon"

// Sizes shows icons take the same class-merge story as every other
// component: size-* controls the box, text-* tints the stroke (stroke
// uses currentColor).
component Sizes() {
	<div class="flex items-center gap-4">
		<uiicon.Star class="size-4"/>
		<uiicon.Star class="size-6"/>
		<uiicon.Star class="size-8 text-primary"/>
		<uiicon.Star class="size-10 text-destructive"/>
	</div>
}
