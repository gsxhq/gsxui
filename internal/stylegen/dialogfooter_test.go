package stylegen

import (
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
)

// dialogFooterBase is the unconditional class upstream gives BOTH DialogFooter
// and AlertDialogFooter, style-invariantly:
// apps/v4/registry/new-york-v4/ui/{dialog,alert-dialog}.tsx and the style-pack
// bases apps/v4/registry/bases/radix/ui/{dialog,alert-dialog}.tsx all carry
// "flex flex-col-reverse gap-2 sm:flex-row sm:justify-end" (AlertDialogFooter
// adds only the group-data-[size=sm] grid pair this port drops — see
// docs/jsx-parity.md's "## alert-dialog"). Everything a style adds on top of
// this has to come from that style's OWN .cn-dialog-footer/.cn-alert-dialog-footer
// rule.
var dialogFooterBase = []string{"flex", "flex-col-reverse", "gap-2", "sm:flex-row", "sm:justify-end"}

// novaFooterChrome is the muted end-bar style-nova.css — and only
// style-nova.css — draws behind both footers: a full-bleed bar (negative
// margins cancelling the content's own padding), its own padding, a top border,
// a muted fill and the content's bottom radius. No other upstream style sheet
// has a .cn-alert-dialog-footer rule, and their .cn-dialog-footer rule (absent
// in lyra) is `gap-2` alone, so in every other style these tokens can only
// have arrived by copying nova's.
//
// They had: a 2026-09-14 parity audit found the whole set on
// .gsxui-recipe-alert-dialog-footer in all seven non-nova styles and on
// .gsxui-recipe-dialog-footer in lyra, turning a plain row of buttons into
// nova's bar everywhere. The failure mode is the one the accordion-root and
// ItemTitle fixes already hit twice — a rule checked against
// registry/styles/nova/ alone and generalised to all eight — so it gets a pin
// rather than a comment.
var novaFooterChrome = []string{"-mx-4", "-mb-4", "rounded-b-xl", "border-t", "bg-muted/50", "p-4"}

// TestDialogFootersCarryNovaChromeOnlyInNova reads the ported source of truth
// (registry/styles/<style>/{dialog,alert-dialog}.css) rather than any generated
// output: the generator faithfully propagates whatever the recipe says, so the
// recipe is where a leak has to be caught.
func TestDialogFootersCarryNovaChromeOnlyInNova(t *testing.T) {
	root := repoRoot(t)

	for _, sheet := range []struct{ file, class string }{
		{"dialog.css", "gsxui-recipe-dialog-footer"},
		{"alert-dialog.css", "gsxui-recipe-alert-dialog-footer"},
	} {
		for _, style := range allStyleDirs(t, root) {
			tokens := recipeTokens(t, filepath.Join(root, "registry", "styles", style, sheet.file), sheet.class)

			for _, want := range dialogFooterBase {
				if !slices.Contains(tokens, want) {
					t.Errorf("%s/%s: .%s must carry upstream's own base token %q, got %v",
						style, sheet.file, sheet.class, want, tokens)
				}
			}

			for _, chrome := range novaFooterChrome {
				has := slices.Contains(tokens, chrome)
				if style == DefaultStyle && !has {
					t.Errorf("%s/%s: .%s must keep nova's own %q (style-nova.css draws the muted end-bar)",
						style, sheet.file, sheet.class, chrome)
				}
				if style != DefaultStyle && has {
					t.Errorf("%s/%s: .%s carries %q, which only style-nova.css defines — %s has no upstream .cn-%s rule at all",
						style, sheet.file, sheet.class, chrome, style, strings.TrimPrefix(sheet.class, "gsxui-recipe-"))
				}
			}
		}
	}
}

// allStyleDirs lists every style pack under registry/styles, read from disk so
// a ninth style cannot be added without this gate seeing it.
func allStyleDirs(t *testing.T, root string) []string {
	t.Helper()
	entries, err := os.ReadDir(filepath.Join(root, "registry", "styles"))
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	for _, e := range entries {
		if e.IsDir() {
			out = append(out, e.Name())
		}
	}
	if len(out) == 0 {
		t.Fatal("no style packs found under registry/styles")
	}
	return out
}

// recipeTokens returns the @apply tokens of one recipe class in one sheet.
func recipeTokens(t *testing.T, path, class string) []string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	re := regexp.MustCompile(`\.` + regexp.QuoteMeta(class) + `\s*\{\s*@apply\s+([^;]+);`)
	m := re.FindSubmatch(content)
	if m == nil {
		t.Fatalf("%s: no .%s rule", path, class)
	}
	return strings.Fields(string(m[1]))
}
