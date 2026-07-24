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
// variant/size are toggleVariants' cva map (default/outline,
// default/sm/lg), each a static class block picked by the JS-resolved prop
// value — no data-[variant=]:/data-[size=]: selectors exist upstream to
// preserve, so both port as a `switch` inside class={}, the same idiom as
// Button's variantClass/sizeClass pair (inlined here, like Item's own pair,
// since only Toggle uses it). The one data-keyed selector upstream DOES
// have, data-[state=on]:bg-accent/text-accent-foreground, carries verbatim
// in the shared base string below — it lights up once toggle.js flips
// data-state on click.
//
// ADAPT: data-variant/data-size are stamped via the house `|> default(...)`
// pattern (see button.gsx/dropdown.gsx, and button-group.gsx's own ADAPT
// note) for consistency with every other cva-backed component in this
// codebase, even though shadcn's own Toggle stamps neither — Radix's
// TogglePrimitive.Root receives only data-slot; toggleVariants resolves
// straight to className, no data attrs of its own to port.
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
		data-slot="toggle"
		data-gsxui-toggle
		data-variant={variant |> default("default")}
		data-size={size |> default("default")}
		data-state={state}
		aria-pressed={pressed}
		class={
			"inline-flex items-center justify-center gap-1 rounded-lg text-sm font-medium whitespace-nowrap transition-[color,box-shadow] outline-none hover:bg-muted hover:text-muted-foreground focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 data-[state=on]:bg-accent data-[state=on]:text-accent-foreground dark:aria-invalid:ring-destructive/40 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4",
			switch variant { case "outline": "border border-input bg-transparent hover:bg-accent hover:text-accent-foreground" default: "bg-transparent" },
			switch size { case "sm": "h-7 min-w-7 rounded-[min(var(--radius-md),12px)] px-2.5 has-[>svg]:px-1.5 text-[0.8rem] [&_svg:not([class*='size-'])]:size-3.5" case "lg": "h-9 min-w-9 px-2.5 has-[>svg]:px-2" default: "h-8 min-w-8 px-2.5 has-[>svg]:px-2" }
		}
		{ attrs... }
	>
		{ children }
	</button>
}
