package recipe

import (
	"bytes"
	"encoding/json"
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
	if strings.Count(string(got), `"values"`) != 1 {
		t.Errorf("shape must be emitted exactly once, got:\n%s", got)
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
