package sidebar

import (
	"github.com/gsxhq/gsxui/ui"
	"github.com/gsxhq/gsxui/ui/icon"
)

// menu is the shared three-item menu every mini demo below reuses, so the
// only thing that visibly differs between them is the variant/collapsible
// axis each one is labeled for.
component menu() {
	<ui.SidebarMenu>
		<ui.SidebarMenuItem>
			<ui.SidebarMenuButton isActive={true} tooltip="Home">
				<icon.House/>
				<span>Home</span>
			</ui.SidebarMenuButton>
		</ui.SidebarMenuItem>
		<ui.SidebarMenuItem>
			<ui.SidebarMenuButton tooltip="Inbox">
				<icon.Inbox/>
				<span>Inbox</span>
			</ui.SidebarMenuButton>
		</ui.SidebarMenuItem>
		<ui.SidebarMenuItem>
			<ui.SidebarMenuButton tooltip="Settings">
				<icon.Settings/>
				<span>Settings</span>
			</ui.SidebarMenuButton>
		</ui.SidebarMenuItem>
	</ui.SidebarMenu>
}

// demo wraps one mini SidebarProvider — every variant/collapsible
// combination needs its own provider (state is per-wrapper), so each cell
// below is a fully independent sidebar, not a shared one re-skinned.
// SidebarRail is only meaningful alongside offcanvas/icon collapsing (it
// positions absolute against the desktop tree's own fixed sidebar-
// container, which collapsible="none"'s flat div never renders — see
// ui/sidebar.gsx's own Sidebar doc comment), so rail gates it off for that
// one demo cell.
component demo(label string, side string, variant string, collapsible string, open bool, rail bool) {
	<div>
		<div class="mb-2 text-sm font-medium">{ label }</div>
		<ui.SidebarProvider open={open} class="h-64 min-h-0 overflow-hidden rounded-lg border">
			<ui.Sidebar open={open} side={side} variant={variant} collapsible={collapsible}>
				<ui.SidebarHeader>
					<div class="px-2 py-1 text-sm font-semibold">Acme Inc</div>
				</ui.SidebarHeader>
				<ui.SidebarContent>
					<ui.SidebarGroup>
						<menu/>
					</ui.SidebarGroup>
				</ui.SidebarContent>
				{ if rail {
					<ui.SidebarRail/>
				} }
			</ui.Sidebar>
			<ui.SidebarInset>
				<header class="flex h-12 items-center gap-2 border-b px-4">
					<ui.SidebarTrigger/>
				</header>
			</ui.SidebarInset>
		</ui.SidebarProvider>
	</div>
}

// Variants demonstrates the three `variant`s (sidebar/floating/inset) and
// the three `collapsible` modes (offcanvas/icon/none) shadcn's own registry
// exposes as separate props on the same component — six independent mini
// sidebars, each its own SidebarProvider.
component Variants() {
	<div class="grid grid-cols-1 gap-6 md:grid-cols-2">
		<demo label="variant=sidebar (default)" side="" variant="" collapsible="" open={true} rail={true}/>
		<demo label="variant=floating" side="" variant="floating" collapsible="" open={true} rail={true}/>
		<demo label="variant=inset" side="" variant="inset" collapsible="" open={true} rail={true}/>
		<demo label="collapsible=offcanvas (default)" side="" variant="" collapsible="offcanvas" open={true} rail={true}/>
		<demo label="collapsible=icon" side="" variant="" collapsible="icon" open={true} rail={true}/>
		<demo label="collapsible=none (always expanded, no rail)" side="" variant="" collapsible="none" open={true} rail={false}/>
		<demo label="right side, collapsed" side="right" variant="" collapsible="offcanvas" open={false} rail={true}/>
		<demo label="icon collapsed" side="" variant="" collapsible="icon" open={false} rail={true}/>
	</div>
}
