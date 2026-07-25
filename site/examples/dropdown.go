package examples

import exampledropdown "github.com/gsxhq/gsxui/site/examples/dropdown"

func init() {
	Register("dropdown", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       exampledropdown.Basic(),
		SourcePath: "dropdown/basic.gsx",
	})
	Register("dropdown", Example{
		Name:       "destructive",
		Title:      "Destructive and disabled items",
		Node:       exampledropdown.Destructive(),
		SourcePath: "dropdown/destructive.gsx",
	})
	Register("dropdown", Example{
		Name:       "checkboxes",
		Title:      "Checkbox items",
		Node:       exampledropdown.Checkboxes(),
		SourcePath: "dropdown/checkboxes.gsx",
	})
	Register("dropdown", Example{
		Name:       "radios",
		Title:      "Radio items",
		Node:       exampledropdown.Radios(),
		SourcePath: "dropdown/radios.gsx",
	})
	Register("dropdown", Example{
		Name:       "submenu",
		Title:      "Submenu",
		Node:       exampledropdown.Submenu(),
		SourcePath: "dropdown/submenu.gsx",
	})
}
