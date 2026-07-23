package examples

import examplebadge "github.com/gsxhq/gsxui/site/examples/badge"

func init() {
	Register("badge", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       examplebadge.Basic(),
		SourcePath: "badge/basic.gsx",
	})
	Register("badge", Example{
		Name:       "variants",
		Title:      "Variants",
		Node:       examplebadge.Variants(),
		SourcePath: "badge/variants.gsx",
	})
}
