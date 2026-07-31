package ui_test

import (
	"html"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/internal/stylegen"
	"github.com/gsxhq/gsxui/merge"
	maiaui "github.com/gsxhq/gsxui/registry/generated/maia"
	novaui "github.com/gsxhq/gsxui/registry/generated/nova"
	"github.com/gsxhq/gsxui/ui"
)

// defaultStyleRecipe loads one component's recipe stylesheet from the default
// style. InputGroup composes three migrated components (Button, Input,
// Textarea), so unlike the single-component helpers in separator_test.go and
// item_test.go this one is parameterised by component name.
var defaultStyleRecipe = sync.OnceValue(func() func(string) recipe.Style {
	cache := map[string]recipe.Style{}
	return func(component string) recipe.Style {
		if style, ok := cache[component]; ok {
			return style
		}
		path := filepath.Join("..", "registry", "styles", stylegen.DefaultStyle, component+".css")
		src, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}
		style, err := recipe.ParseStyle(path, src)
		if err != nil {
			panic(err)
		}
		cache[component] = style
		return style
	}
})

func styleRecipeUtilities(component, class string) []string {
	rule, ok := defaultStyleRecipe()(component).Lookup(class)
	if !ok {
		panic("default style declares no recipe " + class)
	}
	return rule.Utilities
}

func inputGroupRecipeClasses(slot string, extra ...string) []string {
	base := "gsxui-recipe-input-group"
	if slot != "" {
		base += "-" + slot
	}
	classes := append([]string(nil), styleRecipeUtilities("input-group", base)...)
	for _, value := range extra {
		classes = append(classes, styleRecipeUtilities("input-group", base+"-"+value)...)
	}
	return classes
}

// canonicalInputGroupClass is the class attribute one InputGroup slot renders,
// plus any caller classes, merged the way gsx merges class values at runtime.
func canonicalInputGroupClass(slot string, extra []string, caller ...string) string {
	classes := append(inputGroupRecipeClasses(slot, extra...), caller...)
	return `class="` + html.EscapeString(merge.Merge(classes)) + `"`
}

// canonicalInputGroupControlClass is what InputGroupInput/InputGroupTextarea
// render: the composed control's own recipe utilities with InputGroup's
// control slot merged in after them as an ordinary caller class.
//
// tailwind-merge genuinely REMOVES two of the control's own utilities here —
// focus-visible:ring-[3px] loses to focus-visible:ring-0, and dark:bg-input/30
// loses to dark:bg-transparent. That is upstream shadcn's own behaviour
// (registry/new-york-v4/ui/input-group.tsx puts exactly those two utilities on
// InputGroupInput's className and lets cn() resolve them), and it is a change
// from the pre-migration render, where the same two rules sat in
// @layer components and lost to Input's utilities-layer classes instead.
func canonicalInputGroupControlClass(control string, caller ...string) string {
	classes := append([]string(nil), styleRecipeUtilities(control, "gsxui-recipe-"+control)...)
	classes = append(classes, inputGroupRecipeClasses("control")...)
	classes = append(classes, caller...)
	return `class="` + html.EscapeString(merge.Merge(classes)) + `"`
}

func TestInputGroupPinned(t *testing.T) {
	got := render(t, ui.InputGroup(gsx.Raw("x"), nil))
	want := `<div role="group" ` + canonicalInputGroupClass("", nil) + ` data-gsxui-slot-input-group>x</div>`
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
	if strings.Count(got, "max-w-sm") != 1 || strings.Count(got, `class=`) != 1 {
		t.Errorf("caller class must be forwarded exactly once\nin: %s", got)
	}
	if want := canonicalInputGroupClass("", nil, "max-w-sm"); !strings.Contains(got, want) {
		t.Errorf("caller class not merged after the recipe's own\nwant: %s\nin: %s", want, got)
	}
}

// TestInputGroupAddonDefaultPinned pins the zero-value ("inline-start") align.
func TestInputGroupAddonDefaultPinned(t *testing.T) {
	got := render(t, ui.InputGroupAddon("", gsx.Raw("x"), nil))
	want := `<div role="group" data-align="inline-start" ` + canonicalInputGroupClass("addon", []string{"align-inline-start"}) + ` data-gsxui-slot-input-group-addon>x</div>`
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
	want := `<div role="group" data-align="block-start" ` + canonicalInputGroupClass("addon", []string{"align-block-start"}) + ` data-gsxui-slot-input-group-addon>x</div>`
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
	want := `<span ` + canonicalInputGroupClass("text", nil) + ` data-gsxui-slot-input-group-text>x</span>`
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
	want := `<input type="text" ` + canonicalInputGroupControlClass("input") + ` data-gsxui-slot-input-group-control data-gsxui-slot-input>`
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
	want := `<textarea ` + canonicalInputGroupControlClass("textarea") + ` data-gsxui-slot-input-group-control data-gsxui-slot-textarea>hi</textarea>`
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
