// Package stylegen generates consumer component source from authored style recipes.
package stylegen

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"slices"

	parse "github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/css"
)

const RecipePrefix = "gsxui-recipe-"

type Recipe struct {
	Token     string
	Utilities []string
}

type Recipes struct {
	ordered []Recipe
	byToken map[string]Recipe
}

func ParseRecipes(filename string, src []byte) (Recipes, error) {
	parser := css.NewParser(parse.NewInputBytes(src), false)
	var result Recipes
	var componentsLayer bool
	var layerSeen bool
	var rule *recipeRule
	lastContext := "stylesheet"

	for {
		grammar, _, data := parser.Next()
		switch grammar {
		case css.ErrorGrammar:
			if err := parser.Err(); errors.Is(err, io.EOF) {
				if rule != nil {
					return Recipes{}, recipeError(filename, src, parser.Offset(), "malformed CSS in %s", rule.selector)
				}
				if componentsLayer {
					return Recipes{}, recipeError(filename, src, parser.Offset(), "malformed CSS in @layer components")
				}
				if !layerSeen {
					return Recipes{}, recipeError(filename, src, len(src), "missing @layer components")
				}
				if len(result.ordered) == 0 {
					return Recipes{}, recipeError(filename, src, len(src), "@layer components contains no recipe rules")
				}
				if err := validateCSSStructure(src); err != nil {
					return Recipes{}, recipeError(filename, src, len(src), "malformed CSS in %s: %v", lastContext, err)
				}
				return result, nil
			} else if err != nil {
				return Recipes{}, parserRecipeError(filename, src, parser.Offset(), err, malformedContext(rule, componentsLayer))
			}
			return Recipes{}, recipeError(filename, src, parser.Offset(), "malformed CSS")

		case css.CommentGrammar:
			continue

		case css.BeginAtRuleGrammar:
			atRule := atRuleText(data, parser.Values())
			if rule != nil {
				return Recipes{}, recipeError(filename, src, parser.Offset(), "nested at-rule %s in %s", atRule, rule.selector)
			}
			if componentsLayer {
				return Recipes{}, recipeError(filename, src, parser.Offset(), "nested at-rule %s in @layer components", atRule)
			}
			if layerSeen {
				return Recipes{}, recipeError(filename, src, parser.Offset(), "multiple component layers: %s", atRule)
			}
			if !isComponentsLayer(data, parser.Values()) {
				return Recipes{}, recipeError(filename, src, parser.Offset(), "expected @layer components, found %s", atRule)
			}
			layerSeen = true
			componentsLayer = true
			lastContext = "@layer components"

		case css.EndAtRuleGrammar:
			if !componentsLayer || rule != nil {
				return Recipes{}, recipeError(filename, src, parser.Offset(), "malformed @layer components ownership")
			}
			componentsLayer = false

		case css.AtRuleGrammar:
			atRule := atRuleText(data, parser.Values())
			if rule == nil {
				return Recipes{}, recipeError(filename, src, parser.Offset(), "expected @layer components block, found %s", atRule)
			}
			if !bytes.Equal(data, []byte("@apply")) {
				return Recipes{}, recipeError(filename, src, parser.Offset(), "expected @apply in %s, found %s", rule.selector, atRule)
			}
			utilities, err := parseApply(parser.Values())
			if err != nil {
				return Recipes{}, recipeError(filename, src, parser.Offset(), "invalid @apply in %s: %v (%s)", rule.selector, err, atRule)
			}
			rule.utilities = append(rule.utilities, utilities...)

		case css.BeginRulesetGrammar:
			selector := selectorText(parser.Values())
			if !componentsLayer {
				return Recipes{}, recipeError(filename, src, parser.Offset(), "recipe rule %s must be inside @layer components", selector)
			}
			if rule != nil {
				return Recipes{}, recipeError(filename, src, parser.Offset(), "nested rule %s in %s", selector, rule.selector)
			}
			token, ok := recipeSelector(parser.Values())
			if !ok {
				return Recipes{}, recipeError(filename, src, parser.Offset(), "invalid recipe selector %s", selector)
			}
			rule = &recipeRule{selector: selector, token: token}
			lastContext = selector

		case css.EndRulesetGrammar:
			if rule == nil {
				return Recipes{}, recipeError(filename, src, parser.Offset(), "ruleset outside @layer components")
			}
			if len(rule.utilities) == 0 {
				return Recipes{}, recipeError(filename, src, parser.Offset(), "recipe rule %s requires at least one @apply", rule.selector)
			}
			if _, exists := result.byToken[rule.token]; exists {
				return Recipes{}, recipeError(filename, src, parser.Offset(), "duplicate recipe token %s in %s", rule.token, rule.selector)
			}
			if result.byToken == nil {
				result.byToken = make(map[string]Recipe)
			}
			recipe := Recipe{Token: rule.token, Utilities: slices.Clone(rule.utilities)}
			result.ordered = append(result.ordered, recipe)
			result.byToken[recipe.Token] = recipe
			rule = nil

		case css.DeclarationGrammar, css.CustomPropertyGrammar:
			if rule != nil {
				return Recipes{}, recipeError(filename, src, parser.Offset(), "ordinary declaration in recipe rule %s", rule.selector)
			}
			return Recipes{}, recipeError(filename, src, parser.Offset(), "ordinary declaration outside recipe rule")

		default:
			return Recipes{}, recipeError(filename, src, parser.Offset(), "unsupported CSS grammar %s", grammar)
		}
	}
}

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

func (r Recipes) Lookup(token string) (Recipe, bool) {
	recipe, ok := r.byToken[token]
	if !ok {
		return Recipe{}, false
	}
	recipe.Utilities = slices.Clone(recipe.Utilities)
	return recipe, true
}

func (r Recipes) Tokens() []string {
	tokens := make([]string, len(r.ordered))
	for i, recipe := range r.ordered {
		tokens[i] = recipe.Token
	}
	return tokens
}

type recipeRule struct {
	selector  string
	token     string
	utilities []string
}

func isComponentsLayer(name []byte, values []css.Token) bool {
	values = significantTokens(values)
	return bytes.Equal(name, []byte("@layer")) &&
		len(values) == 1 &&
		values[0].TokenType == css.IdentToken &&
		bytes.Equal(values[0].Data, []byte("components"))
}

func recipeSelector(tokens []css.Token) (string, bool) {
	if len(tokens) != 2 ||
		tokens[0].TokenType != css.DelimToken || !bytes.Equal(tokens[0].Data, []byte(".")) ||
		tokens[1].TokenType != css.IdentToken {
		return "", false
	}
	token := string(tokens[1].Data)
	if len(token) == len(RecipePrefix) || !bytes.HasPrefix(tokens[1].Data, []byte(RecipePrefix)) || bytes.ContainsRune(tokens[1].Data, '\\') {
		return "", false
	}
	return token, true
}

func parseApply(tokens []css.Token) ([]string, error) {
	var utilities []string
	var current bytes.Buffer
	var last css.TokenType
	depth := 0

	finish := func() error {
		if current.Len() == 0 {
			return nil
		}
		if last == css.ColonToken || (last == css.DelimToken && bytes.Equal(current.Bytes()[current.Len()-1:], []byte("/"))) {
			return errors.New("unterminated utility")
		}
		utilities = append(utilities, current.String())
		current.Reset()
		last = css.EmptyToken
		return nil
	}

	for _, token := range tokens {
		if token.TokenType == css.WhitespaceToken && depth == 0 {
			if err := finish(); err != nil {
				return nil, err
			}
			continue
		}
		if token.TokenType == css.CommaToken && depth == 0 {
			return nil, errors.New("utilities must be whitespace-separated")
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
		last = token.TokenType
	}
	if depth != 0 {
		return nil, errors.New("unterminated utility")
	}
	if err := finish(); err != nil {
		return nil, err
	}
	if len(utilities) == 0 {
		return nil, errors.New("missing utility")
	}
	return utilities, nil
}

func selectorText(tokens []css.Token) string {
	var selector bytes.Buffer
	for _, token := range tokens {
		selector.Write(token.Data)
	}
	return selector.String()
}

func atRuleText(name []byte, tokens []css.Token) string {
	var atRule bytes.Buffer
	atRule.Write(name)
	tokens = significantTokens(tokens)
	if len(tokens) > 0 {
		atRule.WriteByte(' ')
	}
	for _, token := range tokens {
		atRule.Write(token.Data)
	}
	return atRule.String()
}

func significantTokens(tokens []css.Token) []css.Token {
	first := 0
	for first < len(tokens) && tokens[first].TokenType == css.WhitespaceToken {
		first++
	}
	last := len(tokens)
	for first < last && tokens[last-1].TokenType == css.WhitespaceToken {
		last--
	}
	return tokens[first:last]
}

func malformedContext(rule *recipeRule, componentsLayer bool) string {
	if rule != nil {
		return rule.selector
	}
	if componentsLayer {
		return "@layer components"
	}
	return "stylesheet"
}

func parserRecipeError(filename string, src []byte, offset int, err error, context string) error {
	if parseErr, ok := errors.AsType[*parse.Error](err); ok {
		return fmt.Errorf("%s:%d:%d: malformed CSS in %s: %w", filename, parseErr.Line, parseErr.Column, context, err)
	}
	return recipeError(filename, src, offset, "malformed CSS in %s: %v", context, err)
}

func recipeError(filename string, src []byte, offset int, format string, args ...any) error {
	if offset < 0 {
		offset = 0
	}
	if offset > len(src) {
		offset = len(src)
	}
	line, column, _ := parse.Position(bytes.NewReader(src), offset)
	return fmt.Errorf("%s:%d:%d: %s", filename, line, column, fmt.Sprintf(format, args...))
}
