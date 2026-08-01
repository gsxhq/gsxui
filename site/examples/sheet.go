package examples

import examplesheet "github.com/gsxhq/gsxui/site/examples/sheet"

func init() {
	Register("sheet", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       examplesheet.Basic(),
		SourcePath: "sheet/basic.gsx",
	})
	Register("sheet", Example{
		Name:       "directions",
		Title:      "Directions",
		Node:       examplesheet.Directions(),
		SourcePath: "sheet/directions.gsx",
	})
	Register("sheet", Example{
		Name:       "rtl",
		Title:      "RTL",
		Node:       examplesheet.Rtl(),
		SourcePath: "sheet/rtl.gsx",
	})
}
