package examples

import examplekbd "github.com/gsxhq/gsxui/site/examples/kbd"

func init() {
	Register("kbd", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       examplekbd.Basic(),
		SourcePath: "kbd/basic.gsx",
	})
}
