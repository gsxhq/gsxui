package examples

import exampleradio "github.com/gsxhq/gsxui/site/examples/radio"

func init() {
	Register("radio", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       exampleradio.Basic(),
		SourcePath: "radio/basic.gsx",
	})
	Register("radio", Example{
		Name:       "states",
		Title:      "States",
		Node:       exampleradio.States(),
		SourcePath: "radio/states.gsx",
	})
	Register("radio", Example{
		Name:       "rtl",
		Title:      "RTL",
		Node:       exampleradio.Rtl(),
		SourcePath: "radio/rtl.gsx",
	})
}
