package stylegen

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// authoringPatternAll matches the hand-authored presentation markers that an
// UNMIGRATED component must never carry: data-slot attributes, bare group/ and
// peer/ variant roots, and the group-* / peer-* conditional variant forms
// (group-data/, peer-has/, group-focus/, and so on). This is the pattern set
// make audit used to apply to every ui/*.gsx file except button.gsx.
var authoringPatternAll = regexp.MustCompile(
	`data-slot|group/|peer/|(group|peer)-(data|has|focus|hover|active|disabled|aria|open|checked)[^[:space:]]*/`)

// authoringPatternNarrow is authoringPatternAll with the bare "group/" root
// removed. A migrated component's generated output legitimately carries
// group/<component> (e.g. group/button) so its own group can be targeted by
// its own descendants; that is not hand-authored presentation, it is compiled
// structure. Everything else — data-slot, peer/, and the group-*/peer-*
// conditional variant forms — remains disallowed even for generated output,
// matching what make audit asserted for ui/button.gsx specifically.
var authoringPatternNarrow = regexp.MustCompile(
	`data-slot|peer/|(group|peer)-(data|has|focus|hover|active|disabled|aria|open|checked)[^[:space:]]*/`)

// authoringClassIndent matches a line that is (after leading whitespace)
// exactly a class= attribute continuation, e.g. a class="..." line inside a
// multi-line tag.
var authoringClassIndent = regexp.MustCompile(`^[[:space:]]+class=`)

// authoringClassTag matches an opening tag line that carries a class=
// attribute directly.
var authoringClassTag = regexp.MustCompile(`^[[:space:]]*<[^>]*class=`)

// CheckAuthoring enforces the split make audit used to draw between migrated
// and unmigrated components, but reads the migrated set from
// registry/canonical/shapes.All() instead of a hand-maintained filename
// allowlist in the Makefile:
//
//   - An UNMIGRATED component's ui/<c>.gsx must contain no hand-authored
//     presentation: no data-slot, group/, peer/, group-*/peer-* variant forms,
//     and no class= attribute.
//   - A MIGRATED component's ui/<c>.gsx is generated output. It is expected to
//     carry class= attributes and its own group/<c> root, so those are not
//     flagged. It must still carry no data-slot, peer/, or group-*/peer-*
//     conditional variant forms, and — because it is supposed to be generated
//     output — it must be byte-identical to
//     registry/generated/<DefaultStyle>/<c>.gsx.
//
// The dev/ directory, when present, gets the same unmigrated pattern check
// that ui/ used to apply via $(wildcard dev) — dev is scratch/prototype code,
// never generated output, so the class= check and the generated-output
// comparison do not apply there.
func CheckAuthoring(root string) error {
	var violations []string

	uiViolations, err := checkAuthoringDir(root, filepath.Join(root, "ui"), true)
	if err != nil {
		return err
	}
	violations = append(violations, uiViolations...)

	devDir := filepath.Join(root, "dev")
	if info, err := os.Stat(devDir); err == nil && info.IsDir() {
		devViolations, err := checkAuthoringDir(root, devDir, false)
		if err != nil {
			return err
		}
		violations = append(violations, devViolations...)
	}

	if len(violations) == 0 {
		return nil
	}
	sort.Strings(violations)
	return fmt.Errorf("%d authoring violation(s):\n\n%s", len(violations), strings.Join(violations, "\n"))
}

// checkAuthoringDir walks every .gsx file under dir and applies the migrated
// or unmigrated rule set depending on whether its component name is declared
// in shapes.All(). checkClass gates whether the class= rules apply to this
// directory (they only ever applied to ui/, never to dev/).
func checkAuthoringDir(root, dir string, checkClass bool) ([]string, error) {
	migrated := shapes.All()
	var violations []string

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Ext(path) != ".gsx" {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		name := strings.TrimSuffix(filepath.Base(path), ".gsx")
		_, isMigrated := migrated[name]

		if isMigrated {
			genPath := filepath.Join(root, "registry", "generated", DefaultStyle, name+".gsx")
			genContent, err := os.ReadFile(genPath)
			if err != nil {
				violations = append(violations, fmt.Sprintf(
					"%s: component %q is declared in shapes.All() but has no generated counterpart at %s (%v)",
					rel, name, filepath.Join("registry", "generated", DefaultStyle, name+".gsx"), err))
			} else if !bytes.Equal(content, genContent) {
				violations = append(violations, fmt.Sprintf(
					"%s: migrated component %q must be byte-identical to %s (run go run ./cmd/stylegen)",
					rel, name, filepath.Join("registry", "generated", DefaultStyle, name+".gsx")))
			}
			violations = append(violations, scanAuthoringLines(rel, content, authoringPatternNarrow)...)
			return nil
		}

		violations = append(violations, scanAuthoringLines(rel, content, authoringPatternAll)...)
		if checkClass {
			violations = append(violations, scanAuthoringLines(rel, content, authoringClassIndent)...)
			violations = append(violations, scanAuthoringLines(rel, content, authoringClassTag)...)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return violations, nil
}

// scanAuthoringLines reports every line in content matching pattern, in the
// same 1-indexed, per-line style rg -n uses.
func scanAuthoringLines(rel string, content []byte, pattern *regexp.Regexp) []string {
	var out []string
	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		if pattern.MatchString(line) {
			out = append(out, fmt.Sprintf("%s:%d: %s", rel, i+1, strings.TrimSpace(line)))
		}
	}
	return out
}
