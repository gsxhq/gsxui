package table

import "github.com/gsxhq/gsxui/ui"

// Data renders a compound table: caption, header, and three rows built
// from a slice via a for-range loop; the Unpaid row's status gets a
// conditional class.
component Data() {
	{{ invoices := []struct{ Invoice, Status, Amount string }{{"INV001", "Paid", "$250.00"}, {"INV002", "Pending", "$150.00"}, {"INV003", "Unpaid", "$350.00"}} }}
	<ui.Table>
		<ui.TableCaption>A list of your recent invoices.</ui.TableCaption>
		<ui.TableHeader>
			<ui.TableRow><ui.TableHead>Invoice</ui.TableHead><ui.TableHead>Status</ui.TableHead><ui.TableHead class="text-right">Amount</ui.TableHead></ui.TableRow>
		</ui.TableHeader>
		<ui.TableBody>
			{ for _, inv := range invoices {
				<ui.TableRow>
					<ui.TableCell class="font-medium">{ inv.Invoice }</ui.TableCell>
					<ui.TableCell class={ "text-destructive": inv.Status == "Unpaid" }>{ inv.Status }</ui.TableCell>
					<ui.TableCell class="text-right">{ inv.Amount }</ui.TableCell>
				</ui.TableRow>
			} }
		</ui.TableBody>
	</ui.Table>
}
