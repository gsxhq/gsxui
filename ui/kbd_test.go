package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestKbdDefault(t *testing.T) {
	// The class string carries Tailwind arbitrary-variant selectors built
	// from "&" and "'" — both are HTML-escaped (&amp;, &#39;) in the
	// rendered attribute value, the same ambiguous-ampersand/attribute
	// escaping accordion's ported "[[open]>summary_&]:rotate-180" selector
	// goes through (see accordion_test.go). Browsers decode the entity back
	// to the literal character when parsing the attribute, so Tailwind's
	// selector matching is unaffected.
	got := render(t, ui.Kbd(gsx.Raw("Ctrl"), nil))
	for _, want := range []string{
		`<kbd data-slot="kbd"`,
		"pointer-events-none inline-flex h-5 w-fit min-w-5 items-center justify-center gap-1 rounded-sm bg-muted px-1 font-sans text-xs font-medium text-muted-foreground select-none",
		"[&amp;_svg:not([class*=&#39;size-&#39;])]:size-3",
		"[[data-slot=tooltip-content]_&amp;]:bg-background/20 [[data-slot=tooltip-content]_&amp;]:text-background dark:[[data-slot=tooltip-content]_&amp;]:bg-background/10",
		">Ctrl</kbd>",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

// Caller class merges via tailwind-merge — h-5 (height utility) must be
// replaced, not just appended alongside, by a caller-supplied h-8.
func TestKbdCallerClassMerges(t *testing.T) {
	got := render(t, ui.Kbd(nil, gsx.Attrs{{Key: "class", Value: "h-8"}}))
	if strings.Contains(got, "h-5") {
		t.Errorf("base h-5 should be dropped by caller h-8\nin: %s", got)
	}
	if !strings.Contains(got, "h-8") {
		t.Errorf("missing caller class h-8\nin: %s", got)
	}
}

func TestKbdAttrsFallThrough(t *testing.T) {
	got := render(t, ui.Kbd(nil, gsx.Attrs{{Key: "id", Value: "k1"}, {Key: "aria-label", Value: "control"}}))
	for _, want := range []string{`id="k1"`, `aria-label="control"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestKbdPinned(t *testing.T) {
	// Exact full-render pin. Token-for-token against shadcn's Kbd
	// (registry/new-york-v4/ui/kbd.tsx) — no dropped tokens.
	got := render(t, ui.Kbd(gsx.Raw("Ctrl"), nil))
	want := `<kbd data-slot="kbd" class="pointer-events-none inline-flex h-5 w-fit min-w-5 items-center justify-center gap-1 rounded-sm bg-muted px-1 font-sans text-xs font-medium text-muted-foreground select-none [&amp;_svg:not([class*=&#39;size-&#39;])]:size-3 [[data-slot=tooltip-content]_&amp;]:bg-background/20 [[data-slot=tooltip-content]_&amp;]:text-background dark:[[data-slot=tooltip-content]_&amp;]:bg-background/10">Ctrl</kbd>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestKbdGroupDefault(t *testing.T) {
	got := render(t, ui.KbdGroup(gsx.Raw("<kbd>Ctrl</kbd><kbd>C</kbd>"), nil))
	want := `<kbd data-slot="kbd-group" class="inline-flex items-center gap-1"><kbd>Ctrl</kbd><kbd>C</kbd></kbd>`
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func TestKbdGroupAttrsFallThrough(t *testing.T) {
	got := render(t, ui.KbdGroup(nil, gsx.Attrs{{Key: "id", Value: "g1"}}))
	if !strings.Contains(got, `id="g1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

func TestKbdGroupPinned(t *testing.T) {
	got := render(t, ui.KbdGroup(gsx.Raw("<kbd>Ctrl</kbd>"), nil))
	want := `<kbd data-slot="kbd-group" class="inline-flex items-center gap-1"><kbd>Ctrl</kbd></kbd>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}
