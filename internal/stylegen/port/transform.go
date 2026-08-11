package port

import (
	"fmt"
	"strings"

	"github.com/gsxhq/gsxui/internal/recipe"
)

// Ported is one component's recipe, ready to render.
type Ported struct {
	Component        string
	Base             map[string][]string                       // slot -> utilities
	Values           map[string]map[string]map[string][]string // slot -> dim -> value -> utilities
	Carried          map[string]bool                           // class -> true when carried from nova
	SrcStart, SrcEnd int
}

// Transform applies the three upstream→recipe rules, then the fallback
// policy, to one component's upstream MARK section:
//
//  1. Slot rename: an upstream `.cn-<component>[-<slot>]` rule's own
//     utilities land on the gsxui slot SlotFor resolves it to.
//  2. Dimension re-split: a leading `data-[<dim>=<value>]:` token whose
//     <dim> is a dimension the shape DECLARES on that slot moves to the
//     `(slot, dim, value)` rule; the prefix is stripped (Classify's Rest).
//     Any other token — including a `data-[<dim>=<value>]:` token whose
//     <dim> the shape does NOT declare on that slot — stays on the slot's
//     base rule verbatim (Classify's Raw), unmodified.
//  3. Descendant decomposition: a `**:data-[slot=<x>]:`/`*:data-[slot=<x>]:`
//     token's utility (Rest) moves to whichever gsxui slot SlotFor resolves
//     <x> to, if the shape declares that slot; otherwise the token is
//     unmapped.
//
// Then the fallback: any declared (slot) or (slot, dim, value) that ends up
// with no utilities from rules 1-3 inherits the corresponding rule's
// utilities from fallback (the current nova recipe), recorded in
// Ported.Carried. A declared rule with no upstream utilities AND no nova
// fallback is a hard error — the shape declares something nothing has ever
// styled.
//
// unmapped lists every upstream class or (class, utility) pair Transform
// found no gsxui destination for — reported, never silently dropped. A
// whole rule with no slot at all (SlotFor/shape.Slot both failed) appears
// as its bare class; a single token within an otherwise-mapped rule (Rule
// 3's target slot missing) appears as "<class>: <token>".
func Transform(shape recipe.Shape, sec Section, fallback recipe.Style) (Ported, []string, error) {
	ported := Ported{
		Component: shape.Component,
		Base:      make(map[string][]string),
		Values:    make(map[string]map[string]map[string][]string),
		Carried:   make(map[string]bool),
		SrcStart:  sec.Start,
		SrcEnd:    sec.End,
	}

	var unmapped []string

	for _, rule := range sec.Rules {
		if Ignored(rule.Class) {
			continue
		}

		slot, ok := SlotFor(shape.Component, rule.Class)
		if !ok {
			unmapped = append(unmapped, rule.Class)
			continue
		}
		slotDecl, hasSlot := shape.Slot(slot)
		if !hasSlot {
			unmapped = append(unmapped, rule.Class)
			continue
		}

		for _, token := range rule.Utilities {
			classified := Classify(token)
			switch classified.Kind {
			case KindSlot:
				targetSlot, ok := SlotFor(shape.Component, "cn-"+classified.Slot)
				if !ok {
					unmapped = append(unmapped, fmt.Sprintf("%s: %s", rule.Class, token))
					continue
				}
				if _, hasTarget := shape.Slot(targetSlot); !hasTarget {
					unmapped = append(unmapped, fmt.Sprintf("%s: %s", rule.Class, token))
					continue
				}
				ported.appendBase(targetSlot, classified.Rest)

			case KindDimension:
				if _, hasDimension := slotDecl.Dimension(classified.Dimension); hasDimension {
					ported.appendValue(slot, classified.Dimension, classified.Value, classified.Rest)
				} else {
					ported.appendBase(slot, classified.Raw)
				}

			default: // KindPlain
				ported.appendBase(slot, classified.Raw)
			}
		}
	}

	var missing []string
	for _, slot := range shape.Slots {
		if slot.Base && len(ported.Base[slot.Name]) == 0 {
			class := shape.BaseClass(slot.Name)
			if rule, ok := fallback.Lookup(class); ok {
				ported.Base[slot.Name] = rule.Utilities
				ported.Carried[class] = true
			} else {
				missing = append(missing, fmt.Sprintf("slot %s: no upstream rule and no nova fallback (class %s)", slotLabel(slot.Name), class))
			}
		}
		for _, dimension := range slot.Dimensions {
			for _, value := range dimension.Values {
				if len(ported.Values[slot.Name][dimension.Name][value]) > 0 {
					continue
				}
				class := shape.ValueClass(slot.Name, dimension.Name, value)
				if rule, ok := fallback.Lookup(class); ok {
					ported.setValue(slot.Name, dimension.Name, value, rule.Utilities)
					ported.Carried[class] = true
				} else {
					missing = append(missing, fmt.Sprintf("slot %s dimension %q value %q: no upstream rule and no nova fallback (class %s)", slotLabel(slot.Name), dimension.Name, value, class))
				}
			}
		}
	}
	if len(missing) > 0 {
		return Ported{}, unmapped, fmt.Errorf("%s: %s", shape.Component, strings.Join(missing, "; "))
	}

	return ported, unmapped, nil
}

func (p *Ported) appendBase(slot, utility string) {
	p.Base[slot] = append(p.Base[slot], utility)
}

func (p *Ported) appendValue(slot, dimension, value, utility string) {
	p.ensureValues(slot, dimension)
	p.Values[slot][dimension][value] = append(p.Values[slot][dimension][value], utility)
}

func (p *Ported) setValue(slot, dimension, value string, utilities []string) {
	p.ensureValues(slot, dimension)
	p.Values[slot][dimension][value] = utilities
}

func (p *Ported) ensureValues(slot, dimension string) {
	byDimension, ok := p.Values[slot]
	if !ok {
		byDimension = make(map[string]map[string][]string)
		p.Values[slot] = byDimension
	}
	if _, ok := byDimension[dimension]; !ok {
		byDimension[dimension] = make(map[string][]string)
	}
}

// slotLabel renders a slot name for an error message, naming the root slot
// explicitly rather than leaving it a confusing empty string.
func slotLabel(name string) string {
	if name == "" {
		return "<root>"
	}
	return fmt.Sprintf("%q", name)
}
