package examples

import exampleempty "github.com/gsxhq/gsxui/site/examples/empty"

func init() {
	Register("empty", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       exampleempty.Basic(),
		SourcePath: "empty/basic.gsx",
	})
	Register("empty", Example{
		Name:       "rtl",
		Title:      "RTL",
		Node:       exampleempty.Rtl(),
		SourcePath: "empty/rtl.gsx",
	})
}
