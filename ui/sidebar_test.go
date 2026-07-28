package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func requireContainsAll(t *testing.T, got string, wants ...string) {
	t.Helper()
	for _, want := range wants {
		if !strings.Contains(got, want) {
			t.Errorf("render missing %q\nin: %s", want, got)
		}
	}
}

func TestSidebarProviderReflectsStateWidthsAndBehaviorHook(t *testing.T) {
	open := render(t, ui.SidebarProvider(true, gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "caller"}}))
	requireContainsAll(t, open,
		`data-state="expanded"`,
		`data-gsxui-sidebar-wrapper`,
		`data-gsxui-slot-sidebar-wrapper`,
		`--sidebar-width:16rem`,
		`--sidebar-width-icon:3rem`,
		`class="caller"`,
	)
	if strings.Count(open, `class="caller"`) != 1 {
		t.Fatalf("caller class count = %d, want 1\nin: %s", strings.Count(open, `class="caller"`), open)
	}

	closed := render(t, ui.SidebarProvider(false, gsx.Raw("x"), nil))
	requireContainsAll(t, closed, `data-state="collapsed"`, `data-gsxui-sidebar-wrapper`)
}

func TestSidebarRendersTwoExplicitResponsiveTrees(t *testing.T) {
	got := render(t, ui.Sidebar(
		false,
		"right",
		"floating",
		"icon",
		gsx.Raw("CONTENT"),
		gsx.Attrs{{Key: "class", Value: "caller"}},
	))

	if count := strings.Count(got, "CONTENT"); count != 2 {
		t.Fatalf("children count = %d, want mobile and desktop copies\nin: %s", count, got)
	}
	requireContainsAll(t, got,
		`data-gsxui-slot-sidebar-mobile-root data-gsxui-slot-sheet data-gsxui-slot-dialog`,
		`data-gsxui-slot-sidebar-mobile-content data-gsxui-slot-sidebar data-gsxui-slot-sheet-content data-gsxui-slot-dialog-content`,
		`data-mobile="true"`,
		`data-side="right"`,
		`data-gsxui-sidebar-mobile-dialog`,
		`--sidebar-width:18rem`,
		`data-gsxui-slot-sidebar-mobile-header data-gsxui-slot-sheet-header`,
		`data-gsxui-slot-sidebar-mobile-title data-gsxui-slot-sheet-title`,
		`data-gsxui-slot-sidebar-mobile-description data-gsxui-slot-sheet-description`,
		`data-gsxui-slot-sidebar-mobile-inner`,
		`data-gsxui-slot-sidebar-desktop data-gsxui-slot-sidebar`,
		`data-state="collapsed"`,
		`data-collapsible="icon"`,
		`data-gsxui-sidebar-collapsible="icon"`,
		`data-variant="floating"`,
		`data-gsxui-sidebar-desktop`,
		`data-gsxui-slot-sidebar-gap`,
		`data-gsxui-slot-sidebar-container`,
		`data-gsxui-slot-sidebar-inner`,
	)
	if count := strings.Count(got, `class="caller"`); count != 2 {
		t.Fatalf("caller class count = %d, want one on each responsive tree\nin: %s", count, got)
	}
	if strings.Contains(got, `data-slot=`) || strings.Contains(got, `data-sidebar=`) {
		t.Fatalf("legacy styling hooks remain\nin: %s", got)
	}
}

func TestSidebarExpandedDesktopClearsActiveCollapsibleMode(t *testing.T) {
	got := render(t, ui.Sidebar(true, "", "", "icon", gsx.Raw("x"), nil))
	requireContainsAll(t, got,
		`data-gsxui-slot-sidebar-desktop data-gsxui-slot-sidebar`,
		`data-state="expanded"`,
		`data-collapsible=""`,
		`data-gsxui-sidebar-collapsible="icon"`,
		`data-variant="sidebar"`,
		`data-side="left"`,
	)
}

func TestSidebarCollapsibleNoneIsOneCanonicalFlatTree(t *testing.T) {
	got := render(t, ui.Sidebar(
		false,
		"right",
		"inset",
		"none",
		gsx.Raw("CONTENT"),
		gsx.Attrs{{Key: "class", Value: "caller"}},
	))
	if count := strings.Count(got, "CONTENT"); count != 1 {
		t.Fatalf("children count = %d, want 1\nin: %s", count, got)
	}
	requireContainsAll(t, got,
		`data-gsxui-slot-sidebar`,
		`data-side="right"`,
		`data-variant="inset"`,
		`data-collapsible="none"`,
		`class="caller"`,
	)
	if strings.Contains(got, `sidebar-mobile-root`) || strings.Contains(got, `sidebar-desktop`) {
		t.Fatalf("collapsible=none rendered a responsive state-machine tree\nin: %s", got)
	}
}

func TestSidebarPrimitiveCompositionAndCallerClassPlacement(t *testing.T) {
	cases := []struct {
		name      string
		got       string
		want      string
		wantClass string
	}{
		{
			name:      "trigger button",
			got:       render(t, ui.SidebarTrigger(gsx.Attrs{{Key: "class", Value: "caller"}})),
			want:      `data-gsxui-slot-sidebar-trigger data-gsxui-slot-button`,
			wantClass: canonicalButtonClass("ghost", "icon", "caller"),
		},
		{
			name:      "input",
			got:       render(t, ui.SidebarInput(gsx.Attrs{{Key: "class", Value: "caller"}})),
			want:      `data-gsxui-slot-sidebar-input data-gsxui-slot-input`,
			wantClass: `class="caller"`,
		},
		{
			name:      "separator",
			got:       render(t, ui.SidebarSeparator(gsx.Attrs{{Key: "class", Value: "caller"}})),
			want:      `data-gsxui-slot-sidebar-separator data-gsxui-slot-separator`,
			wantClass: `class="caller"`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			requireContainsAll(t, tc.got, tc.want, tc.wantClass)
			if count := strings.Count(tc.got, tc.wantClass); count != 1 || strings.Count(tc.got, `class=`) != 1 {
				t.Fatalf("exact class count = %d, want 1 for %s\nin: %s", count, tc.wantClass, tc.got)
			}
			if strings.Contains(tc.got, `data-slot=`) || strings.Contains(tc.got, `data-sidebar=`) {
				t.Fatalf("legacy styling hook remains\nin: %s", tc.got)
			}
		})
	}

	trigger := render(t, ui.SidebarTrigger(nil))
	requireContainsAll(t, trigger,
		`data-variant="ghost"`,
		`data-size="icon"`,
		`data-gsxui-sidebar-trigger`,
		`data-gsxui-slot-icon`,
		`data-gsxui-slot-sidebar-trigger-label`,
	)
}

func TestSidebarTriggerComposesPresenceMarkersOnButton(t *testing.T) {
	got := render(t, ui.SidebarTrigger(nil))
	requirePresenceAttributesOnSameTag(t, got, "data-gsxui-sidebar-trigger",
		"data-gsxui-slot-sidebar-trigger",
		"data-gsxui-slot-button",
	)
}

func TestSidebarPlainPartsExposeNamespacedSlotsWithoutPresentationClasses(t *testing.T) {
	parts := []struct {
		name string
		got  string
		slot string
	}{
		{"rail", render(t, ui.SidebarRail(nil)), "sidebar-rail"},
		{"inset", render(t, ui.SidebarInset(gsx.Raw("x"), nil)), "sidebar-inset"},
		{"header", render(t, ui.SidebarHeader(gsx.Raw("x"), nil)), "sidebar-header"},
		{"footer", render(t, ui.SidebarFooter(gsx.Raw("x"), nil)), "sidebar-footer"},
		{"content", render(t, ui.SidebarContent(gsx.Raw("x"), nil)), "sidebar-content"},
		{"group", render(t, ui.SidebarGroup(gsx.Raw("x"), nil)), "sidebar-group"},
		{"group label", render(t, ui.SidebarGroupLabel(gsx.Raw("x"), nil)), "sidebar-group-label"},
		{"group action", render(t, ui.SidebarGroupAction(gsx.Raw("x"), nil)), "sidebar-group-action"},
		{"group content", render(t, ui.SidebarGroupContent(gsx.Raw("x"), nil)), "sidebar-group-content"},
		{"menu", render(t, ui.SidebarMenu(gsx.Raw("x"), nil)), "sidebar-menu"},
		{"menu item", render(t, ui.SidebarMenuItem(gsx.Raw("x"), nil)), "sidebar-menu-item"},
		{"menu badge", render(t, ui.SidebarMenuBadge(gsx.Raw("x"), nil)), "sidebar-menu-badge"},
		{"menu sub", render(t, ui.SidebarMenuSub(gsx.Raw("x"), nil)), "sidebar-menu-sub"},
		{"menu sub item", render(t, ui.SidebarMenuSubItem(gsx.Raw("x"), nil)), "sidebar-menu-sub-item"},
	}
	for _, part := range parts {
		t.Run(part.name, func(t *testing.T) {
			requireContainsAll(t, part.got, `data-gsxui-slot-`+part.slot)
			if strings.Contains(part.got, ` class=`) {
				t.Fatalf("library presentation class remains\nin: %s", part.got)
			}
			if strings.Contains(part.got, `data-slot=`) || strings.Contains(part.got, `data-sidebar=`) {
				t.Fatalf("legacy styling hook remains\nin: %s", part.got)
			}
		})
	}

	rail := render(t, ui.SidebarRail(nil))
	requireContainsAll(t, rail,
		`data-gsxui-sidebar-rail`,
		`aria-label="Toggle Sidebar"`,
		`tabindex="-1"`,
		`title="Toggle Sidebar"`,
	)
}

func TestSidebarMenuButtonReflectsEveryPresentationAxis(t *testing.T) {
	for _, tc := range []struct {
		variant string
		size    string
		active  bool
	}{
		{"", "", false},
		{"default", "sm", true},
		{"outline", "lg", false},
	} {
		got := render(t, ui.SidebarMenuButton(
			tc.active,
			tc.variant,
			tc.size,
			"",
			gsx.Raw("x"),
			gsx.Attrs{{Key: "class", Value: "caller"}},
		))
		variant := tc.variant
		if variant == "" {
			variant = "default"
		}
		size := tc.size
		if size == "" {
			size = "default"
		}
		requireContainsAll(t, got,
			`data-gsxui-slot-sidebar-menu-button`,
			`data-variant="`+variant+`"`,
			`data-size="`+size+`"`,
			`class="caller"`,
		)
		if tc.active {
			requireContainsAll(t, got, ` data-active`)
			if strings.Contains(got, `data-active=`) {
				t.Fatalf("active button must render data-active as a bare marker\nin: %s", got)
			}
		} else if strings.Contains(got, `data-active`) {
			t.Fatalf("inactive button must omit data-active\nin: %s", got)
		}
		if strings.Count(got, `class="caller"`) != 1 {
			t.Fatalf("caller class not forwarded exactly once\nin: %s", got)
		}
	}
}

func TestSidebarMenuControlsPreserveCallerStateAndDisabledAxes(t *testing.T) {
	button := render(t, ui.SidebarMenuButton(
		false,
		"",
		"",
		"",
		gsx.Raw("x"),
		gsx.Attrs{
			{Key: "data-state", Value: "open"},
			{Key: "disabled", Value: true},
			{Key: "aria-disabled", Value: "true"},
		},
	))
	requireContainsAll(t, button,
		`data-state="open"`,
		` disabled`,
		`aria-disabled="true"`,
	)

	action := render(t, ui.SidebarMenuAction(
		true,
		gsx.Raw("x"),
		gsx.Attrs{{Key: "data-state", Value: "open"}},
	))
	requireContainsAll(t, action, `data-state="open"`)

	subButton := render(t, ui.SidebarMenuSubButton(
		"",
		false,
		gsx.Raw("x"),
		gsx.Attrs{
			{Key: "disabled", Value: true},
			{Key: "aria-disabled", Value: "true"},
		},
	))
	requireContainsAll(t, subButton, ` disabled`, `aria-disabled="true"`)
}

func TestSidebarMenuButtonComposesTooltipTokens(t *testing.T) {
	got := render(t, ui.SidebarMenuButton(
		true,
		"outline",
		"lg",
		"Inbox",
		gsx.Raw("Inbox"),
		nil,
	))
	requireContainsAll(t, got,
		`data-gsxui-slot-sidebar-menu-button-tooltip data-gsxui-slot-tooltip`,
		`data-gsxui-slot-sidebar-menu-button`,
		`data-gsxui-tooltip-trigger`,
		`data-gsxui-slot-sidebar-menu-button-tooltip-content data-gsxui-slot-tooltip-content`,
		`data-gsxui-slot-tooltip-arrow`,
		`data-variant="outline"`,
		`data-size="lg"`,
		`data-active`,
		">Inbox<",
	)
	if strings.Contains(got, `data-active=`) {
		t.Fatalf("active tooltip button must render data-active as a bare marker\nin: %s", got)
	}
	if strings.Contains(got, ` class=`) || strings.Contains(got, `data-slot=`) || strings.Contains(got, `data-sidebar=`) {
		t.Fatalf("tooltip composition retained presentation classes or legacy styling hooks\nin: %s", got)
	}
}

func TestSidebarMenuActionReflectsShowOnHover(t *testing.T) {
	for _, show := range []bool{false, true} {
		got := render(t, ui.SidebarMenuAction(show, gsx.Raw("x"), nil))
		requireContainsAll(t, got, `data-gsxui-slot-sidebar-menu-action`)
		if show {
			requireContainsAll(t, got, ` data-show-on-hover`)
			if strings.Contains(got, `data-show-on-hover=`) {
				t.Fatalf("showOnHover=true must render data-show-on-hover as a bare marker\nin: %s", got)
			}
		} else if strings.Contains(got, `data-show-on-hover`) {
			t.Fatalf("showOnHover=false must omit data-show-on-hover\nin: %s", got)
		}
		if strings.Contains(got, ` class=`) {
			t.Fatalf("library presentation class remains\nin: %s", got)
		}
	}
}

func TestSidebarMenuSubButtonReflectsSizeActiveAndCallerClass(t *testing.T) {
	got := render(t, ui.SidebarMenuSubButton(
		"sm",
		true,
		gsx.Raw("x"),
		gsx.Attrs{{Key: "class", Value: "caller"}},
	))
	requireContainsAll(t, got,
		`data-gsxui-slot-sidebar-menu-sub-button`,
		`data-size="sm"`,
		`data-active`,
		`class="caller"`,
	)
	if strings.Contains(got, `data-active=`) {
		t.Fatalf("active sub-button must render data-active as a bare marker\nin: %s", got)
	}
}

func TestSidebarMenuSkeletonComposesSkeletonPartsAndKeepsDynamicWidth(t *testing.T) {
	got := render(t, ui.SidebarMenuSkeleton(true, nil))
	requireContainsAll(t, got,
		`data-gsxui-slot-sidebar-menu-skeleton`,
		`data-gsxui-slot-sidebar-menu-skeleton-icon data-gsxui-slot-skeleton`,
		`data-gsxui-slot-sidebar-menu-skeleton-text data-gsxui-slot-skeleton`,
		`--skeleton-width:`,
	)
	if strings.Contains(got, ` class=`) || strings.Contains(got, `data-sidebar=`) || strings.Contains(got, `data-slot=`) {
		t.Fatalf("skeleton composition retained presentation classes or legacy styling hooks\nin: %s", got)
	}

	without := render(t, ui.SidebarMenuSkeleton(false, nil))
	if strings.Contains(without, "sidebar-menu-skeleton-icon") {
		t.Fatalf("showIcon=false rendered icon part\nin: %s", without)
	}
	requireContainsAll(t, without, `data-gsxui-slot-sidebar-menu-skeleton-text data-gsxui-slot-skeleton`)
}
