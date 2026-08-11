package port

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/gsxhq/gsxui/internal/recipe"
)

func mustParseStyle(t *testing.T, src string) recipe.Style {
	t.Helper()
	style, err := recipe.ParseStyle("fallback.css", []byte(src))
	if err != nil {
		t.Fatalf("ParseStyle:\n%s\n--- error ---\n%v", src, err)
	}
	return style
}

// Rule 1: slot rename. cn-accordion-item's own utilities land on the "item"
// slot's base rule verbatim — nothing in this token needs re-splitting.
func TestTransformRule1SlotRename(t *testing.T) {
	shape := recipe.Shape{
		Component: "accordion",
		Slots: []recipe.Slot{
			{Name: "item", Base: true},
		},
	}
	sec := Section{
		Name: "Accordion",
		Rules: []Rule{
			{Class: "cn-accordion-item", Utilities: []string{"not-last:border-b"}},
		},
	}

	ported, unmapped, err := Transform(shape, sec, recipe.Style{})
	if err != nil {
		t.Fatalf("Transform: %v", err)
	}
	if len(unmapped) != 0 {
		t.Fatalf("unmapped = %v, want none", unmapped)
	}
	if want := []string{"not-last:border-b"}; !reflect.DeepEqual(ported.Base["item"], want) {
		t.Errorf("Base[item] = %v, want %v", ported.Base["item"], want)
	}
}

// Rule 2: dimension re-split. cn-radio-group-thumb exercises radio's
// upstream-prefix quirk (SlotFor("radio", "cn-radio-group-thumb") ==
// "thumb" — "thumb" is a hypothetical slot name that isn't in mapping.go's
// radio-specific override table, so this still hits the plain prefix-strip
// path, not an override) as well as the split itself: the bare "size-4"
// utility stays on the base rule, and "data-[size=sm]:size-3" — size IS a
// dimension this slot declares — moves to the value rule with its prefix
// stripped. The dimension's other declared value ("default") gets nothing
// from upstream in this fixture, so the fallback supplies it; that is
// incidental to this test (the Fallback test below covers it directly), not
// what's asserted.
func TestTransformRule2DimensionResplit(t *testing.T) {
	shape := recipe.Shape{
		Component: "radio",
		Slots: []recipe.Slot{
			{
				Name: "thumb", Base: true,
				Dimensions: []recipe.Dimension{
					{Name: "size", Default: "default", Values: []string{"default", "sm"}},
				},
			},
		},
	}
	sec := Section{
		Name: "Radio Group",
		Rules: []Rule{
			{Class: "cn-radio-group-thumb", Utilities: []string{"size-4", "data-[size=sm]:size-3"}},
		},
	}
	fallback := mustParseStyle(t, fmt.Sprintf(`@layer components {
  .%s {
    @apply size-4;
  }
}
`, shape.ValueClass("thumb", "size", "default")))

	ported, unmapped, err := Transform(shape, sec, fallback)
	if err != nil {
		t.Fatalf("Transform: %v", err)
	}
	if len(unmapped) != 0 {
		t.Fatalf("unmapped = %v, want none", unmapped)
	}
	if want := []string{"size-4"}; !reflect.DeepEqual(ported.Base["thumb"], want) {
		t.Errorf("Base[thumb] = %v, want %v", ported.Base["thumb"], want)
	}
	if want := []string{"size-3"}; !reflect.DeepEqual(ported.Values["thumb"]["size"]["sm"], want) {
		t.Errorf(`Values["thumb"]["size"]["sm"] = %v, want %v`, ported.Values["thumb"]["size"]["sm"], want)
	}
}

// Rule 2 negative: "side" is not a dimension the "content" slot declares
// (only Rule 2's own dimension check decides this, per slot — not any
// global list of known dimension names), so the whole token stays verbatim
// on the base rule, unstripped, and no error is raised.
func TestTransformRule2NegativeUndeclaredDimension(t *testing.T) {
	shape := recipe.Shape{
		Component: "select",
		Slots: []recipe.Slot{
			{Name: "content", Base: true},
		},
	}
	sec := Section{
		Name: "Select",
		Rules: []Rule{
			{Class: "cn-select-content", Utilities: []string{"data-[side=bottom]:slide-in-from-top-2"}},
		},
	}

	ported, unmapped, err := Transform(shape, sec, recipe.Style{})
	if err != nil {
		t.Fatalf("Transform: %v", err)
	}
	if len(unmapped) != 0 {
		t.Fatalf("unmapped = %v, want none", unmapped)
	}
	if want := []string{"data-[side=bottom]:slide-in-from-top-2"}; !reflect.DeepEqual(ported.Base["content"], want) {
		t.Errorf("Base[content] = %v, want %v", ported.Base["content"], want)
	}
	if len(ported.Values["content"]) != 0 {
		t.Errorf("Values[content] = %v, want empty (side is not a declared dimension)", ported.Values["content"])
	}
}

// Rule 3: descendant decomposition. cn-accordion-trigger's own utility
// (p-4) stays on "trigger"; the **:data-[slot=accordion-trigger-icon]:
// token's utility (size-4) moves to the "trigger-icon" slot's own base
// rule, its prefix stripped.
func TestTransformRule3DescendantDecomposition(t *testing.T) {
	shape := recipe.Shape{
		Component: "accordion",
		Slots: []recipe.Slot{
			{Name: "trigger", Base: true},
			{Name: "trigger-icon", Base: true},
		},
	}
	sec := Section{
		Name: "Accordion",
		Rules: []Rule{
			{Class: "cn-accordion-trigger", Utilities: []string{"**:data-[slot=accordion-trigger-icon]:size-4", "p-4"}},
		},
	}

	ported, unmapped, err := Transform(shape, sec, recipe.Style{})
	if err != nil {
		t.Fatalf("Transform: %v", err)
	}
	if len(unmapped) != 0 {
		t.Fatalf("unmapped = %v, want none", unmapped)
	}
	if want := []string{"p-4"}; !reflect.DeepEqual(ported.Base["trigger"], want) {
		t.Errorf("Base[trigger] = %v, want %v", ported.Base["trigger"], want)
	}
	if want := []string{"size-4"}; !reflect.DeepEqual(ported.Base["trigger-icon"], want) {
		t.Errorf("Base[trigger-icon] = %v, want %v", ported.Base["trigger-icon"], want)
	}
}

// Rule 3, unknown slot: the shape declares no "nonesuch" (nor any slot
// SlotFor would resolve "cn-nonesuch" to), so the token is reported in
// unmapped and lands in no rule at all — it must not silently create a
// "nonesuch" entry anywhere, and it must not suppress the rule's other,
// resolvable utility (p-4).
func TestTransformRule3UnknownSlot(t *testing.T) {
	shape := recipe.Shape{
		Component: "accordion",
		Slots: []recipe.Slot{
			{Name: "trigger", Base: true},
		},
	}
	sec := Section{
		Name: "Accordion",
		Rules: []Rule{
			{Class: "cn-accordion-trigger", Utilities: []string{"**:data-[slot=nonesuch]:x", "p-4"}},
		},
	}

	ported, unmapped, err := Transform(shape, sec, recipe.Style{})
	if err != nil {
		t.Fatalf("Transform: %v", err)
	}
	if want := []string{"p-4"}; !reflect.DeepEqual(ported.Base["trigger"], want) {
		t.Errorf("Base[trigger] = %v, want %v", ported.Base["trigger"], want)
	}
	if want := []string{"cn-accordion-trigger: **:data-[slot=nonesuch]:x"}; !reflect.DeepEqual(unmapped, want) {
		t.Errorf("unmapped = %v, want %v", unmapped, want)
	}
	if _, ok := ported.Base["nonesuch"]; ok {
		t.Errorf("Base contains slot %q, want it absent", "nonesuch")
	}
	if _, ok := ported.Base["trigger-icon"]; ok {
		t.Errorf("Base contains slot %q, want it absent", "trigger-icon")
	}
}

// Fallback: a shape slot upstream never mentions (this fixture's Section
// has no cn-accordion-content rule at all) gets the nova rule's utilities,
// marked Carried; the slot upstream DID cover must not be marked Carried.
func TestTransformFallbackCarriesFromNova(t *testing.T) {
	shape := recipe.Shape{
		Component: "accordion",
		Slots: []recipe.Slot{
			{Name: "trigger", Base: true},
			{Name: "content", Base: true},
		},
	}
	sec := Section{
		Name: "Accordion",
		Rules: []Rule{
			{Class: "cn-accordion-trigger", Utilities: []string{"p-4"}},
		},
	}
	fallback := mustParseStyle(t, fmt.Sprintf(`@layer components {
  .%s {
    @apply overflow-hidden;
  }
}
`, shape.BaseClass("content")))

	ported, unmapped, err := Transform(shape, sec, fallback)
	if err != nil {
		t.Fatalf("Transform: %v", err)
	}
	if len(unmapped) != 0 {
		t.Fatalf("unmapped = %v, want none", unmapped)
	}
	if want := []string{"overflow-hidden"}; !reflect.DeepEqual(ported.Base["content"], want) {
		t.Errorf("Base[content] = %v, want %v", ported.Base["content"], want)
	}
	contentClass := shape.BaseClass("content")
	if !ported.Carried[contentClass] {
		t.Errorf("Carried[%q] = false, want true", contentClass)
	}
	triggerClass := shape.BaseClass("trigger")
	if ported.Carried[triggerClass] {
		t.Errorf("Carried[%q] = true, want false (covered by upstream)", triggerClass)
	}
}

// Fallback, missing: the shape declares a "content" slot upstream never
// mentions and the fallback Style has nothing for it either — a hard error,
// not a silently incomplete Ported.
func TestTransformFallbackMissingIsError(t *testing.T) {
	shape := recipe.Shape{
		Component: "accordion",
		Slots: []recipe.Slot{
			{Name: "trigger", Base: true},
			{Name: "content", Base: true},
		},
	}
	sec := Section{
		Name: "Accordion",
		Rules: []Rule{
			{Class: "cn-accordion-trigger", Utilities: []string{"p-4"}},
		},
	}

	_, _, err := Transform(shape, sec, recipe.Style{})
	if err == nil {
		t.Fatal("Transform: want error, got nil")
	}
	if !strings.Contains(err.Error(), "content") {
		t.Errorf("error %q does not mention the unresolvable slot %q", err.Error(), "content")
	}
}
