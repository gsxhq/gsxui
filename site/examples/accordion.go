package examples

import exampleaccordion "github.com/gsxhq/gsxui/site/examples/accordion"

func init() {
	Register("accordion", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       exampleaccordion.Basic(),
		SourcePath: "accordion/basic.gsx",
	})
	Register("accordion", Example{
		Name:       "compact",
		Title:      "Compact",
		Node:       exampleaccordion.Compact(),
		SourcePath: "accordion/compact.gsx",
	})
	Register("accordion", Example{
		Name:       "rtl",
		Title:      "RTL",
		Node:       exampleaccordion.Rtl(),
		SourcePath: "accordion/rtl.gsx",
	})
}
