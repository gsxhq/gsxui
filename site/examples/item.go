package examples

import exampleitem "github.com/gsxhq/gsxui/site/examples/item"

func init() {
	Register("item", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       exampleitem.Basic(),
		SourcePath: "item/basic.gsx",
	})
	Register("item", Example{
		Name:       "rtl",
		Title:      "RTL",
		Node:       exampleitem.Rtl(),
		SourcePath: "item/rtl.gsx",
	})
}
