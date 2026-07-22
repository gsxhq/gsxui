package avatar_test

import (
	"context"
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui/avatar"
)

func render(t *testing.T, n gsx.Node) string {
	t.Helper()
	var sb strings.Builder
	if err := n.Render(context.Background(), &sb); err != nil {
		t.Fatal(err)
	}
	return sb.String()
}

func TestAvatarStructure(t *testing.T) {
	got := render(t, avatar.Avatar(gsx.Fragment(
		avatar.AvatarImage("/broken.jpg", "shadcn", nil),
		avatar.AvatarFallback(gsx.Raw("CN"), nil),
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
	got := render(t, avatar.Avatar(gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "size-12"}}))
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
	// and docs/jsx-parity.md's ADAPT: data-gsxui-avatar-image replaces
	// Radix's load-state context.
	got := render(t, avatar.AvatarImage("/shadcn.jpg", "shadcn", nil))
	want := `<img data-slot="avatar-image" data-gsxui-avatar-image src="/shadcn.jpg" alt="shadcn" class="aspect-square size-full"/>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestAvatarAttrsFallThrough(t *testing.T) {
	got := render(t, avatar.Avatar(gsx.Raw("x"), gsx.Attrs{{Key: "id", Value: "av1"}, {Key: "aria-label", Value: "profile"}}))
	for _, want := range []string{`id="av1"`, `aria-label="profile"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}
