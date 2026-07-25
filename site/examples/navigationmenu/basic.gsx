// Package navigationmenu holds the site's example gsx components for
// ui/navigation-menu.
package navigationmenu

import (
	"github.com/gsxhq/gsxui/ui"
)

// Basic mirrors the shape of shadcn's own navigation-menu-demo.tsx: a plain
// top-level link (Home), a trigger opening a small grid of description
// links (Components), and another plain link (Docs) — the default
// data-viewport="true" mode, so Home/Docs each render via
// ui.NavigationMenuTriggerStyle() directly (no dropdown of their own) while
// Components opens the shared, JS-measured viewport panel.
component Basic() {
	<ui.NavigationMenu>
		<ui.NavigationMenuList>
			<ui.NavigationMenuItem>
				<ui.NavigationMenuLink class={ ui.NavigationMenuTriggerStyle() } href="#">Home</ui.NavigationMenuLink>
			</ui.NavigationMenuItem>
			<ui.NavigationMenuItem>
				<ui.NavigationMenuTrigger>Components</ui.NavigationMenuTrigger>
				<ui.NavigationMenuContent>
					<div class="grid w-80 gap-2">
						<ui.NavigationMenuLink href="#">
							<div class="text-sm font-medium">Dialog</div>
							<div class="text-muted-foreground">A window overlaid on the page, rendered with the native &lt;dialog&gt; element.</div>
						</ui.NavigationMenuLink>
						<ui.NavigationMenuLink href="#">
							<div class="text-sm font-medium">Dropdown Menu</div>
							<div class="text-muted-foreground">Displays a menu of actions or options, triggered by a button.</div>
						</ui.NavigationMenuLink>
						<ui.NavigationMenuLink href="#">
							<div class="text-sm font-medium">Tooltip</div>
							<div class="text-muted-foreground">A popup that displays information related to an element on hover.</div>
						</ui.NavigationMenuLink>
					</div>
				</ui.NavigationMenuContent>
			</ui.NavigationMenuItem>
			<ui.NavigationMenuItem>
				<ui.NavigationMenuLink class={ ui.NavigationMenuTriggerStyle() } href="#">Docs</ui.NavigationMenuLink>
			</ui.NavigationMenuItem>
		</ui.NavigationMenuList>
	</ui.NavigationMenu>
}
