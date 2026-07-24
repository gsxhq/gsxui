package ui

// native-select.gsx backs the shadcn/ui NativeSelect component. Renamed
// 2026-07-24 (was Select/SelectOption/SelectGroup in ui/select.gsx):
// shadcn ships native-select.tsx and select.tsx as permanently-coexisting
// components — the former a styled native <select>, the latter a custom
// Radix-driven listbox — and gsxui now mirrors that split, freeing the
// Select/SelectOption/SelectGroup identifiers for the Tier 3 custom
// listbox port. "select" is also a Go keyword, so a per-component package
// could never have held this file anyway — one of the reasons ui/ is a
// single flat package (see docs/jsx-parity.md packaging entry); as a file
// basename and CLI name (`gsxui add native-select`) it is fine.
// Component identifiers are NativeSelect/NativeSelectOption/NativeSelectGroup.

import (
	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui/icon"
)

// NativeSelect is the shadcn/ui Select, ported (ADAPT, native-select-v1,
// prominent) as a styled native <select>: form-native, mobile-superior,
// zero JS. shadcn's custom listbox (Trigger/Content/Item/portal machinery
// on top of Radix's SelectPrimitive) is post-v1 backlog; the shadcn *look*
// comes from porting SelectTrigger's classes onto the real <select>
// element, minus the Radix-only/dead-selector tokens ledgered in
// docs/jsx-parity.md. The chevron is rendered via ui/icon (icon.ChevronDown)
// — this import is the native-select → icon dependency internal/registry
// derives and internal/registry/registry_test.go pins.
//
// The chevron overlays the <select> from a positioned wrapper (a native
// select can only contain option/optgroup), so the wrapper — not the
// select — must carry the width: it is w-fit (shadcn's trigger default)
// and takes the caller's class (width intent like w-full / w-[180px] maps
// here, where shadcn callers put it on the Trigger), while the select
// fills it with w-full. Putting w-fit on the select inside an unconstrained
// wrapper detaches the absolutely-anchored chevron to the wrapper's far
// edge. Non-class attrs still land on the <select> (name, id, aria-*,
// disabled are form-control concerns).
component NativeSelect(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="native-select-wrapper" class={ "relative w-fit", attrs.Class() }>
		<select
			data-slot="native-select"
			class="flex w-full items-center justify-between gap-2 rounded-lg border border-input bg-transparent pl-2.5 py-1 text-sm whitespace-nowrap transition-[color,box-shadow] outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 h-8 dark:bg-input/30 dark:hover:bg-input/50 dark:aria-invalid:ring-destructive/40 appearance-none pr-8"
			{ attrs.Without("class")... }
		>
			{ children }
		</select>
		<icon.ChevronDown class="pointer-events-none absolute right-2.5 top-1/2 size-4 -translate-y-1/2 opacity-50"/>
	</div>
}

// NativeSelectOption is a native <option>. selected/disabled are HTML
// boolean attributes (gsx.IsBooleanAttr classifies both "selected" and
// "disabled"): zero value (false) renders absent, matching browser
// selectedness/disabled truth — no data-state plumbing needed, unlike
// Radix's SelectItem.
component NativeSelectOption(value string, selected bool, disabled bool, children gsx.Node, attrs gsx.Attrs) {
	<option value={value} selected={selected} disabled={disabled} { attrs... }>{ children }</option>
}

// NativeSelectGroup is a native <optgroup>. shadcn's separate SelectGroup
// (wrapper) + SelectLabel (child text) collapse into the one native element
// that already carries a label as an attribute (ADAPT — see
// docs/jsx-parity.md): <optgroup> has no equivalent of an arbitrary label
// child, only the label attribute, so there is nothing to port SelectLabel's
// own class string onto.
component NativeSelectGroup(label string, children gsx.Node, attrs gsx.Attrs) {
	<optgroup label={label} { attrs... }>{ children }</optgroup>
}
