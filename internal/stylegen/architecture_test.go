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

// shapesImportPath is the leaf package holding every component's shape as pure
// data. It exists to break a bootstrap cycle: the generator must read shapes in
// order to emit the accessors registry/canonical needs to compile at all, so
// the shapes cannot live in registry/canonical.
const shapesImportPath = canonicalImportPath + "/shapes"

// canonicalImporters names the only packages allowed to import registry/canonical:
// just itself, including its own external test package. The generator reads the
// canonical sources as FILES, not as a Go dependency — importing them would
// resurrect the bootstrap cycle the shapes package exists to break.
var canonicalImporters = map[string]bool{
	canonicalImportPath: true,
}

// shapesImporters names the packages allowed to import the shapes leaf package.
// Unlike registry/canonical it is safe to depend on — it is pure data and
// carries no recipe tokens — but it is still registry-private.
//
// internal/stylegen/port needs it for the same reason internal/stylegen
// itself does: composing a recipe class requires the shape's own
// BaseClass/ValueClass methods (Render), and recognizing when an upstream
// `**:data-[slot=x]:` token targets a DIFFERENT gsxui component's root
// rather than a sub-slot of the current one (Transform's cross-component
// KindSlot case) requires knowing every declared component name.
var shapesImporters = map[string]bool{
	"github.com/gsxhq/gsxui/internal/stylegen":      true,
	"github.com/gsxhq/gsxui/internal/stylegen/port": true,
	canonicalImportPath:                             true,
	shapesImportPath:                                true,
}

// TestNothingImportsCanonical fails if any package in the module takes a direct
// dependency on registry/canonical, or if anything outside the registry and the
// generator depends on registry/canonical/shapes. It inspects every package's
// imports, including those of its internal and external test binaries, so a
// test file cannot smuggle the dependency in either.
func TestNothingImportsCanonical(t *testing.T) {
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
		for _, imported := range imports {
			switch {
			case imported == shapesImportPath || strings.HasPrefix(imported, shapesImportPath+"/"):
				if !shapesImporters[importPath] {
					t.Errorf("%s imports %s; the shapes package is registry-private", importPath, shapesImportPath)
				}
			case imported == canonicalImportPath || strings.HasPrefix(imported, canonicalImportPath+"/"):
				if !canonicalImporters[importPath] {
					t.Errorf("%s imports %s; the canonical package must never ship", importPath, canonicalImportPath)
				}
			}
		}
	}
}
