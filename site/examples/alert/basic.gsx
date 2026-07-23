// Package alert holds the site's example gsx components for ui/alert.
package alert

import (
	"github.com/gsxhq/gsxui/ui"
	uiicon "github.com/gsxhq/gsxui/ui/icon"
)

// Basic renders a default Alert with an icon, title, and description.
component Basic() {
	<ui.Alert>
		<uiicon.Terminal/>
		<ui.AlertTitle>Heads up!</ui.AlertTitle>
		<ui.AlertDescription>You can add components to your app using the CLI.</ui.AlertDescription>
	</ui.Alert>
}
