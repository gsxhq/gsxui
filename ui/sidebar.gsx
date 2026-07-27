package ui

import (
	"math/rand"
	"strconv"

	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Sidebar uses one semantic component tree and two responsive render
// branches. A server cannot know the browser viewport, so non-flat sidebars
// render one native Sheet tree for mobile and one fixed tree for desktop.
// CSS exposes exactly one tree at a time. Consequently children render twice;
// callers must not put IDs inside them.
//
// State is server-rendered from open. JavaScript reflects later changes onto
// SidebarProvider and the desktop root, emits gsxui:change, and opens the
// native mobile dialog. Persistence remains caller-owned.
const (
	sidebarWidth       = "16rem"
	sidebarWidthMobile = "18rem"
	sidebarWidthIcon   = "3rem"
)

component SidebarProvider(open bool, children gsx.Node, attrs gsx.Attrs) {
	{{
		state := "collapsed"
		if open {
			state = "expanded"
		}
	}}
	<div
		data-state={state}
		data-gsxui-sidebar-wrapper
		style={ css`--sidebar-width:@{sidebarWidth}`, css`--sidebar-width-icon:@{sidebarWidthIcon}` }
		{ withSlot("sidebar-wrapper", attrs)... }
	>
		{ children }
	</div>
}

component Sidebar(open bool, side string, variant string, collapsible string, children gsx.Node, attrs gsx.Attrs) {
	{{
		s := side
		if s == "" {
			s = "left"
		}
		v := variant
		if v == "" {
			v = "sidebar"
		}
		c := collapsible
		if c == "" {
			c = "offcanvas"
		}
		state := "collapsed"
		activeCollapsible := c
		if open {
			state = "expanded"
			activeCollapsible = ""
		}
	}}
	{ if collapsible == "none" {
		<div
			data-side={s}
			data-variant={v}
			data-collapsible="none"
			{ withSlot("sidebar", attrs)... }
		>
			{ children }
		</div>
	} else {
		<>
			<Sheet { withSlot("sidebar-mobile-root", attrs)... }>
				<SheetContent
					side={s}
					data-mobile="true"
					data-gsxui-sidebar-mobile-dialog
					style=css`--sidebar-width:@{sidebarWidthMobile}`
					{ withSlot("sidebar", withSlot("sidebar-mobile-content", nil))... }
				>
					<SheetHeader { withSlot("sidebar-mobile-header", nil)... }>
						<SheetTitle { withSlot("sidebar-mobile-title", nil)... }>Sidebar</SheetTitle>
						<SheetDescription { withSlot("sidebar-mobile-description", nil)... }>
							Displays the mobile sidebar.
						</SheetDescription>
					</SheetHeader>
					<div { withSlot("sidebar-mobile-inner", nil)... }>{ children }</div>
				</SheetContent>
			</Sheet>
			<div
				data-state={state}
				data-collapsible={activeCollapsible}
				data-gsxui-sidebar-collapsible={c}
				data-variant={v}
				data-side={s}
				data-gsxui-sidebar-desktop
				{ withSlot("sidebar", withSlot("sidebar-desktop", nil))... }
			>
				<div { withSlot("sidebar-gap", nil)... }></div>
				<div { withSlot("sidebar-container", attrs)... }>
					<div { withSlot("sidebar-inner", nil)... }>
						{ children }
					</div>
				</div>
			</div>
		</>
	} }
}

component SidebarTrigger(attrs gsx.Attrs) {
	<Button
		data-gsxui-sidebar-trigger
		variant="ghost"
		size="icon"
		{ withSlot("sidebar-trigger", attrs)... }
	>
		<icon.PanelLeft/>
		<span { withSlot("sidebar-trigger-label", nil)... }>Toggle Sidebar</span>
	</Button>
}

// SidebarRail is a pointer affordance matching the reference: it is not a
// keyboard tab stop and activates the same state transition as the trigger.
component SidebarRail(attrs gsx.Attrs) {
	<button
		type="button"
		data-gsxui-sidebar-rail
		aria-label="Toggle Sidebar"
		tabindex="-1"
		title="Toggle Sidebar"
		{ withSlot("sidebar-rail", attrs)... }
	></button>
}

component SidebarInset(children gsx.Node, attrs gsx.Attrs) {
	<main { withSlot("sidebar-inset", attrs)... }>{ children }</main>
}

component SidebarInput(attrs gsx.Attrs) {
	<Input { withSlot("sidebar-input", attrs)... }/>
}

component SidebarHeader(children gsx.Node, attrs gsx.Attrs) {
	<div { withSlot("sidebar-header", attrs)... }>{ children }</div>
}

component SidebarFooter(children gsx.Node, attrs gsx.Attrs) {
	<div { withSlot("sidebar-footer", attrs)... }>{ children }</div>
}

component SidebarSeparator(attrs gsx.Attrs) {
	<Separator { withSlot("sidebar-separator", attrs)... }/>
}

component SidebarContent(children gsx.Node, attrs gsx.Attrs) {
	<div { withSlot("sidebar-content", attrs)... }>{ children }</div>
}

component SidebarGroup(children gsx.Node, attrs gsx.Attrs) {
	<div { withSlot("sidebar-group", attrs)... }>{ children }</div>
}

component SidebarGroupLabel(children gsx.Node, attrs gsx.Attrs) {
	<div { withSlot("sidebar-group-label", attrs)... }>{ children }</div>
}

component SidebarGroupAction(children gsx.Node, attrs gsx.Attrs) {
	<button type="button" { withSlot("sidebar-group-action", attrs)... }>
		{ children }
	</button>
}

component SidebarGroupContent(children gsx.Node, attrs gsx.Attrs) {
	<div { withSlot("sidebar-group-content", attrs)... }>{ children }</div>
}

component SidebarMenu(children gsx.Node, attrs gsx.Attrs) {
	<ul { withSlot("sidebar-menu", attrs)... }>{ children }</ul>
}

component SidebarMenuItem(children gsx.Node, attrs gsx.Attrs) {
	<li { withSlot("sidebar-menu-item", attrs)... }>{ children }</li>
}

component SidebarMenuButton(isActive bool, variant string, size string, tooltip string, children gsx.Node, attrs gsx.Attrs) {
	{ if tooltip == "" {
		<button
			type="button"
			data-variant={variant |> default("default")}
			data-size={size |> default("default")}
			data-active={isActive}
			{ withSlot("sidebar-menu-button", attrs)... }
		>
			{ children }
		</button>
	} else {
		<Tooltip { withSlot("sidebar-menu-button-tooltip", nil)... }>
			<button
				type="button"
				data-variant={variant |> default("default")}
				data-size={size |> default("default")}
				data-active={isActive}
				data-gsxui-tooltip-trigger
				{ withSlot("sidebar-menu-button", attrs)... }
			>
				{ children }
			</button>
			<TooltipContent { withSlot("sidebar-menu-button-tooltip-content", nil)... }>
				{ tooltip }
			</TooltipContent>
		</Tooltip>
	} }
}

component SidebarMenuAction(showOnHover bool, children gsx.Node, attrs gsx.Attrs) {
	<button
		type="button"
		data-show-on-hover={showOnHover}
		{ withSlot("sidebar-menu-action", attrs)... }
	>
		{ children }
	</button>
}

component SidebarMenuBadge(children gsx.Node, attrs gsx.Attrs) {
	<div { withSlot("sidebar-menu-badge", attrs)... }>{ children }</div>
}

// The randomized width is the one dynamic presentation value in this part.
// It is chosen once per server render, then consumed by style CSS.
component SidebarMenuSkeleton(showIcon bool, attrs gsx.Attrs) {
	{{
		width := strconv.Itoa(rand.Intn(40)+50) + "%"
	}}
	<div { withSlot("sidebar-menu-skeleton", attrs)... }>
		{ if showIcon {
			<Skeleton { withSlot("sidebar-menu-skeleton-icon", nil)... }/>
		} }
		<Skeleton
			style=css`--skeleton-width:@{width}`
			{ withSlot("sidebar-menu-skeleton-text", nil)... }
		/>
	</div>
}

component SidebarMenuSub(children gsx.Node, attrs gsx.Attrs) {
	<ul { withSlot("sidebar-menu-sub", attrs)... }>{ children }</ul>
}

component SidebarMenuSubItem(children gsx.Node, attrs gsx.Attrs) {
	<li { withSlot("sidebar-menu-sub-item", attrs)... }>{ children }</li>
}

component SidebarMenuSubButton(size string, isActive bool, children gsx.Node, attrs gsx.Attrs) {
	<a
		data-size={size |> default("md")}
		data-active={isActive}
		{ withSlot("sidebar-menu-sub-button", attrs)... }
	>
		{ children }
	</a>
}
