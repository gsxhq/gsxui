package examples

import exampleinput "github.com/gsxhq/gsxui/site/examples/input"

func init() {
	Register("input", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       exampleinput.Basic(),
		SourcePath: "input/basic.gsx",
	})
	Register("input", Example{
		Name:       "states",
		Title:      "States",
		Node:       exampleinput.States(),
		SourcePath: "input/states.gsx",
	})
	Register("input", Example{
		Name:       "custom",
		Title:      "Custom class",
		Node:       exampleinput.Custom(),
		SourcePath: "input/custom.gsx",
	})
	Register("input", Example{
		Name:       "form-row",
		Title:      "Form row",
		Node:       exampleinput.FormRow(),
		SourcePath: "input/form-row.gsx",
	})
}
