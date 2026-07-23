package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestEndToEnd exercises the real flow: temp module → init → add dialog →
// go build. Needs network (go get) and the real gsx toolchain; skipped with
// -short.
func TestEndToEnd(t *testing.T) {
	if testing.Short() {
		t.Skip("network-dependent e2e; run without -short")
	}
	dir := t.TempDir()
	mustRun(t, dir, "go", "mod", "init", "example.com/app")
	t.Chdir(dir)

	if err := Run([]string{"init"}); err != nil {
		t.Fatal(err)
	}
	if err := Run([]string{"add", "dialog"}); err != nil {
		t.Fatal(err)
	}
	// generate ran for real: generated files exist
	for _, p := range []string{"ui/dialog.x.go", "ui/button.x.go"} {
		if _, err := os.Stat(filepath.Join(dir, p)); err != nil {
			t.Fatalf("missing generated %s: %v", p, err)
		}
	}
	// vendored .gsx keeps package ui regardless of cfg.UI's basename
	dialogSrc, err := os.ReadFile(filepath.Join(dir, "ui/dialog.gsx"))
	if err != nil {
		t.Fatalf("reading vendored dialog.gsx: %v", err)
	}
	if !strings.Contains(string(dialogSrc), "package ui") {
		t.Fatalf("vendored dialog.gsx missing package ui clause:\n%s", dialogSrc)
	}
	mustRun(t, dir, "go", "build", "./...")

	if err := Run([]string{"add", "select", "tabs"}); err != nil {
		t.Fatal(err)
	}
	// select depends on icon (chevron); icon vendors as its own ui/icon/
	// directory package. tabs is JS-backed: its behavior module lands under
	// the JS root (web/gsxui by default) and the barrel must be regenerated
	// to import it.
	for _, p := range []string{
		"ui/select.x.go",
		"ui/icon/icon.x.go",
		"ui/icon/icon_data.go",
		"ui/tabs.x.go",
		"web/gsxui/tabs.js",
	} {
		if _, err := os.Stat(filepath.Join(dir, p)); err != nil {
			t.Fatalf("missing generated/vendored %s: %v", p, err)
		}
	}
	iconDir, err := os.Stat(filepath.Join(dir, "ui/icon"))
	if err != nil || !iconDir.IsDir() {
		t.Fatalf("ui/icon must vendor as a directory: %v", err)
	}
	barrel, err := os.ReadFile(filepath.Join(dir, "web/gsxui/index.js"))
	if err != nil {
		t.Fatalf("reading barrel: %v", err)
	}
	if !strings.Contains(string(barrel), `"./tabs.js"`) {
		t.Fatalf("barrel index.js missing tabs import:\n%s", barrel)
	}
	mustRun(t, dir, "go", "build", "./...")
}

func mustRun(t *testing.T, dir, name string, args ...string) {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s %v: %v\n%s", name, args, err, out)
	}
}
