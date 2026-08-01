package examples

import examplebreadcrumb "github.com/gsxhq/gsxui/site/examples/breadcrumb"

func init() {
	Register("breadcrumb", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       examplebreadcrumb.Basic(),
		SourcePath: "breadcrumb/basic.gsx",
	})
	Register("breadcrumb", Example{
		Name:       "rtl",
		Title:      "RTL",
		Node:       examplebreadcrumb.Rtl(),
		SourcePath: "breadcrumb/rtl.gsx",
	})
}
