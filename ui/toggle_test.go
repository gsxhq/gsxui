package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestToggleOffPinned(t *testing.T) {
	// Exact full-render pin for the zero-value (unpressed, default variant,
	// default size) render — token-by-token against toggleVariants' base +
	// default variant + default size (registry/new-york-v4/ui/toggle.tsx).
	got := render(t, ui.Toggle(false, "", "", gsx.Raw("Bold"), nil))
	want := `<button type="button" data-slot="toggle" data-gsxui-toggle data-variant="default" data-size="default" data-state="off" aria-pressed="false" class="inline-flex items-center justify-center gap-1 rounded-lg text-sm font-medium whitespace-nowrap transition-[color,box-shadow] outline-none hover:bg-muted hover:text-muted-foreground focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 data-[state=on]:bg-accent data-[state=on]:text-accent-foreground dark:aria-invalid:ring-destructive/40 [&amp;_svg]:pointer-events-none [&amp;_svg]:shrink-0 [&amp;_svg:not([class*=&#39;size-&#39;])]:size-4 bg-transparent h-8 min-w-8 px-2.5 has-[&gt;svg]:px-2">Bold</button>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestTogglePressedPinned(t *testing.T) {
	// Exact full-render pin for pressed={true} — the server-visible initial
	// "on" state (aria-pressed="true" data-state="on"), no click required.
	got := render(t, ui.Toggle(true, "", "", gsx.Raw("Bold"), nil))
	want := `<button type="button" data-slot="toggle" data-gsxui-toggle data-variant="default" data-size="default" data-state="on" aria-pressed="true" class="inline-flex items-center justify-center gap-1 rounded-lg text-sm font-medium whitespace-nowrap transition-[color,box-shadow] outline-none hover:bg-muted hover:text-muted-foreground focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 data-[state=on]:bg-accent data-[state=on]:text-accent-foreground dark:aria-invalid:ring-destructive/40 [&amp;_svg]:pointer-events-none [&amp;_svg]:shrink-0 [&amp;_svg:not([class*=&#39;size-&#39;])]:size-4 bg-transparent h-8 min-w-8 px-2.5 has-[&gt;svg]:px-2">Bold</button>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestToggleOutlineVariant(t *testing.T) {
	got := render(t, ui.Toggle(false, "outline", "", gsx.Raw("x"), nil))
	for _, want := range []string{
		`data-variant="outline"`,
		"border border-input bg-transparent hover:bg-accent hover:text-accent-foreground",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestToggleSizes(t *testing.T) {
	sm := render(t, ui.Toggle(false, "", "sm", gsx.Raw("x"), nil))
	for _, want := range []string{`data-size="sm"`, "h-7 min-w-7 rounded-[min(var(--radius-md),12px)] px-2.5 has-[&gt;svg]:px-1.5 text-[0.8rem]"} {
		if !strings.Contains(sm, want) {
			t.Errorf("sm missing %q\nin: %s", want, sm)
		}
	}

	lg := render(t, ui.Toggle(false, "", "lg", gsx.Raw("x"), nil))
	for _, want := range []string{`data-size="lg"`, "h-9 min-w-9 px-2.5 has-[&gt;svg]:px-2"} {
		if !strings.Contains(lg, want) {
			t.Errorf("lg missing %q\nin: %s", want, lg)
		}
	}

	def := render(t, ui.Toggle(false, "", "default", gsx.Raw("x"), nil))
	for _, want := range []string{`data-size="default"`, "h-8 min-w-8 px-2.5 has-[&gt;svg]:px-2"} {
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
	got := render(t, ui.Toggle(false, "", "", gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "h-12"}}))
	if strings.Contains(got, "h-8") {
		t.Errorf("caller h-12 must drop default h-8\nin: %s", got)
	}
	if !strings.Contains(got, "h-12") || !strings.Contains(got, "inline-flex") {
		t.Errorf("want h-12 plus surviving structural classes\nin: %s", got)
	}
}
