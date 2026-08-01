package canonical

import (
	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Accordion uses the native grouped <details name>/<summary> disclosure
// mechanism. The root's data-name is descriptive; matching item name
// attributes provide exclusive-open behavior without JavaScript.
component Accordion(name string, children gsx.Node, attrs gsx.Attrs) {
	<div data-name={name} { attrs... } data-gsxui-slot-accordion>{ children }</div>
}

component AccordionItem(name string, open bool, children gsx.Node, attrs gsx.Attrs) {
	<details name={name} open={open} class={ accordion.Item() } { attrs... } data-gsxui-slot-accordion-item>
		{ children }
	</details>
}

component AccordionTrigger(children gsx.Node, attrs gsx.Attrs) {
	<summary class={ accordion.Trigger() } { attrs... } data-gsxui-slot-accordion-trigger>
		{ children }
		<icon.ChevronDown class={ accordion.TriggerIcon() } data-gsxui-slot-accordion-trigger-icon/>
	</summary>
}

// The outer token owns disclosure mechanics and non-class attributes. Caller
// classes join the inner padding token so utilities override its defaults on
// the same box.
component AccordionContent(children gsx.Node, attrs gsx.Attrs) {
	<div class={ accordion.Content() } { attrs.Without("class")... } data-gsxui-slot-accordion-content>
		<div class={ accordion.ContentInner(), attrs.Class() } data-gsxui-slot-accordion-content-inner>{ children }</div>
	</div>
}
