package pages

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/site/examples"
	"github.com/gsxhq/gsxui/site/hl"
)

// Component is the /components/{name} page — the harness that proves
// source-shown-is-source-run: every example is a real gsx component
// (registered in site/examples/registry.go), rendered live next to the
// exact source text that produced it. Unknown/unregistered names 404 (see
// Props below and ErrorWithStatus in pages.go).
type Component struct{}

// ComponentProps is Component's Props result.
type ComponentProps struct {
	Name     string
	Title    string
	Examples []exampleProps
}

// exampleProps pairs a registered examples.Example with the key its
// highlighted source is stored under. The source text itself is no longer
// loaded here: site/hl holds every example pre-rendered to highlighted HTML
// (generated from these same files, see site/hl/gen), so the page looks the
// block up by SourcePath instead of reading and escaping it per request.
type exampleProps struct {
	Name       string
	Title      string
	Node       gsx.Node
	SourcePath string
	Isolated   bool
	Previews   []examples.Preview
}

// Props resolves the {name} path param against the examples registry.
// Unregistered names (including real ui/ components Task 3 hasn't wired
// examples for yet) 404 via ErrorWithStatus rather than rendering an empty
// page.
func (Component) Props(r *http.Request) (ComponentProps, error) {
	name := r.PathValue("name")
	exs := examples.For(name)
	if len(exs) == 0 {
		return ComponentProps{}, ErrorWithStatus{
			Status:  http.StatusNotFound,
			Message: "unknown component: " + name,
		}
	}
	eps := make([]exampleProps, len(exs))
	for i, ex := range exs {
		eps[i] = exampleProps{
			Name:       ex.Name,
			Title:      ex.Title,
			Node:       ex.Node,
			SourcePath: ex.SourcePath,
			Isolated:   ex.Isolated,
			Previews:   ex.Previews,
		}
	}
	return ComponentProps{Name: name, Title: capitalize(name), Examples: eps}, nil
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// shadcnSlug maps gsxui component names to the slug shadcn/ui uses under
// ui.shadcn.com/docs/components/{slug} — most names match verbatim, but a
// couple of gsxui components restructure their shadcn source under a
// different name (dropdown ports DropdownMenu; radio ports RadioGroup).
// Names not present here pass through unchanged.
var shadcnSlug = map[string]string{
	"dropdown": "dropdown-menu",
	"radio":    "radio-group",
}

// shadcnName resolves a gsxui component name to its shadcn/ui docs slug,
// passing the name through unchanged when no rename is on record.
func shadcnName(name string) string {
	if slug, ok := shadcnSlug[name]; ok {
		return slug
	}
	return name
}

func examplePreviewURL(componentName string, exampleName string, preview string) string {
	path := "/examples/" + url.PathEscape(componentName) + "/" + url.PathEscape(exampleName)
	if preview == "" {
		return path
	}
	return path + "?" + url.Values{
		examples.PreviewQueryKey: []string{preview},
	}.Encode()
}

func componentTOCItems(examples []exampleProps) []docTOCItem {
	items := make([]docTOCItem, len(examples))
	for i, example := range examples {
		items[i] = docTOCItem{
			ID:    "example-" + example.Name,
			Title: example.Title,
			Depth: 2,
		}
	}
	return items
}

component (c Component) Page(props ComponentProps) {
	{{ toc := componentTOCItems(props.Examples) }}
	<siteLayout title={props.Title} active={props.Name} mode={layoutDocs} toc={toc}>
		<div class="flex flex-col gap-10 py-10">
			<h1 class="text-3xl font-semibold tracking-tight">{ props.Title }</h1>
			{ for i, ex := range props.Examples {
				<section class="flex flex-col gap-3">
					<docHeading item={toc[i]} class="text-sm font-medium uppercase tracking-wide text-muted-foreground"/>
					{ if ex.Isolated && len(ex.Previews) == 0 {
						<iframe
							data-site-isolated-preview
							title={ex.Title + " preview"}
							src={examplePreviewURL(props.Name, ex.Name, "")}
							loading="lazy"
							class="block h-[32rem] w-full rounded-lg border bg-background"
						></iframe>
					} else if ex.Isolated {
						<div class="flex flex-col gap-6">
							{ for _, preview := range ex.Previews {
								<div class="flex flex-col gap-2">
									<div class="text-sm font-medium">{ preview.Title }</div>
									<iframe
										data-site-isolated-preview
										title={preview.Title + " preview"}
										src={examplePreviewURL(props.Name, ex.Name, preview.Name)}
										loading="lazy"
										class="block h-80 w-full rounded-lg border bg-background"
									></iframe>
								</div>
							} }
						</div>
					} else {
						<div class="border rounded-lg p-8 bg-background">
							{ ex.Node }
						</div>
					} }
					<div class="relative" data-site-example>
						<pre
							class="overflow-x-auto rounded-2xl bg-muted/50 px-4 py-3.5 font-mono text-sm"
						><code>{ hl.Node(ex.SourcePath) }</code></pre>
						<button
							type="button"
							data-site-copy
							class="absolute right-2 top-2 rounded-md border border-border bg-background px-2 py-1 text-xs text-muted-foreground transition-colors hover:bg-accent hover:text-accent-foreground"
						>
							Copy
						</button>
					</div>
				</section>
			} }
			<footer class="flex flex-col gap-3 border-t border-border pt-6 text-sm text-muted-foreground">
				<pre
					class="overflow-x-auto rounded-lg border border-border bg-card p-4 text-card-foreground"
				><code>{ "gsxui add " + props.Name }</code></pre>
				{ if props.Name == "icon" {
					<a
						href="https://lucide.dev"
						target="_blank"
						rel="noreferrer"
						class="underline underline-offset-4 hover:text-foreground"
					>
						View the icon set on lucide.dev
					</a>
				} else {
					<a
						href={"https://ui.shadcn.com/docs/components/" + shadcnName(props.Name)}
						target="_blank"
						rel="noreferrer"
						class="underline underline-offset-4 hover:text-foreground"
					>
						View the original on shadcn/ui
					</a>
				} }
			</footer>
		</div>
	</siteLayout>
}
