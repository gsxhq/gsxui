package accordion

import "github.com/gsxhq/gsxui/ui"

// Compact overrides AccordionContent's default pb-4 with a caller pb-2 —
// tailwind-merge drops the base utility, same class-merge story as every
// other component.
component Compact() {
	<ui.Accordion name="compact">
		<ui.AccordionItem name="compact" open>
			<ui.AccordionTrigger>Shipping</ui.AccordionTrigger>
			<ui.AccordionContent class="pb-2">Ships within 2 business days.</ui.AccordionContent>
		</ui.AccordionItem>
		<ui.AccordionItem name="compact">
			<ui.AccordionTrigger>Returns</ui.AccordionTrigger>
			<ui.AccordionContent class="pb-2">30-day return window.</ui.AccordionContent>
		</ui.AccordionItem>
	</ui.Accordion>
}
