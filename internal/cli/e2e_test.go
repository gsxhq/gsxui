package cli

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	gsxui "github.com/gsxhq/gsxui"
	"github.com/gsxhq/gsxui/internal/preset"
)

// TestE2E exercises the real flow: pinned gsx scaffold → compact Maia init →
// Button → production CSS/Vite build → style round-trip → Dialog → Go build.
// It needs network and the real gsx/npm toolchains; skipped with -short.
func TestE2E(t *testing.T) {
	if testing.Short() {
		t.Skip("network-dependent e2e; run without -short")
	}
	// This temporary consumer is intentionally outside the test runner's workspace.
	t.Setenv("GOWORK", "off")
	dir := scaffoldGSXProject(t)
	t.Chdir(dir)

	maiaCode := presetCode(t, preset.Default(preset.StyleMaia))
	if !strings.HasPrefix(maiaCode, "gsxui:p1:") {
		t.Fatalf("Maia preset code = %q, want compact transport", maiaCode)
	}
	if err := Run([]string{"init", "--preset", maiaCode}); err != nil {
		t.Fatal(err)
	}
	if err := Run([]string{"add", "button"}); err != nil {
		t.Fatal(err)
	}
	maiaButton, err := fs.ReadFile(gsxui.Files, "registry/generated/maia/button.gsx")
	if err != nil {
		t.Fatal(err)
	}
	if got, err := os.ReadFile(filepath.Join(dir, "ui", "button.gsx")); err != nil {
		t.Fatal(err)
	} else if !bytes.Equal(got, maiaButton) {
		t.Fatal("e2e compact init did not install exact Maia Button")
	}

	assertScaffoldPackageIntegration(t, dir)
	viteBefore := readFile(t, dir, "vite.config.ts")
	mainBefore := readFile(t, dir, "web/main.js")
	if err := Run([]string{"init", "--preset", maiaCode}); err != nil {
		t.Fatal(err)
	}
	if got := readFile(t, dir, "vite.config.ts"); got != viteBefore {
		t.Fatal("e2e init rerun changed vite.config.ts")
	}
	if got := readFile(t, dir, "web/main.js"); got != mainBefore {
		t.Fatal("e2e init rerun changed web/main.js")
	}

	mustRun(t, dir, "npm", "run", "build")
	compiledCSS := builtCSS(t, dir)
	if !strings.Contains(compiledCSS, ".h-9") {
		t.Fatal("production CSS is missing Button's h-9 utility")
	}
	for _, unresolved := range []string{"@apply", "@theme", `@import "tailwindcss"`} {
		if strings.Contains(compiledCSS, unresolved) {
			t.Fatalf("production CSS contains unresolved Tailwind directive %q", unresolved)
		}
	}
	mustRun(t, dir, "go", "build", "./...")
	if err := Run([]string{
		"apply",
		"--preset", presetCode(t, preset.Default(preset.StyleNova)),
		"--yes",
	}); err != nil {
		t.Fatal(err)
	}
	novaButton, err := fs.ReadFile(gsxui.Files, "registry/generated/nova/button.gsx")
	if err != nil {
		t.Fatal(err)
	}
	if got, err := os.ReadFile(filepath.Join(dir, "ui", "button.gsx")); err != nil {
		t.Fatal(err)
	} else if !bytes.Equal(got, novaButton) {
		t.Fatal("e2e apply did not reinstall exact Nova Button")
	}
	mustRun(t, dir, "go", "build", "./...")
	if err := Run([]string{"add", "dialog"}); err != nil {
		t.Fatal(err)
	}
	// generate ran for real: generated files exist
	for _, p := range []string{
		"ui/dialog.x.go",
		"ui/button.x.go",
		"ui/i18n/i18n.go",
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
	const renderTest = `package main_test

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

func TestVendoredDialogFooterCloseComposesButtonMarker(t *testing.T) {
	var output bytes.Buffer
	if err := ui.DialogFooter(true, gsx.Text(""), nil).Render(context.Background(), &output); err != nil {
		t.Fatal(err)
	}
	html := output.String()
	buttonStart := strings.Index(html, "<button")
	if buttonStart == -1 {
		t.Fatalf("DialogFooter close action missing button: %s", html)
	}
	buttonEnd := strings.Index(html[buttonStart:], ">")
	if buttonEnd == -1 {
		t.Fatalf("DialogFooter close action has unterminated button: %s", html)
	}
	button := html[buttonStart : buttonStart+buttonEnd+1]
	for _, marker := range []string{
		"data-gsxui-slot-dialog-footer-close",
		"data-gsxui-slot-button",
	} {
		if !strings.Contains(button, marker) || strings.Contains(button, marker+"=") {
			t.Fatalf("DialogFooter close button missing bare %s: %s", marker, button)
		}
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

func TestE2ECustomUIPath(t *testing.T) {
	if testing.Short() {
		t.Skip("network-dependent e2e; run without -short")
	}
	// This temporary consumer is intentionally outside the test runner's workspace.
	t.Setenv("GOWORK", "off")
	dir := scaffoldGSXProject(t)
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

func scaffoldGSXProject(t *testing.T) string {
	t.Helper()
	repoRoot := activeModuleDir(t, "github.com/gsxhq/gsxui")
	toolPath := filepath.Join(t.TempDir(), "gsx")
	mustRun(t, repoRoot, "go", "build", "-o", toolPath, "github.com/gsxhq/gsx/cmd/gsx")
	parent := t.TempDir()
	mustRun(
		t,
		parent,
		toolPath,
		"init",
		"--yes",
		"--module",
		"example.com/app",
		"app",
	)
	return filepath.Join(parent, "app")
}

func assertScaffoldPackageIntegration(t *testing.T, dir string) {
	t.Helper()
	var manifest struct {
		DevDependencies map[string]string `json:"devDependencies"`
	}
	if err := json.Unmarshal([]byte(readFile(t, dir, "package.json")), &manifest); err != nil {
		t.Fatal(err)
	}
	for _, dependency := range []string{"tailwindcss", "@tailwindcss/vite", "tw-animate-css"} {
		if manifest.DevDependencies[dependency] == "" {
			t.Errorf("package.json devDependencies missing %s", dependency)
		}
		if !strings.Contains(readFile(t, dir, "package-lock.json"), `"`+dependency+`"`) {
			t.Errorf("package-lock.json missing %s", dependency)
		}
	}

	vite := readFile(t, dir, "vite.config.ts")
	if strings.Count(vite, `import tailwindcss from "@tailwindcss/vite"`) != 1 ||
		strings.Count(vite, "tailwindcss()") != 1 {
		t.Fatalf("vite.config.ts integration is not singular:\n%s", vite)
	}
	main := readFile(t, dir, "web/main.js")
	if strings.Count(main, `import "./gsxui/index.js"`) != 1 ||
		strings.Count(main, `import "./gsxui/index.css"`) != 1 {
		t.Fatalf("web/main.js integration is not singular:\n%s", main)
	}
}

func builtCSS(t *testing.T, dir string) string {
	t.Helper()
	paths, err := filepath.Glob(filepath.Join(dir, "dist", "assets", "*.css"))
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) == 0 {
		t.Fatal("production build emitted no CSS assets")
	}
	var combined strings.Builder
	for _, path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		combined.Write(content)
	}
	return combined.String()
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

// TestE2EInternalMessagesTranslate is the reporter's pseudo-localization
// gate from #31, run for real: a scratch consumer registers a marker
// translator for i18n.T in gsx.toml, vendors every swept component, and
// renders them. No swept English may appear outside the markers, and every
// marked string must be one of the swept messages. Renderers are compiled
// in by gsx generate, which is why this cannot be a unit test in ui/.
func TestE2EInternalMessagesTranslate(t *testing.T) {
	if testing.Short() {
		t.Skip("network-dependent e2e; run without -short")
	}
	t.Setenv("GOWORK", "off")
	dir := scaffoldGSXProject(t)
	t.Chdir(dir)
	if err := Run([]string{"init"}); err != nil {
		t.Fatal(err)
	}

	// The renderer target must exist and be registered before `add` runs
	// gsx generate: the generated component files call it by import.
	// ui/i18n/translate.go is the consumer-owned file the docs promise
	// survives vendoring — `add` writes only ui/i18n/i18n.go beside it.
	if err := os.MkdirAll(filepath.Join(dir, "ui", "i18n"), 0o755); err != nil {
		t.Fatal(err)
	}
	const translator = `package i18n

import "context"

// Translate wraps every message in markers so a test can tell translated
// text from text that bypassed T.
func Translate(_ context.Context, m T) string {
	return "«" + string(m) + "»"
}
`
	if err := os.WriteFile(filepath.Join(dir, "ui", "i18n", "translate.go"), []byte(translator), 0o644); err != nil {
		t.Fatal(err)
	}
	tomlPath := filepath.Join(dir, "gsx.toml")
	existing, err := os.ReadFile(tomlPath)
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	registration := string(existing) + "\n[renderers]\n\"example.com/app/ui/i18n.T\" = \"example.com/app/ui/i18n.Translate\"\n"
	if err := os.WriteFile(tomlPath, []byte(registration), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Run([]string{"add", "dialog", "sheet", "toast", "toaster", "carousel", "sidebar", "breadcrumb", "pagination", "calendar"}); err != nil {
		t.Fatal(err)
	}

	const renderTest = `package main_test

import (
	"bytes"
	"context"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/gsxhq/gsx"
	"example.com/app/ui"
)

var swept = []string{
	"Close", "Notifications", "Previous slide", "Next slide", "Sidebar",
	"Displays the mobile sidebar.", "Toggle Sidebar", "More", "Previous",
	"Next", "More pages", "Previous month", "Next month", "Month", "Year",
}

func renderAll(t *testing.T) string {
	t.Helper()
	nodes := []gsx.Node{
		ui.DialogContent(false, gsx.Text(""), nil),
		ui.DialogFooter(true, gsx.Text(""), nil),
		ui.SheetContent("right", false, gsx.Text(""), nil),
		ui.Toast("default", "t", "d", "a", "c", nil),
		ui.Toaster(nil),
		ui.CarouselPrevious("horizontal", nil),
		ui.CarouselNext("horizontal", nil),
		ui.Sidebar(true, "left", "sidebar", "offcanvas", gsx.Text(""), nil),
		ui.SidebarTrigger(nil),
		ui.BreadcrumbEllipsis(nil),
		ui.PaginationPrevious("#", nil),
		ui.PaginationNext("#", nil),
		ui.PaginationEllipsis(nil),
		ui.Calendar("single", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), nil, time.Time{}, time.Time{},
			time.Sunday, true, "dropdown", 0, 0, time.Time{}, time.Time{}, nil, nil, "", ui.CalendarLocale{}, nil),
	}
	var out bytes.Buffer
	for _, n := range nodes {
		if err := n.Render(context.Background(), &out); err != nil {
			t.Fatal(err)
		}
		out.WriteString("\n")
	}
	return out.String()
}

func TestNoSweptEnglishOutsideMarkers(t *testing.T) {
	html := renderAll(t)
	marked := regexp.MustCompile("«([^»]*)»")
	seen := map[string]bool{}
	for _, m := range marked.FindAllStringSubmatch(html, -1) {
		seen[m[1]] = true
	}
	stripped := marked.ReplaceAllString(html, "")
	for _, s := range swept {
		if !seen[s] {
			t.Errorf("%q was never rendered through i18n.T", s)
		}
		if strings.Contains(stripped, ">"+s+"<") || strings.Contains(stripped, "\""+s+"\"") {
			t.Errorf("%q still appears outside the translator's markers", s)
		}
	}
	for s := range seen {
		known := false
		for _, w := range swept {
			if w == s {
				known = true
			}
		}
		if !known {
			t.Errorf("translator saw %q, which is not in the swept table — add it to ui/i18n_test.go", s)
		}
	}
}
`
	if err := os.WriteFile(filepath.Join(dir, "messages_test.go"), []byte(renderTest), 0o644); err != nil {
		t.Fatal(err)
	}
	mustRun(t, dir, "go", "test", "./...")
}

func activeModuleDir(t *testing.T, module string) string {
	t.Helper()
	cmd := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", module)
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("locate %s: %v", module, err)
	}
	dir := strings.TrimSpace(string(output))
	if dir == "" {
		t.Fatalf("locate %s: empty module directory", module)
	}
	return dir
}
