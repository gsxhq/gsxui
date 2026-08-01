package examples

import exampletable "github.com/gsxhq/gsxui/site/examples/table"

func init() {
	Register("table", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       exampletable.Basic(),
		SourcePath: "table/basic.gsx",
	})
	Register("table", Example{
		Name:       "data",
		Title:      "Data table",
		Node:       exampletable.Data(),
		SourcePath: "table/data.gsx",
	})
	Register("table", Example{
		Name:       "rtl",
		Title:      "RTL",
		Node:       exampletable.Rtl(),
		SourcePath: "table/rtl.gsx",
	})
}
