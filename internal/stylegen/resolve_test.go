package stylegen

import (
	"bytes"
	"fmt"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"

	gsxast "github.com/gsxhq/gsx/ast"
	"github.com/gsxhq/gsx/gen"
	gsxparser "github.com/gsxhq/gsx/parser"

	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical"
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

	got, report, err := resolveTokens(filename, src, resolveTestRecipes(t))
	if err != nil {
		t.Fatalf("resolveTokens() error = %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("resolveTokens() output mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
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
	if bytes.Contains(got, []byte(recipe.Prefix)) {
		t.Errorf("resolveTokens() output still contains %q", recipe.Prefix)
	}
	if !bytes.Contains(got, []byte("// Preserve this comment adjacent to the class switches.")) {
		t.Error("resolveTokens() removed the class-switch comment")
	}
	if !bytes.Contains(got, []byte(`sm:hover:[&>svg]:size-4`)) {
		t.Error("resolveTokens() changed an unrelated arbitrary Tailwind class")
	}
	if !bytes.Contains(got, []byte(`const unrelatedRecipeDocumentation = "recipe names belong only in canonical class expressions"`)) {
		t.Error("resolveTokens() changed an unrelated string literal")
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

func TestResolveAcceptsElementValuedExpressions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		body string
		want string
	}{
		{
			name: "element-valued attribute",
			body: `<div title={<strong>Bold</strong>}></div>`,
			want: `<strong>Bold</strong>`,
		},
		{
			name: "element-valued call argument",
			body: `{wrap(<b>hi</b>)}`,
			want: `wrap(<b>hi</b>)`,
		},
		{
			name: "element-valued call argument with nested class",
			body: `{wrap(<b class={"safe"}>hi</b>)}`,
			want: `wrap(<b class={"safe"}>hi</b>)`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := []byte(resolveComponent(tt.body))
			if _, err := gsxparser.ParseFile(token.NewFileSet(), "element-valued.gsx", src, 0); err != nil {
				t.Fatalf("fixture does not parse: %v", err)
			}
			tokens, err := RecipeTokens("element-valued.gsx", src)
			if err != nil {
				t.Fatalf("RecipeTokens() error = %v", err)
			}
			if len(tokens) != 0 {
				t.Fatalf("RecipeTokens() = %q, want none", tokens)
			}
			wantFormatted, err := gen.Format("element-valued.gsx", src)
			if err != nil {
				t.Fatalf("gen.Format() error = %v", err)
			}
			got, report, err := resolveTokens("element-valued.gsx", src, recipe.Style{})
			if err != nil {
				t.Fatalf("resolveTokens() error = %v", err)
			}
			if len(report.UsedTokens) != 0 {
				t.Errorf("ResolveReport.UsedTokens = %q, want none", report.UsedTokens)
			}
			if !bytes.Equal(got, wantFormatted) {
				t.Errorf("resolveTokens() changed recipe-free element-valued expression\n--- got ---\n%s\n--- want ---\n%s", got, wantFormatted)
			}
			if !bytes.Contains(got, []byte(tt.want)) {
				t.Errorf("resolveTokens() output does not preserve %q:\n%s", tt.want, got)
			}
			if bytes.Contains(got, []byte(recipe.Prefix)) {
				t.Errorf("resolveTokens() output still contains %q", recipe.Prefix)
			}
		})
	}
}

func TestRecipeTokensRejectsRecipeUseInsideElementValuedExpressions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		body string
		want string
	}{
		{
			name: "nested element attribute",
			body: `<div title={<strong data-token={"gsxui-recipe-button"}>Bold</strong>}></div>`,
			want: "data-token",
		},
		{
			name: "nested element interpolation",
			body: `{wrap(<b>{"gsxui-recipe-button"}</b>)}`,
			want: "interpolation",
		},
		{
			name: "go argument beside element",
			body: `{wrap("gsxui-recipe-button", <b>hi</b>)}`,
			want: "interpolation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := "element-valued-recipe.gsx"
			source := []byte(resolveComponent(tt.body))
			if _, err := gsxparser.ParseFile(token.NewFileSet(), filename, source, 0); err != nil {
				t.Fatalf("fixture does not parse: %v", err)
			}
			if _, err := RecipeTokens(filename, source); err == nil {
				t.Fatal("RecipeTokens() error = nil, want structural recipe-use error")
			} else {
				for _, want := range []string{"recipe token", tt.want} {
					if !strings.Contains(err.Error(), want) {
						t.Errorf("RecipeTokens() error %q does not contain %q", err, want)
					}
				}
			}
			got, _, err := resolveTokens(filename, source, recipe.Style{})
			if err == nil {
				t.Fatal("resolveTokens() error = nil, want structural recipe-use error")
			}
			if !strings.Contains(err.Error(), "recipe token") {
				t.Errorf("resolveTokens() error %q does not identify recipe-token misuse", err)
			}
			if bytes.Contains(got, []byte(recipe.Prefix)) {
				t.Errorf("resolveTokens() returned output containing %q: %s", recipe.Prefix, got)
			}
		})
	}
}

func TestRecipeTokensRejectsRecipeContentInNestedElementValuedExpressions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		body string
		want string
	}{
		{
			name: "nested element class literal",
			body: `{wrap(<b class={"gsxui-recipe-button"}>hi</b>)}`,
			want: "nested class expression",
		},
		{
			name: "nested element text",
			body: `{wrap(<b>gsxui-recipe-button</b>)}`,
			want: "nested element text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := "nested-element-recipe.gsx"
			source := []byte(resolveComponent(tt.body))
			if _, err := gsxparser.ParseFile(token.NewFileSet(), filename, source, 0); err != nil {
				t.Fatalf("fixture does not parse: %v", err)
			}

			if _, err := RecipeTokens(filename, source); err == nil {
				t.Fatal("RecipeTokens() error = nil, want nested recipe-use error")
			} else {
				for _, want := range []string{"recipe token", tt.want} {
					if !strings.Contains(err.Error(), want) {
						t.Errorf("RecipeTokens() error %q does not contain %q", err, want)
					}
				}
			}

			recipes := testRecipes(t,
				recipe.Rule{Class: "gsxui-recipe-button", Utilities: []string{"resolved-utility"}},
			)
			got, report, err := resolveTokens(filename, source, recipes)
			if err == nil {
				t.Fatal("resolveTokens() error = nil, want nested recipe-use error")
			}
			for _, want := range []string{"recipe token", tt.want} {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("resolveTokens() error %q does not contain %q", err, want)
				}
			}
			if got != nil {
				t.Errorf("resolveTokens() output = %q, want nil on nested recipe misuse", got)
			}
			if len(report.UsedTokens) != 0 {
				t.Errorf("ResolveReport.UsedTokens = %q, want none on nested recipe misuse", report.UsedTokens)
			}
		})
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
	recipes := testRecipes(t,
		recipe.Rule{Class: "gsxui-recipe-button", Utilities: []string{"gsxui-recipe-leaked"}},
	)
	got, _, err := resolveTokens("introduced-token.gsx", src, recipes)
	if err == nil {
		t.Fatal("resolveTokens() error = nil, want structural output-invariant error")
	}
	if !strings.Contains(err.Error(), "gsxui-recipe-leaked") {
		t.Errorf("resolveTokens() error %q does not identify introduced token", err)
	}
	if bytes.Contains(got, []byte(recipe.Prefix)) {
		t.Errorf("resolveTokens() returned output containing %q: %s", recipe.Prefix, got)
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
	got, _, err := resolveTokens("output-invariant.gsx", src, recipe.Style{})
	if err == nil {
		t.Fatal("resolveTokens() error = nil, want independent no-prefix invariant error")
	}
	if bytes.Contains(got, []byte(recipe.Prefix)) {
		t.Errorf("resolveTokens() returned output containing %q: %s", recipe.Prefix, got)
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
			got, report, err := resolveTokens("dynamic.gsx", src, recipe.Style{})
			if err != nil {
				t.Fatalf("resolveTokens() error = %v", err)
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
		recipes recipe.Style
		want    []string
	}{
		{
			name: "missing used token and unused declaration reported together",
			src:  resolveComponent(`<div class={ "gsxui-recipe-missing" }></div>`),
			recipes: testRecipes(t,
				recipe.Rule{Class: "gsxui-recipe-unused", Utilities: []string{"grid"}},
			),
			want: []string{"missing", "gsxui-recipe-missing", "unused", "gsxui-recipe-unused"},
		},
		{
			name: "unused declaration",
			src:  resolveComponent(`<div class={ "gsxui-recipe-button" }></div>`),
			recipes: testRecipes(t,
				recipe.Rule{Class: "gsxui-recipe-button", Utilities: []string{"flex"}},
				recipe.Rule{Class: "gsxui-recipe-unused", Utilities: []string{"grid"}},
			),
			want: []string{"unused", "gsxui-recipe-unused"},
		},
		// There used to be a "duplicate declaration" case here, constructing a
		// recipe.Style with two rules sharing a class to exercise
		// declaredRecipeTokens's own duplicate check in resolve.go. That
		// required reaching past recipe.Style's exported API to build a
		// state recipe.ParseStyle itself now refuses to produce (it already
		// rejects a duplicate class while parsing CSS - see
		// TestParseRecipesRejectsInvalidGrammar/duplicate_recipe_token in
		// internal/recipe/parse_test.go). Since every recipe.Style reaching
		// Resolve in production comes from recipe.ParseStyle, that duplicate
		// path in declaredRecipeTokens is unreachable via the public API and
		// the fixture could not be expressed as parseable CSS, so the case
		// was removed rather than kept alive via a bypass constructor.
		{
			name:    "static class",
			src:     resolveComponent(`<div class="gsxui-recipe-button"></div>`),
			recipes: testRecipes(t, recipe.Rule{Class: "gsxui-recipe-button", Utilities: []string{"flex"}}),
			want:    []string{"static class", "gsxui-recipe-button"},
		},
		{
			name:    "embedded in larger class token",
			src:     resolveComponent(`<div class={ "prefix-gsxui-recipe-button" }></div>`),
			recipes: testRecipes(t, recipe.Rule{Class: "gsxui-recipe-button", Utilities: []string{"flex"}}),
			want:    []string{"whole", "prefix-gsxui-recipe-button"},
		},
		{
			name:    "assembled by concatenation",
			src:     resolveComponent(`<div class={ "gsxui-recipe-" + suffix }></div>`),
			recipes: testRecipes(t, recipe.Rule{Class: "gsxui-recipe-button", Utilities: []string{"flex"}}),
			want:    []string{"literal", "concatenation"},
		},
		{
			name:    "assembled from split constant concatenation",
			src:     resolveComponent(`<div class={ "gsxui-" + "recipe-button" }></div>`),
			recipes: testRecipes(t, recipe.Rule{Class: "gsxui-recipe-button", Utilities: []string{"flex"}}),
			want:    []string{"literal", "concatenation"},
		},
		{
			name:    "assembled by interpolation",
			src:     resolveComponent("<div class={ f`gsxui-recipe-@{suffix}` }></div>"),
			recipes: testRecipes(t, recipe.Rule{Class: "gsxui-recipe-button", Utilities: []string{"flex"}}),
			want:    []string{"literal", "interpolation"},
		},
		{
			name:    "token in non-class expression",
			src:     resolveComponent(`<div title={ "gsxui-recipe-button" } class={ "safe" }></div>`),
			recipes: testRecipes(t, recipe.Rule{Class: "gsxui-recipe-button", Utilities: []string{"flex"}}),
			want:    []string{"non-class", "gsxui-recipe-button"},
		},
		{
			name:    "non-string class expression containing recipe",
			src:     resolveComponent(`<div class={ choose("gsxui-recipe-button") }></div>`),
			recipes: testRecipes(t, recipe.Rule{Class: "gsxui-recipe-button", Utilities: []string{"flex"}}),
			want:    []string{"string literal", "gsxui-recipe-button"},
		},
		{
			name:    "malformed gsx",
			src:     "package p\ncomponent C() { <div>",
			recipes: testRecipes(t, recipe.Rule{Class: "gsxui-recipe-button", Utilities: []string{"flex"}}),
			want:    []string{"malformed.gsx"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := filepath.Join("testdata", "malformed.gsx")
			_, _, err := resolveTokens(filename, []byte(tt.src), tt.recipes)
			if err == nil {
				t.Fatal("resolveTokens() error = nil, want error")
			}
			for _, want := range tt.want {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("resolveTokens() error %q does not contain %q", err, want)
				}
			}
		})
	}
}

func resolveTestRecipes(t *testing.T) recipe.Style {
	t.Helper()
	return testRecipes(t,
		recipe.Rule{Class: "gsxui-recipe-button-size-default", Utilities: []string{"h-9", "px-4"}},
		recipe.Rule{Class: "gsxui-recipe-button", Utilities: []string{"inline-flex", "items-center"}},
		recipe.Rule{Class: "gsxui-recipe-button-variant-outline", Utilities: []string{"border", "border-input", "bg-background"}},
		recipe.Rule{Class: "gsxui-recipe-button-size-icon", Utilities: []string{"size-9", "[&>svg]:shrink-0"}},
		recipe.Rule{Class: "gsxui-recipe-button-variant-default", Utilities: []string{"bg-primary", "text-primary-foreground"}},
	)
}

// testRecipes builds a Style fixture by authoring the rules as real CSS and
// running it through recipe.ParseStyle, so tests exercise the actual parser
// rather than bypass it. Every stylegen test fixture that can be expressed
// as parseable CSS uses this helper; see the "duplicate declaration" case in
// TestResolveRejectsInvalidRecipeUse for the one case that cannot be
// (ParseStyle itself now rejects a duplicate class), and how that case was
// handled instead.
func testRecipes(t *testing.T, rules ...recipe.Rule) recipe.Style {
	t.Helper()

	var css strings.Builder
	css.WriteString("@layer components {\n")
	for _, rule := range rules {
		fmt.Fprintf(&css, ".%s { @apply %s; }\n", rule.Class, strings.Join(rule.Utilities, " "))
	}
	css.WriteString("}\n")

	style, err := recipe.ParseStyle("testRecipes.css", []byte(css.String()))
	if err != nil {
		t.Fatalf("recipe.ParseStyle(fixture) error = %v", err)
	}
	return style
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

	got, _, err := resolveTokens(filename, source, recipe.Style{})
	if err == nil {
		t.Fatal("resolveTokens() error = nil, want structural recipe-use error")
	}
	if bytes.Contains(got, []byte(recipe.Prefix)) {
		t.Errorf("resolveTokens() returned output containing %q: %s", recipe.Prefix, got)
	}
	if !strings.Contains(err.Error(), want) {
		t.Errorf("resolveTokens() error %q does not contain %q", err, want)
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

// testResolved builds the small Button-shaped fixture the helper-call tests
// desugar against: base {inline-flex items-center}, variant default ->
// {bg-primary}, outline -> {border-border bg-background}. It goes through
// recipe.ParseStyle and recipe.Conform so the fixture is a Resolved the real
// conformance check produced, not a hand-built struct.
func testResolved(t *testing.T) recipe.Resolved {
	t.Helper()

	shape := recipe.Shape{
		Component: "button",
		Base:      true,
		Dimensions: []recipe.Dimension{
			{Name: "variant", Default: "default", Values: []string{"default", "outline"}},
		},
	}
	style := testRecipes(t,
		recipe.Rule{Class: "gsxui-recipe-button", Utilities: []string{"inline-flex", "items-center"}},
		recipe.Rule{Class: "gsxui-recipe-button-variant-default", Utilities: []string{"bg-primary"}},
		recipe.Rule{Class: "gsxui-recipe-button-variant-outline", Utilities: []string{"border-border", "bg-background"}},
	)
	resolved, err := recipe.Conform("testResolved.css", shape, style)
	if err != nil {
		t.Fatalf("recipe.Conform(fixture) error = %v", err)
	}
	return resolved
}

func TestResolveDesugarsHelperCalls(t *testing.T) {
	t.Parallel()
	src := []byte(`package canonical

import "github.com/gsxhq/gsx"

component B(variant string, children gsx.Node) {
	<button class={ "group/button", button.Role(), button.Variant(variant) }>
		{ children }
	</button>
}
`)
	got, err := Resolve("button.gsx", src, testResolved(t))
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	for _, want := range []string{
		`"inline-flex items-center"`,
		`case "outline":`,
		`"border-border bg-background"`,
		`default:`,
		`"bg-primary"`,
		`"group/button"`,
	} {
		if !strings.Contains(string(got), want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
	if strings.Contains(string(got), recipe.Prefix) {
		t.Errorf("resolved source still contains the recipe prefix:\n%s", got)
	}
	if strings.Contains(string(got), "button.Role()") {
		t.Errorf("resolved source still contains a helper call:\n%s", got)
	}
	if _, err := gsxparser.ParseFile(token.NewFileSet(), "button.gsx", got, 0); err != nil {
		t.Fatalf("resolved output does not parse: %v", err)
	}
}

func TestResolveDefaultArmCarriesDeclaredDefault(t *testing.T) {
	t.Parallel()
	// The declared default's utilities must appear in the default arm, not "",
	// so a misspelled variant renders the default rather than nothing.
	src := []byte(`package canonical

import "github.com/gsxhq/gsx"

component B(variant string, children gsx.Node) {
	<button class={ button.Variant(variant) }>{ children }</button>
}
`)
	got, err := Resolve("button.gsx", src, testResolved(t))
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if !regexp.MustCompile(`default:\s*\n?\s*"bg-primary"`).Match(got) {
		t.Errorf("default arm must carry the declared default's utilities\nin: %s", got)
	}
	// The default value must not also get its own redundant case arm.
	if strings.Contains(string(got), `case "default":`) {
		t.Errorf("default value must not get a redundant case arm\nin: %s", got)
	}
}

func TestResolveRejectsUnknownDimension(t *testing.T) {
	t.Parallel()
	src := []byte(`package canonical

import "github.com/gsxhq/gsx"

component B(tone string, children gsx.Node) {
	<button class={ button.Tone(tone) }>{ children }</button>
}
`)
	_, err := Resolve("button.gsx", src, testResolved(t))
	if err == nil {
		t.Fatal("Resolve() = nil error, want error")
	}
	if !strings.Contains(err.Error(), "button.gsx:6:") {
		t.Errorf("error must carry a source position, got %q", err)
	}
}

func TestResolveRejectsMalformedHelperCalls(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		body string
		want string
	}{
		{
			name: "wrong component",
			body: `<div class={ badge.Role() }></div>`,
			want: `component "badge"`,
		},
		{
			name: "role with arguments",
			body: `<div class={ button.Role("x") }></div>`,
			want: "takes no arguments",
		},
		{
			name: "dimension without argument",
			body: `<div class={ button.Variant() }></div>`,
			want: "exactly one argument",
		},
		{
			name: "dimension with literal argument",
			body: `<div class={ button.Variant("outline") }></div>`,
			want: "identifier",
		},
		{
			name: "non-identifier receiver",
			body: `<div class={ pkg.button.Role() }></div>`,
			want: "component identifier",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := Resolve("helper.gsx", []byte(resolveComponent(tt.body)), testResolved(t))
			if err == nil {
				t.Fatal("Resolve() error = nil, want error")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Errorf("Resolve() error %q does not contain %q", err, tt.want)
			}
			if !strings.Contains(err.Error(), "helper.gsx:3:") {
				t.Errorf("Resolve() error %q does not carry a source position", err)
			}
		})
	}
}

func TestHelperCallsReportsEveryCallWithoutAShape(t *testing.T) {
	t.Parallel()

	src := []byte(`package canonical

import "github.com/gsxhq/gsx"

component B(variant string, size string, children gsx.Node) {
	<button class={ "group/button", button.Role(), button.Variant(variant), button.Size(size) }>
		{ children }
	</button>
}
`)
	got, err := HelperCalls("button.gsx", src)
	if err != nil {
		t.Fatalf("HelperCalls() error = %v", err)
	}
	want := []Call{
		{Component: "button"},
		{Component: "button", Dimension: "variant"},
		{Component: "button", Dimension: "size"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("HelperCalls() = %+v, want %+v", got, want)
	}
}

func TestResolveDesugarsTheRealCanonicalButton(t *testing.T) {
	t.Parallel()

	root := resolveRepoRoot(t)
	filename := filepath.Join(root, "registry", "canonical", "button.gsx")
	src, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	stylePath := filepath.Join(root, "registry", "styles", "nova", "button.css")
	styleSource, err := os.ReadFile(stylePath)
	if err != nil {
		t.Fatal(err)
	}
	style, err := recipe.ParseStyle(stylePath, styleSource)
	if err != nil {
		t.Fatalf("recipe.ParseStyle() error = %v", err)
	}
	resolved, err := recipe.Conform(stylePath, canonical.Shapes()["button"], style)
	if err != nil {
		t.Fatalf("recipe.Conform() error = %v", err)
	}

	got, err := Resolve(filename, src, resolved)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if bytes.Contains(got, []byte(recipe.Prefix)) {
		t.Errorf("resolved canonical still contains the recipe prefix:\n%s", got)
	}
	for _, helper := range []string{".Role(", ".Variant(", ".Size("} {
		if bytes.Contains(got, []byte(helper)) {
			t.Errorf("resolved canonical still contains %q:\n%s", helper, got)
		}
	}
	// Every declared value's utilities must survive into the output.
	for dimension, values := range resolved.Values {
		for value := range values {
			want := strings.Join(resolved.Utilities(dimension, value), " ")
			if !bytes.Contains(got, []byte(strconv.Quote(want))) {
				t.Errorf("resolved canonical is missing %s=%s utilities %q", dimension, value, want)
			}
		}
	}
	if _, err := gsxparser.ParseFile(token.NewFileSet(), filename, got, 0); err != nil {
		t.Fatalf("resolved canonical does not parse: %v", err)
	}
}

func resolveRepoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("module root not found")
		}
		dir = parent
	}
}
