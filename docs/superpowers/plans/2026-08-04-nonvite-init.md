# Non-Vite Init Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** `gsxui init` auto-detects projects without the gsx/Vite scaffold and initializes them npm-free (vendoring tw-animate-css as CSS), with matching Getting Started docs.

**Architecture:** A scaffold-presence probe in `internal/cli/scaffold_integration.go` puts init into one of two modes. Vite mode is byte-for-byte today's behavior. Non-Vite mode drops all npm interaction, vendors `animate.css` next to the other CSS assets, rewrites one `@import` line in the vendored `index.css`, and prints different next steps. Docs quote the real CLI output (an existing test enforces snippet↔CLI equality).

**Tech Stack:** Go (CLI in `internal/cli`), gsx site pages (`site/pages/*.gsx`), Makefile, `go:embed` assets.

## Global Constraints

- Non-Vite mode runs **no npm commands** and writes **no package.json / package-lock.json** (spec: "we assume npm may not exist").
- We recommend tools only as examples; init never installs or configures a CSS build tool.
- Vite-mode behavior is unchanged, including its npm command list and summary output.
- `site/snippets/manual-integration.txt` must equal `manualScaffoldIntegrationInstructions` verbatim (`TestManualScaffoldInstructionsMatchDocumentation`).
- Any edit to `site/**/*.gsx` or new snippet requires `make highlight` and committing the regenerated `site/hl/blocks.gen.go` (TestBlocksMatchSourceText).
- Verification gates before claiming done: `go build ./...`, `go test ./... -count=1`, `make audit`, `gofmt -l .` clean. (CSS/Playwright gates are untouched by this work but run `npx playwright test --config jstest/playwright.config.ts` if any site page changed.)
- Commits: small, one per task, message style `feat(cli): …` / `docs(site): …`; end commit messages with the Claude-Session trailer used on this branch.

---

### Task 1: Vendored animate.css asset + refresh target + sync test

tw-animate-css's compiled sheet becomes a committed, embedded asset so non-Vite init can vendor it without npm.

**Files:**
- Create: `assets/css/animate.css` (copied from `node_modules/tw-animate-css/dist/tw-animate.css` + header)
- Modify: `Makefile` (new `generate-animate` target, wired into `generate`)
- Modify: `NOTICE.md` (attribution entry)
- Test: `internal/cli/animate_asset_test.go`

**Interfaces:**
- Produces: embedded file `assets/css/animate.css` readable via `fs.ReadFile(gsxui.Files, "assets/css/animate.css")` (the `//go:embed ui assets registry merge NOTICE.md` directive in `embed.go` already covers `assets/`). Header is exactly two lines, each starting `/*` — Task 2 strips nothing; it uses the file verbatim.

- [ ] **Step 1: Write the failing sync test**

`internal/cli/animate_asset_test.go`:

```go
package cli

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	gsxui "github.com/gsxhq/gsxui"
)

// TestAnimateAssetMatchesInstalledPackage pins assets/css/animate.css to the
// npm package it was copied from, so a tw-animate-css upgrade fails loudly
// until `make generate-animate` refreshes the embedded copy.
func TestAnimateAssetMatchesInstalledPackage(t *testing.T) {
	embedded, err := fs.ReadFile(gsxui.Files, "assets/css/animate.css")
	if err != nil {
		t.Fatal(err)
	}
	upstream, err := os.ReadFile(filepath.Join("..", "..", "node_modules", "tw-animate-css", "dist", "tw-animate.css"))
	if err != nil {
		t.Fatal(err)
	}
	header, body, found := bytes.Cut(embedded, []byte("*/\n"))
	if !found || !bytes.Contains(header, []byte("tw-animate-css")) {
		t.Fatalf("assets/css/animate.css must start with a tw-animate-css attribution comment, got %q", embedded[:min(len(embedded), 80)])
	}
	if !bytes.Equal(bytes.TrimSpace(body), bytes.TrimSpace(upstream)) {
		t.Fatal("assets/css/animate.css drifted from node_modules/tw-animate-css/dist/tw-animate.css — run `make generate-animate`")
	}
}
```

- [ ] **Step 2: Run it to verify it fails**

Run: `go test ./internal/cli -run TestAnimateAssetMatchesInstalledPackage -count=1`
Expected: FAIL (file does not exist).

- [ ] **Step 3: Add the Makefile target and generate the asset**

Makefile (next to the other generate targets; match existing tab/recipe style):

```make
# Refresh the embedded tw-animate-css copy vendored by non-Vite `gsxui init`.
generate-animate:
	printf '/* tw-animate-css (MIT) — https://github.com/Wombosvideo/tw-animate-css\n   Vendored by gsxui for npm-free builds; refresh with `make generate-animate`. */\n' > assets/css/animate.css
	cat node_modules/tw-animate-css/dist/tw-animate.css >> assets/css/animate.css
```

Add `generate-animate` to the existing `generate:` prerequisite list. Run `make generate-animate`.

- [ ] **Step 4: Add the NOTICE.md entry**

Append to `NOTICE.md`, matching its existing entry format:

```
tw-animate-css (assets/css/animate.css)
MIT License — Copyright (c) Wombosvideo
https://github.com/Wombosvideo/tw-animate-css
```

- [ ] **Step 5: Run the test to verify it passes**

Run: `go test ./internal/cli -run TestAnimateAssetMatchesInstalledPackage -count=1`
Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add assets/css/animate.css Makefile NOTICE.md internal/cli/animate_asset_test.go
git commit -m "feat(assets): vendor tw-animate-css as embedded animate.css"
```

---

### Task 2: Scaffold-presence probe and non-Vite refusal guidance

**Files:**
- Modify: `internal/cli/scaffold_integration.go`
- Modify: `site/snippets/manual-integration.txt` (keep drift test green)
- Test: `internal/cli/scaffold_integration_test.go`

**Interfaces:**
- Produces: `func scaffoldAbsent(dir string) (bool, error)` — true iff **neither** `vite.config.ts` nor `web/main.js` exists (any other stat error is returned). Also `const nonViteGuidance string` appended to `manualScaffoldIntegrationInstructions`.
- Consumes: nothing from other tasks.

- [ ] **Step 1: Write failing tests**

Append to `internal/cli/scaffold_integration_test.go`:

```go
func TestScaffoldAbsent(t *testing.T) {
	dir := t.TempDir()
	absent, err := scaffoldAbsent(dir)
	if err != nil || !absent {
		t.Fatalf("empty dir: absent=%v err=%v, want true nil", absent, err)
	}

	withVite := t.TempDir()
	writeFile(t, withVite, "vite.config.ts", "export default {}\n")
	absent, err = scaffoldAbsent(withVite)
	if err != nil || absent {
		t.Fatalf("vite.config.ts present: absent=%v err=%v, want false nil", absent, err)
	}

	withMain := t.TempDir()
	writeFile(t, withMain, "web/main.js", "console.log(1)\n")
	absent, err = scaffoldAbsent(withMain)
	if err != nil || absent {
		t.Fatalf("web/main.js present: absent=%v err=%v, want false nil", absent, err)
	}
}

func TestRefusalIncludesNonViteGuidance(t *testing.T) {
	dir := writeScaffoldFixture(t, []byte("tampered"), testGSXMainPristine)
	_, err := planScaffoldIntegration(dir, DefaultConfig())
	if err == nil || !strings.Contains(err.Error(), "without Vite") {
		t.Fatalf("refusal must mention the non-Vite path, got: %v", err)
	}
}
```

- [ ] **Step 2: Run to verify failure**

Run: `go test ./internal/cli -run 'TestScaffoldAbsent|TestRefusalIncludesNonViteGuidance' -count=1`
Expected: FAIL (`scaffoldAbsent` undefined; guidance missing).

- [ ] **Step 3: Implement**

In `scaffold_integration.go`:

```go
// scaffoldAbsent reports whether the project has no gsx npm/Vite scaffold at
// all — neither vite.config.ts nor web/main.js exists. That state selects
// init's non-Vite mode; a partially present or modified scaffold still goes
// through planScaffoldIntegration's refusal so tampering never silently
// downgrades an existing Vite project.
func scaffoldAbsent(dir string) (bool, error) {
	for _, relative := range []string{"vite.config.ts", "web/main.js"} {
		path, err := artifactPath(dir, relative)
		if err != nil {
			return false, err
		}
		if _, err := os.Stat(path); err == nil {
			return false, nil
		} else if !os.IsNotExist(err) {
			return false, err
		}
	}
	return true, nil
}
```

Extend the instructions const (this exact text also lands in the snippet and docs):

```go
const manualScaffoldIntegrationInstructions = `Automatic integration currently supports an unmodified gsx init npm/Vite scaffold with the default gsxui JS and CSS paths.

Integrate this project manually:
  npm install --save-dev tailwindcss@^4.3.3 @tailwindcss/vite@^4.3.3 tw-animate-css@^1.4.0
  vite.config.ts: import tailwindcss from "@tailwindcss/vite" and add tailwindcss() to plugins
  web/main.js: import "./gsxui/index.js" and import "./gsxui/index.css"

Or, to run without Vite entirely, delete vite.config.ts and web/main.js and
re-run gsxui init: it will vendor everything self-contained (no npm), and you
serve the gsxui JS directory statically, load index.js with one
<script type="module"> tag, and build index.css with any Tailwind v4 tool
(e.g. npx @tailwindcss/cli -i web/gsxui/index.css -o dist.css).`
```

Update `site/snippets/manual-integration.txt` to the same text (trailing newline preserved — the drift test compares `strings.TrimSpace`d snippet content to the const).

- [ ] **Step 4: Run tests**

Run: `go test ./internal/cli -run 'TestScaffoldAbsent|TestRefusalIncludesNonViteGuidance|TestManualScaffoldInstructionsMatchDocumentation' -count=1`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/cli/scaffold_integration.go internal/cli/scaffold_integration_test.go site/snippets/manual-integration.txt
git commit -m "feat(cli): detect scaffold absence; refusal points at the non-Vite path"
```

---

### Task 3: Non-Vite init mode

**Files:**
- Modify: `internal/cli/init.go`
- Test: `internal/cli/init_test.go`

**Interfaces:**
- Consumes: `scaffoldAbsent(dir)` (Task 2); `assets/css/animate.css` (Task 1).
- Produces: non-Vite `runInit` behavior: vendors `web/gsxui/animate.css`; `web/gsxui/index.css` line 2 becomes `@import "./animate.css";`; npm never invoked; no package.json artifacts; summary text below. `initArtifacts` gains a trailing `nonVite bool` parameter.

- [ ] **Step 1: Write the failing test**

Append to `internal/cli/init_test.go` (reuses `writeFile`, `runCommand` stub pattern from `initTestModule`, but on a bare module — no scaffold fixture):

```go
// nonViteTestModule is initTestModule without the npm/Vite scaffold: just a
// Go module in an empty directory.
func nonViteTestModule(t *testing.T) (dir string, commands *[][]string) {
	t.Helper()
	dir = t.TempDir()
	writeFile(t, dir, "go.mod", "module example.com/app\n\ngo 1.26\n")
	var got [][]string
	orig := runCommand
	runCommand = func(dir, name string, args ...string) error {
		got = append(got, append([]string{name}, args...))
		return nil
	}
	t.Cleanup(func() { runCommand = orig })
	t.Chdir(dir)
	return dir, &got
}

func TestInitNonViteVendorsWithoutNPM(t *testing.T) {
	dir, commands := nonViteTestModule(t)
	if err := Run([]string{"init"}); err != nil {
		t.Fatal(err)
	}
	for _, p := range []string{
		"gsxui.json",
		"gsxui.preset.json",
		"web/gsxui/index.css",
		"web/gsxui/animate.css",
		"web/gsxui/foundation.css",
		"web/gsxui/theme.css",
		"web/gsxui/style.css",
		"web/gsxui/gsxui.js",
		"web/gsxui/index.js",
		"ui/merge/merge.go",
		"gsx.toml",
	} {
		if _, err := os.Stat(filepath.Join(dir, p)); err != nil {
			t.Errorf("missing %s: %v", p, err)
		}
	}
	for _, p := range []string{"vite.config.ts", "web/main.js", "package.json", "package-lock.json"} {
		if _, err := os.Stat(filepath.Join(dir, p)); !os.IsNotExist(err) {
			t.Errorf("%s must not be created in non-Vite mode (err=%v)", p, err)
		}
	}
	index := readFile(t, dir, "web/gsxui/index.css")
	if !strings.Contains(index, `@import "./animate.css";`) || strings.Contains(index, `@import "tw-animate-css";`) {
		t.Fatalf("index.css must import ./animate.css instead of tw-animate-css:\n%s", index)
	}
	for _, command := range *commands {
		if command[0] == "npm" {
			t.Fatalf("npm must never run in non-Vite mode, got %v", *commands)
		}
	}
	// go tooling still installed, and generation still runs.
	joined := fmt.Sprint(*commands)
	for _, want := range []string{"go get github.com/gsxhq/gsx@latest", "go get -tool github.com/gsxhq/gsx/cmd/gsx@latest"} {
		if !strings.Contains(joined, want) {
			t.Errorf("missing command %q in %v", want, *commands)
		}
	}
}

func TestInitNonVitePrintsServingInstructions(t *testing.T) {
	_, _ = nonViteTestModule(t)
	var out bytes.Buffer
	origOut := commandStdout
	commandStdout = &out
	t.Cleanup(func() { commandStdout = origOut })
	if err := Run([]string{"init"}); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"no Vite scaffold detected",
		`<script type="module"`,
		"@tailwindcss/cli",
	} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("summary missing %q:\n%s", want, out.String())
		}
	}
}
```

(Add `"fmt"` to the test file imports if absent.)

- [ ] **Step 2: Run to verify failure**

Run: `go test ./internal/cli -run 'TestInitNonVite' -count=1`
Expected: FAIL — today `Run([]string{"init"})` errors with the scaffold refusal.

- [ ] **Step 3: Implement in `runInit`**

In `init.go`, replace the fixed integration/npm block:

```go
	nonVite, err := scaffoldAbsent(dir)
	if err != nil {
		return err
	}
	artifacts, err := initArtifacts(dir, module, cfg, selectedPreset, nonVite)
	if err != nil {
		return err
	}
	if !nonVite {
		integrationArtifacts, err := planScaffoldIntegration(dir, cfg)
		if err != nil {
			return err
		}
		artifacts = append(artifacts, integrationArtifacts...)
		packageArtifacts, err := packageMetadataArtifacts(dir)
		if err != nil {
			return err
		}
		artifacts = append(artifacts, packageArtifacts...)
	}
```

Commands: keep `commands` (the go get list) unconditional; guard only npm:

```go
		func() error {
			if !nonVite {
				if err := runCommand(dir, "npm", npmCommand...); err != nil {
					return fmt.Errorf("npm %v: %w", npmCommand, err)
				}
			}
			for _, command := range commands {
				...
```

`initArtifacts(dir, module string, cfg Config, selected preset.Preset, nonVite bool)`: when `nonVite`, after the css asset loop append the animate artifact and rewrite the index entry. Concretely, inside the `cssAssetTargets` loop, special-case the entry (`asset.source == "assets/css/index.css"`):

```go
		if nonVite && asset.source == "assets/css/index.css" {
			content = bytes.Replace(content, []byte("@import \"tw-animate-css\";"), []byte("@import \"./animate.css\";"), 1)
		}
```

and after the loop:

```go
	if nonVite {
		animate, err := fs.ReadFile(gsxui.Files, "assets/css/animate.css")
		if err != nil {
			return nil, err
		}
		artifacts = append(artifacts, artifact{
			RelativePath: filepath.ToSlash(filepath.Join(filepath.Dir(cfg.CSS), "animate.css")),
			Content:      animate,
			Managed:      true,
		})
	}
```

(Add `"bytes"` to init.go imports.)

Summary output — replace the single `fmt.Fprintf` with:

```go
	if nonVite {
		fmt.Fprintf(
			commandStdout,
			`gsxui initialized (no Vite scaffold detected — npm-free mode).
  css:  %[1]s
  js:   %[2]s/index.js
  next: gsxui add button

Serve and build with your own tooling:
  1. serve %[2]s/ statically and load it with
     <script type="module" src="/your-prefix/index.js"></script>
  2. build %[1]s with any Tailwind v4 tool, for example:
     npx @tailwindcss/cli -i %[1]s -o dist/gsxui.css
     (or the standalone tailwindcss binary)
  3. link the built stylesheet from your pages
`,
			cfg.CSS,
			cfg.JS,
		)
		return nil
	}
	fmt.Fprintf( /* existing Vite summary, unchanged */ ...)
```

- [ ] **Step 4: Run the new tests, then the whole cli package**

Run: `go test ./internal/cli -run 'TestInitNonVite' -count=1` → PASS.
Run: `go test ./internal/cli -count=1`
Expected: all PASS — in particular `TestInitWritesEverything` (Vite mode untouched; its fixture has the scaffold so `nonVite` is false).

- [ ] **Step 5: Commit**

```bash
git add internal/cli/init.go internal/cli/init_test.go
git commit -m "feat(cli): npm-free non-Vite init mode with vendored animate.css"
```

---

### Task 4: `gsxui add` after non-Vite init keeps animate.css intact

Regression guard: `add` regenerates the behavior barrel in `cfg.JS`; `animate.css` lives in `filepath.Dir(cfg.CSS)` which is the same directory under default config. Assert `add` neither deletes nor rewrites it, and that a component's JS lands beside it.

**Files:**
- Test: `internal/cli/add_test.go`

**Interfaces:**
- Consumes: non-Vite init behavior (Task 3), `nonViteTestModule` helper (Task 3).

- [ ] **Step 1: Write the test**

Append to `internal/cli/add_test.go`:

```go
func TestAddAfterNonViteInitPreservesAnimateCSS(t *testing.T) {
	dir, _ := nonViteTestModule(t)
	if err := Run([]string{"init"}); err != nil {
		t.Fatal(err)
	}
	before := readFile(t, dir, "web/gsxui/animate.css")
	if err := Run([]string{"add", "dialog"}); err != nil {
		t.Fatal(err)
	}
	if got := readFile(t, dir, "web/gsxui/animate.css"); got != before {
		t.Fatal("gsxui add must not rewrite animate.css")
	}
	index := readFile(t, dir, "web/gsxui/index.js")
	if !strings.Contains(index, `import "./dialog.js";`) {
		t.Fatalf("barrel missing dialog behavior:\n%s", index)
	}
}
```

- [ ] **Step 2: Run it**

Run: `go test ./internal/cli -run TestAddAfterNonViteInitPreservesAnimateCSS -count=1`
Expected: PASS immediately (add's barrel globbing only touches `*.js`). If it fails instead, the barrel or artifact plan is treating `animate.css` as managed JS — fix in `add.go`'s artifact planning, not the test.

- [ ] **Step 3: Commit**

```bash
git add internal/cli/add_test.go
git commit -m "test(cli): add after non-Vite init leaves animate.css alone"
```

---

### Task 5: Getting Started "Without Vite" docs

**Files:**
- Modify: `site/pages/getting_started.gsx` (new TOC item + section)
- Create: `site/snippets/nonvite-init.output.txt`, `site/snippets/nonvite-serve.go.txt`, `site/snippets/nonvite-css.sh.txt`
- Modify (generated): `site/hl/blocks.gen.go` via `make highlight`
- Test: existing `TestBlocksMatchSourceText` (site/hl) + `TestManualScaffoldInstructionsMatchDocumentation` stay green; site tests in `site/pages`.

**Interfaces:**
- Consumes: exact non-Vite summary text (Task 3) — `nonvite-init.output.txt` must be copied from a real `gsxui init` run in a scratch module, per the page's "real CLI output, not invented" convention.

- [ ] **Step 1: Capture real output snippets**

In a scratch dir (temp Go module, no scaffold): run the freshly built CLI `go run ./cmd/gsxui init` and copy its stdout verbatim into `site/snippets/nonvite-init.output.txt`.

`site/snippets/nonvite-serve.go.txt` (from the verified scenario, trimmed to essentials):

```go
mux := http.NewServeMux()
mux.Handle("/assets/gsxui/",
	http.StripPrefix("/assets/gsxui/", http.FileServer(http.Dir("web/gsxui"))))
```

with the page-side tags shown alongside in the section prose:

```html
<link rel="stylesheet" href="/assets/gsxui.css"/>
<script type="module" src="/assets/gsxui/index.js"></script>
```

`site/snippets/nonvite-css.sh.txt`:

```sh
# any Tailwind v4 tool works; two examples
npx @tailwindcss/cli -i web/gsxui/index.css -o dist/gsxui.css
tailwindcss -i web/gsxui/index.css -o dist/gsxui.css   # standalone binary
```

- [ ] **Step 2: Add the section to `getting_started.gsx`**

Add to `gettingStartedTOCItems` after the manual-integration entry:

```go
	{ID: "without-vite", Title: "Without Vite", Depth: 3},
```

(keep subsequent items; IDs are stable strings, order in the slice is the page order — insert, don't renumber). Then, after the manual-integration `<div>`, in the same section:

```
<div class="mt-4 flex flex-col gap-3">
	<docHeading item={gettingStartedTOCItems[3]}/>
	<p>
		gsxui does not require Vite — or npm. If <code>gsxui init</code> finds neither
		<code>vite.config.ts</code> nor <code>web/main.js</code>, it initializes in npm-free mode: the
		component behaviors are dependency-free native ES modules, and the one CSS dependency
		(<code>tw-animate-css</code>) is vendored as <code>web/gsxui/animate.css</code> so the stylesheet
		builds without <code>node_modules</code>.
	</p>
	<pre><code>{ hl.Node("snippets/nonvite-init.output") }</code></pre>
	<p>
		You own serving and CSS building. Serve the vendored JS directory statically and load the barrel
		with one module script tag — no bundler required, any bundler welcome:
	</p>
	<pre><code>{ hl.Node("snippets/nonvite-serve.go") }</code></pre>
	<p>
		Build the CSS entry with the Tailwind v4 tool of your choice (these are examples — gsxui never
		installs or manages your build tooling), then link the output from your pages:
	</p>
	<pre><code>{ hl.Node("snippets/nonvite-css.sh") }</code></pre>
</div>
```

Adjust the `gettingStartedTOCItems[N]` indices used by later sections (add-components, first-page) for the inserted item. Use `{/* */}` for any comments inside markup, never `//`.

- [ ] **Step 3: Regenerate and verify**

Run: `make highlight` (regenerates `site/hl/blocks.gen.go`; commit it alongside).
Run: `go tool gsx generate && go build ./... && go test ./site/... -count=1`
Expected: build green, `TestBlocksMatchSourceText` PASS.

- [ ] **Step 4: Visual check**

Run the site (`npm run dev` or harness) and confirm /docs/getting-started renders the new section with highlighted snippets and a working TOC link.

- [ ] **Step 5: Commit**

```bash
git add site/pages/getting_started.gsx site/pages/getting_started.x.go site/snippets/nonvite-*.txt site/hl/blocks.gen.go
git commit -m "docs(site): Without Vite section on Getting Started"
```

---

### Task 6: Full verification gates

- [ ] **Step 1: Run the gate list**

```bash
go build ./... && gofmt -l . && go test ./... -count=1
make audit && make verify-generated && make verify-generated-styles
npx playwright test --config jstest/playwright.config.ts
```

Expected: all green; `gofmt -l` prints nothing. (Playwright note: `dialog.spec.ts:453/503` are known parallel-load flakes — retry in isolation before blaming the change.)

- [ ] **Step 2: Real-world smoke (optional but cheap)**

Repeat the scratchpad scenario end-to-end with the new CLI: bare module → `gsxui init` (expect npm-free summary) → `gsxui add dialog` → serve with the plain Go server → dialog opens in browser.

- [ ] **Step 3: Finish**

Use superpowers:finishing-a-development-branch to merge/PR `nonvite-init`.
