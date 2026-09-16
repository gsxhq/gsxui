package resizable

import "github.com/gsxhq/gsxui/site/uirtl"

// Rtl mirrors shadcn's own resizable-rtl demo: the same horizontal
// two-panel / nested-vertical-group shape as this dir's own Basic, labels
// translated to Arabic and wrapped in dir="rtl".
component Rtl() {
	<div dir="rtl" lang="ar">
		<uirtl.ResizablePanelGroup orientation="horizontal" class="max-w-sm rounded-lg border">
			<uirtl.ResizablePanel defaultSize="50%">
				<div class="flex h-[200px] items-center justify-center p-6">
					<span class="font-semibold">واحد</span>
				</div>
			</uirtl.ResizablePanel>
			<uirtl.ResizableHandle orientation="horizontal" withHandle={true}/>
			<uirtl.ResizablePanel defaultSize="50%">
				<uirtl.ResizablePanelGroup orientation="vertical">
					<uirtl.ResizablePanel defaultSize="25%">
						<div class="flex h-full items-center justify-center p-6">
							<span class="font-semibold">اثنان</span>
						</div>
					</uirtl.ResizablePanel>
					<uirtl.ResizableHandle orientation="vertical" withHandle={true}/>
					<uirtl.ResizablePanel defaultSize="75%">
						<div class="flex h-full items-center justify-center p-6">
							<span class="font-semibold">ثلاثة</span>
						</div>
					</uirtl.ResizablePanel>
				</uirtl.ResizablePanelGroup>
			</uirtl.ResizablePanel>
		</uirtl.ResizablePanelGroup>
	</div>
}
