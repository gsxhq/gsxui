// Package textarea holds the site's example gsx components for ui/textarea.
package textarea

import "github.com/gsxhq/gsxui/ui"

// Basic renders a Textarea with placeholder text.
component Basic() {
	<ui.Textarea value="" placeholder="Type your message here."/>
}
