package examples

import examplebutton "github.com/gsxhq/gsxui/site/examples/button"

// Registered here (rather than inline in registry.go) so each component's
// example batch is a self-contained addition — Task 3's batches follow the
// same pattern, one file per component's Register calls.
func init() {
	Register("button", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       examplebutton.Basic(),
		SourcePath: "button/basic.gsx",
	})
	Register("button", Example{
		Name:       "variants",
		Title:      "Variants",
		Node:       examplebutton.Variants(),
		SourcePath: "button/variants.gsx",
	})
}
