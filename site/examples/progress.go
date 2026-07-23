package examples

import exampleprogress "github.com/gsxhq/gsxui/site/examples/progress"

func init() {
	Register("progress", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       exampleprogress.Basic(),
		SourcePath: "progress/basic.gsx",
	})
}
