package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/gsxhq/gsxui/internal/stylegen"
	"github.com/gsxhq/gsxui/internal/stylegen/port"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "port" {
		runPort(os.Args[2:])
		return
	}

	check := flag.Bool("check", false, "check generated style sources without writing")
	checkLayers := flag.Bool("check-layers", false,
		"check that every override of compiled component presentation can win the cascade")
	checkAuthoring := flag.Bool("check-authoring", false,
		"check that unmigrated components carry no hand-authored presentation and migrated components are generated output")
	flag.Parse()
	if flag.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "usage: go run ./cmd/stylegen [--check] [--check-layers] [--check-authoring]")
		os.Exit(2)
	}
	modes := 0
	for _, m := range []bool{*check, *checkLayers, *checkAuthoring} {
		if m {
			modes++
		}
	}
	if modes > 1 {
		fmt.Fprintln(os.Stderr, "stylegen: --check, --check-layers, and --check-authoring are separate gates; run them separately")
		os.Exit(2)
	}

	root, err := repositoryRoot()
	if err == nil {
		switch {
		case *checkLayers:
			err = stylegen.CheckLayerPrecedence(root)
		case *checkAuthoring:
			err = stylegen.CheckAuthoring(root)
		default:
			err = stylegen.GenerateAll(root, *check)
		}
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "stylegen: %v\n", err)
		os.Exit(1)
	}
}

// runPort implements the `stylegen port` subcommand: port one style (or all
// 8) from a local shadcn-ui checkout into registry/styles/, validated by
// internal/stylegen/port.Run before anything is written.
func runPort(args []string) {
	fs := flag.NewFlagSet("port", flag.ExitOnError)
	upstream := fs.String("upstream", "", "path to a local shadcn-ui checkout root (required)")
	style := fs.String("style", "", "style to port: one of "+joinStyles()+", or \"all\" (required)")
	dryRun := fs.Bool("dry-run", false, "compute the port and report it without writing any files")
	allowUnmapped := fs.Bool("allow-unmapped", false,
		"print unmapped upstream classes/utilities instead of failing when Unmapped > 0")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: go run ./cmd/stylegen port --upstream <path> --style maia|all [--dry-run] [--allow-unmapped]")
		fs.PrintDefaults()
	}
	fs.Parse(args)
	if fs.NArg() != 0 || *upstream == "" || *style == "" {
		fs.Usage()
		os.Exit(2)
	}

	var styles []string
	if *style == "all" {
		styles = slices.Clone(port.AllStyles)
	} else if slices.Contains(port.AllStyles, *style) {
		styles = []string{*style}
	} else {
		fmt.Fprintf(os.Stderr, "stylegen: port: unknown style %q (want one of %s, or \"all\")\n", *style, joinStyles())
		os.Exit(2)
	}

	root, err := repositoryRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "stylegen: %v\n", err)
		os.Exit(1)
	}

	report, err := port.Run(root, *upstream, styles, *dryRun)
	if err != nil {
		fmt.Fprintf(os.Stderr, "stylegen: port: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("stylegen: port: written=%d carried=%d unmapped=%d skipped=%d\n",
		report.Written, report.Carried, report.Unmapped, report.Skipped)

	if report.Unmapped == 0 {
		return
	}
	if *allowUnmapped {
		fmt.Println("stylegen: port: unmapped entries (--allow-unmapped set, continuing):")
		for _, u := range report.UnmappedDetail {
			fmt.Println("  " + u)
		}
		return
	}
	fmt.Fprintln(os.Stderr, "stylegen: port: unmapped entries (pass --allow-unmapped to continue anyway):")
	for _, u := range report.UnmappedDetail {
		fmt.Fprintln(os.Stderr, "  "+u)
	}
	os.Exit(1)
}

func joinStyles() string {
	out := ""
	for i, s := range port.AllStyles {
		if i > 0 {
			out += ", "
		}
		out += s
	}
	return out
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
