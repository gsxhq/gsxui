package examples

import examplenavigationmenu "github.com/gsxhq/gsxui/site/examples/navigationmenu"

func init() {
	Register("navigation-menu", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       examplenavigationmenu.Basic(),
		SourcePath: "navigationmenu/basic.gsx",
	})
	Register("navigation-menu", Example{
		Name:       "mega",
		Title:      "Mega menu with a shared, JS-measured viewport",
		Node:       examplenavigationmenu.Mega(),
		SourcePath: "navigationmenu/mega.gsx",
	})
}
