// Package table holds the site's example gsx components for ui/table.
package table

import "github.com/gsxhq/gsxui/ui"

// Basic renders a minimal table: header row plus two data rows.
component Basic() {
	<ui.Table>
		<ui.TableHeader>
			<ui.TableRow>
				<ui.TableHead>Product</ui.TableHead>
				<ui.TableHead>Price</ui.TableHead>
			</ui.TableRow>
		</ui.TableHeader>
		<ui.TableBody>
			<ui.TableRow>
				<ui.TableCell>Widget</ui.TableCell>
				<ui.TableCell>$9.00</ui.TableCell>
			</ui.TableRow>
			<ui.TableRow>
				<ui.TableCell>Gadget</ui.TableCell>
				<ui.TableCell>$19.00</ui.TableCell>
			</ui.TableRow>
		</ui.TableBody>
	</ui.Table>
}
