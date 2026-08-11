package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestInputDefault(t *testing.T) {
	got := render(t, ui.Input(nil))
	for _, want := range []string{
		"<input", `type="text"`, `data-gsxui-slot-input`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestInputPinned(t *testing.T) {
	// Presentation lives in the stylesheet; the render pin covers structure.
	// The leading structural baseline (w-full through disabled:opacity-50)
	// is restored "carried: no upstream counterpart" content — see the
	// style-porter report's "Input/Textarea structural baseline" entry.
	got := render(t, ui.Input(nil))
	want := `<input type="text" class="w-full min-w-0 outline-none file:inline-flex file:border-0 file:bg-transparent file:text-foreground placeholder:text-muted-foreground disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50 dark:bg-input/30 border-input focus-visible:border-ring focus-visible:ring-ring/50 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive dark:aria-invalid:border-destructive/50 disabled:bg-input/50 dark:disabled:bg-input/80 h-8 rounded-lg border bg-transparent px-2.5 py-1 text-base transition-colors file:h-6 file:text-sm file:font-medium focus-visible:ring-3 aria-invalid:ring-3 md:text-sm" data-gsxui-slot-input>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestInputTypeOverridable(t *testing.T) {
	// type="text" is authored BEFORE { attrs... }: caller type="email" wins,
	// not duplicated.
	got := render(t, ui.Input(gsx.Attrs{{Key: "type", Value: "email"}}))
	if !strings.Contains(got, `type="email"`) {
		t.Errorf("caller type=email must override default\nin: %s", got)
	}
	if strings.Contains(got, `type="text"`) {
		t.Errorf("default type should be replaced, not duplicated\nin: %s", got)
	}
}

func TestInputCallerClassMerges(t *testing.T) {
	got := render(t, ui.Input(gsx.Attrs{{Key: "class", Value: "h-12"}}))
	if !strings.Contains(got, "h-12") {
		t.Errorf("caller class must be forwarded\nin: %s", got)
	}
	if strings.Contains(got, "h-8") {
		t.Errorf("caller's h-12 must win over the recipe's h-8\nin: %s", got)
	}
}

func TestInputAttrsFallThrough(t *testing.T) {
	got := render(t, ui.Input(gsx.Attrs{{Key: "id", Value: "email"}, {Key: "placeholder", Value: "you@example.com"}, {Key: "disabled", Value: true}}))
	for _, want := range []string{`id="email"`, `placeholder="you@example.com"`, "disabled"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}
