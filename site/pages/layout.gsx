package pages

import (
	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/internal/registry"
	"github.com/gsxhq/vite"
)

// Layout is the shared page shell every site page renders through: header
// (wordmark + GitHub link), sidebar (component list from the registry —
// derived, so it can never drift from what `ui/` actually ships), and
// footer. active names the component whose sidebar entry should highlight;
// pages outside /components/ pass "".
component Layout(title string, active string, children gsx.Node) {
	<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="UTF-8"/>
			<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
			<title>{ title } · gsxui</title>
			{{ v := vite.FromContext(ctx) }}
			{ if v.Dev() {
				<style>
					html[data-loading] body {
						visibility: hidden;
					}

					html[data-loading] * {
						transition: none !important;
					}
				</style>
				<script>
					// Dev-only FOUC gate. Vite injects CSS via JS after the HTML
					// loads, so hide the page until every module script has run
					// (DOMContentLoaded) and one paint has landed (double rAF),
					// then reveal. Prod ships real <link rel=stylesheet> tags
					// below, so no gate is emitted there.
					document.documentElement.dataset.loading = "true";
					var unhide = function () {
						document.documentElement.removeAttribute("data-loading");
					};
					var reveal = function () {
						requestAnimationFrame(function () { requestAnimationFrame(unhide); });
					};
					if (document.readyState === "loading") {
						document.addEventListener("DOMContentLoaded", reveal);
					} else {
						reveal();
					}
					// Safety net (rAF pauses in background tabs).
					setTimeout(unhide, 5000);
				</script>
			} }
			{{ assets := v.Entry("web/main.js") }}
			{ for _, href := range assets.CSS {
				<link rel="stylesheet" href={href}/>
			} }
			{ for _, src := range assets.Preloads {
				<link rel="modulepreload" href={src}/>
			} }
			{ for _, src := range assets.JS {
				<script type="module" src={src}></script>
			} }
		</head>
		<body class="min-h-svh bg-background text-foreground antialiased">
			<header class="sticky top-0 z-10 border-b border-border bg-background/95 backdrop-blur">
				<div class="mx-auto flex h-14 max-w-6xl items-center justify-between px-4">
					<a href={Home{} |> url} class="font-mono text-sm font-semibold tracking-tight">gsxui</a>
					<nav class="flex items-center gap-4">
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
		</body>
	</html>
}
