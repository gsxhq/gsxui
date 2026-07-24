package sonner

import "github.com/gsxhq/gsxui/ui"

// Basic mirrors shadcn's own sonner-demo.tsx: a single outline button that
// fires a toast with a description and an "Undo" action. It uses the
// declarative data-gsxui-toast trigger (no page script needed) — the same
// zero-JS idiom as the dialog trigger. The toast itself is built by
// ui/sonner.js and mounted into the layout's <ui.Toaster/> region.
component Basic() {
	<ui.Button
		variant="outline"
		data-gsxui-toast="Event has been created"
		data-gsxui-toast-description="Sunday, December 03, 2023 at 9:00 AM"
		data-gsxui-toast-action="Undo"
	>
		Show Toast
	</ui.Button>
}
