package examples

import exampleinputotp "github.com/gsxhq/gsxui/site/examples/inputotp"

func init() {
	Register("input-otp", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       exampleinputotp.Basic(),
		SourcePath: "inputotp/basic.gsx",
	})
	Register("input-otp", Example{
		Name:       "pattern",
		Title:      "Pattern",
		Node:       exampleinputotp.Pattern(),
		SourcePath: "inputotp/pattern.gsx",
	})
	Register("input-otp", Example{
		Name:       "separator",
		Title:      "Separator",
		Node:       exampleinputotp.Separator(),
		SourcePath: "inputotp/separator.gsx",
	})
}
