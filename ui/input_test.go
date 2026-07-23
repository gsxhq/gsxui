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
		"<input", `data-slot="input"`, `type="text"`,
		"h-9 w-full min-w-0 rounded-md border border-input",
		"focus-visible:border-ring focus-visible:ring-[3px]",
		"aria-invalid:border-destructive aria-invalid:ring-destructive/20",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestInputPinned(t *testing.T) {
	// Exact full-render pin, verified token-by-token against shadcn's Input
	// (registry/new-york-v4/ui/input.tsx) and docs/jsx-parity.md — straight
	// port, no ADAPT deviations.
	got := render(t, ui.Input(nil))
	want := `<input data-slot="input" type="text" class="h-9 w-full min-w-0 rounded-md border border-input bg-transparent px-3 py-1 text-base shadow-xs transition-[color,box-shadow] outline-none selection:bg-primary selection:text-primary-foreground file:inline-flex file:h-7 file:border-0 file:bg-transparent file:text-sm file:font-medium file:text-foreground placeholder:text-muted-foreground disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50 md:text-sm dark:bg-input/30 focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40"/>`
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
	if strings.Contains(got, "h-9") {
		t.Errorf("caller h-12 must drop default h-9\nin: %s", got)
	}
	for _, want := range []string{"h-12", "rounded-md"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
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
