package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Button is a single-element component: one root slot carrying both dimensions.
// Its dimensions mirror the public variant and size parameters of button.gsx.
var Button = recipe.Shape{
	Component: "button",
	Slots: []recipe.Slot{{
		Name: "",
		Base: true,
		Dimensions: []recipe.Dimension{
			{Name: "variant", Default: "default", Values: []string{
				"default", "destructive", "outline", "secondary", "ghost", "link"}},
			{Name: "size", Default: "default", Values: []string{
				"default", "xs", "sm", "lg", "icon", "icon-xs", "icon-sm", "icon-lg"}},
		},
	}},
}
