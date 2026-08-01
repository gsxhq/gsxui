package breadcrumb

import "github.com/gsxhq/gsxui/ui"

// Rtl mirrors Basic's trail shape — linked root, collapsed-middle
// ellipsis, linked mid-level page, current page — translated to Arabic
// and wrapped in dir="rtl". shadcn's own breadcrumb-rtl demo swaps the
// middle segment for a DropdownMenu; gsxui's breadcrumb example set only
// demonstrates the static BreadcrumbEllipsis (see Basic), so this example
// keeps that shape rather than introducing dropdown-menu composition here.
component Rtl() {
	<div dir="rtl" lang="ar">
		<ui.Breadcrumb>
			<ui.BreadcrumbList>
				<ui.BreadcrumbItem>
					<ui.BreadcrumbLink href="/">الرئيسية</ui.BreadcrumbLink>
				</ui.BreadcrumbItem>
				<ui.BreadcrumbSeparator/>
				<ui.BreadcrumbItem>
					<ui.BreadcrumbEllipsis/>
				</ui.BreadcrumbItem>
				<ui.BreadcrumbSeparator/>
				<ui.BreadcrumbItem>
					<ui.BreadcrumbLink href="/components">المكونات</ui.BreadcrumbLink>
				</ui.BreadcrumbItem>
				<ui.BreadcrumbSeparator/>
				<ui.BreadcrumbItem>
					<ui.BreadcrumbPage>مسار التنقل</ui.BreadcrumbPage>
				</ui.BreadcrumbItem>
			</ui.BreadcrumbList>
		</ui.Breadcrumb>
	</div>
}
