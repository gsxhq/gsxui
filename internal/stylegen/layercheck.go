package stylegen

import (
	"bytes"
	"errors"
	"fmt"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	gsxast "github.com/gsxhq/gsx/ast"
	gsxparser "github.com/gsxhq/gsx/parser"
	parse "github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/css"

	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/merge"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// slotMarkerPrefix is the attribute-name prefix every slot marker carries.
const slotMarkerPrefix = "data-gsxui-slot-"

// designSpecRef points at the invariant this file enforces.
const designSpecRef = "docs/superpowers/specs/2026-07-29-typed-recipe-model-design.md §9"

// ComponentComposedMarkers returns every data-gsxui-slot-* marker on an element
// that renders through a migrated component. Those markers are the ones whose
// presentation now comes from the utilities layer, so a components-layer rule
// against them is dead.
func ComponentComposedMarkers(root string) ([]string, error) {
	byMarker, err := composedMarkers(root)
	if err != nil {
		return nil, err
	}
	markers := make([]string, 0, len(byMarker))
	for marker := range byMarker {
		markers = append(markers, marker)
	}
	sort.Strings(markers)
	return markers, nil
}

// composedMarkers maps each composed marker to the migrated component whose
// compiled utilities land on the same element. The component name is what the
// diagnostic needs in order to say whose utilities are winning.
func composedMarkers(root string) (map[string]string, error) {
	byMarker := map[string]string{}
	// migrated maps a GSX component tag to the canonical component it renders.
	// The tag is the exported name — matching case-insensitively would make the
	// plain HTML <button> element look like <Button>, which it is not.
	migrated := map[string]string{}
	for name, shape := range shapes.All() {
		migrated[accessorName(name)] = name // "button" matches <ui.Button> and <Button>
		// A migrated component's own slot markers sit on the very elements that
		// carry its compiled utilities, so they are composed by definition.
		for _, slot := range shape.Slots {
			byMarker[slotMarker(name, slot.Name)] = name
		}
	}

	paths, err := filepath.Glob(filepath.Join(root, "ui", "*.gsx"))
	if err != nil {
		return nil, err
	}
	sort.Strings(paths)

	for _, path := range paths {
		src, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		file, err := gsxparser.ParseFile(token.NewFileSet(), path, src, 0)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", path, err)
		}
		var walkErr error
		gsxast.Inspect(file, func(node gsxast.Node) bool {
			element, ok := node.(*gsxast.Element)
			if !ok {
				return true
			}
			component, ok := migrated[strings.TrimPrefix(element.Tag, "ui.")]
			if !ok {
				return true
			}
			for _, attr := range element.Attrs {
				marker := attrName(attr)
				if !strings.HasPrefix(marker, slotMarkerPrefix) {
					continue
				}
				if owner, dup := byMarker[marker]; dup && owner != component {
					walkErr = fmt.Errorf(
						"%s: marker %s renders through both %q and %q",
						path, marker, owner, component)
					return false
				}
				byMarker[marker] = component
			}
			return true
		})
		if walkErr != nil {
			return nil, walkErr
		}
	}
	return byMarker, nil
}

// slotMarker names the marker attribute a component renders for one of its
// declared slots. The root slot carries the bare component marker.
func slotMarker(component, slot string) string {
	if slot == "" {
		return slotMarkerPrefix + component
	}
	return slotMarkerPrefix + component + "-" + slot
}

// attrName reports the attribute name of any named GSX attribute. Spread and
// markup attributes carry no marker name, so they report "".
func attrName(attr gsxast.Attr) string {
	switch a := attr.(type) {
	case *gsxast.StaticAttr:
		return a.Name
	case *gsxast.BoolAttr:
		return a.Name
	case *gsxast.ExprAttr:
		return a.Name
	case *gsxast.ClassAttr:
		return a.Name
	case *gsxast.EmbeddedAttr:
		return a.Name
	case *gsxast.MarkupAttr:
		return a.Name
	default:
		return ""
	}
}

// CheckLayerPrecedence enforces the layer invariant of the design spec §9:
//
//	A rule overriding compiled component presentation must live in
//	@layer utilities AND carry specificity >= (0,1,0).
//
// Both halves fail silently on their own — the rule parses, generation and
// make audit stay green, and the override simply never applies in the browser.
// Every violation is reported, not just the first.
func CheckLayerPrecedence(root string) error {
	markers, err := composedMarkers(root)
	if err != nil {
		return err
	}
	utilities, err := componentUtilities(root, markers)
	if err != nil {
		return err
	}

	path := filepath.Join(root, "assets", "css", "styles", "default.css")
	src, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	rules, err := layeredRules(path, src)
	if err != nil {
		return err
	}

	var violations []string
	for _, rule := range rules {
		violations = append(violations, rule.violations(path, markers, utilities)...)
	}
	if len(violations) == 0 {
		return nil
	}
	return fmt.Errorf(
		"%d rule(s) override compiled component presentation but cannot win the cascade:\n\n%s",
		len(violations), strings.Join(violations, "\n\n"))
}

// componentUtilities collects, per migrated component that owns at least one
// composed marker, every utility that component can render. Two sets are kept
// because a component's `[&_svg…]:` utilities land on its descendants, not on
// the marker element itself.
type utilitySets struct {
	own        []string
	descendant []string
}

func componentUtilities(root string, markers map[string]string) (map[string]utilitySets, error) {
	wanted := map[string]struct{}{}
	for _, component := range markers {
		wanted[component] = struct{}{}
	}
	sets := make(map[string]utilitySets, len(wanted))
	for component := range wanted {
		path := filepath.Join(root, "registry", "styles", DefaultStyle, component+".css")
		src, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read %s %s recipe: %w", DefaultStyle, component, err)
		}
		style, err := recipe.ParseStyle(path, src)
		if err != nil {
			return nil, fmt.Errorf("parse %s %s recipe: %w", DefaultStyle, component, err)
		}
		var set utilitySets
		seen := map[string]struct{}{}
		for _, rule := range style.Rules() {
			for _, utility := range rule.Utilities {
				if _, dup := seen[utility]; dup {
					continue
				}
				seen[utility] = struct{}{}
				if inner, ok := descendantUtility(utility); ok {
					set.descendant = append(set.descendant, inner)
					continue
				}
				set.own = append(set.own, utility)
			}
		}
		sort.Strings(set.own)
		sort.Strings(set.descendant)
		sets[component] = set
	}
	return sets, nil
}

// descendantUtility unwraps Tailwind's arbitrary-descendant variant, so
// `[&_svg:not([class*='size-'])]:size-4` is compared as `size-4` against a rule
// that selects a descendant of the marker.
func descendantUtility(utility string) (string, bool) {
	if !strings.HasPrefix(utility, "[&") {
		return "", false
	}
	depth := 0
	for i, r := range utility {
		switch r {
		case '[':
			depth++
		case ']':
			depth--
			if depth == 0 {
				rest := utility[i+1:]
				if strings.HasPrefix(rest, ":") && len(rest) > 1 {
					return rest[1:], true
				}
				return "", false
			}
		}
	}
	return "", false
}

// layerRule is one ruleset from default.css together with the cascade layer it
// resolved in and the utilities it applies.
type layerRule struct {
	layer     string
	selector  string
	line      int
	utilities []string
}

// layeredRules parses default.css and returns every ruleset that carries at
// least one @apply, tagged with its enclosing @layer.
func layeredRules(filename string, src []byte) ([]layerRule, error) {
	parser := css.NewParser(parse.NewInputBytes(src), false)

	type frame struct {
		layer    string
		selector string
		line     int
		rule     bool
	}
	var stack []frame
	currentLayer := func() string {
		for i := len(stack) - 1; i >= 0; i-- {
			if stack[i].layer != "" {
				return stack[i].layer
			}
		}
		return ""
	}

	var rules []layerRule
	for {
		grammar, _, data := parser.Next()
		switch grammar {
		case css.ErrorGrammar:
			if err := parser.Err(); errors.Is(err, io.EOF) {
				return rules, nil
			} else if err != nil {
				return nil, fmt.Errorf("%s: %w", filename, err)
			}
			return nil, fmt.Errorf("%s: malformed CSS at offset %d", filename, parser.Offset())

		case css.BeginAtRuleGrammar:
			f := frame{}
			if bytes.Equal(data, []byte("@layer")) {
				f.layer = strings.TrimSpace(recipe.SelectorText(parser.Values()))
			}
			stack = append(stack, f)

		case css.BeginRulesetGrammar:
			stack = append(stack, frame{
				selector: recipe.SelectorText(parser.Values()),
				line:     lineAt(src, parser.Offset()),
				rule:     true,
			})

		case css.EndAtRuleGrammar, css.EndRulesetGrammar:
			if len(stack) == 0 {
				return nil, fmt.Errorf("%s: unbalanced block at offset %d", filename, parser.Offset())
			}
			stack = stack[:len(stack)-1]

		case css.AtRuleGrammar:
			if !bytes.Equal(data, []byte("@apply")) || len(stack) == 0 {
				continue
			}
			top := stack[len(stack)-1]
			if !top.rule {
				continue
			}
			applied := strings.Fields(recipe.SelectorText(parser.Values()))
			if len(applied) == 0 {
				continue
			}
			rules = append(rules, layerRule{
				layer:     currentLayer(),
				selector:  top.selector,
				line:      top.line,
				utilities: applied,
			})
		}
	}
}

func lineAt(src []byte, offset int) int {
	if offset > len(src) {
		offset = len(src)
	}
	return bytes.Count(src[:offset], []byte("\n")) + 1
}

// violations reports every way this rule loses to compiled component
// presentation. A rule can name several markers via a selector list, so each
// complex selector is judged on its own.
func (r layerRule) violations(filename string, markers map[string]string, sets map[string]utilitySets) []string {
	var out []string
	for _, complex := range splitSelectorList(r.selector) {
		_, component, descendant, ok := composedTarget(complex, markers)
		if !ok {
			continue
		}
		set := sets[component]
		against := set.own
		if descendant {
			against = set.descendant
		}
		contested := contestedUtilities(r.utilities, against)
		if len(contested) == 0 {
			continue
		}
		spec := selectorSpecificity(complex)
		switch {
		case r.layer != "utilities":
			layer := r.layer
			if layer == "" {
				layer = "no layer"
			}
			out = append(out, fmt.Sprintf(
				"%s:%d: %s applies %s in @layer %s, but that element renders through\n"+
					"  %s, whose own utilities win the layer ordering. Move this rule to\n"+
					"  @layer utilities and give it specificity >= (0,1,0) — see %s.",
				filepath.ToSlash(filename), r.line, strings.TrimSpace(complex),
				strings.Join(contested, " "), layer, component, designSpecRef))
		case !spec.beatsPlainUtility():
			out = append(out, fmt.Sprintf(
				"%s:%d: %s applies %s in @layer utilities, but its specificity is\n"+
					"  %s — inside one layer the cascade falls back to specificity, and a plain\n"+
					"  utility class from %s scores (0,1,0). Drop the :where() wrapper so the\n"+
					"  selector carries specificity >= (0,1,0) — see %s.",
				filepath.ToSlash(filename), r.line, strings.TrimSpace(complex),
				strings.Join(contested, " "), spec, component, designSpecRef))
		}
	}
	return out
}

// contestedUtilities reports which of a rule's utilities collide with one the
// component itself renders. Collision is decided by the repo's own Tailwind
// merger — the same authority that decides which class wins everywhere else —
// so no property table has to be maintained here. Identical utilities are not
// contested: restating a value the component already sets changes nothing.
func contestedUtilities(applied, componentUtilities []string) []string {
	var contested []string
	for _, utility := range applied {
		for _, own := range componentUtilities {
			if own == utility {
				continue
			}
			if merge.Merge([]string{own, utility}) == utility {
				contested = append(contested, utility)
				break
			}
		}
	}
	return contested
}

// composedTarget reports the composed marker a complex selector targets, and
// whether the selected element is the marker itself or a descendant of it.
func composedTarget(complex string, markers map[string]string) (marker, component string, descendant, ok bool) {
	compounds := splitCompounds(complex)
	for i := len(compounds) - 1; i >= 0; i-- {
		for name, owner := range markers {
			if !containsAttributeName(compounds[i], name) {
				continue
			}
			return name, owner, i != len(compounds)-1, true
		}
	}
	return "", "", false, false
}

// containsAttributeName reports whether selector names the given attribute,
// requiring a full-token match so data-gsxui-slot-card never matches
// data-gsxui-slot-card-header.
func containsAttributeName(selector, name string) bool {
	for offset := 0; ; {
		index := strings.Index(selector[offset:], name)
		if index < 0 {
			return false
		}
		index += offset
		after := index + len(name)
		if after >= len(selector) || !isAttributeNameByte(selector[after]) {
			if index == 0 || !isAttributeNameByte(selector[index-1]) {
				return true
			}
		}
		offset = index + 1
	}
}

func isAttributeNameByte(b byte) bool {
	return b == '-' || b == '_' ||
		(b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9')
}

// splitSelectorList splits a selector list on its top-level commas, ignoring
// commas nested inside :where(...)/:is(...) argument lists and attribute values.
func splitSelectorList(selector string) []string {
	var parts []string
	depth := 0
	start := 0
	var quote byte
	for i := 0; i < len(selector); i++ {
		c := selector[i]
		switch {
		case quote != 0:
			if c == quote {
				quote = 0
			}
		case c == '"' || c == '\'':
			quote = c
		case c == '(' || c == '[':
			depth++
		case c == ')' || c == ']':
			depth--
		case c == ',' && depth == 0:
			parts = append(parts, selector[start:i])
			start = i + 1
		}
	}
	parts = append(parts, selector[start:])
	out := parts[:0]
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

// splitCompounds splits one complex selector into its compound selectors,
// dropping the combinators between them. The last compound is the subject.
func splitCompounds(complex string) []string {
	var compounds []string
	depth := 0
	start := 0
	var quote byte
	flush := func(end int) {
		if trimmed := strings.TrimSpace(complex[start:end]); trimmed != "" {
			compounds = append(compounds, trimmed)
		}
	}
	for i := 0; i < len(complex); i++ {
		c := complex[i]
		switch {
		case quote != 0:
			if c == quote {
				quote = 0
			}
		case c == '"' || c == '\'':
			quote = c
		case c == '(' || c == '[':
			depth++
		case c == ')' || c == ']':
			depth--
		case depth == 0 && (c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '>' || c == '+' || c == '~'):
			flush(i)
			start = i + 1
		}
	}
	flush(len(complex))
	if len(compounds) == 0 {
		return []string{strings.TrimSpace(complex)}
	}
	return compounds
}

// specificity is a CSS specificity triple (ids, classes, types).
type specificity struct {
	a, b, c int
}

func (s specificity) String() string { return fmt.Sprintf("(%d,%d,%d)", s.a, s.b, s.c) }

func (s specificity) less(other specificity) bool {
	if s.a != other.a {
		return s.a < other.a
	}
	if s.b != other.b {
		return s.b < other.b
	}
	return s.c < other.c
}

func (s specificity) add(other specificity) specificity {
	return specificity{s.a + other.a, s.b + other.b, s.c + other.c}
}

// beatsPlainUtility reports whether the selector can out-rank a plain Tailwind
// utility class, which scores exactly (0,1,0).
func (s specificity) beatsPlainUtility() bool {
	return !s.less(specificity{0, 1, 0})
}

// selectorSpecificity computes a complex selector's specificity per CSS
// Selectors Level 4: :where() contributes nothing, and :is()/:not()/:has()
// contribute the specificity of their most specific argument.
func selectorSpecificity(selector string) specificity {
	var total specificity
	for i := 0; i < len(selector); {
		switch c := selector[i]; {
		case c == '#':
			i = skipIdent(selector, i+1)
			total.a++
		case c == '.':
			i = skipIdent(selector, i+1)
			total.b++
		case c == '[':
			i = skipBalanced(selector, i, '[', ']')
			total.b++
		case c == ':':
			i, total = pseudoSpecificity(selector, i, total)
		case c == '*':
			i++
		case isIdentStartByte(c):
			i = skipIdent(selector, i)
			total.c++
		default:
			i++
		}
	}
	return total
}

// pseudoSpecificity consumes one pseudo-class or pseudo-element starting at the
// ':' in index and folds its contribution into total.
func pseudoSpecificity(selector string, index int, total specificity) (int, specificity) {
	i := index + 1
	element := false
	if i < len(selector) && selector[i] == ':' {
		element = true
		i++
	}
	nameStart := i
	i = skipIdent(selector, i)
	name := strings.ToLower(selector[nameStart:i])
	if i < len(selector) && selector[i] == '(' {
		end := skipBalanced(selector, i, '(', ')')
		args := selector[i+1 : max(end-1, i+1)]
		i = end
		switch name {
		case "where":
			return i, total
		case "is", "not", "has", "matches", "any":
			var best specificity
			for _, arg := range splitSelectorList(args) {
				if candidate := selectorSpecificity(arg); best.less(candidate) {
					best = candidate
				}
			}
			return i, total.add(best)
		}
	}
	if element {
		total.c++
	} else {
		total.b++
	}
	return i, total
}

func isIdentStartByte(b byte) bool {
	return b == '-' || b == '_' || b == '\\' || b >= 0x80 ||
		(b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

func skipIdent(s string, i int) int {
	for i < len(s) && isAttributeNameByte(s[i]) {
		i++
	}
	return i
}

// skipBalanced returns the index just past the bracket that closes the one at
// open, ignoring brackets inside strings.
func skipBalanced(s string, open int, left, right byte) int {
	depth := 0
	var quote byte
	for i := open; i < len(s); i++ {
		c := s[i]
		switch {
		case quote != 0:
			if c == quote {
				quote = 0
			}
		case c == '"' || c == '\'':
			quote = c
		case c == left:
			depth++
		case c == right:
			depth--
			if depth == 0 {
				return i + 1
			}
		}
	}
	return len(s)
}
