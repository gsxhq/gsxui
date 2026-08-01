package table

import "github.com/gsxhq/gsxui/ui"

// Rtl mirrors shadcn's own table-rtl demo: the same invoices table as Data
// plus a payment-method column and a footer total row, translated to
// Arabic and wrapped in dir="rtl".
component Rtl() {
	{{
		invoices := []struct{ Invoice, Status, Method, Amount string }{
			{"INV001", "مدفوع", "بطاقة ائتمانية", "$250.00"},
			{"INV002", "قيد الانتظار", "PayPal", "$150.00"},
			{"INV003", "غير مدفوع", "تحويل بنكي", "$350.00"},
		}
	}}
	<ui.Table dir="rtl" lang="ar">
		<ui.TableCaption>قائمة بفواتيرك الأخيرة.</ui.TableCaption>
		<ui.TableHeader>
			<ui.TableRow>
				<ui.TableHead class="w-[100px]">الفاتورة</ui.TableHead>
				<ui.TableHead>الحالة</ui.TableHead>
				<ui.TableHead>الطريقة</ui.TableHead>
				<ui.TableHead class="text-right">المبلغ</ui.TableHead>
			</ui.TableRow>
		</ui.TableHeader>
		<ui.TableBody>
			{ for _, inv := range invoices {
				<ui.TableRow>
					<ui.TableCell class="font-medium">{ inv.Invoice }</ui.TableCell>
					<ui.TableCell>{ inv.Status }</ui.TableCell>
					<ui.TableCell>{ inv.Method }</ui.TableCell>
					<ui.TableCell class="text-right">{ inv.Amount }</ui.TableCell>
				</ui.TableRow>
			} }
		</ui.TableBody>
		<ui.TableFooter>
			<ui.TableRow>
				<ui.TableCell colspan="3">المجموع</ui.TableCell>
				<ui.TableCell class="text-right">$2,500.00</ui.TableCell>
			</ui.TableRow>
		</ui.TableFooter>
	</ui.Table>
}
