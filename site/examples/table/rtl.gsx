package table

import "github.com/gsxhq/gsxui/site/uirtl"

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
	<uirtl.Table dir="rtl" lang="ar">
		<uirtl.TableCaption>قائمة بفواتيرك الأخيرة.</uirtl.TableCaption>
		<uirtl.TableHeader>
			<uirtl.TableRow>
				<uirtl.TableHead class="w-[100px]">الفاتورة</uirtl.TableHead>
				<uirtl.TableHead>الحالة</uirtl.TableHead>
				<uirtl.TableHead>الطريقة</uirtl.TableHead>
				<uirtl.TableHead class="text-right">المبلغ</uirtl.TableHead>
			</uirtl.TableRow>
		</uirtl.TableHeader>
		<uirtl.TableBody>
			{ for _, inv := range invoices {
				<uirtl.TableRow>
					<uirtl.TableCell class="font-medium">{ inv.Invoice }</uirtl.TableCell>
					<uirtl.TableCell>{ inv.Status }</uirtl.TableCell>
					<uirtl.TableCell>{ inv.Method }</uirtl.TableCell>
					<uirtl.TableCell class="text-right">{ inv.Amount }</uirtl.TableCell>
				</uirtl.TableRow>
			} }
		</uirtl.TableBody>
		<uirtl.TableFooter>
			<uirtl.TableRow>
				<uirtl.TableCell colspan="3">المجموع</uirtl.TableCell>
				<uirtl.TableCell class="text-right">$2,500.00</uirtl.TableCell>
			</uirtl.TableRow>
		</uirtl.TableFooter>
	</uirtl.Table>
}
