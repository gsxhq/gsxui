package empty

import (
	"github.com/gsxhq/gsxui/ui"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Rtl mirrors shadcn's own empty-rtl demo's Arabic "no projects yet" copy,
// kept in Basic's shape (icon.Inbox media, title, description, one
// EmptyContent Button — shadcn's version swaps in a folder icon and a
// second outline Button plus a "learn more" link, which this directory's
// Basic doesn't compose; DEVIATED from rather than invented).
component Rtl() {
	<div dir="rtl" lang="ar">
		<ui.Empty>
			<ui.EmptyHeader>
				<ui.EmptyMedia>الوسائط الافتراضية</ui.EmptyMedia>
				<ui.EmptyMedia variant="icon">
					<icon.Inbox/>
				</ui.EmptyMedia>
				<ui.EmptyTitle>لا توجد مشاريع بعد</ui.EmptyTitle>
				<ui.EmptyDescription>
					لم تقم بإنشاء أي مشاريع بعد. ابدأ بإنشاء مشروعك الأول.
				</ui.EmptyDescription>
			</ui.EmptyHeader>
			<ui.EmptyContent>
				<ui.Button>إنشاء مشروع</ui.Button>
			</ui.EmptyContent>
		</ui.Empty>
	</div>
}
