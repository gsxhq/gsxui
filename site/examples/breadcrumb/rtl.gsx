package breadcrumb

import "github.com/gsxhq/gsxui/site/uirtl"

// Rtl mirrors Basic's trail shape — linked root, collapsed-middle
// ellipsis, linked mid-level page, current page — translated to Arabic
// and wrapped in dir="rtl". shadcn's own breadcrumb-rtl demo swaps the
// middle segment for a DropdownMenu; gsxui's breadcrumb example set only
// demonstrates the static BreadcrumbEllipsis (see Basic), so this example
// keeps that shape rather than introducing dropdown-menu composition here.
component Rtl() {
	<div dir="rtl" lang="ar">
		<uirtl.Breadcrumb>
			<uirtl.BreadcrumbList>
				<uirtl.BreadcrumbItem>
					<uirtl.BreadcrumbLink href="/">الرئيسية</uirtl.BreadcrumbLink>
				</uirtl.BreadcrumbItem>
				<uirtl.BreadcrumbSeparator/>
				<uirtl.BreadcrumbItem>
					<uirtl.BreadcrumbEllipsis/>
				</uirtl.BreadcrumbItem>
				<uirtl.BreadcrumbSeparator/>
				<uirtl.BreadcrumbItem>
					<uirtl.BreadcrumbLink href="/components">المكونات</uirtl.BreadcrumbLink>
				</uirtl.BreadcrumbItem>
				<uirtl.BreadcrumbSeparator/>
				<uirtl.BreadcrumbItem>
					<uirtl.BreadcrumbPage>مسار التنقل</uirtl.BreadcrumbPage>
				</uirtl.BreadcrumbItem>
			</uirtl.BreadcrumbList>
		</uirtl.Breadcrumb>
	</div>
}
