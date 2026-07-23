package skeleton_test

import (
	"context"
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui/skeleton"
)

func render(t *testing.T, n gsx.Node) string {
	t.Helper()
	var sb strings.Builder
	if err := n.Render(context.Background(), &sb); err != nil {
		t.Fatal(err)
	}
	return sb.String()
}

func TestSkeletonDefault(t *testing.T) {
	got := render(t, skeleton.Skeleton(nil))
	for _, want := range []string{
		`data-slot="skeleton"`,
		"animate-pulse rounded-md bg-accent",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestSkeletonCallerClassMerges(t *testing.T) {
	got := render(t, skeleton.Skeleton(gsx.Attrs{{Key: "class", Value: "h-4 w-full rounded-md"}}))
	for _, want := range []string{"h-4 w-full rounded-md", "animate-pulse"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestSkeletonPinned(t *testing.T) {
	// Exact full-render pin, verified token-by-token against shadcn's
	// Skeleton (registry/new-york-v4/ui/skeleton.tsx) — straight port, no
	// divergences.
	got := render(t, skeleton.Skeleton(nil))
	want := `<div data-slot="skeleton" class="animate-pulse rounded-md bg-accent"></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestSkeletonAttrsFallThrough(t *testing.T) {
	got := render(t, skeleton.Skeleton(gsx.Attrs{{Key: "id", Value: "sk1"}, {Key: "aria-hidden", Value: "true"}}))
	for _, want := range []string{`id="sk1"`, `aria-hidden="true"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}
