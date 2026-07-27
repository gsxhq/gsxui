package ui

import "github.com/gsxhq/gsx"

// Collapsible is a native ungrouped <details>/<summary> disclosure. Its open
// parameter is the server-visible initial state; interaction remains native.
component Collapsible(open bool, children gsx.Node, attrs gsx.Attrs) {
	<details open={open} { withSlot("collapsible", attrs)... }>
		{ children }
	</details>
}

component CollapsibleTrigger(children gsx.Node, attrs gsx.Attrs) {
	<summary { withSlot("collapsible-trigger", attrs)... }>
		{ children }
	</summary>
}

component CollapsibleContent(children gsx.Node, attrs gsx.Attrs) {
	<div { withSlot("collapsible-content", attrs)... }>
		{ children }
	</div>
}
