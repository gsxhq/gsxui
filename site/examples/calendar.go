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
}
