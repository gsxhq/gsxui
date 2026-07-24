package examples

import examplenativeselect "github.com/gsxhq/gsxui/site/examples/nativeselect"

func init() {
	Register("native-select", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       examplenativeselect.Basic(),
		SourcePath: "nativeselect/basic.gsx",
	})
	Register("native-select", Example{
		Name:       "groups",
		Title:      "Groups",
		Node:       examplenativeselect.Groups(),
		SourcePath: "nativeselect/groups.gsx",
	})
}
