package command

import (
	"github.com/gsxhq/gsxui/site/uirtl"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Rtl mirrors shadcn's own command-rtl demo's Arabic strings, kept in
// Basic's shape: inline palette plus the CommandDialog variant.
component Rtl() {
	<div dir="rtl" lang="ar">
		<uirtl.Command class="max-w-md border shadow-md">
			<uirtl.CommandInput placeholder="اكتب أمرًا أو ابحث..."/>
			<uirtl.CommandList>
				<uirtl.CommandEmpty>لم يتم العثور على نتائج.</uirtl.CommandEmpty>
				<uirtl.CommandGroup heading="اقتراحات">
					<uirtl.CommandItem>
						<icon.Calendar/>
						<span>التقويم</span>
					</uirtl.CommandItem>
					<uirtl.CommandItem>
						<icon.Smile/>
						<span>البحث عن الرموز التعبيرية</span>
					</uirtl.CommandItem>
					<uirtl.CommandItem data-disabled="true">
						<icon.Calculator/>
						<span>الآلة الحاسبة</span>
					</uirtl.CommandItem>
				</uirtl.CommandGroup>
				<uirtl.CommandSeparator/>
				<uirtl.CommandGroup heading="الإعدادات">
					<uirtl.CommandItem>
						<icon.User/>
						<span>الملف الشخصي</span>
						<uirtl.CommandShortcut>⌘P</uirtl.CommandShortcut>
					</uirtl.CommandItem>
					<uirtl.CommandItem>
						<icon.CreditCard/>
						<span>الفوترة</span>
						<uirtl.CommandShortcut>⌘B</uirtl.CommandShortcut>
					</uirtl.CommandItem>
					<uirtl.CommandItem>
						<icon.Settings/>
						<span>الإعدادات</span>
						<uirtl.CommandShortcut>⌘S</uirtl.CommandShortcut>
					</uirtl.CommandItem>
				</uirtl.CommandGroup>
			</uirtl.CommandList>
		</uirtl.Command>
		<uirtl.CommandDialog title="لوحة الأوامر" description="البحث عن الأوامر">
			<uirtl.CommandInput placeholder="ابحث عن الأوامر..."/>
			<uirtl.CommandList>
				<uirtl.CommandItem>فتح الإعدادات</uirtl.CommandItem>
			</uirtl.CommandList>
		</uirtl.CommandDialog>
	</div>
}
