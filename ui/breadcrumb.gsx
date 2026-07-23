package ui

import (
	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Breadcrumb and its parts are the shadcn/ui Breadcrumb
// (registry/new-york-v4/ui/breadcrumb.tsx) — no Radix primitive underneath
// the original either; every part is already a plain styled element, Radix's
// Slot used only for BreadcrumbLink's asChild (dropped, see that component's
// own comment and docs/jsx-parity.md). BreadcrumbSeparator's default child
// (lucide's ChevronRight) and BreadcrumbEllipsis's MoreHorizontal both come
// from ui/icon (icon.ChevronRight, icon.Ellipsis — Lucide renamed
// MoreHorizontal to "ellipsis", the same rename precedent as Spinner's
// Loader2Icon/LoaderCircle, see ui/spinner.gsx) — this import is the
// breadcrumb -> icon dependency internal/registry derives and
// internal/registry/registry_test.go pins.
component Breadcrumb(children gsx.Node, attrs gsx.Attrs) {
	<nav aria-label="breadcrumb" data-slot="breadcrumb" { attrs... }>{ children }</nav>
}

component BreadcrumbList(children gsx.Node, attrs gsx.Attrs) {
	<ol
		data-slot="breadcrumb-list"
		class="flex flex-wrap items-center gap-1.5 text-sm break-words text-muted-foreground sm:gap-2.5"
		{ attrs... }
	>
		{ children }
	</ol>
}

component BreadcrumbItem(children gsx.Node, attrs gsx.Attrs) {
	<li data-slot="breadcrumb-item" class="inline-flex items-center gap-1.5" { attrs... }>
		{ children }
	</li>
}

// BreadcrumbLink renders a real <a> unconditionally — shadcn's own default
// (`const Comp = asChild ? Slot.Root : "a"`) already resolves to "a" for the
// dominant/only realistic use; the asChild tag-swap itself is GAP (narrow,
// dropped): no gsx equivalent renders an arbitrary caller component in this
// slot (e.g. a router Link), the same narrow gap as Button's asChild (see
// docs/jsx-parity.md). Behavior-attachment uses of asChild are covered by
// the data-attribute mechanism (see dialog).
component BreadcrumbLink(href string, children gsx.Node, attrs gsx.Attrs) {
	<a data-slot="breadcrumb-link" href={href} class="transition-colors hover:text-foreground" { attrs... }>
		{ children }
	</a>
}

component BreadcrumbPage(children gsx.Node, attrs gsx.Attrs) {
	<span
		data-slot="breadcrumb-page"
		role="link"
		aria-disabled="true"
		aria-current="page"
		class="font-normal text-foreground"
		{ attrs... }
	>
		{ children }
	</span>
}

// BreadcrumbSeparator defaults to a ChevronRight icon when the caller passes
// no children, exactly like shadcn's `{children ?? <ChevronRight />}` — pass
// children to override with any other glyph or text.
component BreadcrumbSeparator(children gsx.Node, attrs gsx.Attrs) {
	<li
		data-slot="breadcrumb-separator"
		role="presentation"
		aria-hidden="true"
		class="[&>svg]:size-3.5"
		{ attrs... }
	>
		{ if children != nil {
			{ children }
		} else {
			<icon.ChevronRight/>
		} }
	</li>
}

// BreadcrumbEllipsis takes no children — like shadcn's own version, its
// content is the fixed MoreHorizontal icon plus a screen-reader-only label,
// not a caller-supplied slot.
component BreadcrumbEllipsis(attrs gsx.Attrs) {
	<span
		data-slot="breadcrumb-ellipsis"
		role="presentation"
		aria-hidden="true"
		class="flex size-9 items-center justify-center"
		{ attrs... }
	>
		<icon.Ellipsis class="size-4"/>
		<span class="sr-only">More</span>
	</span>
}
