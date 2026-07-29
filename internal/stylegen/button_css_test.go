package stylegen

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
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
