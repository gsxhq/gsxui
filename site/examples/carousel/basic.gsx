package carousel

import "github.com/gsxhq/gsxui/ui"

var slides = []int{1, 2, 3, 4, 5}

// Basic mirrors shadcn's own carousel-demo.tsx: single item per view (no
// per-item basis override), five numbered slides each in a bordered
// aspect-square Card, the docs' own baseline look.
component Basic() {
	<ui.Carousel orientation="" class="mx-auto w-full max-w-xs">
		<ui.CarouselContent orientation="">
			{ for _, n := range slides {
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
}
