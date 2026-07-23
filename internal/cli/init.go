package cli

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	gsxui "github.com/gsxhq/gsxui"
)

// writeVendored installs content at path: absent → write, identical → no-op,
// different → error unless overwrite.
func writeVendored(path string, content []byte, overwrite bool) error {
	if existing, err := os.ReadFile(path); err == nil {
		if string(existing) == string(content) {
			return nil
		}
		if !overwrite {
			return fmt.Errorf("%s differs from the gsxui version — pass --overwrite to replace it", path)
		}
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, content, 0o644)
}

func runInit(args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("usage: gsxui init")
	}
	dir, err := os.Getwd()
	if err != nil {
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
		if err := cfg.Save(dir); err != nil {
			return err
		}
	default:
		return err // unparsable or unreadable: never overwrite
	}

	css, err := fs.ReadFile(gsxui.Files, "assets/gsxui.css")
	if err != nil {
		return err
	}
	if err := writeVendored(filepath.Join(dir, cfg.CSS), css, false); err != nil {
		return err
	}

	core, err := fs.ReadFile(gsxui.Files, "ui/gsxui.js")
	if err != nil {
		return err
	}
	if err := writeVendored(filepath.Join(dir, cfg.JS, "gsxui.js"), core, false); err != nil {
		return err
	}
	indexPath := filepath.Join(dir, cfg.JS, "index.js")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		if err := writeVendored(indexPath, Barrel(nil), false); err != nil {
			return err
		}
	}

	merge, err := fs.ReadFile(gsxui.Files, "merge/merge.go")
	if err != nil {
		return err
	}
	if err := writeVendored(filepath.Join(dir, cfg.UI, "merge", "merge.go"), merge, false); err != nil {
		return err
	}

	if err := ensureClassMerger(dir, module, cfg.UI); err != nil {
		return err
	}

	for _, c := range [][]string{
		{"go", "get", "github.com/gsxhq/gsx@latest"},
		{"go", "get", "github.com/jackielii/tailwind-merge-go@latest"},
		{"go", "get", "-tool", "github.com/gsxhq/gsx/cmd/gsx@latest"},
	} {
		if err := runCommand(dir, c[0], c[1:]...); err != nil {
			return fmt.Errorf("%v: %w", c, err)
		}
	}

	fmt.Printf("gsxui initialized.\n  css:  %s (import it from your entry point)\n  js:   %s/index.js (import it from your entry point)\n  next: gsxui add button\n", cfg.CSS, cfg.JS)
	return nil
}

// ensureClassMerger makes gsx.toml name the vendored merger. Top-level keys
// must precede any [table] header, so a missing directive is prepended.
func ensureClassMerger(dir, module, uiDir string) error {
	path := filepath.Join(dir, "gsx.toml")
	line := fmt.Sprintf("class_merger = %q\n", module+"/"+uiDir+"/merge.Merge")
	existing, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return os.WriteFile(path, []byte(line), 0o644)
	}
	if err != nil {
		return err
	}
	if strings.Contains(string(existing), "class_merger") {
		fmt.Println("gsx.toml already sets class_merger — left unchanged")
		return nil
	}
	return os.WriteFile(path, append([]byte(line+"\n"), existing...), 0o644)
}
