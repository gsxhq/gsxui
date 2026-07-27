package ui

import "github.com/gsxhq/gsx"

// Drawer composes Dialog's root and behavior while exposing a directional
// side-anchored content role. Drag-to-dismiss remains outside this component.
component Drawer(children gsx.Node, attrs gsx.Attrs) {
	<Dialog { withSlot("drawer", attrs)... }>{ children }</Dialog>
}

component DrawerTrigger(children gsx.Node, attrs gsx.Attrs) {
	<button
		data-gsxui-dialog-trigger
		type="button"
		aria-haspopup="dialog"
		aria-expanded="false"
		{ withSlot("drawer-trigger", attrs)... }
	>
		{ children }
	</button>
}

component DrawerContent(direction string, children gsx.Node, attrs gsx.Attrs) {
	<dialog
		data-gsxui-dialog-content
		data-state="closed"
		data-side={direction |> default("bottom")}
		{ withSlot("dialog-content", withSlot("drawer-content", attrs))... }
	>
		{ if direction == "" || direction == "bottom" {
			<div { withSlot("drawer-handle", nil)... }></div>
		} }
		{ children }
	</dialog>
}

component DrawerHeader(children gsx.Node, attrs gsx.Attrs) {
	<div { withSlot("drawer-header", attrs)... }>{ children }</div>
}

component DrawerFooter(children gsx.Node, attrs gsx.Attrs) {
	<div { withSlot("drawer-footer", attrs)... }>{ children }</div>
}

component DrawerTitle(children gsx.Node, attrs gsx.Attrs) {
	<h2 data-gsxui-dialog-title { withSlot("drawer-title", attrs)... }>{ children }</h2>
}

component DrawerDescription(children gsx.Node, attrs gsx.Attrs) {
	<p data-gsxui-dialog-description { withSlot("drawer-description", attrs)... }>{ children }</p>
}

component DrawerClose(children gsx.Node, attrs gsx.Attrs) {
	<button data-gsxui-dialog-close type="button" { withSlot("drawer-close", attrs)... }>{ children }</button>
}
