package resizable

import "github.com/gsxhq/gsxui/ui"

// Handle mirrors shadcn's own resizable-handle.tsx: a single
// orientation="horizontal" group (Sidebar 25% / Content 75%) with the
// withHandle grip shown — the canonical "visible handle" demo.
component Handle() {
	<ui.ResizablePanelGroup orientation="horizontal" class="min-h-[200px] max-w-md rounded-lg border md:min-w-[450px]">
		<ui.ResizablePanel defaultSize="25%">
			<div class="flex h-full items-center justify-center p-6">
				<span class="font-semibold">Sidebar</span>
			</div>
		</ui.ResizablePanel>
		<ui.ResizableHandle orientation="horizontal" withHandle={true}/>
		<ui.ResizablePanel defaultSize="75%">
			<div class="flex h-full items-center justify-center p-6">
				<span class="font-semibold">Content</span>
			</div>
		</ui.ResizablePanel>
	</ui.ResizablePanelGroup>
}
