package skeleton

import "github.com/gsxhq/gsxui/ui"

// Card composes Skeletons into a card-shaped loading placeholder — the
// space a real Card would occupy while its data is still loading.
component Card() {
	<div class="flex flex-col gap-3">
		<ui.Skeleton class="h-[125px] w-[250px] rounded-xl"/>
		<div class="grid gap-2">
			<ui.Skeleton class="h-4 w-[250px]"/>
			<ui.Skeleton class="h-4 w-[200px]"/>
		</div>
	</div>
}
