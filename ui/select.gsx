package ui

// select.gsx backs the shadcn/ui Select — the custom Radix listbox (NOT the
// styled native <select>, which ships separately as ui.NativeSelect). It is
// a real listbox built on the SAME native-popover machinery ui.DropdownMenu
// uses (popover="auto" top layer, light dismiss, free Esc, closest()
// proximity wiring, the sync-data-state-before-showPopover flash fix, and the
// hover-is-focus / arrow-key .focus() roving-focus idioms) — see ui/select.js
// and docs/jsx-parity.md ## select. What select.js adds on top of dropdown's
// machinery: a value model (one checked item per root, trigger text update,
// data-placeholder clearing), bespoke 1000ms prefix typeahead (startsWith +
// same-char-repeat cycling; works on the closed trigger too), and a hidden
// native <select> form bridge.
//
// Icon deps: SelectTrigger's chevron reuses icon.ChevronDown; SelectItem's
// check indicator uses icon.Check — this import is the select -> icon edge
// internal/registry derives and registry_test.go pins.
//
// FORM BRIDGE (ADAPT + GAP): when name != "", Select server-renders a real
// visually hidden <select aria-hidden tabindex="-1"> sibling carrying
// name/required/disabled/form, so native form submission / FormData / a
// <label>'s click-through / autofill all see an ordinary working <select>.
// gsx has no context to collect the SelectItem values into it at render time,
// so the server renders the bridge with only a synthetic empty <option>;
// select.js fills in one <option> per DOM item at init (module load, before
// any interaction) and keeps .value synced on every selection. GAP: a no-JS
// form submit therefore carries only the empty value — Radix's own bridge is
// SSR-populated via React context this port has no equivalent for. Ledgered
// in docs/jsx-parity.md ## select.

import (
	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Select is the listbox root: a layout-neutral div carrying the
// data-gsxui-select hook select.js scopes every trigger/content lookup to
// (closest("[data-gsxui-select]"), the same proximity wiring as dropdown's
// data-gsxui-dropdown). When name != "" it also renders the hidden native
// <select> form bridge (see the file header). required/disabled/form mirror
// ui.NativeSelect's own form params so the two components share an
// option-authoring shape.
component Select(name string, required bool, disabled bool, form string, children gsx.Node, attrs gsx.Attrs) {
	<div data-gsxui-select class={ "contents" } { attrs... } data-gsxui-slot-select>
		{ children }
		{ if name != "" {
			<select
				data-gsxui-select-bridge
				aria-hidden="true"
				tabindex="-1"
				name={name}
				required={required}
				disabled={disabled}
				{ if form != "" {
					form={form}
				} }
				data-gsxui-slot-select-bridge
			>
				<option value=""></option>
			</select>
		} }
	</div>
}

// SelectTrigger is the combobox button. It renders the caller's SelectValue
// plus the chevron itself (shadcn's SelectTrigger owns the ChevronDownIcon,
// not the caller). data-size is default|sm (nova metrics: h-8 / h-7 + the sm
// radius override). data-placeholder is server-rendered present (initial
// state = no value, placeholder shown, muted via
// data-[placeholder]:text-muted-foreground); select.js removes it on the
// first selection (and at init if an item is server-rendered checked).
// data-state and aria-expanded start closed/"false"; select.js syncs both,
// aria-controls, and aria-required (copied from the bridge) on open/close.
component SelectTrigger(size string, children gsx.Node, attrs gsx.Attrs) {
	<button
		type="button"
		data-gsxui-select-trigger
		role="combobox"
		aria-expanded="false"
		aria-autocomplete="none"
		data-state="closed"
		data-size={size |> default("default")}
		data-placeholder
		class={
			"flex w-fit items-center justify-between gap-1.5 rounded-lg border border-input bg-transparent py-2 pr-2 pl-2.5 text-sm whitespace-nowrap transition-[color,box-shadow] outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 data-[placeholder]:text-muted-foreground dark:bg-input/30 dark:hover:bg-input/50 dark:aria-invalid:ring-destructive/40 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg]:text-muted-foreground [&_svg:not([class*='size-'])]:size-4 [&>svg]:size-4 [&>svg]:opacity-50",
			switch size { case "sm": "h-7 rounded-[min(var(--radius-md),10px)]" default: "h-8" }
		}
		{ attrs... }
		data-gsxui-slot-select-trigger
	>
		{ children }
		<icon.ChevronDown/>
	</button>
}

// SelectValue displays the current value, or the placeholder when nothing is
// selected. select.js overwrites its text content on selection. The default
// stylesheet applies pointer-events:none.
component SelectValue(placeholder string, attrs gsx.Attrs) {
	<span data-gsxui-select-value class={ "pointer-events-none [[data-gsxui-slot-select-trigger]>&]:line-clamp-1 [[data-gsxui-slot-select-trigger]>&]:flex [[data-gsxui-slot-select-trigger]>&]:items-center [[data-gsxui-slot-select-trigger]>&]:gap-1.5" } { attrs... } data-gsxui-slot-select-value>{ placeholder }</span>
}

// SelectContent is the popover listbox. It rides the exact dropdown-menu.js
// popover machinery: popover="auto" (top layer, light dismiss, free Esc),
// server-rendered data-state="closed" + data-side="bottom" (select.js always
// anchors below the trigger), and the discrete-transition enter/exit block
// ported byte-for-byte from DropdownMenuContent (replacing Radix's
// tw-animate keyframes, per docs/jsx-parity.md ## animations). border is
// kept (no border->ring swap, house convention). No scroll up/down buttons —
// the viewport's own overflow-y-auto scrolls natively (GAP, see the parity
// ledger).
component SelectContent(children gsx.Node, attrs gsx.Attrs) {
	<div
		data-gsxui-select-content
		popover="auto"
		role="listbox"
		tabindex="-1"
		data-state="closed"
		data-side="bottom"
		class={
			"z-50 max-h-96 min-w-36 origin-top-left overflow-x-hidden overflow-y-auto rounded-lg border bg-popover p-1 text-popover-foreground shadow-md opacity-0 scale-95 transition-[opacity,scale,translate,display,overlay] transition-discrete duration-150 [&:popover-open]:opacity-100 [&:popover-open]:scale-100 starting:[&:popover-open]:opacity-0 starting:[&:popover-open]:scale-95 starting:[&:popover-open]:data-[side=bottom]:-translate-y-2 starting:[&:popover-open]:data-[side=left]:translate-x-2 starting:[&:popover-open]:data-[side=right]:-translate-x-2 starting:[&:popover-open]:data-[side=top]:translate-y-2"
		}
		{ attrs... }
		data-gsxui-slot-select-content
	>
		{ children }
	</div>
}

// SelectGroup wraps a set of items under a SelectLabel. Unlike
// ui.NativeSelect's SelectGroup (which collapses onto native <optgroup
// label=...>), the custom listbox can hold an arbitrary styled label child,
// so this ports as a real role="group" div. select.js wires aria-labelledby
// to the contained SelectLabel's generated id at init.
component SelectGroup(children gsx.Node, attrs gsx.Attrs) {
	<div data-gsxui-select-group role="group" { attrs... } data-gsxui-slot-select-group>{ children }</div>
}

// SelectLabel is the group heading (select.js stamps its id and the group's
// aria-labelledby at init).
component SelectLabel(children gsx.Node, attrs gsx.Attrs) {
	<div
		data-gsxui-select-label
		class={ "px-1.5 py-1 text-xs text-muted-foreground" }
		{ attrs... }
		data-gsxui-slot-select-label
	>
		{ children }
	</div>
}

// SelectItem is one option. value is the form value (data-value, synced into
// the hidden bridge); disabled skips it in focus/typeahead/selection;
// selected server-renders the initial value (data-state="checked"). Two
// separate attributes track distinct facts, per the traced Radix contract:
//   - data-state="checked"|"unchecked" tracks the VALUE alone and drives the
//     check indicator's CSS visibility through an ancestor-state gate.
//   - aria-selected is server-rendered "false" and recomputed by select.js
//     as (isValue AND isFocused) on every focus change — an item that IS the
//     value but is not the highlighted item reports aria-selected="false".
// items are always tabindex="-1"; select.js moves real DOM focus among them.
component SelectItem(value string, selected bool, disabled bool, children gsx.Node, attrs gsx.Attrs) {
	<div
		data-gsxui-select-item
		role="option"
		data-value={value}
		{ if selected {
			data-state="checked"
		} else {
			data-state="unchecked"
		} }
		aria-selected="false"
		tabindex="-1"
		{ if disabled {
			data-disabled="true"
			aria-disabled="true"
		} }
		class={
			"relative flex w-full cursor-default items-center gap-1.5 rounded-md py-1 pr-8 pl-1.5 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground data-[disabled=true]:pointer-events-none data-[disabled=true]:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg]:text-muted-foreground [&_svg:not([class*='size-'])]:size-4 data-[state=checked]:[&>[data-gsxui-slot-select-item-indicator]]:flex"
		}
		{ attrs... }
		data-gsxui-slot-select-item
	>
		<span
			class={ "pointer-events-none absolute right-2 hidden size-4 items-center justify-center [&>svg]:size-4" }
			data-gsxui-slot-select-item-indicator
		>
			<icon.Check/>
		</span>
		<span class={ "flex items-center gap-2" } data-gsxui-select-item-text data-gsxui-slot-select-item-text>{ children }</span>
	</div>
}

// SelectSeparator divides groups. aria-hidden per Radix's own SelectSeparator
// (a decorative rule, not a role="separator" like DropdownMenuSeparator).
component SelectSeparator(attrs gsx.Attrs) {
	<div aria-hidden="true" class={ "-mx-1 my-1 h-px bg-border" } { attrs... } data-gsxui-slot-select-separator></div>
}
