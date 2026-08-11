package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

// canonicalSwitchClass is Switch's fully-resolved recipe class, copied
// verbatim from the generated ui/switch.gsx output — Switch migrated onto
// the slot axis, so its root class is no longer empty. Nova's ported recipe
// targets checked/unchecked state via data-checked/data-unchecked attribute
// variants rather than the native :checked pseudo-class (the JS that keeps
// those attributes in sync with the input's live checked state lives in
// web/, handled separately from this style-porting task).
const canonicalSwitchClass = `data-checked:bg-primary data-unchecked:bg-input focus-visible:border-ring focus-visible:ring-ring/50 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive dark:aria-invalid:border-destructive/50 dark:data-unchecked:bg-input/80 shrink-0 rounded-full border border-transparent focus-visible:ring-3 aria-invalid:ring-3 data-[size=default]:h-[18.4px] data-[size=default]:w-[32px] data-[size=sm]:h-[14px] data-[size=sm]:w-[24px] before:bg-background dark:data-unchecked:before:bg-foreground dark:data-checked:before:bg-primary-foreground before:rounded-full group-data-[size=default]/switch:before:size-4 group-data-[size=sm]/switch:before:size-3 group-data-[size=default]/switch:data-checked:before:translate-x-[calc(100%-2px)] group-data-[size=sm]/switch:data-checked:before:translate-x-[calc(100%-2px)] group-data-[size=default]/switch:data-unchecked:before:translate-x-0 group-data-[size=sm]/switch:data-unchecked:before:translate-x-0 inline-flex`

func TestSwitchDefault(t *testing.T) {
	got := render(t, ui.Switch(nil))
	for _, want := range []string{
		`<input type="checkbox"`,
		`role="switch"`,
		`data-gsxui-slot-switch`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestSwitchRolePin(t *testing.T) {
	// role="switch" on the native <input type="checkbox"> is the load-
	// bearing a11y contract standing in for Radix's SwitchPrimitive.Root,
	// which stamps role="switch" itself. Pinned separately from the full
	// render pin so a future edit can't silently drop it.
	got := render(t, ui.Switch(nil))
	if !strings.Contains(got, `<input type="checkbox" role="switch" class="`) {
		t.Errorf("missing role=\"switch\" in expected position\nin: %s", got)
	}
}

func TestSwitchCallerClassMerges(t *testing.T) {
	got := render(t, ui.Switch(gsx.Attrs{{Key: "class", Value: "w-12"}}))
	if strings.Count(got, "w-12") != 1 || strings.Count(got, `class="`) != 1 {
		t.Errorf("caller class must be forwarded exactly once, merged with the recipe class\nin: %s", got)
	}
}

func TestSwitchAttrsFallThrough(t *testing.T) {
	got := render(t, ui.Switch(gsx.Attrs{{Key: "id", Value: "s1"}, {Key: "name", Value: "notify"}, {Key: "aria-label", Value: "Notifications"}}))
	for _, want := range []string{`id="s1"`, `name="notify"`, `aria-label="Notifications"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestSwitchCheckedAttr(t *testing.T) {
	got := render(t, ui.Switch(gsx.Attrs{{Key: "checked", Value: true}}))
	if !strings.Contains(got, " checked") || strings.Contains(got, `checked="`) {
		t.Errorf("checked attr should render bare, not stringified\nin: %s", got)
	}

	got = render(t, ui.Switch(gsx.Attrs{{Key: "checked", Value: false}}))
	if strings.Contains(got, `" checked`) || strings.Contains(got, `checked="false"`) {
		t.Errorf("checked=false should omit the attribute entirely\nin: %s", got)
	}
}

func TestSwitchDisabledAttr(t *testing.T) {
	got := render(t, ui.Switch(gsx.Attrs{{Key: "disabled", Value: true}}))
	if !strings.Contains(got, " disabled") || strings.Contains(got, `disabled="`) {
		t.Errorf("disabled attr should render bare\nin: %s", got)
	}
}

func TestSwitchPinned(t *testing.T) {
	// Switch migrated onto the slot axis: the resolved recipe class is now
	// part of the pinned render, copied verbatim from generated output.
	got := render(t, ui.Switch(nil))
	want := `<input type="checkbox" role="switch" class="` + canonicalSwitchClass + `" data-gsxui-slot-switch>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestSwitchDarkCheckedOverride(t *testing.T) {
	// Switch's dark-mode/checked presentation now lives in the resolved
	// recipe class itself (dark:, data-checked:, dark:data-checked:
	// variants) — this test used to assert those tokens were absent from a
	// marker-only render; now it just confirms they're present, since
	// they're supposed to be part of Switch's own migrated presentation.
	got := render(t, ui.Switch(nil))
	if !strings.Contains(got, "dark:data-checked:before:bg-primary-foreground") {
		t.Errorf("missing dark:data-checked:before presentation\nin: %s", got)
	}
}
