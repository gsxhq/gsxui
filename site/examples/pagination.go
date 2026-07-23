package examples

import examplepagination "github.com/gsxhq/gsxui/site/examples/pagination"

func init() {
	Register("pagination", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       examplepagination.Basic(),
		SourcePath: "pagination/basic.gsx",
	})
}
