package alert

import (
	"github.com/gsxhq/gsxui/ui"
	uiicon "github.com/gsxhq/gsxui/ui/icon"
)

// Variants renders the default and destructive Alert variants together.
component Variants() {
	<div class="flex flex-col gap-4">
		<ui.Alert>
			<uiicon.CircleCheck/>
			<ui.AlertTitle>Success! Your changes have been saved.</ui.AlertTitle>
			<ui.AlertDescription>This is an alert with icon, title, and description.</ui.AlertDescription>
		</ui.Alert>
		<ui.Alert variant="destructive">
			<uiicon.CircleAlert/>
			<ui.AlertTitle>Unable to process your payment.</ui.AlertTitle>
			<ui.AlertDescription>Please verify your billing information and try again.</ui.AlertDescription>
		</ui.Alert>
	</div>
}
