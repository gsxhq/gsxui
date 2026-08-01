package examples

import exampleaspectratio "github.com/gsxhq/gsxui/site/examples/aspectratio"

func init() {
	Register("aspect-ratio", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       exampleaspectratio.Basic(),
		SourcePath: "aspectratio/basic.gsx",
	})
	Register("aspect-ratio", Example{
		Name:       "rtl",
		Title:      "RTL",
		Node:       exampleaspectratio.Rtl(),
		SourcePath: "aspectratio/rtl.gsx",
	})
}
