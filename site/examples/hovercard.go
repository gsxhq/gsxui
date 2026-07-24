package examples

import examplehovercard "github.com/gsxhq/gsxui/site/examples/hovercard"

func init() {
	Register("hover-card", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       examplehovercard.Basic(),
		SourcePath: "hovercard/basic.gsx",
	})
}
