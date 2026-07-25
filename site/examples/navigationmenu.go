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
		Title:      "Mega menu with independently-sized panels",
		Node:       examplenavigationmenu.Mega(),
		SourcePath: "navigationmenu/mega.gsx",
	})
}
