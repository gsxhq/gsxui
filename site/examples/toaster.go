package examples

import exampletoaster "github.com/gsxhq/gsxui/site/examples/toaster"

func init() {
	Register("toaster", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       exampletoaster.Basic(),
		SourcePath: "toaster/basic.gsx",
	})
	Register("toaster", Example{
		Name:       "types",
		Title:      "Types",
		Node:       exampletoaster.Types(),
		SourcePath: "toaster/types.gsx",
	})
}
