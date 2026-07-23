// Package alertdialog holds the site's example gsx components for
// ui/alert-dialog.
package alertdialog

import (
	"github.com/gsxhq/gsxui/ui"
)

// Basic mirrors shadcn's own alert-dialog demo: a destructive-confirm flow
// where Cancel/Action are the only ways out — backdrop click does nothing
// (data-gsxui-dialog-static), Esc still closes. The trigger is a real
// Button carrying data-gsxui-dialog-trigger, the same documented idiom
// ui/dialog's own example uses, rather than the AlertDialogTrigger wrapper
// (see ui/alert-dialog.gsx's AlertDialogTrigger doc comment).
component Basic() {
	<ui.AlertDialog>
		<ui.Button variant="outline" data-gsxui-dialog-trigger>Show dialog</ui.Button>
		<ui.AlertDialogContent>
			<ui.AlertDialogHeader>
				<ui.AlertDialogTitle>Are you absolutely sure?</ui.AlertDialogTitle>
				<ui.AlertDialogDescription>
					This action cannot be undone. This will permanently delete your
					account and remove your data from our servers.
				</ui.AlertDialogDescription>
			</ui.AlertDialogHeader>
			<ui.AlertDialogFooter>
				<ui.AlertDialogCancel>Cancel</ui.AlertDialogCancel>
				<ui.AlertDialogAction>Continue</ui.AlertDialogAction>
			</ui.AlertDialogFooter>
		</ui.AlertDialogContent>
	</ui.AlertDialog>
}
