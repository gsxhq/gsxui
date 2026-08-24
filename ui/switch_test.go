package ui_test

import (
	"html"
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

// canonicalSwitchClass is Switch's fully-resolved recipe class, derived from
// the default style's recipe CSS (registry/styles/nova/switch.css) the same
// way TestCheckboxPinned derives Checkbox's — asserting ui/switch.gsx stays
// in lockstep with the registry instead of re-typing a literal on every
// style edit. Switch migrated onto the slot axis, so its root class is no
// longer empty.
//
// Switch is a real native <input type="checkbox" role="switch">: the
// data-checked/data-unchecked attribute-variant form the 8-style port
// mechanically carried over from upstream's Radix vocabulary was dead on
// this markup (recurring class 2 in the style-porter report) and is
// restored here to the native :checked pseudo-class, matching
// ui/switch.gsx's own doc comment. The thumb is this element's own ::before
// pseudo-element (no separate sibling/group-data-[…]/switch: relationship —
// Switch has no data-size axis at all, a deliberate simplification from the
// mechanical port's dead group/switch-gated size variants), so
// before:content-[&#39;&#39;]/before:pointer-events-none/before:block/
// before:transition-transform and appearance-none/outline-none/
// disabled:cursor-not-allowed/disabled:opacity-50 are restored structural
// chrome — see the report's "Switch — full structural rewrite" entry.
func canonicalSwitchClass() string {
	return html.EscapeString(strings.Join(styleRecipeUtilities("switch", "gsxui-recipe-switch"), " "))
}

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
	// part of the pinned render, derived from the default style's recipe CSS.
	got := render(t, ui.Switch(nil))
	want := `<input type="checkbox" role="switch" class="` + canonicalSwitchClass() + `" data-gsxui-slot-switch>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
	// The data-checked/data-unchecked attribute-variant vocabulary is dead
	// on this native markup (see canonicalSwitchClass's comment); a derived
	// expectation can't pin the translation, so reject it directly.
	if strings.Contains(got, "data-checked") || strings.Contains(got, "data-unchecked") {
		t.Errorf("dead data-checked/data-unchecked vocabulary in render\nin: %s", got)
	}
}

func TestSwitchDarkCheckedOverride(t *testing.T) {
	// Switch's dark-mode/checked presentation now lives in the resolved
	// recipe class itself (dark:, data-checked:, dark:data-checked:
	// variants) — this test used to assert those tokens were absent from a
	// marker-only render; now it just confirms they're present, since
	// they're supposed to be part of Switch's own migrated presentation.
	got := render(t, ui.Switch(nil))
	if !strings.Contains(got, "dark:checked:before:bg-primary-foreground") {
		t.Errorf("missing dark:checked:before presentation\nin: %s", got)
	}
}
