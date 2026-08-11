package port

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

var update = flag.Bool("update", false, "update golden files")

// upstreamMaiaPath is the real, locally-checked-out upstream style file this
// package's golden test transforms a real section from. Pinned SHA per the
// plan's Global Constraints; re-verify with
// `git -C /Users/jackieli/personal/shadcn-ui rev-parse HEAD` if the checkout
// ever moves.
const (
	upstreamMaiaPath    = "/Users/jackieli/personal/shadcn-ui/apps/v4/registry/styles/style-maia.css"
	upstreamMaiaRelPath = "apps/v4/registry/styles/style-maia.css"
	upstreamPinnedSHA   = "41bbc12cfd39ed8d9cb8da04275479ee7ecc0612"
)

// TestRenderAccordionMaiaGolden transforms the real Accordion section out of
// the real upstream style-maia.css and asserts Render's output equals the
// committed golden byte-for-byte, following internal/preset/css_test.go's
// golden-file pattern (dossier §7.1). Run with -update to regenerate the
// golden after a deliberate change.
func TestRenderAccordionMaiaGolden(t *testing.T) {
	if _, err := os.Stat(upstreamMaiaPath); err != nil {
		t.Skipf("upstream shadcn-ui checkout not present at %s", upstreamMaiaPath)
	}

	_, sections, err := ReadUpstream(upstreamMaiaPath)
	if err != nil {
		t.Fatalf("ReadUpstream: %v", err)
	}
	var accordion Section
	found := false
	for _, sec := range sections {
		if sec.Name == "Accordion" {
			accordion = sec
			found = true
			break
		}
	}
	if !found {
		t.Fatal("no Accordion section in style-maia.css")
	}

	shape, ok := shapes.All()["accordion"]
	if !ok {
		t.Fatal("no declared shape for accordion")
	}

	// The fallback policy carries a slot's utilities from the CURRENT nova
	// recipe when upstream has no counterpart for it; carryMissingDisplayUtility
	// carries a single missing display-establishing utility even when upstream
	// DID contribute other utilities to that slot; carryStructuralGap carries
	// Accordion trigger-icon's declared structural gap (rotate/transition,
	// never in any of the 8 style files) the same way. Real Maia's Accordion
	// section covers all five declared slots' THEMED content itself, but
	// never mentions "flex" (trigger's display) or the chevron's rotate/
	// transition (trigger-icon) anywhere in the whole style-maia.css file —
	// confirmed against the real upstream source: both are baked into the
	// shared React elements' own class lists, a layer no style-<name>.css
	// file ever carries. So both are legitimately carried from nova here —
	// this fixture is exercised (Transform requires a fallback.Style
	// argument) and these two entries are the correct, expected outcome,
	// asserted below via Ported.Carried.
	novaPath := filepath.Join("..", "..", "..", "registry", "styles", "nova", "accordion.css")
	novaSrc, err := os.ReadFile(novaPath)
	if err != nil {
		t.Fatalf("read nova fallback %s: %v", novaPath, err)
	}
	nova, err := recipe.ParseStyle(novaPath, novaSrc)
	if err != nil {
		t.Fatalf("parse nova fallback: %v", err)
	}

	ported, unmapped, err := Transform(shape, accordion, nova)
	if err != nil {
		t.Fatalf("Transform: %v", err)
	}
	wantCarried := map[string]bool{"gsxui-recipe-accordion-trigger": true, "gsxui-recipe-accordion-trigger-icon": true}
	for class := range wantCarried {
		if !ported.Carried[class] {
			t.Fatalf("Carried = %v, want %v: trigger's display (flex) and trigger-icon's rotate/transition have no upstream counterpart anywhere in style-maia.css", ported.Carried, wantCarried)
		}
	}
	if len(ported.Carried) != len(wantCarried) {
		t.Fatalf("Carried = %v, want exactly %v", ported.Carried, wantCarried)
	}
	// cn-accordion (the upstream root rule) has no gsxui slot -- Accordion's
	// shape declares no root slot at all (registry/canonical/shapes/accordion.go's
	// doc comment) -- a declared, reviewed non-mapping (mapping.go's
	// slotOverrides), not an open question, so it does not show up here.
	if len(unmapped) != 0 {
		t.Fatalf("unmapped = %v, want none", unmapped)
	}

	got, err := Render(ported, "maia", upstreamPinnedSHA, upstreamMaiaRelPath)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}

	goldenPath := filepath.Join("testdata", "golden-accordion-maia.css")
	if *update {
		if err := os.WriteFile(goldenPath, got, 0o644); err != nil {
			t.Fatalf("write golden: %v", err)
		}
	}
	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("read golden %s: %v", goldenPath, err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("Render output does not match golden\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}

	if err := Validate(goldenPath, shape, got); err != nil {
		t.Fatalf("Validate rejected Render's own output: %v", err)
	}
}

// TestValidateRejectsUnknownClass hand-builds recipe bytes carrying a class
// the shape never declares and asserts Validate errors instead of silently
// accepting it -- the load-bearing property from dossier §3.2: Conform
// demands exact bidirectional coverage, no extra classes tolerated.
func TestValidateRejectsUnknownClass(t *testing.T) {
	shape, ok := shapes.All()["accordion"]
	if !ok {
		t.Fatal("no declared shape for accordion")
	}
	src := []byte(`@layer components {
  .gsxui-recipe-accordion-bogus {
    @apply p-4;
  }
}
`)
	err := Validate("bogus.css", shape, src)
	if err == nil {
		t.Fatal("Validate accepted a class the shape doesn't declare, want error")
	}
	if !strings.Contains(err.Error(), "bogus") {
		t.Fatalf("Validate error = %v, want it to mention the unknown class", err)
	}
}
