package stylegen

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"

	parse "github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/css"

	"github.com/gsxhq/gsxui/internal/recipe"
)

const buttonSlotAttribute = "data-gsxui-slot-button"

type cssOwnershipRule struct {
	selector     string
	attributes   map[string]bool
	declarations map[string]bool
	directChild  bool
}

func TestDefaultStyleHasNoButtonPresentationSelectors(t *testing.T) {
	src := readDefaultStyle(t)
	rules, err := parseCSSOwnershipRules(src)
	if err != nil {
		t.Fatal(err)
	}

	var selectors []string
	for _, rule := range rules {
		if rule.attributes[buttonSlotAttribute] {
			selectors = append(selectors, rule.selector)
		}
	}
	if len(selectors) != 0 {
		t.Fatalf(
			"assets/css/styles/default.css still owns Button presentation through %d parsed rules; first selector: %s",
			len(selectors),
			selectors[0],
		)
	}
}

func TestDefaultStyleRetainsComposedComponentButtonPresentation(t *testing.T) {
	src := readDefaultStyle(t)
	rules, err := parseCSSOwnershipRules(src)
	if err != nil {
		t.Fatal(err)
	}

	var (
		buttonGroupSizeRule bool
		paginationPrevious  bool
		paginationNext      bool
		inputGroupRules     int
	)
	for _, rule := range rules {
		attrs := rule.attributes
		buttonGroupSizeRule = buttonGroupSizeRule ||
			(attrs["data-gsxui-slot-button-group"] && attrs["data-size"] && rule.directChild)
		paginationPrevious = paginationPrevious ||
			(attrs["data-gsxui-slot-pagination-link"] && attrs["data-gsxui-slot-pagination-previous"])
		paginationNext = paginationNext ||
			(attrs["data-gsxui-slot-pagination-link"] && attrs["data-gsxui-slot-pagination-next"])
		if attrs["data-gsxui-slot-input-group-button"] {
			inputGroupRules++
		}
	}

	if !buttonGroupSizeRule {
		t.Error("consumer style lost ButtonGroup direct-child size presentation")
	}
	if !paginationPrevious || !paginationNext {
		t.Errorf(
			"consumer style lost Pagination edge presentation: previous=%t next=%t",
			paginationPrevious,
			paginationNext,
		)
	}
	if inputGroupRules != 9 {
		t.Errorf("consumer style has %d InputGroupButton presentation rules, want 9", inputGroupRules)
	}
}

// TestCalendarRecipeRetainsButtonGeometry is the other half of the check
// above, for a component that has since migrated. Calendar's nav and day
// buttons are raw <button>s stamped with data-gsxui-slot-button, and their
// --cell-size geometry was checked in the consumer's flat style.css while
// Calendar was unmigrated. It now lives on Calendar's own recipe instead, so
// the guarantee moves with it rather than being dropped: a future edit that
// swept this geometry into Button's own recipe would still be caught.
func TestCalendarRecipeRetainsButtonGeometry(t *testing.T) {
	for _, style := range []string{"nova", "maia"} {
		path := filepath.Join("..", "..", "registry", "styles", style, "calendar.css")
		src, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		parsed, err := recipe.ParseStyle(path, src)
		if err != nil {
			t.Fatal(err)
		}
		for _, want := range []struct {
			class     string
			utilities []string
		}{
			{"gsxui-recipe-calendar-nav-button", []string{"w-(--cell-size)", "h-(--cell-size)"}},
			{"gsxui-recipe-calendar-day-button", []string{"min-w-(--cell-size)"}},
		} {
			rule, ok := parsed.Lookup(want.class)
			if !ok {
				t.Errorf("%s: %s declares no rule", style, want.class)
				continue
			}
			for _, utility := range want.utilities {
				if !slices.Contains(rule.Utilities, utility) {
					t.Errorf("%s: %s lost %s (utilities: %q)", style, want.class, utility, rule.Utilities)
				}
			}
		}
	}
}

// readDefaultStyle returns the default style's components layer as ONE sheet:
// the aggregator with each of its relative @import lines replaced by the file it
// names, which is byte-for-byte what internal/cli vendors into a consumer's
// style.css. The tests here ask what a consumer's stylesheet still owns, so they
// have to read the whole composed sheet — reading only the aggregator would see
// the @layer utilities block and nothing else, and would report every component
// as lost.
func readDefaultStyle(t *testing.T) []byte {
	t.Helper()
	dir := filepath.Join("..", "..", "assets", "css", "styles")
	aggregator, err := os.ReadFile(filepath.Join(dir, "default.css"))
	if err != nil {
		t.Fatal(err)
	}
	var out []byte
	for _, line := range strings.SplitAfter(string(aggregator), "\n") {
		match := defaultStyleImport.FindStringSubmatch(strings.TrimSuffix(line, "\n"))
		if match == nil {
			out = append(out, line...)
			continue
		}
		part, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(match[1])))
		if err != nil {
			t.Fatal(err)
		}
		out = append(out, part...)
	}
	return out
}

// defaultStyleImport matches a whole-line relative @import in the aggregator.
var defaultStyleImport = regexp.MustCompile(`^@import "(\./[^"]+)";[ \t]*$`)

func parseCSSOwnershipRules(src []byte) ([]cssOwnershipRule, error) {
	parser := css.NewParser(parse.NewInputBytes(src), false)
	var rules []cssOwnershipRule
	var rule *cssOwnershipRule

	for {
		grammar, _, data := parser.Next()
		switch grammar {
		case css.ErrorGrammar:
			if err := parser.Err(); !errors.Is(err, io.EOF) {
				return nil, err
			}
			if err := recipe.ValidateCSSStructure(src); err != nil {
				return nil, err
			}
			return rules, nil
		case css.BeginRulesetGrammar:
			tokens := parser.Values()
			rule = &cssOwnershipRule{
				selector:     recipe.SelectorText(tokens),
				attributes:   selectorAttributeNames(tokens),
				declarations: make(map[string]bool),
				directChild:  selectorHasDirectChild(tokens),
			}
		case css.DeclarationGrammar, css.CustomPropertyGrammar:
			if rule != nil {
				rule.declarations[strings.ToLower(string(data))] = true
			}
		case css.EndRulesetGrammar:
			if rule != nil {
				rules = append(rules, *rule)
				rule = nil
			}
		}
	}
}

func selectorAttributeNames(tokens []css.Token) map[string]bool {
	attributes := make(map[string]bool)
	for index, token := range tokens {
		if token.TokenType != css.LeftBracketToken {
			continue
		}
		for index++; index < len(tokens); index++ {
			token = tokens[index]
			switch token.TokenType {
			case css.WhitespaceToken, css.CommentToken:
				continue
			case css.IdentToken:
				attributes[strings.ToLower(string(token.Data))] = true
			}
			break
		}
	}
	return attributes
}

func selectorHasDirectChild(tokens []css.Token) bool {
	for _, token := range tokens {
		if token.TokenType == css.DelimToken && bytes.Equal(token.Data, []byte(">")) {
			return true
		}
	}
	return false
}
