package luma

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
// data-gsxui-slot-select hook select.js scopes every trigger/content lookup to
// (closest("[data-gsxui-slot-select]"), the same proximity wiring as dropdown's
// data-gsxui-slot-dropdown-menu). When name != "" it also renders the hidden native
// <select> form bridge (see the file header). required/disabled/form mirror
// ui.NativeSelect's own form params so the two components share an
// option-authoring shape.
component Select(name string, required bool, disabled bool, form string, children gsx.Node, attrs gsx.Attrs) {
	<div class={ "contents" } { attrs... } data-gsxui-slot-select>
		{ children }
		{ if name != "" {
			<select
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
		role="combobox"
		aria-expanded="false"
		aria-autocomplete="none"
		data-state="closed"
		data-size={size |> default("default")}
		data-placeholder
		class={
			"bg-input/50 border-transparent data-placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-ring/30 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive dark:aria-invalid:border-destructive/50 gap-1.5 rounded-3xl border py-2 px-3 text-sm transition-[color,box-shadow,background-color] focus-visible:ring-3 aria-invalid:ring-3 [&_svg:not([class*='size-'])]:size-4 flex",
			switch size { case "sm": "h-8" default: "h-9" }
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
	<span class={ "flex gap-1.5 flex-1 text-left" } { attrs... } data-gsxui-slot-select-value>{ placeholder }</span>
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
		popover="auto"
		role="listbox"
		tabindex="-1"
		data-state="closed"
		data-side="bottom"
		class={
			"transition-none z-50 max-h-[min(--spacing(96),var(--gsxui-available-height,9999px))] overflow-x-hidden overflow-y-auto p-1 bg-popover text-popover-foreground data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 ring-foreground/5 dark:ring-foreground/10 min-w-36 rounded-3xl shadow-lg ring-1 duration-100"
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
	<div role="group" { attrs... } data-gsxui-slot-select-group>{ children }</div>
}

// SelectLabel is the group heading (select.js stamps its id and the group's
// aria-labelledby at init).
component SelectLabel(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "text-muted-foreground px-3 py-2.5 text-xs" } { attrs... } data-gsxui-slot-select-label>
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
			"focus:bg-accent focus:text-accent-foreground not-data-[variant=destructive]:focus:**:text-accent-foreground gap-2.5 rounded-2xl py-2 pr-8 pl-3 text-sm font-medium [&_svg:not([class*='size-'])]:size-4 *:[span]:last:flex *:[span]:last:items-center *:[span]:last:gap-2 flex data-[state=checked]:[&>[data-gsxui-slot-select-item-indicator]]:flex"
		}
		{ attrs... }
		data-gsxui-slot-select-item
	>
		<span
			class={ "pointer-events-none absolute right-2 hidden size-4 items-center justify-center" }
			data-gsxui-slot-select-item-indicator
		>
			<icon.Check/>
		</span>
		<span class={ "flex flex-1 gap-2" } data-gsxui-slot-select-item-text>{ children }</span>
	</div>
}

// SelectSeparator divides groups. aria-hidden per Radix's own SelectSeparator
// (a decorative rule, not a role="separator" like DropdownMenuSeparator).
component SelectSeparator(attrs gsx.Attrs) {
	<div aria-hidden="true" class={ "bg-border -mx-1.5 my-1.5 h-px" } { attrs... } data-gsxui-slot-select-separator></div>
}
