package stylepreview_test

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
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
