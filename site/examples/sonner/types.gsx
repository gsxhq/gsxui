package sonner

import "github.com/gsxhq/gsxui/ui"

// Types mirrors shadcn's own sonner-types.tsx: one button per toast type.
// Default/success/info/warning/error use the declarative data-gsxui-toast
// trigger; the promise button needs the imperative API (toast.promise
// morphs a loading toast in place on settle), reached here through the
// window.gsxui global — inline <script> can't import the barrel, and
// ui/sonner.js exposes window.gsxui.toast for exactly this case.
component Types() {
	<div class="flex flex-wrap gap-2">
		<ui.Button variant="outline" data-gsxui-toast="Event has been created">Default</ui.Button>
		<ui.Button variant="outline" data-gsxui-toast="Changes saved" data-gsxui-toast-type="success">Success</ui.Button>
		<ui.Button variant="outline" data-gsxui-toast="Be advised" data-gsxui-toast-type="info">Info</ui.Button>
		<ui.Button variant="outline" data-gsxui-toast="Careful now" data-gsxui-toast-type="warning">Warning</ui.Button>
		<ui.Button variant="outline" data-gsxui-toast="Something went wrong" data-gsxui-toast-type="error">Error</ui.Button>
		<ui.Button variant="outline" id="sonner-promise-btn">Promise</ui.Button>
		<script>
			document.getElementById("sonner-promise-btn").addEventListener("click", () => {
				window.gsxui.toast.promise(
					new Promise((resolve) => setTimeout(resolve, 2000)),
					{ loading: "Loading...", success: "Settings saved", error: "Failed to save" }
				);
			});
		</script>
	</div>
}
