package examples

import exampletooltip "github.com/gsxhq/gsxui/site/examples/tooltip"

func init() {
	Register("tooltip", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       exampletooltip.Basic(),
		SourcePath: "tooltip/basic.gsx",
	})
	Register("tooltip", Example{
		Name:       "wide",
		Title:      "Custom width",
		Node:       exampletooltip.Wide(),
		SourcePath: "tooltip/wide.gsx",
	})
}
