package pages

import (
	"strings"

	"github.com/gsxhq/gsxui/internal/preset"
	"github.com/gsxhq/gsxui/site/stylepreview/maia"
	"github.com/gsxhq/gsxui/site/stylepreview/nova"
)

// ThemePreview is the isolated same-origin document embedded by the theme
// editor. It renders the copied consumer source for the whole component
// catalogue — each style package's generated gallery composition — not the
// site's canonical recipe-class components.
type ThemePreview struct{}

func themePreviewTokenNames() string {
	return strings.Join(preset.TokenNames(), ",")
}

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
				<section data-theme-preview-style="nova">
					{ nova.Gallery("nova") }
				</section>
				<section data-theme-preview-style="maia" hidden>
					{ maia.Gallery("maia") }
				</section>
			</main>
		</body>
	</html>
}
