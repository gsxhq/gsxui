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
	//
	// aria-invalid:border-destructive/ring-destructive and focus-visible:
	// ring-3 are restored structural chrome (Badge's own base recipe never
	// had a way to become visibly invalid before this fix — see the
	// style-porter report's "dark primitive states" entry); has-[>svg]:
	// px-1.5 replaces the dead has-data-[icon=…] form (Badge never stamps a
	// data-icon attribute) — see the report's "Systemic: has-data-[icon=…]"
	// entry.
	got := render(t, ui.Badge("", gsx.Raw("New"), nil))
	want := `<span data-variant="default" class="aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 focus-visible:ring-3 h-5 gap-1 rounded-4xl border border-transparent px-2 py-0.5 text-xs font-medium transition-all has-[&gt;svg]:px-1.5 [&amp;&gt;svg]:size-3 inline-flex items-center justify-center overflow-hidden whitespace-nowrap w-fit shrink-0 focus-visible:border-ring focus-visible:ring-ring/50 [&amp;&gt;svg]:pointer-events-none bg-primary text-primary-foreground [a]:hover:bg-primary/80" data-gsxui-slot-badge>New</span>`
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
