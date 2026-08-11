package pages

import (
	"strings"

	"github.com/gsxhq/gsxui/internal/preset"
	"github.com/gsxhq/gsxui/site/stylepreview/luma"
	"github.com/gsxhq/gsxui/site/stylepreview/lyra"
	"github.com/gsxhq/gsxui/site/stylepreview/maia"
	"github.com/gsxhq/gsxui/site/stylepreview/mira"
	"github.com/gsxhq/gsxui/site/stylepreview/nova"
	"github.com/gsxhq/gsxui/site/stylepreview/rhea"
	"github.com/gsxhq/gsxui/site/stylepreview/sera"
	"github.com/gsxhq/gsxui/site/stylepreview/vega"
)

// ThemePreview is the isolated same-origin document embedded by the theme
// editor. It renders the copied consumer source for the whole component
// catalogue — each style package's generated gallery composition — not the
// site's canonical recipe-class components.
type ThemePreview struct{}

func themePreviewTokenNames() string {
	return strings.Join(preset.TokenNames(), ",")
}

// component (preview ThemePreview) Page renders all 8 style sections into
// one same-origin document; web/theme-preview.js toggles which one is
// visible entirely by DOM query (`[data-theme-preview-style]`), so this
// list is the only place a newly ported style needs to be added — nothing
// in the JS layer names a style directly.
component (preview ThemePreview) Page() {
	<!DOCTYPE html>
	<html lang="en">
		<siteHead title="Theme preview" entry="web/preview.js"/>
		<body
			data-theme-button-preview
			data-theme-preview-tokens={themePreviewTokenNames()}
			class="min-h-svh overflow-auto bg-background text-foreground antialiased"
		>
			<main aria-label="Component gallery preview">
				<section data-theme-preview-style="vega" hidden>
					{ vega.Gallery("vega") }
				</section>
				<section data-theme-preview-style="nova">
					{ nova.Gallery("nova") }
				</section>
				<section data-theme-preview-style="maia" hidden>
					{ maia.Gallery("maia") }
				</section>
				<section data-theme-preview-style="lyra" hidden>
					{ lyra.Gallery("lyra") }
				</section>
				<section data-theme-preview-style="mira" hidden>
					{ mira.Gallery("mira") }
				</section>
				<section data-theme-preview-style="luma" hidden>
					{ luma.Gallery("luma") }
				</section>
				<section data-theme-preview-style="sera" hidden>
					{ sera.Gallery("sera") }
				</section>
				<section data-theme-preview-style="rhea" hidden>
					{ rhea.Gallery("rhea") }
				</section>
			</main>
		</body>
	</html>
}
