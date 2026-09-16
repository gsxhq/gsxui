// Package navigationmenu holds the site's example gsx components for
// ui/navigation-menu.
package navigationmenu

import (
	"github.com/gsxhq/gsxui/site/uirtl"
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
		<uirtl.NavigationMenu>
			<uirtl.NavigationMenuList>
				<uirtl.NavigationMenuItem>
					<uirtl.NavigationMenuTrigger>البدء</uirtl.NavigationMenuTrigger>
					<uirtl.NavigationMenuContent>
						<div class="grid w-96 gap-2">
							<uirtl.NavigationMenuLink href="#">
								<div class="text-sm font-medium">مقدمة</div>
								<div class="text-muted-foreground">مكونات قابلة لإعادة الاستخدام مبنية باستخدام Tailwind CSS.</div>
							</uirtl.NavigationMenuLink>
							<uirtl.NavigationMenuLink href="#">
								<div class="text-sm font-medium">التثبيت</div>
								<div class="text-muted-foreground">كيفية تثبيت التبعيات وتنظيم تطبيقك.</div>
							</uirtl.NavigationMenuLink>
							<uirtl.NavigationMenuLink href="#">
								<div class="text-sm font-medium">الطباعة</div>
								<div class="text-muted-foreground">أنماط للعناوين والفقرات والقوائم...إلخ</div>
							</uirtl.NavigationMenuLink>
						</div>
					</uirtl.NavigationMenuContent>
				</uirtl.NavigationMenuItem>
				<uirtl.NavigationMenuItem>
					<uirtl.NavigationMenuTrigger>مع أيقونة</uirtl.NavigationMenuTrigger>
					<uirtl.NavigationMenuContent>
						<div class="grid w-48 gap-1">
							<uirtl.NavigationMenuLink href="#" class="flex-row items-center gap-2">
								<icon.CircleAlert/>
								قائمة الانتظار
							</uirtl.NavigationMenuLink>
							<uirtl.NavigationMenuLink href="#" class="flex-row items-center gap-2">
								<icon.CircleDashed/>
								المهام
							</uirtl.NavigationMenuLink>
							<uirtl.NavigationMenuLink href="#" class="flex-row items-center gap-2">
								<icon.CircleCheck/>
								منجز
							</uirtl.NavigationMenuLink>
						</div>
					</uirtl.NavigationMenuContent>
				</uirtl.NavigationMenuItem>
				<uirtl.NavigationMenuItem>
					<uirtl.NavigationMenuLink variant="trigger" href="#">الوثائق</uirtl.NavigationMenuLink>
				</uirtl.NavigationMenuItem>
				<uirtl.NavigationMenuIndicator/>
			</uirtl.NavigationMenuList>
		</uirtl.NavigationMenu>
	</div>
}
