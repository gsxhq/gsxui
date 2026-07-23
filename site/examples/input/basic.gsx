// Package input holds the site's example gsx components for ui/input.
// Each example is a real, compiled gsx component — the exact source below
// is what the component page displays AND what it renders, so source shown
// is source run; the examples_test.go drift test enforces they can't
// diverge.
package input

import uiinput "github.com/gsxhq/gsxui/ui/input"

// Basic renders a default Input.
component Basic() {
	<uiinput.Input type="email" placeholder="you@example.com"/>
}
