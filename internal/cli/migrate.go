package cli

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gsxhq/gsxui/internal/rtl"
)

// runMigrate dispatches `gsxui migrate <name>`. Migrations rewrite files a
// consumer already owns, in place and preserving their edits, which `add
// --overwrite` cannot do.
func runMigrate(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: gsxui migrate <rtl>\n\n  rtl  rewrite vendored components to logical direction classes (requires \"rtl\": true in gsxui.json)")
	}
	switch args[0] {
	case "rtl":
		return runMigrateRTL(args[1:])
	default:
		return fmt.Errorf("unknown migration %q (want rtl)", args[0])
	}
}

// runMigrateRTL applies internal/rtl to every managed .gsx under cfg.UI and
// to the vendored style.css, through the same transaction add uses, so the
// managed hashes move to the migrated content and rollback stays possible.
func runMigrateRTL(args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("usage: gsxui migrate rtl")
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
	if !cfg.RTL {
		return fmt.Errorf(`gsxui.json does not set "rtl": true — run 'gsxui init --rtl' first, then 'gsxui migrate rtl'`)
	}
	styleCSS := filepath.ToSlash(filepath.Join(filepath.Dir(cfg.CSS), "style.css"))
	paths := make([]string, 0, len(cfg.Managed))
	for rel := range cfg.Managed {
		paths = append(paths, rel)
	}
	sort.Strings(paths)

	var artifacts []artifact
	for _, rel := range paths {
		isGSX := strings.HasSuffix(rel, ".gsx")
		if !isGSX && rel != styleCSS {
			continue
		}
		path, err := artifactPath(dir, rel)
		if err != nil {
			return err
		}
		current, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue // a managed file the consumer deleted is theirs to have deleted
			}
			return err
		}
		var next []byte
		if isGSX {
			next, err = rtl.GSX(rel, current)
			if err != nil {
				return err
			}
		} else {
			next = rtl.CSS(current)
		}
		if bytes.Equal(next, current) {
			continue
		}
		artifacts = append(artifacts, artifact{RelativePath: rel, Content: next, Managed: true})
	}
	if len(artifacts) == 0 {
		fmt.Println("migrate rtl: nothing to do")
		return nil
	}
	_, plan, err := artifactPlanWithConfig(cfg, artifacts)
	if err != nil {
		return err
	}
	// overwrite=true: rewriting files the consumer has edited is the point.
	if err := validateArtifactPlan(dir, cfg, plan, true); err != nil {
		return err
	}
	fmt.Printf("migrate rtl: %d file(s)\n", len(artifacts))
	return executeArtifactTransaction(
		dir,
		plan,
		func() error { return generateProject(dir) },
		func() error { return generateProject(dir) },
	)
}
