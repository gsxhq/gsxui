package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestToggleGroupRootPinned(t *testing.T) {
	// Exact full-render pin for the zero-value-variant/size/spacing root of
	// a type="multiple" group — role="toolbar", data-orientation stamped
	// horizontal, the dead upstream shadow selector dropped per this port's
	// ADAPT (nova's own .cn-toggle-group precedent, see docs/jsx-parity.md).
	got := render(t, ui.ToggleGroup("multiple", "", "", "", gsx.Raw("x"), nil))
	want := `<div data-slot="toggle-group" data-gsxui-toggle-group data-variant="default" data-size="default" data-spacing="0" data-orientation="horizontal" role="toolbar" style="--gap: 0" class="group/toggle-group flex w-fit items-center gap-[--spacing(var(--gap))] rounded-lg data-[size=sm]:rounded-[min(var(--radius-md),10px)]">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestToggleGroupRootSingleRole(t *testing.T) {
	got := render(t, ui.ToggleGroup("single", "", "", "", gsx.Raw("x"), nil))
	if !strings.Contains(got, `role="radiogroup"`) {
		t.Errorf("want role=radiogroup\nin: %s", got)
	}
}

func TestToggleGroupItemSinglePinned(t *testing.T) {
	// type="single" item: role="radio" + aria-checked, NOT aria-pressed.
	got := render(t, ui.ToggleGroupItem("single", "", "", "", true, "bold", gsx.Raw("B"), nil))
	want := `<button type="button" data-slot="toggle-group-item" data-gsxui-toggle-group-item data-variant="default" data-size="default" data-spacing="0" data-orientation="horizontal" data-state="on" data-value="bold" role="radio" aria-checked="true" class="inline-flex items-center justify-center gap-1 rounded-lg text-sm font-medium whitespace-nowrap transition-[color,box-shadow] outline-none hover:bg-muted hover:text-muted-foreground focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 data-[state=on]:bg-accent data-[state=on]:text-accent-foreground dark:aria-invalid:ring-destructive/40 [&amp;_svg]:pointer-events-none [&amp;_svg]:shrink-0 [&amp;_svg:not([class*=&#39;size-&#39;])]:size-4 bg-transparent h-8 has-[&gt;svg]:px-2 w-auto min-w-0 shrink-0 px-3 focus:z-10 focus-visible:z-10 data-[spacing=0]:rounded-none data-[spacing=0]:shadow-none data-[spacing=0]:data-[variant=outline]:border-l-0 data-[spacing=0]:data-[variant=outline]:first:border-l data-[orientation=horizontal]:data-[spacing=0]:first:rounded-l-lg data-[orientation=horizontal]:data-[spacing=0]:last:rounded-r-lg">B</button>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
	if strings.Contains(got, "aria-pressed") {
		t.Errorf("single-type item must not stamp aria-pressed\nin: %s", got)
	}
}

func TestToggleGroupItemMultiplePinned(t *testing.T) {
	// type="multiple" item: aria-pressed, no role override.
	got := render(t, ui.ToggleGroupItem("multiple", "", "", "", false, "bold", gsx.Raw("B"), nil))
	if !strings.Contains(got, `aria-pressed="false"`) {
		t.Errorf("want aria-pressed=false\nin: %s", got)
	}
	if strings.Contains(got, `role="radio"`) {
		t.Errorf("multiple-type item must not stamp role=radio\nin: %s", got)
	}
	if !strings.Contains(got, `data-state="off"`) {
		t.Errorf("want data-state=off\nin: %s", got)
	}
}

func TestToggleGroupItemOutlineVariant(t *testing.T) {
	got := render(t, ui.ToggleGroupItem("multiple", "outline", "", "", false, "bold", gsx.Raw("B"), nil))
	for _, want := range []string{
		`data-variant="outline"`,
		"border border-input bg-transparent hover:bg-accent hover:text-accent-foreground",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestToggleGroupItemSpacing(t *testing.T) {
	got := render(t, ui.ToggleGroupItem("multiple", "outline", "sm", "2", false, "bold", gsx.Raw("B"), nil))
	for _, want := range []string{`data-spacing="2"`, `data-size="sm"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestToggleGroupDisabledCascade(t *testing.T) {
	// No group-level context to cascade disabled through (see the package
	// doc comment's GAP note) — the caller passes disabled explicitly to
	// the root (inert on the div, present for caller CSS/JS hooks) and to
	// every item that should actually be non-interactive.
	root := render(t, ui.ToggleGroup("multiple", "", "", "", gsx.Raw("x"), gsx.Attrs{{Key: "disabled", Value: true}}))
	if !strings.Contains(root, "disabled") {
		t.Errorf("root: want disabled attribute\nin: %s", root)
	}
	item := render(t, ui.ToggleGroupItem("multiple", "", "", "", false, "bold", gsx.Raw("B"), gsx.Attrs{{Key: "disabled", Value: true}}))
	if !strings.Contains(item, "disabled") {
		t.Errorf("item: want disabled attribute\nin: %s", item)
	}
}

func TestToggleGroupCallerClassMerges(t *testing.T) {
	got := render(t, ui.ToggleGroup("multiple", "", "", "", gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "gap-8"}}))
	if strings.Contains(got, "w-fit") == false {
		t.Errorf("want surviving structural class w-fit\nin: %s", got)
	}
	if !strings.Contains(got, "gap-8") {
		t.Errorf("want caller gap-8\nin: %s", got)
	}
}

func TestToggleGroupItemCallerClassMerges(t *testing.T) {
	got := render(t, ui.ToggleGroupItem("multiple", "", "", "", false, "bold", gsx.Raw("B"), gsx.Attrs{{Key: "class", Value: "px-8"}}))
	if strings.Contains(got, "px-3") {
		t.Errorf("caller px-8 must drop default px-3\nin: %s", got)
	}
	if !strings.Contains(got, "px-8") || !strings.Contains(got, "inline-flex") {
		t.Errorf("want px-8 plus surviving structural classes\nin: %s", got)
	}
}

func TestToggleGroupAttrsFallThrough(t *testing.T) {
	got := render(t, ui.ToggleGroup("multiple", "", "", "", gsx.Raw("x"), gsx.Attrs{{Key: "aria-label", Value: "Text formatting"}}))
	if !strings.Contains(got, `aria-label="Text formatting"`) {
		t.Errorf("missing aria-label\nin: %s", got)
	}
}
