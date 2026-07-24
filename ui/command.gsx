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
		data-slot="command"
		data-gsxui-command
		class="flex h-full w-full flex-col overflow-hidden rounded-xl bg-popover p-1 text-popover-foreground"
		{ attrs... }
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
	<Dialog data-slot="command-dialog">
		<DialogContent
			data-gsxui-command-dialog
			class="overflow-hidden p-0"
			{ attrs... }
		>
			<DialogHeader class="sr-only">
				<DialogTitle>{ title |> default("Command Palette") }</DialogTitle>
				<DialogDescription>{ description |> default("Search for a command to run...") }</DialogDescription>
			</DialogHeader>
			<Command class="**:data-[slot=command-input-wrapper]:h-12 [&_[data-slot=command-group-heading]]:px-2 [&_[data-slot=command-group-heading]]:font-medium [&_[data-slot=command-group-heading]]:text-muted-foreground [&_[data-slot=command-group]]:px-2 [&_[data-slot=command-group]:not([hidden])_~[data-slot=command-group]]:pt-0 [&_[data-slot=command-input-wrapper]_svg]:h-5 [&_[data-slot=command-input-wrapper]_svg]:w-5 [&_[data-slot=command-input]]:h-12 [&_[data-slot=command-item]]:px-2 [&_[data-slot=command-item]]:py-3 [&_[data-slot=command-item]_svg]:h-5 [&_[data-slot=command-item]_svg]:w-5">
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
	<div data-slot="command-input-wrapper" class="flex h-9 items-center gap-2 border-b px-3">
		<icon.Search class="size-4 shrink-0 opacity-50"/>
		<input
			data-slot="command-input"
			data-gsxui-command-input
			type="text"
			role="combobox"
			aria-expanded="true"
			aria-autocomplete="list"
			autocomplete="off"
			spellcheck="false"
			placeholder={placeholder}
			class="flex h-10 w-full rounded-md bg-transparent py-3 text-sm outline-hidden placeholder:text-muted-foreground disabled:cursor-not-allowed disabled:opacity-50"
			{ attrs... }
		/>
	</div>
}

component CommandList(children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="command-list"
		data-gsxui-command-list
		role="listbox"
		class="max-h-72 scroll-py-1 overflow-x-hidden overflow-y-auto"
		{ attrs... }
	>
		{ children }
	</div>
}

// CommandEmpty is server-rendered hidden; command.js reveals it when a
// query matches nothing (cmdk's Empty renders conditionally — same net
// visual, inverted mechanism since there is no VDOM to unmount).
component CommandEmpty(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="command-empty" hidden class="py-6 text-center text-sm" { attrs... }>{ children }</div>
}

// CommandGroup's heading is a real child div (slot command-group-heading)
// rather than cmdk's heading prop + [cmdk-group-heading] runtime stamp —
// the classes shadcn applies through the group's arbitrary selectors land
// on it via the mapped data-slot selectors (see Command's doc comment).
component CommandGroup(heading string, children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="command-group" role="group" class="overflow-hidden p-1 text-foreground" { attrs... }>
		{ if heading != "" {
			<div data-slot="command-group-heading" class="px-2 py-1.5 text-xs font-medium text-muted-foreground">{ heading }</div>
		} }
		{ children }
	</div>
}

component CommandSeparator(attrs gsx.Attrs) {
	<div data-slot="command-separator" role="separator" class="-mx-1 h-px bg-border" { attrs... }></div>
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
		data-slot="command-item"
		data-gsxui-command-item
		data-value={value}
		role="option"
		aria-selected="false"
		class="relative flex cursor-default items-center gap-2 rounded-sm px-2 py-1.5 text-sm outline-hidden select-none data-[disabled=true]:pointer-events-none data-[disabled=true]:opacity-50 data-[selected=true]:bg-accent data-[selected=true]:text-accent-foreground [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 [&_svg:not([class*='text-'])]:text-muted-foreground"
		{ attrs... }
	>
		{ children }
	</div>
}

component CommandShortcut(children gsx.Node, attrs gsx.Attrs) {
	<span data-slot="command-shortcut" class="ml-auto text-xs tracking-widest text-muted-foreground" { attrs... }>
		{ children }
	</span>
}
