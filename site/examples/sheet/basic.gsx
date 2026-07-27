// Package sheet holds the site's example gsx components for ui/sheet.
package sheet

import (
	"github.com/gsxhq/gsxui/ui"
)

// Basic mirrors shadcn's default sheet demo: a right-side profile editor
// with header, form fields, and footer. The trigger is one real Button
// carrying the sheet trigger style slot plus Dialog's behavior/ARIA
// contract. The footer button is a real Button carrying
// data-gsxui-dialog-close directly rather than wrapped in SheetClose, for
// the same button-in-button reason: SheetClose renders its own <button>,
// and nesting a real Button inside it hits the exact HTML trap
// ui/sheet.gsx's SheetTrigger doc comment warns about (see also
// ui/dialog.gsx's own DialogFooter, which uses the identical
// <ui.Button data-gsxui-dialog-close> idiom rather than nesting inside
// DialogClose).
component Basic() {
	<ui.Sheet>
		<ui.Button
			variant="outline"
			data-gsxui-dialog-trigger
			data-gsxui-slot-sheet-trigger
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
			<div class="grid gap-4 px-4">
				<div class="grid grid-cols-4 items-center gap-4">
					<ui.Label for="sheet-basic-name" class="text-right">Name</ui.Label>
					<ui.Input id="sheet-basic-name" value="Pedro Duarte" class="col-span-3"/>
				</div>
				<div class="grid grid-cols-4 items-center gap-4">
					<ui.Label for="sheet-basic-username" class="text-right">Username</ui.Label>
					<ui.Input id="sheet-basic-username" value="@peduarte" class="col-span-3"/>
				</div>
			</div>
			<ui.SheetFooter>
				<ui.Button data-gsxui-dialog-close>Save changes</ui.Button>
			</ui.SheetFooter>
		</ui.SheetContent>
	</ui.Sheet>
}
