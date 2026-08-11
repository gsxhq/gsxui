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
// initial state, JS re-stamps on interaction" split as Dialog/Tabs. This is
// exactly why pressed/data-state is a runtime Tailwind state variant
// (`data-[state=on]:`) folded onto the recipe's root rule rather than a
// shape dimension: a dimension resolves to a literal class at generation
// time from the params passed at this callsite, but data-state can change
// in the browser with no server round trip at all.
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
// variant and size are shape dimensions, resolved to literal classes below.
// toggle.js flips data-state and aria-pressed on interaction; those (and
// disabled/aria-invalid) are runtime states folded onto the root rule, not
// dimensions — see registry/canonical/shapes/toggle.go.
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
//
// ui/toggle-group.gsx's ToggleGroupItem stamps data-gsxui-slot-toggle
// alongside its own marker WITHOUT calling any accessor here — see
// assets/css/styles/default.css's Toggle marker fallback block for how that
// composition keeps working after this migration.
component Toggle(pressed bool, variant string, size string, children gsx.Node, attrs gsx.Attrs) {
	{{
		state := "off"
		if pressed {
			state = "on"
		}
	}}
	<button
		type="button"
		data-variant={variant |> default("default")}
		data-size={size |> default("default")}
		data-state={state}
		aria-pressed={pressed}
		class={
			"hover:text-foreground aria-pressed:bg-muted focus-visible:border-ring focus-visible:ring-ring/30 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive gap-1 rounded-3xl text-sm font-medium transition-colors [&_svg:not([class*='size-'])]:size-4 inline-flex",
			switch variant { case "outline": "border-input hover:bg-muted border bg-transparent" default: "bg-transparent" },
			switch size {
			case "sm":
				"h-8 min-w-8 px-3 has-data-[icon=inline-end]:pr-2 has-data-[icon=inline-start]:pl-2"
			case "lg":
				"h-10 min-w-10 px-4 has-data-[icon=inline-end]:pr-3 has-data-[icon=inline-start]:pl-3"
			default:
				"h-9 min-w-9 px-3 has-data-[icon=inline-end]:pr-2.5 has-data-[icon=inline-start]:pl-2.5"
			}
		}
		{ attrs... }
		data-gsxui-slot-toggle
	>
		{ children }
	</button>
}
