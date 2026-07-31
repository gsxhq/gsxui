package pages

import "github.com/gsxhq/vite"

// siteHead is the canonical document head for full site pages and isolated
// live previews. entry keeps the shared asset/theme bootstrap while allowing
// previews to load only their handshake code.
component siteHead(title string, entry string) {
	<head>
		<meta charset="UTF-8"/>
		<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
		<title>{ title } · gsxui</title>
		{/* Favicon set — embedded in the binary and served by site/main.go at
		   root paths. SVG for modern browsers, 32px PNG for the ones that
		   ignore SVG favicons (Safari), 180px for iOS home-screen. */}
		<link rel="icon" href="/favicon.svg" type="image/svg+xml"/>
		<link rel="icon" href="/favicon-32.png" sizes="32x32" type="image/png"/>
		<link rel="apple-touch-icon" href="/apple-touch-icon.png"/>
		<script>
			// Theme init — runs before first paint (blocking head script) so
			// a stored dark preference never flashes light. Explicit choice
			// (localStorage) wins; otherwise follow the OS preference. The
			// header's data-site-theme-toggle button (web/site.js) flips the
			// class and stores the choice.
			try {
				var gsxuiTheme = localStorage.getItem("gsxui-theme");
				if (gsxuiTheme === "dark" || (!gsxuiTheme && matchMedia("(prefers-color-scheme: dark)").matches)) {
					document.documentElement.classList.add("dark");
				}
			} catch (e) {}
		</script>
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
		{{ assets := v.Entry(entry) }}
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
}
