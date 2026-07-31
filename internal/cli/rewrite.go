package cli

import (
	"bytes"
	goast "go/ast"
	goparser "go/parser"
	"go/token"
	"path"
	"sort"
	"strconv"
	"strings"

	gsxast "github.com/gsxhq/gsx/ast"
	gsxparser "github.com/gsxhq/gsx/parser"
	gomodule "golang.org/x/mod/module"
)

const gsxuiUIImport = "github.com/gsxhq/gsxui/ui"

type importLiteralEdit struct {
	start int
	end   int
	value string
}

// RewriteGsx retargets exact gsxui UI imports at the user's module. It parses
// the GSX file and the Go chunks that can own import declarations, then edits
// only import-path BasicLits identified by their token positions. Comments,
// ordinary string literals, and decoded lookalike paths are not rewrite
// candidates.
func RewriteGsx(src []byte, modulePath, uiDir string) ([]byte, error) {
	destination, err := rewriteDestination(modulePath, uiDir)
	if err != nil {
		return nil, err
	}

	fset := token.NewFileSet()
	file, err := gsxparser.ParseFile(fset, "rewrite.gsx", src, 0)
	if err != nil {
		return nil, err
	}

	var edits []importLiteralEdit
	for _, declaration := range file.Decls {
		chunk, ok := declaration.(*gsxast.GoChunk)
		if !ok {
			continue
		}
		chunkEdits, err := importEditsForChunk(src, fset, chunk, destination)
		if err != nil {
			return nil, err
		}
		edits = append(edits, chunkEdits...)
	}

	rewritten, err := applyImportLiteralEdits(src, edits)
	if err != nil {
		return nil, err
	}
	if _, err := gsxparser.ParseFile(token.NewFileSet(), "rewrite.gsx", rewritten, 0); err != nil {
		return nil, err
	}
	return rewritten, nil
}

func rewriteDestination(modulePath, uiDir string) (string, error) {
	if modulePath == "" || path.Clean(modulePath) != modulePath {
		return "", &invalidRewriteDestination{field: "module", value: modulePath}
	}
	if err := gomodule.CheckImportPath(modulePath); err != nil {
		return "", &invalidRewriteDestination{field: "module", value: modulePath, cause: err}
	}
	if err := validateManagedPath(uiDir); err != nil {
		return "", &invalidRewriteDestination{field: "UI directory", value: uiDir, cause: err}
	}
	destination := modulePath + "/" + uiDir
	if err := gomodule.CheckImportPath(destination); err != nil {
		return "", &invalidRewriteDestination{field: "UI import path", value: destination, cause: err}
	}
	return destination, nil
}

type invalidRewriteDestination struct {
	field string
	value string
	cause error
}

func (err *invalidRewriteDestination) Error() string {
	if err.cause != nil {
		return "invalid " + err.field + " " + strconv.Quote(err.value) + ": " + err.cause.Error()
	}
	return "invalid " + err.field + " " + strconv.Quote(err.value)
}

func (err *invalidRewriteDestination) Unwrap() error {
	return err.cause
}

func importEditsForChunk(src []byte, gsxFileSet *token.FileSet, chunk *gsxast.GoChunk, destination string) ([]importLiteralEdit, error) {
	const prefix = "package _gsxrewrite\n"
	goFileSet := token.NewFileSet()
	parsed, err := goparser.ParseFile(
		goFileSet,
		"rewrite.gsx",
		prefix+chunk.Src,
		goparser.SkipObjectResolution,
	)
	if err != nil {
		return nil, err
	}

	chunkStart := gsxFileSet.Position(chunk.Pos()).Offset
	var edits []importLiteralEdit
	for _, declaration := range parsed.Decls {
		general, ok := declaration.(*goast.GenDecl)
		if !ok || general.Tok != token.IMPORT {
			continue
		}
		for _, rawSpec := range general.Specs {
			spec := rawSpec.(*goast.ImportSpec)
			importPath, err := strconv.Unquote(spec.Path.Value)
			if err != nil {
				return nil, err
			}
			if importPath != gsxuiUIImport && !strings.HasPrefix(importPath, gsxuiUIImport+"/") {
				continue
			}

			literalStart := goFileSet.Position(spec.Path.Pos()).Offset - len(prefix)
			literalEnd := goFileSet.Position(spec.Path.End()).Offset - len(prefix)
			start := chunkStart + literalStart
			end := chunkStart + literalEnd
			if start < 0 || end < start || end > len(src) ||
				!bytes.Equal(src[start:end], []byte(spec.Path.Value)) {
				return nil, &invalidImportLiteralPosition{start: start, end: end, size: len(src)}
			}
			suffix := strings.TrimPrefix(importPath, gsxuiUIImport)
			edits = append(edits, importLiteralEdit{
				start: start,
				end:   end,
				value: strconv.Quote(destination + suffix),
			})
		}
	}
	return edits, nil
}

type invalidImportLiteralPosition struct {
	start int
	end   int
	size  int
}

func (err *invalidImportLiteralPosition) Error() string {
	return "invalid GSX import literal position " +
		strconv.Itoa(err.start) + ":" + strconv.Itoa(err.end) +
		" for source size " + strconv.Itoa(err.size)
}

func applyImportLiteralEdits(src []byte, edits []importLiteralEdit) ([]byte, error) {
	sorted := append([]importLiteralEdit(nil), edits...)
	sort.Slice(sorted, func(left, right int) bool {
		return sorted[left].start > sorted[right].start
	})
	output := append([]byte(nil), src...)
	previousStart := len(src)
	for _, edit := range sorted {
		if edit.start < 0 || edit.end < edit.start || edit.end > previousStart {
			return nil, &invalidImportLiteralPosition{start: edit.start, end: edit.end, size: len(src)}
		}
		output = append(output[:edit.start], append([]byte(edit.value), output[edit.end:]...)...)
		previousStart = edit.start
	}
	return output, nil
}

// Barrel renders index.js for the installed behavior set. Core (gsxui.js)
// sits beside the behaviors, so no path rewriting is needed anywhere in JS.
func Barrel(behaviors []string) []byte {
	sorted := append([]string(nil), behaviors...)
	sort.Strings(sorted)
	var b bytes.Buffer
	b.WriteString("// Generated by gsxui — edit by adding/removing components, not by hand.\n")
	b.WriteString("export * from \"./gsxui.js\";\n")
	for _, name := range sorted {
		b.WriteString("import \"./")
		b.WriteString(name)
		b.WriteString(".js\";\n")
	}
	return b.Bytes()
}
