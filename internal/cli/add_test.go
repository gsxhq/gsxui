package cli

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	gsxui "github.com/gsxhq/gsxui"
	"github.com/gsxhq/gsxui/internal/preset"
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
	button, err := os.ReadFile(filepath.Join(dir, "ui/button.gsx"))
	if err != nil {
		t.Fatal(err)
	}
	wantButton, err := fs.ReadFile(gsxui.Files, "registry/generated/nova/button.gsx")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(button, wantButton) {
		t.Fatalf("dependency Button is not exact Nova registry source")
	}
}

func TestAddButtonUsesSelectedStyleExactSource(t *testing.T) {
	tests := []struct {
		style preset.Style
	}{
		{style: preset.StyleNova},
		{style: preset.StyleMaia},
	}
	for _, tt := range tests {
		t.Run(string(tt.style), func(t *testing.T) {
			dir, _ := initedModuleWithStyle(t, tt.style)
			if err := Run([]string{"add", "button"}); err != nil {
				t.Fatal(err)
			}
			got, err := os.ReadFile(filepath.Join(dir, "ui/button.gsx"))
			if err != nil {
				t.Fatal(err)
			}
			want, err := fs.ReadFile(gsxui.Files, "registry/generated/"+string(tt.style)+"/button.gsx")
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(got, want) {
				t.Fatalf("button.gsx differs from exact %s registry artifact", tt.style)
			}
			assertManagedHashMatches(t, dir, "ui/button.gsx")
		})
	}
}

func TestAddNovaNonButtonUsesExistingUISource(t *testing.T) {
	dir, _ := initedModuleWithStyle(t, preset.StyleNova)
	if err := Run([]string{"add", "badge"}); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(dir, "ui/badge.gsx"))
	if err != nil {
		t.Fatal(err)
	}
	want, err := fs.ReadFile(gsxui.Files, "ui/badge.gsx")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatal("Nova non-Button source did not come from ui/")
	}
}

// TestAddVendorsSelectedStyle asserts the correctness-trap fix: every
// component (not only Button) vendors from registry/generated/<style>/, so
// `gsxui add <component>` under a non-default style yields that style's
// exact source rather than DefaultStyle's. It doubles as the maia-gate
// removal test — `gsxui add card` under maia (previously refused with
// "style maia supports only the standalone Button") now succeeds.
func TestAddVendorsSelectedStyle(t *testing.T) {
	tests := []struct {
		style     preset.Style
		component string
	}{
		{style: preset.StyleMaia, component: "card"},
		{style: preset.StyleLyra, component: "card"},
	}
	for _, tt := range tests {
		t.Run(string(tt.style), func(t *testing.T) {
			dir, _ := initedModuleWithStyle(t, tt.style)
			if err := Run([]string{"add", tt.component}); err != nil {
				t.Fatal(err)
			}
			relative := "ui/" + tt.component + ".gsx"
			got, err := os.ReadFile(filepath.Join(dir, relative))
			if err != nil {
				t.Fatal(err)
			}
			want, err := fs.ReadFile(gsxui.Files, "registry/generated/"+string(tt.style)+"/"+tt.component+".gsx")
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(got, want) {
				t.Fatalf("%s differs from exact %s registry artifact", relative, tt.style)
			}
			assertManagedHashMatches(t, dir, relative)
		})
	}
}

func TestAddRequiresValidProjectPresetBeforeResolutionOrWrites(t *testing.T) {
	t.Run("missing", func(t *testing.T) {
		dir, commands := initTestModule(t)
		cfg := DefaultConfig()
		if err := cfg.Save(dir); err != nil {
			t.Fatal(err)
		}
		before := snapshotFiles(t, dir)
		err := Run([]string{"add", "button"})
		if err == nil || !strings.Contains(err.Error(), "gsxui.preset.json") {
			t.Fatalf("missing preset error = %v", err)
		}
		if !equalFileSnapshots(before, snapshotFiles(t, dir)) {
			t.Fatal("missing preset changed project files")
		}
		if len(*commands) != 0 {
			t.Fatalf("missing preset ran commands: %v", *commands)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		dir, commands := initedModule(t)
		if err := os.WriteFile(
			filepath.Join(dir, "gsxui.preset.json"),
			[]byte(`{"style":"maia"}`),
			0o644,
		); err != nil {
			t.Fatal(err)
		}
		before := snapshotFiles(t, dir)
		commandCount := len(*commands)
		err := Run([]string{"add", "button"})
		if err == nil || !strings.Contains(err.Error(), "preset") {
			t.Fatalf("invalid preset error = %v", err)
		}
		if !equalFileSnapshots(before, snapshotFiles(t, dir)) {
			t.Fatal("invalid preset changed project files")
		}
		if len(*commands) != commandCount {
			t.Fatalf("invalid preset ran commands: %v", (*commands)[commandCount:])
		}
	})
}

func TestAddDoesNotVendorRetiredSlotHelpers(t *testing.T) {
	dir, _ := initTestModule(t)
	cfg := DefaultConfig()
	cfg.UI = "components/ui"
	if err := cfg.Save(dir); err != nil {
		t.Fatal(err)
	}
	if err := Run([]string{"init"}); err != nil {
		t.Fatal(err)
	}
	if err := Run([]string{"add", "button"}); err != nil {
		t.Fatal(err)
	}

	for _, path := range []string{
		"components/ui/slots.go",
		"components/ui/internal/slotattr/slotattr.go",
	} {
		if _, err := os.Stat(filepath.Join(dir, path)); !os.IsNotExist(err) {
			t.Errorf("retired support file %s was vendored: %v", path, err)
		}
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

func TestAddGenerateFailureRollsBackAndRegeneratesRestoredSources(t *testing.T) {
	dir, _ := initedModule(t)
	before := treeDigest(t, dir)
	orig := runCommand
	var generateCalls int
	runCommand = func(dir, name string, args ...string) error {
		if name == "go" && len(args) >= 2 && args[0] == "tool" && args[1] == "gsx" {
			generateCalls++
		}
		if generateCalls == 1 {
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
	if generateCalls != 2 {
		t.Fatalf("generation calls = %d, want failed apply plus restored-source regeneration", generateCalls)
	}
	if got := treeDigest(t, dir); got != before {
		t.Fatalf("generation failure did not restore exact tree:\n before %s\n  after %s", before, got)
	}
	assertNoTransactionArtifacts(t, dir)
}

func TestAddHardlinkAliasIsPreservedAcrossRefusalAndGenerationFailure(t *testing.T) {
	dir, commands := initedModule(t)
	target := filepath.Join(dir, "ui", "button.gsx")
	alias := filepath.Join(dir, "user-button.gsx")
	const userContent = "package userbutton // user-owned hardlink\n"
	if err := os.WriteFile(alias, []byte(userContent), 0o640); err != nil {
		t.Fatal(err)
	}
	if err := os.Link(alias, target); err != nil {
		t.Fatal(err)
	}

	commandCount := len(*commands)
	beforeRefusal := treeDigest(t, dir)
	err := Run([]string{"add", "button"})
	if err == nil || !strings.Contains(err.Error(), "--overwrite") {
		t.Fatalf("plain add error = %v, want overwrite refusal", err)
	}
	if got := treeDigest(t, dir); got != beforeRefusal {
		t.Fatalf("plain add changed hardlinked project tree:\n before %s\n  after %s", beforeRefusal, got)
	}
	if len(*commands) != commandCount {
		t.Fatalf("plain add ran commands: %v", (*commands)[commandCount:])
	}

	configBefore, err := os.ReadFile(filepath.Join(dir, "gsxui.json"))
	if err != nil {
		t.Fatal(err)
	}
	previousRunner := runCommand
	runCommand = func(_ string, name string, args ...string) error {
		if name == "go" && len(args) >= 2 && args[0] == "tool" && args[1] == "gsx" {
			return fmt.Errorf("forced generation failure")
		}
		return nil
	}
	t.Cleanup(func() { runCommand = previousRunner })

	err = Run([]string{"add", "--overwrite", "button"})
	if err == nil || !strings.Contains(err.Error(), "gsx generate") {
		t.Fatalf("overwrite add error = %v, want generation failure", err)
	}
	gotAlias, err := os.ReadFile(alias)
	if err != nil {
		t.Fatal(err)
	}
	if string(gotAlias) != userContent {
		t.Fatalf("generation failure changed user hardlink alias:\n%s", gotAlias)
	}
	gotTarget, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(gotTarget) != userContent {
		t.Fatal("generation failure did not restore the prior Button bytes")
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
		t.Fatal("generation failure did not restore the prior Button hardlink identity")
	}
	configAfter, err := os.ReadFile(filepath.Join(dir, "gsxui.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(configAfter, configBefore) {
		t.Fatal("generation failure persisted managed config changes")
	}
	assertNoTransactionArtifacts(t, dir)
}

func TestAddRejectsLeafSymlinkWithAndWithoutOverwrite(t *testing.T) {
	for _, overwrite := range []bool{false, true} {
		name := "plain"
		args := []string{"add", "button"}
		if overwrite {
			name = "overwrite"
			args = []string{"add", "--overwrite", "button"}
		}
		t.Run(name, func(t *testing.T) {
			dir, commands := initedModule(t)
			alias := filepath.Join(dir, "user-button.gsx")
			const userContent = "package userbutton // symlink target\n"
			if err := os.WriteFile(alias, []byte(userContent), 0o640); err != nil {
				t.Fatal(err)
			}
			target := filepath.Join(dir, "ui", "button.gsx")
			if err := os.Symlink(alias, target); err != nil {
				t.Fatal(err)
			}
			before := treeDigest(t, dir)
			commandCount := len(*commands)

			err := Run(args)
			if err == nil || !strings.Contains(err.Error(), "unsafe symlink") {
				t.Fatalf("add error = %v, want leaf-symlink rejection", err)
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
	if _, err := os.Stat(filepath.Join(dir, "ui/badge.gsx")); !os.IsNotExist(err) {
		t.Errorf("barrel conflict was discovered after writing badge.gsx: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "ui/NOTICE.md")); !os.IsNotExist(err) {
		t.Errorf("barrel conflict was discovered after writing NOTICE.md: %v", err)
	}
	if err := Run([]string{"add", "--overwrite", "badge"}); err != nil {
		t.Fatal(err)
	}
	got, _ = os.ReadFile(indexPath)
	if !strings.HasPrefix(string(got), barrelHeader) {
		t.Errorf("--overwrite should replace with the generated barrel:\n%s", got)
	}
}

func TestAddPreflightsEveryTargetBeforeWriting(t *testing.T) {
	dir, _ := initedModule(t)
	if err := Run([]string{"add", "button"}); err != nil {
		t.Fatal(err)
	}
	buttonPath := filepath.Join(dir, "ui/button.gsx")
	const local = "package ui // local Button\n"
	if err := os.WriteFile(buttonPath, []byte(local), 0o644); err != nil {
		t.Fatal(err)
	}

	err := Run([]string{"add", "badge", "button"})
	if err == nil || !strings.Contains(err.Error(), "--overwrite") {
		t.Fatalf("modified Button error = %v", err)
	}
	if _, statErr := os.Stat(filepath.Join(dir, "ui/badge.gsx")); !os.IsNotExist(statErr) {
		t.Fatalf("preflight failure wrote badge.gsx: %v", statErr)
	}
	got, readErr := os.ReadFile(buttonPath)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if string(got) != local {
		t.Fatalf("preflight failure changed Button:\n%s", got)
	}
}

func TestAddDiffIsReadOnlyForAbsentIdenticalAndDifferentTargets(t *testing.T) {
	dir, commands := initedModule(t)

	beforeAbsent := snapshotFiles(t, dir)
	commandCount := len(*commands)
	absentOutput, err := captureRunOutput(t, []string{"add", "--diff", "button"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(absentOutput, "+++ ui/button.gsx") ||
		!strings.Contains(absentOutput, "+package ui") {
		t.Fatalf("absent diff output:\n%s", absentOutput)
	}
	if !equalFileSnapshots(beforeAbsent, snapshotFiles(t, dir)) {
		t.Fatal("absent diff changed project files")
	}
	if len(*commands) != commandCount {
		t.Fatalf("absent diff ran commands: %v", (*commands)[commandCount:])
	}

	if err := Run([]string{"add", "button"}); err != nil {
		t.Fatal(err)
	}
	identicalBefore := snapshotFiles(t, dir)
	commandCount = len(*commands)
	identicalOutput, err := captureRunOutput(t, []string{"add", "--diff", "button"})
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(identicalOutput) != "(no changes)" {
		t.Fatalf("identical diff output = %q", identicalOutput)
	}
	if !equalFileSnapshots(identicalBefore, snapshotFiles(t, dir)) {
		t.Fatal("identical diff changed project files")
	}
	if len(*commands) != commandCount {
		t.Fatalf("identical diff ran commands: %v", (*commands)[commandCount:])
	}

	buttonPath := filepath.Join(dir, "ui/button.gsx")
	if err := os.WriteFile(buttonPath, []byte("package ui // local\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	differentBefore := snapshotFiles(t, dir)
	commandCount = len(*commands)
	differentOutput, err := captureRunOutput(t, []string{"add", "--diff", "button"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(differentOutput, "-package ui // local") ||
		!strings.Contains(differentOutput, "+package ui") {
		t.Fatalf("different diff output:\n%s", differentOutput)
	}
	if !equalFileSnapshots(differentBefore, snapshotFiles(t, dir)) {
		t.Fatal("different diff changed project files")
	}
	if len(*commands) != commandCount {
		t.Fatalf("different diff ran commands: %v", (*commands)[commandCount:])
	}
}

func TestAddDiffAndOverwriteAreMutuallyExclusive(t *testing.T) {
	_, _ = initedModule(t)
	err := Run([]string{"add", "--diff", "--overwrite", "button"})
	if err == nil || !strings.Contains(err.Error(), "mutually exclusive") {
		t.Fatalf("diff+overwrite error = %v", err)
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

func initedModuleWithStyle(t *testing.T, style preset.Style) (string, *[][]string) {
	t.Helper()
	dir, commands := initTestModule(t)
	code, err := preset.EncodeShare(preset.Default(style))
	if err != nil {
		t.Fatal(err)
	}
	if err := Run([]string{"init", "--preset", code}); err != nil {
		t.Fatal(err)
	}
	return dir, commands
}

func assertManagedHashMatches(t *testing.T, dir, relative string) {
	t.Helper()
	cfg, err := LoadConfig(dir)
	if err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(relative)))
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(content)
	want := hex.EncodeToString(sum[:])
	if cfg.Managed[relative] != want {
		t.Fatalf("Managed[%q] = %q, want %q", relative, cfg.Managed[relative], want)
	}
}

func captureRunOutput(t *testing.T, args []string) (string, error) {
	t.Helper()
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	original := os.Stdout
	os.Stdout = writer
	output := make(chan []byte, 1)
	go func() {
		content, _ := io.ReadAll(reader)
		output <- content
	}()
	runErr := Run(args)
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	os.Stdout = original
	content := <-output
	if err := reader.Close(); err != nil {
		t.Fatal(err)
	}
	return string(content), runErr
}

func TestAddAfterNonViteInitPreservesAnimateCSS(t *testing.T) {
	dir, _ := nonViteTestModule(t)
	if err := Run([]string{"init"}); err != nil {
		t.Fatal(err)
	}
	before := readFile(t, dir, "web/gsxui/animate.css")
	if err := Run([]string{"add", "dialog"}); err != nil {
		t.Fatal(err)
	}
	if got := readFile(t, dir, "web/gsxui/animate.css"); got != before {
		t.Fatal("gsxui add must not rewrite animate.css")
	}
	index := readFile(t, dir, "web/gsxui/index.js")
	if !strings.Contains(index, `import "./dialog.js";`) {
		t.Fatalf("barrel missing dialog behavior:\n%s", index)
	}
}
