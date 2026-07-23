package ui

// Package selectbox backs the shadcn/ui Select component. It cannot be
// named "select" — that is a reserved Go keyword, illegal as a package
// name AND as an import alias (there is no way to write `import select
// "..."` either), mirroring switchctl's own precedent for "switch". The
// registry/CLI-facing component name stays "selectbox" — Components()/
// Deps() walk directory names, and there is no separate "logical name"
// layer (see docs/jsx-parity.md switch entry). Component identifiers stay
// Select/SelectOption/SelectGroup.

import (
	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Select is the shadcn/ui Select, ported (ADAPT, native-select-v1,
// prominent) as a styled native <select>: form-native, mobile-superior,
// zero JS. shadcn's custom listbox (Trigger/Content/Item/portal machinery
// on top of Radix's SelectPrimitive) is post-v1 backlog; the shadcn *look*
// comes from porting SelectTrigger's classes onto the real <select>
// element, minus the Radix-only/dead-selector tokens ledgered in
// docs/jsx-parity.md. The chevron is rendered via ui/icon (icon.ChevronDown)
// — this import is the selectbox → icon dependency internal/registry
// derives and internal/registry/registry_test.go pins.
component Select(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="select" class="relative">
		<select
			data-slot="select-trigger"
			class="flex w-fit items-center justify-between gap-2 rounded-md border border-input bg-transparent px-3 py-2 text-sm whitespace-nowrap shadow-xs transition-[color,box-shadow] outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 h-9 dark:bg-input/30 dark:hover:bg-input/50 dark:aria-invalid:ring-destructive/40 appearance-none pr-8"
			{ attrs... }
		>{ children }</select>
		<icon.ChevronDown class="pointer-events-none absolute right-3 top-1/2 size-4 -translate-y-1/2 opacity-50"/>
	</div>
}

// SelectOption is a native <option>. selected/disabled are HTML boolean
// attributes (gsx.IsBooleanAttr classifies both "selected" and "disabled"):
// zero value (false) renders absent, matching browser selectedness/disabled
// truth — no data-state plumbing needed, unlike Radix's SelectItem.
component SelectOption(value string, selected bool, disabled bool, children gsx.Node, attrs gsx.Attrs) {
	<option value={ value } selected={ selected } disabled={ disabled } { attrs... }>{ children }</option>
}

// SelectGroup is a native <optgroup>. shadcn's separate SelectGroup
// (wrapper) + SelectLabel (child text) collapse into the one native element
// that already carries a label as an attribute (ADAPT — see
// docs/jsx-parity.md): <optgroup> has no equivalent of an arbitrary label
// child, only the label attribute, so there is nothing to port SelectLabel's
// own class string onto.
component SelectGroup(label string, children gsx.Node, attrs gsx.Attrs) {
	<optgroup label={ label } { attrs... }>{ children }</optgroup>
}
