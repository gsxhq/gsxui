package skeleton

import "github.com/gsxhq/gsxui/site/uirtl"

// Rtl mirrors shadcn's own skeleton-rtl demo: the same circular avatar plus
// two text-line placeholders as Basic, wrapped in dir="rtl" — a skeleton's
// shapes carry no text content to translate, so dir only flips the row's
// reading direction.
component Rtl() {
	<div dir="rtl" lang="ar" class="flex items-center gap-4">
		<uirtl.Skeleton class="size-12 rounded-full"/>
		<div class="grid gap-2">
			<uirtl.Skeleton class="h-4 w-[250px]"/>
			<uirtl.Skeleton class="h-4 w-[200px]"/>
		</div>
	</div>
}
