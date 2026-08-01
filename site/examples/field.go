package examples

import examplefield "github.com/gsxhq/gsxui/site/examples/field"

func init() {
	Register("field", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       examplefield.Basic(),
		SourcePath: "field/basic.gsx",
	})
	Register("field", Example{
		Name:       "invalid",
		Title:      "Invalid",
		Node:       examplefield.Invalid(),
		SourcePath: "field/invalid.gsx",
	})
	Register("field", Example{
		Name:       "rtl",
		Title:      "RTL",
		Node:       examplefield.Rtl(),
		SourcePath: "field/rtl.gsx",
	})
}
