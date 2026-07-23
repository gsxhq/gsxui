package card

import (
	"github.com/gsxhq/gsxui/ui"
)

// Compound composes every Card part: a border-b header with title,
// description and action, content, and a border-t footer.
component Compound() {
	<ui.Card>
		<ui.CardHeader class="border-b">
			<ui.CardTitle>Notifications</ui.CardTitle>
			<ui.CardDescription>You have 3 unread messages.</ui.CardDescription>
			<ui.CardAction>
				<ui.Badge variant="secondary">New</ui.Badge>
			</ui.CardAction>
		</ui.CardHeader>
		<ui.CardContent>
			<p class="text-sm text-muted-foreground">Push notifications are enabled for this device.</p>
		</ui.CardContent>
		<ui.CardFooter class="border-t">
			<ui.Button class="w-full">Mark all as read</ui.Button>
		</ui.CardFooter>
	</ui.Card>
}
