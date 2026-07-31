package recipe

import (
	"fmt"
	"slices"
	"strings"
)

// Resolved is a style proven to be a complete implementation of a shape.
type Resolved struct {
	Shape  Shape
	Base   map[string][]string                       // [slot]
	Values map[string]map[string]map[string][]string // [slot][dimension][value]
}

// BaseUtilities returns the utilities a resolved recipe supplies for a slot's
// base rule.
func (r Resolved) BaseUtilities(slot string) []string { return slices.Clone(r.Base[slot]) }

// Utilities returns the utilities a resolved recipe supplies for one value.
func (r Resolved) Utilities(slot, dimension, value string) []string {
	return slices.Clone(r.Values[slot][dimension][value])
}

// Conform checks that style implements shape exactly: every declared
// (slot, dimension, value) has a rule, and every rule maps to a declared one.
func Conform(filename string, shape Shape, style Style) (Resolved, error) {
	if err := shape.Validate(); err != nil {
		return Resolved{}, err
	}

	resolved := Resolved{
		Shape:  shape,
		Base:   make(map[string][]string, len(shape.Slots)),
		Values: make(map[string]map[string]map[string][]string, len(shape.Slots)),
	}
	for _, slot := range shape.Slots {
		byDimension := make(map[string]map[string][]string, len(slot.Dimensions))
		for _, dimension := range slot.Dimensions {
			byDimension[dimension.Name] = make(map[string][]string, len(dimension.Values))
		}
		resolved.Values[slot.Name] = byDimension
	}

	// Style to shape: every rule must be declared.
	for _, rule := range style.Rules() {
		slot, dimension, value, kind, err := shape.DecodeClass(rule.Class)
		if err != nil {
			return Resolved{}, fmt.Errorf("%s: %w", filename, err)
		}
		if kind == ClassBase {
			declared, _ := shape.Slot(slot)
			if !declared.Base {
				return Resolved{}, fmt.Errorf("%s: slot %q declares no base rule, found %s",
					filename, slot, rule.Class)
			}
			resolved.Base[slot] = slices.Clone(rule.Utilities)
			continue
		}
		resolved.Values[slot][dimension][value] = slices.Clone(rule.Utilities)
	}

	// Shape to style: every declaration must be supplied.
	for _, slot := range shape.Slots {
		if slot.Base {
			if _, ok := resolved.Base[slot.Name]; !ok {
				return Resolved{}, fmt.Errorf("%s: slot %q missing base rule .%s",
					filename, slot.Name, shape.BaseClass(slot.Name))
			}
		}
		for _, dimension := range slot.Dimensions {
			for _, value := range dimension.Values {
				if _, ok := resolved.Values[slot.Name][dimension.Name][value]; !ok {
					return Resolved{}, fmt.Errorf("%s: slot %q dimension %q missing value %q",
						filename, slot.Name, dimension.Name, value)
				}
			}
		}
	}
	return resolved, nil
}

// CheckConflicts reports a utility list that contains a Tailwind conflict or repetition.
// merger is the Tailwind-aware class merger; a list that shortens when merged
// contained a utility that is either superseded by a conflicting one or repeated.
// This function uses multiset logic (counting occurrences) to distinguish between
// superseded utilities (output count=0, a conflict) and repeated utilities
// (output count>0 but less than input, a duplicate), so neither behavior is
// accidental or subject to simplification.
func CheckConflicts(filename string, resolved Resolved, merger func([]string) string) error {
	check := func(class string, utilities []string) error {
		kept := strings.Fields(merger(slices.Clone(utilities)))
		if len(kept) == len(utilities) {
			return nil
		}

		// Count occurrences in the original and merged lists.
		inputCounts := make(map[string]int, len(utilities))
		for _, utility := range utilities {
			inputCounts[utility]++
		}
		outputCounts := make(map[string]int, len(kept))
		for _, utility := range kept {
			outputCounts[utility]++
		}

		// Check for utilities that were removed or reduced.
		// Superseded (output count=0) takes priority over repeated (output count>0).
		seen := make(map[string]struct{})
		for _, utility := range utilities {
			if _, ok := seen[utility]; ok {
				continue // Already processed
			}
			seen[utility] = struct{}{}

			outCount := outputCounts[utility]
			if outCount == 0 {
				// Utility was completely removed: it was superseded by a conflicting one.
				return fmt.Errorf("%s: .%s applies conflicting utilities: %s is superseded",
					filename, class, utility)
			}
		}

		// Check for utilities that were repeated (if no superseded utilities found).
		// Iterate the original slice order, not the inputCounts map, so the
		// reported utility is deterministic when a rule repeats more than one.
		reported := make(map[string]struct{})
		for _, utility := range utilities {
			if _, ok := reported[utility]; ok {
				continue
			}
			reported[utility] = struct{}{}
			inCount := inputCounts[utility]
			outCount := outputCounts[utility]
			if outCount > 0 && outCount < inCount {
				// Utility appeared multiple times in input but fewer in output: it was repeated.
				return fmt.Errorf("%s: .%s repeats utility %s", filename, class, utility)
			}
		}

		return nil
	}

	for _, slot := range resolved.Shape.Slots {
		if slot.Base {
			if err := check(resolved.Shape.BaseClass(slot.Name), resolved.Base[slot.Name]); err != nil {
				return err
			}
		}
		for _, dimension := range slot.Dimensions {
			for _, value := range dimension.Values {
				class := resolved.Shape.ValueClass(slot.Name, dimension.Name, value)
				if err := check(class, resolved.Values[slot.Name][dimension.Name][value]); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
