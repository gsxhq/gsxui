package combobox

import "github.com/gsxhq/gsxui/ui"

// TriggerClear exercises the one composition no other fixture reaches:
// showTrigger AND showClear together. ComboboxInput's own doc comment states
// the contract — "both true renders both buttons, with the trigger hidden by
// the input group's clear-button descendant condition whenever a clear button
// is also present (clear wins visually)", matching shadcn's own composition
// table (source map ## combobox §2) — but clear.gsx passes only showClear and
// basic.gsx only showTrigger, so nothing in the catalogue actually rendered
// both until this fixture, and the rule that hides the trigger was therefore
// invisible to the computed-style sweep and to every pin.
//
// It was in fact broken: the hiding rule is keyed on the input group
// (`:where([data-gsxui-slot-input-group]):has([…combobox-clear])
// :where([…combobox-trigger])`) and scored (0,0,0), while the trigger renders
// through InputGroupButton -> Button, whose migrated recipe emits a real
// `inline-flex` utility at (0,1,0) in the same layer. Both buttons showed.
// Combobox's own migration takes the selector to (0,3,0) as an arbitrary
// variant on the trigger's recipe class and restores upstream's composition;
// jstest/specs/layer-precedence.spec.ts pins it.
component TriggerClear() {
	<ui.Combobox name="language" value="go">
		<ui.ComboboxInput placeholder="Search language..." showTrigger showClear value="Go" class="w-[220px]"/>
		<ui.ComboboxContent>
			<ui.ComboboxList>
				<ui.ComboboxEmpty>No language found.</ui.ComboboxEmpty>
				<ui.ComboboxItem value="go" selected={true}>Go</ui.ComboboxItem>
				<ui.ComboboxItem value="rust" selected={false}>Rust</ui.ComboboxItem>
				<ui.ComboboxItem value="typescript" selected={false}>TypeScript</ui.ComboboxItem>
			</ui.ComboboxList>
		</ui.ComboboxContent>
	</ui.Combobox>
}
