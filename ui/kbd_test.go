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

// novaKbdRecipe loads the default style's Kbd recipe once. ui.Kbd and
// ui.KbdGroup are generated output now, so the classes they render are the
// default style's concrete utilities — same pattern as novaButtonRecipe in
// button_test.go.
var novaKbdRecipe = sync.OnceValue(func() recipe.Style {
	path := filepath.Join("..", "registry", "styles", stylegen.DefaultStyle, "kbd.css")
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

func kbdRecipeUtilities(class string) []string {
	rule, ok := novaKbdRecipe().Lookup(class)
	if !ok {
		panic("default style declares no recipe " + class)
	}
	return rule.Utilities
}

func canonicalKbdClass(caller ...string) string {
	classes := append([]string(nil), kbdRecipeUtilities("gsxui-recipe-kbd")...)
	classes = append(classes, caller...)
	return `class="` + html.EscapeString(merge.Merge(classes)) + `"`
}

func canonicalKbdGroupClass(caller ...string) string {
	classes := append([]string(nil), kbdRecipeUtilities("gsxui-recipe-kbd-group")...)
	classes = append(classes, caller...)
	return `class="` + html.EscapeString(merge.Merge(classes)) + `"`
}

func TestKbdDefault(t *testing.T) {
	got := render(t, ui.Kbd(gsx.Raw("Ctrl"), nil))
	for _, want := range []string{
		`<kbd `, `data-gsxui-slot-kbd`,
		">Ctrl</kbd>",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestKbdCallerClassIsForwardedOnce(t *testing.T) {
	got := render(t, ui.Kbd(nil, gsx.Attrs{{Key: "class", Value: "h-8"}}))
	wantClass := canonicalKbdClass("h-8")
	if strings.Count(got, wantClass) != 1 || strings.Count(got, `class=`) != 1 {
		t.Errorf("caller class must follow exact canonical Kbd utilities once\nwant: %s\nin: %s", wantClass, got)
	}
}

func TestKbdAttrsFallThrough(t *testing.T) {
	got := render(t, ui.Kbd(nil, gsx.Attrs{{Key: "id", Value: "k1"}, {Key: "aria-label", Value: "control"}}))
	for _, want := range []string{`id="k1"`, `aria-label="control"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestKbdPinned(t *testing.T) {
	// Exact full-render pin. Token-for-token against shadcn's Kbd
	// (registry/new-york-v4/ui/kbd.tsx) — no dropped tokens.
	got := render(t, ui.Kbd(gsx.Raw("Ctrl"), nil))
	want := `<kbd ` + canonicalKbdClass() + ` data-gsxui-slot-kbd>Ctrl</kbd>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestKbdGroupDefault(t *testing.T) {
	got := render(t, ui.KbdGroup(gsx.Raw("<kbd>Ctrl</kbd><kbd>C</kbd>"), nil))
	want := `<kbd ` + canonicalKbdGroupClass() + ` data-gsxui-slot-kbd-group><kbd>Ctrl</kbd><kbd>C</kbd></kbd>`
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func TestKbdGroupAttrsFallThrough(t *testing.T) {
	got := render(t, ui.KbdGroup(nil, gsx.Attrs{{Key: "id", Value: "g1"}}))
	if !strings.Contains(got, `id="g1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

func TestKbdGroupPinned(t *testing.T) {
	got := render(t, ui.KbdGroup(gsx.Raw("<kbd>Ctrl</kbd>"), nil))
	want := `<kbd ` + canonicalKbdGroupClass() + ` data-gsxui-slot-kbd-group><kbd>Ctrl</kbd></kbd>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}
