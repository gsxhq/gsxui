package shapes

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/internal/stylecontract"
)

// The style contract is AUTHORITATIVE over the shape.
//
// stylecontract describes the runtime DOM surface: the slot markers that are
// actually stamped onto elements and the data-* attributes those elements
// actually carry. A shape only drives style generation — it can only style what
// the DOM already exposes. So a shape may declare LESS than the contract (the
// contract's `disabled` and `aria-invalid` axes are runtime state nobody varies
// a recipe dimension on), but never MORE and never something different. When
// this test fails, change the shape, unless the runtime markup genuinely
// changed — in which case change the contract first and the shape after.
const authorityNote = "the style contract is authoritative (it describes the rendered DOM); " +
	"fix the shape unless the runtime markup itself changed"

// contractSlotName maps a recipe slot name, which is RELATIVE to its component
// ("" is the root, then "header", "menu-button"), onto the full slot name used
// by the style contract ("card", "card-header").
//
// The mapping is registry name for the root, and registry name + "-" + relative
// name otherwise. That holds for every component in stylecontract.All() except
// Sonner, whose slots are named after "toaster"/"toast" rather than its registry
// name "sonner"; Sonner has no shape today, and if it ever gets one this check
// will fail loudly rather than silently accept a mismatch.
func contractSlotName(registryName, slot string) string {
	if slot == "" {
		return registryName
	}
	return registryName + "-" + slot
}

// axisAttribute maps a recipe dimension name onto the data-* attribute the
// contract declares for it.
func axisAttribute(dimension string) string {
	return "data-" + dimension
}

// agreementViolations reports every way shape disagrees with component. An
// empty result means the shape's whole vocabulary is covered by the contract.
func agreementViolations(shape recipe.Shape, component stylecontract.Component) []string {
	var problems []string
	report := func(format string, args ...any) {
		problems = append(problems, fmt.Sprintf(format, args...))
	}

	// 1. component name == RegistryName.
	if shape.Component != component.RegistryName {
		report("shape component %q != style-contract %s RegistryName %q",
			shape.Component, component.Name, component.RegistryName)
	}

	contractSlots := make(map[string]stylecontract.Slot, len(component.Slots))
	for _, slot := range component.Slots {
		contractSlots[slot.Name] = slot
	}

	for _, slot := range shape.Slots {
		full := contractSlotName(component.RegistryName, slot.Name)

		// 2. slot names map onto the contract's full slot names.
		contractSlot, ok := contractSlots[full]
		if !ok {
			report("shape slot %q maps to contract slot %q, which %s does not declare (contract declares %s)",
				slot.Name, full, component.Name, strings.Join(sortedSlotNames(component), ", "))
			continue
		}

		axes := make(map[string]stylecontract.Axis, len(contractSlot.Axes))
		for _, axis := range contractSlot.Axes {
			axes[axis.Attribute] = axis
		}

		for _, dimension := range slot.Dimensions {
			attribute := axisAttribute(dimension.Name)

			// 3. dimensions map onto declared data-* axes.
			axis, ok := axes[attribute]
			if !ok {
				report("shape slot %q dimension %q maps to axis %q, which contract slot %q does not declare",
					slot.Name, dimension.Name, attribute, full)
				continue
			}

			// 4. dimension values are permitted by that axis.
			permitted := make(map[string]struct{}, len(axis.Values))
			for _, value := range axis.Values {
				permitted[value] = struct{}{}
			}
			for _, value := range dimension.Values {
				if _, ok := permitted[value]; !ok {
					report("shape slot %q dimension %q declares value %q, which axis %q does not permit (permits %s)",
						slot.Name, dimension.Name, value, attribute, strings.Join(axis.Values, ", "))
				}
			}
		}
	}
	return problems
}

func sortedSlotNames(component stylecontract.Component) []string {
	names := make([]string, 0, len(component.Slots))
	for _, slot := range component.Slots {
		names = append(names, slot.Name)
	}
	sort.Strings(names)
	return names
}

// TestShapesAgreeWithStyleContracts is the executable bridge between the two
// hand-authored descriptions of a component's vocabulary. They stay
// independent; this is what stops them drifting apart.
func TestShapesAgreeWithStyleContracts(t *testing.T) {
	t.Parallel()

	byRegistryName := make(map[string]stylecontract.Component)
	for _, component := range stylecontract.All() {
		byRegistryName[component.RegistryName] = component
	}

	checked := 0
	for _, name := range sortedShapeNames() {
		component, ok := byRegistryName[name]
		if !ok {
			// No style contract for this shape at all: a shape whose component
			// name matches no RegistryName already violates equivalence 1.
			t.Errorf("shape %q has no style contract with that RegistryName — %s", name, authorityNote)
			continue
		}
		checked++
		for _, problem := range agreementViolations(All()[name], component) {
			t.Errorf("%s: %s — %s", name, problem, authorityNote)
		}
	}

	if checked == 0 {
		t.Fatal("no shape/style-contract pairs were checked — this test must never pass vacuously")
	}
}

func sortedShapeNames() []string {
	all := All()
	names := make([]string, 0, len(all))
	for name := range all {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// TestAgreementViolationsDetectEachDisagreement mutates a copy of Button's
// shape one equivalence at a time and asserts the check catches each. Without
// it, agreementViolations could silently stop detecting anything.
func TestAgreementViolationsDetectEachDisagreement(t *testing.T) {
	t.Parallel()

	var button stylecontract.Component
	for _, component := range stylecontract.All() {
		if component.RegistryName == "button" {
			button = component
		}
	}
	if button.RegistryName == "" {
		t.Fatal("no style contract with RegistryName button")
	}

	for _, testCase := range []struct {
		name    string
		mutate  func(*recipe.Shape)
		wantSub string
	}{
		{
			name:    "component name disagrees",
			mutate:  func(s *recipe.Shape) { s.Component = "buttonn" },
			wantSub: `shape component "buttonn"`,
		},
		{
			name:    "unknown slot",
			mutate:  func(s *recipe.Shape) { s.Slots[0].Name = "gizmo" },
			wantSub: `maps to contract slot "button-gizmo"`,
		},
		{
			name:    "dimension with no axis",
			mutate:  func(s *recipe.Shape) { s.Slots[0].Dimensions[0].Name = "tone" },
			wantSub: `maps to axis "data-tone"`,
		},
		{
			name: "value the axis does not permit",
			mutate: func(s *recipe.Shape) {
				s.Slots[0].Dimensions[0].Values = append(s.Slots[0].Dimensions[0].Values, "chartreuse")
			},
			wantSub: `declares value "chartreuse"`,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			mutated := cloneShape(Button)
			testCase.mutate(&mutated)
			problems := agreementViolations(mutated, button)
			if len(problems) == 0 {
				t.Fatalf("mutation %q produced no violation", testCase.name)
			}
			joined := strings.Join(problems, "\n")
			if !strings.Contains(joined, testCase.wantSub) {
				t.Errorf("violations do not mention %q:\n%s", testCase.wantSub, joined)
			}
		})
	}
}

func cloneShape(shape recipe.Shape) recipe.Shape {
	clone := recipe.Shape{Component: shape.Component}
	for _, slot := range shape.Slots {
		clonedSlot := recipe.Slot{Name: slot.Name, Base: slot.Base}
		for _, dimension := range slot.Dimensions {
			clonedDimension := recipe.Dimension{Name: dimension.Name, Default: dimension.Default}
			clonedDimension.Values = append(clonedDimension.Values, dimension.Values...)
			clonedSlot.Dimensions = append(clonedSlot.Dimensions, clonedDimension)
		}
		clone.Slots = append(clone.Slots, clonedSlot)
	}
	return clone
}
