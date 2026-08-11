package lyra

import "github.com/gsxhq/gsx"

// Label is the shadcn/ui Label. shadcn wraps Radix's LabelPrimitive.Root,
// which renders a plain <label> with an onMouseDown guard that suppresses
// double-click text selection; gsx ports the markup straight (base class
// already carries select-none, so the JS guard's effect is redundant here —
// see docs/jsx-parity.md).
component Label(children gsx.Node, attrs gsx.Attrs) {
	<label
		class={ "gap-2 text-xs leading-none group-data-[disabled=true]:opacity-50 peer-disabled:opacity-50 flex" }
		{ attrs... }
		data-gsxui-slot-label
	>
		{ children }
	</label>
}
