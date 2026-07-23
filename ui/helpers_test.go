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
