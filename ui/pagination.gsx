package ui

import (
	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Pagination and its parts are the shadcn/ui Pagination
// (registry/new-york-v4/ui/pagination.tsx) — no Radix primitive underneath;
// every part is already a plain styled element (nav/ul/li/a/span). The one
// real dependency is PaginationLink, which composes Button's stable styling
// token and reflected variant/size axes instead of duplicating presentation.
// ChevronLeft/ChevronRight/Ellipsis
// (Lucide's MoreHorizontal, see breadcrumb.gsx) come from ui/icon — the
// pagination -> icon dependency, also derived and pinned.
component Pagination(children gsx.Node, attrs gsx.Attrs) {
	<nav
		role="navigation"
		aria-label="pagination"
		{ withSlot("pagination", attrs)... }
	>
		{ children }
	</nav>
}

component PaginationContent(children gsx.Node, attrs gsx.Attrs) {
	<ul { withSlot("pagination-content", attrs)... }>
		{ children }
	</ul>
}

component PaginationItem(children gsx.Node, attrs gsx.Attrs) {
	<li { withSlot("pagination-item", attrs)... }>{ children }</li>
}

// PaginationLink renders the shadcn/ui PaginationLink onto a real <a>,
// composed from Button's token and axes. isActive
// selects the "outline" variant (else "ghost") and stamps data-active plus,
// when true, aria-current="page" — the conditional-attribute mechanism
// (see docs/guide "conditional attributes") standing in for shadcn's
// `aria-current={isActive ? "page" : undefined}`, an attribute entirely
// absent when false, not merely empty. size defaults to "icon"
// (PaginationLinkProps' own `size = "icon"` default), distinct from
// Button's own "default" zero-value size.
component PaginationLink(href string, isActive bool, size string, children gsx.Node, attrs gsx.Attrs) {
	{{
		variant := "ghost"
		if isActive {
			variant = "outline"
		}
		if size == "" {
			size = "icon"
		}
	}}
	<a
		{ if isActive {
			aria-current="page"
		} }
		data-active={isActive}
		data-variant={variant}
		data-size={size}
		href={href}
		{ withSlot("button", withSlot("pagination-link", attrs))... }
	>
		{ children }
	</a>
}

// PaginationPrevious/PaginationNext hardcode their own content (icon + a
// sm:-only label) exactly like shadcn's versions — there is no children
// slot to override it, matching React's behavior where PaginationLink's
// literal JSX children always win over anything spread from ...props.
component PaginationPrevious(href string, attrs gsx.Attrs) {
	<PaginationLink
		href={href}
		size="default"
		aria-label="Go to previous page"
		{ withSlot("pagination-previous", attrs)... }
	>
		<icon.ChevronLeft/>
		<span { withSlot("pagination-previous-label", nil)... }>Previous</span>
	</PaginationLink>
}

component PaginationNext(href string, attrs gsx.Attrs) {
	<PaginationLink
		href={href}
		size="default"
		aria-label="Go to next page"
		{ withSlot("pagination-next", attrs)... }
	>
		<span { withSlot("pagination-next-label", nil)... }>Next</span>
		<icon.ChevronRight/>
	</PaginationLink>
}

component PaginationEllipsis(attrs gsx.Attrs) {
	<span
		aria-hidden="true"
		{ withSlot("pagination-ellipsis", attrs)... }
	>
		<icon.Ellipsis/>
		<span { withSlot("pagination-ellipsis-label", nil)... }>More pages</span>
	</span>
}
