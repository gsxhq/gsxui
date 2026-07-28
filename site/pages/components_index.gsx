package pages

import (
	"github.com/gsxhq/gsxui/internal/registry"
)

// ComponentsIndex is the /components/ catalog page: every component the
// registry ships, as a grid of cards linking to its docs page. The list is
// derived from the registry (same source as the sidebar), so it can never
// drift from what ui/ actually contains. The home hero's "Browse
// components" button lands here — an in-page anchor was its previous
// target, which on desktop viewports fit entirely above the fold and so
// visibly did nothing.
type ComponentsIndex struct{}

component (c ComponentsIndex) Page() {
	<Layout title="Components" active="">
		<section class="flex flex-col gap-6">
			<div>
				<h1 class="text-3xl font-semibold tracking-tight">Components</h1>
				<p class="mt-2 max-w-2xl text-muted-foreground">
					Every component gsxui ships. Copy any of them into your project with
					<code class="rounded bg-muted px-1.5 py-0.5 font-mono text-sm">gsxui add &lt;name&gt;</code>.
				</p>
			</div>
			{{ names, _ := registry.Components() }}
			<div class="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-4">
				{ for _, name := range names {
					<a
						href={"/components/" + name}
						class="rounded-lg border border-border p-4 text-sm font-medium capitalize transition-colors hover:bg-accent hover:text-accent-foreground"
					>
						{ name }
					</a>
				} }
			</div>
		</section>
	</Layout>
}
