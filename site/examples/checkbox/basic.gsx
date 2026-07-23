// Package checkbox holds the site's example gsx components for ui/checkbox.
package checkbox

import (
	"github.com/gsxhq/gsxui/ui"
)

// Basic pairs a Checkbox with a Label via matching id/for.
component Basic() {
	<div class="flex items-center gap-2">
		<ui.Checkbox id="checkbox-basic-marketing"/>
		<ui.Label for="checkbox-basic-marketing">Send me marketing emails</ui.Label>
	</div>
}
