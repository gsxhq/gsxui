package examples

import exampletogglegroup "github.com/gsxhq/gsxui/site/examples/togglegroup"

func init() {
	Register("toggle-group", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       exampletogglegroup.Basic(),
		SourcePath: "togglegroup/basic.gsx",
	})
	Register("toggle-group", Example{
		Name:       "single",
		Title:      "Single",
		Node:       exampletogglegroup.Single(),
		SourcePath: "togglegroup/single.gsx",
	})
	Register("toggle-group", Example{
		Name:       "spacing",
		Title:      "Spacing",
		Node:       exampletogglegroup.Spacing(),
		SourcePath: "togglegroup/spacing.gsx",
	})
}
