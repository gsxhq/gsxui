package card_test

import (
	"context"
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui/card"
)

func render(t *testing.T, n gsx.Node) string {
	t.Helper()
	var sb strings.Builder
	if err := n.Render(context.Background(), &sb); err != nil {
		t.Fatal(err)
	}
	return sb.String()
}

func TestCardParts(t *testing.T) {
	cases := []struct {
		name string
		node gsx.Node
		want []string
	}{
		{"Card", card.Card(gsx.Raw("x"), nil), []string{`data-slot="card"`, "flex flex-col gap-6 rounded-xl border bg-card"}},
		{"CardHeader", card.CardHeader(gsx.Raw("x"), nil), []string{`data-slot="card-header"`, "@container/card-header"}},
		{"CardTitle", card.CardTitle(gsx.Raw("x"), nil), []string{`data-slot="card-title"`, "leading-none font-semibold"}},
		{"CardDescription", card.CardDescription(gsx.Raw("x"), nil), []string{`data-slot="card-description"`, "text-muted-foreground"}},
		{"CardAction", card.CardAction(gsx.Raw("x"), nil), []string{`data-slot="card-action"`, "col-start-2 row-span-2"}},
		{"CardContent", card.CardContent(gsx.Raw("x"), nil), []string{`data-slot="card-content"`, "px-6"}},
		{"CardFooter", card.CardFooter(gsx.Raw("x"), nil), []string{`data-slot="card-footer"`, "flex items-center px-6"}},
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

func TestCardComposition(t *testing.T) {
	got := render(t, card.Card(
		gsx.Fragment(
			card.CardHeader(card.CardTitle(gsx.Raw("Title"), nil), nil),
			card.CardContent(gsx.Raw("Body"), nil),
		),
		gsx.Attrs{{Key: "class", Value: "py-8"}},
	))
	for _, want := range []string{`data-slot="card-title"`, ">Title<", ">Body<", "py-8"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
	if strings.Contains(got, "py-6") {
		t.Errorf("caller py-8 must drop default py-6\nin: %s", got)
	}
}
