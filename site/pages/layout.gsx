package pages

import (
	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/internal/registry"
	"github.com/gsxhq/gsxui/ui"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Layout is the shared page shell every site page renders through: header
// (wordmark + doc search + GitHub link), sidebar (component list from the
// registry — derived, so it can never drift from what `ui/` actually
// ships), and footer. active names the component whose sidebar entry
// should highlight; pages outside /components/ pass "".
//
// Doc search: an outer ui.Dialog root wires the header trigger button to
// CommandDialog's nested dialog element by proximity (dialog.js's
// root.querySelector reaches through the inner root), and command.js's
// global Cmd-K/Ctrl-K hotkey toggles the same dialog. The search index is
// the registry component list plus the static pages — derived, no manual
// list to drift.
component Layout(title string, active string, children gsx.Node) {
	<!DOCTYPE html>
	<html lang="en">
		<siteHead title={title} entry="web/main.js"/>
		<body class="min-h-svh bg-background text-foreground antialiased">
			<header class="sticky top-0 z-10 border-b border-border bg-background/95 backdrop-blur">
				<div class="mx-auto flex h-14 max-w-6xl items-center justify-between px-4">
					<a href={Home{} |> url} class="flex items-center">
						<siteLogo/>
					</a>
					<nav class="flex items-center gap-4">
						<ui.Dialog>
							<button
								data-gsxui-dialog-trigger
								data-gsxui-slot-dialog-trigger
								type="button"
								aria-haspopup="dialog"
								aria-expanded="false"
								class="hidden h-8 w-56 items-center gap-2 rounded-lg border bg-muted/50 px-2.5 text-sm text-muted-foreground transition-colors hover:bg-muted sm:inline-flex"
							>
								<icon.Search class="size-4"/>
								<span class="flex-1 text-left">Search docs...</span>
								<ui.Kbd>⌘K</ui.Kbd>
							</button>
							<ui.CommandDialog title="Search documentation" description="Search components and pages...">
								<ui.CommandInput placeholder="Search documentation..."/>
								<ui.CommandList>
									<ui.CommandEmpty>No results found.</ui.CommandEmpty>
									<ui.CommandGroup heading="Components">
										{{ searchNames, _ := registry.Components() }}
										{ for _, name := range searchNames {
											<ui.CommandItem data-href={"/components/" + name} class="capitalize">{ name }</ui.CommandItem>
										} }
									</ui.CommandGroup>
									<ui.CommandGroup heading="Pages">
										<ui.CommandItem data-href={Home{} |> url}>Home</ui.CommandItem>
										<ui.CommandItem data-href={ComponentsIndex{} |> url}>Components</ui.CommandItem>
										<ui.CommandItem data-href={GettingStarted{} |> url}>Getting Started</ui.CommandItem>
										<ui.CommandItem data-href={Theming{} |> url}>Theming</ui.CommandItem>
										<ui.CommandItem data-href={Theme{} |> url}>Theme Editor</ui.CommandItem>
									</ui.CommandGroup>
								</ui.CommandList>
							</ui.CommandDialog>
						</ui.Dialog>
						<a
							href={Theme{} |> url}
							class="text-sm text-muted-foreground transition-colors hover:text-foreground"
						>
							Theme
						</a>
						<a
							href="https://github.com/gsxhq/gsxui"
							target="_blank"
							rel="noreferrer"
							class="text-sm text-muted-foreground transition-colors hover:text-foreground"
						>
							GitHub
						</a>
						<button
							type="button"
							data-site-theme-toggle
							aria-label="Toggle theme"
							title="Toggle theme"
							class="inline-flex size-8 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:bg-accent hover:text-accent-foreground"
						>
							<svg
								xmlns="http://www.w3.org/2000/svg"
								width="24"
								height="24"
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="2"
								stroke-linecap="round"
								stroke-linejoin="round"
								class="size-4.5"
							>
								<path d="M12 12m-9 0a9 9 0 1 0 18 0a9 9 0 1 0 -18 0"/>
								<path d="M12 3l0 18"/>
								<path d="M12 9l4.65 -4.65"/>
								<path d="M12 14.3l7.37 -7.37"/>
								<path d="M12 19.6l8.85 -8.85"/>
							</svg>
						</button>
					</nav>
				</div>
			</header>
			<div class="mx-auto flex max-w-6xl gap-10 px-4 py-10">
				<aside class="hidden w-44 shrink-0 md:block">
					<nav class="sticky top-20 flex flex-col gap-4 text-sm">
						<div class="flex flex-col gap-1">
							<h3 class="px-2 text-xs font-semibold uppercase tracking-wide text-muted-foreground">Docs</h3>
							<a
								href={GettingStarted{} |> url}
								class={
									"rounded-md px-2 py-1 text-muted-foreground transition-colors hover:bg-accent hover:text-accent-foreground",
									"bg-accent text-accent-foreground": active == "getting-started"
								}
							>
								Getting Started
							</a>
							<a
								href={Theming{} |> url}
								class={
									"rounded-md px-2 py-1 text-muted-foreground transition-colors hover:bg-accent hover:text-accent-foreground",
									"bg-accent text-accent-foreground": active == "theming"
								}
							>
								Theming
							</a>
						</div>
						<div class="flex flex-col gap-1">
							<h3 class="px-2 text-xs font-semibold uppercase tracking-wide text-muted-foreground">Components</h3>
							{{ names, _ := registry.Components() }}
							{ for _, name := range names {
								<a
									href={"/components/" + name}
									class={
										"rounded-md px-2 py-1 capitalize text-muted-foreground transition-colors hover:bg-accent hover:text-accent-foreground",
										"bg-accent text-accent-foreground": active == name
									}
								>
									{ name }
								</a>
							} }
						</div>
					</nav>
				</aside>
				<main class="min-w-0 flex-1">{ children }</main>
			</div>
			<footer class="border-t border-border">
				<div class="mx-auto max-w-6xl px-4 py-6 text-sm text-muted-foreground">
					gsxui — shadcn-style components for gsx. Copy-in, type-checked, server-rendered.
				</div>
			</footer>
			{/* Mounted once per page: the bottom-right region ui/sonner.js
			   appends every client-constructed toast <li> into. */}
			<ui.Toaster/>
		</body>
	</html>
}
