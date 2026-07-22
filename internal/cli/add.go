package cli

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	gsxui "github.com/gsxhq/gsxui"
	"github.com/gsxhq/gsxui/internal/registry"
)

func runAdd(args []string) error {
	fs2 := flag.NewFlagSet("add", flag.ContinueOnError)
	overwrite := fs2.Bool("overwrite", false, "replace locally modified files")
	if err := fs2.Parse(args); err != nil {
		return err
	}
	names := fs2.Args()
	if len(names) == 0 {
		return fmt.Errorf("usage: gsxui add [--overwrite] <component>...")
	}
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	cfg, err := LoadConfig(dir)
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
	fmt.Printf("adding: %s\n", strings.Join(resolved, " "))

	for _, name := range resolved {
		entries, err := fs.ReadDir(gsxui.Files, "ui/"+name)
		if err != nil {
			return err
		}
		for _, e := range entries {
			fname := e.Name()
			if e.IsDir() || strings.HasSuffix(fname, ".x.go") || strings.HasSuffix(fname, "_test.go") || strings.HasSuffix(fname, ".js") {
				continue
			}
			src, err := fs.ReadFile(gsxui.Files, "ui/"+name+"/"+fname)
			if err != nil {
				return err
			}
			if strings.HasSuffix(fname, ".gsx") || strings.HasSuffix(fname, ".go") {
				src = RewriteGsx(src, module, cfg.UI)
			}
			if err := writeVendored(filepath.Join(dir, cfg.UI, name, fname), src, *overwrite); err != nil {
				return err
			}
		}
		if registry.HasJS(name) {
			src, err := fs.ReadFile(gsxui.Files, "ui/"+name+"/"+name+".js")
			if err != nil {
				return err
			}
			if err := writeVendored(filepath.Join(dir, cfg.JS, name+".js"), RewriteJS(src), *overwrite); err != nil {
				return err
			}
		}
	}

	notice, err := fs.ReadFile(gsxui.Files, "NOTICE.md")
	if err != nil {
		return err
	}
	if err := writeVendored(filepath.Join(dir, cfg.UI, "NOTICE.md"), notice, true); err != nil {
		return err
	}

	if err := regenBarrel(dir, cfg); err != nil {
		return err
	}
	if err := runCommand(dir, "go", "tool", "gsx", "generate"); err != nil {
		return fmt.Errorf("gsx generate: %w", err)
	}
	fmt.Println("done — build with: go build ./...")
	return nil
}

// regenBarrel rewrites index.js from the behaviors present on disk.
func regenBarrel(dir string, cfg Config) error {
	matches, err := filepath.Glob(filepath.Join(dir, cfg.JS, "*.js"))
	if err != nil {
		return err
	}
	var behaviors []string
	for _, m := range matches {
		base := filepath.Base(m)
		if base == "index.js" {
			continue
		}
		behaviors = append(behaviors, strings.TrimSuffix(base, ".js"))
	}
	sort.Strings(behaviors)
	return os.WriteFile(filepath.Join(dir, cfg.JS, "index.js"), Barrel(behaviors), 0o644)
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
