package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Item is a layout composite: a media + content + actions row, with
// group/separator parts for stacked lists and header/footer parts for framing.
//
// item's data-size (default|sm) IS a shape dimension: the 8-style port made
// the ramp per-style (upstream .cn-item-size-default/-sm differ in vega/maia/
// luma/sera/rhea), so the old fold-into-the-base-rule approach (nova-only
// density retarget) can't express it. Upstream's xs value is deliberately not
// declared — gsxui ships no dropdown-embedded Item, xs's only consumer.
//
// One axis the style contract declares is deliberately NOT a shape dimension:
//
//   - item-separator's data-orientation. default.css only ever gives
//     ItemSeparator a my-2 margin, identical for both orientations.
//
// The "content" slot's base rule is deliberately PARTIAL: its flex-1, and the
// adjacent-sibling rule that overrides it with flex-none, both stay behind in
// @layer components under design doc § 10b. See
// assets/css/styles/default/item.css for why the pair has to move together.
var Item = recipe.Shape{
	Component: "item",
	Slots: []recipe.Slot{
		{Name: "", Base: true, Dimensions: []recipe.Dimension{
			{Name: "variant", Default: "default", Values: []string{"default", "outline", "muted"}},
			{Name: "size", Default: "default", Values: []string{"default", "sm"}},
		}},
		{Name: "group", Base: true},
		{Name: "separator", Base: true},
		{Name: "media", Base: true, Dimensions: []recipe.Dimension{
			{Name: "variant", Default: "default", Values: []string{"default", "icon", "image"}},
		}},
		{Name: "content", Base: true},
		{Name: "title", Base: true},
		{Name: "description", Base: true},
		{Name: "actions", Base: true},
		{Name: "header", Base: true},
		{Name: "footer", Base: true},
	},
}
