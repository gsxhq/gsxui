package icon_test

import (
	"context"
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui/icon"
)

func render(t *testing.T, n gsx.Node) (string, error) {
	t.Helper()
	var sb strings.Builder
	err := n.Render(context.Background(), &sb)
	return sb.String(), err
}

func TestIconRenders(t *testing.T) {
	// Exact full-render pin, verified against the generated icon_data.go
	// entry for "chevron-down" and the svg attribute list from the task
	// brief (Lucide 24x24 stroke defaults + gsxui's data-slot/aria-hidden/
	// size-4 additions).
	got, err := render(t, icon.ChevronDown())
	if err != nil {
		t.Fatal(err)
	}
	want := `<svg data-slot="icon" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="size-4"><path d="m6 9 6 6 6-6"/></svg>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestIconUnknownErrors(t *testing.T) {
	_, err := render(t, icon.New("nope")())
	if err == nil {
		t.Fatal("want error for unknown icon name, got nil")
	}
	if !strings.Contains(err.Error(), "unknown icon") {
		t.Errorf("error %q does not contain %q", err.Error(), "unknown icon")
	}
}

func TestIconClassMerge(t *testing.T) {
	// Caller size-6 must WIN over default size-4 via tailwind-merge — the
	// base size-4 token must be dropped, not just appended alongside.
	got, err := render(t, icon.ChevronDown(gsx.Attr{Key: "class", Value: "size-6"}))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(got, "size-4") {
		t.Errorf("base size-4 should be dropped by caller size-6\nin: %s", got)
	}
	if !strings.Contains(got, `class="size-6"`) {
		t.Errorf("missing caller class size-6\nin: %s", got)
	}
}

func TestIconAttrsFallThrough(t *testing.T) {
	got, err := render(t, icon.ChevronDown(gsx.Attr{Key: "id", Value: "i1"}, gsx.Attr{Key: "aria-label", Value: "expand"}))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`id="i1"`, `aria-label="expand"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestIconAriaHiddenOverridable(t *testing.T) {
	// A caller-supplied aria-hidden must win over the default (positional
	// precedence: the literal aria-hidden="true" is authored before
	// { attrs... } in icon.gsx).
	got, err := render(t, icon.ChevronDown(gsx.Attr{Key: "aria-hidden", Value: "false"}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, `aria-hidden="false"`) {
		t.Errorf("caller aria-hidden=false should win\nin: %s", got)
	}
	if strings.Contains(got, `aria-hidden="true"`) {
		t.Errorf("default aria-hidden=true should not also render\nin: %s", got)
	}
}
