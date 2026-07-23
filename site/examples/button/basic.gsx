// Package button holds the site's example gsx components for ui/button.
// Each example is a real, compiled gsx component — the exact source below
// is what the component page displays AND what it renders, so source shown
// is source run; the examples_test.go drift test enforces they can't
// diverge.
package button

import uibutton "github.com/gsxhq/gsxui/ui/button"

// Basic renders the default Button.
component Basic() {
	<uibutton.Button>Button</uibutton.Button>
}
