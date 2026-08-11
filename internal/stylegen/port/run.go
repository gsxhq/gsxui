package port

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// AllStyles are the 8 upstream shadcn-ui styles, exactly (plan's Global
// Constraints).
var AllStyles = []string{"vega", "nova", "maia", "lyra", "mira", "luma", "sera", "rhea"}

// Report summarizes one Run. Written is the number of canonical components
// successfully ported (and, unless dryRun, written to disk) — for a single
// style this is always len(shapes.All()) when the mapping table is complete,
// since every component gets a recipe, including the 4 StyleInvariant ones
// (upstream has no section for them; Transform's fallback policy carries 100%
// of their content from nova, exactly reproducing what a hand-authored
// identical-to-nova recipe already does). Carried counts individual carried
// rules (informational, not a gate). Unmapped counts every upstream class or
// token Transform could not place anywhere in the shape — the CLI gate: a
// clean run has Unmapped == 0. Skipped counts upstream MARK sections that
// have no gsxui component at all (SectionComponent's ignoredSections table) —
// purely informational, distinct from Written/Unmapped's per-component
// accounting.
type Report struct {
	Written        int
	Carried        int
	Unmapped       int
	Skipped        int
	UnmappedDetail []string
}

// Run ports one style (or several) from upstreamRoot's shadcn-ui checkout
// into repoRoot's registry/styles/<style>/. dryRun computes the full report
// without writing anything to disk. Every component's rendered recipe is
// validated (recipe.ParseStyle -> recipe.Conform -> recipe.CheckConflicts,
// via Validate) before it is counted as Written or written to disk — a
// Transform/Render/Validate failure aborts the whole run immediately (dossier
// §3: "the porter's output must pass the repo's own gates before being
// written"), distinct from the Unmapped count, which is an expected,
// quantifiable condition the caller may choose to tolerate (--allow-unmapped
// at the CLI).
func Run(repoRoot, upstreamRoot string, styles []string, dryRun bool) (Report, error) {
	if len(styles) == 0 {
		return Report{}, fmt.Errorf("port: no styles requested")
	}

	for _, style := range styles {
		if !validStyle(style) {
			return Report{}, fmt.Errorf("port: unknown style %q (want one of %s)", style, strings.Join(AllStyles, ", "))
		}
	}

	upstreamSHA, err := gitHeadSHA(upstreamRoot)
	if err != nil {
		return Report{}, fmt.Errorf("port: resolve upstream commit: %w", err)
	}

	allShapes := shapes.All()
	components := make([]string, 0, len(allShapes))
	for name := range allShapes {
		components = append(components, name)
	}
	sort.Strings(components)

	// The fallback source (dossier §4's "fallback policy") is always today's
	// nova recipe, regardless of which style is being ported — read once per
	// component, reused across every requested style.
	novaByComponent := make(map[string]recipe.Style, len(components))
	for _, component := range components {
		novaPath := filepath.Join(repoRoot, "registry", "styles", "nova", component+".css")
		src, err := os.ReadFile(novaPath)
		if err != nil {
			return Report{}, fmt.Errorf("port: read nova fallback %s: %w", novaPath, err)
		}
		style, err := recipe.ParseStyle(novaPath, src)
		if err != nil {
			return Report{}, fmt.Errorf("port: parse nova fallback %s: %w", novaPath, err)
		}
		novaByComponent[component] = style
	}

	var report Report
	type write struct {
		path string
		data []byte
	}
	var pending []write

	for _, style := range styles {
		upstreamPath := filepath.Join(upstreamRoot, "apps", "v4", "registry", "styles", "style-"+style+".css")
		parsedStyleName, sections, err := ReadUpstream(upstreamPath)
		if err != nil {
			return Report{}, fmt.Errorf("port: read upstream %s: %w", upstreamPath, err)
		}
		if parsedStyleName != style {
			return Report{}, fmt.Errorf("port: %s declares style %q, want %q", upstreamPath, parsedStyleName, style)
		}
		upstreamRelPath := fmt.Sprintf("apps/v4/registry/styles/style-%s.css", style)

		sectionByComponent := make(map[string]Section, len(sections))
		for _, sec := range sections {
			component, ok := SectionComponent(sec.Name)
			if !ok {
				report.Skipped++
				continue
			}
			if existing, exists := sectionByComponent[component]; exists {
				// An empty MARK section (no rules before the next MARK
				// comment) never contributes anything, so it can never
				// really conflict with another section mapping to the same
				// component — confirmed against style-mira.css's own
				// duplicated, back-to-back "/* MARK: Sidebar */" comment
				// (the first with nothing between it and the second), an
				// upstream authoring artifact rather than two real
				// sections. Only two NON-empty sections mapping to the same
				// component is the real, still-fatal ambiguity.
				switch {
				case len(existing.Rules) == 0:
					sectionByComponent[component] = sec
				case len(sec.Rules) == 0:
				default:
					return Report{}, fmt.Errorf("port: %s: two MARK sections both map to component %q", upstreamPath, component)
				}
				continue
			}
			sectionByComponent[component] = sec
		}

		for _, component := range components {
			shape := allShapes[component]

			sec, hasSection := sectionByComponent[component]
			if !hasSection {
				if !StyleInvariant(component) && !missingSectionExpected(component, style) {
					return Report{}, fmt.Errorf(
						"port: %s: component %q has no upstream MARK section and is not StyleInvariant — extend the mapping table",
						upstreamPath, component)
				}
				sec = Section{} // no upstream section at all: every declared rule is carried from nova.
			}

			ported, unmapped, err := Transform(shape, sec, novaByComponent[component])
			if err != nil {
				return Report{}, fmt.Errorf("port: %s %s: %w", style, component, err)
			}

			outPath := filepath.Join(repoRoot, "registry", "styles", style, component+".css")
			rendered, err := Render(ported, style, upstreamSHA, upstreamRelPath)
			if err != nil {
				return Report{}, fmt.Errorf("port: render %s %s: %w", style, component, err)
			}
			if err := Validate(outPath, shape, rendered); err != nil {
				return Report{}, fmt.Errorf("port: validate %s %s: %w", style, component, err)
			}

			report.Written++
			report.Carried += len(ported.Carried)
			if len(unmapped) > 0 {
				report.Unmapped += len(unmapped)
				for _, u := range unmapped {
					report.UnmappedDetail = append(report.UnmappedDetail, fmt.Sprintf("%s %s: %s", style, component, u))
				}
			}

			pending = append(pending, write{path: outPath, data: rendered})
		}
	}

	if dryRun {
		return report, nil
	}
	for _, w := range pending {
		if err := os.WriteFile(w.path, w.data, 0o644); err != nil {
			return Report{}, fmt.Errorf("port: write %s: %w", w.path, err)
		}
	}
	return report, nil
}

func validStyle(style string) bool {
	for _, s := range AllStyles {
		if s == style {
			return true
		}
	}
	return false
}

// gitHeadSHA resolves the pinned commit the porter cites in every generated
// header (dossier §1.1: "Source: shadcn-ui/ui@<full-commit-sha>").
func gitHeadSHA(upstreamRoot string) (string, error) {
	cmd := exec.Command("git", "-C", upstreamRoot, "rev-parse", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git rev-parse HEAD in %s: %w", upstreamRoot, err)
	}
	return strings.TrimSpace(string(out)), nil
}
