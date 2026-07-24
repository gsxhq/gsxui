package drawer

import "github.com/gsxhq/gsxui/ui"

// Directions renders one drawer per anchoring direction, mirroring the
// direction demo on shadcn's own drawer docs. Each drawer is its own
// ui.Drawer root so the trigger-to-dialog proximity wiring stays
// unambiguous. Content is kept minimal — this example exists to exercise
// the four per-direction class strings (positioning, free-edge rounding,
// slide axis, handle-bar visibility, header alignment), not to restate the
// goal-counter demo.
component Directions() {
	<div class="flex flex-wrap gap-2">
		{ for _, dir := range []string{"bottom", "top", "left", "right"} {
			<ui.Drawer>
				<ui.Button variant="outline" data-gsxui-dialog-trigger class="capitalize">{ dir }</ui.Button>
				<ui.DrawerContent direction={dir}>
					<ui.DrawerHeader>
						<ui.DrawerTitle class="capitalize">{ dir } drawer</ui.DrawerTitle>
						<ui.DrawerDescription>Anchored to the { dir } edge.</ui.DrawerDescription>
					</ui.DrawerHeader>
					<div class="p-4 pt-0 text-sm text-muted-foreground">
						The free edge is rounded and the panel slides in from the { dir }.
					</div>
					<ui.DrawerFooter>
						<ui.Button variant="outline" data-gsxui-dialog-close>Close</ui.Button>
					</ui.DrawerFooter>
				</ui.DrawerContent>
			</ui.Drawer>
		} }
	</div>
}
