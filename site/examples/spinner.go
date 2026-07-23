package examples

import examplespinner "github.com/gsxhq/gsxui/site/examples/spinner"

func init() {
	Register("spinner", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       examplespinner.Basic(),
		SourcePath: "spinner/basic.gsx",
	})
}
