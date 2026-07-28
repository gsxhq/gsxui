package stylepreview_test

import (
	"context"
	"io"
	"testing"

	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/site/stylepreview"
	"github.com/gsxhq/gsxui/site/stylepreview/maia"
	"github.com/gsxhq/gsxui/site/stylepreview/nova"
	"github.com/gsxhq/gsxui/ui"
)

func BenchmarkButtonCSS(b *testing.B) {
	benchmarkButton(b, ui.Button, nil)
}

func BenchmarkButtonInlineNova(b *testing.B) {
	benchmarkButton(b, nova.Button, nil)
}

func BenchmarkButtonInlineMaia(b *testing.B) {
	benchmarkButton(b, maia.Button, nil)
}

func BenchmarkButtonInlineNovaConflict(b *testing.B) {
	benchmarkButton(b, nova.Button, conflictingButtonAttrs())
}

func BenchmarkButtonInlineMaiaConflict(b *testing.B) {
	benchmarkButton(b, maia.Button, conflictingButtonAttrs())
}

func benchmarkButton(b *testing.B, button stylepreview.ButtonRenderer, attrs gsx.Attrs) {
	b.Helper()
	node := button("default", "default", "", false, gsx.Text("Button"), attrs)
	ctx := context.Background()
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		if err := node.Render(ctx, io.Discard); err != nil {
			b.Fatal(err)
		}
	}
}

func conflictingButtonAttrs() gsx.Attrs {
	return gsx.Attrs{{
		Key:   "class",
		Value: "rounded-none h-20 px-12 bg-warning text-warning-foreground rounded-r-none border-l-0",
	}}
}
