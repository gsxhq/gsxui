package button_test

import (
	"context"
	"regexp"
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui/button"
)

// disabledAttr matches the bare boolean `disabled` HTML attribute, as
// distinct from Tailwind's `disabled:pointer-events-none` / `disabled:opacity-50`
// variant classes that appear verbatim in the button's base class string and
// would otherwise false-positive a plain strings.Contains(got, "disabled").
var disabledAttr = regexp.MustCompile(`disabled(>|\s)`)

func render(t *testing.T, n gsx.Node) string {
	t.Helper()
	var sb strings.Builder
	if err := n.Render(context.Background(), &sb); err != nil {
		t.Fatal(err)
	}
	return sb.String()
}

func TestButtonDefault(t *testing.T) {
	got := render(t, button.Button("", "", "", false, gsx.Raw("Save"), nil))
	for _, want := range []string{
		"<button", `data-slot="button"`, `type="button"`,
		`data-variant="default"`, `data-size="default"`,
		"bg-primary text-primary-foreground", "h-9 px-4 py-2",
		">Save</button>",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
	if disabledAttr.MatchString(got) {
		t.Errorf("unexpected disabled attr\nin: %s", got)
	}
}

func TestButtonPinned(t *testing.T) {
	// Exact full-render pin, verified token-by-token against shadcn's
	// buttonVariants base + default variant + default size
	// (registry/new-york-v4/ui/button.tsx) and docs/jsx-parity.md — no ADAPT
	// deviations apply to the default button.
	got := render(t, button.Button("", "", "", false, gsx.Raw("Save"), nil))
	want := `<button data-slot="button" data-variant="default" data-size="default" type="button" class="inline-flex shrink-0 items-center justify-center gap-2 rounded-md text-sm font-medium whitespace-nowrap transition-all outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 [&amp;_svg]:pointer-events-none [&amp;_svg]:shrink-0 [&amp;_svg:not([class*=&#39;size-&#39;])]:size-4 bg-primary text-primary-foreground hover:bg-primary/90 h-9 px-4 py-2 has-[&gt;svg]:px-3">Save</button>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestButtonVariantSize(t *testing.T) {
	got := render(t, button.Button("outline", "sm", "", false, gsx.Raw("x"), nil))
	for _, want := range []string{
		"border bg-background shadow-xs", "h-8 gap-1.5 rounded-md px-3",
		`data-variant="outline"`, `data-size="sm"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestButtonHrefRendersAnchor(t *testing.T) {
	got := render(t, button.Button("", "", "/docs", false, gsx.Raw("Docs"), nil))
	for _, want := range []string{"<a", `href="/docs"`, `data-slot="button"`, ">Docs</a>"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
	if strings.Contains(got, "<button") {
		t.Errorf("href should render <a>, not <button>\nin: %s", got)
	}
}

func TestButtonDisabled(t *testing.T) {
	// disabled wins over href: render a real disabled <button>.
	got := render(t, button.Button("", "", "/docs", true, gsx.Raw("x"), nil))
	if !strings.Contains(got, "<button") {
		t.Errorf("want disabled <button>\nin: %s", got)
	}
	if !disabledAttr.MatchString(got) {
		t.Errorf("want real disabled attribute\nin: %s", got)
	}
}

func TestButtonTypeIsOverridableDefault(t *testing.T) {
	// type="button" is authored BEFORE { attrs... }: caller type=submit wins.
	got := render(t, button.Button("", "", "", false, gsx.Raw("Go"), gsx.Attrs{{Key: "type", Value: "submit"}}))
	if !strings.Contains(got, `type="submit"`) {
		t.Errorf("caller type=submit must override default\nin: %s", got)
	}
	if strings.Contains(got, `type="button"`) {
		t.Errorf("default type should be replaced, not duplicated\nin: %s", got)
	}
}

func TestButtonCallerClassMerges(t *testing.T) {
	got := render(t, button.Button("", "", "", false, gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "h-12"}}))
	if strings.Contains(got, "h-9") {
		t.Errorf("caller h-12 must drop default h-9\nin: %s", got)
	}
	if !strings.Contains(got, "h-12") || !strings.Contains(got, "inline-flex") {
		t.Errorf("want h-12 plus surviving structural classes\nin: %s", got)
	}
}
