// Package sidebar holds the site's example gsx components for ui/sidebar.
package sidebar

import (
	"github.com/gsxhq/gsxui/ui"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Basic mirrors shadcn's own sidebar-demo.tsx shape: a header, a grouped
// menu (SidebarGroup > SidebarGroupLabel + SidebarMenu), a footer, a
// trigger, and SidebarInset holding the "page" content next to it.
component Basic() {
	<ui.SidebarProvider open={true} class="min-h-[32rem] rounded-lg border">
		<ui.Sidebar open={true}>
			<ui.SidebarHeader>
				<div class="px-2 py-1 text-sm font-semibold">Acme Inc</div>
				<ui.SidebarInput placeholder="Search navigation"/>
			</ui.SidebarHeader>
			<ui.SidebarSeparator/>
			<ui.SidebarContent>
				<ui.SidebarGroup>
					<ui.SidebarGroupLabel>Application</ui.SidebarGroupLabel>
					<ui.SidebarGroupAction aria-label="Add item">
						<icon.Plus/>
					</ui.SidebarGroupAction>
					<ui.SidebarGroupContent>
						<ui.SidebarMenu>
							<ui.SidebarMenuItem>
								<ui.SidebarMenuButton href="#home" isActive={true} tooltip="Home">
									<icon.House/>
									<span>Home</span>
								</ui.SidebarMenuButton>
							</ui.SidebarMenuItem>
							<ui.SidebarMenuItem>
								<ui.SidebarMenuButton variant="outline" size="sm">
									<icon.Inbox/>
									<span>Inbox</span>
								</ui.SidebarMenuButton>
								<ui.SidebarMenuAction showOnHover={true} aria-label="More">
									<icon.Plus/>
								</ui.SidebarMenuAction>
								<ui.SidebarMenuSub>
									<ui.SidebarMenuSubItem>
										<ui.SidebarMenuSubButton size="sm" isActive={true}>Unread</ui.SidebarMenuSubButton>
									</ui.SidebarMenuSubItem>
									<ui.SidebarMenuSubItem>
										<ui.SidebarMenuSubButton size="md">Archived</ui.SidebarMenuSubButton>
									</ui.SidebarMenuSubItem>
								</ui.SidebarMenuSub>
							</ui.SidebarMenuItem>
							<ui.SidebarMenuItem>
								<ui.SidebarMenuButton size="lg">
									<icon.Calendar/>
									<span>Calendar</span>
								</ui.SidebarMenuButton>
								<ui.SidebarMenuBadge>24</ui.SidebarMenuBadge>
							</ui.SidebarMenuItem>
							<ui.SidebarMenuItem>
								<ui.SidebarMenuButton>
									<icon.Search/>
									<span>Search</span>
								</ui.SidebarMenuButton>
								<ui.SidebarMenuAction aria-label="Pin">
									<icon.Plus/>
								</ui.SidebarMenuAction>
							</ui.SidebarMenuItem>
						</ui.SidebarMenu>
						<ui.SidebarMenuSkeleton showIcon={true}/>
						<ui.SidebarMenuSkeleton/>
					</ui.SidebarGroupContent>
				</ui.SidebarGroup>
			</ui.SidebarContent>
			<ui.SidebarFooter>
				<ui.SidebarMenu>
					<ui.SidebarMenuItem>
						<ui.SidebarMenuButton tooltip="Settings">
							<icon.Settings/>
							<span>Settings</span>
						</ui.SidebarMenuButton>
					</ui.SidebarMenuItem>
				</ui.SidebarMenu>
			</ui.SidebarFooter>
			<ui.SidebarRail/>
		</ui.Sidebar>
		<ui.SidebarInset>
			<header class="flex h-12 items-center gap-2 border-b px-4">
				<ui.SidebarTrigger/>
				<span class="text-sm text-muted-foreground">Dashboard</span>
			</header>
			<div class="p-4 text-sm text-muted-foreground">
				Toggle the sidebar with the button above, the rail at its edge, or Cmd/Ctrl+B.
			</div>
		</ui.SidebarInset>
	</ui.SidebarProvider>
}
