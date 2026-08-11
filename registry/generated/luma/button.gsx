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
				"focus-visible:border-ring focus-visible:ring-ring/30 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive dark:aria-invalid:border-destructive/50 rounded-4xl border border-transparent bg-clip-padding text-sm font-medium focus-visible:ring-3 aria-invalid:ring-3 active:not-aria-[haspopup]:translate-y-px [&_svg:not([class*='size-'])]:size-4 inline-flex shrink-0 items-center justify-center whitespace-nowrap transition-all outline-none select-none disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0",
				switch variant {
				case "destructive":
					"bg-destructive/10 hover:bg-destructive/20 focus-visible:ring-destructive/20 dark:focus-visible:ring-destructive/40 dark:bg-destructive/20 text-destructive focus-visible:border-destructive/40 dark:hover:bg-destructive/30"
				case "outline":
					"border-border bg-background dark:bg-transparent hover:bg-muted hover:text-foreground dark:hover:bg-input/30 aria-expanded:bg-muted aria-expanded:text-foreground"
				case "secondary":
					"bg-secondary text-secondary-foreground hover:bg-[color-mix(in_oklch,var(--secondary),var(--foreground)_5%)] aria-expanded:bg-secondary aria-expanded:text-secondary-foreground"
				case "ghost":
					"hover:bg-muted hover:text-foreground dark:hover:bg-muted/50 aria-expanded:bg-muted aria-expanded:text-foreground"
				case "link":
					"text-primary underline-offset-4 hover:underline"
				default:
					"bg-primary text-primary-foreground hover:bg-primary/80"
				},
				switch size {
				case "xs":
					"h-6 gap-1 px-2.5 text-xs has-data-[icon=inline-end]:pr-2 has-data-[icon=inline-start]:pl-2 [&_svg:not([class*='size-'])]:size-3"
				case "sm":
					"h-8 gap-1 px-3 has-data-[icon=inline-end]:pr-2 has-data-[icon=inline-start]:pl-2"
				case "lg":
					"h-10 gap-1.5 px-4 has-data-[icon=inline-end]:pr-3 has-data-[icon=inline-start]:pl-3"
				case "icon":
					"size-9"
				case "icon-xs":
					"size-6 [&_svg:not([class*='size-'])]:size-3"
				case "icon-sm":
					"size-8"
				case "icon-lg":
					"size-10"
				default:
					"h-9 gap-1.5 px-3 has-data-[icon=inline-end]:pr-2.5 has-data-[icon=inline-start]:pl-2.5"
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
				"focus-visible:border-ring focus-visible:ring-ring/30 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive dark:aria-invalid:border-destructive/50 rounded-4xl border border-transparent bg-clip-padding text-sm font-medium focus-visible:ring-3 aria-invalid:ring-3 active:not-aria-[haspopup]:translate-y-px [&_svg:not([class*='size-'])]:size-4 inline-flex shrink-0 items-center justify-center whitespace-nowrap transition-all outline-none select-none disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0",
				switch variant {
				case "destructive":
					"bg-destructive/10 hover:bg-destructive/20 focus-visible:ring-destructive/20 dark:focus-visible:ring-destructive/40 dark:bg-destructive/20 text-destructive focus-visible:border-destructive/40 dark:hover:bg-destructive/30"
				case "outline":
					"border-border bg-background dark:bg-transparent hover:bg-muted hover:text-foreground dark:hover:bg-input/30 aria-expanded:bg-muted aria-expanded:text-foreground"
				case "secondary":
					"bg-secondary text-secondary-foreground hover:bg-[color-mix(in_oklch,var(--secondary),var(--foreground)_5%)] aria-expanded:bg-secondary aria-expanded:text-secondary-foreground"
				case "ghost":
					"hover:bg-muted hover:text-foreground dark:hover:bg-muted/50 aria-expanded:bg-muted aria-expanded:text-foreground"
				case "link":
					"text-primary underline-offset-4 hover:underline"
				default:
					"bg-primary text-primary-foreground hover:bg-primary/80"
				},
				switch size {
				case "xs":
					"h-6 gap-1 px-2.5 text-xs has-data-[icon=inline-end]:pr-2 has-data-[icon=inline-start]:pl-2 [&_svg:not([class*='size-'])]:size-3"
				case "sm":
					"h-8 gap-1 px-3 has-data-[icon=inline-end]:pr-2 has-data-[icon=inline-start]:pl-2"
				case "lg":
					"h-10 gap-1.5 px-4 has-data-[icon=inline-end]:pr-3 has-data-[icon=inline-start]:pl-3"
				case "icon":
					"size-9"
				case "icon-xs":
					"size-6 [&_svg:not([class*='size-'])]:size-3"
				case "icon-sm":
					"size-8"
				case "icon-lg":
					"size-10"
				default:
					"h-9 gap-1.5 px-3 has-data-[icon=inline-end]:pr-2.5 has-data-[icon=inline-start]:pl-2.5"
				}
			}
			{ attrs... }
			data-gsxui-slot-button
		>
			{ children }
		</button>
	} }
}
