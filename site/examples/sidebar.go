package examples

import examplesidebar "github.com/gsxhq/gsxui/site/examples/sidebar"

func init() {
	Register("sidebar", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       examplesidebar.Basic(),
		SourcePath: "sidebar/basic.gsx",
	})
	Register("sidebar", Example{
		Name:       "variants",
		Title:      "Variants",
		Node:       examplesidebar.Variants(),
		SourcePath: "sidebar/variants.gsx",
	})
	Register("sidebar", Example{
		Name:       "persisted",
		Title:      "Persisted (cookie round-trip)",
		Node:       examplesidebar.Persisted(),
		SourcePath: "sidebar/persisted.gsx",
	})
}
