package cli

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	gsxui "github.com/gsxhq/gsxui"
	"github.com/gsxhq/gsxui/internal/preset"
	parse "github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/css"
)

// initTestModule creates an exact supported gsx npm/Vite scaffold and stubs
// the command runner.
func initTestModule(t *testing.T) (dir string, commands *[][]string) {
	t.Helper()
	dir = writeScaffoldFixture(t, testGSXVitePristine, testGSXMainPristine)
	writeFile(t, dir, "go.mod", "module example.com/app\n\ngo 1.26\n")
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
		"gsxui.preset.json",
		"web/gsxui/index.css",
		"web/gsxui/foundation.css",
		"web/gsxui/theme.css",
		"web/gsxui/style.css",
		"web/gsxui/gsxui.js",
		"web/gsxui/index.js",
		"ui/merge/merge.go",
		"gsx.toml",
		"vite.config.ts",
		"web/main.js",
	} {
		if _, err := os.Stat(filepath.Join(dir, p)); err != nil {
			t.Errorf("missing %s: %v", p, err)
		}
	}
	assertInitializedPreset(t, dir, preset.StyleNova)
	assertManagedFilesMatch(t, dir, []string{
		"gsxui.preset.json",
		"ui/merge/merge.go",
		"web/gsxui/foundation.css",
		"web/gsxui/gsxui.js",
		"web/gsxui/index.css",
		"web/gsxui/index.js",
		"web/gsxui/style.css",
		"web/gsxui/theme.css",
	})
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
	// The consumer receives ONE flat sheet even though the default style is
	// authored one file per component. This is the end-to-end half of that
	// contract: whatever composeStyleCSS produces is exactly what lands on disk,
	// byte for byte. Without it, a component file the aggregator does not import
	// would ship a consumer a stylesheet missing that component and every other
	// assertion here would still pass.
	composed, err := composeStyleCSS(gsxui.Files, defaultStyleCSS)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(styleCSS, composed) {
		t.Errorf(
			"vendored style.css is not the composed default style (%d bytes on disk, %d composed)",
			len(styleCSS), len(composed),
		)
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
		"npm install --save-dev tailwindcss@^4.3.3 @tailwindcss/vite@^4.3.3 tw-animate-css@^1.4.0",
		"go get github.com/gsxhq/gsx@latest",
		"go get github.com/jackielii/tailwind-merge-go@latest",
		"go get -tool github.com/gsxhq/gsx/cmd/gsx@latest",
		"go tool gsx generate",
	} {
		if !strings.Contains(joined, want) {
			t.Errorf("missing command %q in:\n%s", want, joined)
		}
	}
}

func TestInitPresetInputsSelectMaia(t *testing.T) {
	maiaJSON, err := preset.CanonicalJSON(preset.Default(preset.StyleMaia))
	if err != nil {
		t.Fatal(err)
	}
	maiaCode, err := preset.EncodeShare(preset.Default(preset.StyleMaia))
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name      string
		argument  func(*testing.T, string) string
		stdin     []byte
		wantInput bool
	}{
		{
			name: "file",
			argument: func(t *testing.T, dir string) string {
				path := filepath.Join(dir, "maia.json")
				if err := os.WriteFile(path, maiaJSON, 0o644); err != nil {
					t.Fatal(err)
				}
				return path
			},
		},
		{
			name: "stdin",
			argument: func(*testing.T, string) string {
				return "-"
			},
			stdin:     append(append([]byte(nil), maiaJSON...), '\n'),
			wantInput: true,
		},
		{
			name: "share code",
			argument: func(*testing.T, string) string {
				return maiaCode
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, _ := initTestModule(t)
			if tt.wantInput {
				original := commandStdin
				commandStdin = bytes.NewReader(tt.stdin)
				t.Cleanup(func() { commandStdin = original })
			}
			argument := tt.argument(t, dir)
			if err := Run([]string{"init", "--preset", argument}); err != nil {
				t.Fatal(err)
			}
			assertInitializedPreset(t, dir, preset.StyleMaia)
		})
	}
}

func TestInitMalformedPresetDoesNotWriteAnything(t *testing.T) {
	dir, commands := initTestModule(t)
	before := snapshotFiles(t, dir)
	err := Run([]string{"init", "--preset", `{"style":"maia"}`})
	if err == nil || !strings.Contains(err.Error(), "preset") {
		t.Fatalf("init malformed preset error = %v", err)
	}
	after := snapshotFiles(t, dir)
	if !equalFileSnapshots(before, after) {
		t.Fatal("malformed preset changed the scaffold")
	}
	if len(*commands) != 0 {
		t.Fatalf("malformed preset ran commands: %v", *commands)
	}
}

func TestInitDifferentPresetRequiresOverwriteBeforeAnyWrite(t *testing.T) {
	dir, _ := initTestModule(t)
	if err := Run([]string{"init"}); err != nil {
		t.Fatal(err)
	}
	before := snapshotFiles(t, dir)
	maiaCode, err := preset.EncodeShare(preset.Default(preset.StyleMaia))
	if err != nil {
		t.Fatal(err)
	}

	err = Run([]string{"init", "--preset", maiaCode})
	if err == nil || !strings.Contains(err.Error(), "--overwrite") ||
		!strings.Contains(err.Error(), "gsxui.preset.json") {
		t.Fatalf("different preset error = %v", err)
	}
	afterRefusal := snapshotFiles(t, dir)
	if !equalFileSnapshots(before, afterRefusal) {
		t.Fatal("different preset refusal changed project files")
	}

	if err := Run([]string{"init", "--preset", maiaCode, "--overwrite"}); err != nil {
		t.Fatal(err)
	}
	assertInitializedPreset(t, dir, preset.StyleMaia)
}

func TestInitRejectsModifiedGeneratedBarrelWithoutAdoptingItsHash(t *testing.T) {
	dir, _ := initTestModule(t)
	if err := Run([]string{"init"}); err != nil {
		t.Fatal(err)
	}
	cfgBefore, err := LoadConfig(dir)
	if err != nil {
		t.Fatal(err)
	}
	indexPath := filepath.Join(dir, "web/gsxui/index.js")
	modified := []byte(barrelHeader + "\n// locally modified but still generated-looking\n")
	if err := os.WriteFile(indexPath, modified, 0o644); err != nil {
		t.Fatal(err)
	}

	err = Run([]string{"init"})
	if err == nil || !strings.Contains(err.Error(), "--overwrite") {
		t.Fatalf("modified generated barrel error = %v", err)
	}
	got, readErr := os.ReadFile(indexPath)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if !bytes.Equal(got, modified) {
		t.Fatalf("plain init changed modified barrel:\n%s", got)
	}
	cfgAfter, loadErr := LoadConfig(dir)
	if loadErr != nil {
		t.Fatal(loadErr)
	}
	if cfgAfter.Managed["web/gsxui/index.js"] != cfgBefore.Managed["web/gsxui/index.js"] {
		t.Fatal("plain init adopted the modified barrel hash")
	}

	if err := Run([]string{"init", "--overwrite"}); err != nil {
		t.Fatal(err)
	}
	got, readErr = os.ReadFile(indexPath)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if !bytes.Equal(got, Barrel(nil)) {
		t.Fatalf("init --overwrite did not restore generated barrel:\n%s", got)
	}
	assertManagedFilesMatch(t, dir, []string{
		"gsxui.preset.json",
		"ui/merge/merge.go",
		"web/gsxui/foundation.css",
		"web/gsxui/gsxui.js",
		"web/gsxui/index.css",
		"web/gsxui/index.js",
		"web/gsxui/style.css",
		"web/gsxui/theme.css",
	})
}

func TestInitRejectsConfigAliasedToPlannedArtifactBeforeWrites(t *testing.T) {
	for _, alias := range []string{"symlink", "hardlink"} {
		t.Run(alias, func(t *testing.T) {
			dir, commands := initTestModule(t)
			themePath := filepath.Join(dir, "web", "gsxui", "theme.css")
			if err := os.MkdirAll(filepath.Dir(themePath), 0o755); err != nil {
				t.Fatal(err)
			}
			legacyConfig := []byte(`{"ui":"ui","js":"web/gsxui","css":"web/gsxui/index.css"}`)
			if err := os.WriteFile(themePath, legacyConfig, 0o644); err != nil {
				t.Fatal(err)
			}
			configPath := filepath.Join(dir, "gsxui.json")
			var err error
			if alias == "symlink" {
				err = os.Symlink(filepath.ToSlash(filepath.Join("web", "gsxui", "theme.css")), configPath)
			} else {
				err = os.Link(themePath, configPath)
			}
			if err != nil {
				t.Fatal(err)
			}
			before := snapshotFiles(t, dir)

			err = Run([]string{"init", "--overwrite"})
			wantError := "same destination"
			if alias == "symlink" {
				wantError = "unsafe symlink"
			}
			if err == nil || !strings.Contains(err.Error(), wantError) {
				t.Fatalf("init error = %v, want config/artifact alias rejection", err)
			}
			if len(*commands) != 0 {
				t.Fatalf("aliased config ran commands: %v", *commands)
			}
			after := snapshotFiles(t, dir)
			if !equalFileSnapshots(before, after) {
				t.Fatal("aliased config rejection changed project files")
			}
		})
	}
}

func TestInitRejectsPlannedFileAncestorBeforeAnyWrite(t *testing.T) {
	dir, commands := initTestModule(t)
	cfg := Config{
		UI:  "generated/merge.go",
		JS:  "web/gsxui",
		CSS: "generated/merge.go",
	}
	if err := cfg.Save(dir); err != nil {
		t.Fatal(err)
	}
	before := treeDigest(t, dir)

	err := Run([]string{"init", "--overwrite"})
	if err == nil || !strings.Contains(err.Error(), "default gsxui JS and CSS paths") {
		t.Errorf("init error = %v, want unsupported-path preflight rejection", err)
	}
	if len(*commands) != 0 {
		t.Errorf("ancestor conflict ran commands: %v", *commands)
	}
	after := treeDigest(t, dir)
	if after != before {
		t.Fatalf("ancestor conflict changed complete project tree:\n before %s\n  after %s", before, after)
	}
}

func TestInitRejectsPortableCaseAliasBeforeAnyWrite(t *testing.T) {
	dir, commands := initTestModule(t)
	cfg := Config{
		UI:  "ui",
		JS:  "Web/GSXUI",
		CSS: "web/gsxui/index.js",
	}
	if err := cfg.Save(dir); err != nil {
		t.Fatal(err)
	}
	before := treeDigest(t, dir)

	err := Run([]string{"init", "--overwrite"})
	if err == nil || !strings.Contains(err.Error(), "default gsxui JS and CSS paths") {
		t.Errorf("init error = %v, want unsupported-path preflight rejection", err)
	}
	if len(*commands) != 0 {
		t.Errorf("portable-alias conflict ran commands: %v", *commands)
	}
	after := treeDigest(t, dir)
	if after != before {
		t.Fatalf("portable-alias conflict changed complete project tree:\n before %s\n  after %s", before, after)
	}
}

func TestInitHardlinkAliasIsPreservedAcrossRefusalAndCommandFailure(t *testing.T) {
	dir, commands := initedModule(t)
	target := filepath.Join(dir, "web", "gsxui", "theme.css")
	alias := filepath.Join(dir, "user-theme.css")
	if err := os.Link(target, alias); err != nil {
		t.Fatal(err)
	}
	const userContent = "/* user-owned shared theme */\n"
	if err := os.WriteFile(alias, []byte(userContent), 0o640); err != nil {
		t.Fatal(err)
	}

	commandCount := len(*commands)
	beforeRefusal := treeDigest(t, dir)
	err := Run([]string{"init"})
	if err == nil || !strings.Contains(err.Error(), "--overwrite") {
		t.Fatalf("plain init error = %v, want overwrite refusal", err)
	}
	if got := treeDigest(t, dir); got != beforeRefusal {
		t.Fatalf("plain init changed hardlinked project tree:\n before %s\n  after %s", beforeRefusal, got)
	}
	if len(*commands) != commandCount {
		t.Fatalf("plain init ran commands: %v", (*commands)[commandCount:])
	}

	configBefore, err := os.ReadFile(filepath.Join(dir, "gsxui.json"))
	if err != nil {
		t.Fatal(err)
	}
	beforeFailure := treeDigest(t, dir)
	previousRunner := runCommand
	runCommand = func(_ string, _ string, _ ...string) error {
		return errors.New("forced dependency command failure")
	}
	t.Cleanup(func() { runCommand = previousRunner })

	err = Run([]string{"init", "--overwrite"})
	if err == nil || !strings.Contains(err.Error(), "forced dependency command failure") {
		t.Fatalf("overwrite init error = %v, want command failure", err)
	}
	gotAlias, err := os.ReadFile(alias)
	if err != nil {
		t.Fatal(err)
	}
	if string(gotAlias) != userContent {
		t.Fatalf("command failure changed user hardlink alias:\n%s", gotAlias)
	}
	gotTarget, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(gotTarget) != userContent {
		t.Fatal("command failure did not restore the prior theme bytes")
	}
	targetInfo, err := os.Stat(target)
	if err != nil {
		t.Fatal(err)
	}
	aliasInfo, err := os.Stat(alias)
	if err != nil {
		t.Fatal(err)
	}
	if !os.SameFile(targetInfo, aliasInfo) {
		t.Fatal("command failure did not restore the prior theme hardlink identity")
	}
	configAfter, err := os.ReadFile(filepath.Join(dir, "gsxui.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(configAfter, configBefore) {
		t.Fatal("command failure persisted managed config changes")
	}
	if got := treeDigest(t, dir); got != beforeFailure {
		t.Fatalf("command failure did not restore exact tree:\n before %s\n  after %s", beforeFailure, got)
	}
	assertNoTransactionArtifacts(t, dir)
}

func TestInitGenerationFailureRollsBackAndRegeneratesRestoredSources(t *testing.T) {
	dir, _ := initTestModule(t)
	before := treeDigest(t, dir)
	originalRunner := runCommand
	var generateCalls int
	runCommand = func(_ string, name string, args ...string) error {
		if name == "go" && len(args) >= 2 && args[0] == "tool" && args[1] == "gsx" {
			generateCalls++
			if generateCalls == 1 {
				return errors.New("injected init generation failure")
			}
		}
		return nil
	}
	t.Cleanup(func() { runCommand = originalRunner })

	err := Run([]string{"init"})
	if err == nil || !strings.Contains(err.Error(), "injected init generation failure") {
		t.Fatalf("init generation error = %v", err)
	}
	if generateCalls != 2 {
		t.Fatalf("generation calls = %d, want failed apply plus restored-source regeneration", generateCalls)
	}
	if got := treeDigest(t, dir); got != before {
		t.Fatalf("init generation failure did not restore exact tree:\n before %s\n  after %s", before, got)
	}
	assertNoTransactionArtifacts(t, dir)
}

func TestInitRejectsLeafSymlinkWithAndWithoutOverwrite(t *testing.T) {
	for _, overwrite := range []bool{false, true} {
		name := "plain"
		args := []string{"init"}
		if overwrite {
			name = "overwrite"
			args = []string{"init", "--overwrite"}
		}
		t.Run(name, func(t *testing.T) {
			dir, commands := initedModule(t)
			target := filepath.Join(dir, "web", "gsxui", "theme.css")
			if err := os.Remove(target); err != nil {
				t.Fatal(err)
			}
			alias := filepath.Join(dir, "user-theme.css")
			const userContent = "/* user-owned symlink target */\n"
			if err := os.WriteFile(alias, []byte(userContent), 0o640); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(alias, target); err != nil {
				t.Fatal(err)
			}
			before := treeDigest(t, dir)
			commandCount := len(*commands)

			err := Run(args)
			if err == nil || !strings.Contains(err.Error(), "unsafe symlink") {
				t.Fatalf("init error = %v, want leaf-symlink rejection", err)
			}
			if got := treeDigest(t, dir); got != before {
				t.Fatalf("leaf-symlink rejection changed project tree:\n before %s\n  after %s", before, got)
			}
			if len(*commands) != commandCount {
				t.Fatalf("leaf-symlink rejection ran commands: %v", (*commands)[commandCount:])
			}
		})
	}
}

func assertInitializedPreset(t *testing.T, dir string, style preset.Style) {
	t.Helper()
	wantPreset, err := preset.CanonicalJSON(preset.Default(style))
	if err != nil {
		t.Fatal(err)
	}
	gotPreset, err := os.ReadFile(filepath.Join(dir, "gsxui.preset.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(gotPreset, wantPreset) {
		t.Fatalf("gsxui.preset.json:\n%s\nwant:\n%s", gotPreset, wantPreset)
	}
	wantTheme, err := preset.ThemeCSS(preset.Default(style))
	if err != nil {
		t.Fatal(err)
	}
	gotTheme, err := os.ReadFile(filepath.Join(dir, "web/gsxui/theme.css"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(gotTheme, wantTheme) {
		t.Fatalf("theme.css:\n%s\nwant:\n%s", gotTheme, wantTheme)
	}
}

func assertManagedFilesMatch(t *testing.T, dir string, wantPaths []string) {
	t.Helper()
	cfg, err := LoadConfig(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Managed) != len(wantPaths) {
		t.Fatalf("Managed has %d paths, want %d: %v", len(cfg.Managed), len(wantPaths), cfg.Managed)
	}
	for _, relative := range wantPaths {
		content, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(relative)))
		if err != nil {
			t.Fatalf("read managed %s: %v", relative, err)
		}
		sum := sha256.Sum256(content)
		wantHash := hex.EncodeToString(sum[:])
		if cfg.Managed[relative] != wantHash {
			t.Errorf("Managed[%q] = %q, want %q", relative, cfg.Managed[relative], wantHash)
		}
	}
}

type fileSnapshot map[string][]byte

func snapshotFiles(t *testing.T, root string) fileSnapshot {
	t.Helper()
	snapshot := make(fileSnapshot)
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		snapshot[filepath.ToSlash(relative)] = content
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return snapshot
}

func equalFileSnapshots(left, right fileSnapshot) bool {
	if len(left) != len(right) {
		return false
	}
	for path, content := range left {
		if !bytes.Equal(content, right[path]) {
			return false
		}
	}
	return true
}

func treeDigest(t *testing.T, root string) string {
	t.Helper()
	hash := sha256.New()
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if _, err := io.WriteString(hash, filepath.ToSlash(relative)+"\x00"+info.Mode().String()+"\x00"); err != nil {
			return err
		}
		switch {
		case info.Mode()&os.ModeSymlink != 0:
			target, err := os.Readlink(path)
			if err != nil {
				return err
			}
			_, err = io.WriteString(hash, target)
			return err
		case info.Mode().IsRegular():
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			_, err = hash.Write(content)
			return err
		default:
			return nil
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func TestInitVendorsCSSAssetsBesideCustomEntry(t *testing.T) {
	dir, _ := initTestModule(t)
	cfg := Config{UI: "ui", JS: "web/gsxui", CSS: "web/styles/brand.css"}
	artifacts, err := initArtifacts(dir, "example.com/app", cfg, preset.Default(preset.StyleNova))
	if err != nil {
		t.Fatal(err)
	}
	byPath := make(map[string]artifact, len(artifacts))
	for _, planned := range artifacts {
		byPath[planned.RelativePath] = planned
	}
	for _, relative := range []string{
		"web/styles/brand.css",
		"web/styles/foundation.css",
		"web/styles/theme.css",
		"web/styles/style.css",
	} {
		if _, ok := byPath[relative]; !ok {
			t.Errorf("initArtifacts missing %s", relative)
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

func TestInitRejectsCustomCSSPathBeforeCommands(t *testing.T) {
	dir, commands := initTestModule(t)
	cfg := Config{UI: "ui", JS: "web/gsxui", CSS: "web/styles/brand.css"}
	if err := cfg.Save(dir); err != nil {
		t.Fatal(err)
	}
	err := Run([]string{"init"})
	if err == nil || !strings.Contains(err.Error(), "default gsxui JS and CSS paths") {
		t.Fatalf("want custom path preflight refusal, got %v", err)
	}
	if len(*commands) != 0 {
		t.Fatalf("custom path refusal ran commands: %v", *commands)
	}
}

func TestInitOverwriteRefreshesOnlyVersionedSupportFiles(t *testing.T) {
	dir, _ := initTestModule(t)
	cfg := Config{
		UI:  "components/ui",
		JS:  "web/gsxui",
		CSS: "web/gsxui/index.css",
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
		{path: "web/gsxui/foundation.css", source: "assets/css/foundation.css"},
		{path: "web/gsxui/theme.css", source: "assets/css/themes/default.css"},
		{path: "web/gsxui/style.css", source: "assets/css/styles/default.css"},
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
		// The default style is authored one file per component, so what the
		// consumer receives is the composition of those files, not the aggregator.
		want, readErr := expectedVendoredCSS(target.source)
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
	if err == nil || !strings.Contains(err.Error(), "usage: gsxui init [--preset <file|code|->] [--overwrite]") {
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
		JS:  "web/gsxui",
		CSS: "web/gsxui/index.css",
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
@import "./web/gsxui/index.css";
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
	vendoredStyle, err := os.ReadFile(filepath.Join(dir, "web/gsxui/style.css"))
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

// expectedVendoredCSS is the exact bytes `gsxui init` writes for one embedded CSS
// source: a copy for every sheet except the default style, which is composed
// from assets/css/styles/default/*.css into the single flat file consumers have
// always received. See composeStyleCSS.
func expectedVendoredCSS(source string) ([]byte, error) {
	if source == defaultStyleCSS {
		return composeStyleCSS(gsxui.Files, source)
	}
	return fs.ReadFile(gsxui.Files, source)
}
