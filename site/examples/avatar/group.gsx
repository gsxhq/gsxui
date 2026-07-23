package avatar

import "github.com/gsxhq/gsxui/ui"

// Group renders three overlapping Avatars — a stacked-avatars pattern
// built entirely from attrs (-space-x-2, a ring), no extra component.
component Group() {
	<div class="flex -space-x-2">
		<ui.Avatar class="ring-2 ring-background">
			<ui.AvatarImage src={avatarSVG} alt="Ada Lovelace"/>
			<ui.AvatarFallback>AL</ui.AvatarFallback>
		</ui.Avatar>
		<ui.Avatar class="ring-2 ring-background">
			<ui.AvatarFallback>GH</ui.AvatarFallback>
		</ui.Avatar>
		<ui.Avatar class="ring-2 ring-background">
			<ui.AvatarFallback>AT</ui.AvatarFallback>
		</ui.Avatar>
	</div>
}
