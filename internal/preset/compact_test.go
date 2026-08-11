package preset

import (
	"strings"
	"testing"
)

var compactTransportStyles = []Style{
	StyleNova, StyleMaia, StyleVega, StyleLyra, StyleMira, StyleLuma, StyleSera, StyleRhea,
}

var compactTransportBaseColors = []string{
	"neutral", "stone", "zinc", "mauve", "olive", "mist", "taupe",
}

var compactTransportThemes = []string{
	"neutral", "stone", "zinc", "mauve", "olive", "mist", "taupe",
	"amber", "blue", "cyan", "emerald", "fuchsia", "green", "indigo",
	"lime", "orange", "pink", "purple", "red", "rose", "sky", "teal",
	"violet", "yellow",
}

var compactTransportRadii = []string{"none", "small", "medium", "large"}

func TestCompactShareExactCodes(t *testing.T) {
	tests := []struct {
		name      string
		style     Style
		selection PaletteSelection
		want      string
	}{
		{
			name:      "zero boundary",
			style:     StyleNova,
			selection: PaletteSelection{BaseColor: "neutral", Theme: "neutral", Radius: "none", MenuAccent: "subtle"},
			want:      "gsxui:p1:0",
		},
		{
			name:      "default nova",
			style:     StyleNova,
			selection: DefaultPaletteSelection(),
			want:      "gsxui:p1:4GG",
		},
		{
			name:      "default maia",
			style:     StyleMaia,
			selection: DefaultPaletteSelection(),
			want:      "gsxui:p1:4GH",
		},
		{
			name:      "upper boundary",
			style:     StyleMaia,
			selection: PaletteSelection{BaseColor: "taupe", Theme: "yellow", Radius: "large", MenuAccent: "bold"},
			want:      "gsxui:p1:Ozx",
		},
		{
			// "H32" (2^16) used to be rejected by TestDecodeCompactShareRejectsInvalidTransport
			// as beyond the pre-Task-1c 16-bit schema; menuAccent's bit now
			// claims that exact position, so it decodes to bold accent with
			// every other axis at its catalog default — a concrete
			// illustration of the append-only rule at compact.go:26-27.
			name:      "menu accent bold only",
			style:     StyleNova,
			selection: PaletteSelection{BaseColor: "neutral", Theme: "neutral", Radius: "none", MenuAccent: "bold"},
			want:      "gsxui:p1:H32",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := ResolvePalette(tt.style, tt.selection)
			if err != nil {
				t.Fatal(err)
			}
			got, err := EncodeShare(value)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Fatalf("EncodeShare = %q, want %q", got, tt.want)
			}
			decoded, err := DecodeShare(got)
			if err != nil {
				t.Fatal(err)
			}
			if !presetsEqual(decoded, value) {
				t.Fatalf("DecodeShare(%q) = %#v, want %#v", got, decoded, value)
			}
		})
	}
}

func TestCompactShareRoundTripsEveryCatalogCombination(t *testing.T) {
	for _, style := range compactTransportStyles {
		for _, baseColor := range compactTransportBaseColors {
			for _, theme := range compactTransportThemes {
				selection := PaletteSelection{
					BaseColor:  baseColor,
					Theme:      theme,
					Radius:     compactTransportRadii[0],
					MenuAccent: "subtle",
				}
				if themeIndex := indexString(compactTransportThemes, theme); themeIndex < len(compactTransportBaseColors) &&
					theme != baseColor {
					continue
				}
				for _, radius := range compactTransportRadii {
					selection.Radius = radius
					value, err := ResolvePalette(style, selection)
					if err != nil {
						t.Fatalf("ResolvePalette(%q, %#v): %v", style, selection, err)
					}
					code, err := EncodeShare(value)
					if err != nil {
						t.Fatalf("EncodeShare(%q, %#v): %v", style, selection, err)
					}
					if !strings.HasPrefix(code, "gsxui:p1:") {
						t.Fatalf("EncodeShare(%q, %#v) = %q, want compact code", style, selection, code)
					}
					if len(code) > 12 {
						t.Fatalf("compact code %q has length %d, want at most 12", code, len(code))
					}
					got, err := DecodeShare(code)
					if err != nil {
						t.Fatalf("DecodeShare(%q): %v", code, err)
					}
					if !presetsEqual(got, value) {
						t.Fatalf("DecodeShare(%q) = %#v, want %#v", code, got, value)
					}
					reencoded, err := EncodeShare(got)
					if err != nil {
						t.Fatalf("EncodeShare(DecodeShare(%q)): %v", code, err)
					}
					if reencoded != code {
						t.Fatalf("EncodeShare(DecodeShare(%q)) = %q, want the original code; the catalog is no longer injective", code, reencoded)
					}
				}
			}
		}
	}
}

// TestCompactShareMenuAccentRoundTrips exercises the menu accent axis
// separately from TestCompactShareRoundTripsEveryCatalogCombination, which
// fixes MenuAccent at "subtle" to keep the base/theme/radius/style
// combinatorial sweep's runtime bounded.
func TestCompactShareMenuAccentRoundTrips(t *testing.T) {
	for _, accent := range compactMenuAccents {
		selection := PaletteSelection{BaseColor: "taupe", Theme: "yellow", Radius: "large", MenuAccent: accent}
		value, err := ResolvePalette(StyleRhea, selection)
		if err != nil {
			t.Fatalf("ResolvePalette(%#v): %v", selection, err)
		}
		code, err := EncodeShare(value)
		if err != nil {
			t.Fatalf("EncodeShare(%#v): %v", selection, err)
		}
		if !strings.HasPrefix(code, "gsxui:p1:") {
			t.Fatalf("EncodeShare(%#v) = %q, want compact code", selection, code)
		}
		got, err := DecodeShare(code)
		if err != nil {
			t.Fatalf("DecodeShare(%q): %v", code, err)
		}
		if !presetsEqual(got, value) {
			t.Fatalf("DecodeShare(%q) = %#v, want %#v", code, got, value)
		}
		if got := MatchPalette(got).MenuAccent; got != accent {
			t.Fatalf("MatchPalette(DecodeShare(%q)).MenuAccent = %q, want %q", code, got, accent)
		}
	}
}

func TestShareTransportSchemaPinsAppendOnlyABI(t *testing.T) {
	schema := ShareTransportSchema()
	if schema.FullPrefix != "gsxui:v1:" {
		t.Errorf("full prefix = %q, want gsxui:v1:", schema.FullPrefix)
	}
	if schema.CompactPrefix != "gsxui:p1:" {
		t.Errorf("compact prefix = %q, want gsxui:p1:", schema.CompactPrefix)
	}
	assertStylesEqual(t, schema.Styles, compactTransportStyles)
	assertStringsEqual(t, "base colors", schema.BaseColors, compactTransportBaseColors)
	assertStringsEqual(t, "themes", schema.Themes, compactTransportThemes)
	assertStringsEqual(t, "radii", schema.Radii, compactTransportRadii)
	assertStringsEqual(t, "menu accents", schema.MenuAccents, []string{"subtle", "bold"})

	schema.Styles[0] = "changed"
	schema.BaseColors[0] = "changed"
	schema.Themes[0] = "changed"
	schema.Radii[0] = "changed"
	schema.MenuAccents[0] = "changed"
	again := ShareTransportSchema()
	if again.Styles[0] != StyleNova || again.BaseColors[0] != "neutral" ||
		again.Themes[0] != "neutral" || again.Radii[0] != "none" || again.MenuAccents[0] != "subtle" {
		t.Fatal("ShareTransportSchema returned shared ABI storage")
	}
}

func TestShareEncodingFallsBackToFullTransportForCustomValues(t *testing.T) {
	customTheme := Default(StyleNova)
	customTheme.Theme.Light["primary"] = "oklch(0.51 0.21 20)"

	customRadius := Default(StyleNova)
	customRadius.Radius = "1rem"

	for _, tt := range []struct {
		name  string
		value Preset
	}{
		{name: "theme", value: customTheme},
		{name: "radius", value: customRadius},
	} {
		t.Run(tt.name, func(t *testing.T) {
			code, err := EncodeShare(tt.value)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.HasPrefix(code, "gsxui:v1:") {
				t.Fatalf("EncodeShare = %q, want full v1 code", code)
			}
			got, err := DecodeShare(code)
			if err != nil {
				t.Fatal(err)
			}
			if !presetsEqual(got, tt.value) {
				t.Fatalf("DecodeShare = %#v, want %#v", got, tt.value)
			}
		})
	}
}

func TestDecodeCompactShareRejectsInvalidTransport(t *testing.T) {
	tests := []struct {
		name string
		code string
		want string
	}{
		{name: "empty payload", code: "gsxui:p1:", want: "empty"},
		{name: "invalid character", code: "gsxui:p1:!", want: "base62"},
		{name: "integer overflow", code: "gsxui:p1:zzzzzzzzzzzz", want: "overflow"},
		{name: "unused style index", code: "gsxui:p1:8", want: "style"},
		{name: "unused base color index", code: "gsxui:p1:1o", want: "base color"},
		{name: "unused theme index", code: "gsxui:p1:1bc", want: "theme"},
		{name: "unused radius index", code: "gsxui:p1:8WW", want: "radius"},
		{name: "unused menu accent index", code: "gsxui:p1:Y64", want: "menu accent"},
		{name: "invalid base and theme pair", code: "gsxui:p1:G", want: "combination"},
		{name: "bits beyond schema", code: "gsxui:p1:16C8", want: "schema"},
		{name: "leading zero", code: "gsxui:p1:04GG", want: "canonical"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := DecodeShare(tt.code); err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("DecodeShare(%q) error = %v, want %q", tt.code, err, tt.want)
			}
		})
	}
}

func indexString(values []string, value string) int {
	for i, candidate := range values {
		if candidate == value {
			return i
		}
	}
	return -1
}

func assertStylesEqual(t *testing.T, got, want []Style) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("styles = %#v, want %#v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("styles = %#v, want %#v", got, want)
		}
	}
}

func assertStringsEqual(t *testing.T, name string, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s = %#v, want %#v", name, got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("%s = %#v, want %#v", name, got, want)
		}
	}
}
