package examples

import exampleavatar "github.com/gsxhq/gsxui/site/examples/avatar"

func init() {
	Register("avatar", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       exampleavatar.Basic(),
		SourcePath: "avatar/basic.gsx",
	})
	Register("avatar", Example{
		Name:       "group",
		Title:      "Group",
		Node:       exampleavatar.Group(),
		SourcePath: "avatar/group.gsx",
	})
}
