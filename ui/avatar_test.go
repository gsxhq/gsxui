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
		`data-slot="avatar"`,
		`data-slot="avatar-image"`,
		`data-gsxui-avatar-image`,
		`src="/broken.jpg"`,
		`alt="shadcn"`,
		`data-slot="avatar-fallback"`, ">CN<",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
	// ADAPT: fallback renders with no `hidden` attribute — load state isn't
	// known server-side; JS (avatar.js) toggles display on load/error.
	fallbackStart := strings.Index(got, `data-slot="avatar-fallback"`)
	fallbackTag := got[fallbackStart : fallbackStart+strings.Index(got[fallbackStart:], ">")]
	if strings.Contains(fallbackTag, "hidden") {
		t.Errorf("fallback must not render with a hidden attribute\nin: %s", fallbackTag)
	}
}

func TestAvatarCallerClassMerges(t *testing.T) {
	got := render(t, ui.Avatar(gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "size-12"}}))
	if strings.Contains(got, "size-8") {
		t.Errorf("base size-8 should be dropped by caller size-12\nin: %s", got)
	}
	if !strings.Contains(got, "size-12") {
		t.Errorf("missing caller class size-12\nin: %s", got)
	}
}

func TestAvatarPinned(t *testing.T) {
	// Exact full-render pin for Avatar > AvatarImage, verified token-by-token
	// against shadcn's Avatar/AvatarImage (registry/new-york-v4/ui/avatar.tsx)
	// and docs/jsx-parity.md's ADAPT: AvatarImage adds absolute inset-0
	// to overlay the fallback (no-JS rendering correct).
	got := render(t, ui.AvatarImage("/shadcn.jpg", "shadcn", nil))
	want := `<img data-slot="avatar-image" data-gsxui-avatar-image src="/shadcn.jpg" alt="shadcn" class="aspect-square size-full absolute inset-0">`
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
