package separator_test

import (
	"context"
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui/separator"
)

func render(t *testing.T, n gsx.Node) string {
	t.Helper()
	var sb strings.Builder
	if err := n.Render(context.Background(), &sb); err != nil {
		t.Fatal(err)
	}
	return sb.String()
}

func TestSeparatorDefault(t *testing.T) {
	got := render(t, separator.Separator("", nil))
	for _, want := range []string{
		`data-slot="separator"`,
		`role="none"`,
		`data-orientation="horizontal"`,
		"shrink-0 bg-border",
		"data-[orientation=vertical]:w-px",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestSeparatorVertical(t *testing.T) {
	got := render(t, separator.Separator("vertical", nil))
	if !strings.Contains(got, `data-orientation="vertical"`) {
		t.Errorf("missing vertical data-orientation stamp\nin: %s", got)
	}
}

func TestSeparatorCallerClassMerges(t *testing.T) {
	got := render(t, separator.Separator("", gsx.Attrs{{Key: "class", Value: "bg-red-500"}}))
	if strings.Contains(got, "bg-border") {
		t.Errorf("base bg-border should be dropped by caller bg-red-500\nin: %s", got)
	}
	if !strings.Contains(got, "bg-red-500") {
		t.Errorf("missing caller class bg-red-500\nin: %s", got)
	}
}

func TestSeparatorPinned(t *testing.T) {
	// Exact full-render pin, verified token-by-token against shadcn's
	// Separator (registry/new-york-v4/ui/separator.tsx) and
	// docs/jsx-parity.md — decorative is dropped (ADAPT), role="none" always.
	got := render(t, separator.Separator("", nil))
	want := `<div data-slot="separator" role="none" data-orientation="horizontal" class="shrink-0 bg-border data-[orientation=horizontal]:h-px data-[orientation=horizontal]:w-full data-[orientation=vertical]:h-full data-[orientation=vertical]:w-px"/>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestSeparatorAttrsFallThrough(t *testing.T) {
	got := render(t, separator.Separator("", gsx.Attrs{{Key: "id", Value: "s1"}, {Key: "aria-hidden", Value: "true"}}))
	for _, want := range []string{`id="s1"`, `aria-hidden="true"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}
