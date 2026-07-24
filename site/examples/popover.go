package examples

import examplepopover "github.com/gsxhq/gsxui/site/examples/popover"

func init() {
	Register("popover", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       examplepopover.Basic(),
		SourcePath: "popover/basic.gsx",
	})
}
