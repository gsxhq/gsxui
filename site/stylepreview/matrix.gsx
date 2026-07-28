package stylepreview

import "github.com/gsxhq/gsx"

// ButtonRenderer is the copied Button contract shared by every generated
// style package. Matrix receives the constructor so its coverage cannot drift
// between styles.
type ButtonRenderer func(
	variant, size, href string,
	disabled bool,
	children gsx.Node,
	attrs gsx.Attrs,
) gsx.Node

type previewButtonCase struct {
	Label   string
	Variant string
	Size    string
}

var previewButtonVariants = []previewButtonCase{
	{Label: "Default", Variant: "default"},
	{Label: "Secondary", Variant: "secondary"},
	{Label: "Destructive", Variant: "destructive"},
	{Label: "Outline", Variant: "outline"},
	{Label: "Ghost", Variant: "ghost"},
	{Label: "Link", Variant: "link"},
}

var previewButtonSizes = []previewButtonCase{
	{Label: "Default", Size: "default"},
	{Label: "Extra small", Size: "xs"},
	{Label: "Small", Size: "sm"},
	{Label: "Large", Size: "lg"},
	{Label: "Icon", Size: "icon"},
	{Label: "Icon extra small", Size: "icon-xs"},
	{Label: "Icon small", Size: "icon-sm"},
	{Label: "Icon large", Size: "icon-lg"},
}

component previewPlusIcon() {
	<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden>
		<path d="M12 5v14M5 12h14"/>
	</svg>
}

component Matrix(button ButtonRenderer) {
	<div data-theme-preview-matrix class="mx-auto flex w-full max-w-4xl flex-col gap-8 p-6 sm:p-8">
		<section class="flex flex-col gap-3">
			<h2 class="text-sm font-medium text-muted-foreground">Variants</h2>
			<div class="flex flex-wrap items-center gap-3">
				{ for _, example := range previewButtonVariants {
					{ button(
						example.Variant,
						"default",
						"",
						false,
						gsx.Text(example.Label),
						gsx.Attrs{{Key: "data-theme-preview-case", Value: "text"}},
					) }
				} }
			</div>
		</section>

		<section class="flex flex-col gap-3">
			<h2 class="text-sm font-medium text-muted-foreground">Sizes</h2>
			<div class="flex flex-wrap items-center gap-3">
				{ for _, example := range previewButtonSizes {
					{ if example.Size == "icon" || example.Size == "icon-xs" || example.Size == "icon-sm" || example.Size == "icon-lg" {
						{ button(
							"default",
							example.Size,
							"",
							false,
							previewPlusIcon(),
							gsx.Attrs{
								{Key: "aria-label", Value: example.Label},
								{Key: "data-theme-preview-case", Value: "icon"},
							},
						) }
					} else {
						{ button(
							"default",
							example.Size,
							"",
							false,
							gsx.Text(example.Label),
							gsx.Attrs{{Key: "data-theme-preview-case", Value: "text"}},
						) }
					} }
				} }
			</div>
		</section>

		<section class="flex flex-col gap-3">
			<h2 class="text-sm font-medium text-muted-foreground">States and composition</h2>
			<div class="flex flex-wrap items-center gap-3">
				{ button(
					"default",
					"default",
					"",
					true,
					gsx.Text("Disabled"),
					gsx.Attrs{{Key: "data-theme-preview-case", Value: "disabled"}},
				) }
				{ button(
					"outline",
					"default",
					"/docs/getting-started",
					false,
					gsx.Text("Link"),
					gsx.Attrs{{Key: "data-theme-preview-case", Value: "link"}},
				) }
				{ button(
					"default",
					"default",
					"",
					false,
					gsx.Text("Focus visible"),
					gsx.Attrs{
						{Key: "autofocus", Value: true},
						{Key: "data-theme-preview-case", Value: "focus-visible"},
					},
				) }
				{ button(
					"default",
					"default",
					"",
					false,
					gsx.Text("Invalid"),
					gsx.Attrs{
						{Key: "aria-invalid", Value: "true"},
						{Key: "data-theme-preview-case", Value: "invalid"},
					},
				) }
				{ button(
					"default",
					"default",
					"",
					false,
					gsx.Text("Caller classes win"),
					gsx.Attrs{
						{Key: "class", Value: "rounded-none h-20 px-12 bg-warning text-warning-foreground rounded-r-none border-l-0"},
						{Key: "data-caller-marker", Value: true},
						{Key: "data-theme-preview-case", Value: "caller-composition"},
					},
				) }
			</div>
		</section>
	</div>
}
