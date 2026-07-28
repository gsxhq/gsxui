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
	flag.Parse()
	if flag.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "usage: go run ./cmd/stylegen [--check]")
		os.Exit(2)
	}

	root, err := repositoryRoot()
	if err == nil {
		err = stylegen.GenerateButton(root, *check)
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
			regularFile(filepath.Join(dir, "ui", "button.gsx")) {
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
