package examples

import examplecollapsible "github.com/gsxhq/gsxui/site/examples/collapsible"

func init() {
	Register("collapsible", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       examplecollapsible.Basic(),
		SourcePath: "collapsible/basic.gsx",
	})
}
