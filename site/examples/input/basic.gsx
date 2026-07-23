// Package input holds the site's example gsx components for ui/input.
package input

import uiinput "github.com/gsxhq/gsxui/ui/input"

// Basic renders a default Input.
component Basic() {
	<uiinput.Input type="email" placeholder="you@example.com"/>
}
