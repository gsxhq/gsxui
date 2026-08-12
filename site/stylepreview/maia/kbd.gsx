package maia

import "github.com/gsxhq/gsx"

// Kbd and KbdGroup are the shadcn/ui Kbd. Straight port: both render onto
// real <kbd> elements — browsers freely nest <kbd> inside <kbd>, which is
// exactly how KbdGroup models a compound shortcut like "Ctrl Shift K" (a
// KbdGroup of Kbds). The style pack uses the stable tooltip-content ancestor
// token for the corresponding nested presentation.
component Kbd(children gsx.Node, attrs gsx.Attrs) {
	<kbd
		class={
			"bg-muted text-muted-foreground [[data-gsxui-slot-tooltip-content]_&]:bg-background/20 [[data-gsxui-slot-tooltip-content]_&]:text-background dark:[[data-gsxui-slot-tooltip-content]_&]:bg-background/10 h-5 w-fit min-w-5 gap-1 rounded-sm px-1 font-sans text-xs font-medium [&_svg:not([class*='size-'])]:size-3 inline-flex pointer-events-none items-center justify-center select-none"
		}
		{ attrs... }
		data-gsxui-slot-kbd
	>
		{ children }
	</kbd>
}

// KbdGroup wraps multiple Kbds to render a compound shortcut. shadcn types
// its props as React.ComponentProps<"div"> but the component itself renders
// a <kbd> element (registry/new-york-v4/ui/kbd.tsx, verified) — ported
// verbatim, tag included (see docs/jsx-parity.md).
component KbdGroup(children gsx.Node, attrs gsx.Attrs) {
	<kbd class={ "gap-1 inline-flex items-center" } { attrs... } data-gsxui-slot-kbd-group>
		{ children }
	</kbd>
}
