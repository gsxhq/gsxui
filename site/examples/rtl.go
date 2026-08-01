package examples

import examplertl "github.com/gsxhq/gsxui/site/examples/rtl"

func init() {
	Register("rtl", Example{
		Name:       "login",
		Title:      "Login",
		Node:       examplertl.Login(),
		SourcePath: "rtl/login.gsx",
	})
}
