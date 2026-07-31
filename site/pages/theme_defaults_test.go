package pages_test

import (
	"context"
	"strings"
	"testing"

	"github.com/gsxhq/gsxui/site/pages"
	"golang.org/x/net/html"
)

func TestThemeEditorRendersPalettePickerDefaults(t *testing.T) {
	var output strings.Builder
	if err := pages.ThemeEditor("/theme/preview").Render(context.Background(), &output); err != nil {
		t.Fatal(err)
	}
	document, err := html.Parse(strings.NewReader(output.String()))
	if err != nil {
		t.Fatal(err)
	}
	defaults := map[string]struct {
		title string
		value string
	}{
		"baseColor": {title: "Neutral", value: "neutral"},
		"theme":     {title: "Neutral", value: "neutral"},
		"radius":    {title: "Medium", value: "medium"},
	}
	for picker, want := range defaults {
		node := findElementWithAttributeValue(document, "data-theme-picker", picker)
		if node == nil {
			t.Fatalf("missing %s picker", picker)
		}
		if !strings.Contains(htmlNodeText(node), want.title) {
			t.Fatalf("%s picker missing default selection", picker)
		}
		input := findElementWithAttributeValue(node, "value", want.value)
		if input == nil || input.Data != "input" {
			t.Fatalf("%s picker missing default radio value %q", picker, want.value)
		}
		if !hasHTMLAttribute(input, "checked") {
			t.Errorf("%s default radio %q is not checked", picker, want.value)
		}
	}
}

func findElementWithAttributeValue(node *html.Node, name, value string) *html.Node {
	if node.Type == html.ElementNode {
		if got, ok := htmlAttribute(node, name); ok && got == value {
			return node
		}
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if found := findElementWithAttributeValue(child, name, value); found != nil {
			return found
		}
	}
	return nil
}

func htmlNodeText(node *html.Node) string {
	var text strings.Builder
	var visit func(*html.Node)
	visit = func(current *html.Node) {
		if current.Type == html.TextNode {
			text.WriteString(current.Data)
		}
		for child := current.FirstChild; child != nil; child = child.NextSibling {
			visit(child)
		}
	}
	visit(node)
	return text.String()
}
