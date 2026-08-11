package preset

import (
	"errors"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/mazznoer/csscolorparser"
	parse "github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/css"
)

const SchemaVersion = 1
const SchemaURL = "https://ui.gsxhq.dev/schemas/preset-v1.json"

const (
	defaultRadius = "0.625rem"
	radiusGroup   = "Structure"
)

type Style string

const (
	StyleVega Style = "vega"
	StyleNova Style = "nova"
	StyleMaia Style = "maia"
	StyleLyra Style = "lyra"
	StyleMira Style = "mira"
	StyleLuma Style = "luma"
	StyleSera Style = "sera"
	StyleRhea Style = "rhea"
)

// StyleChoice is one style's display metadata for a picker UI: Name is the
// Style value itself (as a plain string, matching PaletteChoice/
// RadiusChoice's own shape below), Title is its capitalized display name,
// and Description is upstream's own one-line description
// (apps/v4/registry/styles.tsx in the shadcn-ui checkout).
type StyleChoice struct {
	Name        string
	Title       string
	Description string
}

var styleTitles = map[Style]string{
	StyleVega: "Vega",
	StyleNova: "Nova",
	StyleMaia: "Maia",
	StyleLyra: "Lyra",
	StyleMira: "Mira",
	StyleLuma: "Luma",
	StyleSera: "Sera",
	StyleRhea: "Rhea",
}

var styleDescriptions = map[Style]string{
	StyleVega: "Clean, neutral, and familiar",
	StyleNova: "Reduced padding and margins",
	StyleMaia: "Rounded, with generous spacing.",
	StyleLyra: "Boxy and sharp. For mono fonts.",
	StyleMira: "Made for compact interfaces.",
	StyleLuma: "Fluid, luminous, and soft.",
	StyleSera: "Editorial and typographic.",
	StyleRhea: "Like Luma but compact.",
}

// StyleChoices returns every style's display metadata, in Styles()'s own
// order, for a picker UI to render directly — the same catalogue shape
// BaseColorChoices/RadiusChoices already give their own axes.
func StyleChoices() []StyleChoice {
	choices := make([]StyleChoice, len(styles))
	for i, s := range styles {
		choices[i] = StyleChoice{Name: string(s), Title: styleTitles[s], Description: styleDescriptions[s]}
	}
	return choices
}

type ThemeValues map[string]string

type Theme struct {
	Light ThemeValues
	Dark  ThemeValues
}

type Preset struct {
	Schema        string
	SchemaVersion int
	Style         Style
	Radius        string
	Theme         Theme
}

type TokenDefinition struct {
	Name  string
	Group string
	Light string
	Dark  string
}

// styles lists all 8 styles in the plan's own canonical order (Global
// Constraints: "vega, nova, maia, lyra, mira, luma, sera, rhea"). Unlike
// compact.go's compactStyles, this order has no wire-format ABI to
// preserve — it only governs Styles()'s and StyleChoices()'s own iteration
// order, i.e. picker display order.
var styles = []Style{StyleVega, StyleNova, StyleMaia, StyleLyra, StyleMira, StyleLuma, StyleSera, StyleRhea}

var presentationGroups = []string{
	"Base",
	"Brand",
	"Feedback",
	"Structure",
	"Status and overlay",
	"Sidebar",
}

var tokenDefinitions = []TokenDefinition{
	{Name: "background", Group: "Base", Light: "oklch(1 0 0)", Dark: "oklch(0.145 0 0)"},
	{Name: "foreground", Group: "Base", Light: "oklch(0.145 0 0)", Dark: "oklch(0.985 0 0)"},
	{Name: "card", Group: "Base", Light: "oklch(1 0 0)", Dark: "oklch(0.205 0 0)"},
	{Name: "card-foreground", Group: "Base", Light: "oklch(0.145 0 0)", Dark: "oklch(0.985 0 0)"},
	{Name: "popover", Group: "Base", Light: "oklch(1 0 0)", Dark: "oklch(0.205 0 0)"},
	{Name: "popover-foreground", Group: "Base", Light: "oklch(0.145 0 0)", Dark: "oklch(0.985 0 0)"},
	{Name: "primary", Group: "Brand", Light: "oklch(0.205 0 0)", Dark: "oklch(0.922 0 0)"},
	{Name: "primary-foreground", Group: "Brand", Light: "oklch(0.985 0 0)", Dark: "oklch(0.205 0 0)"},
	{Name: "secondary", Group: "Brand", Light: "oklch(0.97 0 0)", Dark: "oklch(0.269 0 0)"},
	{Name: "secondary-foreground", Group: "Brand", Light: "oklch(0.205 0 0)", Dark: "oklch(0.985 0 0)"},
	{Name: "muted", Group: "Feedback", Light: "oklch(0.97 0 0)", Dark: "oklch(0.269 0 0)"},
	{Name: "muted-foreground", Group: "Feedback", Light: "oklch(0.556 0 0)", Dark: "oklch(0.708 0 0)"},
	{Name: "accent", Group: "Brand", Light: "oklch(0.97 0 0)", Dark: "oklch(0.269 0 0)"},
	{Name: "accent-foreground", Group: "Brand", Light: "oklch(0.205 0 0)", Dark: "oklch(0.985 0 0)"},
	{Name: "destructive", Group: "Feedback", Light: "oklch(0.577 0.245 27.325)", Dark: "oklch(0.704 0.191 22.216)"},
	{Name: "destructive-foreground", Group: "Feedback", Light: "oklch(0.97 0.01 17)", Dark: "oklch(0.58 0.22 27)"},
	{Name: "border", Group: "Structure", Light: "oklch(0.922 0 0)", Dark: "oklch(1 0 0 / 10%)"},
	{Name: "input", Group: "Structure", Light: "oklch(0.922 0 0)", Dark: "oklch(1 0 0 / 15%)"},
	{Name: "ring", Group: "Structure", Light: "oklch(0.708 0 0)", Dark: "oklch(0.556 0 0)"},
	{Name: "sidebar", Group: "Sidebar", Light: "oklch(0.985 0 0)", Dark: "oklch(0.205 0 0)"},
	{Name: "sidebar-foreground", Group: "Sidebar", Light: "oklch(0.145 0 0)", Dark: "oklch(0.985 0 0)"},
	{Name: "sidebar-primary", Group: "Sidebar", Light: "oklch(0.205 0 0)", Dark: "oklch(0.488 0.243 264.376)"},
	{Name: "sidebar-primary-foreground", Group: "Sidebar", Light: "oklch(0.985 0 0)", Dark: "oklch(0.985 0 0)"},
	{Name: "sidebar-accent", Group: "Sidebar", Light: "oklch(0.97 0 0)", Dark: "oklch(0.269 0 0)"},
	{Name: "sidebar-accent-foreground", Group: "Sidebar", Light: "oklch(0.205 0 0)", Dark: "oklch(0.985 0 0)"},
	{Name: "sidebar-border", Group: "Sidebar", Light: "oklch(0.922 0 0)", Dark: "oklch(1 0 0 / 10%)"},
	{Name: "sidebar-ring", Group: "Sidebar", Light: "oklch(0.708 0 0)", Dark: "oklch(0.556 0 0)"},
	{Name: "success", Group: "Status and overlay", Light: "oklch(69.6% 0.17 162.48)", Dark: "oklch(69.6% 0.17 162.48)"},
	{Name: "info", Group: "Status and overlay", Light: "oklch(68.5% 0.169 237.323)", Dark: "oklch(68.5% 0.169 237.323)"},
	{Name: "warning", Group: "Status and overlay", Light: "oklch(76.9% 0.188 70.08)", Dark: "oklch(76.9% 0.188 70.08)"},
	{Name: "overlay", Group: "Status and overlay", Light: "oklch(0% 0 0 / 10%)", Dark: "oklch(0% 0 0 / 10%)"},
	{Name: "contrast", Group: "Status and overlay", Light: "oklch(100% 0 0)", Dark: "oklch(100% 0 0)"},
}

var lengthUnits = []string{
	"px", "cm", "mm", "q", "in", "pc", "pt",
	"em", "rem", "ex", "rex", "cap", "rcap", "ch", "rch", "ic", "ric", "lh", "rlh",
	"vw", "svw", "lvw", "dvw", "vh", "svh", "lvh", "dvh",
	"vi", "svi", "lvi", "dvi", "vb", "svb", "lvb", "dvb",
	"vmin", "svmin", "lvmin", "dvmin", "vmax", "svmax", "lvmax", "dvmax",
}

func Styles() []Style {
	return slices.Clone(styles)
}

func TokenNames() []string {
	names := make([]string, len(tokenDefinitions))
	for i, definition := range tokenDefinitions {
		names[i] = definition.Name
	}
	return names
}

func GroupNames() []string {
	return slices.Clone(presentationGroups)
}

func TokenDefinitions() []TokenDefinition {
	return slices.Clone(tokenDefinitions)
}

func RadiusDefinition() TokenDefinition {
	return TokenDefinition{
		Name:  "radius",
		Group: radiusGroup,
		Light: defaultRadius,
		Dark:  defaultRadius,
	}
}

// RadiusUnits returns the exact CSS length-unit allowlist used by Validate.
// The web editor serializes it into its server-authored schema so browser
// drafts cannot accept a radius the CLI would later reject.
func RadiusUnits() []string {
	return slices.Clone(lengthUnits)
}

func Default(style Style) Preset {
	preset, err := resolvePalette(style, DefaultPaletteSelection())
	if err == nil {
		return preset
	}
	return canonicalDefault(style)
}

func Validate(preset Preset) error {
	if preset.Schema != SchemaURL {
		return fmt.Errorf("$schema: must be %q", SchemaURL)
	}
	if preset.SchemaVersion != SchemaVersion {
		return fmt.Errorf("schemaVersion: must be %d", SchemaVersion)
	}
	if !slices.Contains(styles, preset.Style) {
		return fmt.Errorf("style: unsupported value %q", preset.Style)
	}
	if err := validateRadius(preset.Radius); err != nil {
		return fmt.Errorf("radius: %w", err)
	}
	if err := validateThemeMode("theme.light", preset.Theme.Light); err != nil {
		return err
	}
	if err := validateThemeMode("theme.dark", preset.Theme.Dark); err != nil {
		return err
	}
	return nil
}

func validateThemeMode(path string, values ThemeValues) error {
	for _, definition := range tokenDefinitions {
		value, ok := values[definition.Name]
		if !ok {
			return fmt.Errorf("%s.%s: missing required color", path, definition.Name)
		}
		if _, err := csscolorparser.Parse(value); err != nil {
			return fmt.Errorf("%s.%s: invalid CSS color: %w", path, definition.Name, err)
		}
	}

	var extras []string
	for name := range values {
		if !isTokenName(name) {
			extras = append(extras, name)
		}
	}
	slices.Sort(extras)
	if len(extras) != 0 {
		return fmt.Errorf("%s.%s: unknown token", path, extras[0])
	}
	return nil
}

func isTokenName(name string) bool {
	for _, definition := range tokenDefinitions {
		if definition.Name == name {
			return true
		}
	}
	return false
}

func validateRadius(value string) error {
	lexer := css.NewLexer(parse.NewInputString(value))
	tokenType, token := nextComponentToken(lexer)
	switch tokenType {
	case css.NumberToken:
		number, err := strconv.ParseFloat(string(token), 64)
		if err != nil || number != 0 {
			return errors.New("must be zero or a non-negative CSS length")
		}
	case css.DimensionToken:
		numberLength, unitLength := parse.Dimension(token)
		if numberLength == 0 || numberLength+unitLength != len(token) {
			return errors.New("must be zero or a non-negative CSS length")
		}
		number, err := strconv.ParseFloat(string(token[:numberLength]), 64)
		if err != nil || number < 0 {
			return errors.New("must be zero or a non-negative CSS length")
		}
		unit := strings.ToLower(string(token[numberLength:]))
		if !slices.Contains(lengthUnits, unit) {
			return fmt.Errorf("unsupported CSS length unit %q", unit)
		}
	default:
		return errors.New("must be zero or a non-negative CSS length")
	}

	if trailingType, _ := nextComponentToken(lexer); trailingType != css.ErrorToken ||
		!errors.Is(lexer.Err(), io.EOF) {
		return errors.New("must contain exactly one CSS component value")
	}
	return nil
}

func nextComponentToken(lexer *css.Lexer) (css.TokenType, []byte) {
	for {
		tokenType, token := lexer.Next()
		if tokenType != css.WhitespaceToken && tokenType != css.CommentToken {
			return tokenType, token
		}
	}
}
