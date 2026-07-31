package canonical_test

import (
	"regexp"
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/registry/canonical"
)

var disabledAttr = regexp.MustCompile(`disabled(>|\s)`)

func TestButtonRendersButtonByDefault(t *testing.T) {
	got := render(t, canonical.Button("", "", "", false, gsx.Raw("Save"), nil))
	for _, want := range []string{
		"<button", "data-gsxui-slot-button", `type="button"`,
		`data-variant="default"`, `data-size="default"`, ">Save</button>",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
	if disabledAttr.MatchString(got) {
		t.Errorf("unexpected disabled attr\nin: %s", got)
	}
}

func TestButtonWithHrefRendersAnchor(t *testing.T) {
	got := render(t, canonical.Button("", "", "/docs", false, gsx.Raw("Docs"), nil))
	if !strings.Contains(got, "<a ") || !strings.Contains(got, `href="/docs"`) {
		t.Errorf("want an anchor with href\nin: %s", got)
	}
}

func TestDisabledAlwaysRendersButtonEvenWithHref(t *testing.T) {
	got := render(t, canonical.Button("", "", "/docs", true, gsx.Raw("Docs"), nil))
	if !strings.Contains(got, "<button") {
		t.Errorf("disabled with href must render a button\nin: %s", got)
	}
	if !disabledAttr.MatchString(got) {
		t.Errorf("missing disabled attr\nin: %s", got)
	}
}

func TestCallerAttrsFallThrough(t *testing.T) {
	got := render(t, canonical.Button("", "", "", false, gsx.Raw("x"),
		gsx.Attrs{{Key: "data-testid", Value: "save"}}))
	if !strings.Contains(got, `data-testid="save"`) {
		t.Errorf("missing caller attr\nin: %s", got)
	}
}

func TestUnrecognizedVariantResolvesToDefault(t *testing.T) {
	got := render(t, canonical.Button("destructve", "", "", false, gsx.Raw("x"), nil))
	if !strings.Contains(got, "gsxui-recipe-button-variant-default") {
		t.Errorf("unrecognized variant must resolve to the default\nin: %s", got)
	}
}
