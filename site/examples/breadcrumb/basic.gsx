// Package breadcrumb holds the site's example gsx components for
// ui/breadcrumb.
package breadcrumb

import "github.com/gsxhq/gsxui/ui"

// Basic renders a realistic trail: a linked root, a collapsed-middle
// ellipsis, a linked mid-level page, and the current (non-link) page —
// exercising both BreadcrumbSeparator's default chevron and
// BreadcrumbEllipsis.
component Basic() {
	<ui.Breadcrumb>
		<ui.BreadcrumbList>
			<ui.BreadcrumbItem>
				<ui.BreadcrumbLink href="/">Home</ui.BreadcrumbLink>
			</ui.BreadcrumbItem>
			<ui.BreadcrumbSeparator/>
			<ui.BreadcrumbItem>
				<ui.BreadcrumbEllipsis/>
			</ui.BreadcrumbItem>
			<ui.BreadcrumbSeparator/>
			<ui.BreadcrumbItem>
				<ui.BreadcrumbLink href="/components">Components</ui.BreadcrumbLink>
			</ui.BreadcrumbItem>
			<ui.BreadcrumbSeparator/>
			<ui.BreadcrumbItem>
				<ui.BreadcrumbPage>Breadcrumb</ui.BreadcrumbPage>
			</ui.BreadcrumbItem>
		</ui.BreadcrumbList>
	</ui.Breadcrumb>
}
