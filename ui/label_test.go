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

// novaLabelRecipe loads the default style's Label recipe once. ui.Label is
// generated output now, so the classes it renders are the default style's
// concrete utilities; the expectation is derived from the same recipe the
// generator reads rather than restated as a literal that would have to be
// rewritten on every style edit — same pattern as novaButtonRecipe in
// button_test.go.
var novaLabelRecipe = sync.OnceValue(func() recipe.Style {
	path := filepath.Join("..", "registry", "styles", stylegen.DefaultStyle, "label.css")
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

func labelRecipeUtilities(class string) []string {
	rule, ok := novaLabelRecipe().Lookup(class)
	if !ok {
		panic("default style declares no recipe " + class)
	}
	return rule.Utilities
}

// canonicalLabelClass is the class attribute ui.Label renders, plus any
// caller classes, merged the same way gsx merges class values at runtime.
func canonicalLabelClass(caller ...string) string {
	classes := append([]string(nil), labelRecipeUtilities("gsxui-recipe-label")...)
	classes = append(classes, caller...)
	return `class="` + html.EscapeString(merge.Merge(classes)) + `"`
}

func TestLabelDefault(t *testing.T) {
	got := render(t, ui.Label(gsx.Raw("Email"), nil))
	for _, want := range []string{
		"<label", `data-gsxui-slot-label`,
		">Email</label>",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestLabelPinned(t *testing.T) {
	// Exact full-render pin, verified token-by-token against shadcn's Label
	// (registry/new-york-v4/ui/label.tsx) and docs/jsx-parity.md — straight
	// port, no ADAPT deviations.
	got := render(t, ui.Label(gsx.Raw("Email"), nil))
	want := `<label ` + canonicalLabelClass() + ` data-gsxui-slot-label>Email</label>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestLabelCallerClassIsForwardedOnce(t *testing.T) {
	got := render(t, ui.Label(gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "text-lg"}}))
	wantClass := canonicalLabelClass("text-lg")
	if strings.Count(got, wantClass) != 1 || strings.Count(got, `class=`) != 1 {
		t.Errorf("caller class must follow exact canonical Label utilities once\nwant: %s\nin: %s", wantClass, got)
	}
}

func TestLabelAttrsFallThrough(t *testing.T) {
	got := render(t, ui.Label(gsx.Raw("x"), gsx.Attrs{{Key: "for", Value: "email"}, {Key: "id", Value: "email-label"}}))
	for _, want := range []string{`for="email"`, `id="email-label"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}
