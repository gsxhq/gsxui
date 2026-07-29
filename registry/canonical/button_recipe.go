package canonical

import "github.com/gsxhq/gsxui/internal/recipe"

// buttonShape is Button's style interface. Its dimensions mirror the public
// variant and size parameters of button.gsx; every style must implement all of
// it. Values are ordered as they appear in the generated switch.
var buttonShape = recipe.Shape{
	Component: "button",
	Base:      true,
	Dimensions: []recipe.Dimension{
		{Name: "variant", Default: "default", Values: []string{
			"default", "destructive", "outline", "secondary", "ghost", "link"}},
		{Name: "size", Default: "default", Values: []string{
			"default", "xs", "sm", "lg", "icon", "icon-xs", "icon-sm", "icon-lg"}},
	},
}

// button binds Button's shape to the helper calls button.gsx authors. The
// variable name is the component name: stylegen resolves button.Variant(v) by
// looking "button" up in Shapes().
var button = recipe.Component{Shape: buttonShape}
