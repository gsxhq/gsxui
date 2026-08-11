package port

import (
	"sort"
	"testing"

	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

func TestSectionComponent(t *testing.T) {
	tests := []struct {
		mark          string
		wantComponent string
		wantOK        bool
	}{
		{"Radio Group", "radio", true},
		{"Native Select", "native-select", true},
		{"Accordion", "accordion", true},
		{"Alert Dialog", "alert-dialog", true},
		{"Hover Card", "hover-card", true},
		{"Input OTP", "input-otp", true},
		{"React Aria", "", false},
		{"Chart", "", false},
		{"Menu Translucent", "", false},
		{"Bubble", "", false},
		{"Attachment", "", false},
		{"Marker", "", false},
		{"Questionnaire", "", false},
		{"Message Scroller", "", false},
		{"Message", "", false},
		{"", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.mark, func(t *testing.T) {
			gotComponent, gotOK := SectionComponent(tt.mark)
			if gotComponent != tt.wantComponent || gotOK != tt.wantOK {
				t.Errorf("SectionComponent(%q) = (%q, %v), want (%q, %v)",
					tt.mark, gotComponent, gotOK, tt.wantComponent, tt.wantOK)
			}
		})
	}
}

func TestIgnored(t *testing.T) {
	tests := []struct {
		class string
		want  bool
	}{
		{"cn-select-value-aria", true},
		{"cn-select-item-aria", true},
		{"cn-radio-group-item-aria", true},
		{"cn-select-value", false},
		{"cn-accordion-trigger", false},
	}

	for _, tt := range tests {
		t.Run(tt.class, func(t *testing.T) {
			if got := Ignored(tt.class); got != tt.want {
				t.Errorf("Ignored(%q) = %v, want %v", tt.class, got, tt.want)
			}
		})
	}
}

func TestStyleInvariant(t *testing.T) {
	tests := []struct {
		component string
		want      bool
	}{
		{"spinner", true},
		{"aspect-ratio", true},
		{"collapsible", true},
		{"toaster", true},
		{"button", false},
		{"accordion", false},
	}

	for _, tt := range tests {
		t.Run(tt.component, func(t *testing.T) {
			if got := StyleInvariant(tt.component); got != tt.want {
				t.Errorf("StyleInvariant(%q) = %v, want %v", tt.component, got, tt.want)
			}
		})
	}
}

func TestSlotFor(t *testing.T) {
	tests := []struct {
		name      string
		component string
		class     string
		wantSlot  string
		wantOK    bool
	}{
		{"root slot", "accordion", "cn-accordion", "", true},
		{"plain slot", "accordion", "cn-accordion-trigger", "trigger", true},
		{"multi-word slot", "accordion", "cn-accordion-content-inner", "content-inner", true},
		{"radio upstream prefix quirk", "radio", "cn-radio-group-item", "item", true},
		{"radio root", "radio", "cn-radio-group", "", true},
		{"select scroll-up ignored", "select", "cn-select-scroll-up-button", "", false},
		{"select scroll-down ignored", "select", "cn-select-scroll-down-button", "", false},
		{"dialog overlay fuses onto content", "dialog", "cn-dialog-overlay", "content", true},
		{"alert-dialog overlay fuses onto content", "alert-dialog", "cn-alert-dialog-overlay", "content", true},
		{"drawer overlay fuses onto content", "drawer", "cn-drawer-overlay", "content", true},
		{"sheet overlay fuses onto content", "sheet", "cn-sheet-overlay", "content", true},
		{"unrelated class", "accordion", "cn-select-trigger", "", false},
		{"unrelated prefix collision", "accordion", "cn-accordion2-trigger", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSlot, gotOK := SlotFor(tt.component, tt.class)
			if gotSlot != tt.wantSlot || gotOK != tt.wantOK {
				t.Errorf("SlotFor(%q, %q) = (%q, %v), want (%q, %v)",
					tt.component, tt.class, gotSlot, gotOK, tt.wantSlot, tt.wantOK)
			}
		})
	}
}

// upstreamMarkNames is the full, verified set of `/* MARK: <Name> */`
// section names in the real upstream style-maia.css (`grep -n "MARK:"
// apps/v4/registry/styles/style-maia.css`, 59 sections, matching dossier
// §5.2 exactly). It is hardcoded rather than read from the checkout so the
// coverage test below runs without a local shadcn-ui clone — it is testing
// this package's mapping TABLE, not upstream's file layout.
var upstreamMarkNames = []string{
	"Accordion", "Alert Dialog", "Alert", "Avatar", "Badge", "Breadcrumb",
	"Button", "Button Group", "Calendar", "Card", "Carousel", "Chart",
	"Checkbox", "Combobox", "Command", "Context Menu", "Dialog", "Drawer",
	"Dropdown Menu", "Empty", "Field", "Hover Card", "Input", "Input OTP",
	"Item", "Kbd", "Label", "Menubar", "Navigation Menu", "Native Select",
	"Pagination", "Popover", "Progress", "Radio Group", "Resizable",
	"Scroll Area", "Select", "Separator", "Sheet", "Sidebar", "Skeleton",
	"Slider", "Toast", "Switch", "Table", "Tabs", "Textarea", "Toggle",
	"Toggle Group", "Tooltip", "Input Group", "Menu Translucent", "Bubble",
	"Attachment", "Marker", "Questionnaire", "Message Scroller", "Message",
	"React Aria",
}

// TestMappingComponentCoverage is the important test: every one of the 54
// canonical component shapes must either be reachable from some upstream
// MARK section (via SectionComponent) or be declared StyleInvariant. A
// failure here means the mapping TABLE is wrong (a missing override, or a
// component that should be StyleInvariant but isn't listed) — fix the
// table, never this test.
func TestMappingComponentCoverage(t *testing.T) {
	reachable := make(map[string]bool)
	for _, mark := range upstreamMarkNames {
		if component, ok := SectionComponent(mark); ok {
			reachable[component] = true
		}
	}

	all := shapes.All()
	names := make([]string, 0, len(all))
	for name := range all {
		names = append(names, name)
	}
	sort.Strings(names)

	var missing []string
	for _, name := range names {
		if reachable[name] || StyleInvariant(name) {
			continue
		}
		missing = append(missing, name)
	}
	if len(missing) > 0 {
		t.Fatalf("components with no upstream MARK section and not StyleInvariant: %v", missing)
	}
}

// TestMappingComponentCoverageExactCounts pins the coverage down further
// than "nothing missing": the dossier's full 54-component survey (§6) found
// exactly 49 components identity-mapped from a MARK section, 1 declared
// override (radio ← "Radio Group"), and 4 StyleInvariant components with no
// section at all — 49+1+4 = 54. If this ever drifts, the shape of the drift
// (which bucket gained or lost) is more informative than a bare "missing"
// list.
func TestMappingComponentCoverageExactCounts(t *testing.T) {
	reachableByOverride := map[string]bool{}
	reachableByIdentity := map[string]bool{}
	for _, mark := range upstreamMarkNames {
		component, ok := SectionComponent(mark)
		if !ok {
			continue
		}
		if _, overridden := sectionComponentOverrides[mark]; overridden {
			reachableByOverride[component] = true
			continue
		}
		reachableByIdentity[component] = true
	}

	all := shapes.All()
	var identity, override, invariant, unmatched []string
	for name := range all {
		switch {
		case reachableByIdentity[name]:
			identity = append(identity, name)
		case reachableByOverride[name]:
			override = append(override, name)
		case StyleInvariant(name):
			invariant = append(invariant, name)
		default:
			unmatched = append(unmatched, name)
		}
	}

	if len(unmatched) > 0 {
		t.Fatalf("unmatched components (neither mapped nor StyleInvariant): %v", unmatched)
	}
	if len(identity) != 49 {
		t.Errorf("identity-mapped components = %d, want 49: %v", len(identity), identity)
	}
	if len(override) != 1 {
		t.Errorf("override-mapped components = %d, want 1: %v", len(override), override)
	}
	if len(invariant) != 4 {
		t.Errorf("StyleInvariant components = %d, want 4: %v", len(invariant), invariant)
	}
}
