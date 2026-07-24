package carousel

import "github.com/gsxhq/gsxui/ui"

var sizeSlides = []int{1, 2, 3, 4, 5}

// Sizes mirrors shadcn's own carousel-size.tsx combined with
// carousel-spacing.tsx (the map's own demo-inventory note picks one of the
// two near-duplicates; this folds both features into a single example):
// multi-per-view responsive basis (two per view at md, three at lg) plus a
// tightened gap (CarouselContent's own -ml-1 override, each item's own
// pl-1) proving the gap is entirely caller-controlled via class merge, no
// separate spacing prop.
component Sizes() {
	<ui.Carousel orientation="" class="mx-auto w-full max-w-sm">
		<ui.CarouselContent orientation="" class="-ml-1">
			{ for _, n := range sizeSlides {
				<ui.CarouselItem orientation="" class="pl-1 -scroll-ml-1 md:basis-1/2 lg:basis-1/3">
					<div class="p-1">
						<ui.Card>
							<ui.CardContent class="flex aspect-square items-center justify-center p-6">
								<span class="text-2xl font-semibold">{ n }</span>
							</ui.CardContent>
						</ui.Card>
					</div>
				</ui.CarouselItem>
			} }
		</ui.CarouselContent>
		<ui.CarouselPrevious orientation=""/>
		<ui.CarouselNext orientation=""/>
	</ui.Carousel>
}
