package sonner

import "github.com/gsxhq/gsxui/ui"

// Server demonstrates the real-world flash-message pattern: the SERVER
// renders a toast card and the browser adopts it into the same stacking /
// timer / dismiss lifecycle as a JS-triggered toast — with zero HTMX-specific
// code. Two paths are shown:
//
//  1. A static ui.Toast rendered inline (the markup showcase). On a full page
//     load a server would drain its session/request flashes into <ui.Toaster/>
//     exactly like this; ui/sonner.js adopts every li[data-slot="toast"]
//     present at init.
//  2. A button + inline <script> that clones a pre-rendered server row (the
//     ui.Toast wrapped in the <template data-server-flash-demo> below) and
//     appends it into #gsxui-toaster — precisely what an HTMX out-of-band
//     swap (hx-swap-oob="beforeend:#gsxui-toaster") does on the wire. The
//     MutationObserver on the <ol> adopts it; no page-specific lifecycle code.
//
// This is the one-viewport-per-page model: the server is the single source of
// toast markup, ui.Toaster is mounted once, and appends flow in from anywhere.
component Server() {
	<div class="flex flex-col items-start gap-4">
		<ui.Toast toastType="success" title="Profile updated" description="Your changes have been saved."/>
		<ui.Button variant="outline" id="sonner-server-flash-btn">Append a server flash</ui.Button>
		<template data-server-flash-demo>
			<ui.Toast
				toastType="info"
				title="New message"
				description="A server-rendered row, appended like an HTMX OOB swap."
			/>
		</template>
		<script>
			document.getElementById("sonner-server-flash-btn").addEventListener("click", () => {
				const tpl = document.querySelector("template[data-server-flash-demo]");
				const viewport = document.getElementById("gsxui-toaster");
				if (!tpl || !viewport) return;
				viewport.appendChild(tpl.content.firstElementChild.cloneNode(true));
			});
		</script>
	</div>
}
