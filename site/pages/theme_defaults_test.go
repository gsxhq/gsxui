package pages_test

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/gsxhq/gsxui/internal/preset"
	"github.com/gsxhq/gsxui/site/pages"
	"golang.org/x/net/html"
)

func TestThemeEditorRendersPresetGroupNamesAndDefaults(t *testing.T) {
	t.Parallel()

	var output strings.Builder
	if err := pages.ThemeEditor("/theme/preview/button").Render(context.Background(), &output); err != nil {
		t.Fatalf("ThemeEditor.Render: %v", err)
	}
	document, err := html.Parse(strings.NewReader(output.String()))
	if err != nil {
		t.Fatalf("html.Parse: %v", err)
	}

	got := renderedThemeGroups(document)
	want := presetThemeGroups()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("server-rendered theme groups drifted from preset\n got: %#v\nwant: %#v", got, want)
	}
	radius := findElementWithAttribute(document, "data-theme-field")
	if radius == nil {
		t.Fatal("theme editor has no radius input")
	}
	if field, _ := htmlAttribute(radius, "data-theme-field"); field != "radius" {
		t.Fatalf("first data-theme-field = %q, want radius", field)
	}
	value, _ := htmlAttribute(radius, "value")
	if value != preset.Default(preset.StyleNova).Radius {
		t.Errorf("radius input = %q, want %q", value, preset.Default(preset.StyleNova).Radius)
	}
}

type themeGroupSnapshot struct {
	Name string
	Vars []themeVarSnapshot
}

type themeVarSnapshot struct {
	Name  string
	Light string
	Dark  string
}

func presetThemeGroups() []themeGroupSnapshot {
	defaults := preset.Default(preset.StyleNova)
	groupIndexes := make(map[string]int)
	groups := make([]themeGroupSnapshot, 0, len(preset.GroupNames()))
	for _, name := range preset.GroupNames() {
		groupIndexes[name] = len(groups)
		groups = append(groups, themeGroupSnapshot{Name: name})
	}
	add := func(definition preset.TokenDefinition, light, dark string) {
		index, ok := groupIndexes[definition.Group]
		if !ok {
			panic("preset definition references unknown presentation group " + definition.Group)
		}
		groups[index].Vars = append(groups[index].Vars, themeVarSnapshot{
			Name:  "--" + definition.Name,
			Light: light,
			Dark:  dark,
		})
	}
	for _, definition := range preset.TokenDefinitions() {
		add(
			definition,
			defaults.Theme.Light[definition.Name],
			defaults.Theme.Dark[definition.Name],
		)
	}
	return groups
}

func renderedThemeGroups(document *html.Node) []themeGroupSnapshot {
	var groups []themeGroupSnapshot
	var visit func(*html.Node)
	visit = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "section" {
			if group, ok := renderedThemeGroup(node); ok {
				groups = append(groups, group)
				return
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			visit(child)
		}
	}
	visit(document)
	return groups
}

func renderedThemeGroup(section *html.Node) (themeGroupSnapshot, bool) {
	var group themeGroupSnapshot
	varIndexes := make(map[string]int)
	var visit func(*html.Node)
	visit = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "h2" && group.Name == "" {
			group.Name = strings.TrimSpace(htmlNodeText(node))
		}
		if node.Type == html.ElementNode && node.Data == "input" {
			name, hasName := htmlAttribute(node, "data-theme-var")
			mode, hasMode := htmlAttribute(node, "data-theme-mode")
			value, hasValue := htmlAttribute(node, "value")
			if hasName && hasMode && hasValue {
				index, ok := varIndexes[name]
				if !ok {
					index = len(group.Vars)
					varIndexes[name] = index
					group.Vars = append(group.Vars, themeVarSnapshot{Name: name})
				}
				switch mode {
				case "light":
					group.Vars[index].Light = value
				case "dark":
					group.Vars[index].Dark = value
				}
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			visit(child)
		}
	}
	visit(section)
	return group, len(group.Vars) != 0
}

func htmlNodeText(node *html.Node) string {
	var text strings.Builder
	var visit func(*html.Node)
	visit = func(node *html.Node) {
		if node.Type == html.TextNode {
			text.WriteString(node.Data)
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			visit(child)
		}
	}
	visit(node)
	return text.String()
}

func TestThemeDefaultsExposeEverySidebarTokenInTheEditor(t *testing.T) {
	want := map[string][2]string{
		"--sidebar":                    {"oklch(0.985 0 0)", "oklch(0.205 0 0)"},
		"--sidebar-foreground":         {"oklch(0% 0 0)", "oklch(0.985 0 0)"},
		"--sidebar-primary":            {"oklch(0.205 0 0)", "oklch(0.488 0.243 264.376)"},
		"--sidebar-primary-foreground": {"oklch(0.985 0 0)", "oklch(0.985 0 0)"},
		"--sidebar-accent":             {"oklch(0.97 0 0)", "oklch(0.269 0 0)"},
		"--sidebar-accent-foreground":  {"oklch(0.205 0 0)", "oklch(0.985 0 0)"},
		"--sidebar-border":             {"oklch(0.922 0 0)", "oklch(1 0 0 / 10%)"},
		"--sidebar-ring":               {"oklch(0.708 0 0)", "oklch(0.439 0 0)"},
	}

	got := make(map[string][2]string)
	for _, group := range pages.ThemeGroups() {
		for _, variable := range group.Vars {
			if _, ok := want[variable.Name]; ok {
				got[variable.Name] = [2]string{variable.Light, variable.Dark}
			}
		}
	}
	if len(got) != len(want) {
		t.Fatalf("editor sidebar tokens = %#v, want %#v", got, want)
	}
	for name, values := range want {
		if got[name] != values {
			t.Errorf("%s editor defaults = %#v, want %#v", name, got[name], values)
		}
	}
}
