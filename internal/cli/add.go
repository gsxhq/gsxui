package cli

import (
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	gsxui "github.com/gsxhq/gsxui"
	"github.com/gsxhq/gsxui/internal/preset"
	"github.com/gsxhq/gsxui/internal/registry"
	"github.com/pmezard/go-difflib/difflib"
)

// runAdd vendors the requested components (and their transitive deps) into
// cfg.UI. Vendored .gsx files keep their "package ui" clause regardless of
// cfg.UI's basename — Go allows a package's declared name to differ from
// its directory name, so e.g. `gsxui add button` into cfg.UI = "components"
// still produces components/button.gsx starting with "package ui", and it
// compiles and imports fine as "components".
func runAdd(args []string) error {
	flags := flag.NewFlagSet("add", flag.ContinueOnError)
	diff := flags.Bool("diff", false, "show selected source against local files without writing")
	overwrite := flags.Bool("overwrite", false, "replace locally modified files")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *diff && *overwrite {
		return fmt.Errorf("--diff and --overwrite are mutually exclusive")
	}
	names := flags.Args()
	if len(names) == 0 {
		return fmt.Errorf("usage: gsxui add [--diff|--overwrite] <component>...")
	}
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	if !*diff {
		if err := recoverArtifactTransaction(dir); err != nil {
			return err
		}
	}
	cfg, err := LoadConfig(dir)
	if err != nil {
		return err
	}
	selectedPreset, err := loadProjectPreset(dir)
	if err != nil {
		return err
	}
	module, err := modulePath(dir)
	if err != nil {
		return err
	}
	resolved, err := registry.Resolve(names)
	if err != nil {
		return err
	}

	artifacts, err := addArtifacts(dir, module, cfg, selectedPreset, resolved)
	if err != nil {
		return err
	}
	if *diff {
		return printArtifactDiff(dir, artifacts)
	}
	_, completePlan, err := artifactPlanWithConfig(cfg, artifacts)
	if err != nil {
		return err
	}
	if err := validateArtifactPlan(dir, cfg, completePlan, *overwrite); err != nil {
		return err
	}

	fmt.Printf("adding: %s\n", strings.Join(resolved, " "))
	if err := executeArtifactTransaction(
		dir,
		completePlan,
		func() error { return generateProject(dir) },
		func() error { return generateProject(dir) },
	); err != nil {
		return err
	}
	fmt.Println("done — build with: go build ./...")
	return nil
}

func generateProject(dir string) error {
	if err := runCommand(dir, "go", "tool", "gsx", "generate"); err != nil {
		return fmt.Errorf("gsx generate: %w — if the gsx tool is missing, run 'gsxui init' (or 'go get -tool github.com/gsxhq/gsx/cmd/gsx@latest')", err)
	}
	return nil
}

func loadProjectPreset(dir string) (preset.Preset, error) {
	path := filepath.Join(dir, "gsxui.preset.json")
	content, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return preset.Preset{}, fmt.Errorf("gsxui.preset.json not found — run 'gsxui init' first")
	}
	if err != nil {
		return preset.Preset{}, fmt.Errorf("reading gsxui.preset.json: %w", err)
	}
	selected, err := preset.ParseJSON(content)
	if err != nil {
		return preset.Preset{}, fmt.Errorf("parsing gsxui.preset.json: %w", err)
	}
	return selected, nil
}

func addArtifacts(dir, module string, cfg Config, selected preset.Preset, resolved []string) ([]artifact, error) {
	var artifacts []artifact
	for _, name := range resolved {
		if fi, err := fs.Stat(gsxui.Files, "ui/"+name); err == nil && fi.IsDir() {
			// Directory component (icon): vendors as its own package under
			// <cfg.UI>/<name>/ — it must stay a package (icon.New) and its
			// generated data tables stay out of the user's ui namespace.
			entries, err := fs.ReadDir(gsxui.Files, "ui/"+name)
			if err != nil {
				return nil, err
			}
			for _, e := range entries {
				fname := e.Name()
				if e.IsDir() || strings.HasSuffix(fname, ".x.go") || strings.HasSuffix(fname, "_test.go") {
					continue
				}
				src, err := fs.ReadFile(gsxui.Files, "ui/"+name+"/"+fname)
				if err != nil {
					return nil, err
				}
				if strings.HasSuffix(fname, ".gsx") || strings.HasSuffix(fname, ".go") {
					src, err = RewriteGsx(src, module, cfg.UI)
					if err != nil {
						return nil, fmt.Errorf("rewrite %s/%s: %w", name, fname, err)
					}
				}
				artifacts = append(artifacts, artifact{
					RelativePath: filepath.ToSlash(filepath.Join(cfg.UI, name, fname)),
					Content:      src,
					Managed:      true,
				})
			}
		} else {
			styled, err := selectedStyledArtifact(module, cfg, name, selected.Style)
			if err != nil {
				return nil, err
			}
			artifacts = append(artifacts, styled)
		}
		if registry.HasJS(name) {
			src, err := fs.ReadFile(gsxui.Files, "ui/"+name+".js")
			if err != nil {
				return nil, err
			}
			artifacts = append(artifacts, artifact{
				RelativePath: filepath.ToSlash(filepath.Join(cfg.JS, name+".js")),
				Content:      src,
				Managed:      true,
			})
		}
	}

	notice, err := fs.ReadFile(gsxui.Files, "NOTICE.md")
	if err != nil {
		return nil, err
	}
	artifacts = append(artifacts, artifact{
		RelativePath: filepath.ToSlash(filepath.Join(cfg.UI, "NOTICE.md")),
		Content:      notice,
		Managed:      true,
	})

	barrel, err := behaviorBarrelArtifact(dir, cfg, resolved)
	if err != nil {
		return nil, err
	}
	artifacts = append(artifacts, barrel)
	return artifacts, nil
}

func behaviorBarrelArtifact(dir string, cfg Config, resolved []string) (artifact, error) {
	jsRelative := filepath.ToSlash(cfg.JS)
	jsDir, err := artifactDirectoryPath(dir, jsRelative)
	if err != nil {
		return artifact{}, err
	}
	matches, err := filepath.Glob(filepath.Join(jsDir, "*.js"))
	if err != nil {
		return artifact{}, err
	}
	behaviorSet := make(map[string]bool)
	for _, match := range matches {
		name := strings.TrimSuffix(filepath.Base(match), ".js")
		if name != "index" && name != "gsxui" {
			behaviorSet[name] = true
		}
	}
	for _, name := range resolved {
		if registry.HasJS(name) {
			behaviorSet[name] = true
		}
	}
	behaviors := make([]string, 0, len(behaviorSet))
	for name := range behaviorSet {
		behaviors = append(behaviors, name)
	}
	return artifact{
		RelativePath: filepath.ToSlash(filepath.Join(cfg.JS, "index.js")),
		Content:      Barrel(behaviors),
		Managed:      true,
	}, nil
}

// barrelHeader is the first line Barrel emits — its presence marks index.js
// as gsxui-generated and safe to regenerate freely.
const barrelHeader = "// Generated by gsxui — edit by adding/removing components, not by hand."

func printArtifactDiff(dir string, artifacts []artifact) error {
	var output bytes.Buffer
	changed := false
	for _, planned := range artifacts {
		target, err := artifactPath(dir, planned.RelativePath)
		if err != nil {
			return err
		}
		existing, err := os.ReadFile(target)
		if os.IsNotExist(err) {
			existing = nil
		} else if err != nil {
			return fmt.Errorf("read artifact %s: %w", target, err)
		}
		if bytes.Equal(existing, planned.Content) {
			continue
		}
		diff, err := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
			A:        difflib.SplitLines(string(existing)),
			B:        difflib.SplitLines(string(planned.Content)),
			FromFile: planned.RelativePath + " (local)",
			ToFile:   planned.RelativePath + " (selected)",
			Context:  3,
		})
		if err != nil {
			return fmt.Errorf("diff %s: %w", planned.RelativePath, err)
		}
		output.WriteString(diff)
		changed = true
	}
	if !changed {
		fmt.Println("(no changes)")
		return nil
	}
	fmt.Print(output.String())
	return nil
}

func runList(args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("usage: gsxui list")
	}
	names, err := registry.Components()
	if err != nil {
		return err
	}
	for _, n := range names {
		deps, err := registry.Deps(n)
		if err != nil {
			return err
		}
		line := n
		if registry.HasJS(n) {
			line += " (js)"
		}
		if len(deps) > 0 {
			line += " — requires " + strings.Join(deps, ", ")
		}
		fmt.Println(line)
	}
	return nil
}
