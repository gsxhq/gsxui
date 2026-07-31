package maia

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
	<details name={name} open={open} class={ "border-b last:border-b-0" } { attrs... } data-gsxui-slot-accordion-item>
		{ children }
	</details>
}

component AccordionTrigger(children gsx.Node, attrs gsx.Attrs) {
	<summary
		class={
			"flex list-none items-start justify-between rounded-lg py-2.5 text-left text-sm font-medium transition-all outline-none hover:underline focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:pointer-events-none disabled:opacity-50"
		}
		{ attrs... }
		data-gsxui-slot-accordion-trigger
	>
		{ children }
		<icon.ChevronDown
			class={
				"size-4 shrink-0 text-muted-foreground transition-transform duration-200 [[data-gsxui-slot-accordion-item][open]_&]:rotate-180"
			}
			data-gsxui-slot-accordion-trigger-icon
		/>
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
	<div class={ "overflow-hidden text-sm" } { attrs.Without("class")... } data-gsxui-slot-accordion-content>
		<div class={ "pt-0 pb-2.5" } { innerAttrs... } data-gsxui-slot-accordion-content-inner>{ children }</div>
	</div>
}
