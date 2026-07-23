package separator

import "github.com/gsxhq/gsxui/ui"

// Orientation combines both orientations: a horizontal Separator under a
// heading block, then vertical Separators between an inline row of links.
component Orientation() {
	<div>
		<div class="space-y-1">
			<h4 class="text-sm font-medium">gsxui</h4>
			<p class="text-sm text-muted-foreground">Real gsx components, not JSX facsimiles.</p>
		</div>
		<ui.Separator class="my-4"/>
		<div class="flex h-5 items-center gap-4 text-sm">
			<div>Docs</div>
			<ui.Separator orientation="vertical"/>
			<div>Components</div>
			<ui.Separator orientation="vertical"/>
			<div>GitHub</div>
		</div>
	</div>
}
