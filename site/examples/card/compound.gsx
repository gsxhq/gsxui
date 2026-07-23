package card

import (
	uibadge "github.com/gsxhq/gsxui/ui/badge"
	uibutton "github.com/gsxhq/gsxui/ui/button"
	uicard "github.com/gsxhq/gsxui/ui/card"
)

// Compound composes every Card part: a border-b header with title,
// description and action, content, and a border-t footer.
component Compound() {
	<uicard.Card>
		<uicard.CardHeader class="border-b">
			<uicard.CardTitle>Notifications</uicard.CardTitle>
			<uicard.CardDescription>You have 3 unread messages.</uicard.CardDescription>
			<uicard.CardAction><uibadge.Badge variant="secondary">New</uibadge.Badge></uicard.CardAction>
		</uicard.CardHeader>
		<uicard.CardContent>
			<p class="text-sm text-muted-foreground">Push notifications are enabled for this device.</p>
		</uicard.CardContent>
		<uicard.CardFooter class="border-t">
			<uibutton.Button class="w-full">Mark all as read</uibutton.Button>
		</uicard.CardFooter>
	</uicard.Card>
}
