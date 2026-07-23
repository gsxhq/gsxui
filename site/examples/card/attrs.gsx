package card

import "github.com/gsxhq/gsxui/ui"

// Attrs demonstrates class-merge override: class="py-8" replaces Card's
// default py-6 — tailwind-merge resolves same-property conflicts, caller
// wins, same mechanism every gsxui component uses for its fallthrough attrs.
component Attrs() {
	<ui.Card class="py-8">
		<ui.CardHeader>
			<ui.CardTitle>Roomier card</ui.CardTitle>
			<ui.CardDescription>py-8 overrides py-6; the rest of the class list merges untouched.</ui.CardDescription>
		</ui.CardHeader>
		<ui.CardContent>
			<p class="text-sm text-muted-foreground">Every part accepts attrs the same way — this isn't Card-specific.</p>
		</ui.CardContent>
	</ui.Card>
}
