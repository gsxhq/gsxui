// Package icon holds the site's example gsx components for ui/icon.
package icon

import (
	"github.com/gsxhq/gsx"
	uiicon "github.com/gsxhq/gsxui/ui/icon"
)

// namedIcon pairs a display label with one of the 1,748 tag-callable
// Lucide icon components ui/icon/icon_defs.go generates.
type namedIcon struct {
	Name   string
	Render func(...gsx.Attr) gsx.Node
}

// popular is a hand-picked slice of well-known icons — the full set lives
// in ui/icon/icon_defs.go; search isn't in scope for v1.
var popular = []namedIcon{
	{"House", uiicon.House}, {"Settings", uiicon.Settings}, {"User", uiicon.User},
	{"Users", uiicon.Users}, {"Mail", uiicon.Mail}, {"Search", uiicon.Search},
	{"Heart", uiicon.Heart}, {"Star", uiicon.Star}, {"Trash2", uiicon.Trash2},
	{"Copy", uiicon.Copy}, {"Plus", uiicon.Plus}, {"Minus", uiicon.Minus},
	{"Check", uiicon.Check}, {"X", uiicon.X}, {"ChevronDown", uiicon.ChevronDown},
	{"ArrowRight", uiicon.ArrowRight}, {"Bell", uiicon.Bell}, {"Calendar", uiicon.Calendar},
	{"Camera", uiicon.Camera}, {"Clock", uiicon.Clock}, {"Cloud", uiicon.Cloud},
	{"Download", uiicon.Download}, {"Upload", uiicon.Upload}, {"Pencil", uiicon.Pencil},
	{"ExternalLink", uiicon.ExternalLink}, {"Eye", uiicon.Eye}, {"File", uiicon.File},
	{"Folder", uiicon.Folder}, {"Globe", uiicon.Globe}, {"Image", uiicon.Image},
	{"Info", uiicon.Info}, {"Link", uiicon.Link}, {"Lock", uiicon.Lock},
	{"LogOut", uiicon.LogOut}, {"Menu", uiicon.Menu}, {"MessageCircle", uiicon.MessageCircle},
	{"Moon", uiicon.Moon}, {"Sun", uiicon.Sun}, {"Play", uiicon.Play},
	{"ShoppingCart", uiicon.ShoppingCart},
}

// Gallery renders 40 well-known icons in a static grid — 1,748 icons ship
// in total.
component Gallery() {
	<div class="grid grid-cols-4 gap-6 sm:grid-cols-5 md:grid-cols-8">
		{ for _, ic := range popular {
			<div class="flex flex-col items-center gap-2">
				{ ic.Render(gsx.Attr{Key: "class", Value: "size-5"}) }
				<span class="text-xs text-muted-foreground">{ ic.Name }</span>
			</div>
		} }
	</div>
}
