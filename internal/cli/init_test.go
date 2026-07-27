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
		"web/gsxui/index.css",
		"web/gsxui/foundation.css",
		"web/gsxui/theme.css",
		"web/gsxui/style.css",
		"web/gsxui/gsxui.js",
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
	indexCSS, _ := os.ReadFile(filepath.Join(dir, "web/gsxui/index.css"))
	for _, want := range []string{"./foundation.css", "./theme.css", "./style.css"} {
		if !strings.Contains(string(indexCSS), want) {
			t.Errorf("index.css missing import %q:\n%s", want, indexCSS)
		}
	}
	themeCSS, _ := os.ReadFile(filepath.Join(dir, "web/gsxui/theme.css"))
	if !strings.Contains(string(themeCSS), "--primary") {
		t.Error("theme.css does not look like the token file")
	}
	foundationCSS, _ := os.ReadFile(filepath.Join(dir, "web/gsxui/foundation.css"))
	if !strings.Contains(string(foundationCSS), "@theme inline") {
		t.Error("foundation.css does not contain the Tailwind theme mapping")
	}
	styleCSS, _ := os.ReadFile(filepath.Join(dir, "web/gsxui/style.css"))
	if !strings.Contains(string(styleCSS), "[data-slot=\"scroll-area\"]::-webkit-scrollbar") {
		t.Error("style.css does not contain the ScrollArea pseudo-element rules")
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

func TestInitVendorsCSSAssetsBesideCustomEntry(t *testing.T) {
	dir, _ := initTestModule(t)
	cfg := Config{UI: "ui", JS: "web/gsxui", CSS: "web/styles/brand.css"}
	if err := cfg.Save(dir); err != nil {
		t.Fatal(err)
	}

	if err := Run([]string{"init"}); err != nil {
		t.Fatal(err)
	}

	for _, path := range []string{
		"web/styles/brand.css",
		"web/styles/foundation.css",
		"web/styles/theme.css",
		"web/styles/style.css",
	} {
		if _, err := os.Stat(filepath.Join(dir, path)); err != nil {
			t.Errorf("missing %s: %v", path, err)
		}
	}
}

func TestInitRejectsModifiedCustomCSSSibling(t *testing.T) {
	dir, _ := initTestModule(t)
	cfg := Config{UI: "ui", JS: "web/gsxui", CSS: "web/styles/brand.css"}
	if err := cfg.Save(dir); err != nil {
		t.Fatal(err)
	}
	if err := Run([]string{"init"}); err != nil {
		t.Fatal(err)
	}

	themePath := filepath.Join(dir, "web/styles/theme.css")
	const modified = "/* local theme */\n"
	if err := os.WriteFile(themePath, []byte(modified), 0o644); err != nil {
		t.Fatal(err)
	}
	err := Run([]string{"init"})
	if err == nil || !strings.Contains(err.Error(), themePath) || !strings.Contains(err.Error(), "--overwrite") {
		t.Fatalf("want custom sibling overwrite-refusal error, got %v", err)
	}
	got, err := os.ReadFile(themePath)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != modified {
		t.Errorf("modified sibling was overwritten:\n%s", got)
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
