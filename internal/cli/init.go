package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	gsxui "github.com/gsxhq/gsxui"
	"github.com/gsxhq/gsxui/internal/preset"
)

type cssAssetTarget struct {
	source string
	target string
}

func cssAssetTargets(entry string) []cssAssetTarget {
	dir := filepath.Dir(entry)
	return []cssAssetTarget{
		{source: "assets/css/index.css", target: entry},
		{source: "assets/css/foundation.css", target: filepath.Join(dir, "foundation.css")},
		{source: "assets/css/themes/default.css", target: filepath.Join(dir, "theme.css")},
		{source: "assets/css/styles/default.css", target: filepath.Join(dir, "style.css")},
	}
}

func runInit(args []string) error {
	flags := flag.NewFlagSet("init", flag.ContinueOnError)
	presetInput := flags.String("preset", "", "preset file, share code, raw JSON, URL, or - for stdin")
	overwrite := flags.Bool("overwrite", false, "replace locally modified support files")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if len(flags.Args()) != 0 {
		return fmt.Errorf("usage: gsxui init [--preset <file|code|->] [--overwrite]")
	}
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := recoverArtifactTransaction(dir); err != nil {
		return err
	}
	module, err := modulePath(dir)
	if err != nil {
		return err
	}
	cfg, err := LoadConfig(dir)
	switch {
	case err == nil:
	case errors.Is(err, errConfigNotFound):
		cfg = DefaultConfig()
	default:
		return err // unparsable or unreadable: never overwrite
	}

	selectedPreset, err := resolveInitPreset(dir, *presetInput)
	if err != nil {
		return err
	}
	artifacts, err := initArtifacts(dir, module, cfg, selectedPreset)
	if err != nil {
		return err
	}
	integrationArtifacts, err := planScaffoldIntegration(dir, cfg)
	if err != nil {
		return err
	}
	artifacts = append(artifacts, integrationArtifacts...)
	packageArtifacts, err := packageMetadataArtifacts(dir)
	if err != nil {
		return err
	}
	artifacts = append(artifacts, packageArtifacts...)
	_, completePlan, err := artifactPlanWithConfig(cfg, artifacts)
	if err != nil {
		return err
	}
	if err := validateArtifactPlan(dir, cfg, completePlan, *overwrite); err != nil {
		return err
	}
	npmCommand := []string{
		"install",
		"--save-dev",
		"tailwindcss@^4.3.3",
		"@tailwindcss/vite@^4.3.3",
		"tw-animate-css@^1.4.0",
	}
	commands := [][]string{
		{"go", "get", "github.com/gsxhq/gsx@latest"},
		{"go", "get", "github.com/jackielii/tailwind-merge-go@latest"},
		{"go", "get", "-tool", "github.com/gsxhq/gsx/cmd/gsx@latest"},
	}
	if err := executeArtifactTransaction(
		dir,
		completePlan,
		func() error {
			if err := runCommand(dir, "npm", npmCommand...); err != nil {
				return fmt.Errorf("npm %v: %w", npmCommand, err)
			}
			for _, command := range commands {
				if err := runCommand(dir, command[0], command[1:]...); err != nil {
					return fmt.Errorf("%v: %w", command, err)
				}
			}
			return generateProject(dir)
		},
		func() error { return generateProject(dir) },
	); err != nil {
		return err
	}

	fmt.Fprintf(
		commandStdout,
		"gsxui initialized.\n  css:  %s\n  js:   %s/index.js\n  vite: Tailwind CSS configured\n  next: gsxui add button\n",
		cfg.CSS,
		cfg.JS,
	)
	return nil
}

func resolveInitPreset(dir, input string) (preset.Preset, error) {
	if input != "" {
		resolved, err := (preset.InputResolver{Stdin: commandStdin}).Resolve(context.Background(), input)
		if err != nil {
			return preset.Preset{}, fmt.Errorf("resolve init preset: %w", err)
		}
		return resolved, nil
	}

	path := filepath.Join(dir, "gsxui.preset.json")
	content, err := os.ReadFile(path)
	if err == nil {
		resolved, err := preset.ParseJSON(content)
		if err != nil {
			return preset.Preset{}, fmt.Errorf("load existing gsxui.preset.json: %w", err)
		}
		return resolved, nil
	}
	if !os.IsNotExist(err) {
		return preset.Preset{}, fmt.Errorf("read existing gsxui.preset.json: %w", err)
	}
	return preset.Default(preset.StyleNova), nil
}

func initArtifacts(dir, module string, cfg Config, selected preset.Preset) ([]artifact, error) {
	presetJSON, err := preset.CanonicalJSON(selected)
	if err != nil {
		return nil, err
	}
	themeCSS, err := preset.ThemeCSS(selected)
	if err != nil {
		return nil, err
	}

	artifacts := []artifact{{
		RelativePath: "gsxui.preset.json",
		Content:      presetJSON,
		Managed:      true,
	}}
	for _, asset := range cssAssetTargets(cfg.CSS) {
		var content []byte
		if asset.source == "assets/css/themes/default.css" {
			content = themeCSS
		} else {
			content, err = fs.ReadFile(gsxui.Files, asset.source)
			if err != nil {
				return nil, err
			}
		}
		artifacts = append(artifacts, artifact{
			RelativePath: filepath.ToSlash(asset.target),
			Content:      content,
			Managed:      true,
		})
	}

	core, err := fs.ReadFile(gsxui.Files, "ui/gsxui.js")
	if err != nil {
		return nil, err
	}
	artifacts = append(artifacts, artifact{
		RelativePath: filepath.ToSlash(filepath.Join(cfg.JS, "gsxui.js")),
		Content:      core,
		Managed:      true,
	})

	indexRelative := filepath.ToSlash(filepath.Join(cfg.JS, "index.js"))
	indexPath, err := artifactPath(dir, indexRelative)
	if err != nil {
		return nil, err
	}
	index, err := os.ReadFile(indexPath)
	switch {
	case os.IsNotExist(err):
		barrel, err := behaviorBarrelArtifact(dir, cfg, nil)
		if err != nil {
			return nil, err
		}
		artifacts = append(artifacts, barrel)
	case err != nil:
		return nil, fmt.Errorf("read behavior barrel: %w", err)
	case strings.HasPrefix(string(index), barrelHeader):
		barrel, err := behaviorBarrelArtifact(dir, cfg, nil)
		if err != nil {
			return nil, err
		}
		artifacts = append(artifacts, barrel)
	}

	merge, err := fs.ReadFile(gsxui.Files, "merge/merge.go")
	if err != nil {
		return nil, err
	}
	artifacts = append(artifacts, artifact{
		RelativePath: filepath.ToSlash(filepath.Join(cfg.UI, "merge", "merge.go")),
		Content:      merge,
		Managed:      true,
	})

	classMerger, err := classMergerArtifact(dir, module, cfg.UI)
	if err != nil {
		return nil, err
	}
	if classMerger != nil {
		artifacts = append(artifacts, *classMerger)
	}
	return artifacts, nil
}

func classMergerArtifact(dir, module, uiDir string) (*artifact, error) {
	const relative = "gsx.toml"
	path, err := artifactPath(dir, relative)
	if err != nil {
		return nil, err
	}
	line := fmt.Sprintf("class_merger = %q\n", module+"/"+uiDir+"/merge.Merge")
	existing, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &artifact{RelativePath: relative, Content: []byte(line)}, nil
	}
	if err != nil {
		return nil, err
	}
	if strings.Contains(string(existing), "class_merger") {
		fmt.Println("gsx.toml already sets class_merger — left unchanged")
		return nil, nil
	}
	return &artifact{
		RelativePath: relative,
		Content:      append([]byte(line+"\n"), existing...),
	}, nil
}
