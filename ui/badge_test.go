package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestBadgeDefault(t *testing.T) {
	got := render(t, ui.Badge("", gsx.Raw("New"), nil))
	for _, want := range []string{
		`data-gsxui-slot-badge`,
		`data-variant="default"`,
		">New</span>",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestBadgeVariants(t *testing.T) {
	for _, variant := range []string{"default", "secondary", "destructive", "outline", "ghost", "link"} {
		input := variant
		if variant == "default" {
			input = ""
		}
		got := render(t, ui.Badge(input, gsx.Raw("x"), nil))
		if !strings.Contains(got, `data-variant="`+variant+`"`) {
			t.Errorf("variant %s: missing data-variant stamp\nin: %s", variant, got)
		}
	}
}

func TestBadgeCallerClassIsForwardedOnce(t *testing.T) {
	got := render(t, ui.Badge("", gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "px-4"}}))
	if strings.Count(got, `class="`) != 1 || !strings.Contains(got, "px-4") {
		t.Errorf("caller class must merge into the single class attribute and render once\nin: %s", got)
	}
}

func TestBadgePinned(t *testing.T) {
	// Exact full-render pin, verified token-by-token against shadcn's
	// badgeVariants base + default variant (registry/new-york-v4/ui/badge.tsx)
	// and docs/jsx-parity.md — no ADAPT deviations apply to the default badge.
	got := render(t, ui.Badge("", gsx.Raw("New"), nil))
	want := `<span data-variant="default" class="inline-flex h-5 w-fit shrink-0 items-center justify-center gap-1 overflow-hidden rounded-4xl border border-transparent px-2 py-0.5 text-xs font-medium whitespace-nowrap transition-[color,box-shadow] has-[&gt;svg]:px-1.5 focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 [&amp;&gt;svg]:size-3 [&amp;&gt;svg]:pointer-events-none bg-primary text-primary-foreground [a&amp;]:hover:bg-primary/90" data-gsxui-slot-badge>New</span>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestBadgeAttrsFallThrough(t *testing.T) {
	got := render(t, ui.Badge("", gsx.Raw("x"), gsx.Attrs{{Key: "id", Value: "b1"}, {Key: "aria-label", Value: "badge"}}))
	for _, want := range []string{`id="b1"`, `aria-label="badge"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}
