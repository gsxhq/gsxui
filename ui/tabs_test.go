package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestTabsStructure(t *testing.T) {
	got := render(t, ui.Tabs("a", gsx.Fragment(
		ui.TabsList(gsx.Fragment(
			ui.TabsTrigger("a", true, gsx.Raw("Tab A"), nil),
			ui.TabsTrigger("b", false, gsx.Raw("Tab B"), nil),
		), nil),
		ui.TabsContent("a", true, gsx.Raw("Content A"), nil),
		ui.TabsContent("b", false, gsx.Raw("Content B"), nil),
	), nil))
	for _, want := range []string{
		`data-gsxui-slot-tabs`,         // root hook
		`data-value="a"`,               // root's initial value
		`role="tablist"`,               // list
		`data-gsxui-slot-tabs-trigger`, // trigger hook
		`role="tab"`,                   // trigger role
		`role="tabpanel"`,              // content role
		">Tab A<", ">Tab B<",
		">Content A<", ">Content B<",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

// TestTabsTriggerSelectedStamping covers the explicit-selected GAP: the
// caller resolves value==root-value and passes the bool in; the component
// only stamps the result. Both branches must be exercised since the zero
// value (false) must never accidentally read as active.
func TestTabsTriggerSelectedStamping(t *testing.T) {
	active := render(t, ui.TabsTrigger("a", true, gsx.Raw("x"), nil))
	for _, want := range []string{
		`data-state="active"`,
		`aria-selected="true"`,
		`tabindex="0"`,
	} {
		if !strings.Contains(active, want) {
			t.Errorf("selected trigger missing %q\nin: %s", want, active)
		}
	}

	inactive := render(t, ui.TabsTrigger("a", false, gsx.Raw("x"), nil))
	for _, want := range []string{
		`data-state="inactive"`,
		`aria-selected="false"`,
		`tabindex="-1"`,
	} {
		if !strings.Contains(inactive, want) {
			t.Errorf("unselected trigger missing %q\nin: %s", inactive, want)
		}
	}
}

func TestTabsContentSelectedStamping(t *testing.T) {
	active := render(t, ui.TabsContent("a", true, gsx.Raw("x"), nil))
	if !strings.Contains(active, `data-state="active"`) {
		t.Errorf("selected content missing data-state=active\nin: %s", active)
	}
	if strings.Contains(active, "hidden") {
		t.Errorf("selected content must not be hidden\nin: %s", active)
	}

	inactive := render(t, ui.TabsContent("a", false, gsx.Raw("x"), nil))
	if !strings.Contains(inactive, `data-state="inactive"`) {
		t.Errorf("unselected content missing data-state=inactive\nin: %s", inactive)
	}
	if !strings.Contains(inactive, "hidden") {
		t.Errorf("unselected content must be hidden\nin: %s", inactive)
	}
}

func TestTabsCallerClassMerges(t *testing.T) {
	got := render(t, ui.TabsTrigger("a", false, gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "gap-4"}}))
	if strings.Count(got, "gap-4") != 1 || strings.Count(got, `class="`) != 1 {
		t.Errorf("caller class must be forwarded exactly once, merged with the recipe class\nin: %s", got)
	}
}

func TestTabsAttrsFallThrough(t *testing.T) {
	got := render(t, ui.Tabs("a", gsx.Raw("x"), gsx.Attrs{{Key: "id", Value: "t1"}, {Key: "aria-label", Value: "settings"}}))
	for _, want := range []string{`id="t1"`, `aria-label="settings"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestTabsTriggerPinned(t *testing.T) {
	got := render(t, ui.TabsTrigger("a", true, gsx.Raw("Account"), nil))
	want := `<button type="button" role="tab" data-value="a" data-state="active" aria-selected="true" tabindex="0" class="gap-1.5 rounded-md border border-transparent px-1.5 py-0.5 text-sm font-medium group-data-[variant=default]/tabs-list:data-active:shadow-sm group-data-[variant=line]/tabs-list:data-active:shadow-none [&amp;_svg:not([class*=&#39;size-&#39;])]:size-4 has-data-[icon=inline-end]:pr-1 has-data-[icon=inline-start]:pl-1 inline-flex" data-gsxui-slot-tabs-trigger>Account</button>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}
