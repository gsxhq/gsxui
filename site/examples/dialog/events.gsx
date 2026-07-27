package dialog

import (
	"github.com/gsxhq/gsxui/ui"
)

// Events shows the gsxui:open/gsxui:close CustomEvents dialog.js emits on
// the <dialog> element on every open/close path — trigger, Esc, light
// dismiss, and programmatic showModal()/close() alike.
component Events() {
	<ui.Dialog>
		<ui.Button commandfor="events-dialog" command="show-modal">Open</ui.Button>
		<ui.DialogContent id="events-dialog">
			<ui.DialogTitle>Watched dialog</ui.DialogTitle>
			<ui.DialogDescription>Its open/close events log below.</ui.DialogDescription>
			<ui.Button
				variant="outline"
				commandfor="events-dialog"
				command="request-close"
			>
				Close
			</ui.Button>
		</ui.DialogContent>
		<output id="events-log" class="mt-4 block text-sm text-muted-foreground">closed</output>
		<script>
			document.addEventListener("gsxui:open", (e) => {
				if (e.target.id === "events-dialog") document.getElementById("events-log").textContent = "open";
			});
			document.addEventListener("gsxui:close", (e) => {
				if (e.target.id === "events-dialog") document.getElementById("events-log").textContent = "closed";
			});
		</script>
	</ui.Dialog>
}
