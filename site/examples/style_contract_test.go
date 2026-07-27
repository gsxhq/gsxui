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
	"dropdown":        "dropdown-lifecycle",
	"hover-card":      "hover-card-lifecycle",
	"input-otp":       "input-otp-caret",
	"menubar":         "menubar-lifecycle",
	"navigation-menu": "navigation-menu-lifecycle",
	"popover":         "popover-lifecycle",
	"select":          "select-lifecycle",
	"sheet":           "sheet-lifecycle",
	"sonner":          "sonner-lifecycle",
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
					Attribute: "data-gsxui-slot",
					Value:     slot.Name,
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
	if !slices.Equal(contractComponents, registered) {
		t.Errorf("typed contract components = %v; registry components = %v", contractComponents, registered)
	}

	emittedSlots := make(map[string]struct{})
	emittedValues := make(map[string]map[string]map[string]struct{})
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
			slotText, ok := attribute(node, "data-gsxui-slot")
			if !ok {
				return
			}
			tokens := strings.Fields(slotText)
			for _, token := range tokens {
				emittedSlots[token] = struct{}{}
				if _, ok := declaredSlots[token]; !ok {
					t.Errorf("%s emits undeclared slot %q", source, token)
				}
			}

			for _, attr := range node.Attr {
				if _, isAxis := axisAttributes[attr.Key]; !isAxis {
					continue
				}
				for _, token := range tokens {
					slot, declared := declaredSlots[token]
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
								source, token, attr.Key, attr.Val, axis.Values,
							)
						}
						if emittedValues[token] == nil {
							emittedValues[token] = make(map[string]map[string]struct{})
						}
						if emittedValues[token][attr.Key] == nil {
							emittedValues[token][attr.Key] = make(map[string]struct{})
						}
						emittedValues[token][attr.Key][attr.Val] = struct{}{}
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

func walkElements(node *html.Node, visit func(*html.Node)) {
	if node.Type == html.ElementNode {
		visit(node)
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		walkElements(child, visit)
	}
}

func attribute(node *html.Node, name string) (string, bool) {
	for _, attr := range node.Attr {
		if attr.Key == name {
			return attr.Val, true
		}
	}
	return "", false
}
