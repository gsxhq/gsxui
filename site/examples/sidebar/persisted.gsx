package sidebar

import (
	"github.com/gsxhq/gsxui/ui"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Persisted is the copyable recipe for the one thing ui/sidebar.js
// deliberately does NOT do: remember open/collapsed across page loads.
// shadcn's own React SidebarProvider writes a non-HttpOnly `sidebar_state`
// cookie on every toggle and reads it back via a Server Component; gsxui
// drops that from the component itself (ui/sidebar.gsx's own package doc
// comment, decision 1) because a component that unilaterally picks a
// storage mechanism fights every consumer who'd rather use a Go session,
// Alpine, htmx, or nothing at all. This example is the two halves a real
// app wires up instead:
//
//  1. SERVER: read the `sidebar_state` cookie into `open` before you
//     render SidebarProvider — a plain net/http example (swap in your own
//     router's request/cookie types, the shape is identical):
//
//     func Page(w http.ResponseWriter, r *http.Request) {
//         open := true // no cookie yet: shadcn's own default is expanded
//         if c, err := r.Cookie("sidebar_state"); err == nil {
//             open = c.Value == "true"
//         }
//         ui.SidebarProvider(open, pageBody(), nil).Render(r.Context(), w)
//     }
//
//  2. CLIENT: write the cookie back whenever the sidebar toggles — three
//     lines, listening for the gsxui:change event sidebar.js emits on
//     SidebarProvider's own wrapper (see the <script> below, which is
//     this exact snippet, live):
//
//     document.addEventListener("gsxui:change", (e) => {
//       document.cookie = `sidebar_state=${e.detail.open}; path=/; max-age=604800`;
//     });
//
// 604800 is SIDEBAR_COOKIE_MAX_AGE from the shadcn source (7 days) — a
// plain number here, not a gsxui constant, since the component ships no
// cookie code at all to hang one off of. This demo can't show the actual
// reload-persists round trip (a static example gallery page, not a live
// per-visitor route), but toggling it below DOES write the real cookie in
// your browser — inspect it in devtools, or check `document.cookie` in the
// console.
component Persisted() {
	<div>
		<ui.SidebarProvider open={true} class="h-64 min-h-0 overflow-hidden rounded-lg border">
			<ui.Sidebar collapsible="icon">
				<ui.SidebarHeader>
					<div class="px-2 py-1 text-sm font-semibold">Acme Inc</div>
				</ui.SidebarHeader>
				<ui.SidebarContent>
					<ui.SidebarGroup>
						<ui.SidebarMenu>
							<ui.SidebarMenuItem>
								<ui.SidebarMenuButton isActive={true} tooltip="Home">
									<icon.House/>
									<span>Home</span>
								</ui.SidebarMenuButton>
							</ui.SidebarMenuItem>
							<ui.SidebarMenuItem>
								<ui.SidebarMenuButton tooltip="Settings">
									<icon.Settings/>
									<span>Settings</span>
								</ui.SidebarMenuButton>
							</ui.SidebarMenuItem>
						</ui.SidebarMenu>
					</ui.SidebarGroup>
				</ui.SidebarContent>
				<ui.SidebarRail/>
			</ui.Sidebar>
			<ui.SidebarInset>
				<header class="flex h-12 items-center gap-2 border-b px-4">
					<ui.SidebarTrigger/>
					<span class="text-sm text-muted-foreground">Toggle me, then check document.cookie</span>
				</header>
			</ui.SidebarInset>
		</ui.SidebarProvider>
		<script>
			document.addEventListener("gsxui:change", (e) => {
				document.cookie = `sidebar_state=${e.detail.open}; path=/; max-age=604800`;
			});
		</script>
	</div>
}
