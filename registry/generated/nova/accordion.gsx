package ui

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
	<details name={name} open={open} class={ "not-last:border-b" } { attrs... } data-gsxui-slot-accordion-item>
		{ children }
	</details>
}

component AccordionTrigger(children gsx.Node, attrs gsx.Attrs) {
	<summary
		class={
			"focus-visible:ring-ring/50 focus-visible:border-ring focus-visible:after:border-ring rounded-lg py-2.5 text-left text-sm font-medium hover:underline focus-visible:ring-3 flex items-start outline-none transition-all border border-transparent"
		}
		{ attrs... }
		data-gsxui-slot-accordion-trigger
	>
		{ children }
		<icon.ChevronDown
			class={
				"text-muted-foreground ml-auto size-4 shrink-0 transition-transform duration-200 [[data-gsxui-slot-accordion-item][open]_&]:rotate-180"
			}
			data-gsxui-slot-accordion-trigger-icon
		/>
	</summary>
}

// The outer token owns disclosure mechanics and non-class attributes. Caller
// classes join the inner padding token so utilities override its defaults on
// the same box.
component AccordionContent(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "text-sm overflow-hidden" } { attrs.Without("class")... } data-gsxui-slot-accordion-content>
		<div class={ "pt-0 pb-2.5", attrs.Class() } data-gsxui-slot-accordion-content-inner>{ children }</div>
	</div>
}
