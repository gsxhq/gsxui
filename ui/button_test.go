package ui_test

import (
	"regexp"
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

// disabledAttr matches the bare boolean `disabled` HTML attribute, as
// distinct from Tailwind's `disabled:pointer-events-none` / `disabled:opacity-50`
// variant classes that appear verbatim in the button's base class string and
// would otherwise false-positive a plain strings.Contains(got, "disabled").
var disabledAttr = regexp.MustCompile(`disabled(>|\s)`)

func TestButtonDefault(t *testing.T) {
	got := render(t, ui.Button("", "", "", false, gsx.Raw("Save"), nil))
	for _, want := range []string{
		"<button", `data-gsxui-slot-button`, `type="button"`,
		`data-variant="default"`, `data-size="default"`,
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
	got := render(t, ui.Button("", "", "", false, gsx.Raw("Save"), nil))
	want := `<button data-variant="default" data-size="default" type="button" data-gsxui-slot-button>Save</button>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestButtonVariantAndSizeAxes(t *testing.T) {
	for _, variant := range []string{"default", "destructive", "outline", "secondary", "ghost", "link"} {
		input := variant
		if variant == "default" {
			input = ""
		}
		got := render(t, ui.Button(input, "", "", false, gsx.Raw("x"), nil))
		if !strings.Contains(got, `data-variant="`+variant+`"`) {
			t.Errorf("variant %s: missing reflected value\nin: %s", variant, got)
		}
	}
	for _, size := range []string{"default", "xs", "sm", "lg", "icon", "icon-xs", "icon-sm", "icon-lg"} {
		input := size
		if size == "default" {
			input = ""
		}
		got := render(t, ui.Button("", input, "", false, gsx.Raw("x"), nil))
		if !strings.Contains(got, `data-size="`+size+`"`) {
			t.Errorf("size %s: missing reflected value\nin: %s", size, got)
		}
	}
}

func TestButtonHrefRendersAnchor(t *testing.T) {
	got := render(t, ui.Button("", "", "/docs", false, gsx.Raw("Docs"), nil))
	for _, want := range []string{"<a", `href="/docs"`, `data-gsxui-slot-button`, ">Docs</a>"} {
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
	got := render(t, ui.Button("", "", "/docs", true, gsx.Raw("x"), nil))
	if !strings.Contains(got, "<button") {
		t.Errorf("want disabled <button>\nin: %s", got)
	}
	if !disabledAttr.MatchString(got) {
		t.Errorf("want real disabled attribute\nin: %s", got)
	}
}

func TestButtonTypeIsOverridableDefault(t *testing.T) {
	// type="button" is authored BEFORE { attrs... }: caller type=submit wins.
	got := render(t, ui.Button("", "", "", false, gsx.Raw("Go"), gsx.Attrs{{Key: "type", Value: "submit"}}))
	if !strings.Contains(got, `type="submit"`) {
		t.Errorf("caller type=submit must override default\nin: %s", got)
	}
	if strings.Contains(got, `type="button"`) {
		t.Errorf("default type should be replaced, not duplicated\nin: %s", got)
	}
}

func TestButtonCallerClassIsForwardedOnce(t *testing.T) {
	got := render(t, ui.Button("", "", "", false, gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "h-12"}}))
	if strings.Count(got, `class="h-12"`) != 1 {
		t.Errorf("caller class must be the only class and render once\nin: %s", got)
	}
}

func TestButtonOwnPresenceMarkerWinsCollisionAndKeepsComposedMarker(t *testing.T) {
	got := render(t, ui.Button("", "", "", false, gsx.Raw("x"), gsx.Attrs{
		{Key: "data-gsxui-slot-button", Value: "caller-value"},
		{Key: "data-gsxui-slot-pagination-link", Value: gsx.Toggle(true)},
	}))
	requirePresenceAttributesOnSameTag(t, got, "<button",
		"data-gsxui-slot-pagination-link",
		"data-gsxui-slot-button",
	)
}
