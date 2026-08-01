package carousel

import "github.com/gsxhq/gsxui/ui"

var rtlSlides = []string{"١", "٢", "٣", "٤", "٥"}

// Rtl mirrors shadcn's own carousel-rtl demo: five slides numbered with
// Arabic-Indic numerals, wrapped in dir="rtl". gsxui's Carousel has no
// direction option to mirror shadcn's embla opts.direction (see
// ui/carousel.gsx); the dir="rtl" wrapper alone is what Basic's own
// pattern (and pagination/rtl.gsx) relies on for logical-property flip.
component Rtl() {
	<div dir="rtl" lang="ar" class="mx-auto w-full max-w-xs">
		<ui.Carousel orientation="" class="mx-auto w-full max-w-xs">
			<ui.CarouselContent orientation="">
				{ for _, n := range rtlSlides {
					<ui.CarouselItem orientation="">
						<div class="p-1">
							<ui.Card>
								<ui.CardContent class="flex aspect-square items-center justify-center p-6">
									<span class="text-4xl font-semibold">{ n }</span>
								</ui.CardContent>
							</ui.Card>
						</div>
					</ui.CarouselItem>
				} }
			</ui.CarouselContent>
			<ui.CarouselPrevious orientation=""/>
			<ui.CarouselNext orientation=""/>
		</ui.Carousel>
	</div>
}
