package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestFieldSetPinned(t *testing.T) {
	got := render(t, ui.FieldSet(gsx.Raw("x"), nil))
	want := `<fieldset data-gsxui-slot-field-set>x</fieldset>`
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
	want := `<legend data-variant="legend" data-gsxui-slot-field-legend>x</legend>`
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
	want := `<div data-gsxui-slot-field-group>x</div>`
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
	want := `<div role="group" data-orientation="vertical" data-gsxui-slot-field>x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

// TestFieldHorizontalPinned proves data-orientation remains load-bearing.
func TestFieldHorizontalPinned(t *testing.T) {
	got := render(t, ui.Field("horizontal", gsx.Raw("x"), nil))
	want := `<div role="group" data-orientation="horizontal" data-gsxui-slot-field>x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestFieldResponsivePinned(t *testing.T) {
	got := render(t, ui.Field("responsive", gsx.Raw("x"), nil))
	want := `<div role="group" data-orientation="responsive" data-gsxui-slot-field>x</div>`
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
	if strings.Count(got, `class="gap-8"`) != 1 {
		t.Errorf("caller class must be forwarded exactly once\nin: %s", got)
	}
}

func TestFieldContentPinned(t *testing.T) {
	got := render(t, ui.FieldContent(gsx.Raw("x"), nil))
	want := `<div data-gsxui-slot-field-content>x</div>`
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
	want := `<label data-gsxui-slot-field-label data-gsxui-slot-label>x</label>`
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
	want := `<div data-gsxui-slot-field-title>x</div>`
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
	want := `<p data-gsxui-slot-field-description>x</p>`
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
	want := `<div data-content="false" data-gsxui-slot-field-separator-wrapper><div role="none" data-orientation="horizontal" data-gsxui-slot-field-separator data-gsxui-slot-separator></div></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestFieldSeparatorWithChildrenPinned(t *testing.T) {
	got := render(t, ui.FieldSeparator(gsx.Raw("Or"), nil))
	want := `<div data-content="true" data-gsxui-slot-field-separator-wrapper><div role="none" data-orientation="horizontal" data-gsxui-slot-field-separator data-gsxui-slot-separator></div><span data-gsxui-slot-field-separator-content>Or</span></div>`
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
	want := `<div role="alert" data-gsxui-slot-field-error>This field is required.</div>`
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
