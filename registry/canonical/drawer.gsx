package canonical

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
		data-gsxui-dialog-trigger
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
			drawer.Content(),
			drawer.ContentSide(direction),
		}
		data-gsxui-dialog-content
		data-state="closed"
		data-side={direction |> default("bottom")}
		{ attrs... }
		data-gsxui-slot-drawer-content
		data-gsxui-slot-dialog-content
	>
		{ if direction == "" || direction == "bottom" {
			<div class={ drawer.Handle() } data-gsxui-slot-drawer-handle></div>
		} }
		{ children }
	</dialog>
}

component DrawerHeader(children gsx.Node, attrs gsx.Attrs) {
	<div class={ drawer.Header() } { attrs... } data-gsxui-slot-drawer-header>{ children }</div>
}

component DrawerFooter(children gsx.Node, attrs gsx.Attrs) {
	<div class={ drawer.Footer() } { attrs... } data-gsxui-slot-drawer-footer>{ children }</div>
}

component DrawerTitle(children gsx.Node, attrs gsx.Attrs) {
	<h2 class={ drawer.Title() } data-gsxui-dialog-title { attrs... } data-gsxui-slot-drawer-title>{ children }</h2>
}

component DrawerDescription(children gsx.Node, attrs gsx.Attrs) {
	<p class={ drawer.Description() } data-gsxui-dialog-description { attrs... } data-gsxui-slot-drawer-description>{ children }</p>
}

component DrawerClose(children gsx.Node, attrs gsx.Attrs) {
	<button data-gsxui-dialog-close type="button" { attrs... } data-gsxui-slot-drawer-close>{ children }</button>
}
