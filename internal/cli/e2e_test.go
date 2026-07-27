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
	// This temporary consumer is intentionally outside the test runner's workspace.
	t.Setenv("GOWORK", "off")
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
	for _, p := range []string{
		"ui/dialog.x.go",
		"ui/button.x.go",
	} {
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
	const renderTest = `package app_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/gsxhq/gsx"
	"example.com/app/ui"
)

func TestVendoredButtonRenders(t *testing.T) {
	var output bytes.Buffer
	if err := ui.Button("", "", "", false, gsx.Text("Save"), nil).Render(context.Background(), &output); err != nil {
		t.Fatal(err)
	}
	html := output.String()
	if strings.Contains(html, "data-gsxui-slot=") ||
		!strings.Contains(html, "data-gsxui-slot-button") ||
		!strings.Contains(html, ">Save</button>") {
		t.Fatalf("unexpected vendored Button output: %s", html)
	}
}
`
	if err := os.WriteFile(filepath.Join(dir, "render_test.go"), []byte(renderTest), 0o644); err != nil {
		t.Fatal(err)
	}
	mustRun(t, dir, "go", "test", "./...")

	if err := Run([]string{"add", "native-select", "tabs"}); err != nil {
		t.Fatal(err)
	}
	// native-select depends on icon (chevron); icon vendors as its own
	// ui/icon/ directory package. tabs is JS-backed: its behavior module
	// lands under the JS root (web/gsxui by default) and the barrel must be
	// regenerated to import it.
	for _, p := range []string{
		"ui/native-select.x.go",
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

func TestEndToEndCustomUIPath(t *testing.T) {
	if testing.Short() {
		t.Skip("network-dependent e2e; run without -short")
	}
	// This temporary consumer is intentionally outside the test runner's workspace.
	t.Setenv("GOWORK", "off")
	dir := t.TempDir()
	mustRun(t, dir, "go", "mod", "init", "example.com/app")
	t.Chdir(dir)

	cfg := DefaultConfig()
	cfg.UI = "components/ui"
	if err := cfg.Save(dir); err != nil {
		t.Fatal(err)
	}
	if err := Run([]string{"init"}); err != nil {
		t.Fatal(err)
	}
	if err := Run([]string{"add", "spinner"}); err != nil {
		t.Fatal(err)
	}

	for _, path := range []string{
		"components/ui/spinner.x.go",
		"components/ui/icon/icon.x.go",
	} {
		if _, err := os.Stat(filepath.Join(dir, path)); err != nil {
			t.Fatalf("missing generated/vendored %s: %v", path, err)
		}
	}

	for _, path := range []string{
		"components/ui/slots.go",
		"components/ui/internal/slotattr/slotattr.go",
	} {
		if _, err := os.Stat(filepath.Join(dir, path)); !os.IsNotExist(err) {
			t.Fatalf("retired support file %s was vendored: %v", path, err)
		}
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
