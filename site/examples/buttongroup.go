package examples

import examplebuttongroup "github.com/gsxhq/gsxui/site/examples/buttongroup"

func init() {
	Register("button-group", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       examplebuttongroup.Basic(),
		SourcePath: "buttongroup/basic.gsx",
	})
}
