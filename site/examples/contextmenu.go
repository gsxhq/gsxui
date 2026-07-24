package examples

import examplecontextmenu "github.com/gsxhq/gsxui/site/examples/contextmenu"

func init() {
	Register("context-menu", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       examplecontextmenu.Basic(),
		SourcePath: "contextmenu/basic.gsx",
	})
}
