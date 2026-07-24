package examples

import exampleselectbox "github.com/gsxhq/gsxui/site/examples/selectbox"

func init() {
	Register("select", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       exampleselectbox.Basic(),
		SourcePath: "selectbox/basic.gsx",
	})
	Register("select", Example{
		Name:       "scrollable",
		Title:      "Scrollable",
		Node:       exampleselectbox.Scrollable(),
		SourcePath: "selectbox/scrollable.gsx",
	})
}
