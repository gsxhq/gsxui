// Package empty holds the site's example gsx components for ui/empty.
package empty

import (
	"github.com/gsxhq/gsxui/ui"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Basic renders a realistic empty state: icon media, title, description,
// and a Button action inside EmptyContent.
component Basic() {
	<ui.Empty>
		<ui.EmptyHeader>
			<ui.EmptyMedia variant="icon">
				<icon.Inbox/>
			</ui.EmptyMedia>
			<ui.EmptyTitle>No messages</ui.EmptyTitle>
			<ui.EmptyDescription>
				You're all caught up. New messages will appear here.
			</ui.EmptyDescription>
		</ui.EmptyHeader>
		<ui.EmptyContent>
			<ui.Button>Compose message</ui.Button>
		</ui.EmptyContent>
	</ui.Empty>
}
