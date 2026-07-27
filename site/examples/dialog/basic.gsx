// Package dialog holds the site's example gsx components for ui/dialog.
package dialog

import (
	"github.com/gsxhq/gsxui/ui"
)

// Basic renders a confirm dialog. One real Button carries both the dialog
// trigger style slot and the behavior/ARIA contract, so there is no
// button-in-button wrapper or duplicate control.
component Basic() {
	<ui.Dialog>
		<ui.Button
			variant="outline"
			data-gsxui-dialog-trigger
			data-gsxui-slot-dialog-trigger
			aria-haspopup="dialog"
			aria-expanded="false"
		>
			Delete account
		</ui.Button>
		<ui.DialogContent>
			<ui.DialogHeader>
				<ui.DialogTitle>Are you absolutely sure?</ui.DialogTitle>
				<ui.DialogDescription>This will permanently delete your account.</ui.DialogDescription>
			</ui.DialogHeader>
			<ui.DialogFooter>
				<ui.Button variant="outline" data-gsxui-dialog-close>Cancel</ui.Button>
				<ui.Button variant="destructive">Continue</ui.Button>
			</ui.DialogFooter>
		</ui.DialogContent>
	</ui.Dialog>
}
