.PHONY: generate test check icons site-dev site

generate:
	go tool gsx generate

# icons regenerates ui/icon/icon_data.go and ui/icon/icon_defs.go from a
# local Lucide checkout. See internal/lucidegen and Task 1's brief for the
# clone step (git clone --depth 1 https://github.com/lucide-icons/lucide
# /tmp/lucide-checkout).
icons:
	go run ./internal/lucidegen -src /tmp/lucide-checkout/icons -out ui/icon

test: generate
	go vet ./...
	go test ./...

check: test
	@git diff --exit-code -- '*.x.go' || { echo "error: generated .x.go drifted — commit regenerated output"; exit 1; }
	@test -z "$$(git status --porcelain -- '*.x.go' | grep '^??')" || { echo "error: untracked .x.go files"; exit 1; }
	@test -f site/dist/.gitkeep || { echo "error: site/dist/.gitkeep missing (vite build deletes it — restore before commit)"; exit 1; }
	@for f in $$(find ui -name '*.js'); do node --check $$f || exit 1; done
	gofmt -l . | (! grep .)

# site-dev runs the two-command dev loop: `npm install` once, then this.
# `gsx dev` warm-generates .x.go, builds-then-swaps the site/ binary (see
# gsx.toml [dev]), and runs Vite as the front door (proxying everything but
# its own /__vite/ namespace to the Go server).
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
