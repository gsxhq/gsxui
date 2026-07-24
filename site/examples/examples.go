// Package examples embeds every component's example .gsx source (real,
// compiled gsx components — see registry.go) and serves that exact source
// text for display on the component pages. Source(sourcePath) reads the
// same embedded bytes the examples_test.go drift test compares against the
// on-disk file, so the page can never show source that differs from what's
// actually running.
package examples

import "embed"

//go:embed */*.gsx
var files embed.FS

// Source returns the embedded source at site/examples/{sourcePath} (an
// Example's SourcePath field). SourcePath is keyed by the example's actual
// directory — for switch that's the "switchctl" dir (a Go-keyword
// workaround) and for native-select it's the "nativeselect" dir (a hyphen
// can't appear in a Go package name, same as button-group's buttongroup /
// context-menu's contextmenu), both of which differ from those components'
// registered name ("switch"/"native-select"), so Source can't reconstruct
// the path from the registered component name the way it once could.
func Source(sourcePath string) (string, error) {
	b, err := files.ReadFile(sourcePath)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
