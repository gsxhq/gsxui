// Separate module: this generator links tree-sitter (cgo), which must never
// reach gsxui's own module — site/main.go builds CGO_ENABLED=0 for a
// distroless/static image. The replace directives resolve the gsx grammar
// and highlighter from sibling checkouts; nothing here is needed to build or
// deploy the site, only to regenerate site/hl/blocks.gen.go.
module github.com/gsxhq/gsxui/site/hl/gen

go 1.26.1

require github.com/gsxhq/gsxhl v0.0.0

require (
	github.com/gsxhq/tree-sitter-gsx v0.0.0 // indirect
	github.com/jackielii/go-tree-sitter-highlight v0.1.0 // indirect
	github.com/mattn/go-pointer v0.0.1 // indirect
	github.com/tree-sitter/go-tree-sitter v0.25.0 // indirect
	github.com/tree-sitter/tree-sitter-css v0.25.0 // indirect
	github.com/tree-sitter/tree-sitter-javascript v0.25.0 // indirect
)

replace (
	github.com/gsxhq/gsxhl => ../../../../gsxhl
	github.com/gsxhq/tree-sitter-gsx => ../../../../tree-sitter-gsx
)
