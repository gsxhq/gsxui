package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestRadioDefault(t *testing.T) {
	got := render(t, ui.Radio(nil))
	for _, want := range []string{
		`<input type="radio"`,
		`data-slot="radio"`,
		"peer aspect-square size-4 shrink-0 appearance-none rounded-full",
		"text-primary",
		"disabled:cursor-not-allowed disabled:opacity-50",
		"aria-invalid:border-destructive aria-invalid:ring-destructive/20",
		"checked:bg-[radial-gradient(circle_closest-side,currentColor_45%,transparent_50%)]",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestRadioNoStraySpaceInGradient(t *testing.T) {
	// Same class-merge corruption hazard as checkbox's data-URI — see
	// checkbox_test.go's TestCheckboxNoStraySpaceInDataURI and
	// docs/jsx-parity.md. The radial-gradient's embedded spaces are
	// Tailwind's underscore escape, not literal spaces, so it must survive
	// the tailwind-merge pass (merge.Merge, invoked on every render) intact
	// rather than being torn apart at a bare space token boundary.
	got := render(t, ui.Radio(nil))
	want := "checked:bg-[radial-gradient(circle_closest-side,currentColor_45%,transparent_50%)]"
	if !strings.Contains(got, want) {
		t.Errorf("radial-gradient corrupted by class merge\nwant substring: %s\n in: %s", want, got)
	}
}

func TestRadioCallerClassMerges(t *testing.T) {
	got := render(t, ui.Radio(gsx.Attrs{{Key: "class", Value: "size-6"}}))
	if strings.Contains(got, "size-4") {
		t.Errorf("base size-4 should be dropped by caller size-6\nin: %s", got)
	}
	if !strings.Contains(got, "size-6") {
		t.Errorf("missing caller class size-6\nin: %s", got)
	}
}

func TestRadioAttrsFallThrough(t *testing.T) {
	got := render(t, ui.Radio(gsx.Attrs{{Key: "id", Value: "r1"}, {Key: "name", Value: "plan"}, {Key: "value", Value: "pro"}, {Key: "aria-label", Value: "Pro"}}))
	for _, want := range []string{`id="r1"`, `name="plan"`, `value="pro"`, `aria-label="Pro"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestRadioCheckedAttr(t *testing.T) {
	got := render(t, ui.Radio(gsx.Attrs{{Key: "checked", Value: true}}))
	if !strings.Contains(got, " checked") || strings.Contains(got, `checked="`) {
		t.Errorf("checked attr should render bare, not stringified\nin: %s", got)
	}

	got = render(t, ui.Radio(gsx.Attrs{{Key: "checked", Value: false}}))
	if strings.Contains(got, `" checked`) || strings.Contains(got, `checked="false"`) {
		t.Errorf("checked=false should omit the attribute entirely\nin: %s", got)
	}
}

func TestRadioDisabledAttr(t *testing.T) {
	got := render(t, ui.Radio(gsx.Attrs{{Key: "disabled", Value: true}}))
	if !strings.Contains(got, " disabled") || strings.Contains(got, `disabled="`) {
		t.Errorf("disabled attr should render bare\nin: %s", got)
	}
}

func TestRadioPinned(t *testing.T) {
	// Exact full-render pin. Token-by-token verified against shadcn's
	// RadioGroupItem (registry/new-york-v4/ui/radio-group.tsx) plus the
	// ledgered ADAPTs: native <input type="radio"> replaces Radix
	// Root/Indicator/CircleIcon; the Indicator's fill-primary CircleIcon
	// becomes a checked:bg-[radial-gradient(...)] painted in currentColor
	// (text-primary is what makes currentColor resolve to primary here —
	// load-bearing, not vestigial). See docs/jsx-parity.md.
	got := render(t, ui.Radio(nil))
	want := `<input type="radio" data-slot="radio" class="peer aspect-square size-4 shrink-0 appearance-none rounded-full border border-input transition-[color,box-shadow] outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 dark:bg-input/30 checked:border-primary checked:bg-primary checked:text-primary-foreground dark:checked:bg-primary checked:bg-[radial-gradient(circle_closest-side,currentColor_45%,transparent_50%)]">`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}
