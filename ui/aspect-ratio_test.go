package ui_test

import (
	"html"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/internal/stylegen"
	"github.com/gsxhq/gsxui/merge"
	"github.com/gsxhq/gsxui/ui"
)

// novaAspectRatioRecipe loads the default style's AspectRatio recipe once.
// ui.AspectRatio is generated output now, so the classes it renders are the
// default style's concrete utilities — same pattern as novaButtonRecipe in
// button_test.go. The recipe file itself is named "aspect-ratio.css" (the
// registry/CSS identity stays kebab-case even though the Go identifier
// derived from it, aspectRatioRecipe, is camel).
var novaAspectRatioRecipe = sync.OnceValue(func() recipe.Style {
	path := filepath.Join("..", "registry", "styles", stylegen.DefaultStyle, "aspect-ratio.css")
	src, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	style, err := recipe.ParseStyle(path, src)
	if err != nil {
		panic(err)
	}
	return style
})

func aspectRatioRecipeUtilities(class string) []string {
	rule, ok := novaAspectRatioRecipe().Lookup(class)
	if !ok {
		panic("default style declares no recipe " + class)
	}
	return rule.Utilities
}

// canonicalAspectRatioClass is the class attribute ui.AspectRatio renders,
// plus any caller classes, merged the same way gsx merges class values at
// runtime.
func canonicalAspectRatioClass(caller ...string) string {
	classes := append([]string(nil), aspectRatioRecipeUtilities("gsxui-recipe-aspect-ratio")...)
	classes = append(classes, caller...)
	return `class="` + html.EscapeString(merge.Merge(classes)) + `"`
}

func TestAspectRatioDefault(t *testing.T) {
	got := render(t, ui.AspectRatio("16 / 9", gsx.Raw(`<img src="x.png"/>`), nil))
	for _, want := range []string{
		`data-gsxui-slot-aspect-ratio`,
		`style="aspect-ratio: 16 / 9"`,
		`<img src="x.png"/>`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestAspectRatioNumericRatio(t *testing.T) {
	// aspect-ratio also accepts a bare number, not only the "w / h" form.
	got := render(t, ui.AspectRatio("1.5", nil, nil))
	if !strings.Contains(got, `style="aspect-ratio: 1.5"`) {
		t.Errorf("missing numeric ratio style\nin: %s", got)
	}
}

// TestAspectRatioGrammarForms pins every shape of the CSS aspect-ratio
// grammar the component parses. The "<number> / <number>" and "auto || …"
// forms are split so the slash and the auto keyword are static template
// text; a ratio the grammar does not cover is emitted as a single hole and
// left to gsx's CSS value filter.
func TestAspectRatioGrammarForms(t *testing.T) {
	for _, tt := range []struct {
		name  string
		ratio string
		want  string
	}{
		{"pair", "16 / 9", "16 / 9"},
		{"pair without spaces", "16/9", "16 / 9"},
		{"pair with excess spaces", "  16  /  9  ", "16 / 9"},
		{"bare number", "1.5", "1.5"},
		{"exponent number", "1e2 / 3", "1e2 / 3"},
		{"auto alone", "auto", "auto"},
		{"auto before ratio", "auto 16 / 9", "auto 16 / 9"},
		// `auto || <ratio>` is order-free; both orders mean the same thing,
		// so the trailing form is re-emitted in the canonical leading one.
		{"auto after ratio", "16 / 9 auto", "auto 16 / 9"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := render(t, ui.AspectRatio(tt.ratio, nil, nil))
			if want := `style="aspect-ratio: ` + tt.want + `"`; !strings.Contains(got, want) {
				t.Errorf("missing %q\nin: %s", want, got)
			}
		})
	}
}

// TestAspectRatioHostileRatioIsInert is the reason the component parses the
// grammar instead of wrapping the value in gsx.RawCSS. Every author-supplied
// number lands in its own hole, so gsx's CSS value filter (a port of
// html/template's) still judges it: nothing that is not a recognised number
// can reach the rendered declaration, and a hole can never contribute the
// "/" that would be needed to close a CSS comment or open a second
// declaration.
func TestAspectRatioHostileRatioIsInert(t *testing.T) {
	for _, ratio := range []string{
		"red;background:url(javascript:alert(1))",
		"16 / 9; color:red",
		// A "*" immediately after the static slash would open a CSS comment
		// and swallow the rest of the style attribute if the sides were not
		// checked for being numbers.
		"1 /* x",
		"1 / * x",
		"var(--r)",
		"expression(alert(1))",
		"auto16/9",
	} {
		t.Run(ratio, func(t *testing.T) {
			got := render(t, ui.AspectRatio(ratio, nil, nil))
			if !strings.Contains(got, `style="aspect-ratio: ZgotmplZ"`) {
				t.Errorf("hostile ratio was not neutralised\nin: %s", got)
			}
		})
	}
}

func TestAspectRatioAttrsFallThrough(t *testing.T) {
	got := render(t, ui.AspectRatio("16 / 9", nil, gsx.Attrs{{Key: "id", Value: "ar1"}, {Key: "class", Value: "bg-muted"}}))
	if !strings.Contains(got, `id="ar1"`) {
		t.Errorf("missing fallthrough id\nin: %s", got)
	}
	wantClass := canonicalAspectRatioClass("bg-muted")
	if !strings.Contains(got, wantClass) {
		t.Errorf("missing merged class\nwant: %s\nin: %s", wantClass, got)
	}
}

func TestAspectRatioPinned(t *testing.T) {
	// Exact full-render pin. shadcn's AspectRatio is a bare passthrough onto
	// Radix's padding-hack Root (registry/new-york-v4/ui/aspect-ratio.tsx);
	// this port replaces the two-div padding-percentage mechanism with the
	// CSS aspect-ratio property directly on a single div (ADAPT, see
	// docs/jsx-parity.md).
	got := render(t, ui.AspectRatio("16 / 9", gsx.Raw(`<img src="x.png"/>`), nil))
	// The style group is written in source order, ahead of class: a
	// { switch } in attribute position emits where it is written, unlike the
	// { if } chain it replaced, which gsx hoisted class ahead of. The
	// declaration carries no trailing semicolon — a part's trailing ";" is
	// trimmed when style declarations are joined, so writing one in the css
	// literal (which tree-sitter needs to parse it as CSS) costs nothing here.
	want := `<div style="aspect-ratio: 16 / 9" ` + canonicalAspectRatioClass() + ` data-gsxui-slot-aspect-ratio><img src="x.png"/></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}
