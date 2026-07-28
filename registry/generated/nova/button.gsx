package ui

import "github.com/gsxhq/gsx"

// Button is the canonical Button structure. Its semantic recipe roles resolve
// to concrete style source before consumers compile the generated component.
// A non-empty href on an enabled Button renders an <a> (gsx's answer to
// asChild-wrapping a link); disabled always renders a real disabled <button>.
// type="button" is an overridable default — pass type="submit" at the call
// site to submit forms.
component Button(variant string, size string, href string, disabled bool, children gsx.Node, attrs gsx.Attrs) {
	{ if href != "" && !disabled {
		<a
			data-variant={variant |> default("default")}
			data-size={size |> default("default")}
			href={href}
			class={
				"group/button",
				"inline-flex shrink-0 items-center justify-center rounded-lg border border-transparent bg-clip-padding text-sm font-medium whitespace-nowrap transition-all outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 active:not-aria-[haspopup]:translate-y-px disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 [&_svg]:shrink-0 [&_svg]:pointer-events-none [&_svg:not([class*='size-'])]:size-4",
				switch variant {
				case "", "default":
					"bg-primary text-primary-foreground hover:bg-primary/90"
				case "destructive":
					"bg-destructive text-contrast hover:bg-destructive/90 focus-visible:ring-destructive/20 dark:bg-destructive/60 dark:focus-visible:ring-destructive/40"
				case "outline":
					"border-border bg-background hover:bg-accent hover:text-accent-foreground dark:border-input dark:bg-input/30 dark:hover:bg-input/50"
				case "secondary":
					"bg-secondary text-secondary-foreground hover:bg-secondary/80"
				case "ghost":
					"hover:bg-accent hover:text-accent-foreground dark:hover:bg-accent/50"
				case "link":
					"text-primary underline-offset-4 hover:underline"
				default:
					""
				},
				switch size {
				case "", "default":
					"h-8 gap-1.5 px-2.5 has-[>svg]:px-2"
				case "xs":
					"h-6 gap-1 rounded-[min(var(--radius-md),10px)] px-2 text-xs has-[>svg]:px-1.5 [&_svg:not([class*='size-'])]:size-3"
				case "sm":
					"h-7 gap-1 rounded-[min(var(--radius-md),12px)] px-2.5 text-[0.8rem] has-[>svg]:px-1.5 [&_svg:not([class*='size-'])]:size-3.5"
				case "lg":
					"h-9 gap-1.5 px-2.5 has-[>svg]:px-2"
				case "icon":
					"size-8"
				case "icon-xs":
					"size-6 rounded-[min(var(--radius-md),10px)] [&_svg:not([class*='size-'])]:size-3"
				case "icon-sm":
					"size-7 rounded-[min(var(--radius-md),12px)]"
				case "icon-lg":
					"size-9"
				default:
					""
				}
			}
			{ attrs... }
			data-gsxui-slot-button
		>
			{ children }
		</a>
	} else {
		<button
			data-variant={variant |> default("default")}
			data-size={size |> default("default")}
			type="button"
			disabled={disabled}
			class={
				"group/button",
				"inline-flex shrink-0 items-center justify-center rounded-lg border border-transparent bg-clip-padding text-sm font-medium whitespace-nowrap transition-all outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 active:not-aria-[haspopup]:translate-y-px disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 [&_svg]:shrink-0 [&_svg]:pointer-events-none [&_svg:not([class*='size-'])]:size-4",
				switch variant {
				case "", "default":
					"bg-primary text-primary-foreground hover:bg-primary/90"
				case "destructive":
					"bg-destructive text-contrast hover:bg-destructive/90 focus-visible:ring-destructive/20 dark:bg-destructive/60 dark:focus-visible:ring-destructive/40"
				case "outline":
					"border-border bg-background hover:bg-accent hover:text-accent-foreground dark:border-input dark:bg-input/30 dark:hover:bg-input/50"
				case "secondary":
					"bg-secondary text-secondary-foreground hover:bg-secondary/80"
				case "ghost":
					"hover:bg-accent hover:text-accent-foreground dark:hover:bg-accent/50"
				case "link":
					"text-primary underline-offset-4 hover:underline"
				default:
					""
				},
				switch size {
				case "", "default":
					"h-8 gap-1.5 px-2.5 has-[>svg]:px-2"
				case "xs":
					"h-6 gap-1 rounded-[min(var(--radius-md),10px)] px-2 text-xs has-[>svg]:px-1.5 [&_svg:not([class*='size-'])]:size-3"
				case "sm":
					"h-7 gap-1 rounded-[min(var(--radius-md),12px)] px-2.5 text-[0.8rem] has-[>svg]:px-1.5 [&_svg:not([class*='size-'])]:size-3.5"
				case "lg":
					"h-9 gap-1.5 px-2.5 has-[>svg]:px-2"
				case "icon":
					"size-8"
				case "icon-xs":
					"size-6 rounded-[min(var(--radius-md),10px)] [&_svg:not([class*='size-'])]:size-3"
				case "icon-sm":
					"size-7 rounded-[min(var(--radius-md),12px)]"
				case "icon-lg":
					"size-9"
				default:
					""
				}
			}
			{ attrs... }
			data-gsxui-slot-button
		>
			{ children }
		</button>
	} }
}
