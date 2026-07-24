package examples

import exampledrawer "github.com/gsxhq/gsxui/site/examples/drawer"

func init() {
	Register("drawer", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       exampledrawer.Basic(),
		SourcePath: "drawer/basic.gsx",
	})
}
