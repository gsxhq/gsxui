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
}
