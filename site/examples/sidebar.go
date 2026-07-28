package examples

import examplesidebar "github.com/gsxhq/gsxui/site/examples/sidebar"

func init() {
	Register("sidebar", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       examplesidebar.Basic(),
		SourcePath: "sidebar/basic.gsx",
		Isolated:   true,
	})
	Register("sidebar", Example{
		Name:       "variants",
		Title:      "Variants",
		Node:       examplesidebar.Variants(),
		SourcePath: "sidebar/variants.gsx",
		Isolated:   true,
		Previews: []Preview{
			{Name: "sidebar", Title: "variant=sidebar (default)", Node: examplesidebar.Variants()},
			{Name: "floating", Title: "variant=floating", Node: examplesidebar.VariantFloating()},
			{Name: "inset", Title: "variant=inset", Node: examplesidebar.VariantInset()},
			{Name: "offcanvas", Title: "collapsible=offcanvas (default)", Node: examplesidebar.CollapsibleOffcanvas()},
			{Name: "icon", Title: "collapsible=icon", Node: examplesidebar.CollapsibleIcon()},
			{Name: "none", Title: "collapsible=none (always expanded, no rail)", Node: examplesidebar.CollapsibleNone()},
			{Name: "right-collapsed", Title: "right side, collapsed", Node: examplesidebar.RightCollapsed()},
			{Name: "icon-collapsed", Title: "icon collapsed", Node: examplesidebar.IconCollapsed()},
		},
	})
	Register("sidebar", Example{
		Name:       "persisted",
		Title:      "Persisted (cookie round-trip)",
		Node:       examplesidebar.Persisted(),
		SourcePath: "sidebar/persisted.gsx",
		Isolated:   true,
	})
}
