// Package navigationmenu holds the site's example gsx components for
// ui/navigation-menu.
package navigationmenu

import (
	"github.com/gsxhq/gsxui/ui"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Rtl mirrors shadcn's own navigation-menu-rtl demo: a "Getting started"
// trigger opening a list of doc links, a "With Icon" trigger opening a
// small grid of icon links, and a plain Docs link — translated to Arabic
// and wrapped in dir="rtl". shadcn's rtl demo also composes a "Components"
// trigger reusing the same six primitives as its own non-rtl demo; this
// dir's Basic already covers that shape (see basic.gsx), so it's dropped
// here rather than duplicated.
component Rtl() {
	<div dir="rtl" lang="ar">
		<ui.NavigationMenu>
			<ui.NavigationMenuList>
				<ui.NavigationMenuItem>
					<ui.NavigationMenuTrigger>البدء</ui.NavigationMenuTrigger>
					<ui.NavigationMenuContent>
						<div class="grid w-96 gap-2">
							<ui.NavigationMenuLink href="#">
								<div class="text-sm font-medium">مقدمة</div>
								<div class="text-muted-foreground">مكونات قابلة لإعادة الاستخدام مبنية باستخدام Tailwind CSS.</div>
							</ui.NavigationMenuLink>
							<ui.NavigationMenuLink href="#">
								<div class="text-sm font-medium">التثبيت</div>
								<div class="text-muted-foreground">كيفية تثبيت التبعيات وتنظيم تطبيقك.</div>
							</ui.NavigationMenuLink>
							<ui.NavigationMenuLink href="#">
								<div class="text-sm font-medium">الطباعة</div>
								<div class="text-muted-foreground">أنماط للعناوين والفقرات والقوائم...إلخ</div>
							</ui.NavigationMenuLink>
						</div>
					</ui.NavigationMenuContent>
				</ui.NavigationMenuItem>
				<ui.NavigationMenuItem>
					<ui.NavigationMenuTrigger>مع أيقونة</ui.NavigationMenuTrigger>
					<ui.NavigationMenuContent>
						<div class="grid w-48 gap-1">
							<ui.NavigationMenuLink href="#" class="flex-row items-center gap-2">
								<icon.CircleAlert/>
								قائمة الانتظار
							</ui.NavigationMenuLink>
							<ui.NavigationMenuLink href="#" class="flex-row items-center gap-2">
								<icon.CircleDashed/>
								المهام
							</ui.NavigationMenuLink>
							<ui.NavigationMenuLink href="#" class="flex-row items-center gap-2">
								<icon.CircleCheck/>
								منجز
							</ui.NavigationMenuLink>
						</div>
					</ui.NavigationMenuContent>
				</ui.NavigationMenuItem>
				<ui.NavigationMenuItem>
					<ui.NavigationMenuLink variant="trigger" href="#">الوثائق</ui.NavigationMenuLink>
				</ui.NavigationMenuItem>
				<ui.NavigationMenuIndicator/>
			</ui.NavigationMenuList>
		</ui.NavigationMenu>
	</div>
}
