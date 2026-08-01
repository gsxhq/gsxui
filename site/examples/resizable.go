package examples

import exampleresizable "github.com/gsxhq/gsxui/site/examples/resizable"

func init() {
	Register("resizable", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       exampleresizable.Basic(),
		SourcePath: "resizable/basic.gsx",
	})
	Register("resizable", Example{
		Name:       "vertical",
		Title:      "Vertical",
		Node:       exampleresizable.Vertical(),
		SourcePath: "resizable/vertical.gsx",
	})
	Register("resizable", Example{
		Name:       "handle",
		Title:      "Handle",
		Node:       exampleresizable.Handle(),
		SourcePath: "resizable/handle.gsx",
	})
	Register("resizable", Example{
		Name:       "rtl",
		Title:      "RTL",
		Node:       exampleresizable.Rtl(),
		SourcePath: "resizable/rtl.gsx",
	})
}
