package stylegen

import (
	"errors"
	"os/exec"
	"strings"
	"testing"
)

// canonicalImportPath is the package that must never ship. registry/canonical
// holds the structural component sources the generator compiles against each
// style's CSS recipe. It is a real Go package so it type-checks and can host
// style-independent behavior tests, but its class attributes carry recipe
// tokens rather than utilities, so anything that imported it would render
// unstyled markup.
const canonicalImportPath = "github.com/gsxhq/gsxui/registry/canonical"

// canonicalImporters names the only packages allowed to import it: the
// generator that reads it, and the package itself (its own external test
// package included).
var canonicalImporters = map[string]bool{
	"github.com/gsxhq/gsxui/internal/stylegen": true,
	canonicalImportPath:                        true,
}

// TestNothingOutsideStylegenImportsCanonical fails if any package in the module
// takes a direct dependency on registry/canonical. It inspects every package's
// imports, including those of its internal and external test binaries, so a
// test file cannot smuggle the dependency in either.
func TestNothingOutsideStylegenImportsCanonical(t *testing.T) {
	const format = `{{.ImportPath}}` +
		`{{range .Imports}} {{.}}{{end}}` +
		`{{range .TestImports}} {{.}}{{end}}` +
		`{{range .XTestImports}} {{.}}{{end}}`

	cmd := exec.Command("go", "list", "-f", format, "./...")
	cmd.Dir = repoRoot(t)
	out, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			t.Fatalf("go list: %v\n%s", err, exitErr.Stderr)
		}
		t.Fatalf("go list: %v", err)
	}

	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		importPath, imports := fields[0], fields[1:]
		if canonicalImporters[importPath] {
			continue
		}
		for _, imported := range imports {
			if imported == canonicalImportPath || strings.HasPrefix(imported, canonicalImportPath+"/") {
				t.Errorf("%s imports %s; the canonical package must never ship", importPath, canonicalImportPath)
			}
		}
	}
}
