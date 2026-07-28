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

func TestRecipeTokensRejectsRecipeUseOutsideClassValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		body string
		want string
	}{
		{
			name: "value if condition",
			body: `<div class={
				if mode == "gsxui-recipe-button" {
					"safe"
				} else {
					"other"
				},
			}></div>`,
			want: "value-if condition",
		},
		{
			name: "value if init condition",
			body: `<div class={
				if x := "gsxui-recipe-button"; x != "" {
					"safe"
				} else {
					"other"
				},
			}></div>`,
			want: "value-if condition",
		},
		{
			name: "value switch tag",
			body: `<div class={
				switch "gsxui-recipe-button" {
				default:
					"safe"
				},
			}></div>`,
			want: "switch tag",
		},
		{
			name: "value switch init tag",
			body: `<div class={
				switch x := "gsxui-recipe-button"; x {
				default:
					"safe"
				},
			}></div>`,
			want: "switch tag",
		},
		{
			name: "value switch case expression",
			body: `<div class={
				switch mode {
				case "gsxui-recipe-button":
					"safe"
				default:
					"other"
				},
			}></div>`,
			want: "switch case",
		},
		{
			name: "pipeline stage argument",
			body: `<div class={
				"safe" |> decorate("gsxui-recipe-button", args...),
			}></div>`,
			want: "pipeline stage",
		},
		{
			name: "css segment",
			body: "<div style={ \"display:none\", css`color:gsxui-recipe-button` }></div>",
			want: "CSS segment",
		},
		{
			name: "css interpolation pipeline stage",
			body: "<div style={ \"display:none\", css`color:@{value |> decorate(\"gsxui-recipe-button\")}` }></div>",
			want: "CSS segment pipeline stage",
		},
		{
			name: "embedded interpolation pipeline stage",
			body: "<div class={f`prefix-@{value |> decorate(\"gsxui-recipe-button\")}`}></div>",
			want: "class interpolation pipeline stage",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertRecipeRejectedByBothAPIs(t, resolveComponent(tt.body), tt.want)
		})
	}
}

func TestRecipeTokensRejectsUnparseableSupportedMetadata(t *testing.T) {
	t.Parallel()

	src := []byte(resolveComponent(`<div data-value={x := 1}></div>`))
	if _, err := gsxparser.ParseFile(token.NewFileSet(), "unparseable-metadata.gsx", src, 0); err != nil {
		t.Fatalf("fixture must reach resolver metadata validation: %v", err)
	}
	if _, err := RecipeTokens("unparseable-metadata.gsx", src); err == nil {
		t.Fatal("RecipeTokens() error = nil, want metadata parse error")
	} else if !strings.Contains(err.Error(), "parse") {
		t.Errorf("RecipeTokens() error %q does not identify metadata parse failure", err)
	}
	if got, _, err := Resolve("unparseable-metadata.gsx", src, Recipes{}); err == nil {
		t.Fatalf("Resolve() = %s, nil error; want metadata parse error", got)
	}
}

func TestRecipeTokensRejectsPartiallyDynamicRecipeConcatenation(t *testing.T) {
	t.Parallel()

	src := resolveComponent(`<div class={ "gsxui-" + "recipe-" + suffix }></div>`)
	assertRecipeRejectedByBothAPIs(t, src, "concatenation")
}

func TestRecipeTokensRejectsRecipeAssembledAcrossEmbeddedSegments(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		src  string
	}{
		{
			name: "constant expression segment",
			src:  "<div class={f`gsxui-@{\"recipe-button\"}`}></div>",
		},
		{
			name: "partially dynamic expression segment",
			src:  "<div class={f`gsxui-@{\"recipe-\" + suffix}`}></div>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertRecipeRejectedByBothAPIs(t, resolveComponent(tt.src), "interpolation")
		})
	}
}

func TestResolveRejectsRecipeTokenIntroducedByUtility(t *testing.T) {
	t.Parallel()

	src := []byte(resolveComponent(`<div class={ "gsxui-recipe-button" }></div>`))
	recipes := testRecipes(
		Recipe{Token: "gsxui-recipe-button", Utilities: []string{"gsxui-recipe-leaked"}},
	)
	got, _, err := Resolve("introduced-token.gsx", src, recipes)
	if err == nil {
		t.Fatal("Resolve() error = nil, want structural output-invariant error")
	}
	if !strings.Contains(err.Error(), "gsxui-recipe-leaked") {
		t.Errorf("Resolve() error %q does not identify introduced token", err)
	}
	if bytes.Contains(got, []byte(RecipePrefix)) {
		t.Errorf("Resolve() returned output containing %q: %s", RecipePrefix, got)
	}
}

func TestResolveNoPrefixInvariantIsIndependentOfRecipeScanner(t *testing.T) {
	t.Parallel()

	src := []byte(resolveComponent(`{/* gsxui-recipe-output-invariant */}
		<div class={ "safe" }></div>`))
	tokens, err := RecipeTokens("output-invariant.gsx", src)
	if err != nil {
		t.Fatalf("RecipeTokens() error = %v", err)
	}
	if len(tokens) != 0 {
		t.Fatalf("RecipeTokens() = %q, want none for source-only comment", tokens)
	}
	got, _, err := Resolve("output-invariant.gsx", src, Recipes{})
	if err == nil {
		t.Fatal("Resolve() error = nil, want independent no-prefix invariant error")
	}
	if bytes.Contains(got, []byte(RecipePrefix)) {
		t.Errorf("Resolve() returned output containing %q: %s", RecipePrefix, got)
	}
}

func TestResolveAcceptsUnrelatedDynamicClassSurfaces(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		body string
	}{
		{
			name: "class expressions and metadata",
			body: `<div class={
				className,
				"safe" |> decorate(fallback),
				if active := enabled; active {
					enabledClass
				} else {
					"disabled"
				},
				switch selected := mode; selected {
				case alternate:
					alternateClass |> decorate(fallback)
				default:
					"other"
				},
			}></div>`,
		},
		{
			name: "interpolated class",
			body: "<div class={f`prefix-@{suffix |> decorate(fallback)}`}></div>",
		},
		{
			name: "css segments",
			body: "<div style={ \"display:none\", css`color:@{color |> decorate(fallback)}` }></div>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := []byte(resolveComponent(tt.body))
			if _, err := gsxparser.ParseFile(token.NewFileSet(), "dynamic.gsx", src, 0); err != nil {
				t.Fatalf("fixture does not parse: %v", err)
			}
			tokens, err := RecipeTokens("dynamic.gsx", src)
			if err != nil {
				t.Fatalf("RecipeTokens() error = %v", err)
			}
			if len(tokens) != 0 {
				t.Fatalf("RecipeTokens() = %q, want none", tokens)
			}
			got, report, err := Resolve("dynamic.gsx", src, Recipes{})
			if err != nil {
				t.Fatalf("Resolve() error = %v", err)
			}
			if len(report.UsedTokens) != 0 {
				t.Errorf("ResolveReport.UsedTokens = %q, want none", report.UsedTokens)
			}
			if _, err := gsxparser.ParseFile(token.NewFileSet(), "dynamic.gsx", got, 0); err != nil {
				t.Fatalf("resolved output does not parse: %v", err)
			}
		})
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

func assertRecipeRejectedByBothAPIs(t *testing.T, src, want string) {
	t.Helper()

	filename := "invalid-recipe-use.gsx"
	source := []byte(src)
	if _, err := gsxparser.ParseFile(token.NewFileSet(), filename, source, 0); err != nil {
		t.Fatalf("fixture does not parse: %v", err)
	}

	if _, err := RecipeTokens(filename, source); err == nil {
		t.Fatal("RecipeTokens() error = nil, want structural recipe-use error")
	} else if !strings.Contains(err.Error(), want) {
		t.Errorf("RecipeTokens() error %q does not contain %q", err, want)
	}

	got, _, err := Resolve(filename, source, Recipes{})
	if err == nil {
		t.Fatal("Resolve() error = nil, want structural recipe-use error")
	}
	if bytes.Contains(got, []byte(RecipePrefix)) {
		t.Errorf("Resolve() returned output containing %q: %s", RecipePrefix, got)
	}
	if !strings.Contains(err.Error(), want) {
		t.Errorf("Resolve() error %q does not contain %q", err, want)
	}
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
