package examples_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gsxhq/gsxui/site/examples"
)

// TestExampleSourceMatchesFile is the drift guard for
// source-shown-is-source-run: for every registered example, the bytes
// examples.Source returns (read from the embedded FS) must be byte-
// identical to the on-disk .gsx file. If they ever diverge — e.g. someone
// edits the embedded copy without touching examples.go's //go:embed
// pattern, or vice versa — the component page would display source that
// isn't what's actually running.
func TestExampleSourceMatchesFile(t *testing.T) {
	for _, component := range examples.Components() {
		for _, ex := range examples.For(component) {
			t.Run(component+"/"+ex.Name, func(t *testing.T) {
				got, err := examples.Source(component, ex.Name)
				if err != nil {
					t.Fatalf("examples.Source(%q, %q): %v", component, ex.Name, err)
				}

				want, err := os.ReadFile(filepath.FromSlash(ex.SourcePath))
				if err != nil {
					t.Fatalf("reading on-disk file %q: %v", ex.SourcePath, err)
				}

				if got != string(want) {
					t.Errorf("embedded Source(%q, %q) != on-disk %q\n--- embedded ---\n%s\n--- on-disk ---\n%s",
						component, ex.Name, ex.SourcePath, got, string(want))
				}
			})
		}
	}
}
