package sera

import "github.com/gsxhq/gsx"

// Label is the shadcn/ui Label. shadcn wraps Radix's LabelPrimitive.Root,
// which renders a plain <label> with an onMouseDown guard that suppresses
// double-click text selection; gsx ports the markup straight (base class
// already carries select-none, so the JS guard's effect is redundant here —
// see docs/jsx-parity.md).
component Label(children gsx.Node, attrs gsx.Attrs) {
	<label
		class={
			"gap-2 font-semibold peer-[[data-gsxui-slot-checkbox]]:text-sm peer-[[data-gsxui-slot-radio]]:text-sm peer-[[data-gsxui-slot-switch]]:text-sm group-data-[disabled=true]:opacity-50 peer-disabled:opacity-50 peer-[[data-gsxui-slot-checkbox]]:font-normal peer-[[data-gsxui-slot-radio]]:font-normal peer-[[data-gsxui-slot-switch]]:font-normal uppercase text-xs tracking-wide peer-[[data-gsxui-slot-checkbox]]:normal-case peer-[[data-gsxui-slot-checkbox]]:tracking-normal peer-[[data-gsxui-slot-radio]]:normal-case peer-[[data-gsxui-slot-radio]]:tracking-normal peer-[[data-gsxui-slot-switch]]:normal-case peer-[[data-gsxui-slot-switch]]:tracking-normal flex leading-relaxed items-center select-none group-data-[disabled=true]:pointer-events-none peer-disabled:cursor-not-allowed"
		}
		{ attrs... }
		data-gsxui-slot-label
	>
		{ children }
	</label>
}
