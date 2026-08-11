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

// authoringPatternNarrow is what a MIGRATED component's generated ui/<c>.gsx
// is checked against: data-slot alone.
//
// This used to also forbid peer/ and the group-*/peer-* conditional variant
// forms, on the theory that only button.gsx — the one migrated component that
// existed when this check was written — would ever need to reference its own
// group/<c> root, and nothing migrated would need MORE than that bare root.
// The 8-style port's mechanical re-port of Switch, Tabs, Sidebar, Item and
// Field (registry/generated/*/{switch,tabs,sidebar,item,field}.gsx — see
// each one's own recipe under registry/styles/<style>/) showed that
// assumption was too narrow: their compiled output legitimately carries
// group-data-[…]/switch:, peer-hover/menu-button: and similar forms —
// descendants and siblings reacting to ANOTHER slot's state within the SAME
// component's own multi-slot structure, upstream's own composition idiom
// (e.g. Switch's thumb keying off group-data-[size=…]/switch:, matching
// apps/v4/registry/bases/base/ui/switch.tsx), not hand-authored presentation
// coupling one component to another's markup.
//
// The distinction that actually matters for a MIGRATED component is not
// "does this line use group/peer" — it is "was this line generated or
// typed by hand", and TestCheckAuthoringRejectsMigratedDivergingFromGenerated
// already enforces that a migrated ui/<c>.gsx is byte-identical to
// registry/generated/<DefaultStyle>/<c>.gsx. Nothing that check passes can be
// hand-authored, group/peer or otherwise. data-slot alone stays checked here
// because it is not a "was this hand-typed" question: it is upstream's OWN
// raw attribute name (gsxui's is data-gsxui-slot-*), and upstream's source
// CSS embeds it inside arbitrary-variant selector STRINGS the porter's
// class-name-based re-slotting does not parse into
// (`has-[[data-slot=input-group-control]:focus-visible]:...` in
// style-nova.css, for one) — a leak there would compile clean and pass byte-
// identity too, so it needs its own, permanent check.
var authoringPatternNarrow = regexp.MustCompile(`data-slot`)

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
//     carry class= attributes and whatever group/peer forms its own compiled
//     recipe legitimately needs (see authoringPatternNarrow), so only
//     data-slot remains checked — and, because it is supposed to be generated
//     output, it must be byte-identical to
//     registry/generated/<DefaultStyle>/<c>.gsx, which is what actually rules
//     out hand-authoring for everything else.
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
