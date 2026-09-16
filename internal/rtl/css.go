package rtl

import "regexp"

var applyList = regexp.MustCompile(`(@apply\s+)([^;{}]+?)(\s*;)`)

// CSS rewrites the utility list of every @apply. Selectors, properties and
// comments are untouched; the pass is a fixed point like Classes.
func CSS(src []byte) []byte {
	return applyList.ReplaceAllFunc(src, func(m []byte) []byte {
		parts := applyList.FindSubmatch(m)
		return append(append(append([]byte{}, parts[1]...), Classes(string(parts[2]))...), parts[3]...)
	})
}
