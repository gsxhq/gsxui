// Package drawer holds the site's example gsx components for ui/drawer.
package drawer

import (
	"github.com/gsxhq/gsxui/ui"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Basic adapts shadcn's own drawer-demo.tsx (a bottom-anchored goal-counter
// drawer): the recharts BarChart sparkline has no gsxui equivalent (the
// roadmap defers charting entirely — "chart — defer until demanded"), so
// it is replaced with a static decorative bar row of plain divs standing
// in for the chart's shape, not a literal port. The +/- stepper buttons
// are likewise static markup — shadcn's demo wires them to useState, which
// has no server-rendered gsx equivalent; they render with the demo's own
// styling (ui.Button variant="outline" size="icon" rounded-full) as a
// visual stand-in around the fixed goal number.
//
// The trigger is a real Button carrying data-gsxui-dialog-trigger — the
// documented idiom for a styled trigger, no DrawerTrigger wrapper needed
// (see ui/drawer.gsx's DrawerTrigger doc comment, itself following
// SheetTrigger's own). The footer's Cancel button is a real Button
// carrying data-gsxui-dialog-close directly rather than wrapped in
// DrawerClose, the same button-in-button reasoning as
// site/examples/sheet/basic.gsx's own Save-changes button.
component Basic() {
	<ui.Drawer>
		<ui.Button variant="outline" data-gsxui-dialog-trigger>Open Drawer</ui.Button>
		<ui.DrawerContent direction="">
			<div class="mx-auto w-full max-w-sm">
				<ui.DrawerHeader>
					<ui.DrawerTitle>Move Goal</ui.DrawerTitle>
					<ui.DrawerDescription>Set your daily activity goal.</ui.DrawerDescription>
				</ui.DrawerHeader>
				<div class="p-4 pb-0">
					<div class="flex items-center justify-center gap-2">
						<ui.Button variant="outline" size="icon" class="size-8 shrink-0 rounded-full" aria-label="Decrease goal">
							<icon.Minus/>
						</ui.Button>
						<div class="flex-1 text-center">
							<div class="text-7xl font-bold tracking-tighter">350</div>
							<div class="text-[0.70rem] text-muted-foreground uppercase">Calories/day</div>
						</div>
						<ui.Button variant="outline" size="icon" class="size-8 shrink-0 rounded-full" aria-label="Increase goal">
							<icon.Plus/>
						</ui.Button>
					</div>
					<div class="mt-3 flex h-[120px] items-end gap-1">
						<div class="h-[40%] flex-1 rounded-sm bg-foreground/20"></div>
						<div class="h-[65%] flex-1 rounded-sm bg-foreground/20"></div>
						<div class="h-[45%] flex-1 rounded-sm bg-foreground/20"></div>
						<div class="h-[80%] flex-1 rounded-sm bg-foreground"></div>
						<div class="h-[55%] flex-1 rounded-sm bg-foreground/20"></div>
						<div class="h-[70%] flex-1 rounded-sm bg-foreground/20"></div>
						<div class="h-[35%] flex-1 rounded-sm bg-foreground/20"></div>
						<div class="h-[60%] flex-1 rounded-sm bg-foreground/20"></div>
						<div class="h-[50%] flex-1 rounded-sm bg-foreground/20"></div>
						<div class="h-[75%] flex-1 rounded-sm bg-foreground/20"></div>
						<div class="h-[42%] flex-1 rounded-sm bg-foreground/20"></div>
						<div class="h-[58%] flex-1 rounded-sm bg-foreground/20"></div>
						<div class="h-[38%] flex-1 rounded-sm bg-foreground/20"></div>
					</div>
				</div>
				<ui.DrawerFooter>
					<ui.Button>Submit</ui.Button>
					<ui.Button variant="outline" data-gsxui-dialog-close>Cancel</ui.Button>
				</ui.DrawerFooter>
			</div>
		</ui.DrawerContent>
	</ui.Drawer>
}
