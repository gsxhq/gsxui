package maia

import "github.com/gsxhq/gsx"

// Drawer composes Dialog's root and behavior while exposing a directional
// side-anchored content role. Drag-to-dismiss remains outside this component.
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
			"m-0 flex-col gap-4 shadow-lg transition ease-in-out bg-popover text-popover-foreground fixed z-50 text-sm duration-200 outline-none data-[state=closed]:animate-out data-[state=open]:animate-in backdrop:bg-[var(--overlay)] backdrop:backdrop-blur-xs backdrop:duration-200 data-[state=open]:backdrop:animate-in data-[state=open]:backdrop:fade-in-0 data-[state=closed]:backdrop:animate-out data-[state=closed]:backdrop:fade-out-0 open:flex",
			switch direction {
			case "left":
				"inset-y-0 start-0 end-auto h-full max-h-none w-3/4 rounded-e-xl border-e sm:max-w-sm data-[state=closed]:slide-out-to-start data-[state=open]:slide-in-from-start md:[&_[data-gsxui-slot-drawer-header]]:text-start"
			case "right":
				"inset-y-0 end-0 start-auto h-full max-h-none w-3/4 rounded-s-xl border-s sm:max-w-sm data-[state=closed]:slide-out-to-end data-[state=open]:slide-in-from-end md:[&_[data-gsxui-slot-drawer-header]]:text-start"
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
			<div class={ "mx-auto mt-4 h-1 w-[100px] shrink-0 rounded-full bg-muted" } data-gsxui-slot-drawer-handle></div>
		} }
		{ children }
	</dialog>
}

component DrawerHeader(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "flex flex-col gap-0.5 p-4" } { attrs... } data-gsxui-slot-drawer-header>{ children }</div>
}

component DrawerFooter(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "mt-auto flex flex-col gap-2 p-4" } { attrs... } data-gsxui-slot-drawer-footer>{ children }</div>
}

component DrawerTitle(children gsx.Node, attrs gsx.Attrs) {
	<h2 class={ "font-medium text-foreground" } { attrs... } data-gsxui-slot-drawer-title data-gsxui-slot-dialog-title>
		{ children }
	</h2>
}

component DrawerDescription(children gsx.Node, attrs gsx.Attrs) {
	<p
		class={ "text-sm text-muted-foreground" }
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
