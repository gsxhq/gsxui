package examples

import exampledrawer "github.com/gsxhq/gsxui/site/examples/drawer"

func init() {
	Register("drawer", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       exampledrawer.Basic(),
		SourcePath: "drawer/basic.gsx",
	})
	Register("drawer", Example{
		Name:       "directions",
		Title:      "Directions",
		Node:       exampledrawer.Directions(),
		SourcePath: "drawer/directions.gsx",
	})
}
