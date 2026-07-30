package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	maiaui "github.com/gsxhq/gsxui/registry/generated/maia"
	novaui "github.com/gsxhq/gsxui/registry/generated/nova"
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
// Button seam: Button contributes the default style's resolved presentation
// and its public axes. InputGroupButton defaults and forwards Button's size
// axis to xs.
func TestInputGroupButtonDefaultPinned(t *testing.T) {
	got := render(t, ui.InputGroupButton("", "", gsx.Raw("x"), nil))
	want := `<button data-variant="ghost" data-size="xs" type="button" ` + canonicalButtonClass("ghost", "xs") + ` data-gsxui-slot-input-group-button data-gsxui-slot-button>x</button>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestInputGroupButtonSizeAxisMatchesCanonicalButtonClass(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "default", input: "", want: "xs"},
		{name: "xs", input: "xs", want: "xs"},
		{name: "sm", input: "sm", want: "sm"},
		{name: "icon-xs", input: "icon-xs", want: "icon-xs"},
		{name: "icon-sm", input: "icon-sm", want: "icon-sm"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := render(t, ui.InputGroupButton("", tt.input, gsx.Raw("x"), nil))
			for _, want := range []string{
				`data-size="` + tt.want + `"`,
				canonicalButtonClass("ghost", tt.want),
			} {
				if !strings.Contains(got, want) {
					t.Errorf("missing aligned Button size contract %q\nin: %s", want, got)
				}
			}
		})
	}
}

func TestGeneratedButtonSizeContractsForInputGroupComposition(t *testing.T) {
	tests := []struct {
		name       string
		node       gsx.Node
		dataSize   string
		wantClass  string
		rejectSize string
	}{
		{name: "nova/xs", node: novaui.Button("ghost", "xs", "", false, gsx.Raw("x"), nil), dataSize: "xs", wantClass: "h-6", rejectSize: "h-8"},
		{name: "nova/sm", node: novaui.Button("ghost", "sm", "", false, gsx.Raw("x"), nil), dataSize: "sm", wantClass: "h-7", rejectSize: "h-8"},
		{name: "nova/icon-xs", node: novaui.Button("ghost", "icon-xs", "", false, gsx.Raw("x"), nil), dataSize: "icon-xs", wantClass: "size-6", rejectSize: "size-8"},
		{name: "nova/icon-sm", node: novaui.Button("ghost", "icon-sm", "", false, gsx.Raw("x"), nil), dataSize: "icon-sm", wantClass: "size-7", rejectSize: "size-8"},
		{name: "maia/xs", node: maiaui.Button("ghost", "xs", "", false, gsx.Raw("x"), nil), dataSize: "xs", wantClass: "h-6", rejectSize: "h-9"},
		{name: "maia/sm", node: maiaui.Button("ghost", "sm", "", false, gsx.Raw("x"), nil), dataSize: "sm", wantClass: "h-8", rejectSize: "h-9"},
		{name: "maia/icon-xs", node: maiaui.Button("ghost", "icon-xs", "", false, gsx.Raw("x"), nil), dataSize: "icon-xs", wantClass: "size-6", rejectSize: "size-9"},
		{name: "maia/icon-sm", node: maiaui.Button("ghost", "icon-sm", "", false, gsx.Raw("x"), nil), dataSize: "icon-sm", wantClass: "size-8", rejectSize: "size-9"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := render(t, tt.node)
			for _, want := range []string{`data-size="` + tt.dataSize + `"`, tt.wantClass} {
				if !strings.Contains(got, want) {
					t.Errorf("generated Button lacks concrete size contract %q\nin: %s", want, got)
				}
			}
			if strings.Contains(got, tt.rejectSize) {
				t.Errorf("generated Button retained contradictory default geometry %q\nin: %s", tt.rejectSize, got)
			}
		})
	}
}

// TestInputGroupButtonVariantOverride proves variant is forwarded as Button's
// public styling axis, and that the forwarded variant selects the default
// style's outline presentation rather than the default variant's.
func TestInputGroupButtonVariantOverride(t *testing.T) {
	got := render(t, ui.InputGroupButton("outline", "", gsx.Raw("x"), nil))
	for _, want := range []string{`data-variant="outline"`, canonicalButtonClass("outline", "xs")} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
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
	want := canonicalButtonClass("ghost", "xs", "w-full")
	if strings.Count(got, want) != 1 || strings.Count(got, `class=`) != 1 {
		t.Errorf("caller class must follow exact canonical Button roles once\nwant: %s\nin: %s", want, got)
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
	want := `<textarea class="field-sizing-content flex min-h-16 w-full rounded-lg border border-input bg-transparent px-2.5 py-2 text-base transition-[color,box-shadow] outline-none placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:bg-input/30 dark:aria-invalid:ring-destructive/40 md:text-sm" data-gsxui-slot-input-group-control data-gsxui-slot-textarea>hi</textarea>`
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
