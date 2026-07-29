package recipe

import (
	"reflect"
	"slices"
	"strings"
	"testing"
)

func conformShape() Shape {
	return Shape{
		Component: "button",
		Slots: []Slot{{
			Name: "", Base: true,
			Dimensions: []Dimension{
				{Name: "variant", Default: "default", Values: []string{"default", "outline"}},
			},
		}},
	}
}

const conformCSS = `@layer components {
  .gsxui-recipe-button { @apply inline-flex items-center; }
  .gsxui-recipe-button-variant-default { @apply bg-primary; }
  .gsxui-recipe-button-variant-outline { @apply border-border bg-background; }
}`

func mustParse(t *testing.T, src string) Style {
	t.Helper()
	style, err := ParseStyle("nova/button.css", []byte(src))
	if err != nil {
		t.Fatalf("ParseStyle() error = %v", err)
	}
	return style
}

func TestConformAcceptsCompleteStyle(t *testing.T) {
	t.Parallel()
	resolved, err := Conform("nova/button.css", conformShape(), mustParse(t, conformCSS))
	if err != nil {
		t.Fatalf("Conform() error = %v", err)
	}
	if got, want := resolved.BaseUtilities(""), []string{"inline-flex", "items-center"}; !reflect.DeepEqual(got, want) {
		t.Errorf("BaseUtilities(root) = %q, want %q", got, want)
	}
	if got, want := resolved.Utilities("", "variant", "outline"), []string{"border-border", "bg-background"}; !reflect.DeepEqual(got, want) {
		t.Errorf("Utilities() = %q, want %q", got, want)
	}
}

func TestConformRejectsMissingValue(t *testing.T) {
	t.Parallel()
	src := strings.Replace(conformCSS,
		"  .gsxui-recipe-button-variant-outline { @apply border-border bg-background; }\n", "", 1)
	_, err := Conform("maia/button.css", conformShape(), mustParse(t, src))
	if err == nil {
		t.Fatal("Conform() = nil error, want error")
	}
	want := `maia/button.css: slot "" dimension "variant" missing value "outline"`
	if got := err.Error(); got != want {
		t.Errorf("Conform() = %q, want %q", got, want)
	}
}

func TestConformRejectsMissingBase(t *testing.T) {
	t.Parallel()
	src := strings.Replace(conformCSS,
		"  .gsxui-recipe-button { @apply inline-flex items-center; }\n", "", 1)
	_, err := Conform("maia/button.css", conformShape(), mustParse(t, src))
	if err == nil {
		t.Fatal("Conform() = nil error, want error")
	}
	want := `maia/button.css: slot "" missing base rule .gsxui-recipe-button`
	if got := err.Error(); got != want {
		t.Errorf("Conform() = %q, want %q", got, want)
	}
}

func TestConformRejectsUnknownRule(t *testing.T) {
	t.Parallel()
	src := strings.Replace(conformCSS, ".gsxui-recipe-button-variant-outline { @apply border-border bg-background; }\n",
		".gsxui-recipe-button-variant-outline { @apply border-border bg-background; }\n  .gsxui-recipe-button-variant-plain { @apply bg-muted; }\n", 1)
	_, err := Conform("maia/button.css", conformShape(), mustParse(t, src))
	if err == nil {
		t.Fatal("Conform() = nil error, want error")
	}
	if got := err.Error(); !strings.Contains(got, `does not declare value "plain"`) {
		t.Errorf("Conform() = %q, want a message about the undeclared value", got)
	}
}

func TestConformRejectsInvalidShape(t *testing.T) {
	t.Parallel()
	shape := conformShape()
	shape.Slots[0].Dimensions[0].Default = "ghost"
	_, err := Conform("nova/button.css", shape, mustParse(t, conformCSS))
	if err == nil {
		t.Fatal("Conform() = nil error, want error")
	}
	if got, want := err.Error(), `button: slot "" dimension "variant" default "ghost" is not one of its values`; got != want {
		t.Errorf("Conform() = %q, want %q", got, want)
	}
}

// noConflictMerger keeps every class, standing in for a merge with no conflicts.
func noConflictMerger(classes []string) string { return strings.Join(classes, " ") }

func TestCheckConflictsAcceptsCleanLists(t *testing.T) {
	t.Parallel()
	resolved, err := Conform("nova/button.css", conformShape(), mustParse(t, conformCSS))
	if err != nil {
		t.Fatal(err)
	}
	if err := CheckConflicts("nova/button.css", resolved, noConflictMerger); err != nil {
		t.Fatalf("CheckConflicts() = %v, want nil", err)
	}
}

func TestCheckConflictsRejectsSupersededUtility(t *testing.T) {
	t.Parallel()
	src := strings.Replace(conformCSS,
		"@apply inline-flex items-center;", "@apply rounded-lg rounded-md;", 1)
	resolved, err := Conform("nova/button.css", conformShape(), mustParse(t, src))
	if err != nil {
		t.Fatal(err)
	}
	// Stub merger: drops "rounded-lg" when "rounded-md" follows it.
	stub := func(classes []string) string {
		var kept []string
		for _, class := range classes {
			if class == "rounded-lg" && slices.Contains(classes, "rounded-md") {
				continue
			}
			kept = append(kept, class)
		}
		return strings.Join(kept, " ")
	}
	err = CheckConflicts("nova/button.css", resolved, stub)
	if err == nil {
		t.Fatal("CheckConflicts() = nil, want error")
	}
	want := "nova/button.css: .gsxui-recipe-button applies conflicting utilities: rounded-lg is superseded"
	if got := err.Error(); got != want {
		t.Errorf("CheckConflicts() = %q, want %q", got, want)
	}
}

func TestCheckConflictsRejectsRepeatedUtility(t *testing.T) {
	t.Parallel()
	src := strings.Replace(conformCSS,
		"  .gsxui-recipe-button { @apply inline-flex items-center; }", "  .gsxui-recipe-button { @apply inline-flex inline-flex; }", 1)
	resolved, err := Conform("nova/button.css", conformShape(), mustParse(t, src))
	if err != nil {
		t.Fatal(err)
	}
	// Stub merger: deduplicates repeated utilities.
	stub := func(classes []string) string {
		seen := make(map[string]struct{})
		var kept []string
		for _, class := range classes {
			if _, ok := seen[class]; !ok {
				seen[class] = struct{}{}
				kept = append(kept, class)
			}
		}
		return strings.Join(kept, " ")
	}
	err = CheckConflicts("nova/button.css", resolved, stub)
	if err == nil {
		t.Fatal("CheckConflicts() = nil, want error")
	}
	want := "nova/button.css: .gsxui-recipe-button repeats utility inline-flex"
	if got := err.Error(); got != want {
		t.Errorf("CheckConflicts() = %q, want %q", got, want)
	}
}

// TestCheckConflictsCombinedPrefersSuperseded covers a rule that both
// repeats a utility and carries a superseded one. The superseded/conflict
// error must win, deterministically, regardless of Go's randomized map
// iteration order — CheckConflicts walks the original utilities slice, not
// its internal count maps, for exactly this reason.
func TestCheckConflictsCombinedPrefersSuperseded(t *testing.T) {
	t.Parallel()
	src := strings.Replace(conformCSS,
		"  .gsxui-recipe-button { @apply inline-flex items-center; }",
		"  .gsxui-recipe-button { @apply rounded-lg inline-flex rounded-md inline-flex; }", 1)
	resolved, err := Conform("nova/button.css", conformShape(), mustParse(t, src))
	if err != nil {
		t.Fatal(err)
	}
	// Stub merger: drops "rounded-lg" when "rounded-md" follows it, and
	// deduplicates repeated utilities otherwise.
	stub := func(classes []string) string {
		var kept []string
		seen := make(map[string]struct{})
		for _, class := range classes {
			if class == "rounded-lg" && slices.Contains(classes, "rounded-md") {
				continue
			}
			if _, ok := seen[class]; ok {
				continue
			}
			seen[class] = struct{}{}
			kept = append(kept, class)
		}
		return strings.Join(kept, " ")
	}
	for i := 0; i < 20; i++ {
		err = CheckConflicts("nova/button.css", resolved, stub)
		if err == nil {
			t.Fatal("CheckConflicts() = nil, want error")
		}
		want := "nova/button.css: .gsxui-recipe-button applies conflicting utilities: rounded-lg is superseded"
		if got := err.Error(); got != want {
			t.Fatalf("CheckConflicts() = %q, want %q", got, want)
		}
	}
}

func slotConformShape() Shape {
	return Shape{
		Component: "card",
		Slots: []Slot{
			{Name: "", Base: true},
			{Name: "header", Base: true, Dimensions: []Dimension{
				{Name: "variant", Default: "default", Values: []string{"default", "muted"}},
			}},
		},
	}
}

const slotConformCSS = `@layer components {
  .gsxui-recipe-card { @apply rounded-xl border; }
  .gsxui-recipe-card-header { @apply flex gap-1.5; }
  .gsxui-recipe-card-header-variant-default { @apply text-foreground; }
  .gsxui-recipe-card-header-variant-muted { @apply text-muted-foreground; }
}`

func TestConformAcceptsSlottedStyle(t *testing.T) {
	t.Parallel()
	resolved, err := Conform("nova/card.css", slotConformShape(), mustParse(t, slotConformCSS))
	if err != nil {
		t.Fatalf("Conform() error = %v", err)
	}
	if got, want := resolved.BaseUtilities(""), []string{"rounded-xl", "border"}; !reflect.DeepEqual(got, want) {
		t.Errorf("BaseUtilities(root) = %q, want %q", got, want)
	}
	if got, want := resolved.BaseUtilities("header"), []string{"flex", "gap-1.5"}; !reflect.DeepEqual(got, want) {
		t.Errorf("BaseUtilities(header) = %q, want %q", got, want)
	}
	if got, want := resolved.Utilities("header", "variant", "muted"), []string{"text-muted-foreground"}; !reflect.DeepEqual(got, want) {
		t.Errorf("Utilities() = %q, want %q", got, want)
	}
}

func TestConformRejectsMissingSlotBase(t *testing.T) {
	t.Parallel()
	src := strings.Replace(slotConformCSS,
		"  .gsxui-recipe-card-header { @apply flex gap-1.5; }\n", "", 1)
	_, err := Conform("maia/card.css", slotConformShape(), mustParse(t, src))
	if err == nil {
		t.Fatal("Conform() = nil, want error")
	}
	want := `maia/card.css: slot "header" missing base rule .gsxui-recipe-card-header`
	if got := err.Error(); got != want {
		t.Errorf("Conform() = %q, want %q", got, want)
	}
}

func TestConformRejectsMissingSlotValue(t *testing.T) {
	t.Parallel()
	src := strings.Replace(slotConformCSS,
		"  .gsxui-recipe-card-header-variant-muted { @apply text-muted-foreground; }\n", "", 1)
	_, err := Conform("maia/card.css", slotConformShape(), mustParse(t, src))
	if err == nil {
		t.Fatal("Conform() = nil, want error")
	}
	want := `maia/card.css: slot "header" dimension "variant" missing value "muted"`
	if got := err.Error(); got != want {
		t.Errorf("Conform() = %q, want %q", got, want)
	}
}

func TestCheckConflictsNamesTheSlot(t *testing.T) {
	t.Parallel()
	src := strings.Replace(slotConformCSS,
		"  .gsxui-recipe-card-header { @apply flex gap-1.5; }",
		"  .gsxui-recipe-card-header { @apply rounded-lg rounded-md; }", 1)
	resolved, err := Conform("nova/card.css", slotConformShape(), mustParse(t, src))
	if err != nil {
		t.Fatal(err)
	}
	stub := func(classes []string) string {
		var kept []string
		for _, class := range classes {
			if class == "rounded-lg" && slices.Contains(classes, "rounded-md") {
				continue
			}
			kept = append(kept, class)
		}
		return strings.Join(kept, " ")
	}
	err = CheckConflicts("nova/card.css", resolved, stub)
	if err == nil {
		t.Fatal("CheckConflicts() = nil, want error")
	}
	want := "nova/card.css: .gsxui-recipe-card-header applies conflicting utilities: rounded-lg is superseded"
	if got := err.Error(); got != want {
		t.Errorf("CheckConflicts() = %q, want %q", got, want)
	}
}
