package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestBreadcrumbPinned(t *testing.T) {
	got := render(t, ui.Breadcrumb(gsx.Raw("x"), nil))
	want := `<nav aria-label="breadcrumb" data-slot="breadcrumb">x</nav>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestBreadcrumbAttrsFallThrough(t *testing.T) {
	got := render(t, ui.Breadcrumb(nil, gsx.Attrs{{Key: "id", Value: "b1"}}))
	if !strings.Contains(got, `id="b1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

func TestBreadcrumbListPinned(t *testing.T) {
	got := render(t, ui.BreadcrumbList(gsx.Raw("x"), nil))
	want := `<ol data-slot="breadcrumb-list" class="flex flex-wrap items-center gap-1.5 text-sm break-words text-muted-foreground sm:gap-2.5">x</ol>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestBreadcrumbListCallerClassMerges(t *testing.T) {
	got := render(t, ui.BreadcrumbList(nil, gsx.Attrs{{Key: "class", Value: "gap-4"}}))
	if strings.Contains(got, "gap-1.5") {
		t.Errorf("base gap-1.5 should be dropped by caller gap-4\nin: %s", got)
	}
	if !strings.Contains(got, "gap-4") {
		t.Errorf("missing caller class gap-4\nin: %s", got)
	}
}

func TestBreadcrumbItemPinned(t *testing.T) {
	got := render(t, ui.BreadcrumbItem(gsx.Raw("x"), nil))
	want := `<li data-slot="breadcrumb-item" class="inline-flex items-center gap-1.5">x</li>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestBreadcrumbLinkPinned(t *testing.T) {
	got := render(t, ui.BreadcrumbLink("/docs", gsx.Raw("Docs"), nil))
	want := `<a data-slot="breadcrumb-link" href="/docs" class="transition-colors hover:text-foreground">Docs</a>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestBreadcrumbLinkAttrsFallThrough(t *testing.T) {
	got := render(t, ui.BreadcrumbLink("/docs", nil, gsx.Attrs{{Key: "id", Value: "l1"}}))
	if !strings.Contains(got, `id="l1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

func TestBreadcrumbPagePinned(t *testing.T) {
	got := render(t, ui.BreadcrumbPage(gsx.Raw("Settings"), nil))
	want := `<span data-slot="breadcrumb-page" role="link" aria-disabled="true" aria-current="page" class="font-normal text-foreground">Settings</span>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestBreadcrumbSeparatorDefaultPinned(t *testing.T) {
	// No children: renders the default ChevronRight icon, mirroring shadcn's
	// `{children ?? <ChevronRight />}`.
	got := render(t, ui.BreadcrumbSeparator(nil, nil))
	want := `<li data-slot="breadcrumb-separator" role="presentation" aria-hidden="true" class="[&amp;&gt;svg]:size-3.5"><svg data-slot="icon" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="size-4"><path d="m9 18 6-6-6-6"/></svg></li>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestBreadcrumbSeparatorChildrenOverride(t *testing.T) {
	got := render(t, ui.BreadcrumbSeparator(gsx.Raw("/"), nil))
	want := `<li data-slot="breadcrumb-separator" role="presentation" aria-hidden="true" class="[&amp;&gt;svg]:size-3.5">/</li>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestBreadcrumbEllipsisPinned(t *testing.T) {
	got := render(t, ui.BreadcrumbEllipsis(nil))
	want := `<span data-slot="breadcrumb-ellipsis" role="presentation" aria-hidden="true" class="flex size-9 items-center justify-center"><svg data-slot="icon" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="size-4"><circle cx="12" cy="12" r="1"/><circle cx="19" cy="12" r="1"/><circle cx="5" cy="12" r="1"/></svg><span class="sr-only">More</span></span>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestBreadcrumbEllipsisAttrsFallThrough(t *testing.T) {
	got := render(t, ui.BreadcrumbEllipsis(gsx.Attrs{{Key: "id", Value: "e1"}}))
	if !strings.Contains(got, `id="e1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

// Full realistic trail: proves the parts compose the way the breadcrumb
// example does, separator + ellipsis included.
func TestBreadcrumbFullTrail(t *testing.T) {
	got := render(t, ui.Breadcrumb(
		ui.BreadcrumbList(
			gsx.Fragment(
				ui.BreadcrumbItem(ui.BreadcrumbLink("/", gsx.Raw("Home"), nil), nil),
				ui.BreadcrumbSeparator(nil, nil),
				ui.BreadcrumbItem(ui.BreadcrumbEllipsis(nil), nil),
				ui.BreadcrumbSeparator(nil, nil),
				ui.BreadcrumbItem(ui.BreadcrumbPage(gsx.Raw("Settings"), nil), nil),
			),
			nil,
		),
		nil,
	))
	for _, want := range []string{
		`data-slot="breadcrumb-list"`,
		`href="/" class="transition-colors hover:text-foreground">Home</a>`,
		`data-slot="breadcrumb-separator"`,
		`data-slot="breadcrumb-ellipsis"`,
		`aria-current="page" class="font-normal text-foreground">Settings</span>`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}
