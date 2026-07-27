package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

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
	if !strings.Contains(got, `<input type="checkbox" role="switch" data-gsxui-slot-switch`) {
		t.Errorf("missing role=\"switch\" in expected position\nin: %s", got)
	}
}

func TestSwitchCallerClassMerges(t *testing.T) {
	got := render(t, ui.Switch(gsx.Attrs{{Key: "class", Value: "w-12"}}))
	if strings.Count(got, `class="w-12"`) != 1 {
		t.Errorf("caller class must be forwarded exactly once\nin: %s", got)
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
	// Presentation lives in the stylesheet; the render pin covers structure.
	got := render(t, ui.Switch(nil))
	want := `<input type="checkbox" role="switch" data-gsxui-slot-switch>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestSwitchDarkCheckedOverride(t *testing.T) {
	got := render(t, ui.Switch(nil))
	if strings.Contains(got, "dark:") {
		t.Errorf("presentation classes must not be rendered\nin: %s", got)
	}
}
