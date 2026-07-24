package resizable

import "github.com/gsxhq/gsxui/ui"

// Vertical mirrors shadcn's own resizable-vertical.tsx: a single
// orientation="vertical" group (Header 25% / Content 75%), no withHandle
// grip.
component Vertical() {
	<ui.ResizablePanelGroup orientation="vertical" class="min-h-[200px] max-w-md rounded-lg border md:min-w-[450px]">
		<ui.ResizablePanel defaultSize="25%">
			<div class="flex h-full items-center justify-center p-6">
				<span class="font-semibold">Header</span>
			</div>
		</ui.ResizablePanel>
		<ui.ResizableHandle orientation="vertical" withHandle={false}/>
		<ui.ResizablePanel defaultSize="75%">
			<div class="flex h-full items-center justify-center p-6">
				<span class="font-semibold">Content</span>
			</div>
		</ui.ResizablePanel>
	</ui.ResizablePanelGroup>
}
