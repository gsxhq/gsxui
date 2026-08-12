package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestToggleOffPinned(t *testing.T) {
	// Exact full-render pin for the zero-value (unpressed, default variant,
	// default size) render. Toggle is migrated onto the slot axis: its class
	// attribute now carries the recipe classes resolved from
	// registry/canonical/toggle.gsx's toggle.Root()/Variant()/Size()
	// accessor calls, rather than being absent as it was pre-migration. The
	// computed presentation is unchanged — verified by the sweep (make
	// sweep-compare) — only the class attribute's literal contents changed.
	// has-data-[icon=…] (dead selector) replaced with has-[>svg]: — see the
	// style-porter report's "Systemic: has-data-[icon=…]" entry.
	got := render(t, ui.Toggle(false, "", "", gsx.Raw("Bold"), nil))
	want := `<button type="button" data-variant="default" data-size="default" data-state="off" aria-pressed="false" class="hover:text-foreground aria-pressed:bg-muted focus-visible:border-ring focus-visible:ring-ring/50 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive disabled:pointer-events-none disabled:opacity-50 data-[state=on]:bg-muted gap-1 rounded-lg text-sm font-medium transition-all [&amp;_svg:not([class*=&#39;size-&#39;])]:size-4 inline-flex items-center justify-center whitespace-nowrap outline-none focus-visible:ring-[3px] [&amp;_svg]:pointer-events-none [&amp;_svg]:shrink-0 bg-transparent h-8 min-w-8 px-2.5 has-[&gt;svg]:px-2" data-gsxui-slot-toggle>Bold</button>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestTogglePressedPinned(t *testing.T) {
	// Exact full-render pin for pressed={true} — the server-visible initial
	// "on" state (aria-pressed="true" data-state="on"), no click required.
	// See TestToggleOffPinned for why the class attribute is now present.
	got := render(t, ui.Toggle(true, "", "", gsx.Raw("Bold"), nil))
	want := `<button type="button" data-variant="default" data-size="default" data-state="on" aria-pressed="true" class="hover:text-foreground aria-pressed:bg-muted focus-visible:border-ring focus-visible:ring-ring/50 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive disabled:pointer-events-none disabled:opacity-50 data-[state=on]:bg-muted gap-1 rounded-lg text-sm font-medium transition-all [&amp;_svg:not([class*=&#39;size-&#39;])]:size-4 inline-flex items-center justify-center whitespace-nowrap outline-none focus-visible:ring-[3px] [&amp;_svg]:pointer-events-none [&amp;_svg]:shrink-0 bg-transparent h-8 min-w-8 px-2.5 has-[&gt;svg]:px-2" data-gsxui-slot-toggle>Bold</button>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestToggleOutlineVariant(t *testing.T) {
	got := render(t, ui.Toggle(false, "outline", "", gsx.Raw("x"), nil))
	for _, want := range []string{
		`data-variant="outline"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestToggleSizes(t *testing.T) {
	sm := render(t, ui.Toggle(false, "", "sm", gsx.Raw("x"), nil))
	for _, want := range []string{`data-size="sm"`} {
		if !strings.Contains(sm, want) {
			t.Errorf("sm missing %q\nin: %s", want, sm)
		}
	}

	lg := render(t, ui.Toggle(false, "", "lg", gsx.Raw("x"), nil))
	for _, want := range []string{`data-size="lg"`} {
		if !strings.Contains(lg, want) {
			t.Errorf("lg missing %q\nin: %s", want, lg)
		}
	}

	def := render(t, ui.Toggle(false, "", "default", gsx.Raw("x"), nil))
	for _, want := range []string{`data-size="default"`} {
		if !strings.Contains(def, want) {
			t.Errorf("default missing %q\nin: %s", want, def)
		}
	}
}

func TestToggleDisabledFallsThrough(t *testing.T) {
	// disabled is not a declared param — it flows through attrs like any
	// other plain boolean HTML attribute (no href/disabled interplay to
	// resolve server-side the way Button has).
	got := render(t, ui.Toggle(false, "", "", gsx.Raw("x"), gsx.Attrs{{Key: "disabled", Value: true}}))
	if !strings.Contains(got, "disabled") {
		t.Errorf("want disabled attribute\nin: %s", got)
	}
}

func TestToggleAttrsFallThrough(t *testing.T) {
	got := render(t, ui.Toggle(false, "", "", gsx.Raw("x"), gsx.Attrs{{Key: "aria-label", Value: "Toggle bold"}}))
	if !strings.Contains(got, `aria-label="Toggle bold"`) {
		t.Errorf("missing aria-label\nin: %s", got)
	}
}

func TestToggleCallerClassMerges(t *testing.T) {
	// The recipe now also emits a class attribute, so the caller's "h-12" is
	// merged into it rather than being the attribute's sole content — assert
	// containment of the token, not an exact `class="h-12"` attribute, the
	// same treatment TestCardComposition gives Card's caller class.
	got := render(t, ui.Toggle(false, "", "", gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "h-12"}}))
	if strings.Count(got, "h-12") != 1 {
		t.Errorf("caller class must be forwarded exactly once\nin: %s", got)
	}
	if strings.Count(got, `class="`) != 1 {
		t.Errorf("expected exactly one class attribute\nin: %s", got)
	}
}
