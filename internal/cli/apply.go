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
	replacedComponents    []string // relative paths of component sources replaced wholesale (not merged)
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

	if len(axes) == 0 {
		return applyPlan{}, nil
	}

	installed, err := installedStyledComponents(dir, cfg)
	if err != nil {
		return applyPlan{}, err
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
	var replacedComponents []string
	if slices.Contains(axes, "style") {
		for _, name := range installed {
			has, err := hasStyledArtifact(selected.Style, name)
			if err != nil {
				return applyPlan{}, err
			}
			if !has {
				// No per-style artifact (e.g. icon) — vendors identically
				// under every style, so there is nothing to replace.
				continue
			}
			styled, err := selectedStyledArtifact(module, cfg, name, selected.Style)
			if err != nil {
				return applyPlan{}, err
			}
			planned = append(planned, styled)
			replacedComponents = append(replacedComponents, styled.RelativePath)
		}
	}

	modifiedManaged, unmanaged, err := applyOwnershipConflicts(
		dir,
		cfg,
		planned,
		replacedComponents,
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
	replacedSet := make(map[string]bool, len(replacedComponents))
	for _, path := range replacedComponents {
		replacedSet[path] = true
	}
	files := make([]string, len(changed))
	var actuallyReplaced []string
	for index, artifact := range changed {
		files[index] = artifact.RelativePath
		if replacedSet[artifact.RelativePath] {
			actuallyReplaced = append(actuallyReplaced, artifact.RelativePath)
		}
	}
	return applyPlan{
		axes:                  axes,
		artifacts:             complete,
		files:                 files,
		modifiedManaged:       modifiedManaged,
		unmanagedReplacements: unmanaged,
		replacedComponents:    actuallyReplaced,
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

// selectedStyledArtifact vendors name's exact source for style — every
// component (not only Button) has real per-style output at
// registry/generated/<style>/<name>.gsx now that all 8 styles are ported.
func selectedStyledArtifact(module string, cfg Config, name string, style preset.Style) (artifact, error) {
	sourcePath := "registry/generated/" + string(style) + "/" + name + ".gsx"
	source, err := fs.ReadFile(gsxui.Files, sourcePath)
	if err != nil {
		return artifact{}, err
	}
	rewritten, err := RewriteGsx(source, module, cfg.UI)
	if err != nil {
		return artifact{}, fmt.Errorf("rewrite %s: %w", sourcePath, err)
	}
	return artifact{
		RelativePath: filepath.ToSlash(filepath.Join(cfg.UI, name+".gsx")),
		Content:      rewritten,
		Managed:      true,
	}, nil
}

// hasStyledArtifact reports whether registry/generated/<style>/<name>.gsx
// exists. It is false for directory components like icon, which have no
// per-style output and vendor identically under every style.
func hasStyledArtifact(style preset.Style, name string) (bool, error) {
	_, err := fs.Stat(gsxui.Files, "registry/generated/"+string(style)+"/"+name+".gsx")
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func installedStyledComponents(dir string, cfg Config) ([]string, error) {
	components := stylecontract.All()
	installed := make([]string, 0, len(components))
	for _, component := range components {
		name := component.RegistryName
		sourceInfo, directoryErr := fs.Stat(gsxui.Files, "ui/"+name)
		directoryComponent := directoryErr == nil && sourceInfo.IsDir()
		if !directoryComponent {
			var fileErr error
			sourceInfo, fileErr = fs.Stat(gsxui.Files, "ui/"+name+".gsx")
			if fileErr != nil {
				return nil, fmt.Errorf("inspect registry component %s: %w", name, fileErr)
			}
		}
		var (
			target  string
			pathErr error
		)
		if directoryComponent {
			relative := filepath.ToSlash(filepath.Join(cfg.UI, name))
			target, pathErr = artifactDirectoryPath(dir, relative)
		} else {
			relative := filepath.ToSlash(filepath.Join(cfg.UI, name+".gsx"))
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
	requireOwnedComponents []string,
) (modifiedManaged, unmanaged []string, err error) {
	requireOwned := make(map[string]bool, len(requireOwnedComponents))
	for _, path := range requireOwnedComponents {
		requireOwned[path] = true
	}
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
		if requireOwned[planned.RelativePath] && !managed {
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
	replaced := make(map[string]bool, len(plan.replacedComponents))
	for _, relative := range plan.replacedComponents {
		replaced[relative] = true
	}
	for _, relative := range plan.files {
		if replaced[relative] {
			fmt.Printf("    %s (component source replacement)\n", relative)
		} else {
			fmt.Printf("    %s\n", relative)
		}
	}
	if len(plan.replacedComponents) != 0 {
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
