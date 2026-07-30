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
	"github.com/gsxhq/gsxui/ui"
)

// novaFieldRecipe loads the default style's Field recipe once — same pattern as
// novaSeparatorRecipe in separator_test.go.
var novaFieldRecipe = sync.OnceValue(func() recipe.Style {
	path := filepath.Join("..", "registry", "styles", stylegen.DefaultStyle, "field.css")
	src, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	style, err := recipe.ParseStyle(path, src)
	if err != nil {
		panic(err)
	}
	return style
})

func fieldRecipeUtilities(class string) []string {
	rule, ok := novaFieldRecipe().Lookup(class)
	if !ok {
		panic("default style declares no recipe " + class)
	}
	return rule.Utilities
}

// fieldRecipeClasses is the utility list one Field slot contributes: its base
// rule, then any dimension-value rules, in accessor-call order.
func fieldRecipeClasses(slot string, extra ...string) []string {
	base := "gsxui-recipe-field"
	if slot != "" {
		base += "-" + slot
	}
	classes := append([]string(nil), fieldRecipeUtilities(base)...)
	for _, value := range extra {
		classes = append(classes, fieldRecipeUtilities(base+"-"+value)...)
	}
	return classes
}

// canonicalFieldClass is the class attribute one Field slot renders, plus any
// caller classes, merged the way gsx merges class values at runtime.
func canonicalFieldClass(slot string, extra []string, caller ...string) string {
	classes := append(fieldRecipeClasses(slot, extra...), caller...)
	return `class="` + html.EscapeString(merge.Merge(classes)) + `"`
}

// canonicalFieldLabelClass is what FieldLabel renders: ui.Label's own recipe
// utilities with Field's field-label slot merged in after them as an ordinary
// caller class. tailwind-merge is what resolves leading-snug against Label's
// leading-none here — the reason the old @layer utilities promotion in
// assets/css/styles/default.css could be retired.
func canonicalFieldLabelClass(caller ...string) string {
	classes := append([]string(nil), labelRecipeUtilities("gsxui-recipe-label")...)
	classes = append(classes, fieldRecipeClasses("label")...)
	classes = append(classes, caller...)
	return `class="` + html.EscapeString(merge.Merge(classes)) + `"`
}

// canonicalFieldSeparatorClass is what the <Separator> inside FieldSeparator
// renders: Separator's own horizontal recipe plus Field's field-separator slot.
func canonicalFieldSeparatorClass(caller ...string) string {
	classes := append([]string(nil), separatorRecipeUtilities("gsxui-recipe-separator")...)
	classes = append(classes, separatorRecipeUtilities("gsxui-recipe-separator-orientation-horizontal")...)
	classes = append(classes, fieldRecipeClasses("separator")...)
	classes = append(classes, caller...)
	return `class="` + html.EscapeString(merge.Merge(classes)) + `"`
}

func TestFieldSetPinned(t *testing.T) {
	got := render(t, ui.FieldSet(gsx.Raw("x"), nil))
	want := `<fieldset ` + canonicalFieldClass("set", nil) + ` data-gsxui-slot-field-set>x</fieldset>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestFieldSetAttrsFallThrough(t *testing.T) {
	got := render(t, ui.FieldSet(nil, gsx.Attrs{{Key: "id", Value: "fs1"}}))
	if !strings.Contains(got, `id="fs1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

// TestFieldLegendDefaultPinned pins the zero-value ("legend") variant.
func TestFieldLegendDefaultPinned(t *testing.T) {
	got := render(t, ui.FieldLegend("", gsx.Raw("x"), nil))
	want := `<legend data-variant="legend" ` + canonicalFieldClass("legend", []string{"variant-legend"}) + ` data-gsxui-slot-field-legend>x</legend>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestFieldLegendLabelVariantPinned(t *testing.T) {
	got := render(t, ui.FieldLegend("label", gsx.Raw("x"), nil))
	if !strings.Contains(got, `data-variant="label"`) {
		t.Errorf("missing data-variant=label\nin: %s", got)
	}
}

func TestFieldLegendAttrsFallThrough(t *testing.T) {
	got := render(t, ui.FieldLegend("", nil, gsx.Attrs{{Key: "id", Value: "fl1"}}))
	if !strings.Contains(got, `id="fl1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

func TestFieldGroupPinned(t *testing.T) {
	got := render(t, ui.FieldGroup(gsx.Raw("x"), nil))
	want := `<div ` + canonicalFieldClass("group", nil) + ` data-gsxui-slot-field-group>x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestFieldGroupAttrsFallThrough(t *testing.T) {
	got := render(t, ui.FieldGroup(nil, gsx.Attrs{{Key: "id", Value: "fg1"}}))
	if !strings.Contains(got, `id="fg1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

// TestFieldDefaultPinned pins the zero-value ("vertical") orientation.
func TestFieldDefaultPinned(t *testing.T) {
	got := render(t, ui.Field("", gsx.Raw("x"), nil))
	want := `<div role="group" data-orientation="vertical" ` + canonicalFieldClass("", []string{"orientation-vertical"}) + ` data-gsxui-slot-field>x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

// TestFieldHorizontalPinned proves data-orientation remains load-bearing.
func TestFieldHorizontalPinned(t *testing.T) {
	got := render(t, ui.Field("horizontal", gsx.Raw("x"), nil))
	want := `<div role="group" data-orientation="horizontal" ` + canonicalFieldClass("", []string{"orientation-horizontal"}) + ` data-gsxui-slot-field>x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestFieldResponsivePinned(t *testing.T) {
	got := render(t, ui.Field("responsive", gsx.Raw("x"), nil))
	want := `<div role="group" data-orientation="responsive" ` + canonicalFieldClass("", []string{"orientation-responsive"}) + ` data-gsxui-slot-field>x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestFieldAttrsFallThrough(t *testing.T) {
	got := render(t, ui.Field("", nil, gsx.Attrs{{Key: "id", Value: "f1"}}))
	if !strings.Contains(got, `id="f1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

func TestFieldCallerClassMerges(t *testing.T) {
	got := render(t, ui.Field("", nil, gsx.Attrs{{Key: "class", Value: "gap-8"}}))
	if strings.Count(got, "gap-8") != 1 || strings.Count(got, `class=`) != 1 {
		t.Errorf("caller class must be forwarded exactly once\nin: %s", got)
	}
	// gap-8 displaces the base rule's own gap-2 through tailwind-merge —
	// plain-vs-plain, so the caller wins before the cascade is consulted.
	if want := canonicalFieldClass("", []string{"orientation-vertical"}, "gap-8"); !strings.Contains(got, want) {
		t.Errorf("caller class not merged after the recipe's own\nwant: %s\nin: %s", want, got)
	}
}

func TestFieldContentPinned(t *testing.T) {
	got := render(t, ui.FieldContent(gsx.Raw("x"), nil))
	want := `<div ` + canonicalFieldClass("content", nil) + ` data-gsxui-slot-field-content>x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestFieldContentAttrsFallThrough(t *testing.T) {
	got := render(t, ui.FieldContent(nil, gsx.Attrs{{Key: "id", Value: "fc1"}}))
	if !strings.Contains(got, `id="fc1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

// TestFieldLabelPinned proves FieldLabel composes the Label token.
func TestFieldLabelPinned(t *testing.T) {
	got := render(t, ui.FieldLabel(gsx.Raw("x"), nil))
	want := `<label ` + canonicalFieldLabelClass() + ` data-gsxui-slot-field-label data-gsxui-slot-label>x</label>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestFieldLabelAttrsFallThrough(t *testing.T) {
	got := render(t, ui.FieldLabel(nil, gsx.Attrs{{Key: "for", Value: "email"}}))
	if !strings.Contains(got, `for="email"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

// FieldTitle has its own styling token.
func TestFieldTitlePinned(t *testing.T) {
	got := render(t, ui.FieldTitle(gsx.Raw("x"), nil))
	want := `<div ` + canonicalFieldClass("title", nil) + ` data-gsxui-slot-field-title>x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestFieldTitleAttrsFallThrough(t *testing.T) {
	got := render(t, ui.FieldTitle(nil, gsx.Attrs{{Key: "id", Value: "ft1"}}))
	if !strings.Contains(got, `id="ft1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

func TestFieldDescriptionPinned(t *testing.T) {
	got := render(t, ui.FieldDescription(gsx.Raw("x"), nil))
	want := `<p ` + canonicalFieldClass("description", nil) + ` data-gsxui-slot-field-description>x</p>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestFieldDescriptionAttrsFallThrough(t *testing.T) {
	got := render(t, ui.FieldDescription(nil, gsx.Attrs{{Key: "id", Value: "fd1"}}))
	if !strings.Contains(got, `id="fd1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

// FieldSeparator keeps wrapper structure separate from its composed
// Separator token.
func TestFieldSeparatorNoChildrenPinned(t *testing.T) {
	got := render(t, ui.FieldSeparator(nil, nil))
	want := `<div ` + canonicalFieldClass("separator-wrapper", nil) + ` data-gsxui-slot-field-separator-wrapper><div role="none" data-orientation="horizontal" ` + canonicalFieldSeparatorClass() + ` data-gsxui-slot-field-separator data-gsxui-slot-separator></div></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestFieldSeparatorWithChildrenPinned(t *testing.T) {
	got := render(t, ui.FieldSeparator(gsx.Raw("Or"), nil))
	want := `<div data-content ` + canonicalFieldClass("separator-wrapper", nil) + ` data-gsxui-slot-field-separator-wrapper><div role="none" data-orientation="horizontal" ` + canonicalFieldSeparatorClass() + ` data-gsxui-slot-field-separator data-gsxui-slot-separator></div><span ` + canonicalFieldClass("separator-content", nil) + ` data-gsxui-slot-field-separator-content>Or</span></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestFieldSeparatorAttrsFallThrough(t *testing.T) {
	got := render(t, ui.FieldSeparator(nil, gsx.Attrs{{Key: "id", Value: "fsep1"}}))
	if !strings.Contains(got, `id="fsep1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

// TestFieldErrorNilRendersNothing is the gsx equivalent of shadcn's `if
// (!content) return null` — now driven by children alone (the errors prop
// ADAPT, see the file-level comment in ui/field.gsx).
func TestFieldErrorNilRendersNothing(t *testing.T) {
	got := render(t, ui.FieldError(nil, nil))
	if got != "" {
		t.Errorf("want empty render for nil children, got %q", got)
	}
}

func TestFieldErrorPinned(t *testing.T) {
	got := render(t, ui.FieldError(gsx.Raw("This field is required."), nil))
	want := `<div role="alert" ` + canonicalFieldClass("error", nil) + ` data-gsxui-slot-field-error>This field is required.</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestFieldErrorAttrsFallThrough(t *testing.T) {
	got := render(t, ui.FieldError(gsx.Raw("x"), gsx.Attrs{{Key: "id", Value: "fe1"}}))
	if !strings.Contains(got, `id="fe1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

// Realistic composition: a small form — FieldSet + FieldGroup with two
// vertical Fields (label/content/description), the site example's own
// shape.
func TestFieldFormComposition(t *testing.T) {
	got := render(t, ui.FieldSet(
		gsx.Fragment(
			ui.FieldLegend("", gsx.Raw("Profile"), nil),
			ui.FieldGroup(
				gsx.Fragment(
					ui.Field("", gsx.Fragment(
						ui.FieldLabel(gsx.Raw("Name"), gsx.Attrs{{Key: "for", Value: "name"}}),
						gsx.Raw(`<input id="name"/>`),
						ui.FieldDescription(gsx.Raw("Your full name."), nil),
					), nil),
					ui.FieldSeparator(nil, nil),
					ui.Field("", gsx.Fragment(
						ui.FieldLabel(gsx.Raw("Email"), gsx.Attrs{{Key: "for", Value: "email"}}),
						gsx.Raw(`<input id="email" aria-invalid="true"/>`),
						ui.FieldError(gsx.Raw("Enter a valid email."), nil),
					), nil),
				),
				nil,
			),
		),
		nil,
	))
	for _, want := range []string{
		`data-gsxui-slot-field-set`,
		`data-gsxui-slot-field-legend`,
		`>Profile</legend>`,
		`data-gsxui-slot-field-group`,
		`data-gsxui-slot-field`,
		`data-gsxui-slot-field-label data-gsxui-slot-label`,
		`>Name</label>`,
		`for="name"`,
		`data-gsxui-slot-field-description`,
		`>Your full name.</p>`,
		`data-gsxui-slot-field-separator-wrapper`,
		`>Email</label>`,
		`for="email"`,
		`data-gsxui-slot-field-error`,
		`>Enter a valid email.</div>`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}
