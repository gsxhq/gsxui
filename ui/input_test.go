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
		"<input", `type="text"`, `data-gsxui-slot="input"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestInputPinned(t *testing.T) {
	// Presentation lives in the stylesheet; the render pin covers structure.
	got := render(t, ui.Input(nil))
	want := `<input type="text" data-gsxui-slot="input">`
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
	if strings.Count(got, `class="h-12"`) != 1 {
		t.Errorf("caller class must be forwarded exactly once\nin: %s", got)
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
