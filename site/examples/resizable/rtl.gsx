package resizable

import "github.com/gsxhq/gsxui/ui"

// Rtl mirrors shadcn's own resizable-rtl demo: the same horizontal
// two-panel / nested-vertical-group shape as this dir's own Basic, labels
// translated to Arabic and wrapped in dir="rtl".
component Rtl() {
	<div dir="rtl" lang="ar">
		<ui.ResizablePanelGroup orientation="horizontal" class="max-w-sm rounded-lg border">
			<ui.ResizablePanel defaultSize="50%">
				<div class="flex h-[200px] items-center justify-center p-6">
					<span class="font-semibold">واحد</span>
				</div>
			</ui.ResizablePanel>
			<ui.ResizableHandle orientation="horizontal" withHandle={true}/>
			<ui.ResizablePanel defaultSize="50%">
				<ui.ResizablePanelGroup orientation="vertical">
					<ui.ResizablePanel defaultSize="25%">
						<div class="flex h-full items-center justify-center p-6">
							<span class="font-semibold">اثنان</span>
						</div>
					</ui.ResizablePanel>
					<ui.ResizableHandle orientation="vertical" withHandle={true}/>
					<ui.ResizablePanel defaultSize="75%">
						<div class="flex h-full items-center justify-center p-6">
							<span class="font-semibold">ثلاثة</span>
						</div>
					</ui.ResizablePanel>
				</ui.ResizablePanelGroup>
			</ui.ResizablePanel>
		</ui.ResizablePanelGroup>
	</div>
}
