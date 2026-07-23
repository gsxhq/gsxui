package table

import uitable "github.com/gsxhq/gsxui/ui/table"

// Data renders a compound table: caption, header, and three rows built
// from a slice via a for-range loop; the Unpaid row's status gets a
// conditional class.
component Data() {
	{{ invoices := []struct{ Invoice, Status, Amount string }{{"INV001", "Paid", "$250.00"}, {"INV002", "Pending", "$150.00"}, {"INV003", "Unpaid", "$350.00"}} }}
	<uitable.Table>
		<uitable.TableCaption>A list of your recent invoices.</uitable.TableCaption>
		<uitable.TableHeader>
			<uitable.TableRow><uitable.TableHead>Invoice</uitable.TableHead><uitable.TableHead>Status</uitable.TableHead><uitable.TableHead class="text-right">Amount</uitable.TableHead></uitable.TableRow>
		</uitable.TableHeader>
		<uitable.TableBody>
			{ for _, inv := range invoices {
				<uitable.TableRow>
					<uitable.TableCell class="font-medium">{ inv.Invoice }</uitable.TableCell>
					<uitable.TableCell class={ "text-destructive": inv.Status == "Unpaid" }>{ inv.Status }</uitable.TableCell>
					<uitable.TableCell class="text-right">{ inv.Amount }</uitable.TableCell>
				</uitable.TableRow>
			} }
		</uitable.TableBody>
	</uitable.Table>
}
