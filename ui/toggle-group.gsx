package ui

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
// exactly one item (the pressed one for type="single"; the first
// non-disabled item if type="multiple" or nothing is pressed) and
// tabindex="-1" on the rest at init, then maintains that invariant on every
// arrow move and click — see its own header comment.
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
		data-slot="toggle-group"
		data-gsxui-toggle-group
		data-variant={variant |> default("default")}
		data-size={size |> default("default")}
		data-spacing={sp}
		data-orientation="horizontal"
		role={role}
		style=css`--gap: @{sp}`
		class="group/toggle-group flex w-fit items-center gap-[--spacing(var(--gap))] rounded-lg data-[size=sm]:rounded-[min(var(--radius-md),10px)]"
		{ attrs... }
	>
		{ children }
	</div>
}

// ToggleGroupItem composes toggle.gsx's own toggleBase/toggleVariantClass/
// toggleSizeClass — the identical nova-retargeted toggleVariants(variant,
// size) computation Toggle itself uses, not a re-derivation — plus the
// toggle-group-item-only class additions (join-pill layout: w-auto px-3,
// spacing=0 squared-off corners, spacing=0 outline hairline collapse,
// horizontal end-corner rounding). Both type="single" and type="multiple"
// share data-state="on"|"off" (what the shared data-[state=on]:bg-accent
// selector in toggleBase keys off); only the ARIA attribute pair differs —
// see the package doc comment above ToggleGroup.
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
	}}
	<button
		type="button"
		data-slot="toggle-group-item"
		data-gsxui-toggle-group-item
		data-variant={variant |> default("default")}
		data-size={size |> default("default")}
		data-spacing={sp}
		data-orientation="horizontal"
		data-state={state}
		data-value={value}
		{ if groupType == "single" {
			role="radio"
			aria-checked={pressed}
		} else {
			aria-pressed={pressed}
		} }
		class={
			toggleBase,
			toggleVariantClass(variant),
			toggleSizeClass(size),
			"w-auto min-w-0 shrink-0 px-3 focus:z-10 focus-visible:z-10 data-[spacing=0]:rounded-none data-[spacing=0]:shadow-none data-[spacing=0]:data-[variant=outline]:border-l-0 data-[spacing=0]:data-[variant=outline]:first:border-l data-[orientation=horizontal]:data-[spacing=0]:first:rounded-l-lg data-[orientation=horizontal]:data-[spacing=0]:last:rounded-r-lg",
		}
		{ attrs... }
	>
		{ children }
	</button>
}
