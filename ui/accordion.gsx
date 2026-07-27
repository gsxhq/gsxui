package ui

import (
	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Accordion uses the native grouped <details name>/<summary> disclosure
// mechanism. The root's data-name is descriptive; matching item name
// attributes provide exclusive-open behavior without JavaScript.
component Accordion(name string, children gsx.Node, attrs gsx.Attrs) {
	<div data-name={name} { withSlot("accordion", attrs)... }>{ children }</div>
}

component AccordionItem(name string, open bool, children gsx.Node, attrs gsx.Attrs) {
	<details name={name} open={open} { withSlot("accordion-item", attrs)... }>
		{ children }
	</details>
}

component AccordionTrigger(children gsx.Node, attrs gsx.Attrs) {
	<summary { withSlot("accordion-trigger", attrs)... }>
		{ children }
		<icon.ChevronDown { withSlot("accordion-trigger-icon", nil)... }/>
	</summary>
}

// The outer token owns disclosure mechanics and non-class attributes. Caller
// classes join the inner padding token so utilities override its defaults on
// the same box.
component AccordionContent(children gsx.Node, attrs gsx.Attrs) {
	{{
		var innerAttrs gsx.Attrs
		if class, ok := attrs.Get("class"); ok {
			innerAttrs = gsx.Attrs{{Key: "class", Value: class}}
		}
	}}
	<div { withSlot("accordion-content", attrs.Without("class"))... }>
		<div { withSlot("accordion-content-inner", innerAttrs)... }>{ children }</div>
	</div>
}
