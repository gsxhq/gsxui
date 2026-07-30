package examples_test

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/internal/registry"
	"github.com/gsxhq/gsxui/internal/stylecontract"
	"github.com/gsxhq/gsxui/site/examples"
	"github.com/gsxhq/gsxui/ui"
	"golang.org/x/net/html"
)

type runtimeContractEntry struct {
	Component string `json:"component"`
	Slot      string `json:"slot"`
	Attribute string `json:"attribute"`
	Value     string `json:"value"`
	Scenario  string `json:"scenario"`
}

var runtimeScenarioByComponent = map[string]string{
	"alert-dialog":    "alert-dialog-lifecycle",
	"combobox":        "combobox-lifecycle",
	"command":         "command-runtime",
	"context-menu":    "context-menu-lifecycle",
	"dialog":          "dialog-lifecycle",
	"drawer":          "drawer-lifecycle",
	"dropdown-menu":   "dropdown-lifecycle",
	"hover-card":      "hover-card-lifecycle",
	"input-otp":       "input-otp-caret",
	"menubar":         "menubar-lifecycle",
	"navigation-menu": "navigation-menu-lifecycle",
	"popover":         "popover-lifecycle",
	"select":          "select-lifecycle",
	"sheet":           "sheet-lifecycle",
	"toast":           "sonner-lifecycle",
	"toaster":         "sonner-lifecycle",
	"tooltip":         "tooltip-lifecycle",
}

func TestRuntimeStyleContractManifestMatchesTypedContract(t *testing.T) {
	contracts := stylecontract.All()
	want := make([]runtimeContractEntry, 0)
	usedScenarios := make(map[string]struct{})
	for _, component := range contracts {
		scenario := runtimeScenarioByComponent[component.RegistryName]
		for _, slot := range component.Slots {
			if slot.Runtime {
				if scenario == "" {
					t.Fatalf("runtime slot %q on %q has no explicit browser scenario", slot.Name, component.RegistryName)
				}
				usedScenarios[scenario] = struct{}{}
				want = append(want, runtimeContractEntry{
					Component: component.RegistryName,
					Slot:      slot.Name,
					Attribute: stylecontract.SlotAttribute(slot.Name),
					Value:     "",
					Scenario:  scenario,
				})
			}
			for _, axis := range slot.Axes {
				for _, value := range axis.RuntimeValues {
					if scenario == "" {
						t.Fatalf("runtime value %s=%q on %q/%q has no explicit browser scenario", axis.Attribute, value, component.RegistryName, slot.Name)
					}
					usedScenarios[scenario] = struct{}{}
					want = append(want, runtimeContractEntry{
						Component: component.RegistryName,
						Slot:      slot.Name,
						Attribute: axis.Attribute,
						Value:     value,
						Scenario:  scenario,
					})
				}
			}
		}
	}
	for component, scenario := range runtimeScenarioByComponent {
		if _, ok := usedScenarios[scenario]; !ok {
			t.Errorf("runtime scenario %q for component %q covers no typed runtime state", scenario, component)
		}
	}
	sortRuntimeEntries(want)

	path := filepath.Join("..", "..", "jstest", "runtime-style-contract.json")
	if os.Getenv("UPDATE_RUNTIME_STYLE_CONTRACT") == "1" {
		data, err := json.MarshalIndent(want, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		data = append(data, '\n')
		if err := os.WriteFile(path, data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var got []runtimeContractEntry
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
	if !slices.IsSortedFunc(got, compareRuntimeEntry) {
		t.Fatalf("%s must be sorted by component, slot, attribute, value, scenario", path)
	}
	if !slices.Equal(got, want) {
		t.Fatalf("runtime style contract manifest drift\n got: %#v\nwant: %#v", got, want)
	}
}

func sortRuntimeEntries(entries []runtimeContractEntry) {
	slices.SortFunc(entries, compareRuntimeEntry)
}

func compareRuntimeEntry(a, b runtimeContractEntry) int {
	for _, comparison := range []int{
		strings.Compare(a.Component, b.Component),
		strings.Compare(a.Slot, b.Slot),
		strings.Compare(a.Attribute, b.Attribute),
		strings.Compare(a.Value, b.Value),
		strings.Compare(a.Scenario, b.Scenario),
	} {
		if comparison != 0 {
			return comparison
		}
	}
	return 0
}

func TestRegisteredExamplesCoverStyleContract(t *testing.T) {
	contracts := stylecontract.All()
	if err := stylecontract.Validate(contracts); err != nil {
		t.Fatalf("stylecontract.Validate(stylecontract.All()): %v", err)
	}

	contractComponents := make([]string, 0, len(contracts))
	declaredSlots := make(map[string]stylecontract.Slot)
	axisAttributes := make(map[string]struct{})
	for _, component := range contracts {
		if component.RegistryName == "" {
			t.Errorf("contract component %q has no exact registry identity", component.Name)
			continue
		}
		contractComponents = append(contractComponents, component.RegistryName)
		for _, slot := range component.Slots {
			declaredSlots[slot.Name] = slot
			for _, axis := range slot.Axes {
				axisAttributes[axis.Attribute] = struct{}{}
			}
		}
	}
	slices.Sort(contractComponents)

	registered, err := registry.Components()
	if err != nil {
		t.Fatal(err)
	}
	// ui/sonner.gsx is one vendorable file backing two typed style-contract
	// components, "toaster" and "toast" (see internal/stylecontract/contracts_toast.go);
	// registry.Components() stays file-derived, so fold the pair back to
	// "sonner" before comparing against it.
	foldedContractComponents := make([]string, 0, len(contractComponents))
	for _, name := range contractComponents {
		if name == "toaster" || name == "toast" {
			continue
		}
		foldedContractComponents = append(foldedContractComponents, name)
	}
	foldedContractComponents = append(foldedContractComponents, "sonner")
	slices.Sort(foldedContractComponents)
	if !slices.Equal(foldedContractComponents, registered) {
		t.Errorf("typed contract components (toast/toaster folded to sonner) = %v; registry components = %v", foldedContractComponents, registered)
	}

	emittedSlots := make(map[string]struct{})
	emittedValues := make(map[string]map[string]map[string]struct{})
	slotAttributePrefix := stylecontract.SlotAttribute("")
	recordNode := func(source string, node gsx.Node) {
		var rendered bytes.Buffer
		if err := node.Render(context.Background(), &rendered); err != nil {
			t.Errorf("%s render: %v", source, err)
			return
		}
		document, err := html.Parse(strings.NewReader(rendered.String()))
		if err != nil {
			t.Errorf("%s parse: %v", source, err)
			return
		}
		walkElements(document, func(node *html.Node) {
			slots := make([]string, 0)
			for _, attr := range node.Attr {
				slot, valid := stylecontract.SlotName(attr.Key)
				if !valid {
					if strings.HasPrefix(attr.Key, slotAttributePrefix) && slot == "" {
						t.Errorf("%s emits empty slot marker %q", source, attr.Key)
					} else if strings.HasPrefix(attr.Key, slotAttributePrefix) {
						t.Errorf("%s emits malformed slot marker %q", source, attr.Key)
					}
					continue
				}
				if attr.Val != "" {
					t.Errorf("%s emits valued slot marker %s=%q", source, attr.Key, attr.Val)
				}
				slots = append(slots, slot)
				emittedSlots[slot] = struct{}{}
				if _, ok := declaredSlots[slot]; !ok {
					t.Errorf("%s emits undeclared slot %q", source, slot)
				}
			}
			if len(slots) == 0 {
				return
			}

			for _, attr := range node.Attr {
				if _, isAxis := axisAttributes[attr.Key]; !isAxis {
					continue
				}
				if rejectsUnownedCustomAxis(slots, attr.Key, axisAttributes, declaredSlots) {
					t.Errorf(
						"%s emits globally known custom axis %s=%q on slots %q, but none owns the attribute",
						source, attr.Key, attr.Val, slots,
					)
				}
				for _, name := range slots {
					slot, declared := declaredSlots[name]
					if !declared {
						continue
					}
					for _, axis := range slot.Axes {
						if axis.Attribute != attr.Key {
							continue
						}
						if axis.Values != nil && !slices.Contains(axis.Values, attr.Val) {
							t.Errorf(
								"%s slot %q emits %s=%q; declared values are %v",
								source, name, attr.Key, attr.Val, axis.Values,
							)
						}
						if emittedValues[name] == nil {
							emittedValues[name] = make(map[string]map[string]struct{})
						}
						if emittedValues[name][attr.Key] == nil {
							emittedValues[name][attr.Key] = make(map[string]struct{})
						}
						emittedValues[name][attr.Key][attr.Val] = struct{}{}
					}
				}
			}
		})
	}

	// Toaster is intentionally mounted once by the real site and harness
	// shells, not repeated inside an example body.
	recordNode("application shell Toaster", ui.Toaster(nil))

	for _, component := range registered {
		registeredExamples := examples.For(component)
		if len(registeredExamples) == 0 {
			t.Errorf("registry component %q has no real registered example", component)
			continue
		}
		for _, example := range registeredExamples {
			recordNode(component+"/"+example.Name, example.Node)
			for _, preview := range example.Previews {
				recordNode(component+"/"+example.Name+"/preview/"+preview.Name, preview.Node)
			}
		}
	}

	for name, slot := range declaredSlots {
		if _, ok := emittedSlots[name]; !ok && !slot.Runtime {
			t.Errorf("declared slot %q is not emitted by any real registered example", name)
		}
		for _, axis := range slot.Axes {
			for _, value := range axis.Values {
				if slices.Contains(axis.RuntimeValues, value) {
					continue
				}
				if _, ok := emittedValues[name][axis.Attribute][value]; !ok {
					t.Errorf(
						"declared axis value %s=%q on slot %q is not emitted by any real registered example",
						axis.Attribute, value, name,
					)
				}
			}
		}
	}
}

// rejectsUnownedCustomAxis enforces the library-owned custom state namespace
// globally. Native, ARIA, and role attributes also carry intrinsic semantics,
// so they are value-checked only when the emitting slot declares them.
func rejectsUnownedCustomAxis(
	tokens []string,
	attribute string,
	axisAttributes map[string]struct{},
	declaredSlots map[string]stylecontract.Slot,
) bool {
	if _, known := axisAttributes[attribute]; !known || !strings.HasPrefix(attribute, "data-") {
		return false
	}
	for _, token := range tokens {
		for _, axis := range declaredSlots[token].Axes {
			if axis.Attribute == attribute {
				return false
			}
		}
	}
	return true
}

func TestRejectsUnownedCustomAxisGlobally(t *testing.T) {
	knownAxes := map[string]struct{}{
		"data-size":     {},
		"data-state":    {},
		"aria-expanded": {},
	}
	slots := map[string]stylecontract.Slot{
		"first": {
			Name: "first",
			Axes: []stylecontract.Axis{{Attribute: "data-size", Values: []string{"sm"}}},
		},
		"second": {
			Name: "second",
			Axes: []stylecontract.Axis{{Attribute: "data-state", Values: []string{"closed", "open"}}},
		},
	}

	if !rejectsUnownedCustomAxis([]string{"first"}, "data-state", knownAxes, slots) {
		t.Fatal("globally known data-state on an unrelated contracted slot was accepted")
	}
	if rejectsUnownedCustomAxis([]string{"second"}, "data-state", knownAxes, slots) {
		t.Fatal("slot-owned data-state was rejected")
	}
	if rejectsUnownedCustomAxis([]string{"first"}, "data-unknown", knownAxes, slots) {
		t.Fatal("unknown custom attribute was treated as a contract axis")
	}
	if rejectsUnownedCustomAxis([]string{"first"}, "aria-expanded", knownAxes, slots) {
		t.Fatal("intrinsic ARIA semantics were treated as globally owned custom state")
	}
}

func walkElements(node *html.Node, visit func(*html.Node)) {
	if node.Type == html.ElementNode {
		visit(node)
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		walkElements(child, visit)
	}
}
