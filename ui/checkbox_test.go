package ui_test

import (
	"encoding/base64"
	"encoding/xml"
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
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestCheckboxDataURIDecodesToValidSVG(t *testing.T) {
	// The check glyph is a data-URI SVG in a checked:bg-[url(...)] arbitrary
	// value, and its payload must be base64: every richer encoding lost to
	// some layer of the toolchain in turn — literal spaces are class-token
	// boundaries (torn by tailwind-merge), Tailwind's `_` escape is NOT
	// converted inside url() (reached the browser as <svg_xmlns=...>, the
	// shipped-broken-checkmark bug), and percent-encoded payloads with
	// parens broke the vite postcss parse of Tailwind's emitted CSS. Base64
	// is [A-Za-z0-9+/=] only: nothing for any layer to split, convert, or
	// mis-parse. Prove the browser's view: extract each rendered URI,
	// base64-decode it, and parse it as XML.
	// TWO check URIs: the light one strokes white (primary is near-black),
	// the dark: variant strokes the dark theme's --primary-foreground value
	// (primary flips near-white, where a white check would vanish).
	got := render(t, ui.Checkbox(nil))
	const pre, post = "data:image/svg+xml;base64,", "&#39;)]"
	var strokes []string
	rest := got
	for {
		start := strings.Index(rest, pre)
		if start < 0 {
			break
		}
		payload := rest[start+len(pre):]
		end := strings.Index(payload, post)
		if end < 0 {
			t.Fatalf("unterminated data-URI in render\nin: %s", got)
		}
		rest = payload[end:]
		payload = payload[:end]
		decodedBytes, err := base64.StdEncoding.DecodeString(payload)
		if err != nil {
			t.Fatalf("data-URI payload is not clean base64 (any other encoding gets mangled by some toolchain layer): %v\nuri: %s", err, payload)
		}
		decoded := string(decodedBytes)
		var svg struct {
			XMLName xml.Name
			Stroke  string `xml:"stroke,attr"`
		}
		if err := xml.Unmarshal([]byte(decoded), &svg); err != nil {
			t.Fatalf("decoded data-URI is not well-formed XML (the browser would drop the checkmark): %v\nsvg: %s", err, decoded)
		}
		if svg.XMLName.Local != "svg" {
			t.Errorf("decoded root = <%s>, want <svg>\nsvg: %s", svg.XMLName.Local, decoded)
		}
		strokes = append(strokes, svg.Stroke)
	}
	want := []string{"white", "oklch(0.205 0 0)"}
	if len(strokes) != len(want) || strokes[0] != want[0] || strokes[1] != want[1] {
		t.Errorf("check-glyph strokes = %q, want %q (light white-on-primary, dark primary-foreground-on-primary)", strokes, want)
	}
}

func TestCheckboxDarkCheckedOverrides(t *testing.T) {
	// shadcn's dark:data-[state=checked]:bg-primary is NOT redundant: the
	// dark custom variant (:is(.dark *)) carries class specificity (0,2,0),
	// which beats plain :checked (0,1,1) — without the explicit dark:checked
	// overrides, dark:bg-input/30 wins in dark mode and a checked box
	// renders 4.5%-alpha instead of primary (found live, dark-mode sweep).
	got := render(t, ui.Checkbox(nil))
	for _, want := range []string{"dark:checked:bg-primary", "dark:checked:bg-[url("} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
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
	want := `<input type="checkbox" data-slot="checkbox" class="peer size-4 shrink-0 appearance-none rounded-[4px] border border-input shadow-xs transition-shadow outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 dark:bg-input/30 checked:bg-primary checked:border-primary checked:bg-[url(&#39;data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAyNCAyNCIgZmlsbD0ibm9uZSIgc3Ryb2tlPSJ3aGl0ZSIgc3Ryb2tlLXdpZHRoPSIzIiBzdHJva2UtbGluZWNhcD0icm91bmQiIHN0cm9rZS1saW5lam9pbj0icm91bmQiPjxwYXRoIGQ9Ik0yMCA2IDkgMTdsLTUtNSIvPjwvc3ZnPg==&#39;)] checked:bg-center checked:bg-no-repeat checked:bg-[length:12px_12px] dark:checked:bg-primary dark:checked:border-primary dark:checked:bg-[url(&#39;data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAyNCAyNCIgZmlsbD0ibm9uZSIgc3Ryb2tlPSJva2xjaCgwLjIwNSAwIDApIiBzdHJva2Utd2lkdGg9IjMiIHN0cm9rZS1saW5lY2FwPSJyb3VuZCIgc3Ryb2tlLWxpbmVqb2luPSJyb3VuZCI+PHBhdGggZD0iTTIwIDYgOSAxN2wtNS01Ii8+PC9zdmc+&#39;)]">`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}
