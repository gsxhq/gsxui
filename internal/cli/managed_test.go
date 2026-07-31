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
	if err == nil || !strings.Contains(err.Error(), "unsafe symlink") {
		t.Fatalf("artifactPath error = %v, want symlink escape", err)
	}
}

func TestManagedArtifactPathRejectsNonDirectoryAncestor(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "generated"), []byte("not a directory"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := artifactPath(root, "generated/button.gsx")
	if err == nil || !strings.Contains(err.Error(), "non-directory ancestor") {
		t.Fatalf("artifactPath error = %v, want non-directory ancestor rejection", err)
	}
}

func TestManagedArtifactPathRejectsSymlinkAncestorAndLeaf(t *testing.T) {
	for _, location := range []string{"ancestor", "leaf"} {
		t.Run(location, func(t *testing.T) {
			root := t.TempDir()
			realDir := filepath.Join(root, "real")
			if err := os.Mkdir(realDir, 0o755); err != nil {
				t.Fatal(err)
			}
			var link, target, relative string
			if location == "ancestor" {
				link = filepath.Join(root, "ui")
				target = realDir
				relative = "ui/button.gsx"
			} else {
				target = filepath.Join(realDir, "button.gsx")
				if err := os.WriteFile(target, []byte("user bytes"), 0o644); err != nil {
					t.Fatal(err)
				}
				if err := os.Mkdir(filepath.Join(root, "ui"), 0o755); err != nil {
					t.Fatal(err)
				}
				link = filepath.Join(root, "ui", "button.gsx")
				relative = "ui/button.gsx"
			}
			if err := os.Symlink(target, link); err != nil {
				t.Fatal(err)
			}

			_, err := artifactPath(root, relative)
			if err == nil || !strings.Contains(err.Error(), "symlink") {
				t.Fatalf("artifactPath error = %v, want %s symlink rejection", err, location)
			}
		})
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
	if err == nil || !strings.Contains(err.Error(), "unsafe symlink") {
		t.Fatalf("validateArtifactPlan error = %v, want aliased-destination rejection", err)
	}
}

func TestManagedArtifactPlanAncestorDetectionUsesWholePathSegments(t *testing.T) {
	root := t.TempDir()
	cfg := DefaultConfig()
	first := artifact{RelativePath: "generated/merge.go", Content: []byte("file"), Managed: true}
	descendant := artifact{
		RelativePath: "generated/merge.go/merge/merge.go",
		Content:      []byte("descendant"),
		Managed:      true,
	}

	for _, artifacts := range [][]artifact{
		{first, descendant},
		{descendant, first},
	} {
		err := validateArtifactPlan(root, cfg, artifacts, true)
		if err == nil || !strings.Contains(err.Error(), "ancestor") {
			t.Fatalf("validateArtifactPlan error = %v, want ancestor rejection", err)
		}
	}

	siblingPrefix := artifact{
		RelativePath: "generated/merge.go-backup/merge.go",
		Content:      []byte("sibling"),
		Managed:      true,
	}
	if err := validateArtifactPlan(root, cfg, []artifact{first, siblingPrefix}, true); err != nil {
		t.Fatalf("segment-distinct paths were treated as ancestors: %v", err)
	}
}

func TestManagedArtifactPlanRejectsPortableExactAliases(t *testing.T) {
	tests := []struct {
		name  string
		first string
		last  string
	}{
		{
			name:  "case only",
			first: "Web/GSXUI/index.js",
			last:  "web/gsxui/index.js",
		},
		{
			name:  "canonical Unicode",
			first: "web/caf\u00e9/index.js",
			last:  "web/cafe\u0301/index.js",
		},
		{
			name:  "full case fold expansion",
			first: "web/Stra\u00dfe/index.js",
			last:  "web/STRASSE/index.js",
		},
		{
			name:  "folded final sigma",
			first: "web/\u03a3/index.js",
			last:  "web/\u03c2/index.js",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			artifacts := []artifact{
				{RelativePath: tt.first, Content: []byte("first"), Managed: true},
				{RelativePath: tt.last, Content: []byte("last"), Managed: true},
			}
			err := validateArtifactPlan(t.TempDir(), DefaultConfig(), artifacts, true)
			if err == nil || !strings.Contains(err.Error(), "portable") {
				t.Fatalf("validateArtifactPlan error = %v, want portable exact-alias rejection", err)
			}
		})
	}
}

func TestManagedArtifactPlanRejectsPortableAncestorsInBothOrders(t *testing.T) {
	ancestor := artifact{
		RelativePath: "Web/Stra\u00dfe",
		Content:      []byte("ancestor"),
		Managed:      true,
	}
	descendant := artifact{
		RelativePath: "web/STRASSE/cafe\u0301/button.gsx",
		Content:      []byte("descendant"),
		Managed:      true,
	}
	for _, artifacts := range [][]artifact{
		{ancestor, descendant},
		{descendant, ancestor},
	} {
		err := validateArtifactPlan(t.TempDir(), DefaultConfig(), artifacts, true)
		if err == nil || !strings.Contains(err.Error(), "portable") ||
			!strings.Contains(err.Error(), "ancestor") {
			t.Fatalf("validateArtifactPlan error = %v, want portable ancestor rejection", err)
		}
	}
}

func TestManagedArtifactPlanPortableIdentityIncludesConfigDestination(t *testing.T) {
	cfg := DefaultConfig()
	_, complete, err := artifactPlanWithConfig(cfg, []artifact{{
		RelativePath: "GSXUI.JSON",
		Content:      []byte("component bytes"),
		Managed:      true,
	}})
	if err != nil {
		t.Fatal(err)
	}
	err = validateArtifactPlan(t.TempDir(), cfg, complete, true)
	if err == nil || !strings.Contains(err.Error(), "portable") {
		t.Fatalf("validateArtifactPlan error = %v, want portable config alias rejection", err)
	}
}

func TestManagedArtifactPlanAllowsPortablePrefixSiblings(t *testing.T) {
	artifacts := []artifact{
		{RelativePath: "Web/Stra\u00dfe/file", Content: []byte("first"), Managed: true},
		{RelativePath: "web/STRASSEN/file", Content: []byte("last"), Managed: true},
	}
	if err := validateArtifactPlan(t.TempDir(), DefaultConfig(), artifacts, true); err != nil {
		t.Fatalf("non-equivalent prefix siblings were rejected: %v", err)
	}
}

func TestWriteArtifactPlanAtomicallyReplacesOnlyOneHardlinkAndPreservesMode(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "ui", "button.gsx")
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	const oldContent = "user-owned shared bytes\n"
	if err := os.WriteFile(target, []byte(oldContent), 0o751); err != nil {
		t.Fatal(err)
	}
	alias := filepath.Join(root, "user-button.gsx")
	if err := os.Link(target, alias); err != nil {
		t.Fatal(err)
	}

	const newContent = "package ui\n"
	if err := writeArtifactPlan(root, []artifact{{
		RelativePath: "ui/button.gsx",
		Content:      []byte(newContent),
		Managed:      true,
	}}); err != nil {
		t.Fatal(err)
	}

	gotTarget, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(gotTarget) != newContent {
		t.Fatalf("target bytes = %q, want %q", gotTarget, newContent)
	}
	gotAlias, err := os.ReadFile(alias)
	if err != nil {
		t.Fatal(err)
	}
	if string(gotAlias) != oldContent {
		t.Fatalf("hardlink alias bytes = %q, want preserved %q", gotAlias, oldContent)
	}
	targetInfo, err := os.Stat(target)
	if err != nil {
		t.Fatal(err)
	}
	aliasInfo, err := os.Stat(alias)
	if err != nil {
		t.Fatal(err)
	}
	if os.SameFile(targetInfo, aliasInfo) {
		t.Fatal("atomic replacement left target sharing the user-owned inode")
	}
	if targetInfo.Mode().Perm() != 0o751 {
		t.Fatalf("target mode = %v, want preserved 0751", targetInfo.Mode().Perm())
	}
}
