// Package textarea holds the site's example gsx components for ui/textarea.
package textarea

import uitextarea "github.com/gsxhq/gsxui/ui/textarea"

// Basic renders a Textarea with placeholder text.
component Basic() {
	<uitextarea.Textarea value="" placeholder="Type your message here."/>
}
