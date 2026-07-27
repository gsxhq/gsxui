package slotattr

import (
	"strings"

	gsx "github.com/gsxhq/gsx"
)

const attribute = "data-gsxui-slot"

func With(name string, attrs gsx.Attrs) gsx.Attrs {
	if name == "" {
		panic("gsxui: empty slot name")
	}

	tokens := []string{name}
	seen := map[string]struct{}{name: {}}
	if value, ok := attrs.Get(attribute); ok {
		if slotValue, err := gsx.AttrString(value); err == nil {
			for token := range strings.FieldsSeq(slotValue) {
				if _, ok := seen[token]; ok {
					continue
				}
				seen[token] = struct{}{}
				tokens = append(tokens, token)
			}
		}
	}

	out := make(gsx.Attrs, 1, len(attrs)+1)
	out[0] = gsx.Attr{Key: attribute, Value: strings.Join(tokens, " ")}
	return append(out, attrs.Without(attribute)...)
}
