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

// novaSeparatorRecipe loads the default style's Separator recipe once.
// ui.Separator is generated output now, so the classes it renders are the
// default style's concrete utilities — same pattern as novaButtonRecipe in
// button_test.go.
var novaSeparatorRecipe = sync.OnceValue(func() recipe.Style {
	path := filepath.Join("..", "registry", "styles", stylegen.DefaultStyle, "separator.css")
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

func separatorRecipeUtilities(class string) []string {
	rule, ok := novaSeparatorRecipe().Lookup(class)
	if !ok {
		panic("default style declares no recipe " + class)
	}
	return rule.Utilities
}

// canonicalSeparatorClass is the class attribute ui.Separator renders for one
// orientation, plus any caller classes, merged the same way gsx merges class
// values at runtime.
func canonicalSeparatorClass(orientation string, caller ...string) string {
	classes := append([]string(nil), separatorRecipeUtilities("gsxui-recipe-separator")...)
	classes = append(classes, separatorRecipeUtilities("gsxui-recipe-separator-orientation-"+orientation)...)
	classes = append(classes, caller...)
	return `class="` + html.EscapeString(merge.Merge(classes)) + `"`
}

func TestSeparatorDefault(t *testing.T) {
	got := render(t, ui.Separator("", nil))
	for _, want := range []string{
		`data-gsxui-slot-separator`,
		`role="none"`,
		`data-orientation="horizontal"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestSeparatorVertical(t *testing.T) {
	got := render(t, ui.Separator("vertical", nil))
	if !strings.Contains(got, `data-orientation="vertical"`) {
		t.Errorf("missing vertical data-orientation stamp\nin: %s", got)
	}
	wantClass := canonicalSeparatorClass("vertical")
	if !strings.Contains(got, wantClass) {
		t.Errorf("missing vertical canonical class\nwant: %s\nin: %s", wantClass, got)
	}
}

func TestSeparatorCallerClassIsForwardedOnce(t *testing.T) {
	got := render(t, ui.Separator("", gsx.Attrs{{Key: "class", Value: "bg-red-500"}}))
	wantClass := canonicalSeparatorClass("horizontal", "bg-red-500")
	if strings.Count(got, wantClass) != 1 || strings.Count(got, `class=`) != 1 {
		t.Errorf("caller class must follow exact canonical Separator utilities once\nwant: %s\nin: %s", wantClass, got)
	}
}

func TestSeparatorPinned(t *testing.T) {
	// Exact full-render pin, verified token-by-token against shadcn's
	// Separator (registry/new-york-v4/ui/separator.tsx) and
	// docs/jsx-parity.md — decorative is dropped (ADAPT), role="none" always.
	got := render(t, ui.Separator("", nil))
	want := `<div role="none" data-orientation="horizontal" ` + canonicalSeparatorClass("horizontal") + ` data-gsxui-slot-separator></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestSeparatorAttrsFallThrough(t *testing.T) {
	got := render(t, ui.Separator("", gsx.Attrs{{Key: "id", Value: "s1"}, {Key: "aria-hidden", Value: "true"}}))
	for _, want := range []string{`id="s1"`, `aria-hidden="true"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}
