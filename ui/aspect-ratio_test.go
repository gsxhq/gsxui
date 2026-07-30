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
	// gsx renders the class attribute before style regardless of source
	// attribute order (class merging happens ahead of other attrs).
	want := `<div ` + canonicalAspectRatioClass() + ` style="aspect-ratio: 16 / 9" data-gsxui-slot-aspect-ratio><img src="x.png"/></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}
