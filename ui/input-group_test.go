package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestInputGroupPinned(t *testing.T) {
	got := render(t, ui.InputGroup(gsx.Raw("x"), nil))
	want := `<div role="group" data-gsxui-slot-input-group>x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestInputGroupAttrsFallThrough(t *testing.T) {
	got := render(t, ui.InputGroup(nil, gsx.Attrs{{Key: "id", Value: "ig1"}}))
	if !strings.Contains(got, `id="ig1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

func TestInputGroupCallerClassMerges(t *testing.T) {
	got := render(t, ui.InputGroup(nil, gsx.Attrs{{Key: "class", Value: "max-w-sm"}}))
	if strings.Count(got, `class="max-w-sm"`) != 1 {
		t.Errorf("caller class must be forwarded exactly once\nin: %s", got)
	}
}

// TestInputGroupAddonDefaultPinned pins the zero-value ("inline-start") align.
func TestInputGroupAddonDefaultPinned(t *testing.T) {
	got := render(t, ui.InputGroupAddon("", gsx.Raw("x"), nil))
	want := `<div role="group" data-align="inline-start" data-gsxui-slot-input-group-addon>x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestInputGroupAddonInlineEndPinned(t *testing.T) {
	got := render(t, ui.InputGroupAddon("inline-end", gsx.Raw("x"), nil))
	for _, want := range []string{`data-align="inline-end"`, `data-gsxui-slot-input-group-addon`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestInputGroupAddonBlockStartPinned(t *testing.T) {
	got := render(t, ui.InputGroupAddon("block-start", gsx.Raw("x"), nil))
	want := `<div role="group" data-align="block-start" data-gsxui-slot-input-group-addon>x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestInputGroupAddonBlockEndPinned(t *testing.T) {
	got := render(t, ui.InputGroupAddon("block-end", gsx.Raw("x"), nil))
	for _, want := range []string{`data-align="block-end"`, `data-gsxui-slot-input-group-addon`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestInputGroupAddonAttrsFallThrough(t *testing.T) {
	got := render(t, ui.InputGroupAddon("", nil, gsx.Attrs{{Key: "id", Value: "a1"}}))
	if !strings.Contains(got, `id="a1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

// TestInputGroupButtonDefaultPinned proves InputGroupButton composes the
// CSS-owned Button seam: Button contributes its stable token and public axes,
// while InputGroupButton's caller classes remain the only rendered classes.
func TestInputGroupButtonDefaultPinned(t *testing.T) {
	got := render(t, ui.InputGroupButton("", "", gsx.Raw("x"), nil))
	want := `<button data-variant="ghost" type="button" data-size="xs" data-gsxui-slot-input-group-button data-gsxui-slot-button>x</button>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
	if strings.Contains(got, "focus-visible:ring") || strings.Contains(got, "hover:bg-accent") {
		t.Errorf("Button presentation classes must be owned by CSS\nin: %s", got)
	}
}

func TestInputGroupButtonSmPinned(t *testing.T) {
	got := render(t, ui.InputGroupButton("", "sm", gsx.Raw("x"), nil))
	for _, want := range []string{`data-size="sm"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestInputGroupButtonIconXsPinned(t *testing.T) {
	got := render(t, ui.InputGroupButton("", "icon-xs", gsx.Raw("x"), nil))
	for _, want := range []string{`data-size="icon-xs"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestInputGroupButtonIconSmPinned(t *testing.T) {
	got := render(t, ui.InputGroupButton("", "icon-sm", gsx.Raw("x"), nil))
	for _, want := range []string{`data-size="icon-sm"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

// TestInputGroupButtonVariantOverride proves variant is forwarded as Button's
// public styling axis. The variant's presentation is owned by CSS.
func TestInputGroupButtonVariantOverride(t *testing.T) {
	got := render(t, ui.InputGroupButton("outline", "", gsx.Raw("x"), nil))
	if !strings.Contains(got, `data-variant="outline"`) {
		t.Errorf("missing data-variant=outline override\nin: %s", got)
	}
	if strings.Contains(got, "dark:border-input dark:bg-input/30") {
		t.Errorf("outline variant presentation must not render inline\nin: %s", got)
	}
}

func TestInputGroupButtonAttrsFallThrough(t *testing.T) {
	got := render(t, ui.InputGroupButton("", "", nil, gsx.Attrs{{Key: "id", Value: "b1"}}))
	if !strings.Contains(got, `id="b1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

func TestInputGroupButtonCallerClassMerges(t *testing.T) {
	got := render(t, ui.InputGroupButton("", "", nil, gsx.Attrs{{Key: "class", Value: "w-full"}}))
	if !strings.Contains(got, "w-full") {
		t.Errorf("missing caller class w-full\nin: %s", got)
	}
}

func TestInputGroupTextPinned(t *testing.T) {
	got := render(t, ui.InputGroupText(gsx.Raw("x"), nil))
	want := `<span data-gsxui-slot-input-group-text>x</span>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestInputGroupTextAttrsFallThrough(t *testing.T) {
	got := render(t, ui.InputGroupText(nil, gsx.Attrs{{Key: "id", Value: "t1"}}))
	if !strings.Contains(got, `id="t1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

// InputGroupInput composes Input's token with the group-control token.
func TestInputGroupInputPinned(t *testing.T) {
	got := render(t, ui.InputGroupInput(nil))
	want := `<input type="text" data-gsxui-slot-input-group-control data-gsxui-slot-input>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestInputGroupInputAttrsFallThrough(t *testing.T) {
	got := render(t, ui.InputGroupInput(gsx.Attrs{{Key: "placeholder", Value: "Search..."}}))
	if !strings.Contains(got, `placeholder="Search..."`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

// InputGroupTextarea composes Textarea's token with the group-control token.
func TestInputGroupTextareaPinned(t *testing.T) {
	got := render(t, ui.InputGroupTextarea("hi", nil))
	want := `<textarea data-gsxui-slot-input-group-control data-gsxui-slot-textarea>hi</textarea>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestInputGroupTextareaAttrsFallThrough(t *testing.T) {
	got := render(t, ui.InputGroupTextarea("", gsx.Attrs{{Key: "id", Value: "ta1"}}))
	if !strings.Contains(got, `id="ta1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

// Realistic composition: a search-style InputGroup (leading icon addon,
// InputGroupInput, trailing InputGroupButton) — the site example's own
// shape.
func TestInputGroupSearchComposition(t *testing.T) {
	got := render(t, ui.InputGroup(
		gsx.Fragment(
			ui.InputGroupAddon("", gsx.Raw("<svg/>"), nil),
			ui.InputGroupInput(gsx.Attrs{{Key: "placeholder", Value: "Search..."}}),
			ui.InputGroupAddon("inline-end", ui.InputGroupButton("", "icon-xs", gsx.Raw("<svg/>"), gsx.Attrs{{Key: "aria-label", Value: "Send"}}), nil),
		),
		nil,
	))
	for _, want := range []string{
		`data-gsxui-slot-input-group`,
		`data-align="inline-start"`,
		`data-gsxui-slot-input-group-control data-gsxui-slot-input`,
		`placeholder="Search..."`,
		`data-align="inline-end"`,
		`data-gsxui-slot-input-group-button data-gsxui-slot-button`,
		`data-size="icon-xs"`,
		`aria-label="Send"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}
