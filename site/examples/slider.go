package examples

import exampleslider "github.com/gsxhq/gsxui/site/examples/slider"

func init() {
	Register("slider", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       exampleslider.Basic(),
		SourcePath: "slider/basic.gsx",
	})
	Register("slider", Example{
		Name:       "rtl",
		Title:      "RTL",
		Node:       exampleslider.Rtl(),
		SourcePath: "slider/rtl.gsx",
	})
}
