package pages

import (
	"net/http"

	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/site/examples"
)

// ExamplePreview is the minimal same-origin document used for application-
// shell examples whose viewport-owned layout cannot render inline in the
// documentation page.
type ExamplePreview struct{}

type examplePreviewProps struct {
	Title string
	Node  gsx.Node
}

func (ExamplePreview) Props(r *http.Request) (examplePreviewProps, error) {
	component := r.PathValue("component")
	name := r.PathValue("example")
	title, node, ok := examples.Find(component, name, r.URL.Query())
	if !ok {
		return examplePreviewProps{}, ErrorWithStatus{
			Status:  http.StatusNotFound,
			Message: "unknown example: " + component + "/" + name,
		}
	}
	return examplePreviewProps{
		Title: component + " · " + title,
		Node:  node,
	}, nil
}

component (preview ExamplePreview) Page(props examplePreviewProps) {
	<!DOCTYPE html>
	<html lang="en">
		<siteHead title={props.Title}/>
		<body
			data-site-isolated-document
			class="h-svh overflow-auto bg-background text-foreground antialiased"
		>
			{ props.Node }
		</body>
	</html>
}
