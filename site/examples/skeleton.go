package examples

import exampleskeleton "github.com/gsxhq/gsxui/site/examples/skeleton"

func init() {
	Register("skeleton", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       exampleskeleton.Basic(),
		SourcePath: "skeleton/basic.gsx",
	})
	Register("skeleton", Example{
		Name:       "card",
		Title:      "Card composition",
		Node:       exampleskeleton.Card(),
		SourcePath: "skeleton/card.gsx",
	})
}
