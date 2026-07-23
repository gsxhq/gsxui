// Package dialog holds the site's example gsx components for ui/dialog.
package dialog

import (
	uibutton "github.com/gsxhq/gsxui/ui/button"
	uidialog "github.com/gsxhq/gsxui/ui/dialog"
)

// Basic renders a confirm dialog. The trigger is a real Button carrying
// data-gsxui-dialog-trigger — the documented idiom for a styled trigger,
// no DialogTrigger wrapper needed (see docs/jsx-parity.md).
component Basic() {
	<uidialog.Dialog>
		<uibutton.Button variant="outline" data-gsxui-dialog-trigger>Delete account</uibutton.Button>
		<uidialog.DialogContent>
			<uidialog.DialogHeader>
				<uidialog.DialogTitle>Are you absolutely sure?</uidialog.DialogTitle>
				<uidialog.DialogDescription>This will permanently delete your account.</uidialog.DialogDescription>
			</uidialog.DialogHeader>
			<uidialog.DialogFooter>
				<uibutton.Button variant="outline" data-gsxui-dialog-close>Cancel</uibutton.Button>
				<uibutton.Button variant="destructive">Continue</uibutton.Button>
			</uidialog.DialogFooter>
		</uidialog.DialogContent>
	</uidialog.Dialog>
}
