package examples

import examplecard "github.com/gsxhq/gsxui/site/examples/card"

func init() {
	Register("card", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       examplecard.Basic(),
		SourcePath: "card/basic.gsx",
	})
	Register("card", Example{
		Name:       "compound",
		Title:      "Compound",
		Node:       examplecard.Compound(),
		SourcePath: "card/compound.gsx",
	})
	Register("card", Example{
		Name:       "attrs",
		Title:      "Attrs",
		Node:       examplecard.Attrs(),
		SourcePath: "card/attrs.gsx",
	})
}
