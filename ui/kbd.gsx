package ui

import "github.com/gsxhq/gsx"

// Kbd and KbdGroup are the shadcn/ui Kbd. Straight port: both render onto
// real <kbd> elements — browsers freely nest <kbd> inside <kbd>, which is
// exactly how KbdGroup models a compound shortcut like "Ctrl Shift K" (a
// KbdGroup of Kbds). The [[data-slot=tooltip-content]_&]:... tokens are a
// real, exercisable selector: nesting a Kbd inside a ui.TooltipContent
// (data-slot="tooltip-content") activates them, ported as-is.
component Kbd(children gsx.Node, attrs gsx.Attrs) {
	<kbd
		data-slot="kbd"
		class={
			"pointer-events-none inline-flex h-5 w-fit min-w-5 items-center justify-center gap-1 rounded-sm bg-muted px-1 font-sans text-xs font-medium text-muted-foreground select-none",
			"[&_svg:not([class*='size-'])]:size-3",
			"[[data-slot=tooltip-content]_&]:bg-background/20 [[data-slot=tooltip-content]_&]:text-background dark:[[data-slot=tooltip-content]_&]:bg-background/10"
		}
		{ attrs... }
	>
		{ children }
	</kbd>
}

// KbdGroup wraps multiple Kbds to render a compound shortcut. shadcn types
// its props as React.ComponentProps<"div"> but the component itself renders
// a <kbd> element (registry/new-york-v4/ui/kbd.tsx, verified) — ported
// verbatim, tag included (see docs/jsx-parity.md).
component KbdGroup(children gsx.Node, attrs gsx.Attrs) {
	<kbd data-slot="kbd-group" class="inline-flex items-center gap-1" { attrs... }>
		{ children }
	</kbd>
}
