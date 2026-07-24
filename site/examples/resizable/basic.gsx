package resizable

import "github.com/gsxhq/gsxui/ui"

// Basic mirrors shadcn's own resizable-demo-with-handle.tsx: a horizontal
// two-panel group (50%/50%) whose right panel nests a second, vertical
// group (25%/75%) — both the outer and the nested handle show the
// withHandle grip.
component Basic() {
	<ui.ResizablePanelGroup orientation="horizontal" class="max-w-md rounded-lg border md:min-w-[450px]">
		<ui.ResizablePanel defaultSize="50%">
			<div class="flex h-[200px] items-center justify-center p-6">
				<span class="font-semibold">One</span>
			</div>
		</ui.ResizablePanel>
		<ui.ResizableHandle orientation="horizontal" withHandle={true}/>
		<ui.ResizablePanel defaultSize="50%">
			<ui.ResizablePanelGroup orientation="vertical">
				<ui.ResizablePanel defaultSize="25%">
					<div class="flex h-full items-center justify-center p-6">
						<span class="font-semibold">Two</span>
					</div>
				</ui.ResizablePanel>
				<ui.ResizableHandle orientation="vertical" withHandle={true}/>
				<ui.ResizablePanel defaultSize="75%">
					<div class="flex h-full items-center justify-center p-6">
						<span class="font-semibold">Three</span>
					</div>
				</ui.ResizablePanel>
			</ui.ResizablePanelGroup>
		</ui.ResizablePanel>
	</ui.ResizablePanelGroup>
}
