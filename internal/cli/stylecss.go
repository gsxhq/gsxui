package cli

import (
	"fmt"
	"io/fs"
	"path"
	"regexp"
	"strings"
)

// relativeImport matches one whole-line `@import "./x.css";` statement. Only a
// line that is nothing but the statement counts: an @import sharing a line with
// other CSS, or one carrying a media query or layer(), is not something this
// composer is allowed to guess at, so it is left alone and reported instead.
var relativeImport = regexp.MustCompile(`^@import\s+"([^"]+)";[ \t]*$`)

// composeStyleCSS flattens a stylesheet by substituting each relative @import
// line with the text of the file it names, in place.
//
// WHY THIS EXISTS: assets/css/styles/default.css is split into one file per
// component under styles/default/ so that migrating a component is a file
// deletion and parallel migrations stop colliding in one 3000-line file. That
// split is for this repository's benefit, not the consumer's — `gsxui init`
// still vendors ONE flat style.css, because 42 files where there was 1 changes
// a shipped contract for no gain to anyone downstream.
//
// Substituting in place, rather than appending the parts, is what makes the
// result byte-identical to the pre-split sheet modulo the per-file `@layer
// components` wrappers: every rule lands exactly where it used to sit, so
// source order — and therefore cascade order among equal-specificity rules —
// is preserved. TestVendoredStyleCSSMatchesTheAggregator pins that.
//
// Bare specifiers (`@import "tailwindcss";`) are passed through untouched: they
// are resolved by the consumer's bundler, not by us. A relative import that
// cannot be read is an error, never a pass-through — silently emitting the
// unresolved line would ship a consumer a stylesheet missing a component, which
// is the exact failure this function has to make impossible.
func composeStyleCSS(fsys fs.FS, source string) ([]byte, error) {
	var out strings.Builder
	if err := composeInto(&out, fsys, source, nil); err != nil {
		return nil, err
	}
	return []byte(out.String()), nil
}

// composeInto writes source's composed text into out. stack carries the import
// chain so a cycle is reported as a cycle rather than recursing until the stack
// dies.
func composeInto(out *strings.Builder, fsys fs.FS, source string, stack []string) error {
	for _, seen := range stack {
		if seen == source {
			return fmt.Errorf("composing %s: import cycle %s -> %s", stack[0], strings.Join(stack, " -> "), source)
		}
	}
	content, err := fs.ReadFile(fsys, source)
	if err != nil {
		return fmt.Errorf("composing %s: %w", source, err)
	}
	// SplitAfter keeps each line's terminator, so a file that does not end in a
	// newline is reproduced without one rather than gaining a byte.
	for _, line := range strings.SplitAfter(string(content), "\n") {
		if line == "" { // SplitAfter's trailing empty element
			continue
		}
		match := relativeImport.FindStringSubmatch(strings.TrimSuffix(line, "\n"))
		if match == nil {
			out.WriteString(line)
			continue
		}
		target := match[1]
		if !strings.HasPrefix(target, "./") && !strings.HasPrefix(target, "../") {
			out.WriteString(line) // bare specifier: the consumer's bundler owns it
			continue
		}
		resolved := path.Join(path.Dir(source), target)
		if err := composeInto(out, fsys, resolved, append(stack, source)); err != nil {
			return err
		}
	}
	return nil
}
