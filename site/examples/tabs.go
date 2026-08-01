package examples

import exampletabs "github.com/gsxhq/gsxui/site/examples/tabs"

func init() {
	Register("tabs", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       exampletabs.Basic(),
		SourcePath: "tabs/basic.gsx",
	})
	Register("tabs", Example{
		Name:       "icons",
		Title:      "With icons",
		Node:       exampletabs.Icons(),
		SourcePath: "tabs/icons.gsx",
	})
	Register("tabs", Example{
		Name:       "rtl",
		Title:      "RTL",
		Node:       exampletabs.Rtl(),
		SourcePath: "tabs/rtl.gsx",
	})
}
