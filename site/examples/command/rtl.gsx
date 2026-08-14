package command

import (
	"github.com/gsxhq/gsxui/ui"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Rtl mirrors shadcn's own command-rtl demo's Arabic strings, kept in
// Basic's shape: inline palette plus the CommandDialog variant.
component Rtl() {
	<div dir="rtl" lang="ar">
		<ui.Command class="max-w-md border shadow-md">
			<ui.CommandInput placeholder="اكتب أمرًا أو ابحث..."/>
			<ui.CommandList>
				<ui.CommandEmpty>لم يتم العثور على نتائج.</ui.CommandEmpty>
				<ui.CommandGroup heading="اقتراحات">
					<ui.CommandItem>
						<icon.Calendar/>
						<span>التقويم</span>
					</ui.CommandItem>
					<ui.CommandItem>
						<icon.Smile/>
						<span>البحث عن الرموز التعبيرية</span>
					</ui.CommandItem>
					<ui.CommandItem data-disabled="true">
						<icon.Calculator/>
						<span>الآلة الحاسبة</span>
					</ui.CommandItem>
				</ui.CommandGroup>
				<ui.CommandSeparator/>
				<ui.CommandGroup heading="الإعدادات">
					<ui.CommandItem>
						<icon.User/>
						<span>الملف الشخصي</span>
						<ui.CommandShortcut>⌘P</ui.CommandShortcut>
					</ui.CommandItem>
					<ui.CommandItem>
						<icon.CreditCard/>
						<span>الفوترة</span>
						<ui.CommandShortcut>⌘B</ui.CommandShortcut>
					</ui.CommandItem>
					<ui.CommandItem>
						<icon.Settings/>
						<span>الإعدادات</span>
						<ui.CommandShortcut>⌘S</ui.CommandShortcut>
					</ui.CommandItem>
				</ui.CommandGroup>
			</ui.CommandList>
		</ui.Command>
		<ui.CommandDialog title="لوحة الأوامر" description="البحث عن الأوامر">
			<ui.CommandInput placeholder="ابحث عن الأوامر..."/>
			<ui.CommandList>
				<ui.CommandItem>فتح الإعدادات</ui.CommandItem>
			</ui.CommandList>
		</ui.CommandDialog>
	</div>
}
