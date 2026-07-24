package examples

import examplecommand "github.com/gsxhq/gsxui/site/examples/command"

func init() {
	Register("command", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       examplecommand.Basic(),
		SourcePath: "command/basic.gsx",
	})
}
