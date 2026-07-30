package maia

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
		data-gsxui-toggle
		data-variant={variant |> default("default")}
		data-size={size |> default("default")}
		data-state={state}
		aria-pressed={pressed}
		class={
			"inline-flex items-center justify-center gap-1 rounded-lg bg-transparent text-sm font-medium whitespace-nowrap transition-[color,box-shadow] outline-none hover:bg-muted hover:text-muted-foreground focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 data-[state=on]:bg-accent data-[state=on]:text-accent-foreground [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4",
			switch variant {
			case "outline":
				"border border-input hover:bg-accent hover:text-accent-foreground"
			default:
				"bg-transparent"
			},
			switch size {
			case "sm":
				"h-7 min-w-7 rounded-[min(var(--radius-md),12px)] px-2.5 text-[0.8rem] has-[>svg]:px-1.5 [&_svg:not([class*='size-'])]:size-3.5"
			case "lg":
				"h-9 min-w-9 px-2.5 has-[>svg]:px-2"
			default:
				"h-8 min-w-8 px-2.5 has-[>svg]:px-2"
			}
		}
		{ attrs... }
		data-gsxui-slot-toggle
	>
		{ children }
	</button>
}
