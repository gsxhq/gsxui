package pages

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/gsxhq/gsx"
)

func TestDocHeadingPreservesSemanticsAndCallerPresentation(t *testing.T) {
	tests := []struct {
		name  string
		depth int
		tag   string
	}{
		{name: "depth two", depth: 2, tag: "h2"},
		{name: "depth three", depth: 3, tag: "h3"},
		{name: "unknown depth falls back to h2", depth: 99, tag: "h2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := docHeading(
				docTOCItem{ID: "explicit-id", Title: "Section title", Depth: tt.depth},
				gsx.Attrs{
					{Key: "ID", Value: "caller-id"},
					{Key: "DATA-DOC-HEADING", Value: "caller-marker"},
					{Key: "TaBiNdEx", Value: "0"},
					{Key: "class", Value: "text-sm font-medium"},
					{Key: "data-caller", Value: true},
				},
			)
			var output bytes.Buffer
			if err := node.Render(context.Background(), &output); err != nil {
				t.Fatalf("render docHeading: %v", err)
			}
			got := output.String()

			for _, want := range []string{
				"<" + tt.tag + " ",
				`id="explicit-id"`,
				`data-doc-heading`,
				`tabindex="-1"`,
				`class="scroll-mt-20 text-sm font-medium"`,
				`data-caller`,
				">Section title</" + tt.tag + ">",
			} {
				if !strings.Contains(got, want) {
					t.Errorf("docHeading(depth=%d) missing %q\nin: %s", tt.depth, want, got)
				}
			}
			lower := strings.ToLower(got)
			for _, forbidden := range []string{`id="caller-id"`, `tabindex="0"`, `data-doc-heading="caller-marker"`} {
				if strings.Contains(lower, forbidden) {
					t.Errorf("docHeading(depth=%d) rendered caller override %q\nin: %s", tt.depth, forbidden, got)
				}
			}
			for _, attribute := range []string{"id", "tabindex", "data-doc-heading"} {
				if count := strings.Count(lower, " "+attribute); count != 1 {
					t.Errorf("docHeading(depth=%d) rendered %d case-insensitive %s attributes, want exactly 1\nin: %s", tt.depth, count, attribute, got)
				}
			}
			otherTag := "h3"
			if tt.tag == "h3" {
				otherTag = "h2"
			}
			if strings.Contains(got, "<"+otherTag+" ") {
				t.Errorf("docHeading(depth=%d) unexpectedly rendered <%s>\nin: %s", tt.depth, otherTag, got)
			}
		})
	}
}
