package scrollarea

import "github.com/gsxhq/gsxui/ui"

// artist names shadcn's own scroll-area-horizontal-demo.tsx's `works` array
// keys off (artist, art) pairs; only the artist survives here since `art`
// is a remote Unsplash image URL and the site has no remote-image assets
// (see the styled-div placeholder below, same "no external images" call as
// site/examples/avatar/basic.gsx's own inlined SVG stand-in).
var artists = []string{
	"Ornella Binni", "Tom Byrom", "Vladimir Malyavko",
	"Fahrul Azmi", "Christian Holzinger", "Marcus Loke",
	"Sindre Fs", "Jeremy Bishop",
}

// Horizontal mirrors shadcn's own scroll-area-horizontal-demo.tsx: a w-96
// bordered box, whitespace-nowrap, a flex w-max row of fixed-width
// "artwork" cards — orientation="horizontal" (overflow-x-auto) proves the
// axis switch actually works, per the source map's own demo-inventory call
// ("horizontal requires the explicit opt-in every time under this ADAPT's
// collapsed-single-div design, so both must be demoed").
component Horizontal() {
	<ui.ScrollArea orientation="horizontal" class="w-96 rounded-md border whitespace-nowrap">
		<div class="flex w-max gap-4 p-4">
			{ for _, artist := range artists {
				<figure class="w-[150px] shrink-0">
					<div class="flex aspect-[3/4] items-center justify-center rounded-md bg-muted text-xs text-muted-foreground">
						Photo
					</div>
					<figcaption class="pt-2 text-xs text-muted-foreground">
						Photo by <span class="font-semibold text-foreground">{ artist }</span>
					</figcaption>
				</figure>
			} }
		</div>
	</ui.ScrollArea>
}
