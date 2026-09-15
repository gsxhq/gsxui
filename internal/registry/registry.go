// Package registry derives the vendorable component set from the embedded
// filesystem — the component list, inter-component dependencies, and
// behavior-JS presence. Derived, never declared: it cannot drift from the
// code it describes.
//
// ui/ is one flat package, so a component is a .gsx file basename:
// ui/button.gsx is component "button". A directory under ui/ is its own
// vendored package: ui/icon (a component, so icon.New reads as a name) and
// ui/i18n (a helper, the T message type). Dependencies come from two
// sources: a ui/<sub> import in .gsx source, and — because intra-package
// references have no import to scan — identifiers in a component's
// generated .x.go that another component's .x.go declares, resolved with
// go/parser against a declaration index.
package registry

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"regexp"
	"slices"
	"sort"
	"strings"
	"sync"

	gsxui "github.com/gsxhq/gsxui"

	gsxast "github.com/gsxhq/gsx/ast"
	gsxparser "github.com/gsxhq/gsx/parser"
)

// uiSubpackageImportRe captures the name of a ui/ sub-package a component's
// .gsx imports — ui/icon, ui/i18n. It is the only dependency edge a
// directory package can produce, because its identifiers are reached
// through the import rather than bare.
var uiSubpackageImportRe = regexp.MustCompile(`"github\.com/gsxhq/gsxui/ui/([^"/]+)"`)

// classified splits every vendorable entry under ui/ into components (they
// declare at least one component) and helpers (Go-only, such as the ui/i18n
// package holding the T message type). A directory is a component when a
// .gsx directly inside it declares one (ui/icon); a file is a component
// when it declares one. The embedded tree is immutable, so this runs once.
var classified = sync.OnceValues(func() (classification, error) {
	entries, err := fs.ReadDir(gsxui.Files, "ui")
	if err != nil {
		return classification{}, err
	}
	c := classification{directories: map[string]bool{}}
	for _, e := range entries {
		name := e.Name()
		var declares bool
		if e.IsDir() {
			declares, err = directoryDeclaresComponent(name)
			if err != nil {
				return classification{}, err
			}
			c.directories[name] = true
		} else {
			var ok bool
			name, ok = strings.CutSuffix(name, ".gsx")
			if !ok {
				continue
			}
			declares, err = fileDeclaresComponent(e.Name())
			if err != nil {
				return classification{}, err
			}
		}
		if declares {
			c.components = append(c.components, name)
		} else {
			c.helpers = append(c.helpers, name)
		}
	}
	sort.Strings(c.components)
	sort.Strings(c.helpers)
	return c, nil
})

type classification struct {
	components  []string
	helpers     []string
	directories map[string]bool
}

// isDirectoryPackage reports whether name vendors as its own package under
// ui/<name>/ rather than as a file in the flat ui package.
func isDirectoryPackage(name string) bool {
	c, err := classified()
	if err != nil {
		return false
	}
	return c.directories[name]
}

func fileDeclaresComponent(fileName string) (bool, error) {
	src, err := fs.ReadFile(gsxui.Files, "ui/"+fileName)
	if err != nil {
		return false, err
	}
	file, err := gsxparser.ParseFile(token.NewFileSet(), fileName, src, 0)
	if err != nil {
		return false, fmt.Errorf("parse ui/%s: %w", fileName, err)
	}
	return declaresComponent(file.Decls), nil
}

func directoryDeclaresComponent(dir string) (bool, error) {
	entries, err := fs.ReadDir(gsxui.Files, "ui/"+dir)
	if err != nil {
		return false, err
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".gsx") {
			continue
		}
		declares, err := fileDeclaresComponent(dir + "/" + e.Name())
		if err != nil {
			return false, err
		}
		if declares {
			return true, nil
		}
	}
	return false, nil
}

func declaresComponent(decls []gsxast.Decl) bool {
	for _, d := range decls {
		if _, ok := d.(*gsxast.Component); ok {
			return true
		}
	}
	return false
}

// Components lists what `gsxui add <name>` is documented to accept and what
// the site renders a page for. Helpers are excluded: they vendor only as
// dependencies (see Helpers).
func Components() ([]string, error) {
	c, err := classified()
	if err != nil {
		return nil, err
	}
	return slices.Clone(c.components), nil
}

// Helpers lists the Go-only entries under ui/. They have no component,
// no examples and no style recipe; they exist to be depended on.
func Helpers() ([]string, error) {
	c, err := classified()
	if err != nil {
		return nil, err
	}
	return slices.Clone(c.helpers), nil
}

// isVendorable reports whether name is a component or a helper — excluding
// non-.gsx files like index.js/gsxui.js by construction.
func isVendorable(name string) bool {
	c, err := classified()
	if err != nil {
		return false
	}
	return slices.Contains(c.components, name) || slices.Contains(c.helpers, name)
}

// parseX parses component name's committed generated source.
func parseX(name string) (*ast.File, error) {
	src, err := fs.ReadFile(gsxui.Files, "ui/"+name+".x.go")
	if err != nil {
		return nil, err
	}
	return parser.ParseFile(token.NewFileSet(), name+".x.go", src, parser.SkipObjectResolution)
}

// declIndex maps every top-level identifier declared in a flat component's
// .x.go to the component (file basename) that declares it. Exported and
// unexported names are both indexed: ui/ is one flat package, so gsx's
// generated code freely calls other components' unexported render helpers
// (e.g. dialog.x.go invokes button.x.go's _gsxrenderButton, not Button) —
// and since all these files compile together as a single package, two
// components cannot legally declare the same top-level unexported name, so
// the index stays injective.
func declIndex() (map[string]string, error) {
	components, err := Components()
	if err != nil {
		return nil, err
	}
	helpers, err := Helpers()
	if err != nil {
		return nil, err
	}
	comps := append(slices.Clone(components), helpers...)
	idx := map[string]string{}
	for _, c := range comps {
		if isDirectoryPackage(c) {
			continue
		}
		f, err := parseX(c)
		if err != nil {
			return nil, err
		}
		for _, decl := range f.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				if d.Recv == nil {
					idx[d.Name.Name] = c
				}
			case *ast.GenDecl:
				for _, spec := range d.Specs {
					switch s := spec.(type) {
					case *ast.TypeSpec:
						idx[s.Name.Name] = c
					case *ast.ValueSpec:
						for _, n := range s.Names {
							idx[n.Name] = c
						}
					}
				}
			}
		}
	}
	return idx, nil
}

func Deps(name string) ([]string, error) {
	if !isVendorable(name) {
		return nil, fmt.Errorf("unknown component %q (run 'gsxui list')", name)
	}
	if isDirectoryPackage(name) {
		return nil, nil
	}
	seen := map[string]bool{}
	var deps []string
	add := func(dep string) {
		if dep != name && !seen[dep] {
			seen[dep] = true
			deps = append(deps, dep)
		}
	}
	src, err := fs.ReadFile(gsxui.Files, "ui/"+name+".gsx")
	if err != nil {
		return nil, err
	}
	for _, m := range uiSubpackageImportRe.FindAllSubmatch(src, -1) {
		if sub := string(m[1]); isDirectoryPackage(sub) {
			add(sub)
		}
	}
	idx, err := declIndex()
	if err != nil {
		return nil, err
	}
	f, err := parseX(name)
	if err != nil {
		return nil, err
	}
	ast.Inspect(f, func(n ast.Node) bool {
		if sel, ok := n.(*ast.SelectorExpr); ok {
			// Only the base of a selector can be a package-level component
			// ident; the .Sel side (icon.New, sb.WriteString) never is.
			ast.Inspect(sel.X, func(m ast.Node) bool {
				if id, ok := m.(*ast.Ident); ok {
					if owner, ok := idx[id.Name]; ok {
						add(owner)
					}
				}
				return true
			})
			return false
		}
		if id, ok := n.(*ast.Ident); ok {
			if owner, ok := idx[id.Name]; ok {
				add(owner)
			}
		}
		return true
	})
	sort.Strings(deps)
	return deps, nil
}

// HasJS reports whether name is a component with companion behavior JS.
// The isVendorable guard matters: ui/gsxui.js and ui/index.js are real
// files under ui/ but aren't any component's behavior JS.
func HasJS(name string) bool {
	if !isVendorable(name) {
		return false
	}
	_, err := fs.Stat(gsxui.Files, "ui/"+name+".js")
	return err == nil
}

func Resolve(names []string) ([]string, error) {
	seen := map[string]bool{}
	var visit func(string) error
	visit = func(n string) error {
		if seen[n] {
			return nil
		}
		deps, err := Deps(n)
		if err != nil {
			return err
		}
		seen[n] = true
		for _, d := range deps {
			if err := visit(d); err != nil {
				return err
			}
		}
		return nil
	}
	for _, n := range names {
		if err := visit(n); err != nil {
			return nil, err
		}
	}
	out := make([]string, 0, len(seen))
	for n := range seen {
		out = append(out, n)
	}
	sort.Strings(out)
	return out, nil
}
