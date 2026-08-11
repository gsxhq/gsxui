package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestPaginationPinned(t *testing.T) {
	got := render(t, ui.Pagination(gsx.Raw("x"), nil))
	want := `<nav role="navigation" aria-label="pagination" class="mx-auto flex w-full justify-center" data-gsxui-slot-pagination>x</nav>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestPaginationAttrsFallThrough(t *testing.T) {
	got := render(t, ui.Pagination(nil, gsx.Attrs{{Key: "id", Value: "pg"}}))
	if !strings.Contains(got, `id="pg"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

func TestPaginationContentPinned(t *testing.T) {
	got := render(t, ui.PaginationContent(gsx.Raw("x"), nil))
	want := `<ul class="gap-0.5 flex flex-row" data-gsxui-slot-pagination-content>x</ul>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestPaginationItemPinned(t *testing.T) {
	got := render(t, ui.PaginationItem(gsx.Raw("x"), nil))
	want := `<li data-gsxui-slot-pagination-item>x</li>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

// TestPaginationLinkDefaultPinned is the token-for-token pin against shadcn's
// PaginationLink: isActive false -> "ghost" variant, size zero-value ->
// "icon" (PaginationLinkProps' own size="icon" default, not Button's
// "default"), no aria-current, and no data-active presence marker.
func TestPaginationLinkDefaultPinned(t *testing.T) {
	got := render(t, ui.PaginationLink("/p/1", false, "", gsx.Raw("1"), nil))
	want := `<a data-variant="ghost" data-size="icon" href="/p/1" data-gsxui-slot-pagination-link data-gsxui-slot-button>1</a>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

// TestPaginationLinkActivePinned proves isActive selects the "outline"
// variant, stamps bare data-active, and adds aria-current="page" — the
// conditional-attribute mechanism replacing shadcn's
// `aria-current={isActive ? "page" : undefined}`.
func TestPaginationLinkActivePinned(t *testing.T) {
	got := render(t, ui.PaginationLink("/p/2", true, "", gsx.Raw("2"), nil))
	want := `<a aria-current="page" data-active data-variant="outline" data-size="icon" href="/p/2" data-gsxui-slot-pagination-link data-gsxui-slot-button>2</a>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestPaginationLinkInactiveNoAriaCurrent(t *testing.T) {
	got := render(t, ui.PaginationLink("/p/1", false, "", gsx.Raw("1"), nil))
	if strings.Contains(got, "aria-current") {
		t.Errorf("inactive link must not stamp aria-current at all\nin: %s", got)
	}
}

func TestPaginationLinkExplicitSize(t *testing.T) {
	got := render(t, ui.PaginationLink("/p/1", false, "default", gsx.Raw("1"), nil))
	if !strings.Contains(got, `data-size="default"`) {
		t.Errorf("explicit size=default must be reflected\nin: %s", got)
	}
}

func TestPaginationLinkSizeAxis(t *testing.T) {
	for _, size := range []string{"default", "xs", "sm", "lg", "icon", "icon-xs", "icon-sm", "icon-lg"} {
		got := render(t, ui.PaginationLink("/p/1", false, size, gsx.Raw("1"), nil))
		if !strings.Contains(got, `data-size="`+size+`"`) {
			t.Errorf("size %s: missing reflected value\nin: %s", size, got)
		}
	}
}

func TestPaginationLinkAttrsFallThrough(t *testing.T) {
	got := render(t, ui.PaginationLink("/p/1", false, "", gsx.Raw("1"), gsx.Attrs{{Key: "id", Value: "l1"}}))
	if !strings.Contains(got, `id="l1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

func TestPaginationLinkCallerClassIsForwardedOnce(t *testing.T) {
	got := render(t, ui.PaginationLink("/p/1", false, "", gsx.Raw("1"), gsx.Attrs{{Key: "class", Value: "size-10"}}))
	if strings.Count(got, `class="size-10"`) != 1 {
		t.Errorf("caller class must be the only class and render once\nin: %s", got)
	}
}

// TestPaginationPreviousPinned pins PaginationPrevious's fixed content
// (ChevronLeft + sm:-only "Previous" label) and its own class/size/aria-label
// overrides — mirrors shadcn's hardcoded JSX children, which always win over
// anything a caller could otherwise pass.
func TestPaginationPreviousPinned(t *testing.T) {
	got := render(t, ui.PaginationPrevious("/p/1", nil))
	want := `<a data-variant="ghost" data-size="default" href="/p/1" aria-label="Go to previous page" data-gsxui-slot-pagination-previous data-gsxui-slot-pagination-link data-gsxui-slot-button><svg aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="rtl:rotate-180" data-gsxui-slot-icon><path d="m15 18-6-6 6-6"/></svg><span class="hidden sm:block" data-gsxui-slot-pagination-previous-label>Previous</span></a>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestPaginationPreviousComposesAllPresenceMarkersOnLink(t *testing.T) {
	got := render(t, ui.PaginationPrevious("/p/1", nil))
	requirePresenceAttributesOnSameTag(t, got, `href="/p/1"`,
		"data-gsxui-slot-pagination-previous",
		"data-gsxui-slot-pagination-link",
		"data-gsxui-slot-button",
	)
}

func TestPaginationNextPinned(t *testing.T) {
	got := render(t, ui.PaginationNext("/p/3", nil))
	want := `<a data-variant="ghost" data-size="default" href="/p/3" aria-label="Go to next page" data-gsxui-slot-pagination-next data-gsxui-slot-pagination-link data-gsxui-slot-button><span class="hidden sm:block" data-gsxui-slot-pagination-next-label>Next</span><svg aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="rtl:rotate-180" data-gsxui-slot-icon><path d="m9 18 6-6-6-6"/></svg></a>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestPaginationEllipsisPinned(t *testing.T) {
	got := render(t, ui.PaginationEllipsis(nil))
	want := `<span aria-hidden="true" class="size-8 items-center justify-center [&amp;_svg:not([class*=&#39;size-&#39;])]:size-4 flex" data-gsxui-slot-pagination-ellipsis><svg aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" data-gsxui-slot-icon><circle cx="12" cy="12" r="1"/><circle cx="19" cy="12" r="1"/><circle cx="5" cy="12" r="1"/></svg><span data-gsxui-slot-pagination-ellipsis-label>More pages</span></span>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestPaginationEllipsisAttrsFallThrough(t *testing.T) {
	got := render(t, ui.PaginationEllipsis(gsx.Attrs{{Key: "id", Value: "e1"}}))
	if !strings.Contains(got, `id="e1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

// TestPaginationIconDependency proves ChevronLeft/ChevronRight/Ellipsis are
// actually wired in at render time, not just imported — the pagination ->
// icon import internal/registry derives Deps("pagination") from.
func TestPaginationIconDependency(t *testing.T) {
	for _, tc := range []struct {
		name string
		got  string
	}{
		{"previous", render(t, ui.PaginationPrevious("#", nil))},
		{"next", render(t, ui.PaginationNext("#", nil))},
		{"ellipsis", render(t, ui.PaginationEllipsis(nil))},
	} {
		if !strings.Contains(tc.got, `data-gsxui-slot-icon`) {
			t.Errorf("%s: expected an icon svg in render\nin: %s", tc.name, tc.got)
		}
	}
}

// Full realistic trail: prev / 1 / 2 (active) / 3 / ellipsis / next,
// exercising the whole compound composition the way the site example does.
func TestPaginationFullTrail(t *testing.T) {
	got := render(t, ui.Pagination(
		ui.PaginationContent(
			gsx.Fragment(
				ui.PaginationItem(ui.PaginationPrevious("#", nil), nil),
				ui.PaginationItem(ui.PaginationLink("#", false, "", gsx.Raw("1"), nil), nil),
				ui.PaginationItem(ui.PaginationLink("#", true, "", gsx.Raw("2"), nil), nil),
				ui.PaginationItem(ui.PaginationLink("#", false, "", gsx.Raw("3"), nil), nil),
				ui.PaginationItem(ui.PaginationEllipsis(nil), nil),
				ui.PaginationItem(ui.PaginationNext("#", nil), nil),
			),
			nil,
		),
		nil,
	))
	for _, want := range []string{
		`data-gsxui-slot-pagination-content`,
		`aria-label="Go to previous page"`,
		`>1</a>`,
		`aria-current="page" data-active data-variant="outline" data-size="icon"`,
		`>3</a>`,
		`data-gsxui-slot-pagination-ellipsis`,
		`aria-label="Go to next page"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}
