package examples

import exampleslider "github.com/gsxhq/gsxui/site/examples/slider"

func init() {
	Register("slider", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       exampleslider.Basic(),
		SourcePath: "slider/basic.gsx",
	})
}
