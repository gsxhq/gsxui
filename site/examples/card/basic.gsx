// Package card holds the site's example gsx components for ui/card.
package card

import "github.com/gsxhq/gsxui/ui"

// Basic renders a Card with a header and content.
component Basic() {
	<ui.Card>
		<ui.CardHeader>
			<ui.CardTitle>Create project</ui.CardTitle>
			<ui.CardDescription>Deploy your new project in one click.</ui.CardDescription>
		</ui.CardHeader>
		<ui.CardContent>
			<p class="text-sm text-muted-foreground">Choose a template to get started.</p>
		</ui.CardContent>
	</ui.Card>
}
