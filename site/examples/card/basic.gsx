// Package card holds the site's example gsx components for ui/card.
package card

import uicard "github.com/gsxhq/gsxui/ui/card"

// Basic renders a Card with a header and content.
component Basic() {
	<uicard.Card>
		<uicard.CardHeader>
			<uicard.CardTitle>Create project</uicard.CardTitle>
			<uicard.CardDescription>Deploy your new project in one click.</uicard.CardDescription>
		</uicard.CardHeader>
		<uicard.CardContent>
			<p class="text-sm text-muted-foreground">Choose a template to get started.</p>
		</uicard.CardContent>
	</uicard.Card>
}
