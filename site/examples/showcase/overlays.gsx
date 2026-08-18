package showcase

import "github.com/gsxhq/gsxui/ui"

// OverlaysCard proves the interactive components work server-rendered:
// Tabs switch between an overlays pane (Dialog, DropdownMenu, Tooltip)
// and a feedback pane (a server-flash Toast appended into ui.Toaster,
// which siteLayout mounts once per page).
component OverlaysCard() {
	<ui.Card>
		<ui.CardHeader>
			<ui.CardTitle>Interactive, no framework</ui.CardTitle>
			<ui.CardDescription>
				Dialogs, menus, tooltips and toasts — server-rendered, hydrated by tiny shims.
			</ui.CardDescription>
		</ui.CardHeader>
		<ui.CardContent>
			<ui.Tabs value="overlays">
				<ui.TabsList>
					<ui.TabsTrigger value="overlays" selected>Overlays</ui.TabsTrigger>
					<ui.TabsTrigger value="feedback">Feedback</ui.TabsTrigger>
				</ui.TabsList>
				<ui.TabsContent value="overlays" selected>
					<div class="flex flex-wrap items-center gap-3 pt-2">
						<ui.Dialog>
							<ui.Button
								variant="outline"
								data-gsxui-slot-dialog-trigger
								aria-haspopup="dialog"
								aria-expanded="false"
							>
								Open dialog
							</ui.Button>
							<ui.DialogContent>
								<ui.DialogHeader>
									<ui.DialogTitle>Edit profile</ui.DialogTitle>
									<ui.DialogDescription>
										Rendered by ui/dialog on the native &lt;dialog&gt; element — no client framework required.
									</ui.DialogDescription>
								</ui.DialogHeader>
								<ui.DialogFooter showCloseButton={true}></ui.DialogFooter>
							</ui.DialogContent>
						</ui.Dialog>
						<ui.DropdownMenu>
							<ui.Button
								variant="outline"
								data-gsxui-slot-dropdown-menu-trigger
								aria-haspopup="menu"
								aria-expanded="false"
							>
								Options
							</ui.Button>
							<ui.DropdownMenuContent>
								<ui.DropdownMenuGroup>
									<ui.DropdownMenuLabel>My Account</ui.DropdownMenuLabel>
									<ui.DropdownMenuSeparator/>
									<ui.DropdownMenuItem>Profile</ui.DropdownMenuItem>
									<ui.DropdownMenuItem>Billing</ui.DropdownMenuItem>
									<ui.DropdownMenuItem>
										Settings <ui.DropdownMenuShortcut>⌘,</ui.DropdownMenuShortcut>
									</ui.DropdownMenuItem>
								</ui.DropdownMenuGroup>
							</ui.DropdownMenuContent>
						</ui.DropdownMenu>
						<ui.Tooltip>
							<ui.Button
								variant="outline"
								data-gsxui-slot-tooltip-trigger
							>
								Hover me
							</ui.Button>
							<ui.TooltipContent>Server-rendered tooltip</ui.TooltipContent>
						</ui.Tooltip>
					</div>
				</ui.TabsContent>
				<ui.TabsContent value="feedback">
					<div class="flex flex-col items-start gap-3 pt-2">
						<p class="text-sm text-muted-foreground">
							Server flash pattern: a pre-rendered toast row appended into the page's toaster, exactly like an HTMX
							out-of-band swap.
						</p>
						<ui.Button variant="outline" id="home-showcase-toast-btn">Show toast</ui.Button>
						<template data-home-showcase-toast>
							<ui.Toast
								toastType="success"
								title="Saved"
								description="A server-rendered toast, adopted by ui/toaster."
							/>
						</template>
						<script>
							document.getElementById("home-showcase-toast-btn").addEventListener("click", () => {
								const tpl = document.querySelector("template[data-home-showcase-toast]");
								const viewport = document.getElementById("gsxui-toaster");
								if (!tpl || !viewport) return;
								viewport.appendChild(tpl.content.firstElementChild.cloneNode(true));
							});
						</script>
					</div>
				</ui.TabsContent>
			</ui.Tabs>
		</ui.CardContent>
	</ui.Card>
}
