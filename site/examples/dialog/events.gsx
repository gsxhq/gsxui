package dialog

import (
	uibutton "github.com/gsxhq/gsxui/ui/button"
	uidialog "github.com/gsxhq/gsxui/ui/dialog"
)

// Events shows the gsxui:open/gsxui:close CustomEvents dialog.js emits on
// the <dialog> element on every open/close path — trigger, Esc, light
// dismiss, and programmatic showModal()/close() alike.
component Events() {
	<uidialog.Dialog>
		<uibutton.Button data-gsxui-dialog-trigger>Open</uibutton.Button>
		<uidialog.DialogContent id="events-dialog">
			<uidialog.DialogTitle>Watched dialog</uidialog.DialogTitle>
			<uidialog.DialogDescription>Its open/close events log below.</uidialog.DialogDescription>
		</uidialog.DialogContent>
		<output id="events-log" class="mt-4 block text-sm text-muted-foreground">closed</output>
		<script>
			document.addEventListener("gsxui:open", (e) => {
				if (e.target.id === "events-dialog") document.getElementById("events-log").textContent = "open";
			});
			document.addEventListener("gsxui:close", (e) => {
				if (e.target.id === "events-dialog") document.getElementById("events-log").textContent = "closed";
			});
		</script>
	</uidialog.Dialog>
}
