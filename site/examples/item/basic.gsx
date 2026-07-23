// Package item holds the site's example gsx components for ui/item.
package item

import (
	"github.com/gsxhq/gsxui/ui"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Basic renders a standalone Item (media + content + actions row) followed
// by an ItemGroup of a couple items split by an ItemSeparator.
component Basic() {
	<div class="flex flex-col gap-6">
		<ui.Item variant="outline">
			<ui.ItemMedia variant="icon">
				<icon.Bell/>
			</ui.ItemMedia>
			<ui.ItemContent>
				<ui.ItemTitle>New comment on your post</ui.ItemTitle>
				<ui.ItemDescription>Sarah replied to "Launch week recap".</ui.ItemDescription>
			</ui.ItemContent>
			<ui.ItemActions>
				<ui.Button variant="outline" size="sm">View</ui.Button>
			</ui.ItemActions>
		</ui.Item>
		<ui.ItemGroup>
			<ui.Item>
				<ui.ItemMedia variant="icon">
					<icon.User/>
				</ui.ItemMedia>
				<ui.ItemContent>
					<ui.ItemTitle>Jamie Lee</ui.ItemTitle>
					<ui.ItemDescription>jamie@example.com</ui.ItemDescription>
				</ui.ItemContent>
				<ui.ItemActions>
					<ui.Button variant="ghost" size="sm">Remove</ui.Button>
				</ui.ItemActions>
			</ui.Item>
			<ui.ItemSeparator/>
			<ui.Item>
				<ui.ItemMedia variant="icon">
					<icon.User/>
				</ui.ItemMedia>
				<ui.ItemContent>
					<ui.ItemTitle>Alex Chen</ui.ItemTitle>
					<ui.ItemDescription>alex@example.com</ui.ItemDescription>
				</ui.ItemContent>
				<ui.ItemActions>
					<ui.Button variant="ghost" size="sm">Remove</ui.Button>
				</ui.ItemActions>
			</ui.Item>
		</ui.ItemGroup>
	</div>
}
