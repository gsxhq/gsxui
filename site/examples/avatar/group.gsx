package avatar

import uiavatar "github.com/gsxhq/gsxui/ui/avatar"

// Group renders three overlapping Avatars — a stacked-avatars pattern
// built entirely from attrs (-space-x-2, a ring), no extra component.
component Group() {
	<div class="flex -space-x-2">
		<uiavatar.Avatar class="ring-2 ring-background">
			<uiavatar.AvatarImage src={avatarSVG} alt="Ada Lovelace"/>
			<uiavatar.AvatarFallback>AL</uiavatar.AvatarFallback>
		</uiavatar.Avatar>
		<uiavatar.Avatar class="ring-2 ring-background">
			<uiavatar.AvatarFallback>GH</uiavatar.AvatarFallback>
		</uiavatar.Avatar>
		<uiavatar.Avatar class="ring-2 ring-background">
			<uiavatar.AvatarFallback>AT</uiavatar.AvatarFallback>
		</uiavatar.Avatar>
	</div>
}
