package examples

import examplehovercard "github.com/gsxhq/gsxui/site/examples/hovercard"

func init() {
	Register("hover-card", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       examplehovercard.Basic(),
		SourcePath: "hovercard/basic.gsx",
	})
	Register("hover-card", Example{
		Name:       "rtl",
		Title:      "RTL",
		Node:       examplehovercard.Rtl(),
		SourcePath: "hovercard/rtl.gsx",
	})
}
