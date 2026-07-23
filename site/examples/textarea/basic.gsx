// Package textarea holds the site's example gsx components for
// ui/textarea. Each example is a real, compiled gsx component — the exact
// source below is what the component page displays AND what it renders, so
// source shown is source run; the examples_test.go drift test enforces
// they can't diverge.
package textarea

import uitextarea "github.com/gsxhq/gsxui/ui/textarea"

// Basic renders a Textarea with placeholder text.
component Basic() {
	<uitextarea.Textarea value="" placeholder="Type your message here."/>
}
