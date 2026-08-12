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
	<div
		class={
			"group/native-select",
			"relative w-fit [&>svg]:pointer-events-none [&>svg]:absolute [&>svg]:opacity-50 [&>svg]:text-muted-foreground [&>svg]:top-1/2 [&>svg]:right-2.5 [&>svg]:size-4 [&>svg]:-translate-y-1/2 has-[select:disabled]:opacity-50",
			attrs.Class()
		}
		data-gsxui-slot-native-select-wrapper
	>
		<select
			class={
				"bg-input/50 border-transparent placeholder:text-muted-foreground selection:bg-primary selection:text-primary-foreground focus-visible:border-ring focus-visible:ring-ring/30 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive dark:aria-invalid:border-destructive/50 h-9 w-full min-w-0 appearance-none rounded-3xl border py-1 pr-8 pl-3 text-sm transition-[color,box-shadow,background-color] select-none focus-visible:ring-3 aria-invalid:ring-3 data-[size=sm]:h-8 outline-none disabled:pointer-events-none disabled:cursor-not-allowed"
			}
			{ attrs.Without("class")... }
			data-gsxui-slot-native-select
		>
			{ children }
		</select>
		<icon.ChevronDown/>
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
