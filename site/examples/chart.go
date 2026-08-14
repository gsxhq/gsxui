package examples

import examplechart "github.com/gsxhq/gsxui/site/examples/chart"

func init() {
	Register("chart", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       examplechart.Basic(),
		SourcePath: "chart/basic.gsx",
	})
}
