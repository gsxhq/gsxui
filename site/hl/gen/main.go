// Command gen renders every docs code block to highlighted HTML and writes
// site/hl/blocks.gen.go.
//
// This lives in its own module (see go.mod beside this file) on purpose.
// Highlighting needs tree-sitter, which is C; the site binary is built
// CGO_ENABLED=0 for a distroless/static image. Keeping the generator out of
// gsxui's module means those cgo dependencies — and the sibling `replace`
// directives that resolve them from a local checkout — never touch the
// server build or its Docker context. The generated output is committed, so
// a normal build never runs this at all.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"html"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/gsxhq/gsxhl"
)

// snippetExt reports the language-bearing extension of a doc snippet.
//
// Snippets carry a doubled extension — first-page.gsx.txt, theme-tokens.css.txt
// — because `gsx generate` walks the whole tree and compiles every .gsx file
// it finds, with no ignore mechanism. These are illustrative fragments (one
// is a bare <ui.Button> with no package clause), so a real .gsx name makes
// the build fail. The trailing .txt keeps every tool out while the inner
// extension still says what the block is.
func snippetExt(path string) string {
	if base, ok := strings.CutSuffix(path, ".txt"); ok {
		if inner := filepath.Ext(base); inner != "" {
			return inner
		}
	}
	return filepath.Ext(path)
}

// langForPath maps a file to its highlighter language. Go maps to LangGSX
// deliberately: the gsx grammar parses Go natively, so one configuration
// covers both.
//
// Anything unlisted (.sh, .txt — shell commands and captured CLI output)
// renders as escaped text with no spans. No shell grammar is wired in;
// adding a dependency to color three-word `go install` lines earns nothing.
func langForPath(path string) (string, bool) {
	switch snippetExt(path) {
	case ".gsx", ".go":
		return gsxhl.LangGSX, true
	case ".css":
		return gsxhl.LangCSS, true
	case ".js":
		return gsxhl.LangJS, true
	default:
		return "", false
	}
}

type block struct {
	key  string
	html string
	lang string
}

func main() {
	root := flag.String("root", "../../..", "gsxui repo root, relative to this directory")
	out := flag.String("out", "../blocks.gen.go", "generated file to write")
	flag.Parse()

	h, err := gsxhl.New()
	if err != nil {
		log.Fatalf("gen: %v", err)
	}

	var blocks []block
	// Component examples, keyed by the SourcePath the registry already uses
	// ("button/basic.gsx"), so component.gsx can look them up with no extra
	// bookkeeping.
	examplesDir := filepath.Join(*root, "site", "examples")
	if err := collect(h, examplesDir, "", &blocks); err != nil {
		log.Fatalf("gen: examples: %v", err)
	}
	// Doc snippets, namespaced under snippets/ so they cannot collide with a
	// component directory name.
	snippetsDir := filepath.Join(*root, "site", "snippets")
	if _, err := os.Stat(snippetsDir); err == nil {
		if err := collect(h, snippetsDir, "snippets/", &blocks); err != nil {
			log.Fatalf("gen: snippets: %v", err)
		}
	}

	if len(blocks) == 0 {
		log.Fatal("gen: no source files found — wrong -root?")
	}
	sort.Slice(blocks, func(i, j int) bool { return blocks[i].key < blocks[j].key })

	src, err := render(blocks)
	if err != nil {
		log.Fatalf("gen: %v", err)
	}
	if err := os.WriteFile(*out, src, 0o644); err != nil {
		log.Fatalf("gen: %v", err)
	}

	var plain int
	for _, b := range blocks {
		if b.lang == "" {
			plain++
		}
	}
	fmt.Printf("gsxui/hl: %d blocks (%d highlighted, %d plain text) → %s\n",
		len(blocks), len(blocks)-plain, plain, filepath.Clean(*out))
}

// collect walks dir and highlights every file it finds, keying each by its
// path relative to dir, with prefix prepended.
func collect(h *gsxhl.Highlighter, dir, prefix string, out *[]block) error {
	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		// .x.go files are gsx's generated Go output sitting next to each
		// example; the page shows the .gsx source, never the compiled form.
		if strings.HasSuffix(path, ".x.go") {
			return nil
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		// Only nested example files (button/basic.gsx) are examples; the
		// package's own .go files sit at the top level and are not shown.
		if prefix == "" && !strings.Contains(rel, string(filepath.Separator)) {
			return nil
		}
		source, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		// Keys drop the .txt guard extension, so pages reference the block by
		// what it actually is: "snippets/first-page.gsx".
		key := prefix + strings.TrimSuffix(filepath.ToSlash(rel), ".txt")
		lang, ok := langForPath(path)
		if !ok {
			*out = append(*out, block{key: key, html: html.EscapeString(string(source))})
			return nil
		}
		hl, err := h.Highlight(lang, source)
		if err != nil {
			return fmt.Errorf("%s: %w", key, err)
		}
		*out = append(*out, block{key: key, html: hl, lang: lang})
		return nil
	})
}

func render(blocks []block) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("// Code generated by site/hl/gen. DO NOT EDIT.\n\n")
	buf.WriteString("package hl\n\n")
	buf.WriteString("// blocks maps a source-relative key to its pre-rendered highlighted\n")
	buf.WriteString("// HTML (inner markup only — the page supplies <pre><code>).\n")
	buf.WriteString("var blocks = map[string]string{\n")
	for _, b := range blocks {
		fmt.Fprintf(&buf, "\t%q: %q,\n", b.key, b.html)
	}
	buf.WriteString("}\n\n")

	langs := make([]string, 0, 4)
	for _, b := range blocks {
		if b.lang != "" && !slices.Contains(langs, b.lang) {
			langs = append(langs, b.lang)
		}
	}
	sort.Strings(langs)
	fmt.Fprintf(&buf, "// generatedLanguages records which highlighters produced the blocks\n")
	fmt.Fprintf(&buf, "// above, for the staleness test.\n")
	fmt.Fprintf(&buf, "var generatedLanguages = %#v\n", langs)

	return format.Source(buf.Bytes())
}
