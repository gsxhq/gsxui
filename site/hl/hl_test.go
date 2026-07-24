package hl_test

import (
	"bytes"
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gsxhq/gsxui/site/hl"
)

// sourceKeys walks the same two roots site/hl/gen does and returns the keys
// it would generate. Kept deliberately independent of the generator (which
// lives in its own module and cannot be imported here) so this test checks
// the committed output rather than re-deriving it from the same code.
func sourceKeys(t *testing.T) map[string]string {
	t.Helper()
	keys := map[string]string{}

	walk := func(dir, prefix string, nestedOnly bool) {
		err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() || strings.HasSuffix(path, ".x.go") {
				return nil
			}
			rel, err := filepath.Rel(dir, path)
			if err != nil {
				return err
			}
			if nestedOnly && !strings.Contains(rel, string(filepath.Separator)) {
				return nil
			}
			keys[prefix+strings.TrimSuffix(filepath.ToSlash(rel), ".txt")] = path
			return nil
		})
		if err != nil {
			t.Fatalf("walk %s: %v", dir, err)
		}
	}
	walk("../examples", "", true)
	walk("../snippets", "snippets/", false)
	return keys
}

// TestBlocksCoverSources is the staleness guard. blocks.gen.go is committed
// (the generator needs cgo and lives in its own module, so a normal build
// never runs it), which means an example added or renamed without running
// `make highlight` would otherwise ship a page rendering a "missing
// highlighted block" placeholder.
func TestBlocksCoverSources(t *testing.T) {
	sources := sourceKeys(t)
	for key := range sources {
		if !hl.Has(key) {
			t.Errorf("no highlighted block for %s — run `make highlight`", key)
		}
	}
	for _, key := range hl.Keys() {
		if _, ok := sources[key]; !ok {
			t.Errorf("highlighted block %s has no source file — run `make highlight`", key)
		}
	}
}

// TestBlocksMatchSourceText asserts the generated HTML still carries the
// current source text: strip the markup, unescape, and it must equal the file
// on disk. This is what catches an edited example whose block was never
// regenerated — the case TestBlocksCoverSources cannot see, because the key
// is still present and only the content went stale.
func TestBlocksMatchSourceText(t *testing.T) {
	for key, path := range sourceKeys(t) {
		want, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		got := unescape(stripTags(string(renderBlock(t, key))))
		// The highlighter drops carriage returns and trailing invalid UTF-8;
		// compare with CRs removed so a CRLF checkout does not fail here.
		if normalize(got) != normalize(string(want)) {
			t.Errorf("block %s is stale — run `make highlight`\n got: %q\nwant: %q",
				key, truncate(got), truncate(string(want)))
		}
	}
}

// TestUnknownKeyIsVisible pins the deliberate fallback: an unknown key must
// render something a reader would notice, not silently render nothing.
func TestUnknownKeyIsVisible(t *testing.T) {
	got := string(renderBlock(t, "definitely/not/a/block.gsx"))
	if !strings.Contains(got, "missing highlighted block") {
		t.Errorf("unknown key rendered %q, expected a visible placeholder", got)
	}
}

// TestPlaceholderEscapesKey guards the placeholder path against injecting
// markup from a key.
func TestPlaceholderEscapesKey(t *testing.T) {
	got := string(renderBlock(t, "<script>alert(1)</script>"))
	if strings.Contains(got, "<script>") {
		t.Errorf("placeholder did not escape its key: %q", got)
	}
}

// renderBlock renders hl.Node(key) the way a page would.
func renderBlock(t *testing.T, key string) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := hl.Node(key).Render(context.Background(), &buf); err != nil {
		t.Fatalf("render %s: %v", key, err)
	}
	return buf.Bytes()
}

func normalize(s string) string {
	return strings.TrimRight(strings.ReplaceAll(s, "\r", ""), "\n")
}

func truncate(s string) string {
	if len(s) > 120 {
		return s[:120] + "…"
	}
	return s
}

func stripTags(s string) string {
	var b strings.Builder
	depth := 0
	for _, r := range s {
		switch {
		case r == '<':
			depth++
		case r == '>' && depth > 0:
			depth--
		case depth == 0:
			b.WriteRune(r)
		}
	}
	return b.String()
}

func unescape(s string) string {
	return strings.NewReplacer(
		"&lt;", "<", "&gt;", ">", "&#34;", `"`, "&#39;", "'", "&amp;", "&",
	).Replace(s)
}
