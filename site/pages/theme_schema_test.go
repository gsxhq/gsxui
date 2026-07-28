package pages_test

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/gsxhq/gsxui/internal/preset"
	"github.com/gsxhq/gsxui/site/pages"
	"golang.org/x/net/html"
)

func TestThemeEditorSchemaMatchesPresetAuthority(t *testing.T) {
	t.Parallel()

	var output strings.Builder
	if err := pages.ThemeEditor("/theme/preview/button").Render(context.Background(), &output); err != nil {
		t.Fatalf("ThemeEditor.Render: %v", err)
	}
	document, err := html.Parse(strings.NewReader(output.String()))
	if err != nil {
		t.Fatalf("html.Parse: %v", err)
	}
	schemaNode := findElementWithAttribute(document, "data-theme-schema")
	if schemaNode == nil {
		t.Fatal("ThemeEditor has no data-theme-schema script")
	}

	var schema struct {
		Schema            string                     `json:"schema"`
		SchemaVersion     int                        `json:"schemaVersion"`
		TransportPrefix   string                     `json:"transportPrefix"`
		TokenNames        []string                   `json:"tokenNames"`
		RadiusUnits       []string                   `json:"radiusUnits"`
		Styles            []string                   `json:"styles"`
		Defaults          map[string]json.RawMessage `json:"defaults"`
		CanonicalDefaults map[string]string          `json:"canonicalDefaults"`
	}
	if err := json.Unmarshal([]byte(htmlNodeText(schemaNode)), &schema); err != nil {
		t.Fatalf("parse data-theme-schema JSON: %v\n%s", err, htmlNodeText(schemaNode))
	}
	if schema.Schema != preset.SchemaURL {
		t.Errorf("schema URL = %q, want %q", schema.Schema, preset.SchemaURL)
	}
	if schema.SchemaVersion != preset.SchemaVersion {
		t.Errorf("schema version = %d, want %d", schema.SchemaVersion, preset.SchemaVersion)
	}
	if schema.TransportPrefix != "gsxui:v1:" {
		t.Errorf("transport prefix = %q, want gsxui:v1:", schema.TransportPrefix)
	}
	if !reflect.DeepEqual(schema.TokenNames, preset.TokenNames()) {
		t.Errorf("token names = %#v, want %#v", schema.TokenNames, preset.TokenNames())
	}
	if !reflect.DeepEqual(schema.RadiusUnits, preset.RadiusUnits()) {
		t.Errorf("radius units = %#v, want %#v", schema.RadiusUnits, preset.RadiusUnits())
	}
	wantStyles := []string{"nova", "maia"}
	if !reflect.DeepEqual(schema.Styles, wantStyles) {
		t.Errorf("styles = %#v, want %#v", schema.Styles, wantStyles)
	}
	for _, style := range preset.Styles() {
		raw, ok := schema.Defaults[string(style)]
		if !ok {
			t.Errorf("schema defaults missing %s", style)
			continue
		}
		got, err := preset.ParseJSON(raw)
		if err != nil {
			t.Errorf("schema default %s: %v", style, err)
			continue
		}
		if !reflect.DeepEqual(got, preset.Default(style)) {
			t.Errorf("schema default %s = %#v, want %#v", style, got, preset.Default(style))
		}
		canonical, err := preset.CanonicalJSON(got)
		if err != nil {
			t.Fatal(err)
		}
		if schema.CanonicalDefaults[string(style)] != string(canonical) {
			t.Errorf("schema default %s is not canonical JSON", style)
		}
	}
}

func findElementWithAttribute(node *html.Node, name string) *html.Node {
	if node.Type == html.ElementNode && hasHTMLAttribute(node, name) {
		return node
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if found := findElementWithAttribute(child, name); found != nil {
			return found
		}
	}
	return nil
}
