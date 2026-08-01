package sidebar

import (
	"github.com/gsxhq/gsxui/ui"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Rtl mirrors Basic's own shape (header, grouped menu, footer, trigger,
// inset content) but under dir="rtl" with side="right" — shadcn's own RTL
// convention keeps the physical side on the visual right so it still reads
// as inline-start in a right-to-left document. Everything else (icons,
// spacing, keyboard order) adapts through the same logical properties and
// rtl:rotate-180 rules Basic already exercises.
component Rtl() {
	<div dir="rtl">
		<ui.SidebarProvider open={true} class="min-h-[32rem] rounded-lg border">
			<ui.Sidebar open={true} side="right">
				<ui.SidebarHeader>
					<div class="px-2 py-1 text-sm font-semibold">شركة أكمي</div>
				</ui.SidebarHeader>
				<ui.SidebarSeparator/>
				<ui.SidebarContent>
					<ui.SidebarGroup>
						<ui.SidebarGroupLabel>التطبيق</ui.SidebarGroupLabel>
						<ui.SidebarGroupContent>
							<ui.SidebarMenu>
								<ui.SidebarMenuItem>
									<ui.SidebarMenuButton isActive={true} tooltip="الرئيسية">
										<icon.House/>
										<span>الرئيسية</span>
									</ui.SidebarMenuButton>
								</ui.SidebarMenuItem>
								<ui.SidebarMenuItem>
									<ui.SidebarMenuButton tooltip="البريد الوارد">
										<icon.Inbox/>
										<span>البريد الوارد</span>
									</ui.SidebarMenuButton>
								</ui.SidebarMenuItem>
								<ui.SidebarMenuItem>
									<ui.SidebarMenuButton tooltip="الإعدادات">
										<icon.Settings/>
										<span>الإعدادات</span>
									</ui.SidebarMenuButton>
								</ui.SidebarMenuItem>
							</ui.SidebarMenu>
						</ui.SidebarGroupContent>
					</ui.SidebarGroup>
				</ui.SidebarContent>
				<ui.SidebarRail/>
			</ui.Sidebar>
			<ui.SidebarInset>
				<header class="flex h-12 items-center gap-2 border-b px-4">
					<ui.SidebarTrigger/>
					<span class="text-sm text-muted-foreground">لوحة التحكم</span>
				</header>
				<div class="p-4 text-sm text-muted-foreground">
					بدّل الشريط الجانبي بالزر أعلاه، أو المقبض على حافته، أو Cmd/Ctrl+B.
				</div>
			</ui.SidebarInset>
		</ui.SidebarProvider>
	</div>
}
