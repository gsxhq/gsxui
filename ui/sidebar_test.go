package ui_test

import (
	"html"
	"regexp"
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/merge"
	"github.com/gsxhq/gsxui/ui"
)

// sidebarRecipeUtilities reads one Sidebar slot's compiled utilities out of
// the default style. ui/sidebar.gsx is generated output now, so these are
// exactly the classes it renders. defaultStyleRecipe's loader
// (input-group_test.go) already caches every component's recipe, and Sidebar
// composes four other migrated components, so this reuses it rather than
// adding a fifth per-component loader.
func sidebarRecipeUtilities(slot string) []string {
	return styleRecipeUtilities("sidebar", "gsxui-recipe-"+slot)
}

// canonicalSidebarClass is the class attribute one Sidebar slot renders on its
// own, plus any caller classes, merged the way gsx merges class values at
// runtime.
func canonicalSidebarClass(slot string, caller ...string) string {
	classes := append([]string(nil), sidebarRecipeUtilities(slot)...)
	classes = append(classes, caller...)
	return `class="` + html.EscapeString(merge.Merge(classes)) + `"`
}

var classAttributePattern = regexp.MustCompile(`(?s) class="[^"]*"`)

// withoutClassAttributes strips every class attribute from a render. Sidebar's
// recipe utilities contain arbitrary variants written against its own
// data-attributes (`[&[data-active]]:…`, `md:[&[data-show-on-hover]]:…`), so a
// bare strings.Contains for an ATTRIBUTE now matches inside the class value
// too. Every assertion about attribute presence reads the stripped render.
func withoutClassAttributes(render string) string {
	return classAttributePattern.ReplaceAllString(render, "")
}

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
		`data-gsxui-slot-sidebar-wrapper`,
		`style="--sidebar-width:16rem;--sidebar-width-icon:3rem"`,
		canonicalSidebarClass("sidebar-wrapper", "caller"),
	)
	if got, want := strings.Count(open, canonicalSidebarClass("sidebar-wrapper", "caller")), 1; got != want {
		t.Fatalf("wrapper class count = %d, want %d\nin: %s", got, want, open)
	}

	closed := render(t, ui.SidebarProvider(false, gsx.Raw("x"), nil))
	requireContainsAll(t, closed, `data-state="collapsed"`, `data-gsxui-slot-sidebar-wrapper`)
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
		`data-gsxui-slot-sidebar-desktop`,
		`data-gsxui-slot-sidebar-gap`,
		`data-gsxui-slot-sidebar-container`,
		`data-gsxui-slot-sidebar-inner`,
	)
	// The mobile tree's root is <ui.Sheet>, which composes <ui.Dialog>
	// directly — Dialog is migrated to the slot axis, so its own "contents"
	// recipe class merges with the caller's "caller" there instead of caller
	// rendering alone. The desktop tree's container is Sidebar's own
	// container slot, which since Sidebar's migration carries that slot's
	// recipe utilities ahead of the caller class.
	if count := strings.Count(got, `class="contents caller"`); count != 1 {
		t.Fatalf(`class="contents caller" count = %d, want exactly 1 (mobile Sheet/Dialog root)\nin: %s`, count, got)
	}
	container := canonicalSidebarClass("sidebar-container", "caller")
	if count := strings.Count(got, container); count != 1 {
		t.Fatalf("container class count = %d, want exactly 1 for %s\nin: %s", count, container, got)
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
		canonicalSidebarClass("sidebar", "caller"),
	)
	if strings.Contains(got, `sidebar-mobile-root`) || strings.Contains(got, `sidebar-desktop`) {
		t.Fatalf("collapsible=none rendered a responsive state-machine tree\nin: %s", got)
	}
}

func TestSidebarPrimitiveCompositionAndCallerClassPlacement(t *testing.T) {
	cases := []struct {
		name string
		got  string
		want string
		// wantClass is the class attribute the COMPOSED element renders.
		wantClass string
		// classAttributes is how many class attributes the whole render
		// carries. The trigger's is 3 because SidebarTrigger also renders its
		// own sr-only label span (a Sidebar slot of its own) and its
		// PanelLeft icon carries a literal rtl:rotate-180 class so the
		// trigger glyph flips in RTL contexts.
		classAttributes int
	}{
		{
			// Sidebar's trigger slot rides into <ui.Button>'s class
			// attribute, where merge.Merge drops Button's size-8 icon arm in
			// favour of size-7 — the promotion in default.css this retired.
			name:            "trigger button",
			got:             render(t, ui.SidebarTrigger(gsx.Attrs{{Key: "class", Value: "caller"}})),
			want:            `data-gsxui-slot-sidebar-trigger data-gsxui-slot-button`,
			wantClass:       canonicalButtonClass("ghost", "icon", append(sidebarRecipeUtilities("sidebar-trigger"), "caller")...),
			classAttributes: 3,
		},
		{
			name: "input",
			got:  render(t, ui.SidebarInput(gsx.Attrs{{Key: "class", Value: "caller"}})),
			want: `data-gsxui-slot-sidebar-input data-gsxui-slot-input`,
			wantClass: canonicalComposedClass(
				styleRecipeUtilities("input", "gsxui-recipe-input"),
				sidebarRecipeUtilities("sidebar-input"),
				[]string{"caller"},
			),
			classAttributes: 1,
		},
		{
			name:            "separator",
			got:             render(t, ui.SidebarSeparator(gsx.Attrs{{Key: "class", Value: "caller"}})),
			want:            `data-gsxui-slot-sidebar-separator data-gsxui-slot-separator`,
			wantClass:       canonicalSeparatorClass("horizontal", append(sidebarRecipeUtilities("sidebar-separator"), "caller")...),
			classAttributes: 1,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			requireContainsAll(t, tc.got, tc.want, tc.wantClass)
			if count := strings.Count(tc.got, tc.wantClass); count != 1 {
				t.Fatalf("exact class count = %d, want 1 for %s\nin: %s", count, tc.wantClass, tc.got)
			}
			if count := strings.Count(tc.got, `class=`); count != tc.classAttributes {
				t.Fatalf("class attribute count = %d, want %d\nin: %s", count, tc.classAttributes, tc.got)
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
		`data-gsxui-slot-sidebar-trigger`,
		`data-gsxui-slot-icon`,
		`data-gsxui-slot-sidebar-trigger-label`,
	)
}

func TestSidebarTriggerComposesPresenceMarkersOnButton(t *testing.T) {
	got := render(t, ui.SidebarTrigger(nil))
	requirePresenceAttributesOnSameTag(t, got, "data-gsxui-slot-sidebar-trigger",
		"data-gsxui-slot-button",
	)
}

// TestSidebarPlainPartsExposeNamespacedSlotsWithTheirRecipeClass replaces the
// pre-migration "…WithoutPresentationClasses". Sidebar is on the slot axis
// now, so each of these parts legitimately renders exactly its own slot's
// compiled utilities and nothing else. The legacy-hook bans are unchanged.
func TestSidebarPlainPartsExposeNamespacedSlotsWithTheirRecipeClass(t *testing.T) {
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
			requireContainsAll(t, part.got, `data-gsxui-slot-`+part.slot, canonicalSidebarClass(part.slot))
			if count := strings.Count(part.got, ` class=`); count != 1 {
				t.Fatalf("class attribute count = %d, want exactly 1\nin: %s", count, part.got)
			}
			if strings.Contains(part.got, `data-slot=`) || strings.Contains(part.got, `data-sidebar=`) {
				t.Fatalf("legacy styling hook remains\nin: %s", part.got)
			}
		})
	}

	rail := render(t, ui.SidebarRail(nil))
	requireContainsAll(t, rail,
		`data-gsxui-slot-sidebar-rail`,
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
			canonicalSidebarClass("sidebar-menu-button", "caller"),
		)
		attributes := withoutClassAttributes(got)
		if tc.active {
			requireContainsAll(t, attributes, ` data-active`)
			if strings.Contains(attributes, `data-active=`) {
				t.Fatalf("active button must render data-active as a bare marker\nin: %s", got)
			}
		} else if strings.Contains(attributes, `data-active`) {
			t.Fatalf("inactive button must omit data-active\nin: %s", got)
		}
		if strings.Count(got, canonicalSidebarClass("sidebar-menu-button", "caller")) != 1 {
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
	if strings.Contains(withoutClassAttributes(got), `data-active=`) {
		t.Fatalf("active tooltip button must render data-active as a bare marker\nin: %s", got)
	}
	// Tooltip's root/content/arrow carry their own recipe classes, and since
	// Sidebar's own migration the <button> carries the menu-button slot's.
	// Four class attributes, one per styled element — the tooltip wrapper
	// itself still contributes none of Sidebar's presentation.
	if strings.Count(got, `class=`) != 4 {
		t.Fatalf("expected exactly 4 class= attributes (menu button + tooltip root, content, arrow)\nin: %s", got)
	}
	requireContainsAll(t, got, canonicalSidebarClass("sidebar-menu-button"))
	if strings.Contains(got, `data-slot=`) || strings.Contains(got, `data-sidebar=`) {
		t.Fatalf("tooltip composition retained legacy styling hooks\nin: %s", got)
	}
}

func TestSidebarMenuActionReflectsShowOnHover(t *testing.T) {
	for _, show := range []bool{false, true} {
		got := render(t, ui.SidebarMenuAction(show, gsx.Raw("x"), nil))
		requireContainsAll(t, got, `data-gsxui-slot-sidebar-menu-action`, canonicalSidebarClass("sidebar-menu-action"))
		attributes := withoutClassAttributes(got)
		if show {
			requireContainsAll(t, attributes, ` data-show-on-hover`)
			if strings.Contains(attributes, `data-show-on-hover=`) {
				t.Fatalf("showOnHover=true must render data-show-on-hover as a bare marker\nin: %s", got)
			}
		} else if strings.Contains(attributes, `data-show-on-hover`) {
			t.Fatalf("showOnHover=false must omit data-show-on-hover\nin: %s", got)
		}
		if count := strings.Count(got, ` class=`); count != 1 {
			t.Fatalf("class attribute count = %d, want exactly 1\nin: %s", count, got)
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
		canonicalSidebarClass("sidebar-menu-sub-button", "caller"),
	)
	if strings.Contains(withoutClassAttributes(got), `data-active=`) {
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
	// SidebarMenuSkeleton composes ui.Skeleton, which is migrated too, so the
	// two child elements carry Skeleton's utilities merged with the two
	// skeleton slots' own. Assert the WRAPPER carries exactly its own slot,
	// and keep the legacy-hook bans absolute.
	wrapper := got[:strings.Index(got, ">")+1]
	if !strings.Contains(wrapper, canonicalSidebarClass("sidebar-menu-skeleton")) {
		t.Fatalf("sidebar-menu-skeleton wrapper missing its recipe class\nin: %s", got)
	}
	requireContainsAll(t, got, canonicalComposedClass(
		styleRecipeUtilities("skeleton", "gsxui-recipe-skeleton"),
		sidebarRecipeUtilities("sidebar-menu-skeleton-icon"),
	))
	if strings.Contains(got, `data-sidebar=`) || strings.Contains(got, `data-slot=`) {
		t.Fatalf("skeleton composition retained legacy styling hooks\nin: %s", got)
	}

	without := render(t, ui.SidebarMenuSkeleton(false, nil))
	if strings.Contains(without, "sidebar-menu-skeleton-icon") {
		t.Fatalf("showIcon=false rendered icon part\nin: %s", without)
	}
	requireContainsAll(t, without, `data-gsxui-slot-sidebar-menu-skeleton-text data-gsxui-slot-skeleton`)
}

// canonicalComposedClass merges several utility lists in render order into the
// class attribute the composed element renders. Sidebar composes Button,
// Input, Separator and Skeleton, and in every case its own slot utilities
// arrive AFTER the composed component's, which is what lets merge.Merge settle
// the properties the two contest.
func canonicalComposedClass(lists ...[]string) string {
	var classes []string
	for _, list := range lists {
		classes = append(classes, list...)
	}
	return `class="` + html.EscapeString(merge.Merge(classes)) + `"`
}
