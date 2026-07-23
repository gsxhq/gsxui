package alert

import (
	uialert "github.com/gsxhq/gsxui/ui/alert"
	uiicon "github.com/gsxhq/gsxui/ui/icon"
)

// Variants renders the default and destructive Alert variants together.
component Variants() {
	<div class="flex flex-col gap-4">
		<uialert.Alert>
			<uiicon.CircleCheck/>
			<uialert.AlertTitle>Success! Your changes have been saved.</uialert.AlertTitle>
			<uialert.AlertDescription>This is an alert with icon, title, and description.</uialert.AlertDescription>
		</uialert.Alert>
		<uialert.Alert variant="destructive">
			<uiicon.CircleAlert/>
			<uialert.AlertTitle>Unable to process your payment.</uialert.AlertTitle>
			<uialert.AlertDescription>Please verify your billing information and try again.</uialert.AlertDescription>
		</uialert.Alert>
	</div>
}
