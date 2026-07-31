package sheet

import "github.com/gsxhq/gsxui/ui"

// Directions mirrors shadcn's sheet-side demo: every side uses the same
// complete profile-edit composition so direction changes only the panel's
// anchoring and motion, not the quality or spacing of its contents.
component Directions() {
	<div class="grid grid-cols-2 gap-2">
		{ for _, side := range []string{"top", "right", "bottom", "left"} {
			<ui.Sheet>
				<ui.Button
					variant="outline"
					data-gsxui-slot-sheet-trigger
					aria-haspopup="dialog"
					aria-expanded="false"
					class="capitalize"
				>
					{ side }
				</ui.Button>
				<ui.SheetContent side={side}>
					<ui.SheetHeader>
						<ui.SheetTitle>Edit profile</ui.SheetTitle>
						<ui.SheetDescription>
							Make changes to your profile here. Click save when you're done.
						</ui.SheetDescription>
					</ui.SheetHeader>
					<div class="grid gap-4 px-4">
						<div class="grid grid-cols-4 items-center gap-4">
							<ui.Label for={"sheet-" + side + "-name"} class="text-right">Name</ui.Label>
							<ui.Input id={"sheet-" + side + "-name"} value="Pedro Duarte" class="col-span-3"/>
						</div>
						<div class="grid grid-cols-4 items-center gap-4">
							<ui.Label for={"sheet-" + side + "-username"} class="text-right">Username</ui.Label>
							<ui.Input id={"sheet-" + side + "-username"} value="@peduarte" class="col-span-3"/>
						</div>
					</div>
					<ui.SheetFooter>
						<ui.Button data-gsxui-dialog-close>Save changes</ui.Button>
					</ui.SheetFooter>
				</ui.SheetContent>
			</ui.Sheet>
		} }
	</div>
}
