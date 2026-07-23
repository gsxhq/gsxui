// Package inputgroup holds the site's example gsx components for
// ui/input-group.
package inputgroup

import (
	"github.com/gsxhq/gsxui/ui"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Basic renders a search-style InputGroup (leading icon addon +
// InputGroupInput) and a second InputGroup pairing InputGroupInput with a
// trailing InputGroupButton addon.
component Basic() {
	<div class="flex w-full max-w-sm flex-col gap-4">
		<ui.InputGroup>
			<ui.InputGroupAddon>
				<icon.Search class="size-4"/>
			</ui.InputGroupAddon>
			<ui.InputGroupInput placeholder="Search..."/>
		</ui.InputGroup>
		<ui.InputGroup>
			<ui.InputGroupInput placeholder="Email address" type="email"/>
			<ui.InputGroupAddon align="inline-end">
				<ui.InputGroupButton aria-label="Send">
					<icon.Send/>
				</ui.InputGroupButton>
			</ui.InputGroupAddon>
		</ui.InputGroup>
	</div>
}
