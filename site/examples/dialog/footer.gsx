package dialog

import (
	uibutton "github.com/gsxhq/gsxui/ui/button"
	uidialog "github.com/gsxhq/gsxui/ui/dialog"
)

// Footer shows DialogFooter's showCloseButton: true appends an outline
// Close button automatically, no explicit DialogClose needed.
component Footer() {
	<uidialog.Dialog>
		<uibutton.Button data-gsxui-dialog-trigger>Share</uibutton.Button>
		<uidialog.DialogContent>
			<uidialog.DialogHeader>
				<uidialog.DialogTitle>Share link</uidialog.DialogTitle>
				<uidialog.DialogDescription>Anyone with this link can view.</uidialog.DialogDescription>
			</uidialog.DialogHeader>
			<uidialog.DialogFooter showCloseButton>
				<uibutton.Button>Copy link</uibutton.Button>
			</uidialog.DialogFooter>
		</uidialog.DialogContent>
	</uidialog.Dialog>
}
