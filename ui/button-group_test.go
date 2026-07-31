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

// novaButtonGroupRecipe loads the default style's ButtonGroup recipe once —
// same pattern as novaSeparatorRecipe in separator_test.go.
var novaButtonGroupRecipe = sync.OnceValue(func() recipe.Style {
	path := filepath.Join("..", "registry", "styles", stylegen.DefaultStyle, "button-group.css")
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

func buttonGroupRecipeUtilities(class string) []string {
	rule, ok := novaButtonGroupRecipe().Lookup(class)
	if !ok {
		panic("default style declares no recipe " + class)
	}
	return rule.Utilities
}

// canonicalButtonGroupClass is the class attribute ui.ButtonGroup renders for
// one orientation, plus any caller classes, merged the way gsx merges class
// values at runtime.
func canonicalButtonGroupClass(orientation string, caller ...string) string {
	classes := append([]string(nil), buttonGroupRecipeUtilities("gsxui-recipe-button-group")...)
	classes = append(classes, buttonGroupRecipeUtilities("gsxui-recipe-button-group-orientation-"+orientation)...)
	classes = append(classes, caller...)
	return `class="` + html.EscapeString(merge.Merge(classes)) + `"`
}

func canonicalButtonGroupTextClass(caller ...string) string {
	classes := append([]string(nil), buttonGroupRecipeUtilities("gsxui-recipe-button-group-text")...)
	classes = append(classes, caller...)
	return `class="` + html.EscapeString(merge.Merge(classes)) + `"`
}

// canonicalButtonGroupSeparatorClass is what ButtonGroupSeparator renders:
// ui.Separator's own recipe utilities for the orientation, with ButtonGroup's
// separator-slot utilities merged in AFTER them as an ordinary caller class —
// which is how bg-input now displaces Separator's bg-border (tailwind-merge
// dedupes the pair outright) instead of having to outrank it from an
// unwrapped @layer utilities rule, as it did before ButtonGroup migrated.
func canonicalButtonGroupSeparatorClass(orientation string, caller ...string) string {
	classes := append([]string(nil), separatorRecipeUtilities("gsxui-recipe-separator")...)
	classes = append(classes, separatorRecipeUtilities("gsxui-recipe-separator-orientation-"+orientation)...)
	classes = append(classes, buttonGroupRecipeUtilities("gsxui-recipe-button-group-separator")...)
	classes = append(classes, caller...)
	return `class="` + html.EscapeString(merge.Merge(classes)) + `"`
}

// TestButtonGroupDefaultPinned pins the zero-value (horizontal) orientation:
// data-orientation defaults via the house |> default pattern, and the class
// list picks the "horizontal" cva block (rounded-l-none/border-l-0/
// rounded-r-none on non-edge children) via a switch value-form, the same
// idiom badge.gsx uses for its own variant map.
func TestButtonGroupDefaultPinned(t *testing.T) {
	got := render(t, ui.ButtonGroup("", gsx.Raw("x"), nil))
	want := `<div role="group" data-orientation="horizontal" ` + canonicalButtonGroupClass("horizontal") + ` data-gsxui-slot-button-group>x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestButtonGroupVerticalPinned(t *testing.T) {
	got := render(t, ui.ButtonGroup("vertical", gsx.Raw("x"), nil))
	want := `<div role="group" data-orientation="vertical" ` + canonicalButtonGroupClass("vertical") + ` data-gsxui-slot-button-group>x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestButtonGroupAttrsFallThrough(t *testing.T) {
	got := render(t, ui.ButtonGroup("", nil, gsx.Attrs{{Key: "id", Value: "bg1"}}))
	if !strings.Contains(got, `id="bg1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

func TestButtonGroupCallerClassIsForwardedOnce(t *testing.T) {
	got := render(t, ui.ButtonGroup("", nil, gsx.Attrs{{Key: "class", Value: "gap-2"}}))
	if strings.Count(got, "gap-2") != 1 || strings.Count(got, `class=`) != 1 {
		t.Errorf("caller class must merge in exactly once\nin: %s", got)
	}
	if want := canonicalButtonGroupClass("horizontal", "gap-2"); !strings.Contains(got, want) {
		t.Errorf("caller class not merged after the recipe's own\nwant: %s\nin: %s", want, got)
	}
}

func TestButtonGroupTextPinned(t *testing.T) {
	got := render(t, ui.ButtonGroupText(gsx.Raw("x"), nil))
	want := `<div ` + canonicalButtonGroupTextClass() + ` data-gsxui-slot-button-group-text>x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestButtonGroupTextAttrsFallThrough(t *testing.T) {
	got := render(t, ui.ButtonGroupText(nil, gsx.Attrs{{Key: "id", Value: "t1"}}))
	if !strings.Contains(got, `id="t1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

// TestButtonGroupSeparatorDefaultPinned proves the "vertical" default (the
// opposite of Separator's own "horizontal" default). Separator is migrated
// onto the slot axis now, so its own utilities render as a real class here;
// ButtonGroup is migrated too now, so its separator-slot utilities arrive as
// an ordinary caller class on <Separator>: tailwind-merge drops Separator's
// own bg-border in favour of bg-input outright, and data-[orientation=vertical]:h-auto
// rides along to outrank h-full at render time. The unwrapped
// @layer utilities promotion this pin used to describe is retired.
func TestButtonGroupSeparatorDefaultPinned(t *testing.T) {
	got := render(t, ui.ButtonGroupSeparator("", nil))
	want := `<div role="none" data-orientation="vertical" ` + canonicalButtonGroupSeparatorClass("vertical") + ` data-gsxui-slot-button-group-separator data-gsxui-slot-separator></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestButtonGroupSeparatorOrientationOverride(t *testing.T) {
	got := render(t, ui.ButtonGroupSeparator("horizontal", nil))
	if !strings.Contains(got, `data-orientation="horizontal"`) {
		t.Errorf("missing data-orientation=horizontal override\nin: %s", got)
	}
}

func TestButtonGroupSeparatorAttrsFallThrough(t *testing.T) {
	got := render(t, ui.ButtonGroupSeparator("", gsx.Attrs{{Key: "id", Value: "sep1"}}))
	if !strings.Contains(got, `id="sep1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

// TestButtonGroupSeparatorComposesSeparator proves ButtonGroupSeparator
// actually renders through ui.Separator (role="none" + the base
// data-[orientation=...] selectors both come from Separator itself), the
// button-group -> separator dependency internal/registry derives.
func TestButtonGroupSeparatorComposesSeparator(t *testing.T) {
	got := render(t, ui.ButtonGroupSeparator("", nil))
	for _, want := range []string{`role="none"`, `data-gsxui-slot-button-group-separator data-gsxui-slot-separator`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q (expected from ui.Separator)\nin: %s", want, got)
		}
	}
}

// Realistic composition: two buttons split by a separator, the
// button-group-separator demo shape.
func TestButtonGroupWithSeparator(t *testing.T) {
	got := render(t, ui.ButtonGroup("",
		gsx.Fragment(
			ui.Button("secondary", "sm", "", false, gsx.Raw("Copy"), nil),
			ui.ButtonGroupSeparator("", nil),
			ui.Button("secondary", "sm", "", false, gsx.Raw("Paste"), nil),
		),
		nil,
	))
	for _, want := range []string{
		`data-gsxui-slot-button-group`,
		`>Copy</button>`,
		`data-gsxui-slot-button-group-separator data-gsxui-slot-separator`,
		`>Paste</button>`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}
