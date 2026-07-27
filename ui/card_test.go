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
		{"Card", ui.Card(gsx.Raw("x"), nil), []string{`data-gsxui-slot="card"`}},
		{"CardHeader", ui.CardHeader(gsx.Raw("x"), nil), []string{`data-gsxui-slot="card-header"`}},
		{"CardTitle", ui.CardTitle(gsx.Raw("x"), nil), []string{`data-gsxui-slot="card-title"`}},
		{"CardDescription", ui.CardDescription(gsx.Raw("x"), nil), []string{`data-gsxui-slot="card-description"`}},
		{"CardAction", ui.CardAction(gsx.Raw("x"), nil), []string{`data-gsxui-slot="card-action"`}},
		{"CardContent", ui.CardContent(gsx.Raw("x"), nil), []string{`data-gsxui-slot="card-content"`}},
		{"CardFooter", ui.CardFooter(gsx.Raw("x"), nil), []string{`data-gsxui-slot="card-footer"`}},
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
	// Exact full-render pin, verified token-by-token against shadcn's Card
	// (registry/new-york-v4/ui/card.tsx) and docs/jsx-parity.md — a straight
	// port, no divergences.
	got := render(t, ui.Card(gsx.Raw("Content"), nil))
	want := `<div data-gsxui-slot="card">Content</div>`
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
	for _, want := range []string{`data-gsxui-slot="card-title"`, ">Title<", ">Body<", `class="py-8"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}
