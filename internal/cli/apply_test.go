package cli

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	gsxui "github.com/gsxhq/gsxui"
	"github.com/gsxhq/gsxui/internal/preset"
)

func TestApplyNoOpDoesNotPromptWriteOrGenerate(t *testing.T) {
	dir, commands := initedModuleWithStyle(t, preset.StyleNova)
	before := treeDigest(t, dir)
	commandCount := len(*commands)
	code := presetCode(t, preset.Default(preset.StyleNova))
	originalStdin := commandStdin
	commandStdin = strings.NewReader("yes\n")
	t.Cleanup(func() { commandStdin = originalStdin })

	output, err := captureRunOutput(t, []string{"apply", "--preset", code})
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(output) != "(no changes)" {
		t.Fatalf("no-op output = %q", output)
	}
	if got := treeDigest(t, dir); got != before {
		t.Fatalf("no-op apply changed tree:\n before %s\n  after %s", before, got)
	}
	if len(*commands) != commandCount {
		t.Fatalf("no-op apply ran commands: %v", (*commands)[commandCount:])
	}
}

func TestApplyFullNovaToMaiaReplacesOnlyInstalledButtonSource(t *testing.T) {
	dir, _ := initedModuleWithStyle(t, preset.StyleNova)
	if err := Run([]string{"add", "button"}); err != nil {
		t.Fatal(err)
	}
	themeBefore := readProjectFile(t, dir, "web/gsxui/theme.css")

	output, err := captureRunOutput(t, []string{
		"apply",
		"--preset", presetCode(t, preset.Default(preset.StyleMaia)),
		"--yes",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output, "axes: style") ||
		!strings.Contains(output, "ui/button.gsx (component source replacement)") {
		t.Fatalf("style apply summary:\n%s", output)
	}
	assertProjectPreset(t, dir, preset.Default(preset.StyleMaia))
	button := readProjectFile(t, dir, "ui/button.gsx")
	wantButton, err := fs.ReadFile(gsxui.Files, "registry/generated/maia/button.gsx")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(button, wantButton) {
		t.Fatal("full style apply did not install exact Maia Button source")
	}
	if themeAfter := readProjectFile(t, dir, "web/gsxui/theme.css"); !bytes.Equal(themeAfter, themeBefore) {
		t.Fatal("style-only difference rewrote unchanged theme bytes")
	}
	assertManagedHashMatches(t, dir, "ui/button.gsx")
	assertNoTransactionArtifacts(t, dir)
}

func TestApplyFullThemeChangeWritesCanonicalPresetAndTheme(t *testing.T) {
	dir, _ := initedModuleWithStyle(t, preset.StyleNova)
	incoming := preset.Default(preset.StyleNova)
	incoming.Radius = "0.75rem"
	incoming.Theme.Light["primary"] = "oklch(0.5 0.2 120)"
	incoming.Theme.Dark["background"] = "oklch(0.2 0.01 240)"

	if err := Run([]string{"apply", "--preset", presetCode(t, incoming), "--yes"}); err != nil {
		t.Fatal(err)
	}
	assertProjectPreset(t, dir, incoming)
	wantTheme, err := preset.ThemeCSS(incoming)
	if err != nil {
		t.Fatal(err)
	}
	if got := readProjectFile(t, dir, "web/gsxui/theme.css"); !bytes.Equal(got, wantTheme) {
		t.Fatalf("theme.css:\n%s\nwant:\n%s", got, wantTheme)
	}
	assertManagedHashMatches(t, dir, "gsxui.preset.json")
	assertManagedHashMatches(t, dir, "web/gsxui/theme.css")
}

func TestApplyOnlyThemePreservesStyleAndButtonBytes(t *testing.T) {
	dir, _ := initedModuleWithStyle(t, preset.StyleNova)
	if err := Run([]string{"add", "button"}); err != nil {
		t.Fatal(err)
	}
	buttonBefore := readProjectFile(t, dir, "ui/button.gsx")
	incoming := preset.Default(preset.StyleMaia)
	incoming.Radius = "1rem"
	incoming.Theme.Light["primary"] = "oklch(0.55 0.2 120)"

	output, err := captureRunOutput(t, []string{
		"apply",
		"--preset", presetCode(t, incoming),
		"--only", "theme",
		"--yes",
	})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(output, "axes: style") || strings.Contains(output, "button.gsx") {
		t.Fatalf("theme-only summary included style/Button:\n%s", output)
	}
	want := incoming
	want.Style = preset.StyleNova
	assertProjectPreset(t, dir, want)
	if got := readProjectFile(t, dir, "ui/button.gsx"); !bytes.Equal(got, buttonBefore) {
		t.Fatal("--only theme changed Button bytes")
	}
}

func TestApplyOnlyStylePreservesThemeAndRadius(t *testing.T) {
	current := preset.Default(preset.StyleNova)
	current.Radius = "0.25rem"
	current.Theme.Light["primary"] = "oklch(0.45 0.18 255)"
	dir, _ := initedModuleWithPreset(t, current)
	if err := Run([]string{"add", "button"}); err != nil {
		t.Fatal(err)
	}
	themeBefore := readProjectFile(t, dir, "web/gsxui/theme.css")

	if err := Run([]string{
		"apply",
		"--preset", presetCode(t, preset.Default(preset.StyleMaia)),
		"--only", "style",
		"--yes",
	}); err != nil {
		t.Fatal(err)
	}
	want := current
	want.Style = preset.StyleMaia
	assertProjectPreset(t, dir, want)
	if got := readProjectFile(t, dir, "web/gsxui/theme.css"); !bytes.Equal(got, themeBefore) {
		t.Fatal("--only style changed theme bytes")
	}
}

func TestApplyRefusesModifiedManagedTargetsAndOverwriteReplacesThem(t *testing.T) {
	tests := []struct {
		name     string
		relative string
	}{
		{name: "Button", relative: "ui/button.gsx"},
		{name: "theme", relative: "web/gsxui/theme.css"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, commands := initedModuleWithStyle(t, preset.StyleNova)
			if err := Run([]string{"add", "button"}); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(
				filepath.Join(dir, filepath.FromSlash(tt.relative)),
				[]byte("locally modified\n"),
				0o640,
			); err != nil {
				t.Fatal(err)
			}
			incoming := preset.Default(preset.StyleMaia)
			incoming.Theme.Light["primary"] = "oklch(0.5 0.2 120)"
			before := treeDigest(t, dir)
			commandCount := len(*commands)

			err := Run([]string{"apply", "--preset", presetCode(t, incoming), "--yes"})
			if err == nil || !strings.Contains(err.Error(), "--overwrite") ||
				!strings.Contains(err.Error(), tt.relative) {
				t.Fatalf("modified target error = %v", err)
			}
			if got := treeDigest(t, dir); got != before {
				t.Fatalf("modified target refusal changed tree:\n before %s\n  after %s", before, got)
			}
			if len(*commands) != commandCount {
				t.Fatalf("modified target refusal ran commands: %v", (*commands)[commandCount:])
			}

			output, err := captureRunOutput(t, []string{
				"apply",
				"--preset", presetCode(t, incoming),
				"--yes",
				"--overwrite",
			})
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(output, "modified managed conflicts: "+tt.relative) {
				t.Fatalf("overwrite summary omitted modified target:\n%s", output)
			}
			assertProjectPreset(t, dir, incoming)
			assertNoTransactionArtifacts(t, dir)
		})
	}
}

func TestApplyStyleChangeRequiresOverwriteForUnmanagedButton(t *testing.T) {
	dir, commands := initedModuleWithStyle(t, preset.StyleNova)
	buttonPath := filepath.Join(dir, "ui", "button.gsx")
	button, err := fs.ReadFile(gsxui.Files, "registry/generated/nova/button.gsx")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(buttonPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(buttonPath, button, 0o644); err != nil {
		t.Fatal(err)
	}
	before := treeDigest(t, dir)
	commandCount := len(*commands)

	err = Run([]string{
		"apply",
		"--preset", presetCode(t, preset.Default(preset.StyleMaia)),
		"--yes",
	})
	if err == nil || !strings.Contains(err.Error(), "unmanaged") ||
		!strings.Contains(err.Error(), "--overwrite") {
		t.Fatalf("unmanaged Button error = %v", err)
	}
	if got := treeDigest(t, dir); got != before {
		t.Fatalf("unmanaged Button refusal changed tree:\n before %s\n  after %s", before, got)
	}
	if len(*commands) != commandCount {
		t.Fatalf("unmanaged Button refusal ran commands: %v", (*commands)[commandCount:])
	}

	if err := Run([]string{
		"apply",
		"--preset", presetCode(t, preset.Default(preset.StyleMaia)),
		"--yes",
		"--overwrite",
	}); err != nil {
		t.Fatal(err)
	}
	assertProjectPreset(t, dir, preset.Default(preset.StyleMaia))
}

func TestApplyRejectsInvalidConfigBeforePresetResolutionOrWrites(t *testing.T) {
	dir, commands := initedModuleWithStyle(t, preset.StyleNova)
	configPath := filepath.Join(dir, "gsxui.json")
	if err := os.WriteFile(configPath, []byte(`{"ui":"../escape"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	before := treeDigest(t, dir)
	commandCount := len(*commands)

	err := Run([]string{
		"apply",
		"--preset", presetCode(t, preset.Default(preset.StyleMaia)),
		"--yes",
	})
	if err == nil || !strings.Contains(err.Error(), "gsxui.json") {
		t.Fatalf("invalid config error = %v", err)
	}
	if got := treeDigest(t, dir); got != before {
		t.Fatalf("invalid config changed tree:\n before %s\n  after %s", before, got)
	}
	if len(*commands) != commandCount {
		t.Fatalf("invalid config ran commands: %v", (*commands)[commandCount:])
	}
}

func TestApplyInteractiveSummaryDeclineAndConfirmation(t *testing.T) {
	dir, commands := initedModuleWithStyle(t, preset.StyleNova)
	if err := Run([]string{"add", "button"}); err != nil {
		t.Fatal(err)
	}
	incoming := preset.Default(preset.StyleMaia)
	incoming.Radius = "0.75rem"
	incoming.Theme.Light["primary"] = "oklch(0.5 0.2 120)"
	before := treeDigest(t, dir)
	commandCount := len(*commands)
	originalStdin := commandStdin
	commandStdin = strings.NewReader("no\n")

	output, err := captureRunOutput(t, []string{"apply", "--preset", presetCode(t, incoming)})
	commandStdin = originalStdin
	if err != nil {
		t.Fatal(err)
	}
	want := `apply preset:
  axes: style, radius, theme.light
  files:
    gsxui.preset.json
    web/gsxui/theme.css
    ui/button.gsx (component source replacement)
    gsxui.json
  component source is replaced, not merged.
  commit or stash local changes before applying.
continue? [y/N] apply declined.
`
	if output != want {
		t.Fatalf("interactive summary:\n%s\nwant:\n%s", output, want)
	}
	if got := treeDigest(t, dir); got != before {
		t.Fatalf("declined apply changed tree:\n before %s\n  after %s", before, got)
	}
	if len(*commands) != commandCount {
		t.Fatalf("declined apply ran commands: %v", (*commands)[commandCount:])
	}

	commandStdin = strings.NewReader("yes\n")
	t.Cleanup(func() { commandStdin = originalStdin })
	if _, err := captureRunOutput(t, []string{"apply", "--preset", presetCode(t, incoming)}); err != nil {
		t.Fatal(err)
	}
	assertProjectPreset(t, dir, incoming)
}

// TestApplyStyleChangeVendorsEveryInstalledStyledComponent asserts the
// correctness-trap fix on the apply side: a style-axis change re-vendors
// *every* installed styled component from registry/generated/<style>/, not
// only Button. It also doubles as the maia-gate removal test — applying
// maia while a non-Button component (badge) is installed (previously
// refused with "style maia supports only the standalone Button") now
// succeeds and replaces both sources exactly.
func TestApplyStyleChangeVendorsEveryInstalledStyledComponent(t *testing.T) {
	dir, _ := initedModuleWithStyle(t, preset.StyleNova)
	if err := Run([]string{"add", "button"}); err != nil {
		t.Fatal(err)
	}
	if err := Run([]string{"add", "badge"}); err != nil {
		t.Fatal(err)
	}
	if err := Run([]string{"add", "icon"}); err != nil {
		t.Fatal(err)
	}
	iconBefore := readProjectFile(t, dir, "ui/icon/icon.gsx")

	output, err := captureRunOutput(t, []string{
		"apply",
		"--preset", presetCode(t, preset.Default(preset.StyleMaia)),
		"--yes",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output, "ui/button.gsx (component source replacement)") ||
		!strings.Contains(output, "ui/badge.gsx (component source replacement)") {
		t.Fatalf("style apply summary missing component replacements:\n%s", output)
	}
	assertProjectPreset(t, dir, preset.Default(preset.StyleMaia))

	for _, name := range []string{"button", "badge"} {
		got := readProjectFile(t, dir, "ui/"+name+".gsx")
		want, err := fs.ReadFile(gsxui.Files, "registry/generated/maia/"+name+".gsx")
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("full style apply did not install exact Maia %s source", name)
		}
		assertManagedHashMatches(t, dir, "ui/"+name+".gsx")
	}

	// icon has no per-style artifact — it vendors identically under every
	// style, so a style-axis change must leave it untouched.
	if got := readProjectFile(t, dir, "ui/icon/icon.gsx"); !bytes.Equal(got, iconBefore) {
		t.Fatal("style apply rewrote style-invariant icon source")
	}
	assertNoTransactionArtifacts(t, dir)
}

func TestApplyNovaAllowsInstalledNonButtonComponents(t *testing.T) {
	dir, _ := initedModuleWithStyle(t, preset.StyleNova)
	if err := Run([]string{"add", "badge"}); err != nil {
		t.Fatal(err)
	}
	incoming := preset.Default(preset.StyleNova)
	incoming.Radius = "0.75rem"
	if err := Run([]string{"apply", "--preset", presetCode(t, incoming), "--yes"}); err != nil {
		t.Fatal(err)
	}
	assertProjectPreset(t, dir, incoming)
	if _, err := os.Stat(filepath.Join(dir, "ui", "badge.gsx")); err != nil {
		t.Fatalf("Nova apply lost installed Badge: %v", err)
	}
}

func TestApplyMalformedPresetDoesNotWriteOrGenerate(t *testing.T) {
	dir, commands := initedModuleWithStyle(t, preset.StyleNova)
	before := treeDigest(t, dir)
	commandCount := len(*commands)
	err := Run([]string{"apply", "--preset", `{"style":"maia"}`, "--yes"})
	if err == nil || !strings.Contains(err.Error(), "preset") {
		t.Fatalf("malformed preset error = %v", err)
	}
	if got := treeDigest(t, dir); got != before {
		t.Fatalf("malformed preset changed tree:\n before %s\n  after %s", before, got)
	}
	if len(*commands) != commandCount {
		t.Fatalf("malformed preset ran commands: %v", (*commands)[commandCount:])
	}
}

func TestRecoverBeforeMutationPlanning(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "init", args: []string{"init", "--preset", `{"style":"maia"}`}},
		{name: "add", args: []string{"add", "not-a-component"}},
		{name: "apply", args: []string{"apply", "--preset", `{"style":"maia"}`, "--yes"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, _ := initedModuleWithStyle(t, preset.StyleNova)
			before := treeDigest(t, dir)
			interruptTransactionAfterReplacement(t, dir)
			if _, err := os.Stat(filepath.Join(dir, transactionJournalName)); err != nil {
				t.Fatalf("interrupted transaction journal: %v", err)
			}

			if err := Run(tt.args); err == nil {
				t.Fatal("planning failure unexpectedly succeeded")
			}
			if got := treeDigest(t, dir); got != before {
				t.Fatalf(
					"%s did not recover before planning failure:\n before %s\n  after %s",
					tt.name,
					before,
					got,
				)
			}
			assertNoTransactionArtifacts(t, dir)
		})
	}
}

func TestApplyRejectsInvalidOnlyAndMissingPresetUsage(t *testing.T) {
	_, _ = initedModuleWithStyle(t, preset.StyleNova)
	for _, args := range [][]string{
		{"apply"},
		{"apply", "--preset", "unused", "--only", "everything"},
		{"apply", "--preset", "unused", "unexpected"},
	} {
		err := Run(args)
		if err == nil || !strings.Contains(err.Error(), "usage: gsxui apply") {
			t.Fatalf("Run(%v) error = %v, want apply usage", args, err)
		}
	}
}

func initedModuleWithPreset(t *testing.T, selected preset.Preset) (string, *[][]string) {
	t.Helper()
	dir, commands := initTestModule(t)
	if err := Run([]string{"init", "--preset", presetCode(t, selected)}); err != nil {
		t.Fatal(err)
	}
	return dir, commands
}

func presetCode(t *testing.T, selected preset.Preset) string {
	t.Helper()
	code, err := preset.EncodeShare(selected)
	if err != nil {
		t.Fatal(err)
	}
	return code
}

func readProjectFile(t *testing.T, dir, relative string) []byte {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(relative)))
	if err != nil {
		t.Fatal(err)
	}
	return content
}

func assertProjectPreset(t *testing.T, dir string, want preset.Preset) {
	t.Helper()
	content := readProjectFile(t, dir, "gsxui.preset.json")
	got, err := preset.ParseJSON(content)
	if err != nil {
		t.Fatal(err)
	}
	gotJSON, err := preset.CanonicalJSON(got)
	if err != nil {
		t.Fatal(err)
	}
	wantJSON, err := preset.CanonicalJSON(want)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(gotJSON, wantJSON) {
		t.Fatalf("project preset:\n%s\nwant:\n%s", gotJSON, wantJSON)
	}
}

func interruptTransactionAfterReplacement(t *testing.T, dir string) {
	t.Helper()
	originalHook := transactionCheckpointHook
	transactionCheckpointHook = func(checkpoint transactionCheckpoint) {
		if checkpoint.phase == transactionPhaseReplaced {
			panic("injected planning-boundary interruption")
		}
	}
	defer func() {
		transactionCheckpointHook = originalHook
		if recovered := recover(); recovered != "injected planning-boundary interruption" {
			t.Fatalf("transaction panic = %v", recovered)
		}
	}()
	_ = executeArtifactTransaction(
		dir,
		[]artifact{{
			RelativePath: "web/gsxui/theme.css",
			Content:      []byte("interrupted theme\n"),
			Managed:      true,
		}},
		nil,
		nil,
	)
}

func TestApplyGenerationFailureRestoresTreeAndReportsBothFailures(t *testing.T) {
	dir, _ := initedModuleWithStyle(t, preset.StyleNova)
	if err := Run([]string{"add", "button"}); err != nil {
		t.Fatal(err)
	}
	before := treeDigest(t, dir)
	originalRunner := runCommand
	runCommand = func(_ string, name string, args ...string) error {
		if name == "go" && len(args) >= 2 && args[0] == "tool" && args[1] == "gsx" {
			return fmt.Errorf("injected generation failure")
		}
		return nil
	}
	t.Cleanup(func() { runCommand = originalRunner })

	err := Run([]string{
		"apply",
		"--preset", presetCode(t, preset.Default(preset.StyleMaia)),
		"--yes",
	})
	if err == nil || strings.Count(err.Error(), "injected generation failure") != 2 ||
		!strings.Contains(err.Error(), "regenerate restored sources") {
		t.Fatalf("apply generation error = %v, want both failures", err)
	}
	if got := treeDigest(t, dir); got != before {
		t.Fatalf("generation failure did not restore exact tree:\n before %s\n  after %s", before, got)
	}
	assertNoTransactionArtifacts(t, dir)
}

func TestRecoverAfterGenerationInterruptionRegeneratesRestoredSources(t *testing.T) {
	dir, _ := initedModuleWithStyle(t, preset.StyleNova)
	if err := Run([]string{"add", "button"}); err != nil {
		t.Fatal(err)
	}
	generatedPath := filepath.Join(dir, "ui", "button.x.go")
	const restoredGenerated = "restored generated bytes\n"
	if err := os.WriteFile(generatedPath, []byte(restoredGenerated), 0o644); err != nil {
		t.Fatal(err)
	}
	before := treeDigest(t, dir)
	originalRunner := runCommand
	t.Cleanup(func() { runCommand = originalRunner })
	runCommand = func(_ string, name string, args ...string) error {
		if name == "go" && len(args) >= 2 && args[0] == "tool" && args[1] == "gsx" {
			if err := os.WriteFile(generatedPath, []byte("partial generated bytes\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			panic("injected generation interruption")
		}
		return nil
	}

	func() {
		defer func() {
			if recovered := recover(); recovered != "injected generation interruption" {
				t.Fatalf("generation panic = %v", recovered)
			}
		}()
		_ = Run([]string{
			"apply",
			"--preset", presetCode(t, preset.Default(preset.StyleMaia)),
			"--yes",
		})
	}()
	if got := string(readProjectFile(t, dir, "ui/button.x.go")); got != "partial generated bytes\n" {
		t.Fatalf("generation interruption probe = %q", got)
	}

	var recoveryGenerationCalls int
	runCommand = func(_ string, name string, args ...string) error {
		if name == "go" && len(args) >= 2 && args[0] == "tool" && args[1] == "gsx" {
			recoveryGenerationCalls++
			if err := os.WriteFile(generatedPath, []byte(restoredGenerated), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		return nil
	}
	err := Run([]string{"apply", "--preset", `{"style":"maia"}`, "--yes"})
	if err == nil || !strings.Contains(err.Error(), "preset") {
		t.Fatalf("post-recovery planning error = %v, want malformed preset", err)
	}
	if recoveryGenerationCalls != 1 {
		t.Fatalf("recovery generation calls = %d, want 1", recoveryGenerationCalls)
	}
	if got := treeDigest(t, dir); got != before {
		t.Fatalf("generation interruption recovery did not restore exact tree:\n before %s\n  after %s", before, got)
	}
	assertNoTransactionArtifacts(t, dir)
}
