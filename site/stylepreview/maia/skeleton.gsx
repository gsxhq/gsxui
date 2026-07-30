package maia

import "github.com/gsxhq/gsx"

// Skeleton is the shadcn/ui Skeleton loading placeholder. Straight port; no
// divergences.
component Skeleton(attrs gsx.Attrs) {
	<div class={ "animate-pulse rounded-md bg-accent" } { attrs... } data-gsxui-slot-skeleton></div>
}
