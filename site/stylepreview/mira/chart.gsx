// Adapted from templui/shadcn-templ v2 (github.com/templui/templui @ 9ec720c03909, MIT) — see NOTICE.md.
package mira

import (
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/gsxhq/gsx"
)

// ChartConfig is the pendant of shadcn's ChartConfig: label and color per
// series key, keyed by ChartSeries.Key. Chart turns it into --color-<key>
// custom properties scoped to its own data-chart attribute (styleBlock).
type ChartConfig []ChartSeries

// ChartSeries is one entry of a ChartConfig. Key identifies the series both
// here and in the data rows a later task's chart roots read (rows are
// matched against Key, never Label). Theme, present, takes precedence over
// Color for both color schemes — see ChartConfig.styleBlock. Icon replaces
// the color swatch in the legend and tooltip parts a later task adds; it is
// unused by the container itself.
type ChartSeries struct {
	Key   string
	Label string
	Color string
	Theme *ChartSeriesTheme
	Icon  gsx.Node
}

// ChartSeriesTheme is a color pair for the light and dark color schemes,
// the pendant of the `theme` field in shadcn's own ChartConfig entries.
type ChartSeriesTheme struct {
	Light string
	Dark  string
}

// styleBlock is the ChartStyle pendant: one rule per color scheme, each
// scoped by [data-chart=<id>]. Ported byte-faithful (semantics, not syntax)
// from templui/shadcn-templ v2's chart.templ Config.styleBlock (see this
// file's header credit):
//
//   - a theme color wins over the plain Color for both schemes;
//   - a theme with no value for a given scheme falls back to Color, like
//     itemConfig.theme?.[theme] ?? itemConfig.color upstream;
//   - a series with neither Color nor Theme contributes no line, but does
//     not block the other series' lines;
//   - when NO series in the whole config carries a color, styleBlock
//     returns "" — no <style> tag at all — matching ChartStyle returning
//     null rather than an empty one upstream.
func (c ChartConfig) styleBlock(id string) string {
	hasColor := false
	for _, s := range c {
		if s.Color != "" || s.Theme != nil {
			hasColor = true
			break
		}
	}
	if !hasColor {
		return ""
	}
	var sb strings.Builder
	for _, scheme := range []struct{ Name, Prefix string }{{"light", ""}, {"dark", ".dark "}} {
		fmt.Fprintf(&sb, "%s[data-chart=%s] {\n", scheme.Prefix, id)
		for _, s := range c {
			color := s.Color
			if s.Theme != nil {
				themed := s.Theme.Light
				if scheme.Name == "dark" {
					themed = s.Theme.Dark
				}
				if themed != "" {
					color = themed
				}
			}
			if color != "" {
				fmt.Fprintf(&sb, "  --color-%s: %s;\n", s.Key, color)
			}
		}
		sb.WriteString("}\n")
	}
	return sb.String()
}

// chartInstanceCount backs every Chart's data-chart scoping key.
//
// gsxui has no pre-existing per-render unique-id helper: an exhaustive grep
// across registry/canonical, internal/, and the vendored gsx module found
// none. Every other component that needs a stable id takes it from the
// caller instead — an explicit param (e.g. Calendar's name), a fixed
// singleton default a caller can override through attrs (Toaster's
// id="gsxui-toaster"), or a caller-supplied prefix at the call site
// (site/stylepreview/gallery.gsx.src's idp-prefixed ids). None of those fit
// here: data-chart is not a caller-facing identity a consumer would ever
// want to set (unlike a DOM id), it is purely the selector key that pairs
// this <div> with the <style> block styleBlock emits right below it in the
// SAME render call, so it has to be computed inside Chart, before attrs is
// spread, with nothing to key off but the render call itself.
//
// This is a process-lifetime monotonic counter: deterministic (no math/rand
// or crypto/rand — the brief is explicit that pinned tests need
// determinism), safe under concurrent renders (atomic.Uint64), and unique
// for as long as the package's process runs, which is what actually matters
// for styleBlock's [data-chart=...] selector not leaking one chart's colors
// onto a sibling chart's elements. Every package this file's content gets
// copied into by stylegen (registry/generated/<style>, ui) gets its own
// counter, since each is a distinct Go package — same per-render-tree scope
// every other id-like mechanism in this codebase already has.
var chartInstanceCount atomic.Uint64

// nextChartID returns the next deterministic, process-unique data-chart
// scoping key: "chart-1", "chart-2", ....
func nextChartID() string {
	return "chart-" + strconv.FormatUint(chartInstanceCount.Add(1), 10)
}

// Chart is the shadcn/ui ChartContainer pendant: an aspect-video box that
// exposes its ChartConfig as --color-<key> CSS custom properties (see
// styleBlock) for the chart roots and parts a later task adds as children.
// templui's own class carries a wall of [&_.recharts-*] arbitrary-variant
// selectors (axis tick fill, grid stroke, cursor fill) that reach into the
// client-rendered SVG upstream's Recharts leaves behind; gsxui's client
// renderer (Task 5) stamps gsxui's own classes on those nodes instead, so
// that job moves to the renderer-emitted attributes reading the same
// --color-<key>/--chart-N variables, and the selector wall is dropped here
// entirely rather than ported dead.
component Chart(config ChartConfig, children gsx.Node, attrs gsx.Attrs) {
	{{ id := nextChartID() }}
	<div
		data-chart={id}
		class={ "flex aspect-video justify-center text-xs" }
		{ attrs... }
		data-gsxui-slot-chart
	>
		{ if block := config.styleBlock(id); block != "" {
			{ gsx.Raw("<style>" + block + "</style>") }
		} }
		{ children }
	</div>
}
