// This file pins the *generated* Nova output of ui.Button: the concrete
// utilities the default style's recipe compiles to, plus the gsx-level class
// merging and attribute precedence that only the generated component exhibits.
//
// Style-independent Button behavior — href renders an anchor, disabled renders
// a real button even with an href, caller attributes fall through — lives in
// registry/canonical/button_test.go and is deliberately NOT restated here.
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

func TestButtonPinned(t *testing.T) {
	// Exact full-render pin, verified token-by-token against shadcn's
	// buttonVariants base + default variant + default size
	// (registry/new-york-v4/ui/button.tsx) and docs/jsx-parity.md. The
	// *default* variant/size pinned here carries no ADAPT deviation; the
	// destructive variant does (nova's dark:hover:bg-destructive/90 — see
	// docs/jsx-parity.md's `## button` entry).
	got := render(t, ui.Button("", "", "", false, gsx.Raw("Save"), nil))
	want := `<button data-variant="default" data-size="default" type="button" ` +
		canonicalButtonClass("default", "default") +
		` data-gsxui-slot-button>Save</button>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestButtonVariantAndSizeAxes(t *testing.T) {
	for _, variant := range []string{"default", "destructive", "outline", "secondary", "ghost", "link"} {
		input := variant
		if variant == "default" {
			input = ""
		}
		got := render(t, ui.Button(input, "", "", false, gsx.Raw("x"), nil))
		if !strings.Contains(got, `data-variant="`+variant+`"`) {
			t.Errorf("variant %s: missing reflected value\nin: %s", variant, got)
		}
		if !strings.Contains(got, canonicalButtonClass(variant, "default")) {
			t.Errorf("variant %s: missing exact canonical role classes\nin: %s", variant, got)
		}
	}
	for _, size := range []string{"default", "xs", "sm", "lg", "icon", "icon-xs", "icon-sm", "icon-lg"} {
		input := size
		if size == "default" {
			input = ""
		}
		got := render(t, ui.Button("", input, "", false, gsx.Raw("x"), nil))
		if !strings.Contains(got, `data-size="`+size+`"`) {
			t.Errorf("size %s: missing reflected value\nin: %s", size, got)
		}
		if !strings.Contains(got, canonicalButtonClass("default", size)) {
			t.Errorf("size %s: missing exact canonical role classes\nin: %s", size, got)
		}
	}
}

func TestButtonTypeIsOverridableDefault(t *testing.T) {
	// type="button" is authored BEFORE { attrs... }: caller type=submit wins.
	got := render(t, ui.Button("", "", "", false, gsx.Raw("Go"), gsx.Attrs{{Key: "type", Value: "submit"}}))
	if !strings.Contains(got, `type="submit"`) {
		t.Errorf("caller type=submit must override default\nin: %s", got)
	}
	if strings.Contains(got, `type="button"`) {
		t.Errorf("default type should be replaced, not duplicated\nin: %s", got)
	}
}

func TestButtonCallerClassFollowsCanonicalRolesOnce(t *testing.T) {
	got := render(t, ui.Button("", "", "", false, gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "h-12"}}))
	want := canonicalButtonClass("default", "default", "h-12")
	if strings.Count(got, want) != 1 || strings.Count(got, `class=`) != 1 {
		t.Errorf("caller class must follow exact canonical roles and render once\nwant: %s\nin: %s", want, got)
	}
}

func TestButtonOwnPresenceMarkerWinsCollisionAndKeepsComposedMarker(t *testing.T) {
	got := render(t, ui.Button("", "", "", false, gsx.Raw("x"), gsx.Attrs{
		{Key: "data-gsxui-slot-button", Value: "caller-value"},
		{Key: "data-gsxui-slot-pagination-link", Value: gsx.Toggle(true)},
	}))
	requirePresenceAttributesOnSameTag(t, got, "<button",
		"data-gsxui-slot-pagination-link",
		"data-gsxui-slot-button",
	)
}

// novaButtonRecipe loads the default style's Button recipe once. ui.Button is
// generated output now, so the classes it renders are the default style's
// concrete utilities; the expectation is derived from the same recipe the
// generator reads rather than restated as a literal that would have to be
// rewritten on every style edit.
var novaButtonRecipe = sync.OnceValue(func() recipe.Style {
	path := filepath.Join("..", "registry", "styles", stylegen.DefaultStyle, "button.css")
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

func recipeUtilities(class string) []string {
	rule, ok := novaButtonRecipe().Lookup(class)
	if !ok {
		panic("default style declares no recipe " + class)
	}
	return rule.Utilities
}

// canonicalButtonClass is the class attribute ui.Button renders for one
// variant/size pair, plus any caller classes. gsx merges the whole class value
// through merge.Merge, so the expectation runs the same merger: a caller
// utility that conflicts with a role utility replaces it rather than following
// it.
func canonicalButtonClass(variant, size string, caller ...string) string {
	classes := []string{"group/button"}
	classes = append(classes, recipeUtilities("gsxui-recipe-button")...)
	classes = append(classes, recipeUtilities("gsxui-recipe-button-variant-"+variant)...)
	classes = append(classes, recipeUtilities("gsxui-recipe-button-size-"+size)...)
	classes = append(classes, caller...)
	return `class="` + html.EscapeString(merge.Merge(classes)) + `"`
}

func assertButtonCallerAttrsOnce(t *testing.T, got, variant, size string) {
	t.Helper()

	wantClass := canonicalButtonClass(variant, size, "caller-only")
	if strings.Count(got, wantClass) != 1 || strings.Count(got, `class=`) != 1 {
		t.Errorf("caller class must follow exact canonical Button roles once\nwant: %s\nin: %s", wantClass, got)
	}
	if strings.Count(got, `id="caller-id"`) != 1 {
		t.Errorf("caller id must render exactly once\nin: %s", got)
	}
}
