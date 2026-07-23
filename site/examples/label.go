package examples

import examplelabel "github.com/gsxhq/gsxui/site/examples/label"

func init() {
	Register("label", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       examplelabel.Basic(),
		SourcePath: "label/basic.gsx",
	})
	Register("label", Example{
		Name:       "disabled",
		Title:      "Disabled peer",
		Node:       examplelabel.Disabled(),
		SourcePath: "label/disabled.gsx",
	})
}
