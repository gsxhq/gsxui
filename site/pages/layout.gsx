package pages

import (
	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/internal/registry"
	"github.com/gsxhq/gsxui/ui"
	"github.com/gsxhq/gsxui/ui/icon"
)

type layoutMode string

const (
	layoutMarketing layoutMode = "marketing"
	layoutDocs      layoutMode = "docs"
	layoutWorkspace layoutMode = "workspace"
)

component docsNavigation(active string) {
	<div class="flex flex-col gap-4 text-sm">
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
	</div>
}

component compactDocsNavigation(active string) {
	<ui.Popover data-site-docs-mobile-nav class="lg:hidden">
		<ui.PopoverTrigger
			aria-label="Open documentation navigation"
			class="inline-flex h-8 items-center gap-2 rounded-md px-2 text-sm font-medium text-muted-foreground transition-colors hover:bg-accent hover:text-accent-foreground"
		>
			<icon.Menu class="size-4"/>
			<span>Docs</span>
		</ui.PopoverTrigger>
		<ui.PopoverContent class="max-h-[calc(100svh-5rem)] w-56 overflow-y-auto p-3">
			<nav aria-label="Documentation navigation">
				<docsNavigation active={active}/>
			</nav>
		</ui.PopoverContent>
	</ui.Popover>
}

// siteLayout is the shared page shell. Its explicit mode selects the page's
// spatial responsibilities without inferring them from the request path.
// active names the documentation navigation entry to highlight.
//
// Doc search: CommandDialog owns both its named trigger and content within one
// Dialog root, matching dialog.js's nearest-root ownership. command.js's global
// Cmd-K/Ctrl-K hotkey toggles the same dialog. The search index is the registry
// component list plus the static pages — derived, no manual list to drift.
component siteLayout(title string, active string, mode layoutMode, toc []docTOCItem, children gsx.Node) {
	{{
		headerContainerClass := "mx-auto flex h-14 items-center justify-between"
		headerClass := "sticky top-0 z-10 border-b border-border bg-background/95 backdrop-blur"
		contentContainerClass := "mx-auto w-full py-10"
		mainClass := "min-w-0 flex-1"
		bodyClass := "min-h-svh bg-background text-foreground antialiased"
		footerContainerClass := "mx-auto px-4 py-6 text-sm text-muted-foreground"

		switch mode {
		case layoutDocs:
			headerContainerClass += " max-w-[1568px] px-4 sm:px-6 lg:px-8"
			contentContainerClass += " max-w-[1568px] grid grid-cols-1 px-4 sm:px-6 lg:grid-cols-[288px_minmax(0,1fr)] lg:px-8 xl:grid-cols-[288px_minmax(0,1fr)_288px]"
			mainClass = "mx-auto w-full min-w-0 max-w-[640px]"
			footerContainerClass += " max-w-[1568px] sm:px-6 lg:px-8"
		case layoutWorkspace:
			headerContainerClass += " max-w-none px-4"
			headerClass += " shrink-0"
			contentContainerClass += " flex max-w-none px-4 lg:min-h-0 lg:flex-1 lg:py-4"
			mainClass += " lg:min-h-0 lg:overflow-y-auto"
			bodyClass += " lg:flex lg:h-svh lg:flex-col lg:overflow-hidden"
		case layoutMarketing:
			headerContainerClass += " max-w-6xl px-4"
			contentContainerClass += " flex max-w-6xl px-4"
			footerContainerClass += " max-w-6xl"
		}
	}}
	<!DOCTYPE html>
	<html lang="en">
		<siteHead title={title} entry="web/main.js"/>
		<body
			data-site-layout={mode}
			class={bodyClass}
		>
			<header class={headerClass}>
				<div class={headerContainerClass}>
					<div class="flex items-center gap-2">
						<a href={Home{} |> url} class="flex items-center">
							<siteLogo/>
						</a>
						{ if mode == layoutDocs {
							<compactDocsNavigation active={active}/>
						} }
					</div>
					<nav class="flex items-center gap-4">
						<ui.CommandDialog
							title="Search documentation"
							description="Search components and pages..."
							trigger={ <ui.DialogTrigger
								class="hidden h-8 w-56 items-center gap-2 rounded-lg border bg-muted/50 px-2.5 text-sm text-muted-foreground transition-colors hover:bg-muted sm:inline-flex"
							>
								<icon.Search class="size-4"/>
								<span class="flex-1 text-left">Search docs...</span>
								<ui.Kbd>⌘K</ui.Kbd>
							</ui.DialogTrigger> }
						>
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
						<a
							href={GettingStarted{} |> url}
							class="text-sm text-muted-foreground transition-colors hover:text-foreground"
						>
							Docs
						</a>
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
			<div class={contentContainerClass}>
				{ if mode == layoutDocs {
					<aside data-site-docs-sidebar class="hidden min-w-0 lg:block">
						<nav
							aria-label="Documentation navigation"
							class="sticky top-24 max-h-[calc(100svh-7rem)] overflow-y-auto pb-1 pr-16"
						>
							<docsNavigation active={active}/>
						</nav>
					</aside>
				} }
				<main
					{ if mode == layoutDocs {
						data-site-docs-article
					} }
					class={mainClass}
				>
					{ children }
				</main>
				{ if mode == layoutDocs && len(toc) > 0 {
					<aside data-site-docs-toc class="hidden min-w-0 xl:block">
						<div class="sticky top-24 max-h-[calc(100svh-7rem)] overflow-y-auto pb-1 pl-16">
							<docTableOfContents items={toc}/>
						</div>
					</aside>
				} }
			</div>
			{ if mode != layoutWorkspace {
				<footer data-site-footer class="border-t border-border">
					<div class={footerContainerClass}>
						gsxui — shadcn-style components for gsx. Copy-in, type-checked, server-rendered.
					</div>
				</footer>
			} }
			{/* Mounted once per page: the bottom-right region ui/toaster.js
			   appends every client-constructed toast <li> into. */}
			<ui.Toaster/>
		</body>
	</html>
}
