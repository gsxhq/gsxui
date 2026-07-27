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
		<ui.InputGroup data-disabled="true">
			<ui.InputGroupAddon align="block-start"><ui.InputGroupText>URL</ui.InputGroupText></ui.InputGroupAddon>
			<ui.InputGroupInput disabled placeholder="Disabled"/>
			<ui.InputGroupAddon align="block-end"><ui.InputGroupText>Required</ui.InputGroupText></ui.InputGroupAddon>
		</ui.InputGroup>
		<div class="flex flex-wrap gap-2">
			<ui.InputGroupButton variant="default">Default</ui.InputGroupButton>
			<ui.InputGroupButton variant="destructive">Destructive</ui.InputGroupButton>
			<ui.InputGroupButton variant="outline">Outline</ui.InputGroupButton>
			<ui.InputGroupButton variant="secondary">Secondary</ui.InputGroupButton>
			<ui.InputGroupButton variant="link">Link</ui.InputGroupButton>
			<ui.InputGroupButton size="sm">Small</ui.InputGroupButton>
			<ui.InputGroupButton size="icon-sm" aria-label="Icon"><icon.Send/></ui.InputGroupButton>
		</div>
		<ui.InputGroup>
			<ui.InputGroupInput aria-invalid="true" value="Invalid value"/>
		</ui.InputGroup>
	</div>
}
