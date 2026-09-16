package rtl_test

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

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

func TestCSSRewritesApplyListsOnly(t *testing.T) {
	src := []byte(".a {\n  @apply pl-1.5 text-left;\n}\n/* pl-1.5 in a comment stays */\n.b { color: left; }\n")
	want := ".a {\n  @apply ps-1.5 text-start;\n}\n/* pl-1.5 in a comment stays */\n.b { color: left; }\n"
	if got := string(rtl.CSS(src)); got != want {
		t.Fatalf("got\n%s\nwant\n%s", got, want)
	}
	if got := string(rtl.CSS([]byte(want))); got != want {
		t.Fatal("CSS is not idempotent")
	}
}

// Every shipped component transforms, reparses, is a fixed point, and
// carries no physical direction class outside side-keyed arms and variants.
func TestGSXAcrossTheRegistry(t *testing.T) {
	files, err := filepath.Glob(filepath.Join("..", "..", "ui", "*.gsx"))
	if err != nil || len(files) == 0 {
		t.Fatalf("no ui/*.gsx found: %v", err)
	}
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
