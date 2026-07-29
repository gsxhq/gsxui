package stylecontract

// Toaster and Toast are declared separately because they are separate
// gsx components (ui/sonner.gsx exports both Toaster and Toast, and gsxui
// renders the toast DOM itself rather than delegating to the sonner npm
// package as upstream shadcn's sonner.tsx does). Upstream shadcn has since
// added a matching standalone "toast" component (apps/v4/registry/bases/base/ui/toast.tsx,
// PR #11266) using the same slot vocabulary rendered here. Splitting the
// contract this way lets each component's slots follow the
// "<RegistryName>-<relative>" naming rule with no exception, per
// registry/canonical/shapes/agreement_test.go.
//
// Both components still live in the single ui/sonner.gsx/.js/.x.go source
// file; this split is contract-only.
var toastContracts = []Component{
	{
		Name:         "Toaster",
		RegistryName: "toaster",
		Slots: []Slot{
			{Name: "toaster", Axes: []Axis{{Attribute: "data-expanded", Values: []string{"false", "true"}, RuntimeValues: []string{"false", "true"}}}},
		},
	},
	{
		Name:         "Toast",
		RegistryName: "toast",
		Slots: []Slot{
			{Name: "toast", Axes: []Axis{
				{Attribute: "data-type", Values: []string{"default", "success", "info", "warning", "error", "loading"}},
				{Attribute: "data-state", Values: []string{"closed", "open"}, RuntimeValues: []string{"closed", "open"}},
				{Attribute: "data-visible", Values: []string{"false", "true"}, RuntimeValues: []string{"false", "true"}},
			}},
			{Name: "toast-icon"},
			{Name: "toast-content"},
			{Name: "toast-title"},
			{Name: "toast-description"},
			{Name: "toast-action"},
			{Name: "toast-cancel"},
			{Name: "toast-close"},
			{Name: "toast-close-icon"},
		},
	},
}
