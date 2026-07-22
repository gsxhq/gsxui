package checkbox_test

import (
	"context"
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui/checkbox"
)

func render(t *testing.T, n gsx.Node) string {
	t.Helper()
	var sb strings.Builder
	if err := n.Render(context.Background(), &sb); err != nil {
		t.Fatal(err)
	}
	return sb.String()
}

func TestCheckboxDefault(t *testing.T) {
	got := render(t, checkbox.Checkbox(nil))
	for _, want := range []string{
		`<input type="checkbox"`,
		`data-slot="checkbox"`,
		"peer size-4 shrink-0 appearance-none rounded-[4px]",
		"disabled:cursor-not-allowed disabled:opacity-50",
		"aria-invalid:border-destructive aria-invalid:ring-destructive/20",
		"checked:bg-primary checked:border-primary",
		"checked:bg-[url(&#39;data:image/svg+xml",
		"checked:bg-center checked:bg-no-repeat checked:bg-[length:12px_12px]",
		"/>",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestCheckboxNoStraySpaceInDataURI(t *testing.T) {
	// The data-URI's SVG attribute/path separators must survive the
	// tailwind-merge pass intact — they are authored as Tailwind's
	// underscore escape (`_`) for whitespace inside a bracketed arbitrary
	// value, never a literal space (a literal space is a class-token
	// boundary and gets torn apart by any whitespace-splitting class
	// tooling, corrupting the embedded SVG). See docs/jsx-parity.md. (The
	// quotes render as &#39; — gsx attribute-escapes every value.)
	got := render(t, checkbox.Checkbox(nil))
	want := "checked:bg-[url(&#39;data:image/svg+xml;charset=utf-8,%3Csvg_xmlns=%22http://www.w3.org/2000/svg%22_viewBox=%220_0_24_24%22_fill=%22none%22_stroke=%22white%22_stroke-width=%223%22_stroke-linecap=%22round%22_stroke-linejoin=%22round%22%3E%3Cpath_d=%22M20_6_9_17l-5-5%22/%3E%3C/svg%3E&#39;)]"
	if !strings.Contains(got, want) {
		t.Errorf("data-URI corrupted by class merge\nwant substring: %s\n in: %s", want, got)
	}
}

func TestCheckboxCallerClassMerges(t *testing.T) {
	got := render(t, checkbox.Checkbox(gsx.Attrs{{Key: "class", Value: "size-6"}}))
	if strings.Contains(got, "size-4") {
		t.Errorf("base size-4 should be dropped by caller size-6\nin: %s", got)
	}
	if !strings.Contains(got, "size-6") {
		t.Errorf("missing caller class size-6\nin: %s", got)
	}
}

func TestCheckboxAttrsFallThrough(t *testing.T) {
	got := render(t, checkbox.Checkbox(gsx.Attrs{{Key: "id", Value: "c1"}, {Key: "name", Value: "terms"}, {Key: "aria-label", Value: "Accept"}}))
	for _, want := range []string{`id="c1"`, `name="terms"`, `aria-label="Accept"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestCheckboxCheckedAttr(t *testing.T) {
	// checked is an HTML boolean (presence-only) attribute: a bool value on
	// it must render bare — no checked="true" — matching browser :checked
	// truth (gsx.IsBooleanAttr classifies "checked").
	got := render(t, checkbox.Checkbox(gsx.Attrs{{Key: "checked", Value: true}}))
	if !strings.Contains(got, " checked") || strings.Contains(got, `checked="`) {
		t.Errorf("checked attr should render bare, not stringified\nin: %s", got)
	}

	got = render(t, checkbox.Checkbox(gsx.Attrs{{Key: "checked", Value: false}}))
	// The class carries "checked:"-variant tokens regardless, so assert on
	// the bare-attribute shape specifically, not the substring "checked".
	if strings.Contains(got, `" checked`) || strings.Contains(got, `checked="false"`) {
		t.Errorf("checked=false should omit the attribute entirely\nin: %s", got)
	}
}

func TestCheckboxDisabledAttr(t *testing.T) {
	got := render(t, checkbox.Checkbox(gsx.Attrs{{Key: "disabled", Value: true}}))
	if !strings.Contains(got, " disabled") || strings.Contains(got, `disabled="`) {
		t.Errorf("disabled attr should render bare\nin: %s", got)
	}
}

func TestCheckboxPinned(t *testing.T) {
	// Exact full-render pin. Token-by-token verified against shadcn's
	// Checkbox (registry/new-york-v4/ui/checkbox.tsx) plus the ledgered
	// ADAPT: Radix Root/Indicator/CheckIcon replaced by a native
	// <input type="checkbox"> whose checked-state visuals move to a
	// checked:bg-[url(...)] data-URI background. See docs/jsx-parity.md.
	got := render(t, checkbox.Checkbox(nil))
	want := `<input type="checkbox" data-slot="checkbox" class="peer size-4 shrink-0 appearance-none rounded-[4px] border border-input shadow-xs transition-shadow outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 dark:bg-input/30 checked:bg-primary checked:border-primary checked:bg-[url(&#39;data:image/svg+xml;charset=utf-8,%3Csvg_xmlns=%22http://www.w3.org/2000/svg%22_viewBox=%220_0_24_24%22_fill=%22none%22_stroke=%22white%22_stroke-width=%223%22_stroke-linecap=%22round%22_stroke-linejoin=%22round%22%3E%3Cpath_d=%22M20_6_9_17l-5-5%22/%3E%3C/svg%3E&#39;)] checked:bg-center checked:bg-no-repeat checked:bg-[length:12px_12px]"/>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}
