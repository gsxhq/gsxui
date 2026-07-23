// Package accordion holds the site's example gsx components for
// ui/accordion.
package accordion

import uiaccordion "github.com/gsxhq/gsxui/ui/accordion"

// Basic renders three items sharing name="faq" — native <details name>
// gives the group exclusive-open behavior, no JS. The first item opens.
component Basic() {
	<uiaccordion.Accordion name="faq">
		<uiaccordion.AccordionItem name="faq" open>
			<uiaccordion.AccordionTrigger>Is it accessible?</uiaccordion.AccordionTrigger>
			<uiaccordion.AccordionContent>Yes, it follows the WAI-ARIA design pattern.</uiaccordion.AccordionContent>
		</uiaccordion.AccordionItem>
		<uiaccordion.AccordionItem name="faq">
			<uiaccordion.AccordionTrigger>Is it styled?</uiaccordion.AccordionTrigger>
			<uiaccordion.AccordionContent>Yes, it comes with default styles.</uiaccordion.AccordionContent>
		</uiaccordion.AccordionItem>
		<uiaccordion.AccordionItem name="faq">
			<uiaccordion.AccordionTrigger>Is it animated?</uiaccordion.AccordionTrigger>
			<uiaccordion.AccordionContent>Yes, height is animated by default.</uiaccordion.AccordionContent>
		</uiaccordion.AccordionItem>
	</uiaccordion.Accordion>
}
