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
	// The leading structural baseline is restored "carried: no upstream
	// counterpart" content — see the style-porter report's "Input/Textarea
	// structural baseline" entry.
	got := render(t, ui.Textarea("", nil))
	want := `<textarea class="field-sizing-content min-h-16 w-full outline-none placeholder:text-muted-foreground disabled:cursor-not-allowed disabled:opacity-50 border-input dark:bg-input/30 focus-visible:border-ring focus-visible:ring-ring/50 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive dark:aria-invalid:border-destructive/50 disabled:bg-input/50 dark:disabled:bg-input/80 rounded-lg border bg-transparent px-2.5 py-2 text-base transition-colors focus-visible:ring-3 aria-invalid:ring-3 md:text-sm flex" data-gsxui-slot-textarea></textarea>`
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
	if !strings.Contains(got, "min-h-32") {
		t.Errorf("caller class must be forwarded\nin: %s", got)
	}
	if strings.Contains(got, "min-h-16") {
		t.Errorf("caller's min-h-32 must win over the recipe's min-h-16\nin: %s", got)
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
