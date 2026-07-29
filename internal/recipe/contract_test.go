package recipe

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestBuildContractEmitsSharedShapeAndPerStyleUtilities(t *testing.T) {
	t.Parallel()
	shape := conformShape()
	resolved, err := Conform("nova/button.css", shape, mustParse(t, conformCSS))
	if err != nil {
		t.Fatal(err)
	}
	contract := BuildContract(
		map[string]Shape{"button": shape},
		map[string]map[string]Resolved{"nova": {"button": resolved}},
	)
	got, err := contract.MarshalIndent()
	if err != nil {
		t.Fatal(err)
	}
	var round map[string]any
	if err := json.Unmarshal(got, &round); err != nil {
		t.Fatalf("contract is not valid JSON: %v", err)
	}
	for _, want := range []string{
		`"version": 1`,
		`"slots"`,
		`"components"`,
		`"default": "default"`,
		`"outline"`,
		`"styles"`,
		`"nova"`,
		`"border-border"`,
	} {
		if !strings.Contains(string(got), want) {
			t.Errorf("contract missing %q\nin: %s", want, got)
		}
	}
	// The shape must appear once, under components — never repeated per style.
	// One "values" key per declared dimension, derived from the fixture rather
	// than hardcoded, so a second dimension does not silently pass.
	dimensionCount := 0
	for _, slot := range shape.Slots {
		dimensionCount += len(slot.Dimensions)
	}
	if count, want := strings.Count(string(got), `"values"`), dimensionCount; count != want {
		t.Errorf("shape emitted %d times, want %d:\n%s", count, want, got)
	}
}

func TestContractIsDeterministic(t *testing.T) {
	t.Parallel()
	shape := conformShape()
	resolved, err := Conform("nova/button.css", shape, mustParse(t, conformCSS))
	if err != nil {
		t.Fatal(err)
	}
	build := func() []byte {
		contract := BuildContract(
			map[string]Shape{"button": shape},
			map[string]map[string]Resolved{"nova": {"button": resolved}, "maia": {"button": resolved}},
		)
		out, err := contract.MarshalIndent()
		if err != nil {
			t.Fatal(err)
		}
		return out
	}
	if !bytes.Equal(build(), build()) {
		t.Error("contract marshalling must be deterministic across map iteration orders")
	}
}

func TestBuildContractGroupsSlotsUnderComponent(t *testing.T) {
	t.Parallel()
	shape := slotConformShape()
	resolved, err := Conform("nova/card.css", shape, mustParse(t, slotConformCSS))
	if err != nil {
		t.Fatal(err)
	}
	out, err := BuildContract(
		map[string]Shape{"card": shape},
		map[string]map[string]Resolved{"nova": {"card": resolved}},
	).MarshalIndent()
	if err != nil {
		t.Fatal(err)
	}

	var parsed Contract
	if err := json.Unmarshal(out, &parsed); err != nil {
		t.Fatalf("contract is not valid JSON: %v", err)
	}
	if parsed.Version != 2 {
		t.Errorf("Version = %d, want 2", parsed.Version)
	}
	card, ok := parsed.Components["card"]
	if !ok {
		t.Fatal("components.card missing")
	}
	if _, ok := card.Slots["header"]; !ok {
		t.Error("components.card.slots.header missing — slots must group under their component")
	}
	if got := card.Slots["header"].Dimensions["variant"].Default; got != "default" {
		t.Errorf("header variant default = %q, want %q", got, "default")
	}
	// The shape must appear once, under components — never repeated per style.
	if got, want := strings.Count(string(out), `"values"`), 1; got != want {
		t.Errorf("shape emitted %d times, want %d:\n%s", got, want, out)
	}
	nova := parsed.Styles["nova"]["card"]
	if got, want := nova.Slots["header"].Dimensions["variant"]["muted"], []string{"text-muted-foreground"}; !reflect.DeepEqual(got, want) {
		t.Errorf("nova header muted = %q, want %q", got, want)
	}
}

func TestContractIsDeterministicWithSlots(t *testing.T) {
	t.Parallel()
	shape := slotConformShape()
	resolved, err := Conform("nova/card.css", shape, mustParse(t, slotConformCSS))
	if err != nil {
		t.Fatal(err)
	}
	build := func() []byte {
		out, err := BuildContract(
			map[string]Shape{"card": shape},
			map[string]map[string]Resolved{"nova": {"card": resolved}, "maia": {"card": resolved}},
		).MarshalIndent()
		if err != nil {
			t.Fatal(err)
		}
		return out
	}
	if !bytes.Equal(build(), build()) {
		t.Error("contract marshalling must be deterministic across map iteration orders")
	}
}
