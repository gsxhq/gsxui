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
// group/<c> root (e.g. group/button) and that alone must not be flagged, only
// data-slot / peer/ / group-*/peer-* conditional variant forms would be.
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
