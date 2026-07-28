package stylegen

import (
	"bytes"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	gsxast "github.com/gsxhq/gsx/ast"
	"github.com/gsxhq/gsx/gen"
	gsxparser "github.com/gsxhq/gsx/parser"
)

func TestResolveGolden(t *testing.T) {
	t.Parallel()

	fixture := filepath.Join("testdata", "resolve-input.gsx.txt")
	filename := strings.TrimSuffix(fixture, ".txt")
	src, err := os.ReadFile(fixture)
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile(filepath.Join("testdata", "resolve-nova.golden.gsx.txt"))
	if err != nil {
		t.Fatal(err)
	}

	got, report, err := Resolve(filename, src, resolveTestRecipes())
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("Resolve() output mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
	if want := []string{
		"gsxui-recipe-button",
		"gsxui-recipe-button-size-default",
		"gsxui-recipe-button-size-icon",
		"gsxui-recipe-button-variant-default",
		"gsxui-recipe-button-variant-outline",
	}; !reflect.DeepEqual(report.UsedTokens, want) {
		t.Errorf("ResolveReport.UsedTokens = %q, want %q", report.UsedTokens, want)
	}
	if bytes.Contains(got, []byte(RecipePrefix)) {
		t.Errorf("Resolve() output still contains %q", RecipePrefix)
	}
	if !bytes.Contains(got, []byte("// Preserve this comment adjacent to the class switches.")) {
		t.Error("Resolve() removed the class-switch comment")
	}
	if !bytes.Contains(got, []byte(`sm:hover:[&>svg]:size-4`)) {
		t.Error("Resolve() changed an unrelated arbitrary Tailwind class")
	}
	if !bytes.Contains(got, []byte(`const unrelatedRecipeDocumentation = "recipe names belong only in canonical class expressions"`)) {
		t.Error("Resolve() changed an unrelated string literal")
	}

	formatted, err := gen.Format(filename, got)
	if err != nil {
		t.Fatalf("gen.Format(resolved) error = %v", err)
	}
	if !bytes.Equal(formatted, got) {
		t.Errorf("resolved output is not idempotently formatted\n--- formatted ---\n%s\n--- resolved ---\n%s", formatted, got)
	}
	if _, err := gsxparser.ParseFile(token.NewFileSet(), filename, got, 0); err != nil {
		t.Fatalf("parser.ParseFile(resolved) error = %v", err)
	}

	anchorClass, buttonClass := resolvedElementClassExpressions(t, filename, got)
	if !bytes.Equal(anchorClass, buttonClass) {
		t.Errorf("resolved anchor/button class expressions differ\n--- anchor ---\n%s\n--- button ---\n%s", anchorClass, buttonClass)
	}
}

func TestRecipeTokensReturnsSortedUniqueCanonicalUses(t *testing.T) {
	t.Parallel()

	src := []byte(`package p
component C() {
	<div class={
		"gsxui-recipe-z arbitrary",
		switch mode {
		case "a":
			"gsxui-recipe-a gsxui-recipe-z"
		default:
			"gsxui-recipe-m"
		},
	}></div>
}
`)
	got, err := RecipeTokens("tokens.gsx", src)
	if err != nil {
		t.Fatalf("RecipeTokens() error = %v", err)
	}
	want := []string{"gsxui-recipe-a", "gsxui-recipe-m", "gsxui-recipe-z"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("RecipeTokens() = %q, want %q", got, want)
	}
}

func TestResolveRejectsInvalidRecipeUse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		src     string
		recipes Recipes
		want    []string
	}{
		{
			name: "missing used token and unused declaration reported together",
			src:  resolveComponent(`<div class={ "gsxui-recipe-missing" }></div>`),
			recipes: testRecipes(
				Recipe{Token: "gsxui-recipe-unused", Utilities: []string{"grid"}},
			),
			want: []string{"missing", "gsxui-recipe-missing", "unused", "gsxui-recipe-unused"},
		},
		{
			name: "unused declaration",
			src:  resolveComponent(`<div class={ "gsxui-recipe-button" }></div>`),
			recipes: testRecipes(
				Recipe{Token: "gsxui-recipe-button", Utilities: []string{"flex"}},
				Recipe{Token: "gsxui-recipe-unused", Utilities: []string{"grid"}},
			),
			want: []string{"unused", "gsxui-recipe-unused"},
		},
		{
			name: "duplicate declaration",
			src:  resolveComponent(`<div class={ "gsxui-recipe-button" }></div>`),
			recipes: Recipes{
				ordered: []Recipe{
					{Token: "gsxui-recipe-button", Utilities: []string{"flex"}},
					{Token: "gsxui-recipe-button", Utilities: []string{"grid"}},
				},
				byToken: map[string]Recipe{
					"gsxui-recipe-button": {Token: "gsxui-recipe-button", Utilities: []string{"grid"}},
				},
			},
			want: []string{"duplicate", "gsxui-recipe-button"},
		},
		{
			name:    "static class",
			src:     resolveComponent(`<div class="gsxui-recipe-button"></div>`),
			recipes: testRecipes(Recipe{Token: "gsxui-recipe-button", Utilities: []string{"flex"}}),
			want:    []string{"static class", "gsxui-recipe-button"},
		},
		{
			name:    "embedded in larger class token",
			src:     resolveComponent(`<div class={ "prefix-gsxui-recipe-button" }></div>`),
			recipes: testRecipes(Recipe{Token: "gsxui-recipe-button", Utilities: []string{"flex"}}),
			want:    []string{"whole", "prefix-gsxui-recipe-button"},
		},
		{
			name:    "assembled by concatenation",
			src:     resolveComponent(`<div class={ "gsxui-recipe-" + suffix }></div>`),
			recipes: testRecipes(Recipe{Token: "gsxui-recipe-button", Utilities: []string{"flex"}}),
			want:    []string{"literal", "concatenation"},
		},
		{
			name:    "assembled from split constant concatenation",
			src:     resolveComponent(`<div class={ "gsxui-" + "recipe-button" }></div>`),
			recipes: testRecipes(Recipe{Token: "gsxui-recipe-button", Utilities: []string{"flex"}}),
			want:    []string{"literal", "concatenation"},
		},
		{
			name:    "assembled by interpolation",
			src:     resolveComponent("<div class={ f`gsxui-recipe-@{suffix}` }></div>"),
			recipes: testRecipes(Recipe{Token: "gsxui-recipe-button", Utilities: []string{"flex"}}),
			want:    []string{"literal", "interpolation"},
		},
		{
			name:    "token in non-class expression",
			src:     resolveComponent(`<div title={ "gsxui-recipe-button" } class={ "safe" }></div>`),
			recipes: testRecipes(Recipe{Token: "gsxui-recipe-button", Utilities: []string{"flex"}}),
			want:    []string{"non-class", "gsxui-recipe-button"},
		},
		{
			name:    "non-string class expression containing recipe",
			src:     resolveComponent(`<div class={ choose("gsxui-recipe-button") }></div>`),
			recipes: testRecipes(Recipe{Token: "gsxui-recipe-button", Utilities: []string{"flex"}}),
			want:    []string{"string literal", "gsxui-recipe-button"},
		},
		{
			name:    "malformed gsx",
			src:     "package p\ncomponent C() { <div>",
			recipes: testRecipes(Recipe{Token: "gsxui-recipe-button", Utilities: []string{"flex"}}),
			want:    []string{"malformed.gsx"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := filepath.Join("testdata", "malformed.gsx")
			_, _, err := Resolve(filename, []byte(tt.src), tt.recipes)
			if err == nil {
				t.Fatal("Resolve() error = nil, want error")
			}
			for _, want := range tt.want {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("Resolve() error %q does not contain %q", err, want)
				}
			}
		})
	}
}

func resolveTestRecipes() Recipes {
	return testRecipes(
		Recipe{Token: "gsxui-recipe-button-size-default", Utilities: []string{"h-9", "px-4"}},
		Recipe{Token: "gsxui-recipe-button", Utilities: []string{"inline-flex", "items-center"}},
		Recipe{Token: "gsxui-recipe-button-variant-outline", Utilities: []string{"border", "border-input", "bg-background"}},
		Recipe{Token: "gsxui-recipe-button-size-icon", Utilities: []string{"size-9", "[&>svg]:shrink-0"}},
		Recipe{Token: "gsxui-recipe-button-variant-default", Utilities: []string{"bg-primary", "text-primary-foreground"}},
	)
}

func testRecipes(recipes ...Recipe) Recipes {
	result := Recipes{
		ordered: append([]Recipe(nil), recipes...),
		byToken: make(map[string]Recipe, len(recipes)),
	}
	for _, recipe := range recipes {
		result.byToken[recipe.Token] = recipe
	}
	return result
}

func resolveComponent(body string) string {
	return "package p\ncomponent C() {\n\t" + body + "\n}\n"
}

func resolvedElementClassExpressions(t *testing.T, filename string, src []byte) ([]byte, []byte) {
	t.Helper()

	fset := token.NewFileSet()
	file, err := gsxparser.ParseFile(fset, filename, src, 0)
	if err != nil {
		t.Fatalf("parser.ParseFile() error = %v", err)
	}
	classes := make(map[string][]byte)
	gsxast.Inspect(file, func(node gsxast.Node) bool {
		element, ok := node.(*gsxast.Element)
		if !ok || (element.Tag != "a" && element.Tag != "button") {
			return true
		}
		for _, attr := range element.Attrs {
			class, ok := attr.(*gsxast.ClassAttr)
			if !ok || class.Name != "class" {
				continue
			}
			start := fset.Position(class.Pos()).Offset
			end := fset.Position(class.End()).Offset
			classes[element.Tag] = append([]byte(nil), src[start:end]...)
		}
		return true
	})
	if classes["a"] == nil || classes["button"] == nil {
		t.Fatalf("class expressions found = %v, want a and button", classes)
	}
	return classes["a"], classes["button"]
}
