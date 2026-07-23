// Package examples embeds every component's example .gsx source (real,
// compiled gsx components — see registry.go) and serves that exact source
// text for display on the component pages. Source(component, example)
// reads the same embedded bytes the examples_test.go drift test compares
// against the on-disk file, so the page can never show source that differs
// from what's actually running.
package examples

import "embed"

//go:embed */*.gsx
var files embed.FS

// Source returns the embedded source of site/examples/{component}/{example}.gsx.
func Source(component, example string) (string, error) {
	b, err := files.ReadFile(component + "/" + example + ".gsx")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
