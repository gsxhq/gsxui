// Package alert holds the site's example gsx components for ui/alert.
package alert

import (
	uialert "github.com/gsxhq/gsxui/ui/alert"
	uiicon "github.com/gsxhq/gsxui/ui/icon"
)

// Basic renders a default Alert with an icon, title, and description.
component Basic() {
	<uialert.Alert>
		<uiicon.Terminal/>
		<uialert.AlertTitle>Heads up!</uialert.AlertTitle>
		<uialert.AlertDescription>You can add components to your app using the CLI.</uialert.AlertDescription>
	</uialert.Alert>
}
