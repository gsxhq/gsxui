package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gsxhq/gsxui/internal/stylegen"
)

func main() {
	check := flag.Bool("check", false, "check generated style sources without writing")
	checkLayers := flag.Bool("check-layers", false,
		"check that every override of compiled component presentation can win the cascade")
	flag.Parse()
	if flag.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "usage: go run ./cmd/stylegen [--check] [--check-layers]")
		os.Exit(2)
	}
	if *check && *checkLayers {
		fmt.Fprintln(os.Stderr, "stylegen: --check and --check-layers are separate gates; run them separately")
		os.Exit(2)
	}

	root, err := repositoryRoot()
	if err == nil {
		if *checkLayers {
			err = stylegen.CheckLayerPrecedence(root)
		} else {
			err = stylegen.GenerateAll(root, *check)
		}
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "stylegen: %v\n", err)
		os.Exit(1)
	}
}

func repositoryRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir, err = filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	for {
		if regularFile(filepath.Join(dir, "go.mod")) &&
			regularFile(filepath.Join(dir, "registry", "canonical", "button.gsx")) {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("repository root not found from %s", dir)
		}
		dir = parent
	}
}

func regularFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}
