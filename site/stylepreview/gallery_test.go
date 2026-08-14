package stylepreview_test

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
)

// galleryTags maps every component key in registry/generated/recipes.json to
// the representative tag the gallery must render for it. A component whose
// key is missing here fails the test: adding a component to the catalogue
// requires either adding it to the gallery (and naming its tag here) or
// deliberately excluding it in galleryExcluded with a reason.
var galleryTags = map[string]string{
	"accordion":       "Accordion",
	"alert":           "Alert",
	"alert-dialog":    "AlertDialogContent",
	"aspect-ratio":    "AspectRatio",
	"avatar":          "Avatar",
	"badge":           "Badge",
	"breadcrumb":      "Breadcrumb",
	"button":          "Button",
	"button-group":    "ButtonGroup",
	"calendar":        "Calendar",
	"card":            "Card",
	"carousel":        "Carousel",
	"checkbox":        "Checkbox",
	"collapsible":     "Collapsible",
	"combobox":        "Combobox",
	"command":         "Command",
	"context-menu":    "ContextMenu",
	"dialog":          "Dialog",
	"drawer":          "Drawer",
	"dropdown-menu":   "DropdownMenu",
	"empty":           "Empty",
	"field":           "Field",
	"hover-card":      "HoverCard",
	"input":           "Input",
	"input-group":     "InputGroup",
	"input-otp":       "InputOTP",
	"item":            "Item",
	"kbd":             "Kbd",
	"label":           "Label",
	"menubar":         "Menubar",
	"native-select":   "NativeSelect",
	"navigation-menu": "NavigationMenu",
	"pagination":      "Pagination",
	"popover":         "Popover",
	"progress":        "Progress",
	"radio":           "Radio",
	"resizable":       "ResizablePanelGroup",
	"scroll-area":     "ScrollArea",
	"select":          "Select",
	"separator":       "Separator",
	"sheet":           "Sheet",
	"sidebar":         "Sidebar",
	"skeleton":        "Skeleton",
	"slider":          "Slider",
	"spinner":         "Spinner",
	"switch":          "Switch",
	"table":           "Table",
	"tabs":            "Tabs",
	"textarea":        "Textarea",
	"toast":           "Toast",
	"toggle":          "Toggle",
	"toggle-group":    "ToggleGroup",
	"tooltip":         "Tooltip",
}

// galleryExcluded names components deliberately absent from the gallery's
// visual layout, each with a reason. An exclusion here must be a conscious
// decision, never a silent skip.
var galleryExcluded = map[string]string{
	// Toaster is the fixed-position toast REGION: invisible chrome at rest,
	// so it adds nothing a gallery can show — the Toast card carries the
	// themed chrome — and it stamps the fixed id="gsxui-toaster" (its
	// caller-visible HTMX oob contract, deliberately unprefixable), which
	// would appear twice in a document that renders the gallery per style.
	"toaster": "singleton id gsxui-toaster; the preview document renders the gallery twice, and the Toast card already shows the themed chrome",
	// Chart ships its container in this task (data-chart, --color-<key>
	// variables, ChartConfig/ChartSeries), not a chart CARD: there is no
	// data model builder, no client renderer, and no ChartTooltip/
	// ChartLegend yet (those land in Tasks 3-7), so nothing visual exists
	// for the gallery to show. Task 9 (docs/superpowers/plans/
	// 2026-08-14-chart-component.md) adds the real gallery chart card once
	// BarChart/AreaChart exist.
	"chart": "no chart model/renderer yet — Task 9 adds the real gallery chart card once BarChart/AreaChart exist",
}

// TestGalleryCoversEveryComponent asserts that every component the recipe
// contract declares appears in the authored gallery composition at least
// once, so a future component cannot silently miss the theme preview.
func TestGalleryCoversEveryComponent(t *testing.T) {
	contractSource, err := os.ReadFile("../../registry/generated/recipes.json")
	if err != nil {
		t.Fatal(err)
	}
	var contract struct {
		Components map[string]json.RawMessage `json:"components"`
	}
	if err := json.Unmarshal(contractSource, &contract); err != nil {
		t.Fatal(err)
	}
	if len(contract.Components) == 0 {
		t.Fatal("recipes.json declares no components")
	}

	gallery, err := os.ReadFile("gallery.gsx.src")
	if err != nil {
		t.Fatal(err)
	}

	for component := range contract.Components {
		if reason, ok := galleryExcluded[component]; ok {
			if reason == "" {
				t.Errorf("component %q is excluded without a reason", component)
			}
			if _, both := galleryTags[component]; both {
				t.Errorf("component %q is both excluded and mapped to a gallery tag", component)
			}
			continue
		}
		tag, ok := galleryTags[component]
		if !ok {
			t.Errorf("component %q has no gallery tag mapping; add it to the gallery or exclude it with a reason", component)
			continue
		}
		pattern := regexp.MustCompile(fmt.Sprintf(`<%s[\s/>]`, regexp.QuoteMeta(tag)))
		if !pattern.Match(gallery) {
			t.Errorf("gallery.gsx.src renders no <%s> for component %q", tag, component)
		}
	}

	for component := range galleryTags {
		if _, ok := contract.Components[component]; !ok {
			t.Errorf("galleryTags names %q, which recipes.json does not declare", component)
		}
	}
	for component := range galleryExcluded {
		if _, ok := contract.Components[component]; !ok {
			t.Errorf("galleryExcluded names %q, which recipes.json does not declare", component)
		}
	}
}

// galleryPageCards lists the card component identifiers each preview page
// (theme-creator-parity Task 4c) must render, in DOM order. It backs
// TestGalleryHasTwoPagesWithExpectedCards below and is mirrored by
// jstest/specs/style-visual.spec.ts's own per-page card list for the
// visual snapshot gate — keep the two in sync.
var galleryPageCards = map[string][]string{
	"components": {
		"galleryButtonsCard", "galleryLoginCard", "gallerySettingsCard",
		"galleryCalendarCard", "galleryMenusCard", "galleryFeedbackCard",
		"galleryTeamCard", "galleryTableCard", "galleryChartCard",
		"galleryTabsCard", "galleryNavigationCard", "galleryControlsCard",
		"galleryOverlaysCard", "galleryEmptyCard", "galleryMediaCard",
		"gallerySidebarCard",
	},
	"product": {
		"galleryPricingCard", "galleryStatsCard", "galleryChatCard",
		"galleryRolesCard", "galleryBookingCard", "galleryActivityCard",
		"galleryUploadsCard", "galleryNotificationsCard", "gallerySearchCard",
		"galleryProfileCard", "galleryOnboardingCard", "galleryProjectsEmptyCard",
		"galleryVerifyCard", "galleryNotificationPrefsCard", "galleryBillingCard",
	},
}

// TestGalleryHasTwoPagesWithExpectedCards asserts the theme-creator-parity
// Task 4c page split stays structurally intact: both page wrapper divs
// exist, "product" starts hidden (matching web/theme-state.js's
// createThemeState default of state.page = "components"), and every card
// galleryPageCards names for a page is actually referenced inside that
// page's own component body — not merely present somewhere in the file,
// which TestGalleryCoversEveryComponent above already tolerates (a card
// moved to the wrong page, or dropped from both, would still pass that
// test as long as its tag appears once anywhere).
func TestGalleryHasTwoPagesWithExpectedCards(t *testing.T) {
	gallery, err := os.ReadFile("gallery.gsx.src")
	if err != nil {
		t.Fatal(err)
	}
	source := string(gallery)

	if !strings.Contains(source, `data-theme-preview-page="components"`) {
		t.Error(`gallery.gsx.src has no data-theme-preview-page="components" wrapper`)
	}
	if !strings.Contains(source, `data-theme-preview-page="product"`) {
		t.Error(`gallery.gsx.src has no data-theme-preview-page="product" wrapper`)
	}
	// The product page must start hidden: theme.gsx's page tab pair
	// defaults to "components" for its own JS-less first paint
	// (themePagePressedAttr), and web/theme-state.js's createThemeState
	// defaults state.page the same way. A markup/state default mismatch
	// here would only ever surface as a visible flash before JS runs.
	productWrapper := regexp.MustCompile(`data-theme-preview-page="product"[^>]*>`)
	match := productWrapper.FindString(source)
	if match == "" {
		t.Fatal("could not find the product page's opening tag")
	}
	if !strings.Contains(match, "hidden") {
		t.Errorf("product page wrapper is not hidden by default: %q", match)
	}

	pageBody := func(component string) string {
		t.Helper()
		start := strings.Index(source, "component "+component+"(idp string) {")
		if start == -1 {
			t.Fatalf("gallery.gsx.src has no component %s(idp string)", component)
		}
		// Both page components are flat — a wrapping div plus self-closed
		// card tags, no nested control flow — so the first "\n}\n" after
		// the opening line closes the component body.
		end := strings.Index(source[start:], "\n}\n")
		if end == -1 {
			t.Fatalf("could not find the end of component %s", component)
		}
		return source[start : start+end]
	}

	pageComponent := map[string]string{
		"components": "galleryComponentsPage",
		"product":    "galleryProductPage",
	}
	for page, cards := range galleryPageCards {
		body := pageBody(pageComponent[page])
		for _, card := range cards {
			if !strings.Contains(body, "<"+card) {
				t.Errorf("%s does not render <%s> — expected on the %q page", pageComponent[page], card, page)
			}
		}
	}
}
