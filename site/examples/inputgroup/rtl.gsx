package inputgroup

import (
	"github.com/gsxhq/gsxui/site/uirtl"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Rtl keeps Basic's own shape (search InputGroup, email + trailing-button
// InputGroup, block-start/block-end addon group, variant-button row,
// invalid group), translated to Arabic and wrapped in dir="rtl".
component Rtl() {
	<div dir="rtl" lang="ar">
		<div class="flex w-full max-w-sm flex-col gap-4">
			<uirtl.InputGroup>
				<uirtl.InputGroupAddon>
					<icon.Search class="size-4"/>
				</uirtl.InputGroupAddon>
				<uirtl.InputGroupInput placeholder="بحث..."/>
			</uirtl.InputGroup>
			<uirtl.InputGroup>
				<uirtl.InputGroupInput placeholder="البريد الإلكتروني" type="email"/>
				<uirtl.InputGroupAddon align="inline-end">
					<uirtl.InputGroupButton aria-label="إرسال">
						<icon.Send/>
					</uirtl.InputGroupButton>
				</uirtl.InputGroupAddon>
			</uirtl.InputGroup>
			<uirtl.InputGroup data-disabled="true">
				<uirtl.InputGroupAddon align="block-start">
					<uirtl.InputGroupText>الرابط</uirtl.InputGroupText>
				</uirtl.InputGroupAddon>
				<uirtl.InputGroupInput disabled placeholder="معطل"/>
				<uirtl.InputGroupAddon align="block-end">
					<uirtl.InputGroupText>مطلوب</uirtl.InputGroupText>
				</uirtl.InputGroupAddon>
			</uirtl.InputGroup>
			<div class="flex flex-wrap gap-2">
				<uirtl.InputGroupButton variant="default">افتراضي</uirtl.InputGroupButton>
				<uirtl.InputGroupButton variant="destructive">مدمر</uirtl.InputGroupButton>
				<uirtl.InputGroupButton variant="outline">مخطط</uirtl.InputGroupButton>
				<uirtl.InputGroupButton variant="secondary">ثانوي</uirtl.InputGroupButton>
				<uirtl.InputGroupButton variant="link">رابط</uirtl.InputGroupButton>
				<uirtl.InputGroupButton size="sm">صغير</uirtl.InputGroupButton>
				<uirtl.InputGroupButton size="icon-sm" aria-label="أيقونة">
					<icon.Send/>
				</uirtl.InputGroupButton>
			</div>
			<uirtl.InputGroup>
				<uirtl.InputGroupInput aria-invalid="true" value="قيمة غير صالحة"/>
			</uirtl.InputGroup>
		</div>
	</div>
}
