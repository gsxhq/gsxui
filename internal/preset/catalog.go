package preset

import (
	"crypto/sha256"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/mazznoer/csscolorparser"
)

const CustomChoice = "custom"

type PaletteChoice struct {
	Name   string
	Title  string
	Swatch string
}

type RadiusChoice struct {
	Name  string
	Title string
	Value string
}

type PaletteSelection struct {
	BaseColor string
	Theme     string
	Radius    string
}

type baseColorDefinition struct {
	name   string
	title  string
	swatch string
	light  ThemeValues
	dark   ThemeValues
}

type themeDefinition struct {
	name   string
	title  string
	swatch string
	light  ThemeValues
	dark   ThemeValues
}

var baseLayerTokens = map[string]bool{
	"background": true, "foreground": true,
	"card": true, "card-foreground": true,
	"popover": true, "popover-foreground": true,
	"secondary": true, "secondary-foreground": true,
	"muted": true, "muted-foreground": true,
	"accent": true, "accent-foreground": true,
	"destructive": true, "border": true, "input": true, "ring": true,
	"sidebar": true, "sidebar-foreground": true,
	"sidebar-accent": true, "sidebar-accent-foreground": true,
	"sidebar-border": true, "sidebar-ring": true,
}

var themeLayerTokens = map[string]bool{
	"primary": true, "primary-foreground": true,
	"secondary": true, "secondary-foreground": true,
	"sidebar-primary": true, "sidebar-primary-foreground": true,
}

func BaseColorChoices() []PaletteChoice {
	choices := make([]PaletteChoice, len(baseColorCatalog))
	for i, definition := range baseColorCatalog {
		choices[i] = PaletteChoice{
			Name: definition.name, Title: definition.title, Swatch: definition.swatch,
		}
	}
	return choices
}

func ThemeChoices(baseColor string) ([]PaletteChoice, error) {
	if baseColorDefinitionByName(baseColor) == nil {
		return nil, fmt.Errorf("base color %q is not in the palette catalog", baseColor)
	}

	choices := make([]PaletteChoice, 0, len(themeCatalog))
	for _, definition := range themeCatalog {
		if definition.name == baseColor {
			choices = append(choices, paletteChoice(definition))
			break
		}
	}
	for _, definition := range themeCatalog {
		if isAccentTheme(definition.name) {
			choices = append(choices, paletteChoice(definition))
		}
	}
	return choices, nil
}

func RadiusChoices() []RadiusChoice {
	return slices.Clone(radiusCatalog)
}

func DefaultPaletteSelection() PaletteSelection {
	return PaletteSelection{BaseColor: "neutral", Theme: "neutral", Radius: "medium"}
}

func ResolvePalette(style Style, selection PaletteSelection) (Preset, error) {
	if !slices.Contains(styles, style) {
		return Preset{}, fmt.Errorf("style: unsupported value %q", style)
	}
	preset, err := resolvePalette(style, selection)
	if err != nil {
		return Preset{}, err
	}
	if err := Validate(preset); err != nil {
		return Preset{}, fmt.Errorf("resolve palette: %w", err)
	}
	return preset, nil
}

func MatchPalette(preset Preset) PaletteSelection {
	selection := PaletteSelection{
		BaseColor: CustomChoice,
		Theme:     CustomChoice,
		Radius:    CustomChoice,
	}
	if Validate(preset) != nil {
		return selection
	}

	for _, base := range baseColorCatalog {
		choices, err := ThemeChoices(base.name)
		if err != nil {
			return selection
		}
		for _, theme := range choices {
			candidate, err := resolvePalette(preset.Style, PaletteSelection{
				BaseColor: base.name,
				Theme:     theme.Name,
				Radius:    "medium",
			})
			if err == nil && themeValuesEqual(candidate.Theme.Light, preset.Theme.Light) &&
				themeValuesEqual(candidate.Theme.Dark, preset.Theme.Dark) {
				selection.BaseColor = base.name
				selection.Theme = theme.Name
				break
			}
		}
		if selection.BaseColor != CustomChoice {
			break
		}
	}
	for _, choice := range radiusCatalog {
		if choice.Value == preset.Radius {
			selection.Radius = choice.Name
			break
		}
	}
	return selection
}

func canonicalDefault(style Style) Preset {
	light := make(ThemeValues, len(tokenDefinitions))
	dark := make(ThemeValues, len(tokenDefinitions))
	for _, definition := range tokenDefinitions {
		light[definition.Name] = definition.Light
		dark[definition.Name] = definition.Dark
	}
	return Preset{
		Schema:        SchemaURL,
		SchemaVersion: SchemaVersion,
		Style:         style,
		Radius:        defaultRadius,
		Theme:         Theme{Light: light, Dark: dark},
	}
}

func resolvePalette(style Style, selection PaletteSelection) (Preset, error) {
	base := baseColorDefinitionByName(selection.BaseColor)
	if base == nil {
		return Preset{}, fmt.Errorf("base color %q is not in the palette catalog", selection.BaseColor)
	}
	theme := themeDefinitionByName(selection.Theme)
	if theme == nil {
		return Preset{}, fmt.Errorf("theme %q is not in the palette catalog", selection.Theme)
	}
	if theme.name != selection.BaseColor && !isAccentTheme(theme.name) {
		return Preset{}, fmt.Errorf("theme %q is not available for base color %q", selection.Theme, selection.BaseColor)
	}
	radius := radiusChoiceByName(selection.Radius)
	if radius == nil {
		return Preset{}, fmt.Errorf("radius %q is not in the palette catalog", selection.Radius)
	}

	preset := canonicalDefault(style)
	applyThemeValues(preset.Theme.Light, base.light)
	applyThemeValues(preset.Theme.Dark, base.dark)
	applyThemeValues(preset.Theme.Light, theme.light)
	applyThemeValues(preset.Theme.Dark, theme.dark)
	preset.Radius = radius.Value
	return preset, nil
}

func paletteChoice(definition themeDefinition) PaletteChoice {
	return PaletteChoice{Name: definition.name, Title: definition.title, Swatch: definition.swatch}
}

func baseColorDefinitionByName(name string) *baseColorDefinition {
	for i := range baseColorCatalog {
		if baseColorCatalog[i].name == name {
			return &baseColorCatalog[i]
		}
	}
	return nil
}

func themeDefinitionByName(name string) *themeDefinition {
	for i := range themeCatalog {
		if themeCatalog[i].name == name {
			return &themeCatalog[i]
		}
	}
	return nil
}

func radiusChoiceByName(name string) *RadiusChoice {
	for i := range radiusCatalog {
		if radiusCatalog[i].Name == name {
			return &radiusCatalog[i]
		}
	}
	return nil
}

func isAccentTheme(name string) bool {
	for _, definition := range baseColorCatalog {
		if definition.name == name {
			return false
		}
	}
	return themeDefinitionByName(name) != nil
}

func applyThemeValues(destination, overrides ThemeValues) {
	maps.Copy(destination, overrides)
}

func themeValuesEqual(a, b ThemeValues) bool {
	for _, name := range TokenNames() {
		if a[name] != b[name] {
			return false
		}
	}
	return true
}

func init() {
	if err := validateCatalog(); err != nil {
		panic(fmt.Sprintf("preset palette catalog: %v", err))
	}
}

func validateCatalog() error {
	if err := validateBaseColorCatalog(); err != nil {
		return err
	}
	if err := validateThemeCatalog(); err != nil {
		return err
	}
	if err := validateRadiusCatalog(); err != nil {
		return err
	}

	seen := make(map[[32]byte]PaletteSelection)
	for _, base := range baseColorCatalog {
		choices, err := ThemeChoices(base.name)
		if err != nil {
			return err
		}
		for _, theme := range choices {
			preset, err := resolvePalette(StyleNova, PaletteSelection{
				BaseColor: base.name,
				Theme:     theme.Name,
				Radius:    radiusCatalog[0].Name,
			})
			if err != nil {
				return err
			}
			digest := paletteDigest(preset)
			if previous, ok := seen[digest]; ok {
				return fmt.Errorf("duplicate palette resolutions %q/%q and %q/%q", previous.BaseColor, previous.Theme, base.name, theme.Name)
			}
			seen[digest] = PaletteSelection{BaseColor: base.name, Theme: theme.Name}
		}
	}
	return nil
}

func validateBaseColorCatalog() error {
	seen := make(map[string]bool, len(baseColorCatalog))
	for _, definition := range baseColorCatalog {
		if definition.name == "" || definition.title == "" || definition.swatch == "" {
			return fmt.Errorf("base color has incomplete metadata")
		}
		if seen[definition.name] {
			return fmt.Errorf("duplicate base color %q", definition.name)
		}
		seen[definition.name] = true
		if err := validateLayer(definition.name+".light", definition.light, baseLayerTokens); err != nil {
			return err
		}
		if err := validateLayer(definition.name+".dark", definition.dark, baseLayerTokens); err != nil {
			return err
		}
	}
	return nil
}

func validateThemeCatalog() error {
	seen := make(map[string]bool, len(themeCatalog))
	for _, definition := range themeCatalog {
		if definition.name == "" || definition.title == "" || definition.swatch == "" {
			return fmt.Errorf("theme has incomplete metadata")
		}
		if seen[definition.name] {
			return fmt.Errorf("duplicate theme %q", definition.name)
		}
		seen[definition.name] = true
		if err := validateLayer(definition.name+".light", definition.light, themeLayerTokens); err != nil {
			return err
		}
		if err := validateLayer(definition.name+".dark", definition.dark, themeLayerTokens); err != nil {
			return err
		}
	}
	for _, base := range baseColorCatalog {
		if !seen[base.name] {
			return fmt.Errorf("missing theme definition for base color %q", base.name)
		}
	}
	return nil
}

func validateRadiusCatalog() error {
	seen := make(map[string]bool, len(radiusCatalog))
	seenValues := make(map[string]bool, len(radiusCatalog))
	for _, choice := range radiusCatalog {
		if choice.Name == "" || choice.Title == "" || choice.Value == "" {
			return fmt.Errorf("radius has incomplete metadata")
		}
		if seen[choice.Name] {
			return fmt.Errorf("duplicate radius %q", choice.Name)
		}
		if seenValues[choice.Value] {
			return fmt.Errorf("duplicate radius value %q", choice.Value)
		}
		seen[choice.Name] = true
		seenValues[choice.Value] = true
		if err := validateRadius(choice.Value); err != nil {
			return fmt.Errorf("radius %q: %w", choice.Name, err)
		}
	}
	return nil
}

func validateLayer(path string, values ThemeValues, allowed map[string]bool) error {
	for _, name := range TokenNames() {
		if allowed[name] {
			if _, ok := values[name]; ok {
				continue
			}
			return fmt.Errorf("%s.%s: missing required token", path, name)
		}
	}
	var extra []string
	for name, value := range values {
		if !isTokenName(name) {
			extra = append(extra, name)
			continue
		}
		if !allowed[name] {
			extra = append(extra, name)
			continue
		}
		if _, err := csscolorparser.Parse(value); err != nil {
			return fmt.Errorf("%s.%s: invalid CSS color: %w", path, name, err)
		}
	}
	slices.Sort(extra)
	if len(extra) != 0 {
		return fmt.Errorf("%s.%s: token is outside its layer", path, extra[0])
	}
	return nil
}

func paletteDigest(preset Preset) [32]byte {
	var values strings.Builder
	for _, mode := range []ThemeValues{preset.Theme.Light, preset.Theme.Dark} {
		for _, name := range TokenNames() {
			fmt.Fprintf(&values, "%s=%s\n", name, mode[name])
		}
	}
	return sha256.Sum256([]byte(values.String()))
}
