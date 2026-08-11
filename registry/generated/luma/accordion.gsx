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
	<details
		name={name}
		open={open}
		class={ "data-[state=open]:bg-muted/50 not-last:border-b" }
		{ attrs... }
		data-gsxui-slot-accordion-item
	>
		{ children }
	</details>
}

component AccordionTrigger(children gsx.Node, attrs gsx.Attrs) {
	<summary
		class={ "gap-6 p-4 text-left text-sm font-medium hover:underline flex" }
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
	<div
		class={ "data-[state=open]:animate-accordion-down data-[state=closed]:animate-accordion-up px-4 text-sm" }
		{ attrs.Without("class")... }
		data-gsxui-slot-accordion-content
	>
		<div class={ "pt-0 pb-4", attrs.Class() } data-gsxui-slot-accordion-content-inner>{ children }</div>
	</div>
}
