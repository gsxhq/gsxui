// Package registry derives the vendorable component set from the embedded
// filesystem — the component list, inter-component dependencies (parsed from
// .gsx imports), and behavior-JS presence. Derived, never declared: it
// cannot drift from the code it describes.
package registry

import (
	"fmt"
	"io/fs"
	"regexp"
	"sort"

	gsxui "github.com/gsxhq/gsxui"
)

// importRe matches gsxui-internal imports in .gsx source, capturing the
// component package name.
var importRe = regexp.MustCompile(`"github\.com/gsxhq/gsxui/ui/([a-z]+)"`)

func Components() ([]string, error) {
	entries, err := fs.ReadDir(gsxui.Files, "ui")
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() && e.Name() != "core" {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	return names, nil
}

func Deps(name string) ([]string, error) {
	entries, err := fs.ReadDir(gsxui.Files, "ui/"+name)
	if err != nil {
		return nil, fmt.Errorf("unknown component %q (run 'gsxui list')", name)
	}
	seen := map[string]bool{}
	var deps []string
	for _, e := range entries {
		if e.IsDir() || !isGsx(e.Name()) {
			continue
		}
		src, err := fs.ReadFile(gsxui.Files, "ui/"+name+"/"+e.Name())
		if err != nil {
			return nil, err
		}
		for _, m := range importRe.FindAllStringSubmatch(string(src), -1) {
			if dep := m[1]; dep != name && !seen[dep] {
				seen[dep] = true
				deps = append(deps, dep)
			}
		}
	}
	sort.Strings(deps)
	return deps, nil
}

func HasJS(name string) bool {
	_, err := fs.Stat(gsxui.Files, "ui/"+name+"/"+name+".js")
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

func isGsx(name string) bool {
	const ext = ".gsx"
	return len(name) > len(ext) && name[len(name)-len(ext):] == ext
}
