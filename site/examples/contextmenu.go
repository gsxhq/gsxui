package examples

import examplecontextmenu "github.com/gsxhq/gsxui/site/examples/contextmenu"

func init() {
	Register("context-menu", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       examplecontextmenu.Basic(),
		SourcePath: "contextmenu/basic.gsx",
	})
	Register("context-menu", Example{
		Name:       "full",
		Title:      "Checkboxes, radios, and a submenu",
		Node:       examplecontextmenu.Full(),
		SourcePath: "contextmenu/full.gsx",
	})
	Register("context-menu", Example{
		Name:       "rtl",
		Title:      "RTL",
		Node:       examplecontextmenu.Rtl(),
		SourcePath: "contextmenu/rtl.gsx",
	})
}
