package examples

import examplebreadcrumb "github.com/gsxhq/gsxui/site/examples/breadcrumb"

func init() {
	Register("breadcrumb", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       examplebreadcrumb.Basic(),
		SourcePath: "breadcrumb/basic.gsx",
	})
}
