package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestEmptyPinned(t *testing.T) {
	got := render(t, ui.Empty(gsx.Raw("x"), nil))
	want := `<div data-slot="empty" class="flex min-w-0 flex-1 flex-col items-center justify-center gap-6 rounded-lg border-dashed p-6 text-center text-balance md:p-12">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestEmptyAttrsFallThrough(t *testing.T) {
	got := render(t, ui.Empty(nil, gsx.Attrs{{Key: "id", Value: "e1"}}))
	if !strings.Contains(got, `id="e1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

func TestEmptyCallerClassMerges(t *testing.T) {
	got := render(t, ui.Empty(nil, gsx.Attrs{{Key: "class", Value: "gap-2"}}))
	if strings.Contains(got, "gap-6") {
		t.Errorf("base gap-6 should be dropped by caller gap-2\nin: %s", got)
	}
	if !strings.Contains(got, "gap-2") {
		t.Errorf("missing caller class gap-2\nin: %s", got)
	}
}

func TestEmptyHeaderPinned(t *testing.T) {
	got := render(t, ui.EmptyHeader(gsx.Raw("x"), nil))
	want := `<div data-slot="empty-header" class="flex max-w-sm flex-col items-center gap-2 text-center">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestEmptyHeaderAttrsFallThrough(t *testing.T) {
	got := render(t, ui.EmptyHeader(nil, gsx.Attrs{{Key: "id", Value: "eh1"}}))
	if !strings.Contains(got, `id="eh1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

// TestEmptyMediaDefaultPinned pins the zero-value ("default") variant.
// shadcn's own emptyMediaVariants cva map picks between two entirely static
// class blocks by the JS-resolved variant value — no data-[variant=...]
// selectors to preserve — so this ports as a switch inside class={}, the
// same idiom as badge/button-group.
func TestEmptyMediaDefaultPinned(t *testing.T) {
	got := render(t, ui.EmptyMedia("", gsx.Raw("x"), nil))
	want := `<div data-slot="empty-icon" data-variant="default" class="mb-2 flex shrink-0 items-center justify-center [&amp;_svg]:pointer-events-none [&amp;_svg]:shrink-0 bg-transparent">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestEmptyMediaIconPinned(t *testing.T) {
	got := render(t, ui.EmptyMedia("icon", gsx.Raw("x"), nil))
	want := `<div data-slot="empty-icon" data-variant="icon" class="mb-2 [&amp;_svg]:pointer-events-none [&amp;_svg]:shrink-0 flex size-10 shrink-0 items-center justify-center rounded-lg bg-muted text-foreground [&amp;_svg:not([class*=&#39;size-&#39;])]:size-6">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestEmptyMediaAttrsFallThrough(t *testing.T) {
	got := render(t, ui.EmptyMedia("", nil, gsx.Attrs{{Key: "id", Value: "em1"}}))
	if !strings.Contains(got, `id="em1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

func TestEmptyTitlePinned(t *testing.T) {
	got := render(t, ui.EmptyTitle(gsx.Raw("x"), nil))
	want := `<div data-slot="empty-title" class="text-lg font-medium tracking-tight">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestEmptyTitleAttrsFallThrough(t *testing.T) {
	got := render(t, ui.EmptyTitle(nil, gsx.Attrs{{Key: "id", Value: "et1"}}))
	if !strings.Contains(got, `id="et1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

// TestEmptyDescriptionPinned proves the port renders a <div>, matching
// shadcn's own actual returned element (its TS prop type says "p" but its
// JSX returns a div — see ui/empty.gsx's own comment and docs/jsx-parity.md).
func TestEmptyDescriptionPinned(t *testing.T) {
	got := render(t, ui.EmptyDescription(gsx.Raw("x"), nil))
	want := `<div data-slot="empty-description" class="text-sm/relaxed text-muted-foreground [&amp;&gt;a]:underline [&amp;&gt;a]:underline-offset-4 [&amp;&gt;a:hover]:text-primary">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestEmptyDescriptionAttrsFallThrough(t *testing.T) {
	got := render(t, ui.EmptyDescription(nil, gsx.Attrs{{Key: "id", Value: "ed1"}}))
	if !strings.Contains(got, `id="ed1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

func TestEmptyContentPinned(t *testing.T) {
	got := render(t, ui.EmptyContent(gsx.Raw("x"), nil))
	want := `<div data-slot="empty-content" class="flex w-full max-w-sm min-w-0 flex-col items-center gap-4 text-sm text-balance">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestEmptyContentAttrsFallThrough(t *testing.T) {
	got := render(t, ui.EmptyContent(nil, gsx.Attrs{{Key: "id", Value: "ec1"}}))
	if !strings.Contains(got, `id="ec1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

// Realistic composition: the icon-media empty state shape the site example
// also renders — header (media + title + description) plus a content slot
// for an action.
func TestEmptyFullComposition(t *testing.T) {
	got := render(t, ui.Empty(
		gsx.Fragment(
			ui.EmptyHeader(
				gsx.Fragment(
					ui.EmptyMedia("icon", gsx.Raw("<svg/>"), nil),
					ui.EmptyTitle(gsx.Raw("No results"), nil),
					ui.EmptyDescription(gsx.Raw("Try a different search."), nil),
				),
				nil,
			),
			ui.EmptyContent(gsx.Raw("<button>Clear filters</button>"), nil),
		),
		nil,
	))
	for _, want := range []string{
		`data-slot="empty"`,
		`data-slot="empty-header"`,
		`data-slot="empty-icon" data-variant="icon"`,
		`data-slot="empty-title"`,
		`>No results</div>`,
		`data-slot="empty-description"`,
		`>Try a different search.</div>`,
		`data-slot="empty-content"`,
		`<button>Clear filters</button>`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}
