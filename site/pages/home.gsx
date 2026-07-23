package pages

import (
	"github.com/gsxhq/gsxui/ui"
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
					href="https://gsxhq.github.io/"
					target="_blank"
					rel="noreferrer"
					class="underline underline-offset-4 hover:text-foreground"
				>
					gsx
				</a>
				: copy-in, type-checked, server-rendered. Own the code, style it with Tailwind, ship it.
			</p>
			<div class="flex flex-wrap items-center gap-3 pt-2">
				<ui.Button size="lg" href="#components">Browse components</ui.Button>
				<ui.Button
					size="lg"
					variant="outline"
					href="https://github.com/gsxhq/gsxui"
					target="_blank"
					rel="noreferrer"
				>
					View on GitHub
				</ui.Button>
				<ui.Button
					size="lg"
					variant="ghost"
					href="https://gsxhq.github.io/"
					target="_blank"
					rel="noreferrer"
				>
					gsx documentation ↗
				</ui.Button>
			</div>
		</section>
		<section class="pb-10">
			<pre
				class="overflow-x-auto rounded-lg border border-border bg-card p-4 text-sm text-card-foreground"
			><code>{ installSnippet }</code></pre>
		</section>
		<section id="components" class="flex flex-col gap-10 border-t border-border py-10">
			<div>
				<h2 class="text-sm font-medium uppercase tracking-wide text-muted-foreground">Button</h2>
				<div class="mt-4 flex flex-wrap items-center gap-3">
					<ui.Button>Default</ui.Button>
					<ui.Button variant="secondary">Secondary</ui.Button>
					<ui.Button variant="destructive">Destructive</ui.Button>
					<ui.Button variant="outline">Outline</ui.Button>
					<ui.Button variant="ghost">Ghost</ui.Button>
					<ui.Button variant="link">Link</ui.Button>
				</div>
			</div>
			<div>
				<h2 class="text-sm font-medium uppercase tracking-wide text-muted-foreground">Badge</h2>
				<div class="mt-4 flex flex-wrap items-center gap-3">
					<ui.Badge>Default</ui.Badge>
					<ui.Badge variant="secondary">Secondary</ui.Badge>
					<ui.Badge variant="destructive">Destructive</ui.Badge>
					<ui.Badge variant="outline">Outline</ui.Badge>
				</div>
			</div>
			<div>
				<h2 class="text-sm font-medium uppercase tracking-wide text-muted-foreground">Dialog</h2>
				<div class="mt-4">
					<ui.Dialog>
						<ui.Button variant="outline" data-gsxui-dialog-trigger>Open dialog</ui.Button>
						<ui.DialogContent>
							<ui.DialogHeader>
								<ui.DialogTitle>Edit profile</ui.DialogTitle>
								<ui.DialogDescription>
									Rendered by ui/dialog on the native &lt;dialog&gt; element — no client framework required.
								</ui.DialogDescription>
							</ui.DialogHeader>
							<ui.DialogFooter showCloseButton={true}></ui.DialogFooter>
						</ui.DialogContent>
					</ui.Dialog>
				</div>
			</div>
		</section>
	</Layout>
}
