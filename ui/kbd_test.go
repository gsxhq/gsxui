package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestKbdDefault(t *testing.T) {
	got := render(t, ui.Kbd(gsx.Raw("Ctrl"), nil))
	for _, want := range []string{
		`<kbd data-gsxui-slot="kbd"`,
		">Ctrl</kbd>",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestKbdCallerClassIsForwardedOnce(t *testing.T) {
	got := render(t, ui.Kbd(nil, gsx.Attrs{{Key: "class", Value: "h-8"}}))
	if strings.Count(got, `class="h-8"`) != 1 {
		t.Errorf("caller class must be the only class and render once\nin: %s", got)
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
	want := `<kbd data-gsxui-slot="kbd">Ctrl</kbd>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestKbdGroupDefault(t *testing.T) {
	got := render(t, ui.KbdGroup(gsx.Raw("<kbd>Ctrl</kbd><kbd>C</kbd>"), nil))
	want := `<kbd data-gsxui-slot="kbd-group"><kbd>Ctrl</kbd><kbd>C</kbd></kbd>`
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
	want := `<kbd data-gsxui-slot="kbd-group"><kbd>Ctrl</kbd></kbd>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}
