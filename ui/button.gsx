package ui

import "github.com/gsxhq/gsx"

// Button is the canonical Button structure. Its semantic recipe roles resolve
// to concrete style source before consumers compile the generated component.
// A non-empty href on an enabled Button renders an <a> (gsx's answer to
// asChild-wrapping a link); disabled always renders a real disabled <button>.
// type="button" is an overridable default — pass type="submit" at the call
// site to submit forms.
component Button(variant string, size string, href string, disabled bool, children gsx.Node, attrs gsx.Attrs) {
	{ if href != "" && !disabled {
		<a
			data-variant={variant |> default("default")}
			data-size={size |> default("default")}
			href={href}
			class={
				"gsxui-recipe-button",
				switch variant {
				case "", "default":
					"gsxui-recipe-button-variant-default"
				case "destructive":
					"gsxui-recipe-button-variant-destructive"
				case "outline":
					"gsxui-recipe-button-variant-outline"
				case "secondary":
					"gsxui-recipe-button-variant-secondary"
				case "ghost":
					"gsxui-recipe-button-variant-ghost"
				case "link":
					"gsxui-recipe-button-variant-link"
				default:
					""
				},
				switch size {
				case "", "default":
					"gsxui-recipe-button-size-default"
				case "xs":
					"gsxui-recipe-button-size-xs"
				case "sm":
					"gsxui-recipe-button-size-sm"
				case "lg":
					"gsxui-recipe-button-size-lg"
				case "icon":
					"gsxui-recipe-button-size-icon"
				case "icon-xs":
					"gsxui-recipe-button-size-icon-xs"
				case "icon-sm":
					"gsxui-recipe-button-size-icon-sm"
				case "icon-lg":
					"gsxui-recipe-button-size-icon-lg"
				default:
					""
				},
			}
			{ attrs... }
			data-gsxui-slot-button
		>
			{ children }
		</a>
	} else {
		<button
			data-variant={variant |> default("default")}
			data-size={size |> default("default")}
			type="button"
			disabled={disabled}
			class={
				"gsxui-recipe-button",
				switch variant {
				case "", "default":
					"gsxui-recipe-button-variant-default"
				case "destructive":
					"gsxui-recipe-button-variant-destructive"
				case "outline":
					"gsxui-recipe-button-variant-outline"
				case "secondary":
					"gsxui-recipe-button-variant-secondary"
				case "ghost":
					"gsxui-recipe-button-variant-ghost"
				case "link":
					"gsxui-recipe-button-variant-link"
				default:
					""
				},
				switch size {
				case "", "default":
					"gsxui-recipe-button-size-default"
				case "xs":
					"gsxui-recipe-button-size-xs"
				case "sm":
					"gsxui-recipe-button-size-sm"
				case "lg":
					"gsxui-recipe-button-size-lg"
				case "icon":
					"gsxui-recipe-button-size-icon"
				case "icon-xs":
					"gsxui-recipe-button-size-icon-xs"
				case "icon-sm":
					"gsxui-recipe-button-size-icon-sm"
				case "icon-lg":
					"gsxui-recipe-button-size-icon-lg"
				default:
					""
				},
			}
			{ attrs... }
			data-gsxui-slot-button
		>
			{ children }
		</button>
	} }
}
