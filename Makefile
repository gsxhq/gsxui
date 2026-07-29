.PHONY: generate generate-styles verify-generated verify-generated-styles test test-js test-theme-state test-css-audit audit check ci icons site-dev site highlight

audit-source-dirs := ui registry/canonical site/examples site/pages web $(wildcard dev)
audit-css-source-dirs := assets/css site web $(wildcard dev)

generate: generate-styles
	go tool gsx generate

generate-styles:
	go run ./cmd/stylegen

verify-generated:
	go run ./internal/generatedcheck/cmd

verify-generated-styles:
	go run ./cmd/stylegen --check

# highlight regenerates site/hl/blocks.gen.go — every component example and
# doc snippet pre-rendered to highlighted HTML. Run it after adding, renaming
# or editing anything under site/examples/ or site/snippets/; TestBlocksCoverSources
# and TestBlocksMatchSourceText fail when the committed output goes stale.
#
# The generator is a SEPARATE module (site/hl/gen/go.mod) because tree-sitter
# is C: keeping it out of gsxui's module is what lets site/main.go keep
# building CGO_ENABLED=0 into a distroless/static image. It resolves the
# grammar and highlighter from sibling checkouts via replace directives, so
# it needs ../tree-sitter-gsx and ../gsxhl next to this repo — nothing else
# does, including CI and the Docker build, which consume the committed output.
highlight:
	cd site/hl/gen && go run .

# icons regenerates ui/icon/icon_data.go and ui/icon/icon_defs.go from a
# local Lucide checkout. See internal/lucidegen and Task 1's brief for the
# clone step (git clone --depth 1 https://github.com/lucide-icons/lucide
# /tmp/lucide-checkout).
icons:
	go run ./internal/lucidegen -src /tmp/lucide-checkout/icons -out ui/icon

test: generate
	go vet ./...
	go test ./...

# test-js runs the Playwright suite in jstest/ against real Chromium. It is
# deliberately NOT part of `make test` — a Go-only edit should not pay for a
# browser boot. Playwright's globalSetup writes the example manifest and
# compiles web/site.css into jstest/.tmp/ (gitignored), and its webServer
# block starts jstest/harness on 127.0.0.1:7799 — clear of the dev loop's
# 7777 and Vite's 5173, so `make site-dev` can keep running.
#
# First run on a new machine needs the browser: `npx playwright install chromium`.
test-js:
	npx playwright test --config jstest/playwright.config.ts

test-css-audit:
	node --test jstest/support/compiled-css-audit.test.ts

test-theme-state:
	node --test web/theme-state.test.js

audit:
	@! rg -n '^[[:space:]]*<[^>]*data-slot=|^[[:space:]]+data-slot=' $(audit-source-dirs) -g '!*.x.go' -g '!*.gen.go'
	@! rg -n 'data-slot|className[[:space:]]*=|[.]classList|setAttribute[(][^)]*class|[.]className[[:space:]]*=' ui $(wildcard dev) -g '*.js'
	@! rg -n 'data-slot|group/|peer/|(group|peer)-(data|has|focus|hover|active|disabled|aria|open|checked)[^[:space:]]*/' ui $(wildcard dev) -g '*.gsx' -g '!button.gsx'
	@! rg -n 'data-slot|peer/|(group|peer)-(data|has|focus|hover|active|disabled|aria|open|checked)[^[:space:]]*/' ui/button.gsx
	@! rg -n -P '\bwithSlot[[:space:]]*\(|\bslotattr[[:space:]]*\.[[:space:]]*With[[:space:]]*\(' ui $(wildcard dev) -g '*.gsx'
	@! rg -n -U -P '(?s)<[^>]*\bdata-gsxui-slot(?=[[:space:]]*(?:=|/?>))' $(audit-source-dirs) -g '*.gsx' -g '!*.x.go' -g '!*.gen.go'
	@! rg -n -P '(?:\bKey[[:space:]]*:[[:space:]]*|[.]SetAttribute\([[:space:]]*|\bsetAttribute\([[:space:]]*)["\x27\x60]data-gsxui-slot["\x27\x60]' $(audit-source-dirs) -g '*.go' -g '*.js' -g '!*.x.go' -g '!*.gen.go' -g '!*_test.go'
	@! rg -n -P '\[[[:space:]]*data-gsxui-slot(?=[[:space:]]*(?:[~|^$*]?=|\]))' $(audit-css-source-dirs) -g '*.css'
	@! rg -n -P '[\w./:\[\]&>-]+!' ui -g '*.gsx'
	@! rg -n 'gsxui-recipe-' ui -g '*.gsx'
	@! rg -n '^[[:space:]]+class=' ui -g '*.gsx' -g '!button.gsx'
	@! rg -n '^[[:space:]]*<[^>]*class=' ui -g '*.gsx' -g '!button.gsx'
	@! rg -n 'gsxui-recipe-' registry/canonical -g '*.gsx'
	@! rg -n '!important' assets/css/foundation.css assets/css/styles/default.css

check: audit test-css-audit test-theme-state verify-generated-styles
	@$(MAKE) --no-print-directory verify-generated
	go vet ./...
	go test ./...
	npx playwright test --config jstest/playwright.config.ts
	@test -f site/dist/.gitkeep || { echo "error: site/dist/.gitkeep missing (vite build deletes it — restore before commit)"; exit 1; }
	@for f in $$(find ui jstest web -name '*.js'); do node --check $$f || exit 1; done
	gofmt -l . | (! grep .)

# ci is the authoritative uncached gate. It mirrors check without reusing
# Go's test-result cache and keeps the browser, generation, syntax,
# structural, and formatting checks in the same run.
ci: audit test-css-audit test-theme-state verify-generated-styles
	@$(MAKE) --no-print-directory verify-generated
	go vet ./...
	go test -count=1 ./...
	npx playwright test --config jstest/playwright.config.ts
	@test -f site/dist/.gitkeep || { echo "error: site/dist/.gitkeep missing (vite build deletes it — restore before commit)"; exit 1; }
	@for f in $$(find ui jstest web -name '*.js'); do node --check $$f || exit 1; done
	gofmt -l . | (! grep .)

# site-dev runs the two-command dev loop: `npm install` once, then this.
# `gsx dev` warm-generates .x.go, builds-then-swaps the site/ binary (see
# gsx.toml [dev]), and runs Vite as the front door (proxying everything but
# its own /__vite/ namespace to the Go server).
# Backend output lands in tmp/dev.log via gsx.toml's [dev] log key (gsx dev
# tees it itself — no shell pipe, so the console keeps its exit status too).
site-dev:
	go tool gsx dev

# site builds the production bundle (Vite assets embedded by site/main.go)
# and runs the server in prod mode (no VITE_DEV_URL → gsxhq/vite serves the
# embedded dist/ instead of proxying to a dev server).
site:
	npx vite build
	@touch site/dist/.gitkeep
	go tool gsx generate
	go run ./site
