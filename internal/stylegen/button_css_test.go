package stylegen

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	parse "github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/css"
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
		calendarNavGeometry bool
		calendarDayGeometry bool
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
		calendarNavGeometry = calendarNavGeometry ||
			(attrs["data-gsxui-slot-calendar-nav-button"] &&
				rule.declarations["width"] &&
				rule.declarations["height"])
		calendarDayGeometry = calendarDayGeometry ||
			(attrs["data-gsxui-slot-calendar-day-button"] &&
				rule.declarations["min-width"])
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
	if !calendarNavGeometry || !calendarDayGeometry {
		t.Errorf(
			"consumer style lost Calendar button geometry: nav=%t day=%t",
			calendarNavGeometry,
			calendarDayGeometry,
		)
	}
}

func readDefaultStyle(t *testing.T) []byte {
	t.Helper()
	path := filepath.Join("..", "..", "assets", "css", "styles", "default.css")
	src, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return src
}

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
			if err := validateCSSStructure(src); err != nil {
				return nil, err
			}
			return rules, nil
		case css.BeginRulesetGrammar:
			tokens := parser.Values()
			rule = &cssOwnershipRule{
				selector:     selectorText(tokens),
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

func selectorText(tokens []css.Token) string {
	var selector bytes.Buffer
	for _, token := range tokens {
		selector.Write(token.Data)
	}
	return selector.String()
}

// validateCSSStructure is a local copy of the same-named helper in
// internal/recipe/parse.go: it checks that CSS constructs (parentheses,
// brackets, braces, strings, comments) are balanced and well-formed. It is
// duplicated here rather than exported from internal/recipe because this
// test's ownership-rule scan is unrelated to the typed recipe model.
func validateCSSStructure(src []byte) error {
	lexer := css.NewLexer(parse.NewInputBytes(src))
	var blocks []css.TokenType
	for {
		tokenType, data := lexer.Next()
		switch tokenType {
		case css.ErrorToken:
			if err := lexer.Err(); !errors.Is(err, io.EOF) {
				return err
			}
			if len(blocks) != 0 {
				return fmt.Errorf("unclosed %s", blocks[len(blocks)-1])
			}
			return nil
		case css.BadStringToken, css.BadURLToken:
			return fmt.Errorf("invalid %s", tokenType)
		case css.CommentToken:
			if !bytes.HasSuffix(data, []byte("*/")) {
				return errors.New("unterminated comment")
			}
		case css.FunctionToken, css.LeftParenthesisToken:
			blocks = append(blocks, css.LeftParenthesisToken)
		case css.LeftBracketToken:
			blocks = append(blocks, css.LeftBracketToken)
		case css.LeftBraceToken:
			blocks = append(blocks, css.LeftBraceToken)
		case css.RightParenthesisToken:
			if !popCSSBlock(&blocks, css.LeftParenthesisToken) {
				return errors.New("unexpected right parenthesis")
			}
		case css.RightBracketToken:
			if !popCSSBlock(&blocks, css.LeftBracketToken) {
				return errors.New("unexpected right bracket")
			}
		case css.RightBraceToken:
			if !popCSSBlock(&blocks, css.LeftBraceToken) {
				return errors.New("unexpected right brace")
			}
		}
	}
}

func popCSSBlock(blocks *[]css.TokenType, want css.TokenType) bool {
	if len(*blocks) == 0 || (*blocks)[len(*blocks)-1] != want {
		return false
	}
	*blocks = (*blocks)[:len(*blocks)-1]
	return true
}
