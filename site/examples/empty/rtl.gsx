package empty

import (
	"github.com/gsxhq/gsxui/site/uirtl"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Rtl mirrors shadcn's own empty-rtl demo's Arabic "no projects yet" copy,
// kept in Basic's shape (icon.Inbox media, title, description, one
// EmptyContent Button — shadcn's version swaps in a folder icon and a
// second outline Button plus a "learn more" link, which this directory's
// Basic doesn't compose; DEVIATED from rather than invented).
component Rtl() {
	<div dir="rtl" lang="ar">
		<uirtl.Empty>
			<uirtl.EmptyHeader>
				<uirtl.EmptyMedia>الوسائط الافتراضية</uirtl.EmptyMedia>
				<uirtl.EmptyMedia variant="icon">
					<icon.Inbox/>
				</uirtl.EmptyMedia>
				<uirtl.EmptyTitle>لا توجد مشاريع بعد</uirtl.EmptyTitle>
				<uirtl.EmptyDescription>
					لم تقم بإنشاء أي مشاريع بعد. ابدأ بإنشاء مشروعك الأول.
				</uirtl.EmptyDescription>
			</uirtl.EmptyHeader>
			<uirtl.EmptyContent>
				<uirtl.Button>إنشاء مشروع</uirtl.Button>
			</uirtl.EmptyContent>
		</uirtl.Empty>
	</div>
}
