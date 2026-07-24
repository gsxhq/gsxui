package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestFieldSetPinned(t *testing.T) {
	got := render(t, ui.FieldSet(gsx.Raw("x"), nil))
	want := `<fieldset data-slot="field-set" class="flex flex-col gap-4 has-[&gt;[data-slot=checkbox-group]]:gap-3 has-[&gt;[data-slot=radio-group]]:gap-3">x</fieldset>`
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

// TestFieldLegendDefaultPinned pins the zero-value ("legend") variant. Both
// data-[variant=legend]:text-base and data-[variant=label]:text-sm are
// always present in the class string (WIN: no switch needed, the same
// data-attribute-driven-selector shape as Separator's data-orientation, not
// a static-block cva switch — see ui/field.gsx's own comment).
func TestFieldLegendDefaultPinned(t *testing.T) {
	got := render(t, ui.FieldLegend("", gsx.Raw("x"), nil))
	want := `<legend data-slot="field-legend" data-variant="legend" class="mb-1.5 font-medium data-[variant=legend]:text-base data-[variant=label]:text-sm">x</legend>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestFieldLegendLabelVariantPinned(t *testing.T) {
	got := render(t, ui.FieldLegend("label", gsx.Raw("x"), nil))
	if !strings.Contains(got, `data-variant="label"`) {
		t.Errorf("missing data-variant=label\nin: %s", got)
	}
	if !strings.Contains(got, "data-[variant=label]:text-sm") {
		t.Errorf("missing data-[variant=label]:text-sm selector\nin: %s", got)
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
	want := `<div data-slot="field-group" class="group/field-group @container/field-group flex w-full flex-col gap-5 data-[slot=checkbox-group]:gap-3 [&amp;&gt;[data-slot=field-group]]:gap-4">x</div>`
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
	want := `<div role="group" data-slot="field" data-orientation="vertical" class="group/field flex w-full gap-2 data-[invalid=true]:text-destructive flex-col [&amp;&gt;*]:w-full [&amp;&gt;.sr-only]:w-auto">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

// TestFieldHorizontalPinned also proves data-orientation is genuinely
// load-bearing: FieldDescription's own group-has-[[data-orientation=
// horizontal]]/field:text-balance selector (see its own test) depends on
// this exact attribute value existing on an ancestor.
func TestFieldHorizontalPinned(t *testing.T) {
	got := render(t, ui.Field("horizontal", gsx.Raw("x"), nil))
	want := `<div role="group" data-slot="field" data-orientation="horizontal" class="group/field flex w-full gap-2 data-[invalid=true]:text-destructive flex-row items-center [&amp;&gt;[data-slot=field-label]]:flex-auto has-[&gt;[data-slot=field-content]]:items-start has-[&gt;[data-slot=field-content]]:[&amp;&gt;[role=checkbox],[role=radio]]:mt-px">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestFieldResponsivePinned(t *testing.T) {
	got := render(t, ui.Field("responsive", gsx.Raw("x"), nil))
	want := `<div role="group" data-slot="field" data-orientation="responsive" class="group/field flex w-full gap-2 data-[invalid=true]:text-destructive flex-col @md/field-group:flex-row @md/field-group:items-center [&amp;&gt;*]:w-full @md/field-group:[&amp;&gt;*]:w-auto [&amp;&gt;.sr-only]:w-auto @md/field-group:[&amp;&gt;[data-slot=field-label]]:flex-auto @md/field-group:has-[&gt;[data-slot=field-content]]:items-start @md/field-group:has-[&gt;[data-slot=field-content]]:[&amp;&gt;[role=checkbox],[role=radio]]:mt-px">x</div>`
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
	if !strings.Contains(got, "gap-8") {
		t.Errorf("missing caller class gap-8\nin: %s", got)
	}
}

func TestFieldContentPinned(t *testing.T) {
	got := render(t, ui.FieldContent(gsx.Raw("x"), nil))
	want := `<div data-slot="field-content" class="group/field-content flex flex-1 flex-col gap-0.5 leading-snug">x</div>`
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

// TestFieldLabelPinned proves FieldLabel actually composes ui.Label
// (Label's own base classes come through) and that data-slot is overridden
// from Label's own "label" to "field-label", and that leading-snug (this
// overlay) wins its tailwind-merge conflict against Label's own
// leading-none.
func TestFieldLabelPinned(t *testing.T) {
	got := render(t, ui.FieldLabel(gsx.Raw("x"), nil))
	want := `<label class="items-center text-sm font-medium select-none group-data-[disabled=true]:pointer-events-none group-data-[disabled=true]:opacity-50 peer-disabled:cursor-not-allowed peer-disabled:opacity-50 group/field-label peer/field-label flex w-fit gap-2 leading-snug group-data-[disabled=true]/field:opacity-50 has-[&gt;[data-slot=field]]:w-full has-[&gt;[data-slot=field]]:flex-col has-[&gt;[data-slot=field]]:rounded-lg has-[&gt;[data-slot=field]]:border [&amp;&gt;*]:data-[slot=field]:p-2.5 has-data-[state=checked]:border-primary has-data-[state=checked]:bg-primary/5 dark:has-data-[state=checked]:bg-primary/10" data-slot="field-label">x</label>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
	if strings.Contains(got, "leading-none") {
		t.Errorf("leading-none should be dropped by leading-snug\nin: %s", got)
	}
}

func TestFieldLabelAttrsFallThrough(t *testing.T) {
	got := render(t, ui.FieldLabel(nil, gsx.Attrs{{Key: "for", Value: "email"}}))
	if !strings.Contains(got, `for="email"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

// TestFieldTitlePinned proves FieldTitle shares FieldLabel's data-slot
// ("field-label") — a shadcn source verbatim quirk, ported as-is (see
// ui/field.gsx's own comment).
func TestFieldTitlePinned(t *testing.T) {
	got := render(t, ui.FieldTitle(gsx.Raw("x"), nil))
	want := `<div data-slot="field-label" class="flex w-fit items-center gap-2 text-sm leading-snug font-medium group-data-[disabled=true]/field:opacity-50">x</div>`
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
	want := `<p data-slot="field-description" class="text-sm leading-normal font-normal text-muted-foreground group-has-[[data-orientation=horizontal]]/field:text-balance last:mt-0 [[data-variant=legend]+&amp;]:-mt-1.5 [&amp;&gt;a]:underline [&amp;&gt;a]:underline-offset-4 [&amp;&gt;a:hover]:text-primary">x</p>`
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

// TestFieldSeparatorNoChildrenPinned proves FieldSeparator composes
// ui.Separator (role="none" and Separator's own data-[orientation=...] base
// classes both come through, horizontal default) and stamps
// data-content="false" — gsx's bool-to-"true"/"false" attribute rendering,
// the same mechanism as pagination.gsx's data-active (see ui/pagination.gsx)
// — when no children are given; no field-separator-content span renders.
func TestFieldSeparatorNoChildrenPinned(t *testing.T) {
	got := render(t, ui.FieldSeparator(nil, nil))
	want := `<div data-slot="field-separator" data-content="false" class="relative -my-2 h-5 text-sm group-data-[variant=outline]/field-group:-mb-2"><div data-slot="separator" role="none" data-orientation="horizontal" class="shrink-0 bg-border data-[orientation=horizontal]:h-px data-[orientation=horizontal]:w-full data-[orientation=vertical]:h-full data-[orientation=vertical]:w-px absolute inset-0 top-1/2"></div></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestFieldSeparatorWithChildrenPinned(t *testing.T) {
	got := render(t, ui.FieldSeparator(gsx.Raw("Or"), nil))
	want := `<div data-slot="field-separator" data-content="true" class="relative -my-2 h-5 text-sm group-data-[variant=outline]/field-group:-mb-2"><div data-slot="separator" role="none" data-orientation="horizontal" class="shrink-0 bg-border data-[orientation=horizontal]:h-px data-[orientation=horizontal]:w-full data-[orientation=vertical]:h-full data-[orientation=vertical]:w-px absolute inset-0 top-1/2"></div><span class="relative mx-auto block w-fit bg-background px-2 text-muted-foreground" data-slot="field-separator-content">Or</span></div>`
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
	want := `<div role="alert" data-slot="field-error" class="text-sm font-normal text-destructive">This field is required.</div>`
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
		`data-slot="field-set"`,
		`data-slot="field-legend"`,
		`>Profile</legend>`,
		`data-slot="field-group"`,
		`data-slot="field"`,
		`data-slot="field-label"`,
		`>Name</label>`,
		`for="name"`,
		`data-slot="field-description"`,
		`>Your full name.</p>`,
		`data-slot="field-separator"`,
		`>Email</label>`,
		`for="email"`,
		`data-slot="field-error"`,
		`>Enter a valid email.</div>`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}
