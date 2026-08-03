# Home Page Component Showcase Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the landing page's flat Button/Badge/Dialog section with a 2×2 bento grid of four composed demo cards (sign-in, settings, stats, interactive overlays) built from existing ui components.

**Architecture:** New Go/gsx package `site/examples/showcase/` exposing four exported components (`SignInCard`, `SettingsCard`, `StatsCard`, `OverlaysCard`), each a `ui.Card` composition. `site/pages/home.gsx` imports the package and renders the grid in place of the current `#components` section. No changes to `ui/` or the registry.

**Tech Stack:** gsx (Go server-rendered components), Tailwind classes, Go stdlib tests, Playwright (jstest harness).

## Global Constraints

- Spec: `docs/superpowers/specs/2026-08-03-home-showcase-design.md`.
- Comments inside gsx markup MUST be `{/* ... */}` — never `//` (that's Go-code-position only).
- All element `id=` / `for=` values on these cards are prefixed `home-showcase-` (the home page shares a document with nav/search; ids must be unique page-wide).
- After creating or editing any `.gsx` file, run `make generate` (runs `go tool gsx generate`) to produce/update the sibling `.x.go` files. Commit the generated `.x.go` alongside the `.gsx`.
- Never `git add -A`; add files explicitly.
- Do NOT register these components in `site/examples/registry.go` — they are landing-page compositions, not per-component docs examples. (The `//go:embed */*.gsx` in `site/examples/examples.go` will embed the new sources; that is harmless and requires no action.)
- Playwright MUST run as `npx playwright test --config jstest/playwright.config.ts` — without `--config` no harness starts and every test fails with "Cannot navigate to invalid URL".
- Final verification gates (Task 6): `go build ./...`, `go test ./... -count=1` (allow 600s), `make audit`, `make verify-generated`, `gofmt -l .` (expect empty), Playwright as above.

---

### Task 1: showcase package + SignInCard

**Files:**
- Create: `site/examples/showcase/signin.gsx`
- Test: `site/examples/showcase/showcase_test.go`

**Interfaces:**
- Produces: `component SignInCard()` in package `showcase` (`github.com/gsxhq/gsxui/site/examples/showcase`). Zero-arg, renders one `ui.Card`.

- [ ] **Step 1: Write the failing test**

Create `site/examples/showcase/showcase_test.go`:

```go
package showcase_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/gsxhq/gsxui/site/examples/showcase"
)

// render executes a showcase card and returns its HTML.
func render(t *testing.T, node interface {
	Render(ctx context.Context, w *bytes.Buffer) error
}) string {
	t.Helper()
	var buf bytes.Buffer
	if err := node.Render(context.Background(), &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}

func TestSignInCard(t *testing.T) {
	html := render(t, showcase.SignInCard())
	for _, want := range []string{
		"home-showcase-email",
		"home-showcase-password",
		"home-showcase-remember",
		"Sign in",
	} {
		if !strings.Contains(html, want) {
			t.Errorf("SignInCard output missing %q", want)
		}
	}
}
```

Note: if `node.Render`'s writer parameter is `io.Writer` in the generated code (check a neighbor like `site/examples/card/basic.x.go` for the exact signature), use `io.Writer` in the `render` helper instead of `*bytes.Buffer`. Match what the generated `.x.go` actually exposes.

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./site/examples/showcase/ -run TestSignInCard -count=1`
Expected: FAIL — package doesn't exist yet / `showcase.SignInCard` undefined.

- [ ] **Step 3: Write the component**

Create `site/examples/showcase/signin.gsx`:

```go
// Package showcase holds the composed demo cards rendered on the site's
// landing page. Unlike the per-component packages next door, these are
// not registered docs examples — they exist to show several ui components
// working together in one realistic piece of UI.
package showcase

import "github.com/gsxhq/gsxui/ui"

// SignInCard composes Card, Label, Input, Checkbox and Button into the
// classic email/password sign-in form.
component SignInCard() {
	<ui.Card>
		<ui.CardHeader>
			<ui.CardTitle>Sign in</ui.CardTitle>
			<ui.CardDescription>Enter your email below to sign in to your account.</ui.CardDescription>
		</ui.CardHeader>
		<ui.CardContent>
			<form>
				<div class="flex flex-col gap-6">
					<div class="grid gap-2">
						<ui.Label for="home-showcase-email">Email</ui.Label>
						<ui.Input id="home-showcase-email" type="email" placeholder="m@example.com" required/>
					</div>
					<div class="grid gap-2">
						<div class="flex items-center">
							<ui.Label for="home-showcase-password">Password</ui.Label>
							<a href="#" class="ms-auto inline-block text-sm underline-offset-4 hover:underline">Forgot password?</a>
						</div>
						<ui.Input id="home-showcase-password" type="password" required/>
					</div>
					<div class="flex items-center gap-2">
						<ui.Checkbox id="home-showcase-remember"/>
						<ui.Label for="home-showcase-remember">Remember me</ui.Label>
					</div>
				</div>
			</form>
		</ui.CardContent>
		<ui.CardFooter class="flex-col gap-2">
			<ui.Button class="w-full">Sign in</ui.Button>
			<ui.Button variant="ghost" class="w-full">Create an account</ui.Button>
		</ui.CardFooter>
	</ui.Card>
}
```

- [ ] **Step 4: Generate and run the test**

Run: `make generate && go test ./site/examples/showcase/ -run TestSignInCard -count=1`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add site/examples/showcase/signin.gsx site/examples/showcase/signin.x.go site/examples/showcase/showcase_test.go
git commit -m "feat(site): add showcase package with SignInCard"
```

---

### Task 2: SettingsCard

**Files:**
- Create: `site/examples/showcase/settings.gsx`
- Modify (append): `site/examples/showcase/showcase_test.go`

**Interfaces:**
- Consumes: test helper `render(t, node)` from Task 1's `showcase_test.go`.
- Produces: `component SettingsCard()` in package `showcase`.

- [ ] **Step 1: Write the failing test**

Append to `site/examples/showcase/showcase_test.go`:

```go
func TestSettingsCard(t *testing.T) {
	html := render(t, showcase.SettingsCard())
	for _, want := range []string{
		"home-showcase-notifications",
		"home-showcase-autosave",
		"home-showcase-theme",
		"home-showcase-density",
	} {
		if !strings.Contains(html, want) {
			t.Errorf("SettingsCard output missing %q", want)
		}
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./site/examples/showcase/ -run TestSettingsCard -count=1`
Expected: FAIL — `showcase.SettingsCard` undefined.

- [ ] **Step 3: Write the component**

Create `site/examples/showcase/settings.gsx`:

```go
package showcase

import "github.com/gsxhq/gsxui/ui"

// SettingsCard composes Switch, NativeSelect, Slider and Separator into a
// small preferences panel.
component SettingsCard() {
	<ui.Card>
		<ui.CardHeader>
			<ui.CardTitle>Preferences</ui.CardTitle>
			<ui.CardDescription>Manage how the workspace behaves.</ui.CardDescription>
		</ui.CardHeader>
		<ui.CardContent class="flex flex-col gap-4">
			<div class="flex items-center justify-between gap-4">
				<ui.Label for="home-showcase-notifications">Notifications</ui.Label>
				<ui.Switch id="home-showcase-notifications" checked/>
			</div>
			<div class="flex items-center justify-between gap-4">
				<ui.Label for="home-showcase-autosave">Autosave</ui.Label>
				<ui.Switch id="home-showcase-autosave"/>
			</div>
			<ui.Separator/>
			<div class="grid gap-2">
				<ui.Label for="home-showcase-theme">Theme</ui.Label>
				<ui.NativeSelect id="home-showcase-theme" name="home-showcase-theme">
					<ui.NativeSelectOption value="system" selected>System</ui.NativeSelectOption>
					<ui.NativeSelectOption value="light">Light</ui.NativeSelectOption>
					<ui.NativeSelectOption value="dark">Dark</ui.NativeSelectOption>
				</ui.NativeSelect>
			</div>
			<ui.Separator/>
			<div class="grid gap-2">
				<ui.Label for="home-showcase-density">Density</ui.Label>
				<ui.Slider id="home-showcase-density" value={40} min={0} max={100} step={10} aria-label="Density"/>
			</div>
		</ui.CardContent>
	</ui.Card>
}
```

Note: if `make generate` rejects `id=` on `ui.Slider` or `ui.NativeSelect` (attrs spread differences), drop the `id` from the component and the `for` from its Label, keeping the `aria-label`/`name` — then update the test's expected strings to match (`name="home-showcase-theme"` and `aria-label="Density"` remain greppable).

- [ ] **Step 4: Generate and run the test**

Run: `make generate && go test ./site/examples/showcase/ -run TestSettingsCard -count=1`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add site/examples/showcase/settings.gsx site/examples/showcase/settings.x.go site/examples/showcase/showcase_test.go
git commit -m "feat(site): add SettingsCard to showcase"
```

---

### Task 3: StatsCard

**Files:**
- Create: `site/examples/showcase/stats.gsx`
- Modify (append): `site/examples/showcase/showcase_test.go`

**Interfaces:**
- Consumes: test helper `render(t, node)` from Task 1.
- Produces: `component StatsCard()` in package `showcase`.

- [ ] **Step 1: Write the failing test**

Append to `showcase_test.go`:

```go
func TestStatsCard(t *testing.T) {
	html := render(t, showcase.StatsCard())
	for _, want := range []string{
		"Storage",
		"Bandwidth",
		"Ada Lovelace",
		"Active",
	} {
		if !strings.Contains(html, want) {
			t.Errorf("StatsCard output missing %q", want)
		}
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./site/examples/showcase/ -run TestStatsCard -count=1`
Expected: FAIL — `showcase.StatsCard` undefined.

- [ ] **Step 3: Write the component**

Create `site/examples/showcase/stats.gsx`. The avatar uses the same inline-SVG → `dataURL` pattern as `site/examples/avatar/basic.gsx` (a data: URI never hits the network, keeping the site's clean-console invariant):

```go
package showcase

import "github.com/gsxhq/gsxui/ui"

// Stand-in portrait for the usage card, same inline-SVG-to-data-URL
// technique as site/examples/avatar.
var showcaseAvatarSVG = []byte("<svg xmlns='http://www.w3.org/2000/svg' width='64' height='64'><rect width='64' height='64' fill='#6d28d9'/><text x='32' y='34' text-anchor='middle' dominant-baseline='central' font-family='sans-serif' font-weight='600' font-size='26' fill='#fff'>AL</text></svg>")

// StatsCard composes Avatar, Badge and Progress into a usage summary.
component StatsCard() {
	<ui.Card>
		<ui.CardHeader>
			<ui.CardTitle>Usage</ui.CardTitle>
			<ui.CardDescription>Your plan resets in 12 days.</ui.CardDescription>
			<ui.CardAction>
				<ui.Badge variant="secondary">Active</ui.Badge>
			</ui.CardAction>
		</ui.CardHeader>
		<ui.CardContent class="flex flex-col gap-4">
			<div class="flex items-center gap-3">
				<ui.Avatar>
					<ui.AvatarImage src={showcaseAvatarSVG |> dataURL("image/svg+xml")} alt="Ada Lovelace"/>
					<ui.AvatarFallback>AL</ui.AvatarFallback>
				</ui.Avatar>
				<div class="flex flex-col">
					<span class="text-sm font-medium">Ada Lovelace</span>
					<span class="text-sm text-muted-foreground">Pro plan</span>
				</div>
			</div>
			<div class="grid gap-2">
				<div class="flex items-center justify-between text-sm">
					<span>Storage</span>
					<span class="text-muted-foreground">72%</span>
				</div>
				<ui.Progress value={72}/>
			</div>
			<div class="grid gap-2">
				<div class="flex items-center justify-between text-sm">
					<span>Bandwidth</span>
					<span class="text-muted-foreground">31%</span>
				</div>
				<ui.Progress value={31}/>
			</div>
		</ui.CardContent>
	</ui.Card>
}
```

- [ ] **Step 4: Generate and run the test**

Run: `make generate && go test ./site/examples/showcase/ -run TestStatsCard -count=1`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add site/examples/showcase/stats.gsx site/examples/showcase/stats.x.go site/examples/showcase/showcase_test.go
git commit -m "feat(site): add StatsCard to showcase"
```

---

### Task 4: OverlaysCard

**Files:**
- Create: `site/examples/showcase/overlays.gsx`
- Modify (append): `site/examples/showcase/showcase_test.go`

**Interfaces:**
- Consumes: test helper `render(t, node)` from Task 1.
- Produces: `component OverlaysCard()` in package `showcase`. Contains exactly ONE `data-gsxui-slot-dialog-trigger` (home.gsx's current demo dialog moves here in Task 5, keeping the page's dialog-trigger count at 2 for `validateDialogTriggerSlotMarkers`).

- [ ] **Step 1: Write the failing test**

Append to `showcase_test.go`:

```go
func TestOverlaysCard(t *testing.T) {
	html := render(t, showcase.OverlaysCard())
	for _, want := range []string{
		"data-gsxui-slot-dialog-trigger",
		"data-gsxui-slot-dropdown-menu-trigger",
		"data-gsxui-slot-tooltip-trigger",
		"home-showcase-toast-btn",
	} {
		if !strings.Contains(html, want) {
			t.Errorf("OverlaysCard output missing %q", want)
		}
	}
	if got := strings.Count(html, "data-gsxui-slot-dialog-trigger"); got != 1 {
		t.Errorf("OverlaysCard has %d dialog triggers, want exactly 1", got)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./site/examples/showcase/ -run TestOverlaysCard -count=1`
Expected: FAIL — `showcase.OverlaysCard` undefined.

- [ ] **Step 3: Write the component**

Create `site/examples/showcase/overlays.gsx`. Trigger-slot attribute patterns come verbatim from `site/examples/dropdown/basic.gsx`, `site/examples/tooltip/basic.gsx`, the toast template pattern from `site/examples/toast/server.gsx`, and the dialog markup from the current `site/pages/home.gsx:85-103`:

```go
package showcase

import "github.com/gsxhq/gsxui/ui"

// OverlaysCard proves the interactive components work server-rendered:
// Tabs switch between an overlays pane (Dialog, DropdownMenu, Tooltip)
// and a feedback pane (a server-flash Toast appended into ui.Toaster,
// which siteLayout mounts once per page).
component OverlaysCard() {
	<ui.Card>
		<ui.CardHeader>
			<ui.CardTitle>Interactive, no framework</ui.CardTitle>
			<ui.CardDescription>Dialogs, menus, tooltips and toasts — server-rendered, hydrated by tiny shims.</ui.CardDescription>
		</ui.CardHeader>
		<ui.CardContent>
			<ui.Tabs value="overlays">
				<ui.TabsList>
					<ui.TabsTrigger value="overlays" selected>Overlays</ui.TabsTrigger>
					<ui.TabsTrigger value="feedback">Feedback</ui.TabsTrigger>
				</ui.TabsList>
				<ui.TabsContent value="overlays" selected>
					<div class="flex flex-wrap items-center gap-3 pt-2">
						<ui.Dialog>
							<ui.Button
								variant="outline"
								data-gsxui-slot-dialog-trigger
								aria-haspopup="dialog"
								aria-expanded="false"
							>
								Open dialog
							</ui.Button>
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
						<ui.DropdownMenu>
							<ui.Button
								variant="outline"
								data-gsxui-slot-dropdown-menu-trigger
								aria-haspopup="menu"
								aria-expanded="false"
							>
								Options
							</ui.Button>
							<ui.DropdownMenuContent>
								<ui.DropdownMenuGroup>
									<ui.DropdownMenuLabel>My Account</ui.DropdownMenuLabel>
									<ui.DropdownMenuSeparator/>
									<ui.DropdownMenuItem>Profile</ui.DropdownMenuItem>
									<ui.DropdownMenuItem>Billing</ui.DropdownMenuItem>
									<ui.DropdownMenuItem>
										Settings <ui.DropdownMenuShortcut>⌘,</ui.DropdownMenuShortcut>
									</ui.DropdownMenuItem>
								</ui.DropdownMenuGroup>
							</ui.DropdownMenuContent>
						</ui.DropdownMenu>
						<ui.Tooltip>
							<ui.Button
								variant="outline"
								data-gsxui-slot-tooltip-trigger
							>
								Hover me
							</ui.Button>
							<ui.TooltipContent>Server-rendered tooltip</ui.TooltipContent>
						</ui.Tooltip>
					</div>
				</ui.TabsContent>
				<ui.TabsContent value="feedback">
					<div class="flex flex-col items-start gap-3 pt-2">
						<p class="text-sm text-muted-foreground">Server flash pattern: a pre-rendered toast row appended into the page's toaster, exactly like an HTMX out-of-band swap.</p>
						<ui.Button variant="outline" id="home-showcase-toast-btn">Show toast</ui.Button>
						<template data-home-showcase-toast>
							<ui.Toast
								toastType="success"
								title="Saved"
								description="A server-rendered toast, adopted by ui/toaster."
							/>
						</template>
						<script>
							document.getElementById("home-showcase-toast-btn").addEventListener("click", () => {
								const tpl = document.querySelector("template[data-home-showcase-toast]");
								const viewport = document.getElementById("gsxui-toaster");
								if (!tpl || !viewport) return;
								viewport.appendChild(tpl.content.firstElementChild.cloneNode(true));
							});
						</script>
					</div>
				</ui.TabsContent>
			</ui.Tabs>
		</ui.CardContent>
	</ui.Card>
}
```

- [ ] **Step 4: Generate and run the test**

Run: `make generate && go test ./site/examples/showcase/ -count=1`
Expected: PASS (all four card tests)

- [ ] **Step 5: Commit**

```bash
git add site/examples/showcase/overlays.gsx site/examples/showcase/overlays.x.go site/examples/showcase/showcase_test.go
git commit -m "feat(site): add OverlaysCard to showcase"
```

---

### Task 5: Wire the grid into home.gsx

**Files:**
- Modify: `site/pages/home.gsx:61-106` (the `#components` section)
- Modify: `site/pages/pages_test.go` (extend `TestSiteRoutes`)

**Interfaces:**
- Consumes: `showcase.SignInCard()`, `showcase.SettingsCard()`, `showcase.StatsCard()`, `showcase.OverlaysCard()` from `github.com/gsxhq/gsxui/site/examples/showcase`.

- [ ] **Step 1: Write the failing test**

In `site/pages/pages_test.go`, inside `TestSiteRoutes` after the existing `validateDialogTriggerSlotMarkers` check, add:

```go
	for _, want := range []string{
		"Built with gsxui",
		"home-showcase-email",         // SignInCard
		"home-showcase-notifications", // SettingsCard
		"Ada Lovelace",                // StatsCard
		"home-showcase-toast-btn",     // OverlaysCard
	} {
		if !strings.Contains(rec.Body.String(), want) {
			t.Errorf("home response missing showcase marker %q", want)
		}
	}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./site/pages/ -run TestSiteRoutes -count=1`
Expected: FAIL — showcase markers missing from home response.

- [ ] **Step 3: Rewrite the section**

In `site/pages/home.gsx`:
1. Add `"github.com/gsxhq/gsxui/site/examples/showcase"` to the import block.
2. Replace the entire `#components` section (lines 61–106, from `<section id="components"` through its closing `</section>`) with:

```go
		<section id="components" class="flex flex-col gap-6 border-t border-border py-10">
			<div class="flex flex-wrap items-end justify-between gap-3">
				<h2 class="text-2xl font-semibold tracking-tight">Built with gsxui</h2>
				<a
					href={ComponentsIndex{} |> url}
					class="text-sm text-muted-foreground underline-offset-4 hover:text-foreground hover:underline"
				>
					Browse all components →
				</a>
			</div>
			<div class="grid gap-4 md:grid-cols-2">
				<showcase.SignInCard/>
				<showcase.SettingsCard/>
				<showcase.StatsCard/>
				<showcase.OverlaysCard/>
			</div>
		</section>
```

If gsx requires imported-package component invocation in a different form than `<showcase.SignInCard/>`, follow whatever form `home.gsx` already uses for `<ui.Button>` — dotted package-qualified tags are the established pattern.

- [ ] **Step 4: Generate and run the tests**

Run: `make generate && go test ./site/pages/ -count=1`
Expected: PASS — including the pre-existing `validateDialogTriggerSlotMarkers(document, 2)` (the demo dialog moved into OverlaysCard, so the page still has exactly 2 dialog triggers). If that assertion fails on count, inspect the rendered home body for `data-gsxui-slot-dialog-trigger` occurrences and reconcile — do not just bump the number without understanding which trigger appeared or vanished.

- [ ] **Step 5: Commit**

```bash
git add site/pages/home.gsx site/pages/home.x.go site/pages/pages_test.go
git commit -m "feat(site): replace flat component rows with showcase bento grid"
```

---

### Task 6: Playwright coverage + full gates

**Files:**
- Create: `jstest/specs/home-showcase.spec.ts`

**Interfaces:**
- Consumes: the harness's site server (specs navigate with relative URLs against the configured `baseURL`, e.g. `page.goto("/")` — same pattern as `theme-editor.spec.ts`'s `page.goto("/theme")`).

- [ ] **Step 1: Write the spec**

Create `jstest/specs/home-showcase.spec.ts`. Before writing, skim `jstest/specs/dialog.spec.ts` and `jstest/specs/dropdown.spec.ts` for the harness's import path of `test`/`expect` and any shared setup helpers, and mirror them. The spec body:

```ts
import { test, expect } from "@playwright/test";

test.describe("home showcase", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/");
  });

  test("renders all four showcase cards", async ({ page }) => {
    const section = page.locator("#components");
    await expect(section.getByText("Sign in", { exact: true }).first()).toBeVisible();
    await expect(section.getByText("Preferences")).toBeVisible();
    await expect(section.getByText("Usage")).toBeVisible();
    await expect(section.getByText("Interactive, no framework")).toBeVisible();
  });

  test("dialog opens and closes from the overlays card", async ({ page }) => {
    const trigger = page.locator("#components [data-gsxui-slot-dialog-trigger]");
    await trigger.click();
    const dialog = page.locator("#components dialog[open]");
    await expect(dialog).toBeVisible();
    await expect(dialog.getByText("Edit profile")).toBeVisible();
    await page.keyboard.press("Escape");
    await expect(dialog).toBeHidden();
  });

  test("dropdown menu opens from the overlays card", async ({ page }) => {
    const trigger = page.locator("#components [data-gsxui-slot-dropdown-menu-trigger]");
    await trigger.click();
    await expect(page.getByText("My Account")).toBeVisible();
  });

  test("toast appends into the toaster viewport", async ({ page }) => {
    await page.getByRole("tab", { name: "Feedback" }).click();
    await page.locator("#home-showcase-toast-btn").click();
    await expect(page.locator("#gsxui-toaster [data-gsxui-toast]").first()).toBeVisible();
  });
});
```

Adjust selectors to the harness's real DOM if an assertion fails (e.g. the dialog element or toast row attribute names) — inspect the rendered page rather than guessing; the existing dialog/toast specs show the canonical selectors.

- [ ] **Step 2: Run the new spec**

Run: `npx playwright test --config jstest/playwright.config.ts jstest/specs/home-showcase.spec.ts`
Expected: PASS (4 tests). Known trap: without `--config` every test fails "Cannot navigate to invalid URL".

- [ ] **Step 3: Run the full gate set**

```bash
go build ./... && \
go test ./... -count=1 && \
make audit && \
make verify-generated && \
gofmt -l . && \
npx playwright test --config jstest/playwright.config.ts
```

Expected: all pass; `gofmt -l .` prints nothing. Known flake: `dialog.spec.ts:453`/`503` are animation-timing flakes under parallel load — re-run in isolation before blaming the change.

- [ ] **Step 4: Visual sanity check**

Start the site (`make site-dev` or the harness), load `/`, and confirm: 2×2 grid on desktop, stacked on mobile width, dark mode looks right (toggle the theme control in the nav). No code change expected; fix any layout issue found and amend.

- [ ] **Step 5: Commit**

```bash
git add jstest/specs/home-showcase.spec.ts
git commit -m "test(site): add Playwright coverage for home showcase"
```
