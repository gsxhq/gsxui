// Package calendar holds the site's example gsx components for ui/calendar.
package calendar

import (
	"time"

	"github.com/gsxhq/gsxui/ui"
)

// DefaultMonth is the fixed month Basic renders when no month is given —
// 2026-01, deliberately never time.Now(): jstest/specs/calendar.spec.ts
// asserts exact grid contents (data-date values, the "Go and JS agree on …"
// diffs), and a test like that cannot depend on the day it runs.
var DefaultMonth = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

// Basic renders a single-mode calendar pinned to DefaultMonth by default. A
// zero month falls back to DefaultMonth, so the zero-arg registration
// (site/examples/calendar.go's static Node) and jstest/harness's own
// per-request re-render (its Query hook, driven by ?month=) share one
// component instead of two.
component Basic(month time.Time) {
	{{
		if month.IsZero() {
			month = DefaultMonth
		}
	}}
	<ui.Calendar mode="single" month={month} weekStartsOn={time.Sunday} showOutsideDays={true} captionLayout="label"/>
}
