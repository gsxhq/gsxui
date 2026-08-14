// Package command holds the site's example gsx components for ui/command.
package command

import (
	"github.com/gsxhq/gsxui/ui"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Basic mirrors shadcn's own command demo (registry/new-york-v4/examples/
// command-demo.tsx): an inline (non-dialog) palette with a suggestions
// group, a separator, and a settings group with shortcuts. Type to see
// score-ranked filtering; ArrowUp/ArrowDown move the selection while focus
// stays in the input; Enter or click activates (these demo items emit
// gsxui:select only — they carry no data-href).
component Basic() {
	<div>
		<ui.Command class="max-w-md border shadow-md">
			<ui.CommandInput placeholder="Type a command or search..."/>
			<ui.CommandList>
				<ui.CommandEmpty>No results found.</ui.CommandEmpty>
				<ui.CommandGroup heading="Suggestions">
					<ui.CommandItem>
						<icon.Calendar/>
						<span>Calendar</span>
					</ui.CommandItem>
					<ui.CommandItem>
						<icon.Smile/>
						<span>Search Emoji</span>
					</ui.CommandItem>
					<ui.CommandItem data-disabled="true">
						<icon.Calculator/>
						<span>Calculator</span>
					</ui.CommandItem>
				</ui.CommandGroup>
				<ui.CommandSeparator/>
				<ui.CommandGroup heading="Settings">
					<ui.CommandItem>
						<icon.User/>
						<span>Profile</span>
						<ui.CommandShortcut>⌘P</ui.CommandShortcut>
					</ui.CommandItem>
					<ui.CommandItem>
						<icon.CreditCard/>
						<span>Billing</span>
						<ui.CommandShortcut>⌘B</ui.CommandShortcut>
					</ui.CommandItem>
					<ui.CommandItem>
						<icon.Settings/>
						<span>Settings</span>
						<ui.CommandShortcut>⌘S</ui.CommandShortcut>
					</ui.CommandItem>
				</ui.CommandGroup>
			</ui.CommandList>
		</ui.Command>
		<ui.CommandDialog title="Command palette" description="Search commands">
			<ui.CommandInput placeholder="Search commands..."/>
			<ui.CommandList>
				<ui.CommandItem>Open settings</ui.CommandItem>
			</ui.CommandList>
		</ui.CommandDialog>
	</div>
}
