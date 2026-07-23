package ui

import "github.com/gsxhq/gsx"

// Tabs and its parts are the shadcn/ui Tabs, minus Radix's client context —
// each part is a plain sibling component, no shared React-tree state. The
// root still needs to know which trigger/panel is active at first paint, so
// Tabs stamps data-value from the caller's value; TabsTrigger/TabsContent
// each take an explicit selected bool (the caller compares its own value to
// the group's value) and stamp aria-selected/data-state/tabindex/hidden from
// it — the gsx answer to "no context", same shape as the switch/checkbox
// explicit-state ADAPTs. ui/tabs/tabs.js takes over from there: click and
// roving-arrow-key activation, re-stamping state on every trigger/panel and
// emitting gsxui:change on the root. Requires the tabs behavior module
// (ui/tabs/tabs.js).
//
// ADAPT: shadcn's `orientation` (horizontal/vertical) and TabsList's
// `variant` (default/line) cva axis are both dropped — out of task scope, no
// param for either. Every class token whose sole purpose was to key off one
// of those two Radix-only accessed states is dropped as dead weight, same
// "drop the selector, don't ship dead CSS" call as avatar's size prop and
// dialog's close-button data-[state=open]: pair: the two
// group-data-[orientation=vertical]/tabs: tokens on Tabs' root and
// TabsTrigger, the data-[variant=line]/group-data-[variant=line]/tabs-list:
// family on TabsList/TabsTrigger (rounding, background, the after:
// indicator — invisible under the only variant we ship), and
// group-data-[variant=default]/tabs-list:data-[state=active]:shadow-sm
// unwraps to an unconditional data-[state=active]:shadow-sm. Root and list
// no longer stamp data-orientation/orientation/data-variant — nothing reads
// them. See docs/jsx-parity.md.
component Tabs(value string, children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="tabs" data-gsxui-tabs data-value={ value } class="flex flex-col gap-2" { attrs... }>{ children }</div>
}

component TabsList(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="tabs-list" role="tablist" class="bg-muted text-muted-foreground inline-flex h-9 w-fit items-center justify-center rounded-lg p-[3px]" { attrs... }>{ children }</div>
}

// TabsTrigger's selected bool is the explicit, server-visible stand-in for
// "does my value match the root's" — the caller (which already has both
// values in scope when building the tree) resolves the comparison; this
// component only renders the result. Zero value (false) renders the
// inactive state, matching a caller who forgets to pass it — never
// accidentally active.
component TabsTrigger(value string, selected bool, children gsx.Node, attrs gsx.Attrs) {
	{{
		state := "inactive"
		tabindex := -1
		if selected {
			state, tabindex = "active", 0
		}
	}}
	<button
		type="button"
		role="tab"
		data-slot="tabs-trigger"
		data-gsxui-tabs-trigger
		data-value={ value }
		data-state={ state }
		aria-selected={ selected }
		tabindex={ tabindex }
		class="relative inline-flex h-[calc(100%-1px)] flex-1 items-center justify-center gap-1.5 rounded-md border border-transparent px-2 py-1 text-sm font-medium whitespace-nowrap text-foreground/60 transition-all hover:text-foreground focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 focus-visible:outline-1 focus-visible:outline-ring disabled:pointer-events-none disabled:opacity-50 data-[state=active]:shadow-sm dark:text-muted-foreground dark:hover:text-foreground [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 data-[state=active]:bg-background data-[state=active]:text-foreground dark:data-[state=active]:border-input dark:data-[state=active]:bg-input/30 dark:data-[state=active]:text-foreground"
		{ attrs... }
	>{ children }</button>
}

// TabsContent's selected bool mirrors TabsTrigger's — same value-comparison
// contract, same zero-value-is-inactive default.
component TabsContent(value string, selected bool, children gsx.Node, attrs gsx.Attrs) {
	{{
		state := "inactive"
		if selected {
			state = "active"
		}
	}}
	<div
		role="tabpanel"
		data-slot="tabs-content"
		data-value={ value }
		data-state={ state }
		hidden={ !selected }
		class="flex-1 outline-none"
		{ attrs... }
	>{ children }</div>
}
