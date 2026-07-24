package ui

import (
	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Carousel and its parts are the shadcn/ui Carousel (registry/new-york-v4/
// ui/carousel.tsx). shadcn's version wraps embla-carousel-react, a
// JS-transform carousel (it drags a flex track with
// `transform: translate3d(...)`, driven by its own pointer-physics engine —
// no CSS scroll-snap, no native `overflow: auto` anywhere). Per the roadmap's
// own decision (2026-07-24 controls source map, `## carousel`), this port
// substitutes real native CSS scroll-snap for embla entirely:
// `overflow-x-auto snap-x snap-mandatory` (or the `-y`/`snap-y` pair,
// vertical) on CarouselContent's own viewport div, `snap-start` on every
// CarouselItem, plain JS only for the prev/next buttons and disabled-state/
// current-index bookkeeping (ui/carousel.js). This is a genuine behavior
// substitution, not a strict downgrade: native scroll brings real touch/
// trackpad momentum and rubber-banding for free (a WIN embla's transform
// approach never had) but drops embla's `loop` (infinite wraparound) —
// GAP, no native-scroll analog without slide-cloning, and no docs demo
// exercises it, so it's an accepted v1 gap, not silently dropped.
//
// GAP (no context, orientation/spacing passed explicitly): Radix-less —
// there is no CarouselContext broadcasting `orientation` from Carousel down
// to CarouselContent/CarouselItem/CarouselPrevious/CarouselNext the way
// shadcn's own `useCarousel()` does. Every part that needs it takes its own
// `orientation` param instead, the same explicit-prop-instead-of-context
// shape `## toggle-group`'s groupType/variant/size/spacing already
// establishes. Spacing (embla's `-ml-4`/`pl-4` default gap, overridden to
// `-ml-1`/`pl-1` by `carousel-spacing.tsx`) stays caller-controlled via the
// ordinary class-merge mechanism — no separate spacing param, matching the
// upstream demo's own proof that gap is entirely a `className` override.
component Carousel(orientation string, children gsx.Node, attrs gsx.Attrs) {
	<div
		role="region"
		aria-roledescription="carousel"
		data-slot="carousel"
		data-gsxui-carousel
		data-orientation={orientation |> default("horizontal")}
		class="relative"
		{ attrs... }
	>
		{ children }
	</div>
}

// CarouselContent renders BOTH divs from shadcn's own source: the outer div
// is embla's `carouselRef` viewport target, ported here as the REAL native
// scroll container (`overflow-x-auto`/`-y-auto` + `snap-x`/`snap-y
// snap-mandatory`, replacing embla's own bare `overflow-hidden` — a scroll
// container needs `overflow: auto`, not `hidden`, to scroll at all); the
// inner div stays the plain flex track shadcn's own version already has,
// `{ attrs... }` landing there exactly as it does on shadcn's own inner div
// (`carousel-spacing.tsx`'s `-ml-1` override targets this same inner div).
//
// `[scrollbar-width:none] [&::-webkit-scrollbar]:hidden` on the outer div is
// NEW — present in neither shadcn's source nor nova's own CSS — needed only
// because this port is now on a real native scroll container that would
// otherwise show a visible scrollbar embla's transform-based approach never
// had anything analogous to.
component CarouselContent(orientation string, children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="carousel-content"
		class={
			if orientation == "vertical" {
				"overflow-y-auto snap-y snap-mandatory [scrollbar-width:none] [&::-webkit-scrollbar]:hidden"
			} else {
				"overflow-x-auto snap-x snap-mandatory [scrollbar-width:none] [&::-webkit-scrollbar]:hidden"
			},
		}
	>
		<div
			class={
				"flex",
				if orientation == "vertical" { "-mt-4 flex-col" } else { "-ml-4" },
			}
			{ attrs... }
		>
			{ children }
		</div>
	</div>
}

// CarouselItem adds `snap-start` to shadcn's own class string — also NEW,
// not in shadcn's source, required for native scroll-snap to have any snap
// points at all (embla needed none: it never scrolls, it transforms).
component CarouselItem(orientation string, children gsx.Node, attrs gsx.Attrs) {
	<div
		role="group"
		aria-roledescription="slide"
		data-slot="carousel-item"
		class={
			"min-w-0 shrink-0 grow-0 basis-full snap-start",
			if orientation == "vertical" { "pt-4" } else { "pl-4" },
		}
		{ attrs... }
	>
		{ children }
	</div>
}

// CarouselPrevious/CarouselNext compose Button (variant="outline"
// size="icon") exactly like shadcn's own versions, plus
// data-gsxui-carousel-prev/-next for carousel.js's delegated click wiring.
// shadcn computes `disabled={!canScrollPrev}`/`!canScrollNext` from embla's
// live scroll-progress state, unavailable at Go render time — the initial
// server-rendered `disabled` value is chosen per button since the two are
// NOT symmetric unknowns: a freshly mounted scroll container always starts
// at `scrollLeft`/`scrollTop` 0 (a real DOM invariant, not a guess), so
// Previous genuinely has nowhere to scroll back to and starts disabled;
// whether Next has anywhere to scroll forward TO depends on rendered
// content/viewport widths carousel.gsx cannot measure, so it starts enabled
// (the permissive default — a functional button that turns out to have
// nothing to do is harmless, a disabled one that turns out to be wrong is
// not). carousel.js's own init pass recomputes and corrects both from the
// real DOM immediately on load either way — see its own header comment.
component CarouselPrevious(orientation string, attrs gsx.Attrs) {
	<Button
		data-slot="carousel-previous"
		data-gsxui-carousel-prev
		variant="outline"
		size="icon"
		disabled={true}
		class={
			"absolute size-8 rounded-full",
			if orientation == "vertical" {
				"-top-12 left-1/2 -translate-x-1/2 rotate-90"
			} else {
				"top-1/2 -left-12 -translate-y-1/2"
			},
		}
		{ attrs... }
	>
		<icon.ArrowLeft/>
		<span class="sr-only">Previous slide</span>
	</Button>
}

component CarouselNext(orientation string, attrs gsx.Attrs) {
	<Button
		data-slot="carousel-next"
		data-gsxui-carousel-next
		variant="outline"
		size="icon"
		class={
			"absolute size-8 rounded-full",
			if orientation == "vertical" {
				"-bottom-12 left-1/2 -translate-x-1/2 rotate-90"
			} else {
				"top-1/2 -right-12 -translate-y-1/2"
			},
		}
		{ attrs... }
	>
		<icon.ArrowRight/>
		<span class="sr-only">Next slide</span>
	</Button>
}
