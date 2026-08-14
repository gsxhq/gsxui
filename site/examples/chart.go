package examples

import examplechart "github.com/gsxhq/gsxui/site/examples/chart"

func init() {
	Register("chart", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       examplechart.Basic(),
		SourcePath: "chart/basic.gsx",
	})
	Register("chart", Example{
		Name:       "area",
		Title:      "Area",
		Node:       examplechart.Area(),
		SourcePath: "chart/area.gsx",
	})
	Register("chart", Example{
		Name:       "line",
		Title:      "Line",
		Node:       examplechart.Line(),
		SourcePath: "chart/line.gsx",
	})
	Register("chart", Example{
		Name:       "pie",
		Title:      "Pie",
		Node:       examplechart.Pie(),
		SourcePath: "chart/pie.gsx",
	})
	Register("chart", Example{
		Name:       "radar",
		Title:      "Radar",
		Node:       examplechart.Radar(),
		SourcePath: "chart/radar.gsx",
	})
	Register("chart", Example{
		Name:       "radial",
		Title:      "Radial",
		Node:       examplechart.Radial(),
		SourcePath: "chart/radial.gsx",
	})
}
