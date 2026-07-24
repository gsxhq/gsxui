package ui_test

import (
	"regexp"
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

// hasDisabledAttr matches the standalone boolean `disabled` HTML attribute,
// not the `disabled:...` Tailwind variant prefix that appears throughout
// Button's own base class string (e.g. `disabled:pointer-events-none`).
var hasDisabledAttr = regexp.MustCompile(`\sdisabled(\s|>)`)

func TestCarouselRootHorizontalPinned(t *testing.T) {
	got := render(t, ui.Carousel("", gsx.Raw("x"), nil))
	want := `<div role="region" aria-roledescription="carousel" data-slot="carousel" data-gsxui-carousel data-orientation="horizontal" class="relative">x</div>`
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
	want := `<div data-slot="carousel-content" class="overflow-x-auto snap-x snap-mandatory [scrollbar-width:none] [&amp;::-webkit-scrollbar]:hidden"><div class="flex -ml-4">x</div></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestCarouselContentVerticalPinned(t *testing.T) {
	got := render(t, ui.CarouselContent("vertical", gsx.Raw("x"), nil))
	want := `<div data-slot="carousel-content" class="overflow-y-auto snap-y snap-mandatory [scrollbar-width:none] [&amp;::-webkit-scrollbar]:hidden"><div class="flex -mt-4 flex-col">x</div></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestCarouselContentCallerClassMerges(t *testing.T) {
	got := render(t, ui.CarouselContent("", gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "-ml-1"}}))
	if strings.Contains(got, "-ml-4") {
		t.Errorf("caller -ml-1 must drop default -ml-4\nin: %s", got)
	}
	if !strings.Contains(got, "-ml-1") || !strings.Contains(got, "flex") {
		t.Errorf("want -ml-1 plus surviving structural classes\nin: %s", got)
	}
}

func TestCarouselItemHorizontalPinned(t *testing.T) {
	got := render(t, ui.CarouselItem("", gsx.Raw("x"), nil))
	want := `<div role="group" aria-roledescription="slide" data-slot="carousel-item" class="min-w-0 shrink-0 grow-0 basis-full snap-start last:snap-end pl-4 -scroll-ml-4">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestCarouselItemVerticalPinned(t *testing.T) {
	got := render(t, ui.CarouselItem("vertical", gsx.Raw("x"), nil))
	want := `<div role="group" aria-roledescription="slide" data-slot="carousel-item" class="min-w-0 shrink-0 grow-0 basis-full snap-start last:snap-end pt-4 -scroll-mt-4">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestCarouselItemCallerClassMerges(t *testing.T) {
	got := render(t, ui.CarouselItem("", gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "md:basis-1/2"}}))
	if !strings.Contains(got, "md:basis-1/2") || !strings.Contains(got, "basis-full") {
		t.Errorf("want md:basis-1/2 plus surviving structural classes\nin: %s", got)
	}
}

func TestCarouselPreviousHorizontalPinned(t *testing.T) {
	got := render(t, ui.CarouselPrevious("", nil))
	want := `<button data-variant="outline" data-size="icon" type="button" class="inline-flex shrink-0 items-center justify-center border bg-clip-padding text-sm font-medium whitespace-nowrap transition-all outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 [&amp;_svg]:pointer-events-none [&amp;_svg]:shrink-0 [&amp;_svg:not([class*=&#39;size-&#39;])]:size-4 border-border bg-background hover:bg-accent hover:text-accent-foreground dark:border-input dark:bg-input/30 dark:hover:bg-input/50 absolute size-8 rounded-full top-1/2 -left-12 -translate-y-1/2 active:not-aria-[haspopup]:translate-y-[calc(1px_-_50%)]" disabled data-slot="carousel-previous" data-gsxui-carousel-prev="true"><svg data-slot="icon" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="size-4"><path d="m12 19-7-7 7-7"/><path d="M19 12H5"/></svg><span class="sr-only">Previous slide</span></button>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestCarouselPreviousVerticalPositioning(t *testing.T) {
	got := render(t, ui.CarouselPrevious("vertical", nil))
	for _, want := range []string{"-top-12", "left-1/2", "-translate-x-1/2", "rotate-90", "data-gsxui-carousel-prev", `data-variant="outline"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
	if strings.Contains(got, "-left-12") {
		t.Errorf("vertical previous must not carry horizontal positioning\nin: %s", got)
	}
}

func TestCarouselNextHorizontalPinned(t *testing.T) {
	got := render(t, ui.CarouselNext("", nil))
	want := `<button data-variant="outline" data-size="icon" type="button" class="inline-flex shrink-0 items-center justify-center border bg-clip-padding text-sm font-medium whitespace-nowrap transition-all outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 [&amp;_svg]:pointer-events-none [&amp;_svg]:shrink-0 [&amp;_svg:not([class*=&#39;size-&#39;])]:size-4 border-border bg-background hover:bg-accent hover:text-accent-foreground dark:border-input dark:bg-input/30 dark:hover:bg-input/50 absolute size-8 rounded-full top-1/2 -right-12 -translate-y-1/2 active:not-aria-[haspopup]:translate-y-[calc(1px_-_50%)]" data-slot="carousel-next" data-gsxui-carousel-next="true"><svg data-slot="icon" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="size-4"><path d="M5 12h14"/><path d="m12 5 7 7-7 7"/></svg><span class="sr-only">Next slide</span></button>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
	if hasDisabledAttr.MatchString(got) {
		t.Errorf("next must not start disabled (permissive default, see package doc comment)\nin: %s", got)
	}
}

func TestCarouselNextVerticalPositioning(t *testing.T) {
	got := render(t, ui.CarouselNext("vertical", nil))
	for _, want := range []string{"-bottom-12", "left-1/2", "-translate-x-1/2", "rotate-90", "data-gsxui-carousel-next"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
	if strings.Contains(got, "-right-12") {
		t.Errorf("vertical next must not carry horizontal positioning\nin: %s", got)
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
