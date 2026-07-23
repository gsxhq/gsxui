package card

import uicard "github.com/gsxhq/gsxui/ui/card"

// Attrs demonstrates class-merge override: class="py-8" replaces Card's
// default py-6 — tailwind-merge resolves same-property conflicts, caller
// wins, same mechanism every gsxui component uses for its fallthrough attrs.
component Attrs() {
	<uicard.Card class="py-8">
		<uicard.CardHeader>
			<uicard.CardTitle>Roomier card</uicard.CardTitle>
			<uicard.CardDescription>py-8 overrides py-6; the rest of the class list merges untouched.</uicard.CardDescription>
		</uicard.CardHeader>
		<uicard.CardContent>
			<p class="text-sm text-muted-foreground">Every part accepts attrs the same way — this isn't Card-specific.</p>
		</uicard.CardContent>
	</uicard.Card>
}
