package avatar

import "github.com/gsxhq/gsxui/ui"

// Rtl mirrors shadcn's own avatar-rtl demo shape: a loaded avatar, and an
// overlapping group with a "+N" count tile, wrapped in dir="rtl". gsxui
// has no AvatarBadge/AvatarGroup/AvatarGroupCount components (see
// ui/avatar.gsx), so the status-badge avatar is dropped and the group's
// count tile is approximated with a fourth Avatar whose fallback reads
// the Arabic-Indic "+٣" shadcn's own translations table uses for ar.
component Rtl() {
	<div dir="rtl" lang="ar" class="flex flex-row flex-wrap items-center gap-6">
		<ui.Avatar>
			<ui.AvatarImage src={avatarSVG |> dataURL("image/svg+xml")} alt="Ada Lovelace"/>
			<ui.AvatarFallback>AL</ui.AvatarFallback>
		</ui.Avatar>
		<div class="flex -space-x-2">
			<ui.Avatar class="ring-2 ring-background">
				<ui.AvatarImage src={avatarSVG |> dataURL("image/svg+xml")} alt="Ada Lovelace"/>
				<ui.AvatarFallback>AL</ui.AvatarFallback>
			</ui.Avatar>
			<ui.Avatar class="ring-2 ring-background">
				<ui.AvatarFallback>GH</ui.AvatarFallback>
			</ui.Avatar>
			<ui.Avatar class="ring-2 ring-background">
				<ui.AvatarFallback>AT</ui.AvatarFallback>
			</ui.Avatar>
			<ui.Avatar class="ring-2 ring-background">
				<ui.AvatarFallback>{ "+٣" }</ui.AvatarFallback>
			</ui.Avatar>
		</div>
	</div>
}
