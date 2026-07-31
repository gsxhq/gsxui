// Package alertdialog holds the site's example gsx components for
// ui/alert-dialog.
package alertdialog

import (
	"github.com/gsxhq/gsxui/ui"
)

// Basic mirrors shadcn's own alert-dialog demo: a destructive-confirm flow
// where Cancel/Action are the only ways out — backdrop click does nothing
// (data-gsxui-dialog-static), Esc still closes. The trigger is one real
// Button carrying both the alert-dialog trigger style slot and Dialog's
// behavior/ARIA contract — gsxui's CSS-only equivalent of shadcn's
// AlertDialogTrigger asChild composition.
component Basic() {
	<ui.AlertDialog>
		<ui.Button
			variant="outline"
			data-gsxui-dialog-trigger
			data-gsxui-slot-alert-dialog-trigger
			aria-haspopup="dialog"
			aria-expanded="false"
		>
			Show dialog
		</ui.Button>
		<ui.AlertDialogContent>
			<ui.AlertDialogHeader>
				<ui.AlertDialogTitle>Are you absolutely sure?</ui.AlertDialogTitle>
				<ui.AlertDialogDescription>
					This action cannot be undone. This will permanently delete your account and remove your data from our servers.
				</ui.AlertDialogDescription>
			</ui.AlertDialogHeader>
			<ui.AlertDialogFooter>
				<ui.AlertDialogCancel>Cancel</ui.AlertDialogCancel>
				<ui.AlertDialogAction>Continue</ui.AlertDialogAction>
			</ui.AlertDialogFooter>
		</ui.AlertDialogContent>
	</ui.AlertDialog>
}
