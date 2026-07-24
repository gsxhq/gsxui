package ui

import (
	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Combobox and its parts are the shadcn/ui Combobox (registry/new-york-v4/
// ui/combobox.tsx) — a Base UI (@base-ui/react) filtered listbox, NOT the
// legacy Command+Popover "combobox" recipe some registries also ship
// (registry/new-york-v4/examples/combobox-*.tsx is that older pattern and is
// NOT this component's demo; the real demos, verified against
// content/docs/components/base/combobox.mdx, live at apps/v4/examples/base/
// combobox-*.tsx — see 2026-07-24 tier4 source map `## combobox`). Twelve
// parts ship here: Combobox, ComboboxInput, ComboboxTrigger, ComboboxClear,
// ComboboxContent, ComboboxList, ComboboxItem, ComboboxGroup, ComboboxLabel,
// ComboboxEmpty, ComboboxSeparator, ComboboxValue.
//
// MECHANISM (open/close + value model, ported from ui/select.gsx/select.js):
// ComboboxContent rides the same popover machinery as SelectContent —
// popover="auto" (top layer, light dismiss, free Esc), position:fixed
// anchoring below the input (ui/combobox.js), the discrete-transition
// enter/exit class block copied BYTE-IDENTICAL from ui/select.gsx (see
// below), and data-state="open" stamped synchronously BEFORE
// showPopover() — the flash-avoidance rule (docs/jsx-parity.md
// ## animations). MECHANISM (form bridge): when name != "", Combobox
// server-renders a real hidden sr-only <input type="text" name value>
// sibling — a non-JS form post, FormData collection, and a <label>'s
// click-through all see an ordinary text control — identical in shape to
// select.gsx's hidden <select> bridge (ui/select.gsx's own header), just a
// plain <input> here since a combobox's value is one string, not an
// enumerated <option> set. combobox.js keeps its .value synced on every
// commit/clear and dispatches a bubbling native "change", exactly like
// select.js:64-66 — Combobox and Select are interchangeable to a consumer
// listening for gsxui:select or the bridge's change event.
//
// ADAPT (filter: Base UI's `contains`, web-verified, NOT read from package
// source): @base-ui/react is absent from this checkout, so the filter
// algorithm can't be traced from dist/ the way select.js's typeahead or
// command.js's commandScore were. Base UI's published docs and issue
// tracker (web-verified, not a source read) name the Combobox default as
// `useFilter().contains` — an Intl.Collator-backed BOOLEAN match at
// sensitivity: "base" (case- and accent-insensitive), ported verbatim in
// ui/combobox.js. It returns a boolean, not a score: there is NO ranking
// and NO reordering, a real divergence from cmdk's command-score fuzzy
// engine (## command's own filter). combobox.js therefore only
// hides/shows items on every keystroke — it does NOT port command.js's
// DOM-reordering pass.
//
// GAP (no shared JS with ui/command.js, registry.Deps reasoning): a JS
// `import` from combobox.js into command.js would be invisible to
// internal/registry.Deps, which derives a component's dependency list by
// go/parser over the generated ui/*.x.go (Go source only) — command.js's
// own commandScore ranking engine has no Go representation to scan, so the
// CLI vendoring path would silently miss it. combobox.js duplicates only
// what it needs (its own `contains` matcher, its own highlight-tracking
// loop) rather than import command.js at all.
//
// GAP (deferred, not built — multi-select/chips and two React-only
// helpers): ComboboxChips, ComboboxChip, ComboboxChipsInput, ChipRemove
// (the multi-select chip-input composition) are not ported — this is a
// single-select port (Combobox's own `value` param is one string).
// ComboboxCollection (a React items-array-to-JSX mapping helper) is not
// needed: a gsx caller just writes an ordinary `for` loop over its own
// slice, calling ComboboxItem per element. useComboboxAnchor (`() =>
// React.useRef<HTMLDivElement|null>(null)`, existing only so a multi-select
// caller can share one ref between ComboboxChips and ComboboxContent) has
// no gsx equivalent to build — a gsx caller passes an anchor selector
// string instead, when/if chips ship.
//
// ADAPT (Content/List drop the anchor-width/available-height CSS vars):
// Base UI's Positioner primitive sets --anchor-width/--available-width/
// --available-height/--transform-origin custom properties at runtime for
// ComboboxContent/List to consume (registry/new-york-v4/ui/combobox.tsx's
// own w-(--anchor-width), max-w-(--available-width),
// origin-(--transform-origin), and List's
// calc(var(--available-height)-...) tokens). combobox.js implements fixed
// positioning below the input instead (ui/select.js's own openContent
// model, not a full anchor-positioning engine), so those vars are never
// set; the tokens that depend on them are dropped in favor of plain
// nova-retargeted values (min-w-36 origin-top-left on Content, a flat
// max-h-72 on List) — the exact same substitution ## select's own ADAPT
// entry documents for SelectContent's min-w-[8rem]->min-w-36. duration-100
// (nova's one-off value for this component) is dropped the same way ##
// select's own entry rejects it: duration-150 is the popover family's
// shared standard, supplied by the discrete-transition block below.
component Combobox(name string, value string, children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="combobox" data-gsxui-combobox class="contents" { attrs... }>
		{ children }
		{ if name != "" {
			<input
				data-gsxui-combobox-bridge
				aria-hidden="true"
				tabindex="-1"
				class="sr-only"
				type="text"
				name={name}
				value={value}
			/>
		} }
	</div>
}

// ComboboxInput composes ui.InputGroup/ui.InputGroupInput/ui.InputGroupAddon/
// ui.InputGroupButton directly rather than calling ComboboxTrigger/
// ComboboxClear as nested components: shadcn's own source polymorphically
// merges ComboboxTrigger's props onto InputGroupButton's own <button> via
// Base UI's `render` prop (one DOM button, two logical components layered
// on it) — gsx has no clone-element merge, so nesting two real
// <button>/<InputGroupButton> calls would emit an invalid button-in-button.
// This inlines the same single button shadcn's merge produces instead,
// stamping combobox-trigger's/combobox-clear's own data-slot via the
// explicit-non-parameter-attribute override mechanism (the same one
// ui/input-group.gsx's InputGroupButton doc comment documents for
// data-size). The standalone ComboboxTrigger/ComboboxClear parts below
// still exist for direct use (the "popup" composition variant, out of this
// port's shipped examples) — see their own doc comments.
//
// data-slot="input-group-input" (not InputGroupInput's own default
// "input-group-control") is this port's own naming call for the actual
// text control, since ARIA anatomy here has no traced Base UI source to
// copy a slot name from (`derived-not-read`) — the underlying
// data-slot="input-group-control" InputGroupInput itself stamps is
// overridden the same way. showTrigger/showClear are independent booleans:
// both true renders both buttons, with the trigger hidden via
// group-has-data-[slot=combobox-clear]/input-group:hidden whenever a clear
// button is also present (clear wins visually), matching shadcn's own
// composition table exactly (source map ## combobox §2).
//
// ADAPT (attrs split class vs. everything else, matching source's own
// `{className, ...props}` destructure): the source sends `className` into
// `cn("w-auto", className)` on InputGroup, and spreads every OTHER extra
// prop (`...props`) onto the input itself. gsx's attrs bag doesn't
// destructure on its own, so this port does it explicitly — the same
// class/rest split ui/native-select.gsx's own wrapper-vs-control doc
// comment documents (attrs.Class() to the wrapper, attrs.Without("class")
// to the control): a caller's `class="w-[220px]"` sizes the InputGroup
// (merged with the "w-auto" default via the usual tailwind-merge
// precedence), while `id`/`aria-invalid`/etc. land on the actual `<input>`,
// e.g. for a <label for> pairing (site/examples/combobox/form.gsx).
component ComboboxInput(placeholder string, showTrigger bool, showClear bool, disabled bool, children gsx.Node, attrs gsx.Attrs) {
	<InputGroup class={ "w-auto", attrs.Class() }>
		<InputGroupInput
			data-slot="input-group-input"
			data-gsxui-combobox-input
			type="text"
			role="combobox"
			aria-expanded="false"
			aria-autocomplete="list"
			autocomplete="off"
			spellcheck="false"
			placeholder={placeholder}
			disabled={disabled}
			{ attrs.Without("class")... }
		/>
		<InputGroupAddon align="inline-end">
			{ if showTrigger {
				<InputGroupButton
					size="icon-xs"
					variant="ghost"
					data-slot="combobox-trigger"
					data-gsxui-combobox-trigger
					class="group-has-data-[slot=combobox-clear]/input-group:hidden data-pressed:bg-transparent [&_svg:not([class*='size-'])]:size-4"
					disabled={disabled}
				>
					<icon.ChevronDown data-slot="combobox-trigger-icon" class="pointer-events-none size-4 text-muted-foreground"/>
				</InputGroupButton>
			} }
			{ if showClear {
				<InputGroupButton
					size="icon-xs"
					variant="ghost"
					data-slot="combobox-clear"
					data-gsxui-combobox-clear
					disabled={disabled}
				>
					<icon.X class="pointer-events-none"/>
				</InputGroupButton>
			} }
		</InputGroupAddon>
		{ children }
	</InputGroup>
}

// ComboboxTrigger is the standalone chevron-toggle button — used directly
// only by the "trigger a popup from a button" composition variant (combobox
// composed with a ui.Button and the search input moved inside
// ComboboxContent instead of ComboboxInput's own addon; out of this port's
// shipped examples, see the GAP list above). ComboboxInput's own addon does
// NOT call this — see its doc comment for why.
component ComboboxTrigger(attrs gsx.Attrs) {
	<button
		type="button"
		data-slot="combobox-trigger"
		data-gsxui-combobox-trigger
		class="[&_svg:not([class*='size-'])]:size-4"
		{ attrs... }
	>
		<icon.ChevronDown data-slot="combobox-trigger-icon" class="pointer-events-none size-4 text-muted-foreground"/>
	</button>
}

// ComboboxClear is the standalone clear button — see ComboboxTrigger's own
// doc comment for why ComboboxInput's addon inlines the equivalent markup
// instead of calling this.
component ComboboxClear(attrs gsx.Attrs) {
	<button
		type="button"
		data-slot="combobox-clear"
		data-gsxui-combobox-clear
		class="[&_svg:not([class*='size-'])]:size-4"
		{ attrs... }
	>
		<icon.X class="pointer-events-none"/>
	</button>
}

// ComboboxContent is the popover listbox surface. See the package doc
// comment for the popover-machinery MECHANISM and the CSS-var-drop ADAPT.
// group/combobox-content is the named group ComboboxEmpty's own
// group-data-empty/combobox-content:flex selector targets;
// combobox.js stamps data-empty on this element (and, separately, on
// ComboboxList, whose own data-empty:p-0 selector is a different
// attribute-presence target) whenever a filter pass leaves zero visible
// items. ring-1 ring-foreground/10 (not border) is the shadcn source's own
// original token here, not a nova border->ring swap to reject — new-york-v4
// never had a border on this part (source map ## combobox §1).
component ComboboxContent(children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="combobox-content"
		data-gsxui-combobox-content
		popover="auto"
		data-state="closed"
		data-side="bottom"
		class={
			"group/combobox-content relative z-50 max-h-72 min-w-36 origin-top-left overflow-hidden rounded-lg bg-popover text-popover-foreground shadow-md ring-1 ring-foreground/10 *:data-[slot=input-group]:m-1 *:data-[slot=input-group]:mb-0 *:data-[slot=input-group]:h-8 *:data-[slot=input-group]:border-input/30 *:data-[slot=input-group]:bg-input/30 *:data-[slot=input-group]:shadow-none",
			"opacity-0 scale-95 transition-[opacity,scale,translate,display,overlay] transition-discrete duration-150 open:opacity-100 open:scale-100 starting:open:opacity-0 starting:open:scale-95",
			"data-[side=bottom]:starting:open:-translate-y-2 data-[side=left]:starting:open:translate-x-2 data-[side=right]:starting:open:-translate-x-2 data-[side=top]:starting:open:translate-y-2"
		}
		{ attrs... }
	>
		{ children }
	</div>
}

// ComboboxList is the scrolling item container (unlike ui.Select, which has
// no separate list part — SelectContent itself scrolls — Combobox's List is
// its own element, matching the raw source's own separate part). role and
// tabindex="-1" implement the WAI-ARIA combobox pattern's listbox half;
// combobox.js never moves real focus here (command.js's model — see the
// package doc comment), only stamps data-highlighted on the current item
// and aria-activedescendant on the INPUT. data-empty:p-0 collapses this
// element's own padding when combobox.js reports zero visible items (a
// second, independent data-empty target from ComboboxContent's own group
// selector — see ComboboxContent's doc comment).
component ComboboxList(children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="combobox-list"
		data-gsxui-combobox-list
		role="listbox"
		tabindex="-1"
		class="no-scrollbar max-h-72 scroll-py-1 overflow-y-auto p-1 data-empty:p-0"
		{ attrs... }
	>
		{ children }
	</div>
}

// ComboboxItem is one option. value is the form/filter value (data-value,
// synced into the hidden bridge on commit); selected server-renders the
// initial pick. Two attributes track two different facts, the same split
// ## select documents for SelectItem: data-state="checked"|"unchecked"
// tracks the VALUE alone and drives the checkmark's CSS visibility
// (group-data-[state=checked]/combobox-item: gating, the same
// mount-gating-in-CSS MECHANISM ## select's own entry names); aria-selected
// tracks the value too here (NOT isValue&&isFocused like Select's own
// aria-selected — Combobox's separate highlight cursor is
// data-highlighted + the input's aria-activedescendant, a materially
// different split since focus never leaves the input). Items are never tab
// stops (command.js's model): no tabindex is stamped at all.
component ComboboxItem(value string, selected bool, children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="combobox-item"
		data-gsxui-combobox-item
		role="option"
		data-value={value}
		{ if selected {
			data-state="checked"
		} else {
			data-state="unchecked"
		} }
		aria-selected={selected}
		class="group/combobox-item relative flex w-full cursor-default items-center gap-2 rounded-md py-1 pr-8 pl-1.5 text-sm outline-hidden select-none data-highlighted:bg-accent data-highlighted:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4"
		{ attrs... }
	>
		{ children }
		<span data-slot="combobox-item-indicator" class="pointer-events-none absolute right-2 hidden size-4 items-center justify-center group-data-[state=checked]/combobox-item:flex">
			<icon.Check class="size-4 pointer-coarse:size-5"/>
		</span>
	</div>
}

// ComboboxGroup wraps a set of items under a ComboboxLabel — real structural
// parts, not a native <optgroup> collapse, the same call ## select makes
// for SelectGroup/SelectLabel (a custom listbox can hold an arbitrary
// styled label child; a native <select> can't). combobox.js wires
// aria-labelledby to the contained ComboboxLabel's generated id at init,
// mirroring select.js's own group-labelling MECHANISM exactly.
component ComboboxGroup(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="combobox-group" data-gsxui-combobox-group role="group" { attrs... }>{ children }</div>
}

// ComboboxLabel is the group heading. pointer-coarse: variants are a real
// new (non-metric) addition nova's own CSS carries for this part — ported
// verbatim, not a retarget.
component ComboboxLabel(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="combobox-label" class="px-2 py-1.5 text-xs text-muted-foreground pointer-coarse:px-3 pointer-coarse:py-2 pointer-coarse:text-sm" { attrs... }>{ children }</div>
}

// ComboboxEmpty is server-rendered hidden; combobox.js reveals it via the
// group-data-empty/combobox-content: selector when a filter pass leaves
// zero visible items — see ComboboxContent's own doc comment for exactly
// which element combobox.js stamps data-empty on.
component ComboboxEmpty(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="combobox-empty" class="hidden w-full justify-center py-2 text-center text-sm text-muted-foreground group-data-empty/combobox-content:flex" { attrs... }>{ children }</div>
}

// ComboboxSeparator divides groups. aria-hidden, matching ## select's own
// SelectSeparator (a decorative rule, not a role="separator").
component ComboboxSeparator(attrs gsx.Attrs) {
	<div data-slot="combobox-separator" aria-hidden="true" class="-mx-1 my-1 h-px bg-border" { attrs... }></div>
}

// ComboboxValue is a plain display slot for the current value/label —
// used by the "trigger a popup from a button" composition variant
// (ComboboxTrigger wraps a ui.Button, ComboboxValue supplies its visible
// text) rather than ComboboxInput's own addon, which reads/writes the
// input's own .value directly instead. No placeholder param (unlike
// ui.SelectValue): children is caller-supplied initial content;
// combobox.js does not touch this element at all in the shipped
// ComboboxInput composition.
component ComboboxValue(children gsx.Node, attrs gsx.Attrs) {
	<span data-slot="combobox-value" { attrs... }>{ children }</span>
}
