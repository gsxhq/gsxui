package skeleton

import uiskeleton "github.com/gsxhq/gsxui/ui/skeleton"

// Card composes Skeletons into a card-shaped loading placeholder — the
// space a real Card would occupy while its data is still loading.
component Card() {
	<div class="flex flex-col gap-3">
		<uiskeleton.Skeleton class="h-[125px] w-[250px] rounded-xl"/>
		<div class="grid gap-2">
			<uiskeleton.Skeleton class="h-4 w-[250px]"/>
			<uiskeleton.Skeleton class="h-4 w-[200px]"/>
		</div>
	</div>
}
