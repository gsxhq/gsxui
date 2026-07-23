// Package skeleton holds the site's example gsx components for
// ui/skeleton.
package skeleton

import uiskeleton "github.com/gsxhq/gsxui/ui/skeleton"

// Basic renders a circular avatar placeholder beside two text-line
// placeholders — the shape, not the content, is what a skeleton loads.
component Basic() {
	<div class="flex items-center gap-4">
		<uiskeleton.Skeleton class="size-12 rounded-full"/>
		<div class="grid gap-2">
			<uiskeleton.Skeleton class="h-4 w-[250px]"/>
			<uiskeleton.Skeleton class="h-4 w-[200px]"/>
		</div>
	</div>
}
