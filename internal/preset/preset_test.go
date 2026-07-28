package preset

import (
	"slices"
	"strings"
	"testing"

	"github.com/mazznoer/csscolorparser"
)

var expectedTokenNames = []string{
	"background",
	"foreground",
	"card",
	"card-foreground",
	"popover",
	"popover-foreground",
	"primary",
	"primary-foreground",
	"secondary",
	"secondary-foreground",
	"muted",
	"muted-foreground",
	"accent",
	"accent-foreground",
	"destructive",
	"destructive-foreground",
	"border",
	"input",
	"ring",
	"sidebar",
	"sidebar-foreground",
	"sidebar-primary",
	"sidebar-primary-foreground",
	"sidebar-accent",
	"sidebar-accent-foreground",
	"sidebar-border",
	"sidebar-ring",
	"success",
	"info",
	"warning",
	"overlay",
	"contrast",
}

func TestCSSColorParserAcceptsShippedThemeColors(t *testing.T) {
	t.Parallel()

	colors := []string{
		"oklch(1 0 0)",
		"oklch(0% 0 0)",
		"oklch(0.145 0 0)",
		"oklch(0.205 0 0)",
		"oklch(0.985 0 0)",
		"oklch(0.922 0 0)",
		"oklch(0.97 0 0)",
		"oklch(0.269 0 0)",
		"oklch(0.371 0 0)",
		"oklch(0.556 0 0)",
		"oklch(0.708 0 0)",
		"oklch(0.577 0.245 27.325)",
		"oklch(0.704 0.191 22.216)",
		"oklch(0.97 0.01 17)",
		"oklch(0.58 0.22 27)",
		"oklch(1 0 0 / 10%)",
		"oklch(1 0 0 / 15%)",
		"oklch(69.6% 0.17 162.48)",
		"oklch(68.5% 0.169 237.323)",
		"oklch(76.9% 0.188 70.08)",
		"oklch(0% 0 0 / 10%)",
		"oklch(100% 0 0)",
		"oklch(0.488 0.243 264.376)",
		"oklch(0.439 0 0)",
	}

	for _, value := range colors {
		t.Run(value, func(t *testing.T) {
			t.Parallel()
			if _, err := csscolorparser.Parse(value); err != nil {
				t.Fatalf("Parse(%q): %v", value, err)
			}
		})
	}
}

func TestStylesAndTokenNamesReturnOrderedCopies(t *testing.T) {
	t.Parallel()

	styles := Styles()
	if want := []Style{StyleNova, StyleMaia}; !slices.Equal(styles, want) {
		t.Fatalf("Styles() = %v, want %v", styles, want)
	}
	styles[0] = "changed"
	if got := Styles()[0]; got != StyleNova {
		t.Fatalf("Styles returned shared storage: first style = %q", got)
	}

	names := TokenNames()
	if !slices.Equal(names, expectedTokenNames) {
		t.Fatalf("TokenNames() = %#v, want %#v", names, expectedTokenNames)
	}
	names[0] = "changed"
	if got := TokenNames()[0]; got != "background" {
		t.Fatalf("TokenNames returned shared storage: first name = %q", got)
	}
}

func TestDefaultReturnsIndependentValidPresets(t *testing.T) {
	t.Parallel()

	for _, style := range []Style{StyleNova, StyleMaia} {
		t.Run(string(style), func(t *testing.T) {
			t.Parallel()

			first := Default(style)
			if err := Validate(first); err != nil {
				t.Fatalf("Validate(Default(%q)): %v", style, err)
			}
			if first.Schema != SchemaURL {
				t.Fatalf("Schema = %q, want %q", first.Schema, SchemaURL)
			}
			if first.SchemaVersion != SchemaVersion {
				t.Fatalf("SchemaVersion = %d, want %d", first.SchemaVersion, SchemaVersion)
			}
			if first.Style != style {
				t.Fatalf("Style = %q, want %q", first.Style, style)
			}
			if first.Radius != "0.625rem" {
				t.Fatalf("Radius = %q, want 0.625rem", first.Radius)
			}

			first.Theme.Light["primary"] = "red"
			first.Theme.Dark["primary"] = "blue"
			again := Default(style)
			if got := again.Theme.Light["primary"]; got != "oklch(0% 0 0)" {
				t.Fatalf("light map was shared: primary = %q", got)
			}
			if got := again.Theme.Dark["primary"]; got != "oklch(0.922 0 0)" {
				t.Fatalf("dark map was shared: primary = %q", got)
			}
		})
	}
}

func TestNovaAndMaiaShareSemanticDefaults(t *testing.T) {
	t.Parallel()

	nova := Default(StyleNova)
	maia := Default(StyleMaia)
	nova.Style = maia.Style
	if !presetsEqual(nova, maia) {
		t.Fatal("Nova and Maia defaults differ beyond style")
	}
}

func TestValidateRejectsInvalidPresetFields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		edit func(*Preset)
		path string
	}{
		{
			name: "schema",
			edit: func(p *Preset) { p.Schema = "https://example.test/schema.json" },
			path: "$schema",
		},
		{
			name: "schema version",
			edit: func(p *Preset) { p.SchemaVersion = 2 },
			path: "schemaVersion",
		},
		{
			name: "unknown style",
			edit: func(p *Preset) { p.Style = "future" },
			path: "style",
		},
		{
			name: "missing light token",
			edit: func(p *Preset) { delete(p.Theme.Light, "primary") },
			path: "theme.light.primary",
		},
		{
			name: "missing dark token",
			edit: func(p *Preset) { delete(p.Theme.Dark, "sidebar-border") },
			path: "theme.dark.sidebar-border",
		},
		{
			name: "extra light token",
			edit: func(p *Preset) { p.Theme.Light["future-token"] = "red" },
			path: "theme.light.future-token",
		},
		{
			name: "extra dark token",
			edit: func(p *Preset) { p.Theme.Dark["future-token"] = "red" },
			path: "theme.dark.future-token",
		},
		{
			name: "invalid light color",
			edit: func(p *Preset) { p.Theme.Light["primary"] = "not-a-color" },
			path: "theme.light.primary",
		},
		{
			name: "invalid dark color",
			edit: func(p *Preset) { p.Theme.Dark["sidebar-border"] = "var(--border)" },
			path: "theme.dark.sidebar-border",
		},
		{
			name: "negative radius",
			edit: func(p *Preset) { p.Radius = "-1rem" },
			path: "radius",
		},
		{
			name: "multiple radius values",
			edit: func(p *Preset) { p.Radius = "1rem 2rem" },
			path: "radius",
		},
		{
			name: "radius percentage",
			edit: func(p *Preset) { p.Radius = "25%" },
			path: "radius",
		},
		{
			name: "radius keyword",
			edit: func(p *Preset) { p.Radius = "inherit" },
			path: "radius",
		},
		{
			name: "radius function",
			edit: func(p *Preset) { p.Radius = "calc(1rem + 1px)" },
			path: "radius",
		},
		{
			name: "radius variable",
			edit: func(p *Preset) { p.Radius = "var(--radius)" },
			path: "radius",
		},
		{
			name: "unsupported radius unit",
			edit: func(p *Preset) { p.Radius = "1fr" },
			path: "radius",
		},
		{
			name: "trailing radius garbage",
			edit: func(p *Preset) { p.Radius = "1rem;" },
			path: "radius",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			preset := Default(StyleNova)
			tt.edit(&preset)
			err := Validate(preset)
			if err == nil {
				t.Fatal("Validate returned nil")
			}
			if !strings.Contains(err.Error(), tt.path) {
				t.Fatalf("error = %q, want field path %q", err, tt.path)
			}
		})
	}
}

func TestValidateAcceptsSchemaV1RadiusLengths(t *testing.T) {
	t.Parallel()

	valid := []string{
		"0",
		"0.0",
		"-0",
		"1px",
		"2CM",
		"3Q",
		".5rem",
		"1e2vh",
		"0svmin",
		"4dvb",
		"5RCH",
	}
	for _, radius := range valid {
		t.Run(radius, func(t *testing.T) {
			t.Parallel()
			preset := Default(StyleNova)
			preset.Radius = radius
			if err := Validate(preset); err != nil {
				t.Fatalf("Validate(radius %q): %v", radius, err)
			}
		})
	}
}

func presetsEqual(a, b Preset) bool {
	if a.Schema != b.Schema ||
		a.SchemaVersion != b.SchemaVersion ||
		a.Style != b.Style ||
		a.Radius != b.Radius {
		return false
	}
	return mapsEqual(a.Theme.Light, b.Theme.Light) &&
		mapsEqual(a.Theme.Dark, b.Theme.Dark)
}

func mapsEqual(a, b ThemeValues) bool {
	if len(a) != len(b) {
		return false
	}
	for name, value := range a {
		if b[name] != value {
			return false
		}
	}
	return true
}
