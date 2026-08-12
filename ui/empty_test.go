package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestEmptyPinned(t *testing.T) {
	got := render(t, ui.Empty(gsx.Raw("x"), nil))
	want := `<div class="gap-4 rounded-xl border-dashed p-6 flex flex-col w-full min-w-0 flex-1 items-center justify-center text-center text-balance" data-gsxui-slot-empty>x</div>`
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

func TestEmptyCallerClassIsForwardedOnce(t *testing.T) {
	got := render(t, ui.Empty(nil, gsx.Attrs{{Key: "class", Value: "gap-2"}}))
	if strings.Count(got, "gap-2") != 1 || strings.Count(got, `class="`) != 1 {
		t.Errorf("caller class must be forwarded exactly once, merged with the recipe class\nin: %s", got)
	}
}

func TestEmptyHeaderPinned(t *testing.T) {
	got := render(t, ui.EmptyHeader(gsx.Raw("x"), nil))
	want := `<div class="gap-2 flex flex-col max-w-sm items-center text-center" data-gsxui-slot-empty-header>x</div>`
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
	want := `<div data-variant="default" class="mb-2 flex shrink-0 items-center justify-center [&amp;_svg]:shrink-0 [&amp;_svg]:pointer-events-none bg-transparent" data-gsxui-slot-empty-icon>x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestEmptyMediaIconPinned(t *testing.T) {
	got := render(t, ui.EmptyMedia("icon", gsx.Raw("x"), nil))
	want := `<div data-variant="icon" class="mb-2 [&amp;_svg]:shrink-0 [&amp;_svg]:pointer-events-none bg-muted text-foreground flex size-8 shrink-0 items-center justify-center rounded-lg [&amp;_svg:not([class*=&#39;size-&#39;])]:size-4" data-gsxui-slot-empty-icon>x</div>`
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
	want := `<div class="text-sm font-medium tracking-tight" data-gsxui-slot-empty-title>x</div>`
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
	want := `<div class="text-sm/relaxed text-muted-foreground [&amp;&gt;a]:underline [&amp;&gt;a]:underline-offset-4 [&amp;&gt;a:hover]:text-primary" data-gsxui-slot-empty-description>x</div>`
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
	want := `<div class="gap-2.5 text-sm flex flex-col w-full max-w-sm min-w-0 items-center text-balance" data-gsxui-slot-empty-content>x</div>`
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
		`data-gsxui-slot-empty`,
		`data-gsxui-slot-empty-header`,
		`data-variant="icon" class="`,
		`data-gsxui-slot-empty-icon`,
		`data-gsxui-slot-empty-title`,
		`>No results</div>`,
		`data-gsxui-slot-empty-description`,
		`>Try a different search.</div>`,
		`data-gsxui-slot-empty-content`,
		`<button>Clear filters</button>`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}
