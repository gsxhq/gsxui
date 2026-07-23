package examples

import examplesheet "github.com/gsxhq/gsxui/site/examples/sheet"

func init() {
	Register("sheet", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       examplesheet.Basic(),
		SourcePath: "sheet/basic.gsx",
	})
}
