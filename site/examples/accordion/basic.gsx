// Package accordion holds the site's example gsx components for
// ui/accordion.
package accordion

import "github.com/gsxhq/gsxui/ui"

// Basic renders three items sharing name="faq" — native <details name>
// gives the group exclusive-open behavior, no JS. The first item opens.
component Basic() {
	<ui.Accordion name="faq">
		<ui.AccordionItem name="faq" open>
			<ui.AccordionTrigger>Is it accessible?</ui.AccordionTrigger>
			<ui.AccordionContent>Yes, it follows the WAI-ARIA design pattern.</ui.AccordionContent>
		</ui.AccordionItem>
		<ui.AccordionItem name="faq">
			<ui.AccordionTrigger>Is it styled?</ui.AccordionTrigger>
			<ui.AccordionContent>Yes, it comes with default styles.</ui.AccordionContent>
		</ui.AccordionItem>
		<ui.AccordionItem name="faq">
			<ui.AccordionTrigger>Is it animated?</ui.AccordionTrigger>
			<ui.AccordionContent>Yes, height is animated by default.</ui.AccordionContent>
		</ui.AccordionItem>
	</ui.Accordion>
}
