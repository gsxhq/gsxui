package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestSkeletonDefault(t *testing.T) {
	got := render(t, ui.Skeleton(nil))
	for _, want := range []string{
		`data-gsxui-slot-skeleton`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestSkeletonCallerClassIsForwardedOnce(t *testing.T) {
	got := render(t, ui.Skeleton(gsx.Attrs{{Key: "class", Value: "h-4 w-full rounded-md"}}))
	if strings.Count(got, `class="`) != 1 || !strings.Contains(got, "h-4") || !strings.Contains(got, "w-full") {
		t.Errorf("caller class must merge into the single class attribute and render once\nin: %s", got)
	}
}

func TestSkeletonPinned(t *testing.T) {
	// Exact full-render pin, verified token-by-token against shadcn's
	// Skeleton (registry/new-york-v4/ui/skeleton.tsx) — straight port, no
	// divergences. Now carries the recipe's resolved class (slot axis
	// migration); the markup shape is otherwise unchanged.
	got := render(t, ui.Skeleton(nil))
	want := `<div class="animate-pulse rounded-md bg-accent" data-gsxui-slot-skeleton></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestSkeletonAttrsFallThrough(t *testing.T) {
	got := render(t, ui.Skeleton(gsx.Attrs{{Key: "id", Value: "sk1"}, {Key: "aria-hidden", Value: "true"}}))
	for _, want := range []string{`id="sk1"`, `aria-hidden="true"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}
