package tabs

import (
	"github.com/gsxhq/gsxui/ui"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Icons pairs each trigger with a lucide icon alongside its label.
component Icons() {
	<ui.Tabs value="grid">
		<ui.TabsList>
			<ui.TabsTrigger value="grid" selected><icon.LayoutGrid class="size-4"/>Grid</ui.TabsTrigger>
			<ui.TabsTrigger value="list"><icon.List class="size-4"/>List</ui.TabsTrigger>
		</ui.TabsList>
		<ui.TabsContent value="grid" selected>Items shown as a grid.</ui.TabsContent>
		<ui.TabsContent value="list">Items shown as a list.</ui.TabsContent>
	</ui.Tabs>
}
