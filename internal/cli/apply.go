package cli

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	gsxui "github.com/gsxhq/gsxui"
	"github.com/gsxhq/gsxui/internal/preset"
	"github.com/gsxhq/gsxui/internal/stylecontract"
)

type applyPlan struct {
	axes                  []string
	artifacts             []artifact
	files                 []string
	modifiedManaged       []string
	unmanagedReplacements []string
	buttonReplacement     bool
}

func runApply(args []string) error {
	flags := flag.NewFlagSet("apply", flag.ContinueOnError)
	presetInput := flags.String("preset", "", "preset file, share code, raw JSON, URL, or - for stdin")
	only := flags.String("only", "", "apply only theme or style")
	yes := flags.Bool("yes", false, "apply without interactive confirmation")
	overwrite := flags.Bool("overwrite", false, "replace locally modified managed files")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *presetInput == "" || len(flags.Args()) != 0 ||
		(*only != "" && *only != "theme" && *only != "style") {
		return fmt.Errorf("usage: gsxui apply --preset <file|code|-> [--only theme|style] [--yes] [--overwrite]")
	}

	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := recoverArtifactTransaction(dir); err != nil {
		return err
	}
	cfg, err := LoadConfig(dir)
	if err != nil {
		return err
	}
	current, err := loadProjectPreset(dir)
	if err != nil {
		return err
	}
	incoming, err := (preset.InputResolver{Stdin: commandStdin}).Resolve(context.Background(), *presetInput)
	if err != nil {
		return fmt.Errorf("resolve apply preset: %w", err)
	}
	module, err := modulePath(dir)
	if err != nil {
		return err
	}
	plan, err := planPresetApply(dir, module, cfg, current, incoming, *only, *overwrite)
	if err != nil {
		return err
	}
	if len(plan.axes) == 0 {
		fmt.Println("(no changes)")
		return nil
	}

	printApplySummary(plan)
	if !*yes {
		confirmed, err := confirmApply()
		if err != nil {
			return err
		}
		if !confirmed {
			fmt.Println("apply declined.")
			return nil
		}
	}
	if err := executeArtifactTransaction(
		dir,
		plan.artifacts,
		func() error { return generateProject(dir) },
		func() error { return generateProject(dir) },
	); err != nil {
		return err
	}
	fmt.Println("preset applied.")
	return nil
}

func planPresetApply(
	dir, module string,
	cfg Config,
	current, incoming preset.Preset,
	only string,
	overwrite bool,
) (applyPlan, error) {
	selected := mergePresetAxes(current, incoming, only)
	axes := changedPresetAxes(current, selected)

	installed, err := installedStyledComponents(dir, cfg)
	if err != nil {
		return applyPlan{}, err
	}
	if selected.Style == preset.StyleMaia {
		unsupported := make([]string, 0, len(installed))
		for _, name := range installed {
			if name != "button" {
				unsupported = append(unsupported, name)
			}
		}
		if len(unsupported) != 0 {
			return applyPlan{}, fmt.Errorf(
				"style maia supports only the standalone Button; installed styled components: %s",
				strings.Join(unsupported, ", "),
			)
		}
	}
	if len(axes) == 0 {
		return applyPlan{}, nil
	}

	presetJSON, err := preset.CanonicalJSON(selected)
	if err != nil {
		return applyPlan{}, err
	}
	planned := []artifact{{
		RelativePath: "gsxui.preset.json",
		Content:      presetJSON,
		Managed:      true,
	}}
	if slices.Contains(axes, "radius") ||
		slices.Contains(axes, "theme.light") ||
		slices.Contains(axes, "theme.dark") {
		themeCSS, err := preset.ThemeCSS(selected)
		if err != nil {
			return applyPlan{}, err
		}
		planned = append(planned, artifact{
			RelativePath: themeArtifactRelativePath(cfg),
			Content:      themeCSS,
			Managed:      true,
		})
	}
	buttonInstalled := slices.Contains(installed, "button")
	if slices.Contains(axes, "style") && buttonInstalled {
		button, err := selectedButtonArtifact(module, cfg, selected.Style)
		if err != nil {
			return applyPlan{}, err
		}
		planned = append(planned, button)
	}

	modifiedManaged, unmanaged, err := applyOwnershipConflicts(
		dir,
		cfg,
		planned,
		slices.Contains(axes, "style") && buttonInstalled,
	)
	if err != nil {
		return applyPlan{}, err
	}
	if !overwrite && (len(modifiedManaged) != 0 || len(unmanaged) != 0) {
		var conflicts []string
		if len(modifiedManaged) != 0 {
			conflicts = append(conflicts, "modified managed files: "+strings.Join(modifiedManaged, ", "))
		}
		if len(unmanaged) != 0 {
			conflicts = append(conflicts, "unmanaged replacements: "+strings.Join(unmanaged, ", "))
		}
		return applyPlan{}, fmt.Errorf(
			"%s — pass --overwrite to replace them; commit or stash local changes first",
			strings.Join(conflicts, "; "),
		)
	}

	_, complete, err := artifactPlanWithConfig(cfg, planned)
	if err != nil {
		return applyPlan{}, err
	}
	if err := validateArtifactPlan(dir, cfg, complete, true); err != nil {
		return applyPlan{}, err
	}
	changed, err := changedArtifacts(dir, complete)
	if err != nil {
		return applyPlan{}, err
	}
	files := make([]string, len(changed))
	buttonReplacement := false
	buttonRelative := filepath.ToSlash(filepath.Join(cfg.UI, "button.gsx"))
	for index, artifact := range changed {
		files[index] = artifact.RelativePath
		if artifact.RelativePath == buttonRelative {
			buttonReplacement = true
		}
	}
	return applyPlan{
		axes:                  axes,
		artifacts:             complete,
		files:                 files,
		modifiedManaged:       modifiedManaged,
		unmanagedReplacements: unmanaged,
		buttonReplacement:     buttonReplacement,
	}, nil
}

func mergePresetAxes(current, incoming preset.Preset, only string) preset.Preset {
	switch only {
	case "theme":
		incoming.Style = current.Style
		return incoming
	case "style":
		current.Style = incoming.Style
		return current
	default:
		return incoming
	}
}

func changedPresetAxes(current, selected preset.Preset) []string {
	var axes []string
	if current.Style != selected.Style {
		axes = append(axes, "style")
	}
	if current.Radius != selected.Radius {
		axes = append(axes, "radius")
	}
	if !maps.Equal(current.Theme.Light, selected.Theme.Light) {
		axes = append(axes, "theme.light")
	}
	if !maps.Equal(current.Theme.Dark, selected.Theme.Dark) {
		axes = append(axes, "theme.dark")
	}
	return axes
}

func themeArtifactRelativePath(cfg Config) string {
	return filepath.ToSlash(filepath.Join(filepath.Dir(cfg.CSS), "theme.css"))
}

func selectedButtonArtifact(module string, cfg Config, style preset.Style) (artifact, error) {
	sourcePath := "registry/generated/" + string(style) + "/button.gsx"
	source, err := fs.ReadFile(gsxui.Files, sourcePath)
	if err != nil {
		return artifact{}, err
	}
	rewritten, err := RewriteGsx(source, module, cfg.UI)
	if err != nil {
		return artifact{}, fmt.Errorf("rewrite %s: %w", sourcePath, err)
	}
	return artifact{
		RelativePath: filepath.ToSlash(filepath.Join(cfg.UI, "button.gsx")),
		Content:      rewritten,
		Managed:      true,
	}, nil
}

// styledComponentSourceBasename returns the ui/ file (or directory) basename
// that vendors a style-contract component's RegistryName. This is the
// identity function for every component except "toast" and "toaster", which
// are two typed style-contract components backed by the single
// ui/sonner.gsx — gsxui renders both the Toaster and the Toast DOM itself,
// so they're split in the contract even though they still ship as one gsx
// source file (see internal/stylecontract/contracts_toast.go).
func styledComponentSourceBasename(name string) string {
	switch name {
	case "toast", "toaster":
		return "sonner"
	default:
		return name
	}
}

func installedStyledComponents(dir string, cfg Config) ([]string, error) {
	components := stylecontract.All()
	installed := make([]string, 0, len(components))
	for _, component := range components {
		name := component.RegistryName
		source := styledComponentSourceBasename(name)
		sourceInfo, directoryErr := fs.Stat(gsxui.Files, "ui/"+source)
		directoryComponent := directoryErr == nil && sourceInfo.IsDir()
		if !directoryComponent {
			var fileErr error
			sourceInfo, fileErr = fs.Stat(gsxui.Files, "ui/"+source+".gsx")
			if fileErr != nil {
				return nil, fmt.Errorf("inspect registry component %s: %w", name, fileErr)
			}
		}
		var (
			target  string
			pathErr error
		)
		if directoryComponent {
			relative := filepath.ToSlash(filepath.Join(cfg.UI, source))
			target, pathErr = artifactDirectoryPath(dir, relative)
		} else {
			relative := filepath.ToSlash(filepath.Join(cfg.UI, source+".gsx"))
			target, pathErr = artifactPath(dir, relative)
		}
		if pathErr != nil {
			return nil, pathErr
		}
		info, err := os.Stat(target)
		switch {
		case os.IsNotExist(err):
			continue
		case err != nil:
			return nil, fmt.Errorf("inspect installed component %s: %w", name, err)
		case directoryComponent != info.IsDir():
			return nil, fmt.Errorf("installed component target %s has the wrong file type", target)
		default:
			installed = append(installed, name)
		}
	}
	slices.Sort(installed)
	return installed, nil
}

func applyOwnershipConflicts(
	dir string,
	cfg Config,
	artifacts []artifact,
	requireOwnedButton bool,
) (modifiedManaged, unmanaged []string, err error) {
	buttonRelative := filepath.ToSlash(filepath.Join(cfg.UI, "button.gsx"))
	for _, planned := range artifacts {
		target, pathErr := artifactPath(dir, planned.RelativePath)
		if pathErr != nil {
			return nil, nil, pathErr
		}
		existing, readErr := os.ReadFile(target)
		switch {
		case os.IsNotExist(readErr):
			continue
		case readErr != nil:
			return nil, nil, fmt.Errorf("read apply target %s: %w", target, readErr)
		}
		recorded, managed := cfg.Managed[planned.RelativePath]
		if requireOwnedButton && planned.RelativePath == buttonRelative && !managed {
			unmanaged = append(unmanaged, planned.RelativePath)
			continue
		}
		if bytes.Equal(existing, planned.Content) {
			continue
		}
		if !managed {
			unmanaged = append(unmanaged, planned.RelativePath)
			continue
		}
		if contentHash(existing) != recorded {
			modifiedManaged = append(modifiedManaged, planned.RelativePath)
		}
	}
	return modifiedManaged, unmanaged, nil
}

func printApplySummary(plan applyPlan) {
	fmt.Println("apply preset:")
	fmt.Printf("  axes: %s\n", strings.Join(plan.axes, ", "))
	fmt.Println("  files:")
	for _, relative := range plan.files {
		if plan.buttonReplacement && strings.HasSuffix(relative, "/button.gsx") {
			fmt.Printf("    %s (component source replacement)\n", relative)
		} else {
			fmt.Printf("    %s\n", relative)
		}
	}
	if plan.buttonReplacement {
		fmt.Println("  component source is replaced, not merged.")
	}
	if len(plan.modifiedManaged) != 0 {
		fmt.Printf("  modified managed conflicts: %s\n", strings.Join(plan.modifiedManaged, ", "))
	}
	if len(plan.unmanagedReplacements) != 0 {
		fmt.Printf("  unmanaged replacements: %s\n", strings.Join(plan.unmanagedReplacements, ", "))
	}
	fmt.Println("  commit or stash local changes before applying.")
}

func confirmApply() (bool, error) {
	fmt.Print("continue? [y/N] ")
	line, err := bufio.NewReader(commandStdin).ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return false, fmt.Errorf("read apply confirmation: %w", err)
	}
	switch strings.ToLower(strings.TrimSpace(line)) {
	case "y", "yes":
		return true, nil
	default:
		return false, nil
	}
}
