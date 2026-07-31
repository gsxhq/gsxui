package ui_test

import (
	"html"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/internal/stylegen"
	"github.com/gsxhq/gsxui/merge"
	"github.com/gsxhq/gsxui/ui"
)

// novaCarouselRecipe loads the default style's Carousel recipe once — same
// pattern as novaButtonRecipe, which these tests already lean on for the two
// Button-composing arrow parts.
var novaCarouselRecipe = sync.OnceValue(func() recipe.Style {
	path := filepath.Join("..", "registry", "styles", stylegen.DefaultStyle, "carousel.css")
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

func carouselRecipeUtilities(class string) []string {
	rule, ok := novaCarouselRecipe().Lookup(class)
	if !ok {
		panic("default style declares no recipe " + class)
	}
	return rule.Utilities
}

func canonicalCarouselClass(classes ...string) string {
	var utilities []string
	for _, class := range classes {
		utilities = append(utilities, carouselRecipeUtilities(class)...)
	}
	return `class="` + html.EscapeString(merge.Merge(utilities)) + `"`
}

// canonicalCarouselArrowClass is the class attribute CarouselPrevious /
// CarouselNext render: Button's own outline/icon roles with Carousel's slot
// utilities appended as a caller class, merged once. Carousel's rounded-full
// replaces Button's rounded-lg, and (horizontal only) Carousel's press token
// repeats Button's modifier chain so it replaces translate-y-px rather than
// stacking with it.
func canonicalCarouselArrowClass(slot, orientation string) string {
	caller := append([]string(nil), carouselRecipeUtilities("gsxui-recipe-carousel-"+slot)...)
	caller = append(caller, carouselRecipeUtilities("gsxui-recipe-carousel-"+slot+"-orientation-"+orientation)...)
	return canonicalButtonClass("outline", "icon", caller...)
}

// hasDisabledAttr matches the standalone boolean `disabled` HTML attribute.
var hasDisabledAttr = regexp.MustCompile(`\sdisabled(\s|>)`)

func TestCarouselRootHorizontalPinned(t *testing.T) {
	got := render(t, ui.Carousel("", gsx.Raw("x"), nil))
	want := `<div role="region" aria-roledescription="carousel" data-orientation="horizontal" ` + canonicalCarouselClass("gsxui-recipe-carousel") + ` data-gsxui-slot-carousel>x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestCarouselRootVerticalOrientation(t *testing.T) {
	got := render(t, ui.Carousel("vertical", gsx.Raw("x"), nil))
	if !strings.Contains(got, `data-orientation="vertical"`) {
		t.Errorf("want data-orientation=vertical\nin: %s", got)
	}
}

func TestCarouselContentHorizontalPinned(t *testing.T) {
	got := render(t, ui.CarouselContent("", gsx.Raw("x"), nil))
	want := `<div data-orientation="horizontal" ` + canonicalCarouselClass("gsxui-recipe-carousel-content") + ` data-gsxui-slot-carousel-content><div data-orientation="horizontal" ` + canonicalCarouselClass("gsxui-recipe-carousel-track-orientation-horizontal") + ` data-gsxui-slot-carousel-track>x</div></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestCarouselContentVerticalPinned(t *testing.T) {
	got := render(t, ui.CarouselContent("vertical", gsx.Raw("x"), nil))
	want := `<div data-orientation="vertical" ` + canonicalCarouselClass("gsxui-recipe-carousel-content") + ` data-gsxui-slot-carousel-content><div data-orientation="vertical" ` + canonicalCarouselClass("gsxui-recipe-carousel-track-orientation-vertical") + ` data-gsxui-slot-carousel-track>x</div></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestCarouselContentCallerClassMerges(t *testing.T) {
	got := render(t, ui.CarouselContent("", gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "-ml-1"}}))
	// The track's own -ml-4 and the caller's -ml-1 are both plain utilities
	// on the same element now, so merge.Merge replaces rather than stacks.
	want := `class="-ml-1" data-gsxui-slot-carousel-track`
	if !strings.Contains(got, want) {
		t.Errorf("caller class must displace the track's own spacing\nwant: %s\nin: %s", want, got)
	}
}

func TestCarouselItemHorizontalPinned(t *testing.T) {
	got := render(t, ui.CarouselItem("", gsx.Raw("x"), nil))
	want := `<div role="group" aria-roledescription="slide" data-orientation="horizontal" ` + canonicalCarouselClass("gsxui-recipe-carousel-item-orientation-horizontal") + ` data-gsxui-slot-carousel-item>x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestCarouselItemVerticalPinned(t *testing.T) {
	got := render(t, ui.CarouselItem("vertical", gsx.Raw("x"), nil))
	want := `<div role="group" aria-roledescription="slide" data-orientation="vertical" ` + canonicalCarouselClass("gsxui-recipe-carousel-item-orientation-vertical") + ` data-gsxui-slot-carousel-item>x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestCarouselItemCallerClassMerges(t *testing.T) {
	got := render(t, ui.CarouselItem("", gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "md:basis-1/2"}}))
	// md:basis-1/2 contests nothing the item's own pl-4/-scroll-ml-4 set, so
	// both survive the merge — the caller class follows the slot's own.
	want := canonicalCarouselClass("gsxui-recipe-carousel-item-orientation-horizontal")
	want = want[:len(want)-1] + ` md:basis-1/2" data-gsxui-slot-carousel-item`
	if !strings.Contains(got, want) {
		t.Errorf("caller class must merge onto the item's own\nwant: %s\nin: %s", want, got)
	}
}

func TestCarouselPreviousHorizontalPinned(t *testing.T) {
	got := render(t, ui.CarouselPrevious("", nil))
	want := `<button data-variant="outline" data-size="icon" type="button" disabled ` + canonicalCarouselArrowClass("previous", "horizontal") + ` data-orientation="horizontal" data-gsxui-slot-carousel-previous data-gsxui-slot-button><svg aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" data-gsxui-slot-icon><path d="m12 19-7-7 7-7"/><path d="M19 12H5"/></svg><span data-gsxui-slot-carousel-control-label>Previous slide</span></button>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestCarouselPreviousVerticalPositioning(t *testing.T) {
	got := render(t, ui.CarouselPrevious("vertical", nil))
	for _, want := range []string{`data-orientation="vertical"`, "data-gsxui-slot-carousel-previous", `data-variant="outline"`, `data-gsxui-slot-carousel-previous data-gsxui-slot-button`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
	if !strings.Contains(got, canonicalCarouselArrowClass("previous", "vertical")) {
		t.Errorf("vertical previous lost exact canonical Button roles\nin: %s", got)
	}
}

func TestCarouselNextHorizontalPinned(t *testing.T) {
	got := render(t, ui.CarouselNext("", nil))
	want := `<button data-variant="outline" data-size="icon" type="button" ` + canonicalCarouselArrowClass("next", "horizontal") + ` data-orientation="horizontal" data-gsxui-slot-carousel-next data-gsxui-slot-button><svg aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" data-gsxui-slot-icon><path d="M5 12h14"/><path d="m12 5 7 7-7 7"/></svg><span data-gsxui-slot-carousel-control-label>Next slide</span></button>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
	if hasDisabledAttr.MatchString(got) {
		t.Errorf("next must not start disabled (permissive default, see package doc comment)\nin: %s", got)
	}
}

func TestCarouselNextVerticalPositioning(t *testing.T) {
	got := render(t, ui.CarouselNext("vertical", nil))
	for _, want := range []string{`data-orientation="vertical"`, "data-gsxui-slot-carousel-next", `data-gsxui-slot-carousel-next data-gsxui-slot-button`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
	if !strings.Contains(got, canonicalCarouselArrowClass("next", "vertical")) {
		t.Errorf("vertical next lost exact canonical Button roles\nin: %s", got)
	}
}

func TestCarouselPreviousStartsDisabled(t *testing.T) {
	// A freshly mounted scroll container always starts at scrollLeft/
	// scrollTop 0 — Previous genuinely has nowhere to go, see the package
	// doc comment.
	got := render(t, ui.CarouselPrevious("", nil))
	if !hasDisabledAttr.MatchString(got) {
		t.Errorf("want disabled attribute\nin: %s", got)
	}
}

func TestCarouselAttrsFallThrough(t *testing.T) {
	got := render(t, ui.Carousel("", gsx.Raw("x"), gsx.Attrs{{Key: "id", Value: "gallery"}}))
	if !strings.Contains(got, `id="gallery"`) {
		t.Errorf("missing id\nin: %s", got)
	}
}
