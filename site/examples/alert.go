package examples

import examplealert "github.com/gsxhq/gsxui/site/examples/alert"

func init() {
	Register("alert", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       examplealert.Basic(),
		SourcePath: "alert/basic.gsx",
	})
	Register("alert", Example{
		Name:       "variants",
		Title:      "Variants",
		Node:       examplealert.Variants(),
		SourcePath: "alert/variants.gsx",
	})
}
