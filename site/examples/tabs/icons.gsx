package tabs

import (
	"github.com/gsxhq/gsxui/ui/icon"
	uitabs "github.com/gsxhq/gsxui/ui/tabs"
)

// Icons pairs each trigger with a lucide icon alongside its label.
component Icons() {
	<uitabs.Tabs value="grid">
		<uitabs.TabsList>
			<uitabs.TabsTrigger value="grid" selected><icon.LayoutGrid class="size-4"/>Grid</uitabs.TabsTrigger>
			<uitabs.TabsTrigger value="list"><icon.List class="size-4"/>List</uitabs.TabsTrigger>
		</uitabs.TabsList>
		<uitabs.TabsContent value="grid" selected>Items shown as a grid.</uitabs.TabsContent>
		<uitabs.TabsContent value="list">Items shown as a list.</uitabs.TabsContent>
	</uitabs.Tabs>
}
