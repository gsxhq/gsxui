package ui

import (
	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui/internal/slotattr"
)

func withSlot(name string, attrs gsx.Attrs) gsx.Attrs {
	return slotattr.With(name, attrs)
}
