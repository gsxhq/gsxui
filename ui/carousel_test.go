package ui_test

import (
	"regexp"
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

// hasDisabledAttr matches the standalone boolean `disabled` HTML attribute.
var hasDisabledAttr = regexp.MustCompile(`\sdisabled(\s|>)`)

func TestCarouselRootHorizontalPinned(t *testing.T) {
	got := render(t, ui.Carousel("", gsx.Raw("x"), nil))
	want := `<div role="region" aria-roledescription="carousel" data-gsxui-carousel data-orientation="horizontal" data-gsxui-slot-carousel>x</div>`
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
	want := `<div data-gsxui-carousel-content data-orientation="horizontal" data-gsxui-slot-carousel-content><div data-orientation="horizontal" data-gsxui-slot-carousel-track>x</div></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestCarouselContentVerticalPinned(t *testing.T) {
	got := render(t, ui.CarouselContent("vertical", gsx.Raw("x"), nil))
	want := `<div data-gsxui-carousel-content data-orientation="vertical" data-gsxui-slot-carousel-content><div data-orientation="vertical" data-gsxui-slot-carousel-track>x</div></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestCarouselContentCallerClassMerges(t *testing.T) {
	got := render(t, ui.CarouselContent("", gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "-ml-1"}}))
	if !strings.Contains(got, `class="-ml-1" data-gsxui-slot-carousel-track`) {
		t.Errorf("caller class must remain on the track and be the only class\nin: %s", got)
	}
}

func TestCarouselItemHorizontalPinned(t *testing.T) {
	got := render(t, ui.CarouselItem("", gsx.Raw("x"), nil))
	want := `<div role="group" aria-roledescription="slide" data-gsxui-carousel-item data-orientation="horizontal" data-gsxui-slot-carousel-item>x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestCarouselItemVerticalPinned(t *testing.T) {
	got := render(t, ui.CarouselItem("vertical", gsx.Raw("x"), nil))
	want := `<div role="group" aria-roledescription="slide" data-gsxui-carousel-item data-orientation="vertical" data-gsxui-slot-carousel-item>x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestCarouselItemCallerClassMerges(t *testing.T) {
	got := render(t, ui.CarouselItem("", gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "md:basis-1/2"}}))
	if !strings.Contains(got, `class="md:basis-1/2" data-gsxui-slot-carousel-item`) {
		t.Errorf("caller class must be the only rendered class\nin: %s", got)
	}
}

func TestCarouselPreviousHorizontalPinned(t *testing.T) {
	got := render(t, ui.CarouselPrevious("", nil))
	want := `<button data-variant="outline" data-size="icon" type="button" disabled ` + canonicalButtonClass("outline", "icon") + ` data-gsxui-carousel-prev data-orientation="horizontal" data-gsxui-slot-carousel-previous data-gsxui-slot-button><svg aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" data-gsxui-slot-icon><path d="m12 19-7-7 7-7"/><path d="M19 12H5"/></svg><span data-gsxui-slot-carousel-control-label>Previous slide</span></button>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestCarouselPreviousVerticalPositioning(t *testing.T) {
	got := render(t, ui.CarouselPrevious("vertical", nil))
	for _, want := range []string{`data-orientation="vertical"`, "data-gsxui-carousel-prev", `data-variant="outline"`, `data-gsxui-slot-carousel-previous data-gsxui-slot-button`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
	if !strings.Contains(got, canonicalButtonClass("outline", "icon")) {
		t.Errorf("vertical previous lost exact canonical Button roles\nin: %s", got)
	}
}

func TestCarouselNextHorizontalPinned(t *testing.T) {
	got := render(t, ui.CarouselNext("", nil))
	want := `<button data-variant="outline" data-size="icon" type="button" ` + canonicalButtonClass("outline", "icon") + ` data-gsxui-carousel-next data-orientation="horizontal" data-gsxui-slot-carousel-next data-gsxui-slot-button><svg aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" data-gsxui-slot-icon><path d="M5 12h14"/><path d="m12 5 7 7-7 7"/></svg><span data-gsxui-slot-carousel-control-label>Next slide</span></button>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
	if hasDisabledAttr.MatchString(got) {
		t.Errorf("next must not start disabled (permissive default, see package doc comment)\nin: %s", got)
	}
}

func TestCarouselNextVerticalPositioning(t *testing.T) {
	got := render(t, ui.CarouselNext("vertical", nil))
	for _, want := range []string{`data-orientation="vertical"`, "data-gsxui-carousel-next", `data-gsxui-slot-carousel-next data-gsxui-slot-button`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
	if !strings.Contains(got, canonicalButtonClass("outline", "icon")) {
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
