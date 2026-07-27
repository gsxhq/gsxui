package ui

import (
	"strings"

	gsx "github.com/gsxhq/gsx"
)

const slotAttribute = "data-gsxui-slot"

func withSlot(name string, attrs gsx.Attrs) gsx.Attrs {
	if name == "" {
		panic("gsxui: empty slot name")
	}

	tokens := []string{name}
	seen := map[string]struct{}{name: {}}
	if value, ok := attrs.Get(slotAttribute); ok {
		slotValue, err := gsx.AttrString(value)
		if err != nil {
			panic(err)
		}
		for token := range strings.FieldsSeq(slotValue) {
			if _, ok := seen[token]; ok {
				continue
			}
			seen[token] = struct{}{}
			tokens = append(tokens, token)
		}
	}

	out := make(gsx.Attrs, 1, len(attrs)+1)
	out[0] = gsx.Attr{Key: slotAttribute, Value: strings.Join(tokens, " ")}
	return append(out, attrs.Without(slotAttribute)...)
}
