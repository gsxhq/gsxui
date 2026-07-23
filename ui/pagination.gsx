package ui

import (
	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Pagination and its parts are the shadcn/ui Pagination
// (registry/new-york-v4/ui/pagination.tsx) — no Radix primitive underneath;
// every part is already a plain styled element (nav/ul/li/a/span). The one
// real dependency is PaginationLink, which composes button.gsx's own
// package-private base/variantClass/sizeClass helpers (flat package, directly
// callable) instead of duplicating buttonVariants — this is the
// pagination -> button dependency internal/registry derives from those
// identifiers and registry_test.go pins. ChevronLeft/ChevronRight/Ellipsis
// (Lucide's MoreHorizontal, see breadcrumb.gsx) come from ui/icon — the
// pagination -> icon dependency, also derived and pinned.
component Pagination(children gsx.Node, attrs gsx.Attrs) {
	<nav
		role="navigation"
		aria-label="pagination"
		data-slot="pagination"
		class="mx-auto flex w-full justify-center"
		{ attrs... }
	>
		{ children }
	</nav>
}

component PaginationContent(children gsx.Node, attrs gsx.Attrs) {
	<ul data-slot="pagination-content" class="flex flex-row items-center gap-1" { attrs... }>
		{ children }
	</ul>
}

component PaginationItem(children gsx.Node, attrs gsx.Attrs) {
	<li data-slot="pagination-item" { attrs... }>{ children }</li>
}

// PaginationLink renders the shadcn/ui PaginationLink onto a real <a>,
// composed from button.gsx's base/variantClass/sizeClass — the same
// buttonVariants({variant, size}) computation Button itself uses. isActive
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
		data-slot="pagination-link"
		data-active={isActive}
		href={href}
		class={ base, variantClass(variant), sizeClass(size) }
		{ attrs... }
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
		class="gap-1 px-2.5 sm:pl-2.5"
		{ attrs... }
	>
		<icon.ChevronLeft/>
		<span class="hidden sm:block">Previous</span>
	</PaginationLink>
}

component PaginationNext(href string, attrs gsx.Attrs) {
	<PaginationLink
		href={href}
		size="default"
		aria-label="Go to next page"
		class="gap-1 px-2.5 sm:pr-2.5"
		{ attrs... }
	>
		<span class="hidden sm:block">Next</span>
		<icon.ChevronRight/>
	</PaginationLink>
}

component PaginationEllipsis(attrs gsx.Attrs) {
	<span
		aria-hidden="true"
		data-slot="pagination-ellipsis"
		class="flex size-9 items-center justify-center"
		{ attrs... }
	>
		<icon.Ellipsis class="size-4"/>
		<span class="sr-only">More pages</span>
	</span>
}
