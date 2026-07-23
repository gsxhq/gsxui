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
		`data-gsxui-tabs`,         // root hook
		`data-value="a"`,          // root's initial value
		`role="tablist"`,          // list
		`data-gsxui-tabs-trigger`, // trigger hook
		`role="tab"`,              // trigger role
		`role="tabpanel"`,         // content role
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
	if strings.Contains(got, "gap-1.5") {
		t.Errorf("base gap-1.5 should be dropped by caller gap-4\nin: %s", got)
	}
	if !strings.Contains(got, "gap-4") {
		t.Errorf("missing caller class gap-4\nin: %s", got)
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
	// Exact full-render pin for the active TabsTrigger, verified token-by-
	// token against shadcn's TabsTrigger (registry/new-york-v4/ui/tabs.tsx)
	// with the orientation/variant ADAPTs from docs/jsx-parity.md applied:
	// group-data-[orientation=vertical]/tabs: tokens dropped (no orientation
	// param), group-data-[variant=*]/tabs-list: unwrapped to the
	// unconditional data-[state=active]:shadow-sm (no variant param, default
	// is the only shipped variant), the after: line-indicator dropped
	// entirely (invisible under the default variant).
	got := render(t, ui.TabsTrigger("a", true, gsx.Raw("Account"), nil))
	want := `<button type="button" role="tab" data-slot="tabs-trigger" data-gsxui-tabs-trigger data-value="a" data-state="active" aria-selected="true" tabindex="0" class="relative inline-flex h-[calc(100%-1px)] flex-1 items-center justify-center gap-1.5 rounded-md border border-transparent px-2 py-1 text-sm font-medium whitespace-nowrap text-foreground/60 transition-all hover:text-foreground focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 focus-visible:outline-1 focus-visible:outline-ring disabled:pointer-events-none disabled:opacity-50 data-[state=active]:shadow-sm dark:text-muted-foreground dark:hover:text-foreground [&amp;_svg]:pointer-events-none [&amp;_svg]:shrink-0 [&amp;_svg:not([class*=&#39;size-&#39;])]:size-4 data-[state=active]:bg-background data-[state=active]:text-foreground dark:data-[state=active]:border-input dark:data-[state=active]:bg-input/30 dark:data-[state=active]:text-foreground">Account</button>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}
