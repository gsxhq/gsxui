package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestPaginationPinned(t *testing.T) {
	got := render(t, ui.Pagination(gsx.Raw("x"), nil))
	want := `<nav role="navigation" aria-label="pagination" data-slot="pagination" class="mx-auto flex w-full justify-center">x</nav>`
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
	want := `<ul data-slot="pagination-content" class="flex flex-row items-center gap-0.5">x</ul>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestPaginationItemPinned(t *testing.T) {
	got := render(t, ui.PaginationItem(gsx.Raw("x"), nil))
	want := `<li data-slot="pagination-item">x</li>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

// TestPaginationLinkDefaultPinned is the token-for-token pin against shadcn's
// PaginationLink: isActive false -> "ghost" variant, size zero-value ->
// "icon" (PaginationLinkProps' own size="icon" default, not Button's
// "default"), no aria-current, data-active="false".
func TestPaginationLinkDefaultPinned(t *testing.T) {
	got := render(t, ui.PaginationLink("/p/1", false, "", gsx.Raw("1"), nil))
	want := `<a data-slot="pagination-link" data-active="false" href="/p/1" class="inline-flex shrink-0 items-center justify-center rounded-lg border border-transparent bg-clip-padding text-sm font-medium whitespace-nowrap transition-all outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 active:not-aria-[haspopup]:translate-y-px disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 [&amp;_svg]:pointer-events-none [&amp;_svg]:shrink-0 [&amp;_svg:not([class*=&#39;size-&#39;])]:size-4 hover:bg-accent hover:text-accent-foreground dark:hover:bg-accent/50 size-8">1</a>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

// TestPaginationLinkActivePinned proves isActive selects the "outline"
// variant, stamps data-active="true", and adds aria-current="page" — the
// conditional-attribute mechanism replacing shadcn's
// `aria-current={isActive ? "page" : undefined}`.
func TestPaginationLinkActivePinned(t *testing.T) {
	got := render(t, ui.PaginationLink("/p/2", true, "", gsx.Raw("2"), nil))
	want := `<a aria-current="page" data-slot="pagination-link" data-active="true" href="/p/2" class="inline-flex shrink-0 items-center justify-center rounded-lg border bg-clip-padding text-sm font-medium whitespace-nowrap transition-all outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 active:not-aria-[haspopup]:translate-y-px disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 [&amp;_svg]:pointer-events-none [&amp;_svg]:shrink-0 [&amp;_svg:not([class*=&#39;size-&#39;])]:size-4 border-border bg-background hover:bg-accent hover:text-accent-foreground dark:border-input dark:bg-input/30 dark:hover:bg-input/50 size-8">2</a>`
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
	if !strings.Contains(got, "h-8 gap-1.5 px-2.5 has-[&gt;svg]:px-2") {
		t.Errorf("explicit size=default should use Button's default size class\nin: %s", got)
	}
	if strings.Contains(got, "size-8") {
		t.Errorf("explicit size=default should not carry the icon size class\nin: %s", got)
	}
}

func TestPaginationLinkAttrsFallThrough(t *testing.T) {
	got := render(t, ui.PaginationLink("/p/1", false, "", gsx.Raw("1"), gsx.Attrs{{Key: "id", Value: "l1"}}))
	if !strings.Contains(got, `id="l1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

func TestPaginationLinkCallerClassMerges(t *testing.T) {
	got := render(t, ui.PaginationLink("/p/1", false, "", gsx.Raw("1"), gsx.Attrs{{Key: "class", Value: "size-10"}}))
	if strings.Contains(got, "size-8") {
		t.Errorf("base size-8 should be dropped by caller size-10\nin: %s", got)
	}
	if !strings.Contains(got, "size-10") {
		t.Errorf("missing caller class size-10\nin: %s", got)
	}
}

// TestPaginationPreviousPinned pins PaginationPrevious's fixed content
// (ChevronLeft + sm:-only "Previous" label) and its own class/size/aria-label
// overrides — mirrors shadcn's hardcoded JSX children, which always win over
// anything a caller could otherwise pass.
func TestPaginationPreviousPinned(t *testing.T) {
	got := render(t, ui.PaginationPrevious("/p/1", nil))
	want := `<a data-slot="pagination-link" data-active="false" href="/p/1" class="inline-flex shrink-0 items-center justify-center rounded-lg border border-transparent bg-clip-padding text-sm font-medium whitespace-nowrap transition-all outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 active:not-aria-[haspopup]:translate-y-px disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 [&amp;_svg]:pointer-events-none [&amp;_svg]:shrink-0 [&amp;_svg:not([class*=&#39;size-&#39;])]:size-4 hover:bg-accent hover:text-accent-foreground dark:hover:bg-accent/50 h-8 gap-1.5 px-2.5 has-[&gt;svg]:px-2 pl-1.5!" aria-label="Go to previous page"><svg data-slot="icon" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="size-4"><path d="m15 18-6-6 6-6"/></svg><span class="hidden sm:block">Previous</span></a>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestPaginationNextPinned(t *testing.T) {
	got := render(t, ui.PaginationNext("/p/3", nil))
	want := `<a data-slot="pagination-link" data-active="false" href="/p/3" class="inline-flex shrink-0 items-center justify-center rounded-lg border border-transparent bg-clip-padding text-sm font-medium whitespace-nowrap transition-all outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 active:not-aria-[haspopup]:translate-y-px disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 [&amp;_svg]:pointer-events-none [&amp;_svg]:shrink-0 [&amp;_svg:not([class*=&#39;size-&#39;])]:size-4 hover:bg-accent hover:text-accent-foreground dark:hover:bg-accent/50 h-8 gap-1.5 px-2.5 has-[&gt;svg]:px-2 pr-1.5!" aria-label="Go to next page"><span class="hidden sm:block">Next</span><svg data-slot="icon" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="size-4"><path d="m9 18 6-6-6-6"/></svg></a>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestPaginationEllipsisPinned(t *testing.T) {
	got := render(t, ui.PaginationEllipsis(nil))
	want := `<span aria-hidden="true" data-slot="pagination-ellipsis" class="flex size-8 items-center justify-center"><svg data-slot="icon" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="size-4"><circle cx="12" cy="12" r="1"/><circle cx="19" cy="12" r="1"/><circle cx="5" cy="12" r="1"/></svg><span class="sr-only">More pages</span></span>`
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
		if !strings.Contains(tc.got, `data-slot="icon"`) {
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
		`data-slot="pagination-content"`,
		`aria-label="Go to previous page"`,
		`>1</a>`,
		`aria-current="page" data-slot="pagination-link" data-active="true"`,
		`>3</a>`,
		`data-slot="pagination-ellipsis"`,
		`aria-label="Go to next page"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}
