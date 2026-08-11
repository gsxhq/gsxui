package lyra

import "github.com/gsxhq/gsx"

// Skeleton is the shadcn/ui Skeleton loading placeholder. Straight port; no
// divergences.
component Skeleton(attrs gsx.Attrs) {
	<div class={ "bg-muted rounded-none" } { attrs... } data-gsxui-slot-skeleton></div>
}
