package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// InputOTP is a compound container: five styled slots, no dimensions. Its
// two "caret" markers (data-gsxui-slot-input-otp-caret-overlay,
// data-gsxui-slot-input-otp-caret) are NOT declared here — ui/input-otp.js
// creates both elements at runtime with no class attribute at all, and
// neither has any call site in registry/canonical/input-otp.gsx for a
// SlotClass accessor to resolve into; declaring them as slots here would
// fail checkShapeCoverage ("the component never renders it"), the same
// reasoning Accordion's root and Table's table-header/table-body already
// established for an unstyled/unreachable slot. Their rules stay behind in
// assets/css/styles/default.css's @layer utilities block as the escape
// hatch — not because they're untranslatable CSS, but because nothing in
// this migration's mechanism can carry a resolved class to a JS-created
// element. This is another HYPHENATED component ("input-otp" kebab
// identity, derived camel Go identifier `inputOtpRecipe`/`inputOtp` per
// componentIdentifier's segment-title-case-then-lower-initial rule — "otp"
// title-cases to "Otp", not an all-caps "OTP").
var InputOTP = recipe.Shape{
	Component: "input-otp",
	Slots: []recipe.Slot{
		{Name: "", Base: true},
		{Name: "input", Base: true},
		{Name: "group", Base: true},
		{Name: "slot", Base: true},
		{Name: "separator", Base: true},
	},
}
