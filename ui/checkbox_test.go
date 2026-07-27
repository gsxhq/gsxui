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
	cssBytes, err := os.ReadFile("../assets/css/styles/default.css")
	if err != nil {
		t.Fatal(err)
	}
	got := string(cssBytes)
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
	cssBytes, err := os.ReadFile("../assets/css/styles/default.css")
	if err != nil {
		t.Fatal(err)
	}
	got := string(cssBytes)
	for _, want := range []string{
		`.dark :where([data-gsxui-slot-checkbox]):checked`,
		`background-image: url("data:image/svg+xml;base64,`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestCheckboxCallerClassMerges(t *testing.T) {
	got := render(t, ui.Checkbox(gsx.Attrs{{Key: "class", Value: "size-6"}}))
	if strings.Count(got, `class="size-6"`) != 1 {
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
	// Presentation lives in the stylesheet; the render pin covers structure.
	got := render(t, ui.Checkbox(nil))
	want := `<input type="checkbox" data-gsxui-slot-checkbox>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}
