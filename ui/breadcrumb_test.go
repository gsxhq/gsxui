package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestBreadcrumbPinned(t *testing.T) {
	got := render(t, ui.Breadcrumb(gsx.Raw("x"), nil))
	want := `<nav aria-label="breadcrumb" data-gsxui-slot-breadcrumb>x</nav>`
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
	want := `<ol class="flex flex-wrap items-center gap-1.5 text-sm break-words text-muted-foreground" data-gsxui-slot-breadcrumb-list>x</ol>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestBreadcrumbListCallerClassIsForwardedOnce(t *testing.T) {
	// Not a bare class="gap-4" match: BreadcrumbList now carries its own
	// recipe class, so the caller's gap-4 merges in (replacing the
	// recipe's own gap-1.5, same-property override) alongside it rather
	// than rendering as the entire attribute value.
	got := render(t, ui.BreadcrumbList(nil, gsx.Attrs{{Key: "class", Value: "gap-4"}}))
	if strings.Count(got, "gap-4") != 1 {
		t.Errorf("caller class must merge in exactly once\nin: %s", got)
	}
	if strings.Count(got, "class=") != 1 {
		t.Errorf("expected exactly one class= attribute\nin: %s", got)
	}
}

func TestBreadcrumbItemPinned(t *testing.T) {
	got := render(t, ui.BreadcrumbItem(gsx.Raw("x"), nil))
	want := `<li class="inline-flex items-center gap-1" data-gsxui-slot-breadcrumb-item>x</li>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestBreadcrumbLinkPinned(t *testing.T) {
	got := render(t, ui.BreadcrumbLink("/docs", gsx.Raw("Docs"), nil))
	want := `<a href="/docs" class="transition-colors hover:text-foreground" data-gsxui-slot-breadcrumb-link>Docs</a>`
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
	want := `<span role="link" aria-disabled="true" aria-current="page" class="font-normal text-foreground" data-gsxui-slot-breadcrumb-page>Settings</span>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestBreadcrumbSeparatorDefaultPinned(t *testing.T) {
	// No children: renders the default ChevronRight icon, mirroring shadcn's
	// `{children ?? <ChevronRight />}`.
	got := render(t, ui.BreadcrumbSeparator(nil, nil))
	want := `<li role="presentation" aria-hidden="true" class="[&amp;&gt;svg]:size-3.5 [&amp;&gt;svg]:rtl:rotate-180" data-gsxui-slot-breadcrumb-separator><svg aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" data-gsxui-slot-icon><path d="m9 18 6-6-6-6"/></svg></li>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestBreadcrumbSeparatorChildrenOverride(t *testing.T) {
	got := render(t, ui.BreadcrumbSeparator(gsx.Raw("/"), nil))
	want := `<li role="presentation" aria-hidden="true" class="[&amp;&gt;svg]:size-3.5 [&amp;&gt;svg]:rtl:rotate-180" data-gsxui-slot-breadcrumb-separator>/</li>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestBreadcrumbEllipsisPinned(t *testing.T) {
	got := render(t, ui.BreadcrumbEllipsis(nil))
	want := `<span role="presentation" aria-hidden="true" class="flex size-5 items-center justify-center [&amp;&gt;svg]:size-4" data-gsxui-slot-breadcrumb-ellipsis><svg aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" data-gsxui-slot-icon><circle cx="12" cy="12" r="1"/><circle cx="19" cy="12" r="1"/><circle cx="5" cy="12" r="1"/></svg><span data-gsxui-slot-breadcrumb-ellipsis-label>More</span></span>`
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
		`data-gsxui-slot-breadcrumb-list`,
		// Not anchored to href="/" immediately preceding the marker: Link
		// now carries its own recipe class between them.
		`href="/" class="transition-colors hover:text-foreground" data-gsxui-slot-breadcrumb-link>Home</a>`,
		`data-gsxui-slot-breadcrumb-separator`,
		`data-gsxui-slot-breadcrumb-ellipsis`,
		`aria-current="page" class="font-normal text-foreground" data-gsxui-slot-breadcrumb-page>Settings</span>`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}
