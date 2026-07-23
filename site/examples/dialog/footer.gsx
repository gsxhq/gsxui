package dialog

import (
	"github.com/gsxhq/gsxui/ui"
)

// Footer shows DialogFooter's showCloseButton: true appends an outline
// Close button automatically, no explicit DialogClose needed.
component Footer() {
	<ui.Dialog>
		<ui.Button data-gsxui-dialog-trigger>Share</ui.Button>
		<ui.DialogContent>
			<ui.DialogHeader>
				<ui.DialogTitle>Share link</ui.DialogTitle>
				<ui.DialogDescription>Anyone with this link can view.</ui.DialogDescription>
			</ui.DialogHeader>
			<ui.DialogFooter showCloseButton>
				<ui.Button>Copy link</ui.Button>
			</ui.DialogFooter>
		</ui.DialogContent>
	</ui.Dialog>
}
