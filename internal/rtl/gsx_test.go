package rtl_test

import (
	"bytes"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	gsxast "github.com/gsxhq/gsx/ast"
	gsxparser "github.com/gsxhq/gsx/parser"

	"github.com/gsxhq/gsxui/internal/rtl"
)

// The fixtures are named .gsx.txt, not .gsx: internal/generatedcheck requires
// a committed .x.go beside every .gsx in the tree, and `gsx generate` skips
// testdata, so a .gsx fixture would fail `make verify-generated` (and have the
// gsx LSP drop a stray .x.go beside it). The .want file is still gsx fmt's own
// output — GSX ends in gen.Format, and the test below reformats it again.
func TestGSXRewritesOnlyClassLiterals(t *testing.T) {
	src, err := os.ReadFile("testdata/composed.gsx.txt")
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile("testdata/composed.want.gsx.txt")
	if err != nil {
		t.Fatal(err)
	}
	got, err := rtl.GSX("composed.gsx", src)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("transform mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
	again, err := rtl.GSX("composed.gsx", got)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(again, got) {
		t.Fatal("GSX is not idempotent")
	}
}

func TestGSXLeavesAFileWithoutDirectionClassesByteIdentical(t *testing.T) {
	src := []byte("package ui\n\ncomponent Plain(attrs gsx.Attrs) {\n\t<p class={ \"flex gap-2\", attrs.Class() }>x</p>\n}\n")
	got, err := rtl.GSX("plain.gsx", src)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, src) {
		t.Fatalf("unchanged file was rewritten:\n%s", got)
	}
}

// TestCSSRewritesApplyListsOnly pins CSS's actual reach: the @apply regex
// has no comment awareness and no side-keyed notion at all, so an @apply
// inside a /* … */ comment is rewritten exactly like a live one. That is
// safe today only because TestShippedSheetsHaveNoSideKeyedApply confirms no
// shipped sheet has an @apply in a comment, or a physical-side selector on
// the same line as an @apply, to be miscategorised.
func TestCSSRewritesApplyListsOnly(t *testing.T) {
	src := []byte(".a {\n  @apply pl-1.5 text-left;\n}\n/* pl-1.5 in a comment stays */\n.b { color: left; }\n")
	want := ".a {\n  @apply ps-1.5 text-start;\n}\n/* pl-1.5 in a comment stays */\n.b { color: left; }\n"
	if got := string(rtl.CSS(src)); got != want {
		t.Fatalf("got\n%s\nwant\n%s", got, want)
	}
	if got := string(rtl.CSS([]byte(want))); got != want {
		t.Fatal("CSS is not idempotent")
	}

	// A comment's own @apply is rewritten too — pinned, not desired: CSS
	// scans bytes for the @apply pattern with no notion of a /* … */ span.
	commentSrc := []byte("/* @apply pl-2; */\n.c { @apply pr-2; }\n")
	commentWant := "/* @apply ps-2; */\n.c { @apply pe-2; }\n"
	if got := string(rtl.CSS(commentSrc)); got != commentWant {
		t.Fatalf("got\n%s\nwant\n%s", got, commentWant)
	}
}

// TestShippedSheetsHaveNoSideKeyedApply is the other half of the pin above:
// it keeps CSS's blindness to comments and side-keyed selectors harmless by
// asserting no shipped stylesheet gives it anything to miscategorise — no
// line pairs a physical-side selector with an @apply.
func TestShippedSheetsHaveNoSideKeyedApply(t *testing.T) {
	root := filepath.Join("..", "..", "assets", "css")
	sideKeyed := regexp.MustCompile(`data-side=|data-\[side=`)
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() || filepath.Ext(path) != ".css" {
			return err
		}
		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for i, line := range strings.Split(string(src), "\n") {
			if sideKeyed.MatchString(line) && strings.Contains(line, "@apply") {
				t.Errorf("%s:%d: side-keyed selector shares a line with @apply: %q", path, i+1, line)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

// registryGSXFiles is the file set both registry sweeps below share: the
// canonical components in ui/, plus their vendored copies across every
// generated registry style — everything `gsxui migrate rtl` and `gsxui add
// --style` can hand to GSX.
func registryGSXFiles(t *testing.T) []string {
	t.Helper()
	uiFiles, err := filepath.Glob(filepath.Join("..", "..", "ui", "*.gsx"))
	if err != nil || len(uiFiles) == 0 {
		t.Fatalf("no ui/*.gsx found: %v", err)
	}
	generatedFiles, err := filepath.Glob(filepath.Join("..", "..", "registry", "generated", "*", "*.gsx"))
	if err != nil || len(generatedFiles) == 0 {
		t.Fatalf("no registry/generated/*/*.gsx found: %v", err)
	}
	files := append(uiFiles, generatedFiles...)
	sort.Strings(files)
	return files
}

// Every shipped component — ui/'s canonical source and every generated
// registry style's vendored copy — transforms, reparses, is a fixed point,
// and carries no physical direction class outside side-keyed arms and
// variants.
func TestGSXAcrossTheRegistry(t *testing.T) {
	files := registryGSXFiles(t)
	physical := regexp.MustCompile(`^-?(m|p|border|rounded|scroll-m|scroll-p)?[lr]-|^-?(left|right)-|^text-(left|right)$|^border-[lr]$|^float-(left|right)$|^clear-(left|right)$|^origin-(top-|bottom-)?(left|right)$`)
	for _, file := range files {
		src, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		out, err := rtl.GSX(filepath.Base(file), src)
		if err != nil {
			t.Fatalf("%s: %v", file, err)
		}
		again, err := rtl.GSX(filepath.Base(file), out)
		if err != nil || !bytes.Equal(again, out) {
			t.Fatalf("%s: not idempotent (%v)", file, err)
		}
		for _, list := range rtl.ClassListsForTest(filepath.Base(file), out) {
			if list.SideKeyed {
				continue
			}
			for _, class := range strings.Fields(list.Value) {
				variant, value := rtl.SplitForTest(class)
				if rtl.IsSideKeyedForTest(variant) {
					continue
				}
				if physical.MatchString(value) {
					t.Errorf("%s: physical class %q survived in %q", filepath.Base(file), class, list.Value)
				}
			}
		}
	}
}

// TestRegistryHasNoConstructsTheWalkerSkips is a registry-only assertion,
// not a production walkErr in classLists: a consumer's own vendored source
// can legally hold class={ f"…" } or a tagless switch inside a class list
// (classLists simply does not rewrite what it cannot see inside), and
// `migrate rtl` must not fail on a file that isn't ours to police. What
// this test guards is narrower — that the shipped registry never hands the
// walker one of those blind spots, so TestGSXAcrossTheRegistry's "no
// physical class survived" sweep is not quietly skipping part of a file.
// It also pins the side-keyed ruling's textual scope: the only class-list
// switch tags or if conditions naming side or direction live in sheet.gsx
// (Sheet, Sidebar) or drawer.gsx (Drawer) — any other file using either
// identifier there would be a class list the sweep above treats as
// side-keyed (and so exempts from the physical-class check) without the
// physical-side ruling actually applying.
func TestRegistryHasNoConstructsTheWalkerSkips(t *testing.T) {
	files := registryGSXFiles(t)
	for _, file := range files {
		base := filepath.Base(file)
		src, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		fset := token.NewFileSet()
		f, err := gsxparser.ParseFile(fset, base, src, 0)
		if err != nil {
			t.Fatalf("%s: %v", file, err)
		}

		var stack []gsxast.Node
		inClassAttr := func() bool {
			for i := len(stack) - 1; i >= 0; i-- {
				if ca, ok := stack[i].(*gsxast.ComposedAttr); ok {
					return ca.Name == "class"
				}
			}
			return false
		}

		gsxast.Inspect(f, func(n gsxast.Node) bool {
			if n == nil {
				stack = stack[:len(stack)-1]
				return true
			}
			stack = append(stack, n)
			switch a := n.(type) {
			case *gsxast.EmbeddedAttr:
				if a.Name == "class" {
					t.Errorf("%s: class=css`…`/f`…` (EmbeddedAttr) — the walker only rewrites StaticAttr/ComposedAttr class literals", base)
				}
			case *gsxast.ExprAttr:
				if a.Name == "class" {
					t.Errorf("%s: class={expr} (ExprAttr) — the walker only rewrites literals, not a bare expression", base)
				}
			case *gsxast.ComposedPart:
				if !inClassAttr() {
					break
				}
				if len(a.LiteralSegments) > 0 {
					t.Errorf("%s: class part has LiteralSegments (a css`…` literal part inside a class list)", base)
				}
			case *gsxast.ValueArm:
				if !inClassAttr() {
					break
				}
				if len(a.Segments) > 0 {
					t.Errorf("%s: class value arm has Segments (an f`…` literal arm inside a class list)", base)
				}
			case *gsxast.ValueSwitch:
				if !inClassAttr() {
					break
				}
				if a.Tag == "" {
					t.Errorf("%s: tagless switch inside a class list (classLists' side-keyed check reads only the tag)", base)
					break
				}
				if rtl.MentionsPlacementForTest(a.Tag) && base != "sheet.gsx" && base != "drawer.gsx" {
					t.Errorf("%s: switch %q inside a class list mentions side/direction outside sheet.gsx/drawer.gsx", base, a.Tag)
				}
			case *gsxast.ValueIf:
				if !inClassAttr() {
					break
				}
				if rtl.MentionsPlacementForTest(a.Cond) && base != "sheet.gsx" && base != "drawer.gsx" {
					t.Errorf("%s: if %q inside a class list mentions side/direction outside sheet.gsx/drawer.gsx", base, a.Cond)
				}
			}
			return true
		})
	}
}
