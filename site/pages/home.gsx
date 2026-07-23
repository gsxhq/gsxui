package pages

import (
	"github.com/gsxhq/gsxui/ui/badge"
	"github.com/gsxhq/gsxui/ui/button"
	"github.com/gsxhq/gsxui/ui/dialog"
)

// Home is the site's landing page.
type Home struct{}

const installSnippet = `go install github.com/gsxhq/gsxui/cmd/gsxui@latest
gsxui init
gsxui add button`

component (h Home) Page() {
	<Layout title="gsxui" active="">
		<section class="flex flex-col gap-4 py-10">
			<h1 class="text-4xl font-semibold tracking-tight sm:text-5xl">
				Components for modern web frontends in Go.
			</h1>
			<p class="max-w-2xl text-lg text-muted-foreground">
				gsxui is a shadcn-style component set for
				<a
					href="https://github.com/gsxhq/gsx"
					target="_blank"
					rel="noreferrer"
					class="underline underline-offset-4 hover:text-foreground"
				>
					gsx
				</a>
				: copy-in, type-checked, server-rendered. Own the code, style it with Tailwind, ship it.
			</p>
			<div class="flex flex-wrap items-center gap-3 pt-2">
				<button.Button size="lg" href="#components">Browse components</button.Button>
				<button.Button
					size="lg"
					variant="outline"
					href="https://github.com/gsxhq/gsxui"
					target="_blank"
					rel="noreferrer"
				>
					View on GitHub
				</button.Button>
			</div>
		</section>
		<section class="pb-10">
			<pre class="overflow-x-auto rounded-lg border border-border bg-card p-4 text-sm text-card-foreground"><code>
				{ installSnippet }
			</code></pre>
		</section>
		<section id="components" class="flex flex-col gap-10 border-t border-border py-10">
			<div>
				<h2 class="text-sm font-medium uppercase tracking-wide text-muted-foreground">Button</h2>
				<div class="mt-4 flex flex-wrap items-center gap-3">
					<button.Button>Default</button.Button>
					<button.Button variant="secondary">Secondary</button.Button>
					<button.Button variant="destructive">Destructive</button.Button>
					<button.Button variant="outline">Outline</button.Button>
					<button.Button variant="ghost">Ghost</button.Button>
					<button.Button variant="link">Link</button.Button>
				</div>
			</div>
			<div>
				<h2 class="text-sm font-medium uppercase tracking-wide text-muted-foreground">Badge</h2>
				<div class="mt-4 flex flex-wrap items-center gap-3">
					<badge.Badge>Default</badge.Badge>
					<badge.Badge variant="secondary">Secondary</badge.Badge>
					<badge.Badge variant="destructive">Destructive</badge.Badge>
					<badge.Badge variant="outline">Outline</badge.Badge>
				</div>
			</div>
			<div>
				<h2 class="text-sm font-medium uppercase tracking-wide text-muted-foreground">Dialog</h2>
				<div class="mt-4">
					<dialog.Dialog>
						<dialog.DialogTrigger>
							<button.Button variant="outline">Open dialog</button.Button>
						</dialog.DialogTrigger>
						<dialog.DialogContent>
							<dialog.DialogHeader>
								<dialog.DialogTitle>Edit profile</dialog.DialogTitle>
								<dialog.DialogDescription>
									Rendered by ui/dialog on the native &lt;dialog&gt; element — no client framework required.
								</dialog.DialogDescription>
							</dialog.DialogHeader>
							<dialog.DialogFooter showCloseButton={true}></dialog.DialogFooter>
						</dialog.DialogContent>
					</dialog.Dialog>
				</div>
			</div>
		</section>
	</Layout>
}
