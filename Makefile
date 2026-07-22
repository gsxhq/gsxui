.PHONY: generate test check

generate:
	go tool gsx generate

test: generate
	go vet ./...
	go test ./...

check: test
	@git diff --exit-code -- '*.x.go' || { echo "error: generated .x.go drifted — commit regenerated output"; exit 1; }
	@test -z "$$(git status --porcelain -- '*.x.go' | grep '^??')" || { echo "error: untracked .x.go files"; exit 1; }
	@for f in $$(find ui -name '*.js'); do node --check $$f || exit 1; done
	gofmt -l . | (! grep .)
