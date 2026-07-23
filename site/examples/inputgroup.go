package examples

import exampleinputgroup "github.com/gsxhq/gsxui/site/examples/inputgroup"

func init() {
	Register("input-group", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       exampleinputgroup.Basic(),
		SourcePath: "inputgroup/basic.gsx",
	})
}
