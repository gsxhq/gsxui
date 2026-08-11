package vega

import "github.com/gsxhq/gsx"

// Drawer composes Dialog's root and behavior while exposing a directional
// side-anchored content role. Drag-to-dismiss remains outside this component.
//
// direction is physical, matching shadcn: direction="left"/"right" always
// anchors to that physical viewport edge (position, border, corner radius
// and slide direction), in both dir="ltr" and dir="rtl" documents.
// direction="top"/"bottom" is unaffected by text direction. Interior
// presentation — header/footer padding and title/description text
// alignment — is logical and follows dir normally.
//
// DrawerContent deliberately does NOT reuse Dialog's content presentation —
// it must not inherit the plain-modal box — so it carries the shared
// <dialog> chrome as its own recipe utilities instead. Every class attribute
// below is resolved to concrete utilities at generation time.
component Drawer(children gsx.Node, attrs gsx.Attrs) {
	<Dialog { attrs... } data-gsxui-slot-drawer>{ children }</Dialog>
}

component DrawerTrigger(children gsx.Node, attrs gsx.Attrs) {
	<button
		type="button"
		aria-haspopup="dialog"
		aria-expanded="false"
		{ attrs... }
		data-gsxui-slot-drawer-trigger
	>
		{ children }
	</button>
}

component DrawerContent(direction string, children gsx.Node, attrs gsx.Attrs) {
	<dialog
		class={
			"backdrop:bg-black/10 supports-backdrop-filter:backdrop:backdrop-blur-xs bg-popover text-popover-foreground flex h-auto flex-col text-sm data-[vaul-drawer-direction=bottom]:rounded-t-xl data-[vaul-drawer-direction=bottom]:border-t data-[vaul-drawer-direction=left]:rounded-r-xl data-[vaul-drawer-direction=left]:border-r data-[vaul-drawer-direction=right]:rounded-l-xl data-[vaul-drawer-direction=right]:border-l data-[vaul-drawer-direction=top]:rounded-b-xl data-[vaul-drawer-direction=top]:border-b m-0 gap-4 shadow-lg transition ease-in-out fixed z-50 duration-200 outline-none data-[state=closed]:animate-out data-[state=open]:animate-in backdrop:backdrop-blur-xs backdrop:duration-200 data-[state=open]:backdrop:animate-in data-[state=open]:backdrop:fade-in-0 data-[state=closed]:backdrop:animate-out data-[state=closed]:backdrop:fade-out-0 open:flex",
			switch direction {
			case "left":
				"inset-y-0 left-0 right-auto h-full max-h-none w-3/4 rounded-r-xl border-r sm:max-w-sm data-[state=closed]:slide-out-to-left data-[state=open]:slide-in-from-left md:[&_[data-gsxui-slot-drawer-header]]:text-start"
			case "right":
				"inset-y-0 right-0 left-auto h-full max-h-none w-3/4 rounded-l-xl border-l sm:max-w-sm data-[state=closed]:slide-out-to-right data-[state=open]:slide-in-from-right md:[&_[data-gsxui-slot-drawer-header]]:text-start"
			case "top":
				"inset-x-0 top-0 bottom-auto mb-24 h-auto max-h-[80vh] w-full max-w-none rounded-b-xl border-b data-[state=closed]:slide-out-to-top data-[state=open]:slide-in-from-top [&_[data-gsxui-slot-drawer-header]]:text-center"
			default:
				"inset-x-0 bottom-0 top-auto mt-24 h-auto max-h-[80vh] w-full max-w-none rounded-t-xl border-t data-[state=closed]:slide-out-to-bottom data-[state=open]:slide-in-from-bottom [&_[data-gsxui-slot-drawer-header]]:text-center"
			}
		}
		data-state="closed"
		data-side={direction |> default("bottom")}
		{ attrs... }
		data-gsxui-slot-drawer-content
		data-gsxui-slot-dialog-content
	>
		{ if direction == "" || direction == "bottom" {
			<div
				class={ "bg-muted mx-auto mt-4 hidden h-1.5 w-[100px] shrink-0 rounded-full" }
				data-gsxui-slot-drawer-handle
			></div>
		} }
		{ children }
	</dialog>
}

component DrawerHeader(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "gap-0.5 p-4 md:gap-1.5 md:text-left flex flex-col" } { attrs... } data-gsxui-slot-drawer-header>
		{ children }
	</div>
}

component DrawerFooter(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "gap-2 p-4 flex flex-col" } { attrs... } data-gsxui-slot-drawer-footer>{ children }</div>
}

component DrawerTitle(children gsx.Node, attrs gsx.Attrs) {
	<h2 class={ "text-foreground font-medium" } { attrs... } data-gsxui-slot-drawer-title data-gsxui-slot-dialog-title>
		{ children }
	</h2>
}

component DrawerDescription(children gsx.Node, attrs gsx.Attrs) {
	<p
		class={ "text-muted-foreground text-sm" }
		{ attrs... }
		data-gsxui-slot-drawer-description
		data-gsxui-slot-dialog-description
	>
		{ children }
	</p>
}

component DrawerClose(children gsx.Node, attrs gsx.Attrs) {
	<button data-gsxui-dialog-close type="button" { attrs... } data-gsxui-slot-drawer-close>{ children }</button>
}
