package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestItemGroupPinned(t *testing.T) {
	got := render(t, ui.ItemGroup(gsx.Raw("x"), nil))
	want := `<div role="list" data-slot="item-group" class="group/item-group flex flex-col">x</div>`
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

func TestItemGroupCallerClassMerges(t *testing.T) {
	got := render(t, ui.ItemGroup(nil, gsx.Attrs{{Key: "class", Value: "gap-4"}}))
	if !strings.Contains(got, "gap-4") {
		t.Errorf("missing caller class gap-4\nin: %s", got)
	}
}

// TestItemSeparatorDefaultPinned pins the zero-value ("horizontal")
// orientation and proves ItemSeparator actually composes ui.Separator (the
// item -> separator dependency internal/registry derives): role="none" and
// Separator's own data-[orientation=...] base classes both come through.
func TestItemSeparatorDefaultPinned(t *testing.T) {
	got := render(t, ui.ItemSeparator("", nil))
	want := `<div role="none" data-orientation="horizontal" class="shrink-0 bg-border data-[orientation=horizontal]:h-px data-[orientation=horizontal]:w-full data-[orientation=vertical]:h-full data-[orientation=vertical]:w-px my-2" data-slot="item-separator"></div>`
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
	want := `<div data-slot="item" data-variant="default" data-size="default" class="group/item flex flex-wrap items-center rounded-lg border border-transparent text-sm transition-colors duration-100 outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 [a]:transition-colors [a]:hover:bg-accent/50 bg-transparent gap-2.5 px-3 py-2.5">x</div>`
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
	want := `<div data-slot="item" data-variant="outline" data-size="sm" class="group/item flex flex-wrap items-center rounded-lg border text-sm transition-colors duration-100 outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 [a]:transition-colors [a]:hover:bg-accent/50 border-border gap-2.5 px-3 py-2.5">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
	if strings.Contains(got, "border-transparent") {
		t.Errorf("border-transparent should be dropped by border-border\nin: %s", got)
	}
}

func TestItemMutedPinned(t *testing.T) {
	got := render(t, ui.Item("muted", "", gsx.Raw("x"), nil))
	if !strings.Contains(got, `data-variant="muted"`) || !strings.Contains(got, "bg-muted/50") {
		t.Errorf("missing muted variant class\nin: %s", got)
	}
}

func TestItemAttrsFallThrough(t *testing.T) {
	got := render(t, ui.Item("", "", nil, gsx.Attrs{{Key: "id", Value: "i1"}}))
	if !strings.Contains(got, `id="i1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

func TestItemCallerClassMerges(t *testing.T) {
	got := render(t, ui.Item("", "", nil, gsx.Attrs{{Key: "class", Value: "gap-8"}}))
	if !strings.Contains(got, "gap-8") {
		t.Errorf("missing caller class gap-8\nin: %s", got)
	}
}

func TestItemMediaDefaultPinned(t *testing.T) {
	got := render(t, ui.ItemMedia("", gsx.Raw("x"), nil))
	want := `<div data-slot="item-media" data-variant="default" class="flex shrink-0 items-center justify-center gap-2 group-has-[[data-slot=item-description]]/item:translate-y-0.5 group-has-[[data-slot=item-description]]/item:self-start [&amp;_svg]:pointer-events-none bg-transparent">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestItemMediaIconPinned(t *testing.T) {
	got := render(t, ui.ItemMedia("icon", gsx.Raw("x"), nil))
	want := `<div data-slot="item-media" data-variant="icon" class="flex shrink-0 items-center justify-center gap-2 group-has-[[data-slot=item-description]]/item:translate-y-0.5 group-has-[[data-slot=item-description]]/item:self-start [&amp;_svg]:pointer-events-none size-8 rounded-sm border bg-muted [&amp;_svg:not([class*=&#39;size-&#39;])]:size-4">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestItemMediaImagePinned(t *testing.T) {
	got := render(t, ui.ItemMedia("image", gsx.Raw("x"), nil))
	want := `<div data-slot="item-media" data-variant="image" class="flex shrink-0 items-center justify-center gap-2 group-has-[[data-slot=item-description]]/item:translate-y-0.5 group-has-[[data-slot=item-description]]/item:self-start [&amp;_svg]:pointer-events-none size-10 overflow-hidden rounded-sm [&amp;_img]:size-full [&amp;_img]:object-cover">x</div>`
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
	want := `<div data-slot="item-content" class="flex flex-1 flex-col gap-1 [&amp;+[data-slot=item-content]]:flex-none">x</div>`
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
	want := `<div data-slot="item-title" class="flex w-fit items-center gap-2 text-sm leading-snug font-medium">x</div>`
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
	want := `<p data-slot="item-description" class="line-clamp-2 text-sm leading-normal font-normal text-balance text-muted-foreground [&amp;&gt;a]:underline [&amp;&gt;a]:underline-offset-4 [&amp;&gt;a:hover]:text-primary">x</p>`
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
	want := `<div data-slot="item-actions" class="flex items-center gap-2">x</div>`
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
	want := `<div data-slot="item-header" class="flex basis-full items-center justify-between gap-2">x</div>`
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
	want := `<div data-slot="item-footer" class="flex basis-full items-center justify-between gap-2">x</div>`
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
		`data-slot="item-group"`,
		`data-slot="item" data-variant="outline"`,
		`data-slot="item-media" data-variant="icon"`,
		`data-slot="item-content"`,
		`data-slot="item-title"`,
		`>Invoice #1234</div>`,
		`data-slot="item-description"`,
		`>Paid on Jan 4</p>`,
		`data-slot="item-actions"`,
		`<button>View</button>`,
		`data-slot="item-separator"`,
		`second item`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}
