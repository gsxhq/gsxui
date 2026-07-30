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

// novaItemRecipe loads the default style's Item recipe once — same pattern as
// novaSeparatorRecipe in separator_test.go.
var novaItemRecipe = sync.OnceValue(func() recipe.Style {
	path := filepath.Join("..", "registry", "styles", stylegen.DefaultStyle, "item.css")
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

func itemRecipeUtilities(class string) []string {
	rule, ok := novaItemRecipe().Lookup(class)
	if !ok {
		panic("default style declares no recipe " + class)
	}
	return rule.Utilities
}

// canonicalItemClass builds the class attribute one Item slot renders: the
// slot's base utilities, then any dimension-value utilities, then the caller's
// own classes, merged the way gsx merges class values at runtime.
func canonicalItemClass(slot string, extra []string, caller ...string) string {
	base := "gsxui-recipe-item"
	if slot != "" {
		base += "-" + slot
	}
	classes := append([]string(nil), itemRecipeUtilities(base)...)
	for _, value := range extra {
		classes = append(classes, itemRecipeUtilities(base+"-"+value)...)
	}
	classes = append(classes, caller...)
	return `class="` + html.EscapeString(merge.Merge(classes)) + `"`
}

// canonicalItemSeparatorClass is what ItemSeparator renders: ui.Separator's
// own recipe utilities for the orientation, with Item's separator-slot my-2
// merged in after them as an ordinary caller class.
func canonicalItemSeparatorClass(orientation string, caller ...string) string {
	classes := append([]string(nil), separatorRecipeUtilities("gsxui-recipe-separator")...)
	classes = append(classes, separatorRecipeUtilities("gsxui-recipe-separator-orientation-"+orientation)...)
	classes = append(classes, itemRecipeUtilities("gsxui-recipe-item-separator")...)
	classes = append(classes, caller...)
	return `class="` + html.EscapeString(merge.Merge(classes)) + `"`
}

func TestItemGroupPinned(t *testing.T) {
	got := render(t, ui.ItemGroup(gsx.Raw("x"), nil))
	want := `<div ` + canonicalItemClass("group", nil) + ` role="list" data-gsxui-slot-item-group>x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestItemGroupAttrsFallThrough(t *testing.T) {
	got := render(t, ui.ItemGroup(nil, gsx.Attrs{{Key: "id", Value: "ig1"}}))
	if !strings.Contains(got, `id="ig1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

func TestItemGroupCallerClassIsForwardedOnce(t *testing.T) {
	got := render(t, ui.ItemGroup(nil, gsx.Attrs{{Key: "class", Value: "gap-4"}}))
	if strings.Count(got, "gap-4") != 1 || strings.Count(got, `class=`) != 1 {
		t.Errorf("caller class must merge in exactly once\nin: %s", got)
	}
	if want := canonicalItemClass("group", nil, "gap-4"); !strings.Contains(got, want) {
		t.Errorf("caller class not merged after the recipe's own\nwant: %s\nin: %s", want, got)
	}
}

// TestItemSeparatorDefaultPinned pins the zero-value ("horizontal")
// orientation and proves ItemSeparator actually composes ui.Separator (the
// item -> separator dependency internal/registry derives): role="none" and
// Separator's own data-[orientation=...] base classes both come through.
func TestItemSeparatorDefaultPinned(t *testing.T) {
	got := render(t, ui.ItemSeparator("", nil))
	want := `<div role="none" data-orientation="horizontal" ` + canonicalItemSeparatorClass("horizontal") + ` data-gsxui-slot-item-separator data-gsxui-slot-separator></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

// TestItemSeparatorOrientationOverride proves the orientation Go param
// actually overrides Separator's own default, the same competing-defaults
// mechanism as ButtonGroupSeparator (see ui/item.gsx's own comment for why
// this must be a real Go param rather than left to attrs).
func TestItemSeparatorOrientationOverride(t *testing.T) {
	got := render(t, ui.ItemSeparator("vertical", nil))
	if !strings.Contains(got, `data-orientation="vertical"`) {
		t.Errorf("missing data-orientation=vertical override\nin: %s", got)
	}
}

func TestItemSeparatorAttrsFallThrough(t *testing.T) {
	got := render(t, ui.ItemSeparator("", gsx.Attrs{{Key: "id", Value: "is1"}}))
	if !strings.Contains(got, `id="is1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

// TestItemDefaultPinned pins the zero-value (variant="default", size="default").
func TestItemDefaultPinned(t *testing.T) {
	got := render(t, ui.Item("", "", gsx.Raw("x"), nil))
	want := `<div data-variant="default" data-size="default" ` + canonicalItemClass("", []string{"variant-default"}) + ` data-gsxui-slot-item>x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

// TestItemOutlineSmPinned proves both independent switches (variant and
// size) apply together, and that "border-border" wins its tailwind-merge
// conflict against the base's own "border-transparent" (both border-color
// utilities — border-transparent is dropped, the bare "border" width
// utility is untouched).
func TestItemOutlineSmPinned(t *testing.T) {
	got := render(t, ui.Item("outline", "sm", gsx.Raw("x"), nil))
	want := `<div data-variant="outline" data-size="sm" ` + canonicalItemClass("", []string{"variant-outline"}) + ` data-gsxui-slot-item>x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestItemVariantAndSizeAxes(t *testing.T) {
	for _, variant := range []string{"default", "outline", "muted"} {
		input := variant
		if variant == "default" {
			input = ""
		}
		got := render(t, ui.Item(input, "", gsx.Raw("x"), nil))
		if !strings.Contains(got, `data-variant="`+variant+`"`) {
			t.Errorf("variant %s: missing reflected value\nin: %s", variant, got)
		}
	}
	for _, size := range []string{"default", "sm"} {
		input := size
		if size == "default" {
			input = ""
		}
		got := render(t, ui.Item("", input, gsx.Raw("x"), nil))
		if !strings.Contains(got, `data-size="`+size+`"`) {
			t.Errorf("size %s: missing reflected value\nin: %s", size, got)
		}
	}
}

func TestItemAttrsFallThrough(t *testing.T) {
	got := render(t, ui.Item("", "", nil, gsx.Attrs{{Key: "id", Value: "i1"}}))
	if !strings.Contains(got, `id="i1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

func TestItemCallerClassIsForwardedOnce(t *testing.T) {
	got := render(t, ui.Item("", "", nil, gsx.Attrs{{Key: "class", Value: "gap-8"}}))
	if strings.Count(got, "gap-8") != 1 || strings.Count(got, `class=`) != 1 {
		t.Errorf("caller class must merge in exactly once\nin: %s", got)
	}
	// gap-8 displaces the base rule's own gap-2.5 through tailwind-merge —
	// plain-vs-plain, so the caller wins before the cascade is consulted.
	if want := canonicalItemClass("", []string{"variant-default"}, "gap-8"); !strings.Contains(got, want) {
		t.Errorf("caller class not merged after the recipe's own\nwant: %s\nin: %s", want, got)
	}
}

func TestItemMediaDefaultPinned(t *testing.T) {
	got := render(t, ui.ItemMedia("", gsx.Raw("x"), nil))
	want := `<div data-variant="default" ` + canonicalItemClass("media", []string{"variant-default"}) + ` data-gsxui-slot-item-media>x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestItemMediaIconPinned(t *testing.T) {
	got := render(t, ui.ItemMedia("icon", gsx.Raw("x"), nil))
	want := `<div data-variant="icon" ` + canonicalItemClass("media", []string{"variant-icon"}) + ` data-gsxui-slot-item-media>x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestItemMediaImagePinned(t *testing.T) {
	got := render(t, ui.ItemMedia("image", gsx.Raw("x"), nil))
	want := `<div data-variant="image" ` + canonicalItemClass("media", []string{"variant-image"}) + ` data-gsxui-slot-item-media>x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestItemMediaAttrsFallThrough(t *testing.T) {
	got := render(t, ui.ItemMedia("", nil, gsx.Attrs{{Key: "id", Value: "im1"}}))
	if !strings.Contains(got, `id="im1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

func TestItemContentPinned(t *testing.T) {
	got := render(t, ui.ItemContent(gsx.Raw("x"), nil))
	want := `<div ` + canonicalItemClass("content", nil) + ` data-gsxui-slot-item-content>x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestItemContentAttrsFallThrough(t *testing.T) {
	got := render(t, ui.ItemContent(nil, gsx.Attrs{{Key: "id", Value: "ic1"}}))
	if !strings.Contains(got, `id="ic1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

func TestItemTitlePinned(t *testing.T) {
	got := render(t, ui.ItemTitle(gsx.Raw("x"), nil))
	want := `<div ` + canonicalItemClass("title", nil) + ` data-gsxui-slot-item-title>x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestItemTitleAttrsFallThrough(t *testing.T) {
	got := render(t, ui.ItemTitle(nil, gsx.Attrs{{Key: "id", Value: "it1"}}))
	if !strings.Contains(got, `id="it1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

// TestItemDescriptionPinned proves ItemDescription renders a real <p>,
// matching shadcn's own source exactly (unlike EmptyDescription, whose type
// says "p" but whose actual element is a div — see ui/item.gsx's comment).
func TestItemDescriptionPinned(t *testing.T) {
	got := render(t, ui.ItemDescription(gsx.Raw("x"), nil))
	want := `<p ` + canonicalItemClass("description", nil) + ` data-gsxui-slot-item-description>x</p>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestItemDescriptionAttrsFallThrough(t *testing.T) {
	got := render(t, ui.ItemDescription(nil, gsx.Attrs{{Key: "id", Value: "id1"}}))
	if !strings.Contains(got, `id="id1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

func TestItemActionsPinned(t *testing.T) {
	got := render(t, ui.ItemActions(gsx.Raw("x"), nil))
	want := `<div ` + canonicalItemClass("actions", nil) + ` data-gsxui-slot-item-actions>x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestItemActionsAttrsFallThrough(t *testing.T) {
	got := render(t, ui.ItemActions(nil, gsx.Attrs{{Key: "id", Value: "ia1"}}))
	if !strings.Contains(got, `id="ia1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

func TestItemHeaderPinned(t *testing.T) {
	got := render(t, ui.ItemHeader(gsx.Raw("x"), nil))
	want := `<div ` + canonicalItemClass("header", nil) + ` data-gsxui-slot-item-header>x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestItemHeaderAttrsFallThrough(t *testing.T) {
	got := render(t, ui.ItemHeader(nil, gsx.Attrs{{Key: "id", Value: "ih1"}}))
	if !strings.Contains(got, `id="ih1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

func TestItemFooterPinned(t *testing.T) {
	got := render(t, ui.ItemFooter(gsx.Raw("x"), nil))
	want := `<div ` + canonicalItemClass("footer", nil) + ` data-gsxui-slot-item-footer>x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestItemFooterAttrsFallThrough(t *testing.T) {
	got := render(t, ui.ItemFooter(nil, gsx.Attrs{{Key: "id", Value: "if1"}}))
	if !strings.Contains(got, `id="if1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

// Realistic composition: media + content (title/description) + actions row,
// inside an ItemGroup with a separator between two items — the site
// example's own shape.
func TestItemGroupWithSeparatorComposition(t *testing.T) {
	got := render(t, ui.ItemGroup(
		gsx.Fragment(
			ui.Item("outline", "",
				gsx.Fragment(
					ui.ItemMedia("icon", gsx.Raw("<svg/>"), nil),
					ui.ItemContent(
						gsx.Fragment(
							ui.ItemTitle(gsx.Raw("Invoice #1234"), nil),
							ui.ItemDescription(gsx.Raw("Paid on Jan 4"), nil),
						),
						nil,
					),
					ui.ItemActions(gsx.Raw("<button>View</button>"), nil),
				),
				nil,
			),
			ui.ItemSeparator("", nil),
			ui.Item("outline", "", gsx.Raw("second item"), nil),
		),
		nil,
	))
	for _, want := range []string{
		`data-gsxui-slot-item-group`,
		`data-variant="outline" data-size="default" ` + canonicalItemClass("", []string{"variant-outline"}) + ` data-gsxui-slot-item`,
		`data-variant="icon" ` + canonicalItemClass("media", []string{"variant-icon"}) + ` data-gsxui-slot-item-media`,
		`data-gsxui-slot-item-content`,
		`data-gsxui-slot-item-title`,
		`>Invoice #1234</div>`,
		`data-gsxui-slot-item-description`,
		`>Paid on Jan 4</p>`,
		`data-gsxui-slot-item-actions`,
		`<button>View</button>`,
		`data-gsxui-slot-item-separator data-gsxui-slot-separator`,
		`second item`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}
