package pages

// GettingStarted is the /docs/getting-started page: install → init → add →
// a minimal first page, expanded from README.md with real CLI output
// (copied from internal/cli/init.go / add.go's actual printed strings, not
// invented).
type GettingStarted struct{}

const gsInstallSnippet = `go install github.com/gsxhq/gsxui/cmd/gsxui@latest`

const gsInitSnippet = `gsxui init`

const gsInitOutput = `gsxui initialized.
  css:  web/gsxui.css (import it from your entry point)
  js:   web/gsxui/index.js (import it from your entry point)
  next: gsxui add button`

const gsAddSnippet = `gsxui add button card`

const gsAddOutput = `adding: button card
done — build with: go build ./...`

const gsPageGsx = `package main

import "example.com/myapp/ui"

component Home() {
	<html lang="en">
		<head>
			<meta charset="UTF-8"/>
			<link rel="stylesheet" href="/web/gsxui.css"/>
		</head>
		<body class="flex min-h-svh items-center justify-center bg-background p-8 text-foreground">
			<ui.Card class="w-full max-w-sm">
				<ui.CardHeader>
					<ui.CardTitle>Hello, gsxui</ui.CardTitle>
					<ui.CardDescription>Your first page.</ui.CardDescription>
				</ui.CardHeader>
				<ui.CardContent>
					<ui.Button>Click me</ui.Button>
				</ui.CardContent>
			</ui.Card>
		</body>
	</html>
}`

const gsMainGo = `package main

import (
	"log"
	"net/http"
)

func main() {
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Home().Render(r.Context(), w)
	})
	log.Println("listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}`

const gsGenerateSnippet = `go tool gsx generate`

const gsRunSnippet = `go run .`

const gsRunOutput = `2026/07/23 09:00:00 listening on http://localhost:8080`

component (g GettingStarted) Page() {
	<Layout title="Getting Started" active="getting-started">
		<div data-doc="getting-started" class="flex max-w-3xl flex-col gap-10 py-10">
			<div class="flex flex-col gap-4">
				<h1 class="text-3xl font-semibold tracking-tight">Getting Started</h1>
				<p class="text-muted-foreground">
					gsxui components are copy-in: the CLI vendors real <code>.gsx</code> source into
					your own module, so what you build against is code you own and can edit —
					not a package you import and can't touch.
				</p>
			</div>

			<section class="flex flex-col gap-3">
				<h2 class="text-sm font-medium uppercase tracking-wide text-muted-foreground">1. Install the CLI</h2>
				<pre class="overflow-x-auto rounded-lg border border-border bg-card p-4 text-sm text-card-foreground"><code>{ gsInstallSnippet }</code></pre>
			</section>

			<section class="flex flex-col gap-3">
				<h2 class="text-sm font-medium uppercase tracking-wide text-muted-foreground">2. Initialize your project</h2>
				<p class="text-muted-foreground">In your project (a Go module):</p>
				<pre class="overflow-x-auto rounded-lg border border-border bg-card p-4 text-sm text-card-foreground"><code>{ gsInitSnippet }</code></pre>
				<pre class="overflow-x-auto rounded-lg border border-border bg-card p-4 text-sm text-card-foreground"><code>{ gsInitOutput }</code></pre>
				<p class="text-muted-foreground">
					This vendors the theme tokens (<code>web/gsxui.css</code>), the JS runtime
					and behavior barrel (<code>web/gsxui/</code>), and the class merger
					(<code>ui/merge/merge.go</code>), then points <code>gsx.toml</code>'s
					<code>class_merger</code> at it — the seam that makes caller-class-merge
					work (see <a href={ Theming{} |> url } class="underline underline-offset-4 hover:text-foreground">Theming</a>). It also
					<code>go get</code>s <code>gsx</code> and <code>tailwind-merge-go</code>, and
					installs the <code>gsx</code> tool via <code>go get -tool</code>.
				</p>
			</section>

			<section class="flex flex-col gap-3">
				<h2 class="text-sm font-medium uppercase tracking-wide text-muted-foreground">3. Add components</h2>
				<pre class="overflow-x-auto rounded-lg border border-border bg-card p-4 text-sm text-card-foreground"><code>{ gsAddSnippet }</code></pre>
				<pre class="overflow-x-auto rounded-lg border border-border bg-card p-4 text-sm text-card-foreground"><code>{ gsAddOutput }</code></pre>
				<p class="text-muted-foreground">
					<code>card</code> has no dependencies of its own, but a component that does
					(e.g. <code>select</code>, which needs <code>icon</code>) pulls its
					dependency in automatically — <code>gsxui add select</code> vendors
					<code>icon</code> too. You own every file this writes:
					<code>gsxui add</code> never touches one you've already modified unless you
					pass <code>--overwrite</code>. After upgrading the <code>gsxui</code> binary,
					re-run <code>gsxui add &lt;name&gt; --overwrite</code> to refresh vendored
					components — that discards local edits to those files.
				</p>
			</section>

			<section class="flex flex-col gap-3">
				<h2 class="text-sm font-medium uppercase tracking-wide text-muted-foreground">4. Your first page</h2>
				<p class="text-muted-foreground">
					A tiny two-file app: <code>home.gsx</code> renders a <code>Card</code> around a
					<code>Button</code>, and <code>main.go</code> serves it.
				</p>
				<pre class="overflow-x-auto rounded-lg border border-border bg-card p-4 text-sm text-card-foreground"><code>{ gsPageGsx }</code></pre>
				<pre class="overflow-x-auto rounded-lg border border-border bg-card p-4 text-sm text-card-foreground"><code>{ gsMainGo }</code></pre>
				<p class="text-muted-foreground">Compile the <code>.gsx</code> file to plain Go, then run it:</p>
				<pre class="overflow-x-auto rounded-lg border border-border bg-card p-4 text-sm text-card-foreground"><code>{ gsGenerateSnippet }</code></pre>
				<p class="text-muted-foreground">
					(silent on success — it writes <code>home.x.go</code> next to
					<code>home.gsx</code> and exits 0)
				</p>
				<pre class="overflow-x-auto rounded-lg border border-border bg-card p-4 text-sm text-card-foreground"><code>{ gsRunSnippet }</code></pre>
				<pre class="overflow-x-auto rounded-lg border border-border bg-card p-4 text-sm text-card-foreground"><code>{ gsRunOutput }</code></pre>
				<p class="text-muted-foreground">
					Open <code>http://localhost:8080</code> — a styled Card with a Button inside,
					rendered with gsxui's default light theme. Next:
					<a href={ Theming{} |> url } class="underline underline-offset-4 hover:text-foreground">restyle it</a>.
				</p>
			</section>
		</div>
	</Layout>
}
