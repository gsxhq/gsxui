package canonical

import "github.com/gsxhq/gsx"

// Kbd and KbdGroup are the shadcn/ui Kbd. Straight port: both render onto
// real <kbd> elements — browsers freely nest <kbd> inside <kbd>, which is
// exactly how KbdGroup models a compound shortcut like "Ctrl Shift K" (a
// KbdGroup of Kbds). The style pack uses the stable tooltip-content ancestor
// token for the corresponding nested presentation.
component Kbd(children gsx.Node, attrs gsx.Attrs) {
	<kbd class={ kbd.Root() } { attrs... } data-gsxui-slot-kbd>
		{ children }
	</kbd>
}

// KbdGroup wraps multiple Kbds to render a compound shortcut. shadcn types
// its props as React.ComponentProps<"div"> but the component itself renders
// a <kbd> element (registry/new-york-v4/ui/kbd.tsx, verified) — ported
// verbatim, tag included (see docs/jsx-parity.md).
component KbdGroup(children gsx.Node, attrs gsx.Attrs) {
	<kbd class={ kbd.Group() } { attrs... } data-gsxui-slot-kbd-group>
		{ children }
	</kbd>
}
