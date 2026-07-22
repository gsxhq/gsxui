package badge_test

import (
	"context"
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui/badge"
)

func render(t *testing.T, n gsx.Node) string {
	t.Helper()
	var sb strings.Builder
	if err := n.Render(context.Background(), &sb); err != nil {
		t.Fatal(err)
	}
	return sb.String()
}

func TestBadgeDefault(t *testing.T) {
	got := render(t, badge.Badge("", gsx.Raw("New"), nil))
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
		got := render(t, badge.Badge(variant, gsx.Raw("x"), nil))
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
	got := render(t, badge.Badge("", gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "px-4"}}))
	if strings.Contains(got, "px-2") {
		t.Errorf("base px-2 should be dropped by caller px-4\nin: %s", got)
	}
	for _, want := range []string{"px-4", "rounded-full"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestBadgeAttrsFallThrough(t *testing.T) {
	got := render(t, badge.Badge("", gsx.Raw("x"), gsx.Attrs{{Key: "id", Value: "b1"}, {Key: "aria-label", Value: "badge"}}))
	for _, want := range []string{`id="b1"`, `aria-label="badge"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}
