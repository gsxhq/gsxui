package carousel

import "github.com/gsxhq/gsxui/site/uirtl"

var rtlSlides = []string{"١", "٢", "٣", "٤", "٥"}

// Rtl mirrors shadcn's own carousel-rtl demo: five slides numbered with
// Arabic-Indic numerals, wrapped in dir="rtl". gsxui's Carousel has no
// direction option to mirror shadcn's embla opts.direction (see
// ui/carousel.gsx); the dir="rtl" wrapper alone is what Basic's own
// pattern (and pagination/rtl.gsx) relies on for logical-property flip.
component Rtl() {
	<div dir="rtl" lang="ar" class="mx-auto w-full max-w-xs">
		<uirtl.Carousel orientation="" class="mx-auto w-full max-w-xs">
			<uirtl.CarouselContent orientation="">
				{ for _, n := range rtlSlides {
					<uirtl.CarouselItem orientation="">
						<div class="p-1">
							<uirtl.Card>
								<uirtl.CardContent class="flex aspect-square items-center justify-center p-6">
									<span class="text-4xl font-semibold">{ n }</span>
								</uirtl.CardContent>
							</uirtl.Card>
						</div>
					</uirtl.CarouselItem>
				} }
			</uirtl.CarouselContent>
			<uirtl.CarouselPrevious orientation=""/>
			<uirtl.CarouselNext orientation=""/>
		</uirtl.Carousel>
	</div>
}
