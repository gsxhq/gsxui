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
// dialog's close-button open-state pair: the ancestor-orientation rules on
// Tabs' root and TabsTrigger, the line-variant family on
// TabsList/TabsTrigger (rounding, background, the after-element indicator
// — invisible under the only variant we ship), and the default-list active
// shadow rule, which unwraps to an unconditional active-state shadow. Root and list
// no longer stamp data-orientation/orientation/data-variant — nothing reads
// them. See docs/jsx-parity.md.
//
// Retargeted to nova density (2026-07-24 nova density map, `## tabs`).
// ADAPT: nova keys the trigger's directional icon padding off
// `data-icon="inline-start|inline-end"` stamps this component doesn't emit;
// ported instead onto gsxui's existing has-[>svg]:px-* selector mechanism
// (the same substitution button.gsx's sizeClass and toggle.gsx make — see
// their own doc comments), collapsing nova's matching inline-start/
// inline-end value (both px-1) into one has-[>svg]:px-1.
component Tabs(value string, children gsx.Node, attrs gsx.Attrs) {
	<div data-value={value} class={ "gap-2 flex flex-col" } { attrs... } data-gsxui-slot-tabs>{ children }</div>
}

component TabsList(children gsx.Node, attrs gsx.Attrs) {
	<div
		role="tablist"
		class={ "p-1 group-data-horizontal/tabs:h-10 inline-flex" }
		{ attrs... }
		data-gsxui-slot-tabs-list
	>
		{ children }
	</div>
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
		data-value={value}
		data-state={state}
		aria-selected={selected}
		tabindex={tabindex}
		class={
			"gap-2 border border-transparent px-4 py-1.5 text-xs font-semibold group-data-vertical/tabs:px-4 group-data-vertical/tabs:py-2 [&_svg:not([class*='size-'])]:size-3.5 has-data-[icon=inline-end]:pr-2.5 has-data-[icon=inline-start]:pl-2.5 uppercase tracking-wider inline-flex"
		}
		{ attrs... }
		data-gsxui-slot-tabs-trigger
	>
		{ children }
	</button>
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
		data-value={value}
		data-state={state}
		hidden={!selected}
		class={ "text-sm" }
		{ attrs... }
		data-gsxui-slot-tabs-content
	>
		{ children }
	</div>
}
