// Package port implements the shadcn-ui upstream style-file porter: reading
// upstream style-<name>.css files, classifying their Tailwind variant
// grammar, and (in later tasks) transforming them into gsxui recipe files.
package port

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"regexp"
	"strings"

	parse "github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/css"
)

// Rule is one upstream `.cn-*` rule.
type Rule struct {
	Class     string   // e.g. "cn-accordion-trigger" (no leading dot)
	Utilities []string // @apply tokens, in source order
	Line      int      // 1-based line of the selector, for provenance
}

// Section is one `/* MARK: <Name> */` block.
type Section struct {
	Name       string // e.g. "Radio Group"
	Rules      []Rule
	Start, End int // 1-based line range, for the provenance header
}

// ReadUpstream parses an upstream shadcn-ui style-<name>.css file — the
// single `.style-<name> { ... }` wrapper containing every component's
// `.cn-*` rules — into its style name and its MARK-delimited sections.
//
// Rules appearing before the first `/* MARK: ... */` comment belong to a
// synthetic section named "", which is omitted entirely when empty.
func ReadUpstream(path string) (style string, sections []Section, err error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return "", nil, fmt.Errorf("read upstream %s: %w", path, err)
	}

	marks, err := scanMarks(src)
	if err != nil {
		return "", nil, fmt.Errorf("scan MARK comments in %s: %w", path, err)
	}

	style, rules, err := parseWrapper(path, src)
	if err != nil {
		return "", nil, err
	}

	return style, groupSections(marks, rules), nil
}

// upstreamRule is a Rule plus the line its closing brace appeared on, used
// only to compute a Section's End line.
type upstreamRule struct {
	Rule
	endLine int
}

// parseWrapper walks the `.style-<name> { .cn-a {...} .cn-b {...} }`
// structure using the same tdewolff/parse/v2/css tokenizer idiom as
// internal/recipe/parse.go's ParseStyle.
func parseWrapper(filename string, src []byte) (string, []upstreamRule, error) {
	parser := css.NewParser(parse.NewInputBytes(src), false)

	var style string
	var haveWrapper, inWrapper bool
	var rules []upstreamRule
	var current *upstreamRule

	for {
		grammar, _, data := parser.Next()
		switch grammar {
		case css.ErrorGrammar:
			if err := parser.Err(); errors.Is(err, io.EOF) {
				if !haveWrapper {
					return "", nil, upstreamError(filename, src, len(src), "no .style-<name> wrapper found")
				}
				if inWrapper || current != nil {
					return "", nil, upstreamError(filename, src, len(src), "unterminated .style-%s wrapper", style)
				}
				return style, rules, nil
			} else if err != nil {
				return "", nil, upstreamError(filename, src, parser.Offset(), "malformed CSS: %v", err)
			}
			return "", nil, upstreamError(filename, src, parser.Offset(), "malformed CSS")

		case css.CommentGrammar:
			continue

		case css.BeginRulesetGrammar:
			if !haveWrapper {
				name, ok := wrapperSelector(parser.Values())
				if !ok {
					return "", nil, upstreamError(filename, src, parser.Offset(), "expected .style-<name> wrapper, found %s", selectorText(parser.Values()))
				}
				style = name
				haveWrapper = true
				inWrapper = true
				continue
			}
			if !inWrapper {
				return "", nil, upstreamError(filename, src, parser.Offset(), "unexpected rule %s outside .style-%s wrapper", selectorText(parser.Values()), style)
			}
			if current != nil {
				return "", nil, upstreamError(filename, src, parser.Offset(), "nested rule %s in %s", selectorText(parser.Values()), current.Class)
			}
			class, ok := cnSelector(parser.Values())
			if !ok {
				return "", nil, upstreamError(filename, src, parser.Offset(), "expected .cn-* rule inside .style-%s, found %s", style, selectorText(parser.Values()))
			}
			line, _, _ := parse.Position(bytes.NewReader(src), parser.Offset())
			current = &upstreamRule{Rule: Rule{Class: class, Line: line}}

		case css.AtRuleGrammar:
			if current == nil {
				return "", nil, upstreamError(filename, src, parser.Offset(), "unexpected at-rule outside a .cn-* rule")
			}
			if !bytes.Equal(data, []byte("@apply")) {
				return "", nil, upstreamError(filename, src, parser.Offset(), "expected @apply in %s, found %s", current.Class, data)
			}
			utilities, err := parseApply(parser.Values())
			if err != nil {
				return "", nil, upstreamError(filename, src, parser.Offset(), "invalid @apply in %s: %v", current.Class, err)
			}
			current.Utilities = append(current.Utilities, utilities...)

		case css.EndRulesetGrammar:
			if current != nil {
				line, _, _ := parse.Position(bytes.NewReader(src), parser.Offset())
				current.endLine = line
				rules = append(rules, *current)
				current = nil
				continue
			}
			if inWrapper {
				inWrapper = false
				continue
			}
			return "", nil, upstreamError(filename, src, parser.Offset(), "unexpected closing brace")

		case css.DeclarationGrammar, css.CustomPropertyGrammar:
			return "", nil, upstreamError(filename, src, parser.Offset(), "unexpected plain declaration")

		default:
			return "", nil, upstreamError(filename, src, parser.Offset(), "unsupported CSS grammar %s", grammar)
		}
	}
}

// mark is one `/* MARK: <Name> */` comment's name and 1-based line.
type mark struct {
	name string
	line int
}

var markPattern = regexp.MustCompile(`^/\*\s*MARK:\s*(.+?)\s*\*/$`)

// scanMarks finds every `/* MARK: <Name> */` comment in src, regardless of
// nesting depth. It uses the raw css.Lexer (not css.Parser): the parser's
// grammar stream only ever surfaces comments at the top of the stylesheet
// (nesting depth 0) and silently discards comments nested inside a ruleset
// — exactly where every MARK comment in an upstream style file lives, one
// level inside the `.style-<name> { ... }` wrapper — so scanning tokens
// directly is required to see them at all.
func scanMarks(src []byte) ([]mark, error) {
	lexer := css.NewLexer(parse.NewInputBytes(src))
	var marks []mark
	offset := 0
	for {
		tt, data := lexer.Next()
		if tt == css.ErrorToken {
			if err := lexer.Err(); !errors.Is(err, io.EOF) {
				return nil, err
			}
			return marks, nil
		}
		if tt == css.CommentToken {
			if m := markPattern.FindSubmatch(data); m != nil {
				line, _, _ := parse.Position(bytes.NewReader(src), offset)
				marks = append(marks, mark{name: string(m[1]), line: line})
			}
		}
		offset += len(data)
	}
}

// groupSections partitions rules (in source order) among marks (in source
// order) by comparing each rule's line against the marks' line boundaries.
func groupSections(marks []mark, rules []upstreamRule) []Section {
	var sections []Section

	ruleIdx := 0
	firstMarkLine := math.MaxInt
	if len(marks) > 0 {
		firstMarkLine = marks[0].line
	}
	var synthetic []upstreamRule
	for ruleIdx < len(rules) && rules[ruleIdx].Line < firstMarkLine {
		synthetic = append(synthetic, rules[ruleIdx])
		ruleIdx++
	}
	if len(synthetic) > 0 {
		sections = append(sections, buildSection("", synthetic, 0))
	}

	for i, m := range marks {
		next := math.MaxInt
		if i+1 < len(marks) {
			next = marks[i+1].line
		}
		var group []upstreamRule
		for ruleIdx < len(rules) && rules[ruleIdx].Line < next {
			group = append(group, rules[ruleIdx])
			ruleIdx++
		}
		sections = append(sections, buildSection(m.name, group, m.line))
	}

	return sections
}

// buildSection converts rules into a Section. When markLine is non-zero
// (a real MARK section, as opposed to the synthetic leading "" section) it
// wins over the group's own first-rule line as the Start, since the MARK
// comment itself is the more precise provenance anchor.
func buildSection(name string, rules []upstreamRule, markLine int) Section {
	sec := Section{Name: name, Rules: make([]Rule, len(rules))}
	for i, r := range rules {
		sec.Rules[i] = r.Rule
	}
	if len(rules) > 0 {
		sec.Start = rules[0].Line
		sec.End = rules[len(rules)-1].endLine
	}
	if markLine != 0 {
		sec.Start = markLine
		if sec.End == 0 {
			sec.End = markLine
		}
	}
	return sec
}

// wrapperSelector matches the single outer `.style-<name>` selector and
// returns <name>.
func wrapperSelector(tokens []css.Token) (string, bool) {
	class, ok := singleClassSelector(tokens)
	if !ok {
		return "", false
	}
	return strings.CutPrefix(class, "style-")
}

// cnSelector matches a `.cn-<...>` selector and returns the class (without
// the leading dot).
func cnSelector(tokens []css.Token) (string, bool) {
	class, ok := singleClassSelector(tokens)
	if !ok || !strings.HasPrefix(class, "cn-") {
		return "", false
	}
	return class, true
}

// selectorText renders a CSS selector's tokens back to source text, for
// error messages.
func selectorText(tokens []css.Token) string {
	var selector bytes.Buffer
	for _, token := range tokens {
		selector.Write(token.Data)
	}
	return selector.String()
}

// singleClassSelector requires the selector to be exactly one `.` delim
// token followed by one ident token — no compound selectors, combinators,
// pseudo-classes, attribute selectors or comma lists.
func singleClassSelector(tokens []css.Token) (string, bool) {
	if len(tokens) != 2 ||
		tokens[0].TokenType != css.DelimToken || !bytes.Equal(tokens[0].Data, []byte(".")) ||
		tokens[1].TokenType != css.IdentToken {
		return "", false
	}
	return string(tokens[1].Data), true
}

// parseApply splits an `@apply` rule's tokens into whitespace-separated
// utilities, tracking bracket/paren/function depth so a utility carrying an
// arbitrary value (e.g. `data-[slot=x]:size-4`) is never split internally.
func parseApply(tokens []css.Token) ([]string, error) {
	var utilities []string
	var current bytes.Buffer
	depth := 0

	finish := func() {
		if current.Len() == 0 {
			return
		}
		utilities = append(utilities, current.String())
		current.Reset()
	}

	for _, token := range tokens {
		if token.TokenType == css.WhitespaceToken && depth == 0 {
			finish()
			continue
		}
		switch token.TokenType {
		case css.LeftBracketToken, css.LeftParenthesisToken, css.FunctionToken:
			depth++
		case css.RightBracketToken, css.RightParenthesisToken:
			depth--
			if depth < 0 {
				return nil, errors.New("unbalanced utility delimiters")
			}
		}
		current.Write(token.Data)
	}
	if depth != 0 {
		return nil, errors.New("unterminated utility")
	}
	finish()
	if len(utilities) == 0 {
		return nil, errors.New("missing utility")
	}
	return utilities, nil
}

func upstreamError(filename string, src []byte, offset int, format string, args ...any) error {
	if offset < 0 {
		offset = 0
	}
	if offset > len(src) {
		offset = len(src)
	}
	line, column, _ := parse.Position(bytes.NewReader(src), offset)
	return fmt.Errorf("%s:%d:%d: %s", filename, line, column, fmt.Sprintf(format, args...))
}
