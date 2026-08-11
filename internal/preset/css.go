package preset

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	parse "github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/css"
)

type themeMode uint8

const (
	themeModeUnsupported themeMode = iota
	themeModeLight
	themeModeDark
)

type cssScope struct {
	atRule bool
	mode   themeMode
}

func ThemeCSS(preset Preset) ([]byte, error) {
	if err := Validate(preset); err != nil {
		return nil, fmt.Errorf("theme CSS: %w", err)
	}

	var buffer bytes.Buffer
	buffer.WriteString(":root {\n")
	writeCSSProperty(&buffer, "radius", preset.Radius)
	// FontSans/FontHeading are not per-mode (see the Preset field's own doc
	// comment), so they're written once in :root — .dark inherits them like
	// any other unset custom property.
	writeCSSProperty(&buffer, "font-sans", preset.FontSans)
	writeCSSProperty(&buffer, "font-heading", preset.FontHeading)
	for _, name := range TokenNames() {
		writeCSSProperty(&buffer, name, preset.Theme.Light[name])
	}
	buffer.WriteString("}\n\n.dark {\n")
	for _, name := range TokenNames() {
		writeCSSProperty(&buffer, name, preset.Theme.Dark[name])
	}
	buffer.WriteString("}\n")
	return buffer.Bytes(), nil
}

func ImportThemeCSS(base Preset, src []byte) (Preset, error) {
	if err := Validate(base); err != nil {
		return Preset{}, fmt.Errorf("theme CSS import base: %w", err)
	}
	if err := validateCSSStructure(src); err != nil {
		return Preset{}, fmt.Errorf("theme CSS import: malformed CSS: %w", err)
	}

	candidate := clonePreset(base)
	parser := css.NewParser(parse.NewInputBytes(src), false)
	var scopes []cssScope
	seen := map[themeMode]map[string]bool{
		themeModeLight: {},
		themeModeDark:  {},
	}
	var radiusValue string
	var radiusSeen bool
	var fontSansValue string
	var fontSansSeen bool
	var fontHeadingValue string
	var fontHeadingSeen bool

	for {
		grammar, _, data := parser.Next()
		switch grammar {
		case css.ErrorGrammar:
			if err := parser.Err(); !errors.Is(err, io.EOF) {
				return Preset{}, fmt.Errorf("theme CSS import: malformed CSS: %w", err)
			}
			if len(scopes) != 0 {
				return Preset{}, errors.New("theme CSS import: malformed CSS: unclosed rule")
			}
			if err := Validate(candidate); err != nil {
				return Preset{}, fmt.Errorf("theme CSS import: %w", err)
			}
			return clonePreset(candidate), nil

		case css.BeginAtRuleGrammar:
			scopes = append(scopes, cssScope{atRule: true})
		case css.EndAtRuleGrammar:
			if !popCSSScope(&scopes, true) {
				return Preset{}, errors.New("theme CSS import: malformed CSS at-rule ownership")
			}
		case css.BeginRulesetGrammar:
			scopes = append(scopes, cssScope{mode: classifyThemeSelector(parser.Values())})
		case css.EndRulesetGrammar:
			if !popCSSScope(&scopes, false) {
				return Preset{}, errors.New("theme CSS import: malformed CSS ruleset ownership")
			}
		case css.CustomPropertyGrammar:
			name, recognized := recognizedThemeProperty(string(data))
			if !recognized {
				continue
			}
			mode, err := declarationThemeMode(scopes, data)
			if err != nil {
				return Preset{}, fmt.Errorf("theme CSS import: %w", err)
			}
			if seen[mode][name] {
				return Preset{}, fmt.Errorf("theme CSS import: duplicate recognized declaration --%s in %s", name, mode.path())
			}
			seen[mode][name] = true

			value := customPropertyValue(parser.Values())
			switch name {
			case "radius":
				if radiusSeen && radiusValue != value {
					return Preset{}, fmt.Errorf("theme CSS import: conflicting radius declarations %q and %q", radiusValue, value)
				}
				radiusSeen = true
				radiusValue = value
				candidate.Radius = value
			case "font-sans":
				if fontSansSeen && fontSansValue != value {
					return Preset{}, fmt.Errorf("theme CSS import: conflicting font-sans declarations %q and %q", fontSansValue, value)
				}
				fontSansSeen = true
				fontSansValue = value
				candidate.FontSans = value
			case "font-heading":
				if fontHeadingSeen && fontHeadingValue != value {
					return Preset{}, fmt.Errorf("theme CSS import: conflicting font-heading declarations %q and %q", fontHeadingValue, value)
				}
				fontHeadingSeen = true
				fontHeadingValue = value
				candidate.FontHeading = value
			default:
				if mode == themeModeLight {
					candidate.Theme.Light[name] = value
				} else {
					candidate.Theme.Dark[name] = value
				}
			}
		}
	}
}

func validateCSSStructure(src []byte) error {
	lexer := css.NewLexer(parse.NewInputBytes(src))
	var blocks []css.TokenType
	for {
		tokenType, _ := lexer.Next()
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

func writeCSSProperty(buffer *bytes.Buffer, name, value string) {
	buffer.WriteString("  --")
	buffer.WriteString(name)
	buffer.WriteString(": ")
	buffer.WriteString(value)
	buffer.WriteString(";\n")
}

func classifyThemeSelector(tokens []css.Token) themeMode {
	significant := make([]css.Token, 0, len(tokens))
	for _, token := range tokens {
		if token.TokenType != css.WhitespaceToken && token.TokenType != css.CommentToken {
			significant = append(significant, token)
		}
	}
	if len(significant) != 2 {
		return themeModeUnsupported
	}
	if significant[0].TokenType == css.ColonToken &&
		bytes.Equal(significant[0].Data, []byte(":")) &&
		significant[1].TokenType == css.IdentToken &&
		bytes.EqualFold(significant[1].Data, []byte("root")) {
		return themeModeLight
	}
	if significant[0].TokenType == css.DelimToken &&
		bytes.Equal(significant[0].Data, []byte(".")) &&
		significant[1].TokenType == css.IdentToken &&
		bytes.Equal(significant[1].Data, []byte("dark")) {
		return themeModeDark
	}
	return themeModeUnsupported
}

func recognizedThemeProperty(property string) (string, bool) {
	if !strings.HasPrefix(property, "--") {
		return "", false
	}
	name := strings.TrimPrefix(property, "--")
	return name, name == "radius" || name == "font-sans" || name == "font-heading" || isTokenName(name)
}

func declarationThemeMode(scopes []cssScope, property []byte) (themeMode, error) {
	if len(scopes) == 0 {
		return themeModeUnsupported, fmt.Errorf("recognized property %s outside :root or .dark", property)
	}
	if len(scopes) != 1 || scopes[0].atRule {
		return themeModeUnsupported, fmt.Errorf("recognized property %s has nested declaration ownership", property)
	}
	if scopes[0].mode == themeModeUnsupported {
		return themeModeUnsupported, fmt.Errorf("recognized property %s is outside supported selector :root or .dark", property)
	}
	return scopes[0].mode, nil
}

func customPropertyValue(tokens []css.Token) string {
	var value bytes.Buffer
	for _, token := range tokens {
		value.Write(token.Data)
	}
	return strings.TrimSpace(value.String())
}

func popCSSScope(scopes *[]cssScope, atRule bool) bool {
	if len(*scopes) == 0 {
		return false
	}
	last := len(*scopes) - 1
	if (*scopes)[last].atRule != atRule {
		return false
	}
	*scopes = (*scopes)[:last]
	return true
}

func (mode themeMode) path() string {
	if mode == themeModeLight {
		return "theme.light"
	}
	return "theme.dark"
}
