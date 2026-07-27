package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"

	"github.com/gsxhq/gsxui/ui"
)

// shellTmpl is the minimal page every harness route renders into. It
// deliberately carries only what a component needs: the compiled stylesheet
// and one module script. No site chrome, no web/site.js, no theme script —
// those are site code, not library code, and loading them would put
// untested JS in the way of the JS under test.
//
// The body classes match site/pages/layout.gsx so theme tokens resolve the
// same way they do in production.
var shellTmpl = template.Must(template.New("shell").Parse(
	`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{{.Title}}</title>
<link rel="stylesheet" href="{{.Stylesheet}}">
<script type="module" src="{{.Script}}"></script>
</head>
<body class="min-h-svh bg-background text-foreground antialiased">
<main data-harness-root class="p-8">{{.Body}}</main>
{{.Toaster}}
</body>
</html>
`))

type shellData struct {
	Title      string
	Stylesheet string
	Script     string
	Body       template.HTML
	Toaster    template.HTML
}

// renderShell writes the shell around already-rendered markup. body is
// trusted: it comes from a gsx component's own Render, which escapes its
// own interpolations.
func renderShell(w io.Writer, title, stylesheet, script string, body template.HTML) error {
	var toaster bytes.Buffer
	if err := ui.Toaster(nil).Render(context.Background(), &toaster); err != nil {
		return fmt.Errorf("rendering toaster: %w", err)
	}
	return shellTmpl.Execute(w, shellData{
		Title:      title,
		Stylesheet: stylesheet,
		Script:     script,
		Body:       body,
		Toaster:    template.HTML(toaster.String()),
	})
}
