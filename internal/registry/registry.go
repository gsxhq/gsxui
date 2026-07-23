// Package registry derives the vendorable component set from the embedded
// filesystem — the component list, inter-component dependencies, and
// behavior-JS presence. Derived, never declared: it cannot drift from the
// code it describes.
//
// ui/ is one flat package, so a component is a .gsx file basename:
// ui/button.gsx is component "button". ui/icon is the one directory
// component — it stays a package so icon.New reads as a name. Dependencies
// come from two sources: the icon import in .gsx source, and — because
// intra-package references have no import to scan — identifiers in a
// component's generated .x.go that another component's .x.go declares,
// resolved with go/parser against a declaration index.
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

	gsxui "github.com/gsxhq/gsxui"
)

var iconImportRe = regexp.MustCompile(`"github\.com/gsxhq/gsxui/ui/icon"`)

func Components() ([]string, error) {
	entries, err := fs.ReadDir(gsxui.Files, "ui")
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			if e.Name() == "icon" {
				names = append(names, e.Name())
			}
			continue
		}
		if name, ok := strings.CutSuffix(e.Name(), ".gsx"); ok {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return names, nil
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
	comps, err := Components()
	if err != nil {
		return nil, err
	}
	idx := map[string]string{}
	for _, c := range comps {
		if c == "icon" {
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
	if !isComponent(name) {
		return nil, fmt.Errorf("unknown component %q (run 'gsxui list')", name)
	}
	if name == "icon" {
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
	if iconImportRe.Match(src) {
		add("icon")
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
// The isComponent guard matters: ui/gsxui.js and ui/index.js are real files
// under ui/ but aren't any component's behavior JS.
func HasJS(name string) bool {
	if !isComponent(name) {
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

// isComponent reports whether name is a member of Components() — excluding
// non-component files like index.js/gsxui.js (not .gsx) by construction.
func isComponent(name string) bool {
	names, err := Components()
	if err != nil {
		return false
	}
	return slices.Contains(names, name)
}
