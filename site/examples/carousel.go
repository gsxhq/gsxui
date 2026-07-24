package examples

import examplecarousel "github.com/gsxhq/gsxui/site/examples/carousel"

func init() {
	Register("carousel", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       examplecarousel.Basic(),
		SourcePath: "carousel/basic.gsx",
	})
	Register("carousel", Example{
		Name:       "sizes",
		Title:      "Sizes",
		Node:       examplecarousel.Sizes(),
		SourcePath: "carousel/sizes.gsx",
	})
	Register("carousel", Example{
		Name:       "api",
		Title:      "API",
		Node:       examplecarousel.Api(),
		SourcePath: "carousel/api.gsx",
	})
}
