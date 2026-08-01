package accordion

import "github.com/gsxhq/gsxui/ui"

// Compact tightens both boxes with caller classes — the trigger's default
// py-2.5 and the content's default pb-2.5 — tailwind-merge drops each base
// utility, same class-merge story as every other component.
component Compact() {
	<ui.Accordion name="compact">
		<ui.AccordionItem name="compact" open>
			<ui.AccordionTrigger class="py-1.5">Shipping</ui.AccordionTrigger>
			<ui.AccordionContent class="pb-1">Ships within 2 business days.</ui.AccordionContent>
		</ui.AccordionItem>
		<ui.AccordionItem name="compact">
			<ui.AccordionTrigger class="py-1.5">Returns</ui.AccordionTrigger>
			<ui.AccordionContent class="pb-1">30-day return window.</ui.AccordionContent>
		</ui.AccordionItem>
	</ui.Accordion>
}
