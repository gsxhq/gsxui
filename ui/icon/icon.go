// Package icon is the Lucide icon set as tag-callable gsx components:
// <icon.ChevronDown class="size-4"/>. Icon data (ui/icon/icon_data.go) and
// vars (ui/icon/icon_defs.go) are generated — see internal/lucidegen. Lucide
// is ISC-licensed — see NOTICE.md.
package icon

import (
	"context"
	"fmt"
	"io"

	"github.com/gsxhq/gsx"
)

// New returns a tag-callable icon component for a Lucide icon name. An
// unknown name is a render-time error, never a silently empty node.
func New(name string) func(attrs ...gsx.Attr) gsx.Node {
	return func(attrs ...gsx.Attr) gsx.Node {
		inner, ok := data[name]
		if !ok {
			return gsx.Func(func(_ context.Context, _ io.Writer) error {
				return fmt.Errorf("ui/icon: unknown icon %q", name)
			})
		}
		return svgIcon(name, gsx.Raw(inner), gsx.Attrs(attrs))
	}
}
