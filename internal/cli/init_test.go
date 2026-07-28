package cli

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	gsxui "github.com/gsxhq/gsxui"
	parse "github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/css"
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
	if !strings.Contains(string(styleCSS), "[data-gsxui-slot-scroll-area]::-webkit-scrollbar") {
		t.Error("style.css does not contain the ScrollArea pseudo-element rules")
	}
	if _, err := os.Stat(filepath.Join(dir, "web/gsxui/site-button.css")); !os.IsNotExist(err) {
		t.Errorf("init vendored the site-only Button fallback: %v", err)
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

func TestCSSAssetTargetsExcludeSiteButtonFallback(t *testing.T) {
	targets := cssAssetTargets("web/styles/brand.css")
	got := make(map[string]string, len(targets))
	for _, target := range targets {
		got[target.source] = target.target
	}
	want := map[string]string{
		"assets/css/index.css":          "web/styles/brand.css",
		"assets/css/foundation.css":     "web/styles/foundation.css",
		"assets/css/themes/default.css": "web/styles/theme.css",
		"assets/css/styles/default.css": "web/styles/style.css",
	}
	if len(got) != len(want) {
		t.Fatalf("cssAssetTargets() has %d entries, want %d: %v", len(got), len(want), got)
	}
	for source, target := range want {
		if got[source] != target {
			t.Errorf("cssAssetTargets()[%q] = %q, want %q", source, got[source], target)
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

func TestInitOverwriteRefreshesOnlyVersionedSupportFiles(t *testing.T) {
	dir, _ := initTestModule(t)
	cfg := Config{
		UI:  "components/ui",
		JS:  "web/behavior",
		CSS: "web/styles/brand.css",
	}
	if err := cfg.Save(dir); err != nil {
		t.Fatal(err)
	}
	indexPath := filepath.Join(dir, cfg.JS, "index.js")
	if err := os.MkdirAll(filepath.Dir(indexPath), 0o755); err != nil {
		t.Fatal(err)
	}
	const customBarrel = "// consumer-owned behavior barrel\n"
	if err := os.WriteFile(indexPath, []byte(customBarrel), 0o644); err != nil {
		t.Fatal(err)
	}
	const customGsxTOML = "class_merger = \"example.com/custom.Merge\"\n\n[minify]\ncss = true\n"
	if err := os.WriteFile(filepath.Join(dir, "gsx.toml"), []byte(customGsxTOML), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Run([]string{"init"}); err != nil {
		t.Fatal(err)
	}

	type supportTarget struct {
		path   string
		source string
	}
	targets := []supportTarget{
		{path: cfg.CSS, source: "assets/css/index.css"},
		{path: "web/styles/foundation.css", source: "assets/css/foundation.css"},
		{path: "web/styles/theme.css", source: "assets/css/themes/default.css"},
		{path: "web/styles/style.css", source: "assets/css/styles/default.css"},
		{path: filepath.Join(cfg.JS, "gsxui.js"), source: "ui/gsxui.js"},
		{path: filepath.Join(cfg.UI, "merge", "merge.go"), source: "merge/merge.go"},
	}
	for _, target := range targets {
		if err := os.WriteFile(
			filepath.Join(dir, target.path),
			[]byte("locally modified "+target.path+"\n"),
			0o644,
		); err != nil {
			t.Fatal(err)
		}
	}

	preservedPaths := []string{"gsxui.json", cfg.JS + "/index.js", "gsx.toml"}
	preserved := make(map[string][]byte, len(preservedPaths))
	for _, path := range preservedPaths {
		content, err := os.ReadFile(filepath.Join(dir, path))
		if err != nil {
			t.Fatal(err)
		}
		preserved[path] = content
	}

	err := Run([]string{"init"})
	if err == nil || !strings.Contains(err.Error(), "--overwrite") {
		t.Fatalf("plain init must refuse locally modified support files, got %v", err)
	}
	for _, target := range targets {
		got, readErr := os.ReadFile(filepath.Join(dir, target.path))
		if readErr != nil {
			t.Fatal(readErr)
		}
		want := "locally modified " + target.path + "\n"
		if string(got) != want {
			t.Errorf("plain init changed %s:\n%s", target.path, got)
		}
	}

	if err := Run([]string{"init", "--overwrite"}); err != nil {
		t.Fatalf("init --overwrite: %v", err)
	}
	for _, target := range targets {
		want, readErr := fs.ReadFile(gsxui.Files, target.source)
		if readErr != nil {
			t.Fatal(readErr)
		}
		got, readErr := os.ReadFile(filepath.Join(dir, target.path))
		if readErr != nil {
			t.Fatal(readErr)
		}
		if string(got) != string(want) {
			t.Errorf("init --overwrite did not refresh %s from %s", target.path, target.source)
		}
	}
	for _, path := range preservedPaths {
		got, readErr := os.ReadFile(filepath.Join(dir, path))
		if readErr != nil {
			t.Fatal(readErr)
		}
		if string(got) != string(preserved[path]) {
			t.Errorf("init --overwrite changed preserved %s:\n got: %q\nwant: %q", path, got, preserved[path])
		}
	}
}

func TestInitRejectsUnexpectedArgumentsWithOverwriteUsage(t *testing.T) {
	_, _ = initTestModule(t)
	err := Run([]string{"init", "unexpected"})
	if err == nil || !strings.Contains(err.Error(), "usage: gsxui init [--overwrite]") {
		t.Fatalf("unexpected init argument error = %v", err)
	}
}

func TestInitializedConsumerCSSOmitsButtonPresentationAndKeepsRemainingStyle(t *testing.T) {
	repoRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	tailwind := filepath.Join(repoRoot, "node_modules", ".bin", "tailwindcss")
	postcss := filepath.Join(repoRoot, "node_modules", "postcss", "package.json")
	if _, err := exec.LookPath("node"); err != nil {
		t.Skip("requires npm install")
	}
	if _, err := os.Stat(tailwind); err != nil {
		t.Skip("requires npm install")
	}
	if _, err := os.Stat(postcss); err != nil {
		t.Skip("requires npm install")
	}
	dir, _ := initTestModule(t)
	cfg := Config{
		UI:  "components/ui",
		JS:  "web/behavior",
		CSS: "web/styles/brand.css",
	}
	if err := cfg.Save(dir); err != nil {
		t.Fatal(err)
	}
	if err := Run([]string{"init"}); err != nil {
		t.Fatal(err)
	}
	if err := Run([]string{"add", "button-group", "combobox", "command", "menubar"}); err != nil {
		t.Fatal(err)
	}

	nodeModules := filepath.Join(dir, "node_modules")
	if err := os.Symlink(filepath.Join(repoRoot, "node_modules"), nodeModules); err != nil {
		t.Fatal(err)
	}
	input := filepath.Join(dir, "consumer.css")
	if err := os.WriteFile(input, []byte(`
@import "./web/styles/brand.css";
@source "./components/ui/**/*.gsx";
`), 0o644); err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(dir, "consumer.out.css")
	command := exec.Command(tailwind, "-i", input, "-o", output)
	command.Dir = dir
	if combined, runErr := command.CombinedOutput(); runErr != nil {
		t.Fatalf("compile clean consumer CSS: %v\n%s", runErr, combined)
	}

	audit := filepath.Join(repoRoot, "jstest", "support", "compiled-css-audit.ts")
	command = exec.Command("node", audit, output)
	command.Dir = dir
	if combined, runErr := command.CombinedOutput(); runErr != nil {
		t.Fatalf("audit clean consumer CSS: %v\n%s", runErr, combined)
	}

	compiled, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	if len(compiled) == 0 {
		t.Fatal("compiled consumer CSS is empty")
	}
	vendoredStyle, err := os.ReadFile(filepath.Join(dir, "web/styles/style.css"))
	if err != nil {
		t.Fatal(err)
	}
	attributes := parsedSelectorAttributes(t, vendoredStyle)
	if attributes["data-gsxui-slot-button"] {
		t.Error("initialized consumer CSS still contains Button presentation")
	}
	for _, attribute := range []string{
		"data-gsxui-slot-combobox-content",
		"data-gsxui-slot-command",
		"data-gsxui-slot-menubar",
	} {
		if !attributes[attribute] {
			t.Errorf("clean consumer CSS missing canonical selector [%s]", attribute)
		}
	}
}

func parsedSelectorAttributes(t *testing.T, src []byte) map[string]bool {
	t.Helper()
	attributes := make(map[string]bool)
	parser := css.NewParser(parse.NewInputBytes(src), false)
	for {
		grammar, _, _ := parser.Next()
		switch grammar {
		case css.ErrorGrammar:
			if err := parser.Err(); !errors.Is(err, io.EOF) {
				t.Fatalf("parse compiled consumer CSS: %v", err)
			}
			return attributes
		case css.BeginRulesetGrammar:
			tokens := parser.Values()
			for index, token := range tokens {
				if token.TokenType != css.LeftBracketToken {
					continue
				}
				for index++; index < len(tokens); index++ {
					token = tokens[index]
					switch token.TokenType {
					case css.WhitespaceToken, css.CommentToken:
						continue
					case css.IdentToken:
						attributes[strings.ToLower(string(token.Data))] = true
					}
					break
				}
			}
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
