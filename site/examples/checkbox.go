package examples

import examplecheckbox "github.com/gsxhq/gsxui/site/examples/checkbox"

func init() {
	Register("checkbox", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       examplecheckbox.Basic(),
		SourcePath: "checkbox/basic.gsx",
	})
	Register("checkbox", Example{
		Name:       "states",
		Title:      "States",
		Node:       examplecheckbox.States(),
		SourcePath: "checkbox/states.gsx",
	})
}
