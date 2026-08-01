package pagination

import "github.com/gsxhq/gsxui/ui"

// Rtl mirrors Basic's prev/1/2/3/ellipsis/next trail wrapped in dir="rtl".
// Page numbers keep Latin digits — shadcn's own convention — while the
// prev/next chevrons flip via rtl:rotate-180 and the trail itself reads
// right-to-left.
component Rtl() {
	<div dir="rtl">
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
	</div>
}
