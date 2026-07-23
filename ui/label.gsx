package ui

import "github.com/gsxhq/gsx"

// Label is the shadcn/ui Label. shadcn wraps Radix's LabelPrimitive.Root,
// which renders a plain <label> with an onMouseDown guard that suppresses
// double-click text selection; gsx ports the markup straight (base class
// already carries select-none, so the JS guard's effect is redundant here —
// see docs/jsx-parity.md).
component Label(children gsx.Node, attrs gsx.Attrs) {
	<label
		data-slot="label"
		class="flex items-center gap-2 text-sm leading-none font-medium select-none group-data-[disabled=true]:pointer-events-none group-data-[disabled=true]:opacity-50 peer-disabled:cursor-not-allowed peer-disabled:opacity-50"
		{ attrs... }
	>
		{ children }
	</label>
}
