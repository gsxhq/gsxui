// Package pagination holds the site's example gsx components for
// ui/pagination.
package pagination

import "github.com/gsxhq/gsxui/ui"

// Basic renders a realistic prev/1/2/3/ellipsis/next trail with page 2
// active — the same shape as shadcn's own pagination-demo.tsx.
component Basic() {
	<ui.Pagination>
		<ui.PaginationContent>
			<ui.PaginationItem>
				<ui.PaginationPrevious href="#"/>
			</ui.PaginationItem>
			<ui.PaginationItem>
				<ui.PaginationLink href="#">1</ui.PaginationLink>
			</ui.PaginationItem>
			<ui.PaginationItem>
				<ui.PaginationLink href="#" isActive>2</ui.PaginationLink>
			</ui.PaginationItem>
			<ui.PaginationItem>
				<ui.PaginationLink href="#">3</ui.PaginationLink>
			</ui.PaginationItem>
			<ui.PaginationItem>
				<ui.PaginationEllipsis/>
			</ui.PaginationItem>
			<ui.PaginationItem>
				<ui.PaginationNext href="#"/>
			</ui.PaginationItem>
		</ui.PaginationContent>
	</ui.Pagination>
}
