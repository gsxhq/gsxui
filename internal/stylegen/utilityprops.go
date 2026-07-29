package stylegen

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// A components-layer rule that overrides compiled presentation with a PLAIN
// DECLARATION — `display: flex`, not `@apply flex` — cannot be judged by
// comparing class names: there is no class. It has to be judged on the CSS
// property it sets, against the properties the component's own utilities set.
//
// Nothing in Go knows what properties `size-4` or `has-[>svg]:px-2` expand to;
// only Tailwind does, and only against this repo's theme. So the utilities are
// compiled — once per check, one probe class per utility — and the property
// names are read straight out of Tailwind's output. Deriving them from a
// hand-written property table instead would be exactly the silent-false-
// negative machine this gate exists to replace.
//
// The compile is LAZY. It happens only when a rule that actually targets a
// composed marker carries plain declarations, so a tree whose foundation rules
// all name unmigrated components never shells out at all.

// probeClassPrefix names the throwaway class each probed utility is applied to.
const probeClassPrefix = "gsxui-layercheck-probe-"

var probeClassPattern = regexp.MustCompile(`\.` + regexp.QuoteMeta(probeClassPrefix) + `(\d+)\b`)

// utilityPropertyResolver answers "which CSS properties does this set of
// utilities declare?", compiling on first use and caching for the rest of the
// check.
type utilityPropertyResolver struct {
	root string
	sets map[string]utilitySets

	byUtility map[string]map[string]struct{}
	err       error
	done      bool
}

// propertiesOf returns the union of the properties declared by every utility in
// the set.
func (r *utilityPropertyResolver) propertiesOf(utilities []string) (map[string]struct{}, error) {
	if err := r.resolve(); err != nil {
		return nil, err
	}
	union := map[string]struct{}{}
	for _, utility := range utilities {
		for property := range r.byUtility[utility] {
			union[property] = struct{}{}
		}
	}
	return union, nil
}

func (r *utilityPropertyResolver) resolve() error {
	if r.done {
		return r.err
	}
	r.done = true
	seen := map[string]struct{}{}
	var all []string
	for _, set := range r.sets {
		for _, utility := range append(append([]string{}, set.own...), set.descendant...) {
			if _, dup := seen[utility]; dup {
				continue
			}
			seen[utility] = struct{}{}
			all = append(all, utility)
		}
	}
	sort.Strings(all)
	r.byUtility, r.err = compileUtilityProperties(r.root, all)
	return r.err
}

// compileUtilityProperties runs Tailwind over one probe class per utility and
// reads back the properties each one declares.
func compileUtilityProperties(root string, utilities []string) (map[string]map[string]struct{}, error) {
	byUtility := make(map[string]map[string]struct{}, len(utilities))
	for _, utility := range utilities {
		byUtility[utility] = map[string]struct{}{}
	}
	if len(utilities) == 0 {
		return byUtility, nil
	}

	binary := filepath.Join(root, "node_modules", ".bin", "tailwindcss")
	if _, err := os.Stat(binary); err != nil {
		return nil, fmt.Errorf(
			"the layer gate must resolve compiled utilities to CSS properties, which needs the "+
				"Tailwind CLI at %s — run `npm install`: %w", binary, err)
	}

	// The probe lives inside root so that `@import \"tailwindcss\"` resolves
	// through root/node_modules, and it imports foundation.css because that is
	// where this repo's @theme lives — without it `bg-primary` is unknown.
	dir, err := os.MkdirTemp(root, ".layercheck-probe-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)

	var probe strings.Builder
	probe.WriteString("@import \"tailwindcss\" source(none);\n")
	probe.WriteString("@import \"tw-animate-css\";\n")
	fmt.Fprintf(&probe, "@import %q;\n", filepath.ToSlash(filepath.Join(root, "assets", "css", "foundation.css")))
	for index, utility := range utilities {
		fmt.Fprintf(&probe, ".%s%d { @apply %s; }\n", probeClassPrefix, index, utility)
	}

	input := filepath.Join(dir, "probe.css")
	output := filepath.Join(dir, "probe.out.css")
	if err := os.WriteFile(input, []byte(probe.String()), 0o644); err != nil {
		return nil, err
	}
	command := exec.Command(binary, "-i", input, "-o", output)
	command.Dir = root
	if combined, err := command.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("compile component utilities with Tailwind: %w\n%s", err, combined)
	}
	compiled, err := os.ReadFile(output)
	if err != nil {
		return nil, err
	}

	found, err := probedProperties(output, compiled)
	if err != nil {
		return nil, err
	}
	for index, properties := range found {
		if index < 0 || index >= len(utilities) {
			continue
		}
		for property := range properties {
			byUtility[utilities[index]][property] = struct{}{}
		}
	}
	for _, utility := range utilities {
		if len(byUtility[utility]) == 0 {
			return nil, fmt.Errorf(
				"utility %q compiled to no declarations — the layer gate cannot tell which "+
					"properties it contests, so it would silently pass every plain-declaration "+
					"override of it", utility)
		}
	}
	return byUtility, nil
}

// probedProperties maps each probe index to the properties its compiled rules
// declare. Tailwind's output nests @supports INSIDE rulesets and wraps rulesets
// in @media, so attribution is the union of the probe classes named anywhere in
// the enclosing block headers. A CSS grammar parser is the wrong tool here —
// tdewolff's rejects a nested at-rule inside a qualified rule — and only block
// structure and property names are needed, so this scans for exactly those.
func probedProperties(filename string, src []byte) (map[int]map[string]struct{}, error) {
	out := map[int]map[string]struct{}{}
	var stack [][]int

	record := func(segment []byte) {
		declaration := strings.TrimSpace(string(segment))
		if declaration == "" || strings.HasPrefix(declaration, "@") {
			return
		}
		colon := strings.IndexByte(declaration, ':')
		if colon < 0 {
			return
		}
		property := strings.ToLower(strings.TrimSpace(declaration[:colon]))
		if property == "" || strings.ContainsAny(property, " \t\n(){}") {
			return
		}
		for _, indices := range stack {
			for _, index := range indices {
				if out[index] == nil {
					out[index] = map[string]struct{}{}
				}
				out[index][property] = struct{}{}
			}
		}
	}

	segment := 0
	for i := 0; i < len(src); i++ {
		switch c := src[i]; c {
		case '\'', '"':
			for i++; i < len(src); i++ {
				if src[i] == '\\' {
					i++
					continue
				}
				if src[i] == c {
					break
				}
			}
		case '/':
			if i+1 < len(src) && src[i+1] == '*' {
				if end := bytes.Index(src[i+2:], []byte("*/")); end >= 0 {
					i += 2 + end + 1
				} else {
					i = len(src)
				}
			}
		case '{':
			var indices []int
			for _, match := range probeClassPattern.FindAllSubmatch(src[segment:i], -1) {
				index, err := strconv.Atoi(string(match[1]))
				if err != nil {
					continue
				}
				indices = append(indices, index)
			}
			stack = append(stack, indices)
			segment = i + 1
		case '}':
			record(src[segment:i])
			if len(stack) == 0 {
				return nil, fmt.Errorf("%s: unbalanced block at offset %d", filename, i)
			}
			stack = stack[:len(stack)-1]
			segment = i + 1
		case ';':
			record(src[segment:i])
			segment = i + 1
		}
	}
	if len(stack) != 0 {
		return nil, fmt.Errorf("%s: %d unclosed block(s)", filename, len(stack))
	}
	return out, nil
}
