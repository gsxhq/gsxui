package tooltip

import "github.com/gsxhq/gsxui/ui"

// Rtl mirrors shadcn's own tooltip-rtl demo's translated trigger/content
// pair, wrapped in dir="rtl". shadcn's rtl demo also sweeps every
// side ("left"/"top"/"bottom"/"right"/"inline-start"/"inline-end") across
// six Tooltip instances; this port's ui.TooltipContent has no side prop at
// all — its placement is fixed to top (see ui/tooltip.gsx's own doc
// comment: "its arrow is static because placement is always top") — so
// only the single Basic-shaped trigger/content pair is reproduced here
// (DEVIATION, no side sweep).
component Rtl() {
	<div dir="rtl" lang="ar">
		<ui.Tooltip>
			<ui.Button
				variant="outline"
				data-gsxui-slot-tooltip-trigger
			>
				حوم فوقي
			</ui.Button>
			<ui.TooltipContent>إضافة إلى المكتبة</ui.TooltipContent>
		</ui.Tooltip>
	</div>
}
