package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// initTestModule creates a fake module root and stubs the command runner.
func initTestModule(t *testing.T) (dir string, commands *[][]string) {
	t.Helper()
	dir = t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/app\n\ngo 1.26\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var got [][]string
	orig := runCommand
	runCommand = func(dir, name string, args ...string) error {
		got = append(got, append([]string{name}, args...))
		return nil
	}
	t.Cleanup(func() { runCommand = orig })
	t.Chdir(dir)
	return dir, &got
}

func TestInitWritesEverything(t *testing.T) {
	dir, commands := initTestModule(t)
	if err := Run([]string{"init"}); err != nil {
		t.Fatal(err)
	}
	for _, p := range []string{
		"gsxui.json",
		"web/gsxui.css",
		"web/gsxui/core/gsxui.js",
		"web/gsxui/index.js",
		"ui/merge/merge.go",
		"gsx.toml",
	} {
		if _, err := os.Stat(filepath.Join(dir, p)); err != nil {
			t.Errorf("missing %s: %v", p, err)
		}
	}
	toml, _ := os.ReadFile(filepath.Join(dir, "gsx.toml"))
	if want := `class_merger = "example.com/app/ui/merge.Merge"`; !strings.Contains(string(toml), want) {
		t.Errorf("gsx.toml missing %q:\n%s", want, toml)
	}
	css, _ := os.ReadFile(filepath.Join(dir, "web/gsxui.css"))
	if !strings.Contains(string(css), "--primary") {
		t.Error("css does not look like the token file")
	}
	// dependency commands went through the seam
	joined := ""
	for _, c := range *commands {
		joined += strings.Join(c, " ") + "\n"
	}
	for _, want := range []string{
		"go get github.com/gsxhq/gsx@latest",
		"go get github.com/jackielii/tailwind-merge-go@latest",
		"go get -tool github.com/gsxhq/gsx/cmd/gsx@latest",
	} {
		if !strings.Contains(joined, want) {
			t.Errorf("missing command %q in:\n%s", want, joined)
		}
	}
}

func TestInitPreservesExistingGsxToml(t *testing.T) {
	dir, _ := initTestModule(t)
	existing := "[minify]\ncss = true\n"
	os.WriteFile(filepath.Join(dir, "gsx.toml"), []byte(existing), 0o644)
	if err := Run([]string{"init"}); err != nil {
		t.Fatal(err)
	}
	toml, _ := os.ReadFile(filepath.Join(dir, "gsx.toml"))
	s := string(toml)
	if !strings.HasPrefix(s, `class_merger = "example.com/app/ui/merge.Merge"`) {
		t.Errorf("class_merger must be prepended top-level:\n%s", s)
	}
	if !strings.Contains(s, "[minify]") {
		t.Errorf("existing content lost:\n%s", s)
	}
}

func TestInitOutsideModuleRoot(t *testing.T) {
	t.Chdir(t.TempDir())
	err := Run([]string{"init"})
	if err == nil || !strings.Contains(err.Error(), "go.mod not found") {
		t.Fatalf("want module-root error, got %v", err)
	}
}

func TestInitDoesNotClobberUnparsableConfig(t *testing.T) {
	dir, _ := initTestModule(t)
	const broken = `{"ui": "ui", "js": "web/gsxui", "css": "web/gsxui.css",}` // trailing comma
	if err := os.WriteFile(filepath.Join(dir, "gsxui.json"), []byte(broken), 0o644); err != nil {
		t.Fatal(err)
	}
	err := Run([]string{"init"})
	if err == nil || !strings.Contains(err.Error(), "gsxui.json") {
		t.Fatalf("want error mentioning gsxui.json, got %v", err)
	}
	got, readErr := os.ReadFile(filepath.Join(dir, "gsxui.json"))
	if readErr != nil {
		t.Fatal(readErr)
	}
	if string(got) != broken {
		t.Errorf("gsxui.json was modified:\n%s", got)
	}
}
