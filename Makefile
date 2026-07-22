.PHONY: generate test check icons

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
	@for f in $$(find ui -name '*.js'); do node --check $$f || exit 1; done
	gofmt -l . | (! grep .)
