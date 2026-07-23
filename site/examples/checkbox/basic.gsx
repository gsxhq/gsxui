// Package checkbox holds the site's example gsx components for
// ui/checkbox. Each example is a real, compiled gsx component — the exact
// source below is what the component page displays AND what it renders, so
// source shown is source run; the examples_test.go drift test enforces
// they can't diverge.
package checkbox

import (
	uicheckbox "github.com/gsxhq/gsxui/ui/checkbox"
	uilabel "github.com/gsxhq/gsxui/ui/label"
)

// Basic pairs a Checkbox with a Label via matching id/for.
component Basic() {
	<div class="flex items-center gap-2">
		<uicheckbox.Checkbox id="checkbox-basic-marketing"/>
		<uilabel.Label for="checkbox-basic-marketing">Send me marketing emails</uilabel.Label>
	</div>
}
