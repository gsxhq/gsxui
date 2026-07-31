package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// NativeSelect is the styled native <select> (the shadcn native-select.tsx
// port), not the custom Radix-driven listbox that ships as Select. It has
// exactly two slots and no dimensions.
//
// The registry/CSS identity is the hyphenated "native-select"
// (data-gsxui-slot-native-select, gsxui-recipe-native-select); the Go
// receiver and type are derived camel — nativeSelect / nativeSelectRecipe —
// by internal/stylegen/identifier.go, the same split aspect-ratio and
// input-group already exercise.
//
// Slot order is document order: the positioned wrapper <div> encloses the
// <select>, which is the root slot. The wrapper — not the select — carries
// the caller's class (see registry/canonical/native-select.gsx's header for
// why the width intent has to land there), so its `w-fit` is a plain
// utility a caller's `w-full` displaces through merge.Merge, exactly as it
// displaced the old zero-specificity :where() rule.
//
// disabled, :focus-visible and aria-invalid are runtime state, not
// caller-chosen variants — same call as Input's and Textarea's shapes — so
// they fold onto the root slot's one rule as native Tailwind variants. The
// chevron rule keyed on the wrapper's `> svg` child is a relational form of
// the shape Select's trigger already proved ([&>svg]:), not a slot of its
// own: the chevron is icon.ChevronDown and carries no native-select marker.
var NativeSelect = recipe.Shape{
	Component: "native-select",
	Slots: []recipe.Slot{
		{Name: "wrapper", Base: true},
		{Name: "", Base: true},
	},
}
