package examples

import examplemenubar "github.com/gsxhq/gsxui/site/examples/menubar"

func init() {
	Register("menubar", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       examplemenubar.Basic(),
		SourcePath: "menubar/basic.gsx",
	})
	Register("menubar", Example{
		Name:       "full",
		Title:      "Full menu with checkboxes, radios, and a submenu",
		Node:       examplemenubar.Full(),
		SourcePath: "menubar/full.gsx",
	})
}
