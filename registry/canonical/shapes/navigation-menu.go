package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// NavigationMenu owns its own dedicated stylesheet
// (assets/css/styles/default/navigation-menu.css) except for the discrete
// popover-transition block it shares with DropdownMenu/ContextMenu/Menubar/
// the not-yet-migrated Combobox in menu-popover-transition.css. Every marked
// element has real CSS, so every one is declared here.
var NavigationMenu = recipe.Shape{
	Component: "navigation-menu",
	Slots: []recipe.Slot{
		{Name: "", Base: true},
		{Name: "list", Base: true},
		{Name: "item", Base: true},
		{Name: "trigger", Base: true},
		{Name: "trigger-icon", Base: true},
		{Name: "content", Base: true},
		{
			// variant="trigger" shares its ENTIRE class string, verbatim, with the
			// Trigger slot's own base rule (shadcn's own navigationMenuTriggerStyle,
			// reused for a stand-alone top-level link rendered with no dropdown —
			// see ui/navigation-menu.gsx's own doc comment). Duplicated here rather
			// than composed by calling navigationMenu.Trigger() from inside a
			// conditional class expression (no precedent for that in this codebase's
			// .gsx dialect); a verbatim copy within the SAME component/shape is not
			// the cross-component restatement spec section 10b's sibling-precedent
			// warns against, since there is no other, unmigrated owner to race.
			Name: "link", Base: true,
			Dimensions: []recipe.Dimension{
				{Name: "variant", Default: "default", Values: []string{"default", "trigger"}},
				{Name: "active", Default: "false", Values: []string{"false", "true"}},
			},
		},
		{Name: "indicator", Base: true},
		{Name: "indicator-arrow", Base: true},
	},
}
