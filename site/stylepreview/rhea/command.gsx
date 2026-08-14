package rhea

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
		class={ "size-full overflow-hidden", "bg-popover text-popover-foreground rounded-3xl p-1 flex flex-col" }
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
			class={ "rounded-3xl p-0 [dialog&]:overflow-hidden top-1/3 translate-y-0" }
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
		class={
			"mx-1 mt-1 pl-2 bg-input/50 border border-transparent h-8 rounded-2xl [&>svg]:size-4 [&>svg]:shrink-0 [&>svg]:opacity-50 flex items-center gap-2"
		}
		data-gsxui-slot-command-input-wrapper
	>
		<icon.Search/>
		<input
			type="text"
			role="combobox"
			aria-expanded="true"
			aria-autocomplete="list"
			autocomplete="off"
			spellcheck="false"
			placeholder={placeholder}
			class={ "w-full text-sm flex" }
			{ attrs... }
			data-gsxui-slot-command-input
		/>
	</div>
}

component CommandList(children gsx.Node, attrs gsx.Attrs) {
	<div
		role="listbox"
		class={
			"overflow-x-hidden overflow-y-auto",
			"[scrollbar-width:none] [&::-webkit-scrollbar]:hidden max-h-72 scroll-py-1 outline-none"
		}
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
	<div hidden class={ "py-6 text-center text-sm" } { attrs... } data-gsxui-slot-command-empty>{ children }</div>
}

// CommandGroup's heading is a real child div (slot command-group-heading)
// rather than cmdk's heading prop + [cmdk-group-heading] runtime stamp —
// the styles shadcn applies through the group's descendant selectors land
// on it via the mapped public slot selectors (see Command's doc comment).
component CommandGroup(heading string, children gsx.Node, attrs gsx.Attrs) {
	<div
		role="group"
		class={
			"text-foreground **:[[cmdk-group-heading]]:text-muted-foreground overflow-hidden p-1 **:[[cmdk-group-heading]]:px-2 **:[[cmdk-group-heading]]:py-1.5 **:[[cmdk-group-heading]]:text-xs **:[[cmdk-group-heading]]:font-medium"
		}
		{ attrs... }
		data-gsxui-slot-command-group
	>
		{ if heading != "" {
			<div class={ "px-2 py-1.5 text-xs font-medium text-muted-foreground" } data-gsxui-slot-command-group-heading>
				{ heading }
			</div>
		} }
		{ children }
	</div>
}

component CommandSeparator(attrs gsx.Attrs) {
	<div role="separator" class={ "bg-border/50 my-1 h-px" } { attrs... } data-gsxui-slot-command-separator></div>
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
		data-value={value}
		role="option"
		aria-selected="false"
		class={
			"group/command-item",
			"data-disabled:pointer-events-none data-selected:bg-muted data-selected:text-foreground data-selected:*:[svg]:text-foreground relative flex cursor-default items-center gap-2 min-h-7 rounded-xl px-2 py-1.5 text-sm outline-hidden select-none [[data-gsxui-slot-dialog-content]_&]:rounded-2xl [&_svg:not([class*='size-'])]:size-4 [&_svg]:pointer-events-none [&_svg]:shrink-0"
		}
		{ attrs... }
		data-gsxui-slot-command-item
	>
		{ children }
	</div>
}

component CommandShortcut(children gsx.Node, attrs gsx.Attrs) {
	<span
		class={ "text-muted-foreground group-data-selected/command-item:text-foreground ml-auto text-xs tracking-widest" }
		{ attrs... }
		data-gsxui-slot-command-shortcut
	>
		{ children }
	</span>
}
