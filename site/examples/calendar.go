package examples

import (
	"net/url"
	"time"

	"github.com/gsxhq/gsx"
	examplecalendar "github.com/gsxhq/gsxui/site/examples/calendar"
)

func init() {
	Register("calendar", Example{
		Name:       "basic",
		Title:      "Basic",
		Node:       examplecalendar.Basic(examplecalendar.DefaultMonth),
		SourcePath: "calendar/basic.gsx",
		// Query answers jstest/harness's ?month= override (the Go/JS
		// agreement tests in jstest/specs/calendar.spec.ts need the server
		// to render an arbitrary month, not just DefaultMonth) — this is
		// where "month" gets interpreted; the harness itself only ever
		// forwards url.Values generically (registry.go's Example.Query doc
		// comment). An invalid or missing ?month= falls back to
		// DefaultMonth, same as the static Node above.
		Query: func(q url.Values) gsx.Node {
			month := examplecalendar.DefaultMonth
			if v := q.Get("month"); v != "" {
				if t, err := time.Parse("2006-01", v); err == nil {
					month = t
				}
			}
			return examplecalendar.Basic(month)
		},
	})
	Register("calendar", Example{
		Name:       "bounded",
		Title:      "Bounded",
		Node:       examplecalendar.Bounded(),
		SourcePath: "calendar/bounded.gsx",
		// No Query hook: the bounded-navigation browser tests never re-render
		// the server at a different month, they only click prev/next and
		// assert the client-side bound-enforcement calendar.js now has.
	})
	Register("calendar", Example{
		Name:       "loaded",
		Title:      "Loaded",
		Node:       examplecalendar.Loaded(examplecalendar.LoadedDefaultMonth),
		SourcePath: "calendar/loaded.gsx",
		// Query answers ?month= the same way basic's own hook does — Task 4
		// review's Important finding 2 needs a genuine client-vs-server
		// agreement diff for an example that actually carries a selection
		// and disabled rules, not just basic's empty ones.
		Query: func(q url.Values) gsx.Node {
			month := examplecalendar.LoadedDefaultMonth
			if v := q.Get("month"); v != "" {
				if t, err := time.Parse("2006-01", v); err == nil {
					month = t
				}
			}
			return examplecalendar.Loaded(month)
		},
	})
	Register("calendar", Example{
		Name:       "loadedrange",
		Title:      "Loaded range",
		Node:       examplecalendar.LoadedRange(examplecalendar.LoadedRangeDefaultMonth),
		SourcePath: "calendar/loaded.gsx",
		// Same reason as loaded's own Query hook, for mode="range" instead —
		// Task 4 review's Important finding 2 also named range flags
		// specifically as zero-coverage.
		Query: func(q url.Values) gsx.Node {
			month := examplecalendar.LoadedRangeDefaultMonth
			if v := q.Get("month"); v != "" {
				if t, err := time.Parse("2006-01", v); err == nil {
					month = t
				}
			}
			return examplecalendar.LoadedRange(month)
		},
	})
	Register("calendar", Example{
		Name:       "range",
		Title:      "Range",
		Node:       examplecalendar.Range(examplecalendar.RangeDefaultMonth),
		SourcePath: "calendar/range.gsx",
		// Query mirrors basic's own hook — Task 5's browser tests don't use
		// it directly (they click through the two-click range machine
		// starting from RangeDefaultMonth), but keeping the same ?month=
		// override every other example carries means a later task can add a
		// Go/JS agreement test for range selection without touching this
		// file again.
		Query: func(q url.Values) gsx.Node {
			month := examplecalendar.RangeDefaultMonth
			if v := q.Get("month"); v != "" {
				if t, err := time.Parse("2006-01", v); err == nil {
					month = t
				}
			}
			return examplecalendar.Range(month)
		},
	})
	Register("calendar", Example{
		Name:       "multiple",
		Title:      "Multiple",
		Node:       examplecalendar.Multiple(examplecalendar.MultipleDefaultMonth),
		SourcePath: "calendar/multiple.gsx",
		// Same reason as range's own Query hook above.
		Query: func(q url.Values) gsx.Node {
			month := examplecalendar.MultipleDefaultMonth
			if v := q.Get("month"); v != "" {
				if t, err := time.Parse("2006-01", v); err == nil {
					month = t
				}
			}
			return examplecalendar.Multiple(month)
		},
	})
}
