// Package slider holds the site's example gsx components for ui/slider.
package slider

import "github.com/gsxhq/gsxui/ui"

// Basic mirrors shadcn's own slider-demo.tsx — the only slider example in
// the stated source directory (defaultValue={[50]} max={100} step={1},
// single thumb, w-[60%] wrapper). Range/Multiple-thumb/Vertical demos from
// the newer bases/radix examples tree are out of scope: this port's single
// <input type=range> ADAPT has no analog for them at all (GAP, see
// docs/jsx-parity.md `## slider`).
component Basic() {
	<ui.Slider value={50} min={0} max={100} step={1} class="w-[60%]" aria-label="Volume"/>
}
