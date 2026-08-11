package preset

import (
	"errors"
	"fmt"
	"math"
	"slices"
)

const (
	fullSharePrefix    = "gsxui:v1:"
	compactSharePrefix = "gsxui:p1:"
	base62Alphabet     = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	compactStyleBits     = 4
	compactBaseColorBits = 4
	compactThemeBits     = 5
	compactRadiusBits    = 3

	compactBaseColorShift = compactStyleBits
	compactThemeShift     = compactBaseColorShift + compactBaseColorBits
	compactRadiusShift    = compactThemeShift + compactThemeBits
	compactSchemaBits     = compactRadiusShift + compactRadiusBits
)

// These arrays are the append-only p1 transport ABI. They intentionally do
// not derive their order from presentation catalogues.
var compactStyles = []Style{
	StyleNova,
	StyleMaia,
	StyleVega,
	StyleLyra,
	StyleMira,
	StyleLuma,
	StyleSera,
	StyleRhea,
}

var compactBaseColors = []string{
	"neutral",
	"stone",
	"zinc",
	"mauve",
	"olive",
	"mist",
	"taupe",
}

var compactThemes = []string{
	"neutral",
	"stone",
	"zinc",
	"mauve",
	"olive",
	"mist",
	"taupe",
	"amber",
	"blue",
	"cyan",
	"emerald",
	"fuchsia",
	"green",
	"indigo",
	"lime",
	"orange",
	"pink",
	"purple",
	"red",
	"rose",
	"sky",
	"teal",
	"violet",
	"yellow",
}

var compactRadii = []string{
	"none",
	"small",
	"medium",
	"large",
}

type ShareSchema struct {
	FullPrefix    string
	CompactPrefix string
	Styles        []Style
	BaseColors    []string
	Themes        []string
	Radii         []string
}

func ShareTransportSchema() ShareSchema {
	return ShareSchema{
		FullPrefix:    fullSharePrefix,
		CompactPrefix: compactSharePrefix,
		Styles:        slices.Clone(compactStyles),
		BaseColors:    slices.Clone(compactBaseColors),
		Themes:        slices.Clone(compactThemes),
		Radii:         slices.Clone(compactRadii),
	}
}

func encodeCompactShare(value Preset) (string, bool, error) {
	if err := Validate(value); err != nil {
		return "", false, err
	}
	selection := MatchPalette(value)
	if selection.BaseColor == CustomChoice || selection.Theme == CustomChoice ||
		selection.Radius == CustomChoice {
		return "", false, nil
	}

	styleIndex := indexStyle(compactStyles, value.Style)
	baseColorIndex := indexStringValue(compactBaseColors, selection.BaseColor)
	themeIndex := indexStringValue(compactThemes, selection.Theme)
	radiusIndex := indexStringValue(compactRadii, selection.Radius)
	if styleIndex == -1 || baseColorIndex == -1 || themeIndex == -1 || radiusIndex == -1 {
		return "", false, nil
	}

	packed := uint64(styleIndex) |
		uint64(baseColorIndex)<<compactBaseColorShift |
		uint64(themeIndex)<<compactThemeShift |
		uint64(radiusIndex)<<compactRadiusShift
	return compactSharePrefix + encodeBase62(packed), true, nil
}

func decodeCompactShare(payload string) (Preset, error) {
	if payload == "" {
		return Preset{}, errors.New("compact payload is empty")
	}
	packed, err := decodeBase62(payload)
	if err != nil {
		return Preset{}, err
	}
	if encodeBase62(packed) != payload {
		return Preset{}, errors.New("compact base62 spelling is not canonical")
	}
	if packed>>compactSchemaBits != 0 {
		return Preset{}, errors.New("compact payload uses bits beyond the p1 schema")
	}

	styleIndex := int(packed & compactMask(compactStyleBits))
	baseColorIndex := int((packed >> compactBaseColorShift) & compactMask(compactBaseColorBits))
	themeIndex := int((packed >> compactThemeShift) & compactMask(compactThemeBits))
	radiusIndex := int((packed >> compactRadiusShift) & compactMask(compactRadiusBits))

	if styleIndex >= len(compactStyles) {
		return Preset{}, fmt.Errorf("compact payload uses unused style index %d", styleIndex)
	}
	if baseColorIndex >= len(compactBaseColors) {
		return Preset{}, fmt.Errorf("compact payload uses unused base color index %d", baseColorIndex)
	}
	if themeIndex >= len(compactThemes) {
		return Preset{}, fmt.Errorf("compact payload uses unused theme index %d", themeIndex)
	}
	if radiusIndex >= len(compactRadii) {
		return Preset{}, fmt.Errorf("compact payload uses unused radius index %d", radiusIndex)
	}

	selection := PaletteSelection{
		BaseColor: compactBaseColors[baseColorIndex],
		Theme:     compactThemes[themeIndex],
		Radius:    compactRadii[radiusIndex],
	}
	value, err := ResolvePalette(compactStyles[styleIndex], selection)
	if err != nil {
		return Preset{}, fmt.Errorf(
			"compact payload has invalid base-color/theme combination %q/%q: %w",
			selection.BaseColor,
			selection.Theme,
			err,
		)
	}
	return value, nil
}

func compactMask(bits int) uint64 {
	return (1 << bits) - 1
}

func encodeBase62(value uint64) string {
	if value == 0 {
		return "0"
	}
	var reversed [11]byte
	cursor := len(reversed)
	for value > 0 {
		cursor--
		reversed[cursor] = base62Alphabet[value%62]
		value /= 62
	}
	return string(reversed[cursor:])
}

func decodeBase62(value string) (uint64, error) {
	var decoded uint64
	for _, character := range []byte(value) {
		digit := indexByte(base62Alphabet, character)
		if digit == -1 {
			return 0, fmt.Errorf("compact payload contains invalid base62 character %q", character)
		}
		if decoded > (math.MaxUint64-uint64(digit))/62 {
			return 0, errors.New("compact base62 integer overflow")
		}
		decoded = decoded*62 + uint64(digit)
	}
	return decoded, nil
}

func indexStyle(values []Style, value Style) int {
	for i, candidate := range values {
		if candidate == value {
			return i
		}
	}
	return -1
}

func indexStringValue(values []string, value string) int {
	for i, candidate := range values {
		if candidate == value {
			return i
		}
	}
	return -1
}

func indexByte(values string, value byte) int {
	for i := range len(values) {
		if values[i] == value {
			return i
		}
	}
	return -1
}
