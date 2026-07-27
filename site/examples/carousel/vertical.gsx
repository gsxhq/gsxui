package carousel

import "github.com/gsxhq/gsxui/ui"

var verticalSlides = []int{1, 2, 3, 4, 5}

// Vertical mirrors shadcn's vertical carousel example. It is also the
// production proof that every Carousel part reflects the vertical
// orientation declared by the public style contract.
component Vertical() {
	<ui.Carousel orientation="vertical" class="mx-auto w-full max-w-xs">
		<ui.CarouselContent orientation="vertical" class="-mt-1 h-[200px]">
			{ for _, n := range verticalSlides {
				<ui.CarouselItem orientation="vertical" class="pt-1 -scroll-mt-1 md:basis-1/2">
					<div class="p-1">
						<ui.Card>
							<ui.CardContent class="flex items-center justify-center p-6">
								<span class="text-3xl font-semibold">{ n }</span>
							</ui.CardContent>
						</ui.Card>
					</div>
				</ui.CarouselItem>
			} }
		</ui.CarouselContent>
		<ui.CarouselPrevious orientation="vertical"/>
		<ui.CarouselNext orientation="vertical"/>
	</ui.Carousel>
}
