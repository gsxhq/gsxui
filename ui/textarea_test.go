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
		"<textarea", `data-slot="textarea"`,
		"flex field-sizing-content min-h-16 w-full rounded-md border border-input",
		"focus-visible:border-ring focus-visible:ring-[3px]",
		"disabled:cursor-not-allowed disabled:opacity-50",
		"></textarea>",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestTextareaPinned(t *testing.T) {
	// Exact full-render pin, verified token-by-token against shadcn's
	// Textarea (registry/new-york-v4/ui/textarea.tsx) and docs/jsx-parity.md
	// — the only ADAPT is the value param (see docs/jsx-parity.md ## textarea).
	got := render(t, ui.Textarea("", nil))
	want := `<textarea data-slot="textarea" class="flex field-sizing-content min-h-16 w-full rounded-md border border-input bg-transparent px-3 py-2 text-base shadow-xs transition-[color,box-shadow] outline-none placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 md:text-sm dark:bg-input/30 dark:aria-invalid:ring-destructive/40"></textarea>`
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
	if strings.Contains(got, "min-h-16") {
		t.Errorf("caller min-h-32 must drop default min-h-16\nin: %s", got)
	}
	for _, want := range []string{"min-h-32", "rounded-md"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
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
