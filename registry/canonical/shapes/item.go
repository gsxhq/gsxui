package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Item is a layout composite: a media + content + actions row, with
// group/separator parts for stacked lists and header/footer parts for framing.
//
// Two axes the style contract declares are deliberately NOT shape dimensions:
//
//   - item's data-size (default|sm). gsxui's nova density retarget folded the
//     whole size ramp into the base rule (gap-2.5 px-3 py-2.5), so
//     assets/css/styles/default/item.css has no size-keyed rule at all and a
//     size dimension would have no content for either value. Conform requires
//     every declared value to author at least one utility; a shape may declare
//     LESS than the contract (see agreement_test.go), so it declares nothing
//     here rather than inventing a ramp. Upstream shadcn does ship one
//     (registry/new-york-v4/ui/item.tsx's itemVariants.size) — restoring it is
//     a density decision, not a migration one.
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
