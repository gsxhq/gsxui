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

// MenuAccentChoice is a picker-metadata entry for the menu accent axis
// (catalog.go's own analogue of PaletteChoice/RadiusChoice above): unlike
// those, it carries no Swatch — "subtle" and "bold" are not hues, they are
// which token pair (base accent vs. primary) accent/accent-foreground
// borrow from. See applyMenuAccent.
type MenuAccentChoice struct {
	Name  string
	Title string
}

var menuAccentCatalog = []MenuAccentChoice{
	{Name: "subtle", Title: "Subtle"},
	{Name: "bold", Title: "Bold"},
}

// FontChoice is a picker-metadata entry for the font axis (catalog.go's own
// analogue of PaletteChoice/RadiusChoice/MenuAccentChoice above): Stack is
// the literal CSS font-family value the picker commits into
// Preset.FontSans/FontHeading when this choice is selected — see
// resolvePalette's font handling below. Unlike PaletteChoice it carries no
// Swatch: a font isn't a hue.
type FontChoice struct {
	Name  string
	Title string
	Stack string
}

// fontCatalog is gsxui's own curation, not upstream-audited data (unlike
// baseColorCatalog/themeCatalog in catalog_data.go): six self-hosted
// variable fonts chosen to cover the 8 shipped styles' typographic
// character (geometric sans, humanist sans, mono, and a serif for sera's
// editorial style) without the ~26-font/next-font-google/CDN sprawl
// upstream's own picker offers — see docs/superpowers/plans/
// 2026-08-11-theme-creator-parity.md Task 3b. "geist" is first (index 0)
// deliberately: compact.go's wire format relies on index 0 being the
// default choice so pre-Task-3 share codes, whose bits for this axis never
// existed, decode to today's actual default rather than an arbitrary
// entry — see compact.go's own comment on compactFonts.
//
// Each Stack names a family already registered by a hand-authored
// @font-face rule in assets/fonts/fonts.css (self-hosted variable WOFF2,
// font-display: swap — no npm, no CDN, per the plan's Global Constraints).
// The trailing generic-family fallback is Tailwind's own default stack for
// that font's character (ui-sans-serif/ui-monospace/ui-serif), so a
// project that never loads fonts.css still renders a sane system font
// instead of an invalid declaration.
var fontCatalog = []FontChoice{
	{Name: "geist", Title: "Geist", Stack: `"Geist Variable", ui-sans-serif, system-ui, sans-serif`},
	{Name: "inter", Title: "Inter", Stack: `"Inter Variable", ui-sans-serif, system-ui, sans-serif`},
	{Name: "figtree", Title: "Figtree", Stack: `"Figtree Variable", ui-sans-serif, system-ui, sans-serif`},
	{Name: "jetbrains-mono", Title: "JetBrains Mono", Stack: `"JetBrains Mono Variable", ui-monospace, SFMono-Regular, monospace`},
	{Name: "noto-sans", Title: "Noto Sans", Stack: `"Noto Sans Variable", ui-sans-serif, system-ui, sans-serif`},
	{Name: "playfair-display", Title: "Playfair Display", Stack: `"Playfair Display Variable", ui-serif, Georgia, serif`},
}

// FontChoices returns the font axis's picker metadata, in fontCatalog's own
// order, for the theme editor's two font pickers (body/FontSans and
// heading/FontHeading both draw from this same list — mirrors upstream's
// own single font list backing both pickers, dossier §1).
func FontChoices() []FontChoice {
	return slices.Clone(fontCatalog)
}

func fontChoiceByName(name string) *FontChoice {
	for i := range fontCatalog {
		if fontCatalog[i].Name == name {
			return &fontCatalog[i]
		}
	}
	return nil
}

// matchFont looks for a catalog entry whose Stack exactly equals stack,
// returning CustomChoice if none matches (e.g. a foreign theme.css or raw
// JSON import named an arbitrary font-family not in fontCatalog) — mirrors
// matchChartColor's exact-equality catalog scan above.
func matchFont(stack string) string {
	for _, choice := range fontCatalog {
		if choice.Stack == stack {
			return choice.Name
		}
	}
	return CustomChoice
}

// MenuAccentChoices returns the menu accent axis's picker metadata. This is
// the cheap half of the upstream menu-chrome axis (registry/config.ts's
// menuAccent): a pure token override, ported directly. The color
// (default/inverted) and appearance (solid/translucent) halves are
// deliberately out of scope — they are DOM-class re-theming machinery
// (a live MutationObserver toggling classes across every menu-bearing
// component in every style), not a token swap, and gsxui's preview
// architecture only ever pushes CSS custom property values today. See
// docs/superpowers/plans/2026-08-11-theme-creator-parity.md Task 1c.
func MenuAccentChoices() []MenuAccentChoice {
	return slices.Clone(menuAccentCatalog)
}

func menuAccentChoiceByName(name string) *MenuAccentChoice {
	for i := range menuAccentCatalog {
		if menuAccentCatalog[i].Name == name {
			return &menuAccentCatalog[i]
		}
	}
	return nil
}

// accentOverrideTokens are the two tokens applyMenuAccent overwrites when
// MenuAccent is "bold". They are excluded from MatchPalette's structural
// base/theme comparison (below) because a bold override can move them away
// from the values the matched base color would otherwise have contributed,
// independent of which base color or theme is actually selected.
var accentOverrideTokens = map[string]bool{"accent": true, "accent-foreground": true}

func applyMenuAccent(values ThemeValues, accent string) {
	if accent != "bold" {
		return
	}
	values["accent"] = values["primary"]
	values["accent-foreground"] = values["primary-foreground"]
}

// chartLayerTokens is the five-token layer PaletteSelection.ChartColor
// overlays, mirroring baseLayerTokens/themeLayerTokens above. Like
// accentOverrideTokens, it is excluded from MatchPalette's base/theme
// structural comparison: ChartColor may name a different hue than Theme
// (mirrors upstream, registry/config.ts:719-729), so the chart tokens stop
// carrying the base-color/theme signal once a non-default chart color is
// chosen.
var chartLayerTokens = map[string]bool{
	"chart-1": true, "chart-2": true, "chart-3": true, "chart-4": true, "chart-5": true,
}

type PaletteSelection struct {
	BaseColor   string
	Theme       string
	Radius      string
	ChartColor  string
	MenuAccent  string
	FontSans    string
	FontHeading string
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
	// chartLight/chartDark hold ONLY the five chart-1..5 tokens for this
	// hue — a separate layer from light/dark above (the theme/brand layer:
	// primary, secondary, sidebar-primary). ResolvePalette overlays them
	// independently via PaletteSelection.ChartColor, which may name a
	// different hue than PaletteSelection.Theme (see chartLayerTokens).
	chartLight ThemeValues
	chartDark  ThemeValues
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

// ChartColorChoices returns every hue's picker metadata for the chart-color
// axis, unfiltered by base color — unlike ThemeChoices(baseColor), which
// upstream's own ChartColorPicker also filters (dossier §1: "splits
// base-color-named themes from pure accent themes into two picker
// groups"), this is the full 24-entry catalogue, useful for callers (like
// the theme editor's schema) that need every hue's resolved chart values
// regardless of which base color happens to be selected right now.
func ChartColorChoices() []PaletteChoice {
	choices := make([]PaletteChoice, len(themeCatalog))
	for i, definition := range themeCatalog {
		choices[i] = paletteChoice(definition)
	}
	return choices
}

// ResolveChartColor returns the chart-1..5 token values for the given hue
// name, independent of base color, theme, or style — the same overlay
// resolvePalette applies for PaletteSelection.ChartColor, exposed in
// isolation for callers building a picker's value table.
func ResolveChartColor(name string) (light, dark ThemeValues, err error) {
	theme := themeDefinitionByName(name)
	if theme == nil {
		return nil, nil, fmt.Errorf("chart color %q is not in the palette catalog", name)
	}
	return maps.Clone(theme.chartLight), maps.Clone(theme.chartDark), nil
}

func DefaultPaletteSelection() PaletteSelection {
	return PaletteSelection{BaseColor: "neutral", Theme: "neutral", Radius: "medium", ChartColor: "neutral", MenuAccent: "subtle", FontSans: "geist", FontHeading: "geist"}
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
		BaseColor:   CustomChoice,
		Theme:       CustomChoice,
		Radius:      CustomChoice,
		ChartColor:  CustomChoice,
		MenuAccent:  "subtle",
		FontSans:    CustomChoice,
		FontHeading: CustomChoice,
	}
	if Validate(preset) != nil {
		return selection
	}
	selection.MenuAccent = matchMenuAccent(preset)
	selection.ChartColor = matchChartColor(preset)
	selection.FontSans = matchFont(preset.FontSans)
	selection.FontHeading = matchFont(preset.FontHeading)

	// accentOverrideTokens and chartLayerTokens are excluded from this
	// comparison: a bold menu accent (just detected above) overwrites
	// accent/accent-foreground, and a non-default chart color overwrites
	// chart-1..5 — both independent of base color/theme, so neither set
	// carries the base color's own signal once that's happened. The
	// remaining ~25 compared tokens are already proven unique per
	// base-color/theme combination by validateCatalog's own digest check.
	for _, base := range baseColorCatalog {
		choices, err := ThemeChoices(base.name)
		if err != nil {
			return selection
		}
		for _, theme := range choices {
			candidate, err := resolvePalette(preset.Style, PaletteSelection{
				BaseColor:   base.name,
				Theme:       theme.Name,
				Radius:      "medium",
				ChartColor:  "neutral",
				MenuAccent:  "subtle",
				FontSans:    "geist",
				FontHeading: "geist",
			})
			if err == nil &&
				themeValuesEqualExcluding(candidate.Theme.Light, preset.Theme.Light, accentOverrideTokens, chartLayerTokens) &&
				themeValuesEqualExcluding(candidate.Theme.Dark, preset.Theme.Dark, accentOverrideTokens, chartLayerTokens) {
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

// matchMenuAccent has no "custom" outcome, unlike the other axes: bold is
// defined structurally (accent/accent-foreground literally equal
// primary/primary-foreground in both modes, applyMenuAccent's own
// postcondition), so every preset is either bold or — by default — subtle.
func matchMenuAccent(preset Preset) string {
	bold := preset.Theme.Light["accent"] == preset.Theme.Light["primary"] &&
		preset.Theme.Light["accent-foreground"] == preset.Theme.Light["primary-foreground"] &&
		preset.Theme.Dark["accent"] == preset.Theme.Dark["primary"] &&
		preset.Theme.Dark["accent-foreground"] == preset.Theme.Dark["primary-foreground"]
	if bold {
		return "bold"
	}
	return "subtle"
}

// matchChartColor looks for a hue whose chart-1..5 values (light and dark)
// exactly match the preset's, independent of which base color or theme
// also matches — ChartColor is its own axis, so this can and does return
// "custom" even when BaseColor/Theme match a catalog combination exactly.
func matchChartColor(preset Preset) string {
	for _, definition := range themeCatalog {
		if containsThemeValues(preset.Theme.Light, definition.chartLight) &&
			containsThemeValues(preset.Theme.Dark, definition.chartDark) {
			return definition.name
		}
	}
	return CustomChoice
}

// containsThemeValues reports whether every key in subset is present in
// full with an equal value. Used to match a partial layer (like a chart
// color's five tokens) against a full preset's theme map.
func containsThemeValues(full, subset ThemeValues) bool {
	for name, value := range subset {
		if full[name] != value {
			return false
		}
	}
	return true
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
		FontSans:      defaultFontFamily,
		FontHeading:   defaultFontFamily,
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
	chartColor := themeDefinitionByName(selection.ChartColor)
	if chartColor == nil {
		return Preset{}, fmt.Errorf("chart color %q is not in the palette catalog", selection.ChartColor)
	}
	radius := radiusChoiceByName(selection.Radius)
	if radius == nil {
		return Preset{}, fmt.Errorf("radius %q is not in the palette catalog", selection.Radius)
	}
	if menuAccentChoiceByName(selection.MenuAccent) == nil {
		return Preset{}, fmt.Errorf("menu accent %q is not in the palette catalog", selection.MenuAccent)
	}
	fontSans := fontChoiceByName(selection.FontSans)
	if fontSans == nil {
		return Preset{}, fmt.Errorf("font %q is not in the palette catalog", selection.FontSans)
	}
	fontHeading := fontChoiceByName(selection.FontHeading)
	if fontHeading == nil {
		return Preset{}, fmt.Errorf("font %q is not in the palette catalog", selection.FontHeading)
	}

	// Order mirrors upstream's own buildRegistryTheme (registry/config.ts:
	// 697-741): base + theme merge, then the chart overlay, then the menu
	// accent overlay last (so bold reads whatever primary the theme layer
	// just set, not a stale value). Font is independent of all of the
	// above — it never touches Theme.Light/Dark — so its position here
	// doesn't matter; it's applied last for readability only.
	preset := canonicalDefault(style)
	applyThemeValues(preset.Theme.Light, base.light)
	applyThemeValues(preset.Theme.Dark, base.dark)
	applyThemeValues(preset.Theme.Light, theme.light)
	applyThemeValues(preset.Theme.Dark, theme.dark)
	applyThemeValues(preset.Theme.Light, chartColor.chartLight)
	applyThemeValues(preset.Theme.Dark, chartColor.chartDark)
	applyMenuAccent(preset.Theme.Light, selection.MenuAccent)
	applyMenuAccent(preset.Theme.Dark, selection.MenuAccent)
	preset.Radius = radius.Value
	preset.FontSans = fontSans.Stack
	preset.FontHeading = fontHeading.Stack
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

func themeValuesEqualExcluding(a, b ThemeValues, excluded ...map[string]bool) bool {
	for _, name := range TokenNames() {
		if slices.ContainsFunc(excluded, func(set map[string]bool) bool { return set[name] }) {
			continue
		}
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
	if err := validateFontCatalog(); err != nil {
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
				BaseColor:   base.name,
				Theme:       theme.Name,
				Radius:      radiusCatalog[0].Name,
				ChartColor:  "neutral",
				MenuAccent:  "subtle",
				FontSans:    fontCatalog[0].Name,
				FontHeading: fontCatalog[0].Name,
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
	seenChart := make(map[[32]byte]string, len(themeCatalog))
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
		if err := validateLayer(definition.name+".chartLight", definition.chartLight, chartLayerTokens); err != nil {
			return err
		}
		if err := validateLayer(definition.name+".chartDark", definition.chartDark, chartLayerTokens); err != nil {
			return err
		}
		chartDigest := layerDigest(definition.chartLight, definition.chartDark)
		if previous, ok := seenChart[chartDigest]; ok {
			return fmt.Errorf("duplicate chart color resolutions %q and %q", previous, definition.name)
		}
		seenChart[chartDigest] = definition.name
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

func validateFontCatalog() error {
	if len(fontCatalog) == 0 || fontCatalog[0].Name != "geist" {
		return fmt.Errorf("font catalog: index 0 must be %q for compact.go's backward-compatible default", "geist")
	}
	seen := make(map[string]bool, len(fontCatalog))
	seenStacks := make(map[string]bool, len(fontCatalog))
	for _, choice := range fontCatalog {
		if choice.Name == "" || choice.Title == "" || choice.Stack == "" {
			return fmt.Errorf("font has incomplete metadata")
		}
		if seen[choice.Name] {
			return fmt.Errorf("duplicate font %q", choice.Name)
		}
		if seenStacks[choice.Stack] {
			return fmt.Errorf("duplicate font stack %q", choice.Stack)
		}
		seen[choice.Name] = true
		seenStacks[choice.Stack] = true
		if err := validateFontFamily(choice.Stack); err != nil {
			return fmt.Errorf("font %q: %w", choice.Name, err)
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

// layerDigest hashes an arbitrary light/dark token layer in a stable key
// order, independent of TokenNames()'s own ordering — used to detect
// duplicate chart color resolutions, which would otherwise make
// matchChartColor's catalog scan ambiguous.
func layerDigest(light, dark ThemeValues) [32]byte {
	var values strings.Builder
	for _, mode := range []ThemeValues{light, dark} {
		names := slices.Sorted(maps.Keys(mode))
		for _, name := range names {
			fmt.Fprintf(&values, "%s=%s\n", name, mode[name])
		}
	}
	return sha256.Sum256([]byte(values.String()))
}
