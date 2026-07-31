package stylepreview_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/site/stylepreview/maia"
	"github.com/gsxhq/gsxui/site/stylepreview/nova"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// buttonRenderer is the copied Button contract shared by every generated
// style package, restated locally now that stylepreview.Matrix is gone: the
// tests below exercise the generated nova/maia Button sources directly.
type buttonRenderer func(
	variant, size, href string,
	disabled bool,
	children gsx.Node,
	attrs gsx.Attrs,
) gsx.Node

func TestButtonInlineStylesMergeCallerAttrs(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name   string
		button buttonRenderer
	}{
		{name: "nova", button: nova.Button},
		{name: "maia", button: maia.Button},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			rendered := renderButton(t, test.button(
				"default",
				"default",
				"",
				false,
				gsx.Text("Caller"),
				gsx.Attrs{
					{Key: "class", Value: "rounded-none h-20 px-12 bg-warning text-warning-foreground rounded-r-none border-l-0"},
					{Key: "data-caller-marker", Value: true},
				},
			))
			element := firstElement(t, rendered)
			classes := strings.Fields(htmlAttr(t, element, "class"))
			for _, want := range []string{
				"rounded-none",
				"h-20",
				"px-12",
				"bg-warning",
				"text-warning-foreground",
				"rounded-r-none",
				"border-l-0",
			} {
				if !contains(classes, want) {
					t.Errorf("class missing caller token %q: %q", want, strings.Join(classes, " "))
				}
			}
			for _, conflictingDefault := range []string{
				"rounded-lg",
				"rounded-4xl",
				"h-8",
				"h-9",
				"px-2.5",
				"px-4",
				"bg-primary",
				"text-primary-foreground",
			} {
				if contains(classes, conflictingDefault) {
					t.Errorf("class retained conflicting default %q: %q", conflictingDefault, strings.Join(classes, " "))
				}
			}
			if value, ok := attr(element, "data-caller-marker"); !ok || value != "" {
				t.Errorf("data-caller-marker = %q, %v; want bare presence", value, ok)
			}
			if value, ok := attr(element, "data-gsxui-slot-button"); !ok || value != "" {
				t.Errorf("data-gsxui-slot-button = %q, %v; want bare presence", value, ok)
			}
		})
	}
}

func TestButtonInlineStylesPreserveElementSemantics(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name   string
		button buttonRenderer
	}{
		{name: "nova", button: nova.Button},
		{name: "maia", button: maia.Button},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			link := firstElement(t, renderButton(t, test.button(
				"outline",
				"default",
				"/docs/getting-started",
				false,
				gsx.Text("Docs"),
				gsx.Attrs{{Key: "data-link-marker", Value: true}},
			)))
			if link.Data != "a" {
				t.Errorf("enabled href rendered <%s>, want <a>", link.Data)
			}
			if got := htmlAttr(t, link, "href"); got != "/docs/getting-started" {
				t.Errorf("href = %q, want /docs/getting-started", got)
			}
			if value, ok := attr(link, "data-link-marker"); !ok || value != "" {
				t.Errorf("data-link-marker = %q, %v; want bare presence", value, ok)
			}

			disabled := firstElement(t, renderButton(t, test.button(
				"default",
				"default",
				"/ignored-while-disabled",
				true,
				gsx.Text("Disabled"),
				nil,
			)))
			if disabled.Data != "button" {
				t.Errorf("disabled href rendered <%s>, want <button>", disabled.Data)
			}
			if value, ok := attr(disabled, "disabled"); !ok || value != "" {
				t.Errorf("disabled = %q, %v; want bare presence", value, ok)
			}
			if _, ok := attr(disabled, "href"); ok {
				t.Error("disabled button unexpectedly retained href")
			}
		})
	}
}

func renderButton(t *testing.T, node gsx.Node) string {
	t.Helper()
	var output bytes.Buffer
	if err := node.Render(context.Background(), &output); err != nil {
		t.Fatalf("Render: %v", err)
	}
	return output.String()
}

func firstElement(t *testing.T, rendered string) *html.Node {
	t.Helper()
	nodes, err := html.ParseFragment(strings.NewReader(rendered), &html.Node{
		Type:     html.ElementNode,
		Data:     "div",
		DataAtom: atom.Div,
	})
	if err != nil {
		t.Fatalf("html.ParseFragment: %v", err)
	}
	for _, node := range nodes {
		if node.Type == html.ElementNode {
			return node
		}
	}
	t.Fatalf("rendered output has no element: %q", rendered)
	return nil
}

func htmlAttr(t *testing.T, node *html.Node, name string) string {
	t.Helper()
	value, ok := attr(node, name)
	if !ok {
		t.Fatalf("<%s> missing %s", node.Data, name)
	}
	return value
}

func attr(node *html.Node, name string) (string, bool) {
	for _, attribute := range node.Attr {
		if attribute.Key == name {
			return attribute.Val, true
		}
	}
	return "", false
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
