package examples

import examplescrollarea "github.com/gsxhq/gsxui/site/examples/scrollarea"

func init() {
	Register("scroll-area", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       examplescrollarea.Basic(),
		SourcePath: "scrollarea/basic.gsx",
	})
	Register("scroll-area", Example{
		Name:       "horizontal",
		Title:      "Horizontal",
		Node:       examplescrollarea.Horizontal(),
		SourcePath: "scrollarea/horizontal.gsx",
	})
}
