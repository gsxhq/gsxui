package accordion

import uiaccordion "github.com/gsxhq/gsxui/ui/accordion"

// Compact overrides AccordionContent's default pb-4 with a caller pb-2 —
// tailwind-merge drops the base utility, same class-merge story as every
// other component.
component Compact() {
	<uiaccordion.Accordion name="compact">
		<uiaccordion.AccordionItem name="compact" open>
			<uiaccordion.AccordionTrigger>Shipping</uiaccordion.AccordionTrigger>
			<uiaccordion.AccordionContent class="pb-2">Ships within 2 business days.</uiaccordion.AccordionContent>
		</uiaccordion.AccordionItem>
		<uiaccordion.AccordionItem name="compact">
			<uiaccordion.AccordionTrigger>Returns</uiaccordion.AccordionTrigger>
			<uiaccordion.AccordionContent class="pb-2">30-day return window.</uiaccordion.AccordionContent>
		</uiaccordion.AccordionItem>
	</uiaccordion.Accordion>
}
