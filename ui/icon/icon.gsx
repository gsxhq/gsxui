package icon

import "github.com/gsxhq/gsx"

// svgIcon renders a Lucide icon's <svg> wrapper: 24x24 viewBox, the Lucide
// stroke defaults, default class="size-4" (caller-overridable via tailwind
// class merge), data-slot="icon", and aria-hidden="true" unless the caller
// already supplies an aria-hidden attribute — literal attributes authored
// before { attrs... } only render when attrs doesn't already set that exact
// key, so a caller's own aria-hidden (or, via the same fallthrough, other
// aria-* attributes alongside it) wins.
component svgIcon(name string, inner gsx.Node, attrs gsx.Attrs) {
	<svg
		data-slot="icon"
		aria-hidden="true"
		xmlns="http://www.w3.org/2000/svg"
		viewBox="0 0 24 24"
		fill="none"
		stroke="currentColor"
		stroke-width="2"
		stroke-linecap="round"
		stroke-linejoin="round"
		class="size-4"
		{ attrs... }
	>{ inner }</svg>
}
