// Package sheet holds the site's example gsx components for ui/sheet.
package sheet

import (
	"github.com/gsxhq/gsxui/ui"
)

// Basic mirrors shadcn's own sheet demo shape: a right-side (the default
// side) profile-edit drawer with a header/title/description and a footer
// close button. The trigger is a real Button carrying
// data-gsxui-dialog-trigger — the documented idiom for a styled trigger, no
// SheetTrigger wrapper needed (see ui/sheet.gsx's SheetTrigger doc
// comment). The footer button is a real Button carrying
// data-gsxui-dialog-close directly rather than wrapped in SheetClose, for
// the same button-in-button reason: SheetClose renders its own <button>,
// and nesting a real Button inside it hits the exact HTML trap
// ui/sheet.gsx's SheetTrigger doc comment warns about (see also
// ui/dialog.gsx's own DialogFooter, which uses the identical
// <ui.Button data-gsxui-dialog-close> idiom rather than nesting inside
// DialogClose).
component Basic() {
	<ui.Sheet>
		<ui.Button variant="outline" data-gsxui-dialog-trigger>Edit Profile</ui.Button>
		<ui.SheetContent side="" hideCloseButton={false}>
			<ui.SheetHeader>
				<ui.SheetTitle>Edit profile</ui.SheetTitle>
				<ui.SheetDescription>
					Make changes to your profile here. Click save when you're done.
				</ui.SheetDescription>
			</ui.SheetHeader>
			<ui.SheetFooter>
				<ui.Button variant="outline" data-gsxui-dialog-close>Save changes</ui.Button>
			</ui.SheetFooter>
		</ui.SheetContent>
	</ui.Sheet>
}
