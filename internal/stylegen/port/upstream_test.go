package port

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestReadUpstream(t *testing.T) {
	style, sections, err := ReadUpstream(filepath.Join("testdata", "upstream-sample.css"))
	if err != nil {
		t.Fatalf("ReadUpstream: %v", err)
	}

	if style != "demo" {
		t.Fatalf("style = %q, want %q", style, "demo")
	}

	if len(sections) != 2 {
		t.Fatalf("len(sections) = %d, want 2 (sections: %+v)", len(sections), sections)
	}

	accordion := sections[0]
	if accordion.Name != "Accordion" {
		t.Fatalf("sections[0].Name = %q, want %q", accordion.Name, "Accordion")
	}
	radioGroup := sections[1]
	if radioGroup.Name != "Radio Group" {
		t.Fatalf("sections[1].Name = %q, want %q", radioGroup.Name, "Radio Group")
	}

	if len(accordion.Rules) != 2 {
		t.Fatalf("len(accordion.Rules) = %d, want 2 (rules: %+v)", len(accordion.Rules), accordion.Rules)
	}
	if got := accordion.Rules[0].Class; got != "cn-accordion" {
		t.Errorf("accordion.Rules[0].Class = %q, want %q", got, "cn-accordion")
	}
	if got := accordion.Rules[1].Class; got != "cn-accordion-trigger" {
		t.Errorf("accordion.Rules[1].Class = %q, want %q", got, "cn-accordion-trigger")
	}

	wantUtilities := []string{"overflow-hidden", "rounded-2xl", "border"}
	if got := accordion.Rules[0].Utilities; !reflect.DeepEqual(got, wantUtilities) {
		t.Errorf("cn-accordion Utilities = %#v, want %#v", got, wantUtilities)
	}

	if got := accordion.Rules[0].Line; got != 3 {
		t.Errorf("cn-accordion Line = %d, want 3", got)
	}

	if len(radioGroup.Rules) != 1 {
		t.Fatalf("len(radioGroup.Rules) = %d, want 1 (rules: %+v)", len(radioGroup.Rules), radioGroup.Rules)
	}
	radioLine := radioGroup.Rules[0].Line
	if radioLine < radioGroup.Start || radioLine > radioGroup.End {
		t.Errorf("Radio Group line range [%d, %d] does not cover its rule's line %d", radioGroup.Start, radioGroup.End, radioLine)
	}
}

func TestReadUpstreamRealMaia(t *testing.T) {
	const upstreamRoot = "/Users/jackieli/personal/shadcn-ui"
	if _, err := os.Stat(upstreamRoot); err != nil {
		t.Skipf("upstream checkout not present at %s: %v", upstreamRoot, err)
	}

	path := filepath.Join(upstreamRoot, "apps", "v4", "registry", "styles", "style-maia.css")
	style, sections, err := ReadUpstream(path)
	if err != nil {
		t.Fatalf("ReadUpstream(%s): %v", path, err)
	}

	if style != "maia" {
		t.Fatalf("style = %q, want %q", style, "maia")
	}

	if len(sections) < 55 {
		t.Errorf("len(sections) = %d, want >= 55", len(sections))
	}

	// The plan's step calls for ">600 rules total", but the dossier's own
	// §5.3 measurement at this exact pinned commit is "420 total .cn-*
	// rules" (`grep -c '^  \.cn-' style-maia.css`), which this test
	// reproduces exactly. >600 does not match reality; asserting >400
	// keeps the smoke test's intent (many rules parsed) while matching
	// verified ground truth.
	var totalRules int
	for _, sec := range sections {
		totalRules += len(sec.Rules)
	}
	if totalRules <= 400 {
		t.Errorf("total rules = %d, want > 400", totalRules)
	}
}
