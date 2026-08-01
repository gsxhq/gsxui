package inputgroup

import (
	"github.com/gsxhq/gsxui/ui"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Rtl keeps Basic's own shape (search InputGroup, email + trailing-button
// InputGroup, block-start/block-end addon group, variant-button row,
// invalid group), translated to Arabic and wrapped in dir="rtl".
component Rtl() {
	<div dir="rtl" lang="ar">
		<div class="flex w-full max-w-sm flex-col gap-4">
			<ui.InputGroup>
				<ui.InputGroupAddon>
					<icon.Search class="size-4"/>
				</ui.InputGroupAddon>
				<ui.InputGroupInput placeholder="بحث..."/>
			</ui.InputGroup>
			<ui.InputGroup>
				<ui.InputGroupInput placeholder="البريد الإلكتروني" type="email"/>
				<ui.InputGroupAddon align="inline-end">
					<ui.InputGroupButton aria-label="إرسال">
						<icon.Send/>
					</ui.InputGroupButton>
				</ui.InputGroupAddon>
			</ui.InputGroup>
			<ui.InputGroup data-disabled="true">
				<ui.InputGroupAddon align="block-start">
					<ui.InputGroupText>الرابط</ui.InputGroupText>
				</ui.InputGroupAddon>
				<ui.InputGroupInput disabled placeholder="معطل"/>
				<ui.InputGroupAddon align="block-end">
					<ui.InputGroupText>مطلوب</ui.InputGroupText>
				</ui.InputGroupAddon>
			</ui.InputGroup>
			<div class="flex flex-wrap gap-2">
				<ui.InputGroupButton variant="default">افتراضي</ui.InputGroupButton>
				<ui.InputGroupButton variant="destructive">مدمر</ui.InputGroupButton>
				<ui.InputGroupButton variant="outline">مخطط</ui.InputGroupButton>
				<ui.InputGroupButton variant="secondary">ثانوي</ui.InputGroupButton>
				<ui.InputGroupButton variant="link">رابط</ui.InputGroupButton>
				<ui.InputGroupButton size="sm">صغير</ui.InputGroupButton>
				<ui.InputGroupButton size="icon-sm" aria-label="أيقونة">
					<icon.Send/>
				</ui.InputGroupButton>
			</div>
			<ui.InputGroup>
				<ui.InputGroupInput aria-invalid="true" value="قيمة غير صالحة"/>
			</ui.InputGroup>
		</div>
	</div>
}
