package ui_test

import (
	"context"
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
)

// render renders n to a string, failing the test on error. Shared by every
// component's test file — this is the single copy.
func render(t *testing.T, n gsx.Node) string {
	t.Helper()
	var sb strings.Builder
	if err := n.Render(context.Background(), &sb); err != nil {
		t.Fatal(err)
	}
	return sb.String()
}

func openingTagContaining(t *testing.T, html, marker string) string {
	t.Helper()
	markerAt := strings.Index(html, marker)
	if markerAt < 0 {
		t.Fatalf("render missing %q\nin: %s", marker, html)
	}
	start := markerAt
	if html[markerAt] != '<' {
		start = strings.LastIndex(html[:markerAt], "<")
	}
	end := strings.Index(html[markerAt:], ">")
	if start < 0 || end < 0 {
		t.Fatalf("could not isolate opening tag containing %q\nin: %s", marker, html)
	}
	return html[start : markerAt+end+1]
}

func requirePresenceAttributesOnSameTag(t *testing.T, html, anchor string, attributes ...string) {
	t.Helper()
	tag := openingTagContaining(t, html, anchor)
	for _, attribute := range attributes {
		if !strings.Contains(tag, " "+attribute) {
			t.Errorf("opening tag containing %q missing presence attribute %q\ntag: %s", anchor, attribute, tag)
		}
		if strings.Contains(tag, attribute+`="`) {
			t.Errorf("presence attribute %q has a value\ntag: %s", attribute, tag)
		}
	}
}

func countPresenceAttribute(html, attribute string) int {
	count := 0
	for rest := html; ; {
		at := strings.Index(rest, " "+attribute)
		if at < 0 {
			return count
		}
		after := at + 1 + len(attribute)
		if after == len(rest) || rest[after] == ' ' || rest[after] == '>' {
			count++
		}
		rest = rest[after:]
	}
}
