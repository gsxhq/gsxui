package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestManagedContentHashIsLowercaseSHA256(t *testing.T) {
	got := contentHash([]byte("abc"))
	const want = "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"
	if got != want {
		t.Fatalf("contentHash(abc) = %q, want %q", got, want)
	}
}

func TestManagedArtifactPathRejectsSymlinkEscape(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(root, "ui")); err != nil {
		t.Fatal(err)
	}

	_, err := artifactPath(root, "ui/button.gsx")
	if err == nil || !strings.Contains(err.Error(), "escapes module root") {
		t.Fatalf("artifactPath error = %v, want symlink escape", err)
	}
}

func TestManagedArtifactPathAcceptsMissingNormalizedPath(t *testing.T) {
	root := t.TempDir()
	got, err := artifactPath(root, "components/ui/button.gsx")
	if err != nil {
		t.Fatal(err)
	}
	resolvedRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(resolvedRoot, "components", "ui", "button.gsx")
	if got != want {
		t.Fatalf("artifactPath = %q, want %q", got, want)
	}
}

func TestManagedArtifactPlanRejectsSymlinkAliasedDestinations(t *testing.T) {
	root := t.TempDir()
	assetDir := filepath.Join(root, "web", "gsxui")
	if err := os.MkdirAll(assetDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(assetDir, "style.css"), []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("style.css", filepath.Join(assetDir, "theme.css")); err != nil {
		t.Fatal(err)
	}

	artifacts := []artifact{
		{RelativePath: "web/gsxui/theme.css", Content: []byte("theme"), Managed: true},
		{RelativePath: "web/gsxui/style.css", Content: []byte("style"), Managed: true},
	}
	err := validateArtifactPlan(root, DefaultConfig(), artifacts, true)
	if err == nil || !strings.Contains(err.Error(), "same destination") {
		t.Fatalf("validateArtifactPlan error = %v, want aliased-destination rejection", err)
	}
}
