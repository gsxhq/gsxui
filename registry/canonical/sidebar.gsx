package canonical

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
		class={ sidebar.Wrapper() }
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
			class={ sidebar.Root() }
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
					class={ sidebar.Root(), sidebar.MobileContent() }
					side={s}
					data-mobile="true"
					style=css`--sidebar-width:@{sidebarWidthMobile}`
					data-gsxui-slot-sidebar-mobile-content
					data-gsxui-slot-sidebar
				>
					<SheetHeader class={ sidebar.MobileHeader() } data-gsxui-slot-sidebar-mobile-header>
						<SheetTitle data-gsxui-slot-sidebar-mobile-title>Sidebar</SheetTitle>
						<SheetDescription data-gsxui-slot-sidebar-mobile-description>
							Displays the mobile sidebar.
						</SheetDescription>
					</SheetHeader>
					<div data-gsxui-slot-sidebar-mobile-inner>{ children }</div>
				</SheetContent>
			</Sheet>
			<div
				class={ sidebar.Root(), sidebar.Desktop() }
				data-state={state}
				data-collapsible={activeCollapsible}
				data-gsxui-sidebar-collapsible={c}
				data-variant={v}
				data-side={s}
				data-gsxui-slot-sidebar-desktop
				data-gsxui-slot-sidebar
			>
				<div class={ sidebar.Gap() } data-gsxui-slot-sidebar-gap></div>
				<div class={ sidebar.Container() } { attrs... } data-gsxui-slot-sidebar-container>
					<div class={ sidebar.Inner() } data-gsxui-slot-sidebar-inner>
						{ children }
					</div>
				</div>
			</div>
		</>
	} }
}

component SidebarTrigger(attrs gsx.Attrs) {
	<Button
		class={ sidebar.Trigger() }
		variant="ghost"
		size="icon"
		{ attrs... }
		data-gsxui-slot-sidebar-trigger
	>
		<icon.PanelLeft class={ "rtl:rotate-180" }/>
		<span class={ sidebar.TriggerLabel() } data-gsxui-slot-sidebar-trigger-label>Toggle Sidebar</span>
	</Button>
}

// SidebarRail is a pointer affordance matching the reference: it is not a
// keyboard tab stop and activates the same state transition as the trigger.
component SidebarRail(attrs gsx.Attrs) {
	<button
		class={ sidebar.Rail() }
		type="button"
		aria-label="Toggle Sidebar"
		tabindex="-1"
		title="Toggle Sidebar"
		{ attrs... }
		data-gsxui-slot-sidebar-rail
	></button>
}

component SidebarInset(children gsx.Node, attrs gsx.Attrs) {
	<main class={ sidebar.Inset() } { attrs... } data-gsxui-slot-sidebar-inset>{ children }</main>
}

component SidebarInput(attrs gsx.Attrs) {
	<Input class={ sidebar.Input() } { attrs... } data-gsxui-slot-sidebar-input/>
}

component SidebarHeader(children gsx.Node, attrs gsx.Attrs) {
	<div class={ sidebar.Header() } { attrs... } data-gsxui-slot-sidebar-header>{ children }</div>
}

component SidebarFooter(children gsx.Node, attrs gsx.Attrs) {
	<div class={ sidebar.Footer() } { attrs... } data-gsxui-slot-sidebar-footer>{ children }</div>
}

component SidebarSeparator(attrs gsx.Attrs) {
	<Separator class={ sidebar.Separator() } { attrs... } data-gsxui-slot-sidebar-separator/>
}

component SidebarContent(children gsx.Node, attrs gsx.Attrs) {
	<div class={ sidebar.Content() } { attrs... } data-gsxui-slot-sidebar-content>{ children }</div>
}

component SidebarGroup(children gsx.Node, attrs gsx.Attrs) {
	<div class={ sidebar.Group() } { attrs... } data-gsxui-slot-sidebar-group>{ children }</div>
}

component SidebarGroupLabel(children gsx.Node, attrs gsx.Attrs) {
	<div class={ sidebar.GroupLabel() } { attrs... } data-gsxui-slot-sidebar-group-label>{ children }</div>
}

component SidebarGroupAction(children gsx.Node, attrs gsx.Attrs) {
	<button class={ sidebar.GroupAction() } type="button" { attrs... } data-gsxui-slot-sidebar-group-action>
		{ children }
	</button>
}

component SidebarGroupContent(children gsx.Node, attrs gsx.Attrs) {
	<div class={ sidebar.GroupContent() } { attrs... } data-gsxui-slot-sidebar-group-content>{ children }</div>
}

component SidebarMenu(children gsx.Node, attrs gsx.Attrs) {
	<ul class={ sidebar.Menu() } { attrs... } data-gsxui-slot-sidebar-menu>{ children }</ul>
}

component SidebarMenuItem(children gsx.Node, attrs gsx.Attrs) {
	<li class={ sidebar.MenuItem() } { attrs... } data-gsxui-slot-sidebar-menu-item>{ children }</li>
}

// SidebarMenuButton renders a nav entry. A non-empty href renders an <a>
// instead of a <button> — the pendant of shadcn's canonical
// `<SidebarMenuButton asChild><a href=…>` composition (gsx has no asChild;
// ui.Button's href param sets the precedent). disabled always renders a
// real disabled <button>, even with href, because aria-disabled on an <a>
// only blocks pointer input — keyboard activation would still navigate. The
// active anchor also carries aria-current="page" (PaginationLink precedent);
// the tooltip wrapper marks the inner element as its trigger via the bare
// data-gsxui-tooltip-trigger attribute.
component SidebarMenuButton(isActive bool, variant string, size string, tooltip string, href string, disabled bool, children gsx.Node, attrs gsx.Attrs) {
	{{ link := href != "" && !disabled }}
	{ if tooltip == "" {
		<sidebarMenuButtonRoot
			isActive={isActive}
			variant={variant}
			size={size}
			href={href}
			link={link}
			disabled={disabled}
			tooltipTrigger={false}
			{ attrs... }
		>
			{ children }
		</sidebarMenuButtonRoot>
	} else {
		<Tooltip data-gsxui-slot-sidebar-menu-button-tooltip>
			<sidebarMenuButtonRoot
				isActive={isActive}
				variant={variant}
				size={size}
				href={href}
				link={link}
				disabled={disabled}
				tooltipTrigger={true}
				{ attrs... }
			>
				{ children }
			</sidebarMenuButtonRoot>
			<TooltipContent data-gsxui-slot-sidebar-menu-button-tooltip-content>
				{ tooltip }
			</TooltipContent>
		</Tooltip>
	} }
}

// sidebarMenuButtonRoot is the one element SidebarMenuButton renders — an
// <a> when link, else a <button> — so the shared attribute set is written
// once rather than per tag × tooltip combination.
component sidebarMenuButtonRoot(isActive bool, variant string, size string, href string, link bool, disabled bool, tooltipTrigger bool, children gsx.Node, attrs gsx.Attrs) {
	{ if link {
		<a
			class={ sidebar.MenuButton() }
			href={href}
			{ if isActive {
				aria-current="page"
			} }
			data-variant={variant |> default("default")}
			data-size={size |> default("default")}
			data-active={isActive}
			data-gsxui-tooltip-trigger={tooltipTrigger}
			{ attrs... }
			data-gsxui-slot-sidebar-menu-button
		>
			{ children }
		</a>
	} else {
		<button
			class={ sidebar.MenuButton() }
			type="button"
			disabled={disabled}
			data-variant={variant |> default("default")}
			data-size={size |> default("default")}
			data-active={isActive}
			data-gsxui-tooltip-trigger={tooltipTrigger}
			{ attrs... }
			data-gsxui-slot-sidebar-menu-button
		>
			{ children }
		</button>
	} }
}

component SidebarMenuAction(showOnHover bool, children gsx.Node, attrs gsx.Attrs) {
	<button
		class={ sidebar.MenuAction() }
		type="button"
		data-show-on-hover={showOnHover}
		{ attrs... }
		data-gsxui-slot-sidebar-menu-action
	>
		{ children }
	</button>
}

component SidebarMenuBadge(children gsx.Node, attrs gsx.Attrs) {
	<div class={ sidebar.MenuBadge() } { attrs... } data-gsxui-slot-sidebar-menu-badge>{ children }</div>
}

// The randomized width is the one dynamic presentation value in this part.
// It is chosen once per server render, then consumed by style CSS.
component SidebarMenuSkeleton(showIcon bool, attrs gsx.Attrs) {
	{{ width := strconv.Itoa(rand.Intn(40)+50) + "%" }}
	<div class={ sidebar.MenuSkeleton() } { attrs... } data-gsxui-slot-sidebar-menu-skeleton>
		{ if showIcon {
			<Skeleton class={ sidebar.MenuSkeletonIcon() } data-gsxui-slot-sidebar-menu-skeleton-icon/>
		} }
		<Skeleton
			class={ sidebar.MenuSkeletonText() }
			style=css`--skeleton-width:@{width}`
			data-gsxui-slot-sidebar-menu-skeleton-text
		/>
	</div>
}

component SidebarMenuSub(children gsx.Node, attrs gsx.Attrs) {
	<ul class={ sidebar.MenuSub() } { attrs... } data-gsxui-slot-sidebar-menu-sub>{ children }</ul>
}

component SidebarMenuSubItem(children gsx.Node, attrs gsx.Attrs) {
	<li class={ sidebar.MenuSubItem() } { attrs... } data-gsxui-slot-sidebar-menu-sub-item>{ children }</li>
}

component SidebarMenuSubButton(size string, isActive bool, children gsx.Node, attrs gsx.Attrs) {
	<a
		class={ sidebar.MenuSubButton() }
		data-size={size |> default("md")}
		data-active={isActive}
		{ attrs... }
		data-gsxui-slot-sidebar-menu-sub-button
	>
		{ children }
	</a>
}
