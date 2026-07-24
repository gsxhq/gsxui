package examples

import exampletoggle "github.com/gsxhq/gsxui/site/examples/toggle"

func init() {
	Register("toggle", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       exampletoggle.Basic(),
		SourcePath: "toggle/basic.gsx",
	})
}
