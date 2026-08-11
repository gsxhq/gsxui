package preset

import (
	"bytes"
	"strings"
	"testing"
)

func TestThemeCSSMatchesCanonicalGoldenAndRepositoryAsset(t *testing.T) {
	t.Parallel()

	got, err := ThemeCSS(Default(StyleNova))
	if err != nil {
		t.Fatalf("ThemeCSS: %v", err)
	}
	golden := readTestFile(t, "testdata/default-theme.css")
	if !bytes.Equal(got, golden) {
		t.Fatalf("ThemeCSS mismatch\n--- got ---\n%s\n--- want ---\n%s", got, golden)
	}
	asset := readTestFile(t, "../../assets/css/themes/default.css")
	if !bytes.Equal(got, asset) {
		t.Fatalf("repository default theme drifted\n--- generated ---\n%s\n--- asset ---\n%s", got, asset)
	}
}

func TestThemeCSSRejectsInvalidPreset(t *testing.T) {
	t.Parallel()

	preset := Default(StyleNova)
	preset.Theme.Dark["primary"] = "not-a-color"
	if _, err := ThemeCSS(preset); err == nil ||
		!strings.Contains(err.Error(), "theme.dark.primary") {
		t.Fatalf("ThemeCSS error = %v, want theme.dark.primary", err)
	}
}

func TestImportThemeCSSUpdatesRecognizedValuesAndPreservesBase(t *testing.T) {
	t.Parallel()

	base := Default(StyleMaia)
	src := []byte(`
/* compatible theme files may include unrelated CSS */
:root {
  color-scheme: light;
  --third-party: ignored;
  --primary: oklch(0.6 0.2 280);
}
.dark {
  --primary: oklch(0.7 0.2 280);
  --third-party-dark: ignored;
}
`)
	got, err := ImportThemeCSS(base, src)
	if err != nil {
		t.Fatalf("ImportThemeCSS: %v", err)
	}
	if got.Style != StyleMaia {
		t.Fatalf("Style = %q, want maia", got.Style)
	}
	if got.Theme.Light["primary"] != "oklch(0.6 0.2 280)" {
		t.Fatalf("light primary = %q", got.Theme.Light["primary"])
	}
	if got.Theme.Dark["primary"] != "oklch(0.7 0.2 280)" {
		t.Fatalf("dark primary = %q", got.Theme.Dark["primary"])
	}
	for _, name := range TokenNames() {
		if name == "primary" {
			continue
		}
		if got.Theme.Light[name] != base.Theme.Light[name] {
			t.Fatalf("light %s changed: got %q, want %q", name, got.Theme.Light[name], base.Theme.Light[name])
		}
		if got.Theme.Dark[name] != base.Theme.Dark[name] {
			t.Fatalf("dark %s changed: got %q, want %q", name, got.Theme.Dark[name], base.Theme.Dark[name])
		}
	}
	if base.Theme.Light["primary"] != "oklch(0.205 0 0)" ||
		base.Theme.Dark["primary"] != "oklch(0.922 0 0)" {
		t.Fatal("ImportThemeCSS mutated the base preset")
	}

	got.Theme.Light["background"] = "red"
	if base.Theme.Light["background"] != "oklch(1 0 0)" {
		t.Fatal("ImportThemeCSS returned a theme map sharing base storage")
	}
}

func TestImportThemeCSSUpdatesRadiusFromEitherModeWhenConsistent(t *testing.T) {
	t.Parallel()

	got, err := ImportThemeCSS(Default(StyleNova), []byte(`
:root { --radius: 1.25REM; }
.dark { --radius: 1.25REM; }
`))
	if err != nil {
		t.Fatalf("ImportThemeCSS: %v", err)
	}
	if got.Radius != "1.25REM" {
		t.Fatalf("Radius = %q, want 1.25REM", got.Radius)
	}
}

func TestImportThemeCSSUpdatesFontsFromEitherModeWhenConsistent(t *testing.T) {
	t.Parallel()

	got, err := ImportThemeCSS(Default(StyleNova), []byte(`
:root { --font-sans: "Inter Variable", sans-serif; }
.dark { --font-heading: "Playfair Display Variable", serif; }
`))
	if err != nil {
		t.Fatalf("ImportThemeCSS: %v", err)
	}
	if got.FontSans != `"Inter Variable", sans-serif` {
		t.Fatalf("FontSans = %q", got.FontSans)
	}
	if got.FontHeading != `"Playfair Display Variable", serif` {
		t.Fatalf("FontHeading = %q", got.FontHeading)
	}
}

func TestImportThemeCSSRejectsConflictingFontDeclarations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		src  string
	}{
		{
			name: "font-sans",
			src: `
:root { --font-sans: Inter; }
.dark { --font-sans: Geist; }
`,
		},
		{
			name: "font-heading",
			src: `
:root { --font-heading: Inter; }
.dark { --font-heading: Geist; }
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := ImportThemeCSS(Default(StyleNova), []byte(tt.src))
			if err == nil || !strings.Contains(err.Error(), "conflicting") {
				t.Fatalf("ImportThemeCSS error = %v, want conflicting %s", err, tt.name)
			}
		})
	}
}

func TestImportThemeCSSIgnoresUnrelatedRulesAndReturnsCopy(t *testing.T) {
	t.Parallel()

	base := Default(StyleNova)
	got, err := ImportThemeCSS(base, []byte(`
body { color: red; --other: blue; }
@media (width > 10px) { .elsewhere { --other: green; } }
`))
	if err != nil {
		t.Fatalf("ImportThemeCSS: %v", err)
	}
	if !presetsEqual(got, base) {
		t.Fatal("unrelated CSS changed preset")
	}
	got.Theme.Dark["primary"] = "red"
	if base.Theme.Dark["primary"] == "red" {
		t.Fatal("ImportThemeCSS returned shared theme map storage")
	}
}

func TestImportThemeCSSRejectsInvalidOrAmbiguousImportsWithoutPartialPreset(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		src  string
		want string
	}{
		{
			name: "malformed CSS",
			src:  `:root { --primary: oklch(0.6 0.2 280);`,
			want: "CSS",
		},
		{
			name: "malformed unrelated CSS",
			src:  `body { color red; }`,
			want: "CSS",
		},
		{
			name: "duplicate light declaration",
			src:  `:root { --primary: red; --primary: blue; }`,
			want: "primary",
		},
		{
			name: "duplicate dark declaration",
			src:  `.dark { --primary: red; --primary: blue; }`,
			want: "primary",
		},
		{
			name: "recognized property outside mode",
			src:  `body { --primary: red; }`,
			want: "outside",
		},
		{
			name: "invalid color",
			src:  `:root { --primary: definitely-not-a-color; }`,
			want: "theme.light.primary",
		},
		{
			name: "invalid radius",
			src:  `:root { --radius: -1rem; }`,
			want: "radius",
		},
		{
			name: "conflicting radius modes",
			src:  `:root { --radius: 1rem; } .dark { --radius: 2rem; }`,
			want: "conflicting",
		},
		{
			name: "comma root selector",
			src:  `:root, body { --primary: red; }`,
			want: "selector",
		},
		{
			name: "comma dark selector",
			src:  `.dark, [data-theme=dark] { --primary: red; }`,
			want: "selector",
		},
		{
			name: "resembles dark",
			src:  `.dark-mode { --primary: red; }`,
			want: "selector",
		},
		{
			name: "compound dark selector",
			src:  `html.dark { --primary: red; }`,
			want: "selector",
		},
		{
			name: "nested ruleset changes ownership",
			src:  `:root { .nested { --primary: red; } }`,
			want: "nested",
		},
		{
			name: "at-rule changes ownership",
			src:  `@media (width > 10px) { :root { --primary: red; } }`,
			want: "nested",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ImportThemeCSS(Default(StyleNova), []byte(tt.src))
			if err == nil {
				t.Fatal("ImportThemeCSS returned nil error")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error = %q, want %q", err, tt.want)
			}
			if !presetsEqual(got, Preset{}) {
				t.Fatalf("ImportThemeCSS returned partial preset on error: %#v", got)
			}
		})
	}
}
