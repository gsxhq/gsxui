package avatar

import "github.com/gsxhq/gsxui/site/uirtl"

// Rtl mirrors shadcn's own avatar-rtl demo shape: a loaded avatar, and an
// overlapping group with a "+N" count tile, wrapped in dir="rtl". gsxui
// has no AvatarBadge/AvatarGroup/AvatarGroupCount components (see
// ui/avatar.gsx), so the status-badge avatar is dropped and the group's
// count tile is approximated with a fourth Avatar whose fallback reads
// the Arabic-Indic "+٣" shadcn's own translations table uses for ar.
component Rtl() {
	<div dir="rtl" lang="ar" class="flex flex-row flex-wrap items-center gap-6">
		<uirtl.Avatar>
			<uirtl.AvatarImage src={avatarSVG |> dataURL("image/svg+xml")} alt="Ada Lovelace"/>
			<uirtl.AvatarFallback>AL</uirtl.AvatarFallback>
		</uirtl.Avatar>
		<div class="flex -space-x-2">
			<uirtl.Avatar class="ring-2 ring-background">
				<uirtl.AvatarImage src={avatarSVG |> dataURL("image/svg+xml")} alt="Ada Lovelace"/>
				<uirtl.AvatarFallback>AL</uirtl.AvatarFallback>
			</uirtl.Avatar>
			<uirtl.Avatar class="ring-2 ring-background">
				<uirtl.AvatarFallback>GH</uirtl.AvatarFallback>
			</uirtl.Avatar>
			<uirtl.Avatar class="ring-2 ring-background">
				<uirtl.AvatarFallback>AT</uirtl.AvatarFallback>
			</uirtl.Avatar>
			<uirtl.Avatar class="ring-2 ring-background">
				<uirtl.AvatarFallback>{ "+٣" }</uirtl.AvatarFallback>
			</uirtl.Avatar>
		</div>
	</div>
}
