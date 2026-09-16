package pagination

import "github.com/gsxhq/gsxui/site/uirtl"

// Rtl mirrors Basic's prev/1/2/3/ellipsis/next trail wrapped in dir="rtl".
// Page numbers keep Latin digits — shadcn's own convention — while the
// prev/next chevrons flip via rtl:rotate-180 and the trail itself reads
// right-to-left.
component Rtl() {
	<div dir="rtl">
		<uirtl.Pagination>
			<uirtl.PaginationContent>
				<uirtl.PaginationItem>
					<uirtl.PaginationPrevious href="#"/>
				</uirtl.PaginationItem>
				<uirtl.PaginationItem>
					<uirtl.PaginationLink href="#">1</uirtl.PaginationLink>
				</uirtl.PaginationItem>
				<uirtl.PaginationItem>
					<uirtl.PaginationLink href="#" isActive>2</uirtl.PaginationLink>
				</uirtl.PaginationItem>
				<uirtl.PaginationItem>
					<uirtl.PaginationLink href="#">3</uirtl.PaginationLink>
				</uirtl.PaginationItem>
				<uirtl.PaginationItem>
					<uirtl.PaginationEllipsis/>
				</uirtl.PaginationItem>
				<uirtl.PaginationItem>
					<uirtl.PaginationNext href="#"/>
				</uirtl.PaginationItem>
			</uirtl.PaginationContent>
		</uirtl.Pagination>
	</div>
}
