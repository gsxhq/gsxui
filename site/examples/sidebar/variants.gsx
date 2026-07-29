package sidebar

import (
	"github.com/gsxhq/gsxui/ui"
	"github.com/gsxhq/gsxui/ui/icon"
)

// brand uses the same overflow behavior as the rest of an icon-collapsible
// menu: the fixed A mark remains visible while the full company name is
// clipped by SidebarMenuButton's compact state.
component brand() {
	<ui.SidebarMenu>
		<ui.SidebarMenuItem>
			<ui.SidebarMenuButton
				size="lg"
				tooltip="Acme Inc"
				aria-label="Acme Inc"
				data-sidebar-example-brand
			>
				<span
					class="flex aspect-square size-8 shrink-0 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground"
					data-sidebar-example-brand-mark
				>
					A
				</span>
				<span class="truncate font-semibold" data-sidebar-example-brand-name>Acme Inc</span>
			</ui.SidebarMenuButton>
		</ui.SidebarMenuItem>
	</ui.SidebarMenu>
}

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

// demo renders one SidebarProvider. The gallery runs every exported case
// below in its own iframe because each desktop Sidebar intentionally owns
// its document viewport; mounting several cases in one document would make
// their fixed containers overlap.
// SidebarRail is only meaningful alongside offcanvas/icon collapsing (it
// positions absolute against the desktop tree's own fixed sidebar-
// container, which collapsible="none"'s flat div never renders — see
// ui/sidebar.gsx's own Sidebar doc comment), so rail gates it off for that
// one case.
component demo(side string, variant string, collapsible string, open bool, rail bool) {
	<ui.SidebarProvider open={open}>
		<ui.Sidebar open={open} side={side} variant={variant} collapsible={collapsible}>
			<ui.SidebarHeader>
				<brand/>
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
}

// Variants demonstrates the three `variant`s (sidebar/floating/inset) and
// the three `collapsible` modes (offcanvas/icon/none) shadcn's own registry
// exposes as separate props on the same component. It is the default case;
// the gallery registers the sibling cases below as named preview documents
// under this same source block.
component Variants() {
	<demo side="" variant="" collapsible="" open={true} rail={true}/>
}

component VariantFloating() {
	<demo side="" variant="floating" collapsible="" open={true} rail={true}/>
}

component VariantInset() {
	<demo side="" variant="inset" collapsible="" open={true} rail={true}/>
}

component CollapsibleOffcanvas() {
	<demo side="" variant="" collapsible="offcanvas" open={true} rail={true}/>
}

component CollapsibleIcon() {
	<demo side="" variant="" collapsible="icon" open={true} rail={true}/>
}

component CollapsibleNone() {
	<demo side="" variant="" collapsible="none" open={true} rail={false}/>
}

component RightCollapsed() {
	<demo side="right" variant="" collapsible="offcanvas" open={false} rail={true}/>
}

component IconCollapsed() {
	<demo side="" variant="" collapsible="icon" open={false} rail={true}/>
}
