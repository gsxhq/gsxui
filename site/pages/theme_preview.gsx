package pages

import (
	"strings"

	"github.com/gsxhq/gsxui/internal/preset"
	"github.com/gsxhq/gsxui/site/stylepreview"
	"github.com/gsxhq/gsxui/site/stylepreview/maia"
	"github.com/gsxhq/gsxui/site/stylepreview/nova"
)

// ThemePreviewButton is the isolated same-origin document embedded by the
// theme editor. It renders copied consumer source, not the site's canonical
// recipe-class Button.
type ThemePreviewButton struct{}

func themePreviewTokenNames() string {
	return strings.Join(preset.TokenNames(), ",")
}

component (preview ThemePreviewButton) Page() {
	<!DOCTYPE html>
	<html lang="en">
		<siteHead title="Button theme preview"/>
		<body
			data-theme-button-preview
			data-theme-preview-tokens={themePreviewTokenNames()}
			class="min-h-svh overflow-auto bg-background text-foreground antialiased"
		>
			<main aria-label="Button style preview">
				<section data-theme-preview-style="nova">
					{ stylepreview.Matrix(nova.Button) }
				</section>
				<section data-theme-preview-style="maia" hidden>
					{ stylepreview.Matrix(maia.Button) }
				</section>
			</main>
		</body>
	</html>
}
