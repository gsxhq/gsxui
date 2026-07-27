package ui

import "github.com/gsxhq/gsx"

// Toggle is the shadcn/ui Toggle (registry/new-york-v4/ui/toggle.tsx),
// ported as a plain <button type="button"> carrying a server-visible
// `pressed` bool in place of Radix's TogglePrimitive.Root uncontrolled
// pressed/defaultPressed state — the same explicit-initial-state ADAPT as
// Tabs' `selected`/Accordion's `open`. Zero value (false) renders unpressed
// (aria-pressed="false" data-state="off"), matching Radix's own default;
// `pressed={true}` renders already-on. ui/toggle.js takes over from there:
// one capture-delegated click listener flips both attributes on every
// click (see its own header comment) — the same "server renders the
// initial state, JS re-stamps on interaction" split as Dialog/Tabs.
//
// A native <input type="checkbox"> port — the same "real form control, zero
// JS" idiom Switch and Checkbox use — was considered and rejected: input is
// a VOID element (no children, no closing tag), while Toggle's entire
// visible surface IS its children (an icon, text, or both — see upstream's
// toggle-demo/-outline examples), rendered inside the pressable control
// itself, not a label sibling to a hidden control. A <button> is the only
// element shape that can both be the toggle and hold arbitrary child
// content, so this port needs a real (small) behavior module the same way
// Tabs does, rather than riding free on browser :checked state the way
// Switch/Checkbox do.
//
// variant, size, and state are public CSS axes. toggle.js flips data-state
// and aria-pressed on interaction.
//
// ADAPT: data-variant/data-size are stamped via the house `|> default(...)`
// pattern for consistency with every other variant-backed component.
//
// Retargeted to nova density (2026-07-24 nova density map, `## toggle`).
// ADAPT: nova keys directional icon padding off `data-icon="inline-start|
// inline-end"` stamps this component doesn't emit; ported instead onto
// gsxui's existing has-[>svg]:px-* selector mechanism (the same substitution
// button.gsx's sizeClass makes — see its own doc comment), collapsing
// nova's matching inline-start/inline-end values into one has-[>svg]:px-*
// per size.
component Toggle(pressed bool, variant string, size string, children gsx.Node, attrs gsx.Attrs) {
	{{
		state := "off"
		if pressed {
			state = "on"
		}
	}}
	<button
		type="button"
		data-gsxui-toggle
		data-variant={variant |> default("default")}
		data-size={size |> default("default")}
		data-state={state}
		aria-pressed={pressed}
		{ attrs... } data-gsxui-slot-toggle
	>
		{ children }
	</button>
}
