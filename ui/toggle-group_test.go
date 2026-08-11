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
	// group/toggle-group was added (registry/canonical/toggle-group.gsx) so
	// ToggleGroupItem's own group-data-[…]/toggle-group: selectors have a
	// real ancestor to match — see the style-porter report's "missing
	// group/toggle-group marker class" entry.
	got := render(t, ui.ToggleGroup("multiple", "", "", "", gsx.Raw("x"), nil))
	want := `<div data-variant="default" data-size="default" data-spacing="0" data-orientation="horizontal" role="toolbar" style="--gap: 0" class="group/toggle-group flex rounded-lg" data-gsxui-slot-toggle-group>x</div>`
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
	//
	// The class string now also composes Toggle's own base/variant/size
	// presentation (hover:text-foreground .. size-4, per style) — see
	// registry/styles/nova/toggle-group.css's own comment. rounded-lg (from
	// that composed base) is not visible below: tailwind-merge correctly
	// drops it in favour of the spacing="0" arm's own unconditional
	// rounded-none (same class group, same modifier, later in the merge
	// order), which is the pre-existing, unchanged squared-off-pill
	// behaviour this pin already covered.
	// has-data-[icon=…] (dead selector, never stamped by our markup) is
	// replaced with has-[>svg]: throughout, and the base rule's own
	// icon-only padding is no longer duplicated per group-data-[spacing=0]
	// arm — see the style-porter report's "Systemic: has-data-[icon=…]"
	// entry.
	got := render(t, ui.ToggleGroupItem("single", "", "", "", true, "bold", gsx.Raw("B"), nil))
	want := `<button type="button" data-variant="default" data-size="default" data-spacing="0" data-orientation="horizontal" data-state="on" data-value="bold" role="radio" aria-checked="true" class="hover:text-foreground focus-visible:border-ring focus-visible:ring-ring/50 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive disabled:pointer-events-none disabled:opacity-50 data-[state=on]:bg-muted data-[variant=outline]:border-input data-[variant=outline]:hover:bg-muted data-[variant=outline]:border gap-1 text-sm font-medium transition-all [&amp;_svg:not([class*=&#39;size-&#39;])]:size-4 group-data-[spacing=0]/toggle-group:rounded-none group-data-[spacing=0]/toggle-group:px-2 group-data-horizontal/toggle-group:data-[spacing=0]:first:rounded-l-lg group-data-vertical/toggle-group:data-[spacing=0]:first:rounded-t-lg group-data-horizontal/toggle-group:data-[spacing=0]:last:rounded-r-lg group-data-vertical/toggle-group:data-[spacing=0]:last:rounded-b-lg h-8 min-w-8 px-2.5 has-[&gt;svg]:px-2 group-data-[spacing=0]/toggle-group:has-[&gt;svg]:px-2 rounded-none shadow-none data-[variant=outline]:border-s-0 data-[variant=outline]:first:border-s data-[spacing=0]:first:rounded-s-lg data-[spacing=0]:last:rounded-e-lg" data-gsxui-slot-toggle-group-item data-gsxui-slot-toggle>B</button>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
	if strings.Contains(got, "aria-pressed") {
		t.Errorf("single-type item must not stamp aria-pressed\nin: %s", got)
	}
}

func TestToggleGroupItemMultiplePinned(t *testing.T) {
	// type="multiple" item: aria-pressed, no role override — identical
	// class string to the single-type pinned case above (variant/size
	// unchanged); only data-state/data-value/the ARIA attribute pair differ.
	// See TestToggleGroupItemSinglePinned for why rounded-lg is composed in
	// but not visible in the merged string.
	// Same has-[>svg]: replacement as TestToggleGroupItemSinglePinned above.
	got := render(t, ui.ToggleGroupItem("multiple", "", "", "", false, "bold", gsx.Raw("B"), nil))
	want := `<button type="button" data-variant="default" data-size="default" data-spacing="0" data-orientation="horizontal" data-state="off" data-value="bold" aria-pressed="false" class="hover:text-foreground focus-visible:border-ring focus-visible:ring-ring/50 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive disabled:pointer-events-none disabled:opacity-50 data-[state=on]:bg-muted data-[variant=outline]:border-input data-[variant=outline]:hover:bg-muted data-[variant=outline]:border gap-1 text-sm font-medium transition-all [&amp;_svg:not([class*=&#39;size-&#39;])]:size-4 group-data-[spacing=0]/toggle-group:rounded-none group-data-[spacing=0]/toggle-group:px-2 group-data-horizontal/toggle-group:data-[spacing=0]:first:rounded-l-lg group-data-vertical/toggle-group:data-[spacing=0]:first:rounded-t-lg group-data-horizontal/toggle-group:data-[spacing=0]:last:rounded-r-lg group-data-vertical/toggle-group:data-[spacing=0]:last:rounded-b-lg h-8 min-w-8 px-2.5 has-[&gt;svg]:px-2 group-data-[spacing=0]/toggle-group:has-[&gt;svg]:px-2 rounded-none shadow-none data-[variant=outline]:border-s-0 data-[variant=outline]:first:border-s data-[spacing=0]:first:rounded-s-lg data-[spacing=0]:last:rounded-e-lg" data-gsxui-slot-toggle-group-item data-gsxui-slot-toggle>B</button>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
	if strings.Contains(got, `role="radio"`) {
		t.Errorf("multiple-type item must not stamp role=radio\nin: %s", got)
	}
}

func TestToggleGroupItemOutlineVariant(t *testing.T) {
	got := render(t, ui.ToggleGroupItem("multiple", "outline", "", "", false, "bold", gsx.Raw("B"), nil))
	for _, want := range []string{
		`data-variant="outline"`,
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
	// Not a bare class="gap-8" match: ToggleGroup now carries its own
	// recipe class, so the caller's gap-8 merges in alongside it rather
	// than rendering as the entire attribute value.
	got := render(t, ui.ToggleGroup("multiple", "", "", "", gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "gap-8"}}))
	if strings.Count(got, "gap-8") != 1 {
		t.Errorf("caller class must merge in exactly once\nin: %s", got)
	}
	if strings.Count(got, "class=") != 1 {
		t.Errorf("expected exactly one class= attribute\nin: %s", got)
	}
}

func TestToggleGroupItemCallerClassMerges(t *testing.T) {
	// Not a bare class="px-8" match: ToggleGroupItem now carries its own
	// recipe class (including its own px-3), so the caller's px-8 merges
	// in (replacing the recipe's own px-3, same-property override)
	// alongside it.
	got := render(t, ui.ToggleGroupItem("multiple", "", "", "", false, "bold", gsx.Raw("B"), gsx.Attrs{{Key: "class", Value: "px-8"}}))
	if strings.Count(got, "px-8") != 1 {
		t.Errorf("caller class must merge in exactly once\nin: %s", got)
	}
	if strings.Count(got, "class=") != 1 {
		t.Errorf("expected exactly one class= attribute\nin: %s", got)
	}
}

func TestToggleGroupAttrsFallThrough(t *testing.T) {
	got := render(t, ui.ToggleGroup("multiple", "", "", "", gsx.Raw("x"), gsx.Attrs{{Key: "aria-label", Value: "Text formatting"}}))
	if !strings.Contains(got, `aria-label="Text formatting"`) {
		t.Errorf("missing aria-label\nin: %s", got)
	}
}
