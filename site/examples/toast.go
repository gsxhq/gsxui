package examples

import exampletoast "github.com/gsxhq/gsxui/site/examples/toast"

func init() {
	// Toast is registered separately from Toaster because the card is a
	// standalone component: this example renders <ui.Toast> directly, both
	// as inert inline markup and as the server-appended row an HTMX
	// out-of-band swap delivers.
	Register("toast", Example{
		Name:       "server",
		Title:      "Server flashes",
		Node:       exampletoast.Server(),
		SourcePath: "toast/server.gsx",
	})
}
