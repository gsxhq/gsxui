package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// initedModule = initTestModule + a completed init (stubbed runner).
func initedModule(t *testing.T) (string, *[][]string) {
	t.Helper()
	dir, commands := initTestModule(t)
	if err := Run([]string{"init"}); err != nil {
		t.Fatal(err)
	}
	return dir, commands
}

func TestAddVendorsWithDeps(t *testing.T) {
	dir, commands := initedModule(t)
	if err := Run([]string{"add", "dialog"}); err != nil {
		t.Fatal(err)
	}
	// dialog pulls button transitively
	for _, p := range []string{
		"ui/dialog.gsx",
		"ui/button.gsx",
		"web/gsxui/dialog.js",
		"ui/NOTICE.md",
	} {
		if _, err := os.Stat(filepath.Join(dir, p)); err != nil {
			t.Errorf("missing %s", p)
		}
	}
	// no generated or test files vendored
	if _, err := os.Stat(filepath.Join(dir, "ui/dialog.x.go")); err == nil {
		t.Error("dialog.x.go must not be vendored")
	}
	// package clause kept as-is; no unrewritten gsxui-internal refs remain
	// (dialog.gsx has no cross-package import to rewrite — it's flat, so
	// dialog's use of Button is an intra-package identifier reference; the
	// icon-import rewrite path is covered by TestRewriteGsxIcon and the e2e
	// test's ui/icon vendoring)
	gsx, _ := os.ReadFile(filepath.Join(dir, "ui/dialog.gsx"))
	if strings.Contains(string(gsx), "gsxhq/gsxui") {
		t.Errorf("unrewritten import remains:\n%s", gsx)
	}
	if !strings.Contains(string(gsx), "package ui") {
		t.Errorf("vendored dialog.gsx missing package ui clause:\n%s", gsx)
	}
	js, _ := os.ReadFile(filepath.Join(dir, "web/gsxui/dialog.js"))
	if !strings.Contains(string(js), `"./gsxui.js"`) {
		t.Errorf("core import not present:\n%s", js)
	}
	// barrel updated
	index, _ := os.ReadFile(filepath.Join(dir, "web/gsxui/index.js"))
	if !strings.Contains(string(index), `import "./dialog.js";`) {
		t.Errorf("barrel missing dialog:\n%s", index)
	}
	if strings.Contains(string(index), `import "./gsxui.js";`) {
		t.Errorf("core gsxui.js must not be listed as a behavior import:\n%s", index)
	}
	// gsx generate ran
	joined := ""
	for _, c := range *commands {
		joined += strings.Join(c, " ") + "\n"
	}
	if !strings.Contains(joined, "go tool gsx generate") {
		t.Errorf("generate not invoked:\n%s", joined)
	}
}

func TestAddRefusesModifiedFile(t *testing.T) {
	dir, _ := initedModule(t)
	if err := Run([]string{"add", "badge"}); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(dir, "ui/badge.gsx")
	os.WriteFile(target, []byte("package badge // locally hacked\n"), 0o644)
	err := Run([]string{"add", "badge"})
	if err == nil || !strings.Contains(err.Error(), "--overwrite") {
		t.Fatalf("want overwrite-refusal error, got %v", err)
	}
	if err := Run([]string{"add", "--overwrite", "badge"}); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(target)
	if strings.Contains(string(got), "locally hacked") {
		t.Error("overwrite did not replace the file")
	}
}

func TestAddIsIdempotent(t *testing.T) {
	_, _ = initedModule(t)
	if err := Run([]string{"add", "card"}); err != nil {
		t.Fatal(err)
	}
	if err := Run([]string{"add", "card"}); err != nil {
		t.Fatalf("re-add of unmodified component must succeed: %v", err)
	}
}

func TestAddUnknown(t *testing.T) {
	_, _ = initedModule(t)
	err := Run([]string{"add", "nope"})
	if err == nil || !strings.Contains(err.Error(), "unknown component") {
		t.Fatalf("want unknown-component error, got %v", err)
	}
}

func TestAddGenerateFailureHint(t *testing.T) {
	_, _ = initedModule(t)
	orig := runCommand
	runCommand = func(dir, name string, args ...string) error {
		if name == "go" && len(args) > 0 && args[0] == "tool" {
			return fmt.Errorf("exit status 1")
		}
		return nil
	}
	t.Cleanup(func() { runCommand = orig })
	err := Run([]string{"add", "badge"})
	if err == nil {
		t.Fatal("want error when gsx generate fails")
	}
	if !strings.Contains(err.Error(), "gsx generate:") || !strings.Contains(err.Error(), "gsxui init") {
		t.Fatalf("want actionable hint, got %v", err)
	}
}

func TestAddRejectsCore(t *testing.T) {
	_, _ = initedModule(t)
	err := Run([]string{"add", "core"})
	if err == nil || !strings.Contains(err.Error(), "unknown component") {
		t.Fatalf("want unknown-component error for core, got %v", err)
	}
}

func TestAddRefusesCustomBarrel(t *testing.T) {
	dir, _ := initedModule(t)
	indexPath := filepath.Join(dir, "web/gsxui/index.js")
	custom := "// hand-written, thanks\nexport * from \"./gsxui.js\";\n"
	if err := os.WriteFile(indexPath, []byte(custom), 0o644); err != nil {
		t.Fatal(err)
	}
	err := Run([]string{"add", "badge"})
	if err == nil || !strings.Contains(err.Error(), "--overwrite") {
		t.Fatalf("want overwrite-refusal error, got %v", err)
	}
	got, _ := os.ReadFile(indexPath)
	if string(got) != custom {
		t.Errorf("custom index.js was modified:\n%s", got)
	}
	if err := Run([]string{"add", "--overwrite", "badge"}); err != nil {
		t.Fatal(err)
	}
	got, _ = os.ReadFile(indexPath)
	if !strings.HasPrefix(string(got), barrelHeader) {
		t.Errorf("--overwrite should replace with the generated barrel:\n%s", got)
	}
}

func TestAddRegeneratesGeneratedBarrelWithoutOverwrite(t *testing.T) {
	dir, _ := initedModule(t)
	// index.js from init already carries the generated header; adding a
	// component with JS must regenerate it without needing --overwrite.
	if err := Run([]string{"add", "dialog"}); err != nil {
		t.Fatal(err)
	}
	index, _ := os.ReadFile(filepath.Join(dir, "web/gsxui/index.js"))
	if !strings.Contains(string(index), `import "./dialog.js";`) {
		t.Errorf("generated barrel not regenerated:\n%s", index)
	}
}
