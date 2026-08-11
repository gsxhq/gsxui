package ui_test

import (
	"encoding/base64"
	"encoding/xml"
	"os"
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestCheckboxDefault(t *testing.T) {
	got := render(t, ui.Checkbox(nil))
	for _, want := range []string{
		`<input type="checkbox"`,
		`data-gsxui-slot-checkbox`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestCheckboxDataURIDecodesToValidSVG(t *testing.T) {
	// The check glyph now lives in the stylesheet. Its payload remains
	// base64 so CSS parsing cannot corrupt the SVG.
	// TWO check URIs: the light one strokes white (primary is near-black),
	// the dark: variant strokes the dark theme's --primary-foreground value
	// (primary flips near-white, where a white check would vanish).
	got := checkboxCSS(t)
	const pre, post = "data:image/svg+xml;base64,", `")`
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
	got := checkboxCSS(t)
	for _, want := range []string{
		`.dark [data-gsxui-slot-checkbox]:checked`,
		`background-image: url("data:image/svg+xml;base64,`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestCheckboxCallerClassMerges(t *testing.T) {
	got := render(t, ui.Checkbox(gsx.Attrs{{Key: "class", Value: "size-6"}}))
	if !strings.Contains(got, "size-6") || strings.Count(got, "class=") != 1 {
		t.Errorf("caller class must be forwarded exactly once\nin: %s", got)
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
	// data-checked:/aria-invalid:aria-checked: were dead selectors — our
	// <input type="checkbox"> is native, so :checked:/aria-invalid:checked:
	// are what actually match. See the style-porter report's "Radio/Checkbox
	// data-checked: -> native :checked:" entry.
	got := render(t, ui.Checkbox(nil))
	want := `<input type="checkbox" class="border-input dark:bg-input/30 checked:bg-primary checked:text-primary-foreground dark:checked:bg-primary checked:border-primary aria-invalid:checked:border-primary aria-invalid:border-destructive dark:aria-invalid:border-destructive/50 focus-visible:border-ring focus-visible:ring-ring/50 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 flex size-4 items-center justify-center rounded-[4px] border transition-colors group-has-disabled/field:opacity-50 focus-visible:ring-3 aria-invalid:ring-3" data-gsxui-slot-checkbox>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

// checkboxCSS returns Checkbox's authored presentation. The default style's
// components layer is one file per component under assets/css/styles/default/,
// so this is Checkbox's file rather than a 3000-line sheet — and when Checkbox
// migrates to the recipe model that file is deleted, which fails these tests
// loudly instead of letting them keep passing against a stale block.
func checkboxCSS(t *testing.T) string {
	t.Helper()
	cssBytes, err := os.ReadFile("../assets/css/styles/default.css")
	if err != nil {
		t.Fatal(err)
	}
	return string(cssBytes)
}
