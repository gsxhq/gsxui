package examples

import exampletextarea "github.com/gsxhq/gsxui/site/examples/textarea"

func init() {
	Register("textarea", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       exampletextarea.Basic(),
		SourcePath: "textarea/basic.gsx",
	})
	Register("textarea", Example{
		Name:       "states",
		Title:      "States",
		Node:       exampletextarea.States(),
		SourcePath: "textarea/states.gsx",
	})
}
