.PHONY: generate test check

generate:
	go tool gsx generate

test: generate
	go vet ./...
	go test ./...

check: test
	@for f in $$(find ui -name '*.js'); do node --check $$f || exit 1; done
	gofmt -l . | (! grep .)
