package examples

import exampledialog "github.com/gsxhq/gsxui/site/examples/dialog"

func init() {
	Register("dialog", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       exampledialog.Basic(),
		SourcePath: "dialog/basic.gsx",
	})
	Register("dialog", Example{
		Name:       "footer",
		Title:      "Footer close button",
		Node:       exampledialog.Footer(),
		SourcePath: "dialog/footer.gsx",
	})
	Register("dialog", Example{
		Name:       "events",
		Title:      "Events",
		Node:       exampledialog.Events(),
		SourcePath: "dialog/events.gsx",
	})
}
