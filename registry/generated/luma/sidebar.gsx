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
//
// side is physical, matching shadcn: side="left"/"right" always anchors the
// desktop rail and the mobile Sheet to that physical viewport edge, in both
// dir="ltr" and dir="rtl" documents — the container's left/right offset,
// the border between the rail and the inset content, and the rail's own
// resize-cursor/position are all data-side-keyed physical geometry and stay
// that way. Interior presentation — menu padding, badge/action offsets,
// group-label/menu-button text alignment, the SidebarTrigger icon's
// rtl:rotate-180 flip — is logical and follows dir normally. The mobile
// Sheet needs no explicit dir plumbing: it composes <ui.Dialog>, which
// renders inline in the server-rendered document, so dir inherits from the
// ancestor <html dir="..."> the same as everywhere else.
//
// Every class attribute below is resolved to concrete utilities at generation
// time. Six markers carry no class because nothing styles them — see
// registry/canonical/shapes/sidebar.go for which and why.
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
		class={ "has-[>[data-gsxui-slot-sidebar-desktop][data-variant=inset]]:bg-sidebar" }
		data-state={state}
		style=css`--sidebar-width:@{sidebarWidth};--sidebar-width-icon:@{sidebarWidthIcon}`
		{ attrs... }
		data-gsxui-slot-sidebar-wrapper
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
			class={ "data-[collapsible=none]:bg-sidebar data-[collapsible=none]:text-sidebar-foreground" }
			data-side={s}
			data-variant={v}
			data-collapsible="none"
			{ attrs... }
			data-gsxui-slot-sidebar
		>
			{ children }
		</div>
	} else {
		<>
			<Sheet { attrs... } data-gsxui-slot-sidebar-mobile-root>
				<SheetContent
					class={
						"data-[collapsible=none]:bg-sidebar data-[collapsible=none]:text-sidebar-foreground",
						"w-[var(--sidebar-width)] max-w-none p-0 bg-sidebar text-sidebar-foreground sm:max-w-none [&>[data-gsxui-slot-sheet-close-button]]:hidden"
					}
					side={s}
					data-mobile="true"
					style=css`--sidebar-width:@{sidebarWidthMobile}`
					data-gsxui-slot-sidebar-mobile-content
					data-gsxui-slot-sidebar
				>
					<SheetHeader class={ "sr-only p-0" } data-gsxui-slot-sidebar-mobile-header>
						<SheetTitle data-gsxui-slot-sidebar-mobile-title>Sidebar</SheetTitle>
						<SheetDescription data-gsxui-slot-sidebar-mobile-description>
							Displays the mobile sidebar.
						</SheetDescription>
					</SheetHeader>
					<div data-gsxui-slot-sidebar-mobile-inner>{ children }</div>
				</SheetContent>
			</Sheet>
			<div
				class={
					"data-[collapsible=none]:bg-sidebar data-[collapsible=none]:text-sidebar-foreground",
					"text-sidebar-foreground"
				}
				data-state={state}
				data-collapsible={activeCollapsible}
				data-gsxui-sidebar-collapsible={c}
				data-variant={v}
				data-side={s}
				data-gsxui-slot-sidebar-desktop
				data-gsxui-slot-sidebar
			>
				<div class={ "transition-[width] duration-200 ease-linear" } data-gsxui-slot-sidebar-gap></div>
				<div
					class={
						"transition-[left,right,width] duration-200 ease-linear [[data-gsxui-slot-sidebar-desktop][data-variant=sidebar][data-side=left]>&]:border-r [[data-gsxui-slot-sidebar-desktop][data-variant=sidebar][data-side=left]>&]:border-sidebar-border [[data-gsxui-slot-sidebar-desktop][data-variant=sidebar][data-side=right]>&]:border-l [[data-gsxui-slot-sidebar-desktop][data-variant=sidebar][data-side=right]>&]:border-sidebar-border"
					}
					{ attrs... }
					data-gsxui-slot-sidebar-container
				>
					<div
						class={
							"bg-sidebar group-data-[variant=floating]:ring-sidebar-border group-data-[variant=floating]:rounded-2xl group-data-[variant=floating]:shadow-sm group-data-[variant=floating]:ring-1"
						}
						data-gsxui-slot-sidebar-inner
					>
						{ children }
					</div>
				</div>
			</div>
		</>
	} }
}

component SidebarTrigger(attrs gsx.Attrs) {
	<Button
		class={ "size-7" }
		variant="ghost"
		size="icon"
		{ attrs... }
		data-gsxui-slot-sidebar-trigger
	>
		<icon.PanelLeft class={ "rtl:rotate-180" }/>
		<span class={ "sr-only" } data-gsxui-slot-sidebar-trigger-label>Toggle Sidebar</span>
	</Button>
}

// SidebarRail is a pointer affordance matching the reference: it is not a
// keyboard tab stop and activates the same state transition as the trigger.
component SidebarRail(attrs gsx.Attrs) {
	<button
		class={ "hover:after:bg-sidebar-border" }
		type="button"
		aria-label="Toggle Sidebar"
		tabindex="-1"
		title="Toggle Sidebar"
		{ attrs... }
		data-gsxui-slot-sidebar-rail
	></button>
}

component SidebarInset(children gsx.Node, attrs gsx.Attrs) {
	<main
		class={
			"bg-background md:peer-data-[variant=inset]:m-2 md:peer-data-[variant=inset]:ml-0 md:peer-data-[variant=inset]:rounded-2xl md:peer-data-[variant=inset]:shadow-sm md:peer-data-[variant=inset]:peer-data-[state=collapsed]:ml-2"
		}
		{ attrs... }
		data-gsxui-slot-sidebar-inset
	>
		{ children }
	</main>
}

component SidebarInput(attrs gsx.Attrs) {
	<Input class={ "bg-input/50 h-8 w-full shadow-none" } { attrs... } data-gsxui-slot-sidebar-input/>
}

component SidebarHeader(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "gap-2 p-2 [--radius:var(--radius-xl)] flex flex-col" } { attrs... } data-gsxui-slot-sidebar-header>
		{ children }
	</div>
}

component SidebarFooter(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "gap-2 p-2 flex flex-col" } { attrs... } data-gsxui-slot-sidebar-footer>{ children }</div>
}

component SidebarSeparator(attrs gsx.Attrs) {
	<Separator class={ "bg-sidebar-border mx-2" } { attrs... } data-gsxui-slot-sidebar-separator/>
}

component SidebarContent(children gsx.Node, attrs gsx.Attrs) {
	<div
		class={ "[scrollbar-width:none] [&::-webkit-scrollbar]:hidden gap-2 [--radius:var(--radius-xl)]" }
		{ attrs... }
		data-gsxui-slot-sidebar-content
	>
		{ children }
	</div>
}

component SidebarGroup(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "p-2 flex flex-col" } { attrs... } data-gsxui-slot-sidebar-group>{ children }</div>
}

component SidebarGroupLabel(children gsx.Node, attrs gsx.Attrs) {
	<div
		class={
			"text-sidebar-foreground/70 ring-sidebar-ring h-8 rounded-xl px-3 text-xs font-medium transition-[margin,opacity] duration-200 ease-linear group-data-[collapsible=icon]:-mt-8 group-data-[collapsible=icon]:opacity-0 focus-visible:ring-2 [&>svg]:size-4 flex"
		}
		{ attrs... }
		data-gsxui-slot-sidebar-group-label
	>
		{ children }
	</div>
}

component SidebarGroupAction(children gsx.Node, attrs gsx.Attrs) {
	<button
		class={
			"text-sidebar-foreground ring-sidebar-ring hover:bg-sidebar-accent hover:text-sidebar-accent-foreground absolute top-3.5 right-3 w-5 rounded-xl p-0 focus-visible:ring-2 [&>svg]:size-4 flex"
		}
		type="button"
		{ attrs... }
		data-gsxui-slot-sidebar-group-action
	>
		{ children }
	</button>
}

component SidebarGroupContent(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "text-sm" } { attrs... } data-gsxui-slot-sidebar-group-content>{ children }</div>
}

component SidebarMenu(children gsx.Node, attrs gsx.Attrs) {
	<ul class={ "gap-0.5 flex flex-col" } { attrs... } data-gsxui-slot-sidebar-menu>{ children }</ul>
}

component SidebarMenuItem(children gsx.Node, attrs gsx.Attrs) {
	<li class={ "relative" } { attrs... } data-gsxui-slot-sidebar-menu-item>{ children }</li>
}

component SidebarMenuButton(isActive bool, variant string, size string, tooltip string, children gsx.Node, attrs gsx.Attrs) {
	{ if tooltip == "" {
		<button
			class={
				"ring-sidebar-ring hover:bg-sidebar-accent hover:text-sidebar-accent-foreground active:bg-sidebar-accent active:text-sidebar-accent-foreground [&[data-active]]:bg-sidebar-accent [&[data-active]]:text-sidebar-accent-foreground data-[state=open]:hover:bg-sidebar-accent data-[state=open]:hover:text-sidebar-accent-foreground gap-2 rounded-xl px-3 py-2 text-left text-sm transition-[width,height,padding] group-has-data-[sidebar=menu-action]/menu-item:pr-8 group-data-[collapsible=icon]:size-8 group-data-[collapsible=icon]:p-2 focus-visible:ring-2 [&[data-active]]:font-medium flex"
			}
			type="button"
			data-variant={variant |> default("default")}
			data-size={size |> default("default")}
			data-active={isActive}
			{ attrs... }
			data-gsxui-slot-sidebar-menu-button
		>
			{ children }
		</button>
	} else {
		<Tooltip data-gsxui-slot-sidebar-menu-button-tooltip>
			<button
				class={
					"ring-sidebar-ring hover:bg-sidebar-accent hover:text-sidebar-accent-foreground active:bg-sidebar-accent active:text-sidebar-accent-foreground [&[data-active]]:bg-sidebar-accent [&[data-active]]:text-sidebar-accent-foreground data-[state=open]:hover:bg-sidebar-accent data-[state=open]:hover:text-sidebar-accent-foreground gap-2 rounded-xl px-3 py-2 text-left text-sm transition-[width,height,padding] group-has-data-[sidebar=menu-action]/menu-item:pr-8 group-data-[collapsible=icon]:size-8 group-data-[collapsible=icon]:p-2 focus-visible:ring-2 [&[data-active]]:font-medium flex"
				}
				type="button"
				data-variant={variant |> default("default")}
				data-size={size |> default("default")}
				data-active={isActive}
				data-gsxui-tooltip-trigger
				{ attrs... }
				data-gsxui-slot-sidebar-menu-button
			>
				{ children }
			</button>
			<TooltipContent data-gsxui-slot-sidebar-menu-button-tooltip-content>
				{ tooltip }
			</TooltipContent>
		</Tooltip>
	} }
}

component SidebarMenuAction(showOnHover bool, children gsx.Node, attrs gsx.Attrs) {
	<button
		class={
			"text-sidebar-foreground ring-sidebar-ring hover:bg-sidebar-accent hover:text-sidebar-accent-foreground peer-hover/menu-button:text-sidebar-accent-foreground absolute top-1.5 right-1 aspect-square w-5 rounded-xl p-0 peer-data-[size=default]/menu-button:top-2 peer-data-[size=lg]/menu-button:top-2.5 peer-data-[size=sm]/menu-button:top-1 focus-visible:ring-2 [&>svg]:size-4 md:[&[data-show-on-hover]]:opacity-0 md:[[data-gsxui-slot-sidebar-menu-item]:hover>&[data-show-on-hover]]:opacity-100 md:[[data-gsxui-slot-sidebar-menu-item]:focus-within>&[data-show-on-hover]]:opacity-100 md:[&[data-show-on-hover][data-state=open]]:opacity-100 [[data-gsxui-slot-sidebar-menu-button][data-active]~&[data-show-on-hover]]:text-sidebar-primary-foreground flex"
		}
		type="button"
		data-show-on-hover={showOnHover}
		{ attrs... }
		data-gsxui-slot-sidebar-menu-action
	>
		{ children }
	</button>
}

component SidebarMenuBadge(children gsx.Node, attrs gsx.Attrs) {
	<div
		class={
			"text-sidebar-foreground peer-hover/menu-button:text-sidebar-accent-foreground peer-data-active/menu-button:text-sidebar-accent-foreground pointer-events-none absolute right-1 flex h-5 min-w-5 rounded-xl px-1 text-xs font-medium peer-data-[size=default]/menu-button:top-1.5 peer-data-[size=lg]/menu-button:top-2.5 peer-data-[size=sm]/menu-button:top-1"
		}
		{ attrs... }
		data-gsxui-slot-sidebar-menu-badge
	>
		{ children }
	</div>
}

// The randomized width is the one dynamic presentation value in this part.
// It is chosen once per server render, then consumed by style CSS.
component SidebarMenuSkeleton(showIcon bool, attrs gsx.Attrs) {
	{{ width := strconv.Itoa(rand.Intn(40)+50) + "%" }}
	<div class={ "h-8 gap-2 rounded-xl px-2 flex" } { attrs... } data-gsxui-slot-sidebar-menu-skeleton>
		{ if showIcon {
			<Skeleton class={ "size-4 rounded-xl" } data-gsxui-slot-sidebar-menu-skeleton-icon/>
		} }
		<Skeleton
			class={ "h-4" }
			style=css`--skeleton-width:@{width}`
			data-gsxui-slot-sidebar-menu-skeleton-text
		/>
	</div>
}

component SidebarMenuSub(children gsx.Node, attrs gsx.Attrs) {
	<ul
		class={
			"border-sidebar-border mx-3.5 translate-x-px gap-1 border-l px-2.5 py-0.5 group-data-[collapsible=icon]:hidden flex flex-col"
		}
		{ attrs... }
		data-gsxui-slot-sidebar-menu-sub
	>
		{ children }
	</ul>
}

component SidebarMenuSubItem(children gsx.Node, attrs gsx.Attrs) {
	<li class={ "relative" } { attrs... } data-gsxui-slot-sidebar-menu-sub-item>{ children }</li>
}

component SidebarMenuSubButton(size string, isActive bool, children gsx.Node, attrs gsx.Attrs) {
	<a
		class={
			"text-sidebar-foreground ring-sidebar-ring hover:bg-sidebar-accent hover:text-sidebar-accent-foreground active:bg-sidebar-accent active:text-sidebar-accent-foreground [&>svg]:text-sidebar-accent-foreground [&[data-active]]:bg-sidebar-accent [&[data-active]]:text-sidebar-accent-foreground h-7 gap-2 rounded-xl px-3 focus-visible:ring-2 data-[size=md]:text-sm data-[size=sm]:text-xs [&>svg]:size-4 flex"
		}
		data-size={size |> default("md")}
		data-active={isActive}
		{ attrs... }
		data-gsxui-slot-sidebar-menu-sub-button
	>
		{ children }
	</a>
}
