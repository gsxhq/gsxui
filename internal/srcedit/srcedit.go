// Package srcedit applies byte-offset replacements to a source buffer. It
// is the one edit applier shared by the tools that rewrite gsx source in
// place (internal/stylegen's recipe desugaring, internal/rtl's class
// rewriting): parse, record spans, apply from the end so earlier offsets
// stay valid.
package srcedit

import (
	"bytes"
	"fmt"
	"sort"
)

// Edit replaces src[Start:End] with Value.
type Edit struct {
	Start int
	End   int
	Value string
}

// Apply returns a copy of src with every edit applied. Edits may arrive in
// any order; they must not overlap or leave the buffer.
func Apply(src []byte, edits []Edit) ([]byte, error) {
	sorted := append([]Edit(nil), edits...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Start > sorted[j].Start })
	out := append([]byte(nil), src...)
	lastStart := len(src)
	for _, edit := range sorted {
		if edit.Start < 0 || edit.End < edit.Start || edit.End > len(src) {
			return nil, fmt.Errorf("invalid span [%d:%d] in %d bytes", edit.Start, edit.End, len(src))
		}
		if edit.End > lastStart {
			return nil, fmt.Errorf("overlapping span [%d:%d]", edit.Start, edit.End)
		}
		var next bytes.Buffer
		next.Grow(len(out) - (edit.End - edit.Start) + len(edit.Value))
		next.Write(out[:edit.Start])
		next.WriteString(edit.Value)
		next.Write(out[edit.End:])
		out = next.Bytes()
		lastStart = edit.Start
	}
	return out, nil
}
