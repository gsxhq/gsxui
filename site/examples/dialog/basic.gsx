// Package dialog holds the site's example gsx components for ui/dialog.
package dialog

import (
	"github.com/gsxhq/gsxui/ui"
)

// Basic renders a confirm dialog. The trigger is a real Button carrying
// data-gsxui-dialog-trigger — the documented idiom for a styled trigger,
// no DialogTrigger wrapper needed (see docs/jsx-parity.md).
component Basic() {
	<ui.Dialog>
		<ui.Button variant="outline" data-gsxui-dialog-trigger>Delete account</ui.Button>
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
