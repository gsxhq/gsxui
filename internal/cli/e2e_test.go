package cli

import (
	"os"
	"os/exec"
	"path/filepath"
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
	for _, p := range []string{"ui/dialog/dialog.x.go", "ui/button/button.x.go"} {
		if _, err := os.Stat(filepath.Join(dir, p)); err != nil {
			t.Fatalf("missing generated %s: %v", p, err)
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
