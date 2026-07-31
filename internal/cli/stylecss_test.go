package cli

import (
	"io/fs"
	"path"
	"regexp"
	"strings"
	"testing"
	"testing/fstest"

	gsxui "github.com/gsxhq/gsxui"
)

// sheets builds an in-memory stylesheet tree for the composer's failure cases.
func sheets(files map[string]string) fs.FS {
	mapped := fstest.MapFS{}
	for name, content := range files {
		mapped[name] = &fstest.MapFile{Data: []byte(content)}
	}
	return mapped
}

// aggregatorImports returns the relative @import targets of the default style
// aggregator, in source order.
func aggregatorImports(t *testing.T) []string {
	t.Helper()
	content, err := fs.ReadFile(gsxui.Files, defaultStyleCSS)
	if err != nil {
		t.Fatal(err)
	}
	var targets []string
	for _, line := range strings.Split(string(content), "\n") {
		match := relativeImport.FindStringSubmatch(line)
		if match == nil {
			continue
		}
		targets = append(targets, match[1])
	}
	return targets
}

// TestEveryDefaultStylePartIsImported is the guard that makes the split safe to
// extend. A contributor who adds assets/css/styles/default/<component>.css and
// forgets the @import line would ship every consumer a style.css missing that
// component — the authored tree would look complete, the site would look right
// (Tailwind compiles the same aggregator), and only the vendored output would be
// wrong. Nothing else in the repository notices that.
func TestEveryDefaultStylePartIsImported(t *testing.T) {
	t.Parallel()

	dir := path.Dir(defaultStyleCSS) + "/default"
	entries, err := fs.ReadDir(gsxui.Files, dir)
	if err != nil {
		t.Fatal(err)
	}
	onDisk := map[string]bool{}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".css") {
			t.Errorf("%s/%s is not a .css file; the aggregator only imports stylesheets", dir, entry.Name())
			continue
		}
		onDisk[entry.Name()] = true
	}
	if len(onDisk) == 0 {
		t.Fatalf("%s holds no stylesheets", dir)
	}

	imported := map[string]int{}
	for _, target := range aggregatorImports(t) {
		if !strings.HasPrefix(target, "./default/") {
			t.Errorf("aggregator imports %q; every part lives under ./default/", target)
			continue
		}
		imported[strings.TrimPrefix(target, "./default/")]++
	}

	for name := range onDisk {
		switch imported[name] {
		case 1:
		case 0:
			t.Errorf("%s/%s exists but the aggregator does not import it: every consumer's vendored style.css would be missing it", dir, name)
		default:
			t.Errorf("the aggregator imports %s %d times; a duplicated import duplicates its rules in the vendored sheet", name, imported[name])
		}
	}
	for name := range imported {
		if !onDisk[name] {
			t.Errorf("the aggregator imports ./default/%s, which does not exist", name)
		}
	}
}

// TestComposedStyleCSSIsTheVendoredFlatSheet asserts the composition byte for
// byte against an independent substitution, so a bug in composeStyleCSS's line
// handling (a dropped terminator, a swallowed line, a part inlined at the wrong
// place) fails here rather than in a consumer's build.
func TestComposedStyleCSSIsTheVendoredFlatSheet(t *testing.T) {
	t.Parallel()

	aggregator, err := fs.ReadFile(gsxui.Files, defaultStyleCSS)
	if err != nil {
		t.Fatal(err)
	}
	var want strings.Builder
	for _, line := range strings.SplitAfter(string(aggregator), "\n") {
		if line == "" {
			continue
		}
		match := relativeImport.FindStringSubmatch(strings.TrimSuffix(line, "\n"))
		if match == nil {
			want.WriteString(line)
			continue
		}
		part, err := fs.ReadFile(gsxui.Files, path.Join(path.Dir(defaultStyleCSS), match[1]))
		if err != nil {
			t.Fatal(err)
		}
		want.Write(part)
	}

	got, err := composeStyleCSS(gsxui.Files, defaultStyleCSS)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want.String() {
		t.Errorf("composed style.css is not the aggregator with its parts spliced in place (got %d bytes, want %d)", len(got), want.Len())
	}
}

// TestComposedStyleCSSPreservesImportOrder pins the cascade. The parts' bodies
// must appear in the flat sheet in the aggregator's import order: two
// equal-specificity rules in different parts are ordered by nothing else, so a
// composer that reordered them would change rendering while every rule survived.
func TestComposedStyleCSSPreservesImportOrder(t *testing.T) {
	t.Parallel()

	composed, err := composeStyleCSS(gsxui.Files, defaultStyleCSS)
	if err != nil {
		t.Fatal(err)
	}
	flat := string(composed)
	previous := -1
	previousName := ""
	for _, target := range aggregatorImports(t) {
		part, err := fs.ReadFile(gsxui.Files, path.Join(path.Dir(defaultStyleCSS), target))
		if err != nil {
			t.Fatal(err)
		}
		at := strings.Index(flat, string(part))
		if at < 0 {
			t.Fatalf("%s does not appear verbatim in the composed sheet", target)
		}
		if at < previous {
			t.Errorf("%s is composed before %s but imported after it", target, previousName)
		}
		previous, previousName = at, target
	}
}

// unresolvedImport catches any @import at all surviving composition, including
// forms relativeImport deliberately does not match — an @import with a media
// query or layer() would reach a consumer unresolved and silently drop rules.
var unresolvedImport = regexp.MustCompile(`(?m)^@import\b`)

// TestComposedStyleCSSHasNoUnresolvedImports is what the vendored file's
// self-sufficiency rests on: the consumer gets one file with no siblings, so an
// @import that survives points at a path that does not exist in their project.
func TestComposedStyleCSSHasNoUnresolvedImports(t *testing.T) {
	t.Parallel()

	composed, err := composeStyleCSS(gsxui.Files, defaultStyleCSS)
	if err != nil {
		t.Fatal(err)
	}
	if found := unresolvedImport.FindAllString(string(composed), -1); len(found) != 0 {
		t.Errorf("composed style.css still holds %d @import statement(s); the vendored sheet has no siblings to resolve them against", len(found))
	}
}

// TestComposeStyleCSSRejectsAMissingPart pins the failure mode: an import that
// cannot be read must stop the vendor, never pass through as a literal line that
// a consumer's bundler then fails on (or, worse, ignores).
func TestComposeStyleCSSRejectsAMissingPart(t *testing.T) {
	t.Parallel()

	broken := sheets(map[string]string{
		"a.css": "@import \"./missing.css\";\n",
	})
	if _, err := composeStyleCSS(broken, "a.css"); err == nil {
		t.Fatal("composeStyleCSS() = nil error for an unreadable import, want a failure")
	}
}

// TestComposeStyleCSSRejectsACycle keeps a mistaken self-import from recursing
// until the process dies.
func TestComposeStyleCSSRejectsACycle(t *testing.T) {
	t.Parallel()

	looped := sheets(map[string]string{
		"a.css": "@import \"./b.css\";\n",
		"b.css": "@import \"./a.css\";\n",
	})
	_, err := composeStyleCSS(looped, "a.css")
	if err == nil {
		t.Fatal("composeStyleCSS() = nil error for an import cycle, want a failure")
	}
	if !strings.Contains(err.Error(), "cycle") {
		t.Errorf("error must name the cycle, got %q", err)
	}
}

// TestComposeStyleCSSPassesBareSpecifiersThrough keeps `@import "tailwindcss";`
// intact: it is the consumer's bundler's to resolve, not ours.
func TestComposeStyleCSSPassesBareSpecifiersThrough(t *testing.T) {
	t.Parallel()

	sheet := sheets(map[string]string{
		"a.css": "@import \"tailwindcss\";\n@import \"./b.css\";\n",
		"b.css": ".b {\n  color: red;\n}\n",
	})
	got, err := composeStyleCSS(sheet, "a.css")
	if err != nil {
		t.Fatal(err)
	}
	const want = "@import \"tailwindcss\";\n.b {\n  color: red;\n}\n"
	if string(got) != want {
		t.Errorf("composeStyleCSS() = %q, want %q", got, want)
	}
}
