package sidebar

import (
	"github.com/gsxhq/gsxui/site/uirtl"
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
		<uirtl.SidebarProvider open={true} class="min-h-[32rem] rounded-lg border">
			<uirtl.Sidebar open={true} side="right">
				<uirtl.SidebarHeader>
					<div class="px-2 py-1 text-sm font-semibold">شركة أكمي</div>
				</uirtl.SidebarHeader>
				<uirtl.SidebarSeparator/>
				<uirtl.SidebarContent>
					<uirtl.SidebarGroup>
						<uirtl.SidebarGroupLabel>التطبيق</uirtl.SidebarGroupLabel>
						<uirtl.SidebarGroupContent>
							<uirtl.SidebarMenu>
								<uirtl.SidebarMenuItem>
									<uirtl.SidebarMenuButton isActive={true} tooltip="الرئيسية">
										<icon.House/>
										<span>الرئيسية</span>
									</uirtl.SidebarMenuButton>
								</uirtl.SidebarMenuItem>
								<uirtl.SidebarMenuItem>
									<uirtl.SidebarMenuButton tooltip="البريد الوارد">
										<icon.Inbox/>
										<span>البريد الوارد</span>
									</uirtl.SidebarMenuButton>
								</uirtl.SidebarMenuItem>
								<uirtl.SidebarMenuItem>
									<uirtl.SidebarMenuButton tooltip="الإعدادات">
										<icon.Settings/>
										<span>الإعدادات</span>
									</uirtl.SidebarMenuButton>
								</uirtl.SidebarMenuItem>
							</uirtl.SidebarMenu>
						</uirtl.SidebarGroupContent>
					</uirtl.SidebarGroup>
				</uirtl.SidebarContent>
				<uirtl.SidebarRail/>
			</uirtl.Sidebar>
			<uirtl.SidebarInset>
				<header class="flex h-12 items-center gap-2 border-b px-4">
					<uirtl.SidebarTrigger/>
					<span class="text-sm text-muted-foreground">لوحة التحكم</span>
				</header>
				<div class="p-4 text-sm text-muted-foreground">
					بدّل الشريط الجانبي بالزر أعلاه، أو المقبض على حافته، أو Cmd/Ctrl+B.
				</div>
			</uirtl.SidebarInset>
		</uirtl.SidebarProvider>
	</div>
}
