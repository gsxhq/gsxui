package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestAvatarStructure(t *testing.T) {
	got := render(t, ui.Avatar(gsx.Fragment(
		ui.AvatarImage("/broken.jpg", "shadcn", nil),
		ui.AvatarFallback(gsx.Raw("CN"), nil),
	), nil))
	for _, want := range []string{
		`data-gsxui-slot-avatar`,
		`data-gsxui-slot-avatar-image`,
		`src="/broken.jpg"`,
		`alt="shadcn"`,
		`data-gsxui-slot-avatar-fallback`, ">CN<",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
	// ADAPT: fallback renders with no `hidden` attribute — load state isn't
	// known server-side; JS (avatar.js) toggles display on load/error.
	fallbackStart := strings.Index(got, `data-gsxui-slot-avatar-fallback`)
	if fallbackStart < 0 {
		t.Fatal("missing avatar fallback slot")
	}
	fallbackTag := got[fallbackStart : fallbackStart+strings.Index(got[fallbackStart:], ">")]
	if strings.Contains(fallbackTag, "hidden") {
		t.Errorf("fallback must not render with a hidden attribute\nin: %s", fallbackTag)
	}
}

func TestAvatarCallerClassIsForwardedOnce(t *testing.T) {
	got := render(t, ui.Avatar(gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "size-12"}}))
	if strings.Count(got, `class="`) != 1 || !strings.Contains(got, "size-12") {
		t.Errorf("caller class must merge into the single class attribute and render once\nin: %s", got)
	}
}

func TestAvatarPinned(t *testing.T) {
	// Exact full-render pin for Avatar > AvatarImage, verified token-by-token
	// against shadcn's Avatar/AvatarImage (registry/new-york-v4/ui/avatar.tsx)
	// and docs/jsx-parity.md's ADAPT: AvatarImage adds absolute inset-0
	// to overlay the fallback (no-JS rendering correct). Now carries the
	// recipe's resolved class (slot axis migration).
	got := render(t, ui.AvatarImage("/shadcn.jpg", "shadcn", nil))
	want := `<img src="/shadcn.jpg" alt="shadcn" class="rounded-full" data-gsxui-slot-avatar-image>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestAvatarAttrsFallThrough(t *testing.T) {
	got := render(t, ui.Avatar(gsx.Raw("x"), gsx.Attrs{{Key: "id", Value: "av1"}, {Key: "aria-label", Value: "profile"}}))
	for _, want := range []string{`id="av1"`, `aria-label="profile"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}
