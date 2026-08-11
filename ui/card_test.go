package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestCardParts(t *testing.T) {
	cases := []struct {
		name string
		node gsx.Node
		want []string
	}{
		{"Card", ui.Card(gsx.Raw("x"), nil), []string{`data-gsxui-slot-card`}},
		{"CardHeader", ui.CardHeader(gsx.Raw("x"), nil), []string{`data-gsxui-slot-card-header`}},
		{"CardTitle", ui.CardTitle(gsx.Raw("x"), nil), []string{`data-gsxui-slot-card-title`}},
		{"CardDescription", ui.CardDescription(gsx.Raw("x"), nil), []string{`data-gsxui-slot-card-description`}},
		{"CardAction", ui.CardAction(gsx.Raw("x"), nil), []string{`data-gsxui-slot-card-action`}},
		{"CardContent", ui.CardContent(gsx.Raw("x"), nil), []string{`data-gsxui-slot-card-content`}},
		{"CardFooter", ui.CardFooter(gsx.Raw("x"), nil), []string{`data-gsxui-slot-card-footer`}},
	}
	for _, tc := range cases {
		got := render(t, tc.node)
		for _, want := range tc.want {
			if !strings.Contains(got, want) {
				t.Errorf("%s: missing %q\nin: %s", tc.name, want, got)
			}
		}
	}
}

func TestCardPinned(t *testing.T) {
	// Exact full-render pin. Card is migrated onto the slot axis: its class
	// attribute now carries the recipe classes resolved from
	// registry/canonical/card.gsx's card.Root() accessor call, rather than
	// being absent as it was pre-migration. The computed presentation is
	// unchanged — verified by the sweep (make sweep-compare) — only the
	// class attribute's literal contents changed, from nothing to the
	// resolved Nova recipe for the root slot.
	// group/card was added so Card's own has-[>[data-gsxui-slot-field]]
	// descendants (elsewhere in the recipe) and other slots can key off
	// Card's own ancestry — see the style-porter report's "missing
	// group/<name> marker" entry.
	got := render(t, ui.Card(gsx.Raw("Content"), nil))
	want := `<div class="group/card ring-foreground/10 bg-card text-card-foreground gap-(--card-spacing) overflow-hidden rounded-xl py-(--card-spacing) text-sm ring-1 [--card-spacing:--spacing(4)] has-[[data-gsxui-slot-card-footer]]:pb-0 has-[&gt;img:first-child]:pt-0 data-[size=sm]:[--card-spacing:--spacing(3)] data-[size=sm]:has-[[data-gsxui-slot-card-footer]]:pb-0 *:[img:first-child]:rounded-t-xl *:[img:last-child]:rounded-b-xl flex flex-col" data-gsxui-slot-card>Content</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestCardComposition(t *testing.T) {
	got := render(t, ui.Card(
		gsx.Fragment(
			ui.CardHeader(ui.CardTitle(gsx.Raw("Title"), nil), nil),
			ui.CardContent(gsx.Raw("Body"), nil),
		),
		gsx.Attrs{{Key: "class", Value: "py-8"}},
	))
	// The recipe now also emits a class attribute, so the caller's "py-8" is
	// merged into it rather than being the attribute's sole content — assert
	// containment of the token, not an exact `class="py-8"` attribute.
	for _, want := range []string{`data-gsxui-slot-card-title`, ">Title<", ">Body<", "py-8"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}
