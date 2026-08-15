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
// custom properties scoped to its own data-chart attribute (styleSchemes).
type ChartConfig []ChartSeries

// ChartSeries is one entry of a ChartConfig. Key identifies the series both
// here and in the data rows the chart roots read (rows are matched against
// Key, never Label). Theme, present, takes precedence over Color for both
// color schemes — see ChartConfig.styleSchemes. Icon replaces the color
// swatch in the legend and tooltip.
//
// Color (and Theme's pair) are gsx.RawCSS, not string: they land inside a
// <style> block as CSS values, and the values every shadcn demo uses —
// var(--chart-1), oklch(...), url(#gradient) — are exactly what gsx's CSS
// value filter must reject from untrusted input (parens, `--` runs). The
// type is the trust boundary: an author writing Color: "var(--chart-1)" is
// vouching for CSS they wrote, and a caller deriving a color from user
// input must convert to gsx.RawCSS deliberately, in sight of that decision.
// Key stays a plain string and IS filtered — a hostile key cannot break
// out of its declaration.
type ChartSeries struct {
	Key   string
	Label string
	Color gsx.RawCSS
	Theme *ChartSeriesTheme
	Icon  gsx.Node
}

// ChartSeriesTheme is a color pair for the light and dark color schemes,
// the pendant of the `theme` field in shadcn's own ChartConfig entries.
type ChartSeriesTheme struct {
	Light gsx.RawCSS
	Dark  gsx.RawCSS
}

// styleSchemes is the ChartStyle pendant: one rule per color scheme, each
// scoped by [data-chart=<id>]. Ported byte-faithful (semantics, not syntax)
// from templui/shadcn-templ v2's chart.templ Config.styleBlock (see this
// file's header credit):
//
//   - a theme color wins over the plain Color for both schemes;
//   - a theme with no value for a given scheme falls back to Color, like
//     itemConfig.theme?.[theme] ?? itemConfig.color upstream;
//   - a series with neither Color nor Theme contributes no line, but does
//     not block the other series' lines;
//   - when NO series in the whole config carries a color, styleSchemes
//     returns "" — no <style> tag at all — matching ChartStyle returning
//     null rather than an empty one upstream.
//
// chartStyleScheme is one color scheme's resolved rule body for the Chart
// container's <style> block: prefix selects the scheme ("" light, ".dark "
// dark) and decls is the --color-<key> declaration list, already a
// gsx.RawCSS because it is composed from css` ` literals — the RawCSS
// Color passes verbatim (the author's vouch) while the plain-string key
// runs through gsx's CSS value filter inside the literal.
type chartStyleScheme struct {
	prefix string
	decls  gsx.RawCSS
}

// styleSchemes resolves the ChartStyle pendant: one scheme per color mode,
// theme color winning over the plain color, a theme without a value for
// the scheme falling back to the color. Nil when no series carries any
// color, mirroring ChartStyle returning null upstream.
func (c ChartConfig) styleSchemes() []chartStyleScheme {
	hasColor := false
	for _, s := range c {
		if s.Color != "" || s.Theme != nil {
			hasColor = true
			break
		}
	}
	if !hasColor {
		return nil
	}
	var out []chartStyleScheme
	for _, scheme := range []struct{ Name, Prefix string }{{"light", ""}, {"dark", ".dark "}} {
		var decls gsx.RawCSS
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
				decls += css`--color-@{s.Key}: @{color};`
			}
		}
		out = append(out, chartStyleScheme{prefix: scheme.Prefix, decls: decls})
	}
	return out
}

// chartStyleRules composes the container's <style> body: one
// [data-chart=<id>] rule per scheme. Each piece is a css` ` literal, so
// the result is gsx.RawCSS the <style> block emits verbatim: prefix and id
// are ours (a fixed scheme prefix, a package counter), decls already
// carry each key through the value filter.
func chartStyleRules(id string, schemes []chartStyleScheme) gsx.RawCSS {
	var out gsx.RawCSS
	for _, sc := range schemes {
		out += css`@{gsx.RawCSS(sc.prefix)}[data-chart=@{gsx.RawCSS(id)}]{@{sc.decls}}`
	}
	return out
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
// this <div> with the <style> block the container renders from styleSchemes right below it in the
// SAME render call, so it has to be computed inside Chart, before attrs is
// spread, with nothing to key off but the render call itself.
//
// This is a process-lifetime monotonic counter: deterministic (no math/rand
// or crypto/rand — the brief is explicit that pinned tests need
// determinism), safe under concurrent renders (atomic.Uint64), and unique
// for as long as the package's process runs, which is what actually matters
// for the [data-chart=...] selector not leaking one chart's colors
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
// styleSchemes) for the chart roots and parts a later task adds as children.
// templui's own class also carries a wall of [&_.recharts-*] arbitrary-
// variant selectors (axis tick fill, grid stroke, cursor fill, dot/sector
// rings, radial-bar background) that re-color Recharts' own hardcoded SVG
// sentinels (#666/#ccc/#fff/#eee) — ui/chart.render.js is a byte-faithful
// port of templui's chart.js and keeps emitting those exact sentinels
// verbatim, so this wall is the ONLY mechanism that themes them. It lives
// in the root recipe class the div below stamps (registry/styles/<style>/
// chart.css's root rule), not spelled out here — same split every other
// component's recipe class follows.
//
// children is expected to be exactly ONE chart root (BarChart/LineChart/
// AreaChart/PieChart/RadarChart/RadialBarChart), matching upstream's own
// ChartContainer usage — every shadcn/templui demo wraps a single root.
// Nothing here rejects a second root nested alongside the first, but
// ui/chart.render.js's renderChart(container) only ever looks up the
// FIRST `script[data-gsxui-chart-model]` under this div (querySelector,
// not querySelectorAll), and every root's own per-container bookkeeping
// (container._gsxuiActive/_gsxuiActivePoints, the entrance/update-morph
// state chartRoot's caller reads back across renders) is keyed on this
// SAME container object with no per-root namespacing, so a second root
// would silently never draw and could cross-talk the first root's morph
// state besides. One Chart, one root.
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
		class={
			"flex aspect-video justify-center text-xs min-h-0 [&_.recharts-cartesian-axis-tick_text]:fill-muted-foreground [&_.recharts-cartesian-grid_line[stroke='#ccc']]:stroke-border/50 [&_.recharts-curve.recharts-tooltip-cursor]:stroke-border [&_.recharts-dot[stroke='#fff']]:stroke-transparent [&_.recharts-layer]:outline-hidden [&_.recharts-polar-grid_[stroke='#ccc']]:stroke-border [&_.recharts-radial-bar-background-sector]:fill-muted [&_.recharts-rectangle.recharts-tooltip-cursor]:fill-muted [&_.recharts-reference-line_[stroke='#ccc']]:stroke-border [&_.recharts-sector]:outline-hidden [&_.recharts-sector[stroke='#fff']]:stroke-transparent [&_.recharts-surface]:outline-hidden"
		}
		{ attrs... }
		data-gsxui-slot-chart
	>
		{ if schemes := config.styleSchemes(); schemes != nil {
			<style>
				@{ chartStyleRules(id, schemes) }
			</style>
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
// ...) in ~2100 lines of ported code no one wants to rewrite. The one
// gsxui-specific departure is the script tag's attribute name
// (data-gsxui-chart-model, not data-tui-chart-model). The tooltip's own
// chrome does NOT travel through this model at all — Task 6 renders it as
// real server markup instead (ChartTooltipTemplate, below chartRoot) that
// ui/chart.render.js reads directly from the DOM, so ChartModel carries no
// tooltip-class field of any kind (an earlier draft of this field,
// ChartTooltipModel.TooltipClass, never shipped: see that struct's own
// doc comment).

// ChartDatum is one row of chart data, the pendant of the plain objects in
// shadcn's chartData arrays and templui's own Datum. Rows are read back by
// key (chartNum/chartStr) while building the model; they are not
// serialized into it themselves — Labels/TooltipLabels/per-series Values
// carry what the client needs, exactly like templui's own Model.
type ChartDatum map[string]any

// chartMargin is the chart's plot margin (Recharts' default of 5 on every
// side): a chart root's own marginTop/marginRight/marginBottom/marginLeft
// tag params, resolved through chartMarginOf below. A directly constructed
// chartState (nothing in this file does so with a nil margin anymore, but
// buildChartModel/buildChartRadarModel/buildChartRadialModel still handle
// it) treats a nil *chartMargin the same way: the model's margin fields
// fall back to that same 5/5/5/5.
type chartMargin struct {
	Top, Right, Bottom, Left float64
}

// chartFloatPtrOr converts a flat float64 tag parameter into the *float64
// a struct field that treats zero as "unset, leave the SVG default in
// place" (chartRadarOptions.FillOpacity) still expects internally: zero
// becomes nil, any other value becomes a pointer to itself.
func chartFloatPtrOr(v float64) *float64 {
	if v == 0 {
		return nil
	}
	return &v
}

// chartMarginOf resolves the four flat margin tag parameters into a
// chartState margin, replacing Recharts' own 5/5/5/5 default AS A SET the
// moment any one of the four is customized — never per side independently.
// This mirrors Recharts' own margin prop (ported byte-faithful by templui,
// credited in this file's header): the prop is a single object, so
// supplying `{left: 12, right: 12}` leaves `top`/`bottom` at a literal
// `undefined` (JS) — `0` here — not at Recharts' own class-level default,
// exactly like `<ui.AreaChart marginLeft={12} marginRight={12}>` (this
// file's own area demo) must render marginTop/marginBottom 0, not 5.
// Independent per-side zero-sentinel defaulting (an earlier draft of this
// function) got that wrong: it revived the 5 default per side rather than
// per call, a real 5px parity regression on that exact shipped demo,
// invisible to a byte pin that only ever exercised the all-zero path.
// Always non-nil: buildChartModel/buildChartRadarModel/
// buildChartRadialModel's own `if st.margin != nil` branch (unchanged from
// before this task) is still correct for a caller who constructs a
// chartState directly, it just never sees nil from a chart root's own
// component body anymore, matching the same 5/5/5/5 output their nil
// branch already produced for a caller who set no margin at all.
func chartMarginOf(top, right, bottom, left float64) *chartMargin {
	if top == 0 && right == 0 && bottom == 0 && left == 0 {
		return &chartMargin{Top: 5, Right: 5, Bottom: 5, Left: 5}
	}
	return &chartMargin{Top: top, Right: right, Bottom: bottom, Left: left}
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
	margin      *chartMargin
	stackOffset string
	grid        *chartGridOptions
	x           *chartAxisReg
	y           *chartYAxisReg
	tooltip     *chartTooltipReg
	legend      *chartLegendOptions
	defs        []chartGradientReg
	bars        []chartBarReg
	lines       []chartLineReg
	areas       []chartAreaReg

	// Polar collectors — pie, radar, radial-bar. center (the Label child
	// upstream mounts through Pie or PolarRadiusAxis) is not ported: the
	// Label child itself is out of this task's explicit child list.
	pies       []chartPieReg
	polarGrid  *chartPolarGridOptions
	angleAxis  *chartPolarAngleAxisReg
	radars     []chartRadarReg
	radialBars []chartRadialBarReg
	radiusAxis *chartPolarRadiusAxisOptions

	// RadialBarChart's own geometry, the pendant of upstream's own
	// chartState fields of the same name.
	startAngle  float64
	endAngle    *float64
	innerRadius float64
	outerRadius float64
}

type chartAxisReg struct {
	key  string
	opts chartXAxisOptions
}

type chartYAxisReg struct {
	key  string
	opts chartYAxisOptions
}

type chartTooltipReg struct {
	opts chartTooltipOptions
}

type chartGradientReg struct {
	id, x1, y1, x2, y2, stops string
}

type chartBarReg struct {
	key  string
	opts chartBarOptions
}

type chartLineReg struct {
	key  string
	opts chartLineOptions
}

type chartAreaReg struct {
	key  string
	opts chartAreaOptions
}

// chartPieReg is one registered ChartPie: unlike the cartesian series,
// Recharts' own Pie takes its own data prop rather than reading the
// enclosing chart's rows, so it carries its data alongside its key.
type chartPieReg struct {
	data []ChartDatum
	key  string
	opts chartPieOptions
}

// chartPolarAngleAxisReg is one registered ChartPolarAngleAxis. Recharts'
// Tick render-prop is a Go closure and cannot cross the JSON boundary to
// the client renderer (the same reasoning chartTooltipOptions' dropped
// Formatter uses), so key is all this needs.
type chartPolarAngleAxisReg struct {
	key string
}

type chartRadarReg struct {
	key  string
	opts chartRadarOptions
}

type chartRadialBarReg struct {
	key  string
	opts chartRadialBarOptions
}

// chartGridOptions is the pendant of Recharts' CartesianGrid. Both
// directions default to false (gsxui's own convention, ADAPTed from
// Recharts' own default-true — see ChartCartesianGrid's doc comment):
// a bare `<ui.ChartCartesianGrid/>` draws no lines at all, matching a
// caller having to opt in to each direction explicitly.
type chartGridOptions struct {
	Horizontal bool
	Vertical   bool
}

// chartXAxisOptions is the pendant of Recharts' XAxis, minus the render-prop
// fields (TickFormatter) a JSON model boundary cannot carry.
type chartXAxisOptions struct {
	Hide       bool
	TickLine   bool
	AxisLine   bool
	TickMargin float64
	MinTickGap float64 // defaults to Recharts' 5
}

// chartYAxisOptions is the pendant of Recharts' YAxis.
type chartYAxisOptions struct {
	Hide       bool
	TickLine   bool
	AxisLine   bool
	TickMargin float64
	TickCount  int     // defaults to Recharts' 5
	Width      float64 // defaults to Recharts' 60
}

// chartTooltipOptions is the pendant of ChartTooltip plus
// ChartTooltipContent's non-function fields — Formatter/LabelFormatter are
// deferred (see this task's report: a Go closure cannot cross the JSON
// boundary to the client renderer, and the model already carries the raw
// row data those would have formatted from).
type chartTooltipOptions struct {
	// Cursor defaults to false (gsxui's own ADAPT of Recharts' own
	// default-true — see ChartTooltip's doc comment): a bare
	// `<ui.ChartTooltip/>` draws no hover cursor highlight.
	Cursor        bool
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

// chartLegendOptions is the pendant of ChartLegend plus
// ChartLegendContent's props.
type chartLegendOptions struct {
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

// chartAreaOptions is the pendant of one Recharts Area.
type chartAreaOptions struct {
	// Curve is Recharts' curve prop: "natural", "linear" or "step"; empty
	// renders as Recharts' own "linear" default, unchanged by this task.
	Curve string
	// Fill and Stroke are used verbatim, e.g. "url(#fillDesktop)" or
	// "var(--color-desktop)". Empty falls back to the series color.
	Fill        string
	Stroke      string
	StackID     string
	FillOpacity float64 // 0 uses Recharts' default 0.6
}

// chartBarOptions is the pendant of one Recharts Bar.
type chartBarOptions struct {
	Fill    string
	StackID string
	// Radius is Recharts' radius union: a float64 for all corners or a
	// []float64 of four corners, e.g. []float64{0, 0, 4, 4}.
	Radius      any
	StrokeWidth float64
}

// chartDotOptions is the pendant of Recharts' dot prop's object form.
type chartDotOptions struct {
	R           float64
	FillOpacity float64
	Fill        string
	Size        float64 // icon box, unused without an Icon renderer
}

// chartLineOptions is the pendant of one Recharts Line.
type chartLineOptions struct {
	// Curve is Recharts' curve prop: "natural", "linear" or "step"; empty
	// renders as Recharts' own "linear" default, unchanged by this task.
	Curve       string
	Stroke      string
	StrokeWidth float64
	// Dot draws the per point dots; nil is the demos' dot={false}.
	Dot *chartDotOptions
	// ActiveDotR sizes the hover dot, like Recharts' activeDot prop.
	ActiveDotR float64
}

// Layout (Recharts' vertical bar layout) and AccessibilityLayer stay
// deferred on every chart root below — neither was in the shipped task's
// port list (Margin, curve, StackOffset), and this task only reshapes how
// what already ported is exposed, not what is ported.

// chartPolarGridOptions is the pendant of Recharts' PolarGrid.
type chartPolarGridOptions struct {
	// GridType is "polygon" (default) or "circle".
	GridType string
	// RadialLines defaults to false (gsxui's own ADAPT of Recharts' own
	// default-true — see ChartPolarGrid's doc comment).
	RadialLines bool
	// PolarRadius replaces the radius axis ticks with fixed radii.
	PolarRadius []float64
	// Stroke is the grid color, Recharts' "#ccc"; a radial chart demo
	// passes "none" and paints the rings through Class instead.
	Stroke      string
	StrokeWidth float64
	Class       string
}

// chartPolarRadiusAxisOptions is the pendant of Recharts' PolarRadiusAxis.
// Tick, TickLine and AxisLine all default to Recharts' true upstream, but
// (like every other bool this task ADAPTs) default false here. Upstream
// itself never reads these three fields back out of its own registered
// state either (grep confirms): PolarRadiusAxis exists there only as a
// mount point for a center Label child, which this task doesn't port (not
// in the brief's explicit child list) — so these are captured here for
// structural fidelity but currently as inert as they are upstream, and the
// default flip changes nothing observable.
type chartPolarRadiusAxisOptions struct {
	Tick     bool
	TickLine bool
	AxisLine bool
}

// chartRadarOptions is the pendant of one Recharts Radar.
type chartRadarOptions struct {
	Fill string
	// FillOpacity is a pointer because zero is a meaningful value: an
	// unset option leaves the SVG default of 1 in place — built from
	// ChartRadar's own flat fillOpacity float64 param via
	// chartFloatPtrOr's zero-sentinel convention, nil when the param is
	// left at zero.
	FillOpacity *float64
	Stroke      string
	StrokeWidth float64
	Dot         *chartDotOptions
}

// chartRadialBarOptions is the pendant of one Recharts RadialBar.
type chartRadialBarOptions struct {
	Fill string
	// Background draws the track behind the bar, over the full angle range.
	Background   bool
	CornerRadius float64
	StackID      string
	Class        string
}

// chartPieLabelOptions is the pendant of the label option of a Pie.
type chartPieLabelOptions struct {
	// Fill overrides the slice color the labels inherit.
	Fill string
}

// chartPieOptions is the pendant of one Recharts Pie. ActiveIndex and
// ActiveShape (Recharts' active-sector props) are deferred alongside
// Label/LabelList's own list-registering counterparts — the same
// reasoning chartBarOptions' own dropped ActiveIndex/ActiveBar uses: not
// in the brief's explicit child list for this task.
type chartPieOptions struct {
	NameKey     string
	InnerRadius float64
	// OuterRadius defaults to Recharts' 80% of the available radius.
	OuterRadius float64
	StrokeWidth float64
	// Stroke is the separator between the slices, Recharts' "#fff"; pass
	// "0" to drop it.
	Stroke string
	// Label renders the value labels outside the slices — nil (ChartPie's
	// own flat label bool left false) shows none, matching every shipped
	// demo, none of which passes Recharts' own label prop either.
	Label *chartPieLabelOptions
	// LabelLine draws the connecting line from a slice to its outside
	// label; defaults false (gsxui's own ADAPT of Recharts' own
	// default-true — see ChartPie's doc comment). Inert without Label
	// also set: chart.render.js only draws either when Label is present
	// (see ui/chart.render.js's own pie.label gate), so this default flip
	// changes no shipped demo's rendered pixels.
	LabelLine bool
}

// label resolves key's config entry, the ChartConfig pendant of templui's
// Config.Label: it returns that entry's own Label field verbatim (which
// may be empty) when key matches, falling back to key itself only when NO
// entry names it at all.
// icon returns the series' Icon node for key, nil when the series has none
// (or the key is not configured) — the legend renders it in place.
func (c ChartConfig) icon(key string) gsx.Node {
	for _, s := range c {
		if s.Key == key {
			return s.Icon
		}
	}
	return nil
}

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

// chartPayloadConfigKey is the pendant of templui's payloadConfigKey
// (itself getPayloadConfigFromPayload): the key names the config entry,
// unless the data row carries a string under that key, then its own value
// names it instead.
func chartPayloadConfigKey(row ChartDatum, key string) string {
	if row != nil {
		if v, ok := row[key]; ok {
			if sv, isString := v.(string); isString {
				return sv
			}
		}
	}
	return key
}

// chartHasFillColumn reports whether the rows carry their own "fill",
// which Recharts uses as that row's own fill — the pendant of templui's
// hasFillColumn, byte-faithful.
func chartHasFillColumn(data []ChartDatum) bool {
	for _, row := range data {
		if _, ok := row["fill"]; ok {
			return true
		}
	}
	return false
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
// (chart.templ:1645), covering every kind this file builds (bar, line,
// area, pie, radar, radial-bar) — same field order, same json tags, same
// omitempty-or-not as upstream, including its two dead fields (Radius and
// SliceColors, grep-confirmed unassigned anywhere in chart.templ) kept
// here for structural fidelity only.
//
// The one gsxui-specific field is the script tag's attribute name
// (chartModelScript uses data-gsxui-chart-model, not templui's
// data-tui-chart-model) — every other field, its json tag, its
// omitempty-or-not, and its Recharts default is ported byte-faithful so
// Task 5's adapted chart.js — which reads these exact flat fields
// (marginTop, xAxisHeight, tickMargin, minTickGap, xTickLine, ... per
// chart.templ:1645) unmodified — needs no rewriting. The tooltip's own
// chrome (its per-style recipe class, its indicator swatch markup) never
// travels through this model at all: Task 6 renders it as real server
// markup instead (see ChartTooltipTemplate, below chartRoot), which
// ui/chart.render.js reads directly off the DOM rather than off JSON —
// see ChartTooltipModel's own doc comment for why an earlier draft field
// carrying it through the model (TooltipClass) never shipped.
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
	// Radius is dead upstream too (grep confirms nothing ever assigns
	// m.Radius in chart.templ) — ported for structural fidelity only, per
	// this file's header credit.
	Radius float64 `json:"radius,omitempty"`

	Grid           bool `json:"grid,omitempty"`
	GridHorizontal bool `json:"gridHorizontal,omitempty"`
	GridVertical   bool `json:"gridVertical,omitempty"`

	// Layout "vertical" swaps the axes and draws the bars horizontally —
	// out of this task's scope (no BarChart tag param can set it, see
	// docs/jsx-parity.md's own `## chart` GAP entry), so this always
	// marshals as absent for now.
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
	HasTooltip bool     `json:"hasTooltip,omitempty"`
	Labels     []string `json:"labels"`
	// TooltipLabels is dead for "pie" and "radial" (upstream's own
	// buildPieModel/buildRadialModel never set it either — a pie/radial
	// chart's tooltip names come from each series' own TooltipNames
	// instead, see ChartModelSeries).
	TooltipLabels []string `json:"tooltipLabels,omitempty"`
	// SliceColors is dead upstream too (grep confirms nothing ever assigns
	// m.SliceColors) — ported for structural fidelity only, like Radius
	// above; each pie's own per-slice Colors carries what the client needs
	// (see ChartPieModel).
	SliceColors []string           `json:"sliceColors,omitempty"`
	Series      []ChartModelSeries `json:"series"`
	Tooltip     ChartTooltipModel  `json:"tooltip"`
	// Pies carries the pie geometry for the client renderer.
	Pies []ChartPieModel `json:"pies,omitempty"`
	// Polar carries the radar geometry.
	Polar *ChartPolarModel `json:"polar,omitempty"`
	// Radial carries the RadialBarChart geometry.
	Radial *ChartRadialModel `json:"radial,omitempty"`

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
// chartTooltipOptions.
//
// An earlier draft carried a TooltipClass field here (the compiled
// per-style tooltip recipe's class string, for the client to stamp onto
// the tooltip element) — Task 5 found stylegen's checkShapeCoverage
// structurally rejects a recipe-accessor call made from plain Go code
// (buildChartModel is not a component body), so no non-hack route existed
// to populate it, and it stayed permanently empty (BLOCKED, see Task 5's
// own report). Task 6 drops the field rather than resurrect it: the
// tooltip's chrome is server-rendered real markup instead
// (ChartTooltipTemplate, below chartRoot), which is exactly a component
// body with a real class attribute carrying the tooltip slot's recipe
// accessor call — the JSON model boundary was never the right place for a
// compiled class string to begin with.
// Upstream's TooltipModel.Rows (precomputed per-row formatted display
// strings) is not ported — it is derived from the Formatter closures this
// file also drops, which cannot cross the JSON boundary to the client
// renderer. The renderer formats raw values/labels itself instead.
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
}

// ChartModelSeries is one data series with its resolved color variable —
// the pendant of templui's ModelSeries, narrowed to the fields this file's
// kinds actually need (Cells/ActiveIndex/ActiveBar/LabelLists/LabelList
// need the Cell/LabelList children no task ports yet — not in either
// task's explicit child list).
type ChartModelSeries struct {
	Key         string    `json:"key"`
	Label       string    `json:"label"`
	Color       string    `json:"color"`
	Values      []float64 `json:"values"`
	FillOpacity float64   `json:"fillOpacity,omitempty"` // area: 0 uses Recharts' 0.6
	Curve       string    `json:"curve,omitempty"`       // line/area: "natural" (default), "linear", "step", "monotone"
	// FillOpacityPtr is radar's own opacity: an explicit zero stays,
	// unlike FillOpacity above where zero falls back to Recharts' default.
	FillOpacityPtr *float64 `json:"fillOpacityPtr,omitempty"`
	Icon           string   `json:"icon,omitempty"`   // rendered svg, replaces the legend/tooltip color swatch
	Fill           string   `json:"fill,omitempty"`   // area: verbatim fill, e.g. url(#fillDesktop)
	Stroke         string   `json:"stroke,omitempty"` // area/line/radar: verbatim stroke

	Radius  []float64 `json:"radius,omitempty"`  // bar: corner radii, one or four
	StackID string    `json:"stackId,omitempty"` // bar/area/radial-bar
	// Cells is a fill per data row, from that row's own "fill" column —
	// pure data access, ported for both bar and radial-bar (see
	// chartHasFillColumn); the <Cell>-child-driven form both kinds also
	// support upstream stays deferred, matching this file's own Bar/
	// radial-bar scope-narrowing elsewhere.
	Cells       []string `json:"cells,omitempty"`
	StrokeWidth float64  `json:"strokeWidth,omitempty"` // bar/line/radar

	Dot        *ChartDotModel `json:"dot,omitempty"`        // line/radar: per point dots
	ActiveDotR float64        `json:"activeDotR,omitempty"` // line: hover dot radius

	// radial-bar: the track behind the bar, its corner radius, the class
	// upstream puts on the sectors, and each row's own resolved tooltip
	// name (getPayloadConfigFromPayload).
	Background   bool     `json:"background,omitempty"`
	CornerRadius float64  `json:"cornerRadius,omitempty"`
	Class        string   `json:"class,omitempty"`
	TooltipNames []string `json:"tooltipNames,omitempty"`
}

// ChartDotModel describes the per point dots of a line or radar.
type ChartDotModel struct {
	R           float64 `json:"r,omitempty"`
	FillOpacity float64 `json:"fillOpacity,omitempty"`
	Fill        string  `json:"fill,omitempty"`
	Size        float64 `json:"size,omitempty"`
}

// ChartTickContentModel is a custom axis tick for the client renderer —
// the pendant of templui's TickContentModel. Nothing populates it in this
// task (PolarAngleAxis's Tick render-prop is a Go closure and cannot cross
// the JSON boundary — the same reasoning chartTooltipOptions' dropped
// Formatter uses), but the renderer itself still reads polar.ticks[i]
// separately from its own default-grid fallback (chart.js's own
// `const custom = polar.ticks && polar.ticks[i]`), so the field stays —
// ported for structural fidelity only, like Radius/SliceColors.
type ChartTickContentModel struct {
	Spans      []ChartTickSpanModel `json:"spans"`
	FontSize   float64              `json:"fontSize,omitempty"`
	FontWeight string               `json:"fontWeight,omitempty"`
	OffsetY    float64              `json:"offsetY,omitempty"`
}

// ChartTickSpanModel is one tspan of a custom tick.
type ChartTickSpanModel struct {
	Text     string  `json:"text"`
	Class    string  `json:"class,omitempty"`
	FontSize float64 `json:"fontSize,omitempty"`
	Dy       string  `json:"dy,omitempty"`
	ResetX   bool    `json:"resetX,omitempty"`
}

// ChartPolarModel is the grid and angle-axis setup of a radar or radial
// chart — the pendant of templui's PolarModel. Ticks is always empty in
// this task (see ChartTickContentModel's own doc comment); every other
// field is wired.
type ChartPolarModel struct {
	GridType     string                  `json:"gridType,omitempty"`
	RadialLines  bool                    `json:"radialLines"`
	PolarRadius  []float64               `json:"polarRadius,omitempty"`
	Stroke       string                  `json:"stroke,omitempty"`
	StrokeWidth  float64                 `json:"strokeWidth,omitempty"`
	GridClass    string                  `json:"gridClass,omitempty"`
	Ticks        []ChartTickContentModel `json:"ticks,omitempty"`
	HasGrid      bool                    `json:"hasGrid,omitempty"`
	HasAngleAxis bool                    `json:"hasAngleAxis,omitempty"`
}

// ChartRadialModel is the geometry of a RadialBarChart. Center (the middle
// Label) is not ported — the Label child itself is out of this task's
// explicit child list.
type ChartRadialModel struct {
	StartAngle  float64 `json:"startAngle"`
	EndAngle    float64 `json:"endAngle"`
	InnerRadius float64 `json:"innerRadius,omitempty"`
	OuterRadius float64 `json:"outerRadius,omitempty"`
}

// ChartPieLabelModel is the resolved label option of a Pie.
type ChartPieLabelModel struct {
	Fill string `json:"fill,omitempty"`
}

// ChartPieModel is one rendered Pie with its slices — the pendant of
// templui's PieModel, minus LabelList/ActiveIndex/ActiveShape/Center: the
// LabelList and Label children they depend on are out of this task's
// explicit child list, and ActiveIndex/ActiveShape are deferred the same
// way chartBarOptions' own ActiveIndex/ActiveBar are (see chartPieOptions'
// own doc comment).
type ChartPieModel struct {
	Key          string              `json:"key"`
	SeriesLabel  string              `json:"seriesLabel,omitempty"`
	Values       []float64           `json:"values"`
	Labels       []string            `json:"labels"`
	TooltipNames []string            `json:"tooltipNames,omitempty"`
	NameKey      string              `json:"nameKey,omitempty"`
	NameValues   []string            `json:"nameValues,omitempty"`
	Colors       []string            `json:"colors"`
	InnerRadius  float64             `json:"innerRadius,omitempty"`
	OuterRadius  float64             `json:"outerRadius,omitempty"`
	StrokeWidth  float64             `json:"strokeWidth,omitempty"`
	Stroke       string              `json:"stroke,omitempty"`
	Label        *ChartPieLabelModel `json:"label,omitempty"`
	LabelLine    bool                `json:"labelLine,omitempty"`
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
// client renderer consumes, the pendant of templui's buildModel: pie,
// radar and radial-bar each get their own builder below (matching
// upstream's own early-return dispatch), the rest of this function stays
// the cartesian builder (bar, line, area).
func buildChartModel(ctx context.Context, st *chartState) ChartModel {
	config := chartConfigFromCtx(ctx)
	if st.kind == "pie" {
		return buildChartPieModel(ctx, config, st)
	}
	if st.kind == "radar" {
		return buildChartRadarModel(ctx, config, st)
	}
	if st.kind == "radial" {
		return buildChartRadialModel(ctx, config, st)
	}

	m := ChartModel{Kind: st.kind, StackOffset: st.stackOffset}

	if g := st.grid; g != nil {
		m.Grid = true
		m.GridHorizontal = g.Horizontal
		m.GridVertical = g.Vertical
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
		m.Cursor = tt.opts.Cursor
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
			// A fill column names the row's own color, like Recharts reading
			// the fill off the entry — the cartesian pendant of
			// buildChartRadialModel's own fallback (chartHasFillColumn,
			// above). Upstream's Bar also reads an explicit <Cell>-child
			// form first (chart.templ:1237); that form stays unported here,
			// matching this file's own Bar scope-narrowing elsewhere.
			if chartHasFillColumn(st.data) {
				s.Cells = make([]string, len(st.data))
				for i, row := range st.data {
					s.Cells[i] = chartStr(row["fill"])
				}
			}
			m.Series = append(m.Series, s)
		}
		m.Stacked = stacked
		m.DomainMin = chartSeriesMin(m.Series)
	case "line":
		for _, l := range st.lines {
			s := chartModelSeries(config, l.key, "", 0, st.data)
			s.Curve = l.opts.Curve
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
			s.Curve = a.opts.Curve
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

	// applyIcons: copy the config icons into the model series. Upstream
	// only calls this from the cartesian buildModel — buildRadarModel,
	// buildRadialModel and buildPieModel below never do, so icons stay
	// cartesian-only, byte-faithful to that omission.
	for i := range m.Series {
		for _, sc := range config {
			if sc.Key == m.Series[i].Key && sc.Icon != nil {
				m.Series[i].Icon = chartRenderHTML(ctx, sc.Icon)
			}
		}
	}

	return m
}

// chartPolarModel resolves a registered ChartPolarGrid into the model the
// renderer consumes — the pendant of templui's polarModel.
func chartPolarModel(g *chartPolarGridOptions) ChartPolarModel {
	return ChartPolarModel{
		HasGrid:     true,
		GridType:    g.GridType,
		RadialLines: g.RadialLines,
		PolarRadius: g.PolarRadius,
		Stroke:      g.Stroke,
		StrokeWidth: g.StrokeWidth,
		GridClass:   g.Class,
	}
}

// buildChartRadarModel normalizes a RadarChart into the model, the pendant
// of templui's buildRadarModel.
func buildChartRadarModel(ctx context.Context, config ChartConfig, st *chartState) ChartModel {
	m := ChartModel{Kind: "radar"}
	if st.margin != nil {
		m.MarginTop, m.MarginRight, m.MarginBottom, m.MarginLeft = st.margin.Top, st.margin.Right, st.margin.Bottom, st.margin.Left
	} else {
		m.MarginTop, m.MarginRight, m.MarginBottom, m.MarginLeft = 5, 5, 5, 5
	}
	polar := ChartPolarModel{RadialLines: true}
	if g := st.polarGrid; g != nil {
		polar = chartPolarModel(g)
	}
	m.Labels = make([]string, len(st.data))
	m.TooltipLabels = make([]string, len(st.data))
	if a := st.angleAxis; a != nil {
		polar.HasAngleAxis = true
		for i, d := range st.data {
			m.Labels[i] = chartStr(d[a.key])
			m.TooltipLabels[i] = m.Labels[i]
		}
	}
	m.Polar = &polar
	for _, r := range st.radars {
		s := chartModelSeries(config, r.key, r.opts.Fill, 0, st.data)
		s.FillOpacityPtr = r.opts.FillOpacity
		s.Stroke = r.opts.Stroke
		s.StrokeWidth = r.opts.StrokeWidth
		if d := r.opts.Dot; d != nil {
			s.Dot = &ChartDotModel{R: d.R, Fill: d.Fill, FillOpacity: d.FillOpacity}
		}
		m.Series = append(m.Series, s)
	}
	if st.legend != nil {
		m.LegendHeight = defaultChartLegendHeight
	}
	if tt := st.tooltip; tt != nil {
		m.Cursor = tt.opts.Cursor
		m.HasTooltip = true
		m.Tooltip = ChartTooltipModel{
			Indicator: tt.opts.Indicator, HideLabel: tt.opts.HideLabel, HideIndicator: tt.opts.HideIndicator,
			Width: tt.opts.Class, LabelClassName: tt.opts.LabelClass, Color: tt.opts.Color,
		}
	}
	return m
}

// buildChartRadialModel normalizes a RadialBarChart into the model: the
// polar geometry, one series per radial bar and no center label (the
// Label child that would carry one is out of this task's scope) — the
// pendant of templui's buildRadialModel.
func buildChartRadialModel(ctx context.Context, config ChartConfig, st *chartState) ChartModel {
	m := ChartModel{Kind: "radial"}
	if st.margin != nil {
		m.MarginTop, m.MarginRight, m.MarginBottom, m.MarginLeft = st.margin.Top, st.margin.Right, st.margin.Bottom, st.margin.Left
	} else {
		m.MarginTop, m.MarginRight, m.MarginBottom, m.MarginLeft = 5, 5, 5, 5
	}
	// The RadialBarChart defaults: a full turn from zero, no inner radius
	// and an outer radius of 80%.
	r := ChartRadialModel{
		StartAngle:  st.startAngle,
		EndAngle:    360,
		InnerRadius: st.innerRadius,
		OuterRadius: st.outerRadius,
	}
	if st.endAngle != nil {
		r.EndAngle = *st.endAngle
	}
	if g := st.polarGrid; g != nil {
		polar := chartPolarModel(g)
		m.Polar = &polar
	}
	m.Labels = make([]string, len(st.data))
	for i := range st.data {
		// The radius axis has no data key, so its domain is the row
		// index, which is what a tooltip label would show.
		m.Labels[i] = chartStr(i)
	}
	tt := st.tooltip
	for _, rb := range st.radialBars {
		s := chartModelSeries(config, rb.key, rb.opts.Fill, 0, st.data)
		s.Background = rb.opts.Background
		s.CornerRadius = rb.opts.CornerRadius
		s.StackID = rb.opts.StackID
		s.Class = rb.opts.Class
		// A fill column names the row's own color, like Recharts reading
		// the fill off the entry.
		if chartHasFillColumn(st.data) {
			s.Cells = make([]string, len(st.data))
			for i, row := range st.data {
				s.Cells[i] = chartStr(row["fill"])
			}
		}
		// getPayloadConfigFromPayload: the tooltip names every row through
		// the name key, which the rows answer with their own value.
		nameKey := rb.key
		if tt != nil && tt.opts.NameKey != "" {
			nameKey = tt.opts.NameKey
		}
		s.TooltipNames = make([]string, len(st.data))
		for i, row := range st.data {
			s.TooltipNames[i] = config.label(chartPayloadConfigKey(row, nameKey))
		}
		m.Series = append(m.Series, s)
	}
	m.Radial = &r
	if tt != nil {
		m.Cursor = tt.opts.Cursor
		m.HasTooltip = true
		m.Tooltip = ChartTooltipModel{
			Indicator: tt.opts.Indicator, HideLabel: tt.opts.HideLabel, HideIndicator: tt.opts.HideIndicator,
			Width: tt.opts.Class, LabelClassName: tt.opts.LabelClass, Color: tt.opts.Color,
		}
	}
	return m
}

// buildChartPieModel normalizes a PieChart into the model: one ChartPieModel
// per registered ChartPie — the pendant of templui's buildPieModel. Unlike
// every other kind, it never touches the top-level margin, Labels or
// Series fields (upstream doesn't either), so those marshal as their zero
// value (0 / null / null) for a pie chart.
func buildChartPieModel(ctx context.Context, config ChartConfig, st *chartState) ChartModel {
	m := ChartModel{Kind: "pie", Cursor: false}
	if st.legend != nil {
		m.LegendHeight = defaultChartLegendHeight
	}
	for _, ps := range st.pies {
		pm := ChartPieModel{
			Key:         ps.key,
			SeriesLabel: config.label(ps.key),
			NameKey:     ps.opts.NameKey,
			InnerRadius: ps.opts.InnerRadius,
			OuterRadius: ps.opts.OuterRadius,
			StrokeWidth: ps.opts.StrokeWidth,
			Stroke:      ps.opts.Stroke,
			LabelLine:   ps.opts.LabelLine,
		}
		if ps.opts.Label != nil {
			pm.Label = &ChartPieLabelModel{Fill: ps.opts.Label.Fill}
		}
		pm.Values = make([]float64, len(ps.data))
		pm.Labels = make([]string, len(ps.data))
		pm.TooltipNames = make([]string, len(ps.data))
		pm.NameValues = make([]string, len(ps.data))
		pm.Colors = make([]string, len(ps.data))
		for i, d := range ps.data {
			pm.Values[i] = chartNum(d[ps.key])
			// getTooltipNameProp: the name key value names the slice, the
			// data key stands in when the row has none.
			key := ps.key
			if ps.opts.NameKey != "" {
				if v, ok := d[ps.opts.NameKey]; ok && v != nil {
					key = chartStr(v)
				}
			}
			pm.NameValues[i] = key
			pm.Labels[i] = config.label(key)
			// ChartTooltipContent names the row through nameKey when it is
			// set, otherwise through the slice name.
			nameKey := key
			if st.tooltip != nil && st.tooltip.opts.NameKey != "" {
				nameKey = st.tooltip.opts.NameKey
			}
			pm.TooltipNames[i] = config.label(chartPayloadConfigKey(d, nameKey))
			if f, ok := d["fill"]; ok {
				pm.Colors[i] = chartStr(f)
			} else {
				pm.Colors[i] = chartSeriesColor(key)
			}
		}
		m.Pies = append(m.Pies, pm)
	}
	if tt := st.tooltip; tt != nil {
		m.Cursor = tt.opts.Cursor
		m.HasTooltip = true
		m.Tooltip = ChartTooltipModel{
			Indicator: tt.opts.Indicator, HideLabel: tt.opts.HideLabel, HideIndicator: tt.opts.HideIndicator,
			Width: tt.opts.Class, LabelClassName: tt.opts.LabelClass, Color: tt.opts.Color,
		}
		if tt.opts.LabelKey != "" {
			m.Tooltip.Label = config.label(tt.opts.LabelKey)
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
	color gsx.RawCSS // the series' vouched color (or a pie slice's model color)
	icon  gsx.Node   // rendered in place; never a pre-rendered string
	value string     // the payload value the legend sorts by
}

// chartBuildLegendItems builds the legend payload: one entry per series, or
// one per pie slice, then sorted by value like Recharts' itemSorter
// default of "value" — the pendant of templui's legendItems. A pie without
// a name key gives every entry the same value, so the data order survives
// the stable sort, exactly like upstream.
func chartBuildLegendItems(config ChartConfig, m ChartModel, st *chartState, opts *chartLegendOptions) []chartLegendItem {
	var items []chartLegendItem
	if m.Kind == "pie" {
		for pi, pie := range m.Pies {
			rows := st.pies[pi].data
			for i := range pie.Labels {
				// The label comes from the legend's name key on the data
				// row, like getPayloadConfigFromPayload does.
				name := pie.Labels[i]
				if opts.NameKey != "" && i < len(rows) {
					if v, ok := rows[i][opts.NameKey]; ok && v != nil {
						name = config.label(chartStr(v))
					}
				}
				items = append(items, chartLegendItem{label: name, color: gsx.RawCSS(pie.Colors[i]), value: pie.NameValues[i]})
			}
		}
	} else {
		items = make([]chartLegendItem, 0, len(m.Series))
		for _, s := range m.Series {
			label := s.Label
			if opts.NameKey != "" {
				label = config.label(opts.NameKey)
			}
			items = append(items, chartLegendItem{
				label: label,
				color: gsx.RawCSS(s.Color),
				icon:  config.icon(s.Key),
				value: s.Key,
			})
		}
	}
	sort.SliceStable(items, func(i, j int) bool { return items[i].value < items[j].value })
	return items
}

// ChartLegendContent renders the ChartLegendContent pendant, absolutely
// positioned like Recharts' legend wrapper — the pendant of templui's
// legendContent. Called by chartRoot once a chart's children have
// registered a ChartLegend, not meant to be composed directly by a
// consumer. ChartTooltipTemplate (below chartRoot) is the other part of
// this tree that renders real markup from a non-root registration, the
// same reasoning.
component ChartLegendContent(items []chartLegendItem, opts *chartLegendOptions) {
	{{ top := opts.VerticalAlign == "top" }}
	<div
		style={ css`position:absolute;left:0;right:0`, css`top:5px`: top, css`bottom:5px`: !top }
		data-gsxui-slot-chart-legend
	>
		<div class={ "flex items-center justify-center gap-4", "pt-3": !top, "pb-3": top, opts.Class }>
			{ for _, it := range items {
				<div class="flex items-center gap-1.5 [&>svg]:h-3 [&>svg]:w-3 [&>svg]:text-muted-foreground">
					{ if it.icon != nil && !opts.HideIcon {
						{ it.icon }
					} else {
						<div class="h-2 w-2 shrink-0 rounded-[2px]" style={css`background-color:@{it.color}`}></div>
					} }
					{ it.label }
				</div>
			} }
		</div>
	</div>
}

// chartRoot renders one chart root: templui's AreaChart/BarChart/LineChart/
// RadarChart/RadialBarChart/PieChart wrapper div, its children (which
// register into the chartState seeded here), then the model script and (if
// registered) the server-rendered legend — templui's own chartOutput, run
// after children the same way gsx renders a Node's children before the
// statements that follow { children } in its caller. st arrives already
// populated with its kind and kind-specific fields (margin, data,
// RadialBarChart's own angle/radius geometry, ...); every root below
// builds one and hands it here rather than chartRoot taking every kind's
// fields as its own positional parameters.
func chartRoot(st *chartState, children gsx.Node, attrs gsx.Attrs) gsx.Node {
	return gsx.Func(func(ctx context.Context, w io.Writer) error {
		ctx = context.WithValue(ctx, chartStateCtxKey, st)
		gw := gsx.W(w)
		gw.S("<div")
		gw.StyleMerged("position:relative;width:100%;height:100%", attrs.Style())
		gw.Spread(ctx, "div", attrs, gsx.AttrSinks{}, []string{"style"})
		gw.S(">")
		gw.Node(ctx, children)
		m := buildChartModel(ctx, st)
		gw.S(chartModelScript(m))
		if st.tooltip != nil {
			gw.Node(ctx, ChartTooltipTemplate())
		}
		if st.legend != nil {
			items := chartBuildLegendItems(chartConfigFromCtx(ctx), m, st, st.legend)
			gw.Node(ctx, ChartLegendContent(items, st.legend))
		}
		gw.S("</div>")
		return gw.Err()
	})
}

// ChartTooltipTemplate renders the hover tooltip's server-authored chrome:
// an inert <template> ui/chart.render.js discovers once per panel and
// reads from (never renders itself) the first time a hover needs to show
// a tooltip. Called by chartRoot once a chart's children have registered a
// ChartTooltip, the same pattern ChartLegendContent uses for the legend.
//
// Its shell div carries chart's per-style tooltip recipe class (the
// tooltip slot's own accessor call, below) plus the
// data-gsxui-slot-chart-tooltip marker that ends up on the LIVE tooltip
// too — chart.render.js reads the
// shell's own class attribute once per panel and reuses it verbatim when
// it builds the live tooltip's opening tag, rather than typing the
// compiled class string as a JS literal (this codebase's audit gate bans
// every JS-side way of stamping a class attribute, unconditionally: no
// classList, no className assignment, no setAttribute("class", ...) — see
// ui/chart.render.js's own tooltipWrapper comment). The four indicator
// swatches carry the exact utility literals the tooltip's dot/line/dashed
// indicator needs (border-(--color-border), bg-(--color-bg), h-2.5,
// w-2.5, my-0.5, border-[1.5px], w-0 — templui's own indicatorHTML, ported
// byte-faithful in chart.render.js) as real class="..." markup here
// instead of as JS string literals: Task 5's report found these specific
// classes nowhere in web/site.css's compiled output, because they only
// ever existed as literals inside chart.render.js, and web/site.css's
// @source globs scan .gsx files only, never ui/*.js. chart.render.js reads
// each variant's class the same way it reads the shell's.
//
// The data-gsxui-chart-tooltip-part elements below carry the rest of
// tooltipHTML's own markup classes (the row wrapper and its SVG-child
// sizing utilities, the label, the name/value grid, the value's own
// tabular-nums styling): an earlier draft of this template stopped at the
// four indicator swatches and left those as JS template literals in
// chart.render.js, which broke the same way the indicator swatches would
// have — a class that exists nowhere in a .gsx file never compiles, so a
// ChartSeries.Icon rendered in a tooltip row painted at its raw SVG
// default size instead of the row's intended size. chart.render.js's
// tooltipHTML/indicatorHTML read every one of these parts' class
// attributes off this template exactly like the four indicator swatches,
// composing row/value-wrap variants from the shared base part plus a
// modifier part (items-center/items-end) rather than typing either as a
// literal.
component ChartTooltipTemplate() {
	<template data-gsxui-chart-tooltip-template>
		<div
			class={
				"bg-popover text-popover-foreground ring-foreground/5 dark:ring-foreground/10 gap-1.5 rounded-xl px-2.5 py-1.5 text-xs shadow-lg ring-1 grid min-w-32 items-start"
			}
			data-gsxui-slot-chart-tooltip
		></div>
		<div
			class="shrink-0 rounded-[2px] border-(--color-border) bg-(--color-bg) h-2.5 w-2.5"
			data-gsxui-chart-tooltip-indicator="dot"
		></div>
		<div
			class="shrink-0 rounded-[2px] border-(--color-border) bg-(--color-bg) w-1"
			data-gsxui-chart-tooltip-indicator="line"
		></div>
		<div
			class="shrink-0 rounded-[2px] border-(--color-border) bg-(--color-bg) w-0 border-[1.5px] border-dashed bg-transparent"
			data-gsxui-chart-tooltip-indicator="dashed"
		></div>
		<div
			class="shrink-0 rounded-[2px] border-(--color-border) bg-(--color-bg) w-0 border-[1.5px] border-dashed bg-transparent my-0.5"
			data-gsxui-chart-tooltip-indicator="dashed-nested"
		></div>
		<div class="font-medium" data-gsxui-chart-tooltip-part="label"></div>
		<div
			class="flex w-full flex-wrap items-stretch gap-2 [&>svg]:h-2.5 [&>svg]:w-2.5 [&>svg]:text-muted-foreground"
			data-gsxui-chart-tooltip-part="row"
		></div>
		<div class="items-center" data-gsxui-chart-tooltip-part="items-center"></div>
		<div class="items-end" data-gsxui-chart-tooltip-part="items-end"></div>
		<div class="flex flex-1 justify-between leading-none" data-gsxui-chart-tooltip-part="value-wrap"></div>
		<div class="grid gap-1.5" data-gsxui-chart-tooltip-part="grid"></div>
		<span class="text-muted-foreground" data-gsxui-chart-tooltip-part="name"></span>
		<span
			class="font-mono font-medium text-foreground tabular-nums"
			data-gsxui-chart-tooltip-part="value"
		></span>
	</template>
}

// BarChart is the Recharts BarChart root: its parts declare axes, grid,
// tooltip and bars, the root collects them and emits the model payload for
// the client renderer. marginTop/marginRight/marginBottom/marginLeft
// replace Recharts' own 5/5/5/5 default AS A SET the moment any one of the
// four is customized, matching Recharts (chartMarginOf) — a bare
// `<ui.BarChart data={data}>` gets the 5/5/5/5 default; `<ui.AreaChart
// marginLeft={12} marginRight={12}>` gets marginTop/marginBottom 0, not 5.
component BarChart(data []ChartDatum, marginTop float64, marginRight float64, marginBottom float64, marginLeft float64, children gsx.Node, attrs gsx.Attrs) {
	{ chartRoot(&chartState{kind: "bar", data: data, margin: chartMarginOf(marginTop, marginRight, marginBottom, marginLeft)}, children, attrs) }
}

// LineChart is the Recharts LineChart root.
component LineChart(data []ChartDatum, marginTop float64, marginRight float64, marginBottom float64, marginLeft float64, children gsx.Node, attrs gsx.Attrs) {
	{ chartRoot(&chartState{kind: "line", data: data, margin: chartMarginOf(marginTop, marginRight, marginBottom, marginLeft)}, children, attrs) }
}

// AreaChart is the Recharts AreaChart root. stackOffset "expand" normalizes
// each stack to 100%, the pendant of Recharts' own stackOffset prop.
component AreaChart(data []ChartDatum, marginTop float64, marginRight float64, marginBottom float64, marginLeft float64, stackOffset string, children gsx.Node, attrs gsx.Attrs) {
	{ chartRoot(&chartState{kind: "area", data: data, margin: chartMarginOf(marginTop, marginRight, marginBottom, marginLeft), stackOffset: stackOffset}, children, attrs) }
}

// RadarChart is the Recharts RadarChart root.
component RadarChart(data []ChartDatum, marginTop float64, marginRight float64, marginBottom float64, marginLeft float64, children gsx.Node, attrs gsx.Attrs) {
	{ chartRoot(&chartState{kind: "radar", data: data, margin: chartMarginOf(marginTop, marginRight, marginBottom, marginLeft)}, children, attrs) }
}

// RadialBarChart is the Recharts RadialBarChart root. startAngle is
// Recharts' own zero by default, so the flat param needs no sentinel;
// endAngle defaults to a full turn (360°) when left at zero, via
// chartFloatPtrOr's zero-sentinel convention (a literal 0° sweep — an
// invisible chart — is not expressible this way, a deliberately dropped
// edge case, see this task's report). Margins follow BarChart's own
// all-or-nothing chartMarginOf convention (see there).
component RadialBarChart(data []ChartDatum, marginTop float64, marginRight float64, marginBottom float64, marginLeft float64, startAngle float64, endAngle float64, innerRadius float64, outerRadius float64, children gsx.Node, attrs gsx.Attrs) {
	{ chartRoot(&chartState{
		kind:        "radial",
		data:        data,
		margin:      chartMarginOf(marginTop, marginRight, marginBottom, marginLeft),
		startAngle:  startAngle,
		endAngle:    chartFloatPtrOr(endAngle),
		innerRadius: innerRadius,
		outerRadius: outerRadius,
	}, children, attrs) }
}

// PieChart is the Recharts PieChart root: unlike every other root, it
// carries no data or margin of its own — Recharts' own Pie takes its own
// data prop, so each ChartPie child registers its own rows and geometry
// (see ChartPie); upstream's own PieChartProps is empty too.
component PieChart(children gsx.Node, attrs gsx.Attrs) {
	{ chartRoot(&chartState{kind: "pie"}, children, attrs) }
}

// ChartCartesianGrid registers the grid, the pendant of Recharts'
// CartesianGrid element. It renders nothing; the client draws it.
//
// horizontal/vertical default false — gsxui's own convention (`## chart`
// ADAPT, docs/jsx-parity.md), diverging from Recharts' own default-true for
// both: shadcn's own demos pass an explicit false for whichever direction
// they don't want (nearly always vertical) rather than relying on the
// Recharts default, so spelling both out explicitly here — e.g.
// `<ui.ChartCartesianGrid horizontal/>` — reproduces the exact same
// shadcn demo look a bare tag would otherwise miss.
component ChartCartesianGrid(horizontal bool, vertical bool) {
	{{
		if st := chartStateFromCtx(ctx); st != nil {
			st.grid = &chartGridOptions{Horizontal: horizontal, Vertical: vertical}
		}
	}}
}

// ChartXAxis registers the x axis.
component ChartXAxis(key string, hide bool, tickLine bool, axisLine bool, tickMargin float64, minTickGap float64) {
	{{
		if st := chartStateFromCtx(ctx); st != nil {
			st.x = &chartAxisReg{key: key, opts: chartXAxisOptions{
				Hide: hide, TickLine: tickLine, AxisLine: axisLine,
				TickMargin: tickMargin, MinTickGap: minTickGap,
			}}
		}
	}}
}

// ChartYAxis registers the y axis.
component ChartYAxis(key string, hide bool, tickLine bool, axisLine bool, tickMargin float64, tickCount int, width float64) {
	{{
		if st := chartStateFromCtx(ctx); st != nil {
			st.y = &chartYAxisReg{key: key, opts: chartYAxisOptions{
				Hide: hide, TickLine: tickLine, AxisLine: axisLine,
				TickMargin: tickMargin, TickCount: tickCount, Width: width,
			}}
		}
	}}
}

// ChartTooltip registers the tooltip's config; the enclosing chart root
// renders the tooltip's hidden server-side chrome after its other children
// (see ChartTooltipTemplate), and the client fills it with the hovered
// row's values from the model (m.tooltip, the ChartTooltipModel this
// registration resolves into).
//
// cursor defaults false — gsxui's own convention (`## chart` ADAPT),
// diverging from Recharts' own default-true; shadcn's own demos pass an
// explicit `cursor={false}` on nearly every tooltip, so
// `<ui.ChartTooltip/>` reproduces that shipped look.
//
// hasDefaultIndex/defaultIndex split what upstream carries as one nilable
// number (`defaultIndex?: number`, shows the tooltip pinned at that
// category on mount): a flat `defaultIndex int` param alone cannot tell
// "not set" apart from the legitimate index 0 (Go's zero value for int),
// the same reason ChartPie's Label field became a `label bool` gate below
// rather than a bare fill string. No shipped demo sets this.
component ChartTooltip(cursor bool, indicator string, labelKey string, hideLabel bool, hideIndicator bool, nameKey string, class string, labelClass string, color string, hasDefaultIndex bool, defaultIndex int) {
	{{
		var di *int
		if hasDefaultIndex {
			di = &defaultIndex
		}
		if st := chartStateFromCtx(ctx); st != nil {
			st.tooltip = &chartTooltipReg{opts: chartTooltipOptions{
				Cursor: cursor, Indicator: indicator, LabelKey: labelKey,
				HideLabel: hideLabel, HideIndicator: hideIndicator, NameKey: nameKey,
				Class: class, LabelClass: labelClass, Color: color, DefaultIndex: di,
			}}
		}
	}}
}

// ChartLegend registers the legend; the enclosing chart root renders it
// server-side after its other children (see ChartLegendContent).
component ChartLegend(nameKey string, class string, verticalAlign string, hideIcon bool) {
	{{
		if st := chartStateFromCtx(ctx); st != nil {
			st.legend = &chartLegendOptions{
				NameKey: nameKey, Class: class, VerticalAlign: verticalAlign, HideIcon: hideIcon,
			}
		}
	}}
}

// ChartDefs groups gradient definitions like the svg defs element —
// nothing renders visibly server-side (the client draws the SVG), so this
// only needs to render its children to let any ChartLinearGradient among
// them register.
component ChartDefs(children gsx.Node) {
	{ children }
}

// ChartLinearGradient registers one gradient definition. Its children are
// real SVG `<stop>` markup (not gsx.Raw strings, `## chart` element-
// structure decision): rendered here to a string and captured into the
// registered definition verbatim, like Recharts passes defs children
// through.
component ChartLinearGradient(id string, x1 string, y1 string, x2 string, y2 string, children gsx.Node) {
	{{
		if st := chartStateFromCtx(ctx); st != nil {
			var sb strings.Builder
			if children != nil {
				if err := children.Render(ctx, &sb); err != nil {
					return err
				}
			}
			st.defs = append(st.defs, chartGradientReg{
				id: id, x1: x1, y1: y1, x2: x2, y2: y2,
				stops: strings.TrimSpace(sb.String()),
			})
		}
	}}
}

// ChartArea registers one area series. Declaration order is paint order,
// like in Recharts. curve is Recharts' curve prop ("natural"/"linear"/
// "step"); empty renders as Recharts' own "linear" default, unchanged by
// this task.
component ChartArea(key string, curve string, fill string, stroke string, stackId string, fillOpacity float64) {
	{{
		if st := chartStateFromCtx(ctx); st != nil {
			st.areas = append(st.areas, chartAreaReg{key: key, opts: chartAreaOptions{
				Curve: curve, Fill: fill, Stroke: stroke, StackID: stackId, FillOpacity: fillOpacity,
			}})
		}
	}}
}

// ChartBar registers one bar series. radius is Recharts' radius union: a
// float64 for all corners or a []float64 of four corners, e.g.
// []float64{0, 0, 4, 4}.
component ChartBar(key string, fill string, stackId string, radius any, strokeWidth float64) {
	{{
		if st := chartStateFromCtx(ctx); st != nil {
			st.bars = append(st.bars, chartBarReg{key: key, opts: chartBarOptions{
				Fill: fill, StackID: stackId, Radius: radius, StrokeWidth: strokeWidth,
			}})
		}
	}}
}

// ChartLine registers one line series. curve is Recharts' curve prop
// ("natural"/"linear"/"step"); empty renders as Recharts' own "linear"
// default, unchanged by this task.
//
// dot/dotR/dotFillOpacity/dotFill/dotSize flatten Recharts' dot prop's
// object form: dot defaults false, matching every shipped demo's own
// dot={false} and this task's default-false boolean convention — the
// per-point dots draw only once a caller opts in with dot, at which point
// dotR/dotFillOpacity/dotFill/dotSize shape them (each its own Recharts
// SVG default when left at zero).
component ChartLine(key string, curve string, stroke string, strokeWidth float64, dot bool, dotR float64, dotFillOpacity float64, dotFill string, dotSize float64, activeDotR float64) {
	{{
		var d *chartDotOptions
		if dot {
			d = &chartDotOptions{R: dotR, FillOpacity: dotFillOpacity, Fill: dotFill, Size: dotSize}
		}
		if st := chartStateFromCtx(ctx); st != nil {
			st.lines = append(st.lines, chartLineReg{key: key, opts: chartLineOptions{
				Curve: curve, Stroke: stroke, StrokeWidth: strokeWidth, Dot: d, ActiveDotR: activeDotR,
			}})
		}
	}}
}

// ---------------------------------------------------------------------
// Polar model builder (pie, radar, radial-bar) — adapted from templui's
// chart.templ lines 729-1100 (RadarChart / RadialBarChart / PieChart
// roots, PolarGrid / PolarAngleAxis / PolarRadiusAxis / Radar / RadialBar /
// Pie, buildRadarModel / buildRadialModel / buildPieModel / polarModel,
// legendItems' pie branch). Every part below renders nothing itself,
// exactly like the cartesian parts above: children register into the
// chartState the enclosing root seeded through ctx.
//
// Label (the center donut/radial label) and its Pie/PolarRadiusAxis mount
// points, LabelList, Cell, and every Recharts render-prop closure
// (PolarAngleAxis's Tick) are not ported — none is in this task's explicit
// child list, and a closure cannot cross the JSON boundary to the client
// renderer regardless. Pie's ActiveIndex/ActiveShape are deferred the same
// way chartBarOptions' own ActiveIndex/ActiveBar are (see chartPieOptions'
// own doc comment).

// ChartPolarGrid registers the polar grid, the pendant of Recharts'
// PolarGrid element. It renders nothing; the client draws it.
//
// radialLines defaults false — gsxui's own convention (`## chart` ADAPT),
// diverging from Recharts' own default-true; a radar demo that wants the
// spokes upstream shows by default passes `<ui.ChartPolarGrid
// radialLines/>` explicitly.
component ChartPolarGrid(gridType string, radialLines bool, polarRadius []float64, stroke string, strokeWidth float64, class string) {
	{{
		if st := chartStateFromCtx(ctx); st != nil {
			st.polarGrid = &chartPolarGridOptions{
				GridType: gridType, RadialLines: radialLines, PolarRadius: polarRadius,
				Stroke: stroke, StrokeWidth: strokeWidth, Class: class,
			}
		}
	}}
}

// ChartPolarAngleAxis registers the angle axis, the labels around the
// chart.
component ChartPolarAngleAxis(key string) {
	{{
		if st := chartStateFromCtx(ctx); st != nil {
			st.angleAxis = &chartPolarAngleAxisReg{key: key}
		}
	}}
}

// ChartPolarRadiusAxis registers the radius axis. tick/tickLine/axisLine
// are currently inert (see chartPolarRadiusAxisOptions' own doc comment);
// children is still rendered, matching upstream's own { children... }, so
// a future Label child has somewhere to register through once it is
// ported.
component ChartPolarRadiusAxis(tick bool, tickLine bool, axisLine bool, children gsx.Node) {
	{{
		if st := chartStateFromCtx(ctx); st != nil {
			st.radiusAxis = &chartPolarRadiusAxisOptions{Tick: tick, TickLine: tickLine, AxisLine: axisLine}
		}
	}}
	{ children }
}

// ChartRadar registers one radar series.
//
// dot/dotR/dotFill/dotFillOpacity flatten Recharts' dot prop's object
// form, the same split ChartLine's own dot params use — dot defaults
// false, matching this task's default-false boolean convention (no shipped
// radar demo sets a dot today either way).
//
// fillOpacity 0 leaves the SVG default of 1 in place, via
// chartFloatPtrOr's zero-sentinel convention (chartRadarOptions.FillOpacity
// stays a pointer internally for exactly that reason) — an explicit
// fillOpacity=0 (fully transparent fill) is the one edge case this cannot
// express, the same deliberately dropped case chartAreaOptions.FillOpacity
// already accepted for Area before this task.
component ChartRadar(key string, fill string, fillOpacity float64, stroke string, strokeWidth float64, dot bool, dotR float64, dotFill string, dotFillOpacity float64) {
	{{
		var d *chartDotOptions
		if dot {
			d = &chartDotOptions{R: dotR, Fill: dotFill, FillOpacity: dotFillOpacity}
		}
		if st := chartStateFromCtx(ctx); st != nil {
			st.radars = append(st.radars, chartRadarReg{key: key, opts: chartRadarOptions{
				Fill: fill, FillOpacity: chartFloatPtrOr(fillOpacity), Stroke: stroke, StrokeWidth: strokeWidth, Dot: d,
			}})
		}
	}}
}

// ChartRadialBar registers one radial bar series. background draws the
// track behind the bar, over the full angle range.
component ChartRadialBar(key string, fill string, background bool, cornerRadius float64, stackId string, class string) {
	{{
		if st := chartStateFromCtx(ctx); st != nil {
			st.radialBars = append(st.radialBars, chartRadialBarReg{key: key, opts: chartRadialBarOptions{
				Fill: fill, Background: background, CornerRadius: cornerRadius, StackID: stackId, Class: class,
			}})
		}
	}}
}

// ChartPie registers one pie with its own data and geometry — unlike the
// cartesian series, Recharts' own Pie takes its own data prop rather than
// reading the enclosing chart's rows (PieChart itself carries none).
//
// label/labelFill split what upstream carries as one nilable object
// (Recharts' own `label` prop: falsy hides it, an object shape configures
// it) — label defaults false, matching every shipped demo (none passes
// Recharts' own label prop), and labelFill only has an effect once label
// is set.
//
// labelLine defaults false — gsxui's own convention (`## chart` ADAPT),
// diverging from Recharts' own default-true; inert without label also set
// (chart.render.js only draws either when pie.label is present), so this
// flip changes no shipped demo's rendered pixels.
component ChartPie(data []ChartDatum, key string, nameKey string, innerRadius float64, outerRadius float64, strokeWidth float64, stroke string, label bool, labelFill string, labelLine bool) {
	{{
		var lbl *chartPieLabelOptions
		if label {
			lbl = &chartPieLabelOptions{Fill: labelFill}
		}
		if st := chartStateFromCtx(ctx); st != nil {
			st.pies = append(st.pies, chartPieReg{data: data, key: key, opts: chartPieOptions{
				NameKey: nameKey, InnerRadius: innerRadius, OuterRadius: outerRadius,
				StrokeWidth: strokeWidth, Stroke: stroke, Label: lbl, LabelLine: labelLine,
			}})
		}
	}}
}
