package examples

import examplesonner "github.com/gsxhq/gsxui/site/examples/sonner"

func init() {
	Register("sonner", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       examplesonner.Basic(),
		SourcePath: "sonner/basic.gsx",
	})
	Register("sonner", Example{
		Name:       "types",
		Title:      "Types",
		Node:       examplesonner.Types(),
		SourcePath: "sonner/types.gsx",
	})
}
