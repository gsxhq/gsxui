package ui

import "github.com/gsxhq/gsx"

// Button is the shadcn/ui Button, retargeted to nova density (2026-07-24
// nova density map, `## button`). ADAPT: nova keys directional icon padding
// off `data-icon="inline-start|inline-end"` stamps we don't emit; gsxui kept
// its existing has-[>svg]:px-* selector mechanism and substituted nova's
// inline-start numeric value (e.g. default has-[>svg]:px-3 → px-2). All
// variants now carry a transparent 1px border in the base (box-size
// consistency across variants) — outline just recolors that border and no
// longer changes the box. A non-empty href on an enabled Button renders an
// <a> (gsx's answer to asChild-wrapping a link); disabled always renders a
// real disabled <button>. type="button" is an overridable default — pass
// type="submit" at the call site to submit forms.
component Button(variant string, size string, href string, disabled bool, children gsx.Node, attrs gsx.Attrs) {
	{ if href != "" && !disabled {
		<a
			data-variant={variant |> default("default")}
			data-size={size |> default("default")}
			href={href}
			{ withSlot("button", attrs)... }
		>
			{ children }
		</a>
	} else {
		<button
			data-variant={variant |> default("default")}
			data-size={size |> default("default")}
			type="button"
			disabled={disabled}
			{ withSlot("button", attrs)... }
		>
			{ children }
		</button>
	} }
}
