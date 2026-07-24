package carousel

import "github.com/gsxhq/gsxui/ui"

var apiSlides = []int{1, 2, 3, 4, 5}

// Api mirrors shadcn's own carousel-api.tsx: a "Slide X of Y" indicator
// below the carousel. shadcn drives it from embla's own setApi/api.on
// ("select") client state; this port has no equivalent client API object,
// so the indicator listens for the gsxui:carousel-select CustomEvent
// carousel.js emits on the carousel root instead — the direct proof case
// for that event's {index, count} contract (index is 0-based, matching
// gsxuiCarousel.scrollTo's own indexing, hence the +1 below). The initial
// text reads "Slide 1 of 5" without any JS having run yet, matching the
// same scrollLeft-starts-at-0 invariant CarouselPrevious's own initial
// disabled state relies on (see ui/carousel.gsx's package doc comment).
component Api() {
	<div class="mx-auto max-w-xs">
		<ui.Carousel id="api-carousel" orientation="" class="w-full max-w-xs">
			<ui.CarouselContent orientation="">
				{ for _, n := range apiSlides {
					<ui.CarouselItem orientation="">
						<ui.Card>
							<ui.CardContent class="flex aspect-square items-center justify-center p-6">
								<span class="text-4xl font-semibold">{ n }</span>
							</ui.CardContent>
						</ui.Card>
					</ui.CarouselItem>
				} }
			</ui.CarouselContent>
			<ui.CarouselPrevious orientation=""/>
			<ui.CarouselNext orientation=""/>
		</ui.Carousel>
		<output id="api-indicator" class="block py-2 text-center text-sm text-muted-foreground">Slide 1 of 5</output>
		<script>
			document.addEventListener("gsxui:carousel-select", (e) => {
				if (e.target.id !== "api-carousel") return;
				document.getElementById("api-indicator").textContent =
					"Slide " + (e.detail.index + 1) + " of " + e.detail.count;
			});
		</script>
	</div>
}
