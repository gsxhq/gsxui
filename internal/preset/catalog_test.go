package preset

import (
	"crypto/sha256"
	"fmt"
	"slices"
	"strings"
	"testing"
)

var expectedBaseColorNames = []string{
	"neutral", "stone", "zinc", "mauve", "olive", "mist", "taupe",
}

var expectedAccentNames = []string{
	"amber", "blue", "cyan", "emerald", "fuchsia", "green", "indigo",
	"lime", "orange", "pink", "purple", "red", "rose", "sky", "teal",
	"violet", "yellow",
}

var expectedRadiusChoices = []RadiusChoice{
	{Name: "none", Title: "None", Value: "0"},
	{Name: "small", Title: "Small", Value: "0.45rem"},
	{Name: "medium", Title: "Medium", Value: "0.625rem"},
	{Name: "large", Title: "Large", Value: "0.875rem"},
}

func TestPaletteCatalogCardinality(t *testing.T) {
	if got := len(BaseColorChoices()); got != 7 {
		t.Fatalf("base colors = %d, want 7", got)
	}
	themes, err := ThemeChoices("neutral")
	if err != nil {
		t.Fatal(err)
	}
	if got := len(themes); got != 18 {
		t.Fatalf("neutral themes = %d, want selected base plus 17 accents", got)
	}
}

func TestPaletteCatalogChoices(t *testing.T) {
	baseColors := BaseColorChoices()
	if got := choiceNames(baseColors); !slices.Equal(got, expectedBaseColorNames) {
		t.Fatalf("base color names = %#v, want %#v", got, expectedBaseColorNames)
	}
	for _, choice := range baseColors {
		if choice.Title == "" || choice.Swatch == "" {
			t.Fatalf("base color choice is incomplete: %#v", choice)
		}
	}

	for _, baseColor := range expectedBaseColorNames {
		t.Run(baseColor, func(t *testing.T) {
			themes, err := ThemeChoices(baseColor)
			if err != nil {
				t.Fatal(err)
			}
			want := append([]string{baseColor}, expectedAccentNames...)
			if got := choiceNames(themes); !slices.Equal(got, want) {
				t.Fatalf("theme names = %#v, want %#v", got, want)
			}
			for _, choice := range themes {
				if choice.Title == "" || choice.Swatch == "" {
					t.Fatalf("theme choice is incomplete: %#v", choice)
				}
			}
		})
	}

	if _, err := ThemeChoices("unknown"); err == nil {
		t.Fatal("ThemeChoices(unknown) returned nil error")
	}

	if got := RadiusChoices(); !slices.Equal(got, expectedRadiusChoices) {
		t.Fatalf("radius choices = %#v, want %#v", got, expectedRadiusChoices)
	}
}

func TestPaletteCatalogChoiceSlicesAreIndependent(t *testing.T) {
	baseColors := BaseColorChoices()
	baseColors[0].Name = "changed"
	if got := BaseColorChoices()[0].Name; got != "neutral" {
		t.Fatalf("BaseColorChoices returned shared storage: %q", got)
	}

	themes, err := ThemeChoices("neutral")
	if err != nil {
		t.Fatal(err)
	}
	themes[0].Name = "changed"
	themesAgain, err := ThemeChoices("neutral")
	if err != nil {
		t.Fatal(err)
	}
	if got := themesAgain[0].Name; got != "neutral" {
		t.Fatalf("ThemeChoices returned shared storage: %q", got)
	}

	radii := RadiusChoices()
	radii[0].Name = "changed"
	if got := RadiusChoices()[0].Name; got != "none" {
		t.Fatalf("RadiusChoices returned shared storage: %q", got)
	}
}

func TestCatalogValuesRejectDuplicateTokenPairs(t *testing.T) {
	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatal("values accepted duplicate token pairs")
		}
		if got := fmt.Sprint(recovered); got != `catalog values: duplicate token "primary"` {
			t.Fatalf("values panic = %q", got)
		}
	}()

	values(
		"primary", "oklch(0.2 0 0)",
		"primary", "oklch(0.8 0 0)",
	)
}

func TestPaletteCatalogResolvesAndMatchesEveryCombination(t *testing.T) {
	seenPalettes := make(map[[32]byte]PaletteSelection)
	for _, style := range Styles() {
		for _, baseColor := range expectedBaseColorNames {
			themes, err := ThemeChoices(baseColor)
			if err != nil {
				t.Fatal(err)
			}
			for _, theme := range themes {
				for _, radius := range RadiusChoices() {
					for _, accent := range MenuAccentChoices() {
						// ChartColor is fixed at "neutral" and FontSans/FontHeading
						// at "geist" here to keep this sweep's runtime bounded — it
						// already covers style x baseColor x theme x radius x
						// menuAccent. TestChartColorCatalogRoundTrips and
						// TestFontCatalogRoundTrips below cover those two axes on
						// their own.
						selection := PaletteSelection{
							BaseColor:   baseColor,
							Theme:       theme.Name,
							Radius:      radius.Name,
							ChartColor:  "neutral",
							MenuAccent:  accent.Name,
							FontSans:    "geist",
							FontHeading: "geist",
						}
						resolved, err := ResolvePalette(style, selection)
						if err != nil {
							t.Fatalf("ResolvePalette(%q, %#v): %v", style, selection, err)
						}
						if err := Validate(resolved); err != nil {
							t.Fatalf("Validate(ResolvePalette(%q, %#v)): %v", style, selection, err)
						}
						if got := MatchPalette(resolved); got != selection {
							t.Fatalf("MatchPalette(ResolvePalette(%q, %#v)) = %#v", style, selection, got)
						}

						if style == StyleNova && radius.Name == "none" {
							digest := hashPalette(resolved)
							if prior, ok := seenPalettes[digest]; ok {
								t.Fatalf("duplicate palette resolution: %#v and %#v", prior, selection)
							}
							seenPalettes[digest] = selection
						}
					}
				}
			}
		}
	}
}

func TestDefaultResolvesCanonicalPaletteSelection(t *testing.T) {
	want := PaletteSelection{BaseColor: "neutral", Theme: "neutral", Radius: "medium", ChartColor: "neutral", MenuAccent: "subtle", FontSans: "geist", FontHeading: "geist"}
	if got := DefaultPaletteSelection(); got != want {
		t.Fatalf("DefaultPaletteSelection() = %#v, want %#v", got, want)
	}
	for _, style := range Styles() {
		if got := MatchPalette(Default(style)); got != want {
			t.Fatalf("MatchPalette(Default(%q)) = %#v, want %#v", style, got, want)
		}
	}
}

func TestMatchPaletteReportsIndependentCustomSelections(t *testing.T) {
	preset, err := ResolvePalette(StyleNova, DefaultPaletteSelection())
	if err != nil {
		t.Fatal(err)
	}
	preset.Theme.Light["primary"] = "oklch(0.5 0.2 250)"
	want := PaletteSelection{BaseColor: CustomChoice, Theme: CustomChoice, Radius: "medium", ChartColor: "neutral", MenuAccent: "subtle", FontSans: "geist", FontHeading: "geist"}
	if got := MatchPalette(preset); got != want {
		t.Fatalf("MatchPalette(custom palette) = %#v", got)
	}

	preset, err = ResolvePalette(StyleNova, DefaultPaletteSelection())
	if err != nil {
		t.Fatal(err)
	}
	preset.Radius = "1rem"
	want = PaletteSelection{BaseColor: "neutral", Theme: "neutral", Radius: CustomChoice, ChartColor: "neutral", MenuAccent: "subtle", FontSans: "geist", FontHeading: "geist"}
	if got := MatchPalette(preset); got != want {
		t.Fatalf("MatchPalette(custom radius) = %#v", got)
	}

	preset, err = ResolvePalette(StyleNova, PaletteSelection{BaseColor: "neutral", Theme: "blue", Radius: "medium", ChartColor: "neutral", MenuAccent: "bold", FontSans: "geist", FontHeading: "geist"})
	if err != nil {
		t.Fatal(err)
	}
	want = PaletteSelection{BaseColor: "neutral", Theme: "blue", Radius: "medium", ChartColor: "neutral", MenuAccent: "bold", FontSans: "geist", FontHeading: "geist"}
	if got := MatchPalette(preset); got != want {
		t.Fatalf("MatchPalette(bold accent) = %#v, want %#v", got, want)
	}

	preset, err = ResolvePalette(StyleNova, PaletteSelection{BaseColor: "neutral", Theme: "neutral", Radius: "medium", ChartColor: "violet", MenuAccent: "subtle", FontSans: "geist", FontHeading: "geist"})
	if err != nil {
		t.Fatal(err)
	}
	want = PaletteSelection{BaseColor: "neutral", Theme: "neutral", Radius: "medium", ChartColor: "violet", MenuAccent: "subtle", FontSans: "geist", FontHeading: "geist"}
	if got := MatchPalette(preset); got != want {
		t.Fatalf("MatchPalette(custom chart color) = %#v, want %#v", got, want)
	}
}

// TestChartColorCatalogRoundTrips exercises the chart color axis
// separately from TestPaletteCatalogResolvesAndMatchesEveryCombination,
// which fixes ChartColor at "neutral" to keep that sweep's runtime
// bounded.
func TestChartColorCatalogRoundTrips(t *testing.T) {
	if got := len(ChartColorChoices()); got != 24 {
		t.Fatalf("ChartColorChoices() = %d entries, want 24", got)
	}
	for _, chartColor := range ChartColorChoices() {
		selection := PaletteSelection{BaseColor: "stone", Theme: "blue", Radius: "large", ChartColor: chartColor.Name, MenuAccent: "subtle", FontSans: "geist", FontHeading: "geist"}
		resolved, err := ResolvePalette(StyleSera, selection)
		if err != nil {
			t.Fatalf("ResolvePalette(%#v): %v", selection, err)
		}
		light, dark, err := ResolveChartColor(chartColor.Name)
		if err != nil {
			t.Fatal(err)
		}
		for _, name := range []string{"chart-1", "chart-2", "chart-3", "chart-4", "chart-5"} {
			if resolved.Theme.Light[name] != light[name] {
				t.Fatalf("ResolvePalette(%#v).Theme.Light[%q] = %q, want %q", selection, name, resolved.Theme.Light[name], light[name])
			}
			if resolved.Theme.Dark[name] != dark[name] {
				t.Fatalf("ResolvePalette(%#v).Theme.Dark[%q] = %q, want %q", selection, name, resolved.Theme.Dark[name], dark[name])
			}
		}
		if got := MatchPalette(resolved); got != selection {
			t.Fatalf("MatchPalette(ResolvePalette(%#v)) = %#v", selection, got)
		}
	}
}

// expectedFontNames mirrors fontCatalog's own order in catalog.go — "geist"
// first is load-bearing (compact.go's backward-compatible default index).
var expectedFontNames = []string{
	"geist", "inter", "figtree", "jetbrains-mono", "noto-sans", "playfair-display",
}

// TestFontCatalogRoundTrips exercises the font axis (FontSans and
// FontHeading independently) separately from
// TestPaletteCatalogResolvesAndMatchesEveryCombination, which fixes both at
// "geist" to keep that sweep's runtime bounded — mirrors
// TestChartColorCatalogRoundTrips's own structure for its axis.
func TestFontCatalogRoundTrips(t *testing.T) {
	choices := FontChoices()
	if got := choiceNamesFont(choices); !slices.Equal(got, expectedFontNames) {
		t.Fatalf("FontChoices() names = %#v, want %#v", got, expectedFontNames)
	}
	for _, choice := range choices {
		if choice.Title == "" || choice.Stack == "" {
			t.Fatalf("font choice is incomplete: %#v", choice)
		}
	}

	for _, fontSans := range choices {
		for _, fontHeading := range choices {
			selection := PaletteSelection{
				BaseColor: "stone", Theme: "blue", Radius: "large", ChartColor: "neutral", MenuAccent: "subtle",
				FontSans: fontSans.Name, FontHeading: fontHeading.Name,
			}
			resolved, err := ResolvePalette(StyleSera, selection)
			if err != nil {
				t.Fatalf("ResolvePalette(%#v): %v", selection, err)
			}
			if resolved.FontSans != fontSans.Stack {
				t.Fatalf("ResolvePalette(%#v).FontSans = %q, want %q", selection, resolved.FontSans, fontSans.Stack)
			}
			if resolved.FontHeading != fontHeading.Stack {
				t.Fatalf("ResolvePalette(%#v).FontHeading = %q, want %q", selection, resolved.FontHeading, fontHeading.Stack)
			}
			if got := MatchPalette(resolved); got != selection {
				t.Fatalf("MatchPalette(ResolvePalette(%#v)) = %#v", selection, got)
			}
		}
	}
}

func TestResolvePaletteRejectsUnknownFontNames(t *testing.T) {
	base := DefaultPaletteSelection()

	sans := base
	sans.FontSans = "comic-sans"
	if _, err := ResolvePalette(StyleNova, sans); err == nil {
		t.Fatal("ResolvePalette(unknown FontSans) returned nil error")
	}

	heading := base
	heading.FontHeading = "comic-sans"
	if _, err := ResolvePalette(StyleNova, heading); err == nil {
		t.Fatal("ResolvePalette(unknown FontHeading) returned nil error")
	}
}

func choiceNamesFont(choices []FontChoice) []string {
	names := make([]string, len(choices))
	for i, choice := range choices {
		names[i] = choice.Name
	}
	return names
}

func TestResolveChartColorRejectsUnknownNames(t *testing.T) {
	if _, _, err := ResolveChartColor("plaid"); err == nil {
		t.Fatal("ResolveChartColor(unknown) returned nil error")
	}
}

func TestResolveChartColorReturnsIndependentStorage(t *testing.T) {
	light, _, err := ResolveChartColor("blue")
	if err != nil {
		t.Fatal(err)
	}
	light["chart-1"] = "changed"
	againLight, _, err := ResolveChartColor("blue")
	if err != nil {
		t.Fatal(err)
	}
	if againLight["chart-1"] == "changed" {
		t.Fatal("ResolveChartColor returned shared storage")
	}
}

func choiceNames(choices []PaletteChoice) []string {
	names := make([]string, len(choices))
	for i, choice := range choices {
		names[i] = choice.Name
	}
	return names
}

func hashPalette(preset Preset) [32]byte {
	var data strings.Builder
	for _, mode := range []ThemeValues{preset.Theme.Light, preset.Theme.Dark} {
		for _, name := range TokenNames() {
			fmt.Fprintf(&data, "%s=%s\n", name, mode[name])
		}
	}
	return sha256.Sum256([]byte(data.String()))
}
