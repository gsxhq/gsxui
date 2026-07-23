package ui_test

import (
	"encoding/xml"
	"net/url"
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestCheckboxDefault(t *testing.T) {
	got := render(t, ui.Checkbox(nil))
	for _, want := range []string{
		`<input type="checkbox"`,
		`data-slot="checkbox"`,
		"peer size-4 shrink-0 appearance-none rounded-[4px]",
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

func TestCheckboxDataURIDecodesToValidSVG(t *testing.T) {
	// The check glyph is a data-URI SVG in a checked:bg-[url(...)] arbitrary
	// value. Its whitespace must be percent-encoded (%20): a literal space is
	// a class-token boundary (torn apart by any whitespace-splitting class
	// tooling, tailwind-merge included), and Tailwind's underscore escape is
	// NOT converted to a space inside url() — Tailwind v4 deliberately
	// preserves underscores in URLs, so `_` reaches the browser verbatim and
	// corrupts the SVG (<svg_xmlns=...> — the shipped-broken-checkmark bug).
	// Prove the browser's view: extract the URI payload from the rendered
	// class, percent-decode it, and parse it as XML.
	got := render(t, ui.Checkbox(nil))
	const pre, post = "data:image/svg+xml;charset=utf-8,", "&#39;)]"
	start := strings.Index(got, pre)
	if start < 0 {
		t.Fatalf("data-URI not found in render\nin: %s", got)
	}
	payload := got[start+len(pre):]
	end := strings.Index(payload, post)
	if end < 0 {
		t.Fatalf("unterminated data-URI in render\nin: %s", got)
	}
	payload = payload[:end]
	if strings.ContainsAny(payload, "_ ") {
		t.Errorf("data-URI must percent-encode whitespace (%%20): underscores survive Tailwind's url() handling and literal spaces split the class token\nuri: %s", payload)
	}
	decoded, err := url.PathUnescape(payload)
	if err != nil {
		t.Fatalf("data-URI does not percent-decode: %v\nuri: %s", err, payload)
	}
	var svg struct {
		XMLName xml.Name
		Stroke  string `xml:"stroke,attr"`
	}
	if err := xml.Unmarshal([]byte(decoded), &svg); err != nil {
		t.Fatalf("decoded data-URI is not well-formed XML (the browser would drop the checkmark): %v\nsvg: %s", err, decoded)
	}
	if svg.XMLName.Local != "svg" || svg.Stroke != "white" {
		t.Errorf("decoded root = <%s stroke=%q>, want <svg stroke=\"white\">\nsvg: %s", svg.XMLName.Local, svg.Stroke, decoded)
	}
}

func TestCheckboxCallerClassMerges(t *testing.T) {
	got := render(t, ui.Checkbox(gsx.Attrs{{Key: "class", Value: "size-6"}}))
	if strings.Contains(got, "size-4") {
		t.Errorf("base size-4 should be dropped by caller size-6\nin: %s", got)
	}
	if !strings.Contains(got, "size-6") {
		t.Errorf("missing caller class size-6\nin: %s", got)
	}
}

func TestCheckboxAttrsFallThrough(t *testing.T) {
	got := render(t, ui.Checkbox(gsx.Attrs{{Key: "id", Value: "c1"}, {Key: "name", Value: "terms"}, {Key: "aria-label", Value: "Accept"}}))
	for _, want := range []string{`id="c1"`, `name="terms"`, `aria-label="Accept"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestCheckboxCheckedAttr(t *testing.T) {
	// checked is an HTML boolean (presence-only) attribute: a bool value on
	// it must render bare — no checked="true" — matching browser :checked
	// truth (gsx.IsBooleanAttr classifies "checked").
	got := render(t, ui.Checkbox(gsx.Attrs{{Key: "checked", Value: true}}))
	if !strings.Contains(got, " checked") || strings.Contains(got, `checked="`) {
		t.Errorf("checked attr should render bare, not stringified\nin: %s", got)
	}

	got = render(t, ui.Checkbox(gsx.Attrs{{Key: "checked", Value: false}}))
	// The class carries "checked:"-variant tokens regardless, so assert on
	// the bare-attribute shape specifically, not the substring "checked".
	if strings.Contains(got, `" checked`) || strings.Contains(got, `checked="false"`) {
		t.Errorf("checked=false should omit the attribute entirely\nin: %s", got)
	}
}

func TestCheckboxDisabledAttr(t *testing.T) {
	got := render(t, ui.Checkbox(gsx.Attrs{{Key: "disabled", Value: true}}))
	if !strings.Contains(got, " disabled") || strings.Contains(got, `disabled="`) {
		t.Errorf("disabled attr should render bare\nin: %s", got)
	}
}

func TestCheckboxPinned(t *testing.T) {
	// Exact full-render pin. Token-by-token verified against shadcn's
	// Checkbox (registry/new-york-v4/ui/checkbox.tsx) plus the ledgered
	// ADAPT: Radix Root/Indicator/CheckIcon replaced by a native
	// <input type="checkbox"> whose checked-state visuals move to a
	// checked:bg-[url(...)] data-URI background. See docs/jsx-parity.md.
	got := render(t, ui.Checkbox(nil))
	want := `<input type="checkbox" data-slot="checkbox" class="peer size-4 shrink-0 appearance-none rounded-[4px] border border-input shadow-xs transition-shadow outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 dark:bg-input/30 checked:bg-primary checked:border-primary checked:bg-[url(&#39;data:image/svg+xml;charset=utf-8,%3Csvg%20xmlns=%22http://www.w3.org/2000/svg%22%20viewBox=%220%200%2024%2024%22%20fill=%22none%22%20stroke=%22white%22%20stroke-width=%223%22%20stroke-linecap=%22round%22%20stroke-linejoin=%22round%22%3E%3Cpath%20d=%22M20%206%209%2017l-5-5%22/%3E%3C/svg%3E&#39;)] checked:bg-center checked:bg-no-repeat checked:bg-[length:12px_12px]"/>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}
