package examples

import examplealertdialog "github.com/gsxhq/gsxui/site/examples/alertdialog"

func init() {
	Register("alert-dialog", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       examplealertdialog.Basic(),
		SourcePath: "alertdialog/basic.gsx",
	})
	Register("alert-dialog", Example{
		Name:       "rtl",
		Title:      "RTL",
		Node:       examplealertdialog.Rtl(),
		SourcePath: "alertdialog/rtl.gsx",
	})
}
