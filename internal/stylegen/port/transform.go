package port

import (
	"fmt"
	"slices"
	"strings"

	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/merge"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// Ported is one component's recipe, ready to render.
type Ported struct {
	Component        string
	Base             map[string][]string                       // slot -> utilities
	Values           map[string]map[string]map[string][]string // slot -> dim -> value -> utilities
	Carried          map[string]bool                           // class -> true when carried from nova
	SrcStart, SrcEnd int
}

// Transform applies the upstream→recipe rules, then the fallback policy, to
// one component's upstream MARK section:
//
//  0. Whole-rule dimension decode: upstream sometimes spells a dimension
//     value as its OWN top-level rule — `.cn-button-variant-destructive` —
//     rather than as a `data-[dim=value]:` prefix inside a shared rule
//     (confirmed against button.css, the one hand-verified real port: its
//     .cn-button-variant-*/.cn-button-size-* class names are isomorphic to
//     gsxui's own .gsxui-recipe-button-variant-*/-size-* grammar). When
//     SlotFor's default derivation names a "slot" the shape doesn't declare,
//     decodeUpstreamValue retries it as a whole dimension-value class before
//     giving up.
//  1. Slot rename: an upstream `.cn-<component>[-<slot>]` rule's own
//     utilities land on the gsxui slot (or slot+dimension+value, per 0
//     above) SlotFor/decodeUpstreamValue resolves it to.
//  2. Dimension re-split: within a plain (non-value) rule, a leading
//     `data-[<dim>=<value>]:` token whose <dim> is a dimension the shape
//     DECLARES on that slot moves to the `(slot, dim, value)` rule; the
//     prefix is stripped (Classify's Rest). Any other token — including a
//     `data-[<dim>=<value>]:` token whose <dim> the shape does NOT declare
//     on that slot — stays on the rule's own destination verbatim
//     (Classify's Raw), unmodified.
//  3. Descendant decomposition: a `**:data-[slot=<x>]:`/`*:data-[slot=<x>]:`
//     token names a gsxui slot <x> a DIFFERENT element renders. When the
//     source rule is itself unconditional (rule 0/2 above did not put it
//     behind a dimension value) AND <x> names a slot other than the source's
//     own, the utility (Rest) moves to <x>'s own base rule, unconditionally
//     — <x> always looks this way wherever it's nested (accordion's
//     trigger-icon). Otherwise — the source rule is ITSELF a dimension value
//     (a variant/size-scoped rule reaching into a child, e.g. Alert's
//     destructive-variant painting its description red) or <x> names the
//     SAME slot the rule started on (a component nesting inside itself,
//     e.g. Field's nested field-group) — the styling is conditional on that
//     context, not universal, so it cannot become <x>'s own unconditional
//     rule. It is instead rewritten as an arbitrary descendant-selector
//     utility (`[&_[data-gsxui-slot-<x>]]:`/`[&>[data-gsxui-slot-<x>]]:` —
//     Child decides `_` vs `>`) kept on the SOURCE rule's own destination,
//     matching gsxui's own existing convention for exactly this shape
//     (registry/styles/nova/alert.css's variant-destructive rule,
//     registry/styles/nova/field.css's nested field-group rule).
//
// Then the fallback: any declared (slot) or (slot, dim, value) that ends up
// with no utilities from the rules above inherits the corresponding rule's
// utilities from fallback (the current nova recipe), recorded in
// Ported.Carried. A declared rule with no upstream utilities AND no nova
// fallback is a hard error — the shape declares something nothing has ever
// styled.
//
// unmapped lists every upstream class or (class, utility) pair Transform
// found no gsxui destination for — reported, never silently dropped, unless
// the mapping table itself declared the non-mapping on purpose
// (silentNonMapping). A whole rule with no slot at all appears as its bare
// class; a single token within an otherwise-mapped rule (rule 3's target
// slot missing) appears as "<class>: <token>".
func Transform(shape recipe.Shape, sec Section, fallback recipe.Style) (Ported, []string, error) {
	ported := Ported{
		Component: shape.Component,
		Base:      make(map[string][]string),
		Values:    make(map[string]map[string]map[string][]string),
		Carried:   make(map[string]bool),
		SrcStart:  sec.Start,
		SrcEnd:    sec.End,
	}

	// Every declared component shape, for recognizing a KindSlot token that
	// targets an entirely different gsxui component's root (see the KindSlot
	// case below) rather than a sub-slot of this one.
	allShapes := shapes.All()

	var unmapped []string

	for _, rule := range sec.Rules {
		if Ignored(rule.Class) {
			continue
		}
		// renameClass substitutes a differently-spelled equivalent class
		// BEFORE any resolution, for the rare upstream/gsxui naming mismatch
		// that isn't a simple slot rename (see classRenames): rule.Class
		// itself is kept for unmapped/error reporting throughout, so a
		// renamed class's messages still name what upstream actually wrote.
		effectiveClass := renameClass(shape.Component, rule.Class)

		slot, ok := SlotFor(shape.Component, effectiveClass)
		var dimension, value string
		if ok {
			if _, hasSlot := shape.Slot(slot); !hasSlot {
				if s, d, v, decoded := decodeUpstreamValue(shape, slot); decoded {
					slot, dimension, value = s, d, v
				} else {
					ok = false
				}
			}
		}
		if !ok {
			if !silentNonMapping(shape.Component, effectiveClass) {
				unmapped = append(unmapped, rule.Class)
			}
			continue
		}
		slotDecl, _ := shape.Slot(slot) // guaranteed to exist: either matched directly, or decodeUpstreamValue confirmed it

		put := func(utility string) {
			if dimension == "" {
				ported.appendBase(slot, utility)
			} else {
				ported.appendValue(slot, dimension, value, utility)
			}
		}

		variant := overrideVariant(shape.Component, effectiveClass)
		for _, token := range rule.Utilities {
			token = stripImportantModifier(token)
			if token == noScrollbarUtility {
				for _, u := range noScrollbarExpansion {
					put(u)
				}
				continue
			}
			token = rewriteMarkerVariant(shape.Component, token)
			if slotAttributeDropped(shape.Component, token) {
				continue
			}
			rewritten, resolvedSlots := rewriteSlotAttributeReferences(allShapes, token)
			if !resolvedSlots {
				unmapped = append(unmapped, fmt.Sprintf("%s: %s", rule.Class, token))
				continue
			}
			token = rewritten
			if variant != "" {
				token = insertVariant(token, variant)
			}
			classified := Classify(token)
			switch classified.Kind {
			case KindSlot:
				targetSlot, ok := SlotFor(shape.Component, "cn-"+classified.Slot)
				if ok {
					if _, hasTarget := shape.Slot(targetSlot); !hasTarget {
						ok = false
					}
				}
				if ok {
					if dimension == "" && targetSlot != slot {
						ported.appendBase(targetSlot, classified.Rest)
					} else {
						put(descendantSelector(shape.Component, targetSlot, classified.Child) + ":" + classified.Rest)
					}
					continue
				}
				// classified.Slot may not be a sub-slot of THIS component at
				// all: it may instead be the ROOT of an entirely DIFFERENT
				// gsxui component nested inside this one (Tooltip reaching
				// into a rendered Kbd, Combobox reaching into a rendered
				// InputGroup). That relationship is always conditional on
				// this specific composition — it must never become the
				// OTHER component's own unconditional rule (every Kbd
				// everywhere would get Tooltip's treatment) — so it always
				// rewrites as a descendant-selector utility kept on the
				// SOURCE'S own rule, using the OTHER component's bare name
				// as the marker (no sub-slot: it's that component's root).
				if _, isOtherComponent := allShapes[classified.Slot]; isOtherComponent {
					put(descendantSelector(classified.Slot, "", classified.Child) + ":" + classified.Rest)
					continue
				}
				unmapped = append(unmapped, fmt.Sprintf("%s: %s", rule.Class, token))

			case KindDimension:
				if dimension == "" {
					if _, hasDimension := slotDecl.Dimension(classified.Dimension); hasDimension {
						ported.appendValue(slot, classified.Dimension, classified.Value, classified.Rest)
						continue
					}
				}
				put(classified.Raw)

			default: // KindPlain
				put(gateNativeVisibility(shape.Component, slot, dimension, classified.Raw))
			}
		}

		// A declared override may carry extra, verbatim utilities to append
		// after upstream's own contribution — content that is genuinely
		// style-invariant but lives outside style-<name>.css entirely (see
		// mapping.go's "button" slotOverrides entry).
		for _, extra := range overrideExtra(shape.Component, effectiveClass) {
			put(extra)
		}
	}

	var missing []string
	for _, slot := range shape.Slots {
		if slot.Base {
			class := shape.BaseClass(slot.Name)
			if len(ported.Base[slot.Name]) == 0 {
				if rule, ok := fallback.Lookup(class); ok {
					ported.Base[slot.Name] = rule.Utilities
					ported.Carried[class] = true
				} else {
					missing = append(missing, fmt.Sprintf("slot %s: no upstream rule and no nova fallback (class %s)", slotLabel(slot.Name), class))
				}
			} else {
				ported.carryMissingDisplayUtility(slot.Name, class, fallback)
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

	if err := ported.resolveMergeOrder(); err != nil {
		return Ported{}, unmapped, fmt.Errorf("%s: %w", shape.Component, err)
	}

	return ported, unmapped, nil
}

// resolveMergeOrder cleans up each assembled rule's utilities, when needed,
// so that merge.Merge — the same Tailwind-aware oracle recipe.CheckConflicts
// validates every recipe against — reports no conflict. Two independent
// contributors can land the identical utility string on the same slot (the
// dossier's own predicted case: Rule 1's slot-rename and Rule 3's descendant
// decomposition both reaching the same target, e.g. Select's trigger
// contributing `*:data-[slot=select-value]:flex` on top of select-value's
// own `flex`) — dedupeStable removes the exact repeat, which changes nothing
// about the compiled CSS (a repeated utility string has no additional
// effect). What's left, if still flagged, is a genuine same-rule ORDER
// conflict: Rule 1 preserves upstream's own @apply token order verbatim,
// and Tailwind utility conflicts are order-sensitive in exactly the way
// `@apply` itself is (a later declaration in the same rule overrides an
// earlier one for the same CSS property) — upstream sometimes lists a
// `leading-*` utility BEFORE a `text-<size>` utility that bundles its own
// default line-height, which silently drops the leading value once compiled
// (confirmed empirically: merge.Merge(["leading-snug","text-sm"]) drops
// leading-snug, merge.Merge(["text-sm","leading-snug"]) keeps both — and
// independently confirmed by registry/styles/nova/field.css's own
// field-title rule, whose human author already reordered these two
// utilities for exactly this reason). resolveConflictOrder does not
// hand-encode which Tailwind class groups conflict (that knowledge belongs
// to merge.Merge alone): it treats merge.Merge purely as an oracle and
// tries moving each utility, in turn, to the end of the list, keeping the
// first move that merges cleanly. A rule neither step can fix is a genuine
// unresolvable problem and is left for Validate to report loudly, never
// silently corrected here.
func (p *Ported) resolveMergeOrder() error {
	for slot, utilities := range p.Base {
		resolved, ok := resolveConflictOrder(dedupeStable(utilities))
		if !ok {
			return fmt.Errorf("slot %s: no utility reordering satisfies the Tailwind merger (likely a genuine conflict): %s",
				slotLabel(slot), strings.Join(utilities, " "))
		}
		p.Base[slot] = resolved
	}
	for slot, byDimension := range p.Values {
		for dimension, byValue := range byDimension {
			for value, utilities := range byValue {
				resolved, ok := resolveConflictOrder(dedupeStable(utilities))
				if !ok {
					return fmt.Errorf("slot %s dimension %q value %q: no utility reordering satisfies the Tailwind merger (likely a genuine conflict): %s",
						slotLabel(slot), dimension, value, strings.Join(utilities, " "))
				}
				byValue[value] = resolved
			}
		}
	}
	return nil
}

// dedupeStable drops exact-string repeats, keeping each utility's first
// occurrence and original relative order.
func dedupeStable(utilities []string) []string {
	seen := make(map[string]bool, len(utilities))
	out := make([]string, 0, len(utilities))
	for _, u := range utilities {
		if seen[u] {
			continue
		}
		seen[u] = true
		out = append(out, u)
	}
	return out
}

// resolveConflictOrder returns utilities unchanged if merge.Merge already
// reports no conflict; otherwise it tries moving each utility to the end of
// the list, one at a time, in order, and returns the first arrangement that
// merges cleanly. ok is false when no single move resolves it.
func resolveConflictOrder(utilities []string) (resolved []string, ok bool) {
	if mergesCleanly(utilities) {
		return utilities, true
	}
	for i := range utilities {
		candidate := make([]string, 0, len(utilities))
		candidate = append(candidate, utilities[:i]...)
		candidate = append(candidate, utilities[i+1:]...)
		candidate = append(candidate, utilities[i])
		if mergesCleanly(candidate) {
			return candidate, true
		}
	}
	return utilities, false
}

func mergesCleanly(utilities []string) bool {
	return len(strings.Fields(merge.Merge(slices.Clone(utilities)))) == len(utilities)
}

// decodeUpstreamValue retries a "slot" name SlotFor's default derivation
// produced (a plain prefix strip) that turned out not to be a real declared
// slot, as a whole dimension-value class instead. Upstream's own
// `.cn-<component>-<dim>-<value>` rules use exactly gsxui's own
// `<slot>-<dim>-<value>` hyphenation (button.css's .cn-button-variant-*
// confirmed isomorphic to .gsxui-recipe-button-variant-*), so this builds
// the equivalent synthetic gsxui class and reuses shape.DecodeClass itself
// — the longest-slot-name-first search DecodeClass already implements —
// rather than reimplementing that ambiguity resolution here.
//
// A second, looser upstream form also exists: some single-dimension axes are
// spelled WITHOUT the dimension's own name at all — just the bare value
// (`.cn-separator-horizontal`, not `.cn-separator-orientation-horizontal`),
// confirmed against registry/styles/nova/separator.css's own header
// ("folded in as an orientation dimension... matching the style contract
// and default.css exactly, no discrepancy to resolve"). When the
// dimension-named form doesn't decode, this also searches every declared
// slot+dimension for one whose Values contains bogusSlot's own remainder
// verbatim, accepting only an UNAMBIGUOUS single match — it never guesses
// between two dimensions that happen to share a value spelling.
func decodeUpstreamValue(shape recipe.Shape, bogusSlot string) (slot, dimension, value string, ok bool) {
	synthetic := shape.BaseClass("") + "-" + bogusSlot
	if s, d, v, kind, err := shape.DecodeClass(synthetic); err == nil && kind == recipe.ClassValue {
		return s, d, v, true
	}

	type candidate struct{ slot, dimension, value string }
	var candidates []candidate
	ordered := slices.Clone(shape.Slots)
	slices.SortFunc(ordered, func(a, b recipe.Slot) int { return len(b.Name) - len(a.Name) })
	for _, slotDecl := range ordered {
		remainder := bogusSlot
		if slotDecl.Name != "" {
			trimmed, cut := strings.CutPrefix(bogusSlot, slotDecl.Name+"-")
			if !cut {
				continue
			}
			remainder = trimmed
		}
		for _, dim := range slotDecl.Dimensions {
			if dim.Has(remainder) {
				candidates = append(candidates, candidate{slotDecl.Name, dim.Name, remainder})
			}
		}
	}
	if len(candidates) == 1 {
		return candidates[0].slot, candidates[0].dimension, candidates[0].value, true
	}
	return "", "", "", false
}

// carryMissingDisplayUtility appends a display-establishing utility (flex,
// grid, ...) from the fallback (nova) recipe when the transform DID produce
// some utilities for slot from upstream but is missing one nova proves
// necessary. This is deliberately narrower than the whole-rule fallback
// above it: running the real port and rendering the result found upstream's
// style-<name>.css genuinely never carries a huge swath of the catalog's
// own display-establishing utility at all — confirmed directly against
// upstream source (.cn-command, .cn-badge, .cn-avatar, .cn-item, and
// dozens more never mention flex/inline-flex/grid anywhere in ANY of
// their rules). It is baked into the shared React component's own cva()
// base class list instead (apps/v4/registry/bases/radix/ui/<component>.tsx
// or the plain-element equivalent), a layer the porter has no access to.
// Every one of gsxui's own hand-authored nova recipes already carries the
// correct structural utility — a human composed it against the real
// rendered DOM — and it is style-invariant BY CONSTRUCTION: the same
// display mechanism a component needs in every style, only the surrounding
// spacing/color/radius varies. So this carries just that one utility from
// nova, rather than the whole rule, keeping every upstream-sourced value
// intact while fixing the missing layout mechanism. Confirmed empirically
// against the real rendered page: without it, Field's whole compound
// layout collapses to block flow and FieldSeparator's negative margin
// overlaps the row after it; every other slot this function touches loses
// its flex arrangement (icon+label rows collapse, gap stops applying,
// wrapping breaks) the same way.
func (p *Ported) carryMissingDisplayUtility(slot, class string, fallback recipe.Style) {
	rule, ok := fallback.Lookup(class)
	if !ok {
		return
	}
	have := make(map[string]bool, len(p.Base[slot]))
	for _, u := range p.Base[slot] {
		have[u] = true
	}

	missingDisplay := false
	for _, u := range rule.Utilities {
		if displayEstablishingUtilities[u] && !have[u] {
			missingDisplay = true
			break
		}
	}
	if !missingDisplay {
		return
	}

	// A display utility is missing, so its container companions
	// (flex-direction/flex-wrap/grid-template — see
	// layoutCompanionUtilities's own doc comment) are carried alongside it:
	// upstream never mentions either family for this slot, and "flex"
	// without nova's own "flex-col" is worse than neither, not better.
	for _, u := range rule.Utilities {
		if have[u] {
			continue
		}
		if !displayEstablishingUtilities[u] && !isLayoutCompanionUtility(u) {
			continue
		}
		p.Base[slot] = append(p.Base[slot], u)
		p.Carried[class] = true
	}
}

// descendantSelector builds the arbitrary descendant-selector variant gsxui
// already uses to reach into a specific slot's element from an ancestor's
// rule: `[&_[data-gsxui-slot-<component>[-<slot>]]]:` for any descendant
// (child=false, upstream's `**:`) or `[&>[...]]:` for a direct child
// (child=true, upstream's `*:`).
func descendantSelector(component, targetSlot string, child bool) string {
	marker := "data-gsxui-slot-" + component
	if targetSlot != "" {
		marker += "-" + targetSlot
	}
	if child {
		return "[&>[" + marker + "]]"
	}
	return "[&_[" + marker + "]]"
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
