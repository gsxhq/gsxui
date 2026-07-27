package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestTextareaDefault(t *testing.T) {
	got := render(t, ui.Textarea("", nil))
	for _, want := range []string{
		"<textarea", `data-gsxui-slot-textarea`,
		"></textarea>",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestTextareaPinned(t *testing.T) {
	// Presentation lives in the stylesheet; the render pin covers structure.
	got := render(t, ui.Textarea("", nil))
	want := `<textarea data-gsxui-slot-textarea></textarea>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestTextareaValueRendersAsTextChild(t *testing.T) {
	// ADAPT: native <textarea> takes its initial content as a text child,
	// not a value attribute. value renders as escaped text between the tags.
	got := render(t, ui.Textarea("hello world", nil))
	if !strings.Contains(got, ">hello world</textarea>") {
		t.Errorf("want value rendered as text child\nin: %s", got)
	}
	if strings.Contains(got, `value=`) {
		t.Errorf("value must not render as an attribute\nin: %s", got)
	}
}

func TestTextareaValueEscaped(t *testing.T) {
	got := render(t, ui.Textarea("<script>", nil))
	if strings.Contains(got, "<script>") {
		t.Errorf("value must be escaped, not raw\nin: %s", got)
	}
	if !strings.Contains(got, "&lt;script&gt;") {
		t.Errorf("want escaped value\nin: %s", got)
	}
}

func TestTextareaCallerClassMerges(t *testing.T) {
	got := render(t, ui.Textarea("", gsx.Attrs{{Key: "class", Value: "min-h-32"}}))
	if strings.Count(got, `class="min-h-32"`) != 1 {
		t.Errorf("caller class must be forwarded exactly once\nin: %s", got)
	}
}

func TestTextareaAttrsFallThrough(t *testing.T) {
	got := render(t, ui.Textarea("", gsx.Attrs{{Key: "id", Value: "bio"}, {Key: "placeholder", Value: "Tell us about yourself"}, {Key: "rows", Value: "4"}}))
	for _, want := range []string{`id="bio"`, `placeholder="Tell us about yourself"`, `rows="4"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}
