package canonical

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
				"group/button",
				button.Role(),
				button.Variant(variant),
				button.Size(size),
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
				"group/button",
				button.Role(),
				button.Variant(variant),
				button.Size(size),
			}
			{ attrs... }
			data-gsxui-slot-button
		>
			{ children }
		</button>
	} }
}
