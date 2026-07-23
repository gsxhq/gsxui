package examples

import exampleseparator "github.com/gsxhq/gsxui/site/examples/separator"

func init() {
	Register("separator", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       exampleseparator.Basic(),
		SourcePath: "separator/basic.gsx",
	})
	Register("separator", Example{
		Name:       "orientation",
		Title:      "Orientation",
		Node:       exampleseparator.Orientation(),
		SourcePath: "separator/orientation.gsx",
	})
}
