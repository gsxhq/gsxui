package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

// canonicalSwitchClass is Switch's fully-resolved recipe class, copied
// verbatim from the generated ui/switch.gsx output — Switch migrated onto
// the slot axis, so its root class is no longer empty. render() HTML-escapes
// attribute values, so "'" in before:content-[”] becomes "&#39;" below.
const canonicalSwitchClass = `inline-flex h-[1.15rem] w-8 shrink-0 appearance-none items-center rounded-full border border-transparent bg-input transition-all outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 checked:bg-primary dark:bg-input/80 dark:checked:bg-primary before:content-[&#39;&#39;] before:pointer-events-none before:block before:size-4 before:rounded-full before:bg-background before:transition-transform checked:before:translate-x-[calc(100%-2px)] dark:before:bg-foreground dark:checked:before:bg-primary-foreground`

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
	// recipe class itself (dark:, checked:, dark:checked: variants) — this
	// test used to assert those tokens were absent from a marker-only
	// render; now it just confirms they're present, since they're supposed
	// to be part of Switch's own migrated presentation.
	got := render(t, ui.Switch(nil))
	if !strings.Contains(got, "dark:checked:before:bg-primary-foreground") {
		t.Errorf("missing dark:checked:before presentation\nin: %s", got)
	}
}
