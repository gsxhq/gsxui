// Adapted from templui/shadcn-templ v2 (github.com/templui/templui @ 9ec720c03909, MIT) — see NOTICE.md.
package rhea

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sort"
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
	{{
		id := nextChartID()
		// The config travels to the chart roots and their parts through ctx,
		// the pendant of shadcn's ChartContext — templui's Container does the
		// same (context.WithValue inside its own templ body). Chart roots
		// (BarChart/LineChart/AreaChart, appended below) read it back via
		// chartConfigFromCtx to resolve each series' label and to render the
		// legend server-side.
		ctx = context.WithValue(ctx, chartConfigCtxKey, config)
	}}
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

// ---------------------------------------------------------------------
// Cartesian model builder (bar, line, area) — adapted from templui's
// chart.templ lines 138-1319 (Datum, chartState, ctx keys, AreaChart /
// BarChart / LineChart, CartesianGrid / XAxis / YAxis / Tooltip / Legend /
// Defs / LinearGradient / Area / Bar / Line, buildModel, legendItems,
// legendContent, ModelScript). Every part below renders nothing itself
// (except the roots' wrapping div and the legend, both explicit here):
// children register into a chartState the enclosing root seeds through
// ctx, mirroring templui's own ctx = context.WithValue(...) inside each
// templ body — see NOTICE.md.
//
// The serialized model (ChartModel, below) is a byte-faithful, field-for-
// field port of templui's own flat Model — same json tags, same
// omitempty-or-not, same Recharts defaults — because Task 5 adapts
// templui's chart.js renderer under a "keep geometry/animation code
// unmodified otherwise" constraint, and that renderer reads these exact
// flat fields (marginTop, xAxisHeight, tickMargin, minTickGap, xTickLine,
// ...) in ~2100 lines of ported code no one wants to rewrite. The only two
// gsxui-specific departures are the script tag's attribute name
// (data-gsxui-chart-model, not data-tui-chart-model) and one new field,
// ChartTooltipModel.TooltipClass — see ChartModel's own doc comment.

// ChartDatum is one row of chart data, the pendant of the plain objects in
// shadcn's chartData arrays and templui's own Datum. Rows are read back by
// key (chartNum/chartStr) while building the model; they are not
// serialized into it themselves — Labels/TooltipLabels/per-series Values
// carry what the client needs, exactly like templui's own Model.
type ChartDatum map[string]any

// ChartMargin replaces the chart's plot margin (Recharts' default of 5 on
// every side) entirely when registered. Left unset, the model omits its
// margin field and the client renderer applies the Recharts default
// itself — see ChartModel's own doc comment on that convention.
type ChartMargin struct {
	Top, Right, Bottom, Left float64
}

// ChartCurveType selects a Line's or Area's interpolation, the pendant of
// templui's CurveType (itself Recharts' curve prop).
type ChartCurveType string

// ChartCurveNatural is Recharts' default curve.
const ChartCurveNatural ChartCurveType = "natural"

// ChartBool returns a pointer to v, for options whose Recharts default is
// true (e.g. CartesianGrid's Horizontal/Vertical) — pass ChartBool(false)
// to turn one off, a nil option leaves the Recharts default in place.
func ChartBool(v bool) *bool { return &v }

// ChartFloat returns a pointer to v, for numeric options that carry a
// meaningful zero.
func ChartFloat(v float64) *float64 { return &v }

// chartBoolOr resolves an optional option against its Recharts default.
func chartBoolOr(v *bool, def bool) bool {
	if v == nil {
		return def
	}
	return *v
}

// chartCtxKey namespaces this file's two context values so they cannot
// collide with a key any other package might put on the same ctx.
type chartCtxKey int

const (
	chartConfigCtxKey chartCtxKey = iota
	chartStateCtxKey
)

// chartConfigFromCtx reads back the ChartConfig Chart seeded into ctx.
func chartConfigFromCtx(ctx context.Context) ChartConfig {
	c, _ := ctx.Value(chartConfigCtxKey).(ChartConfig)
	return c
}

// chartStateFromCtx reads back the chartState the enclosing chart root
// seeded into ctx, or nil when called outside one (a part rendered
// without an enclosing BarChart/LineChart/AreaChart registers nothing,
// the same no-op templui's own stateFrom(ctx) == nil path takes).
func chartStateFromCtx(ctx context.Context) *chartState {
	s, _ := ctx.Value(chartStateCtxKey).(*chartState)
	return s
}

// chartState is the collector a chart root shares with its parts through
// ctx: each part appends to it and renders nothing, the root normalizes
// the result into a ChartModel once every child has rendered (gsx renders
// children before the root's own post-children statements run, mirroring
// templui's own children-then-chartOutput order).
type chartState struct {
	kind        string
	data        []ChartDatum
	margin      *ChartMargin
	stackOffset string
	grid        *ChartCartesianGridOptions
	x           *chartAxisReg
	y           *chartYAxisReg
	tooltip     *chartTooltipReg
	legend      *ChartLegendOptions
	defs        []chartGradientReg
	bars        []chartBarReg
	lines       []chartLineReg
	areas       []chartAreaReg
}

type chartAxisReg struct {
	key  string
	opts ChartXAxisOptions
}

type chartYAxisReg struct {
	key  string
	opts ChartYAxisOptions
}

type chartTooltipReg struct {
	opts ChartTooltipOptions
}

type chartGradientReg struct {
	id, x1, y1, x2, y2, stops string
}

type chartBarReg struct {
	key  string
	opts ChartBarOptions
}

type chartLineReg struct {
	key  string
	opts ChartLineOptions
}

type chartAreaReg struct {
	key  string
	opts ChartAreaOptions
}

// ChartCartesianGridOptions is the pendant of Recharts' CartesianGrid.
// Both directions default to Recharts' true, so they are pointers: pass
// ChartBool(false) to turn one off.
type ChartCartesianGridOptions struct {
	Horizontal *bool
	Vertical   *bool
}

// ChartXAxisOptions is the pendant of Recharts' XAxis, minus the render-prop
// fields (TickFormatter) a JSON model boundary cannot carry.
type ChartXAxisOptions struct {
	Hide       bool
	TickLine   bool
	AxisLine   bool
	TickMargin float64
	MinTickGap float64 // defaults to Recharts' 5
}

// ChartYAxisOptions is the pendant of Recharts' YAxis.
type ChartYAxisOptions struct {
	Hide       bool
	TickLine   bool
	AxisLine   bool
	TickMargin float64
	TickCount  int     // defaults to Recharts' 5
	Width      float64 // defaults to Recharts' 60
}

// ChartTooltipOptions is the pendant of ChartTooltip plus
// ChartTooltipContent's non-function fields — Formatter/LabelFormatter are
// deferred (see this task's report: a Go closure cannot cross the JSON
// boundary to the client renderer, and the model already carries the raw
// row data those would have formatted from).
type ChartTooltipOptions struct {
	// Cursor defaults to Recharts' true, so it is a pointer: pass
	// ChartBool(false) to turn it off.
	Cursor        *bool
	Indicator     string // "dot" (default), "line", "dashed"
	LabelKey      string
	HideLabel     bool
	HideIndicator bool
	NameKey       string
	Class         string // extra content class, e.g. "w-[150px]"
	LabelClass    string
	// Color overrides the indicator color for every row.
	Color string
	// DefaultIndex shows the tooltip on mount at that category.
	DefaultIndex *int
}

// ChartLegendOptions is the pendant of ChartLegend plus
// ChartLegendContent's props.
type ChartLegendOptions struct {
	// NameKey picks the config entry for the label, like the nameKey of
	// ChartLegendContent.
	NameKey string
	Class   string
	// VerticalAlign "top" places the legend above the plot; everything
	// else is "bottom".
	VerticalAlign string
	// HideIcon drops the config icon and falls back to the color swatch.
	HideIcon bool
}

// ChartLinearGradientOptions is the pendant of a linearGradient element in
// the chart defs, declared by a caller that fills an area with a
// gradient.
type ChartLinearGradientOptions struct {
	ID             string
	X1, Y1, X2, Y2 string
}

// ChartAreaOptions is the pendant of one Recharts Area.
type ChartAreaOptions struct {
	Type ChartCurveType
	// Fill and Stroke are used verbatim, e.g. "url(#fillDesktop)" or
	// "var(--color-desktop)". Empty falls back to the series color.
	Fill        string
	Stroke      string
	StackID     string
	FillOpacity float64 // 0 uses Recharts' default 0.6
}

// ChartBarOptions is the pendant of one Recharts Bar.
type ChartBarOptions struct {
	Fill    string
	StackID string
	// Radius is Recharts' radius union: a float64 for all corners or a
	// []float64 of four corners, e.g. []float64{0, 0, 4, 4}.
	Radius      any
	StrokeWidth float64
}

// ChartDotOptions is the pendant of Recharts' dot prop's object form.
type ChartDotOptions struct {
	R           float64
	FillOpacity float64
	Fill        string
	Size        float64 // icon box, unused without an Icon renderer
}

// ChartLineOptions is the pendant of one Recharts Line.
type ChartLineOptions struct {
	Type        ChartCurveType
	Stroke      string
	StrokeWidth float64
	// Dot draws the per point dots; nil is the demos' dot={false}.
	Dot *ChartDotOptions
	// ActiveDotR sizes the hover dot, like Recharts' activeDot prop.
	ActiveDotR float64
}

// ChartBarChartOptions overrides BarChart's Recharts defaults. Layout
// (Recharts' vertical bar layout) and AccessibilityLayer are deferred —
// neither is in this task's explicit port list (Margin, CurveType,
// StackOffset, Bool/Float helpers); a later task can extend this struct.
type ChartBarChartOptions struct {
	Margin *ChartMargin
}

// ChartLineChartOptions overrides LineChart's Recharts defaults.
type ChartLineChartOptions struct {
	Margin *ChartMargin
}

// ChartAreaChartOptions overrides AreaChart's Recharts defaults.
type ChartAreaChartOptions struct {
	Margin *ChartMargin
	// StackOffset "expand" normalizes each stack to 100%, the pendant of
	// Recharts' stackOffset prop.
	StackOffset string
}

// label resolves key's config entry, the ChartConfig pendant of templui's
// Config.Label: it returns that entry's own Label field verbatim (which
// may be empty) when key matches, falling back to key itself only when NO
// entry names it at all.
func (c ChartConfig) label(key string) string {
	for _, s := range c {
		if s.Key == key {
			return s.Label
		}
	}
	return key
}

// chartSeriesColor resolves the fill/stroke for a data key: the
// container's generated --color-<key> variable (see ChartConfig.styleBlock).
func chartSeriesColor(key string) string {
	return "var(--color-" + key + ")"
}

// chartNum normalizes a data value to float64, the pendant of templui's
// num — anything not already a numeric kind contributes 0.
func chartNum(v any) float64 {
	switch n := v.(type) {
	case float64:
		return n
	case float32:
		return float64(n)
	case int:
		return float64(n)
	case int64:
		return float64(n)
	}
	return 0
}

// chartRadiusCorners normalizes Recharts' radius union into corner radii:
// a number applies to all four corners, a slice is used as given.
func chartRadiusCorners(v any) []float64 {
	switch r := v.(type) {
	case nil:
		return nil
	case float64:
		return []float64{r}
	case int:
		return []float64{float64(r)}
	case []float64:
		return r
	}
	return nil
}

// chartSeriesMin is the lowest value across every series' own Values, so a
// negative domain keeps its zero baseline — the pendant of templui's
// seriesMin, byte-faithful (it also reads s.Values, not the raw rows).
func chartSeriesMin(series []ChartModelSeries) float64 {
	min := 0.0
	for _, s := range series {
		for _, v := range s.Values {
			if v < min {
				min = v
			}
		}
	}
	return min
}

// chartStr renders a data value as text, the pendant of templui's str —
// fmt.Sprint on whatever the row carries under that key.
func chartStr(v any) string {
	return fmt.Sprint(v)
}

// chartRenderHTML renders a gsx.Node (e.g. a config icon) to raw HTML for
// the legend markup.
func chartRenderHTML(ctx context.Context, n gsx.Node) string {
	if n == nil {
		return ""
	}
	var sb strings.Builder
	if err := n.Render(ctx, &sb); err != nil {
		return ""
	}
	return sb.String()
}

// defaultChartXAxisHeight and defaultChartLegendHeight are Recharts'
// reserved-space defaults, ported from templui's own
// defaultXAxisHeight/defaultLegendHeight (documented in the reference
// chart.js header alongside minTickGap 5, tickCount 5).
const (
	defaultChartXAxisHeight  = 30.0
	defaultChartLegendHeight = 28.0
)

// ChartModel is the payload a chart root serializes for the client
// renderer: a faithful, field-for-field port of templui's own flat Model
// (chart.templ:1645), narrowed to the fields the cartesian kinds this task
// builds (bar, line, area) actually populate — SliceColors/Pies/Polar/Radial
// and the always-unset top-level Radius are pie/radar/radial-only in
// upstream too (grep confirms nothing ever assigns m.Radius there), so
// Task 4 adds them alongside the kinds that need them, the same way this
// task's own buildChartModel switch only covers bar/line/area.
//
// The two gsxui-specific fields are TooltipModel.TooltipClass (Task 5
// populates it from the compiled recipe class; empty here since nothing
// renders a tooltip element yet — see ChartTooltipModel's own doc comment)
// and the script tag's attribute name (chartModelScript uses
// data-gsxui-chart-model, not templui's data-tui-chart-model). Every other
// field, its json tag, its omitempty-or-not, and its Recharts default is
// ported byte-faithful so Task 5's adapted chart.js — which reads these
// exact flat fields (marginTop, xAxisHeight, tickMargin, minTickGap,
// xTickLine, ... per chart.templ:1645) unmodified — needs no rewriting.
type ChartModel struct {
	Kind         string  `json:"kind"` // "bar" | "area" | "line"
	MarginTop    float64 `json:"marginTop"`
	MarginRight  float64 `json:"marginRight"`
	MarginBottom float64 `json:"marginBottom"`
	MarginLeft   float64 `json:"marginLeft"`
	XAxisHeight  float64 `json:"xAxisHeight,omitempty"`
	TickMargin   float64 `json:"tickMargin,omitempty"`
	MinTickGap   float64 `json:"minTickGap,omitempty"`
	YAxisWidth   float64 `json:"yAxisWidth,omitempty"`
	YAxisMargin  float64 `json:"yAxisMargin,omitempty"` // tickMargin of the y axis
	TickCount    int     `json:"tickCount,omitempty"`   // y ticks, Recharts default 5
	XTickLine    bool    `json:"xTickLine,omitempty"`
	XAxisLine    bool    `json:"xAxisLine,omitempty"`
	YTickLine    bool    `json:"yTickLine,omitempty"`
	YAxisLine    bool    `json:"yAxisLine,omitempty"`
	LegendHeight float64 `json:"legendHeight,omitempty"`
	LegendVAlign string  `json:"legendVAlign,omitempty"` // "top" raises the legend above the plot
	CategoryGap  float64 `json:"categoryGap,omitempty"`

	Grid           bool `json:"grid,omitempty"`
	GridHorizontal bool `json:"gridHorizontal,omitempty"`
	GridVertical   bool `json:"gridVertical,omitempty"`

	// Layout "vertical" swaps the axes and draws the bars horizontally —
	// out of this task's scope (see ChartBarChartOptions' own doc
	// comment), so this always marshals as absent for now.
	Layout      string  `json:"layout,omitempty"`
	XAxisHide   bool    `json:"xAxisHide,omitempty"`
	YAxisHide   bool    `json:"yAxisHide,omitempty"`
	DomainMin   float64 `json:"domainMin,omitempty"` // negative values extend the domain
	Stacked     bool    `json:"stacked,omitempty"`
	StackOffset string  `json:"stackOffset,omitempty"` // "expand" normalizes each stack to 1

	Defs []ChartLinearGradientModel `json:"defs,omitempty"`

	Cursor bool `json:"cursor"`
	// HasTooltip marks a declared ChartTooltip: without one the client
	// renders no tooltip at all, so it skips the hover wiring.
	HasTooltip    bool               `json:"hasTooltip,omitempty"`
	Labels        []string           `json:"labels"`
	TooltipLabels []string           `json:"tooltipLabels,omitempty"`
	Series        []ChartModelSeries `json:"series"`
	Tooltip       ChartTooltipModel  `json:"tooltip"`

	// AccessibilityLayer switches on the keyboard layer — also out of this
	// task's scope, always absent for now.
	AccessibilityLayer bool `json:"accessibilityLayer,omitempty"`
}

// ChartLinearGradientModel is a registered ChartLinearGradient, its Stops
// carrying the rendered raw <stop> children verbatim.
type ChartLinearGradientModel struct {
	ID    string `json:"id"`
	X1    string `json:"x1,omitempty"`
	Y1    string `json:"y1,omitempty"`
	X2    string `json:"x2,omitempty"`
	Y2    string `json:"y2,omitempty"`
	Stops string `json:"stops,omitempty"`
}

// ChartTooltipModel mirrors templui's TooltipModel, resolved from
// ChartTooltipOptions, plus this task's one addition: TooltipClass.
type ChartTooltipModel struct {
	Indicator      string `json:"indicator,omitempty"` // "dot" (default) | "line" | "dashed"
	Label          string `json:"label,omitempty"`     // labelKey resolved through the config
	HideLabel      bool   `json:"hideLabel,omitempty"`
	HideIndicator  bool   `json:"hideIndicator,omitempty"`
	Width          string `json:"width,omitempty"` // extra content class, e.g. "w-[150px]"
	LabelClassName string `json:"labelClass,omitempty"`
	Color          string `json:"color,omitempty"` // indicator color override
	// DefaultIndex shows the tooltip on mount at that category.
	DefaultIndex *int `json:"defaultIndex,omitempty"`
	// TooltipClass is the gsxui addition: the compiled per-style tooltip
	// recipe's class string, for the client to stamp on the tooltip
	// element. Empty in this task — nothing declares that recipe slot yet
	// (see this file's own doc comment on ChartTooltip) — Task 5 populates
	// it once the recipe exists.
	TooltipClass string `json:"tooltipClass,omitempty"`
}

// ChartModelSeries is one data series with its resolved color variable —
// the pendant of templui's ModelSeries, narrowed to the cartesian fields
// (Cells/ActiveIndex/ActiveBar/LabelLists/LabelList need the Cell/LabelList
// children this task doesn't port — not in the brief's explicit child
// list — and Background/CornerRadius/Class/TooltipNames are radial-bar-only;
// Task 4 adds whichever of these its own kinds need).
type ChartModelSeries struct {
	Key         string    `json:"key"`
	Label       string    `json:"label"`
	Color       string    `json:"color"`
	Values      []float64 `json:"values"`
	FillOpacity float64   `json:"fillOpacity,omitempty"` // area: 0 uses Recharts' 0.6
	Curve       string    `json:"curve,omitempty"`       // line/area: "natural" (default), "linear", "step", "monotone"
	Icon        string    `json:"icon,omitempty"`        // rendered svg, replaces the legend/tooltip color swatch
	Fill        string    `json:"fill,omitempty"`        // area: verbatim fill, e.g. url(#fillDesktop)
	Stroke      string    `json:"stroke,omitempty"`      // area/line: verbatim stroke

	Radius      []float64 `json:"radius,omitempty"`      // bar: corner radii, one or four
	StackID     string    `json:"stackId,omitempty"`     // bar/area
	StrokeWidth float64   `json:"strokeWidth,omitempty"` // bar/line

	Dot        *ChartDotModel `json:"dot,omitempty"`        // line: per point dots
	ActiveDotR float64        `json:"activeDotR,omitempty"` // line: hover dot radius
}

// ChartDotModel describes the per point dots of a line.
type ChartDotModel struct {
	R           float64 `json:"r,omitempty"`
	FillOpacity float64 `json:"fillOpacity,omitempty"`
	Fill        string  `json:"fill,omitempty"`
	Size        float64 `json:"size,omitempty"`
}

// chartModelSeries resolves the shared fields of one series (color, label,
// values) — the pendant of templui's modelSeries, byte-faithful.
func chartModelSeries(config ChartConfig, key, fill string, fillOpacity float64, data []ChartDatum) ChartModelSeries {
	color := fill
	if color == "" || color == "gradient" {
		color = chartSeriesColor(key)
	}
	values := make([]float64, len(data))
	for i, d := range data {
		values[i] = chartNum(d[key])
	}
	return ChartModelSeries{Key: key, Label: config.label(key), Color: color, Values: values, FillOpacity: fillOpacity}
}

// buildChartModel normalizes the collected chart tree into the model the
// client renderer consumes, the pendant of templui's buildModel narrowed
// to the cartesian kinds this task ports (bar, line, area); a later task
// extends the kind switch for pie/radar/radial.
func buildChartModel(ctx context.Context, st *chartState) ChartModel {
	config := chartConfigFromCtx(ctx)
	m := ChartModel{Kind: st.kind, StackOffset: st.stackOffset}

	if g := st.grid; g != nil {
		m.Grid = true
		m.GridHorizontal = chartBoolOr(g.Horizontal, true)
		m.GridVertical = chartBoolOr(g.Vertical, true)
	}
	if y := st.y; y != nil {
		m.YAxisWidth = y.opts.Width
		if m.YAxisWidth == 0 {
			m.YAxisWidth = 60
		}
		m.YAxisMargin = y.opts.TickMargin
		m.TickCount = y.opts.TickCount
		m.YTickLine = y.opts.TickLine
		m.YAxisLine = y.opts.AxisLine
		m.YAxisHide = y.opts.Hide
		if y.opts.Hide {
			m.YAxisWidth = 0
		}
	}
	if st.margin != nil {
		m.MarginTop, m.MarginRight, m.MarginBottom, m.MarginLeft = st.margin.Top, st.margin.Right, st.margin.Bottom, st.margin.Left
	} else {
		m.MarginTop, m.MarginRight, m.MarginBottom, m.MarginLeft = 5, 5, 5, 5
	}
	if x := st.x; x != nil {
		m.XAxisHeight = defaultChartXAxisHeight
		m.TickMargin = x.opts.TickMargin
		m.MinTickGap = x.opts.MinTickGap
		if m.MinTickGap == 0 {
			m.MinTickGap = 5
		}
		m.XTickLine = x.opts.TickLine
		m.XAxisLine = x.opts.AxisLine
		m.XAxisHide = x.opts.Hide
		if x.opts.Hide {
			m.XAxisHeight = 0
		}
	}
	if st.legend != nil {
		m.LegendHeight = defaultChartLegendHeight
		m.LegendVAlign = st.legend.VerticalAlign
	}
	tt := st.tooltip
	if tt != nil {
		m.Cursor = chartBoolOr(tt.opts.Cursor, true)
		m.HasTooltip = true
		m.Tooltip = ChartTooltipModel{
			Indicator: tt.opts.Indicator, HideLabel: tt.opts.HideLabel, HideIndicator: tt.opts.HideIndicator,
			Width: tt.opts.Class, LabelClassName: tt.opts.LabelClass, Color: tt.opts.Color,
			DefaultIndex: tt.opts.DefaultIndex,
		}
		if tt.opts.LabelKey != "" {
			m.Tooltip.Label = config.label(tt.opts.LabelKey)
		}
	}

	// The category axis carries the labels: XAxis in this task's scope —
	// Layout "vertical" (which would read YAxis instead, like upstream)
	// isn't ported, see ChartModel.Layout's own doc comment.
	dataKey := ""
	if x := st.x; x != nil {
		dataKey = x.key
	}
	m.Labels = make([]string, len(st.data))
	m.TooltipLabels = make([]string, len(st.data))
	if dataKey != "" {
		for i, d := range st.data {
			raw := d[dataKey]
			m.Labels[i] = chartStr(raw)
			m.TooltipLabels[i] = chartStr(raw)
		}
	}

	switch st.kind {
	case "bar":
		m.CategoryGap = 0.1
		stacked := false
		for _, b := range st.bars {
			s := chartModelSeries(config, b.key, b.opts.Fill, 0, st.data)
			s.Radius = chartRadiusCorners(b.opts.Radius)
			s.StackID = b.opts.StackID
			s.StrokeWidth = b.opts.StrokeWidth
			if b.opts.StackID != "" {
				stacked = true
			}
			m.Series = append(m.Series, s)
		}
		m.Stacked = stacked
		m.DomainMin = chartSeriesMin(m.Series)
	case "line":
		for _, l := range st.lines {
			s := chartModelSeries(config, l.key, "", 0, st.data)
			s.Curve = string(l.opts.Type)
			s.Stroke = l.opts.Stroke
			s.StrokeWidth = l.opts.StrokeWidth
			s.ActiveDotR = l.opts.ActiveDotR
			if d := l.opts.Dot; d != nil {
				s.Dot = &ChartDotModel{R: d.R, FillOpacity: d.FillOpacity, Fill: d.Fill, Size: d.Size}
			}
			m.Series = append(m.Series, s)
		}
	case "area":
		stacked := false
		for _, a := range st.areas {
			if a.opts.StackID != "" {
				stacked = true
			}
			s := chartModelSeries(config, a.key, "", a.opts.FillOpacity, st.data)
			s.Curve = string(a.opts.Type)
			s.Fill = a.opts.Fill
			s.Stroke = a.opts.Stroke
			s.StackID = a.opts.StackID
			m.Series = append(m.Series, s)
		}
		m.Stacked = stacked
	}

	// nameKey moves the row label to another config entry, like
	// ChartTooltipContent's nameKey.
	if tt != nil && tt.opts.NameKey != "" {
		label := config.label(tt.opts.NameKey)
		for i := range m.Series {
			m.Series[i].Label = label
		}
	}

	if len(st.defs) > 0 {
		m.Defs = make([]ChartLinearGradientModel, len(st.defs))
		for i, d := range st.defs {
			m.Defs[i] = ChartLinearGradientModel{ID: d.id, X1: d.x1, Y1: d.y1, X2: d.x2, Y2: d.y2, Stops: d.stops}
		}
	}

	// applyIcons: copy the config icons into the model series.
	for i := range m.Series {
		for _, sc := range config {
			if sc.Key == m.Series[i].Key && sc.Icon != nil {
				m.Series[i].Icon = chartRenderHTML(ctx, sc.Icon)
			}
		}
	}

	return m
}

// chartModelScript renders the embedded JSON payload the client reads —
// the pendant of templui's ModelScript. json.Marshal HTML-escapes
// <, >, & by default, so a data row's own text cannot break out of the
// script tag.
func chartModelScript(m ChartModel) string {
	b, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return `<script type="application/json" data-gsxui-chart-model>` + string(b) + `</script>`
}

// chartLegendItem is one rendered legend entry, the pendant of templui's
// LegendItem.
type chartLegendItem struct {
	label string
	color string
	icon  string
	value string // the payload value the legend sorts by
}

// chartBuildLegendItems builds the legend payload: one entry per series,
// sorted by key like Recharts' itemSorter default of "value" — the
// pendant of templui's legendItems narrowed to the cartesian kinds (no
// pie branch; Task 4 extends this).
func chartBuildLegendItems(config ChartConfig, m ChartModel, opts *ChartLegendOptions) []chartLegendItem {
	items := make([]chartLegendItem, 0, len(m.Series))
	for _, s := range m.Series {
		label := s.Label
		if opts.NameKey != "" {
			label = config.label(opts.NameKey)
		}
		items = append(items, chartLegendItem{label: label, color: s.Color, icon: s.Icon, value: s.Key})
	}
	sort.SliceStable(items, func(i, j int) bool { return items[i].value < items[j].value })
	return items
}

// ChartLegendContent renders the ChartLegendContent pendant, absolutely
// positioned like Recharts' legend wrapper — the pendant of templui's
// legendContent. It is the ONE part of this task's cartesian tree that
// renders real markup from a non-root registration, which is why it (and
// only it) gets a data-gsxui-slot-chart-legend marker and a shape slot
// (see shapes/chart.go): ChartTooltip stays data-only in this task. Called
// by chartRoot once a chart's children have registered a ChartLegend, not
// meant to be composed directly by a consumer.
component ChartLegendContent(items []chartLegendItem, opts *ChartLegendOptions) {
	{{
		posStyle := "position:absolute;left:0;right:0;bottom:5px"
		pad := "pt-3"
		if opts.VerticalAlign == "top" {
			posStyle = "position:absolute;left:0;right:0;top:5px"
			pad = "pb-3"
		}
	}}
	<div style={ posStyle } data-gsxui-slot-chart-legend>
		<div class={ "flex items-center justify-center gap-4", pad, opts.Class }>
			{ for _, it := range items {
				<div class="flex items-center gap-1.5 [&>svg]:h-3 [&>svg]:w-3 [&>svg]:text-muted-foreground">
					{ if it.icon != "" && !opts.HideIcon {
						{ gsx.Raw(it.icon) }
					} else {
						<div class="h-2 w-2 shrink-0 rounded-[2px]" style={ "background-color:" + it.color }></div>
					} }
					{ it.label }
				</div>
			} }
		</div>
	</div>
}

// chartRoot renders one chart root: templui's AreaChart/BarChart/LineChart
// wrapper div, its children (which register into the chartState seeded
// here), then the model script and (if registered) the server-rendered
// legend — templui's own chartOutput, run after children the same way gsx
// renders a Node's children before the statements that follow { children }
// in its caller.
func chartRoot(kind string, data []ChartDatum, margin *ChartMargin, stackOffset string, children gsx.Node, attrs gsx.Attrs) gsx.Node {
	return gsx.Func(func(ctx context.Context, w io.Writer) error {
		st := &chartState{kind: kind, data: data, margin: margin, stackOffset: stackOffset}
		ctx = context.WithValue(ctx, chartStateCtxKey, st)
		gw := gsx.W(w)
		gw.S("<div")
		gw.StyleMerged("position:relative;width:100%;height:100%", attrs.Style())
		gw.Spread(ctx, "div", attrs, gsx.AttrSinks{}, []string{"style"})
		gw.S(">")
		gw.Node(ctx, children)
		m := buildChartModel(ctx, st)
		gw.S(chartModelScript(m))
		if st.legend != nil {
			items := chartBuildLegendItems(chartConfigFromCtx(ctx), m, st.legend)
			gw.Node(ctx, ChartLegendContent(items, st.legend))
		}
		gw.S("</div>")
		return gw.Err()
	})
}

// BarChart is the Recharts BarChart root: its parts declare axes, grid,
// tooltip and bars, the root collects them and emits the model payload
// for the client renderer.
func BarChart(data []ChartDatum, opts *ChartBarChartOptions, children gsx.Node, attrs gsx.Attrs) gsx.Node {
	var margin *ChartMargin
	if opts != nil {
		margin = opts.Margin
	}
	return chartRoot("bar", data, margin, "", children, attrs)
}

// LineChart is the Recharts LineChart root.
func LineChart(data []ChartDatum, opts *ChartLineChartOptions, children gsx.Node, attrs gsx.Attrs) gsx.Node {
	var margin *ChartMargin
	if opts != nil {
		margin = opts.Margin
	}
	return chartRoot("line", data, margin, "", children, attrs)
}

// AreaChart is the Recharts AreaChart root.
func AreaChart(data []ChartDatum, opts *ChartAreaChartOptions, children gsx.Node, attrs gsx.Attrs) gsx.Node {
	var margin *ChartMargin
	stackOffset := ""
	if opts != nil {
		margin = opts.Margin
		stackOffset = opts.StackOffset
	}
	return chartRoot("area", data, margin, stackOffset, children, attrs)
}

// ChartCartesianGrid registers the grid, the pendant of Recharts'
// CartesianGrid element. It renders nothing; the client draws it.
func ChartCartesianGrid(opts *ChartCartesianGridOptions) gsx.Node {
	if opts == nil {
		opts = &ChartCartesianGridOptions{}
	}
	return gsx.Func(func(ctx context.Context, _ io.Writer) error {
		if st := chartStateFromCtx(ctx); st != nil {
			st.grid = opts
		}
		return nil
	})
}

// ChartXAxis registers the x axis.
func ChartXAxis(key string, opts *ChartXAxisOptions) gsx.Node {
	if opts == nil {
		opts = &ChartXAxisOptions{}
	}
	return gsx.Func(func(ctx context.Context, _ io.Writer) error {
		if st := chartStateFromCtx(ctx); st != nil {
			st.x = &chartAxisReg{key: key, opts: *opts}
		}
		return nil
	})
}

// ChartYAxis registers the y axis.
func ChartYAxis(key string, opts *ChartYAxisOptions) gsx.Node {
	if opts == nil {
		opts = &ChartYAxisOptions{}
	}
	return gsx.Func(func(ctx context.Context, _ io.Writer) error {
		if st := chartStateFromCtx(ctx); st != nil {
			st.y = &chartYAxisReg{key: key, opts: *opts}
		}
		return nil
	})
}

// ChartTooltip registers the tooltip's config; the client renders the
// actual tooltip markup from the model (see ChartTooltipModel's own doc
// comment on why this task doesn't render it server-side).
func ChartTooltip(opts *ChartTooltipOptions) gsx.Node {
	if opts == nil {
		opts = &ChartTooltipOptions{}
	}
	return gsx.Func(func(ctx context.Context, _ io.Writer) error {
		if st := chartStateFromCtx(ctx); st != nil {
			st.tooltip = &chartTooltipReg{opts: *opts}
		}
		return nil
	})
}

// ChartLegend registers the legend; the enclosing chart root renders it
// server-side after its other children (see ChartLegendContent).
func ChartLegend(opts *ChartLegendOptions) gsx.Node {
	if opts == nil {
		opts = &ChartLegendOptions{}
	}
	return gsx.Func(func(ctx context.Context, _ io.Writer) error {
		if st := chartStateFromCtx(ctx); st != nil {
			st.legend = opts
		}
		return nil
	})
}

// ChartDefs groups gradient definitions like the svg defs element —
// nothing renders visibly server-side (the client draws the SVG), so this
// only needs to render its children to let any ChartLinearGradient among
// them register.
func ChartDefs(children gsx.Node) gsx.Node {
	return gsx.Func(func(ctx context.Context, w io.Writer) error {
		if children == nil {
			return nil
		}
		return children.Render(ctx, w)
	})
}

// ChartLinearGradient registers one gradient definition. Its children are
// the raw stop elements, captured and passed through into the registered
// definition verbatim, like Recharts passes defs children through.
func ChartLinearGradient(opts ChartLinearGradientOptions, stops gsx.Node) gsx.Node {
	return gsx.Func(func(ctx context.Context, _ io.Writer) error {
		st := chartStateFromCtx(ctx)
		if st == nil {
			return nil
		}
		var sb strings.Builder
		if stops != nil {
			if err := stops.Render(ctx, &sb); err != nil {
				return err
			}
		}
		st.defs = append(st.defs, chartGradientReg{
			id: opts.ID, x1: opts.X1, y1: opts.Y1, x2: opts.X2, y2: opts.Y2,
			stops: strings.TrimSpace(sb.String()),
		})
		return nil
	})
}

// ChartArea registers one area series. Declaration order is paint order,
// like in Recharts.
func ChartArea(key string, opts *ChartAreaOptions) gsx.Node {
	if opts == nil {
		opts = &ChartAreaOptions{}
	}
	return gsx.Func(func(ctx context.Context, _ io.Writer) error {
		if st := chartStateFromCtx(ctx); st != nil {
			st.areas = append(st.areas, chartAreaReg{key: key, opts: *opts})
		}
		return nil
	})
}

// ChartBar registers one bar series.
func ChartBar(key string, opts *ChartBarOptions) gsx.Node {
	if opts == nil {
		opts = &ChartBarOptions{}
	}
	return gsx.Func(func(ctx context.Context, _ io.Writer) error {
		if st := chartStateFromCtx(ctx); st != nil {
			st.bars = append(st.bars, chartBarReg{key: key, opts: *opts})
		}
		return nil
	})
}

// ChartLine registers one line series.
func ChartLine(key string, opts *ChartLineOptions) gsx.Node {
	if opts == nil {
		opts = &ChartLineOptions{}
	}
	return gsx.Func(func(ctx context.Context, _ io.Writer) error {
		if st := chartStateFromCtx(ctx); st != nil {
			st.lines = append(st.lines, chartLineReg{key: key, opts: *opts})
		}
		return nil
	})
}
