package stylegen

import (
	goparser "go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/gsxhq/gsxui/internal/recipe"
)

// hyphenatedShape is input-group's shape. Eleven of the fifty-two components
// carry a hyphen in their registry name; this is the one the blocker was found
// on, and its Go identifier can only come from deriving it.
func hyphenatedShape() recipe.Shape {
	return recipe.Shape{
		Component: "input-group",
		Slots: []recipe.Slot{
			{Name: "", Base: true, Dimensions: []recipe.Dimension{
				{Name: "size", Default: "default", Values: []string{"default", "sm"}},
			}},
			{Name: "addon", Base: true},
		},
	}
}

func TestGenerateAccessorsDerivesAGoIdentifierForAHyphenatedComponent(t *testing.T) {
	t.Parallel()
	src, err := GenerateAccessors(hyphenatedShape())
	if err != nil {
		t.Fatalf("GenerateAccessors() error = %v", err)
	}
	// The whole point: the output has to be Go. `input-groupRecipe` does not
	// parse, so this is the assertion the blocker fails.
	if _, err := goparser.ParseFile(token.NewFileSet(), "input-group_recipe.gen.go", src, 0); err != nil {
		t.Fatalf("generated accessors do not parse: %v\n%s", err, src)
	}
	got := string(src)
	for _, want := range []string{
		"type inputGroupRecipe struct{ c recipe.Component }",
		"func (r inputGroupRecipe) Root() string { return r.c.SlotClass(\"\") }",
		"func (r inputGroupRecipe) Addon() string { return r.c.SlotClass(\"addon\") }",
		"func (r inputGroupRecipe) Size(value string) string",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
	// The registry name survives in the slot arguments — those are the CSS
	// identity — but never as a Go identifier.
	if strings.Contains(got, "input-groupRecipe") {
		t.Errorf("generated accessors use the registry name as a Go identifier:\n%s", got)
	}
}

func TestResolveBindsTheDerivedReceiverForAHyphenatedComponent(t *testing.T) {
	t.Parallel()
	shape := hyphenatedShape()
	style := testRecipes(t,
		recipe.Rule{Class: "gsxui-recipe-input-group", Utilities: []string{"flex", "items-center"}},
		recipe.Rule{Class: "gsxui-recipe-input-group-size-default", Utilities: []string{"h-9"}},
		recipe.Rule{Class: "gsxui-recipe-input-group-size-sm", Utilities: []string{"h-8"}},
		recipe.Rule{Class: "gsxui-recipe-input-group-addon", Utilities: []string{"px-2"}},
	)
	resolved, err := recipe.Conform("input-group.css", shape, style)
	if err != nil {
		t.Fatalf("recipe.Conform() error = %v", err)
	}

	// `input-group.Root()` would parse as subtraction, so the canonical source
	// receives the derived identifier.
	src := []byte(`package canonical

import "github.com/gsxhq/gsx"

component InputGroup(size string, children gsx.Node) {
	<div class={ inputGroup.Root(), inputGroup.Size(size) } data-gsxui-slot-input-group>
		<span class={ inputGroup.Addon() } data-gsxui-slot-input-group-addon>{ children }</span>
	</div>
}
`)
	calls, err := HelperCalls("input-group.gsx", src, shape)
	if err != nil {
		t.Fatalf("HelperCalls() error = %v", err)
	}
	// Calls are reported by REGISTRY name, not by the Go receiver: everything
	// downstream (shape lookup, markers, diagnostics) speaks the registry name.
	want := []Call{
		{Component: "input-group"},
		{Component: "input-group", Dimension: "size"},
		{Component: "input-group", Slot: "addon"},
	}
	if len(calls) != len(want) {
		t.Fatalf("HelperCalls() = %+v, want %+v", calls, want)
	}
	for i := range want {
		if calls[i] != want[i] {
			t.Errorf("call %d = %+v, want %+v", i, calls[i], want[i])
		}
	}

	got, err := Resolve("input-group.gsx", src, resolved)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	for _, fragment := range []string{
		`"flex items-center"`,
		`case "sm":`,
		`"h-8"`,
		`"px-2"`,
	} {
		if !strings.Contains(string(got), fragment) {
			t.Errorf("missing %q\nin: %s", fragment, got)
		}
	}
	if strings.Contains(string(got), "inputGroup.") {
		t.Errorf("resolved source still contains an accessor call:\n%s", got)
	}
	// The markers are the CSS identity and must have survived untouched.
	if !strings.Contains(string(got), "data-gsxui-slot-input-group-addon") {
		t.Errorf("resolved source lost the registry-named slot marker:\n%s", got)
	}
}

func TestRegistryNameStillNamesTheRecipeClasses(t *testing.T) {
	t.Parallel()
	shape := hyphenatedShape()
	// Deriving a Go identifier must not touch the CSS identity: the classes a
	// style authors, and therefore every selector in every registry stylesheet,
	// still carry the hyphenated registry name.
	if got, want := shape.BaseClass(""), "gsxui-recipe-input-group"; got != want {
		t.Errorf("BaseClass(\"\") = %q, want %q", got, want)
	}
	if got, want := shape.BaseClass("addon"), "gsxui-recipe-input-group-addon"; got != want {
		t.Errorf("BaseClass(\"addon\") = %q, want %q", got, want)
	}
	if got, want := shape.ValueClass("", "size", "sm"), "gsxui-recipe-input-group-size-sm"; got != want {
		t.Errorf("ValueClass = %q, want %q", got, want)
	}
	if got, want := slotMarker("input-group", "addon"), "data-gsxui-slot-input-group-addon"; got != want {
		t.Errorf("slotMarker = %q, want %q", got, want)
	}
	// And the derived identifier is emphatically not one of them.
	identifier, err := componentIdentifier(shape.Component)
	if err != nil {
		t.Fatalf("componentIdentifier() error = %v", err)
	}
	if identifier != "inputGroup" {
		t.Fatalf("componentIdentifier(%q) = %q, want %q", shape.Component, identifier, "inputGroup")
	}
	if strings.Contains(shape.BaseClass(""), identifier) {
		t.Errorf("recipe class %q leaked the Go identifier %q", shape.BaseClass(""), identifier)
	}
}

func TestComponentIdentifierRejectsIllegalDerivations(t *testing.T) {
	t.Parallel()
	tests := []struct {
		component string
		want      string
	}{
		{component: "2-column", want: "starts with a digit"},
		{component: "-", want: "empty Go identifier"},
		{component: "input group", want: "not a valid Go identifier"},
	}
	for _, tt := range tests {
		t.Run(tt.component, func(t *testing.T) {
			t.Parallel()
			_, err := componentIdentifier(tt.component)
			if err == nil {
				t.Fatalf("componentIdentifier(%q) error = nil, want error", tt.component)
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Errorf("componentIdentifier(%q) error %q does not contain %q", tt.component, err, tt.want)
			}
			if !strings.Contains(err.Error(), tt.component) {
				t.Errorf("componentIdentifier(%q) error %q does not name the component", tt.component, err)
			}
		})
	}
}

// TestComponentIdentifierSuffixesKeywordsAndPredeclaredIdentifiers is the
// replacement for the old rejection behavior: switch/select are real,
// already-shipped registry names (data-gsxui-slot-switch,
// data-gsxui-slot-select) that must not be renamed just because their Go
// identifier collides with a keyword. componentIdentifier instead appends a
// single trailing underscore, Go's own convention for this exact collision.
func TestComponentIdentifierSuffixesKeywordsAndPredeclaredIdentifiers(t *testing.T) {
	t.Parallel()
	tests := []struct {
		component string
		want      string
	}{
		{component: "switch", want: "switch_"},
		{component: "select", want: "select_"},
		{component: "range", want: "range_"},
		{component: "type", want: "type_"},
		{component: "string", want: "string_"},
		{component: "len", want: "len_"},
	}
	for _, tt := range tests {
		t.Run(tt.component, func(t *testing.T) {
			t.Parallel()
			got, err := componentIdentifier(tt.component)
			if err != nil {
				t.Fatalf("componentIdentifier(%q) error = %v, want no error", tt.component, err)
			}
			if got != tt.want {
				t.Errorf("componentIdentifier(%q) = %q, want %q", tt.component, got, tt.want)
			}
			if !token.IsIdentifier(got) || token.IsKeyword(got) {
				t.Errorf("componentIdentifier(%q) = %q is not a legal, non-keyword Go identifier", tt.component, got)
			}
		})
	}
}

func TestGenerateAccessorsSuffixesAKeywordComponent(t *testing.T) {
	t.Parallel()
	// The suffixing has to reach the generator, not just the helper: a
	// component named "range" would otherwise emit `type rangeRecipe` (which
	// parses fine) but `var range = rangeRecipe{...}` and a canonical source
	// authoring `range.Root()` do not — "range" the bare receiver is the
	// actual keyword collision, not the type name.
	shape := recipe.Shape{
		Component: "range",
		Slots:     []recipe.Slot{{Name: "", Base: true}},
	}
	src, err := GenerateAccessors(shape)
	if err != nil {
		t.Fatalf("GenerateAccessors() error = %v, want no error", err)
	}
	if _, err := goparser.ParseFile(token.NewFileSet(), "range_recipe.gen.go", src, 0); err != nil {
		t.Fatalf("generated accessors do not parse: %v\n%s", err, src)
	}
	got := string(src)
	for _, want := range []string{
		"type range_Recipe struct{ c recipe.Component }",
		"func (r range_Recipe) Root() string { return r.c.SlotClass(\"\") }",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

// TestRegistryNameStillNamesTheRecipeClassesForAKeywordComponent proves both
// halves hold simultaneously for switch/select, the way
// TestRegistryNameStillNamesTheRecipeClasses already does for the hyphenated
// case: the Go identifier is suffixed (switch_/select_), but the CSS/registry
// identity — recipe class names and slot markers — is untouched.
func TestRegistryNameStillNamesTheRecipeClassesForAKeywordComponent(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		component  string
		identifier string
	}{
		{component: "switch", identifier: "switch_"},
		{component: "select", identifier: "select_"},
	} {
		t.Run(tt.component, func(t *testing.T) {
			t.Parallel()
			shape := recipe.Shape{
				Component: tt.component,
				Slots:     []recipe.Slot{{Name: "", Base: true}},
			}
			if got, want := shape.BaseClass(""), "gsxui-recipe-"+tt.component; got != want {
				t.Errorf("BaseClass(\"\") = %q, want %q", got, want)
			}
			if got, want := slotMarker(tt.component, ""), "data-gsxui-slot-"+tt.component; got != want {
				t.Errorf("slotMarker = %q, want %q", got, want)
			}
			identifier, err := componentIdentifier(tt.component)
			if err != nil {
				t.Fatalf("componentIdentifier() error = %v", err)
			}
			if identifier != tt.identifier {
				t.Fatalf("componentIdentifier(%q) = %q, want %q", tt.component, identifier, tt.identifier)
			}
			if strings.Contains(shape.BaseClass(""), identifier) {
				t.Errorf("recipe class %q leaked the Go identifier %q", shape.BaseClass(""), identifier)
			}
		})
	}
}

func TestCheckComponentIdentifiersRejectsCollisionsNamingBoth(t *testing.T) {
	t.Parallel()
	tests := [][]string{
		{"input-group", "inputGroup"},
		{"input-group", "input_group"},
		{"a-1b", "a1b"},
		// A hypothetical component literally named "switch_" collides with
		// "switch" the same way "input_group" collides with "input-group":
		// identifierFold strips both hyphens and underscores, so "switch_"
		// folds to the same string "switch" does, independently of the fact
		// that componentIdentifier("switch") itself now derives "switch_".
		{"switch", "switch_"},
	}
	for _, pair := range tests {
		t.Run(strings.Join(pair, "+"), func(t *testing.T) {
			t.Parallel()
			err := checkComponentIdentifiers(pair)
			if err == nil {
				t.Fatalf("checkComponentIdentifiers(%q) error = nil, want error", pair)
			}
			for _, name := range pair {
				if !strings.Contains(err.Error(), name) {
					t.Errorf("checkComponentIdentifiers(%q) error %q does not name %q", pair, err, name)
				}
			}
		})
	}
}

func TestGenerateAccessorsRejectsAComponentCollidingWithTheRegistry(t *testing.T) {
	t.Parallel()
	// "Button" derives the same identifier as the registry's own "button", so
	// both would generate `buttonRecipe` in package canonical and neither
	// could be named unambiguously by an accessor receiver.
	shape := recipe.Shape{
		Component: "Button",
		Slots:     []recipe.Slot{{Name: "", Base: true}},
	}
	_, err := GenerateAccessors(shape)
	if err == nil {
		t.Fatal("GenerateAccessors() error = nil, want error")
	}
	if !strings.Contains(err.Error(), `"Button"`) || !strings.Contains(err.Error(), `"button"`) {
		t.Errorf("GenerateAccessors() error %q does not name both colliding components", err)
	}
}

func TestComponentIdentifierLeavesUnhyphenatedComponentsAlone(t *testing.T) {
	t.Parallel()
	// Button is already migrated. Its derived identifier must be byte-identical
	// to the name it already uses, or the fix rewrites generated output.
	for _, component := range []string{"button", "badge", "card", "avatar"} {
		identifier, err := componentIdentifier(component)
		if err != nil {
			t.Fatalf("componentIdentifier(%q) error = %v", component, err)
		}
		if identifier != component {
			t.Errorf("componentIdentifier(%q) = %q, want it unchanged", component, identifier)
		}
	}
}
