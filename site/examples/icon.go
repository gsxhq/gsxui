package examples

import exampleicon "github.com/gsxhq/gsxui/site/examples/icon"

func init() {
	Register("icon", Example{
		Name:       "gallery",
		Title:      "Gallery",
		Node:       exampleicon.Gallery(),
		SourcePath: "icon/gallery.gsx",
	})
	Register("icon", Example{
		Name:       "sizes",
		Title:      "Sizes and color",
		Node:       exampleicon.Sizes(),
		SourcePath: "icon/sizes.gsx",
	})
}
