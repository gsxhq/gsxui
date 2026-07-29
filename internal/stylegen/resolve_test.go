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

	"github.com/gsxhq/gsx/gen"
	gsxparser "github.com/gsxhq/gsx/parser"

	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

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
			got, err := Resolve("element-valued.gsx", src, testResolved(t))
			if err != nil {
				t.Fatalf("Resolve() error = %v", err)
			}
			if !bytes.Equal(got, wantFormatted) {
				t.Errorf("Resolve() changed recipe-free element-valued expression\n--- got ---\n%s\n--- want ---\n%s", got, wantFormatted)
			}
			if !bytes.Contains(got, []byte(tt.want)) {
				t.Errorf("Resolve() output does not preserve %q:\n%s", tt.want, got)
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
			got, err := Resolve(filename, source, testResolved(t))
			if err == nil {
				t.Fatal("Resolve() error = nil, want structural recipe-use error")
			}
			if !strings.Contains(err.Error(), "recipe token") {
				t.Errorf("Resolve() error %q does not identify recipe-token misuse", err)
			}
			if bytes.Contains(got, []byte(recipe.Prefix)) {
				t.Errorf("Resolve() returned output containing %q: %s", recipe.Prefix, got)
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

			got, err := Resolve(filename, source, testResolved(t))
			if err == nil {
				t.Fatal("Resolve() error = nil, want nested recipe-use error")
			}
			for _, want := range []string{"recipe token", tt.want} {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("Resolve() error %q does not contain %q", err, want)
				}
			}
			if got != nil {
				t.Errorf("Resolve() output = %q, want nil on nested recipe misuse", got)
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

	src := []byte(resolveComponent(`<div class={ button.Root() }></div>`))
	shape := recipe.Shape{
		Component: "button",
		Slots: []recipe.Slot{{Name: "", Base: true, Dimensions: []recipe.Dimension{
			{Name: "variant", Default: "default", Values: []string{"default"}},
		}}},
	}
	style := testRecipes(t,
		recipe.Rule{Class: "gsxui-recipe-button", Utilities: []string{"gsxui-recipe-leaked"}},
		recipe.Rule{Class: "gsxui-recipe-button-variant-default", Utilities: []string{"bg-primary"}},
	)
	resolved, err := recipe.Conform("introduced-token.css", shape, style)
	if err != nil {
		t.Fatalf("recipe.Conform() error = %v", err)
	}
	got, err := Resolve("introduced-token.gsx", src, resolved)
	if err == nil {
		t.Fatal("Resolve() error = nil, want structural output-invariant error")
	}
	if !strings.Contains(err.Error(), "gsxui-recipe-leaked") {
		t.Errorf("Resolve() error %q does not identify introduced token", err)
	}
	if bytes.Contains(got, []byte(recipe.Prefix)) {
		t.Errorf("Resolve() returned output containing %q: %s", recipe.Prefix, got)
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
	got, err := Resolve("output-invariant.gsx", src, testResolved(t))
	if err == nil {
		t.Fatal("Resolve() error = nil, want independent no-prefix invariant error")
	}
	if bytes.Contains(got, []byte(recipe.Prefix)) {
		t.Errorf("Resolve() returned output containing %q: %s", recipe.Prefix, got)
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
			got, err := Resolve("dynamic.gsx", src, testResolved(t))
			if err != nil {
				t.Fatalf("Resolve() error = %v", err)
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
		name string
		src  string
		want []string
	}{
		{
			name: "static class",
			src:  resolveComponent(`<div class="gsxui-recipe-button"></div>`),
			want: []string{"static class", "gsxui-recipe-button"},
		},
		{
			name: "embedded in larger class token",
			src:  resolveComponent(`<div class={ "prefix-gsxui-recipe-button" }></div>`),
			want: []string{"whole", "prefix-gsxui-recipe-button"},
		},
		{
			name: "assembled by concatenation",
			src:  resolveComponent(`<div class={ "gsxui-recipe-" + suffix }></div>`),
			want: []string{"literal", "concatenation"},
		},
		{
			name: "assembled from split constant concatenation",
			src:  resolveComponent(`<div class={ "gsxui-" + "recipe-button" }></div>`),
			want: []string{"literal", "concatenation"},
		},
		{
			name: "assembled by interpolation",
			src:  resolveComponent("<div class={ f`gsxui-recipe-@{suffix}` }></div>"),
			want: []string{"literal", "interpolation"},
		},
		{
			name: "token in non-class expression",
			src:  resolveComponent(`<div title={ "gsxui-recipe-button" } class={ "safe" }></div>`),
			want: []string{"non-class", "gsxui-recipe-button"},
		},
		{
			name: "non-string class expression containing recipe",
			src:  resolveComponent(`<div class={ choose("gsxui-recipe-button") }></div>`),
			want: []string{"string literal", "gsxui-recipe-button"},
		},
		{
			name: "malformed gsx",
			src:  "package p\ncomponent C() { <div>",
			want: []string{"malformed.gsx"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := filepath.Join("testdata", "malformed.gsx")
			_, err := Resolve(filename, []byte(tt.src), testResolved(t))
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

	got, err := Resolve(filename, source, testResolved(t))
	if err == nil {
		t.Fatal("Resolve() error = nil, want structural recipe-use error")
	}
	if bytes.Contains(got, []byte(recipe.Prefix)) {
		t.Errorf("Resolve() returned output containing %q: %s", recipe.Prefix, got)
	}
	if !strings.Contains(err.Error(), want) {
		t.Errorf("Resolve() error %q does not contain %q", err, want)
	}
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
		Slots: []recipe.Slot{{Name: "", Base: true, Dimensions: []recipe.Dimension{
			{Name: "variant", Default: "default", Values: []string{"default", "outline"}},
		}}},
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
	<button class={ "group/button", button.Root(), button.Variant(variant) }>
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
	if strings.Contains(string(got), "button.Root()") {
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
			body: `<div class={ badge.Root() }></div>`,
			want: `receiver "badge"`,
		},
		{
			name: "role with arguments",
			body: `<div class={ button.Root("x") }></div>`,
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
			body: `<div class={ pkg.button.Root() }></div>`,
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

func TestHelperCallsReportsEveryCall(t *testing.T) {
	t.Parallel()

	src := []byte(`package canonical

import "github.com/gsxhq/gsx"

component B(variant string, size string, children gsx.Node) {
	<button class={ "group/button", button.Root(), button.Variant(variant), button.Size(size) }>
		{ children }
	</button>
}
`)
	got, err := HelperCalls("button.gsx", src, buttonScanShape())
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

// buttonScanShape is the shape HelperCalls resolves accessor names against in
// the scan-mode tests. It is the real Button shape, so the accessor names in
// the fixtures are the ones the generator actually emits.
func buttonScanShape() recipe.Shape { return shapes.Button }

// assertNoAccessorResidue fails if a generated accessor call survived
// desugaring. The method name is derived from the shape rather than hardcoded.
func assertNoAccessorResidue(t *testing.T, src []byte, component, method string) {
	t.Helper()
	if residue := component + "." + method + "("; bytes.Contains(src, []byte(residue)) {
		t.Errorf("resolved canonical still contains %q:\n%s", residue, src)
	}
}

// TestResolveDesugarsSlotAccessors pins the multi-slot path: a named slot's
// base accessor must resolve to that slot's utilities, not the root's.
func TestResolveDesugarsSlotAccessors(t *testing.T) {
	t.Parallel()
	src := []byte(`package canonical

import "github.com/gsxhq/gsx"

component C(children gsx.Node) {
	<div class={ card.Root() }>
		<div class={ card.Header() }>{ children }</div>
	</div>
}
`)
	resolved := testSlotResolved(t)
	got, err := Resolve("card.gsx", src, resolved)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	for _, want := range []string{`"rounded-xl border"`, `"flex gap-1.5"`} {
		if !strings.Contains(string(got), want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
	for _, forbidden := range []string{recipe.Prefix, "card.Root()", "card.Header()"} {
		if strings.Contains(string(got), forbidden) {
			t.Errorf("resolved source still contains %q:\n%s", forbidden, got)
		}
	}
}

// testSlotResolved builds a two-slot card fixture through the real ParseStyle
// and Conform, so the Resolved is one the conformance check produced.
func testSlotResolved(t *testing.T) recipe.Resolved {
	t.Helper()
	shape := recipe.Shape{
		Component: "card",
		Slots: []recipe.Slot{
			{Name: "", Base: true},
			{Name: "header", Base: true, Dimensions: []recipe.Dimension{
				{Name: "variant", Default: "default", Values: []string{"default", "muted"}},
			}},
		},
	}
	style := testRecipes(t,
		recipe.Rule{Class: "gsxui-recipe-card", Utilities: []string{"rounded-xl", "border"}},
		recipe.Rule{Class: "gsxui-recipe-card-header", Utilities: []string{"flex", "gap-1.5"}},
		recipe.Rule{Class: "gsxui-recipe-card-header-variant-default", Utilities: []string{"text-foreground"}},
		recipe.Rule{Class: "gsxui-recipe-card-header-variant-muted", Utilities: []string{"text-muted-foreground"}},
	)
	resolved, err := recipe.Conform("testSlotResolved.css", shape, style)
	if err != nil {
		t.Fatalf("recipe.Conform(card fixture) error = %v", err)
	}
	return resolved
}

// TestResolveSplitsSlotAndDimensionByAccessorName is the reason accessor names
// are resolved against the shape rather than split. "MenuButtonSize" is slot
// "menu-button" plus dimension "size", and nothing in the string marks the
// boundary — a ToLower or split-on-capitals would silently mis-resolve it.
func TestResolveSplitsSlotAndDimensionByAccessorName(t *testing.T) {
	t.Parallel()
	// A local fixture, not a registry entry: sidebar has no canonical source
	// yet, and test scaffolding does not belong in the shipped shape data. The
	// slot names mirror the real sidebar's prefix collisions.
	shape := recipe.Shape{
		Component: "sidebar",
		Slots: []recipe.Slot{
			{Name: "", Base: true},
			{Name: "menu", Base: true},
			{Name: "menu-button", Base: true, Dimensions: []recipe.Dimension{
				{Name: "size", Default: "default", Values: []string{"default", "lg"}},
			}},
			{Name: "menu-button-tooltip-content", Base: true},
		},
	}
	calls, err := HelperCalls("sidebar.gsx", []byte(`package canonical

import "github.com/gsxhq/gsx"

component S(size string, children gsx.Node) {
	<div class={ sidebar.MenuButtonSize(size) }>{ children }</div>
}
`), shape)
	if err != nil {
		t.Fatalf("HelperCalls() error = %v", err)
	}
	if len(calls) != 1 {
		t.Fatalf("HelperCalls() = %d calls, want 1", len(calls))
	}
	if got, want := calls[0].Slot, "menu-button"; got != want {
		t.Errorf("Slot = %q, want %q", got, want)
	}
	if got, want := calls[0].Dimension, "size"; got != want {
		t.Errorf("Dimension = %q, want %q", got, want)
	}
}

func TestResolveDesugarsTheRealCanonicalButton(t *testing.T) {
	t.Parallel()

	root := repoRoot(t)
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
	resolved, err := recipe.Conform(stylePath, shapes.Button, style)
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
	// Derived from the shape, so a new slot or dimension is covered without a
	// second edit here.
	for _, slot := range resolved.Shape.Slots {
		if slot.Base {
			assertNoAccessorResidue(t, got, resolved.Shape.Component, accessorName(slot.Name))
		}
		for _, dimension := range slot.Dimensions {
			assertNoAccessorResidue(t, got, resolved.Shape.Component,
				dimensionAccessorName(slot.Name, dimension.Name))
		}
	}
	// Every declared value's utilities must survive into the output.
	for slot, dimensions := range resolved.Values {
		for dimension, values := range dimensions {
			for value := range values {
				want := strings.Join(resolved.Utilities(slot, dimension, value), " ")
				if !bytes.Contains(got, []byte(strconv.Quote(want))) {
					t.Errorf("resolved canonical is missing %s/%s=%s utilities %q", slot, dimension, value, want)
				}
			}
		}
	}
	if _, err := gsxparser.ParseFile(token.NewFileSet(), filename, got, 0); err != nil {
		t.Fatalf("resolved canonical does not parse: %v", err)
	}
}

func repoRoot(t *testing.T) string {
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

// TestResolveRejectsHelperCallInNonPlainClassPart pins the part-shape guard.
// The replacement span runs to part.End(), so without this guard a helper call
// carrying a condition, a pipeline or control flow would have its trailing
// metadata silently deleted rather than rejected.
func TestResolveRejectsHelperCallInNonPlainClassPart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		body string
		want string
		at   string
	}{
		{
			name: "condition",
			body: `<div class={ button.Root() : disabled }></div>`,
			want: "plain class part",
			at:   "part.gsx:3:",
		},
		{
			name: "pipeline",
			body: `<div class={ button.Root() |> upper() }></div>`,
			want: "plain class part",
			at:   "part.gsx:3:",
		},
		{
			name: "dimension call with condition",
			body: `<div class={ button.Variant(variant) : disabled }></div>`,
			want: "plain class part",
			at:   "part.gsx:3:",
		},
		{
			name: "value switch arm",
			body: `<div class={
				switch mode {
				case "a":
					button.Root()
				default:
					"other"
				},
			}></div>`,
			want: "whole class part",
			at:   "part.gsx:6:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			source := []byte(resolveComponent(tt.body))
			if _, err := gsxparser.ParseFile(token.NewFileSet(), "part.gsx", source, 0); err != nil {
				t.Fatalf("fixture does not parse: %v", err)
			}

			got, err := Resolve("part.gsx", source, testResolved(t))
			if err == nil {
				t.Fatalf("Resolve() error = nil, want error; output:\n%s", got)
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Errorf("Resolve() error %q does not contain %q", err, tt.want)
			}
			if !strings.Contains(err.Error(), tt.at) {
				t.Errorf("Resolve() error %q does not carry position %s", err, tt.at)
			}
			if got != nil {
				t.Errorf("Resolve() output = %q, want nil on a rejected part", got)
			}

			// The same guard must run in scan mode, so Task 7's pre-flight
			// check cannot pass a call that Resolve would later refuse.
			calls, err := HelperCalls("part.gsx", source, buttonScanShape())
			if err == nil {
				t.Fatalf("HelperCalls() error = nil, want error; calls = %+v", calls)
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Errorf("HelperCalls() error %q does not contain %q", err, tt.want)
			}
			if !strings.Contains(err.Error(), tt.at) {
				t.Errorf("HelperCalls() error %q does not carry position %s", err, tt.at)
			}
		})
	}
}

// TestResolveRejectsRoleAgainstABaselessShape pins the base-rule guard.
// Shape.Validate does not require Base, so without this check a Base:false
// shape would resolve button.Root() to "" and emit a silently unstyled element.
func TestResolveRejectsRoleAgainstABaselessShape(t *testing.T) {
	t.Parallel()

	shape := recipe.Shape{
		Component: "button",
		Slots: []recipe.Slot{{Name: "", Base: false, Dimensions: []recipe.Dimension{
			{Name: "variant", Default: "default", Values: []string{"default", "outline"}},
		}}},
	}
	style := testRecipes(t,
		recipe.Rule{Class: "gsxui-recipe-button-variant-default", Utilities: []string{"bg-primary"}},
		recipe.Rule{Class: "gsxui-recipe-button-variant-outline", Utilities: []string{"border-border"}},
	)
	resolved, err := recipe.Conform("baseless.css", shape, style)
	if err != nil {
		t.Fatalf("recipe.Conform(baseless) error = %v", err)
	}
	if root, _ := resolved.Shape.Slot(""); root.Base {
		t.Fatal("fixture shape must declare no base rule")
	}

	source := []byte(resolveComponent(`<div class={ button.Root() }></div>`))
	got, err := Resolve("baseless.gsx", source, resolved)
	if err == nil {
		t.Fatalf("Resolve() error = nil, want error; output:\n%s", got)
	}
	if !strings.Contains(err.Error(), "declares no accessor button.Root") {
		t.Errorf("Resolve() error %q does not identify the missing base accessor", err)
	}
	if !strings.Contains(err.Error(), "baseless.gsx:3:") {
		t.Errorf("Resolve() error %q does not carry a source position", err)
	}
	if got != nil {
		t.Errorf("Resolve() output = %q, want nil", got)
	}

	// A baseless shape is still usable for its dimensions.
	dimensionOnly := []byte(resolveComponent(`<div class={ button.Variant(variant) }></div>`))
	resolvedSource, err := Resolve("baseless.gsx", dimensionOnly, resolved)
	if err != nil {
		t.Fatalf("Resolve(dimension only) error = %v", err)
	}
	if !strings.Contains(string(resolvedSource), `"bg-primary"`) {
		t.Errorf("Resolve(dimension only) did not desugar the dimension:\n%s", resolvedSource)
	}
}
