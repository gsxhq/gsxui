package lyra

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
	<nav aria-label="breadcrumb" { attrs... } data-gsxui-slot-breadcrumb>{ children }</nav>
}

component BreadcrumbList(children gsx.Node, attrs gsx.Attrs) {
	<ol class={ "text-muted-foreground gap-1.5 text-xs flex flex-wrap" } { attrs... } data-gsxui-slot-breadcrumb-list>
		{ children }
	</ol>
}

component BreadcrumbItem(children gsx.Node, attrs gsx.Attrs) {
	<li class={ "gap-1 inline-flex" } { attrs... } data-gsxui-slot-breadcrumb-item>
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
	<a href={href} class={ "hover:text-foreground transition-colors" } { attrs... } data-gsxui-slot-breadcrumb-link>
		{ children }
	</a>
}

component BreadcrumbPage(children gsx.Node, attrs gsx.Attrs) {
	<span
		role="link"
		aria-disabled="true"
		aria-current="page"
		class={ "text-foreground font-normal" }
		{ attrs... }
		data-gsxui-slot-breadcrumb-page
	>
		{ children }
	</span>
}

// BreadcrumbSeparator defaults to a ChevronRight icon when the caller passes
// no children, exactly like shadcn's `{children ?? <ChevronRight />}` — pass
// children to override with any other glyph or text.
component BreadcrumbSeparator(children gsx.Node, attrs gsx.Attrs) {
	<li
		role="presentation"
		aria-hidden="true"
		class={ "[&>svg]:size-3.5" }
		{ attrs... }
		data-gsxui-slot-breadcrumb-separator
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
		role="presentation"
		aria-hidden="true"
		class={ "size-5 [&>svg]:size-4 flex" }
		{ attrs... }
		data-gsxui-slot-breadcrumb-ellipsis
	>
		<icon.Ellipsis/>
		<span data-gsxui-slot-breadcrumb-ellipsis-label>More</span>
	</span>
}
