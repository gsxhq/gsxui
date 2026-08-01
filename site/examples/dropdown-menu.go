package examples

import exampledropdown "github.com/gsxhq/gsxui/site/examples/dropdown"

func init() {
	Register("dropdown-menu", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       exampledropdown.Basic(),
		SourcePath: "dropdown/basic.gsx",
	})
	Register("dropdown-menu", Example{
		Name:       "destructive",
		Title:      "Destructive and disabled items",
		Node:       exampledropdown.Destructive(),
		SourcePath: "dropdown/destructive.gsx",
	})
	Register("dropdown-menu", Example{
		Name:       "checkboxes",
		Title:      "Checkbox items",
		Node:       exampledropdown.Checkboxes(),
		SourcePath: "dropdown/checkboxes.gsx",
	})
	Register("dropdown-menu", Example{
		Name:       "radios",
		Title:      "Radio items",
		Node:       exampledropdown.Radios(),
		SourcePath: "dropdown/radios.gsx",
	})
	Register("dropdown-menu", Example{
		Name:       "submenu",
		Title:      "Submenu",
		Node:       exampledropdown.Submenu(),
		SourcePath: "dropdown/submenu.gsx",
	})
	Register("dropdown-menu", Example{
		Name:       "rtl",
		Title:      "RTL",
		Node:       exampledropdown.Rtl(),
		SourcePath: "dropdown/rtl.gsx",
	})
}
