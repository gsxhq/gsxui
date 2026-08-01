package examples

import examplecombobox "github.com/gsxhq/gsxui/site/examples/combobox"

func init() {
	Register("combobox", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       examplecombobox.Basic(),
		SourcePath: "combobox/basic.gsx",
	})
	Register("combobox", Example{
		Name:       "groups",
		Title:      "Groups",
		Node:       examplecombobox.Groups(),
		SourcePath: "combobox/groups.gsx",
	})
	Register("combobox", Example{
		Name:       "clear",
		Title:      "Clear",
		Node:       examplecombobox.Clear(),
		SourcePath: "combobox/clear.gsx",
	})
	Register("combobox", Example{
		Name:       "trigger-clear",
		Title:      "Trigger and clear",
		Node:       examplecombobox.TriggerClear(),
		SourcePath: "combobox/trigger-clear.gsx",
	})
	Register("combobox", Example{
		Name:       "form",
		Title:      "Form",
		Node:       examplecombobox.Form(),
		SourcePath: "combobox/form.gsx",
	})
	Register("combobox", Example{
		Name:       "rtl",
		Title:      "RTL",
		Node:       examplecombobox.Rtl(),
		SourcePath: "combobox/rtl.gsx",
	})
}
