package ui

import (
	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui/i18n"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Pagination and its parts are the shadcn/ui Pagination
// (registry/new-york-v4/ui/pagination.tsx) — no Radix primitive underneath;
// every part is a plain styled element (nav/ul/li/span) except
// PaginationLink, which renders through Button as upstream's does.
// ChevronLeft/ChevronRight/Ellipsis
// (Lucide's MoreHorizontal, see breadcrumb.gsx) come from ui/icon — the
// pagination -> icon dependency, also derived and pinned.
component Pagination(children gsx.Node, attrs gsx.Attrs) {
	<nav
		role="navigation"
		aria-label={i18n.T("pagination")}
		class={ "mx-auto flex w-full justify-center" }
		{ attrs... }
		data-gsxui-slot-pagination
	>
		{ children }
	</nav>
}

component PaginationContent(children gsx.Node, attrs gsx.Attrs) {
	<ul class={ "gap-1 flex items-center" } { attrs... } data-gsxui-slot-pagination-content>
		{ children }
	</ul>
}

component PaginationItem(children gsx.Node, attrs gsx.Attrs) {
	<li { attrs... } data-gsxui-slot-pagination-item>{ children }</li>
}

// PaginationLink is the shadcn/ui PaginationLink: a Button rendered as its
// <a>, as upstream's `<Button asChild>` is. isActive selects the "outline"
// variant (else "ghost") and stamps data-active plus, when true,
// aria-current="page" — the conditional-attribute mechanism (see docs/guide
// "conditional attributes") standing in for shadcn's
// `aria-current={isActive ? "page" : undefined}`, an attribute entirely
// absent when false, not merely empty. size defaults to "icon"
// (PaginationLinkProps' own `size = "icon"` default), distinct from
// Button's own "default" zero-value size. An empty href renders Button's
// <button> instead.
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
	<Button
		href={href}
		variant={variant}
		size={size}
		{ if isActive {
			aria-current="page"
		} }
		data-active={isActive}
		{ attrs... }
		data-gsxui-slot-pagination-link
	>
		{ children }
	</Button>
}

// PaginationPrevious/PaginationNext hardcode their own content (icon + a
// sm:-only label) exactly like shadcn's versions — there is no children
// slot to override it, matching React's behavior where PaginationLink's
// literal JSX children always win over anything spread from ...props.
component PaginationPrevious(href string, attrs gsx.Attrs) {
	<PaginationLink
		href={href}
		size="default"
		aria-label={i18n.T("Go to previous page")}
		class={ "pl-2 has-[>svg]:pl-2" }
		{ attrs... }
		data-gsxui-slot-pagination-previous
	>
		<icon.ChevronLeft class={ "rtl:rotate-180" }/>
		<span class={ "hidden sm:block" } data-gsxui-slot-pagination-previous-label>{ i18n.T("Previous") }</span>
	</PaginationLink>
}

component PaginationNext(href string, attrs gsx.Attrs) {
	<PaginationLink
		href={href}
		size="default"
		aria-label={i18n.T("Go to next page")}
		class={ "pr-2 has-[>svg]:pr-2" }
		{ attrs... }
		data-gsxui-slot-pagination-next
	>
		<span class={ "hidden sm:block" } data-gsxui-slot-pagination-next-label>{ i18n.T("Next") }</span>
		<icon.ChevronRight class={ "rtl:rotate-180" }/>
	</PaginationLink>
}

component PaginationEllipsis(attrs gsx.Attrs) {
	<span
		aria-hidden="true"
		class={ "size-9 items-center justify-center [&_svg:not([class*='size-'])]:size-4 flex" }
		{ attrs... }
		data-gsxui-slot-pagination-ellipsis
	>
		<icon.Ellipsis/>
		<span data-gsxui-slot-pagination-ellipsis-label>{ i18n.T("More pages") }</span>
	</span>
}
