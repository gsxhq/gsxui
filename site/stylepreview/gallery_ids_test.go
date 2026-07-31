package stylepreview_test

import (
	"context"
	"regexp"
	"strings"
	"testing"

	"github.com/gsxhq/gsxui/site/stylepreview/maia"
	"github.com/gsxhq/gsxui/site/stylepreview/nova"
)

// TestPreviewDocumentIdsAreUnique renders the gallery the way the preview
// document does — once per style — and asserts no id appears twice, and
// every label's for= resolves. This is what caught the Toaster singleton.
func TestPreviewDocumentIdsAreUnique(t *testing.T) {
	var b strings.Builder
	if err := nova.Gallery("nova").Render(context.Background(), &b); err != nil {
		t.Fatal(err)
	}
	if err := maia.Gallery("maia").Render(context.Background(), &b); err != nil {
		t.Fatal(err)
	}
	doc := b.String()
	ids := regexp.MustCompile(`\bid="([^"]+)"`).FindAllStringSubmatch(doc, -1)
	seen := map[string]int{}
	for _, m := range ids {
		seen[m[1]]++
	}
	for id, n := range seen {
		if n > 1 {
			t.Errorf("id %q appears %d times in the two-style preview document", id, n)
		}
	}
	for _, m := range regexp.MustCompile(`\bfor="([^"]+)"`).FindAllStringSubmatch(doc, -1) {
		if seen[m[1]] == 0 {
			t.Errorf("label for=%q resolves to no id in the document", m[1])
		}
	}
}
