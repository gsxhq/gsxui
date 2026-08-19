package main

import (
	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

component SidebarContractFrame(open bool, side string, variant string, collapsible string, children gsx.Node, sidebarAttrs gsx.Attrs) {
	<ui.SidebarProvider open={open}>
		<ui.Sidebar
			open={open}
			side={side}
			variant={variant}
			collapsible={collapsible}
			{ sidebarAttrs... }
		>
			{ children }
			<ui.SidebarRail/>
		</ui.Sidebar>
		<ui.SidebarInset>
			<ui.SidebarTrigger/>
			<button type="button" data-gsxui-slot-sidebar-trigger>Custom Sidebar Trigger</button>
		</ui.SidebarInset>
	</ui.SidebarProvider>
}

component SidebarContractMenu() {
	<ui.SidebarHeader>
		<ui.SidebarInput aria-label="Sidebar search"/>
	</ui.SidebarHeader>
	<ui.SidebarContent>
		<ui.SidebarGroup>
			<ui.SidebarGroupLabel tabindex="0">Application</ui.SidebarGroupLabel>
			<ui.SidebarGroupAction aria-label="Add group">+</ui.SidebarGroupAction>
			<ui.SidebarGroupContent>
				<ui.SidebarMenu>
					<ui.SidebarMenuItem data-sidebar-contract="action-item">
						<ui.SidebarMenuButton size="sm" isActive={true}>Small action</ui.SidebarMenuButton>
						<ui.SidebarMenuAction showOnHover={true} aria-label="Small action menu">…</ui.SidebarMenuAction>
						<ui.SidebarMenuBadge>3</ui.SidebarMenuBadge>
					</ui.SidebarMenuItem>
					<ui.SidebarMenuItem>
						<ui.SidebarMenuButton>Default action</ui.SidebarMenuButton>
						<ui.SidebarMenuAction data-sidebar-contract="default-action">…</ui.SidebarMenuAction>
					</ui.SidebarMenuItem>
					<ui.SidebarMenuItem>
						<ui.SidebarMenuButton size="lg">Large badge</ui.SidebarMenuButton>
						<ui.SidebarMenuBadge data-sidebar-contract="large-badge">9</ui.SidebarMenuBadge>
					</ui.SidebarMenuItem>
					<ui.SidebarMenuItem>
						<ui.SidebarMenuButton tooltip="Tooltip item">Tooltip item</ui.SidebarMenuButton>
					</ui.SidebarMenuItem>
					<ui.SidebarMenuItem>
						<ui.SidebarMenuButton
							href="#tooltip-link"
							tooltip="Tooltip link"
							isActive={true}
							data-sidebar-contract="tooltip-link"
						>
							Tooltip link
						</ui.SidebarMenuButton>
					</ui.SidebarMenuItem>
					<ui.SidebarMenuItem>
						<ui.SidebarMenuButton>Parent</ui.SidebarMenuButton>
						<ui.SidebarMenuSub>
							<ui.SidebarMenuSubItem>
								<ui.SidebarMenuSubButton size="sm" isActive={true} href="#">
									Child
								</ui.SidebarMenuSubButton>
							</ui.SidebarMenuSubItem>
						</ui.SidebarMenuSub>
					</ui.SidebarMenuItem>
				</ui.SidebarMenu>
			</ui.SidebarGroupContent>
		</ui.SidebarGroup>
	</ui.SidebarContent>
	<ui.SidebarFooter>Footer</ui.SidebarFooter>
}

component SidebarTokenFixture() {
	<div
		style="--sidebar:rgb(1,2,3);--sidebar-foreground:rgb(4,5,6);--sidebar-primary:rgb(7,8,9);--sidebar-primary-foreground:rgb(10,11,12);--sidebar-accent:rgb(13,14,15);--sidebar-accent-foreground:rgb(16,17,18);--sidebar-border:rgb(19,20,21);--sidebar-ring:rgb(22,23,24)"
	>
		<SidebarContractFrame open={true}>
			<ui.SidebarHeader>
				<ui.SidebarSeparator/>
			</ui.SidebarHeader>
			<ui.SidebarContent>
				<ui.SidebarGroup>
					<ui.SidebarGroupLabel tabindex="0">Tokens</ui.SidebarGroupLabel>
					<ui.SidebarMenu>
						<ui.SidebarMenuItem>
							<ui.SidebarMenuButton isActive={true} data-sidebar-contract="active">
								Active
							</ui.SidebarMenuButton>
							<ui.SidebarMenuBadge data-sidebar-contract="active-badge">1</ui.SidebarMenuBadge>
						</ui.SidebarMenuItem>
						<ui.SidebarMenuItem>
							<ui.SidebarMenuButton data-sidebar-contract="hover">Hover</ui.SidebarMenuButton>
						</ui.SidebarMenuItem>
					</ui.SidebarMenu>
				</ui.SidebarGroup>
			</ui.SidebarContent>
		</SidebarContractFrame>
	</div>
}

component SidebarCallerFixture() {
	<SidebarContractFrame
		open={true}
		sidebarAttrs={gsx.Attrs{{Key: "class", Value: "w-80"}}}
	>
		<ui.SidebarHeader>
			<ui.SidebarInput class="h-12" aria-label="Caller input"/>
		</ui.SidebarHeader>
		<ui.SidebarMenu>
			<ui.SidebarMenuItem>
				<ui.SidebarMenuButton class="h-16 rounded-none">Caller menu button</ui.SidebarMenuButton>
			</ui.SidebarMenuItem>
		</ui.SidebarMenu>
	</SidebarContractFrame>
}

component SidebarContractFixture(name string) {
	{ switch name {
	case "offcanvas-left":
		<SidebarContractFrame open={true} side="left" collapsible="offcanvas">
			<ui.SidebarContent>Left</ui.SidebarContent>
		</SidebarContractFrame>
	case "offcanvas-right":
		<SidebarContractFrame open={false} side="right" collapsible="offcanvas">
			<ui.SidebarContent>Right</ui.SidebarContent>
		</SidebarContractFrame>
	case "none":
		<SidebarContractFrame open={false} side="right" variant="inset" collapsible="none">
			<ui.SidebarContent>Always visible</ui.SidebarContent>
		</SidebarContractFrame>
	case "icon-sidebar":
		<SidebarContractFrame open={false} collapsible="icon">
			<SidebarContractMenu/>
		</SidebarContractFrame>
	case "icon-floating":
		<SidebarContractFrame open={false} variant="floating" collapsible="icon">
			<SidebarContractMenu/>
		</SidebarContractFrame>
	case "inset-collapsed":
		<SidebarContractFrame open={false} variant="inset" collapsible="icon">
			<SidebarContractMenu/>
		</SidebarContractFrame>
	case "menu-icon":
		<SidebarContractFrame open={false} collapsible="icon">
			<SidebarContractMenu/>
		</SidebarContractFrame>
	case "menu-expanded":
		<SidebarContractFrame open={true} collapsible="icon">
			<SidebarContractMenu/>
		</SidebarContractFrame>
	case "tokens":
		<SidebarTokenFixture/>
	case "caller":
		<SidebarCallerFixture/>
	case "mobile":
		<SidebarContractFrame open={true}>
			<SidebarContractMenu/>
		</SidebarContractFrame>
	case "keyboard":
		<div>
			<SidebarContractFrame open={true}>
				<ui.SidebarContent>First</ui.SidebarContent>
			</SidebarContractFrame>
			<SidebarContractFrame open={true}>
				<ui.SidebarContent>Second</ui.SidebarContent>
			</SidebarContractFrame>
			<input aria-label="Typing target"/>
		</div>
	} }
}
