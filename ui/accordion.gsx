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

// The inner token keeps content padding independently styleable while all
// caller attributes, including class, are forwarded once to the public part.
component AccordionContent(children gsx.Node, attrs gsx.Attrs) {
	<div { withSlot("accordion-content", attrs)... }>
		<div { withSlot("accordion-content-inner", nil)... }>{ children }</div>
	</div>
}
