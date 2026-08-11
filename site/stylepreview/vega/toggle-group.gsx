package vega

import "github.com/gsxhq/gsx"

// ToggleGroup and ToggleGroupItem are the shadcn/ui ToggleGroup
// (registry/new-york-v4/ui/toggle-group.tsx) plus the runtime attributes
// Radix's @radix-ui/react-toggle-group / react-roving-focus / react-toggle
// stamp at mount (traced from their dist/index.mjs, not guessed from the
// .tsx — see 2026-07-24 controls source map, `## toggle-group`).
//
// GAP (no context, group->item is explicit params): Radix's
// ToggleGroupContext.Provider broadcasts type/variant/size/spacing from
// ToggleGroup down to every ToggleGroupItem automatically; gsx has no
// context, so the caller — which already has all four values in scope when
// building the tree — passes them explicitly to ToggleGroup AND to every
// ToggleGroupItem, same shape as `## tabs`' selected. Radix's group-level
// `disabled` OR-cascading onto every item (context.disabled ||
// props.disabled) is the same story: no param for it here, disabled flows
// through each element's own attrs bag (native <button disabled> on the
// item, an inert-but-present attribute on the root's <div>) — the caller
// disables the whole group by passing disabled to every item explicitly.
//
// ADAPT (groupType, not type): shadcn's own prop is named `type` — a Go
// keyword, unusable as a component parameter name (unlike `select`/`switch`,
// which only needed a Go-keyword workaround at the file/package level, see
// select.gsx's own doc comment). `groupType == "single"` selects
// role="radiogroup" (root) / role="radio" aria-checked (item); anything else
// (including the Go zero value "") renders role="toolbar" (root) /
// aria-pressed, no role override (item) — the "multiple" behavior. Radix
// itself throws if `type` is omitted; Go can't throw at render-construction
// time the way React throws at mount, so this is a doc-comment API-design
// note (pass groupType explicitly), not a runtime check.
//
// ADAPT (root shadow selector dropped, not ported dead): new-york-v4's root
// class carries `data-[spacing=default]:data-[variant=outline]:shadow-xs`,
// but `data-spacing` is stamped as the literal (numeric) prop value, never
// the string "default" — the selector can never match anything the
// component itself renders, dead CSS in the shadcn source itself (see
// docs/jsx-parity.md FINDING for the full trace). nova's own
// .cn-toggle-group already drops it rather than porting it as dead weight;
// this port follows nova's own precedent instead of the house's usual
// "port dead weight, ledger it" convention — see docs/jsx-parity.md.
//
// ADAPT (horizontal-only v1): data-orientation="horizontal" is stamped on
// both root and item, but only the horizontal corner-rounding selectors are
// ported (data-[orientation=horizontal]:data-[spacing=0]:first:rounded-l-lg
// / last:rounded-r-lg) — vertical is real new functionality nova adds
// (Radix's own new-york-v4 markup never varies rounding by orientation) and
// is out of v1 scope; see docs/jsx-parity.md GAP. toggle-group.js gates its
// arrow-key handling the same way (ArrowLeft/Up and ArrowRight/Down both
// move focus, since orientation is always "horizontal" here).
//
// MECHANISM (roving tabindex, JS-normalized at init): server renders every
// item with no tabindex attribute at all — a graceful no-JS fallback where
// every item is its own tab stop. toggle-group.js sets tabindex="0" on
// exactly one item at init, then maintains that invariant on every arrow
// move and click — see its own header comment. Entry-tabstop priority is
// type-agnostic (matching Radix's own traced RovingFocusGroup: the "active"
// candidate wins regardless of type, else fall back to the first item —
// the map's behavior contract never conditions this on type="single" vs
// "multiple"): the pressed-and-enabled item if one exists, else the first
// non-disabled item — the same rule applies whether groupType is "single"
// or "multiple".
component ToggleGroup(groupType string, variant string, size string, spacing string, children gsx.Node, attrs gsx.Attrs) {
	{{
		sp := spacing
		if sp == "" {
			sp = "0"
		}
		role := "toolbar"
		if groupType == "single" {
			role = "radiogroup"
		}
	}}
	<div
		data-variant={variant |> default("default")}
		data-size={size |> default("default")}
		data-spacing={sp}
		data-orientation="horizontal"
		role={role}
		style=css`--gap: @{sp}`
		class={
			"group/toggle-group",
			"rounded-md data-[spacing=0]:data-[variant=outline]:shadow-xs flex",
			switch size { case "sm": "rounded-[min(var(--radius-md),10px)]" case "lg": "rounded-lg" default: "rounded-lg" }
		}
		{ attrs... }
		data-gsxui-slot-toggle-group
	>
		{ children }
	</div>
}

// ToggleGroupItem composes ordered tokens "toggle toggle-group-item".
// Variant, size, spacing, orientation, and state are public CSS axes; only
// the ARIA attribute pair differs between single and multiple groups.
//
// MECHANISM (single-type replace-on-activate): clicking a new item in a
// type="single" group simply sets a new single value — there is no
// group-level "uncheck the others" loop needed here because ToggleGroup
// itself carries no value state to update; toggle-group.js re-stamps
// data-state="off"/aria-checked="false" on every OTHER item sharing the
// same root when a type="single" item activates, so exactly one item shows
// data-state="on" at a time (Radix's onItemActivate === setValue, restated
// for a stateless server render). Clicking the already-pressed item in a
// single group toggles it off (Radix allows an empty single value unless a
// caller opts otherwise) — port the same replace-on-activate mechanic.
component ToggleGroupItem(groupType string, variant string, size string, spacing string, pressed bool, value string, children gsx.Node, attrs gsx.Attrs) {
	{{
		sp := spacing
		if sp == "" {
			sp = "0"
		}
		state := "off"
		if pressed {
			state = "on"
		}
		orientation := "horizontal"
	}}
	<button
		type="button"
		data-variant={variant |> default("default")}
		data-size={size |> default("default")}
		data-spacing={sp}
		data-orientation={orientation}
		data-state={state}
		data-value={value}
		{ if groupType == "single" {
			role="radio"
			aria-checked={pressed}
		} else {
			aria-pressed={pressed}
		} }
		class={
			"hover:text-foreground focus-visible:border-ring focus-visible:ring-ring/50 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive disabled:pointer-events-none disabled:opacity-50 data-[state=on]:bg-muted data-[variant=outline]:border-input data-[variant=outline]:hover:bg-muted data-[variant=outline]:border data-[variant=outline]:shadow-xs gap-1 rounded-md text-sm font-medium transition-[color,box-shadow] [&_svg:not([class*='size-'])]:size-4 group-data-[spacing=0]/toggle-group:rounded-none group-data-[spacing=0]/toggle-group:px-2 group-data-[spacing=0]/toggle-group:shadow-none group-data-horizontal/toggle-group:data-[spacing=0]:first:rounded-l-md group-data-vertical/toggle-group:data-[spacing=0]:first:rounded-t-md group-data-horizontal/toggle-group:data-[spacing=0]:last:rounded-r-md group-data-vertical/toggle-group:data-[spacing=0]:last:rounded-b-md",
			switch size {
			case "sm":
				"h-8 min-w-8 px-2.5 has-[>svg]:px-1.5 group-data-[spacing=0]/toggle-group:has-[>svg]:px-1.5"
			case "lg":
				"h-10 min-w-10 px-2.5 has-[>svg]:px-2 group-data-[spacing=0]/toggle-group:has-[>svg]:px-2"
			default:
				"h-9 min-w-9 px-2.5 has-[>svg]:px-2 group-data-[spacing=0]/toggle-group:has-[>svg]:px-2"
			},
			switch sp {
			case "2":
				"isolate"
			default:
				"rounded-none shadow-none data-[variant=outline]:border-s-0 data-[variant=outline]:first:border-s"
			},
			switch orientation { default: "data-[spacing=0]:first:rounded-s-lg data-[spacing=0]:last:rounded-e-lg" }
		}
		{ attrs... }
		data-gsxui-slot-toggle-group-item
		data-gsxui-slot-toggle
	>
		{ children }
	</button>
}
