// Package table holds the site's example gsx components for ui/table.
package table

import uitable "github.com/gsxhq/gsxui/ui/table"

// Basic renders a minimal table: header row plus two data rows.
component Basic() {
	<uitable.Table>
		<uitable.TableHeader>
			<uitable.TableRow><uitable.TableHead>Product</uitable.TableHead><uitable.TableHead>Price</uitable.TableHead></uitable.TableRow>
		</uitable.TableHeader>
		<uitable.TableBody>
			<uitable.TableRow><uitable.TableCell>Widget</uitable.TableCell><uitable.TableCell>$9.00</uitable.TableCell></uitable.TableRow>
			<uitable.TableRow><uitable.TableCell>Gadget</uitable.TableCell><uitable.TableCell>$19.00</uitable.TableCell></uitable.TableRow>
		</uitable.TableBody>
	</uitable.Table>
}
