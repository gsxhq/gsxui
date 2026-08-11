package lyra

import "github.com/gsxhq/gsx"

// Collapsible is a native ungrouped <details>/<summary> disclosure. Its open
// parameter is the server-visible initial state; interaction remains native.
component Collapsible(open bool, children gsx.Node, attrs gsx.Attrs) {
	<details open={open} { attrs... } data-gsxui-slot-collapsible>
		{ children }
	</details>
}

component CollapsibleTrigger(children gsx.Node, attrs gsx.Attrs) {
	<summary class={ "list-none" } { attrs... } data-gsxui-slot-collapsible-trigger>
		{ children }
	</summary>
}

component CollapsibleContent(children gsx.Node, attrs gsx.Attrs) {
	<div { attrs... } data-gsxui-slot-collapsible-content>
		{ children }
	</div>
}
