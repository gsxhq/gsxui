package ui

import (
	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Command and its parts are the shadcn/ui Command (registry/new-york-v4/
// ui/command.tsx). shadcn's version wraps the cmdk React library
// (CommandPrimitive); there is no React here, so the primitive's behavior —
// score-ranked filtering, roving selection that keeps FOCUS in the input
// (aria-activedescendant, not tab focus), Enter/click activation, group
// hiding, DOM reordering by score — is reimplemented in ui/command.js,
// including a faithful port of cmdk's own ranking algorithm (command-score,
// MIT). Markup follows shadcn's data-slot structure with cmdk's private
// [cmdk-*] attribute selectors mapped onto the equivalent data-slot
// selectors (ADAPT — cmdk stamps those attributes at runtime; we own the
// markup, so the slots are the stable hooks). Nova density metrics applied
// per the 2026-07-24 retarget (rounded-xl + p-1 root, max-h-72 list).
//
// GAP: cmdk props not ported — shouldFilter/filter (custom filter fn),
// value/onValueChange (controlled selection), loop. Items opt into
// navigation with a data-href attribute (command.js assigns location);
// anything else listens for the gsxui:select CustomEvent on the item.
component Command(children gsx.Node, attrs gsx.Attrs) {
	<div
		data-gsxui-command
		{ withSlot("command", attrs)... }
	>
		{ children }
	</div>
}

// CommandDialog composes Dialog/DialogContent (so command → dialog derives
// for the CLI and dialog.js's machinery — trigger wiring, Esc, exit
// animation, backdrop — is reused whole). The sr-only header lives INSIDE
// DialogContent, unlike shadcn's outside-the-content placement: our
// dialog.js wireA11y looks the title/description up within the <dialog>
// element to stamp aria-labelledby/-describedby (an ADAPT with identical
// semantics — the text is sr-only either way).
//
// data-gsxui-command-dialog on the content is command.js's global-hotkey
// hook: ⌘K/Ctrl-K toggles the first such dialog on the page.
component CommandDialog(title string, description string, children gsx.Node, attrs gsx.Attrs) {
	<Dialog { withSlot("command-dialog", nil)... }>
		<DialogContent
			data-gsxui-command-dialog
			{ withSlot("command-dialog-content", attrs)... }
		>
			<DialogHeader { withSlot("command-dialog-header", nil)... }>
				<DialogTitle>{ title |> default("Command Palette") }</DialogTitle>
				<DialogDescription>{ description |> default("Search for a command to run...") }</DialogDescription>
			</DialogHeader>
			<Command { withSlot("command-dialog-command", nil)... }>
				{ children }
			</Command>
		</DialogContent>
	</Dialog>
}

// CommandInput renders shadcn's search-icon-plus-input wrapper row. The
// input is the palette's single focus target: command.js filters on input,
// moves selection on ArrowUp/ArrowDown, and activates on Enter, all while
// focus stays here (aria-activedescendant tracks the selected option).
component CommandInput(placeholder string, attrs gsx.Attrs) {
	<div data-gsxui-command-input-wrapper { withSlot("command-input-wrapper", nil)... }>
		<icon.Search/>
		<input
			data-gsxui-command-input
			type="text"
			role="combobox"
			aria-expanded="true"
			aria-autocomplete="list"
			autocomplete="off"
			spellcheck="false"
			placeholder={placeholder}
			{ withSlot("command-input", attrs)... }
		/>
	</div>
}

component CommandList(children gsx.Node, attrs gsx.Attrs) {
	<div
		data-gsxui-command-list
		role="listbox"
		{ withSlot("command-list", attrs)... }
	>
		{ children }
	</div>
}

// CommandEmpty is server-rendered hidden; command.js reveals it when a
// query matches nothing (cmdk's Empty renders conditionally — same net
// visual, inverted mechanism since there is no VDOM to unmount).
component CommandEmpty(children gsx.Node, attrs gsx.Attrs) {
	<div data-gsxui-command-empty hidden { withSlot("command-empty", attrs)... }>{ children }</div>
}

// CommandGroup's heading is a real child div (slot command-group-heading)
// rather than cmdk's heading prop + [cmdk-group-heading] runtime stamp —
// the classes shadcn applies through the group's arbitrary selectors land
// on it via the mapped data-slot selectors (see Command's doc comment).
component CommandGroup(heading string, children gsx.Node, attrs gsx.Attrs) {
	<div data-gsxui-command-group role="group" { withSlot("command-group", attrs)... }>
		{ if heading != "" {
			<div data-gsxui-command-group-heading { withSlot("command-group-heading", nil)... }>{ heading }</div>
		} }
		{ children }
	</div>
}

component CommandSeparator(attrs gsx.Attrs) {
	<div data-gsxui-command-separator role="separator" { withSlot("command-separator", attrs)... }></div>
}

// CommandItem is a role="option" div (cmdk's own role), NOT focusable —
// selection is the data-selected stamp command.js manages, focus never
// leaves the input. value seeds the match text; empty value falls back to
// the item's textContent (cmdk's own default). data-[selected=true] styling
// matches shadcn's tokens. Disable with aria-disabled="true" or a
// data-disabled attribute (skipped by filter, selection, and activation —
// the same contract as DropdownMenuItem).
component CommandItem(value string, children gsx.Node, attrs gsx.Attrs) {
	<div
		data-gsxui-command-item
		data-value={value}
		role="option"
		aria-selected="false"
		{ withSlot("command-item", attrs)... }
	>
		{ children }
	</div>
}

component CommandShortcut(children gsx.Node, attrs gsx.Attrs) {
	<span { withSlot("command-shortcut", attrs)... }>
		{ children }
	</span>
}
