package skeleton

import "github.com/gsxhq/gsx"

// Skeleton is the shadcn/ui Skeleton loading placeholder. Straight port; no
// divergences.
component Skeleton(attrs gsx.Attrs) {
	<div data-slot="skeleton" class="animate-pulse rounded-md bg-accent" { attrs... }/>
}
