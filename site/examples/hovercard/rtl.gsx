package hovercard

import "github.com/gsxhq/gsxui/ui"

// Rtl mirrors Basic's own hover-card-demo shape (link trigger, avatar +
// name + description + joined-date card), translated to Arabic. shadcn's
// own hover-card-rtl demo is a product-price card with four physical-side
// and two logical-side triggers, which this directory's Basic doesn't
// compose; DEVIATED from rather than invented.
component Rtl() {
	<div dir="rtl" lang="ar">
		<ui.HoverCard>
			<ui.HoverCardTrigger>
				<ui.Button variant="link">@nextjs</ui.Button>
			</ui.HoverCardTrigger>
			<ui.HoverCardContent class="w-80">
				<div class="flex justify-between gap-4">
					<ui.Avatar>
						<ui.AvatarImage src={avatarSVG |> dataURL("image/svg+xml")} alt="@nextjs"/>
						<ui.AvatarFallback>VC</ui.AvatarFallback>
					</ui.Avatar>
					<div class="space-y-1">
						<h4 class="text-sm font-semibold">@nextjs</h4>
						<p class="text-sm">إطار React – تم إنشاؤه وصيانته بواسطة @vercel.</p>
						<div class="text-xs text-muted-foreground">انضم في ديسمبر 2021</div>
					</div>
				</div>
			</ui.HoverCardContent>
		</ui.HoverCard>
	</div>
}
