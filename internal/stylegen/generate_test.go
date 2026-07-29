package stylegen

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// copyRepoFixture copies the inputs and artifacts GenerateAll touches into dst,
// so a test can corrupt one without mutating the working tree.
func copyRepoFixture(t *testing.T, dst string) {
	t.Helper()
	root := repoRoot(t)
	for _, dir := range []string{"registry", "ui", "assets", filepath.Join("site", "stylepreview")} {
		src := filepath.Join(root, dir)
		err := filepath.WalkDir(src, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			relative, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			target := filepath.Join(dst, relative)
			if entry.IsDir() {
				return os.MkdirAll(target, 0o755)
			}
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			return os.WriteFile(target, content, 0o644)
		})
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestGenerateAllIsIdempotent(t *testing.T) {
	root := t.TempDir()
	copyRepoFixture(t, root)
	if err := GenerateAll(root, false); err != nil {
		t.Fatalf("GenerateAll() error = %v", err)
	}
	if err := GenerateAll(root, true); err != nil {
		t.Fatalf("GenerateAll(check) after write = %v, want nil", err)
	}
}

func TestGenerateAllEmitsEveryStyleAndTheDefaultCopy(t *testing.T) {
	root := t.TempDir()
	copyRepoFixture(t, root)
	if err := GenerateAll(root, false); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{
		"registry/generated/nova/button.gsx",
		"registry/generated/maia/button.gsx",
		"site/stylepreview/nova/button.gsx",
		"site/stylepreview/maia/button.gsx",
		"ui/button.gsx",
	} {
		if _, err := os.Stat(filepath.Join(root, path)); err != nil {
			t.Errorf("missing generated artifact %s: %v", path, err)
		}
	}
	shipped, err := os.ReadFile(filepath.Join(root, "ui", "button.gsx"))
	if err != nil {
		t.Fatal(err)
	}
	fromDefault, err := os.ReadFile(filepath.Join(root, "registry", "generated", DefaultStyle, "button.gsx"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(shipped, fromDefault) {
		t.Error("ui/button.gsx must be a copy of the default style's generated output")
	}

	generated := readGeneratedButtons(t, root)
	for style, src := range generated {
		assertGeneratedButtonSource(t, filepath.Join(root, "registry", "generated", style, "button.gsx"), src)
	}
	assertButtonStylesDifferVisibly(t, generated["nova"], generated["maia"])
}

func TestGeneratedSourcesAreFreeOfRecipeConstructs(t *testing.T) {
	root := t.TempDir()
	copyRepoFixture(t, root)
	if err := GenerateAll(root, false); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"ui/button.gsx", "registry/generated/maia/button.gsx"} {
		src, err := os.ReadFile(filepath.Join(root, path))
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Contains(src, []byte(recipe.Prefix)) {
			t.Errorf("%s contains the recipe prefix", path)
		}
		// Derived from the shape, so a new slot or dimension needs no edit here.
		for _, slot := range shapes.Button.Slots {
			if slot.Base {
				assertNoAccessorResidue(t, src, shapes.Button.Component, accessorName(slot.Name))
			}
			for _, dimension := range slot.Dimensions {
				assertNoAccessorResidue(t, src, shapes.Button.Component,
					dimensionAccessorName(slot.Name, dimension.Name))
			}
		}
	}
}

func TestGenerateAllValidatesBeforeWriting(t *testing.T) {
	root := t.TempDir()
	copyRepoFixture(t, root)
	styleFile := filepath.Join(root, "registry", "styles", "maia", "button.css")
	src, err := os.ReadFile(styleFile)
	if err != nil {
		t.Fatal(err)
	}
	broken := bytes.Replace(src, []byte(".gsxui-recipe-button-size-icon-lg"), []byte(".gsxui-recipe-button-size-icon-xl"), 1)
	if bytes.Equal(broken, src) {
		t.Fatal("fixture does not declare .gsxui-recipe-button-size-icon-lg")
	}
	if err := os.WriteFile(styleFile, broken, 0o644); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(filepath.Join(root, "ui", "button.gsx"))
	if err != nil {
		t.Fatal(err)
	}
	if err := GenerateAll(root, false); err == nil {
		t.Fatal("GenerateAll() = nil, want an error for the undeclared value")
	}
	after, err := os.ReadFile(filepath.Join(root, "ui", "button.gsx"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, after) {
		t.Error("a failing run must not mutate any artifact")
	}
}

// TestGenerateAllRejectsAHelperCallNamingAnUndeclaredDimension pins the
// cross-check GenerateAll performs before it writes anything: a canonical
// calling an accessor its component never declared must fail generation, with
// the offending accessor named, and must leave every artifact untouched.
func TestGenerateAllRejectsAHelperCallNamingAnUndeclaredDimension(t *testing.T) {
	root := t.TempDir()
	copyRepoFixture(t, root)
	canonicalPath := filepath.Join(root, "registry", "canonical", "button.gsx")
	src, err := os.ReadFile(canonicalPath)
	if err != nil {
		t.Fatal(err)
	}
	broken := bytes.Replace(src, []byte("button.Size(size)"), []byte("button.Density(size)"), 1)
	if bytes.Equal(broken, src) {
		t.Fatal("canonical fixture does not call button.Size(size)")
	}
	if err := os.WriteFile(canonicalPath, broken, 0o644); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(filepath.Join(root, "ui", "button.gsx"))
	if err != nil {
		t.Fatal(err)
	}
	err = GenerateAll(root, false)
	if err == nil {
		t.Fatal("GenerateAll() = nil, want an error for the undeclared dimension")
	}
	if !bytes.Contains([]byte(err.Error()), []byte("Density")) {
		t.Errorf("GenerateAll() error %q does not name the undeclared dimension", err)
	}
	after, err := os.ReadFile(filepath.Join(root, "ui", "button.gsx"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, after) {
		t.Error("a failing run must not mutate any artifact")
	}
}

// cardCoverageShape is a two-slot shape with a dimension, used to exercise
// both directions of checkHelperCalls against real GSX source.
func cardCoverageShape() recipe.Shape {
	return recipe.Shape{
		Component: "card",
		Slots: []recipe.Slot{
			{Name: "", Base: true, Dimensions: []recipe.Dimension{
				{Name: "tone", Default: "plain", Values: []string{"plain", "muted"}},
			}},
			{Name: "footer", Base: true},
		},
	}
}

const cardCoverageSource = `package canonical

import "github.com/gsxhq/gsx"

component Card(tone string, children gsx.Node) {
	<div class={ card.Root(), card.Tone(tone) }>
		{ children }
		<div class={ card.Footer() }></div>
	</div>
}
`

// TestCheckHelperCallsAcceptsFullSlotCoverage pins the passing side: every
// declared slot and dimension is reached by an accessor call.
func TestCheckHelperCallsAcceptsFullSlotCoverage(t *testing.T) {
	t.Parallel()
	if err := checkHelperCalls("card.gsx", []byte(cardCoverageSource), cardCoverageShape()); err != nil {
		t.Fatalf("checkHelperCalls() = %v, want nil", err)
	}
}

// TestCheckHelperCallsRejectsAnUndeclaredSlot pins the direction that already
// worked: an accessor call naming a slot the shape does not declare.
func TestCheckHelperCallsRejectsAnUndeclaredSlot(t *testing.T) {
	t.Parallel()
	src := []byte(`package canonical

import "github.com/gsxhq/gsx"

component Card(tone string, children gsx.Node) {
	<div class={ card.Root(), card.Tone(tone) }>
		{ children }
		<div class={ card.Footer() }></div>
		<div class={ card.Header() }></div>
	</div>
}
`)
	err := checkHelperCalls("card.gsx", src, cardCoverageShape())
	if err == nil {
		t.Fatal("checkHelperCalls() = nil, want an error for the undeclared slot")
	}
	if !strings.Contains(err.Error(), "Header") && !strings.Contains(err.Error(), "header") {
		t.Errorf("error must name the undeclared accessor, got %q", err)
	}
}

// TestCheckHelperCallsRejectsAnUnrenderedSlot pins the new direction. Without
// it a shape could declare a slot no .gsx renders: the style would implement
// it, conformance would pass, and the resulting recipe rules would be dead.
func TestCheckHelperCallsRejectsAnUnrenderedSlot(t *testing.T) {
	t.Parallel()
	src := []byte(`package canonical

import "github.com/gsxhq/gsx"

component Card(tone string, children gsx.Node) {
	<div class={ card.Root(), card.Tone(tone) }>{ children }</div>
}
`)
	err := checkHelperCalls("registry/canonical/card.gsx", src, cardCoverageShape())
	if err == nil {
		t.Fatal("checkHelperCalls() = nil, want an error for the unrendered slot")
	}
	want := `registry/canonical/card.gsx: shape declares slot "footer" but the component never renders it`
	if !strings.Contains(err.Error(), want) {
		t.Errorf("checkHelperCalls() = %q, want it to contain %q", err, want)
	}
}

// TestCheckHelperCallsRejectsAnUnusedDimension — a dimension is reached by a
// single accessor call however many values it has, so an unused one is the
// same dead-rule signal as an unused slot, and every one of its values costs
// the style a rule that conformance forces it to author.
func TestCheckHelperCallsRejectsAnUnusedDimension(t *testing.T) {
	t.Parallel()
	src := []byte(`package canonical

import "github.com/gsxhq/gsx"

component Card(tone string, children gsx.Node) {
	<div class={ card.Root() }>
		{ children }
		<div class={ card.Footer() }></div>
	</div>
}
`)
	err := checkHelperCalls("registry/canonical/card.gsx", src, cardCoverageShape())
	if err == nil {
		t.Fatal("checkHelperCalls() = nil, want an error for the unused dimension")
	}
	want := `registry/canonical/card.gsx: shape declares dimension "tone" on slot "" but the component never applies it`
	if !strings.Contains(err.Error(), want) {
		t.Errorf("checkHelperCalls() = %q, want it to contain %q", err, want)
	}
}

// TestCheckHelperCallsRejectsAnUnusedBaseRule closes the last hole: a slot can
// be reached through a dimension accessor while its base rule stays dead.
func TestCheckHelperCallsRejectsAnUnusedBaseRule(t *testing.T) {
	t.Parallel()
	src := []byte(`package canonical

import "github.com/gsxhq/gsx"

component Card(tone string, children gsx.Node) {
	<div class={ card.Tone(tone) }>
		{ children }
		<div class={ card.Footer() }></div>
	</div>
}
`)
	err := checkHelperCalls("registry/canonical/card.gsx", src, cardCoverageShape())
	if err == nil {
		t.Fatal("checkHelperCalls() = nil, want an error for the unused base rule")
	}
	if !strings.Contains(err.Error(), "base rule") {
		t.Errorf("checkHelperCalls() = %q, want it to mention the unused base rule", err)
	}
}

// TestCheckHelperCallsReportsEveryGap keeps the diagnostic a sweep, not a
// bisect: a 38-slot component should not need 38 runs to converge.
func TestCheckHelperCallsReportsEveryGap(t *testing.T) {
	t.Parallel()
	src := []byte(`package canonical

import "github.com/gsxhq/gsx"

component Card(tone string, children gsx.Node) {
	<div class={ card.Root() }>{ children }</div>
}
`)
	err := checkHelperCalls("registry/canonical/card.gsx", src, cardCoverageShape())
	if err == nil {
		t.Fatal("checkHelperCalls() = nil, want errors")
	}
	if !strings.Contains(err.Error(), `slot "footer"`) || !strings.Contains(err.Error(), `dimension "tone"`) {
		t.Errorf("checkHelperCalls() = %q, want every gap reported", err)
	}
}
