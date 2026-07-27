// Package sheet holds the site's example gsx components for ui/sheet.
package sheet

import (
	"github.com/gsxhq/gsxui/ui"
)

// Basic mirrors shadcn's own sheet demo shape: a right-side (the default
// side) profile-edit drawer with a header/title/description and a footer
// close button. Each trigger is one real Button carrying the sheet trigger
// style slot plus Dialog's behavior/ARIA contract. The footer button is a
// real Button carrying
// data-gsxui-dialog-close directly rather than wrapped in SheetClose, for
// the same button-in-button reason: SheetClose renders its own <button>,
// and nesting a real Button inside it hits the exact HTML trap
// ui/sheet.gsx's SheetTrigger doc comment warns about (see also
// ui/dialog.gsx's own DialogFooter, which uses the identical
// <ui.Button data-gsxui-dialog-close> idiom rather than nesting inside
// DialogClose).
component Basic() {
	<div class="flex gap-2">
	<ui.Sheet>
		<ui.Button
			variant="outline"
			data-gsxui-dialog-trigger
			data-gsxui-slot="sheet-trigger"
			aria-haspopup="dialog"
			aria-expanded="false"
		>
			Edit Profile
		</ui.Button>
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
	<ui.Sheet>
		<ui.Button
			variant="outline"
			data-gsxui-dialog-trigger
			data-gsxui-slot="sheet-trigger"
			aria-haspopup="dialog"
			aria-expanded="false"
		>
			Open top sheet
		</ui.Button>
		<ui.SheetContent side="top"><ui.SheetTitle>Top sheet</ui.SheetTitle></ui.SheetContent>
	</ui.Sheet>
	<ui.Sheet>
		<ui.Button
			variant="outline"
			data-gsxui-dialog-trigger
			data-gsxui-slot="sheet-trigger"
			aria-haspopup="dialog"
			aria-expanded="false"
		>
			Open bottom sheet
		</ui.Button>
		<ui.SheetContent side="bottom"><ui.SheetTitle>Bottom sheet</ui.SheetTitle></ui.SheetContent>
	</ui.Sheet>
	</div>
}
