package pages

import (
	"net/http"
	"strings"

	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/site/examples"
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

// exampleProps pairs a registered examples.Example with its loaded source
// text — loaded once in Props (which can error) so Page itself never has
// to handle an error mid-render.
type exampleProps struct {
	Title  string
	Node   gsx.Node
	Source string
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
		src, err := examples.Source(name, ex.Name)
		if err != nil {
			return ComponentProps{}, err
		}
		eps[i] = exampleProps{Title: ex.Title, Node: ex.Node, Source: src}
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
// handful of gsxui components were renamed off their shadcn source (e.g.
// switchctl, since "switch" is a reserved Go keyword — see
// ui/switchctl/switch.gsx) or restructured (selectbox splits shadcn's
// Select; dropdown ports DropdownMenu). Names not present here pass through
// unchanged.
var shadcnSlug = map[string]string{
	"switchctl": "switch",
	"selectbox": "select",
	"dropdown":  "dropdown-menu",
	"radio":     "radio-group",
}

// shadcnName resolves a gsxui component name to its shadcn/ui docs slug,
// passing the name through unchanged when no rename is on record.
func shadcnName(name string) string {
	if slug, ok := shadcnSlug[name]; ok {
		return slug
	}
	return name
}

component (c Component) Page(props ComponentProps) {
	<Layout title={ props.Title } active={ props.Name }>
		<div class="flex flex-col gap-10 py-10">
			<h1 class="text-3xl font-semibold tracking-tight">{ props.Title }</h1>
			{ for _, ex := range props.Examples {
				<section class="flex flex-col gap-3">
					<h2 class="text-sm font-medium uppercase tracking-wide text-muted-foreground">{ ex.Title }</h2>
					<div class="border rounded-lg p-8 bg-background">
						{ ex.Node }
					</div>
					<div class="relative" data-site-example>
						<pre class="overflow-x-auto rounded-lg border border-border bg-card p-4 text-sm text-card-foreground"><code>{ ex.Source }</code></pre>
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
				<pre class="overflow-x-auto rounded-lg border border-border bg-card p-4 text-card-foreground"><code>{ "gsxui add " + props.Name }</code></pre>
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
						href={ "https://ui.shadcn.com/docs/components/" + shadcnName(props.Name) }
						target="_blank"
						rel="noreferrer"
						class="underline underline-offset-4 hover:text-foreground"
					>
						View the original on shadcn/ui
					</a>
				} }
			</footer>
		</div>
	</Layout>
}
