package ui

import "github.com/gsxhq/gsx"

const base = "inline-flex shrink-0 items-center justify-center rounded-lg border border-transparent bg-clip-padding text-sm font-medium whitespace-nowrap transition-all outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 active:not-aria-[haspopup]:translate-y-px disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4"

func variantClass(variant string) string {
	switch variant {
	case "destructive":
		return "bg-destructive text-white hover:bg-destructive/90 focus-visible:ring-destructive/20 dark:bg-destructive/60 dark:focus-visible:ring-destructive/40"
	case "outline":
		return "border bg-background hover:bg-accent hover:text-accent-foreground dark:border-input dark:bg-input/30 dark:hover:bg-input/50"
	case "secondary":
		return "bg-secondary text-secondary-foreground hover:bg-secondary/80"
	case "ghost":
		return "hover:bg-accent hover:text-accent-foreground dark:hover:bg-accent/50"
	case "link":
		return "text-primary underline-offset-4 hover:underline"
	default:
		return "bg-primary text-primary-foreground hover:bg-primary/90"
	}
}

func sizeClass(size string) string {
	switch size {
	case "xs":
		return "h-6 gap-1 rounded-[min(var(--radius-md),10px)] in-data-[slot=button-group]:rounded-lg px-2 text-xs has-[>svg]:px-1.5 [&_svg:not([class*='size-'])]:size-3"
	case "sm":
		return "h-7 gap-1 rounded-[min(var(--radius-md),12px)] px-2.5 text-[0.8rem] has-[>svg]:px-1.5 [&_svg:not([class*='size-'])]:size-3.5"
	case "lg":
		return "h-9 gap-1.5 px-2.5 has-[>svg]:px-2"
	case "icon":
		return "size-8"
	case "icon-xs":
		return "size-6 rounded-[min(var(--radius-md),10px)] in-data-[slot=button-group]:rounded-lg [&_svg:not([class*='size-'])]:size-3"
	case "icon-sm":
		return "size-7 rounded-[min(var(--radius-md),12px)] in-data-[slot=button-group]:rounded-lg"
	case "icon-lg":
		return "size-9"
	default:
		return "h-8 gap-1.5 px-2.5 has-[>svg]:px-2"
	}
}

// Button is the shadcn/ui Button, retargeted to nova density (2026-07-24
// nova density map, `## button`). ADAPT: nova keys directional icon padding
// off `data-icon="inline-start|inline-end"` stamps we don't emit; gsxui kept
// its existing has-[>svg]:px-* selector mechanism and substituted nova's
// inline-start numeric value (e.g. default has-[>svg]:px-3 → px-2). All
// variants now carry a transparent 1px border in the base (box-size
// consistency across variants) — outline just recolors that border and no
// longer changes the box. A non-empty href on an enabled Button renders an
// <a> (gsx's answer to asChild-wrapping a link); disabled always renders a
// real disabled <button>. type="button" is an overridable default — pass
// type="submit" at the call site to submit forms.
component Button(variant string, size string, href string, disabled bool, children gsx.Node, attrs gsx.Attrs) {
	{ if href != "" && !disabled {
		<a
			data-slot="button"
			data-variant={variant |> default("default")}
			data-size={size |> default("default")}
			href={href}
			class={ base, variantClass(variant), sizeClass(size) }
			{ attrs... }
		>
			{ children }
		</a>
	} else {
		<button
			data-slot="button"
			data-variant={variant |> default("default")}
			data-size={size |> default("default")}
			type="button"
			class={ base, variantClass(variant), sizeClass(size) }
			disabled={disabled}
			{ attrs... }
		>
			{ children }
		</button>
	} }
}
