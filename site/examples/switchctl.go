package examples

import exampleswitch "github.com/gsxhq/gsxui/site/examples/switchctl"

func init() {
	Register("switchctl", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       exampleswitch.Basic(),
		SourcePath: "switchctl/basic.gsx",
	})
	Register("switchctl", Example{
		Name:       "states",
		Title:      "States",
		Node:       exampleswitch.States(),
		SourcePath: "switchctl/states.gsx",
	})
}
