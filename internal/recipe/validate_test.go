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
		Base:      true,
		Dimensions: []Dimension{
			{Name: "variant", Default: "default", Values: []string{"default", "outline"}},
		},
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
	if got, want := resolved.Base, []string{"inline-flex", "items-center"}; !reflect.DeepEqual(got, want) {
		t.Errorf("Base = %q, want %q", got, want)
	}
	if got, want := resolved.Utilities("variant", "outline"), []string{"border-border", "bg-background"}; !reflect.DeepEqual(got, want) {
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
	want := `maia/button.css: dimension "variant" missing value "outline"`
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
	want := `maia/button.css: missing base rule .gsxui-recipe-button`
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
	shape.Dimensions[0].Default = "ghost"
	_, err := Conform("nova/button.css", shape, mustParse(t, conformCSS))
	if err == nil {
		t.Fatal("Conform() = nil error, want error")
	}
	if got, want := err.Error(), `button: dimension "variant" default "ghost" is not one of its values`; got != want {
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
