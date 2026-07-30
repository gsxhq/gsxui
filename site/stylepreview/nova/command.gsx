package nova

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
// MIT). Markup follows shadcn's generic slot structure with cmdk's private
// [cmdk-*] attribute selectors mapped onto equivalent public slot
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
		class={ "flex h-full w-full flex-col overflow-hidden rounded-xl bg-popover p-1 text-popover-foreground" }
		{ attrs... }
		data-gsxui-slot-command
	>
		{ children }
	</div>
}

// CommandDialog composes Dialog/DialogContent (so command → dialog derives
// for the CLI and dialog.js's machinery — trigger wiring, Esc, exit
// animation, backdrop — is reused whole). An optional trigger node renders
// inside that same Dialog root; dialog ownership is deliberately nearest-root
// scoped, so callers must not wrap CommandDialog in another Dialog to attach a
// trigger. The sr-only header lives INSIDE DialogContent, unlike shadcn's
// outside-the-content placement: our
// dialog.js wireA11y looks the title/description up within the <dialog>
// element to stamp aria-labelledby/-describedby (an ADAPT with identical
// semantics — the text is sr-only either way).
//
// data-gsxui-command-dialog on the content is command.js's global-hotkey
// hook: ⌘K/Ctrl-K toggles the first such dialog on the page.
component CommandDialog(title string, description string, trigger gsx.Node, children gsx.Node, attrs gsx.Attrs) {
	<Dialog data-gsxui-slot-command-dialog>
		{ trigger }
		<DialogContent
			data-gsxui-command-dialog
			class={ "[dialog&]:overflow-hidden" }
			{ attrs... }
			data-gsxui-slot-command-dialog-content
		>
			<DialogHeader data-gsxui-slot-command-dialog-header>
				<DialogTitle>{ title |> default("Command Palette") }</DialogTitle>
				<DialogDescription>{ description |> default("Search for a command to run...") }</DialogDescription>
			</DialogHeader>
			<Command data-gsxui-slot-command-dialog-command>
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
	<div
		data-gsxui-command-input-wrapper
		class={
			"flex h-9 items-center gap-2 border-b px-3 [&>svg]:size-4 [&>svg]:shrink-0 [&>svg]:opacity-50 [[data-gsxui-slot-command-dialog-content]_&]:h-12 [[data-gsxui-slot-command-dialog-content]_&]:[&_svg]:size-5"
		}
		data-gsxui-slot-command-input-wrapper
	>
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
			class={
				"flex h-10 w-full rounded-md bg-transparent py-3 text-sm outline-hidden placeholder:text-muted-foreground disabled:cursor-not-allowed disabled:opacity-50 [[data-gsxui-slot-command-dialog-content]_&]:h-12"
			}
			{ attrs... }
			data-gsxui-slot-command-input
		/>
	</div>
}

component CommandList(children gsx.Node, attrs gsx.Attrs) {
	<div
		data-gsxui-command-list
		role="listbox"
		class={ "max-h-72 scroll-py-1 overflow-x-hidden overflow-y-auto" }
		{ attrs... }
		data-gsxui-slot-command-list
	>
		{ children }
	</div>
}

// CommandEmpty is server-rendered hidden; command.js reveals it when a
// query matches nothing (cmdk's Empty renders conditionally — same net
// visual, inverted mechanism since there is no VDOM to unmount).
component CommandEmpty(children gsx.Node, attrs gsx.Attrs) {
	<div data-gsxui-command-empty hidden class={ "py-6 text-center text-sm" } { attrs... } data-gsxui-slot-command-empty>
		{ children }
	</div>
}

// CommandGroup's heading is a real child div (slot command-group-heading)
// rather than cmdk's heading prop + [cmdk-group-heading] runtime stamp —
// the styles shadcn applies through the group's descendant selectors land
// on it via the mapped public slot selectors (see Command's doc comment).
component CommandGroup(heading string, children gsx.Node, attrs gsx.Attrs) {
	<div
		data-gsxui-command-group
		role="group"
		class={
			"overflow-hidden p-1 text-foreground [[data-gsxui-slot-command-dialog-content]_&]:px-2 [[data-gsxui-slot-command-dialog-content]_[data-gsxui-slot-command-group]:not([hidden])~&]:pt-0"
		}
		{ attrs... }
		data-gsxui-slot-command-group
	>
		{ if heading != "" {
			<div
				data-gsxui-command-group-heading
				class={ "px-2 py-1.5 text-xs font-medium text-muted-foreground" }
				data-gsxui-slot-command-group-heading
			>
				{ heading }
			</div>
		} }
		{ children }
	</div>
}

component CommandSeparator(attrs gsx.Attrs) {
	<div
		data-gsxui-command-separator
		role="separator"
		class={ "-mx-1 h-px bg-border" }
		{ attrs... }
		data-gsxui-slot-command-separator
	></div>
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
		class={
			"relative flex cursor-default items-center gap-2 rounded-sm px-2 py-1.5 text-sm outline-hidden select-none data-[disabled=true]:pointer-events-none data-[selected=true]:bg-accent data-[selected=true]:text-accent-foreground [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 [&_svg:not([class*='text-'])]:text-muted-foreground [[data-gsxui-slot-command-dialog-content]_&]:px-2 [[data-gsxui-slot-command-dialog-content]_&]:py-3 [[data-gsxui-slot-command-dialog-content]_&]:[&_svg]:size-5"
		}
		{ attrs... }
		data-gsxui-slot-command-item
	>
		{ children }
	</div>
}

component CommandShortcut(children gsx.Node, attrs gsx.Attrs) {
	<span class={ "ml-auto text-xs tracking-widest text-muted-foreground" } { attrs... } data-gsxui-slot-command-shortcut>
		{ children }
	</span>
}
