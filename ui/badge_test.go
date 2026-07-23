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
		`data-slot="badge"`,
		`data-variant="default"`,
		"inline-flex w-fit shrink-0 items-center justify-center",
		"bg-primary text-primary-foreground",
		">New</span>",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestBadgeVariants(t *testing.T) {
	cases := map[string]string{
		"secondary":   "bg-secondary text-secondary-foreground",
		"destructive": "bg-destructive text-white",
		"outline":     "border-border text-foreground",
		"ghost":       "[a&amp;]:hover:bg-accent",
		"link":        "text-primary underline-offset-4",
	}
	for variant, wantClass := range cases {
		got := render(t, ui.Badge(variant, gsx.Raw("x"), nil))
		if !strings.Contains(got, wantClass) {
			t.Errorf("variant %s: missing %q\nin: %s", variant, wantClass, got)
		}
		if !strings.Contains(got, `data-variant="`+variant+`"`) {
			t.Errorf("variant %s: missing data-variant stamp\nin: %s", variant, got)
		}
	}
}

func TestBadgeCallerClassMerges(t *testing.T) {
	// Caller px-4 must WIN over base px-2 via tailwind-merge — and base
	// structural classes must survive. The core customization semantic.
	got := render(t, ui.Badge("", gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "px-4"}}))
	if strings.Contains(got, "px-2") {
		t.Errorf("base px-2 should be dropped by caller px-4\nin: %s", got)
	}
	for _, want := range []string{"px-4", "rounded-full"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestBadgePinned(t *testing.T) {
	// Exact full-render pin, verified token-by-token against shadcn's
	// badgeVariants base + default variant (registry/new-york-v4/ui/badge.tsx)
	// and docs/jsx-parity.md — no ADAPT deviations apply to the default badge.
	got := render(t, ui.Badge("", gsx.Raw("New"), nil))
	want := `<span data-slot="badge" data-variant="default" class="inline-flex w-fit shrink-0 items-center justify-center gap-1 overflow-hidden rounded-full border border-transparent px-2 py-0.5 text-xs font-medium whitespace-nowrap transition-[color,box-shadow] focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 [&amp;&gt;svg]:pointer-events-none [&amp;&gt;svg]:size-3 bg-primary text-primary-foreground [a&amp;]:hover:bg-primary/90">New</span>`
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
