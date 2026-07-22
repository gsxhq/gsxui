package radio_test

import (
	"context"
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui/radio"
)

func render(t *testing.T, n gsx.Node) string {
	t.Helper()
	var sb strings.Builder
	if err := n.Render(context.Background(), &sb); err != nil {
		t.Fatal(err)
	}
	return sb.String()
}

func TestRadioDefault(t *testing.T) {
	got := render(t, radio.Radio(nil))
	for _, want := range []string{
		`<input type="radio"`,
		`data-slot="radio"`,
		"peer size-4 shrink-0 appearance-none rounded-full",
		"disabled:cursor-not-allowed disabled:opacity-50",
		"aria-invalid:border-destructive aria-invalid:ring-destructive/20",
		"checked:bg-primary checked:border-primary",
		"checked:bg-[url(&#39;data:image/svg+xml",
		"checked:bg-center checked:bg-no-repeat checked:bg-[length:12px_12px]",
		"/>",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestRadioNoStraySpaceInDataURI(t *testing.T) {
	// Same class-merge corruption hazard as checkbox — see
	// checkbox_test.go's TestCheckboxNoStraySpaceInDataURI and
	// docs/jsx-parity.md.
	got := render(t, radio.Radio(nil))
	want := "checked:bg-[url(&#39;data:image/svg+xml;charset=utf-8,%3Csvg_xmlns=%22http://www.w3.org/2000/svg%22_viewBox=%220_0_24_24%22%3E%3Ccircle_cx=%2212%22_cy=%2212%22_r=%226%22_fill=%22white%22/%3E%3C/svg%3E&#39;)]"
	if !strings.Contains(got, want) {
		t.Errorf("data-URI corrupted by class merge\nwant substring: %s\n in: %s", want, got)
	}
}

func TestRadioCallerClassMerges(t *testing.T) {
	got := render(t, radio.Radio(gsx.Attrs{{Key: "class", Value: "size-6"}}))
	if strings.Contains(got, "size-4") {
		t.Errorf("base size-4 should be dropped by caller size-6\nin: %s", got)
	}
	if !strings.Contains(got, "size-6") {
		t.Errorf("missing caller class size-6\nin: %s", got)
	}
}

func TestRadioAttrsFallThrough(t *testing.T) {
	got := render(t, radio.Radio(gsx.Attrs{{Key: "id", Value: "r1"}, {Key: "name", Value: "plan"}, {Key: "value", Value: "pro"}, {Key: "aria-label", Value: "Pro"}}))
	for _, want := range []string{`id="r1"`, `name="plan"`, `value="pro"`, `aria-label="Pro"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestRadioCheckedAttr(t *testing.T) {
	got := render(t, radio.Radio(gsx.Attrs{{Key: "checked", Value: true}}))
	if !strings.Contains(got, " checked") || strings.Contains(got, `checked="`) {
		t.Errorf("checked attr should render bare, not stringified\nin: %s", got)
	}

	got = render(t, radio.Radio(gsx.Attrs{{Key: "checked", Value: false}}))
	if strings.Contains(got, `" checked`) || strings.Contains(got, `checked="false"`) {
		t.Errorf("checked=false should omit the attribute entirely\nin: %s", got)
	}
}

func TestRadioDisabledAttr(t *testing.T) {
	got := render(t, radio.Radio(gsx.Attrs{{Key: "disabled", Value: true}}))
	if !strings.Contains(got, " disabled") || strings.Contains(got, `disabled="`) {
		t.Errorf("disabled attr should render bare\nin: %s", got)
	}
}

func TestRadioPinned(t *testing.T) {
	// Exact full-render pin. Token-by-token verified against shadcn's
	// RadioGroupItem (registry/new-york-v4/ui/radio-group.tsx) plus the
	// ledgered ADAPTs: native <input type="radio"> replaces Radix
	// Root/Indicator/CircleIcon; checked-state visuals move to a
	// checked:bg-[url(...)] data-URI circle. See docs/jsx-parity.md.
	got := render(t, radio.Radio(nil))
	want := `<input type="radio" data-slot="radio" class="peer size-4 shrink-0 appearance-none rounded-full border border-input shadow-xs transition-shadow outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 dark:bg-input/30 checked:bg-primary checked:border-primary checked:bg-[url(&#39;data:image/svg+xml;charset=utf-8,%3Csvg_xmlns=%22http://www.w3.org/2000/svg%22_viewBox=%220_0_24_24%22%3E%3Ccircle_cx=%2212%22_cy=%2212%22_r=%226%22_fill=%22white%22/%3E%3C/svg%3E&#39;)] checked:bg-center checked:bg-no-repeat checked:bg-[length:12px_12px]"/>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}
