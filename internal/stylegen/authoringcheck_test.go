package stylegen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestCheckAuthoringRejectsUnmigratedHandAuthoredClass pins the class=
// prohibition: an unmigrated component (no entry in shapes.All()) must never
// carry a hand-authored class= attribute.
func TestCheckAuthoringRejectsUnmigratedHandAuthoredClass(t *testing.T) {
	root := t.TempDir()
	copyRepoFixture(t, root)

	path := filepath.Join(root, "ui", "zzz-unmigrated-class.gsx")
	src := "package ui\n\nfunc TestUnmig() {\n\t<div class=\"foo\">hi</div>\n}\n"
	if err := os.WriteFile(path, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}

	err := CheckAuthoring(root)
	if err == nil {
		t.Fatal("CheckAuthoring() = nil, want error for hand-authored class= in an unmigrated component")
	}
	if !strings.Contains(err.Error(), "zzz-unmigrated-class.gsx") {
		t.Fatalf("CheckAuthoring() error = %v, want it to name the offending file", err)
	}
}

// TestCheckAuthoringRejectsUnmigratedBareGroup pins the group/ prohibition:
// an unmigrated component must never carry a bare group/<name> root.
func TestCheckAuthoringRejectsUnmigratedBareGroup(t *testing.T) {
	root := t.TempDir()
	copyRepoFixture(t, root)

	path := filepath.Join(root, "ui", "zzz-unmigrated-group.gsx")
	src := "package ui\n\n// no class= here, just a bare group/ root\nconst marker = \"group/foo\"\n"
	if err := os.WriteFile(path, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}

	err := CheckAuthoring(root)
	if err == nil {
		t.Fatal("CheckAuthoring() = nil, want error for bare group/foo in an unmigrated component")
	}
	if !strings.Contains(err.Error(), "zzz-unmigrated-group.gsx") {
		t.Fatalf("CheckAuthoring() error = %v, want it to name the offending file", err)
	}
}

// TestCheckAuthoringAcceptsMigratedGroupButton pins the narrow pattern set: a
// migrated component's generated ui/<c>.gsx legitimately carries its own
// group/<c> root (e.g. group/button) and that alone must not be flagged.
func TestCheckAuthoringAcceptsMigratedGroupButton(t *testing.T) {
	root := t.TempDir()
	copyRepoFixture(t, root)

	// ui/button.gsx is real generated output and is already known to carry
	// group/button (see registry/generated/nova/button.gsx). Confirm the
	// fixture actually exercises that before trusting a green result.
	content, err := os.ReadFile(filepath.Join(root, "ui", "button.gsx"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "group/button") {
		t.Fatal("fixture ui/button.gsx does not contain group/button; test no longer exercises the narrow pattern set")
	}

	if err := CheckAuthoring(root); err != nil {
		t.Fatalf("CheckAuthoring() = %v, want nil: a migrated component's own group/<c> root is compiled structure, not hand-authored presentation", err)
	}
}

// TestCheckAuthoringAcceptsMigratedGroupConditionalForms is the sibling case
// authoringPatternNarrow used to get wrong: a migrated component's compiled
// output legitimately carries group-*/peer-* CONDITIONAL variant forms too —
// not just the bare group/<c> root — when one of its own slots reacts to
// another's state. ui/toggle-group.gsx is real generated output that does
// exactly this (an item keys off group-data-[spacing=0]/toggle-group:,
// matching upstream's own toggle-group.tsx composition), so it pins the case
// directly rather than via a synthetic fixture.
//
// Switch pinned this case originally (the thumb keyed off
// group-data-[size=…]/switch:), but Switch's own ADAPT — a real native
// <input type="checkbox" role="switch"> with the thumb as this same
// element's own ::before pseudo-element, not a separate sibling Radix used
// to need a group/peer relationship for — retired that pattern for Switch
// specifically (ui/switch.gsx now has no group-*/peer-* form at all: see its
// own doc comment). ToggleGroup's spacing=0 join behavior still needs the
// real cross-slot relationship this test exists to cover.
func TestCheckAuthoringAcceptsMigratedGroupConditionalForms(t *testing.T) {
	root := t.TempDir()
	copyRepoFixture(t, root)

	content, err := os.ReadFile(filepath.Join(root, "ui", "toggle-group.gsx"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "group-data-") || !strings.Contains(string(content), "/toggle-group:") {
		t.Fatal("fixture ui/toggle-group.gsx does not contain a group-data-…/toggle-group: conditional form; test no longer exercises this case")
	}

	if err := CheckAuthoring(root); err != nil {
		t.Fatalf("CheckAuthoring() = %v, want nil: a migrated component's own group-*/peer-* conditional forms are compiled structure, not hand-authored presentation", err)
	}
}

// TestCheckAuthoringRejectsMigratedDataSlotLeak pins the one pattern that
// stays checked for a migrated component even though it is generated,
// byte-identity-verified output: a literal data-slot (upstream's own raw
// attribute name — gsxui's is data-gsxui-slot-*). Upstream's source CSS
// embeds data-slot= inside arbitrary-variant selector STRINGS the porter's
// class-name-based re-slotting does not parse into (e.g.
// has-[[data-slot=input-group-control]:focus-visible]:... in
// style-nova.css), so a leak there would compile clean and pass the
// byte-identity check too — this is a real, standing risk the migrated path
// must keep catching, not something the byte-identity check subsumes.
func TestCheckAuthoringRejectsMigratedDataSlotLeak(t *testing.T) {
	root := t.TempDir()
	copyRepoFixture(t, root)

	path := filepath.Join(root, "ui", "button.gsx")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	leaked := strings.Replace(string(content), "data-gsxui-slot-button", "data-slot-button data-gsxui-slot-button", 1)
	if leaked == string(content) {
		t.Fatal("fixture ui/button.gsx does not contain data-gsxui-slot-button to mutate")
	}
	// Also regenerate the DefaultStyle counterpart identically, so this test
	// exercises the pattern check specifically rather than tripping the
	// separate byte-identity check first.
	genPath := filepath.Join(root, "registry", "generated", DefaultStyle, "button.gsx")
	if err := os.WriteFile(genPath, []byte(leaked), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(leaked), 0o644); err != nil {
		t.Fatal(err)
	}

	err = CheckAuthoring(root)
	if err == nil {
		t.Fatal("CheckAuthoring() = nil, want error for a leaked data-slot in a migrated component")
	}
	if !strings.Contains(err.Error(), "data-slot") {
		t.Fatalf("CheckAuthoring() error = %v, want it to name the data-slot leak", err)
	}
}

// TestCheckAuthoringRejectsMigratedDivergingFromGenerated pins the
// byte-identity guarantee: a migrated component's ui/<c>.gsx must match
// registry/generated/<DefaultStyle>/<c>.gsx exactly.
func TestCheckAuthoringRejectsMigratedDivergingFromGenerated(t *testing.T) {
	root := t.TempDir()
	copyRepoFixture(t, root)

	path := filepath.Join(root, "ui", "button.gsx")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(content, []byte("\n// mutated\n")...), 0o644); err != nil {
		t.Fatal(err)
	}

	err = CheckAuthoring(root)
	if err == nil {
		t.Fatal("CheckAuthoring() = nil, want error: ui/button.gsx no longer matches registry/generated/nova/button.gsx")
	}
	if !strings.Contains(err.Error(), "byte-identical") {
		t.Fatalf("CheckAuthoring() error = %v, want it to report the byte-identity mismatch", err)
	}
}

// TestCheckAuthoringRejectsDeadCheckedVocabularyInStyles pins the
// registry/styles gate: a hand-ported sheet must never carry upstream's
// data-checked/data-unchecked variant vocabulary — nothing in gsxui stamps
// data-checked, so a leaked token compiles clean and silently never matches.
func TestCheckAuthoringRejectsDeadCheckedVocabularyInStyles(t *testing.T) {
	root := t.TempDir()
	copyRepoFixture(t, root)

	path := filepath.Join(root, "registry", "styles", DefaultStyle, "field.css")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	leaked := strings.Replace(string(content), "has-[input:checked]:bg-primary/5", "has-data-checked:bg-primary/5", 1)
	if leaked == string(content) {
		t.Fatal("fixture field.css does not contain has-[input:checked]:bg-primary/5 to mutate")
	}
	if err := os.WriteFile(path, []byte(leaked), 0o644); err != nil {
		t.Fatal(err)
	}

	err = CheckAuthoring(root)
	if err == nil {
		t.Fatal("CheckAuthoring() = nil, want error for data-checked vocabulary in a registry/styles sheet")
	}
	if !strings.Contains(err.Error(), "data-checked") || !strings.Contains(err.Error(), "field.css") {
		t.Fatalf("CheckAuthoring() error = %v, want it to name the data-checked leak in field.css", err)
	}
}

// TestCheckAuthoringRejectsDeadStateVocabularyInDetailsStyles pins the
// native-<details> half of the registry/styles gate: accordion's and
// collapsible's sheets must never carry upstream's data-[state=…] variant
// vocabulary. Both components are real <details>/<summary> elements with no
// behavior module, so nothing stamps data-state on them and the port's
// translation is the native `open:` variant — a leaked data-[state=open:] token
// compiles clean and silently never matches.
func TestCheckAuthoringRejectsDeadStateVocabularyInDetailsStyles(t *testing.T) {
	root := t.TempDir()
	copyRepoFixture(t, root)

	path := filepath.Join(root, "registry", "styles", "luma", "accordion.css")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	leaked := strings.Replace(string(content), "open:bg-muted/50", "data-[state=open]:bg-muted/50", 1)
	if leaked == string(content) {
		t.Fatal("fixture luma/accordion.css does not contain open:bg-muted/50 to mutate")
	}
	if err := os.WriteFile(path, []byte(leaked), 0o644); err != nil {
		t.Fatal(err)
	}

	err = CheckAuthoring(root)
	if err == nil {
		t.Fatal("CheckAuthoring() = nil, want error for data-[state= vocabulary in accordion.css")
	}
	if !strings.Contains(err.Error(), "data-[state=") || !strings.Contains(err.Error(), "accordion.css") {
		t.Fatalf("CheckAuthoring() error = %v, want it to name the data-[state= leak in accordion.css", err)
	}
}

// TestCheckAuthoringAllowsDataStateOutsideDetailsStyles keeps the new gate
// narrow: dialog, sheet, drawer and every other component with a behavior
// module really do stamp data-state, so their sheets must stay free to use it.
func TestCheckAuthoringAllowsDataStateOutsideDetailsStyles(t *testing.T) {
	root := t.TempDir()
	copyRepoFixture(t, root)

	path := filepath.Join(root, "registry", "styles", DefaultStyle, "dialog.css")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "data-[state=open]") {
		t.Fatal("fixture dialog.css carries no data-[state=open] token to exercise the exemption")
	}

	if err := CheckAuthoring(root); err != nil {
		t.Fatalf("CheckAuthoring() = %v, want nil — data-[state= is live vocabulary for dialog", err)
	}
}
