package navigationmenu

import (
	"github.com/gsxhq/gsxui/ui"
)

// Mega demonstrates two independently-sized floating panels: Products opens
// a wide three-column grid, Solutions opens a narrow single-column list —
// each is its own self-contained, fully-chromed NavigationMenuContent,
// positioned under its own trigger (this port ships shadcn's
// viewport={false} configuration only; see ui/navigation-menu.gsx's own
// file header GAP paragraph for why there is no shared, JS-measured
// viewport panel morphing between the two instead). Resources is a plain
// link with no dropdown at all, reflected as variant="trigger".
component Mega() {
	<ui.NavigationMenu>
		<ui.NavigationMenuList>
			<ui.NavigationMenuItem>
				<ui.NavigationMenuTrigger>Products</ui.NavigationMenuTrigger>
				<ui.NavigationMenuContent>
					<div class="grid w-[36rem] grid-cols-3 gap-2">
						<ui.NavigationMenuLink href="#">
							<div class="text-sm font-medium">Analytics</div>
							<div class="text-muted-foreground">Track usage and conversion across every surface.</div>
						</ui.NavigationMenuLink>
						<ui.NavigationMenuLink href="#">
							<div class="text-sm font-medium">Automation</div>
							<div class="text-muted-foreground">Wire triggers to actions without writing glue code.</div>
						</ui.NavigationMenuLink>
						<ui.NavigationMenuLink href="#">
							<div class="text-sm font-medium">Billing</div>
							<div class="text-muted-foreground">Usage-based invoices, taxes, and dunning handled for you.</div>
						</ui.NavigationMenuLink>
						<ui.NavigationMenuLink href="#">
							<div class="text-sm font-medium">Search</div>
							<div class="text-muted-foreground">Full-text and vector search over your own data.</div>
						</ui.NavigationMenuLink>
						<ui.NavigationMenuLink href="#">
							<div class="text-sm font-medium">Storage</div>
							<div class="text-muted-foreground">Durable object storage with signed upload URLs.</div>
						</ui.NavigationMenuLink>
						<ui.NavigationMenuLink href="#">
							<div class="text-sm font-medium">Workflows</div>
							<div class="text-muted-foreground">Durable, retryable multi-step background jobs.</div>
						</ui.NavigationMenuLink>
					</div>
				</ui.NavigationMenuContent>
			</ui.NavigationMenuItem>
			<ui.NavigationMenuItem>
				<ui.NavigationMenuTrigger>Solutions</ui.NavigationMenuTrigger>
				<ui.NavigationMenuContent>
					<div class="grid w-56 gap-2">
						<ui.NavigationMenuLink href="#">Startups</ui.NavigationMenuLink>
						<ui.NavigationMenuLink href="#">Enterprise</ui.NavigationMenuLink>
						<ui.NavigationMenuLink href="#">Agencies</ui.NavigationMenuLink>
					</div>
				</ui.NavigationMenuContent>
			</ui.NavigationMenuItem>
			<ui.NavigationMenuItem>
				<ui.NavigationMenuLink variant="trigger" href="#">Resources</ui.NavigationMenuLink>
			</ui.NavigationMenuItem>
		</ui.NavigationMenuList>
	</ui.NavigationMenu>
}
