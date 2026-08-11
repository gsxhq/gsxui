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
						selection := PaletteSelection{
							BaseColor:  baseColor,
							Theme:      theme.Name,
							Radius:     radius.Name,
							MenuAccent: accent.Name,
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
	want := PaletteSelection{BaseColor: "neutral", Theme: "neutral", Radius: "medium", MenuAccent: "subtle"}
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
	want := PaletteSelection{BaseColor: CustomChoice, Theme: CustomChoice, Radius: "medium", MenuAccent: "subtle"}
	if got := MatchPalette(preset); got != want {
		t.Fatalf("MatchPalette(custom palette) = %#v", got)
	}

	preset, err = ResolvePalette(StyleNova, DefaultPaletteSelection())
	if err != nil {
		t.Fatal(err)
	}
	preset.Radius = "1rem"
	want = PaletteSelection{BaseColor: "neutral", Theme: "neutral", Radius: CustomChoice, MenuAccent: "subtle"}
	if got := MatchPalette(preset); got != want {
		t.Fatalf("MatchPalette(custom radius) = %#v", got)
	}

	preset, err = ResolvePalette(StyleNova, PaletteSelection{BaseColor: "neutral", Theme: "blue", Radius: "medium", MenuAccent: "bold"})
	if err != nil {
		t.Fatal(err)
	}
	want = PaletteSelection{BaseColor: "neutral", Theme: "blue", Radius: "medium", MenuAccent: "bold"}
	if got := MatchPalette(preset); got != want {
		t.Fatalf("MatchPalette(bold accent) = %#v, want %#v", got, want)
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
