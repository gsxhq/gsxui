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
		`data-slot="switch"`,
		"peer inline-flex shrink-0 items-center appearance-none rounded-full",
		"h-[1.15rem] w-8",
		"bg-input checked:bg-primary dark:bg-input/80",
		"disabled:cursor-not-allowed disabled:opacity-50",
		"before:block before:size-4 before:rounded-full before:bg-background",
		"before:content-[&#39;&#39;]",
		"checked:before:translate-x-[calc(100%-2px)]",
		"/>",
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
	if !strings.Contains(got, `<input type="checkbox" role="switch" data-slot="switch"`) {
		t.Errorf("missing role=\"switch\" in expected position\nin: %s", got)
	}
}

func TestSwitchCallerClassMerges(t *testing.T) {
	got := render(t, ui.Switch(gsx.Attrs{{Key: "class", Value: "w-12"}}))
	if strings.Contains(got, "w-8") {
		t.Errorf("base w-8 should be dropped by caller w-12\nin: %s", got)
	}
	if !strings.Contains(got, "w-12") {
		t.Errorf("missing caller class w-12\nin: %s", got)
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
	// Exact full-render pin. Token-by-token verified against shadcn's
	// Switch (registry/new-york-v4/ui/switch.tsx) plus the ledgered
	// ADAPTs: native <input type="checkbox" role="switch"> replaces Root,
	// and the Thumb span becomes a before: pseudo-element on the same
	// input. shadcn's Thumb has no shadow at all (ring-0, nothing else) —
	// before:shadow-lg is not sourced from anywhere and is intentionally
	// absent here. See docs/jsx-parity.md.
	got := render(t, ui.Switch(nil))
	want := `<input type="checkbox" role="switch" data-slot="switch" class="peer inline-flex shrink-0 items-center appearance-none rounded-full border border-transparent shadow-xs transition-all outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 h-[1.15rem] w-8 bg-input checked:bg-primary dark:bg-input/80 before:pointer-events-none before:block before:size-4 before:rounded-full before:bg-background before:transition-transform before:content-[&#39;&#39;] checked:before:translate-x-[calc(100%-2px)] dark:before:bg-foreground dark:checked:before:bg-primary-foreground"/>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}
